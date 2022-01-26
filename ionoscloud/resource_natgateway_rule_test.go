//go:build natgateway
// +build natgateway

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceNatGatewayRuleResource = NatGatewayRuleResource + "." + NatGatewayRuleTestResource

func TestAccNatGatewayRuleBasic(t *testing.T) {
	var natGatewayRule ionoscloud.NatGatewayRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigBasic, NatGatewayRuleTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayRuleExists(resourceNatGatewayRuleResource, &natGatewayRule),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "name", NatGatewayRuleTestResource),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "type", "SNAT"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "source_subnet", "10.0.1.0/24"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayRuleResource, "public_ip", IpBlockResource+".natgateway_rule_ips", "ips.0"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_subnet", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_port_range.0.start", "500"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_port_range.0.end", "1000"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigUpdate, UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "type", "SNAT"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "source_subnet", "10.3.1.0/24"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayRuleResource, "public_ip", IpBlockResource+".natgateway_rule_ips", "ips.0"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_subnet", "172.31.0.0/24"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_port_range.0.start", "200"),
					resource.TestCheckResourceAttr(resourceNatGatewayRuleResource, "target_port_range.0.end", "1111")),
			},
		},
	})
}

func testAccCheckNatGatewayRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NatGatewayRuleResource {
			continue
		}

		apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured at checking deletion of nat gateway rule %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("nat gateway rule still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNatGatewayRuleExists(n string, natGateway *ionoscloud.NatGatewayRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNatGatewayRuleExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundNatGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching NatGatewayRule: %s", rs.Primary.ID)
		}
		if *foundNatGatewayRule.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		natGateway = &foundNatGatewayRule

		return nil
	}
}

const testAccCheckNatGatewayRuleConfigBasic = `
resource ` + DatacenterResource + ` "natgateway_rule_datacenter" {
  name              = "test_natgateway_rule"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + IpBlockResource + ` "natgateway_rule_ips" {
  location = ` + DatacenterResource + `.natgateway_rule_datacenter.location
  size = 2
  name = "natgateway_rule_ips"
}

resource ` + LanResource + ` "natgateway_rule_lan" {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id 
  public        = false
  name          = "test_natgateway_rule_lan"
}

resource ` + NatGatewayResource + ` "natgateway" { 
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  name          = "test_natgateway_rule_natgateway" 
  public_ips    = [ ` + IpBlockResource + `.natgateway_rule_ips.ips[0], ` + IpBlockResource + `.natgateway_rule_ips.ips[1] ]
  lans {
     id          = ` + LanResource + `.natgateway_rule_lan.id
     gateway_ips = [ "10.11.2.5"] 
  }
}

resource ` + NatGatewayRuleResource + ` ` + NatGatewayRuleTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  natgateway_id = ` + NatGatewayResource + `.natgateway.id
  name          = "%s"
  type          = "SNAT"
  protocol      = "TCP"
  source_subnet = "10.0.1.0/24"
  public_ip     = ` + IpBlockResource + `.natgateway_rule_ips.ips[0]
  target_subnet = "172.16.0.0/24"
  target_port_range {
      start = 500
      end   = 1000
  }
}
`

const testAccCheckNatGatewayRuleConfigUpdate = `
resource ` + DatacenterResource + ` "natgateway_rule_datacenter" {
  name              = "test_natgateway_rule"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + IpBlockResource + ` "natgateway_rule_ips" {
  location = ` + DatacenterResource + `.natgateway_rule_datacenter.location
  size = 2
  name = "natgateway_rule_ips"
}

resource ` + LanResource + ` "natgateway_rule_lan" {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id 
  public        = false
  name          = "test_natgateway_rule_lan"
}

resource ` + NatGatewayResource + ` "natgateway" { 
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  name          = "test_natgateway_rule_natgateway" 
  public_ips    = [ ` + IpBlockResource + `.natgateway_rule_ips.ips[0], ` + IpBlockResource + `.natgateway_rule_ips.ips[1] ]
  lans {
     id          = ` + LanResource + `.natgateway_rule_lan.id
     gateway_ips = [ "10.11.2.5"] 
  }
}

resource ` + NatGatewayRuleResource + ` ` + NatGatewayRuleTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_rule_datacenter.id
  natgateway_id = ` + NatGatewayResource + `.natgateway.id
  name          = "%s"
  type          = "SNAT"
  protocol      = "UDP"
  source_subnet = "10.3.1.0/24"
  public_ip     = ` + IpBlockResource + `.natgateway_rule_ips.ips[0]
  target_subnet = "172.31.0.0/24"
  target_port_range {
      start = 200
      end   = 1111
  }
}`
