package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGatewayRule_Basic(t *testing.T) {
	var natGatewayRule ionoscloud.NatGatewayRule
	natGatewayRuleName := "natGatewayRule"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigBasic, natGatewayRuleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayRuleExists("ionoscloud_natgateway_rule.natgateway_rule", &natGatewayRule),
					resource.TestCheckResourceAttr("ionoscloud_natgateway_rule.natgateway_rule", "name", natGatewayRuleName),
					resource.TestCheckResourceAttrPair("ionoscloud_natgateway_rule.natgateway_rule", "public_ip", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_natgateway_rule.natgateway_rule", "name", "updated"),
					resource.TestCheckResourceAttrPair("ionoscloud_natgateway_rule.natgateway_rule", "public_ip", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
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
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
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

		foundNatGatewayRule, _, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.ID).Execute()

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
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_ipblock" "natgateway_ips" {
  location = ionoscloud_datacenter.natgateway_datacenter.location
  size = 1
  name = "natgateway_ips"
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
  }
}

resource "ionoscloud_natgateway_rule" "natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name          = "%s"
  type          = "SNAT"
  protocol      = "TCP"
  source_subnet = "10.0.1.0/24"
  public_ip     = ionoscloud_ipblock.natgateway_ips.ips[0]
  target_subnet = "10.0.1.0/24"
  target_port_range {
      start = 500
      end   = 1000
  }
}
`

const testAccCheckNatGatewayRuleConfigUpdate = `
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_ipblock" "natgateway_ips" {
  location = ionoscloud_datacenter.natgateway_datacenter.location
  size = 1
  name = "natgateway_ips"
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
  }
}

resource "ionoscloud_natgateway_rule" "natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name          = "updated"
  type          = "SNAT"
  protocol      = "TCP"
  source_subnet = "10.0.1.0/24"
  public_ip     = ionoscloud_ipblock.natgateway_ips.ips[0]
  target_subnet = "10.0.1.0/24"
  target_port_range {
      start = 500
      end   = 1000
  }
}`
