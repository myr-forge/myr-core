// adapters/out/fabric/tests/role_store_test.go
// Tests unitaires de FabricRoleStore (table des rôles publiée sur le ledger,
// ADR-12) — sans connexion réseau.
package fabric_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"myr-core/adapters/out/fabric"
	"myr-core/domain/role"
)

// channelGateway retourne un stubContract distinct par canal.
type channelGateway struct {
	contracts map[string]*stubContract
}

func newChannelGateway(channels ...string) *channelGateway {
	g := &channelGateway{contracts: map[string]*stubContract{}}
	for _, ch := range channels {
		g.contracts[ch] = &stubContract{}
	}
	return g
}

func (g *channelGateway) Contract() fabric.ContractCaller { return g.contracts[""] }
func (g *channelGateway) ContractFor(ch string) fabric.ContractCaller {
	return g.contracts[ch]
}

// / @brief  Save et Delete publient la modification sur chaque canal du réseau (ADR-12, RM48)
// / @input  réseau à deux canaux, passerelle signée par l'administrateur ; Save(consumer) puis Delete(reader)
// / @expect PutRole avec le JSON du rôle, puis DeleteRole reader, soumis sur les deux canaux
func TestFabricRoleStore_RM48_PublishedOnEveryChannel(t *testing.T) {
	gw := newChannelGateway("sandbox", "fabrication")
	store := fabric.NewUserFabricRoleStore(gw, []string{"sandbox", "fabrication"})

	if err := store.Save(&role.Role{Name: "consumer", Permissions: []role.Permission{role.PermRead}, AutoGrant: true}); err != nil {
		t.Fatalf("Save : %v", err)
	}
	for ch, c := range gw.contracts {
		var got role.Role
		if c.lastSubmitName != "PutRole" || json.Unmarshal([]byte(c.lastSubmitArgs[0]), &got) != nil || got.Name != "consumer" || !got.AutoGrant {
			t.Errorf("canal %s : soumission inattendue %s %v", ch, c.lastSubmitName, c.lastSubmitArgs)
		}
	}
	if err := store.Delete(role.NameReader); err != nil {
		t.Fatalf("Delete : %v", err)
	}
	for ch, c := range gw.contracts {
		if c.lastSubmitName != "DeleteRole" || c.lastSubmitArgs[0] != role.NameReader {
			t.Errorf("canal %s : soumission inattendue %s %v", ch, c.lastSubmitName, c.lastSubmitArgs)
		}
	}
}

// / @brief  Un échec de publication sur un canal n'empêche pas les autres et nomme le canal en échec
// / @input  deux canaux, PutRole refusé sur "fabrication"
// / @expect erreur mentionnant fabrication, enveloppant l'erreur du canal ; sandbox a reçu PutRole
func TestFabricRoleStore_ChannelFailure_Reported(t *testing.T) {
	refused := errors.New("endorsement refusé")
	gw := newChannelGateway("sandbox", "fabrication")
	gw.contracts["fabrication"].submitFn = func(string, ...string) ([]byte, error) { return nil, refused }
	store := fabric.NewUserFabricRoleStore(gw, []string{"sandbox", "fabrication"})

	err := store.Save(&role.Role{Name: "consumer"})
	if !errors.Is(err, refused) || !strings.Contains(err.Error(), "fabrication") {
		t.Fatalf("erreur attendue sur fabrication, obtenu %v", err)
	}
	if gw.contracts["sandbox"].lastSubmitName != "PutRole" {
		t.Error("sandbox aurait dû recevoir PutRole")
	}
}

// / @brief  Une passerelle d'organisation ne publie jamais la table (ADR-12)
// / @input  NewFabricRoleStore (lecture seule), Save puis Delete
// / @expect role.ErrSignerRequired, aucune transaction soumise
func TestFabricRoleStore_ReadOnly_Rejected(t *testing.T) {
	gw := newChannelGateway("sandbox")
	store := fabric.NewFabricRoleStore(gw, []string{"sandbox"})
	if err := store.Save(&role.Role{Name: "consumer"}); !errors.Is(err, role.ErrSignerRequired) {
		t.Errorf("Save : ErrSignerRequired attendu, obtenu %v", err)
	}
	if err := store.Delete("consumer"); !errors.Is(err, role.ErrSignerRequired) {
		t.Errorf("Delete : ErrSignerRequired attendu, obtenu %v", err)
	}
	if gw.contracts["sandbox"].lastSubmitName != "" {
		t.Error("aucune transaction ne devait être soumise")
	}
}

// / @brief  FindAll lit la table sur le premier canal
// / @input  deux canaux ; GetRoles du premier retourne admin et consumer
// / @expect 2 rôles décodés, consumer en attribution automatique ; le second canal n'est pas interrogé
func TestFabricRoleStore_FindAll_FirstChannel(t *testing.T) {
	gw := newChannelGateway("sandbox", "fabrication")
	gw.contracts["sandbox"].evaluateFn = func(string, ...string) ([]byte, error) {
		return []byte(`[{"name":"admin","permissions":["admin"]},{"name":"consumer","permissions":["read"],"auto_grant":true}]`), nil
	}
	roles, err := fabric.NewFabricRoleStore(gw, []string{"sandbox", "fabrication"}).FindAll()
	if err != nil {
		t.Fatalf("FindAll : %v", err)
	}
	if len(roles) != 2 || roles[1].Name != "consumer" || !roles[1].AutoGrant {
		t.Errorf("table inattendue : %+v", roles)
	}
	if gw.contracts["sandbox"].lastEvaluateName != "GetRoles" || gw.contracts["fabrication"].lastEvaluateName != "" {
		t.Error("seul le premier canal doit être interrogé")
	}
}
