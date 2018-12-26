package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
)

type arguments struct {
	Token string `json:"token"`
	State string `json:"state"`

	ID         interface{} `json:"id"`
	Name       interface{} `json:"name"`
	Image      interface{} `json:"image"`
	ServerType string      `json:"server_type"`
	UserData   string      `json:"user_data"`
	Datacenter string      `json:"datacenter"`
	Location   string      `json:"location"`
	Rescue     string      `json:"rescue"`
	SSHKeys    interface{} `json:"ssh_keys"`
	ISO        interface{} `json:"iso"`
}

const (
	stateAbsent    = "absent"
	statePresent   = "present"
	stateList      = "list"
	stateRunning   = "running"
	stateStopped   = "stopped"
	stateRestarted = "restarted"
)

type config struct {
	Token string
	State string

	Name       []string
	ID         []int
	Image      *hcloud.Image
	ISO        *hcloud.ISO
	ServerType string
	UserData   string
	Datacenter *hcloud.Datacenter
	Location   *hcloud.Location
	Rescue     string
	SSHKeys    []*hcloud.SSHKey
}

// Server is the module return value of an hcloud.Server
type Server struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ServerType string `json:"server_type"`
	Status     string `json:"status"`
	Datacenter string `json:"datacenter"`
	Location   string `json:"location"`
	ISO        string `json:"iso"`
	PublicIPv4 string `json:"public_ipv4"`
	PublicIPv6 string `json:"public_ipv6"`
}

type module struct {
	args     arguments
	config   config
	client   *hcloud.Client
	waitFn   util.WaitFn
	messages ansible.MessageLog
}

func (m *module) Args() interface{} {
	return &m.args
}

func (m *module) Run() (resp ansible.ModuleResponse, err error) {
	if m.client, err = hcloud.BuildClient(m.args.Token); err != nil {
		return
	}
	return m.run(context.Background())
}

func (m *module) run(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	if m.config, err = m.argsToConfig(ctx); err != nil {
		return
	}
	if m.args.State == "" {
		m.args.State = statePresent
	}

	switch m.config.State {
	case stateAbsent:
		return m.absent(ctx)
	case statePresent:
		return m.present(ctx)
	case stateList:
		return m.list(ctx)
	case stateRunning:
		return m.running(ctx)
	case stateStopped:
		return m.stopped(ctx)
	case stateRestarted:
		return m.restarted(ctx)
	default:
		err = errors.New("invalid state")
		return
	}
}

func (m *module) absent(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var servers []*hcloud.Server
	servers, err = m.servers(ctx)
	for _, server := range servers {
		if _, err = m.client.Server.Delete(ctx, server); err != nil {
			return
		}
		resp.Changed()
	}
	return
}

func (m *module) present(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var (
		errors     []string
		errorsLock sync.Mutex
		wg         sync.WaitGroup
	)
	for _, name := range m.config.Name {
		wg.Add(1)
		go func(ctx context.Context, name string) (err error) {
			defer wg.Done()
			defer func() {
				if err != nil {
					errorsLock.Lock()
					errors = append(errors, err.Error())
					errorsLock.Unlock()
				}
			}()
			var server *hcloud.Server
			if server, err = m.ensureServerExists(ctx, &resp, name); err != nil {
				return
			}
			if err = m.ensureServerState(ctx, &resp, server, name); err != nil {
				return
			}
			return
		}(ctx, name)
	}
	for _, id := range m.config.ID {
		wg.Add(1)
		go func(ctx context.Context, id int) (err error) {
			defer wg.Done()
			defer func() {
				if err != nil {
					errorsLock.Lock()
					errors = append(errors, err.Error())
					errorsLock.Unlock()
				}
			}()
			var server *hcloud.Server
			if server, _, err = m.client.Server.GetByID(ctx, id); err != nil {
				return
			}
			if server == nil {
				err = fmt.Errorf("Server with id %d not found", id)
				return
			}
			if err = m.ensureServerState(ctx, &resp, server, ""); err != nil {
				return
			}
			return
		}(ctx, id)
	}
	wg.Wait()

	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, ", "))
		return
	}

	if err = m.output(ctx, &resp); err != nil {
		return
	}
	return
}

func (m *module) servers(ctx context.Context) (servers []*hcloud.Server, err error) {
	if len(m.config.Name) == 0 && len(m.config.ID) == 0 {
		err = fmt.Errorf("'name' or 'id' is required")
		return
	}

	for _, id := range m.config.ID {
		var server *hcloud.Server
		if server, _, err = m.client.Server.GetByID(ctx, id); err != nil {
			return
		}
		if server == nil && m.config.State != stateAbsent {
			err = fmt.Errorf("Server with id %q not found", id)
			return
		}
		if server != nil {
			servers = append(servers, server)
		}
	}
	for _, name := range m.config.Name {
		var server *hcloud.Server
		if server, _, err = m.client.Server.GetByName(ctx, name); err != nil {
			return
		}
		if server == nil && m.config.State != stateAbsent {
			err = fmt.Errorf("Server with name %q not found", name)
			return
		}
		if server != nil {
			servers = append(servers, server)
		}
	}
	return
}

func (m *module) ensureServerExists(ctx context.Context, resp *ansible.ModuleResponse, name string) (server *hcloud.Server, err error) {
	if server, _, err = m.client.Server.GetByName(ctx, name); err != nil {
		return
	}

	if needsRecreate(server, m.config) {
		if _, err = m.client.Server.Delete(ctx, server); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d deleted (needs recreate)", server.ID))
		server = nil
		resp.Changed()
	}

	if server == nil {
		var errs []string
		if m.config.Image == nil {
			errs = append(errs, "'image' is required")
		}
		if m.config.ServerType == "" {
			errs = append(errs, "'server_type' is required")
		}
		if len(errs) > 0 {
			return nil, fmt.Errorf("Cannot create server '%s': %s", name, strings.Join(errs, ", "))
		}

		resp.Changed()
		opts := hcloud.ServerCreateOpts{
			Name: name,

			ServerType: &hcloud.ServerType{
				Name: m.config.ServerType,
			},
			UserData: m.config.UserData,
			SSHKeys:  m.config.SSHKeys,
			Image:    m.config.Image,
		}
		if m.config.State == stateStopped {
			opts.StartAfterCreate = hcloud.Bool(false)
		}
		if m.config.Datacenter != nil {
			opts.Datacenter = m.config.Datacenter
		}
		if m.config.Location != nil {
			opts.Location = m.config.Location
		}

		var res hcloud.ServerCreateResult
		res, _, err = m.client.Server.Create(ctx, opts)
		if err != nil {
			return
		}
		if err = m.waitFn(ctx, m.client, res.Action); err != nil {
			return nil, err
		}
		server = res.Server
	}
	return
}

func (m *module) ensureServerState(ctx context.Context, resp *ansible.ModuleResponse, server *hcloud.Server, name string) (err error) {
	// mount/dismount ISO BEFORE changing the power state of the server.
	// this allowes to boot from the ISO in one step and
	// prevents the server booting from the ISO if it is detached and restarted in one step
	if server.ISO != nil &&
		(m.config.ISO == nil || m.config.ISO.ID != server.ISO.ID) {
		var action *hcloud.Action
		if action, _, err = m.client.Server.DetachISO(ctx, server); err != nil {
			return
		}
		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d ISO %d detached", server.ID, server.ISO.ID))
		resp.Changed()
		server.ISO = nil
	}
	if server.ISO == nil && m.config.ISO != nil {
		var action *hcloud.Action
		if action, _, err = m.client.Server.AttachISO(ctx, server, m.config.ISO); err != nil {
			return
		}

		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d ISO %d attached", server.ID, m.config.ISO.ID))
		resp.Changed()
		server.ISO = m.config.ISO
	}

	switch m.config.State {
	case stateRunning, stateRestarted:
		if server.Status != hcloud.ServerStatusRunning {
			var action *hcloud.Action
			if action, _, err = m.client.Server.Poweron(ctx, server); err != nil {
				return
			}

			if err = m.waitFn(ctx, m.client, action); err != nil {
				return
			}
			m.messages.Add(fmt.Sprintf("Server %d started", server.ID))
			resp.Changed()
		} else if m.config.State == stateRestarted {
			var action *hcloud.Action
			if action, _, err = m.client.Server.Reboot(ctx, server); err != nil {
				return
			}
			if err = m.waitFn(ctx, m.client, action); err != nil {
				return
			}
			m.messages.Add(fmt.Sprintf("Server %d restarted", server.ID))
			resp.Changed()
		}

	case stateStopped:
		if server.Status != hcloud.ServerStatusOff {
			var action *hcloud.Action
			if action, _, err = m.client.Server.Poweroff(ctx, server); err != nil {
				return
			}

			if err = m.waitFn(ctx, m.client, action); err != nil {
				return
			}
			m.messages.Add(fmt.Sprintf("Server %d stopped", server.ID))
			resp.Changed()
		}
	}

	if name != "" && server.Name != name {
		server, _, err = m.client.Server.Update(ctx, server, hcloud.ServerUpdateOpts{
			Name: name,
		})
		if err != nil {
			return
		}
	}

	var rescueChanged bool
	if server.RescueEnabled && m.config.Rescue == "" {
		var action *hcloud.Action
		if action, _, err = m.client.Server.DisableRescue(ctx, server); err != nil {
			return
		}

		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d disabled rescue mode", server.ID))
		resp.Changed()
		rescueChanged = true
	}
	if !server.RescueEnabled && m.config.Rescue != "" {
		var res hcloud.ServerEnableRescueResult
		res, _, err = m.client.Server.EnableRescue(ctx, server, hcloud.ServerEnableRescueOpts{
			Type:    hcloud.ServerRescueType(m.config.Rescue),
			SSHKeys: m.config.SSHKeys,
		})
		if err != nil {
			return
		}

		if err = m.waitFn(ctx, m.client, res.Action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d enabled rescue mode", server.ID))
		resp.Changed()
		rescueChanged = true
	}

	if rescueChanged && m.config.State != stateStopped {
		var action *hcloud.Action
		if action, _, err = m.client.Server.Reset(ctx, server); err != nil {
			return
		}
		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}
	}

	return
}

func (m *module) output(ctx context.Context, resp *ansible.ModuleResponse) (err error) {
	var (
		servers []*hcloud.Server
		s       []Server
	)
	if servers, err = m.servers(ctx); err != nil {
		return
	}
	for _, server := range servers {
		s = append(s, toServer(server))
	}
	resp.Set("servers", s)
	return
}

func (m *module) list(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var servers []*hcloud.Server
	if len(m.config.ID) != 0 || len(m.config.Name) != 0 {
		if servers, err = m.servers(ctx); err != nil {
			return
		}
	} else {
		if servers, err = m.client.Server.All(ctx); err != nil {
			return
		}
	}

	var s []Server
	for _, server := range servers {
		s = append(s, toServer(server))
	}
	resp.Set("servers", s)
	return
}

func (m *module) running(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	return m.present(ctx)
}

func (m *module) stopped(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	return m.present(ctx)
}

func (m *module) restarted(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	return m.present(ctx)
}

// needsRecreate checks if the server needs to be recreated
func needsRecreate(server *hcloud.Server, config config) bool {
	if server == nil {
		return false
	}
	if config.Image != nil &&
		server.Image.ID != config.Image.ID {
		return true
	}
	if config.ServerType != "" &&
		server.ServerType.Name != config.ServerType {
		return true
	}
	if config.Datacenter != nil &&
		server.Datacenter.Name != config.Datacenter.Name {
		return true
	}
	if config.Location != nil &&
		server.Datacenter.Location.Name != config.Location.Name {
		return true
	}

	return false
}

func validateState(state string) error {
	if state != statePresent &&
		state != stateAbsent &&
		state != stateList &&
		state != stateRestarted &&
		state != stateRunning &&
		state != stateStopped {
		return fmt.Errorf("'state' must be present, absent, running, stopped, list or restarted")
	}
	return nil
}

func (m *module) argsToConfig(ctx context.Context) (c config, err error) {
	c.Token = m.args.Token

	if m.args.State == "" {
		m.args.State = statePresent
	}
	if err = validateState(m.args.State); err != nil {
		return
	}
	c.State = m.args.State

	c.Name = util.GetNames(m.args.Name)
	c.ID = util.GetIDs(m.args.ID)
	c.ServerType = m.args.ServerType
	c.UserData = m.args.UserData
	c.Rescue = m.args.Rescue

	// Image
	if imageID := util.GetID(m.args.Image); imageID != 0 {
		c.Image, _, err = m.client.Image.GetByID(ctx, imageID)
		if err != nil {
			return
		}
		if c.Image == nil {
			err = fmt.Errorf("requested image with id %d not found", imageID)
			return
		}
	}
	if imageName := util.GetName(m.args.Image); imageName != "" {
		c.Image, _, err = m.client.Image.GetByName(ctx, imageName)
		if err != nil {
			return
		}
		if c.Image == nil {
			err = fmt.Errorf("requested image with name %s not found", imageName)
			return
		}
	}
	if c.Image == nil && m.args.Image != nil {
		err = fmt.Errorf("image unknown format: %v", m.args.Image)
		return
	}

	// ISO
	if m.args.ISO != nil {
		if isoID := util.GetID(m.args.ISO); isoID != 0 {
			c.ISO, _, err = m.client.ISO.GetByID(ctx, isoID)
			if err != nil {
				return
			}
			if c.ISO == nil {
				err = fmt.Errorf("requested ISO with id %d not found", isoID)
				return
			}
		}
		if isoName := util.GetName(m.args.ISO); isoName != "" {
			c.ISO, _, err = m.client.ISO.GetByName(ctx, isoName)
			if err != nil {
				return
			}
			if c.ISO == nil {
				err = fmt.Errorf("requested ISO with name %s not found", isoName)
				return
			}
		}
		if c.ISO == nil && m.args.ISO != nil {
			err = fmt.Errorf("iso unknown format: %v", m.args.ISO)
			return
		}
	}

	// Datacenter
	if m.args.Datacenter != "" {
		c.Datacenter, _, err = m.client.Datacenter.Get(ctx, m.args.Datacenter)
		if err != nil {
			return
		}
		if c.Datacenter == nil {
			err = fmt.Errorf("datacenter '%s' not found", m.args.Datacenter)
			return
		}
	}

	// Location
	if m.args.Location != "" {
		c.Location, _, err = m.client.Location.Get(ctx, m.args.Location)
		if err != nil {
			return
		}
		if c.Location == nil {
			err = fmt.Errorf("location '%s' not found", m.args.Location)
			return
		}
	}

	if m.args.SSHKeys != nil {
		ids := util.GetIdentifiers(m.args.SSHKeys)
		for _, id := range ids {
			var sshKey *hcloud.SSHKey
			if sshKey, _, err = m.client.SSHKey.Get(ctx, id); err != nil {
				return
			}
			if sshKey == nil {
				err = fmt.Errorf("SSH Key %q not found", id)
				return
			}
			c.SSHKeys = append(c.SSHKeys, sshKey)
		}
	}

	return
}

func toServer(server *hcloud.Server) Server {
	s := Server{
		ID:         server.ID,
		Name:       server.Name,
		Status:     string(server.Status),
		ServerType: server.ServerType.Name,
		Datacenter: server.Datacenter.Name,
		Location:   server.Datacenter.Location.Name,
		PublicIPv4: server.PublicNet.IPv4.IP.String(),
		PublicIPv6: server.PublicNet.IPv6.Network.String(),
	}
	if server.ISO != nil {
		s.ISO = server.ISO.Name
	}
	return s
}

var flags = pflag.NewFlagSet("hcloud_server", pflag.ContinueOnError)

func init() {
	flags.BoolP("version", "v", false, "Print version and exit")
}

func main() {
	ansible.RunModule(&module{
		waitFn: util.WaitForAction,
	}, flags)
}
