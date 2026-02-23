package containerregistry

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
)

// Client is a wrapper for the sdk client
type Client struct {
	sdkClient cr.APIClient
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: *cr.NewAPIClient(config),
	}
}
