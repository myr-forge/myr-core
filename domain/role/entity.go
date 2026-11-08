// domain/role/entity.go
package role

import "time"

// Permission identifie une action protégée (routes REST / commandes CLI).
// Catalogue dérivé du découpage déjà existant dans adapters/in/rest/server.go
// (groupes contrib/auth/adminOnly) — pas une granularité par endpoint individuel.
type Permission string

const (
	// PermRead couvre la lecture authentifiée (GET /api/* pour un utilisateur connecté).
	PermRead Permission = "read"
	// PermWrite couvre l'écriture des assets (components, connections, modules,
	// interfaces, licences — le groupe "contributor" existant de server.go).
	PermWrite Permission = "write"
	// PermNetworkAdmin couvre la création/gestion de réseaux, organisations et nœuds
	// (myr network/org/node — UCADM01-05).
	PermNetworkAdmin Permission = "network.admin"
	// PermRoleAdmin couvre la gestion des rôles eux-mêmes (myr role ...).
	PermRoleAdmin Permission = "role.admin"
	// PermIdentityAdmin couvre la gestion des identités/comptes (myr identity set-role,
	// validation des demandes de compte).
	PermIdentityAdmin Permission = "identity.admin"
	// PermAdmin couvre les opérations globales (/api/admin/*).
	PermAdmin Permission = "admin"
)

// AllPermissions liste le catalogue complet, utilisé pour valider les entrées CLI.
var AllPermissions = []Permission{PermRead, PermWrite, PermNetworkAdmin, PermRoleAdmin, PermIdentityAdmin, PermAdmin}

// adminPermissions sont les permissions d'administration, jamais accordées
// automatiquement — ni par une inscription, ni par une demande de rôle.
var adminPermissions = []Permission{PermNetworkAdmin, PermRoleAdmin, PermIdentityAdmin, PermAdmin}

// IsValidPermission vérifie qu'une permission appartient au catalogue connu.
func IsValidPermission(p string) bool {
	for _, known := range AllPermissions {
		if string(known) == p {
			return true
		}
	}
	return false
}

// Role représente un rôle attribuable à une identité, avec un ensemble de
// permissions défini par l'administrateur du réseau (RBAC dynamique).
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"` // ex: "reader", "contributor", "consumer"
	Permissions []Permission `json:"permissions"`
	// AutoGrant = true : une identité qui demande ce rôle l'obtient sans
	// approbation de l'administrateur (UCA08) ; jamais pour un rôle portant
	// une permission d'administration (ErrAdminAutomatic).
	AutoGrant bool `json:"auto_grant,omitempty"`
	// BuiltIn = true pour les 4 rôles de la table initiale (DefaultRoles) —
	// modifiables et supprimables comme les autres, sauf admin, immuable.
	BuiltIn   bool      `json:"built_in"`
	CreatedAt time.Time `json:"created_at"`
}

// Has indique si le rôle porte la permission donnée.
func (r *Role) Has(p Permission) bool {
	for _, existing := range r.Permissions {
		if existing == p {
			return true
		}
	}
	return false
}

// hasAdminPermission indique si le rôle porte au moins une permission
// d'administration.
func (r *Role) hasAdminPermission() (Permission, bool) {
	for _, p := range adminPermissions {
		if r.Has(p) {
			return p, true
		}
	}
	return "", false
}

// Noms des 4 rôles de la table initiale (domain/identity.Role{Reader,Contributor,Auditor,Admin}).
const (
	NameReader      = "reader"
	NameContributor = "contributor"
	NameAuditor     = "auditor"
	NameAdmin       = "admin"
)

// defaultPermissions définit la table initiale du réseau : celle qui
// s'applique tant que l'administrateur n'en a publié aucune. Le chaincode
// (chaincode/model/roles.go) en porte une copie, qui doit rester identique.
var defaultPermissions = map[string][]Permission{
	NameReader:      {PermRead},
	NameContributor: {PermRead, PermWrite},
	NameAuditor:     {PermRead},
	NameAdmin:       AllPermissions,
}

// DefaultRoles retourne la table initiale du réseau, triée par nom.
func DefaultRoles() []*Role {
	names := []string{NameAdmin, NameAuditor, NameContributor, NameReader}
	out := make([]*Role, 0, len(names))
	for _, name := range names {
		out = append(out, defaultRole(name))
	}
	return out
}

// defaultRole retourne le rôle name de la table initiale, nil s'il n'en fait
// pas partie.
func defaultRole(name string) *Role {
	perms, ok := defaultPermissions[name]
	if !ok {
		return nil
	}
	return &Role{ID: name, Name: name, Permissions: append([]Permission(nil), perms...), BuiltIn: true}
}
