package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLan_Basic(t *testing.T) {
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
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != LanResource {
			continue
		}

		_, apiResponse, err := client.LansApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

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
		foundLan, _, err := client.LansApi.DatacentersLansFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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
