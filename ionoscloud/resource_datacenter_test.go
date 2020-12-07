package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccDataCenter_Basic(t *testing.T) {
	var datacenter profitbricks.Datacenter
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
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, err := client.GetDatacenter(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("DataCenter still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching DataCenter %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckDatacenterExists(n string, datacenter *profitbricks.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundDC, err := client.GetDatacenter(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching DC: %s", rs.Primary.ID)
		}
		if foundDC.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		datacenter = foundDC

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
