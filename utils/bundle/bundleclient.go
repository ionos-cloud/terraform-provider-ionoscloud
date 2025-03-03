package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ClientOptions - options passed to the terraform clients
type ClientOptions struct {
	shared.ClientOptions
	StorageOptions   StorageOptions
	Version          string
	TerraformVersion string
}

// StorageOptions - options passed to the storage clients
type StorageOptions struct {
	AccessKey string
	SecretKey string
	Region    string
}
