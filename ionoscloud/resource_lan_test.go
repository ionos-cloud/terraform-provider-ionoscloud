//go:build compute || all || lan

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

func TestAccLanBasic(t *testing.T) {
	var lan ionoscloud.Lan
	var privateLAN ionoscloud.Lan

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckLanDestroyCheck,
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
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "name", constant.LanResource+"."+constant.LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "ip_failover.nic_uuid", constant.LanResource+"."+constant.LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "ip_failover.ip", constant.LanResource+"."+constant.LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "pcc", constant.LanResource+"."+constant.LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "public", constant.LanResource+"."+constant.LanTestResource, "public"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "ipv6_cidr_block", constant.LanResource+"."+constant.LanTestResource, "ipv6_cidr_block"),
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
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.LanDataSourceByID, "ipv6_cidr_block", constant.LanResource+"."+constant.LanTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config: privateLANConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanExists(constant.LanResource+"."+constant.PrivateLANTestResource, &privateLAN),
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.PrivateLANTestResource, "name", constant.PrivateLANTestResource),
					resource.TestCheckResourceAttr(constant.LanResource+"."+constant.PrivateLANTestResource, "public", "false"),
					resource.TestCheckResourceAttrSet(constant.LanResource+"."+constant.PrivateLANTestResource, "ipv4_cidr_block"),
				),
			},
			{
				Config: dataSourcePrivateLANMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "name", constant.LanResource+"."+constant.PrivateLANTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "ip_failover.nic_uuid", constant.LanResource+"."+constant.PrivateLANTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "ip_failover.ip", constant.LanResource+"."+constant.PrivateLANTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "pcc", constant.LanResource+"."+constant.PrivateLANTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "public", constant.LanResource+"."+constant.PrivateLANTestResource, "public"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "ipv4_cidr_block", constant.LanResource+"."+constant.PrivateLANTestResource, "ipv4_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.LanResource+"."+constant.PrivateLANDataSourceByName, "ipv6_cidr_block", constant.LanResource+"."+constant.PrivateLANTestResource, "ipv6_cidr_block"),
				),
			},
		},
	})
}

func testAccCheckLanDestroyCheck(s *terraform.State) error {
	config := testAccProvider.Meta().(bundleclient.SdkBundle).CloudAPIConfig
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.LanResource {
			continue
		}

		client := config.NewAPIClient(rs.Primary.Attributes["location"])
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
		config := testAccProvider.Meta().(bundleclient.SdkBundle).CloudAPIConfig

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

		client := config.NewAPIClient(rs.Primary.Attributes["location"])
		foundLan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching Server: %s", rs.Primary.ID)
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
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,2)
}
data ` + constant.LanResource + ` ` + constant.LanDataSourceByID + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
}
`

const testAccDataSourceLanMatchId = testAccCheckLanConfigBasic + `
data ` + constant.LanResource + ` ` + constant.LanDataSourceByID + ` {
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

const dataSourcePrivateLANMatchName = privateLANConfig + `
data ` + constant.LanResource + ` ` + constant.PrivateLANDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.PrivateLANTestResource + `"
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
