package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ListResources returns the list of list resources for the compute package.
//
// ionoscloud_datacenter is still an SDKv2 resource; only the list resource lives on
// the framework side, and it needs sdkv2Provider to read the managed resource's
// protocol schemas. See resource_datacenter_list.go for how the two are connected.
func ListResources(sdkv2Provider *schema.Provider) []func() list.ListResource {
	return []func() list.ListResource{
		func() list.ListResource { return NewDatacenterListResource(sdkv2Provider) },
	}
}
