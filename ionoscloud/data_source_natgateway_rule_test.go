// +build natgateway

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNatGatewayRule(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway_rule.test_natgateway_rule_id", "name", "test_datasource_natgateway_rule"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_natgateway_rule.test_natgateway_rule_id", "public_ip", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway_rule.test_natgateway_rule_name", "name", "test_datasource_natgateway_rule"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_natgateway_rule.test_natgateway_rule_name", "public_ip", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
		},
	})
}

const testAccDataSourceNatGatewayRuleCreateResources = `
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
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.12.1.2/24"]
  }
}
resource "ionoscloud_natgateway_rule" "natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name          = "test_datasource_natgateway_rule"
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

const testAccDataSourceNatGatewayRuleMatchId = `
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
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.12.1.2/24"]
  }
}
resource "ionoscloud_natgateway_rule" "natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name          = "test_datasource_natgateway_rule"
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

data "ionoscloud_natgateway_rule" "test_natgateway_rule_id" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  id			= ionoscloud_natgateway_rule.natgateway_rule.id
}
`

const testAccDataSourceNatGatewayRuleMatchName = `
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
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.12.1.2/24"]
  }
}
resource "ionoscloud_natgateway_rule" "natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name          = "test_datasource_natgateway_rule"
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
data "ionoscloud_natgateway_rule" "test_natgateway_rule_name" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name			= "test_datasource_"
}
`
