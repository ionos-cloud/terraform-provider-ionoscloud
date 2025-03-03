package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ClientOptions - options passed to the terraform clients
type ClientOptions struct {
	shared.ClientOptions
	Version          string
	TerraformVersion string
}
