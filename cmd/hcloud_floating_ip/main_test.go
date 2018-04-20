package main

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud/hcloudtest"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
)

var (
	floatingIP = &hcloud.FloatingIP{
		ID:           123,
		IP:           net.ParseIP("192.168.1.1"),
		HomeLocation: &hcloud.Location{Name: "fsn1"},
	}
)

func TestList(t *testing.T) {
	client := hcloud.NewClient()
	client.FloatingIP = hcloudtest.NewFloatingIPClientMock()

	m := module{
		client: client,
	}

	floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
	floatingIPMock.On("All", mock.Anything).Return([]*hcloud.FloatingIP{floatingIP}, nil)

	resp, err := m.list(context.Background())
	if assert.NoError(t, err) {
		assert.False(t, resp.HasChanged(), "module should not have changed")
		assert.False(t, resp.HasFailed(), "module should not have failed")
		assert.Equal(t, map[string]interface{}{
			"floating_ips": []FloatingIP{toFloatingIP(floatingIP)},
		}, resp.Data())
		floatingIPMock.AssertCalled(t, "All", mock.Anything)
	}
}

func TestAbsent(t *testing.T) {
	t.Run("floatingip exists", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
		floatingIPMock.On("GetByID", mock.Anything, mock.Anything).Return(floatingIP, r, nil)
		floatingIPMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)

		resp, err := m.absent(context.Background())
		if assert.NoError(t, err) {
			assert.True(t, resp.HasChanged(), "module should not have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			assert.Equal(t, map[string]interface{}(nil), resp.Data())
			floatingIPMock.AssertCalled(t, "Delete", mock.Anything, floatingIP)
		}
	})

	t.Run("floatingip does not exist", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		var floatingIP *hcloud.FloatingIP
		floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
		floatingIPMock.On("GetByID", mock.Anything, mock.Anything).Return(floatingIP, r, nil)
		floatingIPMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)

		resp, err := m.absent(context.Background())
		if assert.NoError(t, err) {
			assert.False(t, resp.HasChanged(), "module should not have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			assert.Equal(t, map[string]interface{}(nil), resp.Data())
			floatingIPMock.AssertNotCalled(t, "Delete", mock.Anything, floatingIP)
		}
	})
}

func TestServer(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}
		server, err := m.server(context.Background(), nil)
		assert.NoError(t, err)
		assert.Nil(t, server)
	})

	t.Run("not found", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		var rServer *hcloud.Server
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(rServer, r, nil)

		server, err := m.server(context.Background(), "test")
		assert.Error(t, err)
		assert.Nil(t, server)
	})

	t.Run("with string", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(rServer, r, nil)

		server, err := m.server(context.Background(), "my-server")
		assert.NoError(t, err)
		assert.Exactly(t, server, rServer)
		serverClientMock.AssertCalled(t, "GetByName", mock.Anything, "my-server")
	})

	t.Run("with int", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(rServer, r, nil)

		server, err := m.server(context.Background(), 123)
		assert.NoError(t, err)
		assert.Exactly(t, server, rServer)
		serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
	})

	t.Run("with map[string]interface{}", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(rServer, r, nil)

		server, err := m.server(context.Background(), map[string]interface{}{"id": 123, "name": "test"})
		assert.NoError(t, err)
		assert.Exactly(t, server, rServer)
		serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
	})
}

func TestPresent(t *testing.T) {
	t.Run("assign existing", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()
		client.Action = hcloudtest.NewActionClientMock()

		m := module{
			client: client,
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
			args: arguments{
				Token: "--token--",

				State:       "present",
				ID:          123,
				Description: "changed",
				Server:      123,
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		rFloatingIP := &hcloud.FloatingIP{
			ID: 123, Description: "test", Server: nil,
			HomeLocation: &hcloud.Location{},
		}
		rAction := &hcloud.Action{ID: 123, Status: hcloud.ActionStatusSuccess}

		floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
		floatingIPMock.On("GetByID", mock.Anything, mock.Anything).Return(rFloatingIP, r, nil)
		floatingIPMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(rFloatingIP, r, nil)
		floatingIPMock.On("Assign", mock.Anything, mock.Anything, mock.Anything).Return(rAction, r, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(rServer, r, nil)

		resp, err := m.run()
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")

		floatingIPMock.AssertCalled(t, "Update", mock.Anything, rFloatingIP, mock.Anything)
		floatingIPMock.AssertCalled(t, "Assign", mock.Anything, rFloatingIP, rServer)
	})

	t.Run("unassign existing", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()
		client.Action = hcloudtest.NewActionClientMock()

		m := module{
			client: client,
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
			args: arguments{
				Token: "--token--",

				State:        "present",
				ID:           123,
				Description:  "changed",
				HomeLocation: "fsn1",
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		rFloatingIP := &hcloud.FloatingIP{
			ID: 123, Description: "test",
			HomeLocation: &hcloud.Location{},
			Server:       &hcloud.Server{ID: 456},
		}
		rAction := &hcloud.Action{ID: 123, Status: hcloud.ActionStatusSuccess}

		floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
		floatingIPMock.On("GetByID", mock.Anything, mock.Anything).Return(rFloatingIP, r, nil)
		floatingIPMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(rFloatingIP, r, nil)
		floatingIPMock.On("Unassign", mock.Anything, mock.Anything, mock.Anything).Return(rAction, r, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(rServer, r, nil)

		resp, err := m.run()
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")

		floatingIPMock.AssertCalled(t, "Update", mock.Anything, rFloatingIP, mock.Anything)
		floatingIPMock.AssertCalled(t, "Unassign", mock.Anything, rFloatingIP)
	})

	t.Run("create", func(t *testing.T) {
		client := hcloud.NewClient()
		client.FloatingIP = hcloudtest.NewFloatingIPClientMock()
		client.Server = hcloudtest.NewServerClientMock()
		client.Action = hcloudtest.NewActionClientMock()

		m := module{
			client: client,
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
			args: arguments{
				Token: "--token--",

				State:        "present",
				Description:  "test",
				HomeLocation: "fsn1",
			},
		}

		var r *hcloud.Response
		rServer := &hcloud.Server{ID: 123}
		rFloatingIP := &hcloud.FloatingIP{
			ID: 123, Description: "test",
			HomeLocation: &hcloud.Location{},
			Server:       &hcloud.Server{ID: 456},
		}
		rAction := &hcloud.Action{ID: 123, Status: hcloud.ActionStatusSuccess}

		floatingIPMock := client.FloatingIP.(*hcloudtest.FloatingIPClientMock)
		floatingIPMock.On("GetByID", mock.Anything, mock.Anything).Return(rFloatingIP, r, nil)
		floatingIPMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.FloatingIPCreateResult{
			FloatingIP: rFloatingIP,
		}, r, nil)
		floatingIPMock.On("Unassign", mock.Anything, mock.Anything, mock.Anything).Return(rAction, r, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(rServer, r, nil)

		resp, err := m.run()
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")

		floatingIPMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestValidateArgs(t *testing.T) {
	t.Run("invalid state", func(t *testing.T) {
		err := validateArgs(arguments{
			State: "not-a-valid-state",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'state' must be present, absent or list")
	})

	t.Run("absent", func(t *testing.T) {
		t.Run("id missing", func(t *testing.T) {
			err := validateArgs(arguments{
				State: "absent",
			})
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "'id' is required")
			}
		})
	})

	t.Run("present", func(t *testing.T) {
		t.Run("home_location and server missing", func(t *testing.T) {
			err := validateArgs(arguments{
				State: "present",
			})
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "'home_location' or 'server' must be set")
			}
		})

		t.Run("home_location and server specified", func(t *testing.T) {
			err := validateArgs(arguments{
				State:        "present",
				HomeLocation: "nbg1",
				Server:       "test",
			})
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "'home_location' and 'server' are mutually exclusive")
			}
		})
	})
}
