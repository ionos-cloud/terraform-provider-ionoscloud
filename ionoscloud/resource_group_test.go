package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccGroup_Basic(t *testing.T) {
	var group profitbricks.Group
	groupName := "terraform test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckGroupConfig_basic, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists("ionoscloud_group.group", &group),
					testAccCheckGroupAttributes("ionoscloud_group.group", groupName),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", groupName),
				),
			},
			{
				Config: testAccCheckGroupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupAttributes("ionoscloud_group.group", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		_, err := client.GetGroup(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("group still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Group %s %s", rs.Primary.ID, err)
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
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckGroupExists(n string, group *profitbricks.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundgroup, err := client.GetGroup(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Group: %s", rs.Primary.ID)
		}
		if foundgroup.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		group = foundgroup

		return nil
	}
}

const testAccCheckGroupConfig_basic = `
resource "ionoscloud_group" "group" {
  name = "%s"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}`

const testAccCheckGroupConfig_update = `
resource "ionoscloud_group" "group" {
  name = "updated"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
}
`
