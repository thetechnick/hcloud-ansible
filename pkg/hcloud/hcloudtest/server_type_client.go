package hcloudtest

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// ServerTypeClientMock mocks the ServerTypeClient interface
type ServerTypeClientMock struct {
	mock.Mock
}

// NewServerTypeClientMock creates a new ServerTypeClientMock
func NewServerTypeClientMock() hcloud.ServerTypeClient {
	return &ServerTypeClientMock{}
}

// GetByID mock
func (m *ServerTypeClientMock) GetByID(ctx context.Context, id int) (*hcloud.ServerType, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.ServerType), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Get mock
func (m *ServerTypeClientMock) Get(ctx context.Context, idOrName string) (*hcloud.ServerType, *hcloud.Response, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*hcloud.ServerType), args.Get(1).(*hcloud.Response), args.Error(2)
}

// GetByName mock
func (m *ServerTypeClientMock) GetByName(ctx context.Context, name string) (*hcloud.ServerType, *hcloud.Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*hcloud.ServerType), args.Get(1).(*hcloud.Response), args.Error(2)
}

// List mock
func (m *ServerTypeClientMock) List(ctx context.Context, opts hcloud.ServerTypeListOpts) ([]*hcloud.ServerType, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*hcloud.ServerType), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *ServerTypeClientMock) All(ctx context.Context) ([]*hcloud.ServerType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.ServerType), args.Error(1)
}
