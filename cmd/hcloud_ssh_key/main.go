package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
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
	m.client, err = hcloud.BuildClient(m.args.Token)
	if err != nil {
		return
	}
	if m.args.State == "" {
		m.args.State = statePresent
	}
	if err = validateArgs(m.args); err != nil {
		return
	}

	switch m.args.State {
	case stateList:
		return m.list(ctx)
	case stateAbsent:
		return m.absent(ctx)
	case statePresent:
		return m.present(ctx)
	default:
		err = errors.New("invalid state")
		return
	}
}

func (m *module) Args() interface{} {
	return &m.args
}

// list handles state: list
func (m *module) list(ctx context.Context) (resp ansible.ModuleResponse, err error) {
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

// absent handles state: absent
func (m *module) absent(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var sshKey *hcloud.SSHKey
	if sshKey, err = m.getSSHKey(ctx); err != nil {
		return
	}
	if sshKey == nil {
		resp.Msg("No SSHKey found, nothing to do")
		return
	}

	if _, err = m.client.SSHKey.Delete(ctx, sshKey); err != nil {
		return
	}
	resp.Msg(fmt.Sprintf("SSHKey %d deleted", sshKey.ID)).Changed()
	return
}

// present handles state: present
func (m *module) present(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var sshKey *hcloud.SSHKey
	if sshKey, err = m.getSSHKey(ctx); err != nil {
		return
	}

	var msg []string
	if sshKey != nil {
		var publicKey ssh.PublicKey
		if publicKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(m.args.PublicKey)); err != nil {
			return
		}

		if ssh.FingerprintLegacyMD5(publicKey) != sshKey.Fingerprint {
			if _, err = m.client.SSHKey.Delete(ctx, sshKey); err != nil {
				return
			}
			resp.Changed()
			msg = append(msg, fmt.Sprintf("SSHKey %d deleted (changed fingerprint)", sshKey.ID))
			sshKey = nil
		} else {
			msg = append(msg, fmt.Sprintf("SSHKey %d exists with matching fingerprint", sshKey.ID))
		}
	}

	if sshKey != nil {
		if sshKey.Name != m.args.Name {
			sshKey, _, err = m.client.SSHKey.Update(ctx, sshKey, hcloud.SSHKeyUpdateOpts{
				Name: m.args.Name,
			})
			if err != nil {
				return
			}
		}
	}

	if sshKey == nil {
		sshKey, _, err = m.client.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
			Name:      m.args.Name,
			PublicKey: m.args.PublicKey,
		})
		if err != nil {
			return
		}
		msg = append(msg, fmt.Sprintf("SSHKey %d created", sshKey.ID))
		resp.Changed()
	}

	resp.
		Msg(strings.Join(msg, ", ")).
		Set("ssh_keys", []SSHKey{toSSHKeyData(sshKey)})
	return
}

func (m *module) getSSHKey(ctx context.Context) (sshKey *hcloud.SSHKey, err error) {
	if m.args.ID != 0 {
		if sshKey, _, err = m.client.SSHKey.GetByID(ctx, m.args.ID); err != nil {
			return
		}
		return
	}
	if sshKey, _, err = m.client.SSHKey.GetByName(ctx, m.args.Name); err != nil {
		return
	}
	return
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
		if args.ID != 0 {
			errs = append(errs, "'id' has no effect")
		}
		if args.Name == "" {
			errs = append(errs, "'name' is required")
		}
		if args.PublicKey == "" {
			errs = append(errs, "'public_key' is required")
		}
	}
	if args.State == stateAbsent &&
		args.ID == 0 && args.Name == "" {
		errs = append(errs, "'name' or 'id' is required")
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
