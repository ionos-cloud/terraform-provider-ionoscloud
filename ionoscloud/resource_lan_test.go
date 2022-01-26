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
					testAccCheckLanExists(LanResource+"."+LanTestResource, &lan),
					resource.TestCheckResourceAttr(LanResource+"."+LanTestResource, "name", LanTestResource),
					resource.TestCheckResourceAttr(LanResource+"."+LanTestResource, "public", "true"),
				),
			},
			{
				Config: testAccDataSourceLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "name", LanResource+"."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "ip_failover.nic_uuid", LanResource+"."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "ip_failover.ip", LanResource+"."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "pcc", LanResource+"."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "public", LanResource+"."+LanTestResource, "public"),
				),
			},
			{
				Config: testAccDataSourceLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "name", LanResource+"."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "ip_failover.nic_uuid", LanResource+"."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "ip_failover.ip", LanResource+"."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "pcc", LanResource+"."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "public", LanResource+"."+LanTestResource, "public"),
				),
			},
			{
				Config:      testAccDataSourceLanWrongName,
				ExpectError: regexp.MustCompile(`no lan found with the specified name`),
			},
			{
				Config: testAccCheckLanConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(LanResource+"."+LanTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(LanResource+"."+LanTestResource, "public", "false"),
					resource.TestCheckResourceAttrPair(LanResource+"."+LanTestResource, "pcc", PCCResource+"."+PCCTestResource, "id"),
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
		if rs.Type != LanResource {
			continue
		}

		_, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while looking for lan %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("LAN still exists %s", rs.Primary.ID)
		}

	}

	return nil
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

const testAccCheckLanConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "` + LanTestResource + `"
}`

const testAccCheckLanConfigUpdate = testAccCheckDatacenterConfigBasic + testAccCheckPrivateCrossConnectConfigBasic + `
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = false
  name = "` + UpdatedResources + `"
  pcc = ` + PCCResource + `.` + PCCTestResource + `.id
}`

const testAccDataSourceLanMatchId = testAccCheckLanConfigBasic + `
data ` + LanResource + ` ` + LanDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + LanResource + `.` + LanTestResource + `.id
}
`

const testAccDataSourceLanMatchName = testAccCheckLanConfigBasic + `
data ` + LanResource + ` ` + LanDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + LanTestResource + `"
}
`

const testAccDataSourceLanWrongName = testAccCheckLanConfigBasic + `
data ` + LanResource + ` ` + LanDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
}
`
