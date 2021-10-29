package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const bootCdromImageId = "83f21679-3321-11eb-a681-1e659523cb7b"

func TestAccServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					testImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "image_password", "K3tTj8G14a3EgKyNeeiY"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.backup_unit_id", BackupUnitResource+"."+BackupUnitTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					testImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "image_password", "K3tTj8G14a3EgKyNeeiYsasad"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.backup_unit_id", BackupUnitResource+"."+BackupUnitTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.bus", "IDE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "21"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:18"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
				),
			},
		},
	})
}
func TestAccServerBootCdromNoImage(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBootCdromNoImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
				),
			},
		},
	})
}

func TestAccServerResolveImageName(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					testImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "image_password", "pass123456"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
				),
			},
		},
	})
}

func TestAccServerWithSnapshot(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckServerWithSnapshot),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					testImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
		},
	})
}

func TestAccServerCubeServer(t *testing.T) {

	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "cores", "data.ionoscloud_template."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "ram", "data.ionoscloud_template."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "template_uuid", "data.ionoscloud_template."+ServerTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "CUBE"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.size", "data.ionoscloud_template."+ServerTestResource, "storage_size"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.licence_type", "LINUX"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
		},
	})
}

func testAccCheckServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ServerResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("unable to fetch server %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("server still exists %s", rs.Primary.ID)

		}
	}

	return nil
}

func testAccCheckServerExists(n string, server *ionoscloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckServerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundServer, _, err := client.ServersApi.DatacentersServersFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Server: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		server = &foundServer

		return nil
	}
}

const testAccCheckServerConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
` + testAccCheckBackupUnitConfigBasic + `

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}`

const testAccCheckServerConfigUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
` + testAccCheckBackupUnitConfigBasic + `

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_ipblock" "webserver_ipblock_update" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + UpdatedResources + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiYsasad"
  type = "ENTERPRISE"
  volume {
    name = "` + UpdatedResources + `"
    size = 6
    disk_type = "SSD Standard"
	backup_unit_id = ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "` + UpdatedResources + `"
    dhcp = false
    firewall_active = false
    ips            = [ ionoscloud_ipblock.webserver_ipblock_update.ips[0], ionoscloud_ipblock.webserver_ipblock_update.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + UpdatedResources + `"
      port_range_start = 21
      port_range_end = 23
	  source_mac = "00:0a:95:9d:68:18"
	  source_ip = ionoscloud_ipblock.webserver_ipblock_update.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock_update.ips[3]
	  type = "INGRESS"
    }
  }
}`

const testAccCheckServerConfigBootCdromNoImage = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  boot_cdrom = "` + bootCdromImageId + `" 
  volume {
    name = "` + ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerResolveImageName = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/fra"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE" 
  image_name        = "Ubuntu-20.04-LTS"
  image_password    = "pass123456"
  volume {
    name           = "` + ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ` + LanResource + `.` + LanTestResource + `.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `"
      port_range_start = 22
      port_range_end   = 22
    }
  }
}`

const testAccCheckServerWithSnapshot = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` "webserver" {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
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
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  volume_id = ` + ServerResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  depends_on = [` + SnapshotResource + `.test_snapshot]
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  image_name = "terraform_snapshot"
  volume {
    name = "` + ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
`

const testAccCheckCubeServer = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location   = "de/txl"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  type              = "CUBE"
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  image_password = "K3tTj8G14a3EgKyNeeiY"  
  datacenter_id     = ionoscloud_datacenter.foobar.id
  volume {
    name            = "` + ServerTestResource + `"
    licence_type    = "LINUX" 
    disk_type = "DAS"
	}
  nic {
    lan             = ionoscloud_lan.webserver_lan.id
    name            = "` + ServerTestResource + `"
    dhcp            = true
    firewall_active = true
  }
}`
