//go:build compute || all || firewall

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFirewallBasic(t *testing.T) {
	var firewall ionoscloud.FirewallRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(constant.FirewallResource+"."+constant.FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", constant.FirewallTestResource),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "name", constant.FirewallResource+"."+constant.FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "protocol", constant.FirewallResource+"."+constant.FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "source_mac", constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "source_ip", constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "target_ip", constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "icmp_type", constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "icmp_code", constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_code"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceById, "type", constant.FirewallResource+"."+constant.FirewallTestResource, "type"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "name", constant.FirewallResource+"."+constant.FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "protocol", constant.FirewallResource+"."+constant.FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "source_mac", constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "source_ip", constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "target_ip", constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "icmp_type", constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "icmp_code", constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_code"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.FirewallResource+"."+constant.FirewallDataSourceByName, "type", constant.FirewallResource+"."+constant.FirewallTestResource, "type"),
				),
			},
			{
				Config:      testAccDataSourceFirewallMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one firewall rule found with the specified criteria name"),
			},
			{
				Config:      testAccDataSourceFirewallWrongNameError,
				ExpectError: regexp.MustCompile("no firewall rule found with the specified name"),
			},
			{
				Config: testAccCheckFirewallConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(constant.FirewallResource+"."+constant.FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_type", "2"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_code", "7"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckFirewallSetICMPToZero,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_type", "0"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "icmp_code", "0"),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(constant.FirewallResource+"."+constant.FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", constant.FirewallTestResource),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "UDP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "port_range_start", "80"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "port_range_end", "80"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckFirewallConfigUpdateUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists(constant.FirewallResource+"."+constant.FirewallTestResource, &firewall),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "UDP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "source_ip", "ionoscloud_ipblock.ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.FirewallResource+"."+constant.FirewallTestResource, "target_ip", "ionoscloud_ipblock.ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "port_range_start", "81"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "port_range_end", "82"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "EGRESS"),
				),
			},
		},
	})
}

func testAccCheckFirewallDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.FirewallResource {
			continue
		}

		_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
			rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("error occurent at checking deletion of firewall %s %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("firewall still exists %s %w", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckFirewallExists(n string, firewall *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching Firewall rule: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		firewall = &foundServer

		return nil
	}
}

const testAccCheckFirewallConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
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

` + ServerImagePassword + `

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}

resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + constant.FirewallTestResource + `"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
  type 		= "INGRESS"
}
`

const testAccCheckFirewallConfigUpdate = testAccCheckDatacenterConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
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
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + `  {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + constant.UpdatedResources + `"
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  icmp_type = 2
  icmp_code = 7
  type = "EGRESS"
}
`

const testAccCheckFirewallConfigUDP = testAccCheckDatacenterConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
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
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + `  {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "UDP"
  name = "` + constant.FirewallTestResource + `"
  port_range_start = 80
  port_range_end = 80
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  type = "INGRESS"
}
`

const testAccCheckFirewallConfigUpdateUDP = testAccCheckDatacenterConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
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
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "UDP"
  name = "` + constant.UpdatedResources + `"
  port_range_start = 81
  port_range_end = 82
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  type = "EGRESS"
}
`

const testAccCheckFirewallSetICMPToZero = testAccCheckDatacenterConfigBasic + `
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
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
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}
resource "ionoscloud_ipblock" "ipblock_update" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "firewall_ipblock"
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + `  {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + constant.UpdatedResources + `"
  source_mac = "00:0a:95:9d:68:17"
  source_ip = ionoscloud_ipblock.ipblock_update.ips[0]
  target_ip = ionoscloud_ipblock.ipblock_update.ips[1]
  icmp_type = 0
  icmp_code = 0
}
`
const testAccDataSourceFirewallMatchId = testAccCheckFirewallConfigBasic + `
data ` + constant.FirewallResource + ` ` + constant.FirewallDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  id = ` + constant.FirewallResource + `.` + constant.FirewallTestResource + `.id
}
`

const testAccDataSourceFirewallMatchName = testAccCheckFirewallConfigBasic + `
data ` + constant.FirewallResource + ` ` + constant.FirewallDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "` + constant.FirewallTestResource + `"
}
`

const testAccDataSourceFirewallMultipleResultsError = testAccCheckFirewallConfigBasic + `
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + `_multiple_results  {
datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "` + constant.FirewallTestResource + `"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
  type = "INGRESS"
}

data ` + constant.FirewallResource + ` ` + constant.FirewallDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "` + constant.FirewallTestResource + `"
}
`

const testAccDataSourceFirewallWrongNameError = testAccCheckFirewallConfigBasic + `
data ` + constant.FirewallResource + ` ` + constant.FirewallDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "wrong_name"
}
`
