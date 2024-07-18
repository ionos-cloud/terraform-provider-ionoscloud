package s3

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

func formatXMLError(err error) error {
	var apiErr shared.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if s3Error, ok := apiErr.Model().(s3.Error); ok {
			err = fmt.Errorf("code:%s\nmessage: %s\nhost:%s\nrequest:%s", s3Error.GetCode(), s3Error.GetMessage(), s3Error.GetHostId(), s3Error.GetRequestId())
			return err
		}
	}
	return err
}
