// adapters/out/localstorage/request_store.go — persistance JSON des demandes de compte
package localstorage

import (
	"myr-core/domain/identity"
)

var _ identity.RequestStore = (*JSONRequestStore)(nil)

// JSONRequestStore persiste les demandes de compte dans un fichier JSON,
// partagé par le serveur et la CLI (jsonFile).
type JSONRequestStore struct {
	file *jsonFile
}

func NewJSONRequestStore(path string) *JSONRequestStore {
	return &JSONRequestStore{file: newJSONFile(path, 0o600)}
}

func (s *JSONRequestStore) Save(req *identity.AccountRequest) error {
	var list []*identity.AccountRequest
	return s.file.update(&list, func() error {
		// Remplacer si l'ID existe déjà, sinon ajouter
		for i, r := range list {
			if r.ID == req.ID {
				list[i] = req
				return nil
			}
		}
		list = append(list, req)
		return nil
	})
}

func (s *JSONRequestStore) FindAll() ([]*identity.AccountRequest, error) {
	var list []*identity.AccountRequest
	if err := s.file.view(&list); err != nil {
		return nil, err
	}
	return list, nil
}
