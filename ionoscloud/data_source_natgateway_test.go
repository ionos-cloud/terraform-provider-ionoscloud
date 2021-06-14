package ionoscloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNatGateway_matchId(t *testing.T) {

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
				Config: fmt.Sprintf(testAccDataSourceNatGatewayCreateResources, publicIp1),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchId, publicIp1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway", "name", "test_datasource_natgateway"),
				),
			},
		},
	})
}

func TestAccDataSourceNatGateway_matchName(t *testing.T) {

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
				Config: fmt.Sprintf(testAccDataSourceNatGatewayCreateResources, publicIp1),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchName, publicIp1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway", "name", "test_datasource_natgateway"),
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

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" { 
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "test_datasource_natgateway" 
  public_ips    = [ "%s"]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}
`

const testAccDataSourceNatGatewayMatchId = `
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_datasource_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" { 
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "test_datasource_natgateway" 
  public_ips    = [ "%s" ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}

data "ionoscloud_natgateway" "test_natgateway" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  id			= ionoscloud_natgateway.natgateway.id
}
`

const testAccDataSourceNatGatewayMatchName = `
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_datasource_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" { 
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "test_datasource_natgateway" 
  public_ips    = [ "%s" ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
     gateway_ips = [ "10.11.2.5/32"] 
  }
}

data "ionoscloud_natgateway" "test_natgateway" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name			= "test_datasource_"
}
`
