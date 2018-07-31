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
		Client:     c,
		Action:     &c.Action,
		Server:     &c.Server,
		SSHKey:     &c.SSHKey,
		Image:      &c.Image,
		FloatingIP: &c.FloatingIP,
		Location:   &c.Location,
		Datacenter: &c.Datacenter,
		ISO:        &c.ISO,
	}
}

// String alias of hcloud.String
func String(s string) *string {
	return hcloud.String(s)
}

// Bool alias of hcloud.Bool
func Bool(b bool) *bool {
	return hcloud.Bool(b)
}

// Action alias of hcloud.Action
type Action = hcloud.Action

// ActionStatus alias of hcloud.ActionStatus
type ActionStatus = hcloud.ActionStatus

// ActionStatus
const (
	ActionStatusError   = hcloud.ActionStatusError
	ActionStatusSuccess = hcloud.ActionStatusSuccess
	ActionStatusRunning = hcloud.ActionStatusRunning
)

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

// FloatingIP alias of hcloud.FloatingIP
type FloatingIP = hcloud.FloatingIP

// FloatingIPCreateOpts alias of hcloud.FloatingIPCreateOpts
type FloatingIPCreateOpts = hcloud.FloatingIPCreateOpts

// FloatingIPListOpts alias of hcloud.FloatingIPListOpts
type FloatingIPListOpts = hcloud.FloatingIPListOpts

// FloatingIPCreateResult alias of hcloud.FloatingIPCreateResult
type FloatingIPCreateResult = hcloud.FloatingIPCreateResult

// FloatingIPType alias of hcloud.FloatingIPType
type FloatingIPType = hcloud.FloatingIPType

// FloatingIPUpdateOpts alias of hcloud.FloatingIPUpdateOpts
type FloatingIPUpdateOpts = hcloud.FloatingIPUpdateOpts

// Floating IP types.
const (
	FloatingIPTypeIPv4 = hcloud.FloatingIPTypeIPv4
	FloatingIPTypeIPv6 = hcloud.FloatingIPTypeIPv6
)

// FloatingIPClient interface of hcloud.FloatingIPClient
type FloatingIPClient interface {
	GetByID(ctx context.Context, id int) (*FloatingIP, *Response, error)
	List(ctx context.Context, opts FloatingIPListOpts) ([]*FloatingIP, *Response, error)
	All(ctx context.Context) ([]*FloatingIP, error)
	Create(ctx context.Context, opts FloatingIPCreateOpts) (FloatingIPCreateResult, *Response, error)
	Delete(ctx context.Context, floatingIP *FloatingIP) (*Response, error)
	Assign(ctx context.Context, floatingIP *FloatingIP, server *Server) (*Action, *Response, error)
	Unassign(ctx context.Context, floatingIP *FloatingIP) (*Action, *Response, error)
	Update(ctx context.Context, floatingIP *FloatingIP, opts FloatingIPUpdateOpts) (*FloatingIP, *Response, error)
}

// Image alias of hcloud.Image
type Image = hcloud.Image

// ImageListOpts alias of hcloud.ImageListOpts
type ImageListOpts = hcloud.ImageListOpts

// ImageClient interface of hcloud.ImageClient
type ImageClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Image, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Image, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Image, *hcloud.Response, error)
	List(ctx context.Context, opts hcloud.ImageListOpts) ([]*hcloud.Image, *hcloud.Response, error)
	All(ctx context.Context) ([]*hcloud.Image, error)
}

// ISO alias of hcloud.ISO
type ISO = hcloud.ISO

// ISOClient interface of hcloud.ISOClient
type ISOClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.ISO, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.ISO, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.ISO, *hcloud.Response, error)
}

// Location alias of hcloud.Location
type Location = hcloud.Location

// LocationClient interface of hcloud.LocationClient
type LocationClient interface {
	Get(ctx context.Context, idOrName string) (*hcloud.Location, *hcloud.Response, error)
}

// PricingClient interface of hcloud.PricingClient
type PricingClient interface{}

// Server alias of hcloud.Server
type Server = hcloud.Server

// ServerStatus alias of hcloud.ServerStatus
type ServerStatus = hcloud.ServerStatus

// ServerStatus alias of hcloud.ServerStatus*
var (
	ServerStatusOff     = hcloud.ServerStatusOff
	ServerStatusRunning = hcloud.ServerStatusRunning
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

// ServerPublicNet alias of hcloud.ServerPublicNet
type ServerPublicNet = hcloud.ServerPublicNet

// ServerPublicNetIPv4 alias of hcloud.ServerPublicNetIPv4
type ServerPublicNetIPv4 = hcloud.ServerPublicNetIPv4

// ServerPublicNetIPv6 alias of hcloud.ServerPublicNetIPv6
type ServerPublicNetIPv6 = hcloud.ServerPublicNetIPv6

// ServerCreateResult alias of hcloud.ServerCreateResult
type ServerCreateResult = hcloud.ServerCreateResult

// ServerClient interface of hcloud.ServerClient
type ServerClient interface {
	GetByID(ctx context.Context, id int) (*Server, *Response, error)
	GetByName(ctx context.Context, name string) (*Server, *Response, error)
	Get(ctx context.Context, idOrName string) (*Server, *Response, error)
	// List(ctx context.Context, opts ServerListOpts) ([]*Server, *Response, error)
	All(ctx context.Context) ([]*Server, error)
	Create(ctx context.Context, opts ServerCreateOpts) (ServerCreateResult, *Response, error)
	Delete(ctx context.Context, server *Server) (*Response, error)
	Update(ctx context.Context, server *Server, opts ServerUpdateOpts) (*Server, *Response, error)
	Poweron(ctx context.Context, server *Server) (*Action, *Response, error)
	Reboot(ctx context.Context, server *Server) (*Action, *Response, error)
	Reset(ctx context.Context, server *Server) (*Action, *Response, error)
	// Shutdown(ctx context.Context, server *Server) (*Action, *Response, error)
	Poweroff(ctx context.Context, server *Server) (*Action, *Response, error)
	// ResetPassword(ctx context.Context, server *Server) (ServerResetPasswordResult, *Response, error)
	// CreateImage(ctx context.Context, server *Server, opts *ServerCreateImageOpts) (ServerCreateImageResult, *Response, error)
	AttachISO(ctx context.Context, server *Server, iso *ISO) (*Action, *Response, error)
	DetachISO(ctx context.Context, server *Server) (*Action, *Response, error)
	EnableRescue(ctx context.Context, server *Server, opts ServerEnableRescueOpts) (ServerEnableRescueResult, *Response, error)
	DisableRescue(ctx context.Context, server *Server) (*Action, *Response, error)
	// Rebuild(ctx context.Context, server *Server, opts ServerRebuildOpts) (*Action, *Response, error)
	// AttachISO(ctx context.Context, server *Server, iso *ISO) (*Action, *Response, error)
	// DetachISO(ctx context.Context, server *Server) (*Action, *Response, error)
	EnableBackup(ctx context.Context, server *Server, window string) (*Action, *Response, error)
	DisableBackup(ctx context.Context, server *Server) (*Action, *Response, error)
	// ChangeType(ctx context.Context, server *Server, opts ServerChangeTypeOpts) (*Action, *Response, error)
	// ChangeDNSPtr(ctx context.Context, server *Server, ip string, ptr *string) (*Action, *Response, error)
}

// ServerType alias of hcloud.ServerType
type ServerType = hcloud.ServerType

// ServerTypeListOpts alias of hcloud.ServerTypeListOpts
type ServerTypeListOpts = hcloud.ServerTypeListOpts

// ServerTypeClient interface of hcloud.ServerTypeClient
type ServerTypeClient interface {
	Get(ctx context.Context, idOrName string) (*ServerType, *Response, error)
}

// SSHKey alias of hcloud.SSHKey
type SSHKey = hcloud.SSHKey

// SSHKeyUpdateOpts alias of hcloud.SSHKeyUpdateOpts
type SSHKeyUpdateOpts = hcloud.SSHKeyUpdateOpts

// SSHKeyCreateOpts alias of hcloud.SSHKeyCreateOpts
type SSHKeyCreateOpts = hcloud.SSHKeyCreateOpts

// SSHKeyClient interface of hcloud.SSHKeyClient
type SSHKeyClient interface {
	GetByID(ctx context.Context, id int) (*SSHKey, *Response, error)
	GetByName(ctx context.Context, name string) (*SSHKey, *Response, error)
	Get(ctx context.Context, idOrName string) (*SSHKey, *Response, error)
	All(ctx context.Context) ([]*SSHKey, error)
	Create(ctx context.Context, opts SSHKeyCreateOpts) (*SSHKey, *Response, error)
	Delete(ctx context.Context, sshKey *SSHKey) (*Response, error)
	Update(ctx context.Context, sshKey *SSHKey, opts SSHKeyUpdateOpts) (*SSHKey, *Response, error)
}
