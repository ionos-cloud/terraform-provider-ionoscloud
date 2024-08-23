package s3

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*bucketLifecycleConfiguration)(nil)
	_ resource.ResourceWithConfigure   = (*bucketLifecycleConfiguration)(nil)
)

type bucketLifecycleConfiguration struct {
	client *s3.APIClient
}

type bucketLifecycleConfigurationModel struct {
	Bucket types.String    `tfsdk:"bucket"`
	Rule   []lifecycleRule `tfsdk:"rule"`
}

type lifecycleRule struct {
	ID                             types.String                    `tfsdk:"id"`
	Prefix                         types.String                    `tfsdk:"prefix"`
	Status                         types.String                    `tfsdk:"status"`
	Expiration                     *expiration                     `tfsdk:"expiration"`
	NoncurrentVersionExpiration    *noncurrentVersionExpiration    `tfsdk:"noncurrent_version_expiration"`
	AbortIncompleteMultipartUpload *abortIncompleteMultipartUpload `tfsdk:"abort_incomplete_multipart_upload"`
}

type expiration struct {
	Days                      types.Int64  `tfsdk:"days"`
	Date                      types.String `tfsdk:"date"`
	ExpiredObjectDeleteMarker types.Bool   `tfsdk:"expired_object_delete_marker"`
}

type noncurrentVersionExpiration struct {
	NoncurrentDays types.Int64 `tfsdk:"noncurrent_days"`
}

type abortIncompleteMultipartUpload struct {
	DaysAfterInitiation types.Int64 `tfsdk:"days_after_initiation"`
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
				Description: "The name of the S3 bucket.",
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
								stringvalidator.LengthBetween(1, 1024),
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
									Description: "Indicates whether IONOS S3 Object Storage will remove a delete marker with no noncurrent versions. If set to true, the delete marker will be expired; if set to false the policy takes no operation. This cannot be specified with Days or Date in a Lifecycle Expiration Policy.",
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedBlock{
							Description: "A lifecycle rule for when non-current object versions expire.",
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									Required:    true,
									Description: "Specifies the number of days an object is noncurrent before Amazon S3 can perform the associated action.",
								},
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									Required:    true,
									Description: "Specifies the number of days after which IONOS S3 Object Storage aborts an incomplete multipart upload.",
								},
							},
							Description: "Specifies the days since the initiation of an incomplete multipart upload that IONOS S3 Object Storage will wait before permanently removing all parts of the upload.",
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

// Create creates the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := putWithContentMD5Header(ctx, r.client, data)
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

	var data *bucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.LifecycleApi.GetBucketLifecycle(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketLifecycleConfigurationModelFromAPIResponse(output, data)
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

	var data *bucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := putWithContentMD5Header(ctx, r.client, data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	output, _, err := r.client.LifecycleApi.GetBucketLifecycle(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketLifecycleConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket lifecycle configuration resource.
func (r *bucketLifecycleConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketLifecycleConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.LifecycleApi.DeleteBucketLifecycle(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}

func putWithContentMD5Header(ctx context.Context, client *s3.APIClient, data *bucketLifecycleConfigurationModel) error {
	body := buildBucketLifecycleConfigurationFromModel(data)
	hash, err := getMD5Hash(body)
	if err != nil {
		return fmt.Errorf("failed to generate MD5 sum: %s", err.Error())
	}

	_, err = client.LifecycleApi.PutBucketLifecycle(ctx, data.Bucket.ValueString()).PutBucketLifecycleRequest(body).ContentMD5(base64.StdEncoding.EncodeToString([]byte(hash))).Execute()
	return err
}

func buildBucketLifecycleConfigurationModelFromAPIResponse(output *s3.GetBucketLifecycleOutput, data *bucketLifecycleConfigurationModel) *bucketLifecycleConfigurationModel {
	data.Rule = buildRulesFromAPIResponse(output.Rules)
	return data
}

func buildRulesFromAPIResponse(rules *[]s3.Rule) []lifecycleRule {
	if rules == nil {
		return nil
	}

	result := make([]lifecycleRule, 0, len(*rules))
	for _, r := range *rules {
		result = append(result, lifecycleRule{
			ID:                             types.StringPointerValue(r.ID),
			Prefix:                         types.StringPointerValue(r.Prefix),
			Status:                         types.StringValue(string(*r.Status)),
			Expiration:                     buildExpirationFromAPIResponse(r.Expiration),
			NoncurrentVersionExpiration:    buildNoncurrentVersionExpirationFromAPIResponse(r.NoncurrentVersionExpiration),
			AbortIncompleteMultipartUpload: buildAbortIncompleteMultipartUploadFromAPIResponse(r.AbortIncompleteMultipartUpload),
		})
	}

	return result
}

func buildExpirationFromAPIResponse(exp *s3.LifecycleExpiration) *expiration {
	if exp == nil {
		return nil
	}

	return &expiration{
		Days:                      types.Int64PointerValue(toInt64(exp.Days)),
		Date:                      types.StringPointerValue(exp.Date),
		ExpiredObjectDeleteMarker: types.BoolPointerValue(exp.ExpiredObjectDeleteMarker),
	}
}

func buildNoncurrentVersionExpirationFromAPIResponse(exp *s3.NoncurrentVersionExpiration) *noncurrentVersionExpiration {
	if exp == nil {
		return nil
	}

	return &noncurrentVersionExpiration{
		NoncurrentDays: types.Int64PointerValue(toInt64(exp.NoncurrentDays)),
	}
}

func buildAbortIncompleteMultipartUploadFromAPIResponse(abort *s3.AbortIncompleteMultipartUpload) *abortIncompleteMultipartUpload {
	if abort == nil {
		return nil
	}

	return &abortIncompleteMultipartUpload{
		DaysAfterInitiation: types.Int64PointerValue(toInt64(abort.DaysAfterInitiation)),
	}
}

func buildBucketLifecycleConfigurationFromModel(data *bucketLifecycleConfigurationModel) s3.PutBucketLifecycleRequest {
	return s3.PutBucketLifecycleRequest{
		Rules: buildRulesFromModel(data.Rule),
	}
}

func buildRulesFromModel(rules []lifecycleRule) *[]s3.Rule {
	if rules == nil {
		return nil
	}

	result := make([]s3.Rule, 0, len(rules))
	for _, r := range rules {
		result = append(result, s3.Rule{
			ID:                             r.ID.ValueStringPointer(),
			Prefix:                         r.Prefix.ValueStringPointer(),
			Status:                         s3.ExpirationStatus(r.Status.ValueString()).Ptr(),
			Expiration:                     buildExpirationFromModel(r.Expiration),
			NoncurrentVersionExpiration:    buildNoncurrentVersionExpirationFromModel(r.NoncurrentVersionExpiration),
			AbortIncompleteMultipartUpload: buildAbortIncompleteMultipartUploadFromModel(r.AbortIncompleteMultipartUpload),
		})
	}

	return &result
}

func buildExpirationFromModel(expiration *expiration) *s3.LifecycleExpiration {
	if expiration == nil {
		return nil
	}

	return &s3.LifecycleExpiration{
		Days:                      toInt32(expiration.Days.ValueInt64Pointer()),
		Date:                      expiration.Date.ValueStringPointer(),
		ExpiredObjectDeleteMarker: expiration.ExpiredObjectDeleteMarker.ValueBoolPointer(),
	}
}

func buildNoncurrentVersionExpirationFromModel(expiration *noncurrentVersionExpiration) *s3.NoncurrentVersionExpiration {
	if expiration == nil {
		return nil
	}

	return &s3.NoncurrentVersionExpiration{
		NoncurrentDays: toInt32(expiration.NoncurrentDays.ValueInt64Pointer()),
	}
}

func buildAbortIncompleteMultipartUploadFromModel(abort *abortIncompleteMultipartUpload) *s3.AbortIncompleteMultipartUpload {
	if abort == nil {
		return nil
	}

	return &s3.AbortIncompleteMultipartUpload{
		DaysAfterInitiation: toInt32(abort.DaysAfterInitiation.ValueInt64Pointer()),
	}
}
