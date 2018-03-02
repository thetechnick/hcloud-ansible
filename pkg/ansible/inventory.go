package ansible

import "encoding/json"

// Inventory represents a Ansible Inventory
type Inventory struct {
	hosts  map[string]InventoryHost
	groups map[string]InventoryGroup
}

// InventoryHost represents a Host in an Ansible Inventory
type InventoryHost struct {
	Host   string
	Groups []string
	Vars   map[string]interface{}
}

// InventoryGroup represents a Group in an Ansible Inventory
type InventoryGroup struct {
	Name     string
	Vars     map[string]interface{}
	Children []string
}

// NewInventory creates a new empty Inventory
func NewInventory() Inventory {
	return Inventory{
		hosts:  map[string]InventoryHost{},
		groups: map[string]InventoryGroup{},
	}
}

// AddHost adds the host to the inventory
func (i Inventory) AddHost(host InventoryHost) {
	for _, group := range host.Groups {
		if _, ok := i.groups[group]; !ok {
			i.AddGroup(InventoryGroup{Name: group})
		}
	}
	i.hosts[host.Host] = host
}

// AddGroup adds a group to the inventory
func (i Inventory) AddGroup(group InventoryGroup) {
	i.groups[group.Name] = group
}

func (i Inventory) buildMetaHostvars() (vars map[string]interface{}) {
	vars = map[string]interface{}{}
	for _, host := range i.hosts {
		vars[host.Host] = host.Vars
	}
	return
}

func (i Inventory) buildMeta() (meta map[string]interface{}) {
	meta = map[string]interface{}{
		"hostvars": i.buildMetaHostvars(),
	}
	return
}

func (i Inventory) hostsInGroup(group string) (hosts []string) {
	for _, host := range i.hosts {
		for _, g := range host.Groups {
			if g == group {
				hosts = append(hosts, host.Host)
			}
		}
	}
	return
}

// MarshalJSON converts this object into the ansible format
func (i Inventory) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"_meta": i.buildMeta(),
	}

	for _, group := range i.groups {
		hosts := i.hostsInGroup(group.Name)
		if len(hosts) == 0 {
			continue
		}
		data[group.Name] = inventoryGroup{
			Hosts:    hosts,
			Vars:     group.Vars,
			Children: group.Children,
		}
	}
	return json.Marshal(data)
}

type inventoryGroup struct {
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Children []string               `json:"children,omitempty"`
}
