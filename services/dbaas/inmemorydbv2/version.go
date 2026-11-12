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

// ListVersions retrieves the list of all supported versions, following pagination.
func (c *Client) ListVersions(ctx context.Context) (inmemorydbv3.VersionReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []inmemorydbv3.VersionRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return inmemorydbv3.VersionReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		// A short (or empty) page marks the end of the result set. Links.HasNext() is
		// deliberately not used: it has been observed to stay populated past the last page.
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
