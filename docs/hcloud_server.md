# hcloud_server

Manages Hetzner Cloud servers. This module can be used to create, modify, delete and reboot servers.

## Requirements (on host that executes module)
- ansible >= 2.2.x (binary module support)

## Options
|parameter|required|default|choices|comments|
|---------|--------|-------|-------|--------|
|token|no|||Hetzner Cloud API Token. Can also be specified with `HCLOUD_TOKEN` environment variable. |
|state|no|present|<ul><li>present</li><li>absent</li><li>running</li><li>stopped</li><li>restarted</li><li>list</li></ul>|  |
| id | no | | | A single id or list of ids. Either `id` or `name` must be set. |
| name | no | | | A single name or list of names. Either `id` or `name` must be set. |
| image | no | | | Required when a server needs to be created. |
| server_type | no | | | Required when a server needs to be created. |
| user_data | no | | | cloud-init userdata |
| datacenter | no | | | Mutually exclusive with `location` |
| location | no | | | Mutually exclusive with `datacenter` |
| rescue | no | | <ul><li>linux64</li><li>linux32</li><li>freebsd64</li></ul> | Will make sure the choosen rescue system is enabled. Automatically resets the server to boot into the rescue system if `state != stopped`. |
| ssh_keys | no | | | List of Hetzner Cloud SSHKey ids, names or dict containing the `id` or `name`. |

## Return Values

These values can be used when registering the modules output.

```yaml
servers:
- id: 123
  name: server-name
  image: debian-9
  server_type: cx11
  status: running
  datacenter: fsn1-dc8
  location: fsn1
  public_ipv4: 10.0.0.1
  public_ipv6: 2001:db8::/64
```

## Examples

```yaml

# create a single server
- hcloud_server:
    name: example-server
    image: debian-9
    server_type: cx11
    datacenter: nbg1-dc3
    ssh_keys:
    - user@example-notebook   # by name
    - 1234                    # by id

# create tree servers at once
- hcloud_server:
    name:
    - example-server01
    - example-server02
    - example-server03
    image: debian-9
    server_type: cx11
    location: nbg1
    ssh_keys:
    - user@example-notebook   # by name
    - 1234                    # by id

# ensure server is running (if the server already exists)
- hcloud_server:
    name: example-server
    state: running

- hcloud_server:
    id: 123
    state: running

# ensure the servers name is "web-234" (if the server already exists)
- hcloud_server:
    id: 123
    name: web-234

# enable and boot into rescue system (if the server already exists)
- hcloud_server:
    name: example-server01
    rescue: linux64
    ssh_keys:                 # rescue os will be configured with these ssh keys
    - user@example-notebook   # by name
    - 1234                    # by id
```
