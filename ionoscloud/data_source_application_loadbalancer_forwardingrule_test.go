package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var resourceNameAlbRuleById = DataSource + "." + ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleDataSourceById
var resourceNameAlbRuleByName = DataSource + "." + ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleDataSourceByName

func TestAccDataSourceApplicationLoadBalancerForwardingRule(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic,
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlbRuleById, "name", resourceNameAlbRule, "name"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlbRuleByName, "name", resourceNameAlbRule, "name"),
				),
			},
		},
	})
}

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId = testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic + `
data ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
  id			= ` + ApplicationLoadBalancerForwardingRuleResource + `.` + ApplicationLoadBalancerForwardingRuleTestResource + `.id
}
`

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName = testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic + `
data ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
  name    		= ` + ApplicationLoadBalancerForwardingRuleResource + `.` + ApplicationLoadBalancerForwardingRuleTestResource + `.name
}
`
