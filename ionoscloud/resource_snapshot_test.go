package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSnapshot_Basic(t *testing.T) {
	var snapshot ionoscloud.Snapshot
	snapshotName := "terraform_snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSnapshotConfig_basic, snapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("ionoscloud_snapshot.test_snapshot", &snapshot),
					resource.TestCheckResourceAttr("ionoscloud_snapshot.test_snapshot", "name", snapshotName),
				),
			},
			{
				Config: testAccCheckSnapshotConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_snapshot.test_snapshot", "name", snapshotName),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, _ := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_snapshot" {
			continue
		}

		_, apiResponse, err := client.SnapshotApi.SnapshotsFindById(ctx, rs.Primary.ID).Execute()

		if apiError, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("Snapshot still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Snapshot %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckSnapshotExists(n string, snapshot *ionoscloud.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckSnapshotExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		ctx, _ := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		foundServer, _, err := client.SnapshotApi.SnapshotsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching Snapshot: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		snapshot = &foundServer

		return nil
	}
}

const testAccCheckSnapshotConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "snapshot-test"
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
	image_name = "debian:9"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  volume_id = "${ionoscloud_server.webserver.boot_volume}"
  name = "%s"
}
`

const testAccCheckSnapshotConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "snapshot-test"
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
	image_name = "debian:9"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  volume_id = "${ionoscloud_server.webserver.boot_volume}"
  name = "terraform_snapshot"
}`
