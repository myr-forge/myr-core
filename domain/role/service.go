// domain/role/service.go
package role

import (
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

// Service applique la table des rôles publiée par le réseau (ADR-12). Tant
// qu'aucune n'est publiée, la table initiale (DefaultRoles) s'applique ; si
// la table devient illisible (réseau injoignable), la dernière table lue
// s'applique, à défaut la table initiale.
type Service struct {
	repo Repo // nil = aucun réseau : seule la table initiale s'applique

	// refresh est la durée pendant laquelle la table lue est réutilisée sans
	// relire Repo — 0 relit à chaque appel.
	refresh time.Duration

	mu     sync.Mutex
	known  []*Role // dernière table lue, nil avant la première lecture
	readAt time.Time
}

// NewService construit le service sur la table publiée par repo (nil : table
// initiale seule, aucune modification possible).
func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

// WithRefresh réutilise la table lue pendant d pour éviter une lecture du
// ledger à chaque contrôle d'accès : une modification publiée par un autre
// nœud s'applique ici au plus tard après d.
func (s *Service) WithRefresh(d time.Duration) *Service {
	s.refresh = d
	return s
}

// table retourne la table effective (voir Service).
func (s *Service) table() []*Role {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.known != nil && s.refresh > 0 && time.Since(s.readAt) < s.refresh {
		return s.known
	}
	if s.repo != nil {
		if roles, err := s.repo.FindAll(); err == nil {
			s.remember(roles)
			return s.known
		}
	}
	if s.known == nil {
		s.remember(nil)
	}
	// Table illisible : la prochaine tentative attend refresh, pour ne pas
	// subir l'échec de lecture à chaque contrôle d'accès.
	s.readAt = time.Now()
	return s.known
}

// remember retient roles comme dernière table lue : la table initiale si
// roles est vide, admin toujours rétabli dans sa définition immuable, triée
// par nom. L'appelant détient s.mu.
func (s *Service) remember(roles []*Role) {
	if len(roles) == 0 {
		roles = DefaultRoles()
	}
	table := make([]*Role, 0, len(roles)+1)
	for _, r := range roles {
		if r != nil && r.Name != NameAdmin {
			table = append(table, r)
		}
	}
	table = append(table, defaultRole(NameAdmin))
	sort.Slice(table, func(i, j int) bool { return table[i].Name < table[j].Name })
	s.known, s.readAt = table, time.Now()
}

// find retourne le rôle name de la table effective, nil s'il en est absent.
func (s *Service) find(name string) *Role {
	for _, r := range s.table() {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// publish applique change à la table publiée puis invalide la table lue, pour
// que la lecture suivante reflète la modification.
func (s *Service) publish(change func(Repo) error) error {
	if s.repo == nil {
		return ErrNoTable
	}
	if err := change(s.repo); err != nil {
		return err
	}
	s.mu.Lock()
	s.readAt = time.Time{}
	s.mu.Unlock()
	return nil
}

// Create ajoute un rôle à la table (ex: consumer, manufacturer — cf. Securite.md §5).
func (s *Service) Create(name string, permissions []Permission, autoGrant bool) (*Role, error) {
	if !namePattern.MatchString(name) {
		return nil, fmt.Errorf("%w: %q — attendu [a-z][a-z0-9_-]{1,63}", ErrInvalidName, name)
	}
	if s.find(name) != nil {
		return nil, fmt.Errorf("%w: %q", ErrAlreadyExists, name)
	}
	r := &Role{ID: name, Name: name, Permissions: permissions, AutoGrant: autoGrant, CreatedAt: time.Now()}
	if err := validate(r); err != nil {
		return nil, err
	}
	if err := s.publish(func(repo Repo) error { return repo.Save(r) }); err != nil {
		return nil, err
	}
	return r, nil
}

// Update modifie un rôle existant, sauf admin, immuable.
func (s *Service) Update(name string, change Change) (*Role, error) {
	if name == NameAdmin {
		return nil, ErrAdminImmutable
	}
	existing := s.find(name)
	if existing == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	r := *existing
	if change.Permissions != nil {
		r.Permissions = *change.Permissions
	}
	if change.AutoGrant != nil {
		r.AutoGrant = *change.AutoGrant
	}
	if err := validate(&r); err != nil {
		return nil, err
	}
	if err := s.publish(func(repo Repo) error { return repo.Save(&r) }); err != nil {
		return nil, err
	}
	return &r, nil
}

// Delete retire un rôle de la table, sauf admin, immuable.
func (s *Service) Delete(name string) error {
	if name == NameAdmin {
		return ErrAdminImmutable
	}
	if s.find(name) == nil {
		return fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	return s.publish(func(repo Repo) error { return repo.Delete(name) })
}

func (s *Service) Get(name string) (*Role, error) {
	r := s.find(name)
	if r == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	return r, nil
}

func (s *Service) List() ([]*Role, error) {
	return s.table(), nil
}

// HasPermission indique si le rôle nommé porte la permission donnée.
// Un rôle inconnu ne porte aucune permission (fail-closed).
func (s *Service) HasPermission(roleName string, p Permission) bool {
	r := s.find(roleName)
	return r != nil && r.Has(p)
}

// Assignable vérifie qu'un rôle peut être attribué à une identité : il doit
// figurer dans la table des rôles — un rôle inconnu est refusé, jamais remplacé
// par un autre (RM46). Sans table (roles nil), seuls les rôles de la table
// initiale sont connus.
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
	if p, ok := r.hasAdminPermission(); ok {
		return fmt.Errorf("%w : %q porte %s", ErrAdminAutomatic, name, p)
	}
	return nil
}

// GrantedOnRequest indique si le rôle name, demandé par une identité, lui est
// attribué sans approbation de l'administrateur (RM47) : il doit figurer dans
// la table (RM46) et y être en attribution automatique. Un rôle d'administration
// n'est jamais attribué automatiquement, même si la table publiée l'indique.
func GrantedOnRequest(roles RoleService, name string) (bool, error) {
	r, err := lookup(roles, name)
	if err != nil {
		return false, err
	}
	if _, admin := r.hasAdminPermission(); admin {
		return false, nil
	}
	return r.AutoGrant, nil
}

// lookup retourne le rôle name de la table roles, ou de la table initiale si
// roles est nil.
func lookup(roles RoleService, name string) (*Role, error) {
	if roles != nil {
		return roles.Get(name)
	}
	if r := defaultRole(name); r != nil {
		return r, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrNotFound, name)
}

// validate vérifie les permissions d'un rôle et qu'un rôle en attribution
// automatique ne porte aucune permission d'administration (RM47).
func validate(r *Role) error {
	for _, p := range r.Permissions {
		if !IsValidPermission(string(p)) {
			return fmt.Errorf("%w: %q", ErrInvalidPermission, p)
		}
	}
	if p, ok := r.hasAdminPermission(); ok && r.AutoGrant {
		return fmt.Errorf("%w : %q porte %s", ErrAdminAutomatic, r.Name, p)
	}
	return nil
}
