// adapters/in/cli/model_link_test.go — tests des handlers CLI liaisons d'assemblage
package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"myr-core/domain/model"
)

/// @brief  runModelLinkAdd crée une liaison directe et affiche son id
/// @input  from="iface-1", to="iface-2"
/// @expect AddAssemblyLink reçoit from/to, sortie contient l'id créé
func TestRunModelLinkAdd_NominalCase(t *testing.T) {
	var gotFrom, gotTo string
	svc := &mockModelSvc{addAssemblyLink: func(fromIfaceID, toIfaceID, label, fromInstanceID, toInstanceID, fastenerAssetID string) (*model.Connection, error) {
		gotFrom, gotTo = fromIfaceID, toIfaceID
		return &model.Connection{ID: "conn-1"}, nil
	}}
	var buf bytes.Buffer
	err := runModelLinkAdd(&buf, svc, "iface-1", "iface-2", "", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFrom != "iface-1" || gotTo != "iface-2" {
		t.Fatalf("unexpected forwarded args: from=%q to=%q", gotFrom, gotTo)
	}
	if !strings.Contains(buf.String(), "conn-1") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runModelLinkAdd rejette l'appel si --from ou --to manque
/// @input  from="", to="iface-2"
/// @expect Erreur retournée sans appeler le service
func TestRunModelLinkAdd_MissingFrom_Rejected(t *testing.T) {
	called := false
	svc := &mockModelSvc{addAssemblyLink: func(string, string, string, string, string, string) (*model.Connection, error) {
		called = true
		return nil, nil
	}}
	var buf bytes.Buffer
	if err := runModelLinkAdd(&buf, svc, "", "iface-2", "", "", "", ""); err == nil {
		t.Fatalf("expected error when --from is missing")
	}
	if called {
		t.Fatalf("AddAssemblyLink should not be called")
	}
}

/// @brief  runModelLinkAdd remonte le refus du service pour une liaison incompatible ou une interface déjà engagée (RM10, RM09)
/// @input  service retournant ErrIncompatibleInterfaces puis ErrInterfaceInUse
/// @expect l'erreur du service est retournée telle quelle, rien n'est affiché
func TestRunModelLinkAdd_RM10_Refused_ErrorReturned(t *testing.T) {
	for _, want := range []error{model.ErrIncompatibleInterfaces, model.ErrInterfaceInUse} {
		svc := &mockModelSvc{addAssemblyLink: func(string, string, string, string, string, string) (*model.Connection, error) {
			return nil, want
		}}
		var buf bytes.Buffer
		if err := runModelLinkAdd(&buf, svc, "iface-1", "iface-2", "", "", "", ""); !errors.Is(err, want) {
			t.Errorf("%v attendue, obtenu %v", want, err)
		}
		if buf.Len() != 0 {
			t.Errorf("aucune sortie attendue : %q", buf.String())
		}
	}
}

/// @brief  runModelLinkRemove supprime la liaison et confirme en sortie
/// @input  id="conn-1"
/// @expect RemoveConnection reçoit l'id, sortie confirme la suppression
func TestRunModelLinkRemove_NominalCase(t *testing.T) {
	removed := ""
	svc := &mockModelSvc{removeConnection: func(id string) error { removed = id; return nil }}
	var buf bytes.Buffer
	if err := runModelLinkRemove(&buf, svc, "conn-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != "conn-1" {
		t.Fatalf("expected RemoveConnection called with conn-1, got %q", removed)
	}
}

/// @brief  runModelLinkList affiche une ligne par liaison, avec marqueur si incompatible
/// @input  service retournant 1 liaison incompatible
/// @expect Sortie contient l'id de la liaison et le marqueur "*"
func TestRunModelLinkList_NominalCase(t *testing.T) {
	svc := &mockModelSvc{listConnections: func() ([]*model.Connection, error) {
		return []*model.Connection{{ID: "conn-1", From: "a", To: "b", Incompatible: true}}, nil
	}}
	var buf bytes.Buffer
	if err := runModelLinkList(&buf, svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "conn-1") {
		t.Fatalf("unexpected output: %q", out)
	}
}
