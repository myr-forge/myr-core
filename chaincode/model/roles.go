// chaincode/model/roles.go
package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"sort"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Table des rôles du réseau (specs/3-Conception/Chaincode.md §4.3, ADR-12) :
// chaque canal en porte un exemplaire, publié par l'administrateur, que le
// chaincode applique à toute écriture. Tant qu'aucune n'est publiée, la table
// initiale s'applique ; le rôle admin porte toujours toutes les permissions.

const (
	permWrite     = "write"
	permRoleAdmin = "role.admin"
	roleAdmin     = "admin"

	// roleTableType est le type de la clé composite de la table : comme
	// l'index des empreintes, elle reste invisible pour ListModels.
	roleTableType = "roles"
)

// roleEntry porte les champs d'un rôle que le chaincode lit ; les autres
// champs publiés par le domaine sont conservés tels quels.
type roleEntry struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// defaultRoles est la table initiale, copie de domain/role (DefaultRoles) qui
// doit rester identique : le chaincode, module distinct, ne peut pas
// l'importer.
var defaultRoles = []roleEntry{
	{Name: roleAdmin, Permissions: []string{"read", permWrite, "network.admin", permRoleAdmin, "identity.admin", "admin"}},
	{Name: "auditor", Permissions: []string{"read"}},
	{Name: "contributor", Permissions: []string{"read", permWrite}},
	{Name: "reader", Permissions: []string{"read"}},
}

var validRoleName = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

func roleTableKey(stub shim.ChaincodeStubInterface) (string, error) {
	return stub.CreateCompositeKey(roleTableType, []string{})
}

// roleTable retourne la table publiée sur le canal, la table initiale si
// aucune ne l'est. Chaque entrée est le JSON publié du rôle.
func roleTable(stub shim.ChaincodeStubInterface) ([]json.RawMessage, error) {
	key, err := roleTableKey(stub)
	if err != nil {
		return nil, err
	}
	b, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("lecture de la table des rôles : %w", err)
	}
	if b == nil {
		table := make([]json.RawMessage, 0, len(defaultRoles))
		for _, r := range defaultRoles {
			raw, err := json.Marshal(map[string]any{"id": r.Name, "name": r.Name, "permissions": r.Permissions, "built_in": true})
			if err != nil {
				return nil, err
			}
			table = append(table, raw)
		}
		return table, nil
	}
	var table []json.RawMessage
	if err := json.Unmarshal(b, &table); err != nil {
		return nil, fmt.Errorf("table des rôles illisible : %w", err)
	}
	return table, nil
}

// findRole retourne l'index du rôle name dans table, -1 s'il en est absent.
func findRole(table []json.RawMessage, name string) int {
	for i, raw := range table {
		var e roleEntry
		if json.Unmarshal(raw, &e) == nil && e.Name == name {
			return i
		}
	}
	return -1
}

// roleAllows indique si le rôle name porte perm dans la table du canal ; le
// rôle admin porte toujours toutes les permissions.
func roleAllows(stub shim.ChaincodeStubInterface, name, perm string) (bool, error) {
	if name == roleAdmin {
		return true, nil
	}
	table, err := roleTable(stub)
	if err != nil {
		return false, err
	}
	i := findRole(table, name)
	if i < 0 {
		return false, nil
	}
	var e roleEntry
	if err := json.Unmarshal(table[i], &e); err != nil {
		return false, err
	}
	return slices.Contains(e.Permissions, perm), nil
}

// checkRoleAdmin vérifie que le signataire est actif et que son rôle porte la
// permission de gérer les rôles.
func checkRoleAdmin(stub shim.ChaincodeStubInterface) error {
	s, err := activeSigner(stub)
	if err != nil {
		return err
	}
	return requirePermission(stub, s, permRoleAdmin)
}

// saveRoleTable publie table, triée par nom de rôle.
func saveRoleTable(stub shim.ChaincodeStubInterface, table []json.RawMessage) error {
	sort.SliceStable(table, func(i, j int) bool {
		var a, b roleEntry
		_ = json.Unmarshal(table[i], &a)
		_ = json.Unmarshal(table[j], &b)
		return a.Name < b.Name
	})
	b, err := json.Marshal(table)
	if err != nil {
		return err
	}
	key, err := roleTableKey(stub)
	if err != nil {
		return err
	}
	return stub.PutState(key, b)
}

// putRole crée ou remplace un rôle de la table (roleJSON), sauf admin,
// immuable. La première publication part de la table initiale.
func (s *SmartContract) putRole(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("PutRole attend 1 argument (roleJSON)")
	}
	var e roleEntry
	if err := json.Unmarshal([]byte(args[0]), &e); err != nil {
		return shim.Error(fmt.Sprintf("PutRole : JSON invalide : %v", err))
	}
	if !validRoleName.MatchString(e.Name) {
		return shim.Error(fmt.Sprintf("PutRole : nom de rôle invalide %q", e.Name))
	}
	if e.Name == roleAdmin {
		return shim.Error("PutRole : le rôle admin est immuable")
	}
	if err := checkRoleAdmin(stub); err != nil {
		return shim.Error(fmt.Sprintf("PutRole %s : %v", e.Name, err))
	}
	table, err := roleTable(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if i := findRole(table, e.Name); i >= 0 {
		table[i] = json.RawMessage(args[0])
	} else {
		table = append(table, json.RawMessage(args[0]))
	}
	if err := saveRoleTable(stub, table); err != nil {
		return shim.Error(fmt.Sprintf("PutRole : %v", err))
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// deleteRole retire le rôle name de la table, sauf admin, immuable.
func (s *SmartContract) deleteRole(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("DeleteRole attend 1 argument (name)")
	}
	name := args[0]
	if name == roleAdmin {
		return shim.Error("DeleteRole : le rôle admin est immuable")
	}
	if err := checkRoleAdmin(stub); err != nil {
		return shim.Error(fmt.Sprintf("DeleteRole %s : %v", name, err))
	}
	table, err := roleTable(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	i := findRole(table, name)
	if i < 0 {
		return shim.Error(fmt.Sprintf("DeleteRole : rôle %q absent de la table", name))
	}
	if err := saveRoleTable(stub, slices.Delete(table, i, i+1)); err != nil {
		return shim.Error(fmt.Sprintf("DeleteRole : %v", err))
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// getRoles retourne la table des rôles du canal (tableau JSON), la table
// initiale si aucune n'est publiée.
func (s *SmartContract) getRoles(stub shim.ChaincodeStubInterface) peer.Response {
	table, err := roleTable(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	b, err := json.Marshal(table)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(b)
}
