// domain/role/tests/service_test.go — tests unitaires du service RBAC
package role_test

import (
	"errors"
	"testing"
	"time"

	"myr-core/domain/role"
)

// ── Mock Repo (Pattern A — struct concrète, état en mémoire) ─────────────────

// mockRepo publie une table en mémoire, comme le ledger : vide tant que rien
// n'est publié, partant de la table initiale à la première publication.
type mockRepo struct {
	roles   map[string]*role.Role // nil = aucune table publiée
	readErr error                 // erreur retournée par FindAll (réseau injoignable)
}

func newMockRepo() *mockRepo { return &mockRepo{} }

func (m *mockRepo) published() map[string]*role.Role {
	if m.roles == nil {
		m.roles = make(map[string]*role.Role)
		for _, r := range role.DefaultRoles() {
			m.roles[r.Name] = r
		}
	}
	return m.roles
}

func (m *mockRepo) Save(r *role.Role) error {
	m.published()[r.Name] = r
	return nil
}

func (m *mockRepo) FindAll() ([]*role.Role, error) {
	if m.readErr != nil {
		return nil, m.readErr
	}
	out := make([]*role.Role, 0, len(m.roles))
	for _, r := range m.roles {
		out = append(out, r)
	}
	return out, nil
}

func (m *mockRepo) Delete(name string) error {
	delete(m.published(), name)
	return nil
}

func perms(p ...role.Permission) *[]role.Permission { return &p }
func flag(b bool) *bool                             { return &b }

// ── Table effective ───────────────────────────────────────────────────────────

// / @brief  Sans table publiée, la table initiale s'applique (RM48)
// / @input  mockRepo vide
// / @expect List retourne admin, auditor, contributor, reader, tous BuiltIn=true
func TestTable_RM48_NothingPublished_DefaultRoles(t *testing.T) {
	svc := role.NewService(newMockRepo())
	all, _ := svc.List()
	want := []string{role.NameAdmin, role.NameAuditor, role.NameContributor, role.NameReader}
	if len(all) != len(want) {
		t.Fatalf("got %d rôles, want %d", len(all), len(want))
	}
	for i, r := range all {
		if r.Name != want[i] || !r.BuiltIn {
			t.Errorf("rôle %d : got %q (BuiltIn=%v), want %q intégré", i, r.Name, r.BuiltIn, want[i])
		}
	}
}

// / @brief  Le rôle admin garde toutes les permissions, quelle que soit la table publiée (RM48)
// / @input  table publiée où admin ne porte que read
// / @expect admin porte network.admin et role.admin
func TestTable_RM48_AdminAlwaysImmutable(t *testing.T) {
	repo := newMockRepo()
	repo.published()[role.NameAdmin] = &role.Role{Name: role.NameAdmin, Permissions: []role.Permission{role.PermRead}}
	svc := role.NewService(repo)
	for _, p := range []role.Permission{role.PermNetworkAdmin, role.PermRoleAdmin} {
		if !svc.HasPermission(role.NameAdmin, p) {
			t.Errorf("admin devrait porter %s", p)
		}
	}
}

// / @brief  Table illisible : la dernière table lue s'applique (RM48)
// / @input  table publiée avec "consumer", lue une fois, puis réseau injoignable
// / @expect consumer reste connu et garde ses permissions
func TestTable_RM48_Unreachable_LastKnownTable(t *testing.T) {
	repo := newMockRepo()
	repo.Save(&role.Role{Name: "consumer", Permissions: []role.Permission{role.PermRead}})
	svc := role.NewService(repo)
	if !svc.HasPermission("consumer", role.PermRead) {
		t.Fatal("consumer devrait porter read")
	}
	repo.readErr = errors.New("ledger injoignable")
	if !svc.HasPermission("consumer", role.PermRead) {
		t.Error("la dernière table lue devrait s'appliquer")
	}
}

// / @brief  Table jamais lue et illisible : la table initiale s'applique (RM48)
// / @input  réseau injoignable dès le départ
// / @expect contributor porte write, un rôle personnalisé est inconnu
func TestTable_RM48_Unreachable_NeverRead_DefaultRoles(t *testing.T) {
	repo := newMockRepo()
	repo.readErr = errors.New("ledger injoignable")
	svc := role.NewService(repo)
	if !svc.HasPermission(role.NameContributor, role.PermWrite) {
		t.Error("contributor devrait porter write (table initiale)")
	}
	if _, err := svc.Get("consumer"); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("ErrNotFound attendu, obtenu %v", err)
	}
}

// / @brief  WithRefresh réutilise la table lue pendant le délai, sauf après une modification par ce service
// / @input  refresh 1 h ; un rôle publié par un autre nœud, puis un rôle créé par ce service
// / @expect le rôle de l'autre nœud n'apparaît qu'après la création locale, qui force la relecture
func TestTable_Refresh_ReusesTableWithinDelay(t *testing.T) {
	repo := newMockRepo()
	svc := role.NewService(repo).WithRefresh(time.Hour)
	svc.List()
	repo.Save(&role.Role{Name: "remote", Permissions: []role.Permission{role.PermRead}})
	if _, err := svc.Get("remote"); err == nil {
		t.Fatal("la table lue devrait être réutilisée pendant le délai")
	}
	if _, err := svc.Create("local", []role.Permission{role.PermRead}, false); err != nil {
		t.Fatalf("Create: %v", err)
	}
	for _, name := range []string{"remote", "local"} {
		if _, err := svc.Get(name); err != nil {
			t.Errorf("Get(%q) après relecture : %v", name, err)
		}
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

// / @brief  Create ajoute un rôle à la table, la première publication conservant la table initiale
// / @input  aucune table publiée, nom "consumer", permissions [read]
// / @expect le rôle est créé, non intégré ; la table compte 5 rôles
func TestCreate_FirstPublication_KeepsDefaultRoles(t *testing.T) {
	svc := role.NewService(newMockRepo())
	r, err := svc.Create("consumer", []role.Permission{role.PermRead}, false)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if r.BuiltIn || !r.Has(role.PermRead) {
		t.Errorf("rôle créé inattendu : %+v", r)
	}
	if all, _ := svc.List(); len(all) != 5 {
		t.Errorf("got %d rôles, want 5 (4 initiaux + consumer)", len(all))
	}
}

// / @brief  Create refuse un nom déjà utilisé (y compris un rôle de la table initiale)
// / @input  tentative de Create("reader", ...)
// / @expect retourne ErrAlreadyExists
func TestCreate_AlreadyExists(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Create(role.NameReader, []role.Permission{role.PermRead}, false); !errors.Is(err, role.ErrAlreadyExists) {
		t.Errorf("ErrAlreadyExists attendu, obtenu %v", err)
	}
}

// / @brief  Create refuse un nom invalide (majuscules, caractères interdits)
// / @input  nom "Consumer!"
// / @expect retourne ErrInvalidName
func TestCreate_InvalidName(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Create("Consumer!", []role.Permission{role.PermRead}, false); !errors.Is(err, role.ErrInvalidName) {
		t.Errorf("ErrInvalidName attendu, obtenu %v", err)
	}
}

// / @brief  Create refuse une permission inconnue du catalogue
// / @input  permission "delete-everything"
// / @expect retourne ErrInvalidPermission
func TestCreate_InvalidPermission(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Create("consumer", []role.Permission{"delete-everything"}, false); !errors.Is(err, role.ErrInvalidPermission) {
		t.Errorf("ErrInvalidPermission attendu, obtenu %v", err)
	}
}

// / @brief  Create refuse un rôle d'administration en attribution automatique (RM47)
// / @input  "moderateur" avec identity.admin et autoGrant
// / @expect retourne ErrAdminAutomatic, rien n'est publié
func TestCreate_RM47_AutoGrantAdminRole_Rejected(t *testing.T) {
	repo := newMockRepo()
	svc := role.NewService(repo)
	if _, err := svc.Create("moderateur", []role.Permission{role.PermIdentityAdmin}, true); !errors.Is(err, role.ErrAdminAutomatic) {
		t.Errorf("ErrAdminAutomatic attendu, obtenu %v", err)
	}
	if repo.roles != nil {
		t.Error("rien ne devrait être publié")
	}
}

// / @brief  Sans réseau (Repo absent), la table ne se modifie pas
// / @input  NewService(nil), Create puis Delete
// / @expect ErrNoTable dans les deux cas ; la table initiale reste lisible
func TestCreate_NoTable_Rejected(t *testing.T) {
	svc := role.NewService(nil)
	if _, err := svc.Create("consumer", []role.Permission{role.PermRead}, false); !errors.Is(err, role.ErrNoTable) {
		t.Errorf("Create : ErrNoTable attendu, obtenu %v", err)
	}
	if err := svc.Delete(role.NameReader); !errors.Is(err, role.ErrNoTable) {
		t.Errorf("Delete : ErrNoTable attendu, obtenu %v", err)
	}
	if !svc.HasPermission(role.NameReader, role.PermRead) {
		t.Error("la table initiale devrait s'appliquer")
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

// / @brief  Update change les permissions d'un rôle sans toucher aux champs non fournis
// / @input  "consumer" [read] en attribution automatique, Update(permissions=[read, write])
// / @expect le rôle porte write et reste en attribution automatique
func TestUpdate_Success(t *testing.T) {
	svc := role.NewService(newMockRepo())
	svc.Create("consumer", []role.Permission{role.PermRead}, true)

	updated, err := svc.Update("consumer", role.Change{Permissions: perms(role.PermRead, role.PermWrite)})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !updated.Has(role.PermWrite) || !updated.AutoGrant {
		t.Errorf("rôle mis à jour inattendu : %+v", updated)
	}
	if !svc.HasPermission("consumer", role.PermWrite) {
		t.Error("la table devrait refléter la mise à jour")
	}
}

// / @brief  Un rôle de la table initiale autre qu'admin se modifie (RM48)
// / @input  Update(contributor, permissions=[read])
// / @expect contributor ne porte plus write
func TestUpdate_RM48_DefaultRoleEditable(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Update(role.NameContributor, role.Change{Permissions: perms(role.PermRead)}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if svc.HasPermission(role.NameContributor, role.PermWrite) {
		t.Error("contributor ne devrait plus porter write")
	}
}

// / @brief  Update refuse de modifier admin, immuable (RM48)
// / @input  Update(admin, permissions=[read])
// / @expect retourne ErrAdminImmutable
func TestUpdate_RM48_Admin_Rejected(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Update(role.NameAdmin, role.Change{Permissions: perms(role.PermRead)}); !errors.Is(err, role.ErrAdminImmutable) {
		t.Errorf("ErrAdminImmutable attendu, obtenu %v", err)
	}
}

// / @brief  Update refuse un rôle inexistant
// / @input  Update("ghost", ...)
// / @expect retourne ErrNotFound
func TestUpdate_NotFound(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if _, err := svc.Update("ghost", role.Change{Permissions: perms(role.PermRead)}); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("ErrNotFound attendu, obtenu %v", err)
	}
}

// / @brief  Update refuse de rendre automatique un rôle d'administration, dans un sens comme dans l'autre (RM47)
// / @input  "consumer" automatique auquel on ajoute identity.admin ; "moderateur" d'administration rendu automatique
// / @expect ErrAdminAutomatic dans les deux cas, les rôles restent inchangés
func TestUpdate_RM47_AutoGrantWithAdminPermission_Rejected(t *testing.T) {
	svc := role.NewService(newMockRepo())
	svc.Create("consumer", []role.Permission{role.PermRead}, true)
	svc.Create("moderateur", []role.Permission{role.PermIdentityAdmin}, false)

	if _, err := svc.Update("consumer", role.Change{Permissions: perms(role.PermIdentityAdmin)}); !errors.Is(err, role.ErrAdminAutomatic) {
		t.Errorf("ajout d'une permission d'administration : ErrAdminAutomatic attendu, obtenu %v", err)
	}
	if _, err := svc.Update("moderateur", role.Change{AutoGrant: flag(true)}); !errors.Is(err, role.ErrAdminAutomatic) {
		t.Errorf("attribution automatique : ErrAdminAutomatic attendu, obtenu %v", err)
	}
	if svc.HasPermission("consumer", role.PermIdentityAdmin) {
		t.Error("consumer ne devrait pas avoir été modifié")
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

// / @brief  Delete retire un rôle ; une identité qui le porte n'a plus aucune permission
// / @input  "consumer" créé, puis supprimé
// / @expect Get("consumer") retourne ErrNotFound, HasPermission retourne false
func TestDelete_Success(t *testing.T) {
	svc := role.NewService(newMockRepo())
	svc.Create("consumer", []role.Permission{role.PermRead}, false)

	if err := svc.Delete("consumer"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get("consumer"); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("ErrNotFound attendu, obtenu %v", err)
	}
	if svc.HasPermission("consumer", role.PermRead) {
		t.Error("un rôle supprimé ne devrait porter aucune permission")
	}
}

// / @brief  Un rôle de la table initiale autre qu'admin se supprime (RM48)
// / @input  Delete(reader) sans table publiée
// / @expect reader disparaît, les 3 autres rôles initiaux restent
func TestDelete_RM48_DefaultRoleDeletable(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if err := svc.Delete(role.NameReader); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(role.NameReader); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("ErrNotFound attendu, obtenu %v", err)
	}
	if all, _ := svc.List(); len(all) != 3 {
		t.Errorf("got %d rôles, want 3", len(all))
	}
}

// / @brief  Delete refuse de supprimer admin, immuable (RM48)
// / @input  Delete(admin)
// / @expect retourne ErrAdminImmutable, admin reste présent
func TestDelete_RM48_Admin_Rejected(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if err := svc.Delete(role.NameAdmin); !errors.Is(err, role.ErrAdminImmutable) {
		t.Errorf("ErrAdminImmutable attendu, obtenu %v", err)
	}
	if _, err := svc.Get(role.NameAdmin); err != nil {
		t.Error("admin devrait toujours exister")
	}
}

// ── HasPermission ─────────────────────────────────────────────────────────────

// / @brief  HasPermission reflète les permissions de la table initiale
// / @input  aucune table publiée
// / @expect reader→read oui/write non ; contributor→read+write oui ; admin→tout oui
func TestHasPermission_DefaultTable(t *testing.T) {
	svc := role.NewService(newMockRepo())

	cases := []struct {
		role string
		perm role.Permission
		want bool
	}{
		{role.NameReader, role.PermRead, true},
		{role.NameReader, role.PermWrite, false},
		{role.NameContributor, role.PermWrite, true},
		{role.NameContributor, role.PermNetworkAdmin, false},
		{role.NameAuditor, role.PermRead, true},
		{role.NameAdmin, role.PermNetworkAdmin, true},
		{role.NameAdmin, role.PermRoleAdmin, true},
	}
	for _, c := range cases {
		if got := svc.HasPermission(c.role, c.perm); got != c.want {
			t.Errorf("HasPermission(%q, %q): got %v, want %v", c.role, c.perm, got, c.want)
		}
	}
}

// / @brief  HasPermission retourne false pour un rôle inconnu (fail-closed)
// / @input  rôle "ghost" jamais créé
// / @expect retourne false sans erreur
func TestHasPermission_UnknownRole(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if svc.HasPermission("ghost", role.PermRead) {
		t.Error("un rôle inconnu ne devrait porter aucune permission")
	}
}

// ── Attribution d'un rôle (RM46, RM47) ────────────────────────────────────────

// / @brief  Assignable accepte un rôle de la table, initial ou ajouté (RM46)
// / @input  table avec le rôle ajouté "consumer"
// / @expect aucune erreur pour reader et consumer
func TestAssignable_RM46_KnownRole_OK(t *testing.T) {
	svc := role.NewService(newMockRepo())
	svc.Create("consumer", []role.Permission{role.PermRead}, false)
	for _, name := range []string{role.NameReader, "consumer"} {
		if err := role.Assignable(svc, name); err != nil {
			t.Errorf("Assignable(%q) : %v", name, err)
		}
	}
}

// / @brief  Assignable refuse un rôle absent de la table, sans le remplacer (RM46)
// / @input  rôle "superviseur" inexistant ; puis table absente (nil) et rôle ajouté
// / @expect ErrNotFound dans les deux cas ; sans table, seuls les rôles initiaux sont connus
func TestAssignable_RM46_UnknownRole_Rejected(t *testing.T) {
	svc := role.NewService(newMockRepo())
	if err := role.Assignable(svc, "superviseur"); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("table : ErrNotFound attendu, obtenu %v", err)
	}
	if err := role.Assignable(nil, "consumer"); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("sans table : ErrNotFound attendu, obtenu %v", err)
	}
	if err := role.Assignable(nil, role.NameContributor); err != nil {
		t.Errorf("sans table, rôle initial : %v", err)
	}
}

// / @brief  AssignableByDefault refuse tout rôle portant une permission d'administration (RM46)
// / @input  admin, rôle ajouté "moderateur" avec identity.admin, contributor
// / @expect ErrAdminAutomatic pour admin et moderateur, aucune erreur pour contributor
func TestAssignableByDefault_RM46_AdminPermission_Rejected(t *testing.T) {
	svc := role.NewService(newMockRepo())
	svc.Create("moderateur", []role.Permission{role.PermRead, role.PermIdentityAdmin}, false)
	for _, name := range []string{role.NameAdmin, "moderateur"} {
		if err := role.AssignableByDefault(svc, name); !errors.Is(err, role.ErrAdminAutomatic) {
			t.Errorf("AssignableByDefault(%q) : ErrAdminAutomatic attendu, obtenu %v", name, err)
		}
	}
	if err := role.AssignableByDefault(svc, role.NameContributor); err != nil {
		t.Errorf("AssignableByDefault(contributor) : %v", err)
	}
	if err := role.AssignableByDefault(nil, role.NameAdmin); !errors.Is(err, role.ErrAdminAutomatic) {
		t.Errorf("sans table, admin : ErrAdminAutomatic attendu, obtenu %v", err)
	}
}

// / @brief  GrantedOnRequest suit le réglage d'attribution du rôle, jamais automatique pour un rôle d'administration (RM47)
// / @input  "consumer" automatique, contributor sur approbation, "moderateur" d'administration publié automatique par un autre client, "ghost" inconnu
// / @expect true pour consumer ; false pour contributor et moderateur ; ErrNotFound pour ghost
func TestGrantedOnRequest_RM47(t *testing.T) {
	repo := newMockRepo()
	repo.Save(&role.Role{Name: "moderateur", Permissions: []role.Permission{role.PermIdentityAdmin}, AutoGrant: true})
	svc := role.NewService(repo)
	svc.Create("consumer", []role.Permission{role.PermRead}, true)

	for name, want := range map[string]bool{"consumer": true, role.NameContributor: false, "moderateur": false} {
		got, err := role.GrantedOnRequest(svc, name)
		if err != nil || got != want {
			t.Errorf("GrantedOnRequest(%q) : got %v, %v ; want %v", name, got, err, want)
		}
	}
	if _, err := role.GrantedOnRequest(svc, "ghost"); !errors.Is(err, role.ErrNotFound) {
		t.Errorf("ghost : ErrNotFound attendu, obtenu %v", err)
	}
}
