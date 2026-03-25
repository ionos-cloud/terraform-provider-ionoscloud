package pgsqlv2

import (
	"context"

	pgsqlv2 "github.com/ionos-cloud/pgsqlv2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetVersion retrieves a PostgreSQL version by its ID.
func (c *Client) GetVersion(ctx context.Context, versionID string) (pgsqlv2.PostgresVersionRead, *shared.APIResponse, error) {
	version, apiResponse, err := c.sdkClient.VersionsApi.VersionsFindById(ctx, versionID).Execute()
	apiResponse.LogInfo()
	return version, apiResponse, err
}

// ListVersions retrieves all available PostgreSQL versions.
func (c *Client) ListVersions(ctx context.Context) (pgsqlv2.PostgresVersionReadList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return versions, apiResponse, err
}
