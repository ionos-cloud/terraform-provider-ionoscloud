//go:build compute || all || volume

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVolumeBasic(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				//added to test - #266. crash when using image_alias on volume
				Config:      testAccCheckVolumeConfigBasicErrorNoPassOrSSHPath,
				ExpectError: regexp.MustCompile(`either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided`),
			},
			{
				Config: testAccCheckVolumeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(constant.VolumeResource+"."+constant.VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "name", constant.VolumeTestResource),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.VolumeResource+"."+constant.VolumeTestResource, "image_name"),
					resource.TestCheckResourceAttrPair(constant.VolumeResource+"."+constant.VolumeTestResource, "boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.VolumeResource+"."+constant.VolumeTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					utils.TestImageNotNull(constant.VolumeResource, "image")),
			},
			{
				Config: testAccDataSourceVolumeMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "name", constant.VolumeResource+"."+constant.VolumeTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "image", constant.VolumeResource+"."+constant.VolumeTestResource, "image"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "image_alias", constant.VolumeResource+"."+constant.VolumeTestResource, "image_alias"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "disk_type", constant.VolumeResource+"."+constant.VolumeTestResource, "disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "sshkey", constant.VolumeResource+"."+constant.VolumeTestResource, "sshkey"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "bus", constant.VolumeResource+"."+constant.VolumeTestResource, "bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "availability_zone", constant.VolumeResource+"."+constant.VolumeTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "cpu_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "ram_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "nic_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "nic_hot_unplug", constant.VolumeResource+"."+constant.VolumeTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "disc_virtio_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "disc_virtio_hot_unplug", constant.VolumeResource+"."+constant.VolumeTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "device_number", constant.VolumeResource+"."+constant.VolumeTestResource, "device_number"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceById, "boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
				),
			},
			{
				Config: testAccDataSourceVolumeMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "name", constant.VolumeResource+"."+constant.VolumeTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "image", constant.VolumeResource+"."+constant.VolumeTestResource, "image"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "image_alias", constant.VolumeResource+"."+constant.VolumeTestResource, "image_alias"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "disk_type", constant.VolumeResource+"."+constant.VolumeTestResource, "disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "sshkey", constant.VolumeResource+"."+constant.VolumeTestResource, "sshkey"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "bus", constant.VolumeResource+"."+constant.VolumeTestResource, "bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "availability_zone", constant.VolumeResource+"."+constant.VolumeTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "cpu_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "cpu_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "ram_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "ram_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "nic_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "nic_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "nic_hot_unplug", constant.VolumeResource+"."+constant.VolumeTestResource, "nic_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "disc_virtio_hot_plug", constant.VolumeResource+"."+constant.VolumeTestResource, "disc_virtio_hot_plug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "disc_virtio_hot_unplug", constant.VolumeResource+"."+constant.VolumeTestResource, "disc_virtio_hot_unplug"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "device_number", constant.VolumeResource+"."+constant.VolumeTestResource, "device_number"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.VolumeResource+"."+constant.VolumeDataSourceByName, "boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
				),
			},
			{
				Config:      testAccDataSourceVolumeWrongNameError,
				ExpectError: regexp.MustCompile(`no volume found with the specified criteria: name`),
			},
			{
				Config: testAccCheckVolumeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "size", "6"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.VolumeResource+"."+constant.VolumeTestResource, "image_name"),
					resource.TestCheckResourceAttrPair(constant.VolumeResource+"."+constant.VolumeTestResource, "boot_server", constant.ServerResource+"."+constant.ServerTestResource+"updated", "id"),
					resource.TestCheckResourceAttrPair(constant.VolumeResource+"."+constant.VolumeTestResource, "image_password", constant.RandomPassword+".server_image_password_updated", "result"),
					utils.TestImageNotNull(constant.VolumeResource, "image")),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeConfigNoPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(constant.VolumeResource+"."+constant.VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "name", constant.VolumeTestResource),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "size", "4"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "licence_type", "unknown"),
				)},
			{
				Config: testAccCheckVolumeConfigNoPasswordUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(constant.VolumeResource+"."+constant.VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "licence_type", "other"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(constant.VolumeResource+"."+constant.VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "name", constant.VolumeTestResource),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.VolumeResource+"."+constant.VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(constant.VolumeResource+"."+constant.VolumeTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					utils.TestImageNotNull(constant.VolumeResource, "image"))},
		},
	})
}

func testAccCheckVolumeDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.VolumeResource {
			continue
		}

		_, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
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
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching Volume: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		volume = &foundServer

		return nil
	}
}

const testAccCheckVolumeConfigBasic = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + `{
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
	datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
	server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
	availability_zone = "ZONE_1"
	name = "` + constant.VolumeTestResource + `"
	size = 5
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	image_password = ` + constant.RandomPassword + `.server_image_password.result
	user_data = "foo"
}
` + ServerImagePassword

const testAccCheckVolumeConfigBasicErrorNoPassOrSSHPath = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + `{
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
	datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
	server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
	availability_zone = "ZONE_1"
	name = "` + constant.VolumeTestResource + `"
	size = 5
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	user_data = "foo"
}
` + ServerImagePassword

const testAccCheckVolumeConfigUpdate = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + `updated {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
	datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
	server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `updated.id
	availability_zone = "ZONE_1"
	name = "` + constant.UpdatedResources + `"
	size = 6
	disk_type = "SSD Standard"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	image_password = ` + constant.RandomPassword + `.server_image_password_updated.result
	user_data = "foo"
}
` + ServerImagePassword + ServerImagePasswordUpdated

var testAccDataSourceVolumeMatchId = testAccCheckVolumeConfigBasic + `
data ` + constant.VolumeResource + ` ` + constant.VolumeDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.VolumeResource + `.` + constant.VolumeTestResource + `.id
}
`

var testAccDataSourceVolumeMatchName = testAccCheckVolumeConfigBasic + `
data ` + constant.VolumeResource + ` ` + constant.VolumeDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= ` + constant.VolumeResource + `.` + constant.VolumeTestResource + `.name
}
`

var testAccDataSourceVolumeWrongNameError = testAccCheckVolumeConfigBasic + `
data ` + constant.VolumeResource + ` ` + constant.VolumeDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckVolumeConfigNoPassword = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password =  ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name = "` + constant.VolumeTestResource + `"
  size           = 4
  disk_type      = "HDD"
  licence_type   = "unknown"
}
` + ServerImagePassword

const testAccCheckVolumeConfigNoPasswordUpdate = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name = "` + constant.UpdatedResources + `"
  size           = 5
  disk_type      = "HDD"
  licence_type   = "other"
}
` + ServerImagePassword

const testAccCheckVolumeResolveImageName = testAccCheckLanConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  availability_zone = "ZONE_1"
  name = "` + constant.VolumeTestResource + `"
  size = 5
  disk_type = "SSD Standard"
  bus = "VIRTIO"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
}
` + ServerImagePassword
