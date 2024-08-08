package s3

import (
	"encoding/xml"
	"errors"
	"fmt"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

func formatXMLError(err error) error {
	var apiErr s3.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if s3Error, ok := apiErr.Model().(s3.Error); ok {
			msg := ""
			if s3Error.Code != nil {
				msg += fmt.Sprintf("code:%s\n", *s3Error.Code)
			}
			if s3Error.Message != nil {
				msg += fmt.Sprintf("message:%s\n", *s3Error.Message)
			}
			if s3Error.HostId != nil {
				msg += fmt.Sprintf("host:%s\n", *s3Error.HostId)
			}
			if s3Error.RequestId != nil {
				msg += fmt.Sprintf("request:%s\n", *s3Error.RequestId)
			}
			return errors.New(msg)
		}
	}
	return err
}

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
