package s3

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	tfs3 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"

	s3 "github.com/ionos-cloud/sdk-go-s3"

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
	client *s3.APIClient
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

	client, ok := req.ProviderData.(*s3.APIClient)
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
func (r *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createBucketConfig := s3.CreateBucketConfiguration{
		LocationConstraint: data.Region.ValueStringPointer(),
	}
	createTimeout, diag := data.Timeouts.Create(ctx, utils.DefaultTimeout)
	if diag != nil {
		resp.Diagnostics = diag
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	createReq := r.client.BucketsApi.CreateBucket(ctx, data.Name.ValueString()).CreateBucketConfiguration(createBucketConfig)
	if !data.ObjectLockEnabled.IsNull() {
		createReq = createReq.XAmzBucketObjectLockEnabled(data.ObjectLockEnabled.ValueBool())
	}
	_, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", formatXMLError(err).Error())
		return
	}

	// Wait for bucket creation
	err = backoff.Retry(func() error {
		return BucketExists(ctx, r.client, data.Name.ValueString())
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(createTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", fmt.Sprintf("error verifying bucket creation: %s", err))
		return
	}

	if err = tfs3.CreateBucketTags(ctx, r.client, data.Name.ValueString(), tags.NewFromTFMap(data.Tags)); err != nil {
		resp.Diagnostics.AddError("failed to create tags", err.Error())
		return
	}

	location, _, err := r.client.BucketsApi.GetBucketLocation(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket location", err.Error())
		return
	}
	data.ID = data.Name
	data.Region = types.StringPointerValue(location.LocationConstraint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket.
func (r *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := r.client.BucketsApi.HeadBucket(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Failed to read bucket", formatXMLError(err).Error())
		return
	}

	location, _, err := r.client.BucketsApi.GetBucketLocation(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read bucket location", formatXMLError(err).Error())
		return
	}

	objLockConfig, apiResponse, err := r.client.ObjectLockApi.GetObjectLockConfiguration(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		if !apiResponse.HttpNotFound() {
			resp.Diagnostics.AddError("Failed to read object lock configuration", formatXMLError(err).Error())
			return
		}
		data.ObjectLockEnabled = types.BoolValue(false)
	}

	if objLockConfig != nil && objLockConfig.ObjectLockEnabled != nil {
		data.ObjectLockEnabled = types.BoolValue(*objLockConfig.ObjectLockEnabled == "Enabled")
	}

	tags, diags := getBucketTags(ctx, r.client, data.Name.ValueString())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Tags = tags
	data.Region = types.StringPointerValue(location.GetLocationConstraint())
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
		if err := tfs3.UpdateBucketTags(ctx, r.client, plan.Name.ValueString(), tags.NewFromTFMap(plan.Tags), tags.NewFromTFMap(state.Tags)); err != nil {
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
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diag := data.Timeouts.Delete(ctx, utils.DefaultTimeout)
	if diag != nil {
		resp.Diagnostics = diag
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	apiResponse, err := r.client.BucketsApi.DeleteBucket(ctx, data.Name.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return
	}

	if isBucketNotEmptyError(err) && data.ForceDestroy.ValueBool() {
		_, err := tfs3.EmptyBucket(ctx, r.client, data.Name.ValueString(), data.ObjectLockEnabled.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("failed to empty bucket", err.Error())
			return
		}

		r.Delete(ctx, req, resp)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to delete bucket", formatXMLError(err).Error())
		return
	}

	// Wait for deletion
	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(deleteTimeout))
	err = backoff.Retry(func() error {
		return IsBucketDeleted(ctx, r.client, data.Name.ValueString())
	}, backOff)

	if err != nil {
		resp.Diagnostics.AddError("failed to delete bucket", fmt.Sprintf("error verifying bucket deletion: %s", err))
		return
	}
}

// IsBucketDeleted checks if the bucket is deleted.
func IsBucketDeleted(ctx context.Context, client *s3.APIClient, bucket string) error {
	apiResponse, err := client.BucketsApi.HeadBucket(ctx, bucket).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", err))
	}
	return fmt.Errorf("bucket still exists")
}

// BucketExists checks if the bucket exists.
func BucketExists(ctx context.Context, client *s3.APIClient, bucket string) error {
	apiResponse, err := client.BucketsApi.HeadBucket(ctx, bucket).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return fmt.Errorf("bucket not found")
		}
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", formatXMLError(err)))
	}
	return nil
}

func getBucketTags(ctx context.Context, client *s3.APIClient, bucketName string) (types.Map, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	output, apiResponse, err := client.TaggingApi.GetBucketTagging(ctx, bucketName).Execute()
	if apiResponse.HttpNotFound() {
		return types.MapNull(types.StringType), nil
	}

	if err != nil {
		diags.AddError("failed to get bucket tags", formatXMLError(err).Error())
		return types.MapNull(types.StringType), diags
	}

	tags, diagErr := getTagsFromAPIResponse(ctx, output)
	if diagErr != nil {
		return types.MapNull(types.StringType), diagErr
	}

	return tags, nil
}

func getTagsFromAPIResponse(ctx context.Context, response *s3.GetBucketTaggingOutput) (types.Map, diag.Diagnostics) {
	if response == nil || response.TagSet == nil {
		return types.Map{}, nil

	}

	result := make(map[string]string)
	for _, tag := range *response.TagSet {
		if tag.Key == nil || tag.Value == nil {
			continue
		}

		result[*tag.Key] = *tag.Value
	}

	tfResult, diagErr := types.MapValueFrom(ctx, types.StringType, result)
	if diagErr != nil {
		return types.Map{}, diagErr
	}

	return tfResult, nil
}
