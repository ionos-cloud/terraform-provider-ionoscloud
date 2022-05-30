//go:build compute || all || volume

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

func TestAccVolumeBasic(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				//added to test - #266. crash when using image_alias on volume
				Config:      testAccCheckVolumeConfigBasicErrorNoPassOrSSHPath,
				ExpectError: regexp.MustCompile(`either 'image_password' or 'ssh_key_path' must be provided`),
			},
			{
				Config: testAccCheckVolumeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(VolumeResource+"."+VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", VolumeTestResource),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(VolumeResource+"."+VolumeTestResource, "image_name"),
					resource.TestCheckResourceAttrPair(VolumeResource+"."+VolumeTestResource, "boot_server", ServerResource+"."+ServerTestResource, "id"),
					testImageNotNull(VolumeResource, "image")),
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
				Config: testAccDataSourceVolumePartialMatchName,
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
			{
				Config: testAccDataSourceVolumeMatchServerId,
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
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "server_id", ServerResource+"."+ServerTestResource, "id")),
			},
			{
				Config: testAccDataSourceVolumeMatchServerIdAndVolumeId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "id", VolumeResource+"."+VolumeTestResource, "id"),
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
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+VolumeResource+"."+VolumeDataSourceByName, "server_id", ServerResource+"."+ServerTestResource, "id")),
			},
			{
				Config:      testAccDataSourceVolumeWrongNameError,
				ExpectError: regexp.MustCompile(`no volume found with the specified criteria: name`),
			},
			{
				Config:      testAccDataSourceVolumeWrongPartialNameError,
				ExpectError: regexp.MustCompile(`no volume found with the specified criteria: name`),
			},
			{
				Config: testAccCheckVolumeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "6"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(VolumeResource+"."+VolumeTestResource, "image_name"),
					resource.TestCheckResourceAttrPair(VolumeResource+"."+VolumeTestResource, "boot_server", ServerResource+"."+ServerTestResource+"updated", "id"),
					testImageNotNull(VolumeResource, "image")),
			},
		},
	})
}

func TestAccVolumeNoPassword(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeConfigNoPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(VolumeResource+"."+VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", VolumeTestResource),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "4"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "licence_type", "unknown"),
				)},
			{
				Config: testAccCheckVolumeConfigNoPasswordUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(VolumeResource+"."+VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "licence_type", "other"),
				)},
		},
	})
}

func TestAccVolumeResolveImageName(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(VolumeResource+"."+VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", VolumeTestResource),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "availability_zone", "ZONE_1"),
					testImageNotNull(VolumeResource, "image"))},
		},
	})
}

func testAccCheckVolumeDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != VolumeResource {
			continue
		}

		_, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("volume still exists %s - an error occurred while checking it %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("volume still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVolumeExists(n string, volume *ionoscloud.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckVolumeExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundServer, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Volume: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		volume = &foundServer

		return nil
	}
}

const testAccCheckVolumeConfigBasic = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + `{
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	availability_zone = "ZONE_1"
	name = "` + VolumeTestResource + `"
	size = 5
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	image_password = "K3tTj8G14a3EgKyNeeiY"
	user_data = "foo"
}`

const testAccCheckVolumeConfigBasicErrorNoPassOrSSHPath = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + `{
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	availability_zone = "ZONE_1"
	name = "` + VolumeTestResource + `"
	size = 5
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	user_data = "foo"

}`

//ubuntu-21.10-server-cloudimg-amd64-20220201
const testAccCheckVolumeConfigUpdate = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + `updated {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `updated.id
	availability_zone = "ZONE_1"
	name = "` + UpdatedResources + `"
	size = 6
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	image_password = "K3tTj8G14a3EgKyNeeiYupdated"
	user_data = "foo"
}`

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

var testAccDataSourceVolumeMatchServerId = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= ` + VolumeResource + `.` + VolumeTestResource + `.name
  server_id		= ` + ServerResource + `.` + ServerTestResource + `.id
}
`

var testAccDataSourceVolumeMatchServerIdAndVolumeId = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id		     	= ` + VolumeResource + `.` + VolumeTestResource + `.id
  server_id			= ` + ServerResource + `.` + ServerTestResource + `.id
}
`

var testAccDataSourceVolumePartialMatchName = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + DataSourcePartial + `"
  partial_match = true
}
`

var testAccDataSourceVolumeWrongNameError = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

var testAccDataSourceVolumeWrongPartialNameError = testAccCheckVolumeConfigBasic + `
data ` + VolumeResource + ` ` + VolumeDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
  partial_match = true

}
`

const testAccCheckVolumeConfigNoPassword = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  name = "` + VolumeTestResource + `"
  size           = 4
  disk_type      = "HDD"
  licence_type   = "unknown"
}`

const testAccCheckVolumeConfigNoPasswordUpdate = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  name = "` + UpdatedResources + `"
  size           = 5
  disk_type      = "HDD"
  licence_type   = "other"
}`

const testAccCheckVolumeResolveImageName = testAccCheckLanConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + VolumeResource + ` ` + VolumeTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  availability_zone = "ZONE_1"
  name = "` + VolumeTestResource + `"
  size = 5
  disk_type = "SSD Standard"
  bus = "VIRTIO"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
}
`
