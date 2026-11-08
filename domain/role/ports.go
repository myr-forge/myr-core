// domain/role/ports.go — port de sortie du domaine role
package role

// Repo publie la table des rôles du réseau, partagée par tous ses nœuds
// (ADR-12) : chaque nœud applique la même table.
type Repo interface {
	// FindAll retourne la table publiée — vide si aucune ne l'a encore été.
	FindAll() ([]*Role, error)
	// Save crée ou remplace un rôle. La première publication part de la table
	// initiale (DefaultRoles) : les rôles qu'elle ne touche pas y restent.
	Save(r *Role) error
	// Delete retire un rôle de la table, avec la même règle de première
	// publication que Save.
	Delete(name string) error
}
