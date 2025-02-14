package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

type ClientOptions struct {
	shared.ClientOverrideOptions
	Version          string
	TerraformVersion string
}
