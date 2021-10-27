package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceName = IpBLockResource + ".test_ipblock"

func TestAccDataSourceIpBlock(t *testing.T) {

	dataSourceName := fmt.Sprintf("%s.%s.test_ipblock_data", DataSource, IpBLockResource)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIpBLockResource,
			},
			{
				Config: testAccDataSourceIpBlockMatchProperties,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "location", resourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceName, "size", resourceName, "size"),
				),
			},
		},
	})
}

const testAccDataSourceIpBLockResource = `resource ` + IpBLockResource + ` "test_ipblock" {
  location = "us/las"
  size     = 1
  name     = "test_ipblock_name"
}
`

const testAccDataSourceIpBlockMatchProperties = testAccDataSourceIpBLockResource +
	"data " + IpBLockResource + " test_ipblock_data " +
	"{ id = " + resourceName + ".id }"
