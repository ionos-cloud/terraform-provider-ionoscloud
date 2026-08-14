package inmemorydbv2

import (
	"context"

	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetVersion retrieves a version by ID.
func (c *Client) GetVersion(ctx context.Context, versionID string) (inmemorydbv3.VersionRead, *shared.APIResponse, error) {
	version, apiResponse, err := c.sdkClient.VersionsApi.VersionsFindById(ctx, versionID).Execute()
	apiResponse.LogInfo()
	return version, apiResponse, err
}

// ListVersions retrieves the list of all supported versions.
func (c *Client) ListVersions(ctx context.Context) (inmemorydbv3.VersionReadList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return versions, apiResponse, err
}
