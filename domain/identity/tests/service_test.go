// domain/identity/tests/service_test.go
package identity_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"myr-core/domain/identity"
	"myr-core/domain/network"
)

type stubCA struct {
	registerFn      func(ctx context.Context, req identity.RegisterRequest) (string, error)
	enrollFn        func(ctx context.Context, name, secret string) (string, string, string, error)
	getStatusFn     func(ctx context.Context, name string) (string, error)
	reEnrollFn      func(ctx context.Context, name, certPEM, keyPEM string) (string, error)
	updateAttrsFn   func(ctx context.Context, name string, attrs map[string]string) error
	certAttrs       map[string]string // attributs portés par les certificats émis
	lastRegisterReq identity.RegisterRequest
	lastEnrollName  string
	lastStatusName  string
	lastUpdateName  string
	lastUpdateAttrs map[string]string
}

func (s *stubCA) Register(ctx context.Context, req identity.RegisterRequest) (string, error) {
	s.lastRegisterReq = req
	if s.registerFn != nil {
		return s.registerFn(ctx, req)
	}
	return "secret123", nil
}
func (s *stubCA) Enroll(ctx context.Context, name, secret string) (string, string, string, error) {
	s.lastEnrollName = name
	if s.enrollFn != nil {
		return s.enrollFn(ctx, name, secret)
	}
	return fakeCertPEM, fakeKeyPEM, fakeCACertPEM, nil
}
func (s *stubCA) GetStatus(ctx context.Context, name string) (string, error) {
	s.lastStatusName = name
	if s.getStatusFn != nil {
		return s.getStatusFn(ctx, name)
	}
	return identity.StatusActive, nil
}
func (s *stubCA) ReEnroll(ctx context.Context, name, certPEM, keyPEM string) (string, error) {
	if s.reEnrollFn != nil {
		return s.reEnrollFn(ctx, name, certPEM, keyPEM)
	}
	return fakeCertPEM, nil
}
func (s *stubCA) UpdateAttributes(ctx context.Context, name string, attrs map[string]string) error {
	s.lastUpdateName = name
	s.lastUpdateAttrs = attrs
	if s.updateAttrsFn != nil {
		return s.updateAttrsFn(ctx, name, attrs)
	}
	return nil
}
func (s *stubCA) CertAttributes(certPEM string) (map[string]string, error) {
	attrs := map[string]string{}
	for k, v := range s.certAttrs {
		attrs[k] = v
	}
	return attrs, nil
}

var _ identity.CAPort = (*stubCA)(nil)

const fakeCertPEM = "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0fake==\n-----END CERTIFICATE-----\n"
const fakeKeyPEM = "-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEIFake/key/data==\n-----END EC PRIVATE KEY-----\n"
const fakeCACertPEM = "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAcafake==\n-----END CERTIFICATE-----\n"

// / @brief  Création de compte échoue sans CA configurée
// / @input  service sans CA (nil), pseudo "alice" dans "Org1MSP"
// / @expect retourne une erreur indiquant l'absence de CA
func TestService_CreateAccount_NoCA(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	if _, _, err := svc.CreateAccount(context.Background(), "alice", "Org1MSP", identity.RoleReader, ""); err == nil {
		t.Fatal("erreur attendue sans CA")
	}
}

// / @brief  Création de compte par l'admin : identité active avec le rôle donné
// / @input  stubCA retournant "secret123", pseudo "alice", org "Org1MSP", rôle contributor, email
// / @expect retourne ("alice@Org1", "secret123") ; la CA reçoit le nom complet, le statut active, le rôle et l'email
func TestService_CreateAccount_Success(t *testing.T) {
	ca := &stubCA{}
	svc := identity.NewService(t.TempDir(), ca)
	handle, secret, err := svc.CreateAccount(context.Background(), "alice", "Org1MSP", identity.RoleContributor, "alice@example.com")
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if handle != "alice@Org1" || secret != "secret123" {
		t.Errorf("attendu (alice@Org1, secret123), obtenu (%q, %q)", handle, secret)
	}
	got := ca.lastRegisterReq
	if got.Name != "alice@Org1" || got.OrgID != "Org1MSP" || got.Status != identity.StatusActive ||
		got.Role != identity.RoleContributor || got.Email != "alice@example.com" {
		t.Errorf("requête CA inattendue : %+v", got)
	}
}

// / @brief  Création de compte propage l'erreur retournée par la CA
// / @input  stubCA dont registerFn retourne "CA: identite deja enregistree"
// / @expect retourne une erreur enveloppant celle de la CA
func TestService_CreateAccount_CAError(t *testing.T) {
	expected := errors.New("CA: identite deja enregistree")
	ca := &stubCA{registerFn: func(_ context.Context, _ identity.RegisterRequest) (string, error) { return "", expected }}
	svc := identity.NewService(t.TempDir(), ca)
	_, _, err := svc.CreateAccount(context.Background(), "alice", "Org1MSP", "", "")
	if !errors.Is(err, expected) {
		t.Errorf("erreur attendue %v, obtenu %v", expected, err)
	}
}

// / @brief  Création de compte refuse un pseudo ou une organisation qui ne forme pas un identifiant sûr
// / @input  pseudos "", "..", "a@b", "a/b", "a b" ; organisation "Org:1"
// / @expect ErrInvalidIdentityName dans chaque cas, la CA n'est jamais appelée
func TestService_CreateAccount_InvalidName_Rejected(t *testing.T) {
	ca := &stubCA{registerFn: func(context.Context, identity.RegisterRequest) (string, error) {
		t.Fatal("la CA ne doit pas être appelée")
		return "", nil
	}}
	svc := identity.NewService(t.TempDir(), ca)
	cases := [][2]string{{"", "Org1MSP"}, {"..", "Org1MSP"}, {"a@b", "Org1MSP"}, {"a/b", "Org1MSP"}, {"a b", "Org1MSP"}, {"alice", "Org:1"}}
	for _, c := range cases {
		if _, _, err := svc.CreateAccount(context.Background(), c[0], c[1], "", ""); !errors.Is(err, identity.ErrInvalidIdentityName) {
			t.Errorf("%q/%q : ErrInvalidIdentityName attendu, obtenu %v", c[0], c[1], err)
		}
	}
}

// / @brief  Enrôlement échoue sans CA configurée
// / @input  service sans CA (nil), identité "alice" avec secret "secret123"
// / @expect retourne une erreur indiquant l'absence de CA
func TestService_Enroll_NoCA(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	_, err := svc.Enroll(context.Background(), "alice", "secret123", "Org1MSP")
	if err == nil {
		t.Fatal("erreur attendue sans CA")
	}
}

// / @brief  Enrôlement réussi avec écriture des fichiers MSP sur disque
// / @input  stubCA retournant des PEM valides, répertoire temporaire, identité "alice" dans "Org1MSP"
// / @expect retourne un WalletEntry avec handle "alice@Org1", statut et rôle lus dans le certificat, fichiers cert/key/ca présents
func TestService_Enroll_Success(t *testing.T) {
	walletDir := t.TempDir()
	ca := &stubCA{certAttrs: map[string]string{identity.AttrStatus: identity.StatusPending, identity.AttrRole: identity.RoleReader}}
	svc := identity.NewService(walletDir, ca)
	entry, err := svc.Enroll(context.Background(), "alice", "secret123", "Org1MSP")
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if entry.Handle != "alice@Org1" {
		t.Errorf("handle attendu %q, obtenu %q", "alice@Org1", entry.Handle)
	}
	if entry.Status != identity.StatusPending || entry.Role != identity.RoleReader {
		t.Errorf("attendu statut %q et rôle %q, obtenu %q / %q", identity.StatusPending, identity.RoleReader, entry.Status, entry.Role)
	}
	mspDir := filepath.Join(walletDir, "alice@Org1", "msp")
	for _, path := range []string{
		filepath.Join(mspDir, "signcerts", "cert.pem"),
		filepath.Join(mspDir, "keystore", "key.pem"),
		filepath.Join(mspDir, "cacerts", "ca.pem"),
	} {
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			t.Errorf("fichier absent : %s", path)
		}
	}
}

// / @brief  Enrôlement propage l'erreur retournée par la CA
// / @input  stubCA dont enrollFn retourne "CA: secret invalide"
// / @expect retourne exactement l'erreur de la CA
func TestService_Enroll_CAError(t *testing.T) {
	expected := errors.New("CA: secret invalide")
	ca := &stubCA{enrollFn: func(_ context.Context, _, _ string) (string, string, string, error) { return "", "", "", expected }}
	svc := identity.NewService(t.TempDir(), ca)
	_, err := svc.Enroll(context.Background(), "alice", "wrong", "Org1MSP")
	if !errors.Is(err, expected) {
		t.Errorf("erreur attendue %v, obtenu %v", expected, err)
	}
}

// / @brief  Enrôlement authentifie auprès de la CA avec l'identifiant complet pseudo@org, pas le pseudo seul
// / @input  stubCA, appel Enroll(ctx, "alice", "secret123", "Org1MSP")
// / @expect la CA reçoit "alice@Org1" comme nom — l'identifiant sous lequel AutoRegister a créé l'identité
func TestService_Enroll_UsesFullHandleForCAAuth(t *testing.T) {
	ca := &stubCA{}
	svc := identity.NewService(t.TempDir(), ca)
	if _, err := svc.Enroll(context.Background(), "alice", "secret123", "Org1MSP"); err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if ca.lastEnrollName != "alice@Org1" {
		t.Errorf("nom transmis à la CA : got %q, want %q", ca.lastEnrollName, "alice@Org1")
	}
}

// / @brief  GetStatus échoue sans CA configurée
// / @input  service sans CA (nil), entrée de wallet {Name:"alice"}
// / @expect retourne une erreur indiquant l'absence de CA
func TestService_GetStatus_NoCA(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	_, err := svc.GetStatus(context.Background(), identity.WalletEntry{Name: "alice"})
	if err == nil {
		t.Fatal("erreur attendue sans CA")
	}
}

// / @brief  GetStatus retourne le statut "active" renvoyé par la CA
// / @input  stubCA dont getStatusFn retourne StatusActive, entrée {Name:"alice"}
// / @expect retourne identity.StatusActive sans erreur
func TestService_GetStatus_ReturnsActiveAfterApproval(t *testing.T) {
	ca := &stubCA{getStatusFn: func(_ context.Context, _ string) (string, error) { return identity.StatusActive, nil }}
	svc := identity.NewService(t.TempDir(), ca)
	status, err := svc.GetStatus(context.Background(), identity.WalletEntry{Name: "alice"})
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if status != identity.StatusActive {
		t.Errorf("attendu %q, obtenu %q", identity.StatusActive, status)
	}
}

// / @brief  GetStatus retourne le statut "pending" renvoyé par la CA
// / @input  stubCA dont getStatusFn retourne StatusPending, entrée {Name:"bob"}
// / @expect retourne identity.StatusPending sans erreur
func TestService_GetStatus_Pending(t *testing.T) {
	ca := &stubCA{getStatusFn: func(_ context.Context, _ string) (string, error) { return identity.StatusPending, nil }}
	svc := identity.NewService(t.TempDir(), ca)
	status, err := svc.GetStatus(context.Background(), identity.WalletEntry{Name: "bob"})
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if status != identity.StatusPending {
		t.Errorf("attendu %q, obtenu %q", identity.StatusPending, status)
	}
}

// / @brief  GetStatus propage l'erreur retournée par la CA
// / @input  stubCA dont getStatusFn retourne "CA: identite introuvable", entrée {Name:"ghost"}
// / @expect retourne exactement l'erreur de la CA
func TestService_GetStatus_CAError(t *testing.T) {
	expected := errors.New("CA: identite introuvable")
	ca := &stubCA{getStatusFn: func(_ context.Context, _ string) (string, error) { return "", expected }}
	svc := identity.NewService(t.TempDir(), ca)
	_, err := svc.GetStatus(context.Background(), identity.WalletEntry{Name: "ghost"})
	if !errors.Is(err, expected) {
		t.Errorf("erreur attendue %v, obtenu %v", expected, err)
	}
}

// / @brief  GetStatus interroge la CA avec l'identifiant complet pseudo@org (wallet.Handle), pas wallet.Name seul
// / @input  stubCA, appel GetStatus avec WalletEntry{Handle:"alice@Org1", Name:"alice", OrgID:"Org1MSP"}
// / @expect la CA reçoit "alice@Org1" comme nom
func TestService_GetStatus_UsesHandleNotBareName(t *testing.T) {
	ca := &stubCA{}
	svc := identity.NewService(t.TempDir(), ca)
	_, err := svc.GetStatus(context.Background(), identity.WalletEntry{Handle: "alice@Org1", Name: "alice", OrgID: "Org1MSP"})
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if ca.lastStatusName != "alice@Org1" {
		t.Errorf("nom transmis à la CA : got %q, want %q", ca.lastStatusName, "alice@Org1")
	}
}

// / @brief  Listage des wallets locaux sur un répertoire inexistant retourne nil
// / @input  service avec répertoire "/no/such/dir" inexistant
// / @expect retourne nil sans erreur
func TestService_ListLocalWallets_NonExistentDir(t *testing.T) {
	svc := identity.NewService("/no/such/dir", nil)
	wallets, err := svc.ListLocalWallets()
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if wallets != nil {
		t.Errorf("nil attendu, obtenu %v", wallets)
	}
}

// / @brief  Listage des wallets locaux sur un répertoire vide retourne une liste vide
// / @input  service avec répertoire temporaire vide
// / @expect retourne une liste de longueur 0 sans erreur
func TestService_ListLocalWallets_EmptyDir(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	wallets, err := svc.ListLocalWallets()
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if len(wallets) != 0 {
		t.Errorf("0 wallet attendu, obtenu %d", len(wallets))
	}
}

// / @brief  Listage des wallets locaux retourne les entrées créées après enrôlement
// / @input  répertoire temporaire, deux enrôlements successifs "alice@Org1" et "bob@Org2"
// / @expect retourne 2 wallets avec les handles attendus
func TestService_ListLocalWallets_AfterEnroll(t *testing.T) {
	walletDir := t.TempDir()
	svc := identity.NewService(walletDir, &stubCA{})
	if _, err := svc.Enroll(context.Background(), "alice", "s1", "Org1MSP"); err != nil {
		t.Fatalf("enroll alice: %v", err)
	}
	if _, err := svc.Enroll(context.Background(), "bob", "s2", "Org2MSP"); err != nil {
		t.Fatalf("enroll bob: %v", err)
	}
	wallets, err := svc.ListLocalWallets()
	if err != nil {
		t.Fatalf("inattendu : %v", err)
	}
	if len(wallets) != 2 {
		t.Fatalf("2 wallets attendus, obtenu %d", len(wallets))
	}
	handles := map[string]bool{}
	for _, w := range wallets {
		handles[w.Handle] = true
	}
	for _, expected := range []string{"alice@Org1", "bob@Org2"} {
		if !handles[expected] {
			t.Errorf("wallet %q absent", expected)
		}
	}
}

// / @brief  Flux complet d'inscription avec approbation : demande, enrôlement, approbation, connexion
// / @details Simule une CA qui retient les attributs de chaque identité et les inscrit dans les
// / certificats qu'elle émet — le statut d'une identité ne change qu'au prochain enrôlement
// / @input  réseau ouvert exigeant une approbation, rôle par défaut reader, identité "carol@Org1"
// / @expect avant approbation le certificat porte pending et la connexion est refusée ;
// / après approbation et nouvel enrôlement il porte active et la connexion est permise
func TestAccountRequest_FullFlow(t *testing.T) {
	registry := map[string]map[string]string{} // attributs CA par identité
	ca := &stubCA{}
	ca.registerFn = func(_ context.Context, req identity.RegisterRequest) (string, error) {
		registry[req.Name] = map[string]string{identity.AttrStatus: req.Status, identity.AttrRole: req.Role}
		return "enrollment-secret", nil
	}
	ca.enrollFn = func(_ context.Context, name, secret string) (string, string, string, error) {
		if secret != "enrollment-secret" {
			return "", "", "", errors.New("secret invalide")
		}
		ca.certAttrs = registry[name]
		return fakeCertPEM, fakeKeyPEM, fakeCACertPEM, nil
	}
	ca.updateAttrsFn = func(_ context.Context, name string, attrs map[string]string) error {
		for k, v := range attrs {
			registry[name][k] = v
		}
		return nil
	}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(&mockRequestStore{})
	policy := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true, DefaultRole: identity.RoleReader}

	req, secret, err := svc.RequestAccess(context.Background(),
		identity.AccountRequest{Pseudo: "carol", Email: "carol@example.com", OrgID: "Org1MSP"}, policy)
	if err != nil {
		t.Fatalf("RequestAccess: %v", err)
	}
	entry, err := svc.Enroll(context.Background(), "carol", secret, "Org1MSP")
	if err != nil {
		t.Fatalf("Enroll: %v", err)
	}
	if !errors.Is(entry.CanSignIn(), identity.ErrIdentityPending) {
		t.Fatalf("connexion refusée attendue avant approbation, obtenu %v (statut %q)", entry.CanSignIn(), entry.Status)
	}

	if _, err := svc.ApproveRequest(context.Background(), req.ID); err != nil {
		t.Fatalf("ApproveRequest: %v", err)
	}
	entry, err = svc.Enroll(context.Background(), "carol", secret, "Org1MSP")
	if err != nil {
		t.Fatalf("Enroll après approbation: %v", err)
	}
	if err := entry.CanSignIn(); err != nil {
		t.Fatalf("connexion attendue après approbation, obtenu %v", err)
	}
	if entry.Role != identity.RoleReader {
		t.Errorf("rôle attendu %q, obtenu %q", identity.RoleReader, entry.Role)
	}
}

// ── Mock RequestStore ─────────────────────────────────────────────────────────

// mockRequestStore remplace une demande de même ID, comme le stockage JSON,
// et en conserve une copie pour ne pas partager l'état de l'appelant.
type mockRequestStore struct {
	requests []*identity.AccountRequest
	err      error
}

func (m *mockRequestStore) Save(req *identity.AccountRequest) error {
	if m.err != nil {
		return m.err
	}
	cp := *req
	for i, r := range m.requests {
		if r.ID == req.ID {
			m.requests[i] = &cp
			return nil
		}
	}
	m.requests = append(m.requests, &cp)
	return nil
}

func (m *mockRequestStore) FindAll() ([]*identity.AccountRequest, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.requests, nil
}

// ── ReEnroll ──────────────────────────────────────────────────────────────────

// / @brief  ReEnroll échoue sans CA configurée
// / @input  service sans CA (nil), wallet {Name:"alice"}
// / @expect retourne une erreur indiquant l'absence de CA
func TestReEnroll_NoCA(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	_, err := svc.ReEnroll(context.Background(), identity.WalletEntry{Name: "alice"})
	if err == nil {
		t.Fatal("erreur attendue sans CA")
	}
}

// / @brief  ReEnroll réussi renouvelle le certificat et l'écrit sur disque
// / @input  wallet enrôlé préalablement, stubCA retournant un nouveau PEM
// / @expect WalletEntry conserve le même handle, fichier cert.pem mis à jour
func TestReEnroll_Success(t *testing.T) {
	walletDir := t.TempDir()
	newCert := "-----BEGIN CERTIFICATE-----\nMIIBnewcert==\n-----END CERTIFICATE-----\n"
	ca := &stubCA{
		reEnrollFn: func(_ context.Context, _, _, _ string) (string, error) {
			return newCert, nil
		},
	}
	svc := identity.NewService(walletDir, ca)
	entry, err := svc.Enroll(context.Background(), "alice", "secret", "Org1MSP")
	if err != nil {
		t.Fatalf("Enroll préalable : %v", err)
	}

	updated, err := svc.ReEnroll(context.Background(), entry)
	if err != nil {
		t.Fatalf("ReEnroll : %v", err)
	}
	if updated.Handle != entry.Handle {
		t.Errorf("Handle: got %q, want %q", updated.Handle, entry.Handle)
	}
	data, err := os.ReadFile(filepath.Join(entry.MSPDir, "signcerts", "cert.pem"))
	if err != nil {
		t.Fatalf("lecture cert : %v", err)
	}
	if string(data) != newCert {
		t.Error("cert.pem non mis à jour après ReEnroll")
	}
}

// / @brief  ReEnroll propage l'erreur retournée par la CA
// / @input  stubCA dont reEnrollFn retourne "CA: cert invalide"
// / @expect retourne exactement l'erreur de la CA
func TestReEnroll_CAError(t *testing.T) {
	expected := errors.New("CA: cert invalide")
	walletDir := t.TempDir()
	ca := &stubCA{
		reEnrollFn: func(_ context.Context, _, _, _ string) (string, error) { return "", expected },
	}
	svc := identity.NewService(walletDir, ca)
	entry, _ := svc.Enroll(context.Background(), "alice", "secret", "Org1MSP")
	_, err := svc.ReEnroll(context.Background(), entry)
	if !errors.Is(err, expected) {
		t.Errorf("erreur attendue %v, obtenu %v", expected, err)
	}
}

// ── SetRole ───────────────────────────────────────────────────────────────────

// / @brief  SetRole refuse sans CA configuré
// / @input  service construit sans CAPort (nil)
// / @expect retourne une erreur explicite
func TestSetRole_NoCA(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	err := svc.SetRole(context.Background(), "alice@org1", "contributor")
	if err == nil {
		t.Fatal("erreur attendue sans CA")
	}
}

// / @brief  SetRole transmet le nouvel attribut Myr.role à la CA
// / @input  stubCA, identité "alice@org1", nouveau rôle "contributor"
// / @expect stubCA.UpdateAttributes reçoit le nom et l'attribut Myr.role="contributor"
func TestSetRole_Success(t *testing.T) {
	ca := &stubCA{}
	svc := identity.NewService(t.TempDir(), ca)
	if err := svc.SetRole(context.Background(), "alice@org1", "contributor"); err != nil {
		t.Fatalf("SetRole : %v", err)
	}
	if ca.lastUpdateName != "alice@org1" {
		t.Errorf("nom transmis : got %q, want alice@org1", ca.lastUpdateName)
	}
	if ca.lastUpdateAttrs["Myr.role"] != "contributor" {
		t.Errorf("Myr.role : got %q, want contributor", ca.lastUpdateAttrs["Myr.role"])
	}
}

// / @brief  SetRole propage l'erreur retournée par la CA
// / @input  stubCA dont updateAttrsFn retourne "CA: identité introuvable"
// / @expect retourne exactement l'erreur de la CA
func TestSetRole_CAError(t *testing.T) {
	expected := errors.New("CA: identité introuvable")
	ca := &stubCA{updateAttrsFn: func(_ context.Context, _ string, _ map[string]string) error { return expected }}
	svc := identity.NewService(t.TempDir(), ca)
	err := svc.SetRole(context.Background(), "ghost@org1", "admin")
	if !errors.Is(err, expected) {
		t.Errorf("erreur attendue %v, obtenu %v", expected, err)
	}
}

// / @brief  SetRole refuse un nom d'identité vide
// / @input  stubCA valide, name=""
// / @expect retourne une erreur de validation sans appeler la CA
func TestSetRole_MissingName(t *testing.T) {
	ca := &stubCA{}
	svc := identity.NewService(t.TempDir(), ca)
	err := svc.SetRole(context.Background(), "", "admin")
	if err == nil {
		t.Fatal("erreur attendue pour identité vide")
	}
	if ca.lastUpdateName != "" {
		t.Error("la CA ne devrait pas être appelée")
	}
}

// ── RequestAccess ─────────────────────────────────────────────────────────────

var openPolicy = identity.RegistrationPolicy{AllowAutoRegister: true, DefaultRole: identity.RoleReader}

func bobRequest() identity.AccountRequest {
	return identity.AccountRequest{Pseudo: "bob", Email: "bob@example.com", OrgID: "Org1MSP", DisplayName: "Bob"}
}

// / @brief  RequestAccess refuse toute demande quand le réseau n'accepte pas l'auto-enregistrement
// / @input  politique fermée (AllowAutoRegister=false), demande valide
// / @expect ErrRegistrationClosed ; rien n'est conservé, la CA n'est pas appelée
func TestRequestAccess_RegistrationClosed_Rejected(t *testing.T) {
	ca := &stubCA{registerFn: func(context.Context, identity.RegisterRequest) (string, error) {
		t.Fatal("la CA ne doit pas être appelée")
		return "", nil
	}}
	store := &mockRequestStore{}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	_, _, err := svc.RequestAccess(context.Background(), bobRequest(), identity.RegistrationPolicy{})
	if !errors.Is(err, identity.ErrRegistrationClosed) {
		t.Fatalf("ErrRegistrationClosed attendu, obtenu %v", err)
	}
	if len(store.requests) != 0 {
		t.Errorf("aucune demande ne doit être conservée, obtenu %d", len(store.requests))
	}
}

// / @brief  RequestAccess sans approbation : l'identité est créée active avec le rôle par défaut
// / @input  politique ouverte sans approbation, rôle par défaut reader
// / @expect demande approved conservée ; la CA reçoit "bob@Org1", statut active, rôle reader ; secret renvoyé
func TestRequestAccess_NoApproval_Active(t *testing.T) {
	ca := &stubCA{}
	store := &mockRequestStore{}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	got, secret, err := svc.RequestAccess(context.Background(), bobRequest(), openPolicy)
	if err != nil {
		t.Fatalf("RequestAccess : %v", err)
	}
	if got.ID == "" || got.Status != identity.RequestApproved || secret != "secret123" {
		t.Errorf("attendu demande approved avec ID et secret, obtenu %+v / %q", got, secret)
	}
	reg := ca.lastRegisterReq
	if reg.Name != "bob@Org1" || reg.Status != identity.StatusActive || reg.Role != identity.RoleReader {
		t.Errorf("requête CA inattendue : %+v", reg)
	}
	if len(store.requests) != 1 || store.requests[0].Status != identity.RequestApproved {
		t.Errorf("store : 1 demande approved attendue, obtenu %+v", store.requests)
	}
}

// / @brief  RequestAccess avec approbation : l'identité est créée en attente
// / @input  politique ouverte exigeant une approbation, rôle par défaut contributor
// / @expect demande pending conservée ; la CA reçoit le statut pending et le rôle contributor
func TestRequestAccess_RequireApproval_Pending(t *testing.T) {
	ca := &stubCA{}
	store := &mockRequestStore{}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	policy := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true, DefaultRole: identity.RoleContributor}
	got, _, err := svc.RequestAccess(context.Background(), bobRequest(), policy)
	if err != nil {
		t.Fatalf("RequestAccess : %v", err)
	}
	if got.Status != identity.RequestPending {
		t.Errorf("statut de demande attendu pending, obtenu %q", got.Status)
	}
	if ca.lastRegisterReq.Status != identity.StatusPending || ca.lastRegisterReq.Role != identity.RoleContributor {
		t.Errorf("requête CA inattendue : %+v", ca.lastRegisterReq)
	}
	if len(store.requests) != 1 || store.requests[0].Status != identity.RequestPending {
		t.Errorf("store : 1 demande pending attendue, obtenu %+v", store.requests)
	}
}

// / @brief  RequestAccess avec approbation exige un stockage des demandes
// / @input  politique exigeant une approbation, service sans requestStore
// / @expect erreur ; la CA n'est pas appelée (l'identité resterait en attente sans décision possible)
func TestRequestAccess_RequireApproval_NoStore_Rejected(t *testing.T) {
	ca := &stubCA{registerFn: func(context.Context, identity.RegisterRequest) (string, error) {
		t.Fatal("la CA ne doit pas être appelée")
		return "", nil
	}}
	svc := identity.NewService(t.TempDir(), ca)
	policy := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true}
	if _, _, err := svc.RequestAccess(context.Background(), bobRequest(), policy); err == nil {
		t.Fatal("erreur attendue sans stockage des demandes")
	}
}

// / @brief  RequestAccess refusé si pseudo, email ou org_id est vide, ou si le pseudo est invalide
// / @input  trois demandes incomplètes et un pseudo contenant "/"
// / @expect une erreur pour chaque cas ; ErrInvalidIdentityName pour le pseudo invalide
func TestRequestAccess_InvalidRequest_Rejected(t *testing.T) {
	svc := identity.NewService(t.TempDir(), &stubCA{}).WithRequestStore(&mockRequestStore{})
	cases := []identity.AccountRequest{
		{Pseudo: "", Email: "x@x.com", OrgID: "Org1MSP"},
		{Pseudo: "alice", Email: "", OrgID: "Org1MSP"},
		{Pseudo: "alice", Email: "x@x.com", OrgID: ""},
	}
	for _, c := range cases {
		if _, _, err := svc.RequestAccess(context.Background(), c, openPolicy); err == nil {
			t.Errorf("attendu erreur pour %+v", c)
		}
	}
	bad := identity.AccountRequest{Pseudo: "../alice", Email: "x@x.com", OrgID: "Org1MSP"}
	if _, _, err := svc.RequestAccess(context.Background(), bad, openPolicy); !errors.Is(err, identity.ErrInvalidIdentityName) {
		t.Errorf("ErrInvalidIdentityName attendu, obtenu %v", err)
	}
}

// / @brief  RequestAccess conserve la demande au statut failed si la CA refuse l'enregistrement
// / @input  stubCA dont registerFn échoue
// / @expect erreur enveloppant celle de la CA ; 1 demande au statut failed dans le store
func TestRequestAccess_CAError_Failed(t *testing.T) {
	expected := errors.New("CA: identite deja enregistree")
	ca := &stubCA{registerFn: func(context.Context, identity.RegisterRequest) (string, error) { return "", expected }}
	store := &mockRequestStore{}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	if _, _, err := svc.RequestAccess(context.Background(), bobRequest(), openPolicy); !errors.Is(err, expected) {
		t.Fatalf("erreur attendue %v, obtenu %v", expected, err)
	}
	if len(store.requests) != 1 || store.requests[0].Status != identity.RequestFailed {
		t.Errorf("store : 1 demande failed attendue, obtenu %+v", store.requests)
	}
}

// ── ApproveRequest / RejectRequest ────────────────────────────────────────────

func pendingStore() *mockRequestStore {
	return &mockRequestStore{requests: []*identity.AccountRequest{
		{ID: "req-1", Pseudo: "bob", OrgID: "Org1MSP", Status: identity.RequestPending},
	}}
}

// / @brief  ApproveRequest active l'identité de la demande et enregistre la décision
// / @input  demande pending "req-1" de bob@Org1
// / @expect la CA reçoit Myr.status=active pour "bob@Org1" ; demande approved
func TestApproveRequest_NominalCase(t *testing.T) {
	ca := &stubCA{}
	store := pendingStore()
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	got, err := svc.ApproveRequest(context.Background(), "req-1")
	if err != nil {
		t.Fatalf("ApproveRequest : %v", err)
	}
	if ca.lastUpdateName != "bob@Org1" || ca.lastUpdateAttrs[identity.AttrStatus] != identity.StatusActive {
		t.Errorf("mise à jour CA inattendue : %q %v", ca.lastUpdateName, ca.lastUpdateAttrs)
	}
	if got.Status != identity.RequestApproved || store.requests[0].Status != identity.RequestApproved {
		t.Errorf("demande approved attendue, obtenu %q (store %q)", got.Status, store.requests[0].Status)
	}
}

// / @brief  RejectRequest suspend l'identité sans rien supprimer
// / @input  demande pending "req-1" de bob@Org1
// / @expect la CA reçoit Myr.status=suspended ; demande rejected, toujours présente dans le store
func TestRejectRequest_NominalCase(t *testing.T) {
	ca := &stubCA{}
	store := pendingStore()
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	got, err := svc.RejectRequest(context.Background(), "req-1")
	if err != nil {
		t.Fatalf("RejectRequest : %v", err)
	}
	if ca.lastUpdateAttrs[identity.AttrStatus] != identity.StatusSuspended {
		t.Errorf("statut CA attendu suspended, obtenu %v", ca.lastUpdateAttrs)
	}
	if got.Status != identity.RequestRejected || len(store.requests) != 1 {
		t.Errorf("demande rejected conservée attendue, obtenu %q (%d demandes)", got.Status, len(store.requests))
	}
}

// / @brief  Une décision ne porte que sur une demande existante et encore en attente
// / @input  identifiant inconnu ; demande déjà approuvée
// / @expect ErrRequestNotFound puis ErrRequestNotPending ; la CA n'est jamais modifiée
func TestDecideRequest_NotFoundOrNotPending_Rejected(t *testing.T) {
	ca := &stubCA{updateAttrsFn: func(context.Context, string, map[string]string) error {
		t.Fatal("la CA ne doit pas être modifiée")
		return nil
	}}
	store := &mockRequestStore{requests: []*identity.AccountRequest{
		{ID: "req-1", Pseudo: "bob", OrgID: "Org1MSP", Status: identity.RequestApproved},
	}}
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	if _, err := svc.ApproveRequest(context.Background(), "req-x"); !errors.Is(err, identity.ErrRequestNotFound) {
		t.Errorf("ErrRequestNotFound attendu, obtenu %v", err)
	}
	if _, err := svc.RejectRequest(context.Background(), "req-1"); !errors.Is(err, identity.ErrRequestNotPending) {
		t.Errorf("ErrRequestNotPending attendu, obtenu %v", err)
	}
}

// / @brief  Une décision dont la mise à jour CA échoue laisse la demande en attente
// / @input  demande pending, UpdateAttributes en erreur
// / @expect erreur propagée ; la demande reste pending
func TestApproveRequest_CAError_StaysPending(t *testing.T) {
	expected := errors.New("CA indisponible")
	ca := &stubCA{updateAttrsFn: func(context.Context, string, map[string]string) error { return expected }}
	store := pendingStore()
	svc := identity.NewService(t.TempDir(), ca).WithRequestStore(store)
	if _, err := svc.ApproveRequest(context.Background(), "req-1"); !errors.Is(err, expected) {
		t.Fatalf("erreur attendue %v, obtenu %v", expected, err)
	}
	if store.requests[0].Status != identity.RequestPending {
		t.Errorf("la demande doit rester pending, obtenu %q", store.requests[0].Status)
	}
}

// ── CanSignIn ─────────────────────────────────────────────────────────────────

// / @brief  Seule une identité active dotée d'un rôle peut ouvrir une session
// / @input  combinaisons de statut et de rôle lues dans le certificat
// / @expect nil pour active+rôle ; erreur dédiée sinon (pending, suspended, statut inconnu, sans rôle)
func TestWalletEntry_CanSignIn(t *testing.T) {
	cases := []struct {
		status, role string
		want         error
	}{
		{identity.StatusActive, identity.RoleReader, nil},
		{identity.StatusPending, identity.RoleReader, identity.ErrIdentityPending},
		{identity.StatusSuspended, identity.RoleReader, identity.ErrIdentitySuspended},
		{"", identity.RoleReader, identity.ErrIdentityInactive},
		{"revoked", identity.RoleReader, identity.ErrIdentityInactive},
		{identity.StatusActive, "", identity.ErrIdentityWithoutRole},
	}
	for _, c := range cases {
		err := identity.WalletEntry{Status: c.status, Role: c.role}.CanSignIn()
		if c.want == nil && err != nil || c.want != nil && !errors.Is(err, c.want) {
			t.Errorf("statut %q rôle %q : attendu %v, obtenu %v", c.status, c.role, c.want, err)
		}
	}
}

// / @brief  PolicyOf reprend les règles d'inscription fixées sur le réseau
// / @input  réseau ouvert exigeant une approbation, rôle par défaut contributor ; réseau absent
// / @expect politique identique aux champs du réseau ; politique fermée sans réseau
func TestPolicyOf(t *testing.T) {
	p := &network.NetworkProfile{AllowAutoRegister: true, RequireApproval: true, AutoRegisterRole: identity.RoleContributor}
	want := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true, DefaultRole: identity.RoleContributor}
	if got := identity.PolicyOf(p); got != want {
		t.Errorf("PolicyOf = %+v, attendu %+v", got, want)
	}
	if got := identity.PolicyOf(nil); got.AllowAutoRegister {
		t.Errorf("politique fermée attendue sans réseau, obtenu %+v", got)
	}
}

// ── ListRequests ──────────────────────────────────────────────────────────────

// / @brief  ListRequests retourne une liste vide si aucune demande n'est stockée
// / @input  mockStore vide
// / @expect retourne une liste de longueur 0 sans erreur
func TestListRequests_Empty(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil).WithRequestStore(&mockRequestStore{})
	reqs, err := svc.ListRequests()
	if err != nil {
		t.Fatalf("ListRequests : %v", err)
	}
	if len(reqs) != 0 {
		t.Errorf("attendu 0 demande, got %d", len(reqs))
	}
}

// / @brief  ListRequests retourne toutes les demandes enregistrées
// / @input  mockStore avec 2 demandes pré-chargées
// / @expect retourne une liste de 2 éléments sans erreur
func TestListRequests_WithRequests(t *testing.T) {
	store := &mockRequestStore{
		requests: []*identity.AccountRequest{
			{ID: "req-1", Pseudo: "alice", Status: identity.RequestPending},
			{ID: "req-2", Pseudo: "bob", Status: identity.RequestPending},
		},
	}
	svc := identity.NewService(t.TempDir(), nil).WithRequestStore(store)
	reqs, err := svc.ListRequests()
	if err != nil {
		t.Fatalf("ListRequests : %v", err)
	}
	if len(reqs) != 2 {
		t.Errorf("attendu 2 demandes, got %d", len(reqs))
	}
}

// / @brief  ListRequests retourne nil sans erreur si aucun store n'est configuré
// / @input  service sans requestStore
// / @expect retourne (nil, nil)
func TestListRequests_NoStore(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	reqs, err := svc.ListRequests()
	if err != nil {
		t.Fatalf("ListRequests sans store : %v", err)
	}
	if reqs != nil {
		t.Error("attendu nil sans store")
	}
}

// ── LoadGuestWallet ───────────────────────────────────────────────────────────

// / @brief  LoadGuestWallet crée le wallet invité et les fichiers MSP sous walletDir
// / @input  fichiers cert.pem et key.pem valides, orgID "Org1MSP"
// / @expect handle "guest@Org1", statut active, fichier cert.pem présent dans MSP
func TestLoadGuestWallet_Success(t *testing.T) {
	certFile := filepath.Join(t.TempDir(), "cert.pem")
	keyFile := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(certFile, []byte(fakeCertPEM), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, []byte(fakeKeyPEM), 0600); err != nil {
		t.Fatal(err)
	}

	walletDir := t.TempDir()
	svc := identity.NewService(walletDir, nil)
	entry, err := svc.LoadGuestWallet(certFile, keyFile, "", "Org1MSP")
	if err != nil {
		t.Fatalf("LoadGuestWallet : %v", err)
	}
	if entry.Handle != "guest@Org1" {
		t.Errorf("Handle: got %q, want guest@Org1", entry.Handle)
	}
	if entry.Status != identity.StatusActive {
		t.Errorf("Status: got %q, want active", entry.Status)
	}
	certPath := filepath.Join(walletDir, "guest@Org1", "msp", "signcerts", "cert.pem")
	if _, statErr := os.Stat(certPath); os.IsNotExist(statErr) {
		t.Error("cert.pem absent dans le wallet invité")
	}
}

// / @brief  LoadGuestWallet échoue si le fichier certificat est introuvable
// / @input  chemins inexistants pour cert et key
// / @expect retourne une erreur de lecture de fichier
func TestLoadGuestWallet_MissingCert(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	_, err := svc.LoadGuestWallet("/no/such/cert.pem", "/no/such/key.pem", "", "Org1MSP")
	if err == nil {
		t.Fatal("erreur attendue pour certificat introuvable")
	}
}
