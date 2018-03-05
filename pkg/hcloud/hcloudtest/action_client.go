package hcloudtest

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/mock"
	hcloud_wrapped "github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// ActionClientMock mocks the hcloud.ActionClient interface
type ActionClientMock struct {
	mock.Mock
}

// NewActionClientMock creates a new ActionClient mock
func NewActionClientMock() hcloud_wrapped.ActionClient {
	return &ActionClientMock{}
}

// GetByID mock
func (m *ActionClientMock) GetByID(ctx context.Context, id int) (*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// List mock
func (m *ActionClientMock) List(ctx context.Context, opts hcloud.ActionListOpts) ([]*hcloud.Action, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*hcloud.Action), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *ActionClientMock) All(ctx context.Context) ([]*hcloud.Action, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.Action), args.Error(1)
}

// WatchProgress mock
func (m *ActionClientMock) WatchProgress(ctx context.Context, action *hcloud.Action) (<-chan int, <-chan error) {
	args := m.Called(ctx, action)
	return args.Get(0).(<-chan int), args.Get(1).(<-chan error)
}
