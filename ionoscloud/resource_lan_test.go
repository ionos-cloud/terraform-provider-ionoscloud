//go:build compute || all || lan

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLanBasic(t *testing.T) {
	var lan ionoscloud.Lan

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLanDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLanConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanExists(constant.LanResource+"."+constant.LanTestResource, &lan),
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.LanTestResource, "name", constant.LanTestResource),
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.LanTestResource, "public", "true"),
				),
			},
			{
				Config: testAccDataSourceLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "name", constant.LanResource+"."+constant.LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "ip_failover.nic_uuid", constant.LanResource+"."+constant.LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "ip_failover.ip", constant.LanResource+"."+constant.LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "pcc", constant.LanResource+"."+constant.LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "public", constant.LanResource+"."+constant.LanTestResource, "public"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceById, "ipv6_cidr_block", constant.LanResource+"."+constant.LanTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config: testAccDataSourceLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "name", constant.LanResource+"."+constant.LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "ip_failover.nic_uuid", constant.LanResource+"."+constant.LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "ip_failover.ip", constant.LanResource+"."+constant.LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "pcc", constant.LanResource+"."+constant.LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "public", constant.LanResource+"."+constant.LanTestResource, "public"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByName, "ipv6_cidr_block", constant.LanResource+"."+constant.LanTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config:      testAccDataSourceLanMultipleResultsError,
				ExpectError: regexp.MustCompile(`more than one lan found with the specified criteria name`),
			},
			{
				Config:      testAccDataSourceLanWrongNameError,
				ExpectError: regexp.MustCompile(`no lan found with the specified name`),
			},
			{
				Config: testAccCheckLanConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.LanTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.LanTestResource, "public", "false"),
					resource.TestCheckResourceAttrPair(constant.LanResource+"."+constant.LanTestResource, "pcc", constant.PCCResource+"."+constant.PCCTestResource, "id"),
				),
			},
		},
	})
}

func testAccCheckLanDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.LanResource {
			continue
		}

		_, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while looking for lan %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("LAN still exists %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckLanExists(n string, lan *ionoscloud.Lan) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
		foundLan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

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

const testAccCheckLanConfigUpdate = testAccCheckDatacenterConfigBasic + testAccCheckPrivateCrossConnectConfigBasic + `
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = false
  name = "` + constant.UpdatedResources + `"
  pcc = ` + constant.PCCResource + `.` + constant.PCCTestResource + `.id
}`

const testAccDataSourceLanMatchId = testAccCheckLanConfigBasic + `
data ` + constant.LanResource + ` ` + constant.LanDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.LanResource + `.` + constant.LanTestResource + `.id
}
`

const testAccDataSourceLanMatchName = testAccCheckLanConfigBasic + `
data ` + constant.LanResource + ` ` + constant.LanDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.LanTestResource + `"
}
`

const testAccDataSourceLanMultipleResultsError = testAccCheckLanConfigBasic + `
resource ` + constant.LanResource + ` ` + constant.LanTestResource + `_multiple_results {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "` + constant.LanTestResource + `"
}

data ` + constant.LanResource + ` ` + constant.LanDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.LanTestResource + `"
}
`

const testAccDataSourceLanWrongNameError = testAccCheckLanConfigBasic + `
data ` + constant.LanResource + ` ` + constant.LanDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "wrong_name"
}
`
