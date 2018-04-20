package hcloudtest

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// ImageClientMock mocks the ImageClient interface
type ImageClientMock struct {
	mock.Mock
}

// NewImageClientMock creates a new ImageClientMock
func NewImageClientMock() hcloud.ImageClient {
	return &ImageClientMock{}
}

// GetByID mock
func (m *ImageClientMock) GetByID(ctx context.Context, id int) (*hcloud.Image, *hcloud.Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*hcloud.Image), args.Get(1).(*hcloud.Response), args.Error(2)
}

// Get mock
func (m *ImageClientMock) Get(ctx context.Context, idOrName string) (*hcloud.Image, *hcloud.Response, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*hcloud.Image), args.Get(1).(*hcloud.Response), args.Error(2)
}

// GetByName mock
func (m *ImageClientMock) GetByName(ctx context.Context, name string) (*hcloud.Image, *hcloud.Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*hcloud.Image), args.Get(1).(*hcloud.Response), args.Error(2)
}

// List mock
func (m *ImageClientMock) List(ctx context.Context, opts hcloud.ImageListOpts) ([]*hcloud.Image, *hcloud.Response, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*hcloud.Image), args.Get(1).(*hcloud.Response), args.Error(2)
}

// All mock
func (m *ImageClientMock) All(ctx context.Context) ([]*hcloud.Image, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*hcloud.Image), args.Error(1)
}

// Delete mock
func (m *ImageClientMock) Delete(ctx context.Context, image *hcloud.Image) (*hcloud.Response, error) {
	args := m.Called(ctx, image)
	return args.Get(0).(*hcloud.Response), args.Error(1)
}
