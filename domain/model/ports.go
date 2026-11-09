// domain/model/ports.go — interfaces secondaires (out-ports) du domaine model
package model

// BlockchainPort est le contrat que tout adapter blockchain doit respecter.
// Intentionnellement minimal pour rester agnostique de la technologie sous-jacente.
type BlockchainPort interface {
	// StoreModelRecord enregistre un asset soumis ; un identifiant déjà
	// enregistré n'est jamais réécrit (ADR-12).
	StoreModelRecord(m *Model3D) error
	// UpdateModelMetadata remplace les métadonnées descriptives d'un asset
	// enregistré — seule modification permise après soumission, réservée à son
	// propriétaire (ADR-12).
	UpdateModelMetadata(id, channelID string, meta AssetMetadata) error
	// GetModelRecord récupère un modèle par son ID sur le canal indiqué.
	// channelID="" utilise le canal par défaut configuré dans l'adapter.
	// Un identifiant absent du ledger renvoie une erreur enveloppant
	// ErrNotFound.
	GetModelRecord(id, channelID string) (*Model3D, error)
	ListModelRecords(channelID string) ([]*Model3D, error)
}

// ConnectionStore gère les connexions d'assemblage entre assets.
// Implémenté par les adapters qui supportent le graphe (ex: JSONBlockchain).
// Non requis par FabricBlockchain.
type ConnectionStore interface {
	SaveConnection(c *Connection) error
	UpdateConnection(c *Connection) error
	RemoveConnection(id string) error
	ListConnections() ([]*Connection, error)
}

// ThumbnailStore persiste les miniatures STL (data URL base64).
// Séparé de Model3D pour ne pas alourdir les réponses blockchain.
type ThumbnailStore interface {
	SaveThumbnail(assetID, dataURL string) error
	GetThumbnail(assetID string) (string, error)
}

// OGImageFetcher récupère l'image représentative (balise og:image) d'une page web
// et la retourne encodée en data URL base64. Permet de redériver la miniature d'un
// asset depuis sa source durable (Model3D.Links) quand aucun fichier 3D n'est fourni —
// voir Service.RegenerateThumbnail.
type OGImageFetcher interface {
	FetchOGImage(pageURL string) (dataURL string, err error)
}

// DraftStore persiste les Model3D en brouillon avant leur soumission à la
// blockchain — cycle brouillon → soumission généralisé à toute création ou
// modification (composant, module, interface). Requis dès qu'un composant ou
// module peut être créé : AddFull/CreateModule échouent sans lui.
type DraftStore interface {
	SaveDraft(m *Model3D) error
	// GetDraft renvoie une erreur enveloppant ErrNotFound si aucun brouillon
	// ne porte cet identifiant ; toute autre erreur signale un stockage
	// inutilisable (fichier illisible…).
	GetDraft(id string) (*Model3D, error)
	RemoveDraft(id string) error
	ListDrafts(channelID string) ([]*Model3D, error)
}

// RemovedAssetStore masque localement un asset déjà soumis à la blockchain
// (composant ou module) sans jamais écrire sur le ledger — Fabric ne supporte
// aucune suppression (RM06/RM08). Un asset masqué disparaît des listes
// (List/ListModules) mais reste résolvable par identifiant direct (Get/GetModule).
type RemovedAssetStore interface {
	HideAsset(id string) error
	IsHidden(id string) (bool, error)
}

// InterfaceStore gère les interfaces physiques des assets et le vocabulaire de référence.
// Non requis par Fabric ou NoOpBlockchain — activé uniquement en mode GUI.
type InterfaceStore interface {
	// Interfaces sur un asset
	SaveInterface(iface *AssetInterface) error
	RemoveInterface(id string) error
	ListInterfacesForAsset(assetID string) ([]*AssetInterface, error)
	// GetInterface renvoie une erreur enveloppant ErrNotFound si aucune
	// interface ne porte cet identifiant.
	GetInterface(id string) (*AssetInterface, error)

	// Vocabulaire extensible (catégories, types, unités)
	GetRefs() (*InterfaceRefs, error)
	AddRefCategory(cat string) error
	AddRefType(cat, typeName string) error
	AddRefUnit(cat, unit string) error
}
