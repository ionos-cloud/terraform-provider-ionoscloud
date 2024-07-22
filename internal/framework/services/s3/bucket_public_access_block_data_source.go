package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var _ datasource.DataSourceWithConfigure = (*bucketPublicAccessBlockDataSource)(nil)

// NewBucketPublicAccessBlockDataSource creates a new data source for the bucket public access block resource.
func NewBucketPublicAccessBlockDataSource() datasource.DataSource {
	return &bucketPublicAccessBlockDataSource{}
}

type bucketPublicAccessBlockDataSource struct {
	client *s3.APIClient
}

type bucketPublicAccessBlockDataSourceModel struct {
	Bucket                types.String `tfsdk:"bucket"`
	BlockPublicACLS       types.Bool   `tfsdk:"block_public_acls"`
	BlockPublicPolicy     types.Bool   `tfsdk:"block_public_policy"`
	IgnorePublicACLS      types.Bool   `tfsdk:"ignore_public_acls"`
	RestrictPublicBuckets types.Bool   `tfsdk:"restrict_public_buckets"`
}

// Metadata returns the metadata for the data source.
func (d *bucketPublicAccessBlockDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_access_block"
}

// Configure configures the data source.
func (d *bucketPublicAccessBlockDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Schema returns the schema for the data source.
func (d *bucketPublicAccessBlockDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthBetween(3, 63)},
			},
			"block_public_acls": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"block_public_policy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"ignore_public_acls": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"restrict_public_buckets": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Read reads the data source.
func (d *bucketPublicAccessBlockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPublicAccessBlockDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := GetBucketPublicAccessBlock(ctx, d.client, data.Bucket.ValueString())
	if err != nil {
		if errors.Is(err, ErrBucketPublicAccessBlockNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to retrieve bucket public access block", err.Error())
		return
	}

	data.IgnorePublicACLS = types.BoolPointerValue(response.IgnorePublicAcls)
	data.BlockPublicACLS = types.BoolPointerValue(response.BlockPublicAcls)
	data.BlockPublicPolicy = types.BoolPointerValue(response.BlockPublicPolicy)
	data.RestrictPublicBuckets = types.BoolPointerValue(response.RestrictPublicBuckets)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
