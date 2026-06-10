package pgsqlv2

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns the list of resources for the package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

// ListResources returns the list of list resources for the package.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewClusterListResource,
	}
}
