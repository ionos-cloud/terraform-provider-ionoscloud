package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

type ClientOptions struct {
	fileconfiguration.ClientOverrideOptions
	Version          string
	TerraformVersion string
}
