package userobjectstorage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ObjectCopyResourceModel defines the fields for the Terraform resource model.
type ObjectCopyResourceModel struct {
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
	ObjectLockMode                        types.String `tfsdk:"object_lock_mode"`
	ObjectLockRetainUntilDate             types.String `tfsdk:"object_lock_retain_until_date"`
	ObjectLockLegalHold                   types.String `tfsdk:"object_lock_legal_hold"`
	Etag                                  types.String `tfsdk:"etag"`
	Metadata                              types.Map    `tfsdk:"metadata"`
	Tags                                  types.Map    `tfsdk:"tags"`
	VersionID                             types.String `tfsdk:"version_id"`
	Source                                types.String `tfsdk:"source"`
	CopyIfMatch                           types.String `tfsdk:"copy_if_match"`
	CopyIfModifiedSince                   types.String `tfsdk:"copy_if_modified_since"`
	CopyIfNoneMatch                       types.String `tfsdk:"copy_if_none_match"`
	CopyIfUnmodifiedSince                 types.String `tfsdk:"copy_if_unmodified_since"`
	SourceCustomerAlgorithm               types.String `tfsdk:"source_customer_algorithm"`
	SourceCustomerKey                     types.String `tfsdk:"source_customer_key"`
	SourceCustomerKeyMD5                  types.String `tfsdk:"source_customer_key_md5"`
	MetadataDirective                     types.String `tfsdk:"metadata_directive"`
	TaggingDirective                      types.String `tfsdk:"tagging_directive"`
	LastModified                          types.String `tfsdk:"last_modified"`
	ForceDestroy                          types.Bool   `tfsdk:"force_destroy"`
}

// CopyObject copies an object.
func (c *Client) CopyObject(ctx context.Context, data *ObjectCopyResourceModel) error {
	req := c.client.ObjectsApi.CopyObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	err := fillObjectCopyRequest(&req, data)
	if err != nil {
		return err
	}

	output, apiResponse, err := req.Execute()
	if err != nil {
		return err
	}

	return c.setObjectCopyComputedAttributes(ctx, data, apiResponse, output)
}

// GetObjectCopy gets an object copy.
func (c *Client) GetObjectCopy(ctx context.Context, data *ObjectCopyResourceModel) (*ObjectCopyResourceModel, bool, error) {
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

	if err = c.setObjectCopyCommonAttributes(ctx, data, apiResponse); err != nil {
		return nil, true, err
	}

	return data, true, nil
}

// UpdateObjectCopy updates an object copy.
func (c *Client) UpdateObjectCopy(ctx context.Context, plan, state *ObjectCopyResourceModel) error {
	if hasCopyConditions(plan) || hasObjectCopyContentChanges(plan, state) {
		if err := c.CopyObject(ctx, plan); err != nil {
			return err
		}
	}

	if plan.VersionID.IsUnknown() {
		plan.VersionID = state.VersionID
	}

	if plan.Etag.IsUnknown() {
		plan.Etag = state.Etag
	}

	objData, found, err := c.GetObjectCopy(ctx, plan)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("object not found")
	}

	*plan = *objData
	return nil
}

// DeleteObjectCopy deletes an object copy.
func (c *Client) DeleteObjectCopy(ctx context.Context, data *ObjectCopyResourceModel) error {
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
		_, resp, err = deleteObjectCopyByModel(ctx, c.client, data)
		if resp.HttpNotFound() {
			return nil
		}
	}

	return err
}

func (c *Client) setObjectCopyComputedAttributes(ctx context.Context, data *ObjectCopyResourceModel, apiResponse *shared.APIResponse, output *objstorage.CopyObjectResult) error {
	contentType := apiResponse.Header.Get("Content-Type")
	if contentType != "" {
		data.ContentType = types.StringValue(contentType)
	}

	data.VersionID = types.StringValue(apiResponse.Header.Get("x-amz-version-id"))

	if output.ETag != nil {
		data.Etag = types.StringValue(strings.Trim(*output.ETag, "\""))
	}

	data.LastModified = types.StringValue(output.LastModified.Format(time.RFC3339))
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

func deleteObjectCopyByModel(ctx context.Context, client *objstorage.APIClient, data *ObjectCopyResourceModel) (map[string]interface{}, *shared.APIResponse, error) {
	req := client.ObjectsApi.DeleteObject(ctx, data.Bucket.ValueString(), data.Key.ValueString())
	if !data.VersionID.IsNull() {
		req = req.VersionId(data.VersionID.ValueString())
	}

	if !data.ForceDestroy.IsNull() {
		req = req.XAmzBypassGovernanceRetention(data.ForceDestroy.ValueBool())
	}

	return req.Execute()
}

func hasCopyConditions(plan *ObjectCopyResourceModel) bool {
	return !plan.CopyIfMatch.IsNull() || !plan.CopyIfModifiedSince.IsNull() || !plan.CopyIfNoneMatch.IsNull() || !plan.CopyIfUnmodifiedSince.IsNull()
}

func fillObjectCopyContentData(data *ObjectCopyResourceModel, req *objstorage.ApiCopyObjectRequest) error {
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

func fillObjectCopyServerSideEncryptionData(data *ObjectCopyResourceModel, req *objstorage.ApiCopyObjectRequest) {
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

func fillObjectCopyLockData(data *ObjectCopyResourceModel, req *objstorage.ApiCopyObjectRequest) error {
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

func fillObjectCopyRequest(req *objstorage.ApiCopyObjectRequest, data *ObjectCopyResourceModel) error {
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
func hasObjectCopyContentChanges(plan, state *ObjectCopyResourceModel) bool {
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

func (c *Client) setObjectCopyCommonAttributes(ctx context.Context, data *ObjectCopyResourceModel, apiResponse *shared.APIResponse) error {
	setObjectCopyContentData(data, apiResponse)
	setObjectCopyServerSideEncryptionData(data, apiResponse)
	if err := setObjectCopyObjectLockData(data, apiResponse); err != nil {
		return err
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

func setObjectCopyContentData(data *ObjectCopyResourceModel, apiResponse *shared.APIResponse) {
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

func setObjectCopyServerSideEncryptionData(data *ObjectCopyResourceModel, apiResponse *shared.APIResponse) {
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

func setObjectCopyObjectLockData(data *ObjectCopyResourceModel, apiResponse *shared.APIResponse) error {
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
