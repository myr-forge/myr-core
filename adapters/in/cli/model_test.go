// adapters/in/cli/model_test.go — tests des handlers CLI model (composant)
package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"myr-core/domain/model"
)

// / @brief  runModelAdd doit appeler AddFull avec la requête construite et afficher le
// /         statut brouillon et l'id créé (toute création est un brouillon)
// / @input  AddRequest{Name: "Roue"}, service retournant un Model3D{ID: "abc", Name: "Roue", Status: draft}
// / @expect Sortie contient l'id, le nom, "brouillon" et invite à "myr model submit"
func TestRunModelAdd_NominalCase(t *testing.T) {
	var captured model.AddRequest
	svc := &mockModelSvc{addFull: func(req model.AddRequest) (*model.Model3D, error) {
		captured = req
		return &model.Model3D{ID: "abc", Name: req.Name, Status: model.ModuleDraft}, nil
	}}
	var buf bytes.Buffer
	err := runModelAdd(&buf, svc, model.AddRequest{Name: "Roue", ChannelID: "green"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Name != "Roue" || captured.ChannelID != "green" {
		t.Fatalf("request not forwarded correctly: %+v", captured)
	}
	out := buf.String()
	if !strings.Contains(out, "abc") || !strings.Contains(out, "Roue") ||
		!strings.Contains(out, "brouillon") || !strings.Contains(out, "submit") {
		t.Fatalf("unexpected output: %q", out)
	}
}

// / @brief  runModelAdd doit propager l'erreur du service (ex : incompatibilité de licence RM03)
// / @input  service retournant une erreur
// / @expect L'erreur est retournée telle quelle
func TestRunModelAdd_ServiceError_Rejected(t *testing.T) {
	wantErr := errors.New("incompatibilité de licence")
	svc := &mockModelSvc{addFull: func(model.AddRequest) (*model.Model3D, error) { return nil, wantErr }}
	var buf bytes.Buffer
	err := runModelAdd(&buf, svc, model.AddRequest{Name: "Roue"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

// / @brief  runModelSubmit doit appeler Submit et afficher le statut final
// / @input  id="d1", service retournant Model3D{ID: "d1", Status: submitted}
// / @expect Sortie contient l'id et "submitted", pas d'erreur
func TestRunModelSubmit_NominalCase(t *testing.T) {
	var captured string
	svc := &mockModelSvc{submit: func(assetID, note string) (*model.Model3D, error) {
		captured = assetID
		return &model.Model3D{ID: assetID, Status: model.ModuleSubmitted}, nil
	}}
	var buf bytes.Buffer
	if err := runModelSubmit(&buf, svc, "d1", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured != "d1" {
		t.Fatalf("Submit appelé avec %q, want %q", captured, "d1")
	}
	if !strings.Contains(buf.String(), "d1") || !strings.Contains(buf.String(), "submitted") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

// / @brief  runModelSubmit doit propager l'erreur du service (ex : brouillon introuvable)
// / @input  service retournant une erreur
// / @expect L'erreur est retournée telle quelle
func TestRunModelSubmit_ServiceError_Rejected(t *testing.T) {
	wantErr := errors.New("brouillon introuvable")
	svc := &mockModelSvc{submit: func(string, string) (*model.Model3D, error) { return nil, wantErr }}
	var buf bytes.Buffer
	err := runModelSubmit(&buf, svc, "unknown", "")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

// / @brief  runModelSubmit transmet la note au service quel que soit le type d'asset —
// /         Submit dispatche en interne selon IsModule() (ADR-11)
// / @input  id="d1", note="v1"
// / @expect Submit reçoit id et note, sortie contient le statut "submitted"
func TestRunModelSubmit_WithNote(t *testing.T) {
	var gotID, gotNote string
	svc := &mockModelSvc{submit: func(assetID, note string) (*model.Model3D, error) {
		gotID, gotNote = assetID, note
		return &model.Model3D{ID: assetID, Status: model.ModuleSubmitted}, nil
	}}
	var buf bytes.Buffer
	if err := runModelSubmit(&buf, svc, "d1", "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "d1" || gotNote != "v1" {
		t.Fatalf("unexpected forwarded args: id=%q note=%q", gotID, gotNote)
	}
}

// / @brief  runModelGet doit afficher les métadonnées d'un modèle
// / @input  id="abc", service retournant Model3D{ID: "abc", Name: "Roue"}
// / @expect Sortie contient l'id et le nom
func TestRunModelGet_NominalCase(t *testing.T) {
	svc := &mockModelSvc{get: func(id, _ string) (*model.Model3D, error) {
		return &model.Model3D{ID: id, Name: "Roue"}, nil
	}}
	var buf bytes.Buffer
	if err := runModelGet(&buf, svc, "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "abc") || !strings.Contains(buf.String(), "Roue") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

// / @brief  runModelGet doit propager l'erreur "introuvable"
// / @input  service retournant model.ErrNotFound
// / @expect L'erreur est retournée
func TestRunModelGet_NotFound_Rejected(t *testing.T) {
	svc := &mockModelSvc{get: func(string, string) (*model.Model3D, error) { return nil, model.ErrNotFound }}
	var buf bytes.Buffer
	err := runModelGet(&buf, svc, "missing")
	if !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// / @brief  runModelList doit afficher un tableau avec une ligne par modèle
// / @input  service retournant 2 modèles
// / @expect Les deux ID apparaissent dans la sortie
func TestRunModelList_NominalCase(t *testing.T) {
	svc := &mockModelSvc{list: func(string) ([]*model.Model3D, error) {
		return []*model.Model3D{{ID: "a", Name: "A"}, {ID: "b", Name: "B"}}, nil
	}}
	var buf bytes.Buffer
	if err := runModelList(&buf, svc, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Fatalf("unexpected output: %q", out)
	}
}

// / @brief  runModelList sur un canal vide ne doit pas planter et affiche seulement l'en-tête
// / @input  service retournant une liste vide
// / @expect Pas d'erreur, sortie non vide (en-tête tableau)
func TestRunModelList_Empty(t *testing.T) {
	svc := &mockModelSvc{list: func(string) ([]*model.Model3D, error) { return nil, nil }}
	var buf bytes.Buffer
	if err := runModelList(&buf, svc, "green"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ID") {
		t.Fatalf("expected table header, got: %q", buf.String())
	}
}

// / @brief  Mode dégradé : quand le service retourne des brouillons locaux accompagnés
// /         d'une erreur ErrBlockchainUnreachable, runModelList affiche quand même le
// /         tableau, précédé d'un avertissement explicite — pas une erreur cobra silencieuse
// /         qui masquerait les brouillons disponibles
// / @input  service retournant ([]*model.Model3D{{ID: "local1"}}, model.ErrBlockchainUnreachable)
// / @expect Pas d'erreur retournée ; la sortie contient "Avertissement", le message de
// /         dégradation et l'ID du brouillon
func TestRunModelList_BlockchainUnreachable_WithDrafts_ShowsWarningAndTable(t *testing.T) {
	svc := &mockModelSvc{list: func(string) ([]*model.Model3D, error) {
		return []*model.Model3D{{ID: "local1", Name: "Brouillon local"}}, model.ErrBlockchainUnreachable
	}}
	var buf bytes.Buffer
	if err := runModelList(&buf, svc, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Avertissement") || !strings.Contains(out, model.DegradedListWarning) {
		t.Fatalf("avertissement explicite attendu, got: %q", out)
	}
	if !strings.Contains(out, "local1") {
		t.Fatalf("brouillon local attendu dans le tableau, got: %q", out)
	}
}

// / @brief  Sans brouillon local disponible, runModelList propage l'erreur telle quelle
// /         (pas de faux mode dégradé sur une liste vraiment vide)
// / @input  service retournant (nil, model.ErrBlockchainUnreachable)
// / @expect L'erreur est retournée
func TestRunModelList_BlockchainUnreachable_NoDrafts_Rejected(t *testing.T) {
	svc := &mockModelSvc{list: func(string) ([]*model.Model3D, error) { return nil, model.ErrBlockchainUnreachable }}
	var buf bytes.Buffer
	err := runModelList(&buf, svc, "")
	if !errors.Is(err, model.ErrBlockchainUnreachable) {
		t.Fatalf("expected ErrBlockchainUnreachable, got %v", err)
	}
}

// / @brief  runModelUpdate transmet la requête de patch et affiche l'id mis à jour
// / @input  UpdateRequest{ID: "abc", Description: "new"}
// / @expect Requête transmise telle quelle, sortie contient l'id
func TestRunModelUpdate_NominalCase(t *testing.T) {
	var captured model.UpdateRequest
	svc := &mockModelSvc{updateAsset: func(req model.UpdateRequest) (*model.Model3D, error) {
		captured = req
		return &model.Model3D{ID: req.ID}, nil
	}}
	var buf bytes.Buffer
	err := runModelUpdate(&buf, svc, model.UpdateRequest{ID: "abc", Description: "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Description != "new" {
		t.Fatalf("request not forwarded: %+v", captured)
	}
	if !strings.Contains(buf.String(), "abc") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

// / @brief  runModelRemove supprime l'asset et confirme en sortie
// / @input  id="abc"
// / @expect Sortie confirme la suppression, pas d'erreur
func TestRunModelRemove_NominalCase(t *testing.T) {
	removed := ""
	svc := &mockModelSvc{remove: func(id string) error { removed = id; return nil }}
	var buf bytes.Buffer
	if err := runModelRemove(&buf, svc, "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != "abc" {
		t.Fatalf("expected Remove called with abc, got %q", removed)
	}
	if !strings.Contains(buf.String(), "abc") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

// / @brief  runModelRemove propage une erreur de suppression non supportée (Fabric, règle 9)
// / @input  service retournant une erreur
// / @expect L'erreur est retournée
func TestRunModelRemove_NotSupported_Rejected(t *testing.T) {
	wantErr := errors.New("suppression non supportée")
	svc := &mockModelSvc{remove: func(string) error { return wantErr }}
	var buf bytes.Buffer
	err := runModelRemove(&buf, svc, "abc")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

// / @brief  runModelChildren liste les dérivés d'un composant parent
// / @input  service retournant 1 enfant
// / @expect L'id de l'enfant apparaît dans la sortie
func TestRunModelChildren_NominalCase(t *testing.T) {
	svc := &mockModelSvc{getChildren: func(parentID string) ([]*model.Model3D, error) {
		return []*model.Model3D{{ID: "child-1", ParentID: parentID}}, nil
	}}
	var buf bytes.Buffer
	if err := runModelChildren(&buf, svc, "parent-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "child-1") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}
