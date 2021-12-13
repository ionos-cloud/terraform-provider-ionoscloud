//go:build natgateway
// +build natgateway

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceIdNatGatewayRuleResource = DataSource + "." + NatGatewayRuleResource + "." + NatGatewayDataSourceById
const dataSourceNameNatGatewayRuleResource = DataSource + "." + NatGatewayRuleResource + "." + NatGatewayDataSourceByName

func TestAccDataSourceNatGatewayRule(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigBasic, NatGatewayRuleTestResource),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchId, NatGatewayRuleTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "name", resourceNatGatewayRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "type", resourceNatGatewayRuleResource, "type"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "protocol", resourceNatGatewayRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "source_subnet", resourceNatGatewayRuleResource, "source_subnet"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "public_ip", resourceNatGatewayRuleResource, "public_ip"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "target_subnet", resourceNatGatewayRuleResource, "target_subnet"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "target_port_range.0.start", resourceNatGatewayRuleResource, "target_port_range.0.start"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayRuleResource, "target_port_range.0.end", resourceNatGatewayRuleResource, "target_port_range.0.end"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchName, NatGatewayRuleTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "name", resourceNatGatewayRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "type", resourceNatGatewayRuleResource, "type"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "protocol", resourceNatGatewayRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "source_subnet", resourceNatGatewayRuleResource, "source_subnet"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "public_ip", resourceNatGatewayRuleResource, "public_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "target_subnet", resourceNatGatewayRuleResource, "target_subnet"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "target_port_range.0.start", resourceNatGatewayRuleResource, "target_port_range.0.start"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayRuleResource, "target_port_range.0.end", resourceNatGatewayRuleResource, "target_port_range.0.end"),
				),
			},
		},
	})
}

const testAccDataSourceNatGatewayRuleMatchId = testAccCheckNatGatewayRuleConfigBasic + `
data ` + NatGatewayRuleResource + ` ` + NatGatewayRuleDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  natgateway_id = ` + NatGatewayResource + `.natgateway.id
  id			= ` + NatGatewayRuleResource + `.` + NatGatewayRuleTestResource + `.id
}
`

const testAccDataSourceNatGatewayRuleMatchName = testAccCheckNatGatewayRuleConfigBasic + `
data ` + NatGatewayRuleResource + ` ` + NatGatewayRuleDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  natgateway_id = ` + NatGatewayResource + `.natgateway.id
  name			= ` + NatGatewayRuleResource + `.` + NatGatewayRuleTestResource + `.name
}
`
