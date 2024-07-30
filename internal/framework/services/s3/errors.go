package s3

import (
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
