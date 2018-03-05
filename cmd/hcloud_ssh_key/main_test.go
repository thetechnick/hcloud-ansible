package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud/hcloudtest"
)

func TestValidateArgs(t *testing.T) {
	t.Run("invalid state", func(t *testing.T) {
		err := validateArgs(arguments{
			State: "not-a-valid-state",
		})
		assert.Error(t, err)
	})

	t.Run("present", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			err := validateArgs(arguments{
				State:     statePresent,
				Name:      "my ssh key",
				PublicKey: "---123---",
			})
			assert.NoError(t, err)
		})

		t.Run("id has no effect", func(t *testing.T) {
			err := validateArgs(arguments{
				ID:        123,
				State:     statePresent,
				Name:      "my ssh key",
				PublicKey: "---123---",
			})
			assert.Error(t, err)
		})

		t.Run("missing name", func(t *testing.T) {
			err := validateArgs(arguments{
				State:     statePresent,
				PublicKey: "---123---",
			})
			assert.Error(t, err)
		})

		t.Run("missing public_key", func(t *testing.T) {
			err := validateArgs(arguments{
				State: statePresent,
				ID:    333,
			})
			assert.Error(t, err)
		})
	})

	t.Run("absent", func(t *testing.T) {
		t.Run("success with name", func(t *testing.T) {
			err := validateArgs(arguments{
				State: stateAbsent,
				Name:  "my-ssh-key",
			})
			assert.NoError(t, err)
		})
		t.Run("success with id", func(t *testing.T) {
			err := validateArgs(arguments{
				State: stateAbsent,
				ID:    155,
			})
			assert.NoError(t, err)
		})
		t.Run("missing id and name", func(t *testing.T) {
			err := validateArgs(arguments{
				State: stateAbsent,
			})
			assert.Error(t, err)
		})
	})

	t.Run("list success", func(t *testing.T) {
		err := validateArgs(arguments{
			State: stateList,
		})
		assert.NoError(t, err)
	})
}

func TestList(t *testing.T) {
	client := hcloud.NewClient()
	client.SSHKey = hcloudtest.NewSSHClientMock()

	m := module{
		client: client,
	}

	sshKey := &hcloud.SSHKey{ID: 123}
	sshKeyMock := client.SSHKey.(*hcloudtest.SSHKeyClientMock)
	sshKeyMock.On("All", mock.Anything).Return([]*hcloud.SSHKey{sshKey}, nil)

	resp, err := m.list(context.Background())
	if assert.NoError(t, err) {
		assert.False(t, resp.HasChanged(), "module should not have changed")
		assert.False(t, resp.HasFailed(), "module should not have failed")
		assert.Equal(t, map[string]interface{}{
			"ssh_keys": []SSHKey{toSSHKeyData(sshKey)},
		}, resp.Data())
		sshKeyMock.AssertCalled(t, "All", mock.Anything)
	}
}

func TestAbsent(t *testing.T) {
	t.Run("sshkey exists", func(t *testing.T) {
		client := hcloud.NewClient()
		client.SSHKey = hcloudtest.NewSSHClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		sshKey := &hcloud.SSHKey{ID: 123}
		sshKeyMock := client.SSHKey.(*hcloudtest.SSHKeyClientMock)
		var r *hcloud.Response

		sshKeyMock.On("GetByID", mock.Anything, mock.Anything).Return(sshKey, r, nil)
		sshKeyMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)

		ctx := context.Background()
		resp, err := m.absent(ctx)
		if assert.NoError(t, err) {
			assert.True(t, resp.HasChanged(), "module should have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			sshKeyMock.AssertNumberOfCalls(t, "Delete", 1)
			sshKeyMock.AssertCalled(t, "Delete", mock.Anything, sshKey)
		}
	})

	t.Run("sshkey does not exist", func(t *testing.T) {
		client := hcloud.NewClient()
		client.SSHKey = hcloudtest.NewSSHClientMock()

		m := module{
			client: client,
			args: arguments{
				ID: 123,
			},
		}

		sshKeyMock := client.SSHKey.(*hcloudtest.SSHKeyClientMock)
		var sshKey *hcloud.SSHKey
		var r *hcloud.Response

		sshKeyMock.On("GetByID", mock.Anything, mock.Anything).Return(sshKey, r, nil)
		sshKeyMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)

		ctx := context.Background()
		resp, err := m.absent(ctx)
		if assert.NoError(t, err) {
			assert.False(t, resp.HasChanged(), "module should not have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			sshKeyMock.AssertNotCalled(t, "Delete")
		}
	})
}

var testPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGrhZ30jkQntYeeJNvC5fUVDfi/XmjSRbnOLMLzhyAuq"
var testFingerprint = "a2:94:75:0d:cf:fd:2c:fc:77:81:0e:c6:7a:8d:a2:21"

func TestPresent(t *testing.T) {
	t.Run("sshkey exists with matching fingerprint", func(t *testing.T) {
		client := hcloud.NewClient()
		client.SSHKey = hcloudtest.NewSSHClientMock()

		m := module{
			client: client,
			args: arguments{
				Name:      "my-ssh-key",
				PublicKey: testPublicKey,
			},
		}

		sshKeyMock := client.SSHKey.(*hcloudtest.SSHKeyClientMock)
		sshKey := &hcloud.SSHKey{
			ID:          123,
			Name:        "my-ssh-key",
			Fingerprint: testFingerprint,
		}
		var r *hcloud.Response

		sshKeyMock.On("GetByName", mock.Anything, mock.Anything).Return(sshKey, r, nil)
		sshKeyMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)
		sshKeyMock.On("Create", mock.Anything, mock.Anything).Return(sshKey, r, nil)

		ctx := context.Background()
		resp, err := m.present(ctx)
		if assert.NoError(t, err) {
			assert.False(t, resp.HasChanged(), "module should not have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			sshKeyMock.AssertNotCalled(t, "Delete")
			sshKeyMock.AssertNotCalled(t, "Create")
		}
	})

	t.Run("sshkey exists with invalid fingerprint", func(t *testing.T) {
		client := hcloud.NewClient()
		client.SSHKey = hcloudtest.NewSSHClientMock()

		m := module{
			client: client,
			args: arguments{
				Name:      "my-ssh-key",
				PublicKey: testPublicKey,
			},
		}

		sshKeyMock := client.SSHKey.(*hcloudtest.SSHKeyClientMock)
		sshKey := &hcloud.SSHKey{
			ID:          123,
			Name:        "my-ssh-key",
			Fingerprint: "",
		}
		var r *hcloud.Response

		sshKeyMock.On("GetByName", mock.Anything, mock.Anything).Return(sshKey, r, nil)
		sshKeyMock.On("Delete", mock.Anything, mock.Anything).Return(r, nil)
		sshKeyMock.On("Create", mock.Anything, mock.Anything).Return(sshKey, r, nil)

		ctx := context.Background()
		resp, err := m.present(ctx)
		if assert.NoError(t, err) {
			assert.True(t, resp.HasChanged(), "module should have changed")
			assert.False(t, resp.HasFailed(), "module should not have failed")
			sshKeyMock.AssertCalled(t, "Delete", mock.Anything, sshKey)
			sshKeyMock.AssertCalled(t, "Create", mock.Anything, hcloud.SSHKeyCreateOpts{
				Name:      "my-ssh-key",
				PublicKey: testPublicKey,
			})
		}
	})
}
