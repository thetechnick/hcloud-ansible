# Hetzner Cloud - Ansible
[![GitHub release](https://img.shields.io/github/release/thetechnick/hcloud-ansible.svg)](https://github.com/thetechnick/hcloud-ansible/releases/latest) [![Build Status](https://travis-ci.org/thetechnick/hcloud-ansible.svg?branch=master)](https://travis-ci.org/thetechnick/hcloud-ansible)

The Hetzner Cloud (hcloud) ansible modules and inventory are used to interact with the resources supported by Hetzner Cloud. The modules and inventory need to be configured with an Hetzner Cloud API token before they can be used.

## Inventory Example Usage

```sh
# ping all hosts in Hetzner Cloud Project
ansible -i hcloud_inventory all -m ping

# ping all hosts with cx11 server type
ansible -i hcloud_inventory cx11 -m ping
```

## Modules

- [hcloud_server - Manage Hetzner Cloud Servers](./docs/hcloud_server.md)
- [hcloud_ssh_key - Manage Hetzner Cloud SSH Keys](./docs/hcloud_ssh_key.md)
- [hcloud_floating_ip - Manage Hetzner Cloud Floating IPs](./docs/hcloud_floating_ip.md)

## Installation

Download the binaries for your OS from the releases page and place them into your [Ansible library](http://docs.ansible.com/ansible/latest/intro_configuration.html#library).

## Licence

MIT license
