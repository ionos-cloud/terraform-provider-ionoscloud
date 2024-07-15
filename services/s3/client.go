package s3

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	awsv4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

func NewClient(id, secret, region string) *s3.APIClient {
	cfg := s3.NewConfiguration()
	signer := awsv4.NewSigner(credentials.NewStaticCredentials(id, secret, ""))
	cfg.MiddlewareWithError = func(r *http.Request) error {
		var reader io.ReadSeeker
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			reader = bytes.NewReader(bodyBytes)
		}

		if region == "" {
			region = "eu-central-3"
		}
		_, err := signer.Sign(r, reader, "s3", region, time.Now())
		return err
	}
	return s3.NewAPIClient(cfg)
}
