package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Config: fmt.Sprintf(testAccCheckIPBlockConfigBasic, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "location", location),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckIPBlockConfigUpdate, location),
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
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
			if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking deletion of IPBlock %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("IPBlock still exists %s %s", rs.Primary.ID, err)
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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
		foundIP, _, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching IP Block: %s", rs.Primary.ID)
		}
		if *foundIP.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		ipblock = &foundIP

		return nil
	}
}

const testAccCheckIPBlockConfigBasic = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "ipblock TF test"
}`

const testAccCheckIPBlockConfigUpdate = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "updated"
}`
