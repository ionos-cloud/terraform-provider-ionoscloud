//go:build compute || all || firewall

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewallBasic(t *testing.T) {
	var firewall ionoscloud.FirewallRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(FirewallResource+"."+FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", FirewallTestResource),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock", "ips.1"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_type", "1"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_code", "8"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "name", FirewallResource+"."+FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "protocol", FirewallResource+"."+FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "source_mac", FirewallResource+"."+FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "source_ip", FirewallResource+"."+FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "target_ip", FirewallResource+"."+FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "icmp_type", FirewallResource+"."+FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "icmp_code", FirewallResource+"."+FirewallTestResource, "icmp_code"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "type", FirewallResource+"."+FirewallTestResource, "type"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "name", FirewallResource+"."+FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "protocol", FirewallResource+"."+FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "source_mac", FirewallResource+"."+FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "source_ip", FirewallResource+"."+FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "target_ip", FirewallResource+"."+FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "icmp_type", FirewallResource+"."+FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "icmp_code", FirewallResource+"."+FirewallTestResource, "icmp_code"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "type", FirewallResource+"."+FirewallTestResource, "type"),
				),
			},
			{
				Config:      testAccDataSourceFirewallWrongName,
				ExpectError: regexp.MustCompile("no firewall found with the specified name"),
			},
			{
				Config: testAccCheckFirewallConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(FirewallResource+"."+FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_type", "2"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_code", "7"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckFirewallSetICMPToZero,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_type", "0"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "icmp_code", "0"),
				),
			},
		},
	})
}

func TestAccFirewallUDP(t *testing.T) {
	var firewall ionoscloud.FirewallRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(FirewallResource+"."+FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", FirewallTestResource),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "UDP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock", "ips.1"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "port_range_start", "80"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "port_range_end", "80"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckFirewallConfigUpdateUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(FirewallResource+"."+FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "UDP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(FirewallResource+"."+FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "port_range_start", "81"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "port_range_end", "82"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "EGRESS"),
				),
			},
		},
	})
}

func testAccCheckFirewallDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != FirewallResource {
			continue
		}

		_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
			rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("error occurent at checking deletion of firewall %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("firewall still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckFirewallExists(n string, firewall *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckFirewallExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundServer, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
			rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Firewall rule: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		firewall = &foundServer

		return nil
	}
}

const testAccCheckFirewallConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
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
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}

resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + FirewallTestResource + `"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
  type = "INGRESS"
}
`

const testAccCheckFirewallConfigUpdate = testAccCheckDatacenterConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "test1234"
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
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + FirewallResource + ` ` + FirewallTestResource + `  {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + UpdatedResources + `"
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  icmp_type = 2
  icmp_code = 7
  type = "EGRESS"
}
`

const testAccCheckFirewallConfigUDP = testAccCheckDatacenterConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "test1234"
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
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + FirewallResource + ` ` + FirewallTestResource + `  {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "UDP"
  name = "` + FirewallTestResource + `"
  port_range_start = 80
  port_range_end = 80
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  type = "INGRESS"
}
`

const testAccCheckFirewallConfigUpdateUDP = testAccCheckDatacenterConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "test1234"
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
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "UDP"
  name = "` + UpdatedResources + `"
  port_range_start = 81
  port_range_end = 82
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  type = "EGRESS"
}
`

const testAccCheckFirewallSetICMPToZero = testAccCheckDatacenterConfigBasic + `
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "Ubuntu-20.04"
  image_password = "test1234"
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
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}
resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + FirewallResource + ` ` + FirewallTestResource + `  {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + UpdatedResources + `"
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  icmp_type = 0
  icmp_code = 0
}
`
const testAccDataSourceFirewallMatchId = testAccCheckFirewallConfigBasic + `
data ` + FirewallResource + ` ` + FirewallDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  id = ` + FirewallResource + `.` + FirewallTestResource + `.id
}
`

const testAccDataSourceFirewallMatchName = testAccCheckFirewallConfigBasic + `
data ` + FirewallResource + ` ` + FirewallDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "` + FirewallTestResource + `"
}
`

const testAccDataSourceFirewallWrongName = testAccCheckFirewallConfigBasic + `
data ` + FirewallResource + ` ` + FirewallDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "wrong_name"
}
`
