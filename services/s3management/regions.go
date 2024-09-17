package s3management

import (
	"context"

	s3management "github.com/ionos-cloud/sdk-go-s3-management"
)

func (c *Client) GetRegion(ctx context.Context, regionId string, depth float32) (s3management.Region, *s3management.APIResponse, error) {
	region, apiResponse, err := c.client.RegionsApi.RegionsFindByRegion(ctx, regionId).Execute()
	apiResponse.LogInfo()
	return region, apiResponse, err
}

func (c *Client) ListRegions(ctx context.Context) (s3management.RegionList, *s3management.APIResponse, error) {
	regions, apiResponse, err := c.client.RegionsApi.RegionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return regions, apiResponse, err
}
