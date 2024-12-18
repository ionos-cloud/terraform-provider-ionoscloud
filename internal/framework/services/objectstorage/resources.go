package objectstorage

import "github.com/hashicorp/terraform-plugin-framework/resource"

// Resources returns the list of resources for the objectstorage package.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewBucketResource,
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
