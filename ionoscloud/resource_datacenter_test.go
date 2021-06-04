package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataCenter_Basic(t *testing.T) {
	var datacenter ionoscloud.Datacenter
	dcName := "datacenter-test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckdatacenterconfigBasic, dcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists("ionoscloud_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("ionoscloud_datacenter.foobar", "name", dcName),
				),
			},
			{
				Config: testacccheckdatacenterconfigUpdate,
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
				payload := ""
				if apiResponse != nil {
					payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				}
				return fmt.Errorf("datacenter still exists %s - an error occurred while checking it %s %s", rs.Primary.ID, err, payload)
			}
		} else {
			return fmt.Errorf("datacenter still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDatacenterExists(n string, datacenter *ionoscloud.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundDC, apiResponse, err := client.DataCenterApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching DC: %s %s", rs.Primary.ID, payload)
		}
		if *foundDC.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		datacenter = &foundDC

		return nil
	}
}

const testacccheckdatacenterconfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "%s"
	location = "us/las"
}`

const testacccheckdatacenterconfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       =  "updated"
	location = "us/las"
}`
