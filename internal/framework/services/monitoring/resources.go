package monitoring

import "github.com/hashicorp/terraform-plugin-framework/resource"

// Resources returns the list of resources for the package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewPipelineResource,
	}
}
