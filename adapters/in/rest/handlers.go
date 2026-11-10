// adapters/in/rest/handlers.go — handlers REST (adaptateur entrant HTTP)
package rest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"myr-core/domain/identity"
	"myr-core/domain/model"
	"myr-core/domain/network"
	rbac "myr-core/domain/role"
)

// validFieldID valide les identifiants libres (owner_id, channel_id, parent_id…).
var validFieldID = regexp.MustCompile(`^[a-zA-Z0-9_\-\.@]{1,128}$`)

// Handler traduit les requêtes HTTP en appels aux ports d'entrée du domaine.
// NetworkInfo regroupe les métadonnées du réseau Fabric transmises au frontend.
type NetworkInfo struct {
	Channel  string   // canal par défaut (depuis fabric.env)
	Channels []string // tous les canaux disponibles sur ce réseau
	Peer     string
	Network  string
}

// ipRateLimiter limite le nombre de tentatives par IP sur les routes sensibles.
// Fenêtre glissante : maxAttempts par window.
type ipRateLimiter struct {
	mu          sync.Mutex
	attempts    map[string][]time.Time
	maxAttempts int
	window      time.Duration
}

func newIPRateLimiter(max int, window time.Duration) *ipRateLimiter {
	l := &ipRateLimiter{attempts: make(map[string][]time.Time), maxAttempts: max, window: window}
	go l.cleanup()
	return l
}

func (l *ipRateLimiter) allow(ip string) bool {
	now := time.Now()
	cutoff := now.Add(-l.window)
	l.mu.Lock()
	defer l.mu.Unlock()
	times := l.attempts[ip]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	valid = append(valid, now)
	l.attempts[ip] = valid
	return len(valid) <= l.maxAttempts
}

func (l *ipRateLimiter) cleanup() {
	for range time.Tick(5 * time.Minute) {
		cutoff := time.Now().Add(-l.window)
		l.mu.Lock()
		for ip, times := range l.attempts {
			valid := times[:0]
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(l.attempts, ip)
			} else {
				l.attempts[ip] = valid
			}
		}
		l.mu.Unlock()
	}
}

// blockchainRouter sélectionne le bon BlockchainPort pour un réseau donné.
// Implémenté par fabric.NetworkPool — défini ici pour éviter l'import circulaire
// entre adapters/in et adapters/out.
type blockchainRouter interface {
	// BlockchainFor retourne un accès au réseau en lecture seule.
	BlockchainFor(networkID string) model.BlockchainPort
	// BlockchainForSigner retourne un accès au réseau dont les écritures sont
	// signées par l'identité (MSP, certificat, clé) de l'utilisateur (ADR-12).
	BlockchainForSigner(networkID, mspID, certPath, keyPath string) model.BlockchainPort
}

type Handler struct {
	modelSvc     model.ModelService
	thumbStore   model.ThumbnailStore
	ifaceStore   model.InterfaceStore
	draftStore   model.DraftStore
	removedStore model.RemovedAssetStore
	identitySvc  identity.IdentityService
	networkSvc   network.NetworkService
	roleSvc      rbac.RoleService
	sessions     sessionBackend
	netInfo      NetworkInfo
	authLimiter  *ipRateLimiter

	trustedProxies []netip.Prefix // proxys dont X-Forwarded-For est cru (limits.go)
	maxUpload      int64          // taille maximale d'un envoi de fichier (octets)
	assetCounts    assetCountCache

	// Phase 2 : multi-réseau
	bcRouter         blockchainRouter      // nil en mode legacy (réseau unique)
	connStore        model.ConnectionStore // pour svcFor
	ogImageFetcher   model.OGImageFetcher  // pour svcFor
	defaultNetworkID string                // réseau actif par défaut au démarrage

	fabricConnected bool   // vrai si le gateway Fabric a démarré correctement
	buildDate       string // injecté via ldflags à la compilation
	version         string // injecté via ldflags à la compilation
	debugErrors     bool   // vrai : le détail des erreurs internes est renvoyé au client (serveurs de test)
}

func NewHandler(svc model.ModelService, ts model.ThumbnailStore, is model.InterfaceStore) *Handler {
	return &Handler{
		modelSvc:    svc,
		thumbStore:  ts,
		ifaceStore:  is,
		sessions:    newSessionStore(),
		authLimiter: newIPRateLimiter(10, time.Minute), // 10 tentatives / minute / IP

		trustedProxies: loopbackProxies,
		maxUpload:      DefaultMaxUpload,
	}
}

// WithIdentityService attache le service d'identité au handler (requis en mode blockchain).
func (h *Handler) WithIdentityService(svc identity.IdentityService) *Handler {
	h.identitySvc = svc
	return h
}

// WithDebugErrors renvoie au client le détail des erreurs internes en plus de
// leur référence — réservé aux serveurs de test : le détail peut révéler l'état
// interne du serveur (identités existantes, chemins, erreurs d'infrastructure).
func (h *Handler) WithDebugErrors(on bool) *Handler {
	h.debugErrors = on
	return h
}

// WithNetworkService attache le service réseau au handler (requis en mode blockchain).
func (h *Handler) WithNetworkService(svc network.NetworkService) *Handler {
	h.networkSvc = svc
	return h
}

// WithRoleService attache le service de rôles RBAC au handler, consulté par
// requireRole. Sans cet appel, requireRole retombe sur la hiérarchie de rang
// historique (reader < contributor < admin) pour compatibilité.
func (h *Handler) WithRoleService(svc rbac.RoleService) *Handler {
	h.roleSvc = svc
	return h
}

// WithNetworkInfo injecte les infos du réseau Fabric dans le handler (exposées via /api/status).
func (h *Handler) WithNetworkInfo(ni NetworkInfo) *Handler {
	h.netInfo = ni
	return h
}

// WithBlockchainRouter attache le routeur multi-réseau (fabric.NetworkPool).
func (h *Handler) WithBlockchainRouter(r blockchainRouter) *Handler {
	h.bcRouter = r
	return h
}

// WithConnStore attache le store de connexions pour svcFor.
func (h *Handler) WithConnStore(cs model.ConnectionStore) *Handler {
	h.connStore = cs
	return h
}

// WithDraftStore attache le store de brouillons pour svcFor (composants créés draft:true).
func (h *Handler) WithDraftStore(ds model.DraftStore) *Handler {
	h.draftStore = ds
	return h
}

// WithRemovedStore attache le registre de masquage local (RM08) pour svcFor.
func (h *Handler) WithRemovedStore(rs model.RemovedAssetStore) *Handler {
	h.removedStore = rs
	return h
}

// WithOGImageFetcher attache la source de régénération de miniature (og:image) pour svcFor.
func (h *Handler) WithOGImageFetcher(f model.OGImageFetcher) *Handler {
	h.ogImageFetcher = f
	return h
}

// WithDefaultNetwork définit le réseau Fabric actif par défaut (ID du NetworkProfile).
func (h *Handler) WithDefaultNetwork(networkID string) *Handler {
	h.defaultNetworkID = networkID
	return h
}

// WithFabricConnected indique si la gateway Fabric a démarré avec succès.
// Utilisé par /api/status pour éviter un appel blockchain bloquant à chaque heartbeat.
func (h *Handler) WithFabricConnected(connected bool) *Handler {
	h.fabricConnected = connected
	return h
}

// WithBuildDate transmet la date de compilation injectée via ldflags.
func (h *Handler) WithBuildDate(d string) *Handler {
	h.buildDate = d
	return h
}

// WithVersion transmet le numéro de version injecté via ldflags.
func (h *Handler) WithVersion(v string) *Handler {
	h.version = v
	return h
}

// WithSessionPersistence active la persistance des sessions sur disque.
// Sans effet si le backend est Redis (les sessions sont déjà persistantes côté Redis).
func (h *Handler) WithSessionPersistence(path string) *Handler {
	if ms, ok := h.sessions.(*sessionStore); ok {
		ms.withPersistence(path)
	}
	return h
}

// WithRedisSession remplace le store de sessions en mémoire par un backend Redis.
// url : "redis://[:password@]host[:port][/db]" ou "rediss://..." pour TLS.
// En cas d'erreur de connexion, le store en mémoire est conservé.
func (h *Handler) WithRedisSession(url string) *Handler {
	rs, err := newRedisSessionStore(url)
	if err != nil {
		return h
	}
	h.sessions = rs
	return h
}

// networkFor retourne l'ID du réseau actif pour cette requête.
// Priorité : session utilisateur → defaultNetworkID configuré au démarrage.
func (h *Handler) networkFor(r *http.Request) string {
	if sess := sessionFromCtx(r); sess != nil && sess.NetworkID != "" {
		return sess.NetworkID
	}
	return h.defaultNetworkID
}

// svcFor retourne un model.ModelService configuré pour le réseau actif de la requête.
// En mode legacy (bcRouter non configuré ou networkID vide), retourne le modelSvc pré-wired.
// En mode multi-réseau, crée une instance légère par requête avec le bon blockchain,
// signé par le wallet de l'utilisateur de la session (voir blockchainFor).
func (h *Handler) svcFor(r *http.Request) model.ModelService {
	// Le service agit toujours pour l'identité de la session : seul le
	// propriétaire d'un asset le modifie (FT-001). Sans session (visiteur),
	// aucune identité : lecture seule.
	actor := ""
	if sess := sessionFromCtx(r); sess != nil {
		actor = sess.Handle
	}
	if h.bcRouter == nil {
		return h.modelSvc.As(actor)
	}
	networkID := h.networkFor(r)
	if networkID == "" {
		// Pas de réseau configuré dans le pool — utilise le service wired au démarrage
		return h.modelSvc.As(actor)
	}
	bc := h.blockchainFor(r, networkID)
	svc := model.NewService(bc)
	if h.connStore != nil {
		svc = svc.WithConnStore(h.connStore)
	}
	if h.ogImageFetcher != nil {
		svc = svc.WithOGImageFetcher(h.ogImageFetcher)
	}
	return svc.WithThumbStore(h.thumbStore).WithIfaceStore(h.ifaceStore).WithDraftStore(h.draftStore).WithRemovedStore(h.removedStore).As(actor)
}

// blockchainFor retourne l'accès au réseau networkID pour la requête : signé
// par le wallet de l'identité de la session (ADR-12) s'il est disponible sur
// ce nœud, en lecture seule sinon (invité, session sans identité).
func (h *Handler) blockchainFor(r *http.Request, networkID string) model.BlockchainPort {
	if sess := sessionFromCtx(r); sess != nil && sess.Handle != "" && h.identitySvc != nil {
		w, err := h.identitySvc.Wallet(sess.Handle)
		if err == nil {
			return h.bcRouter.BlockchainForSigner(networkID, w.OrgID, w.CertPath(), w.KeyPath())
		}
		log.Printf("wallet de %s indisponible, accès au ledger en lecture seule : %v", sess.Handle, err)
	}
	return h.bcRouter.BlockchainFor(networkID)
}

// ownerFor retourne le propriétaire des assets créés par la requête :
// l'identité de la session (« pseudo@org »), celle dont le wallet signera leur
// soumission sur le ledger (ADR-12) — la valeur owner_id du formulaire n'est
// retenue que pour une session sans identité.
func ownerFor(r *http.Request, formOwner string) string {
	if sess := sessionFromCtx(r); sess != nil && sess.Handle != "" {
		return sess.Handle
	}
	return formOwner
}

// channelFor retourne le canal actif de la session, ou le canal par défaut du réseau.
func (h *Handler) channelFor(r *http.Request) string {
	if sess := sessionFromCtx(r); sess != nil && sess.Channel != "" {
		return sess.Channel
	}
	return h.netInfo.Channel
}

// requireAuth est un middleware qui vérifie l'authentification.
// Ordre de priorité : 1) JWT Bearer (nouveau), 2) X-Myr-Token (legacy).
func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Myr-Token")
		sess := h.sessions.get(token)
		if sess == nil {
			jsonError(w, "authentification requise", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxSession, sess)
		next(w, r.WithContext(ctx))
	}
}

// requireReader protège une route de lecture du catalogue (composants, modules,
// référentiels, licences) : une session suffit ; sans jeton, une lecture (GET,
// HEAD) est admise comme visiteur quand le réseau actif ouvre la lecture
// publique (AllowAutoGuest, RM49). Un jeton présent mais invalide reste refusé
// (401) : le client doit apprendre que sa session a expiré. Toute écriture
// passe par requireRole.
func (h *Handler) requireReader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Myr-Token") == "" && (r.Method == http.MethodGet || r.Method == http.MethodHead) && h.publicRead() {
			next(w, r.WithContext(context.WithValue(r.Context(), ctxVisitor, true)))
			return
		}
		h.requireAuth(next)(w, r)
	}
}

// publicRead indique si le réseau actif ouvre la lecture du catalogue aux
// visiteurs sans compte (réglage de l'administrateur, désactivé par défaut).
func (h *Handler) publicRead() bool {
	if h.networkSvc == nil {
		return false
	}
	p, err := h.networkSvc.GetActive()
	return err == nil && p != nil && p.AllowAutoGuest
}

// legacyRoleRank est le repli utilisé quand aucun rbac.RoleService n'est injecté
// (compatibilité rétroactive) : hiérarchie reader < contributor < admin.
var legacyRoleRank = map[string]int{"reader": 1, "contributor": 2, "admin": 3}
var legacyPermRank = map[rbac.Permission]int{
	rbac.PermRead: 1, rbac.PermWrite: 2,
	rbac.PermAdmin: 3, rbac.PermNetworkAdmin: 3, rbac.PermRoleAdmin: 3, rbac.PermIdentityAdmin: 3,
}

// hasPermission indique si le rôle porté par la session donne droit à la permission
// demandée. Consulte le rbac.RoleService injecté (WithRoleService) ; à défaut, retombe
// sur la hiérarchie de rang historique (reader < contributor < admin). Partagé entre
// requireRole (gate de route, échoue par 403) et les handlers qui doivent nuancer leur
// réponse selon le rôle plutôt que bloquer la requête entière (ex. handleIdentityWallets).
func (h *Handler) hasPermission(sess *myrSession, perm rbac.Permission) bool {
	if sess == nil {
		return false
	}
	if h.roleSvc != nil {
		return h.roleSvc.HasPermission(sess.Role, perm)
	}
	return legacyRoleRank[sess.Role] >= legacyPermRank[perm]
}

// requireRole est un middleware qui vérifie que la session porte la permission requise.
func (h *Handler) requireRole(perm rbac.Permission, next http.HandlerFunc) http.HandlerFunc {
	return h.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		sess := sessionFromCtx(r)
		if !h.hasPermission(sess, perm) {
			jsonError(w, "droits insuffisants", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

// ── DTOs ─────────────────────────────────────────────────────────────────────

// componentDTO est la ressource REST unique pour Model3D (ADR-11,
// specs/3-Conception/Conception_intro.md) — qu'il porte un fichier ressource
// (Hash), des instances (Instances), ou les deux à la fois. Instances,
// Assemblies et ModuleVersions restent vides ([]) tant que l'asset n'a jamais
// été décomposé ; le client dérive l'affichage « composant »/« produit » de
// len(Instances) > 0, jamais stocké séparément.
type componentDTO struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Category    model.Category `json:"category"`
	OwnerID     string         `json:"owner_id"`
	ParentID    string         `json:"parent_id,omitempty"`
	BlockID     string         `json:"block_id,omitempty"`
	Hash        string         `json:"hash,omitempty"`
	LicenseID   string         `json:"license_id,omitempty"`
	Tags        []string       `json:"tags"`
	Links       []string       `json:"links,omitempty"`
	Thumbnail   string         `json:"thumbnail,omitempty"`
	// Status vaut "draft" (créé via draft:true, pas encore engagé sur la blockchain
	// — voir POST /components/{id}/submit) ou "submitted" (comportement nominal,
	// immédiat). Vide pour un composant créé avant l'introduction du brouillon.
	Status         model.ModuleStatus        `json:"status,omitempty"`
	Assemblies     []string                  `json:"assemblies"`
	Instances      []model.WorkspaceInstance `json:"instances"`
	ModuleVersions []model.ModuleVersion     `json:"module_versions"`
	CreatedAt      time.Time                 `json:"created_at"`
}

type connectionDTO struct {
	ID              string `json:"id"`
	From            string `json:"from"`
	To              string `json:"to"`
	Label           string `json:"label,omitempty"`
	FromIfaceID     string `json:"from_iface_id,omitempty"`
	ToIfaceID       string `json:"to_iface_id,omitempty"`
	FromInstanceID  string `json:"from_instance_id,omitempty"`
	ToInstanceID    string `json:"to_instance_id,omitempty"`
	FastenerAssetID string `json:"fastener_asset_id,omitempty"`
	Incompatible    bool   `json:"incompatible,omitempty"`
}

func toConnectionDTO(c *model.Connection) connectionDTO {
	return connectionDTO{
		ID: c.ID, From: c.From, To: c.To, Label: c.Label,
		FromIfaceID: c.FromIfaceID, ToIfaceID: c.ToIfaceID,
		FromInstanceID: c.FromInstanceID, ToInstanceID: c.ToInstanceID,
		FastenerAssetID: c.FastenerAssetID, Incompatible: c.Incompatible,
	}
}

type graphResponse struct {
	Components  []componentDTO  `json:"components"`
	Connections []connectionDTO `json:"connections"`
	Total       int             `json:"total"`
	// Degraded et Warning ne sont présents que si la blockchain était injoignable au
	// moment de la requête : la réponse ne contient alors que les brouillons locaux
	// (non soumis), pas les assets on-chain — voir listGraph.
	Degraded bool   `json:"degraded,omitempty"`
	Warning  string `json:"warning,omitempty"`
}

type statusResponse struct {
	Mode      string `json:"mode"`
	Connected bool   `json:"connected"`
	Assets    int    `json:"assets"`
	Channel   string `json:"channel,omitempty"`
	Peer      string `json:"peer,omitempty"`
	Network   string `json:"network,omitempty"`
}

func toComponentDTO(m *model.Model3D, thumbnail string) componentDTO {
	tags := m.Tags
	if tags == nil {
		tags = []string{}
	}
	assemblies := m.Assemblies
	if assemblies == nil {
		assemblies = []string{}
	}
	instances := m.WorkspaceInstances
	if instances == nil {
		instances = []model.WorkspaceInstance{}
	}
	versions := m.ModuleVersions
	if versions == nil {
		versions = []model.ModuleVersion{}
	}
	return componentDTO{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		Category:       m.Category,
		OwnerID:        m.OwnerID,
		ParentID:       m.ParentID,
		BlockID:        m.BlockID,
		Hash:           m.Hash,
		LicenseID:      m.LicenseID,
		Tags:           tags,
		Links:          m.Links,
		Thumbnail:      thumbnail,
		Status:         m.Status,
		Assemblies:     assemblies,
		Instances:      instances,
		ModuleVersions: versions,
		CreatedAt:      m.CreatedAt,
	}
}

// ── /api/components ──────────────────────────────────────────────────────────────

func (h *Handler) handleComponents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listGraph(w, r)
	case http.MethodPost:
		h.createAsset(w, r)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// listGraph liste tous les composants, décomposés (anciens "modules") ou non,
// avec filtres serveur (ADR-11, specs/3-Conception/Conception_intro.md).
//
//	@Summary		Lister les composants
//	@Description	Filtre côté serveur — ne renvoie qu'un sous-ensemble pour éviter de tout charger. Inclut tout Model3D, décomposé (instances non vide) ou non.
//	@Tags			components
//	@Produce		json
//	@Param			q			query		string	false	"recherche texte (nom, description, tags)"
//	@Param			categories	query		string	false	"catégories séparées par virgule (ex: base,variation)"
//	@Param			owner_id	query		string	false	"filtre par propriétaire"
//	@Param			parent_id	query		string	false	"filtre par asset parent"
//	@Param			hash		query		string	false	"filtre par hash SHA-256 exact"
//	@Param			tags		query		string	false	"tags séparés par virgule (correspondance sur au moins un tag)"
//	@Param			status		query		string	false	"filtre par statut : draft ou submitted"
//	@Param			limit		query		int		false	"nombre maximum de résultats (défaut 200, max 1000)"
//	@Success		200	{object}	graphResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/components [get]
func (h *Handler) listGraph(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	cats := r.URL.Query().Get("categories") // "base,variation" ou vide = tous
	ownerID := strings.TrimSpace(r.URL.Query().Get("owner_id"))
	parentID := strings.TrimSpace(r.URL.Query().Get("parent_id"))
	hashFilter := strings.TrimSpace(r.URL.Query().Get("hash")) // hash SHA-256 exact
	tagsParam := r.URL.Query().Get("tags")                     // "tag1,tag2" — match si l'asset possède au moins un
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := 200
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 1000 {
		limit = l
	}
	catSet := map[string]bool{}
	if cats != "" {
		for _, c := range strings.Split(cats, ",") {
			if c = strings.TrimSpace(c); c != "" {
				catSet[c] = true
			}
		}
	}
	filterTags := []string{}
	if tagsParam != "" {
		for _, t := range strings.Split(tagsParam, ",") {
			if t = strings.TrimSpace(strings.ToLower(t)); t != "" {
				filterTags = append(filterTags, t)
			}
		}
	}

	all, err := h.svcFor(r).List(h.channelFor(r))
	degraded := false
	if err != nil {
		if !model.IsDegradedListErr(err) || len(all) == 0 {
			h.internalErr(w, err)
			return
		}
		// Mode dégradé : la blockchain est injoignable mais des brouillons locaux
		// existent — on les renvoie quand même, en le signalant explicitement au
		// client plutôt que de faire comme si la liste était complète.
		degraded = true
		log.Printf("liste dégradée (brouillons locaux uniquement) : %v", err)
	}

	// Filtrer
	var filtered []*model.Model3D
	visitor := isVisitor(r)
	for _, m := range all {
		// Un visiteur ne voit que le catalogue publié, jamais les brouillons (RM49).
		if visitor && m.Status != model.ModuleSubmitted {
			continue
		}
		if len(catSet) > 0 && !catSet[string(m.Category)] {
			continue
		}
		if q != "" && !assetMatchesQuery(m, q) {
			continue
		}
		if ownerID != "" && m.OwnerID != ownerID {
			continue
		}
		if parentID != "" && m.ParentID != parentID {
			continue
		}
		if hashFilter != "" && m.Hash != hashFilter {
			continue
		}
		if status != "" && string(m.Status) != status {
			continue
		}
		if len(filterTags) > 0 && !assetHasAnyTag(m, filterTags) {
			continue
		}
		filtered = append(filtered, m)
	}
	total := len(filtered)
	if limit < len(filtered) {
		filtered = filtered[:limit]
	}

	// Construire les DTOs
	loadedIDs := make(map[string]bool, len(filtered))
	dtos := make([]componentDTO, 0, len(filtered))
	for _, m := range filtered {
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		dtos = append(dtos, toComponentDTO(m, thumb))
		loadedIDs[m.ID] = true
	}

	// Ne renvoyer que les connexions entre assets chargés
	conns, _ := h.svcFor(r).ListConnections()
	connDTOs := make([]connectionDTO, 0)
	for _, c := range conns {
		if loadedIDs[c.From] && loadedIDs[c.To] {
			connDTOs = append(connDTOs, toConnectionDTO(c))
		}
	}

	resp := graphResponse{Components: dtos, Connections: connDTOs, Total: total}
	if degraded {
		resp.Degraded = true
		resp.Warning = model.DegradedListWarning
	}
	jsonOK(w, resp)
}

// assetHasAnyTag retourne true si l'asset possède au moins un des tags (comparaison insensible à la casse).
func assetHasAnyTag(m *model.Model3D, tags []string) bool {
	for _, ft := range tags {
		for _, t := range m.Tags {
			if strings.Contains(strings.ToLower(t), ft) {
				return true
			}
		}
	}
	return false
}

// assetMatchesQuery retourne true si l'asset correspond à la recherche textuelle (nom, tags, description).
func assetMatchesQuery(m *model.Model3D, q string) bool {
	if strings.Contains(strings.ToLower(m.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(m.Description), q) {
		return true
	}
	for _, t := range m.Tags {
		if strings.Contains(strings.ToLower(t), q) {
			return true
		}
	}
	return false
}

// createAsset crée un composant à partir d'un formulaire multipart.
//
//	@Summary		Créer un composant
//	@Description	Crée un asset. La catégorie (RM02) et le canal sont obligatoires : aucune valeur par défaut, 400 s'ils manquent. Le fichier CAO 3D est optionnel (champ "file", ou "stl" pour rétrocompatibilité) ; une miniature base64 peut être fournie, le serveur n'en dérive aucune de lui-même. Le composant est toujours créé en brouillon local ("status":"draft", aucune transaction blockchain) — voir POST /components/{id}/submit pour l'engager explicitement sur la blockchain.
//	@Tags			components
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name		formData	string	true	"nom du composant (max 256 caractères)"
//	@Param			description	formData	string	false	"description (max 10000 caractères)"
//	@Param			category	formData	string	true	"catégorie (RM02) : base, amelioration, variation, adaptation, derivation, extension ou regression"
//	@Param			owner_id	formData	string	false	"propriétaire — ignoré quand la session est liée à une identité : le propriétaire est alors cette identité (pseudo@org), qui signera la soumission"
//	@Param			parent_id	formData	string	false	"asset parent (dérivation)"
//	@Param			channel_id	formData	string	true	"canal cible"
//	@Param			license_id	formData	string	false	"licence"
//	@Param			tags		formData	string	false	"tags séparés par virgule"
//	@Param			links		formData	string	false	"liens externes, tableau JSON encodé en chaîne"
//	@Param			file		formData	file	false	"fichier CAO 3D"
//	@Param			thumbnail	formData	string	false	"miniature en data URL base64"
//	@Success		201	{object}	componentDTO
//	@Failure		400	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components [post]
func (h *Handler) createAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		jsonError(w, "parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	ownerID := strings.TrimSpace(r.FormValue("owner_id"))
	channelID := strings.TrimSpace(r.FormValue("channel_id"))

	if name == "" {
		jsonError(w, "le champ 'name' est requis", http.StatusBadRequest)
		return
	}
	if len(name) > 256 {
		jsonError(w, "name trop long (max 256)", http.StatusBadRequest)
		return
	}
	if len(description) > 10000 {
		jsonError(w, "description trop longue (max 10000)", http.StatusBadRequest)
		return
	}
	if ownerID != "" && !validFieldID.MatchString(ownerID) {
		jsonError(w, "owner_id invalide", http.StatusBadRequest)
		return
	}
	ownerID = ownerFor(r, ownerID)
	if channelID != "" && !validFieldID.MatchString(channelID) {
		jsonError(w, "channel_id invalide", http.StatusBadRequest)
		return
	}

	req := model.AddRequest{
		Name:        name,
		Description: description,
		Category:    model.Category(r.FormValue("category")),
		OwnerID:     ownerID,
		ParentID:    strings.TrimSpace(r.FormValue("parent_id")),
		ChannelID:   channelID,
		LicenseID:   strings.TrimSpace(r.FormValue("license_id")),
		Tags:        parseTags(r.FormValue("tags")),
	}
	if raw := r.FormValue("links"); raw != "" {
		_ = json.Unmarshal([]byte(raw), &req.Links)
	}

	// Fichier 3D optionnel — accepte le champ "file" (générique) ou "stl" (rétrocompat)
	fileField := "file"
	if _, _, err := r.FormFile("file"); err != nil {
		fileField = "stl"
	}
	file, header, ferr := r.FormFile(fileField)
	if ferr != nil {
		file, header = nil, nil
	}

	m, err := h.createOneComponent(r, req, file, header)
	if err != nil {
		h.internalErr(w, err)
		return
	}

	// Miniature fournie par le client (rendu côté client) : le serveur n'en
	// produit aucune de lui-même (FT-040).
	if thumb := r.FormValue("thumbnail"); thumb != "" {
		_ = h.svcFor(r).SaveThumbnail(m.ID, thumb)
	}

	thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, toComponentDTO(m, thumb))
}

// createOneComponent crée un composant à partir d'un gabarit AddRequest (Name
// déjà renseigné) et d'un fichier CAO optionnel (nil si aucun fichier).
func (h *Handler) createOneComponent(r *http.Request, req model.AddRequest, file multipart.File, header *multipart.FileHeader) (*model.Model3D, error) {
	if file != nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".bin"
		}
		tmp, err := os.CreateTemp("", "myr-model-*"+ext)
		if err != nil {
			return nil, err
		}
		defer os.Remove(tmp.Name())
		tmp.Close()

		// Réécrire proprement
		if f2, err := os.Create(tmp.Name()); err == nil {
			defer f2.Close()
			buf := make([]byte, 32*1024)
			for {
				n, readErr := file.Read(buf)
				if n > 0 {
					f2.Write(buf[:n])
				}
				if readErr != nil {
					break
				}
			}
		}
		req.FilePath = tmp.Name()
	}
	return h.svcFor(r).AddFull(req)
}

// ── /api/components/:id ──────────────────────────────────────────────────────────

func (h *Handler) handleComponent(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/components/")
	if rest == "" {
		http.NotFound(w, r)
		return
	}
	// Un visiteur ne lit qu'un asset publié et ses sous-ressources (RM49) : un
	// brouillon lui est inconnu.
	if isVisitor(r) {
		id, _, _ := strings.Cut(rest, "/")
		if m, err := h.svcFor(r).Get(id, h.channelFor(r)); err != nil || m.Status != model.ModuleSubmitted {
			jsonError(w, "asset introuvable", http.StatusNotFound)
			return
		}
	}
	// Sous-route /:id/interfaces
	if strings.HasSuffix(rest, "/interfaces") {
		id := strings.TrimSuffix(rest, "/interfaces")
		h.handleComponentInterfaces(w, r, id)
		return
	}
	// Sous-route /:id/tree — arbre d'évolution complet (ancêtres + descendants)
	if strings.HasSuffix(rest, "/tree") {
		id := strings.TrimSuffix(rest, "/tree")
		h.handleComponentTree(w, r, id)
		return
	}
	// Sous-route /:id/thumbnail/regenerate — redérive la miniature depuis le lien source
	if strings.HasSuffix(rest, "/thumbnail/regenerate") {
		id := strings.TrimSuffix(rest, "/thumbnail/regenerate")
		h.regenerateThumbnail(w, r, id)
		return
	}
	// Sous-route /:id/thumbnail — GET récupère, POST remplace la miniature
	if strings.HasSuffix(rest, "/thumbnail") {
		id := strings.TrimSuffix(rest, "/thumbnail")
		switch r.Method {
		case http.MethodGet:
			thumb, _ := h.svcFor(r).GetThumbnail(id)
			jsonOK(w, map[string]string{"thumbnail": thumb})
		case http.MethodPost:
			h.saveComponentThumbnail(w, r, id)
		default:
			http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		}
		return
	}
	// Sous-route /:id/assemblies — GET liste les connexions internes, POST rattache
	// une connexion existante (ADR-11, ex-/api/modules/{id}/assemblies)
	if strings.HasSuffix(rest, "/assemblies") {
		id := strings.TrimSuffix(rest, "/assemblies")
		h.handleComponentAssemblies(w, r, id)
		return
	}
	// Sous-route /:id/assemblies/:connID — détache une connexion (ADR-11)
	if idx := strings.Index(rest, "/assemblies/"); idx != -1 {
		id := rest[:idx]
		connID := rest[idx+len("/assemblies/"):]
		h.removeComponentAssembly(w, r, id, connID)
		return
	}
	// Sous-route /:id/instances — GET liste les instances, POST en ajoute une
	// (ADR-11, ex-/api/modules/{id}/instances) — applicable à n'importe quel id
	// existant, y compris un composant qui n'a encore aucune instance.
	if strings.HasSuffix(rest, "/instances") {
		id := strings.TrimSuffix(rest, "/instances")
		h.handleComponentInstances(w, r, id)
		return
	}
	// Sous-route /:id/instances/:instanceID — retire ou repositionne une instance
	if idx := strings.Index(rest, "/instances/"); idx != -1 {
		id := rest[:idx]
		instanceID := rest[idx+len("/instances/"):]
		h.handleComponentInstance(w, r, id, instanceID)
		return
	}
	// Sous-route /:id/submit — engage un composant brouillon sur la blockchain
	// (route unique : dispatche en interne vers la logique historique de
	// SubmitModule dès que l'asset a des instances — ADR-11)
	if strings.HasSuffix(rest, "/submit") {
		id := strings.TrimSuffix(rest, "/submit")
		h.submitComponent(w, r, id)
		return
	}
	id := rest
	switch r.Method {
	case http.MethodDelete:
		h.deleteComponent(w, r, id)
	case http.MethodGet:
		h.getComponent(w, r, id)
	case http.MethodPatch:
		h.patchComponent(w, r, id)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// regenerateThumbnail redérive la miniature d'un asset depuis sa source
// durable — l'og:image du premier lien externe enregistré (Links). Sans
// lien externe, il n'existe pas de source régénérable côté serveur (le rendu d'un
// modèle 3D est produit par le client GUI, pas par myr) — voir ModelService.RegenerateThumbnail.
//
//	@Summary		Régénérer la miniature d'un asset depuis son lien source
//	@Description	Redérive la miniature depuis l'og:image du premier lien externe enregistré (Links). Échoue si l'asset n'a aucun lien externe — un modèle 3D sans lien n'a pas de source régénérable côté serveur. Répond { "thumbnail": "<dataURL>" } — jamais le componentDTO complet, qu'il faut requêter séparément si besoin d'autre chose que la miniature.
//	@Tags			components
//	@Produce		json
//	@Param			id	path		string	true	"identifiant de l'asset"
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	map[string]string	"aucun lien externe, échec réseau/timeout ou og:image introuvable — un seul code pour ces trois cas"
//	@Failure		503	{object}	map[string]string	"blockchain indisponible (asset non local et non joignable)"
//	@Security		MyrToken
//	@Router			/components/{id}/thumbnail/regenerate [post]
func (h *Handler) regenerateThumbnail(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	dataURL, err := h.svcFor(r).RegenerateThumbnail(id)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	jsonOK(w, map[string]string{"thumbnail": dataURL})
}

// saveComponentThumbnail remplace la miniature d'un composant, brouillon ou
// déjà soumis — la miniature est un champ hors ledger (comme name/description,
// voir patchComponent), modifiable même après soumission puisqu'elle ne touche
// jamais le contenu blockchain déjà commité (règle 7).
//
//	@Summary		Remplacer la miniature d'un composant
//	@Description	Remplace la miniature d'un composant (brouillon ou déjà soumis) par la data URL fournie. Champ hors ledger — modifiable même après soumission, contrairement au contenu blockchain (règle 7).
//	@Tags			components
//	@Accept			json
//	@Param			id		path	string	true	"identifiant du composant"
//	@Param			body	body	object	true	"champ thumbnail requis (data URL base64)"
//	@Success		204	"pas de contenu"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id}/thumbnail [post]
func (h *Handler) saveComponentThumbnail(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Thumbnail string `json:"thumbnail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "body JSON invalide", http.StatusBadRequest)
		return
	}
	if err := h.svcFor(r).SaveThumbnail(id, body.Thumbnail); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// submitComponent engage un asset créé en brouillon sur la blockchain — un
// composant simple comme un asset décomposé (ADR-11, specs/3-Conception/Conception_intro.md) :
// un seul point d'entrée, ModelService.Submit route en interne vers la
// logique historique de SubmitModule (RM17, ModuleVersion) dès que l'asset a
// des instances (IsModule()) — ce handler n'a pas à connaître cette distinction.
//
//	@Summary		Soumettre un composant en brouillon
//	@Description	Committe en une transaction les métadonnées et les interfaces locales du composant, puis passe "status" à "submitted". Pour un asset décomposé, vérifie qu'au moins un assemblage existe (RM17) et crée une ModuleVersion horodatée — le champ "note" optionnel du corps de la requête n'est utilisé que dans ce cas. Un asset déjà soumis n'est jamais re-soumis (RM19, 409) : une nouvelle version est un asset dérivé (parent_id).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"identifiant du composant en brouillon"
//	@Param			body	body	object	false	"note optionnelle, utilisée uniquement pour un asset décomposé (ModuleVersion.Note)"
//	@Success		200	{object}	componentDTO
//	@Failure		400	{object}	map[string]string
//	@Failure		409	{object}	map[string]string	"asset déjà soumis"
//	@Failure		503	{object}	map[string]string	"blockchain indisponible"
//	@Security		MyrToken
//	@Router			/components/{id}/submit [post]
func (h *Handler) submitComponent(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	m, err := h.svcFor(r).Submit(id, body.Note)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
	jsonOK(w, toComponentDTO(m, thumb))
}

// handleComponentAssemblies liste (GET) les connexions internes d'un asset
// décomposé, ou rattache (POST) une connexion déjà créée à sa liste
// d'assemblage — ex-/api/modules/{id}/connections (GET) et
// /api/modules/{id}/assemblies (POST), fusionnés sous une seule route (ADR-11).
//
//	@Summary		Connexions internes d'un asset décomposé
//	@Description	GET renvoie les connexions référencées par Assemblies (liste vide si l'asset n'a aucune instance). POST rattache une connexion déjà créée (myr model link add) à cette liste — brouillon local, sans effet blockchain avant soumission. La connexion doit exister et relier deux instances de l'asset (400 sinon).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"identifiant de l'asset"
//	@Param			body	body		object	false	"connection_id requis — POST uniquement"
//	@Success		200	{object}	componentDTO
//	@Failure		400	{object}	map[string]string	"connexion inexistante ou qui ne relie pas deux instances de l'asset"
//	@Failure		409	{object}	map[string]string	"asset déjà soumis (RM19)"
//	@Security		MyrToken
//	@Router			/components/{id}/assemblies [get]
//	@Router			/components/{id}/assemblies [post]
func (h *Handler) handleComponentAssemblies(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		m, err := h.svcFor(r).Get(id, h.channelFor(r))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		asmSet := map[string]bool{}
		for _, aid := range m.Assemblies {
			asmSet[aid] = true
		}
		all, _ := h.svcFor(r).ListConnections()
		dtos := make([]connectionDTO, 0)
		for _, c := range all {
			if asmSet[c.ID] {
				dtos = append(dtos, toConnectionDTO(c))
			}
		}
		jsonOK(w, dtos)
	case http.MethodPost:
		var body struct {
			ConnectionID string `json:"connection_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ConnectionID == "" {
			jsonError(w, "le champ 'connection_id' est requis", http.StatusBadRequest)
			return
		}
		if err := h.svcFor(r).AddAssemblyToModule(id, body.ConnectionID); err != nil {
			h.internalErr(w, err)
			return
		}
		m, err := h.svcFor(r).Get(id, h.channelFor(r))
		if err != nil {
			h.internalErr(w, err)
			return
		}
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		jsonOK(w, toComponentDTO(m, thumb))
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// removeComponentAssembly détache une connexion de la liste d'assemblage d'un
// asset — ex-DELETE /api/modules/{id}/assemblies/{connID} (ADR-11).
//
//	@Summary		Détacher une connexion d'un asset décomposé
//	@Description	Retire connID de la liste Assemblies de l'asset — brouillon local, sans effet blockchain avant soumission. Ne supprime pas la connexion elle-même (voir DELETE /connections/{id}).
//	@Tags			components
//	@Produce		json
//	@Param			id		path		string	true	"identifiant de l'asset"
//	@Param			connID	path		string	true	"identifiant de la connexion à détacher"
//	@Success		200	{object}	componentDTO
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id}/assemblies/{connID} [delete]
func (h *Handler) removeComponentAssembly(w http.ResponseWriter, r *http.Request, id, connID string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	if err := h.svcFor(r).RemoveAssemblyFromModule(id, connID); err != nil {
		h.internalErr(w, err)
		return
	}
	m, err := h.svcFor(r).Get(id, h.channelFor(r))
	if err != nil {
		h.internalErr(w, err)
		return
	}
	thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
	jsonOK(w, toComponentDTO(m, thumb))
}

// handleComponentInstances liste (GET) ou ajoute (POST) une instance d'un
// asset — ex-/api/modules/{id}/instances, applicable à n'importe quel id
// existant (ADR-11), y compris un composant qui n'a encore aucune instance.
//
//	@Summary		Instances d'un asset
//	@Description	GET renvoie la liste des instances (vide si l'asset n'a jamais été décomposé). POST ajoute un composant existant comme instance — brouillon local, sans effet blockchain avant soumission. Un asset qui est l'asset lui-même ou le contient, même indirectement, est refusé (400).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"identifiant de l'asset"
//	@Param			body	body		object	false	"asset_id requis — POST uniquement"
//	@Success		200	{array}		model.WorkspaceInstance
//	@Success		200	{object}	componentDTO
//	@Failure		400	{object}	map[string]string	"composition cyclique"
//	@Failure		404	{object}	map[string]string	"asset ajouté introuvable"
//	@Failure		409	{object}	map[string]string	"asset déjà soumis (RM19)"
//	@Security		MyrToken
//	@Router			/components/{id}/instances [get]
//	@Router			/components/{id}/instances [post]
func (h *Handler) handleComponentInstances(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		m, err := h.svcFor(r).Get(id, h.channelFor(r))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		instances := m.WorkspaceInstances
		if instances == nil {
			instances = []model.WorkspaceInstance{}
		}
		jsonOK(w, instances)
	case http.MethodPost:
		var body struct {
			AssetID string `json:"asset_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AssetID == "" {
			jsonError(w, "le champ 'asset_id' est requis", http.StatusBadRequest)
			return
		}
		m, err := h.svcFor(r).AddAssetToWorkspace(id, body.AssetID)
		if err != nil {
			h.internalErr(w, err)
			return
		}
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		jsonOK(w, toComponentDTO(m, thumb))
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// handleComponentInstance retire (DELETE) ou repositionne (PATCH) une
// instance d'un asset — ex-/api/modules/{id}/instances/{instanceId} (ADR-11).
//
//	@Summary		Retirer ou repositionner une instance
//	@Description	DELETE retire l'instance et ses connexions en cascade (RM14). PATCH met à jour sa position (x, y).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"identifiant de l'asset"
//	@Param			instanceId	path		string	true	"identifiant de l'instance"
//	@Param			body		body		object	false	"x, y — PATCH uniquement"
//	@Success		200	{object}	componentDTO
//	@Failure		400	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id}/instances/{instanceId} [delete]
//	@Router			/components/{id}/instances/{instanceId} [patch]
func (h *Handler) handleComponentInstance(w http.ResponseWriter, r *http.Request, id, instanceID string) {
	switch r.Method {
	case http.MethodDelete:
		m, err := h.svcFor(r).RemoveAssetFromWorkspace(id, instanceID)
		if err != nil {
			h.internalErr(w, err)
			return
		}
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		jsonOK(w, toComponentDTO(m, thumb))
	case http.MethodPatch:
		var body struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "body JSON invalide", http.StatusBadRequest)
			return
		}
		m, err := h.svcFor(r).UpdateInstancePosition(id, instanceID, body.X, body.Y)
		if err != nil {
			h.internalErr(w, err)
			return
		}
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		jsonOK(w, toComponentDTO(m, thumb))
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// patchComponent modifie un composant existant (champs partiels).
//
//	@Summary		Modifier un composant
//	@Description	Modifie un composant en brouillon. Champs modifiables : name, description, license_id (présente au catalogue, compatible avec la licence du parent — RM03), tags, links. Un composant soumis est immuable (RM19) : 409, toute évolution passe par un composant dérivé (parent_id).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"identifiant du composant"
//	@Param			body	body		object	true	"champs à modifier (name, description, license_id, tags, links)"
//	@Success		200	{object}	componentDTO
//	@Failure		400	{object}	map[string]string	"licence absente du catalogue ou incompatible avec le parent"
//	@Failure		409	{object}	map[string]string	"composant déjà soumis"
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id} [patch]
func (h *Handler) patchComponent(w http.ResponseWriter, r *http.Request, id string) {
	req, err := decodeUpdateRequest(r, id)
	if err != nil {
		jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}
	m, err := h.svcFor(r).UpdateAsset(req)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
	jsonOK(w, toComponentDTO(m, thumb))
}

// decodeUpdateRequest décode le corps JSON partagé par PATCH /components/{id}
// et PATCH /modules/{id} — mêmes champs modifiables (name, description,
// license_id, tags, links) pour tout asset, composant ou module.
func decodeUpdateRequest(r *http.Request, id string) (model.UpdateRequest, error) {
	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		LicenseID   string   `json:"license_id"`
		Tags        []string `json:"tags"`
		Links       []string `json:"links"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return model.UpdateRequest{}, err
	}
	return model.UpdateRequest{
		ID: id, Name: body.Name, Description: body.Description,
		LicenseID: body.LicenseID, Tags: body.Tags, Links: body.Links,
	}, nil
}

// getComponent renvoie le détail d'un composant.
//
//	@Summary	Détail d'un composant
//	@Tags		components
//	@Produce	json
//	@Param		id	path		string	true	"identifiant du composant"
//	@Success	200	{object}	componentDTO
//	@Failure	404	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/components/{id} [get]
func (h *Handler) getComponent(w http.ResponseWriter, r *http.Request, id string) {
	m, err := h.svcFor(r).Get(id, h.channelFor(r))
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
	jsonOK(w, toComponentDTO(m, thumb))
}

// deleteComponent supprime un composant du stockage local.
//
//	@Summary		Supprimer un composant
//	@Description	Un brouillon est réellement retiré du stockage local. Un composant déjà soumis est masqué localement — il disparaît de GET /components mais reste consultable via GET /components/{id} — la blockchain Fabric ne supporte aucune suppression (règle 9) et son enregistrement n'est jamais modifié.
//	@Tags			components
//	@Param			id	path	string	true	"identifiant du composant"
//	@Success		204	"pas de contenu"
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id} [delete]
func (h *Handler) deleteComponent(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.svcFor(r).Remove(id); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleComponentTree renvoie l'arbre d'évolution complet d'un asset.
//
//	@Summary		Arbre d'évolution d'un composant
//	@Description	Tous les ancêtres (via parent_id) et descendants (parcours en largeur) de l'asset, plus les connexions internes à cet ensemble.
//	@Tags			components
//	@Produce		json
//	@Param			id	path		string	true	"identifiant du composant"
//	@Success		200	{object}	graphResponse
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id}/tree [get]
func (h *Handler) handleComponentTree(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	all, err := h.svcFor(r).List(h.channelFor(r))
	if err != nil {
		h.internalErr(w, err)
		return
	}

	// Index par ID
	byID := make(map[string]*model.Model3D, len(all))
	for _, m := range all {
		byID[m.ID] = m
	}

	// Index enfants : parentID → liste d'enfants
	children := make(map[string][]*model.Model3D)
	for _, m := range all {
		if m.ParentID != "" {
			children[m.ParentID] = append(children[m.ParentID], m)
		}
	}

	treeIDs := make(map[string]bool)

	// Ancêtres (chemin vers la racine)
	cur := id
	for cur != "" {
		treeIDs[cur] = true
		if m, ok := byID[cur]; ok {
			cur = m.ParentID
		} else {
			break
		}
	}

	// Descendants (BFS)
	queue := []string{id}
	for len(queue) > 0 {
		cur, queue = queue[0], queue[1:]
		for _, child := range children[cur] {
			if !treeIDs[child.ID] {
				treeIDs[child.ID] = true
				queue = append(queue, child.ID)
			}
		}
	}

	// Construire les DTOs
	dtos := make([]componentDTO, 0, len(treeIDs))
	for aid := range treeIDs {
		m, ok := byID[aid]
		if !ok {
			continue
		}
		thumb, _ := h.svcFor(r).GetThumbnail(m.ID)
		dtos = append(dtos, toComponentDTO(m, thumb))
	}

	// Connexions entre membres de l'arbre
	conns, _ := h.svcFor(r).ListConnections()
	connDTOs := make([]connectionDTO, 0)
	for _, c := range conns {
		if treeIDs[c.From] && treeIDs[c.To] {
			connDTOs = append(connDTOs, toConnectionDTO(c))
		}
	}

	jsonOK(w, graphResponse{Components: dtos, Connections: connDTOs, Total: len(dtos)})
}

// ── /api/connections ─────────────────────────────────────────────────────────

func (h *Handler) handleConnections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listConnections(w, r)
	case http.MethodPost:
		h.createConnection(w, r)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// listConnections liste toutes les connexions — directes ou entre interfaces —,
// celles d'une composition comprises.
//
//	@Summary	Lister les connexions
//	@Tags		connections
//	@Produce	json
//	@Success	200	{array}	connectionDTO
//	@Security	MyrToken
//	@Router		/connections [get]
func (h *Handler) listConnections(w http.ResponseWriter, r *http.Request) {
	conns, err := h.svcFor(r).ListConnections()
	if err != nil {
		h.internalErr(w, err)
		return
	}
	dtos := make([]connectionDTO, 0, len(conns))
	for _, c := range conns {
		dtos = append(dtos, toConnectionDTO(c))
	}
	jsonOK(w, dtos)
}

// createConnection crée une connexion simple entre deux assets.
//
//	@Summary		Créer une connexion
//	@Description	Une connexion identique (même from et même to, sans interface) est refusée : le client la retrouve dans GET /connections.
//	@Tags			connections
//	@Accept			json
//	@Produce		json
//	@Param			body	body		connectionDTO	true	"from et to requis ; label optionnel"
//	@Success		201	{object}	connectionDTO
//	@Failure		400	{object}	map[string]string
//	@Failure		409	{object}	map[string]string	"connexion déjà existante"
//	@Security	MyrToken
//	@Router		/connections [post]
func (h *Handler) createConnection(w http.ResponseWriter, r *http.Request) {
	var body connectionDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.From == "" || body.To == "" {
		jsonError(w, "les champs 'from' et 'to' sont requis", http.StatusBadRequest)
		return
	}
	conn, err := h.svcFor(r).AddConnection(body.From, body.To, body.Label)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, toConnectionDTO(conn))
}

// ── /api/connections/:id ─────────────────────────────────────────────────────

// handleConnection supprime une connexion existante.
//
//	@Summary	Supprimer une connexion
//	@Tags		connections
//	@Param		id	path	string	true	"identifiant de la connexion"
//	@Success	204	"pas de contenu"
//	@Failure	500	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/connections/{id} [delete]
func (h *Handler) handleConnection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/connections/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	if err := h.svcFor(r).RemoveConnection(id); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── /api/assembly-links ──────────────────────────────────────────────────────

type assemblyLinkRequest struct {
	FromIfaceID     string `json:"from_iface_id"`
	ToIfaceID       string `json:"to_iface_id"`
	Label           string `json:"label,omitempty"`
	FromInstanceID  string `json:"from_instance_id,omitempty"`
	ToInstanceID    string `json:"to_instance_id,omitempty"`
	FastenerAssetID string `json:"fastener_asset_id,omitempty"`
}

// handleAssemblyLinks crée une liaison interface→interface (assemblage).
//
//	@Summary		Créer une liaison entre deux interfaces
//	@Description	Sans asset d'accroche, les deux interfaces doivent être compatibles (RM10, critères RM11) ; avec une accroche, c'est elle qui doit l'être avec chacune. Une interface d'une instance donnée ne sert qu'à une liaison (RM09) et n'est jamais reliée à elle-même. Une liaison qui devient incompatible plus tard porte incompatible:true et n'est pas supprimée automatiquement (RM12).
//	@Tags			connections
//	@Accept			json
//	@Produce		json
//	@Param			body	body		assemblyLinkRequest	true	"from_iface_id et to_iface_id requis"
//	@Success		201	{object}	connectionDTO
//	@Failure		400	{object}	map[string]string	"interfaces incompatibles, ou interface reliée à elle-même"
//	@Failure		409	{object}	map[string]string	"interface déjà engagée dans une liaison (RM09)"
//	@Security		MyrToken
//	@Router			/assembly-links [post]
func (h *Handler) handleAssemblyLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body assemblyLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.FromIfaceID == "" || body.ToIfaceID == "" {
		jsonError(w, "les champs 'from_iface_id' et 'to_iface_id' sont requis", http.StatusBadRequest)
		return
	}
	conn, err := h.svcFor(r).AddAssemblyLink(body.FromIfaceID, body.ToIfaceID, body.Label, body.FromInstanceID, body.ToInstanceID, body.FastenerAssetID)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, toConnectionDTO(conn))
}

// ── /api/channels ────────────────────────────────────────────────────────────

// handleChannels lit ou change le canal Fabric actif de la session.
//
//	@Summary		Lire ou changer le canal actif
//	@Description	GET est public. PUT accepte X-Myr-Token ou une session déjà authentifiée et valide le canal contre la liste des canaux connus du réseau actif.
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object	false	"{ \"channel\": \"sandbox\" } — PUT uniquement"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/channels [get]
//	@Router			/channels [put]
func (h *Handler) handleChannels(w http.ResponseWriter, r *http.Request) {
	sess := sessionFromCtx(r)
	active := h.netInfo.Channel
	if sess != nil && sess.Channel != "" {
		active = sess.Channel
	}

	// Liste des canaux : depuis le réseau actif (multi-réseau) ou depuis netInfo (legacy)
	channels := h.channelsForNetwork(r)

	switch r.Method {
	case http.MethodGet:
		jsonOK(w, map[string]any{
			"channels":   channels,
			"active":     active,
			"network_id": h.networkFor(r),
		})

	case http.MethodPut:
		// Accepte X-Myr-Token (legacy) ou session déjà injectée par requireAuth (JWT)
		token := r.Header.Get("X-Myr-Token")
		if sess == nil {
			if h.sessions.get(token) == nil {
				jsonError(w, "authentification requise", http.StatusUnauthorized)
				return
			}
		} else {
			token = sess.Token
		}
		var body struct {
			Channel string `json:"channel"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Channel == "" {
			jsonError(w, "champ 'channel' requis", http.StatusBadRequest)
			return
		}
		// Valider que le canal est dans la liste connue (si liste non vide)
		if len(channels) > 0 {
			known := false
			for _, ch := range channels {
				if ch == body.Channel {
					known = true
					break
				}
			}
			if !known {
				jsonError(w, fmt.Sprintf("canal '%s' inconnu sur ce réseau", body.Channel), http.StatusBadRequest)
				return
			}
		}
		h.sessions.setChannel(token, body.Channel)
		jsonOK(w, map[string]any{"channel": body.Channel})

	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// channelsForNetwork retourne la liste des canaux du réseau actif de la requête.
// Multi-réseau : depuis networkSvc. Legacy : depuis netInfo.
func (h *Handler) channelsForNetwork(r *http.Request) []string {
	if h.networkSvc != nil {
		networkID := h.networkFor(r)
		if profiles, err := h.networkSvc.List(); err == nil {
			for _, p := range profiles {
				if p.ID == networkID {
					return p.Channels
				}
			}
		}
	}
	if len(h.netInfo.Channels) > 0 {
		return h.netInfo.Channels
	}
	// Pas de réseau configuré → canal unique ou aucune contrainte
	if h.netInfo.Channel != "" {
		return []string{h.netInfo.Channel}
	}
	return nil // aucun réseau → pas de validation du canal
}

// ── /api/ping ────────────────────────────────────────────────────────────────

// handlePing est un endpoint de liveness minimal (aucun appel Fabric).
//
//	@Summary	Liveness
//	@Tags		infra
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/ping [get]
func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"ok": true, "version": h.version, "build_date": h.buildDate})
}

// ── /api/status ──────────────────────────────────────────────────────────────

// handleStatus renvoie l'état du serveur (mode, connexion Fabric, nombre d'assets).
//
//	@Summary	Statut du serveur
//	@Tags		infra
//	@Produce	json
//	@Success	200	{object}	statusResponse
//	@Router		/status [get]
func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	channel := h.channelFor(r)
	jsonOK(w, statusResponse{
		Mode:      "blockchain",
		Connected: h.fabricConnected,
		Assets:    h.assetCount(r, channel),
		Channel:   channel,
		Peer:      h.netInfo.Peer,
		Network:   h.netInfo.Network,
	})
}

// statusCountTTL est la durée pendant laquelle /api/status réutilise le nombre
// d'assets d'un canal : la route est publique, chaque calcul parcourt le ledger.
const statusCountTTL = 30 * time.Second

// assetCountCache retient, par canal, le dernier nombre d'assets calculé.
type assetCountCache struct {
	mu     sync.Mutex
	counts map[string]cachedCount
}

type cachedCount struct {
	n  int
	at time.Time
}

// assetCount retourne le nombre d'assets du canal, recalculé au plus une fois
// par statusCountTTL.
func (h *Handler) assetCount(r *http.Request, channel string) int {
	c := &h.assetCounts
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, ok := c.counts[channel]; ok && time.Since(cached.at) < statusCountTTL {
		return cached.n
	}
	n := 0
	if all, err := h.svcFor(r).List(channel); err == nil {
		n = len(all)
	}
	if c.counts == nil {
		c.counts = make(map[string]cachedCount)
	}
	c.counts[channel] = cachedCount{n: n, at: time.Now()}
	return n
}

// ── /api/components/:id/interfaces ───────────────────────────────────────────────

// handleComponentInterfaces liste ou ajoute les interfaces d'un composant.
//
//	@Summary		Interfaces d'un composant
//	@Description	GET renvoie les interfaces physiques directes si l'asset n'a aucune instance, ou les interfaces exposées (non connectées en interne) s'il en a (ADR-11) — GetModuleInterfaces gère les deux cas. Tant que l'asset porteur est en brouillon, les interfaces sont éditées librement en local et ne rejoignent la blockchain qu'à la soumission de l'asset (règle 27). POST : l'asset doit exister (404) et être un brouillon (409) ; l'identifiant de l'interface est attribué par le serveur, un champ id envoyé est ignoré.
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"identifiant du composant"
//	@Param			body	body		model.AssetInterface	false	"interface à ajouter — POST uniquement"
//	@Success		200	{array}		model.AssetInterface
//	@Success		201	{object}	model.AssetInterface
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string	"asset introuvable"
//	@Failure		409	{object}	map[string]string	"asset déjà soumis (RM19)"
//	@Failure		500	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/components/{id}/interfaces [get]
//	@Router			/components/{id}/interfaces [post]
func (h *Handler) handleComponentInterfaces(w http.ResponseWriter, r *http.Request, assetID string) {
	switch r.Method {
	case http.MethodGet:
		ifaces, err := h.svcFor(r).GetModuleInterfaces(assetID)
		if err != nil {
			h.internalErr(w, err)
			return
		}
		if ifaces == nil {
			ifaces = []*model.AssetInterface{}
		}
		jsonOK(w, ifaces)
	case http.MethodPost:
		var iface model.AssetInterface
		if err := json.NewDecoder(r.Body).Decode(&iface); err != nil {
			jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
			return
		}
		iface.AssetID = assetID
		if err := h.svcFor(r).AddInterface(&iface); err != nil {
			h.internalErr(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		jsonOK(w, &iface)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// ── /api/interfaces/:id ──────────────────────────────────────────────────────

// handleInterface modifie ou supprime une interface existante.
//
//	@Summary		Modifier ou supprimer une interface
//	@Description	DELETE maintient l'invariant qu'un asset possède toujours au moins un slot virtuel (règle 13) : un nouveau slot virtuel est recréé si l'interface supprimée était physique et la dernière restante. PATCH ne change jamais l'asset porteur : un asset_id différent est refusé (400).
//	@Tags			components
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"identifiant de l'interface"
//	@Param			body	body		model.AssetInterface	false	"interface mise à jour — PATCH uniquement"
//	@Success		200	{object}	model.AssetInterface
//	@Success		204	"pas de contenu (DELETE)"
//	@Failure		400	{object}	map[string]string	"changement d'asset porteur"
//	@Failure		404	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/interfaces/{id} [patch]
//	@Router			/interfaces/{id} [delete]
func (h *Handler) handleInterface(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/interfaces/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodDelete:
		if err := h.svcFor(r).RemoveInterface(id); err != nil {
			h.internalErr(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPatch:
		existing, err := h.svcFor(r).GetInterface(id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(existing); err != nil {
			jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
			return
		}
		existing.ID = id // prevent ID overwrite
		if err := h.svcFor(r).UpdateInterface(existing); err != nil {
			h.internalErr(w, err)
			return
		}
		jsonOK(w, existing)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// ── /api/refs ────────────────────────────────────────────────────────────────

// handleRefs renvoie le vocabulaire de référence (catégories, types, unités).
//
//	@Summary	Vocabulaire de référence des interfaces
//	@Tags		refs
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/refs [get]
func (h *Handler) handleRefs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	refs, err := h.svcFor(r).GetRefs()
	if err != nil {
		h.internalErr(w, err)
		return
	}
	jsonOK(w, refs)
}

// handleRefCategories ajoute une catégorie au référentiel d'interfaces.
//
//	@Summary	Ajouter une catégorie de référence
//	@Tags		refs
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"name requis"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/refs/categories [post]
func (h *Handler) handleRefCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		jsonError(w, "le champ 'name' est requis", http.StatusBadRequest)
		return
	}
	if err := h.svcFor(r).AddRefCategory(body.Name); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	refs, _ := h.svcFor(r).GetRefs()
	jsonOK(w, refs)
}

// handleRefTypes ajoute un type au référentiel d'interfaces, pour une catégorie donnée.
//
//	@Summary	Ajouter un type de référence
//	@Tags		refs
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"category et name requis"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/refs/types [post]
func (h *Handler) handleRefTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Category == "" || body.Name == "" {
		jsonError(w, "les champs 'category' et 'name' sont requis", http.StatusBadRequest)
		return
	}
	if err := h.svcFor(r).AddRefType(body.Category, body.Name); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	refs, _ := h.svcFor(r).GetRefs()
	jsonOK(w, refs)
}

// handleRefUnits ajoute une unité au référentiel d'interfaces, pour une catégorie donnée.
//
//	@Summary	Ajouter une unité de référence
//	@Tags		refs
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"category et name requis"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/refs/units [post]
func (h *Handler) handleRefUnits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Category == "" || body.Name == "" {
		jsonError(w, "les champs 'category' et 'name' sont requis", http.StatusBadRequest)
		return
	}
	if err := h.svcFor(r).AddRefUnit(body.Category, body.Name); err != nil {
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	refs, _ := h.svcFor(r).GetRefs()
	jsonOK(w, refs)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func readJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func methodNotAllowed(w http.ResponseWriter) {
	jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// internalErr journalise l'erreur complète côté serveur sous une référence courte,
// renvoyée au client (champ "ref" et en-tête X-Myr-Error-Ref) : le testeur la
// communique, l'erreur se retrouve dans le journal du serveur. Distingue
// l'indisponibilité de la blockchain — non configurée (ErrBlockchainUnavailable)
// ou injoignable (ErrBlockchainUnreachable), toutes deux transitoires et non des
// erreurs de données (503) —, l'écriture sur le ledger sans identité signataire
// (ErrSignerRequired, ADR-12, 403), un asset d'une autre identité (ErrNotOwner,
// FT-001, 403), l'asset ou l'interface introuvable
// (ErrNotFound, 404), la modification d'un asset soumis
// (ErrAssetSubmitted, RM19, 409), l'interface déjà engagée (RM09) ou la liaison
// déjà existante (409), la licence refusée (catalogue ou RM03, 400) et la
// liaison ou la composition refusée (RM10, liaison d'une interface à elle-même,
// liaison étrangère au module, composition cyclique, 400) des autres échecs du
// service domaine (500).
// Pour une 500, le client ne reçoit le détail de l'erreur que si le serveur a été
// démarré avec WithDebugErrors.
func (h *Handler) internalErr(w http.ResponseWriter, err error) {
	ref := errorRef()
	log.Printf("erreur interne [%s] : %v", ref, err)
	w.Header().Set("X-Myr-Error-Ref", ref)
	code, msg := http.StatusInternalServerError, "erreur interne"
	switch {
	case errors.Is(err, model.ErrBlockchainUnavailable), errors.Is(err, model.ErrBlockchainUnreachable):
		code, msg = http.StatusServiceUnavailable, err.Error()
	case errors.Is(err, model.ErrNotFound):
		code, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, model.ErrNotOwner):
		code, msg = http.StatusForbidden, err.Error()
	case errors.Is(err, model.ErrSignerRequired):
		code, msg = http.StatusForbidden, err.Error()+" — reconnectez-vous avec votre identité (POST /api/identity/session)"
	case errors.Is(err, model.ErrAssetSubmitted), errors.Is(err, model.ErrInterfaceInUse), errors.Is(err, model.ErrConnectionExists):
		code, msg = http.StatusConflict, err.Error()
	case errors.Is(err, model.ErrThumbnailTooLarge):
		code, msg = http.StatusRequestEntityTooLarge, err.Error()
	case errors.Is(err, model.ErrUnknownLicense), errors.Is(err, model.ErrLicenseIncompatible),
		errors.Is(err, model.ErrIncompatibleInterfaces), errors.Is(err, model.ErrSelfLink),
		errors.Is(err, model.ErrConnectionNotInModule), errors.Is(err, model.ErrCyclicComposition),
		errors.Is(err, model.ErrInterfaceCarrierChange),
		errors.Is(err, model.ErrInvalidCategory), errors.Is(err, model.ErrChannelRequired):
		code, msg = http.StatusBadRequest, err.Error()
	case h.debugErrors:
		msg = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%s (réf. %s)", msg, ref), "ref": ref})
}

// errorRef retourne une référence aléatoire courte (8 caractères hexadécimaux)
// qui relie une réponse d'erreur à sa ligne du journal serveur.
func errorRef() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func parseTags(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// ── /api/licenses ─────────────────────────────────────────────────────────────

// GET  /api/licenses           → liste complète du catalogue
// GET  /api/licenses/:id       → une licence
// POST /api/licenses/check     → { parent_license_id, proposed_license_id } → LicenseCheck
// POST /api/licenses/check-product → { component_license_ids[], proposed_module_license_id } → LicenseCheck

// handleLicenses liste le catalogue complet des licences.
//
//	@Summary	Catalogue des licences
//	@Tags		licenses
//	@Produce	json
//	@Success	200	{array}	model.License
//	@Security	MyrToken
//	@Router		/licenses [get]
func (h *Handler) handleLicenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	jsonOK(w, h.svcFor(r).ListLicenses())
}

// handleLicense route le détail d'une licence ainsi que les vérifications de
// compatibilité (règle 8) sous /api/licenses/{...}.
//
//	@Summary		Détail d'une licence ou vérification de compatibilité
//	@Description	GET /api/licenses/{id} : détail. POST /api/licenses/check : { parent_license_id, proposed_license_id } → compatibilité pour une dérivation de composant. POST /api/licenses/check-product : { component_license_ids[], proposed_module_license_id } → compatibilité pour un module assemblant plusieurs composants.
//	@Tags			licenses
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	false	"identifiant de la licence — GET /api/licenses/{id}"
//	@Param			body	body		object	false	"corps de vérification — POST /api/licenses/check(-product)"
//	@Success		200	{object}	model.License
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Security		MyrToken
//	@Router			/licenses/{id} [get]
//	@Router			/licenses/check [post]
//	@Router			/licenses/check-product [post]
func (h *Handler) handleLicense(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/licenses/")
	if rest == "" {
		http.NotFound(w, r)
		return
	}
	// /api/licenses/check
	if rest == "check" {
		if r.Method != http.MethodPost {
			http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			ParentLicenseID   string `json:"parent_license_id"`
			ProposedLicenseID string `json:"proposed_license_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
			return
		}
		jsonOK(w, h.svcFor(r).CheckLicenseCompatibility(body.ParentLicenseID, body.ProposedLicenseID))
		return
	}
	// /api/licenses/check-product
	if rest == "check-product" {
		if r.Method != http.MethodPost {
			http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			ComponentLicenseIDs     []string `json:"component_license_ids"`
			ProposedModuleLicenseID string   `json:"proposed_module_license_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "JSON invalide: "+err.Error(), http.StatusBadRequest)
			return
		}
		jsonOK(w, h.svcFor(r).CheckModuleLicenseCompatibility(body.ComponentLicenseIDs, body.ProposedModuleLicenseID))
		return
	}
	// /api/licenses/:id
	if r.Method != http.MethodGet {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	l, err := h.svcFor(r).GetLicense(rest)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonOK(w, l)
}
