package objectstorage

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = (*bucketPolicyDataSource)(nil)

// NewBucketPolicyDataSource creates a new data source for the bucket resource.
func NewBucketPolicyDataSource() datasource.DataSource {
	return &bucketPolicyDataSource{}
}

type bucketPolicyDataSource struct {
	client *objectstorage.Client
}

// Metadata returns the metadata for the data source.
func (d *bucketPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_policy"
}

// Configure configures the data source.
func (d *bucketPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*objectstorage.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *objectstorage.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Schema returns the schema for the data source.
func (d *bucketPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
			},
			"policy": schema.StringAttribute{
				CustomType:  jsontypes.NormalizedType{},
				Description: "Text of the policy",
				Computed:    true,
			},
		},
	}
}

// Read reads the data source.
func (d *bucketPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.BucketPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := d.client.GetBucketPolicy(ctx, data.Bucket)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	if !found {
		resp.Diagnostics.AddError("Failed to read resource", "Resource not found")
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
