//go:build all || natgateway
// +build all natgateway

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceIdNatGatewayResource = NatGatewayResource + "." + NatGatewayDataSourceById
const dataSourceNameNatGatewayResource = NatGatewayResource + "." + NatGatewayDataSourceByName

func TestAccDataSourceNatGateway_matchId(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigBasic, NatGatewayTestResource),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchId, NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "name", resourceNatGatewayResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "public_ips.0", resourceNatGatewayResource, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.id", resourceNatGatewayResource, "lans.0.id"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.gateway_ips", resourceNatGatewayResource, "lans.0.gateway_ips"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchName, NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "name", resourceNatGatewayResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "public_ips.0", resourceNatGatewayResource, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "lans.0.id", resourceNatGatewayResource, "lans.0.id"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "lans.0.gateway_ips", resourceNatGatewayResource, "lans.0.gateway_ips"),
				),
			},
		},
	})
}

const testAccDataSourceNatGatewayMatchId = testAccCheckNatGatewayConfigBasic + `
data ` + NatGatewayResource + ` ` + NatGatewayDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  id			= ` + NatGatewayResource + `.` + NatGatewayTestResource + `.id
}
`

const testAccDataSourceNatGatewayMatchName = testAccCheckNatGatewayConfigBasic + `
data ` + NatGatewayResource + ` ` + NatGatewayDataSourceByName + `  {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  name			= ` + NatGatewayResource + `.` + NatGatewayTestResource + `.name
}
`
