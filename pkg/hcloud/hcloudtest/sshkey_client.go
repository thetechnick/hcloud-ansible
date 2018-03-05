package hcloudtest

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/mock"
	hcloud_wrapped "github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// SSHKeyClientMock mocks the hcloud.SSHKeyClient interface
type SSHKeyClientMock struct {
	mock.Mock
}

// NewSSHClientMock creates a new SSHKeyClientMock
func NewSSHClientMock() hcloud_wrapped.SSHKeyClient {
	return &SSHKeyClientMock{}
}

// GetByID mock
func (m *SSHKeyClientMock) GetByID(ctx context.Context, id int) (*hcloud.SSHKey, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.SSHKey), args.Get(1).(*hcloud.Response), args.Error(2)
}

// GetByName mock
func (m *SSHKeyClientMock) GetByName(ctx context.Context, name string) (*hcloud.SSHKey, *hcloud.Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*hcloud.SSHKey), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Get mock
func (m *SSHKeyClientMock) Get(ctx context.Context, idOrName string) (*hcloud.SSHKey, *hcloud.Response, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*hcloud.SSHKey), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *SSHKeyClientMock) All(ctx context.Context) ([]*hcloud.SSHKey, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.SSHKey), args.Error(1)
}

// Create mock
func (m *SSHKeyClientMock) Create(ctx context.Context, opts hcloud.SSHKeyCreateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*hcloud.SSHKey), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Delete mock
func (m *SSHKeyClientMock) Delete(ctx context.Context, sshKey *hcloud.SSHKey) (*hcloud.Response, error) {
	args := m.Called(ctx, sshKey)
	return args.Get(0).(*hcloud.Response), args.Error(1)
}

// Update mock
func (m *SSHKeyClientMock) Update(ctx context.Context, sshKey *hcloud.SSHKey, opts hcloud.SSHKeyUpdateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	args := m.Called(ctx, sshKey, opts)
	return args.Get(0).(*hcloud.SSHKey), args.Get(1).(*hcloud.Response), args.Error(2)
}
