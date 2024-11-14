package objectstorage

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	awsv4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	objstorage "github.com/ionos-cloud/sdk-go-object-storage"
)

// Client is a wrapper around the Object Storage client.
type Client struct {
	client *objstorage.APIClient
}

const ionosAPIURLObjectStorage = "IONOS_API_URL_OBJECT_STORAGE"

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objstorage.APIClient {
	return c.client
}

// NewClient creates a new Object Storage client with the given credentials and region.
func NewClient(id, secret, region, endpoint string) *Client {
	// Set custom endpoint if provided
	if envValue := os.Getenv(ionosAPIURLObjectStorage); envValue != "" {
		endpoint = envValue
	}

	cfg := objstorage.NewConfiguration(endpoint)
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
		if errors.Is(err, credentials.ErrStaticCredentialsEmpty) {
			return errors.New("Object Storage credentials are missing. Please set s3_access_key and s3_secret_key provider attributes or environment variables IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY")
		}
		return err
	}
	return &Client{
		client: objstorage.NewAPIClient(cfg),
	}
}
