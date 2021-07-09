package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTargetGroup_Basic(t *testing.T) {
	var targetGroup ionoscloud.TargetGroup
	targetGroupName := "targetGroup"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTargetGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckTargetGroupConfigBasic, targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists("ionoscloud_target_group.target_group", &targetGroup),
					resource.TestCheckResourceAttr("ionoscloud_target_group.target_group", "name", targetGroupName),
				),
			},
			{
				Config: testAccCheckTargetGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_target_group.target_group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckTargetGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.ID).Execute()

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
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
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

		foundTargetGroup, _, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, rs.Primary.ID).Execute()

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
resource "ionoscloud_target_group" "target_group" {
 name = "%s"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
 health_check {
     connect_timeout = 1000
 }
 http_health_check {
     path = "/monitoring"
     match_type = "RESPONSE_BODY"
     response = "200"
   }
}
`

const testAccCheckTargetGroupConfigUpdate = `
resource "ionoscloud_target_group" "target_group" {
 name = "updated"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
 health_check {
     connect_timeout = 1000
 }
 http_health_check {
     path = "/monitoring"
     match_type = "RESPONSE_BODY"
     response = "200"
   }
}`
