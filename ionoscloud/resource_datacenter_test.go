//go:build compute || all || datacenter

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"regexp"
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
				Config: testAccDataSourceDatacenterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceById, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceByName, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "name", DatacenterResource+"."+DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "location", DatacenterResource+"."+DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "description", DatacenterResource+"."+DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "version", DatacenterResource+"."+DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "features", DatacenterResource+"."+DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DatacenterResource+"."+DatacenterDataSourceMatching, "sec_auth_protection", DatacenterResource+"."+DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config:      testAccDataSourceDatacenterMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongNameError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongLocationError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongNameAndLocationError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
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
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DatacenterResource {
			continue
		}

		_, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if errorBesideNotFound(apiResponse) {
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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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

		foundDC, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
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

const testAccCheckDatacenterConfigUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       =  "` + UpdatedResources + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}`

const testAccDataSourceDatacenterMatchId = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceById + ` {
  id			= ` + DatacenterResource + `.` + DatacenterTestResource + `.id
}`

const testAccDataSourceDatacenterMatchName = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceByName + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
}`

const testAccDataSourceDatacenterMatching = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterMultipleResultsError = testAccCheckDatacenterConfigBasic + `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + `_multiple_results {
	name       = "` + DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}

data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterWrongLocationError = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = ` + DatacenterResource + `.` + DatacenterTestResource + `.name
    location =  "wrong_location"
}`

const testAccDataSourceDatacenterWrongNameAndLocationError = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name =  "wrong_name"
    location =  "wrong_location"
}`
