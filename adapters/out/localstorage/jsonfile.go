// adapters/out/localstorage/jsonfile.go — lecture et écriture sûres des
// fichiers JSON de l'état local, partagées par tous les stores du paquet.
package localstorage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// ErrCorruptStore est retournée quand un fichier de l'état local est
// illisible : il n'est jamais traité comme vide ni réécrit — sinon la première
// écriture effacerait tout son contenu. L'administrateur restaure une
// sauvegarde ou répare le fichier.
var ErrCorruptStore = errors.New("localstorage: fichier de l'état local illisible")

// jsonFile est un fichier JSON de l'état local :
//   - un fichier absent ou vide est un état vide ; un fichier illisible est
//     une erreur (ErrCorruptStore), jamais un état vide ;
//   - une écriture passe par un fichier temporaire renommé : un arrêt brutal ne
//     laisse jamais un fichier à moitié écrit ;
//   - un verrou de fichier (<chemin>.lock) sérialise les accès des processus
//     qui partagent ce fichier (serveur et CLI), en plus du verrou interne.
type jsonFile struct {
	path string
	perm os.FileMode
	mu   sync.Mutex
}

func newJSONFile(path string, perm os.FileMode) *jsonFile {
	return &jsonFile{path: path, perm: perm}
}

// view décode le fichier dans v, sous verrou partagé.
func (f *jsonFile) view(v any) error {
	return f.locked(false, func() error { return f.decode(v) })
}

// update décode le fichier dans v, applique change puis écrit v, le tout sous
// verrou exclusif : aucune autre lecture-écriture ne s'intercale. Si change
// retourne une erreur, rien n'est écrit.
func (f *jsonFile) update(v any, change func() error) error {
	return f.locked(true, func() error {
		if err := f.decode(v); err != nil {
			return err
		}
		if err := change(); err != nil {
			return err
		}
		return f.encode(v)
	})
}

func (f *jsonFile) locked(exclusive bool, fn func() error) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	lf, err := os.OpenFile(f.path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("verrou de %s : %w", f.path, err)
	}
	defer lf.Close()
	if err := lockFile(lf, exclusive); err != nil {
		return fmt.Errorf("verrou de %s : %w", f.path, err)
	}
	defer unlockFile(lf) //nolint:errcheck
	return fn()
}

func (f *jsonFile) decode(v any) error {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%w : %s : %v", ErrCorruptStore, f.path, err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%w : %s : %v — le fichier n'est pas modifié ; restaurez une sauvegarde ou réparez-le", ErrCorruptStore, f.path, err)
	}
	return nil
}

func (f *jsonFile) encode(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(f.path), filepath.Base(f.path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // sans effet après le renommage
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, f.perm); err != nil {
		return err
	}
	return os.Rename(tmpPath, f.path)
}
