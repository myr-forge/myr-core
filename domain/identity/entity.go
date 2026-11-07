// domain/identity/entity.go
package identity

import (
	"path/filepath"
	"strings"
	"time"
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
	Status string // pending | active | suspended (lu depuis le cert)
	MSPDir string // chemin absolu vers le répertoire MSP
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
	DisplayName string
	LegalName   string
	Email       string
	Country     string
}

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

// AccountRequest est une demande d'accès soumise par un utilisateur qui n'a pas encore
// de compte. Elle est stockée localement et traitée par l'administrateur.
// L'admin crée ensuite l'identité dans la CA Fabric et transmet le secret d'enrollment.
type AccountRequest struct {
	ID          string `json:"id"`
	Pseudo      string `json:"pseudo"`       // identifiant souhaité
	DisplayName string `json:"display_name"` // nom affiché
	Email       string `json:"email"`
	OrgID       string `json:"org_id"`            // organisation souhaitée (ex: Org1MSP)
	Message     string `json:"message,omitempty"` // message libre pour l'admin
	Status      string `json:"status"`            // "pending" | "approved" | "rejected"
	CreatedAt   string `json:"created_at"`        // ISO 8601
}

// Statuts d'une demande de compte
const (
	RequestPending  = "pending"
	RequestApproved = "approved"
	RequestRejected = "rejected"
)
