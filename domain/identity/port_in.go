// domain/identity/port_in.go
package identity

import "context"

// IdentityService est le port d'entrée du domaine identity.
// Implémenté par le service, consommé par les adapters in (CLI, REST).
type IdentityService interface {
	ListLocalWallets() ([]WalletEntry, error)
	// Wallet retourne le wallet local de l'identité handle (« pseudo@org »,
	// voir Handle) — celui qui signe ses écritures sur le ledger (ADR-12).
	// Retourne ErrWalletNotFound si l'identité n'a pas de wallet sur ce nœud.
	Wallet(handle string) (WalletEntry, error)
	// CreateAccount crée, à l'initiative de l'admin, une identité active avec
	// le rôle donné — seule voie d'inscription sur un réseau fermé à
	// l'auto-enregistrement. Retourne l'identifiant pseudo@org et le secret.
	CreateAccount(ctx context.Context, pseudo, orgID, role, email string) (handle, secret string, err error)
	Enroll(ctx context.Context, name, secret, orgID string) (WalletEntry, error)
	GetStatus(ctx context.Context, wallet WalletEntry) (string, error)
	ReEnroll(ctx context.Context, wallet WalletEntry) (WalletEntry, error)
	// LoadGuestWallet charge un wallet invité (lecture seule) depuis des fichiers
	// pré-enrôlés par l'admin, sans passer par la CA.
	LoadGuestWallet(certPath, keyPath, caCertPath, orgID string) (WalletEntry, error)
	// RequestAccess traite une demande d'inscription selon les règles du
	// réseau (policy) : refusée si le réseau n'accepte pas l'auto-enregistrement
	// (ErrRegistrationClosed), sinon l'identité est créée dans la CA, active ou
	// en attente d'approbation. Retourne la demande et le secret d'enrôlement.
	RequestAccess(ctx context.Context, req AccountRequest, policy RegistrationPolicy) (*AccountRequest, string, error)
	// RequestRole traite la demande du rôle wanted par l'identité du wallet w
	// (UCA08, RM47) : attribué aussitôt si le rôle est en attribution
	// automatique (demande approuvée), sinon en attente de l'administrateur.
	RequestRole(ctx context.Context, w WalletEntry, wanted string) (*AccountRequest, error)
	// ApproveRequest approuve une demande en attente : l'identité d'une
	// inscription devient active, celle d'une demande de rôle reçoit ce rôle.
	ApproveRequest(ctx context.Context, id string) (*AccountRequest, error)
	// RejectRequest refuse une demande en attente : l'identité d'une
	// inscription est suspendue, celle d'une demande de rôle garde son rôle.
	RejectRequest(ctx context.Context, id string) (*AccountRequest, error)
	// ListRequests retourne toutes les demandes enregistrées (usage admin).
	ListRequests() ([]*AccountRequest, error)
	WalletDir() string
	// SetRole change le rôle enregistré d'une identité existante auprès de la CA.
	// Le nouveau rôle ne s'applique qu'aux prochains certificats émis — un
	// ré-enrôlement (ReEnroll) est nécessaire pour l'appliquer immédiatement.
	SetRole(ctx context.Context, name, newRole string) error
	// SetStatus change le statut (active ou suspended) d'une identité existante
	// auprès de la CA, avec la même prise d'effet que SetRole.
	SetStatus(ctx context.Context, name, status string) error
}
