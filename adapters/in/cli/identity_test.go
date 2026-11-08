// adapters/in/cli/identity_test.go — tests des handlers CLI identité
package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"myr-core/domain/identity"
	"myr-core/domain/network"
)

/// @brief  runIdentityWallets liste les wallets locaux
/// @input  service retournant 1 wallet
/// @expect Sortie contient le handle du wallet
func TestRunIdentityWallets_NominalCase(t *testing.T) {
	svc := &mockIdentitySvc{listLocalWallets: func() ([]identity.WalletEntry, error) {
		return []identity.WalletEntry{{Handle: "alice@Org1MSP", Name: "alice", OrgID: "Org1MSP", Status: identity.StatusActive}}, nil
	}}
	var buf bytes.Buffer
	if err := runIdentityWallets(&buf, svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alice@Org1MSP") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentityStatus découpe le handle pseudo@org et interroge la CA
/// @input  handle="alice@Org1"
/// @expect GetStatus reçoit Name="alice" OrgID="Org1MSP", sortie contient le statut
func TestRunIdentityStatus_NominalCase(t *testing.T) {
	var got identity.WalletEntry
	svc := &mockIdentitySvc{getStatus: func(_ context.Context, wallet identity.WalletEntry) (string, error) {
		got = wallet
		return identity.StatusActive, nil
	}}
	var buf bytes.Buffer
	if err := runIdentityStatus(&buf, svc, "alice@Org1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "alice" || got.OrgID != "Org1MSP" {
		t.Fatalf("unexpected wallet forwarded: %+v", got)
	}
	if !strings.Contains(buf.String(), "active") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentityStatus rejette un handle sans "@"
/// @input  handle="alice"
/// @expect Erreur retournée, GetStatus jamais appelé
func TestRunIdentityStatus_InvalidHandle_Rejected(t *testing.T) {
	called := false
	svc := &mockIdentitySvc{getStatus: func(context.Context, identity.WalletEntry) (string, error) {
		called = true
		return "", nil
	}}
	var buf bytes.Buffer
	if err := runIdentityStatus(&buf, svc, "alice"); err == nil {
		t.Fatalf("expected error for invalid handle")
	}
	if called {
		t.Fatalf("GetStatus should not be called for an invalid handle")
	}
}

/// @brief  runIdentityEnroll enrôle l'identité et affiche le handle + statut du wallet
/// @input  name="alice", secret="s3cr3t", orgID="Org1MSP"
/// @expect Enroll reçoit les 3 valeurs, sortie contient le handle
func TestRunIdentityEnroll_NominalCase(t *testing.T) {
	var gotName, gotSecret, gotOrg string
	svc := &mockIdentitySvc{enroll: func(_ context.Context, name, secret, orgID string) (identity.WalletEntry, error) {
		gotName, gotSecret, gotOrg = name, secret, orgID
		return identity.WalletEntry{Handle: name + "@" + orgID, Status: identity.StatusPending}, nil
	}}
	var buf bytes.Buffer
	err := runIdentityEnroll(&buf, svc, "alice", "s3cr3t", "Org1MSP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "alice" || gotSecret != "s3cr3t" || gotOrg != "Org1MSP" {
		t.Fatalf("unexpected forwarded args: name=%q secret=%q org=%q", gotName, gotSecret, gotOrg)
	}
	if !strings.Contains(buf.String(), "alice@Org1MSP") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentityEnroll propage une erreur de secret invalide
/// @input  service retournant une erreur
/// @expect L'erreur est retournée
func TestRunIdentityEnroll_InvalidSecret_Rejected(t *testing.T) {
	wantErr := errors.New("secret invalide")
	svc := &mockIdentitySvc{enroll: func(context.Context, string, string, string) (identity.WalletEntry, error) {
		return identity.WalletEntry{}, wantErr
	}}
	var buf bytes.Buffer
	err := runIdentityEnroll(&buf, svc, "alice", "bad", "Org1MSP")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

/// @brief  runIdentityReEnroll retrouve le wallet local par handle et le ré-enrôle
/// @input  handle="alice@Org1MSP" présent dans ListLocalWallets
/// @expect ReEnroll est appelé avec ce wallet, sortie confirme le nouveau statut
func TestRunIdentityReEnroll_NominalCase(t *testing.T) {
	wallet := identity.WalletEntry{Handle: "alice@Org1MSP", Name: "alice", OrgID: "Org1MSP", MSPDir: "/tmp/msp"}
	var gotWallet identity.WalletEntry
	svc := &mockIdentitySvc{
		listLocalWallets: func() ([]identity.WalletEntry, error) { return []identity.WalletEntry{wallet}, nil },
		reEnroll: func(_ context.Context, w identity.WalletEntry) (identity.WalletEntry, error) {
			gotWallet = w
			w.Status = identity.StatusActive
			return w, nil
		},
	}
	var buf bytes.Buffer
	if err := runIdentityReEnroll(&buf, svc, "alice@Org1MSP"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotWallet.MSPDir != "/tmp/msp" {
		t.Fatalf("expected the local wallet entry to be forwarded, got %+v", gotWallet)
	}
	if !strings.Contains(buf.String(), "active") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentityReEnroll rejette un handle sans wallet local (pas encore enrôlé)
/// @input  ListLocalWallets ne contient pas le handle demandé
/// @expect Erreur retournée, ReEnroll jamais appelé
func TestRunIdentityReEnroll_NoLocalWallet_Rejected(t *testing.T) {
	called := false
	svc := &mockIdentitySvc{
		listLocalWallets: func() ([]identity.WalletEntry, error) { return nil, nil },
		reEnroll: func(context.Context, identity.WalletEntry) (identity.WalletEntry, error) {
			called = true
			return identity.WalletEntry{}, nil
		},
	}
	var buf bytes.Buffer
	if err := runIdentityReEnroll(&buf, svc, "bob@Org1MSP"); err == nil {
		t.Fatalf("expected error when no local wallet matches the handle")
	}
	if called {
		t.Fatalf("ReEnroll should not be called")
	}
}

/// @brief  runIdentityRequest, aucun réseau actif : l'inscription est fermée
/// @input  nsvc=nil (pas de réseau actif)
/// @expect RequestAccess reçoit une politique fermée, l'erreur ErrRegistrationClosed remonte
func TestRunIdentityRequest_NoNetwork_Rejected(t *testing.T) {
	var got identity.RegistrationPolicy
	svc := &mockIdentitySvc{
		requestAccess: func(_ context.Context, _ identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			got = policy
			return nil, "", identity.ErrRegistrationClosed
		},
	}
	var buf bytes.Buffer
	req := identity.AccountRequest{Pseudo: "bob", Email: "bob@example.com", OrgID: "Org1MSP"}
	err := runIdentityRequest(&buf, svc, nil, req)
	if !errors.Is(err, identity.ErrRegistrationClosed) {
		t.Fatalf("expected ErrRegistrationClosed, got %v", err)
	}
	if got.AllowAutoRegister {
		t.Fatalf("expected a closed policy without active network, got %+v", got)
	}
}

/// @brief  runIdentityRequest transmet les règles d'inscription du réseau actif
/// @input  réseau actif AllowAutoRegister=true, RequireApproval=true, AutoRegisterRole=contributor
/// @expect RequestAccess reçoit ces trois règles ; la sortie donne le secret et la commande d'approbation
func TestRunIdentityRequest_RequireApproval_Pending(t *testing.T) {
	var got identity.RegistrationPolicy
	svc := &mockIdentitySvc{
		requestAccess: func(_ context.Context, req identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			got = policy
			req.ID = "req-1"
			req.Status = identity.RequestPending
			return &req, "generated-secret", nil
		},
	}
	nsvc := &mockNetworkSvc{getActive: func() (*network.NetworkProfile, error) {
		return &network.NetworkProfile{AllowAutoRegister: true, RequireApproval: true, AutoRegisterRole: "contributor"}, nil
	}}
	var buf bytes.Buffer
	req := identity.AccountRequest{Pseudo: "bob", Email: "bob@example.com", OrgID: "Org1MSP"}
	if err := runIdentityRequest(&buf, svc, nsvc, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := identity.RegistrationPolicy{AllowAutoRegister: true, RequireApproval: true, DefaultRole: "contributor"}
	if got != want {
		t.Fatalf("policy = %+v, want %+v", got, want)
	}
	out := buf.String()
	if !strings.Contains(out, "generated-secret") || !strings.Contains(out, "myr identity approve req-1") {
		t.Fatalf("unexpected output: %q", out)
	}
}

/// @brief  runIdentityRequest, réseau ouvert sans approbation : l'identité est créée active
/// @input  réseau actif AllowAutoRegister=true, RequireApproval=false
/// @expect Le secret d'enrôlement apparaît, aucune commande d'approbation n'est suggérée
func TestRunIdentityRequest_AutoRegister_Approved(t *testing.T) {
	svc := &mockIdentitySvc{
		requestAccess: func(_ context.Context, req identity.AccountRequest, _ identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
			req.ID = "req-2"
			req.Status = identity.RequestApproved
			return &req, "generated-secret", nil
		},
	}
	nsvc := &mockNetworkSvc{getActive: func() (*network.NetworkProfile, error) {
		return &network.NetworkProfile{AllowAutoRegister: true, AutoRegisterRole: "reader"}, nil
	}}
	var buf bytes.Buffer
	req := identity.AccountRequest{Pseudo: "carol", Email: "carol@example.com", OrgID: "Org1MSP"}
	if err := runIdentityRequest(&buf, svc, nsvc, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "generated-secret") || !strings.Contains(out, "carol@Org1") {
		t.Fatalf("expected the handle and enrollment secret in output, got: %q", out)
	}
	if strings.Contains(out, "myr identity approve") {
		t.Fatalf("no approval expected for an active identity, got: %q", out)
	}
}

/// @brief  runIdentityDecision approuve une demande via le service
/// @input  ApproveRequest renvoie la demande approuvée de alice@Org1MSP
/// @expect L'identifiant transmis est celui de la demande ; la sortie donne le statut et le handle
func TestRunIdentityDecision_Approve_OK(t *testing.T) {
	var gotID string
	svc := &mockIdentitySvc{approveRequest: func(_ context.Context, id string) (*identity.AccountRequest, error) {
		gotID = id
		return &identity.AccountRequest{ID: id, Pseudo: "alice", OrgID: "Org1MSP", Status: identity.RequestApproved}, nil
	}}
	var buf bytes.Buffer
	if err := runIdentityDecision(&buf, svc.ApproveRequest, "req-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "req-1" {
		t.Fatalf("ApproveRequest called with %q, want req-1", gotID)
	}
	if out := buf.String(); !strings.Contains(out, identity.RequestApproved) || !strings.Contains(out, "alice@Org1") {
		t.Fatalf("unexpected output: %q", out)
	}
}

/// @brief  runIdentityDecision propage le refus du service
/// @input  RejectRequest renvoie ErrRequestNotPending
/// @expect L'erreur remonte telle quelle
func TestRunIdentityDecision_NotPending_Rejected(t *testing.T) {
	svc := &mockIdentitySvc{rejectRequest: func(context.Context, string) (*identity.AccountRequest, error) {
		return nil, identity.ErrRequestNotPending
	}}
	var buf bytes.Buffer
	if err := runIdentityDecision(&buf, svc.RejectRequest, "req-1"); !errors.Is(err, identity.ErrRequestNotPending) {
		t.Fatalf("expected ErrRequestNotPending, got %v", err)
	}
}

/// @brief  runIdentityRegister crée un compte actif via CreateAccount
/// @input  pseudo=alice, org=Org1MSP, rôle=contributor, email
/// @expect CreateAccount reçoit ces valeurs ; la sortie donne le handle et le secret
func TestRunIdentityRegister_NominalCase(t *testing.T) {
	var gotRole string
	svc := &mockIdentitySvc{createAccount: func(_ context.Context, pseudo, orgID, role, _ string) (string, string, error) {
		gotRole = role
		return identity.Handle(pseudo, orgID), "s3cret", nil
	}}
	var buf bytes.Buffer
	if err := runIdentityRegister(&buf, svc, "alice", "Org1MSP", "contributor", "alice@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRole != "contributor" {
		t.Fatalf("role = %q, want contributor", gotRole)
	}
	if out := buf.String(); !strings.Contains(out, "alice@Org1") || !strings.Contains(out, "s3cret") {
		t.Fatalf("unexpected output: %q", out)
	}
}

/// @brief  runIdentityRequests liste les demandes d'accès en attente
/// @input  service retournant 1 demande
/// @expect Sortie contient le pseudo de la demande
func TestRunIdentityRequests_NominalCase(t *testing.T) {
	svc := &mockIdentitySvc{listRequests: func() ([]*identity.AccountRequest, error) {
		return []*identity.AccountRequest{{ID: "req-1", Pseudo: "dave", Status: identity.RequestPending}}, nil
	}}
	var buf bytes.Buffer
	if err := runIdentityRequests(&buf, svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dave") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentityRequestRole transmet le wallet du compte local et distingue rôle attribué et demande en attente (RM47)
/// @input  wallet alice@Org1 (reader) ; demande approuvée puis demande en attente
/// @expect RequestRole reçoit le wallet et le rôle ; sortie « attribué » puis « en attente »
func TestRunIdentityRequestRole_RM47_GrantedOrPending(t *testing.T) {
	wallet := identity.WalletEntry{Handle: "alice@Org1", Role: identity.RoleReader, Status: identity.StatusActive}
	status := identity.RequestApproved
	var gotWallet identity.WalletEntry
	svc := &mockIdentitySvc{
		wallet: func(string) (identity.WalletEntry, error) { return wallet, nil },
		requestRole: func(_ context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error) {
			gotWallet = w
			return &identity.AccountRequest{ID: "req-1", Kind: identity.RequestKindRole, Role: wanted, Status: status}, nil
		},
	}
	for want, s := range map[string]string{"attribué": identity.RequestApproved, "en attente": identity.RequestPending} {
		status = s
		var buf bytes.Buffer
		if err := runIdentityRequestRole(&buf, svc, "alice@Org1", "contributor"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotWallet.Handle != "alice@Org1" || !strings.Contains(buf.String(), want) {
			t.Errorf("statut %s : wallet %q, sortie %q", s, gotWallet.Handle, buf.String())
		}
	}
}

/// @brief  runIdentitySetRole change le rôle et rappelle la nécessité d'un ré-enrôlement
/// @input  id="alice@org1", role="contributor"
/// @expect SetRole reçoit les deux valeurs, sortie mentionne le ré-enrôlement
func TestRunIdentitySetRole_NominalCase(t *testing.T) {
	var gotID, gotRole string
	svc := &mockIdentitySvc{setRole: func(_ context.Context, name, newRole string) error {
		gotID, gotRole = name, newRole
		return nil
	}}
	var buf bytes.Buffer
	if err := runIdentitySetRole(&buf, svc, "alice@org1", "contributor"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "alice@org1" || gotRole != "contributor" {
		t.Fatalf("unexpected forwarded args: id=%q role=%q", gotID, gotRole)
	}
	if !strings.Contains(buf.String(), "ré-enrôlement") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

/// @brief  runIdentitySetStatus transmet l'identité et le statut, et rappelle le ré-enrôlement
/// @input  id="alice@Org1", status="active"
/// @expect SetStatus reçoit les deux valeurs, sortie mentionne le ré-enrôlement
func TestRunIdentitySetStatus_NominalCase(t *testing.T) {
	var gotID, gotStatus string
	svc := &mockIdentitySvc{setStatus: func(_ context.Context, id, status string) error {
		gotID, gotStatus = id, status
		return nil
	}}
	var buf bytes.Buffer
	if err := runIdentitySetStatus(&buf, svc, "alice@Org1", "active"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "alice@Org1" || gotStatus != "active" || !strings.Contains(buf.String(), "ré-enrôlement") {
		t.Errorf("SetStatus(%q, %q), sortie %q", gotID, gotStatus, buf.String())
	}
}
