package s3management

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s3management "github.com/ionos-cloud/sdk-go-s3-management"
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
func (c *Client) GetRegion(ctx context.Context, regionID string, depth float32) (s3management.RegionRead, *s3management.APIResponse, error) {
	region, apiResponse, err := c.client.RegionsApi.RegionsFindByRegion(ctx, regionID).Execute()
	apiResponse.LogInfo()
	return region, apiResponse, err
}

// ListRegions lists all regions
func (c *Client) ListRegions(ctx context.Context) (s3management.RegionReadList, *s3management.APIResponse, error) {
	regions, apiResponse, err := c.client.RegionsApi.RegionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return regions, apiResponse, err
}

// BuildRegionModelFromAPIResponse builds an RegionDataSourceModel from a region SDK object
func BuildRegionModelFromAPIResponse(output *s3management.RegionRead) *RegionDataSourceModel {
	built := &RegionDataSourceModel{}

	if output.Id != nil {
		built.ID = types.StringPointerValue(output.Id)
	}
	if output.Properties != nil {
		if output.Properties.Version != nil {
			built.Version = types.Int32PointerValue(output.Properties.Version)
		}
		if output.Properties.Endpoint != nil {
			built.Endpoint = types.StringPointerValue(output.Properties.Endpoint)
		}
		if output.Properties.Website != nil {
			built.Website = types.StringPointerValue(output.Properties.Website)
		}

		if output.Properties.Capability != nil {
			built.Capability = &capability{
				Iam:      types.BoolPointerValue(output.Properties.Capability.Iam),
				S3select: types.BoolPointerValue(output.Properties.Capability.S3select),
			}
		}

		if output.Properties.StorageClasses != nil {
			built.Storageclasses = make([]types.String, 0, len(*output.Properties.StorageClasses))
			for i := range *output.Properties.StorageClasses {
				built.Storageclasses = append(built.Storageclasses, types.StringPointerValue(&(*output.Properties.StorageClasses)[i]))
			}
		}
	}

	return built
}
