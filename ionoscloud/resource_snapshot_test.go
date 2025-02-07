//go:build compute || all || snapshot

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSnapshotBasic(t *testing.T) {
	var snapshot ionoscloud.Snapshot

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSnapshotConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists(constant.SnapshotResource+"."+constant.SnapshotTestResource, &snapshot),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "name", constant.SnapshotTestResource),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "description", constant.SnapshotTestResource),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "sec_auth_protection", "true"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "licence_type", "LINUX"),
				),
			},
			{
				Config: testAccDataSourceSnapshotMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "name", constant.SnapshotResource+"."+constant.SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "location", constant.SnapshotResource+"."+constant.SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "size", constant.SnapshotResource+"."+constant.SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "description", constant.SnapshotResource+"."+constant.SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "licence_type", constant.SnapshotResource+"."+constant.SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "sec_auth_protection", constant.SnapshotResource+"."+constant.SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "cpu_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "cpu_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "ram_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "ram_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "nic_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "nic_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "disc_virtio_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "disc_virtio_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "disc_scsi_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceById, "disc_scsi_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_unplug"),
				),
			},
			{
				Config: testAccDataSourceSnapshotMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "name", constant.SnapshotResource+"."+constant.SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "location", constant.SnapshotResource+"."+constant.SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "size", constant.SnapshotResource+"."+constant.SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "description", constant.SnapshotResource+"."+constant.SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "licence_type", constant.SnapshotResource+"."+constant.SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "sec_auth_protection", constant.SnapshotResource+"."+constant.SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "cpu_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "cpu_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "ram_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "ram_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "nic_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "nic_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_virtio_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_virtio_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_scsi_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_scsi_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_unplug"),
				),
			},
			{
				Config: testAccDataSourceSnapshotMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "name", constant.SnapshotResource+"."+constant.SnapshotTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "location", constant.SnapshotResource+"."+constant.SnapshotTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "size", constant.SnapshotResource+"."+constant.SnapshotTestResource, "size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "description", constant.SnapshotResource+"."+constant.SnapshotTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "licence_type", constant.SnapshotResource+"."+constant.SnapshotTestResource, "licence_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "sec_auth_protection", constant.SnapshotResource+"."+constant.SnapshotTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "cpu_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "cpu_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "ram_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "ram_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "nic_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "nic_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_virtio_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_virtio_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_scsi_hot_plug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.SnapshotResource+"."+constant.SnapshotDataSourceByName, "disc_scsi_hot_unplug", constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_scsi_hot_unplug"),
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
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "description", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "sec_auth_protection", "false"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "cpu_hot_plug", "false"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "nic_hot_plug", "false"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "ram_hot_plug", "false"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_unplug", "false"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "disc_virtio_hot_plug", "true"),
					resource.TestCheckResourceAttr(constant.SnapshotResource+"."+constant.SnapshotTestResource, "licence_type", "OTHER"),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.SnapshotResource {
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
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching Snapshot: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		snapshot = &foundServer

		return nil
	}
}

const testAccCheckSnapshotConfigBasic = testSnapshotServer + `
resource ` + constant.SnapshotResource + ` ` + constant.SnapshotTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  volume_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.boot_volume
  name = "` + constant.SnapshotTestResource + `"
  description = "` + constant.SnapshotTestResource + `"
  sec_auth_protection = true
  licence_type = "LINUX"
}
`

const testAccCheckSnapshotConfigUpdate = testSnapshotServer + `
resource ` + constant.SnapshotResource + ` ` + constant.SnapshotTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  volume_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.boot_volume
  name = "` + constant.UpdatedResources + `"
  description = "` + constant.UpdatedResources + `"
  sec_auth_protection = false
  cpu_hot_plug = false
  nic_hot_plug = false
  disc_virtio_hot_plug = true
  disc_virtio_hot_unplug = false
  ram_hot_plug = false
  licence_type = "OTHER"
}`

const testAccDataSourceSnapshotMatchId = testAccCheckSnapshotConfigBasic + `
data ` + constant.SnapshotResource + ` ` + constant.SnapshotDataSourceById + ` {
  id = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.id
}`

const testAccDataSourceSnapshotMatching = testAccCheckSnapshotConfigBasic + `
data ` + constant.SnapshotResource + ` ` + constant.SnapshotDataSourceByName + ` {
    name = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.name
    location = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.location
    size = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongNameError = testAccCheckSnapshotConfigBasic + `
data ` + constant.SnapshotResource + ` ` + constant.SnapshotDataSourceByName + ` {
    name = "wrong_name"
    location = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.location
    size = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongLocation = testAccCheckSnapshotConfigBasic + `
data ` + constant.SnapshotResource + ` ` + constant.SnapshotDataSourceByName + ` {
    name = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.name
    location = "wrong_location"
    size = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.size
}`

const testAccDataSourceSnapshotWrongSize = testAccCheckSnapshotConfigBasic + `
data ` + constant.SnapshotResource + ` ` + constant.SnapshotDataSourceByName + ` {
    name = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.name
    location = ` + constant.SnapshotResource + `.` + constant.SnapshotTestResource + `.location
    size = 1234
}`
