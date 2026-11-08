// adapters/in/rest/tests/signer_test.go
// Signature des écritures par le wallet de l'utilisateur de la session
// (ADR-12) : routage vers BlockchainForSigner ou, sans identité, vers un accès
// en lecture seule ; refus d'écriture sans signataire traduit en 403.
package rest_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"myr-core/adapters/in/rest"
	"myr-core/domain/identity"
	"myr-core/domain/model"
)

// emptyLedger est un BlockchainPort vide.
type emptyLedger struct{}

func (emptyLedger) StoreModelRecord(*model.Model3D) error { return nil }
func (emptyLedger) UpdateModelMetadata(_, _ string, _ model.AssetMetadata) error { return nil }
func (emptyLedger) GetModelRecord(_, _ string) (*model.Model3D, error) {
	return nil, model.ErrNotFound
}
func (emptyLedger) ListModelRecords(string) ([]*model.Model3D, error) { return nil, nil }
func (emptyLedger) VerifyIntegrity(_, _, _ string) (bool, error)      { return false, nil }

// recordingRouter enregistre l'accès demandé pour chaque requête.
type recordingRouter struct {
	readOnly []string
	signed   []string // "networkID|mspID|certPath|keyPath"
}

func (r *recordingRouter) BlockchainFor(networkID string) model.BlockchainPort {
	r.readOnly = append(r.readOnly, networkID)
	return emptyLedger{}
}
func (r *recordingRouter) BlockchainForSigner(networkID, mspID, certPath, keyPath string) model.BlockchainPort {
	r.signed = append(r.signed, strings.Join([]string{networkID, mspID, certPath, keyPath}, "|"))
	return emptyLedger{}
}

var testWallet = identity.WalletEntry{Handle: "test@TestMSP", OrgID: "TestMSP", MSPDir: filepath.Join("wallets", "test@TestMSP", "msp")}

// newSignerMux construit un serveur multi-réseau (réseau net1) dont le service
// d'identité retrouve — ou non — le wallet de l'identité de la session de test.
func newSignerMux(t *testing.T, router *recordingRouter, walletFound bool) http.Handler {
	t.Helper()
	store := &mockStore{}
	ids := &mockIdentitySvc{wallet: func(handle string) (identity.WalletEntry, error) {
		if walletFound && handle == testWallet.Handle {
			return testWallet, nil
		}
		return identity.WalletEntry{}, identity.ErrWalletNotFound
	}}
	h := rest.NewHandler(&mockSvc{}, store, store).
		WithIdentityService(ids).
		WithBlockchainRouter(router).
		WithDefaultNetwork("net1")
	mux, err := rest.NewServer(h, "localhost:0").Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

func getComponents(t *testing.T, mux http.Handler, token string) {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/api/components", nil)
	r.Header.Set("X-Myr-Token", token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/components : %d (%s)", w.Code, w.Body.String())
	}
}

// / @brief  La session d'un utilisateur enrôlé accède au ledger avec son wallet (ADR-12)
// / @input  session créée par POST /api/identity/session (identité test@TestMSP, wallet présent), puis GET /api/components
// / @expect chaque accès au ledger de la requête passe par BlockchainForSigner pour net1, avec le MSP, le certificat et la clé du wallet ; aucun accès en lecture seule
func TestSvcFor_SessionWithWallet_SignsWithWallet(t *testing.T) {
	router := &recordingRouter{}
	mux := newSignerMux(t, router, true)

	getComponents(t, mux, loginToken(t, mux))

	want := strings.Join([]string{"net1", "TestMSP", testWallet.CertPath(), testWallet.KeyPath()}, "|")
	if len(router.signed) == 0 {
		t.Fatal("aucun accès signé")
	}
	for _, got := range router.signed {
		if got != want {
			t.Errorf("accès signé = %s, want %s", got, want)
		}
	}
	if len(router.readOnly) != 0 {
		t.Errorf("accès en lecture seule inattendus : %v", router.readOnly)
	}
}

// / @brief  Sans wallet pour l'identité de la session, l'accès au ledger est en lecture seule
// / @input  session test@TestMSP dont le wallet est introuvable sur ce nœud, puis GET /api/components
// / @expect les accès au ledger passent par BlockchainFor (net1) ; aucun accès signé
func TestSvcFor_SessionWithoutWallet_ReadOnly(t *testing.T) {
	router := &recordingRouter{}
	mux := newSignerMux(t, router, false)

	getComponents(t, mux, loginToken(t, mux))

	if len(router.signed) != 0 || len(router.readOnly) == 0 {
		t.Errorf("signés = %v, lecture seule = %v ; want aucun signé et des accès en lecture seule", router.signed, router.readOnly)
	}
}

// / @brief  Une requête sans session (route publique) n'accède au ledger qu'en lecture seule
// / @input  GET /api/status sans X-Myr-Token
// / @expect BlockchainFor appelé ; aucun accès signé
func TestSvcFor_NoSession_ReadOnly(t *testing.T) {
	router := &recordingRouter{}
	mux := newSignerMux(t, router, true)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/status", nil))

	if len(router.signed) != 0 || len(router.readOnly) == 0 {
		t.Errorf("signés = %v, lecture seule = %v ; want aucun signé", router.signed, router.readOnly)
	}
}

// / @brief  Une soumission refusée faute d'identité signataire est un refus d'autorisation
// / @input  POST /api/components/c1/submit, Submit retourne model.ErrSignerRequired
// / @expect HTTP 403, le message invite à se reconnecter avec son identité
func TestSubmit_SignerRequired_Forbidden(t *testing.T) {
	svc := &mockSvc{submit: func(string, string) (*model.Model3D, error) { return nil, model.ErrSignerRequired }}
	mux := newTestMux(t, svc)

	r := httptest.NewRequest(http.MethodPost, "/api/components/c1/submit", strings.NewReader(`{}`))
	r.Header.Set("X-Myr-Token", loginToken(t, mux))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("got %d, want 403 (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "identity/session") {
		t.Errorf("message sans invitation à se reconnecter : %s", w.Body.String())
	}
}
