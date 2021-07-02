package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNic_Basic(t *testing.T) {
	var nic ionoscloud.Nic
	volumeName := "volume"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNicDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNicConfigBasic, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNICExists("ionoscloud_nic.database_nic", &nic),
					testAccCheckNicAttributes("ionoscloud_nic.database_nic", volumeName),
					resource.TestCheckResourceAttrSet("ionoscloud_nic.database_nic", "mac"),
					resource.TestCheckResourceAttr("ionoscloud_nic.database_nic", "name", volumeName),
				),
			},
			{
				Config: testAccCheckNicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNicAttributes("ionoscloud_nic.database_nic", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_nic.database_nic", "name", "updated"),
				),
			},
		},
	})
}

func TestAccNic_Ips(t *testing.T) {
	var nic ionoscloud.Nic
	publicIp1 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_1")
	if publicIp1 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_1 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}
	publicIp2 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_2")
	if publicIp2 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_2 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNicDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testaccchecknicconfigIps, publicIp1, publicIp2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNICExists("ionoscloud_nic.database_nic", &nic),
					resource.TestCheckResourceAttrSet("ionoscloud_nic.database_nic", "mac"),
					resource.TestCheckResourceAttr("ionoscloud_nic.database_nic", "ips.0", publicIp1),
					resource.TestCheckResourceAttr("ionoscloud_nic.database_nic", "ips.1", publicIp2),
				),
			},
		},
	})
}

func testAccCheckNicDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_nic" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		serverId := rs.Primary.Attributes["server_id"]

		_, apiResponse, err := client.NetworkInterfacesApi.
			DatacentersServersNicsFindById(ctx, dcId, serverId, rs.Primary.ID).
			Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of nic %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("nic %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckNicAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckNicAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckNICExists(n string, nic *ionoscloud.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckVolumeExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}
		dcId := rs.Primary.Attributes["datacenter_id"]
		serverId := rs.Primary.Attributes["server_id"]
		foundNic, _, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, serverId, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Volume: %s", rs.Primary.ID)
		}
		if *foundNic.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		nic = &foundNic

		return nil
	}
}

const testAccCheckNicConfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "nic-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04-LTS-server-2021-07-01"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
    disk_type = "SSD"

}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = false
  firewall_active = true
  firewall_type = "INGRESS"
  name = "%s"
}`

const testAccCheckNicConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "nic-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04-LTS-server-2021-07-01"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
    disk_type = "SSD"
}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = false
  firewall_active = true
  firewall_type = "INGRESS"
  name = "updated"
}
`

const testaccchecknicconfigIps = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "nic-test"
	location = "de/fra"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
	image_name ="Ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"

}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = false
  firewall_active = true
  ips            = ["%s","%s"]
  name = "test_nic"
}`
