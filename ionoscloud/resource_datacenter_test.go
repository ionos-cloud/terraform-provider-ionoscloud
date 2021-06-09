package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataCenter_Basic(t *testing.T) {
	var datacenter ionoscloud.Datacenter
	dc_name := "datacenter-test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDatacenterConfig_basic, dc_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists("ionoscloud_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("ionoscloud_datacenter.foobar", "name", dc_name),
				),
			},
			{
				Config: testAccCheckDatacenterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists("ionoscloud_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("ionoscloud_datacenter.foobar", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDatacenterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, apiResponse, err := client.DataCenterApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of datacenter %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("datacenter %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDatacenterExists(n string, datacenter *ionoscloud.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundDC, _, err := client.DataCenterApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching DC: %s", rs.Primary.ID)
		}
		if *foundDC.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		datacenter = &foundDC

		return nil
	}
}

const testAccCheckDatacenterConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "%s"
	location = "us/las"
}`

const testAccCheckDatacenterConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       =  "updated"
	location = "us/las"
}`
