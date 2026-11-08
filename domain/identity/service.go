// domain/identity/service.go
package identity

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"myr-core/domain/role"
)

type Service struct {
	walletDir    string // ~/.Myr/wallets/
	ca           CAPort
	requestStore RequestStore     // nil = demandes non persistées
	roles        role.RoleService // nil = seuls les rôles intégrés sont attribuables
}

func NewService(walletDir string, ca CAPort) *Service {
	return &Service{walletDir: walletDir, ca: ca}
}

// WithRequestStore injecte le store de persistance des demandes de compte.
func (s *Service) WithRequestStore(rs RequestStore) *Service {
	s.requestStore = rs
	return s
}

// WithRoles injecte la table des rôles, référence des rôles attribuables à
// une identité (RM46).
func (s *Service) WithRoles(roles role.RoleService) *Service {
	s.roles = roles
	return s
}

// WalletDir retourne le chemin du répertoire des wallets locaux.
func (s *Service) WalletDir() string { return s.walletDir }

// NewWithCA retourne un nouveau Service avec le même walletDir mais un CA différent.
func (s *Service) NewWithCA(ca CAPort) *Service {
	return &Service{walletDir: s.walletDir, ca: ca, requestStore: s.requestStore, roles: s.roles}
}

// assignableRole retourne le rôle à attribuer — RoleReader si r est vide
// (RM21) — après avoir vérifié qu'il figure dans la table des rôles et, pour
// le rôle par défaut d'une inscription (byDefault), qu'il ne porte aucun droit
// d'administration (RM46).
func (s *Service) assignableRole(r string, byDefault bool) (string, error) {
	if r == "" {
		r = RoleReader
	}
	check := role.Assignable
	if byDefault {
		check = role.AssignableByDefault
	}
	if err := check(s.roles, r); err != nil {
		return "", fmt.Errorf("%w : %w", ErrInvalidRole, err)
	}
	return r, nil
}

// ListLocalWallets parcourt ~/.Myr/wallets/ et retourne les wallets trouvés.
func (s *Service) ListLocalWallets() ([]WalletEntry, error) {
	entries, err := os.ReadDir(s.walletDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var wallets []WalletEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		wallets = append(wallets, s.walletEntry(e.Name())) // ex: alice@Org1
	}
	return wallets, nil
}

// Wallet retourne le wallet local de l'identité handle (« pseudo@org », voir
// Handle). Retourne ErrWalletNotFound si l'identité n'a pas de wallet sur ce
// nœud. handle doit être un nom simple : il ne peut pas désigner un chemin
// hors du répertoire des wallets.
func (s *Service) Wallet(handle string) (WalletEntry, error) {
	if handle == "" || handle == "." || handle == ".." || strings.ContainsAny(handle, `/\`) {
		return WalletEntry{}, fmt.Errorf("identifiant de wallet invalide %q", handle)
	}
	w := s.walletEntry(handle)
	if _, err := os.Stat(w.CertPath()); err != nil {
		if os.IsNotExist(err) {
			return WalletEntry{}, fmt.Errorf("%w : %s", ErrWalletNotFound, handle)
		}
		return WalletEntry{}, err
	}
	return w, nil
}

// walletEntry décrit le wallet de handle sous walletDir, rôle et statut lus
// depuis son certificat local s'il est lisible.
func (s *Service) walletEntry(handle string) WalletEntry {
	name, orgID := parseHandle(handle)
	w := WalletEntry{
		Handle: handle,
		Name:   name,
		OrgID:  orgID,
		MSPDir: filepath.Join(s.walletDir, handle, "msp"),
	}
	if data, err := os.ReadFile(w.CertPath()); err == nil {
		s.applyCertAttributes(&w, string(data))
	}
	return w
}

// applyCertAttributes renseigne le rôle et le statut de w depuis les attributs
// de son certificat (Myr.role, Myr.status), seule source de vérité de
// l'autorisation (ADR-12). Ils restent vides si la CA n'est pas configurée ou
// si le certificat ne les porte pas.
func (s *Service) applyCertAttributes(w *WalletEntry, certPEM string) {
	if s.ca == nil {
		return
	}
	attrs, err := s.ca.CertAttributes(certPEM)
	if err != nil {
		return
	}
	w.Role, w.Status = attrs[AttrRole], attrs[AttrStatus]
}

// CreateAccount crée, à l'initiative de l'administrateur, l'identité
// « pseudo@org » d'un utilisateur, active d'emblée avec le rôle donné (vide =
// RoleReader) — seule voie d'inscription sur un réseau qui n'accepte pas
// l'auto-enregistrement (ADR-12). Le rôle doit exister dans la table des rôles
// (RM46). Retourne l'identifiant de l'identité et le secret d'enrôlement à
// transmettre à l'utilisateur.
func (s *Service) CreateAccount(ctx context.Context, pseudo, orgID, roleName, email string) (string, string, error) {
	if err := validateIdentityName(pseudo, orgID); err != nil {
		return "", "", err
	}
	roleName, err := s.assignableRole(roleName, false)
	if err != nil {
		return "", "", err
	}
	if s.ca == nil {
		return "", "", fmt.Errorf("aucun CA configuré — ajoutez un ca_endpoint au réseau actif")
	}
	handle := Handle(pseudo, orgID)
	secret, err := s.ca.Register(ctx, RegisterRequest{
		Name:   handle,
		OrgID:  orgID,
		Role:   roleName,
		Status: StatusActive,
		Email:  email,
	})
	if err != nil {
		return "", "", fmt.Errorf("enregistrement de l'identité %s auprès de la CA : %w", handle, err)
	}
	return handle, secret, nil
}

// Enroll effectue l'enrollment et sauvegarde le wallet localement.
func (s *Service) Enroll(ctx context.Context, name, secret, orgID string) (WalletEntry, error) {
	if s.ca == nil {
		return WalletEntry{}, fmt.Errorf("aucun CA configuré")
	}
	// L'identifiant CA d'une identité inscrite est "pseudo@org" (voir
	// RequestAccess) — l'authentification HTTP Basic de l'enrollment doit viser
	// ce même identifiant, pas le pseudo seul, sous peine de "user not found"
	// côté CA dès que le client suit la forme documentée {name, secret, org_id}.
	handle := Handle(name, orgID)
	certPEM, keyPEM, caCertPEM, err := s.ca.Enroll(ctx, handle, secret)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("enrollment : %w", err)
	}

	mspDir := filepath.Join(s.walletDir, handle, "msp")
	if err := saveWallet(mspDir, certPEM, keyPEM, caCertPEM); err != nil {
		return WalletEntry{}, fmt.Errorf("sauvegarde wallet : %w", err)
	}

	entry := WalletEntry{
		Handle: handle,
		Name:   name,
		OrgID:  orgID,
		MSPDir: mspDir,
	}
	s.applyCertAttributes(&entry, certPEM)
	return entry, nil
}

// GetStatus interroge la CA pour obtenir le statut courant de l'identité.
func (s *Service) GetStatus(ctx context.Context, wallet WalletEntry) (string, error) {
	if s.ca == nil {
		return "", fmt.Errorf("aucun CA configuré")
	}
	// Comme pour Enroll : l'identifiant CA est "pseudo@org" (wallet.Handle),
	// pas le pseudo seul (wallet.Name) — sans quoi la CA ne retrouve pas
	// l'identité (même cause que le bug d'enrollment corrigé ci-dessus).
	return s.ca.GetStatus(ctx, wallet.Handle)
}

// LoadGuestWallet charge un wallet invité depuis des fichiers de certificat existants
// (pré-enrôlés par l'admin) et le sauvegarde localement sous le handle "guest@<org>".
// Le wallet invité a Myr.role=reader : il permet la consultation sans écriture.
// caCertPath est optionnel (chaîne vide acceptée).
func (s *Service) LoadGuestWallet(certPath, keyPath, caCertPath, orgID string) (WalletEntry, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("certificat invité introuvable (%s) : %w", certPath, err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("clé invité introuvable (%s) : %w", keyPath, err)
	}
	var caCertPEM string
	if caCertPath != "" {
		if data, err := os.ReadFile(caCertPath); err == nil {
			caCertPEM = string(data)
		}
	}

	handle := Handle("guest", orgID)
	mspDir := filepath.Join(s.walletDir, handle, "msp")
	if err := saveWallet(mspDir, string(certPEM), string(keyPEM), caCertPEM); err != nil {
		return WalletEntry{}, fmt.Errorf("sauvegarde wallet invité : %w", err)
	}
	return WalletEntry{
		Handle: handle,
		Name:   "guest",
		OrgID:  orgID,
		Status: StatusActive,
		MSPDir: mspDir,
	}, nil
}

// ReEnroll re-enrôle l'identité pour récupérer un certificat mis à jour.
func (s *Service) ReEnroll(ctx context.Context, wallet WalletEntry) (WalletEntry, error) {
	if s.ca == nil {
		return WalletEntry{}, fmt.Errorf("aucun CA configuré")
	}
	certPath := wallet.CertPath()
	keyPath, err := findKeyFile(filepath.Dir(wallet.KeyPath()))
	if err != nil {
		return WalletEntry{}, fmt.Errorf("clé privée introuvable : %w", err)
	}
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("certificat local introuvable : %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("clé privée illisible : %w", err)
	}

	newCert, err := s.ca.ReEnroll(ctx, wallet.Name, string(certPEM), string(keyPEM))
	if err != nil {
		return WalletEntry{}, fmt.Errorf("re-enrollment : %w", err)
	}

	// Remplacer le certificat
	err = os.WriteFile(certPath, []byte(newCert), 0600)
	if err != nil {
		return WalletEntry{}, fmt.Errorf("mise à jour certificat : %w", err)
	}

	wallet.Role, wallet.Status = "", ""
	s.applyCertAttributes(&wallet, newCert)
	return wallet, nil
}

// RequestAccess traite une demande d'inscription selon les règles du réseau
// (ADR-12) :
//   - sans auto-enregistrement, la demande est refusée (ErrRegistrationClosed)
//     et rien n'est conservé — seul l'administrateur crée des comptes ;
//   - sinon l'identité « pseudo@org » est créée dans la CA avec le rôle par
//     défaut du réseau, active si aucune approbation n'est requise, en attente
//     sinon (activée par ApproveRequest). Un rôle par défaut inconnu ou portant
//     des droits d'administration refuse la demande (ErrInvalidRole, RM46).
//
// Retourne la demande enregistrée et le secret d'enrôlement de l'identité.
// Si la CA refuse l'enregistrement, la demande est conservée au statut
// RequestFailed et l'erreur est retournée.
func (s *Service) RequestAccess(ctx context.Context, req AccountRequest, policy RegistrationPolicy) (*AccountRequest, string, error) {
	if !policy.AllowAutoRegister {
		return nil, "", ErrRegistrationClosed
	}
	if req.Pseudo == "" || req.Email == "" || req.OrgID == "" {
		return nil, "", fmt.Errorf("pseudo, email et org_id sont requis pour une demande de compte")
	}
	if err := validateIdentityName(req.Pseudo, req.OrgID); err != nil {
		return nil, "", err
	}
	defaultRole, err := s.assignableRole(policy.DefaultRole, true)
	if err != nil {
		return nil, "", fmt.Errorf("rôle par défaut du réseau : %w", err)
	}
	if s.ca == nil {
		return nil, "", fmt.Errorf("aucun CA configuré — impossible d'enregistrer l'identité")
	}
	// Sans stockage, une demande en attente ne pourrait jamais être décidée.
	if policy.RequireApproval && s.requestStore == nil {
		return nil, "", fmt.Errorf("stockage des demandes d'inscription non configuré — approbation impossible")
	}

	identityStatus, requestStatus := StatusActive, RequestApproved
	if policy.RequireApproval {
		identityStatus, requestStatus = StatusPending, RequestPending
	}
	req.ID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	req.Status = RequestPending
	req.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.saveRequest(&req); err != nil {
		return nil, "", err
	}

	handle := Handle(req.Pseudo, req.OrgID)
	secret, err := s.ca.Register(ctx, RegisterRequest{
		Name:        handle,
		OrgID:       req.OrgID,
		Role:        defaultRole,
		Status:      identityStatus,
		DisplayName: req.DisplayName,
		Email:       req.Email,
	})
	if err != nil {
		req.Status = RequestFailed
		_ = s.saveRequest(&req)
		return nil, "", fmt.Errorf("enregistrement de l'identité %s auprès de la CA : %w", handle, err)
	}
	req.Status = requestStatus
	if err := s.saveRequest(&req); err != nil {
		return nil, "", err
	}
	return &req, secret, nil
}

// RequestRole traite la demande du rôle wanted par l'identité active du
// wallet w (UCA08, RM47). Elle est refusée si l'identité porte déjà ce rôle
// (ErrRoleAlreadyHeld), si le rôle est absent de la table (ErrInvalidRole,
// RM46) ou si une autre demande de rôle de l'identité est en attente
// (ErrRequestAlreadyPending). Un rôle en attribution automatique est
// attribué aussitôt et la demande est approuvée ; sinon elle reste en attente
// de l'administrateur (ApproveRequest, RejectRequest). Dans les deux cas, le
// nouveau rôle s'applique au prochain enrôlement de l'identité, c'est-à-dire
// à sa prochaine connexion.
func (s *Service) RequestRole(ctx context.Context, w WalletEntry, wanted string) (*AccountRequest, error) {
	if wanted == "" {
		return nil, fmt.Errorf("rôle demandé requis")
	}
	if err := w.CanSignIn(); err != nil {
		return nil, err
	}
	if wanted == w.Role {
		return nil, fmt.Errorf("%w : %s", ErrRoleAlreadyHeld, wanted)
	}
	granted, err := role.GrantedOnRequest(s.roles, wanted)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrInvalidRole, err)
	}
	if s.ca == nil {
		return nil, fmt.Errorf("aucun CA configuré — impossible de modifier le rôle")
	}
	// Sans stockage, une demande en attente ne pourrait jamais être décidée.
	if !granted && s.requestStore == nil {
		return nil, fmt.Errorf("stockage des demandes non configuré — approbation impossible")
	}
	if err := s.checkNoPendingRoleRequest(w.Handle); err != nil {
		return nil, err
	}

	req := &AccountRequest{
		ID:        fmt.Sprintf("req-%d", time.Now().UnixNano()),
		Kind:      RequestKindRole,
		Pseudo:    w.Name,
		OrgID:     w.OrgID,
		Role:      wanted,
		Status:    RequestPending,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if granted {
		if err := s.ca.UpdateAttributes(ctx, w.Handle, map[string]string{AttrRole: wanted}); err != nil {
			req.Status = RequestFailed
			_ = s.saveRequest(req)
			return nil, fmt.Errorf("attribution du rôle %s à %s : %w", wanted, w.Handle, err)
		}
		req.Status = RequestApproved
	}
	if err := s.saveRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

// checkNoPendingRoleRequest refuse une nouvelle demande de rôle tant qu'une
// autre de la même identité attend la décision de l'administrateur.
func (s *Service) checkNoPendingRoleRequest(handle string) error {
	if s.requestStore == nil {
		return nil
	}
	all, err := s.requestStore.FindAll()
	if err != nil {
		return err
	}
	for _, r := range all {
		if r.Kind == RequestKindRole && r.Status == RequestPending && Handle(r.Pseudo, r.OrgID) == handle {
			return fmt.Errorf("%w : %s (rôle %s)", ErrRequestAlreadyPending, r.ID, r.Role)
		}
	}
	return nil
}

// ApproveRequest approuve une demande en attente. Inscription : l'identité
// devient active (Myr.status = active). Demande de rôle : l'identité reçoit le
// rôle demandé (Myr.role), s'il figure toujours dans la table (RM46). Le
// changement s'applique au prochain enrôlement de l'identité, c'est-à-dire à
// sa prochaine connexion.
func (s *Service) ApproveRequest(ctx context.Context, id string) (*AccountRequest, error) {
	return s.decideRequest(ctx, id, true)
}

// RejectRequest refuse une demande en attente. Inscription : l'identité est
// suspendue (Myr.status = suspended) — elle ne peut ni se connecter ni écrire.
// Demande de rôle : l'identité garde son rôle. Rien n'est supprimé :
// l'administrateur peut revenir sur sa décision.
func (s *Service) RejectRequest(ctx context.Context, id string) (*AccountRequest, error) {
	return s.decideRequest(ctx, id, false)
}

// decideRequest applique à l'identité d'une demande en attente la décision
// approve (decisionAttributes), puis l'enregistre.
func (s *Service) decideRequest(ctx context.Context, id string, approve bool) (*AccountRequest, error) {
	if s.ca == nil {
		return nil, fmt.Errorf("aucun CA configuré — impossible de modifier l'identité")
	}
	if s.requestStore == nil {
		return nil, fmt.Errorf("stockage des demandes non configuré")
	}
	all, err := s.requestStore.FindAll()
	if err != nil {
		return nil, err
	}
	var req *AccountRequest
	for _, r := range all {
		if r.ID == id {
			req = r
			break
		}
	}
	if req == nil {
		return nil, fmt.Errorf("%w : %s", ErrRequestNotFound, id)
	}
	if req.Status != RequestPending {
		return nil, fmt.Errorf("%w : %s est %s", ErrRequestNotPending, id, req.Status)
	}
	attrs, err := s.decisionAttributes(req, approve)
	if err != nil {
		return nil, err
	}
	handle := Handle(req.Pseudo, req.OrgID)
	if len(attrs) > 0 {
		if err := s.ca.UpdateAttributes(ctx, handle, attrs); err != nil {
			return nil, fmt.Errorf("mise à jour de l'identité %s : %w", handle, err)
		}
	}
	req.Status = RequestRejected
	if approve {
		req.Status = RequestApproved
	}
	if err := s.saveRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

// decisionAttributes retourne les attributs de certificat que la décision
// approve change pour l'identité de req : son statut pour une inscription,
// son rôle pour une demande de rôle approuvée — revérifié dans la table, qui
// a pu changer depuis la demande (RM46) — et rien pour une demande de rôle
// refusée.
func (s *Service) decisionAttributes(req *AccountRequest, approve bool) (map[string]string, error) {
	if req.Kind != RequestKindRole {
		if approve {
			return map[string]string{AttrStatus: StatusActive}, nil
		}
		return map[string]string{AttrStatus: StatusSuspended}, nil
	}
	if !approve {
		return nil, nil
	}
	r, err := s.assignableRole(req.Role, false)
	if err != nil {
		return nil, err
	}
	return map[string]string{AttrRole: r}, nil
}

// saveRequest persiste req si un stockage des demandes est configuré.
func (s *Service) saveRequest(req *AccountRequest) error {
	if s.requestStore == nil {
		return nil
	}
	if err := s.requestStore.Save(req); err != nil {
		return fmt.Errorf("enregistrement de la demande : %w", err)
	}
	return nil
}

// validateIdentityName refuse un pseudo ou une organisation qui ne peut pas
// former un identifiant « pseudo@org » sûr : vide, réduit à "." ou "..", ou
// contenant un séparateur (@ de l'identifiant, : de l'authentification
// d'enrôlement, / ou \ d'un chemin de wallet) ou un caractère de contrôle.
func validateIdentityName(names ...string) error {
	for _, n := range names {
		bad := n == "" || n == "." || n == ".." || strings.ContainsAny(n, `@:/\`)
		for _, r := range n {
			if unicode.IsControl(r) || unicode.IsSpace(r) {
				bad = true
			}
		}
		if bad {
			return fmt.Errorf("%w : %q", ErrInvalidIdentityName, n)
		}
	}
	return nil
}

// SetRole change le rôle enregistré d'une identité existante auprès de la CA
// (attribut Myr.role). Le nouveau rôle ne s'applique qu'aux certificats émis
// après l'appel — l'identité doit se ré-enrôler (ReEnroll) pour l'obtenir dans
// son certificat actif : c'est une propriété de la Fabric CA, pas une
// limitation de myr. Le nouveau rôle doit exister dans la table des rôles
// (ErrInvalidRole, RM46).
func (s *Service) SetRole(ctx context.Context, name, newRole string) error {
	if s.ca == nil {
		return fmt.Errorf("aucun CA configuré — impossible de modifier le rôle")
	}
	if name == "" {
		return fmt.Errorf("identité requise")
	}
	if newRole == "" {
		return fmt.Errorf("nouveau rôle requis")
	}
	if _, err := s.assignableRole(newRole, false); err != nil {
		return err
	}
	return s.ca.UpdateAttributes(ctx, name, map[string]string{AttrRole: newRole})
}

// ListRequests retourne toutes les demandes enregistrées.
func (s *Service) ListRequests() ([]*AccountRequest, error) {
	if s.requestStore == nil {
		return nil, nil
	}
	return s.requestStore.FindAll()
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func saveWallet(mspDir, certPEM, keyPEM, caCertPEM string) error {
	writes := map[string]string{
		walletCertFile:   certPEM,
		walletKeyFile:    keyPEM,
		walletCACertFile: caCertPEM,
	}
	for rel, data := range writes {
		path := filepath.Join(mspDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			return err
		}
	}
	return nil
}

func parseHandle(handle string) (name, orgID string) {
	parts := strings.SplitN(handle, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1] + "MSP"
	}
	return handle, ""
}

func findKeyFile(keystoreDir string) (string, error) {
	entries, err := os.ReadDir(keystoreDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if !e.IsDir() {
			return filepath.Join(keystoreDir, e.Name()), nil
		}
	}
	return "", fmt.Errorf("aucun fichier de clé dans %s", keystoreDir)
}
