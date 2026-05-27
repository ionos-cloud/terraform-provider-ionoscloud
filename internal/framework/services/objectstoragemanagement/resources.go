package objectstoragemanagement

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns the list of resources for the objectstoragemanagement package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewAccesskeyResource,
	}
}

// ListResources returns the list of list resources for the objectstoragemanagement package.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewAccesskeyListResource,
	}
}
