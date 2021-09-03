package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLan_Basic(t *testing.T) {
	var lan ionoscloud.Lan
	lanName := "lanName"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLanDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckLanConfigBasic, lanName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanExists("ionoscloud_lan.webserver_lan", &lan),
					testAccCheckLanAttributes("ionoscloud_lan.webserver_lan", lanName),
					resource.TestCheckResourceAttr("ionoscloud_lan.webserver_lan", "name", lanName),
				),
			},
			{
				Config: testAccCheckLanConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanAttributes("ionoscloud_lan.webserver_lan", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_lan.webserver_lan", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckLanDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, apiResponse, err := client.LansApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while looking for lan %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("LAN still exists %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckLanAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s != %s", rs.Primary.Attributes["name"], name)
		}

		return nil
	}
}

func testAccCheckLanExists(n string, lan *ionoscloud.Lan) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckLanExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}
		foundLan, _, err := client.LansApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Server: %s", rs.Primary.ID)
		}
		if *foundLan.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		lan = &foundLan

		return nil
	}
}

const testAccCheckLanConfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "%s"
}`

const testAccCheckLanConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "updated"
}`
