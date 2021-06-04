package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIPBlock_Basic(t *testing.T) {
	var ipblock ionoscloud.IpBlock
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckipblockconfigBasic, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "location", location),
				),
			},
			{
				Config: fmt.Sprintf(testacccheckipblockconfigUpdate, location),
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
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipblock" {
			continue
		}

		_, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse != nil && apiResponse.StatusCode != 404 {
				payload := fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				return fmt.Errorf("IPBlock still exists %s - an error occurred while checking it %s %s", rs.Primary.ID, err, payload)
			}
		} else {
			return fmt.Errorf("IPBlock still exists %s", rs.Primary.ID)
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
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["location"])
		}

		return nil
	}
}

func testAccCheckIPBlockExists(n string, ipblock *ionoscloud.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIPBlockExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundIP, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching IP Block: %s %s", rs.Primary.ID, payload)
		}
		if *foundIP.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		ipblock = &foundIP

		return nil
	}
}

const testacccheckipblockconfigBasic = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "ipblock TF test"
}`

const testacccheckipblockconfigUpdate = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "updated"
}`
