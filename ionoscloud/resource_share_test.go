package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccShareBasic(t *testing.T) {
	var share ionoscloud.GroupShare
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckShareDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckshareconfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareExists(ShareResourceFullName, &share),
					resource.TestCheckResourceAttr(ShareResourceFullName, "edit_privilege", "true"),
					resource.TestCheckResourceAttr(ShareResourceFullName, "share_privilege", "true"),
				),
			},
			{
				Config: testacccheckshareconfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ShareResourceFullName, "edit_privilege", "false"),
					resource.TestCheckResourceAttr(ShareResourceFullName, "share_privilege", "false"),
				),
			},
		},
	})
}

func testAccCheckShareDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type != ShareResource {
			continue
		}

		grpId := rs.Primary.Attributes["group_id"]
		resourceId := rs.Primary.Attributes["resource_id"]

		_, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, resourceId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of resource %s in group %s: %s", resourceId, grpId, err)
			}
		} else {
			return fmt.Errorf("share for resource %s still exists in group %s", resourceId, grpId)
		}

	}

	return nil
}

func testAccCheckShareExists(n string, share *ionoscloud.GroupShare) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckShareExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		grpId := rs.Primary.Attributes["group_id"]
		resourceId := rs.Primary.Attributes["resource_id"]
		foundShare, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(context.TODO(), grpId, resourceId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Share of resource  %s in group %s", rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["group_id"])
		}
		if *foundShare.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		share = &foundShare

		return nil
	}
}

const testacccheckshareconfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "terraform test"
	location = "us/las"
}

resource "ionoscloud_group" "group" {
  name = "terraform test"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "ionoscloud_share" "share" {
  group_id = "${ionoscloud_group.group.id}"
  resource_id = "${ionoscloud_datacenter.foobar.id}"
  edit_privilege = true
  share_privilege = true
}`

const testacccheckshareconfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "terraform test"
	location = "us/las"
}

resource "ionoscloud_group" "group" {
  name = "terraform test"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "ionoscloud_share" "share" {
  group_id = "${ionoscloud_group.group.id}"
  resource_id = "${ionoscloud_datacenter.foobar.id}"
  edit_privilege = false
  share_privilege = false
}
`
