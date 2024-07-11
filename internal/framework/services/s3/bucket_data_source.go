package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var _ datasource.DataSourceWithConfigure = (*bucketDataSource)(nil)

// NewBucketDataSource creates a new data source for the bucket resource.
func NewBucketDataSource() datasource.DataSource {
	return &bucketDataSource{}
}

type bucketDataSource struct {
	client *s3.APIClient
}

type bucketDataSourceModel struct {
	Bucket types.String `tfsdk:"bucket"`
	Region types.String `tfsdk:"region"`
}

// Metadata returns the metadata for the data source.
func (d *bucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Configure configures the data source.
func (d *bucketDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *bucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
			},
			"region": schema.StringAttribute{
				Description: "The location or region of the bucket",
				Computed:    true,
			},
		},
	}
}

// Read reads the data source.
func (d *bucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := d.client.BucketsApi.HeadBucket(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.Diagnostics.AddError("Bucket not found", "The specified bucket does not exist")
			return
		}

		resp.Diagnostics.AddError("Failed to read bucket", err.Error())
		return
	}

	location, _, err := d.client.BucketsApi.GetBucketLocation(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read bucket location", err.Error())
		return
	}

	data.Region = types.StringValue(location.GetLocationConstraint())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
