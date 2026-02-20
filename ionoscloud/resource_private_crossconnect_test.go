//go:build compute || all || pcc

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPrivateCrossConnectBasic(t *testing.T) {
	var privateCrossConnect ionoscloud.PrivateCrossConnect

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckPrivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCrossConnectConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateCrossConnectExists(constant.PCCResource+"."+constant.PCCTestResource, &privateCrossConnect),
					resource.TestCheckResourceAttr(constant.PCCResource+"."+constant.PCCTestResource, "name", constant.PCCTestResource),
					resource.TestCheckResourceAttr(constant.PCCResource+"."+constant.PCCTestResource, "description", constant.PCCTestResource),
				),
			},
			{
				Config: testAccDataSourcePccMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceById, "name", constant.PCCResource+"."+constant.PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceById, "description", constant.PCCResource+"."+constant.PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceById, "peers", constant.PCCResource+"."+constant.PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceById, "connectable_datacenters", constant.PCCResource+"."+constant.PCCTestResource, "connectable_datacenters"),
				),
			},
			{
				Config: testAccDataSourcePccMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceByName, "name", constant.PCCResource+"."+constant.PCCTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceByName, "description", constant.PCCResource+"."+constant.PCCTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceByName, "peers", constant.PCCResource+"."+constant.PCCTestResource, "peers"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PCCResource+"."+constant.PCCDataSourceByName, "connectable_datacenters", constant.PCCResource+"."+constant.PCCTestResource, "connectable_datacenters"),
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
					testAccCheckPrivateCrossConnectExists(constant.PCCResource+"."+constant.PCCTestResource, &privateCrossConnect),
					resource.TestCheckResourceAttr(constant.PCCResource+"."+constant.PCCTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.PCCResource+"."+constant.PCCTestResource, "description", constant.UpdatedResources),
				),
			},
		},
	})
}

func testAccCheckPrivateCrossConnectDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PCCResource {
			continue
		}

		_, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking private cross-connect %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("unable to fetch private cross-connect %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckPrivateCrossConnectExists(n string, privateCrossConnect *ionoscloud.PrivateCrossConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")

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
			return fmt.Errorf("error occurred while fetching private cross-connect: %s", rs.Primary.ID)
		}
		if *foundPrivateCrossConnect.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		privateCrossConnect = &foundPrivateCrossConnect

		return nil
	}
}

const testAccCheckPrivateCrossConnectConfigUpdate = `
resource ` + constant.PCCResource + ` ` + constant.PCCTestResource + ` {
  name        = "` + constant.UpdatedResources + `"
  description = "` + constant.UpdatedResources + `"
}`

const testAccDataSourcePccMatchId = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + constant.PCCResource + ` ` + constant.PCCDataSourceById + ` {
  id			= ` + constant.PCCResource + `.` + constant.PCCTestResource + `.id
}
`

const testAccDataSourcePccMatchName = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + constant.PCCResource + ` ` + constant.PCCDataSourceByName + ` {
  name			= "` + constant.PCCTestResource + `"
}
`

const testAccDataSourcePccWrongNameError = testAccCheckPrivateCrossConnectConfigBasic + `
data ` + constant.PCCResource + ` ` + constant.PCCDataSourceByName + ` {
  name			= "wrong_name"
}
`

const testAccDataSourcePccMultipleResultsError = testAccCheckPrivateCrossConnectConfigBasic + `
resource ` + constant.PCCResource + ` ` + constant.PCCTestResource + `_multiple_results {
  name        = "` + constant.PCCTestResource + `"
  description = "` + constant.PCCTestResource + `"
}

data ` + constant.PCCResource + ` ` + constant.PCCDataSourceByName + ` {
  name			= "` + constant.PCCTestResource + `"
}
`
