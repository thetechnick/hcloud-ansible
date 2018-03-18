# hcloud_floating_ip

Manages Hetzner Cloud floating ips. This module can be used to create, modify, assign and delete floating ips.

## Requirements (on host that executes module)
- ansible >= 2.2.x (binary module support)

## Options
|parameter|required|default|choices|comments|
|---------|--------|-------|-------|--------|
|token|no|||Hetzner Cloud API Token. Can also be specified with `HCLOUD_TOKEN` environment variable. |
|state|no|present|<ul><li>present</li><li>absent</li><li>list</li></ul>|  `list` lists all existing floating ips.<br>**NOTICE:**<br> `present` is not idempotent and will create a new floating ip when `id` is not specified. |
| id | no | | | ID of the floating ip.<br>Required when `state=absent`. |
| description | no | | | Description of the floating ip. |
| type | no | ipv4 |<ul><li>ipv4</li><li>ipv6</li></ul>| Required when `state=present` and `id` is not specified. |
| home_location | no | | | Home location of the floating ip.<br> Required when `state=present` and `server` is not specified.<br> Mutually exclusive with `server`. |
| server | no | | | Server to assign the floating ip to.<br> Required when `state=present` and `home_location` is not specified.<br> Mutually exclusive with `home_location`. |

## Return Values

These values can be used when registering the modules output.

```yaml
floating_ips:
- id: 123
  type: ipv4
  ip: 131.232.99.1
  description: Loadbalancer IP
  home_location: fsn1
  server_id: 123
```

## Examples

```yaml
# create a floating ip and assign it to a server
- hcloud_floating_ip:
    description: Loadbalancer IP
    type: ipv4
    server: 123 # by id

# list all floating ips in the Hetzner Cloud Project and
# assign all of them to a server
- hcloud_floating_ip:
    state: list
  register: hcloud_floating_ips

- hcloud_floating_ip:
    id: "{{item.id}}"
    server: 123
  with_items: {{hcloud_floating_ips.floating_ips}}

# assign floating ip 123 to server "loadbalancer"
- hcloud_floating_ip:
    id: 123
    server: "loadbalancer" # by name
```
