//go:build compute || all || pcc

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePcc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCrossConnectConfigBasic,
			},
			{
				Config: testAccDataSourcePccMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "name", PCCResource+"."+PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "description", PCCResource+"."+PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "peers", PCCResource+"."+PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "connectable_datacenters", PCCResource+"."+PCCTestResource, "connectable_datacenters"),
				),
			},
			{
				Config: testAccDataSourcePccMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "name", PCCResource+"."+PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "description", PCCResource+"."+PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "peers", PCCResource+"."+PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "connectable_datacenters", PCCResource+"."+PCCTestResource, "connectable_datacenters"),
				),
			},
		},
	})
}

const testAccDataSourcePccMatchId = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + PCCResource + ` ` + PCCDataSourceById + ` {
  id			= ` + PCCResource + `.` + PCCTestResource + `.id
}
`

const testAccDataSourcePccMatchName = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + PCCResource + ` ` + PCCDataSourceByName + ` {
  name			= "` + PCCTestResource + `"
}
`
