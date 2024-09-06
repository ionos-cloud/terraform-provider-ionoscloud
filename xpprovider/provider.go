package xpprovider

import (
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	internalfwProvider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetProvider is used to get the provider for crossplane upjet provider
func GetProvider() (fwprovider.Provider, *schema.Provider) {
	return internalfwProvider.New(), ionoscloud.Provider()
}
