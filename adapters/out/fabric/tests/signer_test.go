// adapters/out/fabric/tests/signer_test.go
// Signature des écritures par l'identité de l'utilisateur (ADR-12) : passerelle
// dérivée (GatewayClient.ForIdentity) et routage NetworkPool lecture seule /
// signé. Aucun peer requis : la connexion gRPC n'est ouverte qu'au premier
// appel, vers une adresse locale fermée.
package fabric_test

import (
	"errors"
	"path/filepath"
	"testing"

	"myr-core/adapters/out/fabric"
	"myr-core/domain/model"
)

var (
	testCert = filepath.Join("..", "testdata", "client-cert.pem")
	testKey  = filepath.Join("..", "testdata", "client-key.pem")
)

// newTestGateway ouvre une passerelle sur un peer local inexistant, signée
// par l'identité de testdata (identité du profil réseau).
func newTestGateway(t *testing.T) *fabric.GatewayClient {
	t.Helper()
	gw, err := fabric.NewGatewayClient(fabric.Config{
		PeerEndpoint:  "localhost:1",
		GatewayPeer:   "peer0.org1.example.com",
		MSPID:         "Org1MSP",
		CertPath:      testCert,
		KeyPath:       testKey,
		TLSCertPath:   filepath.Join("..", "testdata", "tls-cert.pem"),
		FabricChannel: "ch1",
		ChaincodeName: "myrcc",
	})
	if err != nil {
		t.Fatalf("NewGatewayClient : %v", err)
	}
	return gw
}

// / @brief  ForIdentity dérive une passerelle signée par une autre identité, et refuse une identité illisible
// / @input  passerelle de profil ; ForIdentity avec le certificat et la clé de testdata, puis avec une clé absente
// / @expect première dérivation sans erreur ; seconde en erreur
func TestGatewayClient_ForIdentity(t *testing.T) {
	base := newTestGateway(t)
	defer base.Close()

	derived, err := base.ForIdentity("Org1MSP", testCert, testKey)
	if err != nil {
		t.Fatalf("ForIdentity : %v", err)
	}
	derived.Close() // ne doit pas fermer la connexion de base

	if _, err := base.ForIdentity("Org1MSP", testCert, filepath.Join(t.TempDir(), "absente.pem")); err == nil {
		t.Error("ForIdentity avec une clé absente : erreur attendue")
	}
}

// / @brief  NetworkPool n'écrit qu'avec l'identité d'un utilisateur ; l'identité du profil réseau reste en lecture seule (ADR-12)
// / @input  pool contenant le réseau net1 ; StoreModelRecord via BlockchainFor, via BlockchainForSigner (wallet valide), via BlockchainForSigner (certificat absent)
// / @expect BlockchainFor : ErrSignerRequired ; signataire valide : la transaction est tentée (peer injoignable) ; certificat absent : repli en lecture seule, ErrSignerRequired
func TestNetworkPool_BlockchainForSigner(t *testing.T) {
	pool := fabric.NewNetworkPool()
	pool.Add("net1", newTestGateway(t))
	defer pool.Close()
	asset := &model.Model3D{ID: "a1", ChannelID: "ch1", Status: model.ModuleSubmitted}

	if err := pool.BlockchainFor("net1").StoreModelRecord(asset); !errors.Is(err, model.ErrSignerRequired) {
		t.Errorf("BlockchainFor : err = %v, want ErrSignerRequired", err)
	}

	err := pool.BlockchainForSigner("net1", "Org1MSP", testCert, testKey).StoreModelRecord(asset)
	if errors.Is(err, model.ErrSignerRequired) {
		t.Errorf("signataire valide : écriture refusée faute de signataire")
	}
	if !errors.Is(err, model.ErrBlockchainUnreachable) {
		t.Errorf("signataire valide : err = %v, want une tentative de transaction vers le peer (ErrBlockchainUnreachable)", err)
	}

	missing := filepath.Join(t.TempDir(), "absent.pem")
	if err := pool.BlockchainForSigner("net1", "Org1MSP", missing, testKey).StoreModelRecord(asset); !errors.Is(err, model.ErrSignerRequired) {
		t.Errorf("certificat absent : err = %v, want ErrSignerRequired (repli en lecture seule)", err)
	}
}
