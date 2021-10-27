package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
				Config: testAccCheckVolumeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(VolumeResource+"."+VolumeTestResource, &volume),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", VolumeTestResource),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "5"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "image_name", "Debian-10-cloud-init.qcow2"),
					testImageNotNull(VolumeResource, "image")),
			},
			{
				Config: testAccCheckVolumeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "size", "6"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "bus", "VIRTIO"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(VolumeResource+"."+VolumeTestResource, "image_name", "Debian-10-cloud-init.qcow2"),
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
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != VolumeResource {
			continue
		}

		_, apiResponse, err := client.VolumeApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

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

		foundServer, _, err := client.VolumeApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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

func testImageNotNull(resource, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource {
				continue
			}

			image := rs.Primary.Attributes[attribute]

			if image == "" {
				return fmt.Errorf("%s is empty, expected an UUID", attribute)
			} else if !IsValidUUID(image) {
				return fmt.Errorf("%s should be a valid UUID, got: %#v", attribute, image)
			}

		}
		return nil
	}
}

const testAccCheckVolumeConfigBasic = testAccCheckLanConfigBasic + testAccCheckBackupUnitConfigBasic + `
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
	image_name ="Debian-10-cloud-init.qcow2"
	image_password = "K3tTj8G14a3EgKyNeeiY"
	backup_unit_id = ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
	user_data = "foo"
}`

const testAccCheckVolumeConfigUpdate = testAccCheckBackupUnitConfigBasic + testAccCheckLanConfigBasic + `
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
	image_name ="Debian-10-cloud-init.qcow2"
	image_password = "K3tTj8G14a3EgKyNeeiYupdated"
	backup_unit_id = ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
	user_data = "foo"
}`

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
  size = 6
  disk_type = "SSD Standard"
  bus = "VIRTIO"
  image_name = "Ubuntu-20.04-LTS"
  image_password = "K3tTj8G14a3EgKyNeeiY"
}
`
