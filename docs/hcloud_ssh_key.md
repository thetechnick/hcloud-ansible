# hcloud_ssh_key

Manages Hetzner Cloud SSH Keys. This module can be used to create, list and delete ssh keys.

## Requirements (on host that executes module)
- ansible >= 2.2.x (binary module support)

## Options
|parameter|required|default|choices|comments|
|---------|--------|-------|-------|--------|
|token|no|||Hetzner Cloud API Token. Can also be specified with `HCLOUD_TOKEN` environment variable. |
|state|no|present|<ul><li>present</li><li>absent</li><li>list</li></ul>|  `list` lists all existing ssh keys. |
| id | no | | | ID of the ssh key. (with state: `absent`) |
| name | no | | | Name of the ssh key. Required when state is `present`. |
| public_key | no | | | Required when state is `present`. |

## Return Values

These values can be used when registering the modules output.

```yaml
ssh_keys:
- id: 123
  name: mykey@machine
  fingerprint: a2:94:75:0d:cf:fd:2c:fc:77:81:0e:c6:7a:8d:a2:21
```

## Examples

```yaml
# create an ssh key
- hcloud_ssh_key:
    name: test key
    public_key: "{{lookup('file', '~/.ssh/id_rsa.pub')}}"


# list all ssh keys in the Hetzner Cloud Project and
# create a single server with the fetched ssh keys
- hcloud_ssh_key:
    state: list
  register: hcloud_ssh_keys

- hcloud_server:
    name: example-server
    image: debian-9
    server_type: cx11
    datacenter: nbg1-dc3
    ssh_keys: "{{ hcloud_ssh_keys.ssh_keys }}"
```
