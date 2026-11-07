// adapters/in/cli/model_owner_test.go — propriétaire d'un asset créé en ligne de commande (ADR-12)
package cli

import (
	"strings"
	"testing"

	"myr-core/domain/session"
)

func sessionOf(name, orgID string) *mockSessionSvc {
	return &mockSessionSvc{current: func() (*session.Session, error) {
		return &session.Session{Name: name, OrgID: orgID}, nil
	}}
}

// / @brief  SessionIdentity déduit l'identité pseudo@org du compte local CLI
// / @input  compte local alice / Org1MSP ; compte sans organisation ; aucun compte local
// / @expect "alice@Org1" ; absente sans organisation ; absente sans compte
func TestSessionIdentity(t *testing.T) {
	if h, ok := SessionIdentity(sessionOf("alice", "Org1MSP")); !ok || h != "alice@Org1" {
		t.Errorf("compte alice/Org1MSP : (%q, %v), want (alice@Org1, true)", h, ok)
	}
	if _, ok := SessionIdentity(sessionOf("alice", "")); ok {
		t.Error("compte sans organisation : identité inattendue")
	}
	if _, ok := SessionIdentity(&mockSessionSvc{}); ok {
		t.Error("sans compte local : identité inattendue")
	}
}

// / @brief  resolveOwner attribue l'asset à l'identité du compte local, seule autorisée à le soumettre
// / @input  compte local alice/Org1MSP avec --owner-id vide, égal à alice@Org1, puis égal à bob@Org1 ; sans compte local avec --owner-id user1
// / @expect alice@Org1 ; alice@Org1 ; erreur citant les deux identités ; user1 (repris tel quel)
func TestResolveOwner(t *testing.T) {
	alice := sessionOf("alice", "Org1MSP")
	for _, flag := range []string{"", "alice@Org1"} {
		if got, err := resolveOwner(alice, flag); err != nil || got != "alice@Org1" {
			t.Errorf("--owner-id %q : (%q, %v), want alice@Org1", flag, got, err)
		}
	}
	_, err := resolveOwner(alice, "bob@Org1")
	if err == nil || !strings.Contains(err.Error(), "bob@Org1") || !strings.Contains(err.Error(), "alice@Org1") {
		t.Errorf("--owner-id bob@Org1 : err = %v, want un refus citant bob@Org1 et alice@Org1", err)
	}
	if got, err := resolveOwner(&mockSessionSvc{}, "user1"); err != nil || got != "user1" {
		t.Errorf("sans compte local : (%q, %v), want user1", got, err)
	}
}
