package s3management

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	s3management "github.com/ionos-cloud/sdk-go-s3-management"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var (
	_ resource.ResourceWithImportState = (*accesskeyResource)(nil)
	_ resource.ResourceWithConfigure   = (*accesskeyResource)(nil)
)

// NewAccesskeyResource creates a new resource for the accesskey resource.
func NewAccesskeyResource() resource.Resource {
	return &accesskeyResource{}
}

type accesskeyResource struct {
	client *services.SdkBundle
}

type accesskeyResourceModel struct {
	AccessKey       types.String   `tfsdk:"accesskey"`
	SecretKey       types.String   `tfsdk:"secretkey"`
	CanonicalUserId types.String   `tfsdk:"canonical_user_id"`
	ContractUserId  types.String   `tfsdk:"contract_user_id"`
	Description     types.String   `tfsdk:"description"`
	ID              types.String   `tfsdk:"id"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

// Metadata returns the metadata for the bucket resource.
func (r *accesskeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_accesskey"
}

// Schema returns the schema for the bucket resource.
func (r *accesskeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID (UUID) of the AccessKey.",
			},
			"description": schema.StringAttribute{
				Description: "Description of the Access key.",
				Optional:    true,
			},
			"accesskey": schema.StringAttribute{
				Description: "Access key metadata is a string of 92 characters.",
				Computed:    true,
			},
			"secretkey": schema.StringAttribute{
				Description: "The secret key of the Access key.",
				Computed:    true,
			},
			"canonical_user_id": schema.StringAttribute{
				Description: "The canonical user ID which is valid for user-owned buckets.",
				Computed:    true,
			},
			"contract_user_id": schema.StringAttribute{
				Description: "The contract user ID which is valid for contract-owned buckets",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
	}
}

// Configure configures the bucket resource.
func (r *accesskeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*services.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the bucket.
func (r *accesskeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *accesskeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	var accessKey = s3management.AccessKeyCreate{
		Properties: &s3management.AccessKeyProperties{
			Description: data.Description.ValueStringPointer(),
		},
	}
	accessKeyResponse, _, err := r.client.S3ManagementClient.CreateAccessKey(ctx, accessKey, createTimeout)
	if err != nil {
		resp.Diagnostics.AddError("failed to create accessKey", err.Error())
		return
	}

	data.ID = basetypes.NewStringPointerValue(accessKeyResponse.Id)

	accessKeyRead, _, err := r.client.S3ManagementClient.GetAccessKey(ctx, *accessKeyResponse.Id)
	if err != nil {
		resp.Diagnostics.AddError("Access Key API error", err.Error())
		return
	}

	data.AccessKey = basetypes.NewStringPointerValue(accessKeyRead.Properties.AccessKey)
	data.CanonicalUserId = basetypes.NewStringPointerValue(accessKeyRead.Properties.CanonicalUserId)
	data.ContractUserId = basetypes.NewStringPointerValue(accessKeyRead.Properties.ContractUserId)
	data.Description = basetypes.NewStringPointerValue(accessKeyRead.Properties.Description)
	data.SecretKey = basetypes.NewStringPointerValue(accessKeyRead.Properties.SecretKey)
	data.ID = basetypes.NewStringPointerValue(accessKeyRead.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket.
func (r *accesskeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data accesskeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessKey, _, err := r.client.S3ManagementClient.GetAccessKey(ctx, data.ID.String())
	if err != nil {
		resp.Diagnostics.AddError("Access Key API error", err.Error())
		return
	}

	data.AccessKey = basetypes.NewStringPointerValue(accessKey.Properties.AccessKey)
	data.CanonicalUserId = basetypes.NewStringPointerValue(accessKey.Properties.CanonicalUserId)
	data.ContractUserId = basetypes.NewStringPointerValue(accessKey.Properties.ContractUserId)
	data.Description = basetypes.NewStringPointerValue(accessKey.Properties.Description)
	data.SecretKey = basetypes.NewStringPointerValue(accessKey.Properties.SecretKey)
	data.ID = basetypes.NewStringPointerValue(accessKey.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the accessKey.
func (r *accesskeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the bucket.
func (r *accesskeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *accesskeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data *accesskeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := data.Timeouts.Update(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	var accessKey = s3management.AccessKeyEnsure{
		Properties: &s3management.AccessKeyProperties{
			Description: data.Description.ValueStringPointer(),
		},
	}
	accessKeyResponse, _, err := r.client.S3ManagementClient.UpdateAccessKey(ctx, data.ID.String(), accessKey, updateTimeout)
	if err != nil {
		resp.Diagnostics.AddError("failed to update accessKey", err.Error())
		return
	}

	data.ID = basetypes.NewStringPointerValue(accessKeyResponse.Id)

	accessKeyRead, _, err := r.client.S3ManagementClient.GetAccessKey(ctx, *accessKeyResponse.Id)
	if err != nil {
		resp.Diagnostics.AddError("Access Key API error", err.Error())
		return
	}

	data.AccessKey = basetypes.NewStringPointerValue(accessKeyRead.Properties.AccessKey)
	data.CanonicalUserId = basetypes.NewStringPointerValue(accessKeyRead.Properties.CanonicalUserId)
	data.ContractUserId = basetypes.NewStringPointerValue(accessKeyRead.Properties.ContractUserId)
	data.Description = basetypes.NewStringPointerValue(accessKeyRead.Properties.Description)
	data.SecretKey = basetypes.NewStringPointerValue(accessKeyRead.Properties.SecretKey)
	data.ID = basetypes.NewStringPointerValue(accessKeyRead.Id)

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the accessKey.
func (r *accesskeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *accesskeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	if _, err := r.client.S3ManagementClient.DeleteAccessKey(ctx, data.ID.String(), deleteTimeout); err != nil {
		resp.Diagnostics.AddError("failed to delete bucket", err.Error())
		return
	}
}
