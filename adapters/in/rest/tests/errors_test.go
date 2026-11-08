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
