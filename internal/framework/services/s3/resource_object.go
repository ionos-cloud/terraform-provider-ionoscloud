package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
	"github.com/mitchellh/go-homedir"

	tfs3 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"
)

var (
	_ resource.ResourceWithImportState      = (*objectResource)(nil)
	_ resource.ResourceWithConfigure        = (*objectResource)(nil)
	_ resource.ResourceWithConfigValidators = (*objectResource)(nil)
)

// NewObjectResource creates a new resource for the object resource.
func NewObjectResource() resource.Resource {
	return &objectResource{}
}

type objectResource struct {
	client *s3.APIClient
}

type objectResourceModel struct {
	Bucket                                types.String `tfsdk:"bucket"`
	Key                                   types.String `tfsdk:"key"`
	Source                                types.String `tfsdk:"source"`
	CacheControl                          types.String `tfsdk:"cache_control"`
	Content                               types.String `tfsdk:"content"`
	ContentDisposition                    types.String `tfsdk:"content_disposition"`
	ContentEncoding                       types.String `tfsdk:"content_encoding"`
	ContentLanguage                       types.String `tfsdk:"content_language"`
	ContentType                           types.String `tfsdk:"content_type"`
	Expires                               types.String `tfsdk:"expires"`
	ServerSideEncryption                  types.String `tfsdk:"server_side_encryption"`
	StorageClass                          types.String `tfsdk:"storage_class"`
	WebsiteRedirect                       types.String `tfsdk:"website_redirect"`
	ServerSideEncryptionCustomerAlgorithm types.String `tfsdk:"server_side_encryption_customer_algorithm"`
	ServerSideEncryptionCustomerKey       types.String `tfsdk:"server_side_encryption_customer_key"`
	ServerSideEncryptionCustomerKeyMD5    types.String `tfsdk:"server_side_encryption_customer_key_md5"`
	ServerSideEncryptionContext           types.String `tfsdk:"server_side_encryption_context"`
	RequestPayer                          types.String `tfsdk:"request_payer"`
	ObjectLockMode                        types.String `tfsdk:"object_lock_mode"`
	ObjectLockRetainUntilDate             types.String `tfsdk:"object_lock_retain_until_date"`
	ObjectLockLegalHold                   types.String `tfsdk:"object_lock_legal_hold"`
	Etag                                  types.String `tfsdk:"etag"`
	Metadata                              types.Map    `tfsdk:"metadata"`
	Tags                                  types.Map    `tfsdk:"tags"`
	VersionID                             types.String `tfsdk:"version_id"`
	MFA                                   types.String `tfsdk:"mfa"`
	ForceDestroy                          types.Bool   `tfsdk:"force_destroy"`
}

// Metadata returns the metadata for the object resource.
func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_object"
}

// Schema returns the schema for the object resource.
func (r *objectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "The key of the object",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"source": schema.StringAttribute{
				Description: "The path to the file to upload",
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Description: "The utf-8 content of the object",
				Optional:    true,
			},
			"cache_control": schema.StringAttribute{
				Description: "Can be used to specify caching behavior along the request/reply chain",
				Optional:    true,
			},
			"content_disposition": schema.StringAttribute{
				Description: "Specifies presentational information for the object",
				Optional:    true,
			},
			"content_encoding": schema.StringAttribute{
				Description: "Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field",
				Optional:    true,
			},
			"content_language": schema.StringAttribute{
				Description: "The natural language or languages of the intended audience for the object",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "A standard MIME type describing the format of the contents",
				Optional:    true,
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The date and time at which the object is no longer cacheable",
				Optional:    true,
			},
			"server_side_encryption": schema.StringAttribute{
				Description: "The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"storage_class": schema.StringAttribute{
				Description: "The storage class of the object. Valid value is 'STANDARD'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Validators:  []validator.String{stringvalidator.OneOf("STANDARD")},
			},
			"website_redirect": schema.StringAttribute{
				Description: "If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata",
				Optional:    true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when encrypting the object (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data",
				Optional:    true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"server_side_encryption_context": schema.StringAttribute{
				Description: " Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs.",
				Optional:    true,
				Sensitive:   true,
			},
			"request_payer": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
			},
			"object_lock_mode": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Description: " The date and time when you want this object's Object Lock to expire. Must be formatted as a timestamp parameter.",
				Optional:    true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Description: "Specifies whether a legal hold will be applied to this object.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Description: "An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "The tag-set for the object",
				Optional:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Description: "A map of metadata to store with the object in IONOS S3 Object Storage",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{mapvalidator.ValueStringsAre([]validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9]+$`), "metadata keys must be lowercase alphanumeric characters"),
				}...)},
			},
			"version_id": schema.StringAttribute{
				Description: "The version of the object",
				Computed:    true,
			},
			"mfa": schema.StringAttribute{
				Description: "The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled.",
				Optional:    true,
			},
			"force_destroy": schema.BoolAttribute{
				Description: "Specifies whether to delete the object even if it has a governance-type Object Lock in place. You must explicitly pass a value of true for this parameter to delete the object.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Configure configures the object resource.
func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *objectResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("source"),
			path.MatchRoot("content"),
		),
	}
}

// Create creates the object.
func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *objectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := uploadObject(ctx, r.client, data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create object", formatXMLError(err).Error())
		return
	}

	if err = setComputedAttributes(ctx, data, apiResponse, r.client); err != nil {
		resp.Diagnostics.AddError("failed to set computed attributes", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setComputedAttributes(ctx context.Context, data *objectResourceModel, apiResponse *s3.APIResponse, client *s3.APIClient) error {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	data.VersionID = types.StringValue(apiResponse.Header.Get("x-amz-version-id"))

	etag := apiResponse.Header.Get("ETag")
	if etag != "" {
		data.Etag = types.StringValue(strings.Trim(etag, "\""))
	}

	return setContentType(ctx, data, client)
}

func (r *objectResource) refreshData(ctx context.Context, data *objectResourceModel) diag.Diagnostics {
	diags := diag.Diagnostics{}
	_, apiResponse, err := findObject(ctx, r.client, data)
	if err != nil {
		diags.AddError("failed to read object", formatXMLError(err).Error())
		return diags
	}

	diags = r.setDataModel(ctx, data, apiResponse)
	if diags.HasError() {
		diags.Append(diags...)
		return diags
	}

	return nil
}

// Read reads the object.
func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *objectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := findObject(ctx, r.client, data)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to read object", formatXMLError(err).Error())
		return
	}

	diags := r.setDataModel(ctx, data, apiResponse)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the object.
func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	bucket, key, err := splitImportID(id)
	if err != nil {
		resp.Diagnostics.AddError("invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), bucket)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}

// Update updates the object.
func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *objectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasObjectContentChanges(plan, state) {
		apiResponse, err := uploadObject(ctx, r.client, plan)
		if err != nil {
			resp.Diagnostics.AddError("failed to update object", formatXMLError(err).Error())
			return
		}

		if err = setComputedAttributes(ctx, plan, apiResponse, r.client); err != nil {
			resp.Diagnostics.AddError("failed to set computed attributes", err.Error())
			return
		}

	} else {
		if !plan.ObjectLockLegalHold.Equal(state.ObjectLockLegalHold) {
			_, err := r.client.ObjectLockApi.PutObjectLegalHold(ctx, state.Bucket.ValueString(), state.Key.ValueString()).
				ObjectLegalHoldConfiguration(s3.ObjectLegalHoldConfiguration{Status: plan.ObjectLockLegalHold.ValueStringPointer()}).
				Execute()
			if err != nil {
				resp.Diagnostics.AddError("failed to update object lock legal hold", formatXMLError(err).Error())
				return
			}
		}

		if !plan.ObjectLockMode.Equal(state.ObjectLockMode) || !plan.ObjectLockRetainUntilDate.Equal(state.ObjectLockRetainUntilDate) {
			if _, err := putRetention(ctx, r.client, plan, state); err != nil {
				resp.Diagnostics.AddError("failed to update object lock retention", formatXMLError(err).Error())
				return
			}
		}
	}

	setStateForUnknown(plan, state)
	diags := r.refreshData(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func setStateForUnknown(plan, state *objectResourceModel) {
	if plan.VersionID.IsUnknown() {
		plan.VersionID = state.VersionID
	}

	if plan.Etag.IsUnknown() {
		plan.Etag = state.Etag
	}
}

func getRetentionDate(d types.String) (string, error) {
	if d.IsNull() {
		return time.Time{}.UTC().Format(time.RFC3339), nil
	}

	t, err := time.Parse(time.RFC3339, d.ValueString())
	if err != nil {
		return "", err
	}

	return t.UTC().Format(time.RFC3339), nil
}

func putRetention(ctx context.Context, client *s3.APIClient, plan, state *objectResourceModel) (*s3.APIResponse, error) {
	retentionDate, err := getRetentionDate(plan.ObjectLockRetainUntilDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object lock retain until date: %w", err)
	}

	baseReq := client.ObjectLockApi.PutObjectRetention(ctx, state.Bucket.ValueString(), state.Key.ValueString()).
		PutObjectRetentionRequest(s3.PutObjectRetentionRequest{
			Mode:            plan.ObjectLockMode.ValueStringPointer(),
			RetainUntilDate: &retentionDate,
		})

	if plan.ObjectLockRetainUntilDate != state.ObjectLockRetainUntilDate {
		oldDate := expandObjectDate(state.ObjectLockRetainUntilDate.ValueString())
		newDate := expandObjectDate(plan.ObjectLockRetainUntilDate.ValueString())

		if plan.ObjectLockRetainUntilDate.IsNull() ||
			(!state.ObjectLockRetainUntilDate.IsNull() && newDate.Time.Before(oldDate.Time)) {
			baseReq = baseReq.XAmzBypassGovernanceRetention(true)
		}
	}

	return baseReq.Execute()
}

func expandObjectDate(v string) *s3.IonosTime {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return &s3.IonosTime{Time: time.Time{}.UTC()}
	}

	return &s3.IonosTime{Time: t.UTC()}
}

// Delete deletes the object.
func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data objectResourceModel
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
			resp.Diagnostics.AddError("failed to delete object versions", formatXMLError(err).Error())
			return
		}
	} else {
		_, apiResponse, err := deleteObject(ctx, r.client, &data)
		if err != nil {
			if apiResponse.HttpNotFound() {
				return
			}

			resp.Diagnostics.AddError("failed to delete object", formatXMLError(err).Error())
			return
		}
	}
}

func uploadObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (*s3.APIResponse, error) {
	putReq := client.ObjectsApi.PutObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	err := fillPutObjectRequest(&putReq, data)
	if err != nil {
		return nil, err
	}

	body, err := getBody(data)
	if err != nil {
		return nil, err
	}

	if needsMD5Header(data) {
		if err = addMD5Header(&putReq, body); err != nil {
			return nil, fmt.Errorf("failed to add MD5 header: %w", err)
		}
	}

	defer func() {
		// Remove temp file if content is provided
		if !data.Content.IsNull() {
			err = os.Remove(body.Name())
			if err != nil {
				log.Printf("failed to remove temp file: %s", err.Error())
			}
		}
		// Close the file
		err = body.Close()
		if err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()

	return putReq.Body(body).Execute()
}

func setContentData(data *objectResourceModel, apiResponse *s3.APIResponse) {
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

func setServerSideEncryptionData(data *objectResourceModel, apiResponse *s3.APIResponse) {
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

	serverSideEncryptionContext := apiResponse.Header.Get("x-amz-server-side-encryption-context")
	if serverSideEncryptionContext != "" {
		data.ServerSideEncryptionContext = types.StringValue(serverSideEncryptionContext)
	}
}

func setObjectLockData(data *objectResourceModel, apiResponse *s3.APIResponse) error {
	objectLockMode := apiResponse.Header.Get("x-amz-object-lock-mode")
	if objectLockMode != "" {
		data.ObjectLockMode = types.StringValue(objectLockMode)
	}

	objectLockRetainUntilDate := apiResponse.Header.Get("x-amz-object-lock-retain-until-date")
	if objectLockRetainUntilDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, objectLockRetainUntilDate)
		if err != nil {
			return fmt.Errorf("failed to parse object lock retain until date: %w", err)
		}

		data.ObjectLockRetainUntilDate = types.StringValue(parsedTime.Format(time.RFC3339))
	}

	objectLockLegalHold := apiResponse.Header.Get("x-amz-object-lock-legal-hold")
	if objectLockLegalHold != "" {
		data.ObjectLockLegalHold = types.StringValue(objectLockLegalHold)
	}

	return nil
}

func (r *objectResource) setTagsData(ctx context.Context, data *objectResourceModel) diag.Diagnostics {
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

func setMetadata(ctx context.Context, data *objectResourceModel, apiResponse *s3.APIResponse) diag.Diagnostics {
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

func (r *objectResource) setDataModel(ctx context.Context, data *objectResourceModel, apiResponse *s3.APIResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	setContentData(data, apiResponse)
	setServerSideEncryptionData(data, apiResponse)

	if err := setObjectLockData(data, apiResponse); err != nil {
		diags.AddError("failed to set object lock data", err.Error())
		return diags
	}

	requestPayer := apiResponse.Header.Get("x-amz-request-payer")
	if requestPayer != "" {
		data.RequestPayer = types.StringValue(requestPayer)
	}

	storageClass := apiResponse.Header.Get("x-amz-storage-class")
	if storageClass != "" {
		data.StorageClass = types.StringValue(storageClass)
	}

	websiteRedirect := apiResponse.Header.Get("x-amz-website-redirect-location")
	if websiteRedirect != "" {
		data.WebsiteRedirect = types.StringValue(websiteRedirect)
	}

	if diags = setMetadata(ctx, data, apiResponse); diags.HasError() {
		return diags
	}

	if diags = r.setTagsData(ctx, data); diags.HasError() {
		return diags
	}

	return nil
}

func getTags(ctx context.Context, client *s3.APIClient, bucket, key string) (map[string]string, error) {
	output, apiResponse, err := client.TaggingApi.GetObjectTagging(ctx, bucket, key).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil, nil
		}

		return nil, err
	}

	if output.TagSet == nil {
		return nil, nil
	}

	tagsMap := map[string]string{}
	for _, tag := range *output.TagSet {
		if tag.Key != nil && tag.Value != nil {
			tagsMap[*tag.Key] = *tag.Value
		}
	}

	return tagsMap, nil
}

func getMetadataMapFromHeaders(apiResponse *s3.APIResponse, prefix string) map[string]string {
	metaHeaders := map[string]string{}
	for name, values := range apiResponse.Header {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			if len(values) > 0 {
				metaKey := strings.TrimPrefix(strings.ToLower(name), strings.ToLower(prefix))
				metaHeaders[metaKey] = values[0]
			}
		}
	}

	return metaHeaders
}

func setContentType(ctx context.Context, data *objectResourceModel, client *s3.APIClient) error {
	_, apiResponse, err := findObject(ctx, client, data)
	if err != nil {
		return err
	}

	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	return nil
}

func deleteObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (map[string]interface{}, *s3.APIResponse, error) {
	req := client.ObjectsApi.DeleteObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if !data.VersionID.IsNull() {
		req = req.VersionId(data.VersionID.ValueString())
	}

	if !data.MFA.IsNull() {
		req = req.XAmzMfa(data.MFA.ValueString())
	}

	if !data.ForceDestroy.IsNull() {
		req = req.XAmzBypassGovernanceRetention(data.ForceDestroy.ValueBool())
	}

	return req.Execute()
}

func findObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (*s3.HeadObjectOutput, *s3.APIResponse, error) {
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

func getBody(data *objectResourceModel) (*os.File, error) {
	if !data.Source.IsNull() {
		filePath, err := homedir.Expand(data.Source.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to expand source file path: %s", err.Error())
		}

		file, err := os.Open(filepath.Clean(filePath))
		if err != nil {
			return nil, fmt.Errorf("failed to open source file: %s", err.Error())
		}

		return file, nil
	}

	if !data.Content.IsNull() {
		tempFile, err := createTempFile("temp", data.Content.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %s", err.Error())
		}

		return tempFile, nil
	}

	return nil, nil
}

func fillContentData(data *objectResourceModel, req *s3.ApiPutObjectRequest) error {
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

	return nil
}

func fillServerSideEncryptionData(data *objectResourceModel, req *s3.ApiPutObjectRequest) {
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

	if !data.ServerSideEncryptionContext.IsNull() {
		*req = req.XAmzServerSideEncryptionContext(data.ServerSideEncryptionContext.ValueString())
	}
}

func fillObjectLockData(data *objectResourceModel, req *s3.ApiPutObjectRequest) error {
	if !data.ObjectLockMode.IsNull() {
		*req = req.XAmzObjectLockMode(data.ObjectLockMode.ValueString())
	}

	if !data.ObjectLockRetainUntilDate.IsNull() {
		t, err := time.Parse(time.RFC3339, data.ObjectLockRetainUntilDate.ValueString())
		if err != nil {
			return fmt.Errorf("can't parse object_lock_retain_until_date: %w", err)
		}

		*req = req.XAmzObjectLockRetainUntilDate(t)
	}

	if !data.ObjectLockLegalHold.IsNull() {
		*req = req.XAmzObjectLockLegalHold(data.ObjectLockLegalHold.ValueString())
	}

	return nil
}

func fillPutObjectRequest(req *s3.ApiPutObjectRequest, data *objectResourceModel) error {
	fillServerSideEncryptionData(data, req)
	if err := fillContentData(data, req); err != nil {
		return err
	}

	if err := fillObjectLockData(data, req); err != nil {
		return err
	}

	if !data.StorageClass.IsNull() {
		*req = req.XAmzStorageClass(data.StorageClass.ValueString())
	}

	if !data.WebsiteRedirect.IsNull() {
		*req = req.XAmzWebsiteRedirectLocation(data.WebsiteRedirect.ValueString())
	}

	if !data.RequestPayer.IsNull() {
		*req = req.XAmzRequestPayer(data.RequestPayer.ValueString())
	}

	if !data.Tags.IsNull() {
		tags, err := buildQueryString(data.Tags)
		if err != nil {
			return fmt.Errorf("failed to build tags query string: %s", err.Error())
		}
		*req = req.XAmzTagging(tags)
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

func needsMD5Header(data *objectResourceModel) bool {
	return !data.ObjectLockMode.IsNull() || !data.ObjectLockRetainUntilDate.IsNull() || !data.ObjectLockLegalHold.IsNull()
}

func addMD5Header(req *s3.ApiPutObjectRequest, file io.ReadSeeker) error {
	body, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Reset the file pointer to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	md5Hash, err := getMD5Hash(body)
	if err != nil {
		return fmt.Errorf("failed to get MD5 hash: %w", err)
	}
	*req = req.ContentMD5(md5Hash)
	return nil
}

func createTempFile(fileName, content string) (*os.File, error) {
	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}

	return f, nil
}

// hasObjectContentChanges returns true if the plan has changes to the object content.
func hasObjectContentChanges(plan, state *objectResourceModel) bool {
	needsChange := !(plan.Source.Equal(state.Source) &&
		plan.CacheControl.Equal(state.CacheControl) &&
		plan.ContentDisposition.Equal(state.ContentDisposition) &&
		plan.ContentEncoding.Equal(state.ContentEncoding) &&
		plan.ContentLanguage.Equal(state.ContentLanguage) &&
		plan.ContentType.Equal(state.ContentType) &&
		plan.Content.Equal(state.Content) &&
		plan.Expires.Equal(state.Expires) &&
		plan.ServerSideEncryption.Equal(state.ServerSideEncryption) &&
		plan.StorageClass.Equal(state.StorageClass) &&
		plan.WebsiteRedirect.Equal(state.WebsiteRedirect) &&
		plan.ServerSideEncryptionCustomerAlgorithm.Equal(state.ServerSideEncryptionCustomerAlgorithm) &&
		plan.ServerSideEncryptionCustomerKey.Equal(state.ServerSideEncryptionCustomerKey) &&
		plan.ServerSideEncryptionCustomerKeyMD5.Equal(state.ServerSideEncryptionCustomerKeyMD5) &&
		plan.ServerSideEncryptionContext.Equal(state.ServerSideEncryptionContext) &&
		plan.Metadata.Equal(state.Metadata))
	return needsChange
}

func buildQueryString(m types.Map) (string, error) {
	values := url.Values{}
	for k, v := range m.Elements() {
		if v.IsNull() {
			continue
		}
		strVal, ok := v.(types.String)
		if !ok {
			return "", fmt.Errorf("expected string value, got %T", v)
		}
		values.Add(k, strVal.ValueString())
	}
	return values.Encode(), nil
}

func fromTFMap(t types.Map) (map[string]string, error) {
	m := map[string]string{}
	for k, v := range t.Elements() {
		if v.IsNull() {
			continue
		}
		strVal, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("expected string value, got %T", v)
		}

		m[k] = strVal.ValueString()
	}
	return m, nil
}

func splitImportID(path string) (string, string, error) {
	// Ensure the path is not empty
	if path == "" {
		return "", "", fmt.Errorf("path cannot be empty")
	}

	// Split the path into two parts: bucket and the remaining key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid path format")
	}

	bucket := parts[0]
	key := parts[1]

	return bucket, key, nil
}
