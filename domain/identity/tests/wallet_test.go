// domain/identity/tests/wallet_test.go
// Identité « pseudo@org » et recherche du wallet qui signe les écritures de
// l'utilisateur sur le ledger (ADR-12).
package identity_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"myr-core/domain/identity"
)

// / @brief  Handle construit l'identifiant pseudo@org, que l'organisation soit donnée par son MSP ou seule
// / @input  ("alice", "Org1MSP") et ("alice", "Org1")
// / @expect "alice@Org1" dans les deux cas
func TestHandle_Format(t *testing.T) {
	for _, org := range []string{"Org1MSP", "Org1"} {
		if got := identity.Handle("alice", org); got != "alice@Org1" {
			t.Errorf("Handle(alice, %s) = %q, want alice@Org1", org, got)
		}
	}
}

// / @brief  Wallet retrouve le wallet créé par Enroll, avec les chemins de son certificat et de sa clé
// / @input  Enroll("alice", secret, "Org1MSP") avec un stubCA, puis Wallet("alice@Org1")
// / @expect entrée Handle alice@Org1, OrgID Org1MSP ; CertPath et KeyPath désignent des fichiers existants
func TestWallet_AfterEnroll_Found(t *testing.T) {
	svc := identity.NewService(t.TempDir(), &stubCA{})
	if _, err := svc.Enroll(context.Background(), "alice", "s3cret", "Org1MSP"); err != nil {
		t.Fatalf("Enroll : %v", err)
	}

	w, err := svc.Wallet("alice@Org1")
	if err != nil {
		t.Fatalf("Wallet : %v", err)
	}
	if w.Handle != "alice@Org1" || w.OrgID != "Org1MSP" {
		t.Errorf("entrée = %+v, want Handle alice@Org1 et OrgID Org1MSP", w)
	}
	for _, p := range []string{w.CertPath(), w.KeyPath()} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s : %v", p, err)
		}
	}
}

// / @brief  Wallet signale une identité jamais enrôlée sur ce nœud
// / @input  répertoire de wallets vide, Wallet("bob@Org1")
// / @expect erreur ErrWalletNotFound
func TestWallet_Unknown_Rejected(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	if _, err := svc.Wallet("bob@Org1"); !errors.Is(err, identity.ErrWalletNotFound) {
		t.Errorf("err = %v, want ErrWalletNotFound", err)
	}
}

// / @brief  Wallet refuse un identifiant qui désignerait un chemin hors du répertoire des wallets
// / @input  "..", "../alice@Org1", "a/b", `a\b`, chaîne vide
// / @expect une erreur pour chacun, jamais ErrWalletNotFound (l'identifiant est rejeté avant toute lecture)
func TestWallet_PathLikeHandle_Rejected(t *testing.T) {
	svc := identity.NewService(t.TempDir(), nil)
	for _, h := range []string{"..", "../alice@Org1", "a/b", `a\b`, ""} {
		_, err := svc.Wallet(h)
		if err == nil || errors.Is(err, identity.ErrWalletNotFound) {
			t.Errorf("Wallet(%q) err = %v, want un rejet de l'identifiant", h, err)
		}
	}
}
