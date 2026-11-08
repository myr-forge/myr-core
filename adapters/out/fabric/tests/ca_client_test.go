// adapters/out/fabric/tests/ca_client_test.go
// Tests unitaires du client Fabric CA — serveur CA simulé (httptest).
package fabric_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"myr-core/adapters/out/fabric"
)

// writeRegistrar écrit dans dir un certificat et une clé de registrar, la clé
// au format PEM donné (« EC PRIVATE KEY » SEC1 ou « PRIVATE KEY » PKCS#8).
func writeRegistrar(t *testing.T, dir, keyType string) (certPath, keyPath string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("clé : %v", err)
	}
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "admin"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("certificat : %v", err)
	}
	var keyDER []byte
	if keyType == "PRIVATE KEY" {
		keyDER, err = x509.MarshalPKCS8PrivateKey(key)
	} else {
		keyDER, err = x509.MarshalECPrivateKey(key)
	}
	if err != nil {
		t.Fatalf("clé DER : %v", err)
	}
	certPath, keyPath = filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	_ = os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600)
	_ = os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: keyType, Bytes: keyDER}), 0o600)
	return certPath, keyPath
}

// / @brief  UpdateAttributes signe sa requête avec la clé du registrar, au format SEC1 comme PKCS#8 (clés écrites par fabric-ca-client)
// / @input  CA simulée ; registrar dont la clé est en « EC PRIVATE KEY » puis en « PRIVATE KEY »
// / @expect les deux mises à jour aboutissent ; la CA reçoit un PUT authentifié sur l'identité visée
func TestCAClient_UpdateAttributes_SEC1AndPKCS8Keys(t *testing.T) {
	var gotPath, gotAuth string
	ca := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotAuth = r.Method+" "+r.URL.Path, r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"success":true,"result":{}}`))
	}))
	defer ca.Close()

	for _, keyType := range []string{"EC PRIVATE KEY", "PRIVATE KEY"} {
		certPath, keyPath := writeRegistrar(t, t.TempDir(), keyType)
		c, err := fabric.NewCAClient(fabric.Config{CAEndpoint: ca.URL, CAName: "ca-test", CAAdminCertPath: certPath, CAAdminKeyPath: keyPath})
		if err != nil {
			t.Fatalf("%s : NewCAClient : %v", keyType, err)
		}
		gotPath, gotAuth = "", ""
		if err := c.UpdateAttributes(context.Background(), "alice@Test", map[string]string{"Myr.status": "active"}); err != nil {
			t.Errorf("%s : UpdateAttributes : %v", keyType, err)
			continue
		}
		if gotPath != "PUT /api/v1/identities/alice@Test" || !strings.Contains(gotAuth, ".") {
			t.Errorf("%s : requête CA inattendue %q (autorisation %q)", keyType, gotPath, gotAuth)
		}
	}
}

// / @brief  ReEnroll demande un certificat au nom de l'identifiant d'enrôlement du certificat actuel, quel que soit le nom transmis (exigence de la Fabric CA)
// / @input  certificat actuel de CN alice@Test (clé PKCS#8), ReEnroll appelé avec le nom « alice » ; CA simulée
// / @expect le CSR envoyé porte le CN alice@Test ; le nouveau certificat est retourné
func TestCAClient_ReEnroll_UsesCertificateEnrollmentID(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeRegistrar(t, dir, "PRIVATE KEY")
	certPEM, _ := os.ReadFile(certPath)
	keyPEM, _ := os.ReadFile(keyPath)
	// Le registrar sert aussi d'identité ré-enrôlée : on remplace son CN.
	block, _ := pem.Decode(keyPEM)
	k, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := k.(*ecdsa.PrivateKey)
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "alice@Test"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	var gotCN string
	ca := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CSR string `json:"certificate_request"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if b, _ := pem.Decode([]byte(req.CSR)); b != nil {
			if csr, err := x509.ParseCertificateRequest(b.Bytes); err == nil {
				gotCN = csr.Subject.CommonName
			}
		}
		_, _ = w.Write([]byte(`{"success":true,"result":{"Cert":"bmV3"}}`))
	}))
	defer ca.Close()

	c, err := fabric.NewCAClient(fabric.Config{CAEndpoint: ca.URL, CAName: "ca-test"})
	if err != nil {
		t.Fatalf("NewCAClient : %v", err)
	}
	if _, err := c.ReEnroll(context.Background(), "alice", string(certPEM), string(keyPEM)); err != nil {
		t.Fatalf("ReEnroll : %v", err)
	}
	if gotCN != "alice@Test" {
		t.Errorf("CN du CSR : %q, attendu alice@Test", gotCN)
	}
}
