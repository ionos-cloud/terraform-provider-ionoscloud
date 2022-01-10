//go:build nlb
// +build nlb

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceNetworkLoadBalancerId = DataSource + "." + NetworkLoadBalancerResource + "." + NetworkLoadBalancerDataSourceById
const dataSourceNetworkLoadBalancerName = DataSource + "." + NetworkLoadBalancerResource + "." + NetworkLoadBalancerDataSourceByName

func TestAccDataSourceNetworkLoadBalancer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerConfigBasic,
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "name", dataSourceNetworkLoadBalancerId, "name"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", dataSourceNetworkLoadBalancerId, "listener_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "ips", dataSourceNetworkLoadBalancerId, "ips"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", dataSourceNetworkLoadBalancerId, "target_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "lb_private_ips", dataSourceNetworkLoadBalancerId, "lb_private_ips"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "name", dataSourceNetworkLoadBalancerName, "name"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", dataSourceNetworkLoadBalancerName, "listener_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "ips", dataSourceNetworkLoadBalancerName, "ips"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", dataSourceNetworkLoadBalancerName, "target_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "lb_private_ips", dataSourceNetworkLoadBalancerName, "lb_private_ips"),
				),
			},
		},
	})
}

const testAccDataSourceNetworkLoadBalancerMatchId = testAccCheckNetworkLoadBalancerConfigBasic + `
data ` + NetworkLoadBalancerResource + ` ` + NetworkLoadBalancerDataSourceById + ` {
  datacenter_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.datacenter_id
  id			= ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.id
}
`

const testAccDataSourceNetworkLoadBalancerMatchName = testAccCheckNetworkLoadBalancerConfigBasic + `
data ` + NetworkLoadBalancerResource + ` ` + NetworkLoadBalancerDataSourceByName + ` {
  datacenter_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.datacenter_id
  name			= ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.name
}
`
