package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVolume_Basic(t *testing.T) {
	var volume ionoscloud.Volume
	volumeName := "volume"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckVolumeConfig_basic, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ionoscloud_volume.database_volume", &volume),
					resource.TestCheckResourceAttr("ionoscloud_volume.database_volume", "name", volumeName),
				),
			},
			{
				Config: testAccCheckVolumeConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_volume.database_volume", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckVolumeDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		_, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response.StatusCode != 404 {
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
			return fmt.Errorf("Record not found")
		}

		volume = &foundServer

		return nil
	}
}

const testAccCheckVolumeConfig_basic = `
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
<<<<<<< HEAD
	image = "81e054dd-a347-11eb-b70c-7ade62b52cc0"
=======
	image_name = "Ubuntu-20.04-LTS-server-2021-05-01"
>>>>>>> master
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
  size = 14
  disk_type = "HDD"
  bus = "VIRTIO"
<<<<<<< HEAD
  image = "81e054dd-a347-11eb-b70c-7ade62b52cc0"
=======
  image_name = "Ubuntu-20.04-LTS-server-2021-05-01"
>>>>>>> master
  image_password = "K3tTj8G14a3EgKyNeeiY"
}`

const testAccCheckVolumeConfig_update = `
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
<<<<<<< HEAD
	image = "81e054dd-a347-11eb-b70c-7ade62b52cc0"
=======
	image_name = "Ubuntu-20.04-LTS-server-2021-05-01"
>>>>>>> master
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
  size = 14
  disk_type = "HDD"
  bus = "VIRTIO"
<<<<<<< HEAD
  image = "81e054dd-a347-11eb-b70c-7ade62b52cc0"
=======
  image_name = "Ubuntu-20.04-LTS-server-2021-05-01"
>>>>>>> master
  image_password = "K3tTj8G14a3EgKyNeeiY"
}`
