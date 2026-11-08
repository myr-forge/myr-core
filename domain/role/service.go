// domain/role/service.go
package role

import (
	"fmt"
	"regexp"
	"time"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

// builtinPermissions définit les permissions par défaut des 4 rôles historiques
// (domain/identity.Role{Reader,Contributor,Auditor,Admin}) — équivalentes au
// comportement REST déjà en place (adapters/in/rest/handlers.go requireRole),
// pour ne provoquer aucune régression à l'introduction du RBAC dynamique.
var builtinPermissions = map[string][]Permission{
	NameReader:      {PermRead},
	NameContributor: {PermRead, PermWrite},
	NameAuditor:     {PermRead},
	NameAdmin:       {PermRead, PermWrite, PermNetworkAdmin, PermRoleAdmin, PermIdentityAdmin, PermAdmin},
}

type Service struct {
	repo Repo
}

// NewService construit le service et amorce les 4 rôles intégrés s'ils sont absents
// du store (idempotent — sans effet si déjà présents, ex. redémarrages successifs).
func NewService(repo Repo) *Service {
	s := &Service{repo: repo}
	s.seedBuiltins()
	return s
}

func (s *Service) seedBuiltins() {
	for name, perms := range builtinPermissions {
		existing, _ := s.repo.FindByName(name)
		if existing != nil {
			continue
		}
		_ = s.repo.Save(&Role{ID: name, Name: name, Permissions: perms, BuiltIn: true, CreatedAt: time.Now()})
	}
}

// Create crée un rôle personnalisé (ex: consumer, manufacturer — cf. Securite.md §5).
func (s *Service) Create(name string, permissions []Permission) (*Role, error) {
	if !namePattern.MatchString(name) {
		return nil, fmt.Errorf("%w: %q — attendu [a-z][a-z0-9_-]{1,63}", ErrInvalidName, name)
	}
	if err := validatePermissions(permissions); err != nil {
		return nil, err
	}
	existing, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: %q", ErrAlreadyExists, name)
	}
	r := &Role{ID: name, Name: name, Permissions: permissions, BuiltIn: false, CreatedAt: time.Now()}
	if err := s.repo.Save(r); err != nil {
		return nil, err
	}
	return r, nil
}

// Update remplace l'ensemble de permissions d'un rôle non intégré.
func (s *Service) Update(name string, permissions []Permission) (*Role, error) {
	if err := validatePermissions(permissions); err != nil {
		return nil, err
	}
	existing, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	if existing.BuiltIn {
		return nil, ErrBuiltIn
	}
	existing.Permissions = permissions
	if err := s.repo.Save(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete supprime un rôle personnalisé. Refuse les 4 rôles intégrés.
func (s *Service) Delete(name string) error {
	existing, err := s.repo.FindByName(name)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	if existing.BuiltIn {
		return ErrBuiltIn
	}
	return s.repo.Delete(name)
}

func (s *Service) Get(name string) (*Role, error) {
	r, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	return r, nil
}

func (s *Service) List() ([]*Role, error) {
	return s.repo.FindAll()
}

// HasPermission indique si le rôle nommé porte la permission donnée.
// Un rôle inconnu ne porte aucune permission (fail-closed).
func (s *Service) HasPermission(roleName string, p Permission) bool {
	r, err := s.repo.FindByName(roleName)
	if err != nil || r == nil {
		return false
	}
	return r.Has(p)
}

// adminPermissions sont les permissions d'administration, jamais accordées par
// une inscription automatique.
var adminPermissions = []Permission{PermNetworkAdmin, PermRoleAdmin, PermIdentityAdmin, PermAdmin}

// Assignable vérifie qu'un rôle peut être attribué à une identité : il doit
// figurer dans la table des rôles — un rôle inconnu est refusé, jamais remplacé
// par un autre (RM46). Sans table (roles nil), seuls les 4 rôles intégrés sont
// connus.
func Assignable(roles RoleService, name string) error {
	_, err := lookup(roles, name)
	return err
}

// AssignableByDefault vérifie qu'un rôle peut être attribué d'office à toute
// identité auto-enregistrée : il est attribuable et ne porte aucune permission
// d'administration — une inscription ne donne jamais de droit d'administration
// (RM46).
func AssignableByDefault(roles RoleService, name string) error {
	r, err := lookup(roles, name)
	if err != nil {
		return err
	}
	for _, p := range adminPermissions {
		if r.Has(p) {
			return fmt.Errorf("%w : %q porte %s", ErrAdminByDefault, name, p)
		}
	}
	return nil
}

// lookup retourne le rôle name de la table roles, ou du jeu des rôles intégrés
// si roles est nil.
func lookup(roles RoleService, name string) (*Role, error) {
	if roles != nil {
		return roles.Get(name)
	}
	perms, ok := builtinPermissions[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	return &Role{ID: name, Name: name, Permissions: perms, BuiltIn: true}, nil
}

func validatePermissions(perms []Permission) error {
	for _, p := range perms {
		if !IsValidPermission(string(p)) {
			return fmt.Errorf("%w: %q", ErrInvalidPermission, p)
		}
	}
	return nil
}
