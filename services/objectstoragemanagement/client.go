package objectstoragemanagement

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	objectstoragemanagement "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"
)

// Client is a wrapper around the S3 client.
type Client struct {
	client *objectstoragemanagement.APIClient
}

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objectstoragemanagement.APIClient {
	return c.client
}

// NewClientFromConfig creates a new Object Storage Management client from a pre-built configuration.
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{client: objectstoragemanagement.NewAPIClient(config)}
}
