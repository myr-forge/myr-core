// adapters/out/localstorage/session_store.go
package localstorage

import (
	"errors"
	"os"

	"myr-core/domain/session"
)

var _ session.Store = (*JSONSessionStore)(nil)

// JSONSessionStore persiste le compte local de la CLI (un seul
// enregistrement) dans un fichier JSON lisible par son seul utilisateur.
type JSONSessionStore struct {
	file *jsonFile
}

func NewJSONSessionStore(path string) *JSONSessionStore {
	return &JSONSessionStore{file: newJSONFile(path, 0o600)}
}

func (s *JSONSessionStore) Save(sess *session.Session) error {
	var current *session.Session
	return s.file.update(&current, func() error {
		current = sess
		return nil
	})
}

func (s *JSONSessionStore) Load() (*session.Session, error) {
	var sess *session.Session
	if err := s.file.view(&sess); err != nil {
		return nil, err
	}
	return sess, nil // nil : pas de session active
}

func (s *JSONSessionStore) Clear() error {
	err := os.Remove(s.file.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
