package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var resourceNameAlbById = DataSource + "." + ApplicationLoadBalancerResource + "." + ApplicationLoadBalancerDataSourceById
var resourceNameAlbByName = DataSource + "." + ApplicationLoadBalancerResource + "." + ApplicationLoadBalancerDataSourceByName

func TestAccDataSourceApplicationLoadBalancer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationLoadBalancerConfigBasic,
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlbById, "name", resourceNameAlb, "name"),
				),
			},

			{
				Config: testAccDataSourceApplicationLoadBalancerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlbByName, "name", resourceNameAlb, "name"),
				),
			},
		},
	})
}

const testAccDataSourceApplicationLoadBalancerMatchId = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + ApplicationLoadBalancerResource + ` ` + ApplicationLoadBalancerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  id			= ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
}
`

const testAccDataSourceApplicationLoadBalancerMatchName = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + ApplicationLoadBalancerResource + ` ` + ApplicationLoadBalancerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  name          = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.name
}
`
