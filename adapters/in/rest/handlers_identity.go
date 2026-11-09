// adapters/in/rest/handlers_identity.go — handlers REST pour les identités (demande de compte)
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"myr-core/domain/identity"
	rbac "myr-core/domain/role"
)

// GuestHandle est le nom réservé pour l'identité invité partagée.
const GuestHandle = "guest"

// walletDTO est la représentation JSON d'un wallet local exposée au GUI.
type walletDTO struct {
	Handle string `json:"handle"`        // "alice@Org1"
	Name   string `json:"name"`          // "alice"
	OrgID  string `json:"org_id"`        // "Org1MSP"
	Status string `json:"status"`        // "pending" | "active" | "suspended"
	Guest  bool   `json:"guest,omitempty"` // true = identité invité (lecture seule)
}

// networkPolicyResponse est la réponse de GET /api/identity/policy.
// Le GUI l'interroge au démarrage pour savoir s'il doit afficher le mur d'identité.
type networkPolicyResponse struct {
	Mode              string `json:"mode"`                        // "blockchain"
	AllowAutoGuest    bool   `json:"allow_auto_guest"`            // le réseau délivre un accès reader automatiquement
	AllowAutoRegister bool   `json:"allow_auto_register"`         // le réseau accepte les demandes d'inscription
	RequireApproval   bool   `json:"require_approval"`            // une inscription doit être approuvée par l'admin avant d'être active
	NetworkName       string `json:"network_name,omitempty"`
}

// sessionDTO est le token de session retourné après identification réussie.
type sessionDTO struct {
	Token   string `json:"token"`
	Role    string `json:"role"`    // "reader" | "contributor" | ...
	Pseudo  string `json:"pseudo"`
	Channel string `json:"channel,omitempty"` // canal Fabric actif
	Guest   bool   `json:"guest,omitempty"`
}

func toWalletDTO(w identity.WalletEntry) walletDTO {
	return walletDTO{
		Handle: w.Handle,
		Name:   w.Name,
		OrgID:  w.OrgID,
		Status: w.Status,
		Guest:  w.Name == GuestHandle,
	}
}

// handleIdentityWallets liste les wallets présents localement (~/.Myr/wallets/) sur le
// serveur. Ce répertoire est partagé par tous les utilisateurs enrôlés sur ce nœud (pas
// seulement par l'appelant) : une session sans la permission identity.admin ne voit donc
// que le ou les wallets correspondant à son propre pseudo (Handle "<pseudo>@<org>"), jamais
// ceux des autres utilisateurs.
//
//	@Summary	Lister les wallets locaux
//	@Tags		identity
//	@Produce	json
//	@Success	200	{array}		walletDTO
//	@Failure	503	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/identity/wallets [get]
func (h *Handler) handleIdentityWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	wallets, err := h.identitySvc.ListLocalWallets()
	if err != nil {
		h.internalErr(w, err)
		return
	}
	sess := sessionFromCtx(r)
	self := !h.hasPermission(sess, rbac.PermIdentityAdmin)
	dtos := make([]walletDTO, 0, len(wallets))
	for _, we := range wallets {
		if self && we.Name != sess.Pseudo {
			continue
		}
		dtos = append(dtos, toWalletDTO(we))
	}
	jsonOK(w, dtos)
}

// handleIdentityEnroll enrôle un utilisateur auprès de la CA Fabric avec le
// secret d'enrollment fourni par l'admin, sans créer de session REST.
//
//	@Summary	Enrôlement CA (sans session)
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"name, secret, org_id requis"
//	@Success	200	{object}	walletDTO
//	@Failure	400	{object}	map[string]string
//	@Failure	429	{object}	map[string]string	"trop de tentatives (10 par minute et par adresse)"
//	@Failure	503	{object}	map[string]string
//	@Router		/identity/enroll [post]
func (h *Handler) handleIdentityEnroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Route publique : sans limite, on pourrait y essayer des secrets
	// d'enrôlement par force brute (R-06).
	if !h.authLimiter.allow(h.clientIP(r)) {
		jsonError(w, "trop de tentatives, réessayez dans une minute", http.StatusTooManyRequests)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	var req struct {
		Name   string `json:"name"`
		Secret string `json:"secret"`
		OrgID  string `json:"org_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.OrgID = strings.TrimSpace(req.OrgID)
	req.Secret = strings.TrimSpace(req.Secret)
	if req.Name == "" || req.Secret == "" || req.OrgID == "" {
		jsonError(w, "name, secret et org_id sont requis", http.StatusBadRequest)
		return
	}

	entry, err := h.identitySvc.Enroll(r.Context(), req.Name, req.Secret, req.OrgID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonOK(w, toWalletDTO(entry))
}

// handleIdentityStatus interroge la CA pour le statut courant d'un wallet.
//
//	@Summary	Statut CA d'un wallet
//	@Tags		identity
//	@Produce	json
//	@Param		handle	query		string	true	"identifiant pseudo@org, ex: alice@Org1"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/identity/status [get]
func (h *Handler) handleIdentityStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	handle := strings.TrimSpace(r.URL.Query().Get("handle"))
	if handle == "" {
		jsonError(w, "paramètre handle requis (ex: alice@Org1)", http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(handle, "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		jsonError(w, "handle invalide — format attendu : pseudo@org", http.StatusBadRequest)
		return
	}
	we := identity.WalletEntry{
		Handle: handle,
		Name:   parts[0],
		OrgID:  parts[1] + "MSP",
	}
	status, err := h.identitySvc.GetStatus(r.Context(), we)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	jsonOK(w, map[string]string{"handle": handle, "status": status})
}

// accountRequestDTO est le corps attendu pour POST /api/identity/request.
type accountRequestDTO struct {
	Pseudo      string `json:"pseudo"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email"`
	OrgID       string `json:"org_id"`
	Message     string `json:"message,omitempty"`
}

// handleIdentityRequest traite une demande d'inscription d'un utilisateur sans
// compte, selon les règles du réseau actif (ADR-12, identity.RequestAccess) :
// refusée si le réseau n'accepte pas l'auto-enregistrement ; sinon l'identité
// est créée, active ou en attente d'approbation, et le secret d'enrôlement est
// renvoyé.
//
//	@Summary	Demander un accès au réseau
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		accountRequestDTO	true	"pseudo, email, org_id requis"
//	@Success	201	{object}	map[string]string	"{id, status (approved | pending), pseudo, org_id, secret}"
//	@Failure	400	{object}	map[string]string
//	@Failure	403	{object}	map[string]string	"le réseau n'accepte pas les demandes d'inscription"
//	@Failure	429	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Router		/identity/request [post]
func (h *Handler) handleIdentityRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.authLimiter.allow(h.clientIP(r)) {
		jsonError(w, "trop de tentatives, réessayez dans une minute", http.StatusTooManyRequests)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	var dto accountRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		jsonError(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	dto.Pseudo = strings.TrimSpace(dto.Pseudo)
	dto.Email = strings.TrimSpace(dto.Email)
	dto.OrgID = strings.TrimSpace(dto.OrgID)
	if dto.Pseudo == "" || dto.Email == "" || dto.OrgID == "" {
		jsonError(w, "pseudo, email et org_id sont requis", http.StatusBadRequest)
		return
	}
	req := identity.AccountRequest{
		Pseudo:      dto.Pseudo,
		DisplayName: strings.TrimSpace(dto.DisplayName),
		Email:       dto.Email,
		OrgID:       dto.OrgID,
		Message:     strings.TrimSpace(dto.Message),
	}
	saved, secret, err := h.identitySvc.RequestAccess(r.Context(), req, h.registrationPolicy())
	switch {
	case errors.Is(err, identity.ErrRegistrationClosed):
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	case errors.Is(err, identity.ErrInvalidIdentityName):
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, map[string]string{
		"id":     saved.ID,
		"status": saved.Status,
		"pseudo": saved.Pseudo,
		"org_id": saved.OrgID,
		"secret": secret,
	})
}

// registrationPolicy retourne les règles d'inscription du réseau actif — sans
// réseau actif, l'inscription est fermée.
func (h *Handler) registrationPolicy() identity.RegistrationPolicy {
	if h.networkSvc == nil {
		return identity.RegistrationPolicy{}
	}
	p, err := h.networkSvc.GetActive()
	if err != nil {
		return identity.RegistrationPolicy{}
	}
	return identity.PolicyOf(p)
}

// handleIdentityRequestDecision approuve ou refuse une demande d'inscription
// ou de rôle en attente — réservé aux identités portant la permission
// identity.admin. Inscription : approuver active l'identité, refuser la
// suspend (ADR-12). Demande de rôle : approuver attribue le rôle, refuser le
// laisse inchangé.
//
//	@Summary	Approuver ou refuser une demande d'inscription ou de rôle
//	@Tags		identity
//	@Produce	json
//	@Param		id		path		string	true	"identifiant de la demande"
//	@Param		action	path		string	true	"approve ou reject"
//	@Success	200	{object}	identity.AccountRequest
//	@Failure	403	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	400	{object}	map[string]string	"rôle demandé retiré de la table des rôles depuis la demande"
//	@Failure	409	{object}	map[string]string	"la demande n'est plus en attente"
//	@Security	MyrToken
//	@Router		/identity/requests/{id}/{action} [post]
func (h *Handler) handleIdentityRequestDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	id, action, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/api/identity/requests/"), "/")
	var decide func(context.Context, string) (*identity.AccountRequest, error)
	switch action {
	case "approve":
		decide = h.identitySvc.ApproveRequest
	case "reject":
		decide = h.identitySvc.RejectRequest
	}
	if id == "" || decide == nil {
		http.NotFound(w, r)
		return
	}
	req, err := decide(r.Context(), id)
	switch {
	case errors.Is(err, identity.ErrRequestNotFound):
		jsonError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, identity.ErrRequestNotPending):
		jsonError(w, err.Error(), http.StatusConflict)
	case errors.Is(err, identity.ErrInvalidRole):
		jsonError(w, err.Error(), http.StatusBadRequest)
	case err != nil:
		h.internalErr(w, err)
	default:
		jsonOK(w, req)
	}
}

// roleRequestDTO est le corps attendu pour POST /api/identity/roles/request.
type roleRequestDTO struct {
	Role string `json:"role"`
}

// handleIdentityRoleRequest traite la demande d'un rôle par l'identité de la
// session (UCA08, identity.RequestRole) : un rôle en attribution automatique
// est attribué aussitôt (200), tout autre attend la décision de
// l'administrateur (202). Le nouveau rôle s'applique à la prochaine connexion.
//
//	@Summary	Demander un rôle
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		roleRequestDTO	true	"rôle demandé"
//	@Success	200	{object}	identity.AccountRequest	"rôle attribué (status approved)"
//	@Success	202	{object}	identity.AccountRequest	"demande en attente de l'administrateur (status pending)"
//	@Failure	400	{object}	map[string]string	"rôle déjà détenu ou absent de la table des rôles"
//	@Failure	403	{object}	map[string]string	"identité inactive"
//	@Failure	409	{object}	map[string]string	"une demande de rôle est déjà en attente"
//	@Failure	503	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/identity/roles/request [post]
func (h *Handler) handleIdentityRoleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	sess := sessionFromCtx(r)
	if sess == nil || sess.Handle == "" {
		jsonError(w, "session liée à une identité requise", http.StatusForbidden)
		return
	}
	var dto roleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		jsonError(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	wallet, err := h.identitySvc.Wallet(sess.Handle)
	if err != nil {
		h.internalErr(w, err)
		return
	}
	req, err := h.identitySvc.RequestRole(r.Context(), wallet, strings.TrimSpace(dto.Role))
	switch {
	case errors.Is(err, identity.ErrRoleAlreadyHeld), errors.Is(err, identity.ErrInvalidRole):
		jsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, identity.ErrRequestAlreadyPending):
		jsonError(w, err.Error(), http.StatusConflict)
	case errors.Is(err, identity.ErrIdentityPending), errors.Is(err, identity.ErrIdentitySuspended),
		errors.Is(err, identity.ErrIdentityInactive), errors.Is(err, identity.ErrIdentityWithoutRole):
		jsonError(w, err.Error(), http.StatusForbidden)
	case err != nil:
		h.internalErr(w, err)
	case req.Status == identity.RequestPending:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(req)
	default:
		jsonOK(w, req)
	}
}

// accountDTO est le corps attendu pour POST /api/identity/accounts.
type accountDTO struct {
	Pseudo string `json:"pseudo"`
	OrgID  string `json:"org_id"`
	Role   string `json:"role,omitempty"`
	Email  string `json:"email,omitempty"`
}

// handleIdentityAccounts crée, à l'initiative de l'administrateur, une identité
// active avec le rôle donné — seule voie d'inscription sur un réseau qui
// n'accepte pas l'auto-enregistrement (ADR-12).
//
//	@Summary	Créer un compte (administrateur)
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		accountDTO	true	"pseudo et org_id requis — role vide = reader"
//	@Success	201	{object}	map[string]string	"{handle, secret}"
//	@Failure	400	{object}	map[string]string	"pseudo, organisation ou rôle invalide (rôle absent de la table des rôles)"
//	@Failure	403	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/identity/accounts [post]
func (h *Handler) handleIdentityAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	var dto accountDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		jsonError(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	handle, secret, err := h.identitySvc.CreateAccount(r.Context(),
		strings.TrimSpace(dto.Pseudo), strings.TrimSpace(dto.OrgID),
		strings.TrimSpace(dto.Role), strings.TrimSpace(dto.Email))
	switch {
	case errors.Is(err, identity.ErrInvalidIdentityName), errors.Is(err, identity.ErrInvalidRole):
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		h.internalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, map[string]string{"handle": handle, "secret": secret})
}

// handleIdentityRequests liste les demandes de compte en attente — réservé aux
// identités portant la permission identity.admin (voir domain/role) : les
// demandes exposent des données personnelles (email, message libre) de tous
// les utilisateurs, jamais accessible à un compte non-admin.
//
//	@Summary	Lister les demandes de compte
//	@Tags		identity
//	@Produce	json
//	@Success	200	{array}		identity.AccountRequest
//	@Failure	403	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Security	MyrToken
//	@Router		/identity/requests [get]
func (h *Handler) handleIdentityRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	reqs, err := h.identitySvc.ListRequests()
	if err != nil {
		h.internalErr(w, err)
		return
	}
	if reqs == nil {
		reqs = []*identity.AccountRequest{}
	}
	jsonOK(w, reqs)
}

// handleIdentityPolicy retourne la politique d'accès du réseau actif — un
// client (GUI ou tout autre logiciel tiers) l'interroge au démarrage pour
// savoir s'il doit bloquer ou autoriser l'accès automatique.
//
//	@Summary	Politique d'accès du réseau actif
//	@Tags		identity
//	@Produce	json
//	@Success	200	{object}	networkPolicyResponse
//	@Router		/identity/policy [get]
func (h *Handler) handleIdentityPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.networkSvc == nil {
		jsonOK(w, networkPolicyResponse{Mode: "blockchain"})
		return
	}
	profile, err := h.networkSvc.GetActive()
	if err != nil || profile == nil {
		jsonOK(w, networkPolicyResponse{Mode: "blockchain"})
		return
	}
	jsonOK(w, networkPolicyResponse{
		Mode:              "blockchain",
		AllowAutoGuest:    profile.AllowAutoGuest,
		AllowAutoRegister: profile.AllowAutoRegister,
		RequireApproval:   profile.RequireApproval,
		NetworkName:       profile.Name,
	})
}

// handleIdentitySession crée une session à partir d'un secret d'enrollment
// (utilisateur déjà connu de la CA) et retourne un token opaque à inclure
// dans l'en-tête X-Myr-Token pour les requêtes suivantes.
//
//	@Summary	Créer une session (utilisateur connu)
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"name, secret, org_id requis ; channel optionnel"
//	@Success	200	{object}	sessionDTO
//	@Failure	400	{object}	map[string]string
//	@Failure	401	{object}	map[string]string
//	@Failure	403	{object}	map[string]string	"compte en attente d'approbation, suspendu, ou sans rôle dans son certificat"
//	@Failure	429	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Router		/identity/session [post]
func (h *Handler) handleIdentitySession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.authLimiter.allow(h.clientIP(r)) {
		jsonError(w, "trop de tentatives, réessayez dans une minute", http.StatusTooManyRequests)
		return
	}
	var req struct {
		Name    string `json:"name"`
		Secret  string `json:"secret"`
		OrgID   string `json:"org_id"`
		Channel string `json:"channel,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.OrgID = strings.TrimSpace(req.OrgID)
	req.Secret = strings.TrimSpace(req.Secret)
	req.Channel = strings.TrimSpace(req.Channel)
	if req.Name == "" || req.Secret == "" || req.OrgID == "" {
		jsonError(w, "name, secret et org_id sont requis", http.StatusBadRequest)
		return
	}
	if h.identitySvc == nil {
		jsonError(w, "service identité non configuré", http.StatusServiceUnavailable)
		return
	}
	entry, err := h.identitySvc.Enroll(r.Context(), req.Name, req.Secret, req.OrgID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Le certificat tout juste émis est la seule source de vérité : seule une
	// identité active et dotée d'un rôle ouvre une session (ADR-12).
	if err := entry.CanSignIn(); err != nil {
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	}
	channel := req.Channel
	if channel == "" {
		channel = h.netInfo.Channel // canal par défaut du réseau
	}
	// La session porte le rôle du certificat et retient l'identité enrôlée,
	// dont le wallet signera les écritures de l'utilisateur sur le ledger.
	sess, err := h.sessions.create(req.Name, entry.Handle, entry.Role, channel)
	if err != nil {
		jsonError(w, "erreur interne", http.StatusInternalServerError)
		return
	}
	jsonOK(w, sessionDTO{Token: sess.Token, Role: sess.Role, Pseudo: sess.Pseudo, Channel: sess.Channel})
}

// handleIdentityGuest indique si le réseau ouvre la lecture publique du
// catalogue (AllowAutoGuest, RM49) : 200 avec guest=true et un jeton vide — la
// lecture se fait sans jeton, aucune session n'est créée —, 403 sinon. Route
// conservée pour les clients qui demandaient un jeton invité ; GET
// /api/identity/policy donne la même information. Une inscription passe par
// POST /api/identity/request, selon les règles du réseau.
//
//	@Summary	Accès invité (lecture publique, sans session)
//	@Description	Déprécié : la lecture publique du catalogue ne demande aucun jeton. Répond 200 {guest: true, token: ""} si le réseau l'autorise, 403 sinon.
//	@Tags		identity
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	false	"pseudo, display_name, email, org_id, message — tous optionnels"
//	@Success	200	{object}	sessionDTO
//	@Failure	403	{object}	map[string]string
//	@Failure	429	{object}	map[string]string
//	@Router		/identity/guest [post]
func (h *Handler) handleIdentityGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.authLimiter.allow(h.clientIP(r)) {
		jsonError(w, "trop de tentatives, réessayez dans une minute", http.StatusTooManyRequests)
		return
	}
	var body struct {
		Pseudo      string `json:"pseudo,omitempty"`
		DisplayName string `json:"display_name,omitempty"`
		Email       string `json:"email,omitempty"`
		OrgID       string `json:"org_id,omitempty"`
		Message     string `json:"message,omitempty"`
	}
	// Corps optionnel — un accès guest pur n'en a pas besoin
	json.NewDecoder(r.Body).Decode(&body) //nolint:errcheck

	if h.publicRead() {
		// Lecture publique : aucune session, le visiteur lit sans jeton (RM49).
		pseudo := strings.TrimSpace(body.Pseudo)
		if pseudo == "" {
			pseudo = "guest"
		}
		jsonOK(w, sessionDTO{Role: rbac.NameReader, Pseudo: pseudo, Channel: h.netInfo.Channel, Guest: true})
		return
	}

	// 403 : le réseau n'autorise pas l'accès automatique
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{ //nolint:errcheck
		"error": "ce réseau ne permet pas l'accès automatique — contactez l'administrateur",
	})
}
