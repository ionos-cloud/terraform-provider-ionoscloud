package objectstoragemanagement

import (
	objectstoragemanagement "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// Client is a wrapper around the S3 client.
type Client struct {
	client *objectstoragemanagement.APIClient
	// Diags carries per-configuration error-enrichment context (e.g. contract
	// number). Populated by the SdkBundle constructor.
	Diags *diagutil.Enricher
}

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objectstoragemanagement.APIClient {
	return c.client
}

func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		client: objectstoragemanagement.NewAPIClient(config),
	}
}
