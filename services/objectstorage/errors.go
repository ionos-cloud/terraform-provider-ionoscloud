package objectstorage

import (
	"context"
	"encoding/xml"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func isBucketNotEmptyError(ctx context.Context, err error) bool {
	var apiErr shared.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var objStoreErr objstorage.Error
		if err := xml.Unmarshal(body, &objStoreErr); err != nil {
			tflog.Warn(ctx, "failed to unmarshal error response", map[string]interface{}{"error": err.Error()})
			return false
		}

		if objStoreErr.Code != nil && *objStoreErr.Code == "BucketNotEmpty" {
			return true
		}
	}
	return false
}

func isInvalidStateBucketWithObjectLock(ctx context.Context, err error) bool {
	var apiErr shared.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var objStoreErr objstorage.Error
		if err := xml.Unmarshal(body, &objStoreErr); err != nil {
			tflog.Warn(ctx, "failed to unmarshal error response", map[string]interface{}{"error": err.Error()})
			return false
		}

		if objStoreErr.Code != nil && *objStoreErr.Code == "InvalidBucketState" &&
			objStoreErr.Message != nil && *objStoreErr.Message == "bucket versioning cannot be disabled on buckets with object lock enabled" {
			return true
		}
	}
	return false
}
