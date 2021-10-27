package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSnapshot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckSnapshotConfigBasic,
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
		},
	})

}

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
