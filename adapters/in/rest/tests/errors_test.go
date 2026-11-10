// adapters/in/rest/tests/errors_test.go — réponses d'erreur interne (référence, détail)
package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"myr-core/adapters/in/rest"
	"myr-core/domain/identity"
	"myr-core/domain/model"
)

// failingRequestMux sert une API dont la demande d'inscription échoue avec err,
// détail des erreurs internes activé ou non.
func failingRequestMux(t *testing.T, err error, debug bool) http.Handler {
	t.Helper()
	svc := &mockIdentitySvc{requestAccess: func(context.Context, identity.AccountRequest, identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
		return nil, "", err
	}}
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store).
		WithIdentityService(svc).
		WithNetworkService(openNetwork(false, "")).
		WithDebugErrors(debug)
	mux, herr := rest.NewServer(h, "localhost:0").Handler()
	if herr != nil {
		t.Fatalf("Handler(): %v", herr)
	}
	return mux
}

// postRequest soumet une demande d'inscription valide et décode la réponse d'erreur.
func postRequest(t *testing.T, mux http.Handler) (*httptest.ResponseRecorder, map[string]string) {
	t.Helper()
	body := `{"pseudo":"bob","email":"bob@example.com","org_id":"Org1MSP"}`
	r := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var resp map[string]string
	if err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&resp); err != nil {
		t.Fatalf("réponse JSON invalide : %v (body: %s)", err, w.Body.String())
	}
	return w, resp
}

var errorRefPattern = regexp.MustCompile(`^[0-9a-f]{8}$`)

// / @brief  Une erreur interne renvoie un message générique et une référence, jamais le détail
// / @input  RequestAccess échoue avec « identité bob@Org1 déjà enregistrée », détail désactivé
// / @expect 500 ; ref de 8 caractères hexadécimaux dans le corps et l'en-tête X-Myr-Error-Ref ; le détail est absent
func TestInternalErr_HidesDetailBehindRef(t *testing.T) {
	w, resp := postRequest(t, failingRequestMux(t, errors.New("identité bob@Org1 déjà enregistrée"), false))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("code 500 attendu, obtenu %d", w.Code)
	}
	ref := resp["ref"]
	if !errorRefPattern.MatchString(ref) || w.Header().Get("X-Myr-Error-Ref") != ref {
		t.Errorf("référence attendue dans le corps et l'en-tête, obtenu %q / %q", ref, w.Header().Get("X-Myr-Error-Ref"))
	}
	if strings.Contains(resp["error"], "bob@Org1") || !strings.Contains(resp["error"], ref) {
		t.Errorf("message générique avec référence attendu, obtenu %q", resp["error"])
	}
}

// / @brief  Avec le détail des erreurs activé (serveur de test), le message porte l'erreur et sa référence
// / @input  même échec, WithDebugErrors(true)
// / @expect 500 ; le message contient le détail et la référence
func TestInternalErr_DebugErrorsShowsDetail(t *testing.T) {
	w, resp := postRequest(t, failingRequestMux(t, errors.New("identité bob@Org1 déjà enregistrée"), true))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("code 500 attendu, obtenu %d", w.Code)
	}
	if !strings.Contains(resp["error"], "bob@Org1") || !strings.Contains(resp["error"], resp["ref"]) {
		t.Errorf("détail et référence attendus, obtenu %q", resp["error"])
	}
}

// / @brief  L'indisponibilité de la blockchain reste explicite (503), avec sa référence
// / @input  RequestAccess échoue avec ErrBlockchainUnavailable, détail désactivé
// / @expect 503 ; le message cite l'indisponibilité ; référence présente
func TestInternalErr_BlockchainUnavailable_Explicit(t *testing.T) {
	w, resp := postRequest(t, failingRequestMux(t, model.ErrBlockchainUnavailable, false))
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("code 503 attendu, obtenu %d", w.Code)
	}
	if !strings.Contains(resp["error"], model.ErrBlockchainUnavailable.Error()) || resp["ref"] == "" {
		t.Errorf("message explicite et référence attendus, obtenu %v", resp)
	}
}

// / @brief  Un asset soumis renvoie 409 et une licence refusée 400, au lieu d'une erreur interne (RM19, RM03)
// / @input  PATCH et POST /submit sur un asset soumis (ErrAssetSubmitted) ; PATCH avec licence inconnue puis incompatible
// / @expect 409 pour l'asset soumis, 400 pour la licence, message explicite
func TestDomainErrors_SubmittedAndLicense_MappedToHTTP(t *testing.T) {
	var next error
	svc := &mockSvc{
		updateAsset: func(model.UpdateRequest) (*model.Model3D, error) { return nil, next },
		submit:      func(string, string) (*model.Model3D, error) { return nil, next },
	}
	mux := newTestMux(t, svc)
	cases := []struct {
		err          error
		method, path string
		want         int
	}{
		{model.ErrAssetSubmitted, http.MethodPatch, "/api/components/c1", http.StatusConflict},
		{model.ErrAssetSubmitted, http.MethodPost, "/api/components/c1/submit", http.StatusConflict},
		{model.ErrUnknownLicense, http.MethodPatch, "/api/components/c1", http.StatusBadRequest},
		{model.ErrLicenseIncompatible, http.MethodPatch, "/api/components/c1", http.StatusBadRequest},
	}
	for _, c := range cases {
		next = c.err
		w := do(t, mux, c.method, c.path, `{"name":"x"}`)
		if w.Code != c.want || !strings.Contains(w.Body.String(), strings.TrimPrefix(c.err.Error(), "model: ")) {
			t.Errorf("%s %s (%v) : code %d attendu, obtenu %d — %s", c.method, c.path, c.err, c.want, w.Code, w.Body.String())
		}
	}
}

// / @brief  Une liaison ou une composition refusée par les règles d'assemblage renvoie 409 ou 400, au lieu d'une erreur interne (RM09, RM10, RM50, RM51, RM52)
// / @input  POST /connections, /assembly-links, /virtual-connect, /components/{id}/assemblies et /components/{id}/instances dont le service refuse la demande
// / @expect 409 pour une interface déjà engagée ou une liaison existante, 400 pour une liaison incompatible, réflexive, étrangère au module ou une composition cyclique
func TestDomainErrors_AssemblyRules_MappedToHTTP(t *testing.T) {
	var next error
	svc := &mockSvc{
		addConnection: func(string, string, string) (*model.Connection, error) { return nil, next },
		addAssemblyLink: func(string, string, string, string, string, string) (*model.Connection, error) {
			return nil, next
		},
		connectVirtualToPhysical: func(string, string, model.AssetInterface, string, string) (*model.Connection, error) {
			return nil, next
		},
		addAssemblyToModule: func(string, string) error { return next },
		addAssetToWorkspace: func(string, string) (*model.Model3D, error) { return nil, next },
	}
	mux := newTestMux(t, svc)
	cases := []struct {
		err        error
		path, body string
		want       int
	}{
		{model.ErrConnectionExists, "/api/connections", `{"from":"a","to":"b"}`, http.StatusConflict},
		{model.ErrInterfaceInUse, "/api/assembly-links", `{"from_iface_id":"i1","to_iface_id":"i2"}`, http.StatusConflict},
		{model.ErrIncompatibleInterfaces, "/api/assembly-links", `{"from_iface_id":"i1","to_iface_id":"i2"}`, http.StatusBadRequest},
		{model.ErrSelfLink, "/api/assembly-links", `{"from_iface_id":"i1","to_iface_id":"i1"}`, http.StatusBadRequest},
		{model.ErrInterfaceInUse, "/api/virtual-connect", `{"virtual_iface_id":"v","physical_iface_id":"p"}`, http.StatusConflict},
		{model.ErrConnectionNotInModule, "/api/components/m1/assemblies", `{"connection_id":"c1"}`, http.StatusBadRequest},
		{model.ErrCyclicComposition, "/api/components/m1/instances", `{"asset_id":"m1"}`, http.StatusBadRequest},
		{model.ErrNotFound, "/api/components/m1/instances", `{"asset_id":"fantome"}`, http.StatusNotFound},
	}
	for _, c := range cases {
		next = c.err
		w := do(t, mux, http.MethodPost, c.path, c.body)
		if w.Code != c.want || !strings.Contains(w.Body.String(), strings.TrimPrefix(c.err.Error(), "model: ")) {
			t.Errorf("POST %s (%v) : code %d attendu, obtenu %d — %s", c.path, c.err, c.want, w.Code, w.Body.String())
		}
	}
}

// / @brief  PATCH d'une interface qui changerait d'asset porteur renvoie 400 au lieu d'une erreur interne
// / @input  PATCH /api/interfaces/i1 avec asset_id ; le service refuse (ErrInterfaceCarrierChange)
// / @expect 400 avec le message du refus
func TestDomainErrors_InterfaceCarrierChange_BadRequest(t *testing.T) {
	svc := &mockSvc{
		getInterface: func(id string) (*model.AssetInterface, error) {
			return &model.AssetInterface{ID: id, AssetID: "a1", Category: "ELEC"}, nil
		},
		updateInterface: func(*model.AssetInterface) error { return model.ErrInterfaceCarrierChange },
	}
	w := do(t, newTestMux(t, svc), http.MethodPatch, "/api/interfaces/i1", `{"asset_id":"a2"}`)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "porteur") {
		t.Errorf("400 attendu, obtenu %d — %s", w.Code, w.Body.String())
	}
}

// / @brief  Le service agit pour l'identité de la session, et le refus d'un asset d'une autre identité renvoie 403 (RM54)
// / @input  PATCH /api/components/c1 par la session de « test » ; le service refuse (ErrNotOwner)
// / @expect le service a reçu l'identité de la session ; HTTP 403 avec le message du refus
func TestDomainErrors_RM54_NotOwner_Forbidden(t *testing.T) {
	svc := &mockSvc{updateAsset: func(model.UpdateRequest) (*model.Model3D, error) { return nil, model.ErrNotOwner }}
	w := do(t, newTestMux(t, svc), http.MethodPatch, "/api/components/c1", `{"name":"x"}`)
	if !strings.HasPrefix(svc.actor, "test") {
		t.Errorf("identité de la session attendue, le service a reçu %q", svc.actor)
	}
	if w.Code != http.StatusForbidden || !strings.Contains(w.Body.String(), "propriétaire") {
		t.Errorf("403 attendu, obtenu %d — %s", w.Code, w.Body.String())
	}
}
