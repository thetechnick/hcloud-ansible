package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/ansible"
	"github.com/thetechnick/hcloud-ansible/pkg/util"
	"github.com/thetechnick/hcloud-ansible/pkg/version"
)

var flags = pflag.NewFlagSet("hcloud_inventory", pflag.ContinueOnError)

func init() {
	flags.BoolP("version", "v", false, "Print version and exit")
	flags.Bool("list", false, "Print inventory")
	flags.String("host", "", "Return hostvars of a single server (unsupported)")
}

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}
	if v, _ := flags.GetBool("version"); v {
		version.PrintText()
	}

	if printHost, _ := flags.GetString("host"); printHost != "" {
		fmt.Fprintf(os.Stderr, "--host is unsupported\n")
		os.Exit(1)
	}

	client, err := util.BuildClient("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Hetzner Cloud client: %v\n", err)
		os.Exit(1)
	}

	servers, err := client.Server.All(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing servers: %v\n", err)
		os.Exit(1)
	}

	inventory := ansible.NewInventory()
	for _, server := range servers {
		inventory.AddHost(ansible.InventoryHost{
			Host: server.Name,
			Vars: varsForServer(server),
			Groups: []string{
				server.Datacenter.Name,
				server.Datacenter.Location.Name,
				server.ServerType.Name,
				fmt.Sprintf("status_%s", string(server.Status)),
				imageTag(server),
			},
		})
	}

	json, _ := json.MarshalIndent(inventory, "", "  ")
	fmt.Println(string(json))
}

func imageTag(server *hcloud.Server) string {
	if server.Image.Name != "" {
		return server.Image.Name
	}
	return fmt.Sprintf("image_%d", server.Image.ID)
}

func varsForServer(server *hcloud.Server) map[string]interface{} {
	vars := map[string]interface{}{
		"hcloud_id":          server.ID,
		"hcloud_name":        server.Name,
		"hcloud_public_ipv4": server.PublicNet.IPv4.IP.String(),
		"hcloud_public_ipv6": server.PublicNet.IPv6.IP.String(),
		"hcloud_location":    server.Datacenter.Location.Name,
		"hcloud_datacenter":  server.Datacenter.Name,
		"hcloud_status":      string(server.Status),
		"hcloud_server_type": server.ServerType.Name,

		"ansible_host": server.PublicNet.IPv4.IP.String(),
	}

	if server.Image.Name != "" {
		vars["hcloud_image"] = server.Image.Name
	} else {
		vars["hcloud_image"] = server.Image.ID
	}

	return vars
}
