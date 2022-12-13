//go:build compute || all || share

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Config: fmt.Sprintf(testAccCheckShareConfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareExists(shareResourceFullName, &share),
					resource.TestCheckResourceAttr(shareResourceFullName, "edit_privilege", "true"),
					resource.TestCheckResourceAttr(shareResourceFullName, "share_privilege", "true"),
				),
			},
			{
				Config: testAccDataSourceShareConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(shareResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(shareResourceFullName, "id", DataSource+"."+ShareResource+"."+sourceShareName, "id"),
					resource.TestCheckResourceAttrPair(shareResourceFullName, "edit_privilege",
						DataSource+"."+ShareResource+"."+sourceShareName, "edit_privilege"),
					resource.TestCheckResourceAttrPair(shareResourceFullName, "share_privilege",
						DataSource+"."+ShareResource+"."+sourceShareName, "share_privilege"),
					resource.TestCheckResourceAttr(DataSource+"."+ShareResource+"."+sourceShareName, "edit_privilege", "true"),
				),
			},
			{
				Config: testAccCheckShareConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(shareResourceFullName, "edit_privilege", "false"),
					resource.TestCheckResourceAttr(shareResourceFullName, "share_privilege", "false"),
				),
			},
		},
	})
}

func testAccCheckShareDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of resource %s in group %s: %w", resourceId, grpId, err)
			}
		} else {
			return fmt.Errorf("share for resource %s still exists in group %s", resourceId, grpId)
		}

	}

	return nil
}

func testAccCheckShareExists(n string, share *ionoscloud.GroupShare) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckShareExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		grpId := rs.Primary.Attributes["group_id"]
		resourceId := rs.Primary.Attributes["resource_id"]
		foundshare, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(context.TODO(), grpId, resourceId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Share of resource  %s in group %s", rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["group_id"])
		}
		if *foundshare.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		share = &foundshare

		return nil
	}
}

const testAccCheckShareConfigBasic = `
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

const testAccCheckShareConfigUpdate = `
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

var testAccDataSourceShareConfigBasic = testAccCheckShareConfigBasic + `
data ` + ShareResource + " " + sourceShareName + `{
  group_id    = "${ionoscloud_group.group.id}"
  resource_id = "${ionoscloud_datacenter.foobar.id}"
  id		  = ` + shareResourceFullName + `.id
}
`
