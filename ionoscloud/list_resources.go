package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// ListResources returns the list resources for the managed resources in this package.
//
// List resources can only be implemented with terraform-plugin-framework, so a list
// resource for a resource that is still on SDKv2 is framework code living in this
// package. It has to be: the results are produced by the resource's own state writer,
// which is unexported, and no package under internal/framework can import this one -
// ionoscloud/provider_test.go is an in-package test that imports
// internal/framework/provider, so an import back would be a cycle in the test build.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewDatacenterListResource,
	}
}
