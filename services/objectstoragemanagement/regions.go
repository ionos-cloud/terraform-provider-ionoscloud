package objectstoragemanagement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objectstoragemanagement "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// RegionDataSourceModel is used to represent an region for a data source.
type RegionDataSourceModel struct {
	Version        types.Int32    `tfsdk:"version"`
	Endpoint       types.String   `tfsdk:"endpoint"`
	Website        types.String   `tfsdk:"website"`
	Capability     *capability    `tfsdk:"capability"`
	Storageclasses []types.String `tfsdk:"storage_classes"`
	Location       types.String   `tfsdk:"location"`
	ID             types.String   `tfsdk:"id"`
}

type capability struct {
	Iam      types.Bool `tfsdk:"iam"`
	S3select types.Bool `tfsdk:"s3select"`
}

// GetRegion retrieves a region
func (c *Client) GetRegion(ctx context.Context, regionID string, depth float32) (objectstoragemanagement.RegionRead, *shared.APIResponse, error) {
	c.modifyConfigURL()
	region, apiResponse, err := c.client.RegionsApi.RegionsFindByRegion(ctx, regionID).Execute()
	apiResponse.LogInfo()
	return region, apiResponse, err
}

// ListRegions lists all regions
func (c *Client) ListRegions(ctx context.Context) (objectstoragemanagement.RegionReadList, *shared.APIResponse, error) {
	c.modifyConfigURL()
	regions, apiResponse, err := c.client.RegionsApi.RegionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return regions, apiResponse, err
}

// BuildRegionModelFromAPIResponse builds an RegionDataSourceModel from a region SDK object
func BuildRegionModelFromAPIResponse(output *objectstoragemanagement.RegionRead) *RegionDataSourceModel {
	built := &RegionDataSourceModel{}

	built.ID = types.StringValue(output.Id)
	built.Version = types.Int32Value(output.Properties.Version)
	built.Endpoint = types.StringValue(output.Properties.Endpoint)
	built.Website = types.StringValue(output.Properties.Website)

	built.Capability = &capability{
		Iam:      types.BoolPointerValue(output.Properties.Capability.Iam),
		S3select: types.BoolPointerValue(output.Properties.Capability.S3select),
	}

	if output.Properties.StorageClasses != nil {
		built.Storageclasses = make([]types.String, 0, len(output.Properties.StorageClasses))
		for i := range output.Properties.StorageClasses {
			built.Storageclasses = append(built.Storageclasses, types.StringPointerValue(&(output.Properties.StorageClasses)[i]))
		}
	}

	return built
}
