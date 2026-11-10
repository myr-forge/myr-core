// adapters/out/localstorage/json_blockchain.go
// Implémente ConnectionStore + ThumbnailStore + InterfaceStore
// via un fichier JSON local. Stocke les données GUI locales (connexions,
// miniatures, interfaces) indépendamment de la blockchain Fabric.
package localstorage

import (
	"fmt"

	"myr-core/domain/model"
)

type dbFile struct {
	Thumbnails  map[string]string    `json:"thumbnails"`
	Connections []*model.Connection  `json:"connections"`
	Interfaces  []*storedInterface   `json:"interfaces"`
	Refs        *model.InterfaceRefs `json:"refs"`
	Drafts      []*model.Model3D     `json:"drafts"`
	HiddenIDs   map[string]bool      `json:"hidden_ids"`
}

// storedInterface est une interface telle qu'enregistrée. Les fichiers
// antérieurs au retrait des slots virtuels (FT-040) en contiennent, marqués
// "virtual" : ils sont écartés à la lecture et disparaissent à l'écriture
// suivante.
type storedInterface struct {
	model.AssetInterface
	Virtual bool `json:"virtual,omitempty"`
}

var _ model.RemovedAssetStore = (*JSONBlockchain)(nil)

// JSONBlockchain persiste les données GUI locales (connexions, miniatures, interfaces, brouillons).
// Il implémente ConnectionStore, ThumbnailStore, InterfaceStore, DraftStore et RemovedAssetStore.
// Le fichier est partagé par le serveur et la CLI (jsonFile : verrou, écriture
// atomique, fichier illisible jamais écrasé).
type JSONBlockchain struct {
	file *jsonFile
}

func NewJSONBlockchain(path string) *JSONBlockchain {
	return &JSONBlockchain{file: newJSONFile(path, 0o600)}
}

// ── persistance interne ──────────────────────────────────────────────────────

// withDefaults initialise les champs absents du fichier, pour qu'aucun
// appelant ne manipule une carte ou une liste nulle.
func (db *dbFile) withDefaults() *dbFile {
	if db.Thumbnails == nil {
		db.Thumbnails = map[string]string{}
	}
	if db.Connections == nil {
		db.Connections = []*model.Connection{}
	}
	kept := []*storedInterface{}
	for _, i := range db.Interfaces {
		if !i.Virtual {
			kept = append(kept, i)
		}
	}
	db.Interfaces = kept
	if db.Drafts == nil {
		db.Drafts = []*model.Model3D{}
	}
	if db.HiddenIDs == nil {
		db.HiddenIDs = map[string]bool{}
	}
	// Refs : appliquer les valeurs par défaut si absentes ou sans catégories.
	if db.Refs == nil || len(db.Refs.Categories) == 0 {
		db.Refs = model.DefaultInterfaceRefs()
	}
	return db
}

// read retourne le contenu du fichier.
func (j *JSONBlockchain) read() (*dbFile, error) {
	var db dbFile
	if err := j.file.view(&db); err != nil {
		return nil, err
	}
	return db.withDefaults(), nil
}

// update applique change au contenu du fichier et l'écrit, sous verrou.
func (j *JSONBlockchain) update(change func(db *dbFile) error) error {
	var db dbFile
	return j.file.update(&db, func() error { return change(db.withDefaults()) })
}

// ── model.ConnectionStore ────────────────────────────────────────────────────

// SaveConnection enregistre c, ou remplace la connexion de même ID. Le refus
// d'une liaison en double relève du domaine (RM09, AddConnection).
func (j *JSONBlockchain) SaveConnection(c *model.Connection) error {
	return j.update(func(db *dbFile) error {
		for i, existing := range db.Connections {
			if existing.ID == c.ID {
				db.Connections[i] = c
				return nil
			}
		}
		db.Connections = append(db.Connections, c)
		return nil
	})
}

func (j *JSONBlockchain) UpdateConnection(c *model.Connection) error {
	return j.update(func(db *dbFile) error {
		for i, existing := range db.Connections {
			if existing.ID == c.ID {
				db.Connections[i] = c
				return nil
			}
		}
		return fmt.Errorf("connexion %s introuvable", c.ID)
	})
}

func (j *JSONBlockchain) RemoveConnection(id string) error {
	return j.update(func(db *dbFile) error {
		filtered := db.Connections[:0]
		for _, c := range db.Connections {
			if c.ID != id {
				filtered = append(filtered, c)
			}
		}
		db.Connections = filtered
		return nil
	})
}

func (j *JSONBlockchain) ListConnections() ([]*model.Connection, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	return db.Connections, nil
}

// ── model.ThumbnailStore ─────────────────────────────────────────────────────

func (j *JSONBlockchain) SaveThumbnail(assetID, dataURL string) error {
	return j.update(func(db *dbFile) error {
		db.Thumbnails[assetID] = dataURL
		return nil
	})
}

func (j *JSONBlockchain) GetThumbnail(assetID string) (string, error) {
	db, err := j.read()
	if err != nil {
		return "", err
	}
	return db.Thumbnails[assetID], nil
}

// ── model.InterfaceStore ─────────────────────────────────────────────────────

func (j *JSONBlockchain) SaveInterface(iface *model.AssetInterface) error {
	return j.update(func(db *dbFile) error {
		stored := &storedInterface{AssetInterface: *iface}
		for i, existing := range db.Interfaces {
			if existing.ID == iface.ID {
				db.Interfaces[i] = stored
				return nil
			}
		}
		db.Interfaces = append(db.Interfaces, stored)
		return nil
	})
}

func (j *JSONBlockchain) RemoveInterface(id string) error {
	return j.update(func(db *dbFile) error {
		filtered := db.Interfaces[:0]
		for _, iface := range db.Interfaces {
			if iface.ID != id {
				filtered = append(filtered, iface)
			}
		}
		db.Interfaces = filtered
		return nil
	})
}

func (j *JSONBlockchain) ListInterfacesForAsset(assetID string) ([]*model.AssetInterface, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	var result []*model.AssetInterface
	for _, iface := range db.Interfaces {
		if iface.AssetID == assetID {
			result = append(result, &iface.AssetInterface)
		}
	}
	return result, nil
}

func (j *JSONBlockchain) GetInterface(id string) (*model.AssetInterface, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	for _, iface := range db.Interfaces {
		if iface.ID == id {
			return &iface.AssetInterface, nil
		}
	}
	return nil, fmt.Errorf("%w : interface %q", model.ErrNotFound, id)
}

// ── model.DraftStore ─────────────────────────────────────────────────────────

func (j *JSONBlockchain) SaveDraft(m *model.Model3D) error {
	return j.update(func(db *dbFile) error {
		for i, existing := range db.Drafts {
			if existing.ID == m.ID {
				db.Drafts[i] = m
				return nil
			}
		}
		db.Drafts = append(db.Drafts, m)
		return nil
	})
}

func (j *JSONBlockchain) GetDraft(id string) (*model.Model3D, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	for _, m := range db.Drafts {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%w : brouillon %q", model.ErrNotFound, id)
}

func (j *JSONBlockchain) RemoveDraft(id string) error {
	return j.update(func(db *dbFile) error {
		filtered := db.Drafts[:0]
		for _, m := range db.Drafts {
			if m.ID != id {
				filtered = append(filtered, m)
			}
		}
		db.Drafts = filtered
		return nil
	})
}

func (j *JSONBlockchain) ListDrafts(channelID string) ([]*model.Model3D, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	if channelID == "" {
		return db.Drafts, nil
	}
	var result []*model.Model3D
	for _, m := range db.Drafts {
		if m.ChannelID == channelID {
			result = append(result, m)
		}
	}
	return result, nil
}

// ── model.RemovedAssetStore ──────────────────────────────────────────────────

func (j *JSONBlockchain) HideAsset(id string) error {
	return j.update(func(db *dbFile) error {
		db.HiddenIDs[id] = true
		return nil
	})
}

func (j *JSONBlockchain) IsHidden(id string) (bool, error) {
	db, err := j.read()
	if err != nil {
		return false, err
	}
	return db.HiddenIDs[id], nil
}

func (j *JSONBlockchain) GetRefs() (*model.InterfaceRefs, error) {
	db, err := j.read()
	if err != nil {
		return nil, err
	}
	return db.Refs, nil // toujours non-nil grâce à withDefaults
}

func (j *JSONBlockchain) AddRefCategory(cat string) error {
	return j.update(func(db *dbFile) error {
		for _, c := range db.Refs.Categories {
			if c == cat {
				return nil
			}
		}
		db.Refs.Categories = append(db.Refs.Categories, cat)
		if db.Refs.Types == nil {
			db.Refs.Types = map[string][]string{}
		}
		if db.Refs.Units == nil {
			db.Refs.Units = map[string][]string{}
		}
		if _, ok := db.Refs.Types[cat]; !ok {
			db.Refs.Types[cat] = []string{}
		}
		if _, ok := db.Refs.Units[cat]; !ok {
			db.Refs.Units[cat] = []string{}
		}
		return nil
	})
}

func (j *JSONBlockchain) AddRefType(cat, typeName string) error {
	return j.update(func(db *dbFile) error {
		for _, t := range db.Refs.Types[cat] {
			if t == typeName {
				return nil
			}
		}
		db.Refs.Types[cat] = append(db.Refs.Types[cat], typeName)
		return nil
	})
}

func (j *JSONBlockchain) AddRefUnit(cat, unit string) error {
	return j.update(func(db *dbFile) error {
		for _, u := range db.Refs.Units[cat] {
			if u == unit {
				return nil
			}
		}
		db.Refs.Units[cat] = append(db.Refs.Units[cat], unit)
		return nil
	})
}
