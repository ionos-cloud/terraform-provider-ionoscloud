//go:build compute || all || snapshot

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSnapshotBasic(t *testing.T) {
	var snapshot ionoscloud.Snapshot

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSnapshotConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists(SnapshotResource+"."+SnapshotTestResource, &snapshot),
					resource.TestCheckResourceAttr(SnapshotResource+"."+SnapshotTestResource, "name", SnapshotTestResource),
				),
			},
			{
				Config: testAccCheckSnapshotConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(SnapshotResource+"."+SnapshotTestResource, "name", UpdatedResources),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != SnapshotResource {
			continue
		}

		_, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("unable to fetch snapshot %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("snapshot %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSnapshotExists(n string, snapshot *ionoscloud.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckSnapshotExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}
		foundServer, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Snapshot: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		snapshot = &foundServer

		return nil
	}
}

const testAccCheckSnapshotConfigBasic = testAccCheckServerConfigBasic + `
resource ` + SnapshotResource + ` ` + SnapshotTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  volume_id = ` + ServerResource + `.` + ServerTestResource + `.boot_volume
  name = "` + SnapshotTestResource + `"
}
`

const testAccCheckSnapshotConfigUpdate = testAccCheckServerConfigBasic + `
resource ` + SnapshotResource + ` ` + SnapshotTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  volume_id = ` + ServerResource + `.` + ServerTestResource + `.boot_volume
  name = "` + UpdatedResources + `"
}`
