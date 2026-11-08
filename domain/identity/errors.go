// domain/identity/errors.go
package identity

import "errors"

// ErrNoAdminCredentials est retourné quand l'opération CA nécessite les
// identifiants administrateur (cert + clé) qui ne sont pas configurés dans
// le profil réseau.
//
// Pour corriger : éditez le réseau (touche e dans le menu) et renseignez
// les champs "Cert path" et "Key path" avec le certificat et la clé de
// l'administrateur CA.
var ErrNoAdminCredentials = errors.New(
	"l'enregistrement CA requiert les identifiants admin — " +
		"éditez le réseau (e) et renseignez Cert path + Key path",
)

// ErrWalletNotFound est retourné quand aucun wallet local n'existe pour une
// identité — elle n'a jamais été enrôlée sur ce nœud.
var ErrWalletNotFound = errors.New("identity: aucun wallet local pour cette identité")

// ErrRegistrationClosed est retourné quand le réseau n'autorise pas
// l'auto-enregistrement : seul l'administrateur crée des comptes (ADR-12).
var ErrRegistrationClosed = errors.New("identity: ce réseau n'accepte pas les demandes d'inscription — contactez l'administrateur")

// ErrInvalidIdentityName est retourné quand un pseudo ou une organisation ne
// peut pas former un identifiant « pseudo@org ».
var ErrInvalidIdentityName = errors.New("identity: pseudo ou organisation invalide")

// ErrInvalidRole est retourné quand le rôle à attribuer n'existe pas dans la
// table des rôles, ou porte des droits d'administration alors qu'il serait
// attribué par défaut à une inscription (RM46).
var ErrInvalidRole = errors.New("identity: rôle non attribuable")

// Refus d'ouverture de session (WalletEntry.CanSignIn, ADR-12).
var (
	ErrIdentityPending     = errors.New("identity: compte en attente d'approbation par l'administrateur")
	ErrIdentitySuspended   = errors.New("identity: compte suspendu")
	ErrIdentityInactive    = errors.New("identity: compte inactif")
	ErrIdentityWithoutRole = errors.New("identity: aucun rôle attribué à cette identité")
)

// Décision sur une demande d'inscription ou de rôle (ApproveRequest, RejectRequest).
var (
	ErrRequestNotFound   = errors.New("identity: demande introuvable")
	ErrRequestNotPending = errors.New("identity: la demande n'est plus en attente")
)

// Demande de rôle (RequestRole, UCA08).
var (
	ErrRoleAlreadyHeld       = errors.New("identity: vous possédez déjà ce rôle")
	ErrRequestAlreadyPending = errors.New("identity: une demande de rôle est déjà en attente pour cette identité")
)

// IsNoAdminCredentials rapporte si err est ou enveloppe ErrNoAdminCredentials.
func IsNoAdminCredentials(err error) bool {
	return errors.Is(err, ErrNoAdminCredentials)
}
