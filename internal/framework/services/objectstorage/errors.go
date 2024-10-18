package objectstorage

import (
	"errors"
	"fmt"

	objstorage "github.com/ionos-cloud/sdk-go-s3"
)

func formatXMLError(err error) error {
	var apiErr objstorage.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if objErr, ok := apiErr.Model().(objstorage.Error); ok {
			msg := ""
			if objErr.Code != nil {
				msg += fmt.Sprintf("code:%s\n", *objErr.Code)
			}
			if objErr.Message != nil {
				msg += fmt.Sprintf("message:%s\n", *objErr.Message)
			}
			if objErr.HostId != nil {
				msg += fmt.Sprintf("host:%s\n", *objErr.HostId)
			}
			if objErr.RequestId != nil {
				msg += fmt.Sprintf("request:%s\n", *objErr.RequestId)
			}
			return errors.New(msg)
		}
	}
	return err
}
