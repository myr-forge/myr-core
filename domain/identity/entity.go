// domain/identity/entity.go
package identity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"myr-core/domain/network"
)

// MyrIdentity représente l'identité d'un auteur dans le réseau Myr.
// Les attributs Myr.status, Myr.role et Myr.channels sont écrits dans
// le certificat X.509 par la Fabric CA.
type MyrIdentity struct {
	ID             string    `json:"id"`
	DisplayName    string    `json:"displayName"`
	LegalName      string    `json:"legalName"`
	Email          string    `json:"email"`
	Country        string    `json:"country"`
	Organization   string    `json:"organization"`
	LicenseDefault string    `json:"licenseDefault"`
	IPAgreement    bool      `json:"ipAgreement"`
	IPAgreedAt     time.Time `json:"ipAgreedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	Status         string    `json:"status"`
	Role           string    `json:"role"`
}

// WalletEntry décrit un wallet local situé dans ~/.Myr/wallets/<handle>/msp/
type WalletEntry struct {
	Handle string // "<pseudo>@<org>" — identifiant unique du wallet
	Name   string // pseudo affiché
	OrgID  string // Org1MSP, Org2MSP, Org3MSP
	Status string // attribut Myr.status du certificat (pending | active | suspended) — vide si absent ou illisible
	Role   string // attribut Myr.role du certificat — vide si absent ou illisible
	MSPDir string // chemin absolu vers le répertoire MSP
}

// CanSignIn vérifie qu'une identité peut ouvrir une session (ADR-12) : son
// certificat la déclare active et lui attribue un rôle. Une identité en
// attente d'approbation, suspendue, sans statut ou sans rôle est refusée.
func (w WalletEntry) CanSignIn() error {
	switch w.Status {
	case StatusActive:
	case StatusPending:
		return ErrIdentityPending
	case StatusSuspended:
		return ErrIdentitySuspended
	default:
		return fmt.Errorf("%w (Myr.status %q)", ErrIdentityInactive, w.Status)
	}
	if w.Role == "" {
		return ErrIdentityWithoutRole
	}
	return nil
}

// Disposition d'un wallet sous son répertoire MSP — seule définition, partagée
// par l'écriture (saveWallet) et par les adapters qui signent avec le wallet.
const (
	walletCertFile   = "signcerts/cert.pem"
	walletKeyFile    = "keystore/key.pem"
	walletCACertFile = "cacerts/ca.pem"
)

// CertPath retourne le chemin du certificat de signature du wallet.
func (w WalletEntry) CertPath() string {
	return filepath.Join(w.MSPDir, filepath.FromSlash(walletCertFile))
}

// KeyPath retourne le chemin de la clé privée du wallet.
func (w WalletEntry) KeyPath() string {
	return filepath.Join(w.MSPDir, filepath.FromSlash(walletKeyFile))
}

// Handle construit l'identifiant d'une identité, « pseudo@org » : c'est à la
// fois son identifiant d'enrôlement auprès de la CA, le nom de son wallet et
// le propriétaire (OwnerID) des assets qu'elle soumet (ADR-12). orgID est le
// MSP (Org1MSP) ou l'organisation seule (Org1).
func Handle(name, orgID string) string {
	return name + "@" + strings.TrimSuffix(orgID, "MSP")
}

// RegisterRequest contient les informations pour enregistrer une nouvelle identité
// auprès de la Fabric CA.
type RegisterRequest struct {
	Name        string // pseudo@org — identifiant CA (ex: alice@Org1)
	OrgID       string // ex: Org1MSP
	Password    string // secret d'enrollment (jamais stocké)
	Role        string // ex: RoleReader, RoleContributor — vide = RoleReader par défaut
	Status      string // statut à l'enregistrement (StatusActive, StatusPending) — vide = StatusPending
	DisplayName string
	LegalName   string
	Email       string
	Country     string
}

// RegistrationPolicy décrit les règles d'inscription d'un réseau, fixées par
// son administrateur (ADR-12).
type RegistrationPolicy struct {
	// AllowAutoRegister : sans auto-enregistrement, aucune demande
	// d'inscription n'est acceptée — seul l'administrateur crée des comptes.
	AllowAutoRegister bool
	// RequireApproval : l'identité est créée en attente (StatusPending) et ne
	// devient active qu'à l'approbation par l'administrateur ; sinon elle est
	// créée active.
	RequireApproval bool
	// DefaultRole est le rôle attribué à l'inscription — vide = RoleReader (RM21).
	DefaultRole string
}

// PolicyOf retourne les règles d'inscription fixées par l'administrateur du
// réseau p — sans réseau, l'inscription est fermée.
func PolicyOf(p *network.NetworkProfile) RegistrationPolicy {
	if p == nil {
		return RegistrationPolicy{}
	}
	return RegistrationPolicy{
		AllowAutoRegister: p.AllowAutoRegister,
		RequireApproval:   p.RequireApproval,
		DefaultRole:       p.AutoRegisterRole,
	}
}

// Attributs d'identité portés par le certificat émis par la CA.
const (
	AttrRole   = "Myr.role"
	AttrStatus = "Myr.status"
)

// Constantes de statut Fabric CA
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusSuspended = "suspended"
)

// Rôles Fabric CA
const (
	RoleReader      = "reader"
	RoleContributor = "contributor"
	RoleAuditor     = "auditor"
	RoleAdmin       = "admin"
)

// AccountRequest est une demande d'inscription d'un utilisateur qui n'a pas encore
// de compte, acceptée seulement si le réseau autorise l'auto-enregistrement
// (RegistrationPolicy). L'identité est créée dans la CA dès la demande ; si le
// réseau exige une approbation, la demande reste en attente jusqu'à la décision
// de l'administrateur (ApproveRequest, RejectRequest).
//
// Une demande de rôle (Kind = RequestKindRole, UCA08) suit le même cycle pour
// une identité déjà inscrite : Role est le rôle demandé, Pseudo et OrgID ceux
// de l'identité qui le demande.
type AccountRequest struct {
	ID          string `json:"id"`
	Kind        string `json:"kind,omitempty"` // vide : inscription ; RequestKindRole : demande de rôle
	Pseudo      string `json:"pseudo"`         // identifiant souhaité
	DisplayName string `json:"display_name"`   // nom affiché
	Email       string `json:"email"`
	OrgID       string `json:"org_id"`            // organisation souhaitée (ex: Org1MSP)
	Role        string `json:"role,omitempty"`    // rôle demandé (demande de rôle uniquement)
	Message     string `json:"message,omitempty"` // message libre pour l'admin
	Status      string `json:"status"`            // "pending" | "approved" | "rejected" | "failed"
	CreatedAt   string `json:"created_at"`        // ISO 8601
}

// RequestKindRole identifie une demande de rôle (AccountRequest.Kind).
const RequestKindRole = "role"

// Statuts d'une demande de compte
const (
	RequestPending  = "pending"
	RequestApproved = "approved"
	RequestRejected = "rejected"
	// RequestFailed : l'identité n'a pas pu être créée dans la CA ; la demande
	// est conservée pour trace, sans identité associée.
	RequestFailed = "failed"
)
