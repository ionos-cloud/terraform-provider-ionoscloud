package s3

import (
	"context"
	"fmt"
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
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	s3 "github.com/ionos-cloud/sdk-go-s3"
	"github.com/mitchellh/go-homedir"
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
				Sensitive:   true,
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
			"etag": schema.StringAttribute{
				Description: "An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				// TODO this may be also computed, we need to check if it's sent in the response
				Optional: true,
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

	var data objectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := uploadObject(ctx, r.client, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create object", err.Error())
		return
	}

	etag := strings.Trim(apiResponse.Header.Get("ETag"), "\"")
	if etag != "" {
		data.Etag = types.StringValue(etag)
	}

	versionID := apiResponse.Header.Get("VersionId")
	if versionID != "" {
		data.VersionID = types.StringValue(versionID)
	}

	contentType, err := getContentType(ctx, &data, r.client)
	if err != nil {
		resp.Diagnostics.AddError("failed to get content type", err.Error())
		return
	}
	data.ContentType = types.StringValue(contentType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the object.
func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data objectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := FindObject(ctx, r.client, &data)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to read object", err.Error())
		return
	}

	diags := r.setDataModel(ctx, &data, apiResponse)
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
	var plan, state objectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasObjectContentChanges(plan, state) {
		_, err := uploadObject(ctx, r.client, &plan)
		if err != nil {
			resp.Diagnostics.AddError("failed to update object", err.Error())
			return
		}
	}

	_, apiResponse, err := FindObject(ctx, r.client, &plan)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to read object", err.Error())
		return
	}

	newState := objectResourceModel{
		Key:    state.Key,
		Bucket: state.Bucket,
	}

	diags := r.setDataModel(ctx, &newState, apiResponse)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
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

	_, apiResponse, err := deleteObject(ctx, r.client, &data)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete object", err.Error())
		return
	}
}

func uploadObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (*shared.APIResponse, error) {
	putReq := client.ObjectsApi.PutObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	err := fillPutObjectRequest(&putReq, data)
	if err != nil {
		return nil, err
	}

	body, err := getBody(data)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = body.Close()
		if err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()

	return putReq.Body(body).Execute()
}

func setContentData(data *objectResourceModel, apiResponse *shared.APIResponse) {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	versionID := apiResponse.Header.Get("x-amz-version-id")
	if versionID != "" {
		data.VersionID = types.StringValue(versionID)
	}

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

func setServerSideEncryptionData(data *objectResourceModel, apiResponse *shared.APIResponse) {
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

func (r *objectResource) setTagsData(ctx context.Context, data *objectResourceModel) diag.Diagnostics {
	tagsMap, err := getTags(ctx, r.client, data)
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

func setMetadata(ctx context.Context, data *objectResourceModel, apiResponse *shared.APIResponse) diag.Diagnostics {
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

func (r *objectResource) setDataModel(ctx context.Context, data *objectResourceModel, apiResponse *shared.APIResponse) diag.Diagnostics {
	setContentData(data, apiResponse)
	setServerSideEncryptionData(data, apiResponse)

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

	diags := setMetadata(ctx, data, apiResponse)
	if diags.HasError() {
		return diags
	}

	diags = r.setTagsData(ctx, data)
	if diags.HasError() {
		return diags
	}

	return nil
}

func getTags(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (map[string]string, error) {
	output, apiResponse, err := client.TaggingApi.GetObjectTagging(ctx, data.Bucket.ValueString(), data.Key.ValueString()).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil, nil
		}

		return nil, err
	}

	tagsMap := map[string]string{}
	for _, tag := range output.TagSet {
		tagsMap[tag.Key] = tag.Value
	}

	return tagsMap, nil
}

func getMetadataMapFromHeaders(apiResponse *shared.APIResponse, prefix string) map[string]string {
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

func getContentType(ctx context.Context, data *objectResourceModel, client *s3.APIClient) (string, error) {
	_, apiResponse, err := FindObject(ctx, client, data)
	if err != nil {
		return "", err
	}

	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	return contentType, nil
}

func deleteObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (map[string]interface{}, *shared.APIResponse, error) {
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

// FindObject finds the object.
func FindObject(ctx context.Context, client *s3.APIClient, data *objectResourceModel) (*s3.HeadObjectOutput, *shared.APIResponse, error) {
	req := client.ObjectsApi.HeadObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if !data.Etag.IsNull() {
		req = req.IfMatch(data.Etag.ValueString())
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

	if !data.WebsiteRedirect.IsNull() {
		*req = req.XAmzWebsiteRedirectLocation(data.WebsiteRedirect.ValueString())
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

func fillPutObjectRequest(req *s3.ApiPutObjectRequest, data *objectResourceModel) error {
	fillServerSideEncryptionData(data, req)
	if err := fillContentData(data, req); err != nil {
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
func hasObjectContentChanges(plan, state objectResourceModel) bool {
	needsChange := !(plan.Source.Equal(state.Source) &&
		plan.CacheControl.Equal(state.CacheControl) &&
		plan.ContentDisposition.Equal(state.ContentDisposition) &&
		plan.ContentEncoding.Equal(state.ContentEncoding) &&
		plan.ContentLanguage.Equal(state.ContentLanguage) &&
		plan.ContentType.Equal(state.ContentType) &&
		plan.Content.Equal(state.Content) &&
		plan.Etag.Equal(state.Etag) &&
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
