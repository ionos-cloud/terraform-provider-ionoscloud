package xpprovider

import (
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	internalfwProvider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetProvider exports the two providers sdkv2 and framework
// Used in crossplane upjet ionoscloud provider
func GetProvider() (fwprovider.Provider, *schema.Provider) {
	return internalfwProvider.New(), ionoscloud.Provider()
}
