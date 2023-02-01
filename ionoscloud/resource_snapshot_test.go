//go:build compute || all || snapshot

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"regexp"
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
		ExternalProviders: randomProviderVersion343(),
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
				Config: testAccDataSourceSnapshotMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "name", SnapshotResource+"."+SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "location", SnapshotResource+"."+SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "size", SnapshotResource+"."+SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "description", SnapshotResource+"."+SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "licence_type", SnapshotResource+"."+SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "sec_auth_protection", SnapshotResource+"."+SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "cpu_hot_plug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "cpu_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "ram_hot_plug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "ram_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "nic_hot_plug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "nic_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "disc_virtio_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "disc_virtio_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "disc_scsi_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceById, "disc_scsi_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_unplug"),
				),
			},
			{
				Config: testAccDataSourceSnapshotMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "name", SnapshotResource+"."+SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "location", SnapshotResource+"."+SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "size", SnapshotResource+"."+SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "description", SnapshotResource+"."+SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "licence_type", SnapshotResource+"."+SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "sec_auth_protection", SnapshotResource+"."+SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "cpu_hot_plug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "cpu_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "ram_hot_plug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "ram_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "nic_hot_plug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "nic_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_virtio_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_virtio_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_scsi_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_scsi_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_unplug"),
				),
			},
			{
				Config: testAccDataSourceSnapshotMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "name", SnapshotResource+"."+SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "location", SnapshotResource+"."+SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "size", SnapshotResource+"."+SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "description", SnapshotResource+"."+SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "licence_type", SnapshotResource+"."+SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "sec_auth_protection", SnapshotResource+"."+SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "cpu_hot_plug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "cpu_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "ram_hot_plug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "ram_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "nic_hot_plug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "nic_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_virtio_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_virtio_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_scsi_hot_plug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+SnapshotResource+"."+SnapshotDataSourceByName, "disc_scsi_hot_unplug", SnapshotResource+"."+SnapshotTestResource, "disc_scsi_hot_unplug"),
				),
			},
			{
				Config:      testAccDataSourceSnapshotWrongNameError,
				ExpectError: regexp.MustCompile(`no snapshot found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceSnapshotWrongLocation,
				ExpectError: regexp.MustCompile(`no snapshot found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceSnapshotWrongSize,
				ExpectError: regexp.MustCompile(`no snapshot found with the specified criteria`),
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
			if !httpNotFound(apiResponse) {
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

const testAccDataSourceSnapshotMatchId = testAccCheckSnapshotConfigBasic + `
data ` + SnapshotResource + ` ` + SnapshotDataSourceById + ` {
  id			= ` + SnapshotResource + `.` + SnapshotTestResource + `.id
}`

const testAccDataSourceSnapshotMatching = testAccCheckSnapshotConfigBasic + `
data ` + SnapshotResource + ` ` + SnapshotDataSourceByName + ` {
    name = ` + SnapshotResource + `.` + SnapshotTestResource + `.name
    location = ` + SnapshotResource + `.` + SnapshotTestResource + `.location
    size = ` + SnapshotResource + `.` + SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongNameError = testAccCheckSnapshotConfigBasic + `
data ` + SnapshotResource + ` ` + SnapshotDataSourceByName + ` {
    name = "wrong_name"
    location = ` + SnapshotResource + `.` + SnapshotTestResource + `.location
    size = ` + SnapshotResource + `.` + SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongLocation = testAccCheckSnapshotConfigBasic + `
data ` + SnapshotResource + ` ` + SnapshotDataSourceByName + ` {
    name = ` + SnapshotResource + `.` + SnapshotTestResource + `.name
    location = "wrong_location"
    size = ` + SnapshotResource + `.` + SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongSize = testAccCheckSnapshotConfigBasic + `
data ` + SnapshotResource + ` ` + SnapshotDataSourceByName + ` {
    name = ` + SnapshotResource + `.` + SnapshotTestResource + `.name
    location = ` + SnapshotResource + `.` + SnapshotTestResource + `.location
    size = 1234
}`
