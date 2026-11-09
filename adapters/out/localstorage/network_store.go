// adapters/out/localstorage/network_store.go
package localstorage

import (
	"myr-core/domain/network"
)

var _ network.Repo = (*JSONNetworkStore)(nil)

// JSONNetworkStore persiste les profils réseau dans un fichier JSON, partagé
// par le serveur et la CLI (jsonFile).
type JSONNetworkStore struct {
	file *jsonFile
}

func NewJSONNetworkStore(path string) *JSONNetworkStore {
	return &JSONNetworkStore{file: newJSONFile(path, 0o600)}
}

func (s *JSONNetworkStore) Save(n *network.NetworkProfile) error {
	var list []*network.NetworkProfile
	return s.file.update(&list, func() error {
		for i, existing := range list {
			if existing.ID == n.ID {
				list[i] = n
				return nil
			}
		}
		list = append(list, n)
		return nil
	})
}

func (s *JSONNetworkStore) FindAll() ([]*network.NetworkProfile, error) {
	var list []*network.NetworkProfile
	if err := s.file.view(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *JSONNetworkStore) Delete(id string) error {
	var list []*network.NetworkProfile
	return s.file.update(&list, func() error {
		filtered := list[:0]
		for _, n := range list {
			if n.ID != id {
				filtered = append(filtered, n)
			}
		}
		list = filtered
		return nil
	})
}
