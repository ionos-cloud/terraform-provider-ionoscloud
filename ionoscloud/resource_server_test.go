package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServer_Basic(t *testing.T) {
	var server ionoscloud.Server
	serverName := "webserver"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckServerConfigBasic, serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("ionoscloud_server.webserver", &server),
					testAccCheckServerAttributes("ionoscloud_server.webserver", serverName),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "name", serverName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckServerConfigBasicDep, serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("ionoscloud_server.webserver", &server),
					testAccCheckServerAttributes("ionoscloud_server.webserver", serverName),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "name", serverName),
				),
			},
			{
				Config: testAccCheckServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerAttributes("ionoscloud_server.webserver", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "name", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "nic.0.firewall_active", "false"),
				),
			},
		},
	})
}

func TestAccServer_NoImage(t *testing.T) {
	var server ionoscloud.Server
	serverName := "webserver"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckserverconfigNoImage, serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("ionoscloud_server.webserver", &server),
					testAccCheckServerAttributes("ionoscloud_server.webserver", serverName),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "name", serverName),
				),
			},
		},
	})
}

func TestAccServer_BootCdromNoImage(t *testing.T) {
	var server ionoscloud.Server
	serverName := "webserver"
	bootCdromImageId := "83f21679-3321-11eb-a681-1e659523cb7b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckserverconfigBootCdromNoImage, serverName, bootCdromImageId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("ionoscloud_server.webserver", &server),
					testAccCheckServerAttributes("ionoscloud_server.webserver", serverName),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "name", serverName),
					resource.TestCheckResourceAttr("ionoscloud_server.webserver", "boot_cdrom", bootCdromImageId),
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
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("unable to fetch server %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("server still exists %s", rs.Primary.ID)

		}
	}

	return nil
}

func testAccCheckServerAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckServerAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
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
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "%s"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04-LTS-server-2021-06-01"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 15
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerConfigBasicDep = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "%s"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
	image_name = "Ubuntu-20.04-LTS-server-2021-06-01"
	image_password = "K3tTj8G14a3EgKyNeeiY"
    name = "system"
    size = 15
    disk_type = "HDD"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "updated"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04-LTS-server-2021-06-01"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 20
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = false
    firewall_active = false
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testacccheckserverconfigNoImage = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}
resource "ionoscloud_server" "webserver" {
  name = "%s"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
	licence_type = "OTHER"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testacccheckserverconfigBootCdromNoImage = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location   = "de/fra"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}
resource "ionoscloud_server" "webserver" {
  name = "%s"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  boot_cdrom = "%s" 
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
	licence_type = "OTHER"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

func Test_Update(t *testing.T) {

}
