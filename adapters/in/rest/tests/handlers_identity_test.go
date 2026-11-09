// adapters/in/rest/tests/handlers_identity_test.go — tests externes des handlers d'identité
package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"myr-core/adapters/in/rest"
	"myr-core/domain/identity"
	"myr-core/domain/network"
	rbac "myr-core/domain/role"
)

// ── Mock IdentityService ──────────────────────────────────────────────────────

type mockIdentitySvc struct {
	listWallets     func() ([]identity.WalletEntry, error)
	wallet          func(handle string) (identity.WalletEntry, error)
	createAccount   func(ctx context.Context, pseudo, orgID, role, email string) (string, string, error)
	enroll          func(ctx context.Context, name, secret, orgID string) (identity.WalletEntry, error)
	getStatus       func(ctx context.Context, wallet identity.WalletEntry) (string, error)
	reEnroll        func(ctx context.Context, wallet identity.WalletEntry) (identity.WalletEntry, error)
	loadGuestWallet func(certPath, keyPath, caCertPath, orgID string) (identity.WalletEntry, error)
	requestAccess   func(ctx context.Context, req identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error)
	requestRole     func(ctx context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error)
	approveRequest  func(ctx context.Context, id string) (*identity.AccountRequest, error)
	rejectRequest   func(ctx context.Context, id string) (*identity.AccountRequest, error)
	listRequests    func() ([]*identity.AccountRequest, error)
	walletDir       func() string
	setRole         func(ctx context.Context, name, newRole string) error
	setStatus       func(ctx context.Context, name, status string) error
}

func (m *mockIdentitySvc) ListLocalWallets() ([]identity.WalletEntry, error) {
	if m.listWallets != nil {
		return m.listWallets()
	}
	return nil, nil
}
func (m *mockIdentitySvc) Wallet(handle string) (identity.WalletEntry, error) {
	if m.wallet != nil {
		return m.wallet(handle)
	}
	return identity.WalletEntry{}, identity.ErrWalletNotFound
}
func (m *mockIdentitySvc) CreateAccount(ctx context.Context, pseudo, orgID, role, email string) (string, string, error) {
	if m.createAccount != nil {
		return m.createAccount(ctx, pseudo, orgID, role, email)
	}
	return identity.Handle(pseudo, orgID), "secret", nil
}

// Enroll renvoie par défaut une identité active dotée d'un rôle, comme le
// certificat d'une identité approuvée — celle qui peut ouvrir une session.
func (m *mockIdentitySvc) Enroll(ctx context.Context, name, secret, orgID string) (identity.WalletEntry, error) {
	if m.enroll != nil {
		return m.enroll(ctx, name, secret, orgID)
	}
	return identity.WalletEntry{
		Handle: name + "@" + orgID, Name: name, OrgID: orgID,
		Status: identity.StatusActive, Role: identity.RoleContributor,
	}, nil
}
func (m *mockIdentitySvc) GetStatus(ctx context.Context, wallet identity.WalletEntry) (string, error) {
	if m.getStatus != nil {
		return m.getStatus(ctx, wallet)
	}
	return identity.StatusActive, nil
}
func (m *mockIdentitySvc) ReEnroll(ctx context.Context, wallet identity.WalletEntry) (identity.WalletEntry, error) {
	if m.reEnroll != nil {
		return m.reEnroll(ctx, wallet)
	}
	return wallet, nil
}
func (m *mockIdentitySvc) RequestAccess(ctx context.Context, req identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
	if m.requestAccess != nil {
		return m.requestAccess(ctx, req, policy)
	}
	if !policy.AllowAutoRegister {
		return nil, "", identity.ErrRegistrationClosed
	}
	req.ID = "req-test"
	req.Status = identity.RequestApproved
	if policy.RequireApproval {
		req.Status = identity.RequestPending
	}
	return &req, "auto-secret", nil
}
func (m *mockIdentitySvc) ApproveRequest(ctx context.Context, id string) (*identity.AccountRequest, error) {
	if m.approveRequest != nil {
		return m.approveRequest(ctx, id)
	}
	return &identity.AccountRequest{ID: id, Status: identity.RequestApproved}, nil
}
func (m *mockIdentitySvc) RejectRequest(ctx context.Context, id string) (*identity.AccountRequest, error) {
	if m.rejectRequest != nil {
		return m.rejectRequest(ctx, id)
	}
	return &identity.AccountRequest{ID: id, Status: identity.RequestRejected}, nil
}
func (m *mockIdentitySvc) ListRequests() ([]*identity.AccountRequest, error) {
	if m.listRequests != nil {
		return m.listRequests()
	}
	return nil, nil
}
func (m *mockIdentitySvc) LoadGuestWallet(certPath, keyPath, caCertPath, orgID string) (identity.WalletEntry, error) {
	if m.loadGuestWallet != nil {
		return m.loadGuestWallet(certPath, keyPath, caCertPath, orgID)
	}
	return identity.WalletEntry{Handle: "guest@" + orgID, Name: "guest", OrgID: orgID, Status: identity.StatusActive}, nil
}
func (m *mockIdentitySvc) WalletDir() string {
	if m.walletDir != nil {
		return m.walletDir()
	}
	return "/tmp/wallets"
}
func (m *mockIdentitySvc) RequestRole(ctx context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error) {
	if m.requestRole != nil {
		return m.requestRole(ctx, w, wanted)
	}
	return &identity.AccountRequest{ID: "req-1", Kind: identity.RequestKindRole, Role: wanted, Status: identity.RequestPending}, nil
}
func (m *mockIdentitySvc) SetStatus(ctx context.Context, name, status string) error {
	if m.setStatus != nil {
		return m.setStatus(ctx, name, status)
	}
	return nil
}
func (m *mockIdentitySvc) SetRole(ctx context.Context, name, newRole string) error {
	if m.setRole != nil {
		return m.setRole(ctx, name, newRole)
	}
	return nil
}

var _ identity.IdentityService = (*mockIdentitySvc)(nil)

// ── Mock NetworkService ───────────────────────────────────────────────────────

type mockNetworkSvc struct {
	getActive func() (*network.NetworkProfile, error)
}

func (m *mockNetworkSvc) GetActive() (*network.NetworkProfile, error) {
	if m.getActive != nil {
		return m.getActive()
	}
	return nil, nil
}
func (m *mockNetworkSvc) Add(name, peerEndpoint, gatewayPeer, mspID, certPath, keyPath, tlsCertPath, fabricChannel, chaincodeName, caEndpoint, caName string, channels []string, allowAutoGuest, allowAutoRegister bool, autoRegisterRole string, isProduction bool) (*network.NetworkProfile, error) {
	return nil, nil
}
func (m *mockNetworkSvc) Update(id, name, peerEndpoint, gatewayPeer, mspID, certPath, keyPath, tlsCertPath, fabricChannel, chaincodeName, caEndpoint, caName string, allowAutoGuest, allowAutoRegister bool, autoRegisterRole string, isProduction bool) (*network.NetworkProfile, error) {
	return nil, nil
}
func (m *mockNetworkSvc) SetRequireApproval(id string, require bool) (*network.NetworkProfile, error) {
	return nil, nil
}
func (m *mockNetworkSvc) List() ([]*network.NetworkProfile, error) { return nil, nil }
func (m *mockNetworkSvc) Activate(id string) error                 { return nil }
func (m *mockNetworkSvc) Delete(id string) error                   { return nil }
func (m *mockNetworkSvc) TestConnection(id string) error           { return nil }
func (m *mockNetworkSvc) AddPeer(networkID string, req network.AddPeerRequest) (*network.PeerCredentials, error) {
	return nil, nil
}
func (m *mockNetworkSvc) Create(req network.CreateNetworkRequest) (*network.NetworkProfile, error) {
	return nil, nil
}

var _ network.NetworkService = (*mockNetworkSvc)(nil)

// ── helpers ───────────────────────────────────────────────────────────────────

// newIdentityMux builds a REST mux with an identity service — requireAuth
// always requires a token, obtainable via a real POST /api/identity/session.
func newIdentityMux(t *testing.T, idSvc identity.IdentityService) http.Handler {
	t.Helper()
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	h.WithIdentityService(idSvc)
	srv := rest.NewServer(h, "localhost:0")
	mux, err := srv.Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

// mockRoleSvc is a minimal rbac.RoleService stub for tests that need to control
// HasPermission independently of the legacy role-rank fallback. Only HasPermission
// is exercised by requireRole; the other methods are unused CRUD surface.
type mockRoleSvc struct {
	hasPermission func(roleName string, p rbac.Permission) bool
}

func (m *mockRoleSvc) Create(name string, permissions []rbac.Permission, autoGrant bool) (*rbac.Role, error) {
	return nil, nil
}
func (m *mockRoleSvc) Update(name string, change rbac.Change) (*rbac.Role, error) {
	return nil, nil
}
func (m *mockRoleSvc) Delete(name string) error            { return nil }
func (m *mockRoleSvc) Get(name string) (*rbac.Role, error) { return nil, nil }
func (m *mockRoleSvc) List() ([]*rbac.Role, error)         { return nil, nil }
func (m *mockRoleSvc) HasPermission(roleName string, p rbac.Permission) bool {
	if m.hasPermission != nil {
		return m.hasPermission(roleName, p)
	}
	return false
}

// newIdentityMuxWithRole is like newIdentityMux but injects a rbac.RoleService,
// letting a test grant a permission to whatever role loginToken's session ends up
// with (the role carried by the mock's enrolled certificate) without depending
// on that role directly.
func newIdentityMuxWithRole(t *testing.T, idSvc identity.IdentityService, roleSvc rbac.RoleService) http.Handler {
	t.Helper()
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	h.WithIdentityService(idSvc)
	h.WithRoleService(roleSvc)
	srv := rest.NewServer(h, "localhost:0")
	mux, err := srv.Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

// newBothMux builds a REST mux with identity + network services (blockchain mode).
func newBothMux(t *testing.T, idSvc identity.IdentityService, netSvc network.NetworkService) http.Handler {
	t.Helper()
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	h.WithIdentityService(idSvc)
	h.WithNetworkService(netSvc)
	srv := rest.NewServer(h, "localhost:0")
	mux, err := srv.Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

// ── GET /api/identity/wallets ─────────────────────────────────────────────────

// / @brief  GET /api/identity/wallets ne renvoie que le wallet de l'appelant, jamais
// / celui des autres utilisateurs enrôlés sur le même serveur (walletDir est partagé)
// / @input  session "test" (pseudo fixé par loginToken), 3 wallets dont un seul "test"
// / @expect seul le wallet dont Name == "test" est renvoyé
func TestHandleIdentityWallets_ScopedToCaller(t *testing.T) {
	svc := &mockIdentitySvc{
		listWallets: func() ([]identity.WalletEntry, error) {
			return []identity.WalletEntry{
				{Handle: "alice@Org1", Name: "alice", OrgID: "Org1MSP", Status: identity.StatusActive},
				{Handle: "test@Org1", Name: "test", OrgID: "Org1MSP", Status: identity.StatusActive},
				{Handle: "bob@Org2", Name: "bob", OrgID: "Org2MSP", Status: identity.StatusPending},
			}, nil
		},
	}
	w := do(t, newIdentityMux(t, svc), http.MethodGet, "/api/identity/wallets", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code attendu 200, obtenu %d", w.Code)
	}
	var dtos []struct {
		Handle string `json:"handle"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(w.Body).Decode(&dtos); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("1 wallet attendu (le sien), obtenu %d", len(dtos))
	}
	if dtos[0].Handle != "test@Org1" {
		t.Errorf("handle attendu %q, obtenu %q", "test@Org1", dtos[0].Handle)
	}
}

// / @brief  GET /api/identity/wallets renvoie tous les wallets à une session identity.admin
// / @input  RoleService accordant rbac.PermIdentityAdmin, 3 wallets appartenant à des pseudos différents
// / @expect les 3 wallets sont renvoyés, pas seulement celui de l'appelant
func TestHandleIdentityWallets_Admin_ReturnsAll(t *testing.T) {
	svc := &mockIdentitySvc{
		listWallets: func() ([]identity.WalletEntry, error) {
			return []identity.WalletEntry{
				{Handle: "alice@Org1", Name: "alice", OrgID: "Org1MSP", Status: identity.StatusActive},
				{Handle: "bob@Org2", Name: "bob", OrgID: "Org2MSP", Status: identity.StatusPending},
			}, nil
		},
	}
	roleSvc := &mockRoleSvc{hasPermission: func(string, rbac.Permission) bool { return true }}
	w := do(t, newIdentityMuxWithRole(t, svc, roleSvc), http.MethodGet, "/api/identity/wallets", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code attendu 200, obtenu %d", w.Code)
	}
	var dtos []struct {
		Handle string `json:"handle"`
	}
	if err := json.NewDecoder(w.Body).Decode(&dtos); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if len(dtos) != 2 {
		t.Fatalf("2 wallets attendus (admin voit tout), obtenu %d", len(dtos))
	}
}

func TestHandleIdentityWallets_EmptyList(t *testing.T) {
	w := do(t, newIdentityMux(t, &mockIdentitySvc{}), http.MethodGet, "/api/identity/wallets", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code attendu 200, obtenu %d", w.Code)
	}
	var dtos []struct {
		Handle string `json:"handle"`
	}
	json.NewDecoder(w.Body).Decode(&dtos)
	if len(dtos) != 0 {
		t.Errorf("liste vide attendue, obtenu %d éléments", len(dtos))
	}
}

func TestHandleIdentityWallets_NoService(t *testing.T) {
	// Sans identitySvc, aucune session ne peut être créée : requireAuth bloque
	// la requête (401) avant même d'atteindre le contrôle "identitySvc nil" du
	// handler — ce contrôle ne reste atteignable qu'avec une session déjà
	// valide (créée avant que l'identitySvc devienne indisponible).
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	srv := rest.NewServer(h, "localhost:0")
	mux, _ := srv.Handler()

	r := httptest.NewRequest(http.MethodGet, "/api/identity/wallets", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code attendu 401, obtenu %d", w.Code)
	}
}

func TestHandleIdentityWallets_ServiceError(t *testing.T) {
	svc := &mockIdentitySvc{
		listWallets: func() ([]identity.WalletEntry, error) {
			return nil, errors.New("disque illisible")
		},
	}
	w := do(t, newIdentityMux(t, svc), http.MethodGet, "/api/identity/wallets", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("code attendu 500, obtenu %d", w.Code)
	}
}

func TestHandleIdentityWallets_MethodNotAllowed(t *testing.T) {
	w := do(t, newIdentityMux(t, &mockIdentitySvc{}), http.MethodPost, "/api/identity/wallets", "")
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("code attendu 405, obtenu %d", w.Code)
	}
}

// ── POST /api/identity/enroll ─────────────────────────────────────────────────

func TestHandleIdentityEnroll_Success(t *testing.T) {
	var capturedName, capturedSecret, capturedOrg string
	svc := &mockIdentitySvc{
		enroll: func(_ context.Context, name, secret, orgID string) (identity.WalletEntry, error) {
			capturedName, capturedSecret, capturedOrg = name, secret, orgID
			return identity.WalletEntry{
				Handle: name + "@Org1",
				Name:   name,
				OrgID:  orgID,
				Status: identity.StatusPending,
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]string{"name": "carol", "secret": "s3cr3t", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newIdentityMux(t, svc).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("code attendu 200, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	var dto struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(w.Body).Decode(&dto); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if dto.Status != identity.StatusPending {
		t.Errorf("statut attendu %q, obtenu %q", identity.StatusPending, dto.Status)
	}
	if capturedName != "carol" || capturedSecret != "s3cr3t" || capturedOrg != "Org1MSP" {
		t.Errorf("arguments transmis incorrects : name=%q secret=%q org=%q", capturedName, capturedSecret, capturedOrg)
	}
}

func TestHandleIdentityEnroll_MissingFields(t *testing.T) {
	mux := newIdentityMux(t, &mockIdentitySvc{})
	cases := []map[string]string{
		{"name": "", "secret": "x", "org_id": "Org1MSP"},
		{"name": "alice", "secret": "", "org_id": "Org1MSP"},
		{"name": "alice", "secret": "x", "org_id": ""},
	}
	for _, bodyMap := range cases {
		b, _ := json.Marshal(bodyMap)
		req := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("code 400 attendu pour %v, obtenu %d", bodyMap, w.Code)
		}
	}
}

func TestHandleIdentityEnroll_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	newIdentityMux(t, &mockIdentitySvc{}).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityEnroll_ServiceError(t *testing.T) {
	svc := &mockIdentitySvc{
		enroll: func(_ context.Context, _, _, _ string) (identity.WalletEntry, error) {
			return identity.WalletEntry{}, errors.New("CA: secret invalide")
		},
	}

	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "wrong", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newIdentityMux(t, svc).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityEnroll_NoService(t *testing.T) {
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	srv := rest.NewServer(h, "localhost:0")
	mux, _ := srv.Handler()

	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "x", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("code 503 attendu, obtenu %d", w.Code)
	}
}

// ── GET /api/identity/status ──────────────────────────────────────────────────

func TestHandleIdentityStatus_Active(t *testing.T) {
	var capturedHandle string
	svc := &mockIdentitySvc{
		getStatus: func(_ context.Context, we identity.WalletEntry) (string, error) {
			capturedHandle = we.Handle
			return identity.StatusActive, nil
		},
	}
	w := do(t, newIdentityMux(t, svc), http.MethodGet, "/api/identity/status?handle=alice@Org1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code attendu 200, obtenu %d", w.Code)
	}
	var res map[string]string
	json.NewDecoder(w.Body).Decode(&res)
	if res["status"] != identity.StatusActive {
		t.Errorf("statut attendu %q, obtenu %q", identity.StatusActive, res["status"])
	}
	if capturedHandle != "alice@Org1" {
		t.Errorf("handle attendu %q transmis au service, obtenu %q", "alice@Org1", capturedHandle)
	}
}

func TestHandleIdentityStatus_MissingHandle(t *testing.T) {
	w := do(t, newIdentityMux(t, &mockIdentitySvc{}), http.MethodGet, "/api/identity/status", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityStatus_InvalidHandle(t *testing.T) {
	w := do(t, newIdentityMux(t, &mockIdentitySvc{}), http.MethodGet, "/api/identity/status?handle=noslash", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityStatus_ServiceError(t *testing.T) {
	svc := &mockIdentitySvc{
		getStatus: func(_ context.Context, _ identity.WalletEntry) (string, error) {
			return "", errors.New("CA injoignable")
		},
	}
	w := do(t, newIdentityMux(t, svc), http.MethodGet, "/api/identity/status?handle=alice@Org1", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("code 500 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityStatus_NoService(t *testing.T) {
	// Sans identitySvc, aucune session ne peut être créée : requireAuth bloque
	// la requête (401) avant même d'atteindre le contrôle "identitySvc nil" du
	// handler — voir TestHandleIdentityWallets_NoService pour le même raisonnement.
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	srv := rest.NewServer(h, "localhost:0")
	mux, _ := srv.Handler()

	r := httptest.NewRequest(http.MethodGet, "/api/identity/status?handle=alice@Org1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code 401 attendu, obtenu %d", w.Code)
	}
}

// ── GET /api/identity/policy ──────────────────────────────────────────────────

func TestHandleIdentityPolicy_BlockchainPrivate(t *testing.T) {
	netSvc := &mockNetworkSvc{
		getActive: func() (*network.NetworkProfile, error) {
			return &network.NetworkProfile{Name: "TestNet", MSPID: "Org1MSP", AllowAutoGuest: false}, nil
		},
	}
	w := do(t, newBothMux(t, &mockIdentitySvc{}, netSvc), http.MethodGet, "/api/identity/policy", "")
	var res struct {
		AllowAutoGuest bool   `json:"allow_auto_guest"`
		NetworkName    string `json:"network_name"`
	}
	json.NewDecoder(w.Body).Decode(&res)
	if res.AllowAutoGuest {
		t.Error("allow_auto_guest=false attendu pour un réseau privé")
	}
	if res.NetworkName != "TestNet" {
		t.Errorf("network_name attendu %q, obtenu %q", "TestNet", res.NetworkName)
	}
}

func TestHandleIdentityPolicy_BlockchainPublic(t *testing.T) {
	netSvc := &mockNetworkSvc{
		getActive: func() (*network.NetworkProfile, error) {
			return &network.NetworkProfile{Name: "PublicNet", AllowAutoGuest: true}, nil
		},
	}
	w := do(t, newBothMux(t, &mockIdentitySvc{}, netSvc), http.MethodGet, "/api/identity/policy", "")
	var res struct {
		AllowAutoGuest bool `json:"allow_auto_guest"`
	}
	json.NewDecoder(w.Body).Decode(&res)
	if !res.AllowAutoGuest {
		t.Error("allow_auto_guest=true attendu pour un réseau public")
	}
}

// / @brief  GET /api/identity/policy expose les règles d'inscription du réseau actif
// / @input  réseau ouvert à l'auto-enregistrement avec approbation exigée
// / @expect allow_auto_register=true et require_approval=true
func TestHandleIdentityPolicy_RequireApproval(t *testing.T) {
	w := do(t, newBothMux(t, &mockIdentitySvc{}, openNetwork(true, "reader")), http.MethodGet, "/api/identity/policy", "")
	var res struct {
		AllowAutoRegister bool `json:"allow_auto_register"`
		RequireApproval   bool `json:"require_approval"`
	}
	decodeJSON(t, w, &res)
	if !res.AllowAutoRegister || !res.RequireApproval {
		t.Errorf("allow_auto_register et require_approval attendus à true, obtenu %+v", res)
	}
}

// ── POST /api/identity/guest ──────────────────────────────────────────────────

// / @brief  POST /api/identity/guest sur un réseau ouvert à la lecture publique : aucune session créée (RM49)
// / @input  réseau AllowAutoGuest=true, corps {pseudo: visiteur}
// / @expect 200, guest=true, role reader, jeton vide — la lecture se fait sans jeton
func TestHandleIdentityGuest_RM49_PublicRead_NoSession(t *testing.T) {
	netSvc := &mockNetworkSvc{
		getActive: func() (*network.NetworkProfile, error) {
			return &network.NetworkProfile{AllowAutoGuest: true, MSPID: "Org1MSP"}, nil
		},
	}

	body, _ := json.Marshal(map[string]string{"pseudo": "visiteur"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/guest", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newBothMux(t, &mockIdentitySvc{}, netSvc).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	var sess struct {
		Token string `json:"token"`
		Role  string `json:"role"`
		Guest bool   `json:"guest"`
	}
	json.NewDecoder(w.Body).Decode(&sess)
	if sess.Token != "" {
		t.Errorf("aucun jeton attendu (lecture publique sans session), obtenu %q", sess.Token)
	}
	if sess.Role != "reader" {
		t.Errorf("role reader attendu, obtenu %q", sess.Role)
	}
	if !sess.Guest {
		t.Error("guest=true attendu")
	}
}

func TestHandleIdentityGuest_DisallowedReturnsForbidden(t *testing.T) {
	netSvc := &mockNetworkSvc{
		getActive: func() (*network.NetworkProfile, error) {
			return &network.NetworkProfile{AllowAutoGuest: false}, nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/identity/guest", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newBothMux(t, &mockIdentitySvc{}, netSvc).ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("code 403 attendu, obtenu %d", w.Code)
	}
}

// ── requireAuth middleware ────────────────────────────────────────────────────

func TestRequireAuth_BlockchainWithoutToken(t *testing.T) {
	// Un accès sans token, quel que soit le réseau configuré, → 401.
	mux := newBlockchainMux(t, &mockSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/refs", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code 401 attendu sans token, obtenu %d", w.Code)
	}
}

func TestRequireAuth_BlockchainWithValidToken(t *testing.T) {
	// Obtenir un token via /api/identity/session puis l'utiliser pour /api/refs.
	idSvc := &mockIdentitySvc{} // Enroll par défaut retourne succès
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	h.WithNetworkInfo(rest.NetworkInfo{Network: "peer0.example.com:7051", Channel: "mychannel"})
	h.WithIdentityService(idSvc)
	srv := rest.NewServer(h, "localhost:0")
	mux, err := srv.Handler()
	if err != nil {
		t.Fatal(err)
	}

	// Étape 1 : créer une session via /api/identity/session
	sessionBody, _ := json.Marshal(map[string]string{"name": "alice", "secret": "x", "org_id": "Org1MSP"})
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(sessionBody))
	r.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, r)
	if w1.Code != http.StatusOK {
		t.Fatalf("création session: got %d, want 200 (body: %s)", w1.Code, w1.Body.String())
	}
	var sess struct {
		Token string `json:"token"`
	}
	json.NewDecoder(w1.Body).Decode(&sess)

	// Étape 2 : accès à /api/refs avec token valide
	r2 := httptest.NewRequest(http.MethodGet, "/api/refs", nil)
	r2.Header.Set("X-Myr-Token", sess.Token)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Errorf("code 200 attendu avec token valide, obtenu %d", w2.Code)
	}
}

// ── POST /api/identity/request ────────────────────────────────────────────────

// openNetwork retourne un service réseau dont le réseau actif accepte
// l'auto-enregistrement avec les règles données.
func openNetwork(requireApproval bool, defaultRole string) *mockNetworkSvc {
	return &mockNetworkSvc{getActive: func() (*network.NetworkProfile, error) {
		return &network.NetworkProfile{Name: "TestNet", AllowAutoRegister: true, RequireApproval: requireApproval, AutoRegisterRole: defaultRole}, nil
	}}
}

// withCertRole fait porter le rôle role au certificat actif enrôlé par svc —
// la session ouverte par loginToken porte alors ce rôle.
func withCertRole(svc *mockIdentitySvc, role string) *mockIdentitySvc {
	svc.enroll = func(_ context.Context, name, _, orgID string) (identity.WalletEntry, error) {
		return identity.WalletEntry{Handle: name + "@" + orgID, Name: name, OrgID: orgID, Status: identity.StatusActive, Role: role}, nil
	}
	return svc
}

// asAdmin fait porter le rôle admin au certificat enrôlé par svc — la session
// ouverte par loginToken dispose alors de identity.admin.
func asAdmin(svc *mockIdentitySvc) *mockIdentitySvc {
	return withCertRole(svc, identity.RoleAdmin)
}

// / @brief  POST /api/identity/request transmet la demande et les règles du réseau actif
// / @input  réseau actif ouvert, approbation exigée, rôle par défaut contributor
// / @expect 201 ; RequestAccess reçoit les champs et la politique ; statut pending et secret renvoyés
func TestHandleIdentityRequest_Success(t *testing.T) {
	var captured identity.AccountRequest
	var policy identity.RegistrationPolicy
	svc := &mockIdentitySvc{
		requestAccess: func(_ context.Context, req identity.AccountRequest, p identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			captured, policy = req, p
			req.ID = "req-001"
			req.Status = identity.RequestPending
			return &req, "s3cret", nil
		},
	}

	body, _ := json.Marshal(map[string]string{
		"pseudo": "bob", "email": "bob@exemple.com", "org_id": "Org1MSP", "message": "accès lecture",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newBothMux(t, svc, openNetwork(true, "contributor")).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("code 201 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	if captured.Pseudo != "bob" || captured.Email != "bob@exemple.com" || captured.OrgID != "Org1MSP" {
		t.Errorf("champs transmis incorrects : %+v", captured)
	}
	want := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true, DefaultRole: "contributor"}
	if policy != want {
		t.Errorf("politique transmise %+v, attendu %+v", policy, want)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != identity.RequestPending || resp["secret"] != "s3cret" {
		t.Errorf("réponse inattendue : %v", resp)
	}
}

// / @brief  POST /api/identity/request est refusé quand l'inscription est fermée (RM43)
// / @input  aucun réseau actif (politique fermée), mock appliquant ErrRegistrationClosed
// / @expect 403
func TestHandleIdentityRequest_RegistrationClosed_Rejected(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"pseudo": "bob", "email": "b@b.com", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newIdentityMux(t, &mockIdentitySvc{}).ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("code 403 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
}

// / @brief  POST /api/identity/request rejette un pseudo invalide
// / @input  RequestAccess renvoie ErrInvalidIdentityName
// / @expect 400
func TestHandleIdentityRequest_InvalidName_Rejected(t *testing.T) {
	svc := &mockIdentitySvc{
		requestAccess: func(context.Context, identity.AccountRequest, identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			return nil, "", identity.ErrInvalidIdentityName
		},
	}
	body, _ := json.Marshal(map[string]string{"pseudo": "bob/../x", "email": "b@b.com", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newBothMux(t, svc, openNetwork(false, "")).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityRequest_MissingFields(t *testing.T) {
	mux := newIdentityMux(t, &mockIdentitySvc{})
	cases := []map[string]string{
		{"pseudo": "", "email": "a@b.com", "org_id": "Org1MSP"},
		{"pseudo": "bob", "email": "", "org_id": "Org1MSP"},
		{"pseudo": "bob", "email": "a@b.com", "org_id": ""},
	}
	for _, c := range cases {
		b, _ := json.Marshal(c)
		req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("code 400 attendu pour %v, obtenu %d", c, w.Code)
		}
	}
}

func TestHandleIdentityRequest_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	newIdentityMux(t, &mockIdentitySvc{}).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

func TestHandleIdentityRequest_ServiceError(t *testing.T) {
	svc := &mockIdentitySvc{
		requestAccess: func(context.Context, identity.AccountRequest, identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			return nil, "", errors.New("disque plein")
		},
	}
	body, _ := json.Marshal(map[string]string{"pseudo": "bob", "email": "b@b.com", "org_id": "Org1MSP"})
	req := httptest.NewRequest(http.MethodPost, "/api/identity/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newBothMux(t, svc, openNetwork(false, "")).ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("code 500 attendu, obtenu %d", w.Code)
	}
}

// ── POST /api/identity/requests/{id}/{approve|reject} ─────────────────────────

// / @brief  POST /api/identity/requests/{id}/approve active l'identité de la demande (RM44)
// / @input  session admin, ApproveRequest renvoie la demande approuvée
// / @expect 200 ; ApproveRequest reçoit l'identifiant de l'URL ; RejectRequest jamais appelé
func TestHandleIdentityRequestDecision_Approve_OK(t *testing.T) {
	var gotID string
	svc := asAdmin(&mockIdentitySvc{
		approveRequest: func(_ context.Context, id string) (*identity.AccountRequest, error) {
			gotID = id
			return &identity.AccountRequest{ID: id, Status: identity.RequestApproved}, nil
		},
		rejectRequest: func(context.Context, string) (*identity.AccountRequest, error) {
			t.Fatal("RejectRequest ne doit pas être appelé")
			return nil, nil
		},
	})
	w := do(t, newIdentityMux(t, svc), http.MethodPost, "/api/identity/requests/req-1/approve", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	if gotID != "req-1" {
		t.Errorf("ApproveRequest appelé avec %q, attendu req-1", gotID)
	}
}

// / @brief  POST /api/identity/requests/{id}/reject suspend l'identité de la demande (RM44)
// / @input  session admin, RejectRequest renvoie la demande refusée
// / @expect 200, statut rejected renvoyé
func TestHandleIdentityRequestDecision_Reject_OK(t *testing.T) {
	w := do(t, newIdentityMux(t, asAdmin(&mockIdentitySvc{})), http.MethodPost, "/api/identity/requests/req-1/reject", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	var resp identity.AccountRequest
	decodeJSON(t, w, &resp)
	if resp.Status != identity.RequestRejected {
		t.Errorf("statut attendu %q, obtenu %q", identity.RequestRejected, resp.Status)
	}
}

// / @brief  La décision sur une demande est réservée à identity.admin
// / @input  session contributor (certificat par défaut du mock)
// / @expect 403, ApproveRequest jamais appelé
func TestHandleIdentityRequestDecision_NonAdmin_Forbidden(t *testing.T) {
	svc := &mockIdentitySvc{approveRequest: func(context.Context, string) (*identity.AccountRequest, error) {
		t.Fatal("ApproveRequest ne doit pas être appelé")
		return nil, nil
	}}
	w := do(t, newIdentityMux(t, svc), http.MethodPost, "/api/identity/requests/req-1/approve", "")
	if w.Code != http.StatusForbidden {
		t.Errorf("code 403 attendu, obtenu %d", w.Code)
	}
}

// / @brief  Erreurs du service traduites en codes HTTP, action inconnue → 404
// / @input  demande introuvable, demande déjà décidée, action "delete"
// / @expect 404, 409, 404
func TestHandleIdentityRequestDecision_Errors(t *testing.T) {
	svc := asAdmin(&mockIdentitySvc{
		approveRequest: func(context.Context, string) (*identity.AccountRequest, error) {
			return nil, identity.ErrRequestNotFound
		},
		rejectRequest: func(context.Context, string) (*identity.AccountRequest, error) {
			return nil, identity.ErrRequestNotPending
		},
	})
	mux := newIdentityMux(t, svc)
	cases := []struct {
		path string
		want int
	}{
		{"/api/identity/requests/req-x/approve", http.StatusNotFound},
		{"/api/identity/requests/req-1/reject", http.StatusConflict},
		{"/api/identity/requests/req-1/delete", http.StatusNotFound},
	}
	for _, c := range cases {
		if w := do(t, mux, http.MethodPost, c.path, ""); w.Code != c.want {
			t.Errorf("%s : code %d attendu, obtenu %d", c.path, c.want, w.Code)
		}
	}
}

// ── POST /api/identity/roles/request ──────────────────────────────────────────

// / @brief  POST /api/identity/roles/request transmet le wallet de la session et le rôle demandé (RM47)
// / @input  session test@TestMSP ; rôle consumer attribué aussitôt, rôle auditor mis en attente
// / @expect 200 (approved) puis 202 (pending) ; RequestRole reçoit le wallet de l'identité de la session
func TestHandleIdentityRoleRequest_RM47_GrantedOrPending(t *testing.T) {
	var gotHandle string
	svc := &mockIdentitySvc{
		wallet: func(handle string) (identity.WalletEntry, error) {
			return identity.WalletEntry{Handle: handle, Status: identity.StatusActive, Role: identity.RoleReader}, nil
		},
		requestRole: func(_ context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error) {
			gotHandle = w.Handle
			status := identity.RequestPending
			if wanted == "consumer" {
				status = identity.RequestApproved
			}
			return &identity.AccountRequest{ID: "req-1", Kind: identity.RequestKindRole, Role: wanted, Status: status}, nil
		},
	}
	mux := newIdentityMux(t, svc)
	for role, want := range map[string]int{"consumer": http.StatusOK, "auditor": http.StatusAccepted} {
		w := do(t, mux, http.MethodPost, "/api/identity/roles/request", `{"role":"`+role+`"}`)
		if w.Code != want {
			t.Fatalf("%s : code %d attendu, obtenu %d — body: %s", role, want, w.Code, w.Body.String())
		}
		var resp identity.AccountRequest
		decodeJSON(t, w, &resp)
		if resp.Role != role || gotHandle != "test@TestMSP" {
			t.Errorf("%s : réponse %+v, wallet %q", role, resp, gotHandle)
		}
	}
}

// / @brief  Erreurs d'une demande de rôle traduites en codes HTTP ; la route exige une session
// / @input  rôle déjà détenu, rôle inconnu, demande déjà en attente, identité suspendue ; requête sans jeton
// / @expect 400, 400, 409, 403 ; 401 sans jeton
func TestHandleIdentityRoleRequest_Errors(t *testing.T) {
	errs := map[string]error{
		"reader":      identity.ErrRoleAlreadyHeld,
		"ghost":       identity.ErrInvalidRole,
		"auditor":     identity.ErrRequestAlreadyPending,
		"contributor": identity.ErrIdentitySuspended,
	}
	svc := &mockIdentitySvc{
		wallet: func(handle string) (identity.WalletEntry, error) { return identity.WalletEntry{Handle: handle}, nil },
		requestRole: func(_ context.Context, _ identity.WalletEntry, wanted string) (*identity.AccountRequest, error) {
			return nil, errs[wanted]
		},
	}
	mux := newIdentityMux(t, svc)
	want := map[string]int{"reader": http.StatusBadRequest, "ghost": http.StatusBadRequest, "auditor": http.StatusConflict, "contributor": http.StatusForbidden}
	for role, code := range want {
		if w := do(t, mux, http.MethodPost, "/api/identity/roles/request", `{"role":"`+role+`"}`); w.Code != code {
			t.Errorf("%s : code %d attendu, obtenu %d", role, code, w.Code)
		}
	}
	r := httptest.NewRequest(http.MethodPost, "/api/identity/roles/request", strings.NewReader(`{"role":"auditor"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("sans jeton : code 401 attendu, obtenu %d", w.Code)
	}
}

// ── POST /api/identity/accounts ───────────────────────────────────────────────

// / @brief  POST /api/identity/accounts crée un compte actif à l'initiative de l'admin
// / @input  session admin, corps {pseudo, org_id, role, email}
// / @expect 201 {handle, secret} ; CreateAccount reçoit les champs du corps
func TestHandleIdentityAccounts_Admin_Created(t *testing.T) {
	var got []string
	svc := asAdmin(&mockIdentitySvc{
		createAccount: func(_ context.Context, pseudo, orgID, role, email string) (string, string, error) {
			got = []string{pseudo, orgID, role, email}
			return identity.Handle(pseudo, orgID), "s3cret", nil
		},
	})
	w := do(t, newIdentityMux(t, svc), http.MethodPost, "/api/identity/accounts",
		`{"pseudo":"alice","org_id":"Org1MSP","role":"contributor","email":"alice@example.com"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("code 201 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	if strings.Join(got, ",") != "alice,Org1MSP,contributor,alice@example.com" {
		t.Errorf("champs transmis à CreateAccount : %v", got)
	}
	var resp map[string]string
	decodeJSON(t, w, &resp)
	if resp["handle"] != "alice@Org1" || resp["secret"] != "s3cret" {
		t.Errorf("réponse inattendue : %v", resp)
	}
}

// / @brief  La création de compte est réservée à identity.admin
// / @input  session contributor
// / @expect 403, CreateAccount jamais appelé
func TestHandleIdentityAccounts_NonAdmin_Forbidden(t *testing.T) {
	svc := &mockIdentitySvc{createAccount: func(context.Context, string, string, string, string) (string, string, error) {
		t.Fatal("CreateAccount ne doit pas être appelé")
		return "", "", nil
	}}
	w := do(t, newIdentityMux(t, svc), http.MethodPost, "/api/identity/accounts", `{"pseudo":"alice","org_id":"Org1MSP"}`)
	if w.Code != http.StatusForbidden {
		t.Errorf("code 403 attendu, obtenu %d", w.Code)
	}
}

// / @brief  POST /api/identity/accounts rejette un pseudo invalide ou un rôle absent de la table (RM46)
// / @input  session admin, CreateAccount renvoie ErrInvalidIdentityName puis ErrInvalidRole
// / @expect 400 dans les deux cas
func TestHandleIdentityAccounts_InvalidNameOrRole_Rejected(t *testing.T) {
	for _, want := range []error{identity.ErrInvalidIdentityName, identity.ErrInvalidRole} {
		svc := asAdmin(&mockIdentitySvc{createAccount: func(context.Context, string, string, string, string) (string, string, error) {
			return "", "", want
		}})
		w := do(t, newIdentityMux(t, svc), http.MethodPost, "/api/identity/accounts", `{"pseudo":"alice","org_id":"Org1MSP","role":"x"}`)
		if w.Code != http.StatusBadRequest {
			t.Errorf("%v : code 400 attendu, obtenu %d", want, w.Code)
		}
	}
}

// ── GET /api/identity/requests ────────────────────────────────────────────────

// / @brief  GET /api/identity/requests expose des données personnelles (email, message) de
// / toutes les demandes de compte — réservé à identity.admin, jamais à un simple contributor.
// / @input  session "contributor" (rôle obtenu via loginToken, pas de RoleService injecté)
// / @expect 403 : requireRole retombe sur legacyRoleRank, qui place contributor sous admin
func TestHandleIdentityRequests_NonAdmin_Forbidden(t *testing.T) {
	svc := &mockIdentitySvc{
		listRequests: func() ([]*identity.AccountRequest, error) {
			return []*identity.AccountRequest{
				{ID: "req-1", Pseudo: "alice", Status: identity.RequestPending},
			}, nil
		},
	}
	w := do(t, newIdentityMux(t, svc), http.MethodGet, "/api/identity/requests", "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("code 403 attendu (permission identity.admin manquante), obtenu %d", w.Code)
	}
}

// / @brief  GET /api/identity/requests retourne la liste quand le rôle porte identity.admin
// / @input  RoleService accordant rbac.PermIdentityAdmin à toute session authentifiée
// / @expect 200, liste des demandes renvoyée telle quelle
func TestHandleIdentityRequests_Admin_ReturnsList(t *testing.T) {
	svc := &mockIdentitySvc{
		listRequests: func() ([]*identity.AccountRequest, error) {
			return []*identity.AccountRequest{
				{ID: "req-1", Pseudo: "alice", Status: identity.RequestPending},
				{ID: "req-2", Pseudo: "bob", Status: identity.RequestApproved},
			}, nil
		},
	}
	roleSvc := &mockRoleSvc{hasPermission: func(string, rbac.Permission) bool { return true }}
	w := do(t, newIdentityMuxWithRole(t, svc, roleSvc), http.MethodGet, "/api/identity/requests", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d", w.Code)
	}
	var list []*identity.AccountRequest
	json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 2 {
		t.Errorf("2 demandes attendues, obtenu %d", len(list))
	}
}

// / @brief  GET /api/identity/requests renvoie [] et non null quand aucune demande n'existe
// / @input  RoleService accordant identity.admin, ListRequests retourne (nil, nil)
// / @expect 200, tableau JSON vide
func TestHandleIdentityRequests_Admin_EmptyList(t *testing.T) {
	roleSvc := &mockRoleSvc{hasPermission: func(string, rbac.Permission) bool { return true }}
	w := do(t, newIdentityMuxWithRole(t, &mockIdentitySvc{}, roleSvc), http.MethodGet, "/api/identity/requests", "")
	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d", w.Code)
	}
	var list []*identity.AccountRequest
	json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 0 {
		t.Errorf("liste vide attendue, obtenu %d", len(list))
	}
}

// ── POST /api/identity/session ────────────────────────────────────────────────

// / @brief  POST /api/identity/session crée une session et retourne un token
// / @input  service Enroll retourne succès, champs name/secret/org_id valides
// / @expect HTTP 200, token non vide, role="contributor"
func TestHandleIdentitySession_Success(t *testing.T) {
	svc := &mockIdentitySvc{}
	mux := newIdentityMux(t, svc)

	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "s3cr3t", "org_id": "Org1MSP"})
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("code 200 attendu, obtenu %d — body: %s", w.Code, w.Body.String())
	}
	var sess struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(w.Body).Decode(&sess); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if sess.Token == "" {
		t.Error("token attendu dans la réponse")
	}
	if sess.Role != "contributor" {
		t.Errorf("role attendu contributor, obtenu %q", sess.Role)
	}
}

// / @brief  POST /api/identity/session retourne 400 si des champs obligatoires sont absents
// / @input  cas : name vide, secret vide, org_id vide
// / @expect HTTP 400 dans chaque cas
func TestHandleIdentitySession_MissingFields(t *testing.T) {
	mux := newIdentityMux(t, &mockIdentitySvc{})
	cases := []map[string]string{
		{"name": "", "secret": "x", "org_id": "Org1MSP"},
		{"name": "alice", "secret": "", "org_id": "Org1MSP"},
		{"name": "alice", "secret": "x", "org_id": ""},
	}
	for _, bodyMap := range cases {
		b, _ := json.Marshal(bodyMap)
		r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(b))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("code 400 attendu pour %v, obtenu %d", bodyMap, w.Code)
		}
	}
}

// / @brief  POST /api/identity/session retourne 400 pour un corps JSON invalide
// / @input  corps "{bad"
// / @expect HTTP 400
func TestHandleIdentitySession_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	newIdentityMux(t, &mockIdentitySvc{}).ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("code 400 attendu, obtenu %d", w.Code)
	}
}

// / @brief  POST /api/identity/session retourne 503 si aucun service identité n'est configuré
// / @input  handler sans WithIdentityService, corps valide
// / @expect HTTP 503
func TestHandleIdentitySession_NoIdentityService(t *testing.T) {
	store := &mockStore{}
	h := rest.NewHandler(&mockSvc{}, store, store)
	srv := rest.NewServer(h, "localhost:0")
	mux, _ := srv.Handler()

	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "x", "org_id": "Org1MSP"})
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("code 503 attendu, obtenu %d", w.Code)
	}
}

// / @brief  GET /api/identity/session retourne 405 (seul POST est autorisé)
// / @input  GET /api/identity/session
// / @expect HTTP 405
func TestHandleIdentitySession_MethodNotAllowed(t *testing.T) {
	w := do(t, newIdentityMux(t, &mockIdentitySvc{}), http.MethodGet, "/api/identity/session", "")
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("code 405 attendu, obtenu %d", w.Code)
	}
}

// / @brief  POST /api/identity/session retourne 401 si Enroll échoue
// / @input  mockIdentitySvc.Enroll retourne une erreur "secret invalide"
// / @expect HTTP 401
func TestHandleIdentitySession_EnrollError(t *testing.T) {
	svc := &mockIdentitySvc{
		enroll: func(_ context.Context, _, _, _ string) (identity.WalletEntry, error) {
			return identity.WalletEntry{}, errors.New("secret invalide")
		},
	}
	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "wrong", "org_id": "Org1MSP"})
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newIdentityMux(t, svc).ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code 401 attendu, obtenu %d", w.Code)
	}
}

// / @brief  POST /api/identity/session n'ouvre aucune session pour une identité inactive ou sans rôle (RM45)
// / @input  certificat enrôlé : en attente, suspendu, statut absent, ou actif sans Myr.role
// / @expect HTTP 403 dans chaque cas, aucun token délivré
func TestHandleIdentitySession_CannotSignIn_Forbidden(t *testing.T) {
	cases := map[string]identity.WalletEntry{
		"pending":   {Status: identity.StatusPending, Role: identity.RoleContributor},
		"suspended": {Status: identity.StatusSuspended, Role: identity.RoleContributor},
		"no status": {Role: identity.RoleContributor},
		"no role":   {Status: identity.StatusActive},
	}
	for name, entry := range cases {
		svc := &mockIdentitySvc{enroll: func(_ context.Context, n, _, orgID string) (identity.WalletEntry, error) {
			entry.Handle, entry.Name, entry.OrgID = n+"@"+orgID, n, orgID
			return entry, nil
		}}
		body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "s3cr3t", "org_id": "Org1MSP"})
		r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		newIdentityMux(t, svc).ServeHTTP(w, r)
		if w.Code != http.StatusForbidden {
			t.Errorf("%s : code 403 attendu, obtenu %d — body: %s", name, w.Code, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "token") {
			t.Errorf("%s : aucun token ne doit être délivré, body: %s", name, w.Body.String())
		}
	}
}

// / @brief  La session porte le rôle inscrit dans le certificat enrôlé (RM45)
// / @input  certificat actif portant Myr.role=admin, aucun RoleService injecté
// / @expect role="admin" dans la réponse ; la session accède à une route identity.admin
func TestHandleIdentitySession_RoleFromCertificate(t *testing.T) {
	mux := newIdentityMux(t, asAdmin(&mockIdentitySvc{}))
	body, _ := json.Marshal(map[string]string{"name": "alice", "secret": "s3cr3t", "org_id": "Org1MSP"})
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var sess struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
	decodeJSON(t, w, &sess)
	if sess.Role != identity.RoleAdmin {
		t.Fatalf("role attendu admin, obtenu %q", sess.Role)
	}
	if w := do(t, mux, http.MethodGet, "/api/identity/requests", ""); w.Code != http.StatusOK {
		t.Errorf("code 200 attendu sur une route identity.admin, obtenu %d", w.Code)
	}
}

// / @brief  Une session ouverte avec un certificat reader ne peut pas écrire (RM21)
// / @input  certificat actif portant Myr.role=reader ; lecture puis création d'un composant
// / @expect lecture 200 ; écriture 403 (permission write absente)
func TestHandleIdentitySession_ReaderCannotWrite(t *testing.T) {
	mux := newIdentityMux(t, withCertRole(&mockIdentitySvc{}, identity.RoleReader))
	if w := do(t, mux, http.MethodGet, "/api/refs", ""); w.Code != http.StatusOK {
		t.Fatalf("lecture : code 200 attendu, obtenu %d", w.Code)
	}
	w := do(t, mux, http.MethodPost, "/api/components", `{"name":"bracket","category":"base"}`)
	if w.Code != http.StatusForbidden || !strings.Contains(w.Body.String(), "droits insuffisants") {
		t.Errorf("écriture : 403 « droits insuffisants » attendu pour un reader, obtenu %d — body: %s", w.Code, w.Body.String())
	}
}

// ── /api/channels PUT ─────────────────────────────────────────────────────────

// / @brief  PUT /api/channels sans token retourne 401
// / @input  PUT /api/channels sans en-tête X-Myr-Token
// / @expect HTTP 401
func TestChannels_PUT_WithoutToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodPut, "/api/channels", strings.NewReader(`{"channel":"sandbox"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newIdentityMux(t, &mockIdentitySvc{}).ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code 401 attendu sans token, obtenu %d", w.Code)
	}
}

// / @brief  PUT /api/channels avec token valide change le canal de la session
// / @input  POST /api/identity/session pour obtenir un token, puis PUT /api/channels
// / @expect HTTP 200, channel="sandbox" dans la réponse
func TestChannels_PUT_SwitchesChannel(t *testing.T) {
	mux := newIdentityMux(t, &mockIdentitySvc{})

	// Étape 1 : créer une session
	sessionBody, _ := json.Marshal(map[string]string{"name": "alice", "secret": "x", "org_id": "Org1MSP"})
	r1 := httptest.NewRequest(http.MethodPost, "/api/identity/session", bytes.NewReader(sessionBody))
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, r1)
	if w1.Code != http.StatusOK {
		t.Fatalf("création session: got %d, want 200 (body: %s)", w1.Code, w1.Body.String())
	}
	var sess struct {
		Token string `json:"token"`
	}
	json.NewDecoder(w1.Body).Decode(&sess)
	if sess.Token == "" {
		t.Fatal("token absent dans la réponse de session")
	}

	// Étape 2 : changer le canal avec le token
	r2 := httptest.NewRequest(http.MethodPut, "/api/channels", strings.NewReader(`{"channel":"sandbox"}`))
	r2.Header.Set("Content-Type", "application/json")
	r2.Header.Set("X-Myr-Token", sess.Token)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("PUT channels: got %d, want 200 (body: %s)", w2.Code, w2.Body.String())
	}
	var resp struct {
		Channel string `json:"channel"`
	}
	json.NewDecoder(w2.Body).Decode(&resp)
	if resp.Channel != "sandbox" {
		t.Errorf("channel: got %q, want sandbox", resp.Channel)
	}
}
