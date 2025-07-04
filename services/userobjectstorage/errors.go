package userobjectstorage

import (
	"encoding/xml"
	"errors"
	"log"

	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func isBucketNotEmptyError(err error) bool {
	var apiErr shared.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var objStoreErr objstorage.Error
		if err := xml.Unmarshal(body, &objStoreErr); err != nil {
			log.Printf("failed to unmarshal error response: %v", err)
			return false
		}

		if objStoreErr.Code != nil && *objStoreErr.Code == "BucketNotEmpty" {
			return true
		}
	}
	return false
}

func isInvalidStateBucketWithObjectLock(err error) bool {
	var apiErr shared.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var objStoreErr objstorage.Error
		if err := xml.Unmarshal(body, &objStoreErr); err != nil {
			log.Printf("failed to unmarshal error response: %v", err)
			return false
		}

		if objStoreErr.Code != nil && *objStoreErr.Code == "InvalidBucketState" &&
			objStoreErr.Message != nil && *objStoreErr.Message == "bucket versioning cannot be disabled on buckets with object lock enabled" {
			return true
		}
	}
	return false
}
