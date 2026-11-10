// domain/model/tests/immutability_test.go — immuabilité d'un asset soumis (RM19,
// règle 27, ADR-12) et contrôle de la licence proposée (catalogue, RM03).
package model_test

import (
	"errors"
	"testing"

	"myr-core/domain/model"
)

// submittedFixture prépare un composant et un module soumis, et le service
// avec sa blockchain simulée.
type submittedFixture struct {
	svc       *model.Service
	bc        *mockBC
	conns     *mockConnStore
	component *model.Model3D
	module    *model.Model3D
	instance  string
	assembly  string
}

func newSubmittedFixture(t *testing.T) *submittedFixture {
	t.Helper()
	f := &submittedFixture{bc: newMockBC(), conns: newMockConnStore()}
	f.svc = model.NewService(f.bc).
		WithConnStore(f.conns).
		WithThumbStore(newMockThumbStore()).
		WithIfaceStore(newMockIfaceStore()).
		WithDraftStore(newMockDraftStore())

	c, err := f.svc.AddFull(model.AddRequest{Category: model.CategoryBase, Name: "vis", ChannelID: "ch1", LicenseID: "cc-by-4"})
	if err != nil {
		t.Fatalf("AddFull: %v", err)
	}
	if err := f.svc.AddInterface(&model.AssetInterface{AssetID: c.ID, Category: "MECA", Type: "filetage", Direction: model.IfaceOut}); err != nil {
		t.Fatalf("AddInterface: %v", err)
	}
	if f.component, err = f.svc.Submit(c.ID, ""); err != nil {
		t.Fatalf("Submit composant: %v", err)
	}

	mod, err := f.svc.CreateModule(model.ModuleRequest{Name: "drone", ChannelID: "ch1"})
	if err != nil {
		t.Fatalf("CreateModule: %v", err)
	}
	f.assembly = assembleInModule(t, f.svc, mod.ID, f.component.ID, f.component.ID)
	if mod, err = f.svc.GetModule(mod.ID); err != nil {
		t.Fatalf("GetModule: %v", err)
	}
	f.instance = mod.WorkspaceInstances[0].ID
	if f.module, err = f.svc.SubmitModule(mod.ID, "v1"); err != nil {
		t.Fatalf("SubmitModule: %v", err)
	}
	return f
}

// / @brief  Un module soumis n'est jamais re-soumis : une nouvelle version passe par un fork (RM19, ADR-12)
// / @input  module soumis avec une ModuleVersion, puis SubmitModule et Submit à nouveau
// / @expect ErrAssetSubmitted pour les deux appels ; l'enregistrement du ledger garde une seule version
func TestSubmitModule_RM19_AlreadySubmitted_Rejected(t *testing.T) {
	f := newSubmittedFixture(t)
	if _, err := f.svc.SubmitModule(f.module.ID, "v2"); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("SubmitModule : ErrAssetSubmitted attendu, obtenu %v", err)
	}
	if _, err := f.svc.Submit(f.module.ID, "v2"); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("Submit : ErrAssetSubmitted attendu, obtenu %v", err)
	}
	if n := len(f.bc.records[f.module.ID].ModuleVersions); n != 1 {
		t.Errorf("le ledger doit garder 1 version, obtenu %d", n)
	}
}

// / @brief  Un composant soumis n'est jamais re-soumis (RM19, ADR-12)
// / @input  composant soumis, Submit à nouveau
// / @expect ErrAssetSubmitted
func TestSubmit_RM19_AlreadySubmitted_Rejected(t *testing.T) {
	f := newSubmittedFixture(t)
	if _, err := f.svc.Submit(f.component.ID, ""); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("ErrAssetSubmitted attendu, obtenu %v", err)
	}
}

// / @brief  Un brouillon rouvert sous l'identifiant d'un asset déjà enregistré n'est pas soumis (ADR-12)
// / @input  brouillon dont l'identifiant existe déjà sur le ledger
// / @expect ErrAssetSubmitted ; l'enregistrement du ledger n'est pas réécrit
func TestSubmit_DraftAlreadyOnLedger_Rejected(t *testing.T) {
	bc := newMockBC()
	drafts := newMockDraftStore()
	svc := model.NewService(bc).WithIfaceStore(newMockIfaceStore()).WithDraftStore(drafts)
	m, err := svc.AddFull(model.AddRequest{Category: model.CategoryBase, Name: "vis", ChannelID: "ch1"})
	if err != nil {
		t.Fatalf("AddFull: %v", err)
	}
	original := &model.Model3D{ID: m.ID, Name: "original", Status: model.ModuleSubmitted}
	bc.records[m.ID] = original
	if _, err := svc.Submit(m.ID, ""); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Fatalf("ErrAssetSubmitted attendu, obtenu %v", err)
	}
	if bc.records[m.ID] != original {
		t.Error("l'enregistrement du ledger ne doit pas être réécrit")
	}
}

// / @brief  La licence d'un asset soumis ne change plus (RM19, R-04)
// / @input  composant soumis sous « cc-by-4 », UpdateAsset sur le nom et la licence « proprietary »
// / @expect ErrAssetSubmitted ; nom et licence inchangés sur le ledger
func TestUpdateAsset_RM19_SubmittedLicense_Rejected(t *testing.T) {
	f := newSubmittedFixture(t)
	_, err := f.svc.UpdateAsset(model.UpdateRequest{ID: f.component.ID, Name: "autre", LicenseID: "proprietary"})
	if !errors.Is(err, model.ErrAssetSubmitted) {
		t.Fatalf("ErrAssetSubmitted attendu, obtenu %v", err)
	}
	rec := f.bc.records[f.component.ID]
	if rec.Name != "vis" || rec.LicenseID != "cc-by-4" {
		t.Errorf("ledger modifié : nom %q, licence %q", rec.Name, rec.LicenseID)
	}
}

// / @brief  Les métadonnées descriptives d'un asset soumis changent sur le ledger, le reste reste figé (ADR-12)
// / @input  composant soumis ; UpdateAsset sur le nom, la description, les tags et les liens (licence inchangée)
// / @expect les quatre champs changent sur le ledger ; empreinte, licence et interfaces inchangées ; aucun brouillon rouvert
func TestUpdateAsset_SubmittedMetadata_UpdatedOnLedger(t *testing.T) {
	f := newSubmittedFixture(t)
	before := *f.bc.records[f.component.ID]
	_, err := f.svc.UpdateAsset(model.UpdateRequest{ID: f.component.ID, Name: "vis M4", Description: "inox",
		Tags: []string{"vis"}, Links: []string{"https://shop.example/vis"}, LicenseID: "cc-by-4"})
	if err != nil {
		t.Fatalf("UpdateAsset : %v", err)
	}
	rec := f.bc.records[f.component.ID]
	if rec.Name != "vis M4" || rec.Description != "inox" || len(rec.Tags) != 1 || len(rec.Links) != 1 {
		t.Errorf("métadonnées non mises à jour : %+v", rec)
	}
	if rec.Hash != before.Hash || rec.LicenseID != before.LicenseID || len(rec.Interfaces) != len(before.Interfaces) {
		t.Errorf("champs figés modifiés : %+v", rec)
	}
	if got, _ := f.svc.Get(f.component.ID, ""); got.Status != model.ModuleSubmitted {
		t.Errorf("l'asset doit rester soumis, obtenu %q", got.Status)
	}
}

// / @brief  La composition d'un module soumis ne change plus (RM19, FT-012)
// / @input  module soumis : ajout et retrait d'assemblage, ajout et retrait d'instance, déplacement d'instance
// / @expect ErrAssetSubmitted pour chaque mutation ; ni brouillon créé, ni connexion supprimée
func TestModuleComposition_RM19_Submitted_Rejected(t *testing.T) {
	f := newSubmittedFixture(t)
	id := f.module.ID
	checks := map[string]error{
		"AddAssemblyToModule":      f.svc.AddAssemblyToModule(id, "autre"),
		"RemoveAssemblyFromModule": f.svc.RemoveAssemblyFromModule(id, f.assembly),
	}
	_, checks["AddAssetToWorkspace"] = f.svc.AddAssetToWorkspace(id, f.component.ID)
	_, checks["RemoveAssetFromWorkspace"] = f.svc.RemoveAssetFromWorkspace(id, f.instance)
	_, checks["UpdateInstancePosition"] = f.svc.UpdateInstancePosition(id, f.instance, 10, 20)
	for name, err := range checks {
		if !errors.Is(err, model.ErrAssetSubmitted) {
			t.Errorf("%s : ErrAssetSubmitted attendu, obtenu %v", name, err)
		}
	}
	if got, err := f.svc.Get(id, ""); err != nil || got.Status != model.ModuleSubmitted {
		t.Errorf("le module doit rester soumis (aucun brouillon rouvert) : %+v, %v", got, err)
	}
	if conns, _ := f.conns.ListConnections(); len(conns) != 1 {
		t.Errorf("la connexion de l'assemblage ne doit pas être supprimée, obtenu %d", len(conns))
	}
}

// / @brief  Les interfaces d'un asset soumis ne changent plus : elles sont enregistrées avec lui (règle 27)
// / @input  composant soumis portant une interface : ajout, modification et suppression d'interface
// / @expect ErrAssetSubmitted pour chaque mutation
func TestInterfaces_RM19_SubmittedCarrier_Rejected(t *testing.T) {
	f := newSubmittedFixture(t)
	ifaces, err := f.svc.ListInterfacesForAsset(f.component.ID)
	if err != nil || len(ifaces) == 0 {
		t.Fatalf("interfaces du composant : %v, %v", ifaces, err)
	}
	physical := ifaces[0]
	if err := f.svc.AddInterface(&model.AssetInterface{AssetID: f.component.ID, Category: "ELEC"}); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("AddInterface : ErrAssetSubmitted attendu, obtenu %v", err)
	}
	updated := *physical
	updated.Type = "autre"
	if err := f.svc.UpdateInterface(&updated); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("UpdateInterface : ErrAssetSubmitted attendu, obtenu %v", err)
	}
	if err := f.svc.RemoveInterface(physical.ID); !errors.Is(err, model.ErrAssetSubmitted) {
		t.Errorf("RemoveInterface : ErrAssetSubmitted attendu, obtenu %v", err)
	}
}

// / @brief  Une nouvelle version d'un module soumis est un module dérivé (fork), qui se soumet à son tour (RM19)
// / @input  module soumis, CreateModule avec ParentID, SubmitModule du dérivé
// / @expect nouvel identifiant, composition reprise, soumission acceptée ; le module d'origine reste intact
func TestCreateModule_RM19_ForkOfSubmitted_Submittable(t *testing.T) {
	f := newSubmittedFixture(t)
	fork, err := f.svc.CreateModule(model.ModuleRequest{Name: "drone v2", ChannelID: "ch1", ParentID: f.module.ID})
	if err != nil {
		t.Fatalf("CreateModule (fork) : %v", err)
	}
	if fork.ID == f.module.ID || len(fork.WorkspaceInstances) != 2 || len(fork.Assemblies) != 1 {
		t.Fatalf("fork inattendu : %+v", fork)
	}
	if _, err := f.svc.SubmitModule(fork.ID, "v2"); err != nil {
		t.Fatalf("SubmitModule (fork) : %v", err)
	}
	if n := len(f.bc.records[f.module.ID].ModuleVersions); n != 1 {
		t.Errorf("le module d'origine doit garder 1 version, obtenu %d", n)
	}
}

// / @brief  Une licence proposée doit exister dans le catalogue, à la création comme en modification
// / @input  AddFull, CreateModule et UpdateAsset d'un brouillon avec la licence « apache » (le catalogue connaît « apache-2 »)
// / @expect ErrUnknownLicense dans les trois cas
func TestLicense_UnknownLicense_Rejected(t *testing.T) {
	svc := newFullSvc(t)
	if _, err := svc.AddFull(model.AddRequest{Category: model.CategoryBase, ChannelID: "ch1", Name: "vis", LicenseID: "apache"}); !errors.Is(err, model.ErrUnknownLicense) {
		t.Errorf("AddFull : ErrUnknownLicense attendu, obtenu %v", err)
	}
	if _, err := svc.CreateModule(model.ModuleRequest{ChannelID: "ch1", Name: "drone", LicenseID: "apache"}); !errors.Is(err, model.ErrUnknownLicense) {
		t.Errorf("CreateModule : ErrUnknownLicense attendu, obtenu %v", err)
	}
	d, err := svc.AddFull(model.AddRequest{Category: model.CategoryBase, ChannelID: "ch1", Name: "vis"})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	if _, err := svc.UpdateAsset(model.UpdateRequest{ID: d.ID, LicenseID: "apache"}); !errors.Is(err, model.ErrUnknownLicense) {
		t.Errorf("UpdateAsset : ErrUnknownLicense attendu, obtenu %v", err)
	}
}

// / @brief  Changer la licence d'un brouillon dérivé applique RM03 comme à la création
// / @input  parent soumis sous « cc-by-nd-4 » (pas de dérivé), brouillon dérivé, UpdateAsset vers « mit »
// / @expect ErrLicenseIncompatible ; la licence du brouillon reste inchangée
func TestUpdateAsset_RM03_IncompatibleLicense_Rejected(t *testing.T) {
	bc := newMockBC()
	svc := model.NewService(bc).WithIfaceStore(newMockIfaceStore()).WithDraftStore(newMockDraftStore())
	bc.records["parent"] = &model.Model3D{ID: "parent", LicenseID: "cc-by-nd-4", Status: model.ModuleSubmitted}
	child, err := svc.AddFull(model.AddRequest{ChannelID: "ch1", Name: "dérivé", ParentID: "parent", Category: model.CategoryAmelioration})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	if _, err := svc.UpdateAsset(model.UpdateRequest{ID: child.ID, LicenseID: "mit"}); !errors.Is(err, model.ErrLicenseIncompatible) {
		t.Fatalf("ErrLicenseIncompatible attendu, obtenu %v", err)
	}
	if got, _ := svc.Get(child.ID, ""); got.LicenseID != "" {
		t.Errorf("licence du brouillon modifiée : %q", got.LicenseID)
	}
}

// / @brief  Un asset créé avec parent_id d'un module soumis reprend une copie indépendante de sa composition, puis se soumet (UCMOD01, RM19)
// / @input  module soumis (2 instances, 1 liaison) ; AddFull(ParentID = module, licence compatible) ; Submit du dérivé
// / @expect le dérivé est un brouillon avec 2 instances et 1 liaison de nouveaux identifiants ; le module d'origine est inchangé ; Submit réussit
func TestAddFull_RM19_DerivedFromSubmittedModule_CopiesComposition(t *testing.T) {
	f := newSubmittedFixture(t)
	d, err := f.svc.AddFull(model.AddRequest{Name: "drone v2", ChannelID: "ch1", ParentID: f.module.ID, Category: model.CategoryAmelioration})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	if d.Status != model.ModuleDraft || len(d.WorkspaceInstances) != 2 || len(d.Assemblies) != 1 {
		t.Fatalf("dérivé inattendu : status %s, %d instance(s), %d liaison(s)", d.Status, len(d.WorkspaceInstances), len(d.Assemblies))
	}
	if d.WorkspaceInstances[0].ID == f.instance || d.Assemblies[0] == f.assembly {
		t.Error("la composition copiée doit avoir de nouveaux identifiants")
	}
	if got := f.bc.records[f.module.ID]; len(got.WorkspaceInstances) != 2 || len(got.Assemblies) != 1 {
		t.Errorf("le module d'origine ne doit pas changer : %+v", got)
	}
	if _, err := f.svc.Submit(d.ID, "v2"); err != nil {
		t.Errorf("Submit du dérivé : %v", err)
	}
}
