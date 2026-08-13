package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// ListResources returns the list of list resources for the compute package.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewDatacenterListResource,
	}
}
