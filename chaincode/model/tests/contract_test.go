// chaincode/model/tests/contract_test.go
// Tests du chaincode myrcc : invariants de StoreModel (Chaincode.md §4.1,
// ADR-12) et migration IndexHashes. Aucune infrastructure Fabric : fakeStub
// reproduit la sémantique du ledger utile ici — écritures validées seulement
// si la transaction réussit, GetState aveugle aux écritures de la transaction
// en cours, clés composites absentes des parcours de plage.
package model_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/pkg/attrmgr"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"

	"myrcc/model"
)

const testChannel = "sandbox"

// ── fakeStub ──────────────────────────────────────────────────────────────────

// fakeStub n'implémente que les méthodes utilisées par le contrat ; l'interface
// embarquée (nil) fait paniquer tout appel imprévu, ce qui signale un test à
// compléter plutôt qu'un comportement silencieux.
type fakeStub struct {
	shim.ChaincodeStubInterface
	fn      string
	args    []string
	creator []byte
	state   map[string][]byte // état validé
	writes  map[string][]byte // écritures de la transaction en cours
}

func newStub(t *testing.T, creator []byte) *fakeStub {
	t.Helper()
	return &fakeStub{creator: creator, state: map[string][]byte{}}
}

func (f *fakeStub) GetFunctionAndParameters() (string, []string) { return f.fn, f.args }
func (f *fakeStub) GetTxID() string                              { return "tx-test" }
func (f *fakeStub) GetChannelID() string                         { return testChannel }
func (f *fakeStub) GetCreator() ([]byte, error)                  { return f.creator, nil }

// GetState lit l'état validé : comme Fabric, une transaction ne voit pas ses
// propres PutState.
func (f *fakeStub) GetState(key string) ([]byte, error) { return f.state[key], nil }

func (f *fakeStub) PutState(key string, value []byte) error {
	f.writes[key] = value
	return nil
}

func (f *fakeStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return shim.CreateCompositeKey(objectType, attributes)
}

// GetStateByRange("", "") parcourt les clés simples triées : comme Fabric, les
// clés composites (préfixe \x00) en sont exclues.
func (f *fakeStub) GetStateByRange(_, _ string) (shim.StateQueryIteratorInterface, error) {
	var keys []string
	for k := range f.state {
		if !strings.HasPrefix(k, "\x00") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	it := &fakeIterator{}
	for _, k := range keys {
		it.kvs = append(it.kvs, &queryresult.KV{Key: k, Value: f.state[k]})
	}
	return it, nil
}

type fakeIterator struct {
	kvs []*queryresult.KV
	pos int
}

func (it *fakeIterator) HasNext() bool { return it.pos < len(it.kvs) }
func (it *fakeIterator) Close() error  { return nil }
func (it *fakeIterator) Next() (*queryresult.KV, error) {
	kv := it.kvs[it.pos]
	it.pos++
	return kv, nil
}

// invoke exécute une transaction et ne valide ses écritures qu'en cas de succès.
func (f *fakeStub) invoke(fn string, args ...string) peer.Response {
	f.fn, f.args = fn, args
	f.writes = map[string][]byte{}
	resp := new(model.SmartContract).Invoke(f)
	if resp.Status == shim.OK {
		for k, v := range f.writes {
			f.state[k] = v
		}
	}
	return resp
}

// ── identités et assets ───────────────────────────────────────────────────────

// identity sérialise une identité Fabric (MSP + certificat X.509) dont le
// certificat porte les attributs donnés, au format de la Fabric CA.
func identity(t *testing.T, mspID, enrollmentID string, attrs map[string]string) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("clé : %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: enrollmentID},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	if attrs != nil {
		value, err := json.Marshal(attrmgr.Attributes{Attrs: attrs})
		if err != nil {
			t.Fatalf("attributs : %v", err)
		}
		tmpl.ExtraExtensions = []pkix.Extension{{Id: attrmgr.AttrOID, Value: value}}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("certificat : %v", err)
	}
	sid, err := proto.Marshal(&msp.SerializedIdentity{
		Mspid:   mspID,
		IdBytes: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
	})
	if err != nil {
		t.Fatalf("identité : %v", err)
	}
	return sid
}

func contributor(t *testing.T, enrollmentID string) []byte {
	return identity(t, "Org1MSP", enrollmentID, map[string]string{"Myr.role": "contributor", "Myr.status": "active"})
}

func asset(id, owner, hash string) model.Model3D {
	return model.Model3D{ID: id, Name: "pièce " + id, Hash: hash, ChannelID: testChannel, OwnerID: owner, Status: model.ModuleSubmitted}
}

func storeJSON(t *testing.T, m model.Model3D) string {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal : %v", err)
	}
	return string(b)
}

func hashKey(t *testing.T, hash string) string {
	t.Helper()
	k, err := shim.CreateCompositeKey("hash", []string{hash})
	if err != nil {
		t.Fatalf("clé composite : %v", err)
	}
	return k
}

func expectRejected(t *testing.T, resp peer.Response, want string) {
	t.Helper()
	if resp.Status == shim.OK {
		t.Fatalf("transaction acceptée, want refus contenant %q", want)
	}
	if !strings.Contains(resp.Message, want) {
		t.Errorf("message = %q, want contenant %q", resp.Message, want)
	}
}

// ── StoreModel ────────────────────────────────────────────────────────────────

// / @brief  StoreModel enregistre un asset soumis par son propriétaire et indexe son empreinte
// / @input  asset a1 (OwnerID alice@Org1, hash h1) signé par alice@Org1 (Org1MSP, active, contributor)
// / @expect transaction acceptée ; l'asset et l'entrée d'index h1 → a1 sont enregistrés
func TestStoreModel_NominalCase(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))

	resp := stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1")))
	if resp.Status != shim.OK {
		t.Fatalf("StoreModel refusé : %s", resp.Message)
	}
	if stub.state["a1"] == nil {
		t.Error("asset a1 absent du ledger")
	}
	if got := string(stub.state[hashKey(t, "sha256:h1")]); got != "a1" {
		t.Errorf("index h1 = %q, want a1", got)
	}
}

// / @brief  StoreModel ne réécrit jamais un asset déjà enregistré (RM06, RM19)
// / @input  a1 enregistré, puis seconde soumission de a1 avec un autre nom et un autre hash
// / @expect seconde transaction refusée ; l'enregistrement d'origine est inchangé
func TestStoreModel_ExistingID_Rejected(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))
	first := storeJSON(t, asset("a1", "alice@Org1", "sha256:h1"))
	if resp := stub.invoke("StoreModel", first); resp.Status != shim.OK {
		t.Fatalf("première soumission refusée : %s", resp.Message)
	}

	rewrite := asset("a1", "alice@Org1", "sha256:h2")
	rewrite.Name = "renommé"
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, rewrite)), "déjà enregistré")
	if string(stub.state["a1"]) != first {
		t.Error("l'enregistrement d'origine de a1 a été modifié")
	}
}

// / @brief  StoreModel refuse un fichier déjà enregistré sous un autre asset
// / @input  a1 (hash h1) enregistré, puis a2 avec le même hash h1
// / @expect a2 refusé, le message cite a1 ; a2 absent du ledger
func TestStoreModel_DuplicateHash_Rejected(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))
	stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1")))

	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a2", "alice@Org1", "sha256:h1"))), "a1")
	if stub.state["a2"] != nil {
		t.Error("a2 enregistré malgré une empreinte déjà présente")
	}
}

// / @brief  StoreModel exige que le parent déclaré soit enregistré sur le ledger
// / @input  enfant c1 dont ParentID vaut "fantome" (absent), puis c2 dont le parent p1 est enregistré
// / @expect c1 refusé ; c2 accepté
func TestStoreModel_ParentMustExist(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))

	orphan := asset("c1", "alice@Org1", "sha256:c1")
	orphan.ParentID = "fantome"
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, orphan)), "parent fantome introuvable")

	stub.invoke("StoreModel", storeJSON(t, asset("p1", "alice@Org1", "sha256:p1")))
	child := asset("c2", "alice@Org1", "sha256:c2")
	child.ParentID = "p1"
	if resp := stub.invoke("StoreModel", storeJSON(t, child)); resp.Status != shim.OK {
		t.Errorf("enfant d'un parent enregistré refusé : %s", resp.Message)
	}
}

// / @brief  StoreModel refuse un OwnerID qui n'est pas le signataire (ADR-12)
// / @input  asset déclaré propriété de bob@Org1, signé par alice@Org1
// / @expect transaction refusée
func TestStoreModel_OwnerNotSigner_Rejected(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a1", "bob@Org1", "sha256:h1"))), "différent du signataire")
}

// / @brief  StoreModel refuse un identifiant d'une organisation émis par le MSP d'une autre
// / @input  certificat alice@Org1 présenté par le MSP Org2MSP
// / @expect transaction refusée
func TestStoreModel_ForeignMSP_Rejected(t *testing.T) {
	forged := identity(t, "Org2MSP", "alice@Org1", map[string]string{"Myr.role": "contributor", "Myr.status": "active"})
	stub := newStub(t, forged)
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1"))), "Org2MSP")
}

// / @brief  StoreModel refuse une identité partagée d'organisation, sans identifiant pseudo@org ni attributs Myr
// / @input  asset déclaré propriété de orgadmin.org1, signé par orgadmin.org1 (sans attribut)
// / @expect transaction refusée : seule la signature par l'utilisateur permet d'écrire
func TestStoreModel_SharedOrgIdentity_Rejected(t *testing.T) {
	stub := newStub(t, identity(t, "Org1MSP", "orgadmin.org1", nil))
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a1", "orgadmin.org1", "sha256:h1"))), "pseudo@<org>")
}

// / @brief  StoreModel exige un signataire au statut active
// / @input  alice@Org1 contributor au statut pending
// / @expect transaction refusée
func TestStoreModel_PendingStatus_Rejected(t *testing.T) {
	pending := identity(t, "Org1MSP", "alice@Org1", map[string]string{"Myr.role": "contributor", "Myr.status": "pending"})
	stub := newStub(t, pending)
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1"))), "pending")
}

// / @brief  StoreModel exige un rôle autorisé à écrire
// / @input  alice@Org1 active au rôle reader
// / @expect transaction refusée
func TestStoreModel_ReaderRole_Rejected(t *testing.T) {
	reader := identity(t, "Org1MSP", "alice@Org1", map[string]string{"Myr.role": "reader", "Myr.status": "active"})
	stub := newStub(t, reader)
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1"))), "reader")
}

// / @brief  StoreModel n'enregistre qu'un asset soumis, sur le canal de la transaction, avec un ID valide
// / @input  asset en draft ; asset déclarant un autre canal ; asset d'ID "a/1"
// / @expect les trois transactions sont refusées
func TestStoreModel_InvalidRecord_Rejected(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))

	draft := asset("a1", "alice@Org1", "sha256:h1")
	draft.Status = model.ModuleDraft
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, draft)), "asset soumis")

	otherChannel := asset("a2", "alice@Org1", "sha256:h2")
	otherChannel.ChannelID = "autre"
	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, otherChannel)), "canal")

	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("a/1", "alice@Org1", "sha256:h3"))), "ID invalide")
}

// ── IndexHashes et ListModels ─────────────────────────────────────────────────

// / @brief  IndexHashes indexe les assets enregistrés avant l'unicité du fichier, signale les doublons et reste idempotent
// / @input  ledger avec L1 (h1), L2 (h1, doublon) et L3 (h2) sans index ; IndexHashes deux fois, puis StoreModel d'un asset de hash h2
// / @expect 1er appel : 2 empreintes indexées, doublon {h1, L1, L2} ; 2e appel : 0 indexée, même doublon ; le nouvel asset h2 est refusé
func TestIndexHashes_LegacyAssets_IndexedAndDuplicatesReported(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))
	for _, m := range []model.Model3D{
		asset("L1", "legacy", "sha256:h1"),
		asset("L2", "legacy", "sha256:h1"),
		asset("L3", "legacy", "sha256:h2"),
	} {
		stub.state[m.ID] = []byte(storeJSON(t, m))
	}

	type result struct {
		Indexed    int `json:"indexed"`
		Duplicates []struct {
			Hash, Kept, Duplicate string
		} `json:"duplicates"`
	}
	run := func() result {
		resp := stub.invoke("IndexHashes")
		if resp.Status != shim.OK {
			t.Fatalf("IndexHashes refusé : %s", resp.Message)
		}
		var r result
		if err := json.Unmarshal(resp.Payload, &r); err != nil {
			t.Fatalf("réponse IndexHashes : %v", err)
		}
		return r
	}

	first := run()
	if first.Indexed != 2 {
		t.Errorf("1er appel : indexed = %d, want 2", first.Indexed)
	}
	if len(first.Duplicates) != 1 || first.Duplicates[0].Kept != "L1" || first.Duplicates[0].Duplicate != "L2" {
		t.Errorf("1er appel : duplicates = %+v, want [{h1 L1 L2}]", first.Duplicates)
	}
	if second := run(); second.Indexed != 0 || len(second.Duplicates) != 1 {
		t.Errorf("2e appel : %+v, want indexed 0 et le même doublon", second)
	}

	expectRejected(t, stub.invoke("StoreModel", storeJSON(t, asset("n1", "alice@Org1", "sha256:h2"))), "L3")
}

// / @brief  ListModels ne retourne que des assets, jamais les entrées de l'index des empreintes
// / @input  deux assets enregistrés (donc deux entrées d'index)
// / @expect ListModels retourne exactement 2 éléments, tous des Model3D
func TestListModels_HashIndexHidden(t *testing.T) {
	stub := newStub(t, contributor(t, "alice@Org1"))
	stub.invoke("StoreModel", storeJSON(t, asset("a1", "alice@Org1", "sha256:h1")))
	stub.invoke("StoreModel", storeJSON(t, asset("a2", "alice@Org1", "sha256:h2")))

	resp := stub.invoke("ListModels", testChannel)
	if resp.Status != shim.OK {
		t.Fatalf("ListModels refusé : %s", resp.Message)
	}
	var list []model.Model3D
	if err := json.Unmarshal(resp.Payload, &list); err != nil {
		t.Fatalf("réponse ListModels : %v", err)
	}
	if len(list) != 2 {
		t.Errorf("ListModels = %d éléments, want 2", len(list))
	}
}
