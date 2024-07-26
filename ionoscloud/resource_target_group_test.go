//go:build all || alb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceNameTargetGroup = constant.TargetGroupResource + "." + constant.TargetGroupTestResource
const resourceNameTargetGroupById = constant.DataSource + "." + constant.TargetGroupResource + "." + constant.TargetGroupDataSourceById
const resourceNameTargetGroupByName = constant.DataSource + "." + constant.TargetGroupResource + "." + constant.TargetGroupDataSourceByName

func TestAccTargetGroupBasic(t *testing.T) {
	var targetGroup ionoscloud.TargetGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTargetGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTargetGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(resourceNameTargetGroup, &targetGroup),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "name", constant.TargetGroupTestResource),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "algorithm", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.ip", "22.231.2.2"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.proxy_protocol", "none"),
				),
			},
			{
				Config: testAccCheckTargetGroupConfigUpdateWithAllParameters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "algorithm", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol_version", "HTTP2"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.ip", "22.231.2.2"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.proxy_protocol", "v2"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.health_check_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.maintenance_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.ip", "22.232.2.3"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.port", "8081"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.weight", "124"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.proxy_protocol", "v1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.health_check_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.1.maintenance_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.check_timeout", "5000"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.check_interval", "50000"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.retries", "2"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.path", "/."),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.method", "GET"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.match_type", "STATUS_CODE"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.response", "200"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.regex", "true"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.negate", "true"),
				),
			},
			{
				Config: testAccDataSourceTargetGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "name", resourceNameTargetGroup, "name"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "protocol_version", resourceNameTargetGroup, "protocol_version"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "algorithm", resourceNameTargetGroup, "algorithm"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.ip", resourceNameTargetGroup, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.port", resourceNameTargetGroup, "targets.0.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.weight", resourceNameTargetGroup, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.proxy_protocol", resourceNameTargetGroup, "targets.0.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.health_check_enabled", resourceNameTargetGroup, "targets.0.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.0.maintenance_enabled", resourceNameTargetGroup, "targets.0.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.ip", resourceNameTargetGroup, "targets.1.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.port", resourceNameTargetGroup, "targets.1.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.weight", resourceNameTargetGroup, "targets.1.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.proxy_protocol", resourceNameTargetGroup, "targets.1.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.health_check_enabled", resourceNameTargetGroup, "targets.1.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "targets.1.maintenance_enabled", resourceNameTargetGroup, "targets.1.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "health_check.0.check_timeout", resourceNameTargetGroup, "health_check.0.check_timeout"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "health_check.0.check_interval", resourceNameTargetGroup, "health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "health_check.0.retries", resourceNameTargetGroup, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.path", resourceNameTargetGroup, "http_health_check.0.path"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.method", resourceNameTargetGroup, "http_health_check.0.method"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.match_type", resourceNameTargetGroup, "http_health_check.0.match_type"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.response", resourceNameTargetGroup, "http_health_check.0.response"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.regex", resourceNameTargetGroup, "http_health_check.0.regex"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupById, "http_health_check.0.negate", resourceNameTargetGroup, "http_health_check.0.negate"),
				),
			},
			{
				Config: testAccDataSourceTargetGroupPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "name", resourceNameTargetGroup, "name"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "algorithm", resourceNameTargetGroup, "algorithm"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.ip", resourceNameTargetGroup, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.port", resourceNameTargetGroup, "targets.0.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.weight", resourceNameTargetGroup, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.proxy_protocol", resourceNameTargetGroup, "targets.0.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.health_check_enabled", resourceNameTargetGroup, "targets.0.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.maintenance_enabled", resourceNameTargetGroup, "targets.0.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.ip", resourceNameTargetGroup, "targets.1.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.port", resourceNameTargetGroup, "targets.1.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.weight", resourceNameTargetGroup, "targets.1.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.proxy_protocol", resourceNameTargetGroup, "targets.1.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.health_check_enabled", resourceNameTargetGroup, "targets.1.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.maintenance_enabled", resourceNameTargetGroup, "targets.1.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.check_timeout", resourceNameTargetGroup, "health_check.0.check_timeout"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.check_interval", resourceNameTargetGroup, "health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.retries", resourceNameTargetGroup, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.path", resourceNameTargetGroup, "http_health_check.0.path"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.method", resourceNameTargetGroup, "http_health_check.0.method"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.match_type", resourceNameTargetGroup, "http_health_check.0.match_type"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.response", resourceNameTargetGroup, "http_health_check.0.response"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.regex", resourceNameTargetGroup, "http_health_check.0.regex"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.negate", resourceNameTargetGroup, "http_health_check.0.negate"),
				),
			},
			{
				Config: testAccDataSourceTargetGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "name", resourceNameTargetGroup, "name"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "algorithm", resourceNameTargetGroup, "algorithm"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.ip", resourceNameTargetGroup, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.port", resourceNameTargetGroup, "targets.0.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.weight", resourceNameTargetGroup, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.proxy_protocol", resourceNameTargetGroup, "targets.0.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.health_check_enabled", resourceNameTargetGroup, "targets.0.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.0.maintenance_enabled", resourceNameTargetGroup, "targets.0.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.ip", resourceNameTargetGroup, "targets.1.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.port", resourceNameTargetGroup, "targets.1.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.weight", resourceNameTargetGroup, "targets.1.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.proxy_protocol", resourceNameTargetGroup, "targets.1.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "targets.1.health_check_enabled", resourceNameTargetGroup, "targets.1.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.check_timeout", resourceNameTargetGroup, "health_check.0.check_timeout"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.check_interval", resourceNameTargetGroup, "health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "health_check.0.retries", resourceNameTargetGroup, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.path", resourceNameTargetGroup, "http_health_check.0.path"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.method", resourceNameTargetGroup, "http_health_check.0.method"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.match_type", resourceNameTargetGroup, "http_health_check.0.match_type"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.response", resourceNameTargetGroup, "http_health_check.0.response"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.regex", resourceNameTargetGroup, "http_health_check.0.regex"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByName, "http_health_check.0.negate", resourceNameTargetGroup, "http_health_check.0.negate"),
				),
			},
			{
				Config:      testAccDataSourceTargetGroupWrongNameError,
				ExpectError: regexp.MustCompile("no target group found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceTargetGroupWrongPartialNameError,
				ExpectError: regexp.MustCompile("no target group found with the specified criteria"),
			},
			{
				Config: testAccCheckTargetGroupConfigUpdateAgain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "algorithm", "RANDOM"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.ip", "22.232.2.3"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.port", "8081"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.weight", "124"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.proxy_protocol", "v1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.health_check_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "targets.0.maintenance_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceNameTargetGroup, "targets.1"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.check_timeout", "5500"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.check_interval", "55000"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "health_check.0.retries", "3"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.path", "../."),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.method", "POST"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.match_type", "RESPONSE_BODY"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.response", "Response"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.regex", "false"),
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "http_health_check.0.negate", "false"),
				),
			},
		},
	})
}

func testAccCheckTargetGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ALBResource {
			continue
		}

		apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred at checking deletion of forwarding rule %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("network loadbalancer forwarding rule still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTargetGroupExists(n string, targetGroup *ionoscloud.TargetGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckTargetGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundTargetGroup, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching TargetGroup: %s, %w", rs.Primary.ID, err)
		}
		if *foundTargetGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		targetGroup = &foundTargetGroup

		return nil
	}
}

const testAccCheckTargetGroupConfigBasic = `
resource ` + constant.TargetGroupResource + ` ` + constant.TargetGroupTestResource + ` {
 name = "` + constant.TargetGroupTestResource + `"
 algorithm = "ROUND_ROBIN"
 protocol = "HTTP"
 protocol_version = "HTTP1"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "1"
  }
}
`

const testAccCheckTargetGroupConfigUpdateWithAllParameters = `
resource ` + constant.TargetGroupResource + ` ` + constant.TargetGroupTestResource + ` {
 name = "` + constant.UpdatedResources + `"
 algorithm = "ROUND_ROBIN"
 protocol = "HTTP"
 protocol_version = "HTTP2"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "1"
   proxy_protocol = "v2"
   health_check_enabled = true
   maintenance_enabled = true
 }
 targets {
	ip = "22.232.2.3"
	port = "8081"
	weight = "124"
	proxy_protocol = "v1"
	health_check_enabled = false
	maintenance_enabled = false
  }
 health_check {
     check_timeout = 5000
     check_interval = 50000
     retries = 2
 }
 http_health_check {
     path = "/."
     method = "GET"
     match_type = "STATUS_CODE"
     response = "200"
     regex = true
     negate = true
   }
}
`

const testAccCheckTargetGroupConfigUpdateAgain = `
resource ` + constant.TargetGroupResource + ` ` + constant.TargetGroupTestResource + ` {
 name = "` + constant.UpdatedResources + `"
 algorithm = "RANDOM"
 protocol = "HTTP"
 protocol_version = "HTTP1"
 targets {
   ip = "22.232.2.3"
   port = "8081"
   weight = "124"
   proxy_protocol = "v1"
   health_check_enabled = false
   maintenance_enabled = false
 }
 health_check {
     check_timeout = 5500
     check_interval = 55000
     retries = 3
 }
 http_health_check {
     path = "../."
     method = "POST"
     match_type = "RESPONSE_BODY"
     response = "Response"
     regex = false
     negate = false
   }
}`

const testAccDataSourceTargetGroupMatchId = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
data ` + constant.TargetGroupResource + ` ` + constant.TargetGroupDataSourceById + ` {
  id			= ` + resourceNameTargetGroup + `.id
}
`

const testAccDataSourceTargetGroupMatchName = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
data ` + constant.TargetGroupResource + ` ` + constant.TargetGroupDataSourceByName + ` {
  name			= ` + resourceNameTargetGroup + `.name
}
`

const testAccDataSourceTargetGroupWrongNameError = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
data ` + constant.TargetGroupResource + ` ` + constant.TargetGroupDataSourceByName + ` {
  name			= "wrong name"
}
`

const testAccDataSourceTargetGroupPartialMatchName = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
data ` + constant.TargetGroupResource + ` ` + constant.TargetGroupDataSourceByName + ` {
  name          = "` + constant.DataSourcePartial + `"
  partial_match = true
}
`

const testAccDataSourceTargetGroupWrongPartialNameError = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
data ` + constant.TargetGroupResource + ` ` + constant.TargetGroupDataSourceByName + ` {
  name			= "wrong name"
  partial_match = true
}
`
