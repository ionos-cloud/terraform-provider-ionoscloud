package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resourceNameAlbRule = ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleTestResource

func TestAccApplicationLoadBalancerForwardingRuleBasic(t *testing.T) {
	var applicationLoadBalancerForwardingRule ionoscloud.ApplicationLoadBalancerForwardingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApplicationLoadBalancerForwardingRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationLoadBalancerForwardingRuleExists(resourceNameAlbRule, &applicationLoadBalancerForwardingRule),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "name", ApplicationLoadBalancerForwardingRuleTestResource),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlbRule, "name", UpdatedResources),
				),
			},
		},
	})
}

func testAccCheckApplicationLoadBalancerForwardingRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ApplicationLoadBalancerForwardingRuleResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.Attributes["application_loadbalancer_id"]
		ruleId := rs.Primary.ID

		_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, ruleId).Execute()
		logApiRequestTime(apiResponse)

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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
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

		foundAlbFw, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, ruleId).Execute()
		logApiRequestTime(apiResponse)

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

const testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic = testAccCheckApplicationLoadBalancerConfigBasic + `
resource ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleTestResource + ` {
 datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
 application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
 name = "` + ApplicationLoadBalancerForwardingRuleTestResource + `"
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

const testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate = testAccCheckApplicationLoadBalancerConfigBasic + `
resource ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleTestResource + ` {
 datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
 application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
 name = "` + UpdatedResources + `"
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
