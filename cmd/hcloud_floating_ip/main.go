package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
)

const (
	stateAbsent  = "absent"
	statePresent = "present"
	stateList    = "list"
)

// FloatingIP is the module return value of an hcloud.FloatingIP
type FloatingIP struct {
	ID           int    `json:"id"`
	Description  string `json:"description"`
	IP           string `json:"ip"`
	Type         string `json:"type"`
	ServerID     *int   `json:"server_id"`
	HomeLocation string `json:"home_location"`
}

type arguments struct {
	Token string `json:"token"`
	State string `json:"state"`

	ID           interface{} `json:"id"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Server       interface{} `json:"server"`
	HomeLocation string      `json:"home_location"`
}

type module struct {
	args   arguments
	client *hcloud.Client
	waitFn util.WaitFn
}

func (m *module) Args() interface{} {
	return &m.args
}

func (m *module) Run() (resp ansible.ModuleResponse, err error) {
	m.client, err = hcloud.BuildClient(m.args.Token)
	if err != nil {
		return
	}
	if m.args.State == "" {
		m.args.State = statePresent
	}
	return m.run()
}

func (m *module) run() (resp ansible.ModuleResponse, err error) {
	ctx := context.Background()
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

func (m *module) present(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var floatingIP *hcloud.FloatingIP
	id := util.GetID(m.args.ID)
	if id != 0 {
		if floatingIP, _, err = m.client.FloatingIP.GetByID(ctx, id); err != nil {
			return
		}
	}

	var server *hcloud.Server
	if server, err = m.server(ctx, m.args.Server); err != nil {
		return
	}

	var msg []string
	if floatingIP == nil {
		opts := hcloud.FloatingIPCreateOpts{
			Description: hcloud.String(m.args.Description),
			Server:      server,
			Type:        hcloud.FloatingIPType(m.args.Type),
		}
		if m.args.HomeLocation != "" {
			opts.HomeLocation = &hcloud.Location{
				ID:   util.GetID(m.args.HomeLocation),
				Name: util.GetName(m.args.HomeLocation),
			}
		}

		var r hcloud.FloatingIPCreateResult
		r, _, err = m.client.FloatingIP.Create(ctx, opts)
		if err != nil {
			return
		}
		if r.Action != nil {
			err = util.WaitForAction(ctx, m.client, r.Action)
			if err != nil {
				return
			}
		}
		floatingIP = r.FloatingIP
		msg = append(msg, fmt.Sprintf("FloatingIP %d created", floatingIP.ID))
		resp.Changed()
	}

	if floatingIP.Description != m.args.Description {
		floatingIP, _, err = m.client.FloatingIP.Update(ctx, floatingIP, hcloud.FloatingIPUpdateOpts{
			Description: m.args.Description,
		})
		if err != nil {
			return
		}
		msg = append(msg, fmt.Sprintf("FloatingIP %d description changed", floatingIP.ID))
		resp.Changed()
	}

	if floatingIP.Server != nil && server == nil {
		var action *hcloud.Action
		action, _, err = m.client.FloatingIP.Unassign(ctx, floatingIP)
		if err != nil {
			return
		}

		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}

		msg = append(msg, fmt.Sprintf("FloatingIP %d unassigned", floatingIP.ID))
		resp.Changed()
	}

	if floatingIP.Server == nil && server != nil ||
		server != nil && floatingIP.Server.ID != server.ID {
		var action *hcloud.Action
		action, _, err = m.client.FloatingIP.Assign(ctx, floatingIP, server)
		if err != nil {
			return
		}

		if err = m.waitFn(ctx, m.client, action); err != nil {
			return
		}

		msg = append(msg, fmt.Sprintf("FloatingIP %d assigned to server %d", floatingIP.ID, server.ID))
		resp.Changed()
	}

	resp.
		Msg(strings.Join(msg, ", ")).
		Set("floating_ips", []FloatingIP{toFloatingIP(floatingIP)})

	return
}

func (m *module) absent(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	id := util.GetID(m.args.ID)
	var floatingIP *hcloud.FloatingIP
	if floatingIP, _, err = m.client.FloatingIP.GetByID(ctx, id); err != nil {
		return
	}
	if floatingIP == nil {
		resp.Msg("No FloatingIP found, nothing to do")
		return
	}
	if _, err = m.client.FloatingIP.Delete(ctx, floatingIP); err != nil {
		return
	}
	resp.Msg(fmt.Sprintf("FloatingIP %d deleted", floatingIP.ID)).Changed()
	return
}

func (m *module) list(ctx context.Context) (resp ansible.ModuleResponse, err error) {
	var (
		list        []FloatingIP
		floatingIPs []*hcloud.FloatingIP
	)
	floatingIPs, err = m.client.FloatingIP.All(ctx)
	if err != nil {
		return
	}

	for _, floatingIP := range floatingIPs {
		list = append(list, toFloatingIP(floatingIP))
	}
	resp.Msg("FloatingIPs listed").Set("floating_ips", list)
	return
}

func (m *module) server(ctx context.Context, serverArg interface{}) (server *hcloud.Server, err error) {
	id := util.GetID(serverArg)
	name := util.GetName(serverArg)

	if id == 0 && name == "" {
		return
	}

	if id != 0 {
		server, _, err = m.client.Server.GetByID(ctx, id)
	}

	if server == nil {
		if name != "" {
			server, _, err = m.client.Server.GetByName(ctx, name)
		}
	}

	if err != nil {
		return
	}
	if server == nil {
		err = fmt.Errorf("Server '%v' not found", serverArg)
	}
	return
}

func toFloatingIP(ip *hcloud.FloatingIP) FloatingIP {
	data := FloatingIP{
		ID:           ip.ID,
		Description:  ip.Description,
		IP:           ip.IP.String(),
		Type:         string(ip.Type),
		HomeLocation: ip.HomeLocation.Name,
	}
	if ip.Server != nil {
		data.ServerID = &ip.Server.ID
	}
	return data
}

func validateArgs(args arguments) error {
	errs := []string{}
	if args.State != stateAbsent &&
		args.State != statePresent &&
		args.State != stateList {
		errs = append(errs, "'state' must be present, absent or list")
	}
	if args.State == statePresent && args.ID == nil {
		if args.HomeLocation == "" && args.Server == nil {
			errs = append(errs, "'home_location' or 'server' must be set")
		}
		if args.HomeLocation != "" && args.Server != nil {
			errs = append(errs, "'home_location' and 'server' are mutually exclusive")
		}
	}
	if args.State == stateAbsent &&
		args.ID == nil {
		errs = append(errs, "'id' is required")
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}
	return nil
}

var flags = pflag.NewFlagSet("hcloud_floating_ip", pflag.ContinueOnError)

func init() {
	flags.BoolP("version", "v", false, "Print version and exit")
}

func main() {
	ansible.RunModule(&module{
		waitFn: util.WaitForAction,
	}, flags)
}
