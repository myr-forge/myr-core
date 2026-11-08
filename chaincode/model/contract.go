// chaincode/model/contract.go
package model

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract implémente le chaincode myrcc — périmètre minimal couvrant
// exactement les 5 fonctions attendues par domain/model/ports.go BlockchainPort
// (StoreModelRecord/UpdateModelMetadata/GetModelRecord/ListModelRecords/
// VerifyIntegrity), plus IndexHashes, opération de mise en service de l'index
// des empreintes, et les 3 fonctions de la table des rôles attendues par
// domain/role/ports.go Repo (PutRole/DeleteRole/GetRoles, roles.go). Les noms
// de fonction et le nombre d'arguments reproduisent exactement ce qu'invoque
// adapters/out/fabric/blockchain.go côté client — ce fichier-là fait foi, pas
// les signatures indicatives de specs/3-Conception/Chaincode.md §4.
type SmartContract struct{}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "StoreModel":
		return s.storeModel(stub, args)
	case "UpdateMetadata":
		return s.updateMetadata(stub, args)
	case "GetModel":
		return s.getModel(stub, args)
	case "ListModels":
		return s.listModels(stub, args)
	case "VerifyModel":
		return s.verifyModel(stub, args)
	case "IndexHashes":
		return s.indexHashes(stub)
	case "PutRole":
		return s.putRole(stub, args)
	case "DeleteRole":
		return s.deleteRole(stub, args)
	case "GetRoles":
		return s.getRoles(stub)
	default:
		return shim.Error(fmt.Sprintf("fonction chaincode inconnue : %s", fn))
	}
}

// storeModel enregistre le Model3D reçu tel quel (aucune reformulation du JSON
// soumis), indexé par son ID — voir FabricBlockchain.StoreModelRecord — après
// vérification des invariants de specs/3-Conception/Chaincode.md §4.1 : un
// asset déjà enregistré n'est jamais réécrit, son fichier est unique, son
// parent existe et son propriétaire est le signataire de la transaction.
func (s *SmartContract) storeModel(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("StoreModel attend 1 argument (model3dJSON)")
	}
	var m Model3D
	if err := json.Unmarshal([]byte(args[0]), &m); err != nil {
		return shim.Error(fmt.Sprintf("StoreModel : JSON invalide : %v", err))
	}
	if err := validateModel(stub, &m); err != nil {
		return shim.Error("StoreModel : " + err.Error())
	}
	if err := checkSubmitter(stub, m.OwnerID); err != nil {
		return shim.Error("StoreModel : " + err.Error())
	}

	existing, err := stub.GetState(m.ID)
	if err != nil {
		return shim.Error(fmt.Sprintf("StoreModel : lecture ledger : %v", err))
	}
	if existing != nil {
		return shim.Error(fmt.Sprintf("StoreModel : l'asset %s est déjà enregistré — un asset soumis est immuable, toute évolution passe par un fork (RM06, RM19)", m.ID))
	}

	if m.ParentID != "" {
		parent, err := stub.GetState(m.ParentID)
		if err != nil {
			return shim.Error(fmt.Sprintf("StoreModel : lecture ledger : %v", err))
		}
		if parent == nil {
			return shim.Error(fmt.Sprintf("StoreModel : asset parent %s introuvable sur le ledger", m.ParentID))
		}
	}

	if m.Hash != "" {
		key, err := hashIndexKey(stub, m.Hash)
		if err != nil {
			return shim.Error(fmt.Sprintf("StoreModel : clé d'index : %v", err))
		}
		holder, err := stub.GetState(key)
		if err != nil {
			return shim.Error(fmt.Sprintf("StoreModel : lecture ledger : %v", err))
		}
		if holder != nil {
			return shim.Error(fmt.Sprintf("StoreModel : le fichier %s est déjà enregistré sous l'asset %s", m.Hash, holder))
		}
		if err := stub.PutState(key, []byte(m.ID)); err != nil {
			return shim.Error(fmt.Sprintf("StoreModel : écriture index : %v", err))
		}
	}

	if err := stub.PutState(m.ID, []byte(args[0])); err != nil {
		return shim.Error(fmt.Sprintf("StoreModel : écriture ledger : %v", err))
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// metadataKeys sont les seuls champs d'un asset enregistré qui changent après
// sa soumission : ses métadonnées descriptives (ADR-12). Empreinte,
// propriétaire, licence, parent, interfaces et composition restent figés.
var metadataKeys = []string{"Name", "Description", "Tags", "Links"}

// updateMetadata remplace les métadonnées descriptives d'un asset enregistré,
// à la demande de son seul propriétaire. Le reste de l'enregistrement est
// recopié octet pour octet : les champs que ce chaincode ne connaît pas ne
// sont jamais perdus.
func (s *SmartContract) updateMetadata(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("UpdateMetadata attend 2 arguments (id, metadataJSON)")
	}
	id := args[0]
	raw, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : lecture ledger : %v", err))
	}
	if raw == nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : asset %s introuvable sur le ledger", id))
	}
	var record map[string]json.RawMessage
	var current Model3D
	if json.Unmarshal(raw, &record) != nil || json.Unmarshal(raw, &current) != nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : enregistrement %s illisible", id))
	}
	if err := checkSubmitter(stub, current.OwnerID); err != nil {
		return shim.Error("UpdateMetadata : " + err.Error())
	}
	var meta map[string]json.RawMessage
	if err := json.Unmarshal([]byte(args[1]), &meta); err != nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : JSON invalide : %v", err))
	}
	for k := range meta {
		if !slices.Contains(metadataKeys, k) {
			return shim.Error(fmt.Sprintf("UpdateMetadata : le champ %s ne change pas après soumission — seuls %v sont modifiables", k, metadataKeys))
		}
	}
	for _, k := range metadataKeys {
		if v, ok := meta[k]; ok {
			record[k] = v
		}
	}
	updated, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : encodage : %v", err))
	}
	if err := stub.PutState(id, updated); err != nil {
		return shim.Error(fmt.Sprintf("UpdateMetadata : écriture ledger : %v", err))
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// hashDuplicate signale, lors de IndexHashes, un asset déjà enregistré dont
// l'empreinte était déjà indexée pour un autre asset (enregistré avant que
// l'unicité du fichier ne soit vérifiée).
type hashDuplicate struct {
	Hash      string `json:"hash"`
	Kept      string `json:"kept"`
	Duplicate string `json:"duplicate"`
}

// indexHashes construit l'index empreinte → ID pour les assets enregistrés
// avant l'introduction de l'unicité du fichier. Idempotent : une empreinte
// déjà indexée n'est jamais réattribuée. En cas de doublon existant, le
// premier asset dans l'ordre des clés conserve l'empreinte ; les autres sont
// signalés, jamais modifiés (RM06).
func (s *SmartContract) indexHashes(stub shim.ChaincodeStubInterface) peer.Response {
	iter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("IndexHashes : itération ledger : %v", err))
	}
	defer iter.Close()

	indexed := 0
	duplicates := make([]hashDuplicate, 0)
	// Fabric ne rend pas visibles, dans une transaction, ses propres PutState
	// à GetState : les empreintes indexées ici sont donc suivies en mémoire.
	indexedNow := map[string]string{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("IndexHashes : lecture ledger : %v", err))
		}
		var m Model3D
		if json.Unmarshal(kv.Value, &m) != nil || m.Hash == "" {
			continue
		}
		holder, ok := indexedNow[m.Hash]
		key, err := hashIndexKey(stub, m.Hash)
		if err != nil {
			return shim.Error(fmt.Sprintf("IndexHashes : clé d'index : %v", err))
		}
		if !ok {
			stored, err := stub.GetState(key)
			if err != nil {
				return shim.Error(fmt.Sprintf("IndexHashes : lecture ledger : %v", err))
			}
			holder = string(stored)
		}
		switch {
		case holder == "":
			if err := stub.PutState(key, []byte(kv.Key)); err != nil {
				return shim.Error(fmt.Sprintf("IndexHashes : écriture index : %v", err))
			}
			indexedNow[m.Hash] = kv.Key
			indexed++
		case holder != kv.Key:
			duplicates = append(duplicates, hashDuplicate{Hash: m.Hash, Kept: holder, Duplicate: kv.Key})
		}
	}
	out, err := json.Marshal(map[string]any{"indexed": indexed, "duplicates": duplicates})
	if err != nil {
		return shim.Error(fmt.Sprintf("IndexHashes : sérialisation : %v", err))
	}
	return shim.Success(out)
}

// getModel retourne le JSON stocké tel quel — FabricBlockchain.GetModelRecord
// le désérialise directement dans *model.Model3D côté client.
func (s *SmartContract) getModel(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("GetModel attend 1 argument (id)")
	}
	data, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("GetModel : lecture ledger : %v", err))
	}
	if data == nil {
		return shim.Error(fmt.Sprintf("GetModel : asset %s introuvable", args[0]))
	}
	return shim.Success(data)
}

// listModels retourne un tableau JSON de tous les Model3D du ledger, filtré par
// ChannelID si non vide. FabricBlockchain.ListModelRecords désérialise
// directement la réponse dans []*model.Model3D.
func (s *SmartContract) listModels(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	channelID := ""
	if len(args) > 0 {
		channelID = args[0]
	}
	iter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("ListModels : itération ledger : %v", err))
	}
	defer iter.Close()

	results := make([]json.RawMessage, 0)
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("ListModels : lecture ledger : %v", err))
		}
		if channelID != "" {
			var m Model3D
			if jsonErr := json.Unmarshal(kv.Value, &m); jsonErr != nil {
				continue // enregistrement non conforme — ignoré plutôt que de faire échouer tout le listing
			}
			if m.ChannelID != channelID {
				continue
			}
		}
		results = append(results, json.RawMessage(kv.Value))
	}
	out, err := json.Marshal(results)
	if err != nil {
		return shim.Error(fmt.Sprintf("ListModels : sérialisation : %v", err))
	}
	return shim.Success(out)
}

// verifyModel compare le hash stocké au hash fourni — FabricBlockchain.VerifyIntegrity
// attend "true"/"false" en texte brut, pas du JSON.
func (s *SmartContract) verifyModel(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("VerifyModel attend 2 arguments (id, hash)")
	}
	data, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("VerifyModel : lecture ledger : %v", err))
	}
	if data == nil {
		return shim.Success([]byte("false"))
	}
	var m Model3D
	if err := json.Unmarshal(data, &m); err != nil {
		return shim.Error(fmt.Sprintf("VerifyModel : JSON invalide en ledger : %v", err))
	}
	if m.Hash == args[1] {
		return shim.Success([]byte("true"))
	}
	return shim.Success([]byte("false"))
}
