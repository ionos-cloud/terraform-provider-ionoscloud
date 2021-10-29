package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIPBlockBasic(t *testing.T) {
	var ipblock ionoscloud.IpBlock
	location := "us/las"
	resourceName := IpBLockResource + ".webserver_ip"
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
					testAccCheckIPBlockExists(resourceName, &ipblock),
					testAccCheckIPBlockAttributes(resourceName, location),
					resource.TestCheckResourceAttr(resourceName, "location", location),
					resource.TestCheckResourceAttr(resourceName, "name", "ipblock TF test"),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testacccheckipblockconfigUpdate, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(resourceName, &ipblock),
					testAccCheckIPBlockAttributes(resourceName, location),
					resource.TestCheckResourceAttr(resourceName, "name", "updated"),
					resource.TestCheckResourceAttr(resourceName, "size", "2"),
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
		if rs.Type != IpBLockResource {
			continue
		}

		_, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("IPBlock still exists %s - an error occurred while checking it %s", rs.Primary.ID, err)
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

const testacccheckipblockconfigBasic = `
resource ` + IpBLockResource + ` "webserver_ip" {
  location = "%s"
  size = 1
  name = "ipblock TF test"
}`

const testacccheckipblockconfigUpdate = `
resource ` + IpBLockResource + `"webserver_ip" {
  location = "%s"
  size = 2
  name = "updated"
}`
