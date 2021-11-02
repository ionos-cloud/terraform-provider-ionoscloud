package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var resourceNameTargetGroupById = DataSource + "." + TargetGroupResource + "." + TargetGroupDataSourceById
var resourceNameTargetGroupByName = DataSource + "." + TargetGroupResource + "." + TargetGroupDataSourceByName

func TestAccDataSourceTargetGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTargetGroupConfigBasic,
			},
			{
				Config: testAccDataSourceTargetGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "name", resourceNameTargetGroup, "name"),
				),
			},
			{
				Config: testAccDataSourceTargetGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "name", resourceNameTargetGroup, "name"),
				),
			},
		},
	})
}

var testAccDataSourceTargetGroupMatchId = testAccCheckTargetGroupConfigBasic + `
data ` + TargetGroupResource + ` ` + TargetGroupDataSourceById + ` {
  id			= ` + resourceNameTargetGroup + `.id
}
`

var testAccDataSourceTargetGroupMatchName = testAccCheckTargetGroupConfigBasic + `
data ` + TargetGroupResource + ` ` + TargetGroupDataSourceByName + ` {
  name			= ` + resourceNameTargetGroup + `.name
}
`
