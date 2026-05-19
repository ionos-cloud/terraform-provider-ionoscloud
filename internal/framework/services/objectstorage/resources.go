package objectstorage

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
)

// Resources returns the list of resources for the objectstorage package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		identity.WithIdentity(NewBucketResource),
		NewBucketPolicyResource,
		NewObjectResource,
		NewObjectCopyResource,
		NewBucketPublicAccessBlockResource,
		NewBucketVersioningResource,
		NewObjectLockConfigurationResource,
		NewServerSideEncryptionConfigurationResource,
		NewBucketCorsConfigurationResource,
		NewBucketLifecycleConfigurationResource,
		NewBucketWebsiteConfigurationResource,
	}
}

// ListResources returns the list of list resources for the objectstorage package.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewBucketListResource,
	}
}
