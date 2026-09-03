package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ListResources returns the list of list resources for the compute package.
//
// ionoscloud_datacenter is still an SDKv2 resource, so its list resource needs that
// resource's schema. sdkv2Provider is the only way to reach it: package ionoscloud
// cannot be imported here, because its own in-package tests import the framework
// provider, which would make the import a cycle.
//
// When the schema cannot be reached the list resource is left unregistered rather
// than registered without one - the framework fails GetProviderSchema, terraform's
// first RPC, for a list resource with no schemas, which would take down the whole
// provider instead of just `terraform query` on ionoscloud_datacenter.
func ListResources(sdkv2Provider *schema.Provider) []func() list.ListResource {
	if sdkv2Provider == nil {
		return nil
	}

	datacenterResource, ok := sdkv2Provider.ResourcesMap[datacenterResourceType]
	if !ok {
		return nil
	}

	return []func() list.ListResource{
		func() list.ListResource { return NewDatacenterListResource(datacenterResource) },
	}
}
