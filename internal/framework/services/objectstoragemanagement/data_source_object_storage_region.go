package objectstoragemanagement

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = (*regionDataSource)(nil)

// NewRegionDataSource creates a new data source for the region resource.
func NewRegionDataSource() datasource.DataSource {
	return &regionDataSource{}
}

type regionDataSource struct {
	client *objectStorageManagementService.Client
}

// Metadata returns the metadata for the data source.
func (d *regionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_region"
}

// Configure configures the data source.
func (d *regionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clientbundle, ok := req.ProviderData.(*services.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clientbundle.ObjectStorageManagementClient
}

// Schema returns the schema for the data source.
func (d *regionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the region",
				Required:    true,
			},
			"version": schema.Int32Attribute{
				Description: "The version of the region properties",
				Computed:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "The endpoint URL for the region",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "The website URL for the region",
				Computed:    true,
			},
			"storage_classes": schema.ListAttribute{
				Description: "The available classes in the region",
				Computed:    true,
				ElementType: types.StringType,
			},
			"location": schema.StringAttribute{
				Description: "The data center location of the region as per [Get Location](/docs/cloud/v6/#tag/Locations/operation/locationsGet). *Can't be used as `LocationConstraint` on bucket creation.*",
				Computed:    true,
			},
			"capability": schema.SingleNestedAttribute{
				Description: "The capabilities of the region",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"iam": schema.BoolAttribute{
						Description: "Indicates if IAM policy based access is supported",
						Computed:    true,
					},
					"s3select": schema.BoolAttribute{
						Description: "Indicates if S3 Select is supported",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Read reads the data source.
func (d *regionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("api client not configured", "The provider client is not configured")
		return
	}

	var data *objectStorageManagementService.RegionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region, apiResponse, err := d.client.GetRegion(ctx, data.ID.ValueString(), 1)

	if apiResponse.HttpNotFound() {
		resp.Diagnostics.AddError("region not found", "The region was not found")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to get region", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, objectStorageManagementService.BuildRegionModelFromAPIResponse(&region))...)
}
