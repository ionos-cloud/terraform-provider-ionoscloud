package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolume_Basic(t *testing.T) {
	var volume ionoscloud.Volume
	volumeName := "volume"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckVolumeConfigBasic, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ionoscloud_volume.database_volume", &volume),
					resource.TestCheckResourceAttr("ionoscloud_volume.database_volume", "name", volumeName),
				),
			},
			{
				Config: testAccCheckVolumeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_volume.database_volume", "name", "updated"),
				),
			},
		},
	})
}

func TestAccVolume_NoPassword(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckvolumeconfigNoPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ionoscloud_volume.no_password_volume", &volume),
					resource.TestCheckResourceAttr("ionoscloud_volume.no_password_volume", "name", "no_password"),
				),
			},
		},
	})
}

func TestAccVolume_ResolveImageName(t *testing.T) {
	var volume ionoscloud.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckVolumeResolveImageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ionoscloud_volume.image_name_volume", &volume),
					resource.TestCheckResourceAttr("ionoscloud_volume.image_name_volume", "name", "image_name_volume"),
				),
			},
		},
	})
}

func TestAccVolume_ImageAlias(t *testing.T) {
	var volume ionoscloud.Volume
	volumeName := "volume"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckVolumeConfigImageAlias, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ionoscloud_volume.database_volume", &volume),
					resource.TestCheckResourceAttr("ionoscloud_volume.database_volume", "name", volumeName),
					testImageNotNull("ionoscloud_volume", "image_name"),
					testImageNotNull("ionoscloud_server", "boot_image"),
				),
			},
		},
	})
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

func testAccCheckVolumeDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}
		_, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("unable to fetch volume %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("volume %s still exists", rs.Primary.ID)
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

		foundServer, _, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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

const testAccCheckVolumeConfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
    disk_type = "HDD"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "database_volume" {
	datacenter_id = "${ionoscloud_datacenter.foobar.id}"
	server_id = "${ionoscloud_server.webserver.id}"
	availability_zone = "ZONE_1"
	name = "%s"
	size = 5
	disk_type = "HDD"
	bus = "VIRTIO"
	image_name ="Debian-10-cloud-init.qcow2"
	image_password = "K3tTj8G14a3EgKyNeeiY"
	backup_unit_id = ionoscloud_backup_unit.example.id
	user_data = "foo"
}`

const testAccCheckVolumeConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}


resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "database_volume" {
	datacenter_id = "${ionoscloud_datacenter.foobar.id}"
	server_id = "${ionoscloud_server.webserver.id}"
	availability_zone = "ZONE_1"
	name = "updated"
	size = 5
	disk_type = "HDD"
	bus = "VIRTIO"
	image_name ="Debian-10-cloud-init.qcow2"
	image_password = "K3tTj8G14a3EgKyNeeiY"
	backup_unit_id = ionoscloud_backup_unit.example.id
	user_data = "foo"
}`

const testacccheckvolumeconfigNoPassword = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "Ubuntu-20.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "no_password_volume" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  name = "no_password"
  size           = 4
  disk_type      = "HDD"
  licence_type   =  "other"
}`

const testAccCheckVolumeResolveImageName = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location   = "de/fra"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}
resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
	image_name = "Ubuntu-20.04-LTS"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}
resource "ionoscloud_volume" "image_name_volume" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  availability_zone = "ZONE_1"
  name = "image_name_volume"
  size = 5
  disk_type = "SSD Standard"
  bus = "VIRTIO"
  image_name = "Ubuntu-20.04-LTS"
  image_password = "K3tTj8G14a3EgKyNeeiY"
}
`

const testAccCheckVolumeConfigImageAlias = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location = "us/las"
}


resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
    disk_type = "HDD"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "database_volume" {
	datacenter_id = "${ionoscloud_datacenter.foobar.id}"
	server_id = "${ionoscloud_server.webserver.id}"
	availability_zone = "ZONE_1"
	name = "%s"
	size = 5
	disk_type = "HDD"
	bus = "VIRTIO"
	image_name ="ubuntu:latest"
	image_password = "K3tTj8G14a3EgKyNeeiY"
}`
