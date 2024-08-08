package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*bucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketResource)(nil)
)

type bucketVersioningResource struct {
	client *s3.APIClient
}

type bucketVersioningResourceModel struct {
	Bucket                  types.String             `tfsdk:"bucket"`
	VersioningConfiguration *versioningConfiguration `tfsdk:"versioning_configuration"`
}

type versioningConfiguration struct {
	Status    types.String `tfsdk:"status"`
	MfaDelete types.String `tfsdk:"mfa_delete"`
}

func (v versioningConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status":     types.StringType,
		"mfa_delete": types.StringType,
	}
}

// NewBucketVersioningResource creates a new resource for the bucket versioning resource.
func NewBucketVersioningResource() resource.Resource {
	return &bucketVersioningResource{}
}

// Metadata returns the metadata for the bucket versioning resource.
func (r *bucketVersioningResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_versioning"
}

// Schema returns the schema for the bucket versioning resource.
func (r *bucketVersioningResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "The name of the bucket.",
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"versioning_configuration": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Required:    true,
						Description: "The versioning status of the bucket.",
						Validators:  []validator.String{stringvalidator.OneOf("Enabled", "Suspended")},
					},
					"mfa_delete": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The MFA delete status of the bucket.",
						Default:     stringdefault.StaticString("Disabled"),
					},
				},
				Description: "The versioning configuration of the bucket.",
			},
		},
	}
}

// Configure configures the bucket versioning resource.
func (r *bucketVersioningResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the bucket versioning resource.
func (r *bucketVersioningResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketVersioningResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.VersioningApi.PutBucketVersioning(ctx, data.Bucket.ValueString()).PutBucketVersioningRequest(buildPutRequestFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create bucket versioning", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket versioning resource.
func (r *bucketVersioningResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketVersioningResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.VersioningApi.GetBucketVersioning(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to read bucket versioning", err.Error())
		return
	}

	data = buildModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state for the bucket versioning resource.
func (r *bucketVersioningResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket versioning resource.
func (r *bucketVersioningResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketVersioningResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.VersioningApi.PutBucketVersioning(ctx, data.Bucket.ValueString()).PutBucketVersioningRequest(buildPutRequestFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update bucket versioning", err.Error())
		return
	}

	output, apiResponse, err := r.client.VersioningApi.GetBucketVersioning(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to read bucket versioning", err.Error())
		return
	}

	data = buildModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket versioning resource.
func (r *bucketVersioningResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketVersioningResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Removing S3 bucket versioning for un-versioned bucket from state
	if data.VersioningConfiguration.Status.ValueString() == string(s3.BUCKETVERSIONINGSTATUS_SUSPENDED) {
		return
	}

	_, err := r.client.VersioningApi.PutBucketVersioning(ctx, data.Bucket.ValueString()).
		PutBucketVersioningRequest(s3.PutBucketVersioningRequest{
			Status: s3.BUCKETVERSIONINGSTATUS_SUSPENDED.Ptr(),
		}).Execute()
	if isInvalidStateBucketWithObjectLock(err) {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to create bucket versioning", err.Error())
		return
	}
}

func buildModelFromAPIResponse(output *s3.GetBucketVersioningOutput, data *bucketVersioningResourceModel) *bucketVersioningResourceModel {
	var versioningConfiguration versioningConfiguration
	if output.Status != nil {
		versioningConfiguration.Status = types.StringValue(string(*output.Status))
	}

	if output.MfaDelete != nil {
		versioningConfiguration.MfaDelete = types.StringValue(string(*output.MfaDelete))
	}

	built := bucketVersioningResourceModel{
		Bucket:                  data.Bucket,
		VersioningConfiguration: &versioningConfiguration,
	}

	return &built
}

func buildPutRequestFromModel(data *bucketVersioningResourceModel) s3.PutBucketVersioningRequest {
	var request s3.PutBucketVersioningRequest
	if !data.VersioningConfiguration.Status.IsNull() {
		request.Status = s3.BucketVersioningStatus(data.VersioningConfiguration.Status.ValueString()).Ptr()
	}

	if !data.VersioningConfiguration.MfaDelete.IsNull() {
		request.MfaDelete = s3.MfaDeleteStatus(data.VersioningConfiguration.MfaDelete.ValueString()).Ptr()
	}
	return request
}
