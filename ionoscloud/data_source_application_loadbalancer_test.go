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
					resource.TestCheckResourceAttr("data.ionoscloud_application_loadbalancer.test_alb", "name", "test_datasource_alb"),
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
					resource.TestCheckResourceAttr("data.ionoscloud_application_loadbalancer.test_alb", "name", "test_datasource_alb"),
				),
			},
		},
	})
}

const testAccDataSourceApplicationLoadBalancerCreateResources = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_2" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}
`

const testAccDataSourceApplicationLoadBalancerMatchId = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_2" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}

data "ionoscloud_application_loadbalancer" "test_alb" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  id			= ionoscloud_application_loadbalancer.alb.id
}
`

const testAccDataSourceApplicationLoadBalancerMatchName = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_2" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}

data "ionoscloud_application_loadbalancer" "test_alb" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "test_datasource_"
}
`
