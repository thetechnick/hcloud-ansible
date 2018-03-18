package hcloudtest

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// FloatingIPClientMock mock of hcloud.FloatingIPClient
type FloatingIPClientMock struct {
	mock.Mock
}

// NewFloatingIPClientMock creates a FloatingIPClientMock
func NewFloatingIPClientMock() hcloud.FloatingIPClient {
	return &FloatingIPClientMock{}
}

// GetByID mock
func (m *FloatingIPClientMock) GetByID(ctx context.Context, id int) (*hcloud.FloatingIP, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.FloatingIP), args.Get(1).(*hcloud.Response), args.Error(2)
}

// List mock
func (m *FloatingIPClientMock) List(ctx context.Context, opts hcloud.FloatingIPListOpts) ([]*hcloud.FloatingIP, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*hcloud.FloatingIP), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *FloatingIPClientMock) All(ctx context.Context) ([]*hcloud.FloatingIP, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.FloatingIP), args.Error(1)
}

// Create mock
func (m *FloatingIPClientMock) Create(ctx context.Context, opts hcloud.FloatingIPCreateOpts) (hcloud.FloatingIPCreateResult, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(hcloud.FloatingIPCreateResult), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Delete mock
func (m *FloatingIPClientMock) Delete(ctx context.Context, floatingIP *hcloud.FloatingIP) (*hcloud.Response, error) {
	args := m.Called(ctx, floatingIP)
	return args.Get(0).(*hcloud.Response), args.Error(1)
}

// Assign mock
func (m *FloatingIPClientMock) Assign(ctx context.Context, floatingIP *hcloud.FloatingIP, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, floatingIP, server)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Unassign mock
func (m *FloatingIPClientMock) Unassign(ctx context.Context, floatingIP *hcloud.FloatingIP) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, floatingIP)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Update mock
func (m *FloatingIPClientMock) Update(ctx context.Context, floatingIP *hcloud.FloatingIP, opts hcloud.FloatingIPUpdateOpts) (*hcloud.FloatingIP, *hcloud.Response, error) {
	args := m.Called(ctx, floatingIP, opts)
	return args.Get(0).(*hcloud.FloatingIP), args.Get(1).(*hcloud.Response), args.Error(2)
}
