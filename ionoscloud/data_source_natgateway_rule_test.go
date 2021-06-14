package ionoscloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNatGatewayRule_matchId(t *testing.T) {

	publicIp1 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_1")
	if publicIp1 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_IP_1 not set; please set it to a valid public IP for the de/fra zone")
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleCreateResources, publicIp1, publicIp1),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchId, publicIp1, publicIp1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway_rule.test_natgateway_rule", "name", "test_datasource_natgateway_rule"),
				),
			},
		},
	})
}

func TestAccDataSourceNatGatewayRule_matchName(t *testing.T) {

	publicIp1 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_1")
	if publicIp1 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_IP_1 not set; please set it to a valid public IP for the de/fra zone")
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleCreateResources, publicIp1, publicIp1),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayRuleMatchName, publicIp1, publicIp1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway_rule.test_natgateway_rule", "name", "test_datasource_natgateway_rule"),
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
resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  public        = false
  name          = "test_natgateway_lan"
}
resource "ionoscloud_natgateway" "natgateway" {
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ "%s" ]
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
  public_ip     = "%s"
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
resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  public        = false
  name          = "test_natgateway_lan"
}
resource "ionoscloud_natgateway" "natgateway" {
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ "%s" ]
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
  public_ip     = "%s"
  target_subnet = "10.0.1.0/24"
  target_port_range {
      start = 500
      end   = 1000
  }
}

data "ionoscloud_natgateway_rule" "test_natgateway_rule" {
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
resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  public        = false
  name          = "test_natgateway_lan"
}
resource "ionoscloud_natgateway" "natgateway" {
  depends_on    = [ ionoscloud_lan.natgateway_lan ]
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "natgateway_test"
  public_ips    = [ "%s" ]
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
  public_ip     = "%s"
  target_subnet = "10.0.1.0/24"
  target_port_range {
      start = 500
      end   = 1000
  }
}

data "ionoscloud_natgateway_rule" "test_natgateway_rule" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  natgateway_id = ionoscloud_natgateway.natgateway.id
  name			= "test_datasource_"
}
`
