package objectstorage

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/mitchellh/go-homedir"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
	hash2 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/hash"
)

// ObjectResourceModel is the schema for the ionoscloud_s3_object resource.
type ObjectResourceModel struct {
	Bucket                                types.String `tfsdk:"bucket"`
	Key                                   types.String `tfsdk:"key"`
	CacheControl                          types.String `tfsdk:"cache_control"`
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
	ObjectLockMode                        types.String `tfsdk:"object_lock_mode"`
	ObjectLockRetainUntilDate             types.String `tfsdk:"object_lock_retain_until_date"`
	ObjectLockLegalHold                   types.String `tfsdk:"object_lock_legal_hold"`
	RequestPayer                          types.String `tfsdk:"request_payer"`
	Etag                                  types.String `tfsdk:"etag"`
	Metadata                              types.Map    `tfsdk:"metadata"`
	Tags                                  types.Map    `tfsdk:"tags"`
	VersionID                             types.String `tfsdk:"version_id"`
	Source                                types.String `tfsdk:"source"`
	Content                               types.String `tfsdk:"content"`
	MFA                                   types.String `tfsdk:"mfa"`
	ForceDestroy                          types.Bool   `tfsdk:"force_destroy"`
}

// ObjectDataSourceModel is the schema for the ionoscloud_s3_object data source.
type ObjectDataSourceModel struct {
	Bucket                                types.String `tfsdk:"bucket"`
	Key                                   types.String `tfsdk:"key"`
	CacheControl                          types.String `tfsdk:"cache_control"`
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
	ObjectLockMode                        types.String `tfsdk:"object_lock_mode"`
	ObjectLockRetainUntilDate             types.String `tfsdk:"object_lock_retain_until_date"`
	ObjectLockLegalHold                   types.String `tfsdk:"object_lock_legal_hold"`
	RequestPayer                          types.String `tfsdk:"request_payer"`
	Etag                                  types.String `tfsdk:"etag"`
	Metadata                              types.Map    `tfsdk:"metadata"`
	Tags                                  types.Map    `tfsdk:"tags"`
	VersionID                             types.String `tfsdk:"version_id"`
	ContentLength                         types.Int64  `tfsdk:"content_length"`
	Range                                 types.String `tfsdk:"range"`
	Body                                  types.String `tfsdk:"body"`
}

// UploadObject uploads an object to a bucket.
func (c *Client) UploadObject(ctx context.Context, data *ObjectResourceModel) (*shared.APIResponse, error) {
	putReq := c.client.ObjectsApi.PutObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
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

// GetObject retrieves an object.
func (c *Client) GetObject(ctx context.Context, data *ObjectResourceModel) (*ObjectResourceModel, bool, error) {
	_, apiResponse, err := c.findObject(ctx, &objectFindRequest{
		Bucket:                                data.Bucket,
		Key:                                   data.Key,
		VersionID:                             data.VersionID,
		Etag:                                  data.Etag,
		ServerSideEncryptionCustomerAlgorithm: data.ServerSideEncryptionCustomerAlgorithm,
		ServerSideEncryptionCustomerKey:       data.ServerSideEncryptionCustomerKey,
		ServerSideEncryptionCustomerKeyMD5:    data.ServerSideEncryptionCustomerKeyMD5,
	})

	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	if err = c.setObjectModelData(ctx, apiResponse, data); err != nil {
		return nil, true, err
	}

	return data, true, nil
}

// UpdateObject updates an object.
func (c *Client) UpdateObject(ctx context.Context, plan, state *ObjectResourceModel) error {
	if hasObjectContentChanges(plan, state) {
		resp, err := c.UploadObject(ctx, plan)
		if err != nil {
			return err
		}

		if err = c.SetObjectComputedAttributes(ctx, plan, resp); err != nil {
			return err
		}
	} else {
		if err := c.updateObjectLock(ctx, plan, state); err != nil {
			return err
		}
	}

	if !plan.Tags.Equal(state.Tags) {
		if err := c.UpdateObjectTags(ctx, plan.Bucket.ValueString(), plan.Key.ValueString(), tags.NewFromMap(plan.Tags), tags.NewFromMap(state.Tags)); err != nil {
			return err
		}
	}

	setStateForUnknown(plan, state)
	objData, found, err := c.GetObject(ctx, plan)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("object not found")
	}

	*plan = *objData
	return nil
}

// DeleteObject deletes an object.
func (c *Client) DeleteObject(ctx context.Context, data *ObjectResourceModel) error {
	var (
		err  error
		resp *shared.APIResponse
	)

	if !data.VersionID.IsNull() {
		_, err = DeleteAllObjectVersions(ctx, c.client, &DeleteRequest{
			Bucket:       data.Bucket.ValueString(),
			Key:          data.Key.ValueString(),
			ForceDestroy: data.ForceDestroy.ValueBool(),
			VersionID:    data.VersionID.ValueString(),
		})
	} else {
		_, resp, err = deleteObjectByModel(ctx, c.client, data)
		if resp.HttpNotFound() {
			return nil
		}
	}

	return err
}

// SetObjectComputedAttributes sets computed attributes for an object.
func (c *Client) SetObjectComputedAttributes(ctx context.Context, data *ObjectResourceModel, apiResponse *shared.APIResponse) error {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	data.VersionID = types.StringValue(apiResponse.Header.Get("x-amz-version-id"))

	etag := apiResponse.Header.Get("ETag")
	if etag != "" {
		data.Etag = types.StringValue(strings.Trim(etag, "\""))
	}

	data.ServerSideEncryption = types.StringValue(apiResponse.Header.Get("x-amz-server-side-encryption"))

	contentType, err := c.getContentType(ctx, &objectFindRequest{
		Bucket:                                data.Bucket,
		Key:                                   data.Key,
		VersionID:                             data.VersionID,
		Etag:                                  data.Etag,
		ServerSideEncryptionCustomerAlgorithm: data.ServerSideEncryptionCustomerAlgorithm,
		ServerSideEncryptionCustomerKey:       data.ServerSideEncryptionCustomerKey,
		ServerSideEncryptionCustomerKeyMD5:    data.ServerSideEncryptionCustomerKeyMD5,
	})
	if err != nil {
		return err
	}

	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	return nil
}

func (c *Client) updateObjectLock(ctx context.Context, plan, state *ObjectResourceModel) error {
	if !plan.ObjectLockLegalHold.Equal(state.ObjectLockLegalHold) {
		_, err := c.client.ObjectLockApi.PutObjectLegalHold(ctx, state.Bucket.ValueString(), state.Key.ValueString()).
			ObjectLegalHoldConfiguration(objstorage.ObjectLegalHoldConfiguration{Status: plan.ObjectLockLegalHold.ValueStringPointer()}).
			Execute()
		if err != nil {
			return fmt.Errorf("failed to update object lock legal hold: %w", err)
		}
	}

	if !plan.ObjectLockMode.Equal(state.ObjectLockMode) || !plan.ObjectLockRetainUntilDate.Equal(state.ObjectLockRetainUntilDate) {
		if _, err := c.putRetention(ctx, plan, state); err != nil {
			return fmt.Errorf("failed to update object lock retention: %w", err)
		}
	}

	return nil
}

func (c *Client) putRetention(ctx context.Context, plan, state *ObjectResourceModel) (*shared.APIResponse, error) {
	retentionDate, err := getRetentionDate(plan.ObjectLockRetainUntilDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object lock retain until date: %w", err)
	}

	baseReq := c.client.ObjectLockApi.PutObjectRetention(ctx, state.Bucket.ValueString(), state.Key.ValueString()).
		PutObjectRetentionRequest(objstorage.PutObjectRetentionRequest{
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

func (c *Client) setObjectCommonAttributes(ctx context.Context, data *ObjectResourceModel, apiResponse *shared.APIResponse) error {
	setContentData(data, apiResponse)
	setServerSideEncryptionData(data, apiResponse)
	if err := setObjectLockData(data, apiResponse); err != nil {
		return err
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

	metadata, err := getMetadataFromAPIResponse(ctx, apiResponse)
	if err != nil {
		return err
	}
	data.Metadata = metadata

	tagsMap, err := c.getTags(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if err != nil {
		return err
	}
	data.Tags = tagsMap
	return nil
}

func (c *Client) getTags(ctx context.Context, bucket, key string) (types.Map, error) {
	defaultMap := types.MapNull(types.StringType)
	result, err := c.ListObjectTags(ctx, bucket, key)
	if err != nil {
		return defaultMap, err
	}

	tagsMap, err := result.ToMap(ctx)
	if err != nil {
		return defaultMap, err
	}

	return tagsMap, nil
}

func (c *Client) setObjectModelData(ctx context.Context, apiResponse *shared.APIResponse, data *ObjectResourceModel) error {
	if err := c.setObjectCommonAttributes(ctx, data, apiResponse); err != nil {
		return err
	}

	serverSideEncryptionContext := apiResponse.Header.Get("x-amz-server-side-encryption-context")
	if serverSideEncryptionContext != "" {
		data.ServerSideEncryptionContext = types.StringValue(serverSideEncryptionContext)
	}

	return nil
}

func (c *Client) downloadObject(ctx context.Context, data *ObjectDataSourceModel) (string, error) {
	req := c.client.ObjectsApi.GetObject(ctx, data.Bucket.ValueString(), data.Key.ValueString()).VersionId(data.VersionID.ValueString())
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

type objectFindRequest struct {
	Bucket                                types.String
	Key                                   types.String
	Etag                                  types.String
	VersionID                             types.String
	ServerSideEncryptionCustomerAlgorithm types.String
	ServerSideEncryptionCustomerKey       types.String
	ServerSideEncryptionCustomerKeyMD5    types.String
}

func (c *Client) findObject(ctx context.Context, data *objectFindRequest) (*objstorage.HeadObjectOutput, *shared.APIResponse, error) {
	req := c.client.ObjectsApi.HeadObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
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

func setStateForUnknown(plan, state *ObjectResourceModel) {
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

func expandObjectDate(v string) *objstorage.IonosTime {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return &objstorage.IonosTime{Time: time.Time{}.UTC()}
	}

	return &objstorage.IonosTime{Time: t.UTC()}
}

func setContentData(data *ObjectResourceModel, apiResponse *shared.APIResponse) {
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

func setServerSideEncryptionData(data *ObjectResourceModel, apiResponse *shared.APIResponse) {
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

func setObjectLockData(data *ObjectResourceModel, apiResponse *shared.APIResponse) error {
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

func getMetadataFromAPIResponse(ctx context.Context, apiResponse *shared.APIResponse) (types.Map, error) {
	metadataMap := getMetadataMapFromHeaders(apiResponse, "X-Amz-Meta-")

	if len(metadataMap) > 0 {
		metadata, diagErr := types.MapValueFrom(ctx, types.StringType, metadataMap)
		if diagErr.HasError() {
			return types.MapNull(types.StringType), fmt.Errorf("failed to convert metadata: %v", diagErr)
		}

		return metadata, nil
	}

	return types.MapNull(types.StringType), nil
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

func (c *Client) getContentType(ctx context.Context, data *objectFindRequest) (string, error) {
	_, apiResponse, err := c.findObject(ctx, &objectFindRequest{
		Bucket:                                data.Bucket,
		Key:                                   data.Key,
		VersionID:                             data.VersionID,
		Etag:                                  data.Etag,
		ServerSideEncryptionCustomerAlgorithm: data.ServerSideEncryptionCustomerAlgorithm,
		ServerSideEncryptionCustomerKey:       data.ServerSideEncryptionCustomerKey,
		ServerSideEncryptionCustomerKeyMD5:    data.ServerSideEncryptionCustomerKeyMD5,
	})
	if err != nil {
		return "", err
	}

	return apiResponse.Header.Get("Content-Type"), nil
}

func deleteObjectByModel(ctx context.Context, client *objstorage.APIClient, data *ObjectResourceModel) (map[string]interface{}, *shared.APIResponse, error) {
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

func getBody(data *ObjectResourceModel) (*os.File, error) {
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

func fillContentData(data *ObjectResourceModel, req *objstorage.ApiPutObjectRequest) error {
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

func fillServerSideEncryptionData(data *ObjectResourceModel, req *objstorage.ApiPutObjectRequest) {
	// Since server_side_encryption is both OPTIONAL and COMPUTED, the attribute is set to UNKNOWN in the plan when the
	// value is not provided in the terraform configuration, not NULL like we originally expected.
	if !data.ServerSideEncryption.IsUnknown() {
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

func fillObjectLockData(data *ObjectResourceModel, req *objstorage.ApiPutObjectRequest) error {
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

func fillPutObjectRequest(req *objstorage.ApiPutObjectRequest, data *ObjectResourceModel) error {
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

func needsMD5Header(data *ObjectResourceModel) bool {
	return !data.ObjectLockMode.IsNull() || !data.ObjectLockRetainUntilDate.IsNull() || !data.ObjectLockLegalHold.IsNull()
}

func addMD5Header(req *objstorage.ApiPutObjectRequest, file io.ReadSeeker) error {
	body, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Reset the file pointer to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	md5Hash, err := hash2.MD5(body)
	if err != nil {
		return fmt.Errorf("failed to get MD5 hash: %w", err)
	}
	*req = req.ContentMD5(base64.StdEncoding.EncodeToString([]byte(md5Hash)))
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
func hasObjectContentChanges(plan, state *ObjectResourceModel) bool {
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
