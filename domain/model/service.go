// domain/model/service.go — logique métier, ne connaît QUE les interfaces
package model

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

type Service struct {
	blockchain   BlockchainPort    // nullable — nil si Fabric indisponible (voir ErrBlockchainUnavailable)
	connStore    ConnectionStore   // optionnel — nil hors mode GUI
	ifaceStore   InterfaceStore    // optionnel — nil hors mode GUI
	draftStore   DraftStore        // optionnel — requis uniquement pour AddRequest.Draft / Submit
	removedStore RemovedAssetStore // optionnel — masquage local d'un asset déjà soumis (RM08)
	// actor est l'identité pour laquelle le service agit (As) ; acting indique
	// qu'As a été appelé.
	actor  string
	acting bool
}

// As retourne une copie du service agissant pour actor (voir ModelService.As).
func (s *Service) As(actor string) ModelService {
	cp := *s
	cp.actor, cp.acting = actor, true
	return &cp
}

// checkOwner refuse qu'une autre identité que le propriétaire modifie,
// supprime ou soumette m (FT-001). Un service utilisé sans As n'agit pour
// personne en particulier et ne contrôle rien : les adapters d'entrée
// appellent toujours As.
func (s *Service) checkOwner(m *Model3D) error {
	if !s.acting {
		return nil
	}
	if s.actor == "" {
		return fmt.Errorf("%w : aucune identité (connectez-vous ou créez un compte local)", ErrNotOwner)
	}
	if m.OwnerID != s.actor {
		return fmt.Errorf("%w : %s", ErrNotOwner, m.ID)
	}
	return nil
}

// creator retourne le propriétaire d'un asset créé : l'identité pour laquelle
// le service agit, sinon owner tel que fourni.
func (s *Service) creator(owner string) (string, error) {
	if !s.acting {
		return owner, nil
	}
	if s.actor == "" {
		return "", fmt.Errorf("%w : aucune identité (connectez-vous ou créez un compte local)", ErrNotOwner)
	}
	return s.actor, nil
}

func NewService(bc BlockchainPort) *Service {
	return &Service{blockchain: bc}
}

// SetBlockchain remplace la blockchain active (utilisé après reconnexion).
func (s *Service) SetBlockchain(bc BlockchainPort) { s.blockchain = bc }

// WithConnStore attache un store de connexions (mode GUI).
func (s *Service) WithConnStore(cs ConnectionStore) *Service {
	s.connStore = cs
	return s
}

// WithIfaceStore attache un store d'interfaces physiques (mode GUI).
func (s *Service) WithIfaceStore(is InterfaceStore) *Service {
	s.ifaceStore = is
	return s
}

// WithDraftStore attache un store de brouillons (requis pour AddRequest.Draft et Submit).
func (s *Service) WithDraftStore(ds DraftStore) *Service {
	s.draftStore = ds
	return s
}

// WithRemovedStore attache le registre de masquage local (RM08) — requis pour
// que Remove() puisse masquer un asset déjà soumis au lieu d'échouer.
func (s *Service) WithRemovedStore(rs RemovedAssetStore) *Service {
	s.removedStore = rs
	return s
}

// AddFull crée un asset avec tous ses champs (mode GUI).
// Toute création est un brouillon (DraftStore, aucune transaction blockchain) —
// l'asset ne rejoint la blockchain qu'à un appel explicite à Submit.
func (s *Service) AddFull(req AddRequest) (*Model3D, error) {
	if s.draftStore == nil {
		return nil, fmt.Errorf("stockage de brouillons non configuré")
	}
	if !req.Category.Valid() {
		return nil, fmt.Errorf("%w : %q", ErrInvalidCategory, req.Category)
	}
	if req.ChannelID == "" {
		return nil, ErrChannelRequired
	}
	owner, err := s.creator(req.OwnerID)
	if err != nil {
		return nil, err
	}
	req.OwnerID = owner
	// Licence : présente au catalogue, compatible avec celle du parent s'il est
	// connu de ce nœud (RM03) — le parent peut lui-même être encore un brouillon.
	var parent *Model3D
	if req.ParentID != "" {
		parent, _, _ = s.getAsset(req.ParentID, req.ChannelID)
	}
	if err := checkLicense(req.LicenseID, parent); err != nil {
		return nil, err
	}

	m := &Model3D{
		ID:          generateID(),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		ParentID:    req.ParentID,
		ChannelID:   req.ChannelID,
		OwnerID:     req.OwnerID,
		LicenseID:   req.LicenseID,
		Tags:        req.Tags,
		Links:       req.Links,
		CreatedAt:   time.Now(),
	}

	if req.FilePath != "" {
		// Le fichier n'est jamais conservé par myr-core : seul son empreinte SHA-256
		// est retenue (traçabilité RM01/RM04) — voir RM42 pour la vérification des
		// emplacements externes, seul mécanisme de localisation du fichier réel.
		hash, err := s.hashFile(req.FilePath)
		if err != nil {
			return nil, fmt.Errorf("hashing file: %w", err)
		}
		m.Hash = hash
		m.Versions = []Version{{Number: 1, Hash: hash, CreatedAt: time.Now()}}
	}

	// Dérivé d'un asset décomposé (module) : il reprend une copie indépendante
	// de sa composition — c'est ainsi qu'évolue un module soumis, figé (UCMOD01,
	// RM19).
	if parent != nil && (parent.IsModule() || len(parent.WorkspaceInstances) > 0) {
		m.WorkspaceInstances = []WorkspaceInstance{}
		m.Assemblies = []string{}
		if err := s.copyModuleComposition(parent, m); err != nil {
			return nil, fmt.Errorf("copie de la composition du parent : %w", err)
		}
	}

	m.Status = ModuleDraft
	if err := s.draftStore.SaveDraft(m); err != nil {
		return nil, fmt.Errorf("saving draft: %w", err)
	}
	return m, nil
}

// Submit engage sur la blockchain un composant créé en brouillon (AddRequest.Draft) :
// une seule transaction committe les métadonnées et les interfaces locales
// (Model3D.Interfaces, depuis InterfaceStore), puis Status passe à submitted.
// Point d'entrée unique pour tout Model3D (ADR-11, specs/3-Conception/Conception_intro.md) :
// dès que l'asset a des Assemblies (IsModule()), délègue à SubmitModule (RM17,
// création d'une ModuleVersion) — l'appelant (REST, CLI) n'a pas à savoir
// laquelle des deux logiques invoquer (règle 28 : ne pas dupliquer ce
// dispatch dans les adapters). note n'est utilisée que dans ce cas (ModuleVersion.Note) ;
// ignorée pour un composant simple. Un asset déjà soumis n'est jamais
// re-soumis (ErrAssetSubmitted, RM19, ADR-12).
func (s *Service) Submit(assetID, note string) (*Model3D, error) {
	if s.draftStore == nil {
		return nil, fmt.Errorf("stockage de brouillons non configuré")
	}
	if s.blockchain == nil {
		return nil, ErrBlockchainUnavailable
	}
	m, err := s.submittableDraft(assetID)
	if err != nil {
		return nil, err
	}
	if m.IsModule() {
		return s.SubmitModule(assetID, note)
	}
	unlock := lockDraft(assetID)
	defer unlock()
	if m, err = s.submittableDraft(assetID); err != nil { // relu sous verrou
		return nil, err
	}
	if s.ifaceStore != nil {
		ifaces, err := s.ifaceStore.ListInterfacesForAsset(assetID)
		if err != nil {
			return nil, fmt.Errorf("listing interfaces: %w", err)
		}
		m.Interfaces = make([]AssetInterface, len(ifaces))
		for i, iface := range ifaces {
			m.Interfaces[i] = *iface
		}
	}
	m.Status = ModuleSubmitted
	if err := s.blockchain.StoreModelRecord(m); err != nil {
		return nil, fmt.Errorf("storing on blockchain: %w", err)
	}
	_ = s.draftStore.RemoveDraft(assetID)
	return m, nil
}

// getAsset récupère un Model3D par son ID : d'abord dans DraftStore (état
// mutable, pré-soumission), puis sur la blockchain (état immuable, soumis).
// Généralise le cycle brouillon → soumission (ADR-02) à tout point d'accès du
// domaine — indique en retour si l'asset provient du brouillon, seul état
// modifiable (voir draftAsset).
func (s *Service) getAsset(id, channelID string) (*Model3D, bool, error) {
	if id == "" {
		return nil, false, fmt.Errorf("%w : identifiant d'asset vide", ErrNotFound)
	}
	if s.draftStore != nil {
		m, err := s.draftStore.GetDraft(id)
		if err == nil {
			return m, true, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return nil, false, err // stockage local inutilisable : ne pas le masquer par le ledger
		}
	}
	if s.blockchain == nil {
		return nil, false, ErrBlockchainUnavailable
	}
	m, err := s.blockchain.GetModelRecord(id, channelID)
	return m, false, err
}

// draftAsset retourne le brouillon id pour le modifier : un asset soumis est
// immuable (ErrAssetSubmitted, RM19, ADR-12) — licence, interfaces et
// composition ne changent plus, toute évolution passe par un fork. Seules ses
// métadonnées descriptives restent modifiables (updateSubmittedMetadata).
func (s *Service) draftAsset(id string) (*Model3D, error) {
	m, isDraft, err := s.getAsset(id, "")
	if err != nil {
		return nil, err
	}
	if !isDraft {
		return nil, fmt.Errorf("%w : %s", ErrAssetSubmitted, id)
	}
	if err := s.checkOwner(m); err != nil {
		return nil, err
	}
	return m, nil
}

// submittableDraft retourne le brouillon id prêt à être soumis : ni un asset
// déjà soumis, ni un brouillon dont l'identifiant est déjà enregistré sur le
// ledger (brouillon rouvert sous l'identifiant d'un asset soumis) — le ledger
// n'accepte qu'une écriture par identifiant (ADR-12).
func (s *Service) submittableDraft(id string) (*Model3D, error) {
	m, err := s.draftAsset(id)
	if err != nil {
		if errors.Is(err, ErrAssetSubmitted) {
			return nil, err
		}
		return nil, fmt.Errorf("brouillon introuvable : %w", err)
	}
	if _, err := s.blockchain.GetModelRecord(id, m.ChannelID); err == nil {
		return nil, fmt.Errorf("%w : %s est déjà enregistré sur le ledger", ErrAssetSubmitted, id)
	}
	return m, nil
}

// saveDraft enregistre un brouillon modifié.
func (s *Service) saveDraft(m *Model3D) error {
	if s.draftStore == nil {
		return fmt.Errorf("stockage de brouillons non configuré")
	}
	return s.draftStore.SaveDraft(m)
}

// editDraft applique change au brouillon id puis l'enregistre, sous le verrou
// de ce brouillon (R-21) ; ErrAssetSubmitted pour un asset soumis (RM19). Rien
// n'est enregistré si change échoue.
func (s *Service) editDraft(id string, change func(m *Model3D) error) (*Model3D, error) {
	unlock := lockDraft(id)
	defer unlock()
	m, err := s.draftAsset(id)
	if err != nil {
		return nil, err
	}
	if err := change(m); err != nil {
		return nil, err
	}
	if err := s.saveDraft(m); err != nil {
		return nil, err
	}
	return m, nil
}

// checkLicense vérifie la licence proposée pour un asset : présente au
// catalogue et, pour un asset dérivé dont le parent est connu, compatible
// avec la licence de ce parent (RM03). Une licence vide n'est pas contrôlée.
func checkLicense(licenseID string, parent *Model3D) error {
	if licenseID == "" {
		return nil
	}
	if _, err := GetLicense(licenseID); err != nil {
		return err
	}
	if parent != nil && parent.LicenseID != "" {
		if check := CheckLicenseCompatibility(parent.LicenseID, licenseID); !check.Compatible {
			return fmt.Errorf("%w : %s", ErrLicenseIncompatible, check.Reason)
		}
	}
	return nil
}

// Get récupère un modèle par son ID sur le canal indiqué.
// channelID="" utilise le canal par défaut configuré dans l'adapter.
// Cherche d'abord dans DraftStore (asset pas encore soumis), puis sur la blockchain.
func (s *Service) Get(id, channelID string) (*Model3D, error) {
	m, _, err := s.getAsset(id, channelID)
	return m, err
}

// List retourne les brouillons locaux (DraftStore) suivis des assets de la blockchain.
//
// Mode dégradé : si la blockchain est configurée mais que ListModelRecords échoue (panne
// d'infrastructure), les brouillons locaux déjà récupérés sont quand même retournés — mais
// accompagnés de l'erreur d'origine (non nil), pour que l'appelant sache explicitement que le
// résultat est partiel et le signale à son tour plutôt que de le présenter comme une liste
// complète. IsDegradedListErr(err) permet de distinguer ce cas d'un échec total.
func (s *Service) List(channelID string) ([]*Model3D, error) {
	var drafts []*Model3D
	if s.draftStore != nil {
		drafts, _ = s.draftStore.ListDrafts(channelID)
	}
	if s.blockchain == nil {
		if len(drafts) > 0 {
			return s.filterHidden(drafts), nil
		}
		return nil, ErrBlockchainUnavailable
	}
	onChain, err := s.blockchain.ListModelRecords(channelID)
	if err != nil {
		return s.filterHidden(drafts), err
	}
	return s.filterHidden(mergeDraftsAndLedger(drafts, onChain)), nil
}

// mergeDraftsAndLedger réunit brouillons et assets du ledger sans doublon
// (R-14) : un brouillon resté sous l'identifiant d'un asset déjà enregistré
// (retrait du brouillon échoué après une soumission) est périmé — le ledger
// fait foi.
func mergeDraftsAndLedger(drafts, onChain []*Model3D) []*Model3D {
	onLedger := make(map[string]bool, len(onChain))
	for _, m := range onChain {
		onLedger[m.ID] = true
	}
	merged := make([]*Model3D, 0, len(drafts)+len(onChain))
	for _, d := range drafts {
		if !onLedger[d.ID] {
			merged = append(merged, d)
		}
	}
	return append(merged, onChain...)
}

// filterHidden retire de assets tout ID masqué localement (RM08) — un asset
// masqué reste résolvable via Get()/GetModule(), seules les listes l'excluent.
func (s *Service) filterHidden(assets []*Model3D) []*Model3D {
	if s.removedStore == nil || len(assets) == 0 {
		return assets
	}
	filtered := assets[:0]
	for _, m := range assets {
		hidden, _ := s.removedStore.IsHidden(m.ID)
		if !hidden {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// ── Connexions d'assemblage (requiert connStore) ────────────────────────────

// connectionsLock sérialise la vérification et l'enregistrement des liaisons :
// deux demandes simultanées ne peuvent pas engager la même interface.
const connectionsLock = "\x00connections"

// AddConnection crée une liaison directe entre deux assets. Une liaison
// identique existante est refusée (ErrConnectionExists) : le client la
// retrouve dans la liste des liaisons, le serveur ne la devine pas (RM50).
func (s *Service) AddConnection(from, to, label string) (*Connection, error) {
	if s.connStore == nil {
		return nil, fmt.Errorf("connection store non configuré")
	}
	unlock := lockDraft(connectionsLock)
	defer unlock()
	conns, err := s.connStore.ListConnections()
	if err != nil {
		return nil, err
	}
	for _, c := range conns {
		if c.FromIfaceID == "" && c.ToIfaceID == "" && c.From == from && c.To == to {
			return nil, fmt.Errorf("%w : %s (%s → %s)", ErrConnectionExists, c.ID, from, to)
		}
	}
	conn := &Connection{ID: generateID(), From: from, To: to, Label: label}
	return conn, s.connStore.SaveConnection(conn)
}

// AddAssemblyLink crée un lien d'assemblage précis entre deux interfaces physiques.
// fastenerAssetID est l'asset d'accroche (vis, câble, clip…) ; vide si connexion directe.
// Les assets source et cible sont déduits des interfaces — le domaine reste l'autorité.
func (s *Service) AddAssemblyLink(fromIfaceID, toIfaceID, label, fromInstanceID, toInstanceID, fastenerAssetID string) (*Connection, error) {
	if s.ifaceStore == nil {
		return nil, fmt.Errorf("interface store non configuré")
	}
	if s.connStore == nil {
		return nil, fmt.Errorf("connection store non configuré")
	}
	fromIface, err := s.ifaceStore.GetInterface(fromIfaceID)
	if err != nil {
		return nil, fmt.Errorf("interface source: %w", err)
	}
	toIface, err := s.ifaceStore.GetInterface(toIfaceID)
	if err != nil {
		return nil, fmt.Errorf("interface cible: %w", err)
	}

	// Validation de l'accroche : elle doit avoir au moins une interface
	// compatible avec chacun des deux endpoints.
	if fastenerAssetID != "" {
		if err := s.validateFastener(fastenerAssetID, fromIface, toIface); err != nil {
			return nil, fmt.Errorf("asset d'accroche invalide : %w", err)
		}
	}

	unlock := lockDraft(connectionsLock)
	defer unlock()
	if err := s.checkLinkable(fromIface, toIface, fromInstanceID, toInstanceID, fastenerAssetID != ""); err != nil {
		return nil, err
	}

	conn := &Connection{
		ID:              generateID(),
		From:            fromIface.AssetID,
		To:              toIface.AssetID,
		Label:           label,
		FromIfaceID:     fromIfaceID,
		ToIfaceID:       toIfaceID,
		FromInstanceID:  fromInstanceID,
		ToInstanceID:    toInstanceID,
		FastenerAssetID: fastenerAssetID,
	}
	return conn, s.connStore.SaveConnection(conn)
}

// linkEnd est une extrémité de liaison : une interface d'une instance donnée
// (instance vide pour une liaison entre assets).
type linkEnd struct{ iface, instance string }

// checkLinkable vérifie qu'une liaison entre deux interfaces est permise :
//   - une interface n'est pas reliée à elle-même (ErrSelfLink) ;
//   - sans asset d'accroche, les deux interfaces sont compatibles (RM10,
//     critères RM11) — avec une accroche, c'est elle qui doit l'être avec
//     chacune (validateFastener) ;
//   - aucune des deux extrémités n'est déjà engagée dans une liaison (RM09) :
//     l'usage unique porte sur l'interface d'une instance donnée, la même
//     interface d'un composant restant disponible dans ses autres instances.
func (s *Service) checkLinkable(from, to *AssetInterface, fromInstance, toInstance string, viaFastener bool) error {
	if from.ID == to.ID && fromInstance == toInstance {
		return fmt.Errorf("%w : %s", ErrSelfLink, from.ID)
	}
	if !viaFastener && !ifacesCompatible(from, to) {
		return fmt.Errorf("%w : %s/%s (%s) ↔ %s/%s (%s)", ErrIncompatibleInterfaces,
			from.Category, from.Type, from.Direction, to.Category, to.Type, to.Direction)
	}
	conns, err := s.connStore.ListConnections()
	if err != nil {
		return err
	}
	busy := make(map[linkEnd]string, 2*len(conns))
	for _, c := range conns {
		if c.FromIfaceID != "" {
			busy[linkEnd{c.FromIfaceID, c.FromInstanceID}] = c.ID
		}
		if c.ToIfaceID != "" {
			busy[linkEnd{c.ToIfaceID, c.ToInstanceID}] = c.ID
		}
	}
	for _, e := range []linkEnd{{from.ID, fromInstance}, {to.ID, toInstance}} {
		if id, ok := busy[e]; ok {
			return fmt.Errorf("%w : interface %s (liaison %s)", ErrInterfaceInUse, e.iface, id)
		}
	}
	return nil
}

// validateFastener vérifie qu'un asset d'accroche est physiquement compatible
// avec les deux interfaces qu'il doit relier.
func (s *Service) validateFastener(fastenerID string, fromIface, toIface *AssetInterface) error {
	if s.ifaceStore == nil {
		return nil
	}
	fasIfaces, err := s.ifaceStore.ListInterfacesForAsset(fastenerID)
	if err != nil {
		return fmt.Errorf("interfaces introuvables : %w", err)
	}
	if len(fasIfaces) == 0 {
		return fmt.Errorf("l'asset d'accroche n'a aucune interface définie")
	}
	fromOK, toOK := false, false
	for _, fi := range fasIfaces {
		if ifacesCompatible(fi, fromIface) {
			fromOK = true
		}
		if ifacesCompatible(fi, toIface) {
			toOK = true
		}
	}
	if !fromOK {
		return fmt.Errorf("aucune interface de l'accroche n'est compatible avec l'interface source")
	}
	if !toOK {
		return fmt.Errorf("aucune interface de l'accroche n'est compatible avec l'interface cible")
	}
	return nil
}

// ifacesCompatible retourne true si deux interfaces physiques peuvent être reliées.
// Règles : même catégorie + même type ; sens cohérent (out↔in ou bidir) ; plages de valeurs compatibles.
func ifacesCompatible(a, b *AssetInterface) bool {
	if a.Category != b.Category || a.Type != b.Type {
		return false
	}
	dirOK := (a.Direction == IfaceOut && b.Direction == IfaceIn) ||
		(a.Direction == IfaceIn && b.Direction == IfaceOut) ||
		a.Direction == IfaceBidi || b.Direction == IfaceBidi
	if !dirOK {
		return false
	}
	// Choisir qui est "out" pour le contrôle de plage
	if a.Direction == IfaceOut || a.Direction == IfaceBidi {
		return ifaceValuesOverlap(a, b)
	}
	return ifaceValuesOverlap(b, a)
}

// ifaceValuesOverlap vérifie que les valeurs de deux interfaces se chevauchent
// (RM11) — out → fournisseur, in → consommateur. Une interface sans unité
// n'impose aucune contrainte de valeur ; avec une unité, une valeur simple vaut
// l'intervalle [v ; v] et 0 est une vraie valeur.
func ifaceValuesOverlap(out, in *AssetInterface) bool {
	if out.Unit == "" || in.Unit == "" {
		return true
	}
	outMin, outMax := valueRange(out)
	inMin, inMax := valueRange(in)
	return outMin <= inMax && inMin <= outMax
}

// valueRange retourne l'intervalle de valeurs d'une interface : sa plage, ou
// [v ; v] pour une valeur simple.
func valueRange(i *AssetInterface) (float64, float64) {
	if i.IsRange {
		return i.ValueMin, i.ValueMax
	}
	return i.ValueMin, i.ValueMin
}

func (s *Service) RemoveConnection(id string) error {
	if s.connStore == nil {
		return fmt.Errorf("connection store non configuré")
	}
	return s.connStore.RemoveConnection(id)
}

func (s *Service) ListConnections() ([]*Connection, error) {
	if s.connStore == nil {
		return nil, nil
	}
	return s.connStore.ListConnections()
}

func (s *Service) GetChildren(parentID string) ([]*Model3D, error) {
	if s.blockchain == nil {
		return nil, ErrBlockchainUnavailable
	}
	all, err := s.blockchain.ListModelRecords("")
	if err != nil {
		return nil, err
	}
	var children []*Model3D
	for _, m := range all {
		if m.ParentID == parentID {
			children = append(children, m)
		}
	}
	return children, nil
}

// ── Interfaces physiques (requiert ifaceStore) ───────────────────────────────

// requireDraftCarrier refuse de modifier les interfaces d'un asset soumis :
// elles ont été enregistrées sur le ledger avec lui (règle 27, RM19). Un asset
// que ce nœud ne connaît pas n'est pas contrôlé.
func (s *Service) requireDraftCarrier(assetID string) error {
	m, isDraft, err := s.getAsset(assetID, "")
	if err != nil {
		return nil
	}
	if !isDraft {
		return fmt.Errorf("%w : %s", ErrAssetSubmitted, assetID)
	}
	return s.checkOwner(m)
}

// AddInterface ajoute une interface à un asset en brouillon, qui doit exister
// (ErrNotFound). Son identifiant est toujours attribué ici : un identifiant
// fourni par le client remplacerait l'interface existante qui le porte.
func (s *Service) AddInterface(iface *AssetInterface) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	if _, err := s.draftAsset(iface.AssetID); err != nil {
		return err
	}
	iface.ID = generateID()
	return s.ifaceStore.SaveInterface(iface)
}

// UpdateInterface met à jour une interface d'un asset en brouillon et marque
// les connexions qui l'utilisent comme incompatibles si elles ne satisfont
// plus les contraintes. L'asset porteur ne change pas
// (ErrInterfaceCarrierChange).
func (s *Service) UpdateInterface(iface *AssetInterface) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	existing, err := s.ifaceStore.GetInterface(iface.ID)
	if err != nil {
		return err
	}
	if iface.AssetID != existing.AssetID {
		return fmt.Errorf("%w : %s appartient à %s", ErrInterfaceCarrierChange, iface.ID, existing.AssetID)
	}
	if err := s.requireDraftCarrier(existing.AssetID); err != nil {
		return err
	}
	if err := s.ifaceStore.SaveInterface(iface); err != nil {
		return err
	}
	if s.connStore == nil {
		return nil
	}
	conns, err := s.connStore.ListConnections()
	if err != nil {
		return nil // non bloquant
	}
	for _, c := range conns {
		if c.FromIfaceID != iface.ID && c.ToIfaceID != iface.ID {
			continue
		}
		fromIface, err := s.ifaceStore.GetInterface(c.FromIfaceID)
		if err != nil {
			continue
		}
		toIface, err := s.ifaceStore.GetInterface(c.ToIfaceID)
		if err != nil {
			continue
		}
		wasIncompatible := c.Incompatible
		c.Incompatible = !ifacesCompatible(fromIface, toIface)
		if c.Incompatible != wasIncompatible {
			_ = s.connStore.UpdateConnection(c)
		}
	}
	return nil
}

// RemoveInterface supprime une interface d'un asset en brouillon.
func (s *Service) RemoveInterface(id string) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	if err := s.requireDraftCarrierOf(id); err != nil {
		return err
	}
	return s.ifaceStore.RemoveInterface(id)
}

// requireDraftCarrierOf applique requireDraftCarrier à l'asset qui porte
// l'interface ifaceID, si elle existe déjà.
func (s *Service) requireDraftCarrierOf(ifaceID string) error {
	if existing, err := s.ifaceStore.GetInterface(ifaceID); err == nil {
		return s.requireDraftCarrier(existing.AssetID)
	}
	return nil
}

func (s *Service) ListInterfacesForAsset(assetID string) ([]*AssetInterface, error) {
	if s.ifaceStore == nil {
		return nil, nil
	}
	return s.ifaceStore.ListInterfacesForAsset(assetID)
}

func (s *Service) GetInterface(id string) (*AssetInterface, error) {
	if s.ifaceStore == nil {
		return nil, fmt.Errorf("interface store non configuré")
	}
	return s.ifaceStore.GetInterface(id)
}

func (s *Service) GetRefs() (*InterfaceRefs, error) {
	if s.ifaceStore == nil {
		return DefaultInterfaceRefs(), nil
	}
	return s.ifaceStore.GetRefs()
}

func (s *Service) AddRefCategory(cat string) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	return s.ifaceStore.AddRefCategory(cat)
}

func (s *Service) AddRefType(cat, typeName string) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	return s.ifaceStore.AddRefType(cat, typeName)
}

func (s *Service) AddRefUnit(cat, unit string) error {
	if s.ifaceStore == nil {
		return fmt.Errorf("interface store non configuré")
	}
	return s.ifaceStore.AddRefUnit(cat, unit)
}

// ── Suppression ─────────────────────────────────────────────────────────────

// removeConnectionsFor supprime en cascade les connexions impliquant un asset (règle 15).
func (s *Service) removeConnectionsFor(id string) {
	if s.connStore == nil {
		return
	}
	conns, _ := s.connStore.ListConnections()
	for _, c := range conns {
		if c.From == id || c.To == id {
			_ = s.connStore.RemoveConnection(c.ID)
		}
	}
}

// Remove supprime un asset (composant ou module). Un brouillon est réellement
// retiré du stockage local. Un asset déjà soumis à la blockchain ne peut
// jamais en être retiré (RM06/RM08, Fabric ne supporte aucune suppression) :
// il est seulement masqué localement (RemovedAssetStore.HideAsset) — il
// disparaît de List()/ListModules() mais reste résolvable par Get()/GetModule().
func (s *Service) Remove(id string) error {
	if s.draftStore != nil {
		d, err := s.draftStore.GetDraft(id)
		if err == nil {
			if err := s.checkOwner(d); err != nil {
				return err
			}
			s.removeConnectionsFor(id)
			return s.draftStore.RemoveDraft(id)
		}
		if !errors.Is(err, ErrNotFound) {
			return err
		}
	}
	if s.removedStore == nil {
		return fmt.Errorf("suppression non supportée : registre de masquage non configuré")
	}
	if s.blockchain == nil {
		return ErrBlockchainUnavailable
	}
	if _, err := s.blockchain.GetModelRecord(id, ""); err != nil {
		return fmt.Errorf("asset introuvable: %w", err)
	}
	s.removeConnectionsFor(id)
	return s.removedStore.HideAsset(id)
}

// applyAssetPatch applique sur m les champs non-vides d'un UpdateRequest.
func applyAssetPatch(m *Model3D, req UpdateRequest) {
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Description != "" {
		m.Description = req.Description
	}
	if req.LicenseID != "" {
		m.LicenseID = req.LicenseID
	}
	if req.Tags != nil {
		m.Tags = req.Tags
	}
	if req.Links != nil {
		m.Links = req.Links
	}
}

// UpdateAsset applique un patch partiel. Sur un brouillon, tout champ du patch
// change ; une nouvelle licence est vérifiée comme à la création (catalogue,
// compatibilité avec le parent — RM03). Sur un asset soumis, seules ses
// métadonnées descriptives changent, par une opération dédiée du ledger que
// seul son propriétaire peut signer ; sa licence reste figée
// (ErrAssetSubmitted, RM19, ADR-12).
func (s *Service) UpdateAsset(req UpdateRequest) (*Model3D, error) {
	m, isDraft, err := s.getAsset(req.ID, "")
	if err != nil {
		return nil, fmt.Errorf("asset introuvable: %w", err)
	}
	if !isDraft {
		return s.updateSubmittedMetadata(m, req)
	}
	return s.editDraft(req.ID, func(m *Model3D) error {
		if req.LicenseID != "" {
			var parent *Model3D
			if m.ParentID != "" {
				parent, _, _ = s.getAsset(m.ParentID, m.ChannelID)
			}
			if err := checkLicense(req.LicenseID, parent); err != nil {
				return err
			}
		}
		applyAssetPatch(m, req)
		return nil
	})
}

// updateSubmittedMetadata remplace sur le ledger les métadonnées descriptives
// d'un asset soumis (nom, description, tags, liens) ; tout autre changement
// exige un asset dérivé.
func (s *Service) updateSubmittedMetadata(m *Model3D, req UpdateRequest) (*Model3D, error) {
	if err := s.checkOwner(m); err != nil {
		return nil, err
	}
	if req.LicenseID != "" && req.LicenseID != m.LicenseID {
		return nil, fmt.Errorf("%w : la licence de %s ne change plus après soumission", ErrAssetSubmitted, m.ID)
	}
	applyAssetPatch(m, req)
	meta := AssetMetadata{Name: m.Name, Description: m.Description, Tags: m.Tags, Links: m.Links}
	if err := s.blockchain.UpdateModelMetadata(m.ID, m.ChannelID, meta); err != nil {
		return nil, fmt.Errorf("mise à jour des métadonnées sur le ledger : %w", err)
	}
	return m, nil
}

// ── Modules ───────────────────────────────────────────────────────────────────

// CreateModule crée un module en brouillon (RM16) : aucune transaction
// blockchain avant SubmitModule. Si req.ParentID référence un module existant
// (RM03 vérifiée si une licence est proposée), sa composition — instances et
// liaisons internes — est dupliquée dans le nouveau brouillon, pour reprendre
// un travail déjà soumis sans le reconstruire manuellement (voir
// copyModuleComposition). Le module parent, déjà soumis, n'est jamais modifié.
func (s *Service) CreateModule(req ModuleRequest) (*Model3D, error) {
	if s.draftStore == nil {
		return nil, fmt.Errorf("stockage de brouillons non configuré")
	}
	if req.ChannelID == "" {
		return nil, ErrChannelRequired
	}
	owner, err := s.creator(req.OwnerID)
	if err != nil {
		return nil, err
	}
	req.OwnerID = owner

	var parent *Model3D
	if req.ParentID != "" {
		p, _, err := s.getAsset(req.ParentID, req.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("module parent introuvable: %w", err)
		}
		parent = p
	}
	if err := checkLicense(req.LicenseID, parent); err != nil {
		return nil, err
	}

	m := &Model3D{
		ID:                 generateID(),
		Name:               req.Name,
		Description:        req.Description,
		OwnerID:            req.OwnerID,
		ParentID:           req.ParentID,
		ChannelID:          req.ChannelID,
		LicenseID:          req.LicenseID,
		Status:             ModuleDraft,
		CreatedAt:          time.Now(),
		Assemblies:         []string{},
		WorkspaceInstances: []WorkspaceInstance{},
		ModuleVersions:     []ModuleVersion{},
	}

	if parent != nil && parent.IsModule() {
		if err := s.copyModuleComposition(parent, m); err != nil {
			return nil, fmt.Errorf("copie de la composition du module parent: %w", err)
		}
	}

	if err := s.draftStore.SaveDraft(m); err != nil {
		return nil, fmt.Errorf("saving draft: %w", err)
	}
	return m, nil
}

// copyModuleComposition clone les WorkspaceInstances et les Liaisons internes
// (Assemblies) de parent dans dst — nouveaux identifiants d'instance et de
// connexion pour ne jamais entrer en collision avec ceux du parent. Les
// interfaces référencées (FromIfaceID/ToIfaceID) restent celles de l'asset
// sous-jacent (AssetInterface est indexée par AssetID, pas par instance) ;
// seuls FromInstanceID/ToInstanceID sont remappés vers les nouvelles instances.
func (s *Service) copyModuleComposition(parent, dst *Model3D) error {
	instanceIDMap := make(map[string]string, len(parent.WorkspaceInstances))
	for _, inst := range parent.WorkspaceInstances {
		newID := generateID()
		instanceIDMap[inst.ID] = newID
		dst.WorkspaceInstances = append(dst.WorkspaceInstances, WorkspaceInstance{
			ID: newID, AssetID: inst.AssetID,
		})
	}
	if len(parent.Assemblies) == 0 {
		return nil
	}
	if s.connStore == nil {
		return fmt.Errorf("connection store non configuré")
	}
	parentConnIDs := make(map[string]bool, len(parent.Assemblies))
	for _, id := range parent.Assemblies {
		parentConnIDs[id] = true
	}
	allConns, err := s.connStore.ListConnections()
	if err != nil {
		return err
	}
	for _, c := range allConns {
		if !parentConnIDs[c.ID] {
			continue
		}
		clone := &Connection{
			ID:              generateID(),
			From:            c.From,
			To:              c.To,
			Label:           c.Label,
			FromIfaceID:     c.FromIfaceID,
			ToIfaceID:       c.ToIfaceID,
			FromInstanceID:  remapInstanceID(instanceIDMap, c.FromInstanceID),
			ToInstanceID:    remapInstanceID(instanceIDMap, c.ToInstanceID),
			FastenerAssetID: c.FastenerAssetID,
		}
		if err := s.connStore.SaveConnection(clone); err != nil {
			return err
		}
		dst.Assemblies = append(dst.Assemblies, clone.ID)
	}
	return nil
}

// remapInstanceID retourne l'ID d'instance clonée correspondant à instanceID,
// ou instanceID inchangé s'il ne référence aucune instance clonée (connexion
// directe entre assets, sans WorkspaceInstance).
func remapInstanceID(instanceIDMap map[string]string, instanceID string) string {
	if instanceID == "" {
		return ""
	}
	if remapped, ok := instanceIDMap[instanceID]; ok {
		return remapped
	}
	return instanceID
}

func (s *Service) GetModule(id string) (*Model3D, error) {
	return s.Get(id, "")
}

// ListModules filtre List() aux seuls modules. Même mode dégradé que List() : si err n'est
// pas nil mais que all contient des brouillons locaux, ils sont filtrés et retournés avec
// l'erreur d'origine — voir IsDegradedListErr.
func (s *Service) ListModules(channelID string) ([]*Model3D, error) {
	all, err := s.List(channelID)
	var modules []*Model3D
	for _, m := range all {
		if m.IsModule() {
			modules = append(modules, m)
		}
	}
	return modules, err
}

// AddAssemblyToModule rattache au module brouillon une liaison existante qui
// relie deux de ses instances (RM51). Une liaison inconnue ou étrangère au
// module est refusée (ErrConnectionNotInModule) : un composant sans instance ne
// devient donc jamais module par erreur.
func (s *Service) AddAssemblyToModule(moduleID, connID string) error {
	_, err := s.editDraft(moduleID, func(m *Model3D) error {
		for _, id := range m.Assemblies {
			if id == connID {
				return nil // déjà présent
			}
		}
		conn, err := s.findConnection(connID)
		if err != nil {
			return err
		}
		if !connectsInstancesOf(m, conn) {
			return fmt.Errorf("%w : %s ne relie pas deux instances de %s", ErrConnectionNotInModule, connID, moduleID)
		}
		m.Assemblies = append(m.Assemblies, connID)
		return nil
	})
	return err
}

// findConnection retourne la liaison connID, ErrConnectionNotInModule si elle
// n'existe pas.
func (s *Service) findConnection(connID string) (*Connection, error) {
	if s.connStore == nil {
		return nil, fmt.Errorf("connection store non configuré")
	}
	conns, err := s.connStore.ListConnections()
	if err != nil {
		return nil, err
	}
	for _, c := range conns {
		if c.ID == connID {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%w : liaison %s introuvable", ErrConnectionNotInModule, connID)
}

// connectsInstancesOf indique si la liaison c relie deux instances du module
// m : par leurs identifiants d'instance, ou, pour une liaison entre assets sans
// instance, par les assets de ses instances.
func connectsInstancesOf(m *Model3D, c *Connection) bool {
	instances := make(map[string]bool, len(m.WorkspaceInstances))
	assets := make(map[string]bool, len(m.WorkspaceInstances))
	for _, inst := range m.WorkspaceInstances {
		instances[inst.ID] = true
		assets[inst.AssetID] = true
	}
	if c.FromInstanceID != "" || c.ToInstanceID != "" {
		return instances[c.FromInstanceID] && instances[c.ToInstanceID]
	}
	return assets[c.From] && assets[c.To]
}

func (s *Service) RemoveAssemblyFromModule(moduleID, connID string) error {
	_, err := s.editDraft(moduleID, func(m *Model3D) error {
		filtered := m.Assemblies[:0]
		for _, id := range m.Assemblies {
			if id != connID {
				filtered = append(filtered, id)
			}
		}
		m.Assemblies = filtered
		return nil
	})
	return err
}

// SubmitModule soumet un module brouillon à la blockchain et crée sa version
// immuable (RM18). Un module déjà soumis n'est jamais re-soumis
// (ErrAssetSubmitted) : une nouvelle version est un nouveau module dérivé
// (AddFull avec ParentID, qui reprend sa composition — RM19, ADR-12).
func (s *Service) SubmitModule(moduleID, note string) (*Model3D, error) {
	if s.blockchain == nil {
		return nil, ErrBlockchainUnavailable
	}
	unlock := lockDraft(moduleID)
	defer unlock()
	m, err := s.submittableDraft(moduleID)
	if err != nil {
		return nil, err
	}
	if len(m.Assemblies) == 0 {
		return nil, fmt.Errorf("le module ne contient aucun assemblage")
	}

	hash := computeModuleHash(m.ID, len(m.ModuleVersions)+1, m.Assemblies)
	blockID := simModuleBlockID()

	v := ModuleVersion{
		Number:     len(m.ModuleVersions) + 1,
		Assemblies: append([]string{}, m.Assemblies...),
		Hash:       hash,
		Note:       note,
		CreatedAt:  time.Now(),
		BlockID:    blockID,
	}
	m.ModuleVersions = append(m.ModuleVersions, v)
	m.Status = ModuleSubmitted
	if err := s.blockchain.StoreModelRecord(m); err != nil {
		return nil, fmt.Errorf("storing on blockchain: %w", err)
	}
	_ = s.draftStore.RemoveDraft(moduleID)
	return m, nil
}

// AddAssetToWorkspace ajoute une nouvelle instance de l'asset au module.
// L'asset ajouté doit exister (ErrNotFound) ; s'il est le module ou le
// contient, même indirectement, il est refusé (ErrCyclicComposition, RM52).
func (s *Service) AddAssetToWorkspace(moduleID, assetID string) (*Model3D, error) {
	m, err := s.editDraft(moduleID, func(m *Model3D) error {
		added, _, err := s.getAsset(assetID, "")
		if err != nil {
			return err
		}
		if err := s.checkNotContaining(added, moduleID); err != nil {
			return err
		}
		m.WorkspaceInstances = append(m.WorkspaceInstances, WorkspaceInstance{ID: generateID(), AssetID: assetID})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

// checkNotContaining refuse d'ajouter l'asset added au module moduleID s'il
// est ce module ou le contient, même indirectement : une composition cyclique
// n'aurait pas de fin (GetModuleInterfaces la parcourt récursivement).
func (s *Service) checkNotContaining(added *Model3D, moduleID string) error {
	seen := map[string]bool{}
	var contains func(a *Model3D) bool
	contains = func(a *Model3D) bool {
		if a.ID == moduleID {
			return true
		}
		if seen[a.ID] {
			return false
		}
		seen[a.ID] = true
		for _, inst := range a.WorkspaceInstances {
			if sub, _, err := s.getAsset(inst.AssetID, ""); err == nil && contains(sub) {
				return true
			}
		}
		return false
	}
	if contains(added) {
		return fmt.Errorf("%w : %s contient %s", ErrCyclicComposition, added.ID, moduleID)
	}
	return nil
}

// RemoveAssetFromWorkspace retire une instance spécifique (par instanceID) du module.
// Supprime en cascade les connexions de cette instance dans le module.
func (s *Service) RemoveAssetFromWorkspace(moduleID, instanceID string) (*Model3D, error) {
	return s.editDraft(moduleID, func(m *Model3D) error {
		s.removeInstance(m, instanceID)
		return nil
	})
}

// removeInstance retire l'instance instanceID du module m et supprime en
// cascade ses liaisons internes (règle 15).
func (s *Service) removeInstance(m *Model3D, instanceID string) {
	if s.connStore != nil {
		asmSet := make(map[string]bool, len(m.Assemblies))
		for _, id := range m.Assemblies {
			asmSet[id] = true
		}
		conns, _ := s.connStore.ListConnections()
		survivingAsm := append([]string{}, m.Assemblies...)
		for _, c := range conns {
			if !asmSet[c.ID] {
				continue
			}
			if c.FromInstanceID == instanceID || c.ToInstanceID == instanceID {
				next := survivingAsm[:0]
				for _, id := range survivingAsm {
					if id != c.ID {
						next = append(next, id)
					}
				}
				survivingAsm = next
				_ = s.connStore.RemoveConnection(c.ID)
			}
		}
		m.Assemblies = survivingAsm
	}

	filtered := m.WorkspaceInstances[:0]
	for _, inst := range m.WorkspaceInstances {
		if inst.ID != instanceID {
			filtered = append(filtered, inst)
		}
	}
	m.WorkspaceInstances = filtered
}

// RemoveModule supprime un module — délègue à Remove, qui gère brouillon local,
// suppression en cascade des connexions et suppression blockchain identiquement.
func (s *Service) RemoveModule(id string) error {
	return s.Remove(id)
}

// GetModuleInterfaces retourne les interfaces d'un composant ou d'un module :
// - Composant simple : toutes ses interfaces directes.
// - Module : interfaces exposées = interfaces des sous-composants/modules non reliées en interne.
// La fonction est récursive ; le cache evite les appels blockchain redondants.
func (s *Service) GetModuleInterfaces(moduleID string) ([]*AssetInterface, error) {
	if s.ifaceStore == nil {
		return nil, fmt.Errorf("interface store non configuré")
	}
	cache := make(map[string][]*AssetInterface)
	return s.getModuleInterfacesInto(moduleID, cache)
}

// getModuleInterfacesInto calcule les interfaces de moduleID dans cache. Un
// module en cours de calcul y figure déjà, vide : une composition cyclique
// enregistrée avant RM52 ne fait donc pas récurser sans fin, le cycle
// n'expose simplement aucune interface.
func (s *Service) getModuleInterfacesInto(moduleID string, cache map[string][]*AssetInterface) ([]*AssetInterface, error) {
	if ifaces, ok := cache[moduleID]; ok {
		return ifaces, nil
	}
	mod, _, err := s.getAsset(moduleID, "")
	if err != nil {
		// L'enregistrement est introuvable (brouillon comme blockchain absents,
		// ou Fabric indisponible). On suppose un composant simple : on retourne
		// les interfaces locales directes.
		ifaces, err2 := s.ifaceStore.ListInterfacesForAsset(moduleID)
		if err2 != nil {
			return nil, err // retourner l'erreur blockchain originale
		}
		cache[moduleID] = ifaces
		return ifaces, nil
	}
	// Composant simple : retourner ses interfaces directes.
	if !mod.IsModule() {
		ifaces, err := s.ifaceStore.ListInterfacesForAsset(moduleID)
		if err != nil {
			return nil, err
		}
		cache[moduleID] = ifaces
		return ifaces, err
	}
	cache[moduleID] = nil // en cours de calcul
	// Calculer les interfaces exposées (non connectées en interne).
	// Une connexion porte l'instance de chaque extrémité (FromInstanceID/ToInstanceID) :
	// quand elle est renseignée, seule cette instance précise est marquée interne, ce qui
	// permet à deux instances d'un même asset d'exposer chacune leur propre interface non
	// connectée. Une connexion créée sans contexte d'instance (endpoint hors workspace)
	// marque l'interface interne pour toutes les instances qui la partagent.
	internalGlobal := map[string]bool{}
	internalByInstance := map[string]map[string]bool{}
	markInternal := func(ifaceID, instanceID string) {
		if ifaceID == "" {
			return
		}
		if instanceID == "" {
			internalGlobal[ifaceID] = true
			return
		}
		if internalByInstance[instanceID] == nil {
			internalByInstance[instanceID] = map[string]bool{}
		}
		internalByInstance[instanceID][ifaceID] = true
	}
	isInternal := func(ifaceID, instanceID string) bool {
		return internalGlobal[ifaceID] || internalByInstance[instanceID][ifaceID]
	}
	if s.connStore != nil {
		asmSet := map[string]bool{}
		for _, id := range mod.Assemblies {
			asmSet[id] = true
		}
		conns, _ := s.connStore.ListConnections()
		for _, c := range conns {
			if asmSet[c.ID] {
				markInternal(c.FromIfaceID, c.FromInstanceID)
				markInternal(c.ToIfaceID, c.ToInstanceID)
			}
		}
	}
	var exposed []*AssetInterface
	for _, inst := range mod.WorkspaceInstances {
		ifaces, err := s.getModuleInterfacesInto(inst.AssetID, cache)
		if err != nil {
			continue
		}
		for _, iface := range ifaces {
			if !isInternal(iface.ID, inst.ID) {
				exposed = append(exposed, iface)
			}
		}
	}
	// Inclure les interfaces directes du module lui-même
	ownIfaces, _ := s.ifaceStore.ListInterfacesForAsset(moduleID)
	for _, iface := range ownIfaces {
		if !isInternal(iface.ID, "") {
			exposed = append(exposed, iface)
		}
	}
	cache[moduleID] = exposed
	return exposed, nil
}

func computeModuleHash(productID string, vNum int, assemblies []string) string {
	sorted := append([]string{}, assemblies...)
	sort.Strings(sorted)
	payload := fmt.Sprintf("%s|v%d|%s", productID, vNum, strings.Join(sorted, ","))
	h := sha256.Sum256([]byte(payload))
	return fmt.Sprintf("sha256:%x", h[:8])
}

func simModuleBlockID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("P%08X", b)
}

// ── Licences ──────────────────────────────────────────────────────────────────

func (s *Service) ListLicenses() []*License {
	return ListLicenses()
}

func (s *Service) GetLicense(id string) (*License, error) {
	return GetLicense(id)
}

// CheckLicenseCompatibility expose la règle domaine via le service.
func (s *Service) CheckLicenseCompatibility(parentLicenseID, proposedLicenseID string) *LicenseCheck {
	return CheckLicenseCompatibility(parentLicenseID, proposedLicenseID)
}

// CheckModuleLicenseCompatibility vérifie la compatibilité d'une licence de produit
// vis-à-vis des licences de tous ses composants.
func (s *Service) CheckModuleLicenseCompatibility(componentLicenseIDs []string, proposedProductLicenseID string) *LicenseCheck {
	return CheckModuleLicenseCompatibility(componentLicenseIDs, proposedProductLicenseID)
}

// ── interne ──────────────────────────────────────────────────────────────────

func (s *Service) hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}
