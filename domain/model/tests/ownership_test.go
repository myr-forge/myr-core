// domain/model/tests/ownership_test.go — seul le propriétaire d'un asset le
// modifie, le supprime ou le soumet, administrateur compris (RM54).
package model_test

import (
	"errors"
	"testing"

	"myr-core/domain/model"
)

// ownedFixture prépare un brouillon d'alice@Org1 portant une interface, et
// les vues du service pour alice, bob et sans identité.
type ownedFixture struct {
	bc               *mockBC
	alice, bob, none model.ModelService
	draft            *model.Model3D
	iface            string
}

func newOwnedFixture(t *testing.T) *ownedFixture {
	t.Helper()
	bc := newMockBC()
	svc := model.NewService(bc).
		WithConnStore(newMockConnStore()).
		WithIfaceStore(newMockIfaceStore()).
		WithDraftStore(newMockDraftStore("piece"))
	f := &ownedFixture{bc: bc, alice: svc.As("alice@Org1"), bob: svc.As("bob@Org1"), none: svc.As("")}
	d, err := f.alice.AddFull(model.AddRequest{Name: "vis", ChannelID: "ch1", Category: model.CategoryBase})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	f.draft = d
	iface := &model.AssetInterface{AssetID: d.ID, Category: "MECA", Type: "vis M3", Direction: model.IfaceOut}
	if err := f.alice.AddInterface(iface); err != nil {
		t.Fatalf("AddInterface : %v", err)
	}
	f.iface = iface.ID
	return f
}

// ownerOps liste les opérations qui modifient, suppriment ou soumettent le brouillon d.
func ownerOps(svc model.ModelService, d *model.Model3D, iface string) map[string]error {
	ops := map[string]error{}
	_, ops["UpdateAsset"] = svc.UpdateAsset(model.UpdateRequest{ID: d.ID, Name: "volé"})
	ops["AddInterface"] = svc.AddInterface(&model.AssetInterface{AssetID: d.ID, Category: "ELEC", Type: "USB-C", Direction: model.IfaceIn})
	ops["UpdateInterface"] = svc.UpdateInterface(&model.AssetInterface{ID: iface, AssetID: d.ID, Category: "MECA", Type: "autre", Direction: model.IfaceOut})
	ops["RemoveInterface"] = svc.RemoveInterface(iface)
	_, ops["AddAssetToWorkspace"] = svc.AddAssetToWorkspace(d.ID, "piece")
	ops["AddAssemblyToModule"] = svc.AddAssemblyToModule(d.ID, "c1")
	_, ops["Submit"] = svc.Submit(d.ID, "")
	ops["Remove"] = svc.Remove(d.ID)
	return ops
}

// / @brief  Une autre identité, administrateur compris, ne modifie, ne supprime ni ne soumet le brouillon d'alice (RM54)
// / @input  brouillon d'alice@Org1 ; chaque opération de modification, de suppression et de soumission par bob@Org1
// / @expect ErrNotOwner pour chacune ; le brouillon et son interface sont inchangés
func TestOwnership_RM54_OtherIdentity_Rejected(t *testing.T) {
	f := newOwnedFixture(t)
	for name, err := range ownerOps(f.bob, f.draft, f.iface) {
		if !errors.Is(err, model.ErrNotOwner) {
			t.Errorf("%s : ErrNotOwner attendu, obtenu %v", name, err)
		}
	}
	got, err := f.alice.Get(f.draft.ID, "")
	if err != nil || got.Name != "vis" || got.Status != model.ModuleDraft {
		t.Errorf("le brouillon ne doit pas changer : %+v, %v", got, err)
	}
	if i, _ := f.alice.GetInterface(f.iface); i == nil || i.Type != "vis M3" {
		t.Errorf("l'interface ne doit pas changer : %+v", i)
	}
}

// / @brief  Sans identité, aucune création ni modification n'est permise (RM54)
// / @input  service sans identité : AddFull, puis les opérations sur le brouillon d'alice
// / @expect ErrNotOwner partout
func TestOwnership_RM54_NoIdentity_Rejected(t *testing.T) {
	f := newOwnedFixture(t)
	if _, err := f.none.AddFull(model.AddRequest{Name: "x", ChannelID: "ch1", Category: model.CategoryBase}); !errors.Is(err, model.ErrNotOwner) {
		t.Errorf("AddFull : ErrNotOwner attendu, obtenu %v", err)
	}
	for name, err := range ownerOps(f.none, f.draft, f.iface) {
		if !errors.Is(err, model.ErrNotOwner) {
			t.Errorf("%s : ErrNotOwner attendu, obtenu %v", name, err)
		}
	}
}

// / @brief  Le propriétaire d'un asset créé est l'identité qui le crée, quel que soit le propriétaire demandé (RM54, ADR-12)
// / @input  alice crée un brouillon en demandant « mallory@Org1 » comme propriétaire, puis le modifie et le soumet
// / @expect propriétaire alice@Org1 ; modification et soumission acceptées
func TestOwnership_RM54_CreatorIsOwner(t *testing.T) {
	f := newOwnedFixture(t)
	d, err := f.alice.AddFull(model.AddRequest{Name: "écrou", ChannelID: "ch1", Category: model.CategoryBase, OwnerID: "mallory@Org1"})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	if d.OwnerID != "alice@Org1" {
		t.Errorf("propriétaire : alice@Org1 attendu, obtenu %q", d.OwnerID)
	}
	if _, err := f.alice.UpdateAsset(model.UpdateRequest{ID: d.ID, Name: "écrou M3"}); err != nil {
		t.Errorf("modification par la propriétaire : %v", err)
	}
	if _, err := f.alice.Submit(d.ID, ""); err != nil {
		t.Errorf("soumission par la propriétaire : %v", err)
	}
}

// / @brief  Les métadonnées d'un asset soumis ne changent que par son propriétaire (RM54, ADR-12)
// / @input  asset soumis d'alice ; UpdateAsset du nom par bob
// / @expect ErrNotOwner ; l'enregistrement du ledger garde son nom
func TestOwnership_RM54_SubmittedMetadata_OtherIdentity_Rejected(t *testing.T) {
	f := newOwnedFixture(t)
	if _, err := f.alice.Submit(f.draft.ID, ""); err != nil {
		t.Fatalf("Submit : %v", err)
	}
	if _, err := f.bob.UpdateAsset(model.UpdateRequest{ID: f.draft.ID, Name: "volé"}); !errors.Is(err, model.ErrNotOwner) {
		t.Errorf("ErrNotOwner attendu, obtenu %v", err)
	}
	if rec := f.bc.records[f.draft.ID]; rec == nil || rec.Name != "vis" {
		t.Errorf("le ledger ne doit pas changer : %+v", rec)
	}
}
