package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var (
	_ resource.ResourceWithImportState = (*bucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketResource)(nil)
)

// NewBucketResource creates a new resource for the bucket resource.
func NewBucketResource() resource.Resource {
	return &bucketResource{}
}

type bucketResource struct {
	client *objectstorage.Client
}

type bucketResourceModel struct {
	Name              types.String   `tfsdk:"name"`
	Region            types.String   `tfsdk:"region"`
	ObjectLockEnabled types.Bool     `tfsdk:"object_lock_enabled"`
	ForceDestroy      types.Bool     `tfsdk:"force_destroy"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Tags              types.Map      `tfsdk:"tags"`
	ID                types.String   `tfsdk:"id"`
}

// Metadata returns the metadata for the bucket resource.
func (r *bucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket"
}

// Schema returns the schema for the bucket resource.
func (r *bucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Same value as name",
			},
			"name": schema.StringAttribute{
				Description: "The name of the bucket. It must start and end with a letter or number and contain only lowercase alphanumeric characters, hyphens, periods and underscores.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "The region of the bucket. Defaults to eu-central-3.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("eu-central-3"),
			},
			"object_lock_enabled": schema.BoolAttribute{
				Description: "Whether object lock is enabled for the bucket",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"force_destroy": schema.BoolAttribute{
				Description: "Whether all objects should be deleted from the bucket so that the bucket can be destroyed",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tags": schema.MapAttribute{
				Description: "A mapping of tags to assign to the bucket",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

// Configure configures the bucket resource.
func (r *bucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*services.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = clientBundle.S3Client
}

// Create creates the bucket.
func (r *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	if err := r.client.CreateBucket(ctx, data.Name, data.Region, data.ObjectLockEnabled, data.Tags, createTimeout); err != nil {
		resp.Diagnostics.AddError("failed to create bucket", err.Error())
		return
	}

	// Set computed values
	location, err := r.client.GetBucketLocation(ctx, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket location", err.Error())
		return
	}

	data.ID = data.Name
	data.Region = location
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket.
func (r *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket, found, err := r.client.GetBucket(ctx, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Bucket API error", err.Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Tags = bucket.Tags
	data.Region = bucket.Region
	data.ObjectLockEnabled = bucket.ObjectLockEnabled
	data.ID = data.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket.
func (r *bucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Update updates the bucket.
func (r *bucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *bucketResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Tags.Equal(state.Tags) {
		if err := r.client.UpdateBucketTags(ctx, plan.Name.ValueString(), tags.NewFromMap(plan.Tags), tags.NewFromMap(state.Tags)); err != nil {
			resp.Diagnostics.AddError("failed to update tags", err.Error())
			return
		}
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the bucket.
func (r *bucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *bucketResourceModel
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

	if err := r.client.DeleteBucket(ctx, data.Name, data.ObjectLockEnabled, data.ForceDestroy, deleteTimeout); err != nil {
		resp.Diagnostics.AddError("failed to delete bucket", err.Error())
		return
	}
}
