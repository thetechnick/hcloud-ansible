package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
)

type arguments struct {
	Token string `json:"token"`
	State string `json:"state"`

	ID         interface{} `json:"id"`
	Name       interface{} `json:"name"`
	Image      string      `json:"image"`
	ServerType string      `json:"server_type"`
	UserData   string      `json:"user_data"`
	Datacenter string      `json:"datacenter"`
	Location   string      `json:"location"`
	Rescue     string      `json:"rescue"`
	SSHKeys    interface{} `json:"ssh_keys"`
}

const (
	stateAbsent    = "absent"
	statePresent   = "present"
	stateRunning   = "running"
	stateStopped   = "stopped"
	stateRestarted = "restarted"
)

type config struct {
	State      string
	Name       []string
	ID         []int
	Image      *hcloud.Image
	ServerType string
	UserData   string
	Datacenter string
	Location   string
	Rescue     string
	SSHKeys    []*hcloud.SSHKey
}

type out struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ServerType string `json:"server_type"`
	Status     string `json:"status"`
	Datacenter string `json:"datacenter"`
	Location   string `json:"location"`
	PublicIPv4 string `json:"public_ipv4"`
	PublicIPv6 string `json:"public_ipv6"`
}

type module struct {
	args            arguments
	client          *hcloud.Client
	config          config
	messages        ansible.MessageLog
	actionWatchLock sync.Mutex
}

func (m *module) Args() interface{} {
	return &m.args
}

func (m *module) Run() (resp ansible.ModuleResponse, err error) {
	ctx := context.Background()
	m.client, err = util.BuildClient(m.args.Token)
	if err != nil {
		return
	}
	m.config = config{
		State:      m.args.State,
		ServerType: m.args.ServerType,
		UserData:   m.args.UserData,
		Datacenter: m.args.Datacenter,
		Location:   m.args.Location,
		Rescue:     m.args.Rescue,
	}
	if m.config.State == "" {
		m.config.State = statePresent
	}
	if m.config.Name, err = names(m.args.Name); err != nil {
		return
	}
	if m.config.Image, err = m.image(ctx, m.args.Image); err != nil {
		return
	}
	if m.config.SSHKeys, err = m.sshKeys(ctx, m.args.SSHKeys); err != nil {
		return
	}

	if err = validateConfig(m.config); err != nil {
		return
	}

	var (
		wg        sync.WaitGroup
		stateLock sync.Mutex
		instances []*out
	)

	ensureServer := func(id int, name string) {
		out, changed, err := m.ensureServer(ctx, 0, name)
		stateLock.Lock()
		if err != nil {
			m.messages.Add(err.Error())
			resp.Failed()
			return
		}
		if changed {
			resp.Changed()
		}
		if out != nil {
			instances = append(instances, out)
		}
		stateLock.Unlock()
		wg.Done()
	}

	if len(m.config.ID) > 0 {
		for i, id := range m.config.ID {
			wg.Add(1)
			var name string
			if len(m.config.Name) > i {
				name = m.config.Name[i]
			}
			ensureServer(id, name)
		}
	} else {
		for _, name := range m.config.Name {
			wg.Add(1)
			ensureServer(0, name)
		}
	}

	wg.Wait()
	if len(instances) > 0 {
		resp.Set("servers", instances)
	}
	resp.Msg(m.messages.String())
	return
}

func (m *module) image(ctx context.Context, imageArg string) (image *hcloud.Image, err error) {
	if imageArg == "" {
		return
	}

	if image, _, err = m.client.Image.Get(ctx, imageArg); err != nil {
		return
	}
	if image == nil {
		err = fmt.Errorf("'image' (%s) not found", imageArg)
		return
	}
	return
}

func (m *module) sshKeys(ctx context.Context, sshArgs interface{}) (sshKeys []*hcloud.SSHKey, err error) {
	sshKeysList, ok := sshArgs.([]interface{})
	if !ok {
		// no ssh keys
		return
	}

	for _, sshKeyValue := range sshKeysList {
		var sshKey *hcloud.SSHKey
		switch v := sshKeyValue.(type) {
		case int:
			sshKey, _, err = m.client.SSHKey.GetByID(ctx, v)

		case string:
			sshKey, _, err = m.client.SSHKey.Get(ctx, v)

		case map[string]interface{}:
			var nameOrID string
			if name, ok := v["name"]; ok {
				name, tok := name.(string)
				if !tok {
					err = fmt.Errorf("'ssh_keys' invalid format: 'name' is no string")
					return
				}
				nameOrID = name
			}
			if id, ok := v["id"]; ok {
				id, tok := id.(float64)
				if !tok {
					err = fmt.Errorf("'ssh_keys' invalid format: 'id' is no number")
					return
				}
				nameOrID = strconv.Itoa(int(id))
			}
			sshKey, _, err = m.client.SSHKey.Get(ctx, nameOrID)

		default:
			err = fmt.Errorf("'ssh_keys' unkown format")
			return
		}

		if err != nil {
			return
		}
		if sshKey == nil {
			err = fmt.Errorf("SSH key not found: %v", sshKeyValue)
			return
		}
		sshKeys = append(sshKeys, sshKey)
	}
	return
}

// ids parses the id input property
func ids(ids interface{}) (out []int, err error) {
	switch v := ids.(type) {
	case int:
		out = append(out, v)
		return

	case []interface{}:
		for _, idI := range v {
			if id, ok := idI.(int); ok {
				out = append(out, id)
				continue
			}
			err = fmt.Errorf("'id' unkown format")
			return
		}
		return

	case nil:
		return

	default:
		err = fmt.Errorf("'id' unkown format")
		return
	}
}

// names parses the name input property
func names(names interface{}) (out []string, err error) {
	switch v := names.(type) {
	case string:
		out = append(out, v)
		return

	case []interface{}:
		for _, nameI := range v {
			if name, ok := nameI.(string); ok {
				out = append(out, name)
				continue
			}
			err = fmt.Errorf("'name' unkown format")
			return
		}
		return

	case nil:
		return

	default:
		err = fmt.Errorf("'name' unkown format")
		return
	}
}

func (m *module) ensureServer(ctx context.Context, id int, name string) (out *out, changed bool, err error) {
	var server *hcloud.Server
	if id != 0 {
		if server, _, err = m.client.Server.GetByID(ctx, id); err != nil {
			return
		}
		if server == nil {
			err = fmt.Errorf("Server %d not found", id)
			return
		}
	} else {
		if server, _, err = m.client.Server.GetByName(ctx, name); err != nil {
			return
		}
	}

	if m.args.State == stateAbsent {
		if server != nil {
			if _, err = m.client.Server.Delete(ctx, server); err != nil {
				return
			}

			m.messages.Add(fmt.Sprintf("Server %d deleted", server.ID))
			changed = true
			return
		}
		return
	}

	if server != nil && needsRecreate(server, m.config) {
		if _, err = m.client.Server.Delete(ctx, server); err != nil {
			return
		}

		m.messages.Add(fmt.Sprintf("Server %d deleted (recreate)", server.ID))
		server = nil
		changed = true
	}

	if server == nil {
		if server, err = m.createServer(ctx, name); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d created", server.ID))
		changed = true
	}

	switch m.config.State {
	case stateRunning:
		if server.Status != hcloud.ServerStatusRunning {
			var action *hcloud.Action
			if action, _, err = m.client.Server.Poweron(ctx, server); err != nil {
				return
			}

			if err = m.waitForAction(ctx, action); err != nil {
				return
			}
			m.messages.Add(fmt.Sprintf("Server %d started", server.ID))
			changed = true
		}

	case stateStopped:
		if server.Status != hcloud.ServerStatusOff {
			var action *hcloud.Action
			if action, _, err = m.client.Server.Poweroff(ctx, server); err != nil {
				return
			}

			if err = m.waitForAction(ctx, action); err != nil {
				return
			}
			m.messages.Add(fmt.Sprintf("Server %d stopped", server.ID))
			changed = true
		}

	case stateRestarted:
		var action *hcloud.Action
		if action, _, err = m.client.Server.Reboot(ctx, server); err != nil {
			return
		}

		if err = m.waitForAction(ctx, action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d restarted", server.ID))
		changed = true
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

		if err = m.waitForAction(ctx, action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d disabled rescue mode", server.ID))
		changed = true
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

		if err = m.waitForAction(ctx, res.Action); err != nil {
			return
		}
		m.messages.Add(fmt.Sprintf("Server %d enabled rescue mode", server.ID))
		changed = true
		rescueChanged = true
	}

	if rescueChanged && m.config.State != stateStopped {
		var action *hcloud.Action
		if action, _, err = m.client.Server.Reset(ctx, server); err != nil {
			return
		}
		if err = m.waitForAction(ctx, action); err != nil {
			return
		}
	}

	out = toServerData(server)
	return
}

func (m *module) createServer(ctx context.Context, name string) (*hcloud.Server, error) {
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

	opts := hcloud.ServerCreateOpts{
		Name: name,
		ServerType: &hcloud.ServerType{
			Name: m.config.ServerType,
		},
		UserData: m.config.UserData,
		SSHKeys:  m.config.SSHKeys,
		Image:    m.config.Image,
	}

	if m.config.Datacenter != "" {
		opts.Datacenter = &hcloud.Datacenter{Name: m.config.Datacenter}
	}
	if m.config.Location != "" {
		opts.Location = &hcloud.Location{Name: m.config.Location}
	}

	res, _, err := m.client.Server.Create(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err = m.waitForAction(ctx, res.Action); err != nil {
		return nil, err
	}

	server, _, err := m.client.Server.GetByID(ctx, res.Server.ID)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// waitForAction makes sure we only watch one action at a time
func (m *module) waitForAction(ctx context.Context, action *hcloud.Action) error {
	m.actionWatchLock.Lock()
	defer m.actionWatchLock.Unlock()
	_, errCh := m.client.Action.WatchProgress(ctx, action)
	return <-errCh
}

func validateConfig(config config) error {
	errs := []string{}
	if len(config.Name) == 0 && len(config.ID) == 0 {
		errs = append(errs, "'name' or 'id' is required")
	}
	if config.State != stateAbsent &&
		config.State != statePresent &&
		config.State != stateRunning &&
		config.State != stateStopped &&
		config.State != stateRestarted {
		errs = append(errs, "'state' must be present, absent, running, stopped or restarted")
	}
	if config.Datacenter != "" && config.Location != "" {
		errs = append(errs, "'datacenter' and 'location' are mutually exclusive")
	}
	if len(errs) > 0 {
		return fmt.Errorf("Invalid config: %s", strings.Join(errs, ", "))
	}
	return nil
}

// needsRecreate checks if the server needs to be recreated
func needsRecreate(server *hcloud.Server, config config) bool {
	if config.Image != nil &&
		server.Image.ID != config.Image.ID {
		return true
	}
	if config.ServerType != "" &&
		server.ServerType.Name != config.ServerType {
		return true
	}
	if config.Datacenter != "" &&
		server.Datacenter.Name != config.Datacenter {
		return true
	}
	if config.Location != "" &&
		server.Datacenter.Location.Name != config.Location {
		return true
	}

	return false
}

func toServerData(server *hcloud.Server) *out {
	return &out{
		ID:         server.ID,
		Name:       server.Name,
		Status:     string(server.Status),
		ServerType: server.ServerType.Name,
		Datacenter: server.Datacenter.Name,
		Location:   server.Datacenter.Location.Name,
		PublicIPv4: server.PublicNet.IPv4.IP.String(),
		PublicIPv6: server.PublicNet.IPv6.Network.String()}
}

var flags = pflag.NewFlagSet("hcloud_server", pflag.ContinueOnError)

func init() {
	flags.BoolP("version", "v", false, "Print version and exit")
}

func main() {
	ansible.RunModule(&module{}, flags)
}
