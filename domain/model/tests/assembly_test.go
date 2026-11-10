// domain/model/tests/assembly_test.go — règles d'assemblage : usage unique des
// interfaces (RM09), compatibilité à la création (RM10, RM11), liaison directe
// unique (RM50), assemblage interne au module (RM51), composition sans cycle
// (RM52), brouillon resté après soumission (R-14), modifications concurrentes
// d'un brouillon (R-21).
package model_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"myr-core/domain/model"
)

// usbPair ajoute une interface USB-C sortante sur asset1 et une entrante sur
// asset2, compatibles entre elles.
func usbPair(t *testing.T, svc *model.Service) (out, in string) {
	t.Helper()
	out = addTestIface(t, svc, &model.AssetInterface{AssetID: "asset1", Category: "ELEC", Type: "USB-C", Direction: model.IfaceOut})
	in = addTestIface(t, svc, &model.AssetInterface{AssetID: "asset2", Category: "ELEC", Type: "USB-C", Direction: model.IfaceIn})
	return out, in
}

// connCount retourne le nombre de liaisons enregistrées.
func connCount(t *testing.T, svc *model.Service) int {
	t.Helper()
	conns, err := svc.ListConnections()
	if err != nil {
		t.Fatalf("ListConnections : %v", err)
	}
	return len(conns)
}

// / @brief  Une interface d'une instance donnée ne sert qu'à une seule liaison (RM09)
// / @input  liaison out(inst-a1) → in(inst-b1), puis out(inst-a1) → in(inst-b2) et out(inst-a2) → in(inst-b1)
// / @expect ErrInterfaceInUse pour les deux réutilisations ; une seule liaison enregistrée
func TestAddAssemblyLink_RM09_InterfaceUsedOnce(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	out, in := usbPair(t, svc)
	if _, err := svc.AddAssemblyLink(out, in, "", "inst-a1", "inst-b1", ""); err != nil {
		t.Fatalf("première liaison : %v", err)
	}
	if _, err := svc.AddAssemblyLink(out, in, "", "inst-a1", "inst-b2", ""); !errors.Is(err, model.ErrInterfaceInUse) {
		t.Errorf("source réutilisée : ErrInterfaceInUse attendu, obtenu %v", err)
	}
	if _, err := svc.AddAssemblyLink(out, in, "", "inst-a2", "inst-b1", ""); !errors.Is(err, model.ErrInterfaceInUse) {
		t.Errorf("cible réutilisée : ErrInterfaceInUse attendu, obtenu %v", err)
	}
	if n := connCount(t, svc); n != 1 {
		t.Errorf("une seule liaison attendue, obtenu %d", n)
	}
}

// / @brief  L'usage unique porte sur l'interface d'une instance : la même interface reste libre dans les autres instances du composant (RM09)
// / @input  liaison out(inst-a1) → in(inst-b1), puis out(inst-a2) → in(inst-b2)
// / @expect les deux liaisons sont acceptées
func TestAddAssemblyLink_RM09_SameInterfaceOtherInstance_OK(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	out, in := usbPair(t, svc)
	if _, err := svc.AddAssemblyLink(out, in, "", "inst-a1", "inst-b1", ""); err != nil {
		t.Fatalf("première liaison : %v", err)
	}
	if _, err := svc.AddAssemblyLink(out, in, "", "inst-a2", "inst-b2", ""); err != nil {
		t.Errorf("autres instances : liaison attendue, obtenu %v", err)
	}
}

// / @brief  Sans instance, l'usage unique porte sur l'interface elle-même (RM09)
// / @input  liaison out → in sans instance, puis la même paire sans instance
// / @expect ErrInterfaceInUse pour la seconde
func TestAddAssemblyLink_RM09_NoInstance_InterfaceUsedOnce(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	out, in := usbPair(t, svc)
	if _, err := svc.AddAssemblyLink(out, in, "", "", "", ""); err != nil {
		t.Fatalf("première liaison : %v", err)
	}
	if _, err := svc.AddAssemblyLink(out, in, "", "", "", ""); !errors.Is(err, model.ErrInterfaceInUse) {
		t.Errorf("ErrInterfaceInUse attendu, obtenu %v", err)
	}
}

// / @brief  Une interface n'est jamais reliée à elle-même (RM09)
// / @input  interface bidirectionnelle reliée à elle-même, sur la même instance puis sans instance
// / @expect ErrSelfLink dans les deux cas ; aucune liaison enregistrée
func TestAddAssemblyLink_RM09_SelfLink_Rejected(t *testing.T) {
	svc := newFullSvc(t, "asset1")
	bidi := addTestIface(t, svc, &model.AssetInterface{AssetID: "asset1", Category: "NUM", Type: "bus", Direction: model.IfaceBidi})
	if _, err := svc.AddAssemblyLink(bidi, bidi, "", "inst-a", "inst-a", ""); !errors.Is(err, model.ErrSelfLink) {
		t.Errorf("même instance : ErrSelfLink attendu, obtenu %v", err)
	}
	if _, err := svc.AddAssemblyLink(bidi, bidi, "", "", "", ""); !errors.Is(err, model.ErrSelfLink) {
		t.Errorf("sans instance : ErrSelfLink attendu, obtenu %v", err)
	}
	if n := connCount(t, svc); n != 0 {
		t.Errorf("aucune liaison attendue, obtenu %d", n)
	}
}

// / @brief  Deux instances d'un même composant se relient par la même interface (RM09)
// / @input  interface bidirectionnelle d'un composant, de l'instance inst-a vers l'instance inst-b
// / @expect liaison acceptée
func TestAddAssemblyLink_RM09_SameInterfaceTwoInstances_OK(t *testing.T) {
	svc := newFullSvc(t, "asset1")
	bidi := addTestIface(t, svc, &model.AssetInterface{AssetID: "asset1", Category: "NUM", Type: "bus", Direction: model.IfaceBidi})
	if _, err := svc.AddAssemblyLink(bidi, bidi, "", "inst-a", "inst-b", ""); err != nil {
		t.Errorf("liaison attendue, obtenu %v", err)
	}
}

// / @brief  Une liaison entre interfaces incompatibles est refusée à la création (RM10, RM11)
// / @input  ELEC/DIN sortante ↔ MECA/vis entrante, sans asset d'accroche
// / @expect ErrIncompatibleInterfaces ; aucune liaison enregistrée
func TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	din := addTestIface(t, svc, &model.AssetInterface{AssetID: "asset1", Category: "ELEC", Type: "DIN", Direction: model.IfaceOut})
	vis := addTestIface(t, svc, &model.AssetInterface{AssetID: "asset2", Category: "MECA", Type: "vis M3", Direction: model.IfaceIn})
	if _, err := svc.AddAssemblyLink(din, vis, "", "inst-a", "inst-b", ""); !errors.Is(err, model.ErrIncompatibleInterfaces) {
		t.Errorf("ErrIncompatibleInterfaces attendu, obtenu %v", err)
	}
	if n := connCount(t, svc); n != 0 {
		t.Errorf("aucune liaison attendue, obtenu %d", n)
	}
}

// / @brief  Compatibilité des valeurs : une valeur simple vaut [v ; v], 0 est une vraie valeur, une interface sans unité n'impose rien (RM11, R-13)
// / @input  sortie 5 V vers entrée [3 ; 6] V ; sortie 0 V vers entrée [3 ; 6] V ; sortie 5 sans unité vers entrée [10 ; 20] V ; sortie 7 V vers entrée 7 V
// / @expect liaison acceptée, refusée (ErrIncompatibleInterfaces), acceptée, acceptée
func TestAddAssemblyLink_RM11_ValueRanges(t *testing.T) {
	cases := []struct {
		name    string
		out, in model.AssetInterface
		ok      bool
	}{
		{"valeur simple dans la plage", model.AssetInterface{ValueMin: 5, Unit: "V"}, model.AssetInterface{ValueMin: 3, ValueMax: 6, IsRange: true, Unit: "V"}, true},
		{"zéro hors plage", model.AssetInterface{ValueMin: 0, Unit: "V"}, model.AssetInterface{ValueMin: 3, ValueMax: 6, IsRange: true, Unit: "V"}, false},
		{"sans unité", model.AssetInterface{ValueMin: 5}, model.AssetInterface{ValueMin: 10, ValueMax: 20, IsRange: true, Unit: "V"}, true},
		{"valeurs simples égales", model.AssetInterface{ValueMin: 7, Unit: "V"}, model.AssetInterface{ValueMin: 7, Unit: "V"}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := newFullSvc(t, "asset1", "asset2")
			out, in := c.out, c.in
			out.AssetID, out.Category, out.Type, out.Direction = "asset1", "ELEC", "alim", model.IfaceOut
			in.AssetID, in.Category, in.Type, in.Direction = "asset2", "ELEC", "alim", model.IfaceIn
			_, err := svc.AddAssemblyLink(addTestIface(t, svc, &out), addTestIface(t, svc, &in), "", "inst-a", "inst-b", "")
			if c.ok && err != nil {
				t.Errorf("liaison attendue, obtenu %v", err)
			}
			if !c.ok && !errors.Is(err, model.ErrIncompatibleInterfaces) {
				t.Errorf("ErrIncompatibleInterfaces attendu, obtenu %v", err)
			}
		})
	}
}

// / @brief  Une liaison directe identique est refusée : le client la retrouve dans la liste, le serveur ne devine rien (RM50)
// / @input  AddConnection(a, b) deux fois, puis AddConnection(b, a)
// / @expect ErrConnectionExists pour le doublon ; le sens inverse est une autre liaison ; deux liaisons enregistrées
func TestAddConnection_RM50_Duplicate_Rejected(t *testing.T) {
	svc := newFullSvc(t)
	if _, err := svc.AddConnection("a", "b", "lien"); err != nil {
		t.Fatalf("première liaison : %v", err)
	}
	if _, err := svc.AddConnection("a", "b", "autre"); !errors.Is(err, model.ErrConnectionExists) {
		t.Errorf("ErrConnectionExists attendu, obtenu %v", err)
	}
	if _, err := svc.AddConnection("b", "a", "retour"); err != nil {
		t.Errorf("sens inverse : liaison attendue, obtenu %v", err)
	}
	if n := connCount(t, svc); n != 2 {
		t.Errorf("deux liaisons attendues, obtenu %d", n)
	}
}

// / @brief  Un brouillon resté sous l'identifiant d'un asset du ledger n'apparaît pas en double : le ledger fait foi (R-14)
// / @input  enregistrement soumis « vis » sur le ledger et brouillon périmé de même identifiant
// / @expect List renvoie l'asset une seule fois, dans sa version du ledger
func TestList_R14_LeftoverDraft_ListedOnce(t *testing.T) {
	bc := newMockBC()
	drafts := newMockDraftStore()
	svc := model.NewService(bc).WithDraftStore(drafts)
	bc.records["vis"] = &model.Model3D{ID: "vis", Name: "vis", ChannelID: "ch1", Status: model.ModuleSubmitted}
	drafts.drafts["vis"] = &model.Model3D{ID: "vis", Name: "vis (brouillon)", ChannelID: "ch1", Status: model.ModuleDraft}

	list, err := svc.List("ch1")
	if err != nil {
		t.Fatalf("List : %v", err)
	}
	if len(list) != 1 || list[0].Status != model.ModuleSubmitted {
		t.Errorf("une entrée, celle du ledger, attendue : %+v", list)
	}
}

// / @brief  Seule une liaison existante qui relie deux instances du module se rattache à lui (RM51)
// / @input  module avec deux instances ; liaison inconnue ; liaison entre assets hors du module ; liaison entre instances d'un autre module
// / @expect ErrConnectionNotInModule pour les trois ; le module n'a aucun assemblage
func TestAddAssemblyToModule_RM51_ConnectionNotInModule_Rejected(t *testing.T) {
	svc := newFullSvc(t, "a1", "a2", "asset1", "asset2", "x1", "x2")
	mod, _ := svc.CreateModule(model.ModuleRequest{Name: "drone", ChannelID: "ch1"})
	svc.AddAssetToWorkspace(mod.ID, "a1")
	svc.AddAssetToWorkspace(mod.ID, "a2")
	other, _ := svc.CreateModule(model.ModuleRequest{Name: "autre", ChannelID: "ch1"})
	foreign := assembleInModule(t, svc, other.ID, "x1", "x2")
	out, in := usbPair(t, svc)
	instLink, err := svc.AddAssemblyLink(out, in, "", "inst-ailleurs-1", "inst-ailleurs-2", "")
	if err != nil {
		t.Fatalf("AddAssemblyLink : %v", err)
	}

	for name, connID := range map[string]string{"inconnue": "absente", "assets hors module": foreign, "instances d'un autre module": instLink.ID} {
		if err := svc.AddAssemblyToModule(mod.ID, connID); !errors.Is(err, model.ErrConnectionNotInModule) {
			t.Errorf("%s : ErrConnectionNotInModule attendu, obtenu %v", name, err)
		}
	}
	if got, _ := svc.GetModule(mod.ID); len(got.Assemblies) != 0 {
		t.Errorf("aucun assemblage attendu : %v", got.Assemblies)
	}
}

// / @brief  Une liaison inconnue ne fait pas d'un composant un module (RM51)
// / @input  composant brouillon sans instance ; AddAssemblyToModule avec une liaison inconnue
// / @expect ErrConnectionNotInModule ; le composant reste un composant
func TestAddAssemblyToModule_RM51_Component_StaysComponent(t *testing.T) {
	svc := newFullSvc(t)
	c, err := svc.AddFull(model.AddRequest{Category: model.CategoryBase, Name: "vis", ChannelID: "ch1"})
	if err != nil {
		t.Fatalf("AddFull : %v", err)
	}
	if err := svc.AddAssemblyToModule(c.ID, "absente"); !errors.Is(err, model.ErrConnectionNotInModule) {
		t.Errorf("ErrConnectionNotInModule attendu, obtenu %v", err)
	}
	if got, _ := svc.Get(c.ID, ""); got.IsModule() {
		t.Errorf("le composant ne doit pas devenir un module : %+v", got)
	}
}

// / @brief  Une liaison entre les instances d'un module se rattache à lui, qu'elle passe par des interfaces ou relie les assets (RM51)
// / @input  module avec une instance de asset1 et une de asset2 ; liaison d'interfaces entre ces instances
// / @expect AddAssemblyToModule accepte la liaison
func TestAddAssemblyToModule_RM51_InstanceLink_OK(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	mod, _ := svc.CreateModule(model.ModuleRequest{Name: "drone", ChannelID: "ch1"})
	svc.AddAssetToWorkspace(mod.ID, "asset1")
	mod, _ = svc.AddAssetToWorkspace(mod.ID, "asset2")
	out, in := usbPair(t, svc)
	conn, err := svc.AddAssemblyLink(out, in, "", mod.WorkspaceInstances[0].ID, mod.WorkspaceInstances[1].ID, "")
	if err != nil {
		t.Fatalf("AddAssemblyLink : %v", err)
	}
	if err := svc.AddAssemblyToModule(mod.ID, conn.ID); err != nil {
		t.Errorf("liaison interne attendue acceptée, obtenu %v", err)
	}
}

// / @brief  Un module ne se contient pas lui-même, même indirectement (RM52)
// / @input  module A ; A ajouté à A ; module B contenant A, puis B ajouté à A
// / @expect ErrCyclicComposition dans les deux cas ; A n'a aucune instance
func TestAddAssetToWorkspace_RM52_Cycle_Rejected(t *testing.T) {
	svc := newFullSvc(t)
	a, _ := svc.CreateModule(model.ModuleRequest{Name: "A", ChannelID: "ch1"})
	b, _ := svc.CreateModule(model.ModuleRequest{Name: "B", ChannelID: "ch1"})
	if _, err := svc.AddAssetToWorkspace(a.ID, a.ID); !errors.Is(err, model.ErrCyclicComposition) {
		t.Errorf("auto-inclusion : ErrCyclicComposition attendu, obtenu %v", err)
	}
	if _, err := svc.AddAssetToWorkspace(b.ID, a.ID); err != nil {
		t.Fatalf("A dans B : %v", err)
	}
	if _, err := svc.AddAssetToWorkspace(a.ID, b.ID); !errors.Is(err, model.ErrCyclicComposition) {
		t.Errorf("cycle indirect : ErrCyclicComposition attendu, obtenu %v", err)
	}
	if got, _ := svc.GetModule(a.ID); len(got.WorkspaceInstances) != 0 {
		t.Errorf("A ne doit contenir aucune instance : %+v", got.WorkspaceInstances)
	}
}

// / @brief  Une composition cyclique enregistrée avant RM52 ne fait pas récurser sans fin le calcul des interfaces exposées (RM52)
// / @input  brouillons A et B écrits directement dans le stockage, A contenant B et B contenant A, chacun avec une interface
// / @expect GetModuleInterfaces(A) se termine sans erreur
func TestGetModuleInterfaces_RM52_LegacyCycle_Terminates(t *testing.T) {
	drafts := newMockDraftStore()
	svc := model.NewService(newMockBC()).
		WithDraftStore(drafts).
		WithIfaceStore(newMockIfaceStore()).
		WithConnStore(newMockConnStore())
	drafts.drafts["A"] = &model.Model3D{ID: "A", Status: model.ModuleDraft, Assemblies: []string{},
		WorkspaceInstances: []model.WorkspaceInstance{{ID: "ia", AssetID: "B"}}}
	drafts.drafts["B"] = &model.Model3D{ID: "B", Status: model.ModuleDraft, Assemblies: []string{},
		WorkspaceInstances: []model.WorkspaceInstance{{ID: "ib", AssetID: "A"}}}
	addTestIface(t, svc, &model.AssetInterface{AssetID: "A", Category: "ELEC", Type: "USB-C", Direction: model.IfaceOut})
	addTestIface(t, svc, &model.AssetInterface{AssetID: "B", Category: "ELEC", Type: "USB-C", Direction: model.IfaceIn})

	if _, err := svc.GetModuleInterfaces("A"); err != nil {
		t.Errorf("GetModuleInterfaces : %v", err)
	}
}

// / @brief  Une interface ou une instance ne se rattache qu'à un asset existant
// / @input  AddInterface sur l'asset « fantome » ; AddAssetToWorkspace de « fantome » dans un module brouillon
// / @expect ErrNotFound dans les deux cas ; ni interface ni instance enregistrée
func TestAssembly_UnknownAsset_Rejected(t *testing.T) {
	svc := newFullSvc(t)
	if err := svc.AddInterface(&model.AssetInterface{AssetID: "fantome", Category: "ELEC", Type: "USB-C", Direction: model.IfaceOut}); !errors.Is(err, model.ErrNotFound) {
		t.Errorf("AddInterface : ErrNotFound attendu, obtenu %v", err)
	}
	if ifaces, _ := svc.ListInterfacesForAsset("fantome"); len(ifaces) != 0 {
		t.Errorf("aucune interface attendue : %+v", ifaces)
	}
	mod, _ := svc.CreateModule(model.ModuleRequest{Name: "drone", ChannelID: "ch1"})
	if _, err := svc.AddAssetToWorkspace(mod.ID, "fantome"); !errors.Is(err, model.ErrNotFound) {
		t.Errorf("AddAssetToWorkspace : ErrNotFound attendu, obtenu %v", err)
	}
	if got, _ := svc.GetModule(mod.ID); len(got.WorkspaceInstances) != 0 {
		t.Errorf("aucune instance attendue : %+v", got.WorkspaceInstances)
	}
}

// / @brief  La catégorie d'un asset est obligatoire et doit être reconnue : le serveur n'en choisit pas (RM02)
// / @input  AddFull sans catégorie, puis avec la catégorie inconnue « gadget »
// / @expect ErrInvalidCategory dans les deux cas ; aucun brouillon créé
func TestAddFull_RM02_CategoryRequired_Rejected(t *testing.T) {
	drafts := newMockDraftStore()
	svc := model.NewService(newMockBC()).WithDraftStore(drafts)
	for _, c := range []model.Category{"", "gadget"} {
		if _, err := svc.AddFull(model.AddRequest{Name: "vis", ChannelID: "ch1", Category: c}); !errors.Is(err, model.ErrInvalidCategory) {
			t.Errorf("catégorie %q : ErrInvalidCategory attendu, obtenu %v", c, err)
		}
	}
	if len(drafts.drafts) != 0 {
		t.Errorf("aucun brouillon attendu : %d", len(drafts.drafts))
	}
}

// / @brief  Le canal cible d'une création est obligatoire : le serveur n'en choisit pas
// / @input  AddFull puis CreateModule sans canal
// / @expect ErrChannelRequired dans les deux cas
func TestCreate_ChannelRequired_Rejected(t *testing.T) {
	svc := model.NewService(newMockBC()).WithDraftStore(newMockDraftStore())
	if _, err := svc.AddFull(model.AddRequest{Name: "vis", Category: model.CategoryBase}); !errors.Is(err, model.ErrChannelRequired) {
		t.Errorf("AddFull : ErrChannelRequired attendu, obtenu %v", err)
	}
	if _, err := svc.CreateModule(model.ModuleRequest{Name: "drone"}); !errors.Is(err, model.ErrChannelRequired) {
		t.Errorf("CreateModule : ErrChannelRequired attendu, obtenu %v", err)
	}
}

// failingDraftStore simule un stockage de brouillons illisible.
type failingDraftStore struct{ *mockDraftStore }

var errUnreadable = errors.New("stockage illisible")

func (failingDraftStore) GetDraft(string) (*model.Model3D, error) { return nil, errUnreadable }

// / @brief  Un stockage de brouillons illisible est signalé tel quel, pas masqué par une recherche sur le ledger
// / @input  stockage de brouillons dont toute lecture échoue ; asset présent sur le ledger ; AddInterface et Get
// / @expect l'erreur du stockage, ni ErrNotFound ni l'asset du ledger
func TestGetAsset_UnreadableDraftStore_NotMasked(t *testing.T) {
	bc := newMockBC()
	bc.records["vis"] = &model.Model3D{ID: "vis", Status: model.ModuleSubmitted}
	svc := model.NewService(bc).
		WithIfaceStore(newMockIfaceStore()).
		WithDraftStore(failingDraftStore{newMockDraftStore()})
	if err := svc.AddInterface(&model.AssetInterface{AssetID: "vis", Category: "ELEC"}); !errors.Is(err, errUnreadable) {
		t.Errorf("AddInterface : erreur du stockage attendue, obtenu %v", err)
	}
	if _, err := svc.Get("vis", ""); !errors.Is(err, errUnreadable) {
		t.Errorf("Get : erreur du stockage attendue, obtenu %v", err)
	}
}

// / @brief  L'identifiant d'une interface est attribué par le serveur : un identifiant fourni ne remplace pas l'interface existante qui le porte
// / @input  interface de asset1 ; AddInterface sur asset2 avec l'identifiant de cette interface
// / @expect nouvel identifiant attribué ; l'interface de asset1 est inchangée
func TestAddInterface_ClientID_NoOverwrite(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	out, _ := usbPair(t, svc)
	hijack := &model.AssetInterface{ID: out, AssetID: "asset2", Category: "MECA", Type: "vis", Direction: model.IfaceIn}
	if err := svc.AddInterface(hijack); err != nil {
		t.Fatalf("AddInterface : %v", err)
	}
	if hijack.ID == out {
		t.Error("l'identifiant fourni par le client ne doit pas être repris")
	}
	if got, _ := svc.GetInterface(out); got == nil || got.AssetID != "asset1" || got.Category != "ELEC" {
		t.Errorf("l'interface d'origine ne doit pas changer : %+v", got)
	}
}

// / @brief  Une interface ne change pas d'asset porteur
// / @input  interface de asset1 ; UpdateInterface avec AssetID = asset2
// / @expect ErrInterfaceCarrierChange ; l'interface reste sur asset1
func TestUpdateInterface_CarrierChange_Rejected(t *testing.T) {
	svc := newFullSvc(t, "asset1", "asset2")
	out, _ := usbPair(t, svc)
	moved, _ := svc.GetInterface(out)
	changed := *moved
	changed.AssetID = "asset2"
	if err := svc.UpdateInterface(&changed); !errors.Is(err, model.ErrInterfaceCarrierChange) {
		t.Errorf("ErrInterfaceCarrierChange attendu, obtenu %v", err)
	}
	if got, _ := svc.GetInterface(out); got == nil || got.AssetID != "asset1" {
		t.Errorf("l'interface doit rester sur asset1 : %+v", got)
	}
}

// syncDraftStore est un DraftStore sûr en concurrence qui, comme un stockage
// fichier, rend une copie à chaque lecture.
type syncDraftStore struct {
	mu     sync.Mutex
	drafts map[string][]byte
}

func (s *syncDraftStore) SaveDraft(m *model.Model3D) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drafts[m.ID] = b
	return nil
}

func (s *syncDraftStore) GetDraft(id string) (*model.Model3D, error) {
	s.mu.Lock()
	b, ok := s.drafts[id]
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("%w : brouillon %s", model.ErrNotFound, id)
	}
	var m model.Model3D
	return &m, json.Unmarshal(b, &m)
}

func (s *syncDraftStore) RemoveDraft(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.drafts, id)
	return nil
}

func (s *syncDraftStore) ListDrafts(string) ([]*model.Model3D, error) { return nil, nil }

// / @brief  Des modifications simultanées d'un même brouillon ne se perdent pas (R-21)
// / @input  stockage qui rend une copie à chaque lecture ; 20 ajouts d'instance simultanés sur un module
// / @expect le module compte 20 instances
func TestAddAssetToWorkspace_R21_Concurrent_NoLostUpdate(t *testing.T) {
	store := &syncDraftStore{drafts: map[string][]byte{}}
	svc := model.NewService(nil).WithDraftStore(store)
	mod, err := svc.CreateModule(model.ModuleRequest{Name: "drone", ChannelID: "ch1"})
	if err != nil {
		t.Fatalf("CreateModule : %v", err)
	}
	const n = 20
	for i := 0; i < n; i++ {
		store.SaveDraft(&model.Model3D{ID: fmt.Sprintf("asset%d", i), Status: model.ModuleDraft})
	}
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if _, err := svc.AddAssetToWorkspace(mod.ID, fmt.Sprintf("asset%d", i)); err != nil {
				t.Errorf("ajout %d : %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	if got, _ := svc.Get(mod.ID, ""); len(got.WorkspaceInstances) != n {
		t.Errorf("%d instances attendues, obtenu %d", n, len(got.WorkspaceInstances))
	}
}
