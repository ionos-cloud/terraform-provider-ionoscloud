package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ListResources returns the list of list resources for the compute package.
//
// The Cloud API resources listed here are still SDKv2 resources, so their list
// resources need those resources' schemas. sdkv2Provider is the only way to reach
// them: package ionoscloud cannot be imported here, because its own in-package tests
// import the framework provider, which would make the import a cycle.
//
// When a schema cannot be reached that one list resource is left unregistered rather
// than registered without one - the framework fails GetProviderSchema, terraform's
// first RPC, for a list resource with no schemas, which would take down the whole
// provider instead of just `terraform query` on that type. Each lookup is therefore
// independent: a miss on one resource must not unregister the others.
func ListResources(sdkv2Provider *schema.Provider) []func() list.ListResource {
	// provider.New documents nil as a supported argument (list resources for SDKv2
	// resources then simply cannot be served), so this guard is not defensive padding.
	if sdkv2Provider == nil {
		return nil
	}

	var listResources []func() list.ListResource

	if datacenterResource, ok := sdkv2Provider.ResourcesMap[datacenterResourceType]; ok {
		listResources = append(listResources, func() list.ListResource { return NewDatacenterListResource(datacenterResource) })
	}

	if ipblockResource, ok := sdkv2Provider.ResourcesMap[ipblockResourceType]; ok {
		listResources = append(listResources, func() list.ListResource { return NewIPBlockListResource(ipblockResource) })
	}

	if targetGroupResource, ok := sdkv2Provider.ResourcesMap[targetGroupResourceType]; ok {
		listResources = append(listResources, func() list.ListResource { return NewTargetGroupListResource(targetGroupResource) })
	}

	return listResources
}
