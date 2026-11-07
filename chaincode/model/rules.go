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

// writerRoles liste les rôles autorisés à écrire sur le ledger. Provisoire :
// la cible (ADR-12, Chaincode.md §4.1) est la table des rôles du réseau,
// éditée par l'administrateur et publiée sur le ledger — ce contrôle devra
// la consulter, seul le rôle admin conservant toutes les permissions.
var writerRoles = map[string]bool{"contributor": true, "admin": true}

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

// checkSubmitter vérifie que le signataire de la transaction est le
// propriétaire déclaré, qu'il appartient à l'organisation de son identifiant
// et que son certificat l'autorise à écrire (statut actif, rôle d'écriture).
func checkSubmitter(stub shim.ChaincodeStubInterface, ownerID string) error {
	id, err := cid.New(stub)
	if err != nil {
		return fmt.Errorf("identité du signataire illisible : %w", err)
	}
	cert, err := id.GetX509Certificate()
	if err != nil || cert == nil {
		return fmt.Errorf("certificat du signataire absent : %v", err)
	}
	enrollmentID := cert.Subject.CommonName
	if ownerID != enrollmentID {
		return fmt.Errorf("OwnerID %q différent du signataire %q", ownerID, enrollmentID)
	}
	mspID, err := id.GetMSPID()
	if err != nil {
		return fmt.Errorf("MSP du signataire illisible : %w", err)
	}
	_, org, ok := strings.Cut(enrollmentID, "@")
	if !ok || org == "" || mspID != org+"MSP" {
		return fmt.Errorf("signataire %q émis par %s : l'identifiant doit être de la forme pseudo@<org> avec le MSP <org>MSP", enrollmentID, mspID)
	}
	if status, _, _ := id.GetAttributeValue(attrStatus); status != statusActive {
		return fmt.Errorf("signataire %q : %s vaut %q, %q requis", enrollmentID, attrStatus, status, statusActive)
	}
	if role, _, _ := id.GetAttributeValue(attrRole); !writerRoles[role] {
		return fmt.Errorf("signataire %q : le rôle %q n'autorise pas l'écriture", enrollmentID, role)
	}
	return nil
}

// hashIndexKey retourne la clé composite de l'index pour une empreinte.
func hashIndexKey(stub shim.ChaincodeStubInterface, hash string) (string, error) {
	return stub.CreateCompositeKey(hashIndexType, []string{hash})
}
