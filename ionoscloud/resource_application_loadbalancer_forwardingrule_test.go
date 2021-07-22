package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApplicationLoadBalancerForwardingRule_Basic(t *testing.T) {
	var applicationLoadBalancerForwardingRule ionoscloud.ApplicationLoadBalancerForwardingRule
	applicationLoadBalancerForwardingRuleName := "applicationLoadBalancerForwardingRule"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApplicationLoadBalancerForwardingRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic, applicationLoadBalancerForwardingRuleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationLoadBalancerForwardingRuleExists("ionoscloud_application_loadbalancer_forwardingrule.forwarding_rule", &applicationLoadBalancerForwardingRule),
					resource.TestCheckResourceAttr("ionoscloud_application_loadbalancer_forwardingrule.forwarding_rule", "name", applicationLoadBalancerForwardingRuleName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_application_loadbalancer_forwardingrule.forwarding_rule", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckApplicationLoadBalancerForwardingRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.Attributes["application_loadbalancer_id"]
		ruleId := rs.Primary.ID

		_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, ruleId).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured and checking deletion of application loadbalancer forwarding rule %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("application loadbalancer still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckApplicationLoadBalancerForwardingRuleExists(n string, alb *ionoscloud.ApplicationLoadBalancerForwardingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckApplicationLoadBalancerForwardingRuleExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.Attributes["application_loadbalancer_id"]
		ruleId := rs.Primary.ID

		foundAlbFw, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, ruleId).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Application Loadbalancer Forwarding Rule: %s %s \n\n", rs.Primary.ID, err)
		}

		if *foundAlbFw.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		alb = &foundAlbFw

		return nil
	}
}

const testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic = `
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
 name = "%s"
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
}`

const testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate = `
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
 name = "updated"
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
}`
