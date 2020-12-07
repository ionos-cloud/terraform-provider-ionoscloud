package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIPBlock_Basic(t *testing.T) {
	var ipblock profitbricks.IPBlock
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIPBlockConfig_basic, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "location", location),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckIPBlockConfig_update, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckIPBlockDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipblock" {
			continue
		}

		_, err := client.GetIPBlock(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("IPBlock still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching IPBlock %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIPBlockAttributes(n string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["location"] != location {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["location"])
		}

		return nil
	}
}

func testAccCheckIPBlockExists(n string, ipblock *profitbricks.IPBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIPBlockExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundIP, err := client.GetIPBlock(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching IP Block: %s", rs.Primary.ID)
		}
		if foundIP.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		ipblock = foundIP

		return nil
	}
}

const testAccCheckIPBlockConfig_basic = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "ipblock TF test"
}`

const testAccCheckIPBlockConfig_update = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "updated"
}`
