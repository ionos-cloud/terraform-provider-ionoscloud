package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resourceNameAlbRule = ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleTestResource
var dataSourceNameAlbRuleById = DataSource + "." + ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleDataSourceById
var dataSourceNameAlbRuleByName = DataSource + "." + ApplicationLoadBalancerForwardingRuleResource + "." + ApplicationLoadBalancerForwardingRuleDataSourceByName

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
					resource.TestCheckResourceAttr(resourceNameAlbRule, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_ip", "10.12.118.224"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_port", "8080"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.name", "http_rule"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.type", "REDIRECT"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.location", "www.ionos.com"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.condition", "EQUALS"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.value", "value"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.name", "http_rule"),
					//resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.type", "FORWARD"),
					//resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.type", "SOURCE_IP"),
					//resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.value", "value"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "name", resourceNameAlbRule, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "protocol", resourceNameAlbRule, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "listener_ip", resourceNameAlbRule, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "listener_port", resourceNameAlbRule, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.name", resourceNameAlbRule, "http_rules.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.type", resourceNameAlbRule, "http_rules.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.location", resourceNameAlbRule, "http_rules.0.location"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.conditions.0.type", resourceNameAlbRule, "http_rules.0.conditions.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.conditions.0.condition", resourceNameAlbRule, "http_rules.0.conditions.0.condition"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.conditions.0.value", resourceNameAlbRule, "http_rules.0.conditions.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.0.name", resourceNameAlbRule, "http_rules.0.name"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.1.type", resourceNameAlbRule, "http_rules.1.type"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.1.conditions.0.type", resourceNameAlbRule, "http_rules.1.conditions.0.type"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleById, "http_rules.1.conditions.0.value", resourceNameAlbRule, "http_rules.1.conditions.0.value"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "name", resourceNameAlbRule, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "protocol", resourceNameAlbRule, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "listener_ip", resourceNameAlbRule, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "listener_port", resourceNameAlbRule, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.name", resourceNameAlbRule, "http_rules.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.type", resourceNameAlbRule, "http_rules.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.location", resourceNameAlbRule, "http_rules.0.location"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.conditions.0.type", resourceNameAlbRule, "http_rules.0.conditions.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.conditions.0.condition", resourceNameAlbRule, "http_rules.0.conditions.0.condition"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.conditions.0.value", resourceNameAlbRule, "http_rules.0.conditions.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.0.name", resourceNameAlbRule, "http_rules.0.name"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.1.type", resourceNameAlbRule, "http_rules.1.type"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.1.conditions.0.type", resourceNameAlbRule, "http_rules.1.conditions.0.type"),
					//resource.TestCheckResourceAttrPair(dataSourceNameAlbRuleByName, "http_rules.1.conditions.0.value", resourceNameAlbRule, "http_rules.1.conditions.0.value"),
				),
			},
			{
				Config:      testAccDataSourceApplicationLoadBalancerForwardingRuleWrongNameError,
				ExpectError: regexp.MustCompile("no application load balanacer forwarding rule found with the specified criteria"),
			},
			{
				Config: testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationLoadBalancerForwardingRuleExists(resourceNameAlbRule, &applicationLoadBalancerForwardingRule),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "name", ApplicationLoadBalancerForwardingRuleTestResource),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_ip", "10.12.118.224"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_port", "8080"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "client_timeout", "1000"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.name", "http_rule"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.type", "REDIRECT"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.drop_query", "true"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.location", "www.ionos.com"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.status_code", "301"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.condition", "EQUALS"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.negate", "true"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.key", "key"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.value", "value"),
				),
			},
			{
				Config: testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdateAgain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlbRule, "name", UpdatedResources),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_ip", "10.12.119.224"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "listener_port", "8081"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "client_timeout", "1500"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.name", "http_rule"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.type", "REDIRECT"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.drop_query", "true"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.location", "www.ionos.com"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.status_code", "301"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.condition", "EQUALS"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.negate", "true"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.key", "key"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.0.conditions.0.value", "value"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.name", "http_rule_2"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.type", "STATIC"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.drop_query", "false"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.response_message", "Response"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.status_code", "303"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.type", "QUERY"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.condition", "MATCHES"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.negate", "false"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.key", "key"),
					resource.TestCheckResourceAttr(resourceNameAlbRule, "http_rules.1.conditions.0.value", "value"),
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

const testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic = testAccCheckApplicationLoadBalancerConfigUpdate + `
resource ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleTestResource + ` {
 datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
 application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
 name = "` + ApplicationLoadBalancerForwardingRuleTestResource + `"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   location =  "www.ionos.com"
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     value = "value"
   }
 }
 #http_rules {
 #  name = "http_rule_2"
 #  type = "FORWARD"
 #  conditions {
 #    type = "SOURCE_IP"
 #    value = "10.12.118.224/24"
 #  }
 #}
}`

const testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdate = testAccCheckApplicationLoadBalancerConfigUpdate + `
resource ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleTestResource + ` {
 datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
 application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
 name = "` + ApplicationLoadBalancerForwardingRuleTestResource + `"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 client_timeout = 1000
 ## server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   status_code =  301
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     negate = true
     key = "key"
     value = "value"
   }
 }
}`

const testAccCheckApplicationLoadBalancerForwardingRuleConfigUpdateAgain = testAccCheckApplicationLoadBalancerConfigUpdate + `
resource ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleTestResource + ` {
 datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
 application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
 name = "` + UpdatedResources + `"
 protocol = "HTTP"
 listener_ip = "10.12.119.224"
 listener_port = 8081
 client_timeout = 1500
 ## server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   status_code =  301
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     negate = true
     key = "key"
     value = "value"
   }
 }
 http_rules {
   name = "http_rule_2"
   type = "STATIC"
   drop_query = false
   status_code = 303
   response_message = "Response"
   content_type = "text/plain"
   conditions {
     type = "QUERY"
     condition = "MATCHES"
     negate = false
     key = "key"
     value = "value"
   }
 }
}`

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchId = testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic + `
data ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
  id			= ` + ApplicationLoadBalancerForwardingRuleResource + `.` + ApplicationLoadBalancerForwardingRuleTestResource + `.id
}
`

const testAccDataSourceApplicationLoadBalancerForwardingRuleMatchName = testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic + `
data ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
  name    		= ` + ApplicationLoadBalancerForwardingRuleResource + `.` + ApplicationLoadBalancerForwardingRuleTestResource + `.name
}
`

const testAccDataSourceApplicationLoadBalancerForwardingRuleWrongNameError = testAccCheckApplicationLoadBalancerForwardingRuleConfigBasic + `
data ` + ApplicationLoadBalancerForwardingRuleResource + ` ` + ApplicationLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  application_loadbalancer_id = ` + ApplicationLoadBalancerResource + `.` + ApplicationLoadBalancerTestResource + `.id
  name    		=  "wrong_name"
}
`
