package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkLoadBalancer_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkLoadBalancerCreateResources,
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_networkloadbalancer.test_networkloadbalancer", "name", "test_datasource_nlb"),
				),
			},
		},
	})
}

func TestAccDataSourceNetworkLoadBalancer_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkLoadBalancerCreateResources,
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_networkloadbalancer.test_networkloadbalancer", "name", "test_datasource_nlb"),
				),
			},
		},
	})
}

const testAccDataSourceNetworkLoadBalancerCreateResources = `
resource "ionoscloud_datacenter" "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_lan_1" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_1"
}

resource "ionoscloud_lan" "nlb_lan_2" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_2"
}


resource "ionoscloud_networkloadbalancer" "networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name          = "test_datasource_nlb"
  listener_lan  = ionoscloud_lan.nlb_lan_1.id
  target_lan    = ionoscloud_lan.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}
`

const testAccDataSourceNetworkLoadBalancerMatchId = `
resource "ionoscloud_datacenter" "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_lan_1" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_1"
}

resource "ionoscloud_lan" "nlb_lan_2" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_2"
}


resource "ionoscloud_networkloadbalancer" "networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name          = "test_datasource_nlb"
  listener_lan  = ionoscloud_lan.nlb_lan_1.id
  target_lan    = ionoscloud_lan.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}

data "ionoscloud_networkloadbalancer" "test_networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  id			= ionoscloud_networkloadbalancer.networkloadbalancer.id
}
`

const testAccDataSourceNetworkLoadBalancerMatchName = `
resource "ionoscloud_datacenter" "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_lan_1" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_1"
}

resource "ionoscloud_lan" "nlb_lan_2" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_2"
}


resource "ionoscloud_networkloadbalancer" "networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name          = "test_datasource_nlb"
  listener_lan  = ionoscloud_lan.nlb_lan_1.id
  target_lan    = ionoscloud_lan.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}

data "ionoscloud_networkloadbalancer" "test_networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name			= "test_datasource_"
}
`
