//go:build all || cr
// +build all cr

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceContainerRegistryLocations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckContainerRegistryLocations,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.ContainerRegistryLocationsResource+"."+constant.ContainerRegistryLocationsTest, "locations"),
				),
			},
		},
	})
}

const testAccCheckContainerRegistryLocations = `
	data ` + constant.ContainerRegistryLocationsResource + ` ` + constant.ContainerRegistryLocationsTest + ` {
	}
`
