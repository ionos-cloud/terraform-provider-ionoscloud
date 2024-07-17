package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	s3 "github.com/ionos-cloud/sdk-go-s3"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var (
	_ resource.ResourceWithImportState = (*s3BucketObjectResource)(nil)
	_ resource.ResourceWithConfigure   = (*s3BucketObjectResource)(nil)
)

// NewBucketObjectResource creates a new resource for the bucket resource.
func NewS3BucketObjectResource() resource.Resource {
	return &s3BucketObjectResource{}
}

type s3BucketObjectResource struct {
	client *s3.APIClient
}

type s3BucketObjectModel struct {
	BucketName                                types.String `tfsdk:"bucket_name"`
	ObjectKey                                 types.String `tfsdk:"object_key"`
	Body                                      types.String `tfsdk:"body"`
	BodyFile                                  types.String `tfsdk:"body_file"`
	CacheControl                              types.String `tfsdk:"cache_control"`
	ContentDisposition                        types.String `tfsdk:"content_disposition"`
	ContentEncoding                           types.String `tfsdk:"content_encoding"`
	ContentLanguage                           types.String `tfsdk:"content_language"`
	ContentLength                             types.Int64  `tfsdk:"content_length"`
	ContentMD5                                types.String `tfsdk:"content_md5"`
	ContentType                               types.String `tfsdk:"content_type"`
	Expires                                   types.String `tfsdk:"expires"`
	XAmzServerSideEncryption                  types.String `tfsdk:"x_amz_server_side_encryption"`
	XAmzStorageClass                          types.String `tfsdk:"x_amz_storage_class"`
	XAmzWebsiteRedirectLocation               types.String `tfsdk:"x_amz_website_redirect_location"`
	XAmzServerSideEncryptionCustomerAlgorithm types.String `tfsdk:"x_amz_server_side_encryption_customer_algorithm"`
	XAmzServerSideEncryptionCustomerKey       types.String `tfsdk:"x_amz_server_side_encryption_customer_key"`
	XAmzServerSideEncryptionCustomerKeyMD5    types.String `tfsdk:"x_amz_server_side_encryption_customer_key_md5"`
	XAmzServerSideEncryptionContext           types.String `tfsdk:"x_amz_server_side_encryption_context"`
	XAmzRequestPayer                          types.String `tfsdk:"x_amz_request_payer"`
	XAmzTagging                               types.String `tfsdk:"x_amz_tagging"`
	XAmzObjectLockMode                        types.String `tfsdk:"x_amz_object_lock_mode"`
	XAmzObjectLockRetainUntilDate             types.String `tfsdk:"x_amz_object_lock_retain_until_date"`
	XAmzObjectLockLegalHold                   types.String `tfsdk:"x_amz_object_lock_legal_hold"`
}

// Metadata returns the metadata for the bucket policy resource.
func (r *s3BucketObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_object" // todo: use constant here maybe
}

// Schema returns the schema for the bucket policy resource.
func (r *s3BucketObjectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Description: "Name of the S3 bucket to which the object will be added.",
				Required:    true,
			},
			"object_key": schema.StringAttribute{
				Description: "Key name of the object",
				Required:    true,
			},
			"body_file": schema.StringAttribute{
				Description: "The name of the file containing the body.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body.",
				Computed:    true,
			},
			"cache_control": schema.StringAttribute{
				Description: "Can be used to specify caching behavior along the request/reply chain.",
				Optional:    true,
			},
			"content_disposition": schema.StringAttribute{
				Description: "Specifies presentational information for the object.",
				Optional:    true,
			},
			"content_encoding": schema.StringAttribute{
				Description: "Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.",
				Optional:    true,
			},
			"content_language": schema.StringAttribute{
				Description: "The language the content is in.",
				Optional:    true,
			},
			"content_length": schema.Int64Attribute{
				Description: "Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically.",
				Optional:    true,
			},
			"content_md5": schema.StringAttribute{
				Description: "he base64 encoded MD5 digest of the message (without the headers) according to RFC 1864",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "A standard MIME type describing the format of the contents.",
				Optional:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The date and time at which the object is no longer cacheable.",
				Optional:    true,
			},
			"x_amz_server_side_encryption": schema.StringAttribute{
				Description: "The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("AES256"),
				},
			},
			"x_amz_storage_class": schema.StringAttribute{
				Description: "The storage class. The valid value is STANDARD.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("STANDARD"),
				},
			},
			"x_amz_website_redirect_location": schema.StringAttribute{
				Description: "If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata.",
				Optional:    true,
			},
			"x_amz_server_side_encryption_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when encrypting the object. The valid option is AES256.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("AES256"),
				},
			},
			"x_amz_server_side_encryption_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data.",
				Optional:    true,
			},
			"x_amz_server_side_encryption_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error.",
				Optional:    true,
			},
			"x_amz_server_side_encryption_context": schema.StringAttribute{
				Description: "Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs.",
				Optional:    true,
			},
			"x_amz_request_payer": schema.StringAttribute{
				Description: "The value is requester",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("requester"),
				},
			},
			"x_amz_tagging": schema.StringAttribute{
				Description: "The tag-set for the object. The tag-set must be encoded as URL Query parameters.",
				Optional:    true,
			},
			"x_amz_object_lock_mode": schema.StringAttribute{
				Description: "The Object Lock mode that you want to apply to this object.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE"),
				},
			},
			"x_amz_object_lock_retain_until_date": schema.StringAttribute{
				Description: "The date and time when you want this object's Object Lock to expire. Must be formatted as a timestamp parameter.",
				Optional:    true,
			},
			"x_amz_object_lock_legal_hold": schema.StringAttribute{
				Description: "Specifies whether a legal hold will be applied to this object.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF"),
				},
			},
		},
	}
}

// Configure configures the bucket object resource.
func (r *s3BucketObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the bucket policy.
func (r *s3BucketObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured") // todo: const for this error maybe?
		return
	}

	var data s3BucketObjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	f, err := os.Open(data.BodyFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot open file", err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(f, buf)
	if err != nil {
		resp.Diagnostics.AddError("cannot read from file", err.Error())
		return
	}
	request := r.client.ObjectsApi.PutObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Body(f)
	addRequestHeaders(data, &request, resp)
	_, err = request.Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", err.Error())
		return
	}

	data.Body = types.StringValue(buf.String())
	f.Close()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket policy.
func (r *s3BucketObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data s3BucketObjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, _, err := r.client.ObjectsApi.GetObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve object bucket", err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(obj, buf)
	obj.Close()
	if err != nil {
		resp.Diagnostics.AddError("cannot read from file", err.Error())
		return
	}

	data.Body = types.StringValue(buf.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket policy.
func (r *s3BucketObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("s3_bucket_object"), req, resp)
}

// Update updates the bucket policy.
func (r *s3BucketObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured") // todo: const for this error maybe?
		return
	}

	var data s3BucketObjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	f, err := os.Open(data.BodyFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot open file", err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(f, buf)
	if err != nil {
		resp.Diagnostics.AddError("cannot read from file", err.Error())
		return
	}
	request := r.client.ObjectsApi.PutObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Body(f)
	addRequestHeaders(data, &request, resp)

	_, err = request.Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", err.Error())
		return
	}

	data.Body = types.StringValue(buf.String())
	f.Close()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the object.
func (r *s3BucketObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}
	var data s3BucketObjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := r.client.ObjectsApi.DeleteObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete Bucket Object", err.Error())
		return
	}
}

func addRequestHeaders(data s3BucketObjectModel, request *s3.ApiPutObjectRequest, respInt interface{}) {

	var convertedTime time.Time
	var err error

	if !data.CacheControl.IsNull() {
		*request = request.CacheControl(data.CacheControl.ValueString())
	}
	if !data.ContentDisposition.IsNull() {
		*request = request.ContentDisposition(data.ContentDisposition.ValueString())
	}
	if !data.ContentEncoding.IsNull() {
		*request = request.ContentEncoding(data.ContentEncoding.ValueString())
	}
	if !data.ContentLanguage.IsNull() {
		*request = request.ContentLanguage(data.ContentLanguage.ValueString())
	}
	if !data.ContentLength.IsNull() {
		*request = request.ContentLength(int32(data.ContentLength.ValueInt64()))
	}
	if !data.ContentMD5.IsNull() {
		*request = request.ContentMD5(data.ContentMD5.ValueString())
	}
	if !data.ContentType.IsNull() {
		*request = request.ContentType(data.ContentType.ValueString())
	}
	if !data.Expires.IsNull() {
		if convertedTime, err = time.Parse(constant.DatetimeTZOffsetLayout, data.Expires.ValueString()); err != nil {
			if convertedTime, err = time.Parse(constant.DatetimeZLayout, data.Expires.ValueString()); err != nil {
				if resp, ok := respInt.(*resource.CreateResponse); ok {
					resp.Diagnostics.AddError("an error occurred while converting from IonosTime string to time.Time", err.Error())
				} else {
					if resp, ok := respInt.(*resource.UpdateResponse); ok {
						resp.Diagnostics.AddError("an error occurred while converting from IonosTime string to time.Time", err.Error())
					}
				}
				return
			}
		}
		*request = request.Expires(convertedTime)
	}
	if !data.XAmzServerSideEncryption.IsNull() {
		*request = request.XAmzServerSideEncryption(data.XAmzServerSideEncryption.ValueString())
	}
	if !data.XAmzStorageClass.IsNull() {
		*request = request.XAmzStorageClass(data.XAmzStorageClass.ValueString())
	}
	if !data.XAmzWebsiteRedirectLocation.IsNull() {
		*request = request.XAmzWebsiteRedirectLocation(data.XAmzWebsiteRedirectLocation.ValueString())
	}
	if !data.XAmzServerSideEncryptionCustomerAlgorithm.IsNull() {
		*request = request.XAmzServerSideEncryptionCustomerAlgorithm(data.XAmzServerSideEncryptionCustomerAlgorithm.ValueString())
	}
	if !data.XAmzServerSideEncryptionCustomerKey.IsNull() {
		*request = request.XAmzServerSideEncryptionCustomerKey(data.XAmzServerSideEncryptionCustomerKey.ValueString())
	}
	if !data.XAmzServerSideEncryptionCustomerKeyMD5.IsNull() {
		*request = request.XAmzServerSideEncryptionCustomerKeyMD5(data.XAmzServerSideEncryptionCustomerKeyMD5.ValueString())
	}
	if !data.XAmzServerSideEncryptionContext.IsNull() {
		*request = request.XAmzServerSideEncryptionContext(data.XAmzServerSideEncryptionContext.ValueString())
	}
	if !data.XAmzRequestPayer.IsNull() {
		*request = request.XAmzRequestPayer(data.XAmzRequestPayer.ValueString())
	}
	if !data.XAmzTagging.IsNull() {
		*request = request.XAmzTagging(data.XAmzTagging.ValueString())
	}
	if !data.XAmzObjectLockMode.IsNull() {
		*request = request.XAmzObjectLockMode(data.XAmzObjectLockMode.ValueString())
	}
	if !data.XAmzObjectLockRetainUntilDate.IsNull() {
		if convertedTime, err = time.Parse(constant.DatetimeTZOffsetLayout, data.Expires.ValueString()); err != nil {
			if convertedTime, err = time.Parse(constant.DatetimeZLayout, data.Expires.ValueString()); err != nil {
				if resp, ok := respInt.(*resource.CreateResponse); ok {
					resp.Diagnostics.AddError("an error occurred while converting from IonosTime string to time.Time", err.Error())
				} else {
					if resp, ok := respInt.(*resource.UpdateResponse); ok {
						resp.Diagnostics.AddError("an error occurred while converting from IonosTime string to time.Time", err.Error())
					}
				}
				return
			}
		}
		*request = request.XAmzObjectLockRetainUntilDate(convertedTime)
	}
	if !data.XAmzObjectLockLegalHold.IsNull() {
		*request = request.XAmzObjectLockLegalHold(data.XAmzObjectLockLegalHold.ValueString())
	}
}
