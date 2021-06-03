package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup_Basic(t *testing.T) {
	var group ionoscloud.Group
	groupName := "terraform test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckgroupconfigBasic, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists("ionoscloud_group.group", &group),
					testAccCheckGroupAttributes("ionoscloud_group.group", groupName),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", groupName),
				),
			},
			{
				Config: testacccheckgroupconfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupAttributes("ionoscloud_group.group", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type != "ionoscloud_group" {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse != nil && apiResponse.StatusCode != 404 {
				payload := fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				return fmt.Errorf("group still exists %s - an error occurred while checking it %s %s", rs.Primary.ID, err, payload)
			}
		} else {
			return fmt.Errorf("group still exists %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckGroupAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckGroupAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckGroupExists(n string, group *ionoscloud.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundGroup, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching Group: %s %s", rs.Primary.ID, payload)
		}
		if *foundGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		group = &foundGroup

		return nil
	}
}

const testacccheckgroupconfigBasic = `
resource "ionoscloud_group" "group" {
  name = "%s"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}
`

const testacccheckgroupconfigUpdate = `
resource "ionoscloud_group" "group" {
  name = "updated"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
}
`
