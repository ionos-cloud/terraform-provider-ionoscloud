//go:build all || cr
// +build all cr

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceContainerRegistryLocations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckContainerRegistryLocations,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(ContainerRegistryLocationsResource+"."+ContainerRegistryLocationsTest, "locations"),
				),
			},
		},
	})
}

const testAccCheckContainerRegistryLocations = `
	data ` + ContainerRegistryLocationsResource + ` ` + ContainerRegistryLocationsTest + ` {
	}
`
