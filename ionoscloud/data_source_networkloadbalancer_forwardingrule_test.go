// +build nlb

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkLoadBalancerForwardingRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleCreateResources,
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_networkloadbalancer_forwardingrule.test_nlb_fr_id", "name", "test_datasource_nlb_fr"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_networkloadbalancer_forwardingrule.test_nlb_fr_name", "name", "test_datasource_nlb_fr"),
				),
			},
		},
	})
}

const testAccDataSourceNetworkLoadBalancerForwardingRuleCreateResources = `
resource "ionoscloud_datacenter" "nlb_fr_datacenter" {
  name              = "test_nlb_fr"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_fr_lan_1" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}

resource "ionoscloud_lan" "nlb_fr_lan_2" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}


resource "ionoscloud_networkloadbalancer" "test_nbl_fr" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  name          = "test_nlb_fr"
  listener_lan  = ionoscloud_lan.nlb_fr_lan_1.id
  target_lan    = ionoscloud_lan.nlb_fr_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}

resource "ionoscloud_networkloadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
 networkloadbalancer_id = ionoscloud_networkloadbalancer.test_nbl_fr.id
 name = "test_datasource_nlb_fr"
 algorithm = "SOURCE_IP"
 protocol = "TCP"
 listener_ip = "10.12.118.224"
 listener_port = "8081"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
}
`

const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId = `
resource "ionoscloud_datacenter" "nlb_fr_datacenter" {
  name              = "test_nlb_fr"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_fr_lan_1" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}

resource "ionoscloud_lan" "nlb_fr_lan_2" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}


resource "ionoscloud_networkloadbalancer" "test_nbl_fr" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  name          = "test_nlb_fr"
  listener_lan  = ionoscloud_lan.nlb_fr_lan_1.id
  target_lan    = ionoscloud_lan.nlb_fr_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}

resource "ionoscloud_networkloadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
 networkloadbalancer_id = ionoscloud_networkloadbalancer.test_nbl_fr.id
 name = "test_datasource_nlb_fr"
 algorithm = "SOURCE_IP"
 protocol = "TCP"
 listener_ip = "10.12.118.224"
 listener_port = "8081"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
}

data "ionoscloud_networkloadbalancer_forwardingrule" "test_nlb_fr_id" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  networkloadbalancer_id  = ionoscloud_networkloadbalancer.test_nbl_fr.id
  id			= ionoscloud_networkloadbalancer_forwardingrule.forwarding_rule.id
}
`

const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName = `
resource "ionoscloud_datacenter" "nlb_fr_datacenter" {
  name              = "test_nlb_fr"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_fr_lan_1" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}

resource "ionoscloud_lan" "nlb_fr_lan_2" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  public        = false
  name          = "test_nlb_fr_lan_1"
}


resource "ionoscloud_networkloadbalancer" "test_nbl_fr" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  name          = "test_nlb_fr"
  listener_lan  = ionoscloud_lan.nlb_fr_lan_1.id
  target_lan    = ionoscloud_lan.nlb_fr_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}

resource "ionoscloud_networkloadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
 networkloadbalancer_id = ionoscloud_networkloadbalancer.test_nbl_fr.id
 name = "test_datasource_nlb_fr"
 algorithm = "SOURCE_IP"
 protocol = "TCP"
 listener_ip = "10.12.118.224"
 listener_port = "8081"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
}

data "ionoscloud_networkloadbalancer_forwardingrule" "test_nlb_fr_name" {
  datacenter_id = ionoscloud_datacenter.nlb_fr_datacenter.id
  networkloadbalancer_id  = ionoscloud_networkloadbalancer.test_nbl_fr.id
  name			= "test_datasource_"
}
`
