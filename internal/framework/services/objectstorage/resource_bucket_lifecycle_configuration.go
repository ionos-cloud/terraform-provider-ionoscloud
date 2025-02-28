package objectstorage

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.ResourceWithImportState = (*bucketLifecycleConfiguration)(nil)
	_ resource.ResourceWithConfigure   = (*bucketLifecycleConfiguration)(nil)
)

type bucketLifecycleConfiguration struct {
	client *objectstorage.Client
}

// NewBucketLifecycleConfigurationResource creates a new resource for the bucket lifecycle configuration resource.
func NewBucketLifecycleConfigurationResource() resource.Resource {
	return &bucketLifecycleConfiguration{}
}

// Metadata returns the metadata for the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_lifecycle_configuration"
}

// Schema returns the schema for the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "The name of the bucket.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "A list of lifecycle rules for objects in the bucket.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Unique identifier for the rule.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
						"prefix": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 1024),
							},
							Description: "Object key prefix identifying one or more objects to which the rule applies.",
						},
						"status": schema.StringAttribute{
							Required:    true,
							Description: "Whether the rule is currently being applied. Valid values: Enabled or Disabled.",
							Validators: []validator.String{
								stringvalidator.OneOf("Enabled", "Disabled"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"expiration": schema.SingleNestedBlock{
							Description: "A lifecycle rule for when an object expires.",
							Attributes: map[string]schema.Attribute{
								"days": schema.Int64Attribute{
									Optional:    true,
									Description: "Specifies the number of days after object creation when the object expires. Required if 'date' is not specified.",
									Validators: []validator.Int64{
										int64validator.AtLeast(1),
									},
								},
								"date": schema.StringAttribute{
									Optional:    true,
									Description: "Specifies the date when the object expires. Required if 'days' is not specified.",
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									Optional:    true,
									Description: "Indicates whether IONOS Object Storage Object Storage will remove a delete marker with no noncurrent versions. If set to true, the delete marker will be expired; if set to false the policy takes no operation. This cannot be specified with Days or Date in a Lifecycle Expiration Policy.",
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedBlock{
							Description: "A lifecycle rule for when non-current object versions expire.",
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									Optional:    true,
									Description: "Specifies the number of days an object is noncurrent before IONOS Object Storage can perform the associated action.",
								},
							},
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(path.Expressions{path.MatchRelative().AtName("noncurrent_days")}...),
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									Optional:    true,
									Description: "Specifies the number of days after which IONOS Object Storage Object Storage aborts an incomplete multipart upload.",
								},
							},
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(path.Expressions{path.MatchRelative().AtName("days_after_initiation")}...),
							},
							Description: "Specifies the days since the initiation of an incomplete multipart upload that IONOS Object Storage Object Storage will wait before permanently removing all parts of the upload.",
						},
					},
				},
			},
		},
	}
}

// Configure configures the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = clientBundle.S3Client
}

// Create creates the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.BucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateBucketLifecycle(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.BucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetBucketLifecycle(ctx, data.Bucket)
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

// ImportState imports the state of the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.BucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateBucketLifecycle(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to update resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.BucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBucketLifecycle(ctx, data.Bucket); err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}
