package s3

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

// Metadata returns the metadata for the bucket resource.
func (r *bucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket"
}

// Schema returns the schema for the bucket resource.
func (r *bucketResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "The region of the bucket",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("eu-central-3"),
			},
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
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	_, err := r.client.BucketsApi.CreateBucket(ctx, data.Name.ValueString()).CreateBucketConfiguration(createBucketConfig).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", formatXMLError(err).Error())
		return
	}

	// Wait for bucket creation
	err = backoff.Retry(func() error {
		return bucketExists(ctx, r.client, data.Name.ValueString())
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", fmt.Sprintf("error verifying bucket creation: %s", err))
		return
	}

	location, _, err := r.client.BucketsApi.GetBucketLocation(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to get bucket location", err.Error())
		return
	}

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

	data.Region = types.StringValue(location.GetLocationConstraint())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket.
func (r *bucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Update updates the bucket.
func (r *bucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Nothing to update for now
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	apiResponse, err := r.client.BucketsApi.DeleteBucket(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete bucket", formatXMLError(err).Error())
		return
	}

	// Wait for deletion
	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout))
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

func bucketExists(ctx context.Context, client *s3.APIClient, bucket string) error {
	apiResponse, err := client.BucketsApi.HeadBucket(ctx, bucket).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return fmt.Errorf("bucket not found")
		}
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", formatXMLError(err)))
	}
	return nil
}
