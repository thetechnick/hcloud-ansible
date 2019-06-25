export CGO_ENABLED:=0

VERSION=$(shell ./scripts/git-version)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE=$(shell date +%s)
REPO=github.com/thetechnick/hcloud-ansible
LD_FLAGS="-w -X $(REPO)/pkg/version.VersionTag=$(VERSION) -X $(REPO)/pkg/version.Branch=$(BRANCH) -X $(REPO)/pkg/version.BuildDate=$(BUILD_DATE)"

all: build

build: clean \
	bin/hcloud_ssh_key \
	bin/hcloud_server \
	bin/hcloud_floating_ip \
	bin/hcloud_inventory

bin/%:
	@go build -o bin/$* -v -ldflags $(LD_FLAGS) $(REPO)/cmd/$*

test:
	@./scripts/test

acceptance-test: build
	@rm -rf library
	@cp -a bin library
	ansible-playbook test.yml

clean:
	@rm -rf bin

clean-release:
	@rm -rf _output

release: \
	clean \
	clean-release \
	_output/hcloud-ansible_linux_amd64.zip \
	_output/hcloud-ansible_linux_386.zip \
	_output/hcloud-ansible_darwin_amd64.zip \
	_output/hcloud-ansible_freebsd_amd64.zip

bin/linux_amd64/hcloud_server:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/hcloud_server:  GOARGS = GOOS=linux GOARCH=386
bin/darwin_amd64/hcloud_server:  GOARGS = GOOS=darwin GOARCH=amd64
bin/freebsd_amd64/hcloud_server:  GOARGS = GOOS=freebsd GOARCH=amd64

bin/linux_amd64/hcloud_ssh_key:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/hcloud_ssh_key:  GOARGS = GOOS=linux GOARCH=386
bin/darwin_amd64/hcloud_ssh_key:  GOARGS = GOOS=darwin GOARCH=amd64
bin/freebsd_amd64/hcloud_ssh_key:  GOARGS = GOOS=freebsd GOARCH=amd64

bin/linux_amd64/hcloud_floating_ip:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/hcloud_floating_ip:  GOARGS = GOOS=linux GOARCH=386
bin/darwin_amd64/hcloud_floating_ip:  GOARGS = GOOS=darwin GOARCH=amd64
bin/freebsd_amd64/hcloud_floating_ip:  GOARGS = GOOS=freebsd GOARCH=amd64

bin/linux_amd64/hcloud_inventory:  GOARGS = GOOS=linux GOARCH=amd64
bin/linux_386/hcloud_inventory:  GOARGS = GOOS=linux GOARCH=386
bin/darwin_amd64/hcloud_inventory:  GOARGS = GOOS=darwin GOARCH=amd64
bin/freebsd_amd64/hcloud_inventory:  GOARGS = GOOS=freebsd GOARCH=amd64

bin/%/hcloud_server: clean
	$(GOARGS) go build -o $@ -ldflags $(LD_FLAGS) -a $(REPO)/cmd/hcloud_server

bin/%/hcloud_ssh_key: clean
	$(GOARGS) go build -o $@ -ldflags $(LD_FLAGS) -a $(REPO)/cmd/hcloud_ssh_key

bin/%/hcloud_floating_ip: clean
	$(GOARGS) go build -o $@ -ldflags $(LD_FLAGS) -a $(REPO)/cmd/hcloud_floating_ip

bin/%/hcloud_inventory: clean
	$(GOARGS) go build -o $@ -ldflags $(LD_FLAGS) -a $(REPO)/cmd/hcloud_inventory

_output/hcloud-ansible_%.zip: NAME=hcloud-ansible_$(VERSION)_$*
_output/hcloud-ansible_%.zip: DEST=_output/$(NAME)
_output/hcloud-ansible_%.zip: \
	bin/%/hcloud_server \
	bin/%/hcloud_ssh_key \
	bin/%/hcloud_floating_ip \
	bin/%/hcloud_inventory

	mkdir -p $(DEST)
	cp bin/$*/hcloud_floating_ip bin/$*/hcloud_server bin/$*/hcloud_ssh_key bin/$*/hcloud_inventory README.md LICENSE $(DEST)
	cd $(DEST) && zip -r ../$(NAME).zip .

.PHONY: all build clean test release acceptance-test
