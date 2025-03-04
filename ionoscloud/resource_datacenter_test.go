//go:build compute || all || datacenter

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataCenterBasic(t *testing.T) {
	var datacenter ionoscloud.Datacenter

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatacenterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "name", constant.DatacenterTestResource),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "description", "Test Datacenter Description"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection", "false"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttrSet(constant.DatacenterResource+"."+constant.DatacenterTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "ipv6_cidr_block", constant.DatacenterResource+"."+constant.DatacenterTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "ipv6_cidr_block", constant.DatacenterResource+"."+constant.DatacenterTestResource, "ipv6_cidr_block"),
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
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "description", "Test Datacenter Description Updated"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection", "false"),
				),
			},
		},
	})
}

func testAccCheckDatacenterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DatacenterResource {
			continue
		}

		_, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of datacenter %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("datacenter %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDatacenterExists(n string, datacenter *ionoscloud.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching DC: %s", rs.Primary.ID)
		}
		if *foundDC.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		datacenter = &foundDC

		return nil
	}
}

const testAccCheckDatacenterConfigUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       =  "` + constant.UpdatedResources + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}`

const testAccDataSourceDatacenterMatchId = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceById + ` {
  id			= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}`

const testAccDataSourceDatacenterMatchName = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceByName + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
}`

const testAccDataSourceDatacenterMatching = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterMultipleResultsError = testAccCheckDatacenterConfigBasic + `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + `_multiple_results {
	name       = "` + constant.DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}

data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterWrongLocationError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location =  "wrong_location"
}`

const testAccDataSourceDatacenterWrongNameAndLocationError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name =  "wrong_name"
    location =  "wrong_location"
}`
