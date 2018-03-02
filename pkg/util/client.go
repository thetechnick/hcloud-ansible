package util

import (
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// BuildClient creates and configures an hcloud client
func BuildClient(token string) (*hcloud.Client, error) {
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

	return hcloud.NewClient(opts...), nil
}
