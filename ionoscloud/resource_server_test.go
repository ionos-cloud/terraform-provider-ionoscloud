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
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cores", "1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "ram", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cpu_family", "AMD_OPTERON"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "image_password", "K3tTj8G14a3EgKyNeeiY"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", "system"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "5"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "volume.0.backup_unit_id", "ionoscloud_backup_unit.example", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.name", "system"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", UpdatedResources),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cores", "2"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "ram", "2048"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cpu_family", "AMD_OPTERON"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "image_password", "K3tTj8G14a3EgKyNeeiYsasad"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", UpdatedResources),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "6"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "volume.0.backup_unit_id", "ionoscloud_backup_unit.example", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.bus", "IDE"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.name", UpdatedResources),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name", UpdatedResources),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start", "21"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end", "23"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:18"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.2"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.3"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.type", "INGRESS"),
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
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cores", "1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "ram", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cpu_family", "INTEL_SKYLAKE"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "5"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end", "22"),
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
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cores", "1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "ram", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cpu_family", "INTEL_SKYLAKE"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "image_password", "pass123456"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "5"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end", "22"),
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
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cores", "1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "ram", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "cpu_family", "INTEL_SKYLAKE"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "5"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "true"),
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
					testAccCheckServerExists("ionoscloud_server."+ServerTestResourceName, &server),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "name", ServerTestResourceName),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "cores", "data.ionoscloud_template."+ServerTestResourceName, "cores"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "ram", "data.ionoscloud_template."+ServerTestResourceName, "ram"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "template_uuid", "data.ionoscloud_template."+ServerTestResourceName, "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "availability_zone", "ZONE_2"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "type", "CUBE"),
					testImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "volume.0.size", "data.ionoscloud_template."+ServerTestResourceName, "storage_size"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "volume.0.licence_type", "LINUX"),
					resource.TestCheckResourceAttrPair("ionoscloud_server."+ServerTestResourceName, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.name", ServerTestResourceName),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr("ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active", "true"),
				),
			},
		},
	})
}

func testAccCheckServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
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
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
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
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_ipblock" "webserver_ipblock_update" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + UpdatedResources + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
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
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
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
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location   = "de/fra"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}
resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  boot_cdrom = "` + bootCdromImageId + `" 
  volume {
    name = "` + ServerTestResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResourceName + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerResolveImageName = `
resource "ionoscloud_datacenter" "datacenter" {
  name        = "test_server"
  location    = "de/fra"
  description = "Test datacenter done by TF"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = true
}

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name              = "` + ServerTestResourceName + `"
  datacenter_id     = ionoscloud_datacenter.datacenter.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE" 
  image_name        = "Ubuntu-20.04-LTS"
  image_password    = "pass123456"
  volume {
    name           = "` + ServerTestResourceName + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ionoscloud_lan.webserver_lan.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResourceName + `"
      port_range_start = 22
      port_range_end   = 22
    }
  }
}`

const testAccCheckServerWithSnapshot = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "volume-test"
	location   = "de/fra"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}
resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = ionoscloud_datacenter.foobar.id
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
    lan = ionoscloud_lan.webserver_lan.id
    dhcp = true
    firewall_active = true
  }
}
resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  volume_id = ionoscloud_server.webserver.boot_volume
  name = "terraform_snapshot"
}
resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  depends_on = [ionoscloud_snapshot.test_snapshot]
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  image_name = "terraform_snapshot"
  volume {
    name = "` + ServerTestResourceName + `"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    dhcp = true
    firewall_active = true
  }
}
`

const testAccCheckCubeServer = `
data "ionoscloud_template" ` + ServerTestResourceName + ` {
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

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name              = "` + ServerTestResourceName + `"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  type              = "CUBE"
  template_uuid     = data.ionoscloud_template.` + ServerTestResourceName + `.id
  image_password = "K3tTj8G14a3EgKyNeeiY"  
  datacenter_id     = ionoscloud_datacenter.foobar.id
  volume {
    name            = "` + ServerTestResourceName + `"
    licence_type    = "LINUX" 
    disk_type = "DAS"
	}
  nic {
    lan             = ionoscloud_lan.webserver_lan.id
    name            = "` + ServerTestResourceName + `"
    dhcp            = true
    firewall_active = true
  }
}`
