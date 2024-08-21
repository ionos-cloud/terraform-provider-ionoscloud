package s3

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"

	tfs3 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"
)

var (
	_ resource.ResourceWithImportState = (*objectCopyResource)(nil)
	_ resource.ResourceWithConfigure   = (*objectCopyResource)(nil)
)

// NewObjectCopyResource creates a new resource for the object copy resource.
func NewObjectCopyResource() resource.Resource {
	return &objectCopyResource{}
}

type objectCopyResource struct {
	client *s3.APIClient
}

type objectCopyResourceModel struct {
	Bucket                                types.String `tfsdk:"bucket"`
	Key                                   types.String `tfsdk:"key"`
	Source                                types.String `tfsdk:"source"`
	CacheControl                          types.String `tfsdk:"cache_control"`
	ContentDisposition                    types.String `tfsdk:"content_disposition"`
	ContentEncoding                       types.String `tfsdk:"content_encoding"`
	ContentLanguage                       types.String `tfsdk:"content_language"`
	ContentType                           types.String `tfsdk:"content_type"`
	CopyIfMatch                           types.String `tfsdk:"copy_if_match"`
	CopyIfModifiedSince                   types.String `tfsdk:"copy_if_modified_since"`
	CopyIfNoneMatch                       types.String `tfsdk:"copy_if_none_match"`
	CopyIfUnmodifiedSince                 types.String `tfsdk:"copy_if_unmodified_since"`
	Expires                               types.String `tfsdk:"expires"`
	MetadataDirective                     types.String `tfsdk:"metadata_directive"`
	TaggingDirective                      types.String `tfsdk:"tagging_directive"`
	ServerSideEncryption                  types.String `tfsdk:"server_side_encryption"`
	StorageClass                          types.String `tfsdk:"storage_class"`
	WebsiteRedirect                       types.String `tfsdk:"website_redirect"`
	ServerSideEncryptionCustomerAlgorithm types.String `tfsdk:"server_side_encryption_customer_algorithm"`
	ServerSideEncryptionCustomerKey       types.String `tfsdk:"server_side_encryption_customer_key"`
	ServerSideEncryptionCustomerKeyMD5    types.String `tfsdk:"server_side_encryption_customer_key_md5"`
	SourceCustomerAlgorithm               types.String `tfsdk:"source_customer_algorithm"`
	SourceCustomerKey                     types.String `tfsdk:"source_customer_key"`
	SourceCustomerKeyMD5                  types.String `tfsdk:"source_customer_key_md5"`
	ObjectLockMode                        types.String `tfsdk:"object_lock_mode"`
	ObjectLockRetainUntilDate             types.String `tfsdk:"object_lock_retain_until_date"`
	ObjectLockLegalHold                   types.String `tfsdk:"object_lock_legal_hold"`
	LastModified                          types.String `tfsdk:"last_modified"`
	Etag                                  types.String `tfsdk:"etag"`
	Metadata                              types.Map    `tfsdk:"metadata"`
	Tags                                  types.Map    `tfsdk:"tags"`
	VersionID                             types.String `tfsdk:"version_id"`
	ForceDestroy                          types.Bool   `tfsdk:"force_destroy"`
}

// Metadata returns the metadata for the object copy resource.
func (r *objectCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_object_copy"
}

// Schema returns the schema for the object copy resource.
func (r *objectCopyResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The key of the object copy",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"source": schema.StringAttribute{
				Description: "The key of the source object",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"copy_if_match": schema.StringAttribute{
				Description: "Copies the object if its entity tag (ETag) matches the specified tag",
				Optional:    true,
			},
			"copy_if_modified_since": schema.StringAttribute{
				Description: "Copies the object if it has been modified since the specified time",
				Optional:    true,
			},
			"copy_if_none_match": schema.StringAttribute{
				Description: "Copies the object if its entity tag (ETag) is different than the specified ETag",
				Optional:    true,
			},
			"copy_if_unmodified_since": schema.StringAttribute{
				Description: "Copies the object if it hasn't been modified since the specified time",
				Optional:    true,
			},
			"cache_control": schema.StringAttribute{
				Description: "Can be used to specify caching behavior along the request/reply chain",
				Optional:    true,
			},
			"content_disposition": schema.StringAttribute{
				Description: "Specifies presentational information for the object copy",
				Optional:    true,
			},
			"content_encoding": schema.StringAttribute{
				Description: "Specifies what content encodings have been applied to the object copy and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field",
				Optional:    true,
			},
			"content_language": schema.StringAttribute{
				Description: "The natural language or languages of the intended audience for the object copy",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "A standard MIME type describing the format of the contents",
				Optional:    true,
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The date and time at which the object copy is no longer cacheable",
				Optional:    true,
			},
			"metadata_directive": schema.StringAttribute{
				Description: "Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("COPY", "REPLACE")},
			},
			"tagging_directive": schema.StringAttribute{
				Description: "Specifies whether the object copy tag-set is copied from the source object or replaced with tag-set provided in the request",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("COPY", "REPLACE")},
			},
			"server_side_encryption": schema.StringAttribute{
				Description: "The server-side encryption algorithm used when storing this object copy in IONOS S3 Object Copy Storage (AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"storage_class": schema.StringAttribute{
				Description: "The storage class of the object copy. Valid value is 'STANDARD'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Validators:  []validator.String{stringvalidator.OneOf("STANDARD")},
			},
			"website_redirect": schema.StringAttribute{
				Description: "If the bucket is configured as a website, redirects requests for this object copy to another object copy in the same bucket or to an external URL. IONOS S3 Object Copy Storage stores the value of this header in the object copy metadata",
				Optional:    true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when encrypting the object copy (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data",
				Optional:    true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Copy Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"source_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when decrypting the source object (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"source_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to decrypt the source object",
				Optional:    true,
			},
			"source_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Copy Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"object_lock_mode": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Description: " The date and time when you want this object copy's Object Copy Lock to expire. Must be formatted as a timestamp parameter.",
				Optional:    true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Description: "Specifies whether a legal hold will be applied to this object copy.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Description: "An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.",
				Computed:    true,
			},
			"last_modified": schema.StringAttribute{
				Description: "The date and time at which the object copy was last modified",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "The tag-set for the object copy",
				Optional:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Description: "A map of metadata to store with the object copy in IONOS S3 Object Copy Storage",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{mapvalidator.ValueStringsAre([]validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9]+$`), "metadata keys must be lowercase alphanumeric characters"),
				}...)},
			},
			"version_id": schema.StringAttribute{
				Description: "The version of the object copy",
				Computed:    true,
			},
			"force_destroy": schema.BoolAttribute{
				Description: "Specifies whether to delete the object copy even if it has a governance-type Object Copy Lock in place. You must explicitly pass a value of true for this parameter to delete the object copy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Configure configures the object copy resource.
func (r *objectCopyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the object copy.
func (r *objectCopyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *objectCopyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := copyObject(ctx, r.client, data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create object copy", formatXMLError(err).Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the object copy.
func (r *objectCopyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *objectCopyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := findObjectCopy(ctx, r.client, data)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to read object copy", formatXMLError(err).Error())
		return
	}

	diags := r.setDataModel(ctx, data, apiResponse)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the object copy.
func (r *objectCopyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	bucket, key, err := splitImportID(id)
	if err != nil {
		resp.Diagnostics.AddError("invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), bucket)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}

func hasCopyConditions(plan *objectCopyResourceModel) bool {
	return !plan.CopyIfMatch.IsNull() || !plan.CopyIfModifiedSince.IsNull() || !plan.CopyIfNoneMatch.IsNull() || !plan.CopyIfUnmodifiedSince.IsNull()
}

// Update updates the object copy.
func (r *objectCopyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *objectCopyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasCopyConditions(plan) || hasObjectCopyContentChanges(plan, state) {
		err := copyObject(ctx, r.client, plan)
		if err != nil {
			resp.Diagnostics.AddError("failed to update object copy", formatXMLError(err).Error())
			return
		}
	}

	setObjectCopyStateForUnknown(plan, state)
	diags := r.refreshData(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the object copy.
func (r *objectCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data objectCopyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.VersionID.IsNull() {
		if _, err := tfs3.DeleteAllObjectVersions(ctx, r.client, &tfs3.DeleteRequest{
			Bucket:       data.Bucket.ValueString(),
			Key:          data.Key.ValueString(),
			ForceDestroy: data.ForceDestroy.ValueBool(),
			VersionID:    data.VersionID.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("failed to delete object copy versions", formatXMLError(err).Error())
			return
		}
	} else {
		_, apiResponse, err := deleteObjectCopy(ctx, r.client, &data)
		if err != nil {
			if apiResponse.HttpNotFound() {
				return
			}

			resp.Diagnostics.AddError("failed to delete object copy", formatXMLError(err).Error())
			return
		}
	}
}

func copyObject(ctx context.Context, client *s3.APIClient, data *objectCopyResourceModel) error {
	req := client.ObjectsApi.CopyObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	err := fillObjectCopyRequest(&req, data)
	if err != nil {
		return err
	}

	output, apiResponse, err := req.Execute()
	if err != nil {
		return err
	}

	return setObjectCopyComputedAttributes(ctx, data, apiResponse, output, client)
}

func setObjectCopyComputedAttributes(ctx context.Context, data *objectCopyResourceModel, apiResponse *s3.APIResponse, output *s3.CopyObjectResult, client *s3.APIClient) error {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	data.VersionID = types.StringValue(apiResponse.Header.Get("x-amz-version-id"))

	if output.ETag != nil {
		data.Etag = types.StringValue(strings.Trim(*output.ETag, "\""))
	}

	data.LastModified = types.StringValue(output.LastModified.Format(time.RFC3339))

	return setObjectCopyContentType(ctx, data, client)
}

func (r *objectCopyResource) refreshData(ctx context.Context, data *objectCopyResourceModel) diag.Diagnostics {
	diags := diag.Diagnostics{}
	_, apiResponse, err := findObjectCopy(ctx, r.client, data)
	if err != nil {
		diags.AddError("failed to read object copy", formatXMLError(err).Error())
		return diags
	}

	diags = r.setDataModel(ctx, data, apiResponse)
	if diags.HasError() {
		diags.Append(diags...)
		return diags
	}

	return nil
}

func setObjectCopyStateForUnknown(plan, state *objectCopyResourceModel) {
	if plan.VersionID.IsUnknown() {
		plan.VersionID = state.VersionID
	}

	if plan.Etag.IsUnknown() {
		plan.Etag = state.Etag
	}
}

func setObjectCopyContentData(data *objectCopyResourceModel, apiResponse *s3.APIResponse) {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	data.VersionID = types.StringValue(apiResponse.Header.Get("x-amz-version-id"))

	etag := apiResponse.Header.Get("ETag")
	if etag != "" {
		data.Etag = types.StringValue(strings.Trim(etag, "\""))
	}

	cacheControl := apiResponse.Header.Get("Cache-Control")
	if cacheControl != "" {
		data.CacheControl = types.StringValue(cacheControl)
	}

	contentDisposition := apiResponse.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		data.ContentDisposition = types.StringValue(contentDisposition)
	}

	contentEncoding := apiResponse.Header.Get("Content-Encoding")
	if contentEncoding != "" {
		data.ContentEncoding = types.StringValue(contentEncoding)
	}

	contentLanguage := apiResponse.Header.Get("Content-Language")
	if contentLanguage != "" {
		data.ContentLanguage = types.StringValue(contentLanguage)
	}

	expires := apiResponse.Header.Get("Expires")
	if expires != "" {
		data.Expires = types.StringValue(expires)
	}
}

func setObjectCopyServerSideEncryptionData(data *objectCopyResourceModel, apiResponse *s3.APIResponse) {
	serverSideEncryption := apiResponse.Header.Get("x-amz-server-side-encryption")
	if serverSideEncryption != "" {
		data.ServerSideEncryption = types.StringValue(serverSideEncryption)
	}

	serverSideEncryptionCustomerAlgorithm := apiResponse.Header.Get("x-amz-server-side-encryption-customer-algorithm")
	if serverSideEncryptionCustomerAlgorithm != "" {
		data.ServerSideEncryptionCustomerAlgorithm = types.StringValue(serverSideEncryptionCustomerAlgorithm)
	}

	serverSideEncryptionCustomerKey := apiResponse.Header.Get("x-amz-server-side-encryption-customer-key")
	if serverSideEncryptionCustomerKey != "" {
		data.ServerSideEncryptionCustomerKey = types.StringValue(serverSideEncryptionCustomerKey)
	}

	serverSideEncryptionCustomerKeyMD5 := apiResponse.Header.Get("x-amz-server-side-encryption-customer-key-MD5")
	if serverSideEncryptionCustomerKeyMD5 != "" {
		data.ServerSideEncryptionCustomerKeyMD5 = types.StringValue(serverSideEncryptionCustomerKeyMD5)
	}
}

func setObjectCopyLockData(data *objectCopyResourceModel, apiResponse *s3.APIResponse) error {
	objectCopyLockMode := apiResponse.Header.Get("x-amz-object copy-lock-mode")
	if objectCopyLockMode != "" {
		data.ObjectLockMode = types.StringValue(objectCopyLockMode)
	}

	objectCopyLockRetainUntilDate := apiResponse.Header.Get("x-amz-object copy-lock-retain-until-date")
	if objectCopyLockRetainUntilDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, objectCopyLockRetainUntilDate)
		if err != nil {
			return fmt.Errorf("failed to parse object copy lock retain until date: %w", err)
		}

		data.ObjectLockRetainUntilDate = types.StringValue(parsedTime.Format(time.RFC3339))
	}

	objectCopyLockLegalHold := apiResponse.Header.Get("x-amz-object copy-lock-legal-hold")
	if objectCopyLockLegalHold != "" {
		data.ObjectLockLegalHold = types.StringValue(objectCopyLockLegalHold)
	}

	return nil
}

func (r *objectCopyResource) setTagsData(ctx context.Context, data *objectCopyResourceModel) diag.Diagnostics {
	tagsMap, err := getTags(ctx, r.client, data.Bucket.ValueString(), data.Key.ValueString())
	if err != nil {
		diags := diag.Diagnostics{}
		diags.AddError("failed to get tags", err.Error())
		return diags
	}

	if len(tagsMap) > 0 {
		tags, diagErr := types.MapValueFrom(ctx, types.StringType, tagsMap)
		if diagErr.HasError() {
			return diagErr
		}
		data.Tags = tags
	}

	return nil
}

func setObjectCopyMetadata(ctx context.Context, data *objectCopyResourceModel, apiResponse *s3.APIResponse) diag.Diagnostics {
	metadataMap := getMetadataMapFromHeaders(apiResponse, "X-Amz-Meta-")

	if len(metadataMap) > 0 {
		metadata, diagErr := types.MapValueFrom(ctx, types.StringType, metadataMap)
		if diagErr.HasError() {
			return diagErr
		}
		data.Metadata = metadata
	}

	return nil
}

func (r *objectCopyResource) setDataModel(ctx context.Context, data *objectCopyResourceModel, apiResponse *s3.APIResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	setObjectCopyContentData(data, apiResponse)
	setObjectCopyServerSideEncryptionData(data, apiResponse)

	if err := setObjectCopyLockData(data, apiResponse); err != nil {
		diags.AddError("failed to set object copy lock data", err.Error())
		return diags
	}

	storageClass := apiResponse.Header.Get("x-amz-storage-class")
	if storageClass != "" {
		data.StorageClass = types.StringValue(storageClass)
	}

	websiteRedirect := apiResponse.Header.Get("x-amz-website-redirect-location")
	if websiteRedirect != "" {
		data.WebsiteRedirect = types.StringValue(websiteRedirect)
	}

	if diags = setObjectCopyMetadata(ctx, data, apiResponse); diags.HasError() {
		return diags
	}

	if diags = r.setTagsData(ctx, data); diags.HasError() {
		return diags
	}

	return nil
}

func setObjectCopyContentType(ctx context.Context, data *objectCopyResourceModel, client *s3.APIClient) error {
	_, apiResponse, err := findObjectCopy(ctx, client, data)
	if err != nil {
		return err
	}

	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	return nil
}

func deleteObjectCopy(ctx context.Context, client *s3.APIClient, data *objectCopyResourceModel) (map[string]interface{}, *s3.APIResponse, error) {
	req := client.ObjectsApi.DeleteObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if !data.VersionID.IsNull() {
		req = req.VersionId(data.VersionID.ValueString())
	}

	if !data.ForceDestroy.IsNull() {
		req = req.XAmzBypassGovernanceRetention(data.ForceDestroy.ValueBool())
	}

	return req.Execute()
}

func findObjectCopy(ctx context.Context, client *s3.APIClient, data *objectCopyResourceModel) (*s3.HeadObjectOutput, *s3.APIResponse, error) {
	req := client.ObjectsApi.HeadObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if !data.Etag.IsNull() {
		req = req.IfMatch(data.Etag.ValueString())
	}

	if !data.VersionID.IsNull() {
		req = req.VersionId(data.VersionID.ValueString())
	}

	if !data.ServerSideEncryptionCustomerAlgorithm.IsNull() {
		req = req.XAmzServerSideEncryptionCustomerAlgorithm(data.ServerSideEncryptionCustomerAlgorithm.ValueString())
	}

	if !data.ServerSideEncryptionCustomerKey.IsNull() {
		req = req.XAmzServerSideEncryptionCustomerKey(data.ServerSideEncryptionCustomerKey.ValueString())
	}

	if !data.ServerSideEncryptionCustomerKeyMD5.IsNull() {
		req = req.XAmzServerSideEncryptionCustomerKeyMD5(data.ServerSideEncryptionCustomerKeyMD5.ValueString())
	}

	return req.Execute()
}

func fillObjectCopyContentData(data *objectCopyResourceModel, req *s3.ApiCopyObjectRequest) error {
	if !data.CacheControl.IsNull() {
		*req = req.CacheControl(data.CacheControl.ValueString())
	}

	if !data.ContentDisposition.IsNull() {
		*req = req.ContentDisposition(data.ContentDisposition.ValueString())
	}

	if !data.ContentEncoding.IsNull() {
		*req = req.ContentEncoding(data.ContentEncoding.ValueString())
	}

	if !data.ContentLanguage.IsNull() {
		*req = req.ContentLanguage(data.ContentLanguage.ValueString())
	}

	if !data.ContentType.IsNull() {
		*req = req.ContentType(data.ContentType.ValueString())
	}

	if !data.Expires.IsNull() {
		t, err := time.Parse(time.RFC3339, data.Expires.ValueString())
		if err != nil {
			return fmt.Errorf("failed to parse expires time: %s", err.Error())
		}

		*req = req.Expires(t)
	}

	if !data.Source.IsNull() {
		*req = req.XAmzCopySource(data.Source.ValueString())
	}

	if !data.CopyIfMatch.IsNull() {
		*req = req.XAmzCopySourceIfMatch(data.CopyIfMatch.ValueString())
	}

	if !data.CopyIfModifiedSince.IsNull() {
		t, err := time.Parse(time.RFC3339, data.CopyIfModifiedSince.ValueString())
		if err != nil {
			return fmt.Errorf("failed to parse copy_if_modified_since time: %s", err.Error())
		}

		*req = req.XAmzCopySourceIfModifiedSince(t)
	}

	if !data.CopyIfNoneMatch.IsNull() {
		*req = req.XAmzCopySourceIfNoneMatch(data.CopyIfNoneMatch.ValueString())
	}

	if !data.CopyIfUnmodifiedSince.IsNull() {
		t, err := time.Parse(time.RFC3339, data.CopyIfUnmodifiedSince.ValueString())
		if err != nil {
			return fmt.Errorf("failed to parse copy_if_unmodified_since time: %s", err.Error())
		}

		*req = req.XAmzCopySourceIfUnmodifiedSince(t)
	}

	return nil
}

func fillObjectCopyServerSideEncryptionData(data *objectCopyResourceModel, req *s3.ApiCopyObjectRequest) {
	if !data.ServerSideEncryption.IsNull() {
		*req = req.XAmzServerSideEncryption(data.ServerSideEncryption.ValueString())
	}

	if !data.ServerSideEncryptionCustomerAlgorithm.IsNull() {
		*req = req.XAmzServerSideEncryptionCustomerAlgorithm(data.ServerSideEncryptionCustomerAlgorithm.ValueString())
	}

	if !data.ServerSideEncryptionCustomerKey.IsNull() {
		*req = req.XAmzServerSideEncryptionCustomerKey(data.ServerSideEncryptionCustomerKey.ValueString())
	}

	if !data.ServerSideEncryptionCustomerKeyMD5.IsNull() {
		*req = req.XAmzServerSideEncryptionCustomerKeyMD5(data.ServerSideEncryptionCustomerKeyMD5.ValueString())
	}

	if !data.SourceCustomerAlgorithm.IsNull() {
		*req = req.XAmzCopySourceServerSideEncryptionCustomerAlgorithm(data.SourceCustomerAlgorithm.ValueString())
	}

	if !data.SourceCustomerKey.IsNull() {
		*req = req.XAmzCopySourceServerSideEncryptionCustomerKey(data.SourceCustomerKey.ValueString())
	}

	if !data.SourceCustomerKeyMD5.IsNull() {
		*req = req.XAmzCopySourceServerSideEncryptionCustomerKeyMD5(data.SourceCustomerKeyMD5.ValueString())
	}
}

func fillObjectCopyLockData(data *objectCopyResourceModel, req *s3.ApiCopyObjectRequest) error {
	if !data.ObjectLockMode.IsNull() {
		*req = req.XAmzObjectLockMode(data.ObjectLockMode.ValueString())
	}

	if !data.ObjectLockRetainUntilDate.IsNull() {
		t, err := time.Parse(time.RFC3339, data.ObjectLockRetainUntilDate.ValueString())
		if err != nil {
			return fmt.Errorf("can't parse objectCopy_lock_retain_until_date: %w", err)
		}

		*req = req.XAmzObjectLockRetainUntilDate(t)
	}

	if !data.ObjectLockLegalHold.IsNull() {
		*req = req.XAmzObjectLockLegalHold(data.ObjectLockLegalHold.ValueString())
	}

	return nil
}

func fillObjectCopyRequest(req *s3.ApiCopyObjectRequest, data *objectCopyResourceModel) error {
	fillObjectCopyServerSideEncryptionData(data, req)
	if err := fillObjectCopyContentData(data, req); err != nil {
		return err
	}

	if err := fillObjectCopyLockData(data, req); err != nil {
		return err
	}

	if !data.StorageClass.IsNull() {
		*req = req.XAmzStorageClass(data.StorageClass.ValueString())
	}

	if !data.WebsiteRedirect.IsNull() {
		*req = req.XAmzWebsiteRedirectLocation(data.WebsiteRedirect.ValueString())
	}

	if !data.TaggingDirective.IsNull() {
		*req = req.XAmzTaggingDirective(data.TaggingDirective.ValueString())
	}

	if !data.Tags.IsNull() {
		tags, err := buildQueryString(data.Tags)
		if err != nil {
			return fmt.Errorf("failed to build tags query string: %s", err.Error())
		}
		*req = req.XAmzTagging(tags)
	}

	if !data.MetadataDirective.IsNull() {
		*req = req.XAmzMetadataDirective(data.MetadataDirective.ValueString())
	}

	if !data.Metadata.IsNull() {
		metadata, err := fromTFMap(data.Metadata)
		if err != nil {
			return fmt.Errorf("failed to convert metadata: %s", err.Error())
		}

		*req = req.XAmzMeta(metadata)
	}

	return nil
}

// hasObjectCopyContentChanges returns true if the plan has changes to the object copy content.
func hasObjectCopyContentChanges(plan, state *objectCopyResourceModel) bool {
	needsChange := !(plan.Source.Equal(state.Source) &&
		plan.CacheControl.Equal(state.CacheControl) &&
		plan.ContentDisposition.Equal(state.ContentDisposition) &&
		plan.ContentEncoding.Equal(state.ContentEncoding) &&
		plan.ContentLanguage.Equal(state.ContentLanguage) &&
		plan.ContentType.Equal(state.ContentType) &&
		plan.Expires.Equal(state.Expires) &&
		plan.ServerSideEncryption.Equal(state.ServerSideEncryption) &&
		plan.StorageClass.Equal(state.StorageClass) &&
		plan.WebsiteRedirect.Equal(state.WebsiteRedirect) &&
		plan.ServerSideEncryptionCustomerAlgorithm.Equal(state.ServerSideEncryptionCustomerAlgorithm) &&
		plan.ServerSideEncryptionCustomerKey.Equal(state.ServerSideEncryptionCustomerKey) &&
		plan.ServerSideEncryptionCustomerKeyMD5.Equal(state.ServerSideEncryptionCustomerKeyMD5) &&
		plan.SourceCustomerAlgorithm.Equal(state.SourceCustomerAlgorithm) &&
		plan.SourceCustomerKey.Equal(state.SourceCustomerKey) &&
		plan.SourceCustomerKeyMD5.Equal(state.SourceCustomerKeyMD5) &&
		plan.MetadataDirective.Equal(state.MetadataDirective) &&
		plan.TaggingDirective.Equal(state.TaggingDirective) &&
		plan.Metadata.Equal(state.Metadata) &&
		plan.Source.Equal(state.Source))
	return needsChange
}
