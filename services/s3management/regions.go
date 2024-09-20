package s3management

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s3management "github.com/ionos-cloud/sdk-go-s3-management"
)

// RegionDataSourceModel is used to represent an region for a data source.
type RegionDataSourceModel struct {
	Version        types.String `tfsdk:"verion"`
	Endpoint       types.String `tfsdk:"endpoint"`
	Website        types.String `tfsdk:"website"`
	Capability     types.String `tfsdk:"capability"`
	Storageclasses types.String `tfsdk:"storage_classes"`
	Location       types.String `tfsdk:"location"`
	Country        types.String `tfsdk:"country"`
	City           types.String `tfsdk:"city"`
	ID             types.String `tfsdk:"id"`
}

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
