package objectstorage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetObjectForDataSource retrieves an object for a data source.
func (c *Client) GetObjectForDataSource(ctx context.Context, data *ObjectDataSourceModel) (*ObjectDataSourceModel, bool, error) {
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

	if err = c.setObjectDataSourceModelData(ctx, apiResponse, data); err != nil {
		return nil, true, err
	}

	return data, true, nil
}

func (c *Client) setObjectDataSourceModelData(ctx context.Context, apiResponse *shared.APIResponse, data *ObjectDataSourceModel) error {
	if err := c.setObjectDataSourceCommonAttributes(ctx, data, apiResponse); err != nil {
		return err
	}

	contentLength := apiResponse.Header.Get("Content-Length")
	if contentLength != "" {
		intLength, err := strconv.Atoi(contentLength)
		if err != nil {
			return fmt.Errorf("failed to parse content length: %w", err)
		}
		data.ContentLength = types.Int64Value(int64(intLength))
	}

	if isContentTypeAllowed(apiResponse.Header.Get("Content-Type")) {
		body, err := c.downloadObject(ctx, data)
		if err != nil {
			return err
		}

		data.Body = types.StringValue(body)
	}

	return nil
}

func (c *Client) setObjectDataSourceCommonAttributes(ctx context.Context, data *ObjectDataSourceModel, apiResponse *shared.APIResponse) error {
	setObjectDataSourceContentData(data, apiResponse)
	setObjectDataSourceServerSideEncryptionData(data, apiResponse)
	if err := setObjectDataSourceObjectLockData(data, apiResponse); err != nil {
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

func setObjectDataSourceContentData(data *ObjectDataSourceModel, apiResponse *shared.APIResponse) {
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

func setObjectDataSourceServerSideEncryptionData(data *ObjectDataSourceModel, apiResponse *shared.APIResponse) {
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

func setObjectDataSourceObjectLockData(data *ObjectDataSourceModel, apiResponse *shared.APIResponse) error {
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
