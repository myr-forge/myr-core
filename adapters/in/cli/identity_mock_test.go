// adapters/in/cli/identity_mock_test.go — mock partagé IdentityService (Pattern B — fn-func)
package cli

import (
	"context"

	"myr-core/domain/identity"
)

type mockIdentitySvc struct {
	listLocalWallets func() ([]identity.WalletEntry, error)
	wallet           func(handle string) (identity.WalletEntry, error)
	createAccount    func(ctx context.Context, pseudo, orgID, role, email string) (string, string, error)
	enroll           func(ctx context.Context, name, secret, orgID string) (identity.WalletEntry, error)
	getStatus        func(ctx context.Context, wallet identity.WalletEntry) (string, error)
	reEnroll         func(ctx context.Context, wallet identity.WalletEntry) (identity.WalletEntry, error)
	loadGuestWallet  func(certPath, keyPath, caCertPath, orgID string) (identity.WalletEntry, error)
	requestAccess    func(ctx context.Context, req identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error)
	requestRole      func(ctx context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error)
	approveRequest   func(ctx context.Context, id string) (*identity.AccountRequest, error)
	rejectRequest    func(ctx context.Context, id string) (*identity.AccountRequest, error)
	listRequests     func() ([]*identity.AccountRequest, error)
	walletDir        func() string
	setRole          func(ctx context.Context, name, newRole string) error
	setStatus        func(ctx context.Context, name, status string) error
}

func (m *mockIdentitySvc) ListLocalWallets() ([]identity.WalletEntry, error) {
	if m.listLocalWallets != nil {
		return m.listLocalWallets()
	}
	return nil, nil
}
func (m *mockIdentitySvc) Wallet(handle string) (identity.WalletEntry, error) {
	if m.wallet != nil {
		return m.wallet(handle)
	}
	return identity.WalletEntry{}, identity.ErrWalletNotFound
}
func (m *mockIdentitySvc) CreateAccount(ctx context.Context, pseudo, orgID, role, email string) (string, string, error) {
	if m.createAccount != nil {
		return m.createAccount(ctx, pseudo, orgID, role, email)
	}
	return identity.Handle(pseudo, orgID), "secret-test", nil
}
func (m *mockIdentitySvc) Enroll(ctx context.Context, name, secret, orgID string) (identity.WalletEntry, error) {
	if m.enroll != nil {
		return m.enroll(ctx, name, secret, orgID)
	}
	return identity.WalletEntry{Handle: name + "@" + orgID, Name: name, OrgID: orgID, Status: identity.StatusPending}, nil
}
func (m *mockIdentitySvc) GetStatus(ctx context.Context, wallet identity.WalletEntry) (string, error) {
	if m.getStatus != nil {
		return m.getStatus(ctx, wallet)
	}
	return identity.StatusActive, nil
}
func (m *mockIdentitySvc) ReEnroll(ctx context.Context, wallet identity.WalletEntry) (identity.WalletEntry, error) {
	if m.reEnroll != nil {
		return m.reEnroll(ctx, wallet)
	}
	return wallet, nil
}
func (m *mockIdentitySvc) LoadGuestWallet(certPath, keyPath, caCertPath, orgID string) (identity.WalletEntry, error) {
	if m.loadGuestWallet != nil {
		return m.loadGuestWallet(certPath, keyPath, caCertPath, orgID)
	}
	return identity.WalletEntry{}, nil
}
func (m *mockIdentitySvc) RequestAccess(ctx context.Context, req identity.AccountRequest, policy identity.RegistrationPolicy) (*identity.AccountRequest, string, error) {
	if m.requestAccess != nil {
		return m.requestAccess(ctx, req, policy)
	}
	if !policy.AllowAutoRegister {
		return nil, "", identity.ErrRegistrationClosed
	}
	req.ID = "req-test"
	req.Status = identity.RequestApproved
	if policy.RequireApproval {
		req.Status = identity.RequestPending
	}
	return &req, "secret-test", nil
}
func (m *mockIdentitySvc) ApproveRequest(ctx context.Context, id string) (*identity.AccountRequest, error) {
	if m.approveRequest != nil {
		return m.approveRequest(ctx, id)
	}
	return &identity.AccountRequest{ID: id, Status: identity.RequestApproved}, nil
}
func (m *mockIdentitySvc) RejectRequest(ctx context.Context, id string) (*identity.AccountRequest, error) {
	if m.rejectRequest != nil {
		return m.rejectRequest(ctx, id)
	}
	return &identity.AccountRequest{ID: id, Status: identity.RequestRejected}, nil
}
func (m *mockIdentitySvc) ListRequests() ([]*identity.AccountRequest, error) {
	if m.listRequests != nil {
		return m.listRequests()
	}
	return nil, nil
}
func (m *mockIdentitySvc) WalletDir() string {
	if m.walletDir != nil {
		return m.walletDir()
	}
	return ""
}
func (m *mockIdentitySvc) RequestRole(ctx context.Context, w identity.WalletEntry, wanted string) (*identity.AccountRequest, error) {
	if m.requestRole != nil {
		return m.requestRole(ctx, w, wanted)
	}
	return &identity.AccountRequest{ID: "req-1", Kind: identity.RequestKindRole, Role: wanted, Status: identity.RequestPending}, nil
}
func (m *mockIdentitySvc) SetStatus(ctx context.Context, name, status string) error {
	if m.setStatus != nil {
		return m.setStatus(ctx, name, status)
	}
	return nil
}
func (m *mockIdentitySvc) SetRole(ctx context.Context, name, newRole string) error {
	if m.setRole != nil {
		return m.setRole(ctx, name, newRole)
	}
	return nil
}
