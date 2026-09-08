//go:build all || alb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const resourceNameTargetGroup = constant.TargetGroupResource + "." + constant.TargetGroupTestResource
const resourceNameTargetGroupByID = constant.DataSource + "." + constant.TargetGroupResource + "." + constant.TargetGroupDataSourceById
const resourceNameTargetGroupByName = constant.DataSource + "." + constant.TargetGroupResource + "." + constant.TargetGroupDataSourceByName

func TestAccTargetGroupBasic(t *testing.T) {
	var targetGroup ionoscloud.TargetGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckTargetGroupDestroyCheck,
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
				Config: testAccDataSourceTargetGroupMatchID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "name", resourceNameTargetGroup, "name"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "protocol_version", resourceNameTargetGroup, "protocol_version"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "algorithm", resourceNameTargetGroup, "algorithm"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.ip", resourceNameTargetGroup, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.port", resourceNameTargetGroup, "targets.0.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.weight", resourceNameTargetGroup, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.proxy_protocol", resourceNameTargetGroup, "targets.0.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.health_check_enabled", resourceNameTargetGroup, "targets.0.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.0.maintenance_enabled", resourceNameTargetGroup, "targets.0.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.ip", resourceNameTargetGroup, "targets.1.ip"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.port", resourceNameTargetGroup, "targets.1.port"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.weight", resourceNameTargetGroup, "targets.1.weight"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.proxy_protocol", resourceNameTargetGroup, "targets.1.proxy_protocol"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.health_check_enabled", resourceNameTargetGroup, "targets.1.health_check_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "targets.1.maintenance_enabled", resourceNameTargetGroup, "targets.1.maintenance_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "health_check.0.check_timeout", resourceNameTargetGroup, "health_check.0.check_timeout"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "health_check.0.check_interval", resourceNameTargetGroup, "health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "health_check.0.retries", resourceNameTargetGroup, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.path", resourceNameTargetGroup, "http_health_check.0.path"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.method", resourceNameTargetGroup, "http_health_check.0.method"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.match_type", resourceNameTargetGroup, "http_health_check.0.match_type"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.response", resourceNameTargetGroup, "http_health_check.0.response"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.regex", resourceNameTargetGroup, "http_health_check.0.regex"),
					resource.TestCheckResourceAttrPair(resourceNameTargetGroupByID, "http_health_check.0.negate", resourceNameTargetGroup, "http_health_check.0.negate"),
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

// TestAccTargetGroupQuery covers the ionoscloud_target_group list resource and the
// resource identity it streams, end to end through `terraform query`.
//
// The list resource is served by the plugin-framework half of the provider even though
// the target group resource itself is implemented with SDKv2, so this also covers the mux
// serving the two halves under the same type name. See
// internal/framework/services/compute/resource_target_group_list.go.
func TestAccTargetGroupQuery(t *testing.T) {
	// Its own name and label, NOT constant.TargetGroupTestResource: querycheck.ExpectLength
	// asserts a contract-wide total, so it needs a name no other test in this suite creates.
	// TestAccTargetGroupBasic uses the shared fixture name, and a concurrent or leftover
	// copy of it would make an ExpectLength(1) assertion here flap.
	const (
		targetGroupQueryName  = "tf-test-target-group-query"
		targetGroupQueryLabel = "test_target_group_query"
		targetGroupAddr       = constant.TargetGroupResource + "." + targetGroupQueryLabel
	)

	queryConfig := fmt.Sprintf(`
resource %[1]q %[2]q {
  name             = %[3]q
  algorithm        = "ROUND_ROBIN"
  protocol         = "HTTP"
  protocol_version = "HTTP1"
  targets {
    ip     = "22.231.2.2"
    port   = 8080
    weight = 1
  }
}`, constant.TargetGroupResource, targetGroupQueryLabel, targetGroupQueryName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckTargetGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: queryConfig,
			},
			// List without filters: the target group must show up with its identity. The
			// identity is a lone `id` - target groups have no location.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q %[2]q {
  provider = ionoscloud
}`, constant.TargetGroupResource, targetGroupQueryLabel),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(targetGroupAddr, map[string]knownvalue.Check{
						"id": knownvalue.NotNull(),
					}),
				},
			},
			// Both filter fields the list resource advertises, ANDed. The unique name
			// guarantees exactly one result.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q %[2]q {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",      field_value = %[3]q },
      { field_name = "algorithm", field_value = "ROUND_ROBIN" },
    ]
  }
}`, constant.TargetGroupResource, targetGroupQueryLabel, targetGroupQueryName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(targetGroupAddr, 1),
				},
			},
			// The same name with an algorithm the fixture does not use: proves the
			// algorithm filter is evaluated rather than ignored. (protocol is not a filter
			// field - the API allows only "HTTP", so it could never narrow anything.)
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q %[2]q {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",      field_value = %[3]q },
      { field_name = "algorithm", field_value = "SOURCE_IP" },
    ]
  }
}`, constant.TargetGroupResource, targetGroupQueryLabel, targetGroupQueryName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(targetGroupAddr, 0),
				},
			},
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    targetGroupAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCheckTargetGroupDestroyCheck(s *terraform.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
	if err != nil {
		return err
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

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
		if err != nil {
			return err
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

const testAccDataSourceTargetGroupMatchID = testAccCheckTargetGroupConfigUpdateWithAllParameters + `
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
