package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceApplicationLoadBalancer_matchId(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerMatchId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway", "name", "test_datasource_alb"),
				),
			},
		},
	})
}

func TestAccDataSourceApplicationLoadBalancer_matchName(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerMatchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_natgateway.test_natgateway", "name", "test_datasource_alb"),
				),
			},
		},
	})
}

const testAccDataSourceApplicationLoadBalancerCreateResources = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
}
`

const testAccDataSourceApplicationLoadBalancerMatchId = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
}

data "ionoscloud_application_loadbalancer" "test_alb" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  id			= ionoscloud_application_loadbalancer.alb.id
}
`

const testAccDataSourceApplicationLoadBalancerMatchName = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
}

data "ionoscloud_application_loadbalancer" "test_alb" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_"
}
`
