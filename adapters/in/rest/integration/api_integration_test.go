// adapters/in/rest/integration/api_integration_test.go
//
// Tests d'intégration end-to-end de l'API REST.
// Aucun mock : domaine réel + localstorage réel + serveur HTTP réel (httptest.Server).
// Seules la blockchain et l'autorité de certification sont simulées en mémoire ;
// la CA émet de vrais certificats X.509, lus par le parseur de l'adapter Fabric,
// et chaque test passe par une vraie connexion (POST /api/identity/session).
//
// Lancer avec :
//   go test -tags integration -v -timeout 60s ./adapters/in/rest/integration/...

//go:build integration

package rest_integration_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"myr-core/adapters/in/rest"
	"myr-core/adapters/out/fabric"
	"myr-core/adapters/out/localstorage"
	"myr-core/domain/identity"
	"myr-core/domain/model"
	"myr-core/domain/role"
)

// ── In-memory blockchain (remplace Fabric en test d'intégration) ──────────────
// Implémente BlockchainPort + le port de suppression optionnel.

type memBC struct {
	mu      sync.Mutex
	records map[string]*model.Model3D
}

func newMemBC() *memBC { return &memBC{records: make(map[string]*model.Model3D)} }

func (b *memBC) StoreModelRecord(m *model.Model3D) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	cp := *m
	b.records[m.ID] = &cp
	return nil
}

func (b *memBC) UpdateModelMetadata(id, _ string, meta model.AssetMetadata) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	m, ok := b.records[id]
	if !ok {
		return fmt.Errorf("modèle introuvable : %s", id)
	}
	cp := *m
	cp.Name, cp.Description, cp.Tags, cp.Links = meta.Name, meta.Description, meta.Tags, meta.Links
	b.records[id] = &cp
	return nil
}

func (b *memBC) GetModelRecord(id, _ string) (*model.Model3D, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if m, ok := b.records[id]; ok {
		cp := *m
		return &cp, nil
	}
	return nil, fmt.Errorf("%w : modèle %s", model.ErrNotFound, id)
}

func (b *memBC) ListModelRecords(channelID string) ([]*model.Model3D, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var list []*model.Model3D
	for _, m := range b.records {
		if channelID == "" || m.ChannelID == channelID {
			cp := *m
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (b *memBC) RemoveModelRecord(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.records[id]; !ok {
		return fmt.Errorf("modèle introuvable : %s", id)
	}
	delete(b.records, id)
	return nil
}

// ── Autorité de certification en mémoire (remplace la Fabric CA) ──────────────
// Retient les attributs de chaque identité et les inscrit, comme la Fabric CA,
// dans l'extension X.509 des certificats qu'elle émet ; la lecture de ces
// attributs passe par le parseur réel de l'adapter Fabric.

type memCA struct {
	mu      sync.Mutex
	attrs   map[string]map[string]string
	secrets map[string]string
}

func newMemCA() *memCA {
	return &memCA{attrs: map[string]map[string]string{}, secrets: map[string]string{}}
}

func (c *memCA) Register(_ context.Context, req identity.RegisterRequest) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.attrs[req.Name]; ok {
		return "", fmt.Errorf("identité %s déjà enregistrée", req.Name)
	}
	status := req.Status
	if status == "" {
		status = identity.StatusPending
	}
	c.attrs[req.Name] = map[string]string{identity.AttrRole: req.Role, identity.AttrStatus: status}
	c.secrets[req.Name] = fmt.Sprintf("secret-%d", time.Now().UnixNano())
	return c.secrets[req.Name], nil
}

func (c *memCA) Enroll(_ context.Context, name, secret string) (string, string, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s, ok := c.secrets[name]; !ok || s != secret {
		return "", "", "", fmt.Errorf("secret d'enrôlement invalide pour %s", name)
	}
	return c.issue(name)
}

func (c *memCA) GetStatus(_ context.Context, name string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.attrs[name][identity.AttrStatus], nil
}

func (c *memCA) ReEnroll(_ context.Context, name, _, _ string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cert, _, _, err := c.issue(name)
	return cert, err
}

func (c *memCA) UpdateAttributes(_ context.Context, name string, attrs map[string]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.attrs[name] == nil {
		return fmt.Errorf("identité %s inconnue", name)
	}
	for k, v := range attrs {
		c.attrs[name][k] = v
	}
	return nil
}

func (c *memCA) CertAttributes(certPEM string) (map[string]string, error) {
	return fabric.CertAttributes(certPEM)
}

// issue émet un certificat portant les attributs courants de l'identité, dans
// l'extension utilisée par la Fabric CA (OID 1.2.3.4.5.6.7.8.1).
func (c *memCA) issue(name string) (string, string, string, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", err
	}
	ext, _ := json.Marshal(map[string]any{"attrs": c.attrs[name]})
	tmpl := &x509.Certificate{
		SerialNumber:    big.NewInt(time.Now().UnixNano()),
		Subject:         pkix.Name{CommonName: name},
		NotBefore:       time.Now().Add(-time.Hour),
		NotAfter:        time.Now().Add(time.Hour),
		ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 1}, Value: ext}},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return "", "", "", err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", "", err
	}
	cert := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	return cert, string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})), cert, nil
}

// ── Helpers ────────────────────────────────────────────────────────────────────

// tokenTransport ajoute le jeton de session à chaque requête du client de test.
type tokenTransport struct {
	base  http.RoundTripper
	token string
}

func (t tokenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("X-Myr-Token", t.token)
	return t.base.RoundTrip(r)
}

// newIntegrationServer câble un serveur HTTP complet avec des adapters réels,
// crée un compte contributor et ouvre sa session : srv.Client() envoie ensuite
// le jeton de session à chaque requête.
func newIntegrationServer(t *testing.T) *httptest.Server {
	t.Helper()
	dataDir := t.TempDir()

	bc := newMemBC()
	db := localstorage.NewJSONBlockchain(filepath.Join(dataDir, "db.json"))

	// Même câblage que cmd/api/main.go (hors Fabric, remplacé par memBC).
	svc := model.NewService(bc).
		WithConnStore(db).
		WithThumbStore(db).
		WithIfaceStore(db).
		WithDraftStore(db).
		WithRemovedStore(db)

	roleSvc := role.NewService(nil) // table initiale, comme sans réseau joignable
	idSvc := identity.NewService(filepath.Join(dataDir, "wallets"), newMemCA()).
		WithRequestStore(localstorage.NewJSONRequestStore(filepath.Join(dataDir, "account-requests.json"))).
		WithRoles(roleSvc)

	h := rest.NewHandler(svc, db, db).
		WithIdentityService(idSvc).
		WithRoleService(roleSvc).
		WithConnStore(db).
		WithDraftStore(db).
		WithRemovedStore(db)
	mux, err := rest.NewServer(h, "localhost:0").Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	srv := httptest.NewServer(mux)

	_, secret, err := idSvc.CreateAccount(context.Background(), "tester", "Org1MSP", identity.RoleContributor, "")
	if err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	resp, err := srv.Client().Post(srv.URL+"/api/identity/session", "application/json",
		bytes.NewBufferString(`{"name":"tester","secret":"`+secret+`","org_id":"Org1MSP"}`))
	if err != nil {
		t.Fatalf("POST /api/identity/session: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		t.Fatalf("connexion: got %d, want 200", resp.StatusCode)
	}
	var sess struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
	decodeBody(t, resp, &sess)
	if sess.Role != identity.RoleContributor {
		t.Fatalf("rôle de session: got %q, want contributor (lu dans le certificat)", sess.Role)
	}
	client := srv.Client()
	client.Transport = tokenTransport{base: client.Transport, token: sess.Token}
	return srv
}

// apiDo envoie une requête HTTP et retourne la réponse.
func apiDo(t *testing.T, client *http.Client, method, url, body string) *http.Response {
	t.Helper()
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do %s %s: %v", method, url, err)
	}
	return resp
}

// decodeBody décode le JSON de la réponse dans v.
func decodeBody(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

// multipartBody construit un corps multipart minimal pour les tests de composants :
// nom, catégorie et canal, que le serveur ne choisit jamais à la place du client.
func multipartBody(name string) (string, string) {
	boundary := "testboundary"
	field := func(k, v string) string {
		return fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"%s\"\r\n\r\n%s\r\n", boundary, k, v)
	}
	body := field("name", name) + field("category", "base") + field("channel_id", "ch-integration") + "--" + boundary + "--\r\n"
	ct := "multipart/form-data; boundary=" + boundary
	return body, ct
}

// ── Tests de base ─────────────────────────────────────────────────────────────

/// @brief  GET /api/ping retourne 200 avec ok=true (serveur réel)
func TestAPIIntegration_Ping(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api/ping")
	if err != nil {
		t.Fatalf("GET /api/ping: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want 200", resp.StatusCode)
	}
	var result struct{ Ok bool `json:"ok"` }
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Ok {
		t.Error("ok doit être true")
	}
}

/// @brief  GET /api/status retourne le mode et le nombre d'assets
func TestAPIIntegration_Status(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api/status")
	if err != nil {
		t.Fatalf("GET /api/status: %v", err)
	}
	var status struct {
		Mode   string `json:"mode"`
		Assets int    `json:"assets"`
	}
	decodeBody(t, resp, &status)
	if status.Mode == "" {
		t.Error("mode ne doit pas être vide")
	}
}

/// @brief  GET /api/licenses retourne le catalogue complet (MIT présent)
func TestAPIIntegration_LicenseCatalog(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api/licenses")
	if err != nil {
		t.Fatalf("GET /api/licenses: %v", err)
	}
	var licenses []struct{ ID string `json:"id"` }
	decodeBody(t, resp, &licenses)
	if len(licenses) == 0 {
		t.Fatal("catalogue de licences vide")
	}
	found := false
	for _, l := range licenses {
		if l.ID == "mit" {
			found = true
		}
	}
	if !found {
		t.Error("licence MIT absente du catalogue")
	}
}

// ── Cycle de vie d'un composant ───────────────────────────────────────────────

/// @brief  Cycle complet : POST → GET → PATCH → DELETE /api/components avec vrais adapters
/// @input  Domaine réel + LocalStorage réel + JSONBlockchain réel
/// @expect Chaque étape retourne le code HTTP attendu ; le nom est bien mis à jour ; 404 après suppression
func TestAPIIntegration_AssetLifecycle(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	// Créer
	body, ct := multipartBody("Vis M3")
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/components", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", ct)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /api/components: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		t.Fatalf("création: got %d, want 201", resp.StatusCode)
	}
	var created struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	decodeBody(t, resp, &created)
	if created.ID == "" {
		t.Fatal("ID vide dans la réponse de création")
	}
	if created.Name != "Vis M3" {
		t.Errorf("name: got %q, want Vis M3", created.Name)
	}

	// Lire
	resp2 := apiDo(t, client, http.MethodGet, srv.URL+"/api/components/"+created.ID, "")
	if resp2.StatusCode != http.StatusOK {
		resp2.Body.Close()
		t.Fatalf("lecture: got %d, want 200", resp2.StatusCode)
	}
	var got struct{ Name string `json:"name"` }
	decodeBody(t, resp2, &got)
	if got.Name != "Vis M3" {
		t.Errorf("name lu: got %q, want Vis M3", got.Name)
	}

	// Vérifier dans la liste
	resp3 := apiDo(t, client, http.MethodGet, srv.URL+"/api/components", "")
	var list struct{ Total int `json:"total"` }
	decodeBody(t, resp3, &list)
	if list.Total < 1 {
		t.Errorf("total: got %d, want ≥1", list.Total)
	}

	// Mettre à jour
	resp4 := apiDo(t, client, http.MethodPatch, srv.URL+"/api/components/"+created.ID, `{"name":"Vis M4"}`)
	if resp4.StatusCode != http.StatusOK {
		resp4.Body.Close()
		t.Errorf("PATCH: got %d, want 200", resp4.StatusCode)
	} else {
		var updated struct{ Name string `json:"name"` }
		decodeBody(t, resp4, &updated)
		if updated.Name != "Vis M4" {
			t.Errorf("name après PATCH: got %q, want Vis M4", updated.Name)
		}
	}

	// Supprimer
	resp5 := apiDo(t, client, http.MethodDelete, srv.URL+"/api/components/"+created.ID, "")
	resp5.Body.Close()
	if resp5.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE: got %d, want 204", resp5.StatusCode)
	}

	// Vérifier la suppression
	resp6 := apiDo(t, client, http.MethodGet, srv.URL+"/api/components/"+created.ID, "")
	resp6.Body.Close()
	if resp6.StatusCode != http.StatusNotFound {
		t.Errorf("GET après DELETE: got %d, want 404", resp6.StatusCode)
	}
}

// ── Connexions ────────────────────────────────────────────────────────────────

/// @brief  Cycle complet : POST → DELETE /api/connections avec JSONBlockchain réel
/// @input  Connexion from/to persistée en JSON
/// @expect connexion créée avec ID ; supprimée en 204 ; persiste entre les requêtes
func TestAPIIntegration_ConnectionLifecycle(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	resp := apiDo(t, client, http.MethodPost, srv.URL+"/api/connections",
		`{"from":"asset-A","to":"asset-B","label":"assemblage"}`)
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		t.Fatalf("POST connexion: got %d, want 201", resp.StatusCode)
	}
	var conn struct {
		ID    string `json:"id"`
		From  string `json:"from"`
		To    string `json:"to"`
		Label string `json:"label"`
	}
	decodeBody(t, resp, &conn)
	if conn.ID == "" {
		t.Fatal("connexion ID vide")
	}
	if conn.From != "asset-A" || conn.To != "asset-B" {
		t.Errorf("from/to: got %q/%q", conn.From, conn.To)
	}

	resp2 := apiDo(t, client, http.MethodDelete, srv.URL+"/api/connections/"+conn.ID, "")
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE connexion: got %d, want 204", resp2.StatusCode)
	}
}

// ── Modules ────────────────────────────────────────────────────────────────────

// createComponent crée un composant brouillon (POST multipart /api/components).
func createComponent(t *testing.T, srv *httptest.Server, name string) (id, status string) {
	t.Helper()
	body, ct := multipartBody(name)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/components", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", ct)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("POST /api/components: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		t.Fatalf("création de %s: got %d, want 201", name, resp.StatusCode)
	}
	var c struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	decodeBody(t, resp, &c)
	return c.ID, c.Status
}

/// @brief  Cycle complet d'un module (ADR-11) : composants brouillons → instances → liaison interne → assemblage → soumission (RM17, RM51)
/// @input  Module et composant réels, connexion réelle, soumission blockchain simulée (memBC)
/// @expect statut "draft" à la création, "submitted" après soumission
func TestAPIIntegration_ModuleLifecycle(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	// Créer le module et un composant à y instancier
	modID, status := createComponent(t, srv, "Drone v1")
	if status != string(model.ModuleDraft) {
		t.Errorf("status initial: got %q, want draft", status)
	}
	partID, _ := createComponent(t, srv, "Moteur")
	propID, _ := createComponent(t, srv, "Hélice")

	// Instancier les deux composants dans le module
	for _, id := range []string{partID, propID} {
		resp1 := apiDo(t, client, http.MethodPost, srv.URL+"/api/components/"+modID+"/instances",
			`{"asset_id":"`+id+`"}`)
		resp1.Body.Close()
		if resp1.StatusCode != http.StatusOK && resp1.StatusCode != http.StatusCreated {
			t.Fatalf("POST instances: got %d, want 200/201", resp1.StatusCode)
		}
	}

	// Relier les deux instances : seule une liaison interne au module s'y rattache (RM51)
	resp2 := apiDo(t, client, http.MethodPost, srv.URL+"/api/connections",
		`{"from":"`+partID+`","to":"`+propID+`","label":"lien"}`)
	var conn struct{ ID string `json:"id"` }
	decodeBody(t, resp2, &conn)

	// Ajouter la connexion comme assemblage
	resp3 := apiDo(t, client, http.MethodPost, srv.URL+"/api/components/"+modID+"/assemblies",
		`{"connection_id":"`+conn.ID+`"}`)
	if resp3.StatusCode != http.StatusOK {
		resp3.Body.Close()
		t.Fatalf("POST assemblies: got %d, want 200", resp3.StatusCode)
	}
	resp3.Body.Close()

	// Soumettre
	resp4 := apiDo(t, client, http.MethodPost, srv.URL+"/api/components/"+modID+"/submit",
		`{"note":"première livraison"}`)
	if resp4.StatusCode != http.StatusOK {
		resp4.Body.Close()
		t.Fatalf("POST submit: got %d, want 200", resp4.StatusCode)
	}
	var submitted struct{ Status string `json:"status"` }
	decodeBody(t, resp4, &submitted)
	if submitted.Status != string(model.ModuleSubmitted) {
		t.Errorf("status après submit: got %q, want submitted", submitted.Status)
	}
}

// ── Interfaces physiques ──────────────────────────────────────────────────────

// physicalInterfaces liste les identifiants des interfaces d'un composant.
func physicalInterfaces(t *testing.T, srv *httptest.Server, assetID string) []string {
	t.Helper()
	resp := apiDo(t, srv.Client(), http.MethodGet, srv.URL+"/api/components/"+assetID+"/interfaces", "")
	var ifaces []struct {
		ID string `json:"id"`
	}
	decodeBody(t, resp, &ifaces)
	var ids []string
	for _, i := range ifaces {
		ids = append(ids, i.ID)
	}
	return ids
}

/// @brief  Cycle complet : ajouter interface → lister → supprimer (JSONBlockchain réel)
/// @input  Composant créé ; interface ELEC USB-C Out ajoutée et supprimée
/// @expect interface physique présente dans la liste après ajout, absente après suppression
func TestAPIIntegration_InterfaceLifecycle(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	// Créer un composant
	assetID, _ := createComponent(t, srv, "Capteur")
	asset := struct{ ID string }{assetID}

	// Ajouter une interface
	resp2 := apiDo(t, client, http.MethodPost, srv.URL+"/api/components/"+asset.ID+"/interfaces",
		`{"category":"ELEC","type":"USB-C","direction":"out"}`)
	if resp2.StatusCode != http.StatusCreated {
		resp2.Body.Close()
		t.Fatalf("POST interface: got %d, want 201", resp2.StatusCode)
	}
	var iface struct{ ID string `json:"id"` }
	decodeBody(t, resp2, &iface)
	if iface.ID == "" {
		t.Fatal("interface ID vide")
	}

	// Lister
	if ifaces := physicalInterfaces(t, srv, asset.ID); len(ifaces) != 1 || ifaces[0] != iface.ID {
		t.Errorf("interfaces physiques: got %v, want [%s]", ifaces, iface.ID)
	}

	// Supprimer
	resp4 := apiDo(t, client, http.MethodDelete, srv.URL+"/api/interfaces/"+iface.ID, "")
	resp4.Body.Close()
	if resp4.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE interface: got %d, want 204", resp4.StatusCode)
	}

	// Vérifier l'absence
	if ifaces := physicalInterfaces(t, srv, asset.ID); len(ifaces) != 0 {
		t.Errorf("interface toujours présente après suppression : %v", ifaces)
	}
}

// ── Persistance JSON ──────────────────────────────────────────────────────────

/// @brief  Les connexions persistent sur disque entre deux instances du service
/// @input  Connexion créée via la 1re instance ; relue via une 2e instance pointant sur le même fichier
/// @expect connexion retrouvée après rechargement du JSONBlockchain depuis disque
func TestAPIIntegration_ConnectionPersistsOnDisk(t *testing.T) {
	dataDir := t.TempDir()
	dbPath := filepath.Join(dataDir, "db.json")

	// Instance 1 : créer une connexion
	func() {
		db := localstorage.NewJSONBlockchain(dbPath)
		if err := db.SaveConnection(&model.Connection{
			ID: "c-persist", From: "asset-X", To: "asset-Y", Label: "test",
		}); err != nil {
			t.Fatalf("SaveConnection: %v", err)
		}
	}()

	// Instance 2 : recharger depuis disque et vérifier
	db2 := localstorage.NewJSONBlockchain(dbPath)
	conns, err := db2.ListConnections()
	if err != nil {
		t.Fatalf("ListConnections: %v", err)
	}
	if len(conns) != 1 {
		t.Fatalf("got %d connexions, want 1", len(conns))
	}
	if conns[0].ID != "c-persist" {
		t.Errorf("ID: got %q, want c-persist", conns[0].ID)
	}
}

// ── Requêtes concurrentes ─────────────────────────────────────────────────────

/// @brief  N assets créés en parallèle — aucune course dans le handler HTTP
/// @input  10 goroutines POST simultanées sur le serveur httptest
/// @expect 10 codes 201 ; total=10 dans la liste finale
func TestAPIIntegration_ConcurrentAssetCreation(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	const n = 10
	results := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(idx int) {
			body, ct := multipartBody(fmt.Sprintf("Asset%02d", idx))
			req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/components", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", ct)
			resp, err := client.Do(req)
			if err != nil {
				results <- 0
				return
			}
			resp.Body.Close()
			results <- resp.StatusCode
		}(i)
	}

	for i := 0; i < n; i++ {
		code := <-results
		if code != http.StatusCreated {
			t.Errorf("goroutine %d: got %d, want 201", i, code)
		}
	}

	resp := apiDo(t, client, http.MethodGet, srv.URL+"/api/components", "")
	var list struct{ Total int `json:"total"` }
	decodeBody(t, resp, &list)
	if list.Total != n {
		t.Errorf("total après %d créations parallèles: got %d, want %d", n, list.Total, n)
	}
}

// ── Vocabulaire de référence ──────────────────────────────────────────────────

/// @brief  Ajout d'une catégorie de référence puis lecture via GET /api/refs
/// @input  POST /api/refs/categories {"name":"PNEUM"} ; GET /api/refs
/// @expect catégorie "PNEUM" présente dans refs.categories
func TestAPIIntegration_RefsCategoryPersists(t *testing.T) {
	srv := newIntegrationServer(t)
	defer srv.Close()
	client := srv.Client()

	resp := apiDo(t, client, http.MethodPost, srv.URL+"/api/refs/categories", `{"name":"PNEUM"}`)
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		t.Fatalf("POST refs/categories: got %d, want 201", resp.StatusCode)
	}
	resp.Body.Close()

	resp2 := apiDo(t, client, http.MethodGet, srv.URL+"/api/refs", "")
	var refs struct {
		Categories []string `json:"categories"`
	}
	decodeBody(t, resp2, &refs)
	found := false
	for _, cat := range refs.Categories {
		if cat == "PNEUM" {
			found = true
		}
	}
	if !found {
		t.Errorf("catégorie PNEUM absente dans refs : %v", refs.Categories)
	}
}
