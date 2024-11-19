package objectstorage

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.ResourceWithImportState      = (*objectLockConfiguration)(nil)
	_ resource.ResourceWithConfigure        = (*objectLockConfiguration)(nil)
	_ resource.ResourceWithConfigValidators = (*objectLockConfiguration)(nil)
)

type objectLockConfiguration struct {
	client *objectstorage.Client
}

// NewObjectLockConfigurationResource creates a new resource for the bucket object lock configuration resource.
func NewObjectLockConfigurationResource() resource.Resource {
	return &objectLockConfiguration{}
}

// Metadata returns the metadata for the bucket object lock configuration resource.
func (r *objectLockConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_object_lock_configuration"
}

// Schema returns the schema for the bucket object lock configuration resource.
func (r *objectLockConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"object_lock_enabled": schema.StringAttribute{
				Description: "Specifies whether Object Lock is enabled for the bucket.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString("Enabled"),
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"default_retention": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Optional:   true,
								Validators: []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
							},
							"days": schema.Int64Attribute{
								Optional:   true,
								Validators: []validator.Int64{int64validator.AtLeast(1)},
							},
							"years": schema.Int64Attribute{
								Optional:   true,
								Validators: []validator.Int64{int64validator.AtLeast(1)},
							},
						},
					},
				},
			},
		},
	}
}

// ConfigValidators returns the config validators for the bucket object lock configuration resource.
func (r *objectLockConfiguration) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("rule").AtName("default_retention").AtName("days"),
			path.MatchRoot("rule").AtName("default_retention").AtName("years"),
		),
	}
}

// Configure configures the bucket object lock configuration resource.
func (r *objectLockConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the bucket object lock configuration resource.
func (r *objectLockConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectLockConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateObjectLock(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket object lock configuration resource.
func (r *objectLockConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectLockConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetObjectLock(ctx, data.Bucket)
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

// ImportState imports the state for the bucket object lock configuration resource.
func (r *objectLockConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket object lock configuration resource.
func (r *objectLockConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectLockConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateObjectLock(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to update resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket object lock configuration resource.
func (r *objectLockConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectLockConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Cannot be deleted
}
