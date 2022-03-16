//go:build compute || all || pcc

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

func TestAccPrivateCrossConnectBasic(t *testing.T) {
	var privateCrossConnect ionoscloud.PrivateCrossConnect

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPrivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCrossConnectConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateCrossConnectExists(PCCResource+"."+PCCTestResource, &privateCrossConnect),
					resource.TestCheckResourceAttr(PCCResource+"."+PCCTestResource, "name", PCCTestResource),
					resource.TestCheckResourceAttr(PCCResource+"."+PCCTestResource, "description", PCCTestResource),
				),
			},
			{
				Config: testAccDataSourcePccMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "name", PCCResource+"."+PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "description", PCCResource+"."+PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "peers", PCCResource+"."+PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceById, "connectable_datacenters", PCCResource+"."+PCCTestResource, "connectable_datacenters"),
				),
			},
			{
				Config: testAccDataSourcePccMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "name", PCCResource+"."+PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "description", PCCResource+"."+PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "peers", PCCResource+"."+PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PCCResource+"."+PCCDataSourceByName, "connectable_datacenters", PCCResource+"."+PCCTestResource, "connectable_datacenters"),
				),
			},
			{
				Config:      testAccDataSourcePccMultipleResultsError,
				ExpectError: regexp.MustCompile(`more than one pcc found with the specified criteria: name`),
			},
			{
				Config:      testAccDataSourcePccWrongNameError,
				ExpectError: regexp.MustCompile(`no pcc found with the specified criteria: name`),
			},
			{
				Config: testAccCheckPrivateCrossConnectConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateCrossConnectExists(PCCResource+"."+PCCTestResource, &privateCrossConnect),
					resource.TestCheckResourceAttr(PCCResource+"."+PCCTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(PCCResource+"."+PCCTestResource, "description", UpdatedResources),
				),
			},
		},
	})
}

func testAccCheckPrivateCrossConnectDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != PCCResource {
			continue
		}

		_, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking private cross-connect %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("unable to fetch private cross-connect %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckPrivateCrossConnectExists(n string, privateCrossConnect *ionoscloud.PrivateCrossConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		foundPrivateCrossConnect, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching private cross-connect: %s", rs.Primary.ID)
		}
		if *foundPrivateCrossConnect.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		privateCrossConnect = &foundPrivateCrossConnect

		return nil
	}
}

const testAccCheckPrivateCrossConnectConfigUpdate = `
resource ` + PCCResource + ` ` + PCCTestResource + ` {
  name        = "` + UpdatedResources + `"
  description = "` + UpdatedResources + `"
}`

const testAccDataSourcePccMatchId = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + PCCResource + ` ` + PCCDataSourceById + ` {
  id			= ` + PCCResource + `.` + PCCTestResource + `.id
}
`

const testAccDataSourcePccMatchName = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + PCCResource + ` ` + PCCDataSourceByName + ` {
  name			= "` + PCCTestResource + `"
}
`

const testAccDataSourcePccWrongNameError = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + PCCResource + ` ` + PCCDataSourceByName + ` {
  name			= "wrong_name"
}
`

const testAccDataSourcePccMultipleResultsError = testAccCheckPrivateCrossConnectConfigBasic + `
resource ` + PCCResource + ` ` + PCCTestResource + `_multiple_results {
  name        = "` + PCCTestResource + `"
  description = "` + PCCTestResource + `"
}

data ` + PCCResource + ` ` + PCCDataSourceByName + ` {
  name			= "` + PCCTestResource + `"
}
`
