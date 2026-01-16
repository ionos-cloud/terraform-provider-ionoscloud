//go:build all || gpu

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGpuServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckGpuServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGpuServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGpuServerExists(constant.ServerGPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "hostname", constant.ServerTestHostname),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerGPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Premium"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.expose_serial", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerGPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "INGRESS"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttrSet(constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.require_legacy_bios"),
				),
			},
			{
				Config: testAccDataSourceGpuServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "name", constant.ServerGPUResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "hostname", constant.ServerGPUResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "availability_zone", constant.ServerGPUResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "type", constant.ServerGPUResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.type", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.bus", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.availability_zone", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.boot_server", constant.ServerGPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.lan", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.dhcp", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_active", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_type", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.ips.0", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.protocol", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceById, "volumes.0.require_legacy_bios", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.require_legacy_bios"),
				),
			},
			{
				Config: testAccDataSourceGpuServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "name", constant.ServerGPUResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "hostname", constant.ServerGPUResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "availability_zone", constant.ServerGPUResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "type", constant.ServerGPUResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.type", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.bus", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.availability_zone", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.boot_server", constant.ServerGPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.lan", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.dhcp", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_active", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_type", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.ips.0", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.name", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerGPUResource+"."+constant.ServerDataSourceByName, "volumes.0.require_legacy_bios", constant.ServerGPUResource+"."+constant.ServerTestResource, "volume.0.require_legacy_bios"),
				),
			},
			{
				Config:      testAccDataSourceGpuServerWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},
			{
				Config: testAccCheckGpuServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGpuServerExists(constant.ServerGPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "hostname", "updatedhostname"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerGPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerGPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
				),
			},
			{
				Config: testAccDataSourceGpus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.id"),
					resource.TestCheckResourceAttrSet(constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.name"),
					resource.TestCheckResourceAttrSet(constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.vendor"),
					resource.TestCheckResourceAttrSet(constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.type"),
					resource.TestCheckResourceAttrSet(constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.model"),
				),
			},
			{
				Config: testAccDataSourceGpu,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+".ionoscloud_gpu.test_gpu", "id", constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+".ionoscloud_gpu.test_gpu", "name", constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+".ionoscloud_gpu.test_gpu", "vendor", constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.vendor"),
					resource.TestCheckResourceAttrPair(constant.DataSource+".ionoscloud_gpu.test_gpu", "type", constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+".ionoscloud_gpu.test_gpu", "model", constant.DataSource+".ionoscloud_gpus.test_gpus", "gpus.0.model"),
				),
			},
		},
	})
}

func testAccCheckGpuServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerGPUResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("unable to fetch server %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("server still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGpuServerExists(n string, server *ionoscloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckGpuServerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}

		foundServer, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching Server: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		server = &foundServer

		return nil
	}
}

// note: reserving IPs in de/fra/2 is yet not expected , de/fra is supposed to be used instead.
const testAccCheckGpuServerConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "gpu-server-test"
	location = "de/fra/2"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "de/fra"
  size = 1
  name = "webserver_ipblock"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerGPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  template_uuid = "6913ed82-a143-4c15-89ac-08fb375a97c5"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  vm_state = "RUNNING"

  volume {
    name = "system"
    licence_type = "LINUX"
    disk_type = "SSD Premium"
    bus = "VIRTIO"
    availability_zone = "AUTO"
    expose_serial = true
    require_legacy_bios = false
  }

  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "INGRESS"
    ips = [ionoscloud_ipblock.webserver_ipblock.ips[0]]

    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
      type = "INGRESS"
    }
  }
}
` + ServerImagePassword

const testAccCheckGpuServerConfigUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
 name       = "gpu-server-test"
 location = "de/fra/2"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "de/fra"
  size = 1
  name = "webserver_ipblock"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerGPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  hostname = "updatedhostname"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  template_uuid = "6913ed82-a143-4c15-89ac-08fb375a97c5"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password_updated.result
  vm_state = "RUNNING"

  volume {
    name = "system"
    licence_type = "LINUX"
    disk_type = "SSD Premium"
    bus = "VIRTIO"
    availability_zone = "AUTO"
    expose_serial = true
    require_legacy_bios = false
  }

  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.UpdatedResources + `"
    dhcp = true
    firewall_active = true
    firewall_type = "INGRESS"
    ips = [ionoscloud_ipblock.webserver_ipblock.ips[0]]

    firewall {
      protocol = "TCP"
      name = "` + constant.UpdatedResources + `"
      port_range_start = 22
      port_range_end = 22
      type = "INGRESS"
    }
  }
}
` + ServerImagePasswordUpdated

const testAccDataSourceGpuServerMatchId = testAccCheckGpuServerConfigBasic + `
data ` + constant.ServerGPUResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerGPUResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceGpuServerMatchName = testAccCheckGpuServerConfigBasic + `
data ` + constant.ServerGPUResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name = "` + constant.ServerTestResource + `"
}
`

const testAccDataSourceGpuServerWrongNameError = testAccCheckGpuServerConfigBasic + `
data ` + constant.ServerGPUResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name = "wrong_name"
}
`

const testAccDataSourceGpus = testAccCheckGpuServerConfigBasic + `
data "ionoscloud_gpus" test_gpus {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerGPUResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceGpu = testAccCheckGpuServerConfigBasic + `
data "ionoscloud_gpus" test_gpus {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerGPUResource + `.` + constant.ServerTestResource + `.id
}

data "ionoscloud_gpu" test_gpu {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerGPUResource + `.` + constant.ServerTestResource + `.id
  id = data.ionoscloud_gpus.test_gpus.gpus.0.id
}
`
