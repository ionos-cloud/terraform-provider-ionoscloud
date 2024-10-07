package s3

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithImportState = (*bucketCorsConfiguration)(nil)
	_ resource.ResourceWithConfigure   = (*bucketCorsConfiguration)(nil)
)

type bucketCorsConfiguration struct {
	client *s3.Client
}

// NewBucketCorsConfigurationResource creates a new resource for the bucket CORS configuration resource.
func NewBucketCorsConfigurationResource() resource.Resource {
	return &bucketCorsConfiguration{}
}

// Metadata returns the metadata for the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_cors_configuration"
}

// Schema returns the schema for the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"cors_rule": schema.ListNestedBlock{
				Description: "A configuration for Cross-Origin Resource Sharing (CORS).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Optional:    true,
							Description: "Container for the Contract Number of the owner.",
						},
						"allowed_headers": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Specifies which headers are allowed in a preflight OPTIONS request through the Access-Control-Request-Headers header.",
						},
						"allowed_methods": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "An HTTP method that you allow the origin to execute. Valid values are GET, PUT, HEAD, POST, DELETE.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(stringvalidator.OneOf("GET", "PUT", "HEAD", "POST", "DELETE")),
							},
						},
						"allowed_origins": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "One or more origins you want customers to be able to access the bucket from.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
						"expose_headers": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "One or more headers in the response that you want customers to be able to access from their applications.",
						},
						"max_age_seconds": schema.Int64Attribute{
							Optional:    true,
							Description: "The time in seconds that your browser is to cache the preflight response for the specified resource.",
						},
					},
				},
			},
		},
	}
}

// Configure configures the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*services.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = clientBundle.S3Client
}

// Create creates the bucket CORS configuration.
func (r *bucketCorsConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketCorsConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateBucketCors(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket CORS configuration.
func (r *bucketCorsConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketCorsConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetBucketCors(ctx, data.Bucket)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket CORS configuration.
func (r *bucketCorsConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket CORS configuration.
func (r *bucketCorsConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketCorsConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateBucketCors(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket CORS configuration.
func (r *bucketCorsConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketCorsConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBucketCors(ctx, data.Bucket); err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}
