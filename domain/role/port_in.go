// domain/role/port_in.go — port d'entrée du domaine role
package role

// RoleService est le port d'entrée du domaine role (RBAC dynamique).
type RoleService interface {
	// Create ajoute un rôle à la table du réseau, avec l'ensemble de
	// permissions donné ; autoGrant le rend attribuable sur simple demande
	// (UCA08), jamais s'il porte une permission d'administration.
	Create(name string, permissions []Permission, autoGrant bool) (*Role, error)
	// Update modifie un rôle existant : les champs nil de change restent
	// inchangés. Refuse admin, immuable (ErrAdminImmutable).
	Update(name string, change Change) (*Role, error)
	// Delete retire un rôle de la table. Refuse admin, immuable
	// (ErrAdminImmutable). Une identité qui porte encore ce rôle n'a plus
	// aucune permission.
	Delete(name string) error
	// Get retourne un rôle par son nom.
	Get(name string) (*Role, error)
	// List retourne la table des rôles du réseau.
	List() ([]*Role, error)
	// HasPermission indique si le rôle nommé porte la permission donnée.
	// Un rôle inconnu ne porte aucune permission (retourne false, pas d'erreur) —
	// utilisé directement dans les points de contrôle d'accès (REST requireRole).
	HasPermission(roleName string, p Permission) bool
}

// Change décrit la modification d'un rôle : un champ nil reste inchangé.
type Change struct {
	Permissions *[]Permission
	AutoGrant   *bool
}
