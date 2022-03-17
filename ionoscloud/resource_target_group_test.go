package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resourceNameTargetGroup = TargetGroupResource + "." + TargetGroupTestResource

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
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "name", TargetGroupTestResource),
				),
			},
			{
				Config: testAccCheckTargetGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTargetGroup, "name", UpdatedResources),
				),
			},
		},
	})
}

func testAccCheckTargetGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ApplicationLoadBalancerResource {
			continue
		}

		apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured at checking deletion of forwarding rule %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("network loadbalancer forwarding rule still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckTargetGroupExists(n string, targetGroup *ionoscloud.TargetGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
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
			return fmt.Errorf("error occured while fetching TargetGroup: %s", rs.Primary.ID)
		}
		if *foundTargetGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		targetGroup = &foundTargetGroup

		return nil
	}
}

const testAccCheckTargetGroupConfigBasic = `
resource ` + TargetGroupResource + ` ` + TargetGroupTestResource + ` {
 name = "` + TargetGroupTestResource + `"
 algorithm = "ROUND_ROBIN"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "1"
   health_check {
     check = true
     check_interval = 2000
	 maintenance = true
   }
 }
 health_check {
     connect_timeout = 5000
     target_timeout = 50000
     retries = 2
 }
 http_health_check {
     path = "/."
     method = "GET"
     match_type = "STATUS_CODE"
     response = "200"
     regex = false
     negate = false
   }
}
`

const testAccCheckTargetGroupConfigUpdate = `
resource ` + TargetGroupResource + ` ` + TargetGroupTestResource + ` {
 name = "` + UpdatedResources + `"
 algorithm = "RANDOM"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.3"
   port = "8081"
   weight = "124"
   health_check {
     check = true
     check_interval = 2500
	 maintenance = true
   }
 }
 health_check {
     connect_timeout = 5500
     target_timeout = 55000
     retries = 3
 }
 http_health_check {
     path = "/."
     method = "GET"
     match_type = "STATUS_CODE"
     response = "200"
     regex = false
     negate = false
   }
}`
