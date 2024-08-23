package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*bucketPublicAccessBlockResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketPublicAccessBlockResource)(nil)
)

// ErrBucketPublicAccessBlockNotFound returned for 404
var ErrBucketPublicAccessBlockNotFound = errors.New("s3 bucket public access block not found")

// NewBucketPublicAccessBlockResource creates a new resource for the bucket public access block resource.
func NewBucketPublicAccessBlockResource() resource.Resource {
	return &bucketPublicAccessBlockResource{}
}

type bucketPublicAccessBlockResource struct {
	client *s3.APIClient
}

type bucketPublicAccessBlockResourceModel struct {
	Bucket                types.String `tfsdk:"bucket"`
	BlockPublicACLS       types.Bool   `tfsdk:"block_public_acls"`
	BlockPublicPolicy     types.Bool   `tfsdk:"block_public_policy"`
	IgnorePublicACLS      types.Bool   `tfsdk:"ignore_public_acls"`
	RestrictPublicBuckets types.Bool   `tfsdk:"restrict_public_buckets"`
}

// Metadata returns the metadata for the bucket resource.
func (r *bucketPublicAccessBlockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_public_access_block"
}

// Schema returns the schema for the bucket resource.
func (r *bucketPublicAccessBlockResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"block_public_acls": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"block_public_policy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"ignore_public_acls": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"restrict_public_buckets": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

// Configure configures the bucket resource.
func (r *bucketPublicAccessBlockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create creates the bucket.
func (r *bucketPublicAccessBlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPublicAccessBlockResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	requestInput := putBucketPublicAccessBlockInput(data)
	_, err := r.client.PublicAccessBlockApi.PutPublicAccessBlock(ctx, data.Bucket.ValueString()).BlockPublicAccessPayload(requestInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket public access block", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket.
func (r *bucketPublicAccessBlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPublicAccessBlockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := GetBucketPublicAccessBlock(ctx, r.client, data.Bucket.ValueString())
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

// ImportState imports the state of the bucket.
func (r *bucketPublicAccessBlockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket.
func (r *bucketPublicAccessBlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketPublicAccessBlockResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	requestInput := putBucketPublicAccessBlockInput(data)
	_, err := r.client.PublicAccessBlockApi.PutPublicAccessBlock(ctx, data.Bucket.ValueString()).BlockPublicAccessPayload(requestInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to update bucket public access block", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket.
func (r *bucketPublicAccessBlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPublicAccessBlockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if apiResponse, err := r.client.PublicAccessBlockApi.DeletePublicAccessBlock(ctx, data.Bucket.ValueString()).Execute(); err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete bucket public access block", err.Error())
		return
	}

}

// GetBucketPublicAccessBlock retrieves the public access block for the bucket
func GetBucketPublicAccessBlock(ctx context.Context, client *s3.APIClient, bucketName string) (*s3.BlockPublicAccessOutput, error) {
	response, apiResponse, err := client.PublicAccessBlockApi.GetPublicAccessBlock(ctx, bucketName).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil, ErrBucketPublicAccessBlockNotFound
		}
		return nil, err
	}
	return response, nil
}

func putBucketPublicAccessBlockInput(model bucketPublicAccessBlockResourceModel) s3.BlockPublicAccessPayload {
	input := s3.BlockPublicAccessPayload{
		BlockPublicPolicy:     model.BlockPublicPolicy.ValueBoolPointer(),
		IgnorePublicAcls:      model.IgnorePublicACLS.ValueBoolPointer(),
		BlockPublicAcls:       model.BlockPublicACLS.ValueBoolPointer(),
		RestrictPublicBuckets: model.RestrictPublicBuckets.ValueBoolPointer(),
	}
	return input
}
