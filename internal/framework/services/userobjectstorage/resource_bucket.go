package userobjectstorage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/userobjectstorage"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var (
	_ resource.ResourceWithImportState = (*bucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketResource)(nil)
)

// NewBucketResource creates a new resource for the user object storage bucket.
func NewBucketResource() resource.Resource {
	return &bucketResource{}
}

type bucketResource struct {
	client *userobjectstorage.Client
}

type bucketResourceModel struct {
	ForceDestroy      types.Bool     `tfsdk:"force_destroy"`
	ID                types.String   `tfsdk:"id"`
	Name              types.String   `tfsdk:"name"`
	ObjectLockEnabled types.Bool     `tfsdk:"object_lock_enabled"`
	Region            types.String   `tfsdk:"region"`
	Tags              types.Map      `tfsdk:"tags"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

// Metadata sets the resource type name.
func (r *bucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_object_storage_bucket"
}

// Schema returns the schema for the resource.
func (r *bucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"force_destroy": schema.BoolAttribute{
				Description: "When true, all objects are deleted from the bucket before destroying it, allowing a non-empty bucket to be destroyed.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Same value as name.",
			},
			"name": schema.StringAttribute{
				Description: "The bucket name. Must be 3–63 characters, lowercase alphanumeric, hyphens, periods, or underscores.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object_lock_enabled": schema.BoolAttribute{
				Description: "Whether object lock is enabled for the bucket. Cannot be changed after creation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.MapAttribute{
				Description: "A mapping of tags to assign to the bucket.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"region": schema.StringAttribute{
				Description: "The region where the bucket is created. Valid values: 'de' (Frankfurt), 'eu-central-2' (Berlin), 'eu-south-2' (Logroño).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(userobjectstorage.ValidRegions...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

// Configure wires the provider client.
func (r *bucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = clientBundle.UserS3Client
}

// Create provisions the bucket.
func (r *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("user object storage client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.BucketDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	if err := r.client.CreateBucket(ctx, data.Name.ValueString(), data.Region.ValueString(), data.ObjectLockEnabled.ValueBool(), createTimeout); err != nil {
		resp.Diagnostics.AddError("failed to create bucket", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	if !data.Tags.IsNull() && len(data.Tags.Elements()) > 0 {
		if err := r.client.PutBucketTags(ctx, data.Name.ValueString(), data.Region.ValueString(), tags.NewFromMap(data.Tags)); err != nil {
			resp.Diagnostics.AddError("failed to set bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
			return
		}
	}

	data.ID = data.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the state from the API.
func (r *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("user object storage client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.BucketDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	found, err := r.client.GetBucket(ctx, data.Name.ValueString(), data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Bucket API error", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	objectLockEnabled, err := r.client.GetObjectLockEnabled(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get object lock configuration", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	rawTags, err := r.client.GetBucketTags(ctx, data.Name.ValueString(), data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	tagsMap, err := tags.KeyValueTags(rawTags).ToMap(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
		return
	}

	data.ObjectLockEnabled = types.BoolValue(objectLockEnabled)
	data.Tags = tagsMap
	data.ID = data.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles in-place changes (force_destroy and tags can change without replacement).
func (r *bucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Tags.Equal(state.Tags) {
		if plan.Tags.IsNull() || len(plan.Tags.Elements()) == 0 {
			if err := r.client.DeleteBucketTags(ctx, plan.Name.ValueString(), plan.Region.ValueString()); err != nil {
				resp.Diagnostics.AddError("failed to delete bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: plan.Name.ValueString()}).Error())
				return
			}
		} else {
			if err := r.client.PutBucketTags(ctx, plan.Name.ValueString(), plan.Region.ValueString(), tags.NewFromMap(plan.Tags)); err != nil {
				resp.Diagnostics.AddError("failed to update bucket tags", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: plan.Name.ValueString()}).Error())
				return
			}
		}
	}

	plan.ID = plan.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete destroys the bucket.
func (r *bucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("user object storage client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	if err := r.client.DeleteBucket(ctx, data.Name.ValueString(), data.Region.ValueString(), data.ForceDestroy.ValueBool(), deleteTimeout); err != nil {
		resp.Diagnostics.AddError("failed to delete bucket", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: data.Name.ValueString()}).Error())
	}
}

// ImportState supports `terraform import ionoscloud_user_object_storage_bucket.x name` or `region:name`.
func (r *bucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var name, region string

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("expected 'region:bucket_name'. Got: %q", req.ID))
		return
	}
	region = parts[0]
	name = parts[1]
	if region == "" || name == "" {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("both region and bucket name must be non-empty. Got: %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region"), region)...)
	req.ID = name
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
