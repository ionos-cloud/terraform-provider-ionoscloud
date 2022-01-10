//go:build nlb
// +build nlb

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceNetworkLoadBalancerForwardingRuleId = DataSource + "." + NetworkLoadBalancerForwardingRuleResource + "." + NetworkLoadBalancerForwardingRuleDataSourceById
const dataSourceNetworkLoadBalancerForwardingRuleName = DataSource + "." + NetworkLoadBalancerForwardingRuleResource + "." + NetworkLoadBalancerForwardingRuleDataSourceByName

func TestAccDataSourceNetworkLoadBalancerForwardingRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic,
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "name", networkLoadBalancerForwardingRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "algorithm", networkLoadBalancerForwardingRuleResource, "algorithm"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "protocol", networkLoadBalancerForwardingRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "listener_ip", networkLoadBalancerForwardingRuleResource, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "listener_port", networkLoadBalancerForwardingRuleResource, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.client_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.connect_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.target_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.retries", networkLoadBalancerForwardingRuleResource, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.ip", networkLoadBalancerForwardingRuleResource, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.port", networkLoadBalancerForwardingRuleResource, "targets.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.weight", networkLoadBalancerForwardingRuleResource, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.check", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.check_interval", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.maintenance", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "name", networkLoadBalancerForwardingRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "algorithm", networkLoadBalancerForwardingRuleResource, "algorithm"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "protocol", networkLoadBalancerForwardingRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "listener_ip", networkLoadBalancerForwardingRuleResource, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "listener_port", networkLoadBalancerForwardingRuleResource, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.client_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.connect_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.target_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.retries", networkLoadBalancerForwardingRuleResource, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.ip", networkLoadBalancerForwardingRuleResource, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.port", networkLoadBalancerForwardingRuleResource, "targets.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.weight", networkLoadBalancerForwardingRuleResource, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.check", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.check_interval", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.maintenance", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance"),
				),
			},
		},
	})
}

const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId = testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic + `
data ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleDataSourceById + ` {
  datacenter_id = ` + networkLoadBalancerForwardingRuleResource + `.datacenter_id
  networkloadbalancer_id  = ` + networkLoadBalancerForwardingRuleResource + `.networkloadbalancer_id
  id			= ` + networkLoadBalancerForwardingRuleResource + `.id
}
`

const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName = testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic + `
data ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + networkLoadBalancerForwardingRuleResource + `.datacenter_id
  networkloadbalancer_id  = ` + networkLoadBalancerForwardingRuleResource + `.networkloadbalancer_id
 name			= ` + networkLoadBalancerForwardingRuleResource + `.name
}
`
