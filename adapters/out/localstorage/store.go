// adapters/out/localstorage/store.go
package localstorage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"myr-core/domain/model"
)

var _ model.FileStoragePort = (*LocalStorage)(nil)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, 0755)
	return &LocalStorage{basePath: basePath}
}

func (l *LocalStorage) Upload(filePath string) (string, error) {
	dest := filepath.Join(l.basePath, filepath.Base(filePath))
	src, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open source: %w", err)
	}
	defer src.Close()
	dst, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create dest: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}
	return dest, nil
}

func (l *LocalStorage) Download(hash, destPath string) error {
	src, _ := os.Open(hash)
	defer src.Close()
	dst, _ := os.Create(destPath)
	defer dst.Close()
	_, err := io.Copy(dst, src)
	return err
}

func (l *LocalStorage) Delete(hash string) error {
	return os.Remove(hash)
}

// Verify recalcule le hash SHA-256 du fichier situé à ref (le chemin local
// retourné par Upload) et le compare à expectedHash. Retourne ref comme
// emplacement si le fichier existe, même en cas de divergence de hash — vide
// uniquement si le fichier est introuvable à cet emplacement.
func (l *LocalStorage) Verify(ref, expectedHash string) (string, bool, error) {
	if ref == "" {
		return "", false, nil
	}
	f, err := os.Open(ref)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("verify open: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ref, false, fmt.Errorf("verify hash: %w", err)
	}
	sum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	return ref, sum == expectedHash, nil
}
