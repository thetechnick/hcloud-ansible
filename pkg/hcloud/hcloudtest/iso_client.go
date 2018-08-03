package hcloudtest

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/mock"
	hcloud_wrapped "github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// ISOClientMock mocks the ISOClient interface
type ISOClientMock struct {
	mock.Mock
}

// NewISOClientMock creates a new ISOClientMock
func NewISOClientMock() hcloud_wrapped.ISOClient {
	return &ISOClientMock{}
}

// GetByID mock
func (m *ISOClientMock) GetByID(ctx context.Context, id int) (*hcloud.ISO, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.ISO), args.Get(1).(*hcloud.Response), args.Error(2)
}

// GetByName mock
func (m *ISOClientMock) GetByName(ctx context.Context, name string) (*hcloud.ISO, *hcloud.Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*hcloud.ISO), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Get mock
func (m *ISOClientMock) Get(ctx context.Context, idOrName string) (*hcloud.ISO, *hcloud.Response, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*hcloud.ISO), args.Get(1).(*hcloud.Response), args.Error(2)
}
