// adapters/out/fabric/tests/ca_attributes_test.go
// Lecture des attributs d'identité (Myr.role, Myr.status) inscrits par la CA
// dans l'extension X.509 du certificat — seule source de vérité de
// l'autorisation (ADR-12).
package fabric_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"myr-core/adapters/out/fabric"
	"myr-core/domain/identity"
)

// selfSignedCert émet un certificat portant, si attrsJSON n'est pas vide,
// l'extension d'attributs de la CA (OID 1.2.3.4.5.6.7.8.1).
func selfSignedCert(t *testing.T, attrsJSON string) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "alice@Org1"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	if attrsJSON != "" {
		tmpl.ExtraExtensions = []pkix.Extension{{
			Id:    asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 1},
			Value: []byte(attrsJSON),
		}}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

/// @brief  CertAttributes lit le rôle et le statut inscrits par la CA dans le certificat
/// @input  certificat portant {"attrs":{"Myr.role":"contributor","Myr.status":"active",...}}
/// @expect Myr.role=contributor, Myr.status=active
func TestCertAttributes_ReadsRoleAndStatus(t *testing.T) {
	cert := selfSignedCert(t, `{"attrs":{"Myr.role":"contributor","Myr.status":"active","hf.EnrollmentID":"alice@Org1"}}`)
	attrs, err := fabric.CertAttributes(cert)
	if err != nil {
		t.Fatalf("CertAttributes: %v", err)
	}
	if attrs[identity.AttrRole] != identity.RoleContributor || attrs[identity.AttrStatus] != identity.StatusActive {
		t.Errorf("attributs inattendus : %v", attrs)
	}
}

/// @brief  Un certificat sans extension d'attributs n'en porte aucun
/// @input  certificat auto-signé sans extension 1.2.3.4.5.6.7.8.1
/// @expect map vide, sans erreur — l'identité n'aura ni rôle ni statut et ne pourra pas se connecter
func TestCertAttributes_NoExtension_Empty(t *testing.T) {
	attrs, err := fabric.CertAttributes(selfSignedCert(t, ""))
	if err != nil {
		t.Fatalf("CertAttributes: %v", err)
	}
	if len(attrs) != 0 {
		t.Errorf("aucun attribut attendu, obtenu %v", attrs)
	}
}

/// @brief  CertAttributes refuse un PEM ou une extension illisibles
/// @input  texte non PEM ; certificat dont l'extension n'est pas du JSON
/// @expect une erreur dans chaque cas
func TestCertAttributes_Invalid_Rejected(t *testing.T) {
	if _, err := fabric.CertAttributes("pas un certificat"); err == nil {
		t.Error("erreur attendue pour un PEM invalide")
	}
	if _, err := fabric.CertAttributes(selfSignedCert(t, "{not json")); err == nil {
		t.Error("erreur attendue pour une extension illisible")
	}
}
