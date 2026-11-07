// adapters/out/fabric/network_pool.go — pool de connexions Fabric multi-réseau
package fabric

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"

	"myr-core/domain/model"
	"myr-core/domain/network"
)

// NetworkPool gère un ensemble de GatewayClient ouverts, un par réseau Fabric.
// Thread-safe — peut être utilisé depuis plusieurs goroutines simultanées.
type NetworkPool struct {
	mu       sync.RWMutex
	gateways map[string]*GatewayClient // networkID → gateway (identité du profil réseau)
	// signers met en cache, par réseau puis par certificat d'utilisateur, les
	// passerelles dérivées qui signent avec l'identité de cet utilisateur
	// (ADR-12) — remplacées dès que le contenu du certificat change.
	signers map[string]map[string]signerGateway
}

// signerGateway est une passerelle dérivée et l'empreinte de l'identité
// (MSP + certificat) avec laquelle elle a été ouverte.
type signerGateway struct {
	identitySum string
	gw          *GatewayClient
}

// NewNetworkPool crée un pool vide.
func NewNetworkPool() *NetworkPool {
	return &NetworkPool{
		gateways: make(map[string]*GatewayClient),
		signers:  make(map[string]map[string]signerGateway),
	}
}

// Add ajoute (ou remplace) un gateway dans le pool.
func (p *NetworkPool) Add(networkID string, gw *GatewayClient) {
	p.mu.Lock()
	if old := p.gateways[networkID]; old != nil {
		p.closeSigners(networkID)
		old.Close()
	}
	p.gateways[networkID] = gw
	p.mu.Unlock()
}

// Get retourne le gateway associé au networkID, ou nil s'il est absent.
func (p *NetworkPool) Get(networkID string) *GatewayClient {
	p.mu.RLock()
	gw := p.gateways[networkID]
	p.mu.RUnlock()
	return gw
}

// Remove ferme et supprime un gateway du pool.
func (p *NetworkPool) Remove(networkID string) {
	p.mu.Lock()
	p.closeSigners(networkID)
	if gw := p.gateways[networkID]; gw != nil {
		gw.Close()
	}
	delete(p.gateways, networkID)
	p.mu.Unlock()
}

// IDs retourne les identifiants de tous les réseaux connectés.
func (p *NetworkPool) IDs() []string {
	p.mu.RLock()
	ids := make([]string, 0, len(p.gateways))
	for id := range p.gateways {
		ids = append(ids, id)
	}
	p.mu.RUnlock()
	return ids
}

// Close ferme toutes les connexions et vide le pool.
func (p *NetworkPool) Close() {
	p.mu.Lock()
	for id, gw := range p.gateways {
		p.closeSigners(id)
		gw.Close()
	}
	p.gateways = make(map[string]*GatewayClient)
	p.mu.Unlock()
}

// closeSigners ferme les passerelles dérivées d'un réseau, avant que sa
// connexion gRPC ne soit fermée. L'appelant détient p.mu.
func (p *NetworkPool) closeSigners(networkID string) {
	for _, s := range p.signers[networkID] {
		s.gw.Close()
	}
	delete(p.signers, networkID)
}

// BlockchainFor retourne un model.BlockchainPort connecté au réseau spécifié,
// en lecture seule : il signe avec l'identité du profil réseau, qui n'écrit
// jamais sur le ledger (ADR-12). Si le réseau est absent du pool, retourne un
// blockchain hors-ligne qui log l'erreur.
// Cette méthode implémente l'interface rest.blockchainRouter.
func (p *NetworkPool) BlockchainFor(networkID string) model.BlockchainPort {
	gw := p.Get(networkID)
	if gw == nil {
		return &offlineBlockchain{networkID: networkID}
	}
	return NewFabricBlockchain(gw)
}

// BlockchainForSigner retourne un model.BlockchainPort du réseau dont les
// écritures sont signées par l'identité de l'utilisateur (MSP, certificat et
// clé de son wallet — ADR-12). Si cette identité est inutilisable, l'erreur
// est journalisée et l'accès retourné est en lecture seule (BlockchainFor).
// Cette méthode implémente l'interface rest.blockchainRouter.
func (p *NetworkPool) BlockchainForSigner(networkID, mspID, certPath, keyPath string) model.BlockchainPort {
	gw, err := p.signerGateway(networkID, mspID, certPath, keyPath)
	if err != nil {
		log.Printf("réseau %s : identité %s inutilisable pour signer, accès en lecture seule : %v", networkID, certPath, err)
		return p.BlockchainFor(networkID)
	}
	return NewUserFabricBlockchain(gw)
}

// signerGateway retourne la passerelle du réseau signée par l'identité donnée,
// depuis le cache tant que cette identité (MSP + certificat) est inchangée.
func (p *NetworkPool) signerGateway(networkID, mspID, certPath, keyPath string) (*GatewayClient, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("certificat illisible : %w", err)
	}
	sum := sha256.Sum256(append([]byte(mspID+"\x00"), certPEM...))
	identitySum := hex.EncodeToString(sum[:])

	p.mu.Lock()
	defer p.mu.Unlock()
	base := p.gateways[networkID]
	if base == nil {
		return nil, fmt.Errorf("réseau %q non disponible dans le pool de connexions", networkID)
	}
	if cached, ok := p.signers[networkID][certPath]; ok {
		if cached.identitySum == identitySum {
			return cached.gw, nil
		}
		cached.gw.Close() // certificat renouvelé (ré-enrôlement)
	}
	gw, err := base.ForIdentity(mspID, certPath, keyPath)
	if err != nil {
		return nil, err
	}
	if p.signers[networkID] == nil {
		p.signers[networkID] = make(map[string]signerGateway)
	}
	p.signers[networkID][certPath] = signerGateway{identitySum: identitySum, gw: gw}
	return gw, nil
}

// LoadFromProfiles tente d'ouvrir une connexion pour chaque profil réseau actif.
// Les échecs sont loggés mais non fatals — le pool démarre avec les réseaux disponibles.
func (p *NetworkPool) LoadFromProfiles(profiles []*network.NetworkProfile) {
	for _, np := range profiles {
		p.connect(np)
	}
}

// Sync ouvre (ou réouvre) la connexion pour ce profil réseau. Implémente
// network.PoolSync — appelé par domain/network.Service à chaque création,
// activation ou modification d'un profil, pour que le pool reflète l'état
// courant sans attendre un redémarrage du processus (voir LoadFromProfiles,
// appelé une seule fois au démarrage).
func (p *NetworkPool) Sync(np *network.NetworkProfile) {
	p.connect(np)
}

// connect tente d'ouvrir la connexion Fabric pour un profil réseau donné et
// l'ajoute au pool en cas de succès. Échec loggé, jamais fatal — centralise
// la logique partagée par LoadFromProfiles (démarrage) et Sync (à chaud).
func (p *NetworkPool) connect(np *network.NetworkProfile) {
	cfg := ConfigFromProfile(np)
	gw, err := NewGatewayClient(cfg)
	if err != nil {
		log.Printf("réseau %s (%s) : connexion impossible : %v", np.Name, np.ID, err)
		return
	}
	p.Add(np.ID, gw)
	log.Printf("réseau %s (%s) : connecté → %s", np.Name, np.ID, np.PeerEndpoint)
}

// ── offlineBlockchain ─────────────────────────────────────────────────────────

// offlineBlockchain est retourné quand un réseau demandé n'est pas dans le pool.
type offlineBlockchain struct {
	networkID string
}

func (o *offlineBlockchain) err() error {
	return fmt.Errorf("réseau %q non disponible dans le pool de connexions", o.networkID)
}

func (o *offlineBlockchain) StoreModelRecord(*model.Model3D) error               { return o.err() }
func (o *offlineBlockchain) GetModelRecord(_, _ string) (*model.Model3D, error)  { return nil, o.err() }
func (o *offlineBlockchain) ListModelRecords(_ string) ([]*model.Model3D, error) { return nil, o.err() }
func (o *offlineBlockchain) VerifyIntegrity(_, _, _ string) (bool, error)        { return false, o.err() }
