package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

// ClientOptions - options passed to the terraform clients
type ClientOptions struct {
	fileconfiguration.ClientOverrideOptions
	Version          string
	TerraformVersion string
}
