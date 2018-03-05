package hcloud

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// BuildClient creates and configures an hcloud client
func BuildClient(token string) (*Client, error) {
	if token == "" {
		token = os.Getenv("HCLOUD_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("argument `token` or environment variable `HCLOUD_TOKEN` is required")
	}
	opts := []hcloud.ClientOption{
		hcloud.WithToken(token),
	}

	if endpoint := os.Getenv("HCLOUD_ENDPOINT"); endpoint != "" {
		opts = append(opts, hcloud.WithEndpoint(endpoint))
	}

	return NewClient(opts...), nil
}

// ClientOption alias of hcloud.ClientOption
type ClientOption = hcloud.ClientOption

// WithToken alias of hcloud.WithToken
func WithToken(token string) ClientOption {
	return hcloud.WithToken(token)
}

// WithEndpoint alias of hcloud.WithEndpoint
func WithEndpoint(token string) ClientOption {
	return hcloud.WithToken(token)
}

// Response alias of hcloud.Response
type Response = hcloud.Response

// Client is an alias using interfaces of hcloud.Client
type Client struct {
	*hcloud.Client
	Action     ActionClient
	Datacenter DatacenterClient
	FloatingIP FloatingIPClient
	Image      ImageClient
	ISO        ISOClient
	Location   LocationClient
	Pricing    PricingClient
	Server     ServerClient
	ServerType ServerTypeClient
	SSHKey     SSHKeyClient
}

// NewClient is creates a new wrapped client
func NewClient(options ...hcloud.ClientOption) *Client {
	c := hcloud.NewClient(options...)
	return &Client{
		Client: c,
		Action: &c.Action,
		Server: &c.Server,
		SSHKey: &c.SSHKey,
		Image:  &c.Image,
	}
}

// Action alias of hcloud.Action
type Action = hcloud.Action

// ActionClient interface of hcloud.ActionClient
type ActionClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Action, *hcloud.Response, error)
	List(ctx context.Context, opts hcloud.ActionListOpts) ([]*hcloud.Action, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.Action, error)
	WatchProgress(ctx context.Context, action *hcloud.Action) (<-chan int, <-chan error)
}

// Datacenter alias of hcloud.Datacenter
type Datacenter = hcloud.Datacenter

// DatacenterClient interface of hcloud.DatacenterClient
type DatacenterClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Datacenter, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Datacenter, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Datacenter, *hcloud.Response, error)
	List(ctx context.Context, opts hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.Datacenter, error)
}

// FloatingIPClient interface of hcloud.FloatingIPClient
type FloatingIPClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.FloatingIP, *hcloud.Response, error)
	List(ctx context.Context, opts hcloud.FloatingIPListOpts) ([]*hcloud.FloatingIP, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.FloatingIP, error)
	Create(ctx context.Context, opts hcloud.FloatingIPCreateOpts) (hcloud.FloatingIPCreateResult, *hcloud.Response, error)
	Delete(ctx context.Context, floatingIP *hcloud.FloatingIP) (*hcloud.Response, error)
}

// Image alias of hcloud.Image
type Image = hcloud.Image

// ImageClient interface of hcloud.ImageClient
type ImageClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Image, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Image, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Image, *hcloud.Response, error)
	List(ctx context.Context, opts hcloud.ImageListOpts) ([]*hcloud.Image, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.Image, error)
}

// ISOClient interface of hcloud.ISOClient
type ISOClient interface{}

// Location alias of hcloud.Location
type Location = hcloud.Location

// LocationClient interface of hcloud.LocationClient
type LocationClient interface{}

// PricingClient interface of hcloud.PricingClient
type PricingClient interface{}

// Server alias of hcloud.Server
type Server = hcloud.Server

// ServerStatus alias of hcloud.ServerStatus
type ServerStatus = hcloud.ServerStatus

// ServerStatus alias of hcloud.ServerStatus*
var (
	ServerStatusOff     ServerStatus = hcloud.ServerStatusOff
	ServerStatusRunning ServerStatus = hcloud.ServerStatusRunning
)

// ServerUpdateOpts alias of hcloud.ServerUpdateOpts
type ServerUpdateOpts = hcloud.ServerUpdateOpts

// ServerEnableRescueResult alias of hcloud.ServerEnableRescueResult
type ServerEnableRescueResult = hcloud.ServerEnableRescueResult

// ServerEnableRescueOpts alias of hcloud.ServerEnableRescueOpts
type ServerEnableRescueOpts = hcloud.ServerEnableRescueOpts

// ServerCreateOpts alias of hcloud.ServerCreateOpts
type ServerCreateOpts = hcloud.ServerCreateOpts

// ServerRescueType alias of hcloud.ServerRescueType
type ServerRescueType = hcloud.ServerRescueType

// ServerClient interface of hcloud.ServerClient
type ServerClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Server, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Server, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Server, *hcloud.Response, error)
	// List(ctx context.Context, opts hcloud.ServerListOpts) ([]*hcloud.Server, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.Server, error)
	Create(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error)
	Delete(ctx context.Context, server *hcloud.Server) (*hcloud.Response, error)
	Update(ctx context.Context, server *hcloud.Server, opts hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error)
	Poweron(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Reboot(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Reset(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	// Shutdown(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Poweroff(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	// ResetPassword(ctx context.Context, server *hcloud.Server) (hcloud.ServerResetPasswordResult, *hcloud.Response, error)
	// CreateImage(ctx context.Context, server *hcloud.Server, opts *hcloud.ServerCreateImageOpts) (hcloud.ServerCreateImageResult, *hcloud.Response, error)
	EnableRescue(ctx context.Context, server *hcloud.Server, opts hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error)
	DisableRescue(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	// Rebuild(ctx context.Context, server *hcloud.Server, opts hcloud.ServerRebuildOpts) (*hcloud.Action, *hcloud.Response, error)
	// AttachISO(ctx context.Context, server *hcloud.Server, iso *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error)
	// DetachISO(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	EnableBackup(ctx context.Context, server *hcloud.Server, window string) (*hcloud.Action, *hcloud.Response, error)
	DisableBackup(ctx context.Context, server *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	// ChangeType(ctx context.Context, server *hcloud.Server, opts hcloud.ServerChangeTypeOpts) (*hcloud.Action, *hcloud.Response, error)
	// ChangeDNSPtr(ctx context.Context, server *hcloud.Server, ip string, ptr *string) (*hcloud.Action, *hcloud.Response, error)
}

// ServerType alias of hcloud.ServerType
type ServerType = hcloud.ServerType

// ServerTypeClient interface of hcloud.ServerTypeClient
type ServerTypeClient interface{}

// SSHKey alias of hcloud.SSHKey
type SSHKey = hcloud.SSHKey

// SSHKeyUpdateOpts alias of hcloud.SSHKeyUpdateOpts
type SSHKeyUpdateOpts = hcloud.SSHKeyUpdateOpts

// SSHKeyCreateOpts alias of hcloud.SSHKeyCreateOpts
type SSHKeyCreateOpts = hcloud.SSHKeyCreateOpts

// SSHKeyClient interface of hcloud.SSHKeyClient
type SSHKeyClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.SSHKey, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.SSHKey, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.SSHKey, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.SSHKey, error)
	Create(ctx context.Context, opts hcloud.SSHKeyCreateOpts) (*hcloud.SSHKey, *hcloud.Response, error)
	Delete(ctx context.Context, sshKey *hcloud.SSHKey) (*hcloud.Response, error)
	Update(ctx context.Context, sshKey *hcloud.SSHKey, opts hcloud.SSHKeyUpdateOpts) (*hcloud.SSHKey, *hcloud.Response, error)
}
