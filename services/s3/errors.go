package s3

import (
	"encoding/xml"
	"errors"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

func isBucketNotEmptyError(err error) bool {
	var apiErr s3.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var s3Err s3.Error
		if err := xml.Unmarshal(body, &s3Err); err != nil {
			return false
		}

		if s3Err.Code != nil && *s3Err.Code == "BucketNotEmpty" {
			return true
		}
	}
	return false
}

func isInvalidStateBucketWithObjectLock(err error) bool {
	var apiErr s3.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		body := apiErr.Body()
		var s3Err s3.Error
		if err := xml.Unmarshal(body, &s3Err); err != nil {
			return false
		}

		if s3Err.Code != nil && *s3Err.Code == "InvalidBucketState" &&
			s3Err.Message != nil && *s3Err.Message == "bucket versioning cannot be disabled on buckets with object lock enabled" {
			return true
		}
	}
	return false
}
