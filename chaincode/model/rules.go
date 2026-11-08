// chaincode/model/rules.go
package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// Invariants vérifiés par le chaincode lui-même, quel que soit le nœud qui
// soumet la transaction (specs/3-Conception/Chaincode.md §4.1) : une règle
// qui n'existerait que côté client n'engagerait que le nœud qui l'applique.

// hashIndexType est le type de clé composite de l'index empreinte → ID.
// Les clés composites n'apparaissent jamais dans un GetStateByRange simple :
// l'index reste invisible pour ListModels.
const hashIndexType = "hash"

// Attributs du certificat émis par la Fabric CA (Register, ECert: true).
const (
	attrRole   = "Myr.role"
	attrStatus = "Myr.status"

	statusActive = "active"
)

var validID = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,128}$`)

// validateModel vérifie la forme d'un Model3D avant toute lecture du ledger.
func validateModel(stub shim.ChaincodeStubInterface, m *Model3D) error {
	if !validID.MatchString(m.ID) {
		return fmt.Errorf("ID invalide %q : caractères [A-Za-z0-9_.-], 128 au plus", m.ID)
	}
	if m.Status != ModuleSubmitted {
		return fmt.Errorf("asset %s : seul un asset soumis (status %q) s'enregistre sur le ledger, reçu %q", m.ID, ModuleSubmitted, m.Status)
	}
	if ch := stub.GetChannelID(); m.ChannelID != ch {
		return fmt.Errorf("asset %s : ChannelID %q différent du canal de la transaction %q", m.ID, m.ChannelID, ch)
	}
	return nil
}

// signer est l'identité vérifiée du signataire d'une transaction.
type signer struct {
	id   string // identifiant d'enrôlement, pseudo@org
	role string // attribut Myr.role du certificat
}

// activeSigner vérifie que le signataire de la transaction appartient à
// l'organisation de son identifiant (pseudo@<org>, MSP <org>MSP) et que son
// certificat le déclare actif.
func activeSigner(stub shim.ChaincodeStubInterface) (signer, error) {
	id, err := cid.New(stub)
	if err != nil {
		return signer{}, fmt.Errorf("identité du signataire illisible : %w", err)
	}
	cert, err := id.GetX509Certificate()
	if err != nil || cert == nil {
		return signer{}, fmt.Errorf("certificat du signataire absent : %v", err)
	}
	enrollmentID := cert.Subject.CommonName
	mspID, err := id.GetMSPID()
	if err != nil {
		return signer{}, fmt.Errorf("MSP du signataire illisible : %w", err)
	}
	_, org, ok := strings.Cut(enrollmentID, "@")
	if !ok || org == "" || mspID != org+"MSP" {
		return signer{}, fmt.Errorf("signataire %q émis par %s : l'identifiant doit être de la forme pseudo@<org> avec le MSP <org>MSP", enrollmentID, mspID)
	}
	if status, _, _ := id.GetAttributeValue(attrStatus); status != statusActive {
		return signer{}, fmt.Errorf("signataire %q : %s vaut %q, %q requis", enrollmentID, attrStatus, status, statusActive)
	}
	role, _, _ := id.GetAttributeValue(attrRole)
	return signer{id: enrollmentID, role: role}, nil
}

// requirePermission vérifie que le rôle du signataire porte perm dans la
// table des rôles du canal.
func requirePermission(stub shim.ChaincodeStubInterface, s signer, perm string) error {
	ok, err := roleAllows(stub, s.role, perm)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("signataire %q : le rôle %q n'a pas la permission %s", s.id, s.role, perm)
	}
	return nil
}

// checkSubmitter vérifie que le signataire de la transaction est actif, qu'il
// est le propriétaire déclaré et que son rôle l'autorise à écrire.
func checkSubmitter(stub shim.ChaincodeStubInterface, ownerID string) error {
	s, err := activeSigner(stub)
	if err != nil {
		return err
	}
	if ownerID != s.id {
		return fmt.Errorf("OwnerID %q différent du signataire %q", ownerID, s.id)
	}
	return requirePermission(stub, s, permWrite)
}

// hashIndexKey retourne la clé composite de l'index pour une empreinte.
func hashIndexKey(stub shim.ChaincodeStubInterface, hash string) (string, error) {
	return stub.CreateCompositeKey(hashIndexType, []string{hash})
}
