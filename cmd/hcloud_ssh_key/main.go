package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
	"golang.org/x/crypto/ssh"
)

type arguments struct {
	Token string `json:"token"`
	State string `json:"state"`

	ID        int    `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

const (
	stateAbsent  = "absent"
	statePresent = "present"
	stateList    = "list"
)

// SSHKey is the module return value of an hcloud.SSHKey
type SSHKey struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type module struct {
	args   arguments
	client *hcloud.Client
}

func (m *module) Run() (resp ansible.ModuleResponse, err error) {
	ctx := context.Background()
	m.client, err = util.BuildClient(m.args.Token)
	if err != nil {
		return
	}
	if m.args.State == "" {
		m.args.State = statePresent
	}
	if err = validateArgs(m.args); err != nil {
		return
	}

	if m.args.State == stateList {
		var (
			list    []SSHKey
			sshKeys []*hcloud.SSHKey
		)
		sshKeys, err = m.client.SSHKey.All(ctx)
		if err != nil {
			return
		}
		for _, sshKey := range sshKeys {
			list = append(list, toSSHKeyData(sshKey))
		}
		resp.Msg("SSHKeys listed").Set("ssh_keys", list)
		return
	}

	var sshKey *hcloud.SSHKey
	if m.args.ID != 0 {
		if sshKey, _, err = m.client.SSHKey.GetByID(ctx, m.args.ID); err != nil {
			return
		}
		if sshKey == nil {
			err = fmt.Errorf("SSHKey %d not found", m.args.ID)
			return
		}
	} else {
		if sshKey, _, err = m.client.SSHKey.Get(ctx, m.args.Name); err != nil {
			return
		}
	}

	switch m.args.State {
	case stateAbsent:
		if sshKey != nil {
			if _, err = m.client.SSHKey.Delete(ctx, sshKey); err != nil {
				return
			}

			resp.Msg(fmt.Sprintf("SSHKey %d deleted", sshKey.ID)).Changed()
			return
		}

		resp.Msg(fmt.Sprintf("No SSHKey `%s` exists", m.args.Name))
		return

	case statePresent:
		var msg []string

		if sshKey != nil {
			var publicKey ssh.PublicKey
			if publicKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(m.args.PublicKey)); err != nil {
				return
			}

			if m.args.Name != "" && sshKey.Name != m.args.Name {
				sshKey, _, err = m.client.SSHKey.Update(ctx, sshKey, hcloud.SSHKeyUpdateOpts{
					Name: m.args.Name,
				})
				if err != nil {
					return
				}
				resp.Changed()
				msg = append(msg, fmt.Sprintf("SSHKey %d name changed to %s", sshKey.ID, sshKey.Name))
			}

			if ssh.FingerprintLegacyMD5(publicKey) != sshKey.Fingerprint {
				if _, err = m.client.SSHKey.Delete(ctx, sshKey); err != nil {
					return
				}
				resp.Changed()
				msg = append(msg, fmt.Sprintf("SSHKey %d deleted", sshKey.ID))
			} else {
				resp.Msg(fmt.Sprintf("SSHKey %d exists and up to date", sshKey.ID))
				resp.Set("ssh_keys", []SSHKey{toSSHKeyData(sshKey)})
				return
			}
		}

		var sshKey *hcloud.SSHKey
		sshKey, _, err = m.client.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
			Name:      m.args.Name,
			PublicKey: m.args.PublicKey,
		})
		if err != nil {
			return
		}

		msg = append(msg, fmt.Sprintf("SSHKey %d created", sshKey.ID))
		resp.Msg(strings.Join(msg, ", ")).Changed().Set("ssh_key", []SSHKey{toSSHKeyData(sshKey)})
		return
	}

	return
}

func (m *module) Args() interface{} {
	return &m.args
}

func toSSHKeyData(key *hcloud.SSHKey) SSHKey {
	return SSHKey{
		ID:          key.ID,
		Name:        key.Name,
		Fingerprint: key.Fingerprint,
	}
}

func validateArgs(args arguments) error {
	errs := []string{}
	if args.State != stateAbsent &&
		args.State != statePresent &&
		args.State != stateList {
		errs = append(errs, "'state' must be present, absent or list")
	}
	if args.State == statePresent {
		if args.ID == 0 && args.Name == "" {
			errs = append(errs, "'name' or 'id' is required")
		}
		if args.PublicKey == "" {
			errs = append(errs, "'public_key' is required")
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}
	return nil
}

var flags = pflag.NewFlagSet("hcloud_ssh_key", pflag.ContinueOnError)

func init() {
	flags.BoolP("version", "v", false, "Print version and exit")
}

func main() {
	ansible.RunModule(&module{}, flags)
}
