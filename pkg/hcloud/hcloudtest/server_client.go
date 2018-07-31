package hcloudtest

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/mock"
	hcloud_wrapped "github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// ServerClientMock mocks the ServerClient interface
type ServerClientMock struct {
	mock.Mock
}

// NewServerClientMock creates a new ServerClientMock
func NewServerClientMock() hcloud_wrapped.ServerClient {
	return &ServerClientMock{}
}

// GetByID mock
func (m *ServerClientMock) GetByID(ctx context.Context, id int) (*hcloud.Server, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.Server), args.Get(1).(*hcloud.Response), args.Error(2)
}

// GetByName mock
func (m *ServerClientMock) GetByName(ctx context.Context, name string) (*hcloud.Server, *hcloud.Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*hcloud.Server), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Get mock
func (m *ServerClientMock) Get(ctx context.Context, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*hcloud.Server), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *ServerClientMock) All(ctx context.Context) ([]*hcloud.Server, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.Server), args.Error(1)
}

// Create mock
func (m *ServerClientMock) Create(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(hcloud.ServerCreateResult), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Delete mock
func (m *ServerClientMock) Delete(ctx context.Context, server *hcloud.Server) (*hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Response), args.Error(1)
}

// Update mock
func (m *ServerClientMock) Update(ctx context.Context, server *hcloud.Server, opts hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error) {
	args := m.Called(ctx, server, opts)
	return args.Get(0).(*hcloud.Server), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Poweron mock
func (m *ServerClientMock) Poweron(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Reboot mock
func (m *ServerClientMock) Reboot(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Reset mock
func (m *ServerClientMock) Reset(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Poweroff mock
func (m *ServerClientMock) Poweroff(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// EnableRescue mock
func (m *ServerClientMock) EnableRescue(ctx context.Context, server *hcloud.Server, opts hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error) {
	args := m.Called(ctx, server, opts)
	return args.Get(0).(hcloud.ServerEnableRescueResult), args.Get(1).(*hcloud.Response), args.Error(2)
}

// DisableRescue mock
func (m *ServerClientMock) DisableRescue(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// EnableBackup mock
func (m *ServerClientMock) EnableBackup(ctx context.Context, server *hcloud.Server, window string) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// DisableBackup mock
func (m *ServerClientMock) DisableBackup(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// AttachISO mock
func (m *ServerClientMock) AttachISO(ctx context.Context, server *hcloud.Server, iso *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server, iso)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// DetachISO mock
func (m *ServerClientMock) DetachISO(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}
