package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataCenterBasic(t *testing.T) {
	var datacenter ionoscloud.Datacenter

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatacenterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(DatacenterResource+"."+DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "name", DatacenterTestResource),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "description", "Test Datacenter Description"),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection", "false"),
				),
			},
			{
				Config: testAccCheckDatacenterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(DatacenterResource+"."+DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "description", "Test Datacenter Description Updated"),
					resource.TestCheckResourceAttr(DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection", "false"),
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
		if rs.Type != DatacenterResource {
			continue
		}

		_, apiResponse, err := client.DataCenterApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
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
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching DC: %s", rs.Primary.ID)
		}
		if *foundDC.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		datacenter = &foundDC

		return nil
	}
}

const testAccCheckDatacenterConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "` + DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}`

const testAccCheckDatacenterConfigUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       =  "` + UpdatedResources + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}`
