package main

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud/hcloudtest"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
)

var (
	server      *hcloud.Server
	image       *hcloud.Image
	iso         *hcloud.ISO
	nilServer   *hcloud.Server
	nilResponse *hcloud.Response
)

func init() {
	i, n, _ := net.ParseCIDR("2001:db8::/64")
	image = &hcloud.Image{ID: 123, Name: "debian-9"}

	server = &hcloud.Server{
		ID:   123,
		Name: "test",
		ServerType: &hcloud.ServerType{
			Name: "cx11",
		},
		Datacenter: &hcloud.Datacenter{Location: &hcloud.Location{}},
		Image:      image,
		Status:     hcloud.ServerStatusRunning,
		PublicNet: hcloud.ServerPublicNet{
			IPv4: hcloud.ServerPublicNetIPv4{
				IP: net.ParseIP("192.168.1.2"),
			},
			IPv6: hcloud.ServerPublicNetIPv6{
				IP:      i,
				Network: n,
			},
		},
	}

	iso = &hcloud.ISO{
		ID:   456,
		Name: "test.iso",
	}
}

func TestList(t *testing.T) {
	t.Run("with id", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				Token: "--token--",
				State: stateList,
				ID:    "123",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()
		var (
			response *hcloud.Response
		)
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)

		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(server, response, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.False(t, resp.HasChanged(), "should not have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(server)},
		}, resp.Data())
	})

	t.Run("without params", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				Token: "--token--",
				State: stateList,
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)

		serverClientMock.On("All", mock.Anything).Return([]*hcloud.Server{server}, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.False(t, resp.HasChanged(), "should not have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		serverClientMock.AssertCalled(t, "All", mock.Anything)
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(server)},
		}, resp.Data())
	})
}

func TestPresent(t *testing.T) {
	client := hcloud.NewClient()
	client.Server = hcloudtest.NewServerClientMock()
	client.Image = hcloudtest.NewImageClientMock()
	client.ServerType = hcloudtest.NewServerTypeClientMock()

	m := module{
		client: client,
		args: arguments{
			Token:      "--token--",
			State:      statePresent,
			Name:       "test",
			Image:      "debian-9",
			ServerType: "cx11",
		},
		waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
			return nil
		}),
	}

	ctx := context.Background()

	imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
	imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

	serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
	serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(nilServer, nilResponse, nil).Once()
	serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(server, nilResponse, nil)
	serverClientMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.ServerCreateResult{
		Server: server,
		Action: &hcloud.Action{ID: 123},
	}, nilResponse, nil)

	resp, err := m.run(ctx)
	assert.NoError(t, err)
	assert.True(t, resp.HasChanged(), "should have changed")
	assert.False(t, resp.HasFailed(), "should not have failed")
	assert.Equal(t, map[string]interface{}{
		"servers": []Server{toServer(server)},
	}, resp.Data())

	t.Run("attach ISO", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()
		client.Image = hcloudtest.NewImageClientMock()
		client.ServerType = hcloudtest.NewServerTypeClientMock()
		client.ISO = hcloudtest.NewISOClientMock()

		m := module{
			client: client,
			args: arguments{
				Token:      "--token--",
				State:      statePresent,
				Name:       "test",
				Image:      "debian-9",
				ServerType: "cx11",
				ISO:        "test.iso",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()
		server := *server

		imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
		imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(&server, nilResponse, nil)
		serverClientMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.ServerCreateResult{
			Server: &server,
			Action: &hcloud.Action{ID: 123},
		}, nilResponse, nil)
		serverClientMock.On("AttachISO", mock.Anything, mock.Anything, mock.Anything).Return(&hcloud.Action{ID: 123}, nilResponse, nil)

		isoClientMock := client.ISO.(*hcloudtest.ISOClientMock)
		isoClientMock.On("GetByName", mock.Anything, mock.Anything).Return(iso, nilResponse, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(&server)},
		}, resp.Data())
	})

	t.Run("detach ISO", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()
		client.Image = hcloudtest.NewImageClientMock()
		client.ServerType = hcloudtest.NewServerTypeClientMock()
		client.ISO = hcloudtest.NewISOClientMock()

		m := module{
			client: client,
			args: arguments{
				Token:      "--token--",
				State:      statePresent,
				Name:       "test",
				Image:      "debian-9",
				ServerType: "cx11",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()
		server := *server
		server.ISO = iso

		imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
		imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(&server, nilResponse, nil)
		serverClientMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.ServerCreateResult{
			Server: &server,
			Action: &hcloud.Action{ID: 123},
		}, nilResponse, nil)
		serverClientMock.On("DetachISO", mock.Anything, mock.Anything).Return(&hcloud.Action{ID: 123}, nilResponse, nil)

		isoClientMock := client.ISO.(*hcloudtest.ISOClientMock)
		isoClientMock.On("GetByName", mock.Anything, mock.Anything).Return(iso, nilResponse, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(&server)},
		}, resp.Data())
	})
}

func TestRunning(t *testing.T) {
	client := hcloud.NewClient()
	client.Server = hcloudtest.NewServerClientMock()
	client.Image = hcloudtest.NewImageClientMock()
	client.ServerType = hcloudtest.NewServerTypeClientMock()

	m := &module{
		client: client,
		args: arguments{
			Token:      "--token--",
			State:      stateRunning,
			Name:       "test",
			Image:      "debian-9",
			ServerType: "cx11",
		},
		waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
			return nil
		}),
	}

	ctx := context.Background()

	server := *server
	server.Status = hcloud.ServerStatusOff

	imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
	imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

	serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
	serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(nilServer, nilResponse, nil).Once()
	serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(&server, nilResponse, nil)
	serverClientMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.ServerCreateResult{
		Server: &server,
		Action: &hcloud.Action{ID: 123},
	}, nilResponse, nil)
	serverClientMock.On("Poweron", mock.Anything, mock.Anything).Return(&hcloud.Action{}, nilResponse, nil)

	resp, err := m.run(ctx)
	assert.NoError(t, err)
	assert.True(t, resp.HasChanged(), "should have changed")
	assert.False(t, resp.HasFailed(), "should not have failed")
	assert.Equal(t, map[string]interface{}{
		"servers": []Server{toServer(&server)},
	}, resp.Data())
	serverClientMock.AssertCalled(t, "Poweron", mock.Anything, &server)
}

func TestStopped(t *testing.T) {
	t.Run("server absent", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()
		client.Image = hcloudtest.NewImageClientMock()
		client.ServerType = hcloudtest.NewServerTypeClientMock()

		m := &module{
			client: client,
			args: arguments{
				Token:      "--token--",
				State:      stateStopped,
				Name:       "test",
				Image:      "debian-9",
				ServerType: "cx11",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()

		server := *server
		server.Status = hcloud.ServerStatusOff

		imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
		imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(nilServer, nilResponse, nil).Once()
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(&server, nilResponse, nil)
		serverClientMock.On("Create", mock.Anything, mock.Anything).Return(hcloud.ServerCreateResult{
			Server: &server,
			Action: &hcloud.Action{ID: 123},
		}, nilResponse, nil)
		serverClientMock.On("Poweroff", mock.Anything, mock.Anything).Return(&hcloud.Action{}, nilResponse, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(&server)},
		}, resp.Data())
		serverClientMock.AssertCalled(t, "Create", mock.Anything, hcloud.ServerCreateOpts{
			Name: "test",
			ServerType: &hcloud.ServerType{
				Name: "cx11",
			},
			Image:            image,
			StartAfterCreate: hcloud.Bool(false),
		})
		serverClientMock.AssertNotCalled(t, "Poweroff", mock.Anything, &server)
	})

	t.Run("server absent by id", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()
		client.Image = hcloudtest.NewImageClientMock()
		client.ServerType = hcloudtest.NewServerTypeClientMock()

		m := &module{
			client: client,
			args: arguments{
				Token: "--token--",
				State: stateStopped,
				ID:    "1",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()

		server := *server
		server.Status = hcloud.ServerStatusOff
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(nilServer, nilResponse, nil)

		resp, err := m.run(ctx)
		assert.Error(t, err)
		assert.False(t, resp.HasFailed(), "should not have failed") // ?
		assert.Equal(t, fmt.Errorf("Server with id 1 not found"), err)
		//		"servers": []Server{toServer(&server)},
		//	}, resp.Data())
		serverClientMock.AssertNotCalled(t, "Poweroff", mock.Anything, &server)
	})

	t.Run("server running", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()
		client.Image = hcloudtest.NewImageClientMock()
		client.ServerType = hcloudtest.NewServerTypeClientMock()

		m := &module{
			client: client,
			args: arguments{
				Token:      "--token--",
				State:      stateStopped,
				Name:       "test",
				Image:      "debian-9",
				ServerType: "cx11",
			},
			waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
				return nil
			}),
		}

		ctx := context.Background()

		server := *server
		server.Status = hcloud.ServerStatusRunning

		imageClientMock := client.Image.(*hcloudtest.ImageClientMock)
		imageClientMock.On("GetByName", mock.Anything, mock.Anything).Return(image, nilResponse, nil)

		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByName", mock.Anything, mock.Anything).Return(&server, nilResponse, nil)
		serverClientMock.On("Poweroff", mock.Anything, mock.Anything).Return(&hcloud.Action{}, nilResponse, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.True(t, resp.HasChanged(), "should have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		assert.Equal(t, map[string]interface{}{
			"servers": []Server{toServer(&server)},
		}, resp.Data())
		serverClientMock.AssertCalled(t, "Poweroff", mock.Anything, &server)
	})

}

func TestAbsent(t *testing.T) {
	client := hcloud.NewClient()
	client.Server = hcloudtest.NewServerClientMock()

	m := module{
		client: client,
		args: arguments{
			Token: "--token--",
			State: stateAbsent,
			ID:    "123",
		},
	}

	ctx := context.Background()
	var response *hcloud.Response
	serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
	serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(server, response, nil)
	serverClientMock.On("Delete", mock.Anything, mock.Anything).Return(response, nil)

	resp, err := m.run(ctx)
	assert.NoError(t, err)
	assert.True(t, resp.HasChanged(), "should have changed")
	assert.False(t, resp.HasFailed(), "should not have failed")
	serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
	serverClientMock.AssertCalled(t, "Delete", mock.Anything, server)
	assert.Equal(t, map[string]interface{}(nil), resp.Data())

	t.Run("server not found", func(t *testing.T) {
		client := hcloud.NewClient()
		client.Server = hcloudtest.NewServerClientMock()

		m := module{
			client: client,
			args: arguments{
				Token: "--token--",
				State: stateAbsent,
				ID:    "123",
			},
		}

		ctx := context.Background()
		var response *hcloud.Response
		serverClientMock := client.Server.(*hcloudtest.ServerClientMock)
		serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(nilServer, response, nil)
		serverClientMock.On("Delete", mock.Anything, mock.Anything).Return(response, nil)

		resp, err := m.run(ctx)
		assert.NoError(t, err)
		assert.False(t, resp.HasChanged(), "should not have changed")
		assert.False(t, resp.HasFailed(), "should not have failed")
		serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
		// serverClientMock.AssertCalled(t, "Delete", mock.Anything, server)
		assert.Equal(t, map[string]interface{}(nil), resp.Data())
	})
}

func TestRestarted(t *testing.T) {
	client := hcloud.NewClient()
	client.Server = hcloudtest.NewServerClientMock()

	m := module{
		client: client,
		args: arguments{
			Token: "--token--",
			State: stateRestarted,
			ID:    "123",
		},
		waitFn: util.WaitFn(func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
			return nil
		}),
	}

	ctx := context.Background()
	server := *server
	server.Status = hcloud.ServerStatusRunning
	var (
		response *hcloud.Response
		action   *hcloud.Action
	)
	serverClientMock := client.Server.(*hcloudtest.ServerClientMock)

	serverClientMock.On("GetByID", mock.Anything, mock.Anything).Return(&server, response, nil)
	serverClientMock.On("Reboot", mock.Anything, mock.Anything).Return(action, response, nil)

	resp, err := m.run(ctx)
	assert.NoError(t, err)
	assert.True(t, resp.HasChanged(), "should have changed")
	assert.False(t, resp.HasFailed(), "should not have failed")
	serverClientMock.AssertCalled(t, "GetByID", mock.Anything, 123)
	serverClientMock.AssertCalled(t, "Reboot", mock.Anything, &server)
	assert.Equal(t, map[string]interface{}{
		"servers": []Server{toServer(&server)},
	}, resp.Data())
}

func TestValidateState(t *testing.T) {
	valid := []string{
		statePresent,
		stateAbsent,
		stateList,
		stateRestarted,
		stateRunning,
		stateStopped,
	}
	invalid := []string{
		"hans",
		"123",
	}
	for _, state := range valid {
		t.Run(fmt.Sprintf("valid - %q", state), func(t *testing.T) {
			assert.NoError(t, validateState(state))
		})
	}
	for _, state := range invalid {
		t.Run(fmt.Sprintf("invalid - %q", state), func(t *testing.T) {
			assert.Error(t, validateState(state))
		})
	}
}

func TestArgsToConfig(t *testing.T) {
	client := hcloud.NewClient()
	client.ISO = hcloudtest.NewISOClientMock()
	isoClientMock := client.ISO.(*hcloudtest.ISOClientMock)
	isoClientMock.On("GetByName", mock.Anything, mock.Anything).Return(iso, nilResponse, nil)

	m := &module{
		client: client,
		args: arguments{
			State:      "present",
			Token:      "--token--",
			ID:         1,
			ServerType: "cx11",
			UserData:   "--user data--",
			Rescue:     "linux64",
			ISO:        "test.iso",
		},
	}

	c, err := m.argsToConfig(context.Background())

	if assert.NoError(t, err) {
		assert.Equal(t, "present", c.State)
		assert.Equal(t, "--token--", c.Token)
		assert.Equal(t, "cx11", c.ServerType)
		assert.Equal(t, "--user data--", c.UserData)
		assert.Equal(t, "linux64", c.Rescue)
		assert.Equal(t, iso, c.ISO)
	}
	isoClientMock.AssertCalled(t, "GetByName", mock.Anything, "test.iso")
}
