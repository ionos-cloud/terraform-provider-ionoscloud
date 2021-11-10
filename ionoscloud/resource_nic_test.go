package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNicBasic(t *testing.T) {
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
					testAccCheckNICExists(fullNicResourceName, &nic),
					resource.TestCheckResourceAttrSet(fullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttr(fullNicResourceName, "name", volumeName),
					resource.TestCheckResourceAttr(fullNicResourceName, "dhcp", "true"),
					resource.TestCheckResourceAttrSet(fullNicResourceName, "mac"),
					resource.TestCheckResourceAttr(fullNicResourceName, "firewall_active", "true"),
					resource.TestCheckResourceAttr(fullNicResourceName, "firewall_type", "INGRESS"),
					resource.TestCheckResourceAttrPair(fullNicResourceName, "ips.0", "ionoscloud_ipblock.test_server", "ips.0"),
					resource.TestCheckResourceAttrPair(fullNicResourceName, "ips.1", "ionoscloud_ipblock.test_server", "ips.1"),
				),
			},
			{
				Config: testAccCheckNicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullNicResourceName, "name", "updated"),
					resource.TestCheckResourceAttr(fullNicResourceName, "dhcp", "false"),
					resource.TestCheckResourceAttr(fullNicResourceName, "firewall_active", "false"),
					resource.TestCheckResourceAttr(fullNicResourceName, "firewall_type", "BIDIRECTIONAL"),
				),
			},
		},
	})
}

func testAccCheckNicDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != nicResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		serverId := rs.Primary.Attributes["server_id"]

		_, apiResponse, err := client.NetworkInterfacesApi.
			DatacentersServersNicsFindById(ctx, dcId, serverId, rs.Primary.ID).
			Execute()

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of nic %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("nic %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckNICExists(n string, nic *ionoscloud.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
		foundNic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, serverId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

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

const testCreateDataCenterAndServer = `
resource "ionoscloud_datacenter" "test_datacenter" {
	name       = "nic-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "test_server" {
  location = ionoscloud_datacenter.test_datacenter.location
  size = 2
  name = "test_server_ipblock"
}
resource "ionoscloud_server" "test_server" {
  name = "test_server"
  datacenter_id = "${ionoscloud_datacenter.test_datacenter.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
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
`

const testAccCheckNicConfigBasic = testCreateDataCenterAndServer + `
resource "ionoscloud_nic" "database_nic" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  firewall_active = true
  firewall_type = "INGRESS"
  ips = [ ionoscloud_ipblock.test_server.ips[0], ionoscloud_ipblock.test_server.ips[1] ]
  name = "%s"
}
`

const testAccCheckNicConfigUpdate = testCreateDataCenterAndServer + `

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = false
  firewall_active = false
  firewall_type = "BIDIRECTIONAL"
  ips = [ ionoscloud_ipblock.test_server.ips[0], ionoscloud_ipblock.test_server.ips[1] ]
  name = "updated"
}
`
