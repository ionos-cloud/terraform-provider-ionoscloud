package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var _ datasource.DataSourceWithConfigure = (*bucketPolicyDataSource)(nil)

// NewBucketPolicyDataSource creates a new data source for the bucket resource.
func NewBucketPolicyDataSource() datasource.DataSource {
	return &bucketPolicyDataSource{}
}

type bucketPolicyDataSource struct {
	client *s3.APIClient
}

type bucketPolicyDataSourceModel struct {
	Bucket types.String         `tfsdk:"bucket"`
	Policy jsontypes.Normalized `tfsdk:"policy"`
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
	var diags diag.Diagnostics

	if d.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPolicyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := GetBucketPolicy(ctx, d.client, data.Bucket.ValueString())
	if err != nil {
		if errors.Is(err, ErrBucketPolicyNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to retrieve bucket policy", err.Error())
		return
	}
	var policyData bucketPolicyModel
	if diags = setBucketPolicyData(policy, &policyData); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.Policy = policyData.Policy

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
