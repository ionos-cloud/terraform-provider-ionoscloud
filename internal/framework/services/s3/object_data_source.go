package s3

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ datasource.DataSourceWithConfigure = (*objectDataSource)(nil)
)

// NewObjectDataSource creates a new data source for object.
func NewObjectDataSource() datasource.DataSource {
	return &objectDataSource{}
}

type objectDataSource struct {
	client *s3.APIClient
}

type objectDataSourceModel struct {
	Bucket                                types.String `tfsdk:"bucket"`
	Key                                   types.String `tfsdk:"key"`
	CacheControl                          types.String `tfsdk:"cache_control"`
	ContentDisposition                    types.String `tfsdk:"content_disposition"`
	ContentEncoding                       types.String `tfsdk:"content_encoding"`
	ContentLanguage                       types.String `tfsdk:"content_language"`
	ContentType                           types.String `tfsdk:"content_type"`
	ContentLength                         types.Int64  `tfsdk:"content_length"`
	ContentMD5                            types.String `tfsdk:"content_md5"`
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
	Range                                 types.String `tfsdk:"range"`
	VersionID                             types.String `tfsdk:"version_id"`
	Body                                  types.String `tfsdk:"body"`
}

// Metadata returns the metadata for the object data source.
func (d *objectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_object"
}

// Configure configures the data source.
func (d *objectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = client
}

// Schema returns the schema for the object data source.
func (d *objectDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"body": schema.StringAttribute{
				Computed: true,
			},
			"cache_control": schema.StringAttribute{
				Computed: true,
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"content_disposition": schema.StringAttribute{
				Computed: true,
			},
			"content_encoding": schema.StringAttribute{
				Computed: true,
			},
			"content_language": schema.StringAttribute{
				Computed: true,
			},
			"content_type": schema.StringAttribute{
				Computed: true,
			},
			"content_md5": schema.StringAttribute{
				Computed: true,
			},
			"expires": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption": schema.StringAttribute{
				Computed: true,
			},
			"storage_class": schema.StringAttribute{
				Computed: true,
			},
			"website_redirect": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption_context": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"request_payer": schema.StringAttribute{
				Computed: true,
			},
			"object_lock_mode": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Computed: true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"range": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Read the data source
func (d *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data objectDataSourceModel

	// Read configuration
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := findObjectDataSource(ctx, d.client, data)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to read object", err.Error())
		return
	}

	if isContentTypeAllowed(apiResponse.Header.Get("Content-Type")) {
		var body string
		body, err = downloadObject(ctx, d.client, data)
		if err != nil {
			resp.Diagnostics.AddError("failed to download object", err.Error())
			return
		}

		data.Body = types.StringValue(body)
	}

	diags = d.setDataModel(ctx, &data, apiResponse)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func downloadObject(ctx context.Context, client *s3.APIClient, data objectDataSourceModel) (string, error) {
	req := client.ObjectsApi.GetObject(ctx, data.Bucket.ValueString(), data.Key.ValueString()).VersionId(data.VersionID.ValueString())
	if !data.Range.IsNull() {
		req = req.Range_(data.Range.ValueString())
	}

	resp, _, err := req.Execute()
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read object data: %w", err)
	}

	return string(bytes), nil
}

func isContentTypeAllowed(contentType string) bool {
	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile(`^application/atom\+xml$`),
		regexp.MustCompile(`^application/json$`),
		regexp.MustCompile(`^application/ld\+json$`),
		regexp.MustCompile(`^application/x-csh$`),
		regexp.MustCompile(`^application/x-httpd-php$`),
		regexp.MustCompile(`^application/x-sh$`),
		regexp.MustCompile(`^application/xhtml\+xml$`),
		regexp.MustCompile(`^application/xml$`),
		regexp.MustCompile(`^text/.+`),
	}
	for _, r := range allowedContentTypes {
		if r.MatchString(contentType) {
			return true
		}
	}

	return false
}

func findObjectDataSource(ctx context.Context, client *s3.APIClient, data objectDataSourceModel) (*s3.HeadObjectOutput, *shared.APIResponse, error) {
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

func setContentDataSource(data *objectDataSourceModel, apiResponse *shared.APIResponse) diag.Diagnostics {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	contentLength := apiResponse.Header.Get("Content-Length")
	if contentLength != "" {
		intLength, err := strconv.Atoi(contentLength)
		if err != nil {
			diagErr := diag.Diagnostics{}
			diagErr.AddError("failed to convert content length", err.Error())
			return diagErr
		}
		data.ContentLength = types.Int64Value(int64(intLength))
	}

	etag := apiResponse.Header.Get("ETag")
	if etag != "" {
		data.Etag = types.StringValue(etag)
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

	contentMD5 := apiResponse.Header.Get("Content-MD5")
	if contentMD5 != "" {
		data.ContentMD5 = types.StringValue(contentMD5)
	}

	expires := apiResponse.Header.Get("Expires")
	if expires != "" {
		data.Expires = types.StringValue(expires)
	}

	return nil
}

func setServerSideEncryptionDataSource(data *objectDataSourceModel, apiResponse *shared.APIResponse) {
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

func setObjectLockDataSource(data *objectDataSourceModel, apiResponse *shared.APIResponse) {
	objectLockMode := apiResponse.Header.Get("x-amz-object-lock-mode")
	if objectLockMode != "" {
		data.ObjectLockMode = types.StringValue(objectLockMode)
	}

	objectLockRetainUntilDate := apiResponse.Header.Get("x-amz-object-lock-retain-until-date")
	if objectLockRetainUntilDate != "" {
		data.ObjectLockRetainUntilDate = types.StringValue(objectLockRetainUntilDate)
	}

	objectLockLegalHold := apiResponse.Header.Get("x-amz-object-lock-legal-hold")
	if objectLockLegalHold != "" {
		data.ObjectLockLegalHold = types.StringValue(objectLockLegalHold)
	}
}

func (d *objectDataSource) setTagsData(ctx context.Context, data *objectDataSourceModel) diag.Diagnostics {
	tagsMap, err := getTagsDataSource(ctx, d.client, data)
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

func setMetadataDataSource(ctx context.Context, data *objectDataSourceModel, apiResponse *shared.APIResponse) diag.Diagnostics {
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

func (d *objectDataSource) setDataModel(ctx context.Context, data *objectDataSourceModel, apiResponse *shared.APIResponse) diag.Diagnostics {
	if diags := setContentDataSource(data, apiResponse); diags.HasError() {
		return diags
	}

	setObjectLockDataSource(data, apiResponse)
	setServerSideEncryptionDataSource(data, apiResponse)

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

	if diags := setMetadataDataSource(ctx, data, apiResponse); diags.HasError() {
		return diags
	}

	if diags := d.setTagsData(ctx, data); diags.HasError() {
		return diags
	}

	return nil
}

func getTagsDataSource(ctx context.Context, client *s3.APIClient, data *objectDataSourceModel) (map[string]string, error) {
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
