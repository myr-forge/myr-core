// adapters/out/localstorage/node_store.go
// Bookkeeping local des nœuds (peer/orderer) provisionnés par myr sur un canal.
// Nécessaire car Hyperledger Fabric ne trace pas les peers individuellement dans
// la configuration de canal (seules les organisations y figurent) — seul myr sait
// quels nœuds il a provisionnés, information utilisée pour appliquer le seuil
// minimum de nœuds actifs (RM27) et pour retrouver un nœud à retirer.
package localstorage

import (
	"fmt"
	"time"
)

// ProvisionedNode décrit un nœud (peer ou orderer) provisionné sur un canal donné.
type ProvisionedNode struct {
	ChannelID string    `json:"channel_id"`
	Type      string    `json:"type"` // "peer" | "orderer"
	Addr      string    `json:"addr"`
	OrgMSP    string    `json:"org_msp"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type nodeDB struct {
	Nodes []*ProvisionedNode `json:"nodes"`
}

// JSONNodeStore persiste les nœuds provisionnés dans un fichier JSON, partagé
// par le serveur et la CLI (jsonFile).
type JSONNodeStore struct {
	file *jsonFile
}

func NewJSONNodeStore(path string) *JSONNodeStore {
	return &JSONNodeStore{file: newJSONFile(path, 0o600)}
}

// read retourne les nœuds enregistrés.
func (s *JSONNodeStore) read() ([]*ProvisionedNode, error) {
	var db nodeDB
	if err := s.file.view(&db); err != nil {
		return nil, err
	}
	return db.Nodes, nil
}

// Add enregistre un nœud (upsert par ChannelID+Addr).
func (s *JSONNodeStore) Add(n *ProvisionedNode) error {
	var db nodeDB
	return s.file.update(&db, func() error {
		for i, existing := range db.Nodes {
			if existing.ChannelID == n.ChannelID && existing.Addr == n.Addr {
				db.Nodes[i] = n
				return nil
			}
		}
		if n.CreatedAt.IsZero() {
			n.CreatedAt = time.Now()
		}
		db.Nodes = append(db.Nodes, n)
		return nil
	})
}

// Deactivate marque un nœud comme inactif (retiré du canal) sans supprimer son historique.
func (s *JSONNodeStore) Deactivate(channelID, addr string) error {
	var db nodeDB
	return s.file.update(&db, func() error {
		for _, n := range db.Nodes {
			if n.ChannelID == channelID && n.Addr == addr {
				n.Active = false
				return nil
			}
		}
		return fmt.Errorf("nœud %s introuvable sur le canal %s", addr, channelID)
	})
}

// Find retourne le nœud correspondant, ou nil si absent.
func (s *JSONNodeStore) Find(channelID, addr string) (*ProvisionedNode, error) {
	nodes, err := s.read()
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if n.ChannelID == channelID && n.Addr == addr {
			return n, nil
		}
	}
	return nil, nil
}

// ActiveCount retourne le nombre de nœuds actifs sur un canal.
func (s *JSONNodeStore) ActiveCount(channelID string) (int, error) {
	nodes, err := s.read()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, n := range nodes {
		if n.ChannelID == channelID && n.Active {
			count++
		}
	}
	return count, nil
}

// ListActive retourne les nœuds actifs d'un canal, filtrés par type si fourni ("" = tous).
func (s *JSONNodeStore) ListActive(channelID, nodeType string) ([]*ProvisionedNode, error) {
	nodes, err := s.read()
	if err != nil {
		return nil, err
	}
	var out []*ProvisionedNode
	for _, n := range nodes {
		if n.ChannelID == channelID && n.Active && (nodeType == "" || n.Type == nodeType) {
			out = append(out, n)
		}
	}
	return out, nil
}
