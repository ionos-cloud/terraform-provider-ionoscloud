package userobjectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
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
	Name              types.String `tfsdk:"name"`
	ObjectLockEnabled types.Bool   `tfsdk:"object_lock_enabled"`
	Region            types.String `tfsdk:"region"`
	Tags              types.Map    `tfsdk:"tags"`
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
			"object_lock_enabled": schema.BoolAttribute{
				Description: "Whether Object Lock is enabled for the bucket.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags assigned to the bucket.",
				Computed:    true,
				ElementType: types.StringType,
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

	if data.Region.ValueString() == "" {
		data.Region = types.StringValue(userobjectstorage.DefaultRegion)
	}

	found, err := d.client.GetBucket(ctx, data.Name.ValueString(), data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("bucket not found", fmt.Sprintf("bucket %q was not found in region %q", data.Name.ValueString(), data.Region.ValueString()))
		return
	}

	objectLockEnabled, err := d.client.GetObjectLockEnabled(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get object lock configuration", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	rawTags, err := d.client.GetBucketTags(ctx, data.Name.ValueString(), data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	tagsMap, tagsErr := tags.KeyValueTags(rawTags).ToMap(ctx)
	if tagsErr != nil {
		resp.Diagnostics.AddError("failed to convert bucket tags", diagutil.WrapError(tagsErr, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	data.ObjectLockEnabled = types.BoolValue(objectLockEnabled)
	data.Tags = tagsMap
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
