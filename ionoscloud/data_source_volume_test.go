//go:build compute || all || volume

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceVolume(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeConfigBasic,
			},
			{
				Config: testAccDataSourceVolumeMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "name", VolumeResource+"."+VolumeTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "image", VolumeResource+"."+VolumeTestResource, "image"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "image_alias", VolumeResource+"."+VolumeTestResource, "image_alias"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "disk_type", VolumeResource+"."+VolumeTestResource, "disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "sshkey", VolumeResource+"."+VolumeTestResource, "sshkey"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "bus", VolumeResource+"."+VolumeTestResource, "bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "availability_zone", VolumeResource+"."+VolumeTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "cpu_hot_plug", VolumeResource+"."+VolumeTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "ram_hot_plug", VolumeResource+"."+VolumeTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "nic_hot_plug", VolumeResource+"."+VolumeTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "nic_hot_unplug", VolumeResource+"."+VolumeTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "disc_virtio_hot_plug", VolumeResource+"."+VolumeTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "disc_virtio_hot_unplug", VolumeResource+"."+VolumeTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "device_number", VolumeResource+"."+VolumeTestResource, "device_number"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceById, "boot_server", ServerResource+"."+ServerTestResource, "id"),
				),
			},
			{
				Config: testAccDataSourceVolumeMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "name", VolumeResource+"."+VolumeTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "image", VolumeResource+"."+VolumeTestResource, "image"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "image_alias", VolumeResource+"."+VolumeTestResource, "image_alias"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "disk_type", VolumeResource+"."+VolumeTestResource, "disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "sshkey", VolumeResource+"."+VolumeTestResource, "sshkey"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "bus", VolumeResource+"."+VolumeTestResource, "bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "availability_zone", VolumeResource+"."+VolumeTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "cpu_hot_plug", VolumeResource+"."+VolumeTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "ram_hot_plug", VolumeResource+"."+VolumeTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "nic_hot_plug", VolumeResource+"."+VolumeTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "nic_hot_unplug", VolumeResource+"."+VolumeTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "disc_virtio_hot_plug", VolumeResource+"."+VolumeTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "disc_virtio_hot_unplug", VolumeResource+"."+VolumeTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "device_number", VolumeResource+"."+VolumeTestResource, "device_number"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "boot_server", ServerResource+"."+ServerTestResource, "id")),
			},
		},
	})
}

var testAccDataSourceVolumeMatchId = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + VolumeResource + `.` + VolumeTestResource + `.id
}
`

var testAccDataSourceVolumeMatchName = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= ` + VolumeResource + `.` + VolumeTestResource + `.name
}
`
