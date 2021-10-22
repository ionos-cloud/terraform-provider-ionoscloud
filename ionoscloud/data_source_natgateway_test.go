package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNatGateway_matchId(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway_id", "name", "test_datasource_natgateway"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_natgateway.test_natgateway_id", "public_ips.0", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway_name", "name", "test_datasource_natgateway"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_natgateway.test_natgateway_name", "public_ips.0", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
		},
	})
}

const testAccDataSourceNatGatewayCreateResources = `
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
  name          = "test_datasource_natgateway" 
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}
`

const testAccDataSourceNatGatewayMatchId = `
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
  name          = "test_datasource_natgateway" 
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}

data "ionoscloud_natgateway" "test_natgateway_id" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  id			= ionoscloud_natgateway.natgateway.id
}
`

const testAccDataSourceNatGatewayMatchName = `
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
  name          = "test_datasource_natgateway" 
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}

data "ionoscloud_natgateway" "test_natgateway_name" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name			= "test_datasource_"
}
`
