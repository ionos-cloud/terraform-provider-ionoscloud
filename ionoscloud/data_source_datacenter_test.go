package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDatacenterConfigBasic,
			},
			{
				Config: testAccDataSourceDatacenterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
		},
	})

}

const testAccDataSourceDatacenterMatchId = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceById + ` {
  id			= ` + DatacenterResource + `.` + DatacenterTestResource + `.id
}`

const testAccDataSourceDatacenterMatchName = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceByName + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
}`

const testAccDataSourceDatacenterMatching = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`
