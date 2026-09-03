package xpprovider

import (
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	internalfwprovider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
)

// GetProvider exports the two providers sdkv2 and framework
// Used in crossplane upjet ionoscloud provider
func GetProvider() (fwprovider.Provider, *schema.Provider) {
	sdkv2Provider := ionoscloud.Provider()
	return internalfwprovider.New(sdkv2Provider), sdkv2Provider
}
