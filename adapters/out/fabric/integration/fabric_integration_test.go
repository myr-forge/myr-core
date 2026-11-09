// adapters/out/fabric/integration/fabric_integration_test.go
//
// Tests d'intégration contre un vrai réseau Hyperledger Fabric portant le
// chaincode myrcc : invariants du ledger (ADR-12, Chaincode.md §4), table des
// rôles (§4.3) et services domaine branchés sur le ledger. Ils écrivent sur le
// ledger : à exécuter uniquement sur un réseau de test, jamais en production.
//
// Lancer avec :
//   go test -tags integration -v -timeout 300s ./adapters/out/fabric/integration/
//
// Prérequis dans testdata/ (MYR_FABRIC_TESTDATA_DIR) :
//   testnet.env       — variables de connexion (FABRIC_*, dont FABRIC_CHANNELS)
//   tls-cert.pem      — certificat racine TLS du peer
//   client-cert.pem   — identité d'organisation partagée (sans attribut Myr)
//   client-key.pem
//   identities/<nom>/ — cert.pem + key.pem d'identités enrôlées auprès de la CA,
//                       avec leurs attributs Myr.role / Myr.status :
//                       alice, bob (contributor, active), carol (reader, active),
//                       dave (contributor, pending), root (admin, active)

//go:build integration

package fabric_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"

	"myr-core/adapters/out/fabric"
	"myr-core/adapters/out/localstorage"
	"myr-core/domain/model"
	"myr-core/domain/role"
)

// ── helpers ──────────────────────────────────────────────────────────────────

// testdataDir retourne le chemin absolu vers testdata/ (sibling de integration/,
// dans adapters/out/fabric/testdata/). Priorité à MYR_FABRIC_TESTDATA_DIR :
// un binaire de test compilé (`go test -c`) et exécuté à distance (scripts/test_remote.ps1)
// n'a plus le chemin source d'origine — testdata/ est alors copié à côté du binaire
// et son chemin passé explicitement par la variable d'environnement.
func testdataDir(t *testing.T) string {
	t.Helper()
	if dir := os.Getenv("MYR_FABRIC_TESTDATA_DIR"); dir != "" {
		return dir
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("impossible de déterminer le répertoire courant : %v", err)
	}
	return filepath.Join(wd, "..", "testdata")
}

// loadEnvFile lit un fichier .env simple (KEY=VALUE, commentaires #).
func loadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m, scanner.Err()
}

// testEnv lit testdata/testnet.env.
func testEnv(t *testing.T) map[string]string {
	t.Helper()
	env, err := loadEnvFile(filepath.Join(testdataDir(t), "testnet.env"))
	if err != nil {
		t.Fatalf("impossible de lire testnet.env : %v", err)
	}
	return env
}

// integrationConfig charge la fabric.Config depuis testdata/testnet.env + chemins PEM relatifs.
func integrationConfig(t *testing.T) fabric.Config {
	t.Helper()
	td := testdataDir(t)
	env := testEnv(t)
	return fabric.Config{
		PeerEndpoint:  env["FABRIC_PEER_ENDPOINT"],
		GatewayPeer:   env["FABRIC_GATEWAY_PEER"],
		MSPID:         env["FABRIC_MSP_ID"],
		FabricChannel: env["FABRIC_CHANNEL_NAME"],
		ChaincodeName: env["FABRIC_CHAINCODE_NAME"],
		TLSCertPath:   filepath.Join(td, "tls-cert.pem"),
		CertPath:      filepath.Join(td, "client-cert.pem"),
		KeyPath:       filepath.Join(td, "client-key.pem"),
	}
}

// channels retourne les canaux du réseau de test (FABRIC_CHANNELS), canal par
// défaut en tête.
func channels(t *testing.T) []string {
	t.Helper()
	env := testEnv(t)
	out := []string{env["FABRIC_CHANNEL_NAME"]}
	for _, ch := range strings.Split(env["FABRIC_CHANNELS"], ",") {
		if ch = strings.TrimSpace(ch); ch != "" && !slices.Contains(out, ch) {
			out = append(out, ch)
		}
	}
	return out
}

// newGateway ouvre une connexion réelle signée par l'identité d'organisation
// partagée ; le test est skipé si la config est incomplète.
func newGateway(t *testing.T) *fabric.GatewayClient {
	t.Helper()
	cfg := integrationConfig(t)
	if cfg.PeerEndpoint == "" {
		t.Skip("FABRIC_PEER_ENDPOINT non défini dans testnet.env")
	}
	gw, err := fabric.NewGatewayClient(cfg)
	if err != nil {
		t.Fatalf("connexion Fabric échouée : %v", err)
	}
	t.Cleanup(func() { gw.Close() })
	return gw
}

// userGateway ouvre une passerelle signée par l'identité de test name
// (testdata/identities/<name>), sur la connexion d'une passerelle d'organisation.
func userGateway(t *testing.T, name string) *fabric.GatewayClient {
	t.Helper()
	base := newGateway(t)
	dir := filepath.Join(testdataDir(t), "identities", name)
	gw, err := base.ForIdentity(integrationConfig(t).MSPID, filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem"))
	if err != nil {
		t.Fatalf("passerelle de %s : %v", name, err)
	}
	t.Cleanup(func() { gw.Close() })
	return gw
}

// handle retourne l'identifiant pseudo@org de l'identité de test name.
func handle(t *testing.T, name string) string {
	return name + "@" + strings.TrimSuffix(integrationConfig(t).MSPID, "MSP")
}

// uniqueID génère un identifiant unique pour éviter les collisions entre exécutions.
func uniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// submitted retourne un asset soumis, propriété de owner, sur le canal ch.
func submitted(id, owner, ch string) *model.Model3D {
	return &model.Model3D{
		ID: id, Name: "Integration " + id, ChannelID: ch, OwnerID: owner,
		Hash: "sha256:" + id, LicenseID: "cc-by-4", Status: model.ModuleSubmitted,
	}
}

// detail retourne le message renvoyé par le chaincode dans une erreur du
// gateway Fabric, à défaut l'erreur elle-même.
func detail(err error) string {
	if st, ok := status.FromError(errors.Unwrap(err)); ok {
		for _, d := range st.Details() {
			if ed, ok := d.(*gateway.ErrorDetail); ok && ed.GetMessage() != "" {
				return ed.GetMessage()
			}
		}
	}
	var msgs []string
	for e := err; e != nil; e = errors.Unwrap(e) {
		if st, ok := status.FromError(e); ok {
			for _, d := range st.Details() {
				if ed, ok := d.(*gateway.ErrorDetail); ok {
					msgs = append(msgs, ed.GetMessage())
				}
			}
		}
	}
	if len(msgs) > 0 {
		return strings.Join(msgs, " ; ")
	}
	return err.Error()
}

// expectRefused vérifie qu'une transaction est refusée et que le message du
// chaincode contient want.
func expectRefused(t *testing.T, what string, err error, want string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s : refus attendu, transaction acceptée", what)
		return
	}
	if msg := detail(err); !strings.Contains(msg, want) {
		t.Errorf("%s : refus attendu contenant %q, obtenu %q", what, want, msg)
	} else {
		t.Logf("%s : refusé — %s", what, msg)
	}
}

// ── Connexion et écriture signée par l'utilisateur (ADR-12) ──────────────────

// / @brief  La connexion gRPC + TLS + identité au peer s'établit
// / @input  testnet.env, identité d'organisation
// / @expect aucune erreur
func TestIntegration_GatewayConnect(t *testing.T) {
	newGateway(t)
}

// / @brief  Un asset soumis par son propriétaire est enregistré puis relu (ADR-12)
// / @input  alice signe un asset dont OwnerID = alice@Test, status submitted
// / @expect StoreModelRecord accepté ; GetModelRecord relit le même asset
func TestIntegration_StoreAndGet_UserSigned(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	ch := integrationConfig(t).FabricChannel
	a := submitted(uniqueID("store"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(a); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	got, err := bc.GetModelRecord(a.ID, ch)
	if err != nil {
		t.Fatalf("GetModelRecord : %v", err)
	}
	if got.ID != a.ID || got.OwnerID != a.OwnerID || got.Hash != a.Hash {
		t.Errorf("asset relu inattendu : %+v", got)
	}
}

// / @brief  Une passerelle d'organisation n'écrit jamais : le client refuse avant tout envoi (ADR-12)
// / @input  NewFabricBlockchain (identité partagée), StoreModelRecord
// / @expect model.ErrSignerRequired
func TestIntegration_SharedIdentity_ClientRefuses(t *testing.T) {
	bc := fabric.NewFabricBlockchain(newGateway(t))
	err := bc.StoreModelRecord(submitted(uniqueID("shared"), "User1@Test", integrationConfig(t).FabricChannel))
	if !errors.Is(err, model.ErrSignerRequired) {
		t.Errorf("ErrSignerRequired attendu, obtenu %v", err)
	}
}

// / @brief  Le chaincode refuse lui-même l'identité d'organisation partagée, même soumise directement (ADR-12)
// / @input  passerelle d'organisation forcée en écriture (NewUserFabricBlockchain), OwnerID = son CN
// / @expect transaction refusée par le chaincode
func TestIntegration_SharedIdentity_ChaincodeRefuses(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(newGateway(t))
	err := bc.StoreModelRecord(submitted(uniqueID("shared"), "User1@test.myrtest.local", integrationConfig(t).FabricChannel))
	expectRefused(t, "identité partagée", err, "signataire")
}

// / @brief  Le chaincode refuse un OwnerID différent du signataire (ADR-12)
// / @input  bob signe un asset déclaré propriété d'alice
// / @expect transaction refusée
func TestIntegration_OwnerNotSigner_Rejected(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "bob"))
	err := bc.StoreModelRecord(submitted(uniqueID("owner"), handle(t, "alice"), integrationConfig(t).FabricChannel))
	expectRefused(t, "OwnerID ≠ signataire", err, "différent du signataire")
}

// / @brief  Le chaincode refuse un signataire en attente ou sans permission d'écriture (ADR-12, RM48)
// / @input  dave (contributor, pending) puis carol (reader, active), chacun propriétaire de son asset
// / @expect deux refus : statut pending, puis permission write absente
func TestIntegration_PendingOrReader_Rejected(t *testing.T) {
	ch := integrationConfig(t).FabricChannel
	err := fabric.NewUserFabricBlockchain(userGateway(t, "dave")).StoreModelRecord(submitted(uniqueID("pending"), handle(t, "dave"), ch))
	expectRefused(t, "signataire pending", err, "pending")
	err = fabric.NewUserFabricBlockchain(userGateway(t, "carol")).StoreModelRecord(submitted(uniqueID("reader"), handle(t, "carol"), ch))
	expectRefused(t, "signataire reader", err, "permission write")
}

// ── Invariants de StoreModel (Chaincode.md §4.1) ─────────────────────────────

// / @brief  Un asset enregistré n'est jamais réécrit (RM06, RM19)
// / @input  alice enregistre a, puis le soumet à nouveau sous le même ID
// / @expect seconde transaction refusée ; l'asset relu est inchangé
func TestIntegration_RM06_WriteOnce(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	ch := integrationConfig(t).FabricChannel
	a := submitted(uniqueID("once"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(a); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	again := *a
	again.Name = "réécrit"
	again.Hash = "sha256:autre-" + a.ID
	expectRefused(t, "réécriture", bc.StoreModelRecord(&again), a.ID)
	if got, _ := bc.GetModelRecord(a.ID, ch); got == nil || got.Name != a.Name {
		t.Errorf("asset modifié : %+v", got)
	}
}

// / @brief  Une empreinte déjà enregistrée sous un autre ID est refusée (RM01 préalable)
// / @input  alice enregistre a1 de hash h, puis a2 de même hash h
// / @expect seconde transaction refusée
func TestIntegration_RM01_DuplicateHash_Rejected(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	ch := integrationConfig(t).FabricChannel
	a1 := submitted(uniqueID("hash1"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(a1); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	a2 := submitted(uniqueID("hash2"), handle(t, "alice"), ch)
	a2.Hash = a1.Hash
	expectRefused(t, "empreinte en double", bc.StoreModelRecord(a2), a1.ID)
}

// / @brief  Un asset dérivé n'est accepté que si son parent est enregistré (RM05)
// / @input  dérivé d'un parent absent ; puis dérivé d'un parent enregistré
// / @expect refus puis acceptation
func TestIntegration_RM05_ParentMustExist(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	ch := integrationConfig(t).FabricChannel
	orphan := submitted(uniqueID("orphan"), handle(t, "alice"), ch)
	orphan.ParentID = "absent-" + orphan.ID
	expectRefused(t, "parent absent", bc.StoreModelRecord(orphan), "parent")

	parent := submitted(uniqueID("parent"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(parent); err != nil {
		t.Fatalf("parent : %s", detail(err))
	}
	child := submitted(uniqueID("child"), handle(t, "alice"), ch)
	child.ParentID = parent.ID
	if err := bc.StoreModelRecord(child); err != nil {
		t.Errorf("dérivé d'un parent enregistré refusé : %s", detail(err))
	}
}

// / @brief  Le chaincode n'enregistre qu'un asset soumis (RM16, RM19)
// / @input  alice soumet un asset au statut draft
// / @expect transaction refusée
func TestIntegration_RM19_DraftRejected(t *testing.T) {
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	a := submitted(uniqueID("draft"), handle(t, "alice"), integrationConfig(t).FabricChannel)
	a.Status = model.ModuleDraft
	expectRefused(t, "brouillon", bc.StoreModelRecord(a), "asset soumis")
}

// ── UpdateMetadata (Chaincode.md §4.2) ───────────────────────────────────────

// / @brief  Le propriétaire modifie les métadonnées descriptives d'un asset soumis ; un autre utilisateur ne le peut pas (ADR-12, RM19)
// / @input  a d'alice ; UpdateModelMetadata par alice puis par bob
// / @expect alice : nom, description, tags, liens changés, empreinte et licence inchangées ; bob : refus
func TestIntegration_UpdateMetadata_OwnerOnly(t *testing.T) {
	ch := integrationConfig(t).FabricChannel
	alice := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	a := submitted(uniqueID("meta"), handle(t, "alice"), ch)
	if err := alice.StoreModelRecord(a); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	meta := model.AssetMetadata{Name: "vis M4", Description: "inox", Tags: []string{"vis"}, Links: []string{"https://shop.example/vis"}}
	if err := alice.UpdateModelMetadata(a.ID, ch, meta); err != nil {
		t.Fatalf("UpdateModelMetadata (propriétaire) : %s", detail(err))
	}
	got, err := alice.GetModelRecord(a.ID, ch)
	if err != nil {
		t.Fatalf("GetModelRecord : %v", err)
	}
	if got.Name != "vis M4" || got.Description != "inox" || len(got.Links) != 1 || got.Hash != a.Hash || got.LicenseID != a.LicenseID || got.OwnerID != a.OwnerID {
		t.Errorf("asset après mise à jour : %+v", got)
	}
	err = fabric.NewUserFabricBlockchain(userGateway(t, "bob")).UpdateModelMetadata(a.ID, ch, model.AssetMetadata{Name: "volé"})
	expectRefused(t, "métadonnées par un tiers", err, "différent du signataire")
}

// ── Lecture ──────────────────────────────────────────────────────────────────

// / @brief  ListModelRecords ne liste que des assets
// / @input  asset enregistré ; liste du canal
// / @expect la liste contient l'asset et aucune entrée sans ID (index ou table des rôles)
func TestIntegration_List(t *testing.T) {
	ch := integrationConfig(t).FabricChannel
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	a := submitted(uniqueID("list"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(a); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	list, err := fabric.NewFabricBlockchain(newGateway(t)).ListModelRecords(ch)
	if err != nil {
		t.Fatalf("ListModelRecords : %v", err)
	}
	found := false
	for _, m := range list {
		if m.ID == "" {
			t.Errorf("entrée sans ID dans la liste : %+v", m)
		}
		found = found || m.ID == a.ID
	}
	if !found {
		t.Errorf("asset %s absent de la liste (%d éléments)", a.ID, len(list))
	}
	t.Logf("%d asset(s) sur %s", len(list), ch)
}

// / @brief  GetModelRecord retourne une erreur pour un ID inexistant
// / @input  ID jamais enregistré
// / @expect erreur
func TestIntegration_GetModelRecord_NotFound(t *testing.T) {
	if _, err := fabric.NewFabricBlockchain(newGateway(t)).GetModelRecord("myr-inexistant-000000000", ""); err == nil {
		t.Error("erreur attendue pour un ID inexistant")
	}
}

// / @brief  IndexHashes est idempotent et ne signale aucun doublon sur un ledger déjà indexé (mise en service, ADR-12)
// / @input  deux exécutions successives d'IndexHashes sur le canal par défaut
// / @expect deux transactions acceptées, aucun doublon
func TestIntegration_IndexHashes_Idempotent(t *testing.T) {
	c := userGateway(t, "root").ContractFor(integrationConfig(t).FabricChannel)
	for i := 1; i <= 2; i++ {
		out, err := c.SubmitTransaction("IndexHashes")
		if err != nil {
			t.Fatalf("IndexHashes #%d : %s", i, detail(err))
		}
		if !strings.Contains(string(out), `"duplicates":[]`) {
			t.Errorf("IndexHashes #%d : doublons signalés : %s", i, out)
		}
		t.Logf("IndexHashes #%d : %s", i, out)
	}
}

// ── Table des rôles (Chaincode.md §4.3, RM48) ────────────────────────────────

// / @brief  Un gestionnaire des rôles publie un rôle sur chaque canal ; chaque canal le restitue (RM48)
// / @input  root (admin) Save d'un rôle unique via FabricRoleStore sur tous les canaux
// / @expect le rôle et son réglage auto_grant sont lus sur chaque canal ; admin porte toutes les permissions
func TestIntegration_RM48_RolePublishedOnEveryChannel(t *testing.T) {
	chs := channels(t)
	name := fmt.Sprintf("intg%d", time.Now().UnixNano()%1_000_000_000)
	store := fabric.NewUserFabricRoleStore(userGateway(t, "root"), chs)
	if err := store.Save(&role.Role{ID: name, Name: name, Permissions: []role.Permission{role.PermRead}, AutoGrant: true}); err != nil {
		t.Fatalf("Save : %s", detail(err))
	}
	t.Cleanup(func() { _ = store.Delete(name) })
	for _, ch := range chs {
		roles, err := fabric.NewFabricRoleStore(newGateway(t), []string{ch}).FindAll()
		if err != nil {
			t.Fatalf("FindAll %s : %v", ch, err)
		}
		var got *role.Role
		for _, r := range roles {
			if r.Name == name {
				got = r
			}
			if r.Name == role.NameAdmin && len(r.Permissions) != len(role.AllPermissions) {
				t.Errorf("%s : admin porte %v", ch, r.Permissions)
			}
		}
		if got == nil || !got.AutoGrant {
			t.Errorf("%s : rôle %s absent ou sans auto_grant (%d rôles)", ch, name, len(roles))
		}
	}
}

// / @brief  Seul un rôle portant role.admin publie la table, et admin reste immuable (RM48)
// / @input  alice (contributor) Save d'un rôle ; root Save et Delete du rôle admin
// / @expect trois refus par le chaincode
func TestIntegration_RM48_RoleTableGuards(t *testing.T) {
	chs := channels(t)[:1]
	err := fabric.NewUserFabricRoleStore(userGateway(t, "alice"), chs).Save(&role.Role{Name: "pirate", Permissions: []role.Permission{role.PermAdmin}})
	expectRefused(t, "publication par un contributor", err, "role.admin")
	root := fabric.NewUserFabricRoleStore(userGateway(t, "root"), chs)
	expectRefused(t, "modification d'admin", root.Save(&role.Role{Name: role.NameAdmin, Permissions: []role.Permission{role.PermRead}}), "immuable")
	expectRefused(t, "suppression d'admin", root.Delete(role.NameAdmin), "immuable")
}

// / @brief  L'écriture d'un asset suit la table publiée : write retiré puis rendu au rôle contributor (ADR-12, RM48)
// / @input  root retire write à contributor ; alice soumet ; root rétablit contributor ; alice soumet à nouveau
// / @expect premier StoreModelRecord refusé, second accepté
func TestIntegration_RM48_WritePermissionFollowsTable(t *testing.T) {
	ch := integrationConfig(t).FabricChannel
	roles := fabric.NewUserFabricRoleStore(userGateway(t, "root"), channels(t))
	restore := &role.Role{ID: role.NameContributor, Name: role.NameContributor, Permissions: []role.Permission{role.PermRead, role.PermWrite}, BuiltIn: true}
	if err := roles.Save(&role.Role{ID: role.NameContributor, Name: role.NameContributor, Permissions: []role.Permission{role.PermRead}, BuiltIn: true}); err != nil {
		t.Fatalf("retrait de write : %s", detail(err))
	}
	t.Cleanup(func() { _ = roles.Save(restore) })

	alice := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	expectRefused(t, "contributor sans write", alice.StoreModelRecord(submitted(uniqueID("nowrite"), handle(t, "alice"), ch)), "permission write")

	if err := roles.Save(restore); err != nil {
		t.Fatalf("rétablissement de write : %s", detail(err))
	}
	if err := alice.StoreModelRecord(submitted(uniqueID("write"), handle(t, "alice"), ch)); err != nil {
		t.Errorf("contributor rétabli : %s", detail(err))
	}
}

// ── Services domaine sur le ledger ───────────────────────────────────────────

// / @brief  role.Service publie un rôle sur le ledger, qu'un autre nœud (lecture seule) applique (RM47, RM48)
// / @input  service signé par root : Create d'un rôle automatique ; service en lecture seule, nouvelle instance
// / @expect le service en lecture seule connaît le rôle, sa permission et son attribution automatique ; Delete le retire
func TestIntegration_RoleService_PublishedAndApplied(t *testing.T) {
	chs := channels(t)
	name := fmt.Sprintf("svc%d", time.Now().UnixNano()%1_000_000_000)
	admin := role.NewService(fabric.NewUserFabricRoleStore(userGateway(t, "root"), chs))
	if _, err := admin.Create(name, []role.Permission{role.PermRead}, true); err != nil {
		t.Fatalf("Create : %s", detail(err))
	}
	reader := role.NewService(fabric.NewFabricRoleStore(newGateway(t), chs))
	if !reader.HasPermission(name, role.PermRead) {
		t.Errorf("%s : permission read non appliquée par un autre service", name)
	}
	if granted, err := role.GrantedOnRequest(reader, name); err != nil || !granted {
		t.Errorf("GrantedOnRequest(%s) : %v, %v", name, granted, err)
	}
	if err := admin.Delete(name); err != nil {
		t.Fatalf("Delete : %s", detail(err))
	}
	if role.NewService(fabric.NewFabricRoleStore(newGateway(t), chs)).HasPermission(name, role.PermRead) {
		t.Errorf("%s : rôle supprimé encore appliqué", name)
	}
}

// / @brief  model.Service sur le ledger : métadonnées d'un asset soumis modifiables, licence et re-soumission refusées (RM19, ADR-12)
// / @input  asset soumis d'alice ; UpdateAsset (nom), UpdateAsset (licence), Submit
// / @expect nom mis à jour sur le ledger ; ErrAssetSubmitted pour la licence et la re-soumission
func TestIntegration_ModelService_SubmittedAsset(t *testing.T) {
	ch := integrationConfig(t).FabricChannel
	bc := fabric.NewUserFabricBlockchain(userGateway(t, "alice"))
	a := submitted(uniqueID("svc"), handle(t, "alice"), ch)
	if err := bc.StoreModelRecord(a); err != nil {
		t.Fatalf("StoreModelRecord : %s", detail(err))
	}
	store := localstorage.NewJSONBlockchain(filepath.Join(t.TempDir(), "assets.json"))
	svc := model.NewService(bc, nil).WithDraftStore(store).WithIfaceStore(store)

	if _, err := svc.UpdateAsset(model.UpdateRequest{ID: a.ID, Name: "renommé"}); err != nil {
		t.Fatalf("UpdateAsset (nom) : %s", detail(err))
	}
	if got, _ := bc.GetModelRecord(a.ID, ch); got == nil || got.Name != "renommé" {
		t.Errorf("nom non mis à jour sur le ledger : %+v", got)
	}
	if _, err := svc.UpdateAsset(model.UpdateRequest{ID: a.ID, LicenseID: "cc-by-sa-4"}); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("licence : ErrAssetSubmitted attendu, obtenu %v", err)
	}
	if _, err := svc.Submit(a.ID, ""); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("re-soumission : ErrAssetSubmitted attendu, obtenu %v", err)
	}
}
