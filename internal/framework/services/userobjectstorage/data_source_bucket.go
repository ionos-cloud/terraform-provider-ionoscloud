package userobjectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/userobjectstorage"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*bucketDataSource)(nil)

// NewBucketDataSource creates a new data source for the user object storage bucket.
func NewBucketDataSource() datasource.DataSource {
	return &bucketDataSource{}
}

type bucketDataSource struct {
	client *userobjectstorage.Client
}

type bucketDataSourceModel struct {
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

// Metadata returns the metadata for the data source.
func (d *bucketDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_object_storage_bucket"
}

// Schema returns the schema for the data source.
func (d *bucketDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the bucket.",
				Required:    true,
			},
			"region": schema.StringAttribute{
				Description: "The region of the bucket. Defaults to 'de' (Frankfurt). Valid values: 'de', 'eu-central-2', 'eu-south-2'.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Configure wires the provider client.
func (d *bucketDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = clientBundle.UserS3Client
}

// Read reads the data source.
func (d *bucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("user object storage client not configured", "The provider client is not configured")
		return
	}

	var data bucketDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Region.IsNull() || data.Region.IsUnknown() || data.Region.ValueString() == "" {
		data.Region = types.StringValue(userobjectstorage.DefaultRegion)
	}

	found, err := d.client.GetBucket(ctx, data.Name, data.Region)
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("bucket not found", fmt.Sprintf("bucket %q was not found in region %q", data.Name.ValueString(), data.Region.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
