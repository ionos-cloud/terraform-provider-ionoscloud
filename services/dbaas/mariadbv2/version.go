package mariadbv2

import (
	"context"

	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ListVersions retrieves a list of all supported MariaDB versions, following pagination.
func (c *Client) ListVersions(ctx context.Context) (mariadbv3.MariadbVersionReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mariadbv3.MariadbVersionRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mariadbv3.MariadbVersionReadList{}, apiResponse, err
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
