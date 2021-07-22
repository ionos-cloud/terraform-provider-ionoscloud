package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceApplicationLoadBalancerForwardingRule_matchId(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerForwardingRuleCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_application_loadbalancer_forwardingrule.test_forwarding_rule", "name", "test_datasource_forwarding_rule"),
				),
			},
		},
	})
}

func TestAccDataSourceApplicationLoadBalancerForwardingRule_matchName(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerForwardingRuleCreateResources),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_application_loadbalancer_forwardingrule.test_forwarding_rule", "name", "test_datasource_forwarding_rule"),
				),
			},
		},
	})
}

const testAccDataSourceApplicationLoadBalancerForwardingRuleCreateResources = `
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
  name          = "alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}

resource "ionoscloud_application_loadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.alb_datacenter.id
 application_loadbalancer_id = ionoscloud_application_loadbalancer.alb.id
 name = "test_datasource_forwarding_rule"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 health_check {
     client_timeout = 1000
 }
 ## server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     value = "something"
   }
 }
}
`

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId = `
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
  name          = "alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}

resource "ionoscloud_application_loadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.alb_datacenter.id
 application_loadbalancer_id = ionoscloud_application_loadbalancer.alb.id
 name = "test_datasource_forwarding_rule"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 health_check {
     client_timeout = 1000
 }
 ## server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     value = "something"
   }
 }
}

data "ionoscloud_application_loadbalancer_forwardingrule" "test_forwarding_rule" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.alb.id
  id			= ionoscloud_application_loadbalancer_forwardingrule.forwarding_rule.id
}
`

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName = `
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
  name          = "alb"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}

resource "ionoscloud_application_loadbalancer_forwardingrule" "forwarding_rule" {
 datacenter_id = ionoscloud_datacenter.alb_datacenter.id
 application_loadbalancer_id = ionoscloud_application_loadbalancer.alb.id
 name = "test_datasource_forwarding_rule"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 health_check {
     client_timeout = 1000
 }
 ## server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     value = "something"
   }
 }
}

data "ionoscloud_application_loadbalancer_forwardingrule" "test_forwarding_rule" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.alb.id
  name    		= "test_datasource"
}
`
