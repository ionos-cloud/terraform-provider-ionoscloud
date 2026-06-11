package inmemorydbv2

import "github.com/hashicorp/terraform-plugin-framework/resource"

// Resources returns the list of resources for the inmemorydbv2 package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}
