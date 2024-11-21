//go:build compute || all || nic

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

func TestAccNicBasic(t *testing.T) {
	var nic ionoscloud.Nic
	name := "testNic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNicDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNicConfigBasic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNICExists(constant.FullNicResourceName, &nic),
					resource.TestCheckResourceAttrSet(constant.FullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "name", name),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "dhcp", "true"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "dhcpv6", "true"),
					resource.TestCheckResourceAttrSet(constant.FullNicResourceName, "ipv6_cidr_block"),
					resource.TestCheckResourceAttrSet(constant.FullNicResourceName, "mac"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "firewall_type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "security_groups_ids.#", "2"),
					resource.TestCheckResourceAttrPair(constant.FullNicResourceName, "ips.0", "ionoscloud_ipblock.test_server", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.FullNicResourceName, "ips.1", "ionoscloud_ipblock.test_server", "ips.1"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.name", "test_flowlog"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.action", "ALL"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.direction", "INGRESS"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.bucket", constant.FlowlogBucket),
				),
			},
			{
				Config: testAccDataSourceNicMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "name", constant.FullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "dhcp", constant.FullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "dhcpv6", constant.FullNicResourceName, "dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_cidr_block", constant.FullNicResourceName, "ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "firewall_active", constant.FullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "firewall_type", constant.FullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "mac", constant.FullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "pci_slot", constant.FullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "lan", constant.FullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ips", constant.FullNicResourceName, "ips"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_ips", constant.FullNicResourceName, "ipv6_ips"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.name", constant.FullNicResourceName, "flowlog.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.action", constant.FullNicResourceName, "flowlog.0.action"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.direction", constant.FullNicResourceName, "flowlog.0.direction"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.bucket", constant.FullNicResourceName, "flowlog.0.bucket"),
				),
			},
			{
				Config: testAccDataSourceNicMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "name", constant.FullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "dhcp", constant.FullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "dhcpv6", constant.FullNicResourceName, "dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_cidr_block", constant.FullNicResourceName, "ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "firewall_active", constant.FullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "firewall_type", constant.FullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "mac", constant.FullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "pci_slot", constant.FullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "lan", constant.FullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ips", constant.FullNicResourceName, "ips"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_ips", constant.FullNicResourceName, "ipv6_ips"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.name", constant.FullNicResourceName, "flowlog.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.action", constant.FullNicResourceName, "flowlog.0.action"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.direction", constant.FullNicResourceName, "flowlog.0.direction"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "flowlog.0.bucket", constant.FullNicResourceName, "flowlog.0.bucket"),
				),
			},
			{
				Config:      testAccDataSourceNicMatchNameError,
				ExpectError: regexp.MustCompile(`no nic found with the specified criteria: name`),
			},
			{
				Config:      testAccDataSourceNicMatchIdAndNameError,
				ExpectError: regexp.MustCompile(`does not match expected name`),
			},
			{
				Config:      testAccDataSourceNicMultipleResultsError,
				ExpectError: regexp.MustCompile(`more than one nic found with the specified criteria: name`),
			},
			{
				Config: testAccCheckNicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "name", "updated"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "dhcp", "false"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "dhcpv6", "false"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "firewall_active", "false"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "security_groups_ids.#", "1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_cidr_block", constant.FullNicResourceName, "ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_ips", constant.FullNicResourceName, "ipv6_ips"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.name", "test_flowlog_updated"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.action", "REJECTED"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.direction", "EGRESS"),
					resource.TestCheckResourceAttr(constant.FullNicResourceName, "flowlog.0.bucket", constant.FlowlogBucketUpdated),
				),
			},
			{
				Config: testAccCheckNicConfigUpdateIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_cidr_block", constant.FullNicResourceName, "ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+dataSourceNicById, "ipv6_ips", constant.FullNicResourceName, "ipv6_ips"),
					resource.TestCheckNoResourceAttr(constant.FullNicResourceName, "flowlog.%"),
				),
			},
		},
	})
}

func testAccCheckNicDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NicResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		serverId := rs.Primary.Attributes["server_id"]

		_, apiResponse, err := client.NetworkInterfacesApi.
			DatacentersServersNicsFindById(ctx, dcId, serverId, rs.Primary.ID).
			Execute()

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of nic %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("nic %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckNICExists(n string, nic *ionoscloud.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching nic: %s %w", rs.Primary.ID, err)
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
  name = "nic-test"
  location = "us/las"
}
resource "ionoscloud_ipblock" "test_server" {
  location = ionoscloud_datacenter.test_datacenter.location
  size = 2
  name = "test_server_ipblock"
}
resource "ionoscloud_lan" "test_lan_1" {
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  name = "Lan 1"
  ipv6_cidr_block = "AUTO"
}
resource "ionoscloud_lan" "test_lan_2" {
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  name = "Lan 2"
  ipv6_cidr_block = "AUTO"
}
resource "ionoscloud_server" "test_server" {
  name = "test_server"
  datacenter_id = "${ionoscloud_datacenter.test_datacenter.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
}
  nic {
    lan = "${ionoscloud_lan.test_lan_1.id}"
    dhcp = true
    dhcpv6 = true
    firewall_active = true
  }
}
` + ServerImagePassword

const testAccCheckNicConfigBasic = testCreateDataCenterAndServer + `
resource ` + constant.NicResource + ` "database_nic" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = "${ionoscloud_lan.test_lan_2.id}"
  dhcpv6 = true
  firewall_active = true
  firewall_type = "INGRESS"
  ips = [ ionoscloud_ipblock.test_server.ips[0], ionoscloud_ipblock.test_server.ips[1] ]
  name = "%s"
  security_groups_ids   = [ionoscloud_nsg.example_1.id, ionoscloud_nsg.example_2.id]
  flowlog {
    name = "test_flowlog"
    action = "ALL"
    direction = "INGRESS"
    bucket = "` + constant.FlowlogBucket + `"
  }
}
` + testSecurityGroups

const testAccCheckNicConfigUpdate = testCreateDataCenterAndServer + `
resource ` + constant.NicResource + ` "database_nic" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = "${ionoscloud_lan.test_lan_2.id}"
  dhcp = false
  dhcpv6 = false
  firewall_active = false
  firewall_type = "BIDIRECTIONAL"
  ips = [ ionoscloud_ipblock.test_server.ips[0], ionoscloud_ipblock.test_server.ips[1] ]
  name = "updated"
  security_groups_ids   = [ionoscloud_nsg.example_1.id]
  ipv6_cidr_block = cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,12)
  ipv6_ips = [ 
                cidrhost(cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,12),1),
                cidrhost(cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,12),2),
                cidrhost(cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,12),3)
             ]
  flowlog {
    name = "test_flowlog_updated"
    action = "REJECTED"
    direction = "EGRESS"
    bucket = "` + constant.FlowlogBucketUpdated + `"
  }
}

data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  id = ` + constant.FullNicResourceName + `.id
}
` + testSecurityGroups

const testAccCheckNicConfigUpdateIpv6 = testCreateDataCenterAndServer + `
resource ` + constant.NicResource + ` "database_nic" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = "${ionoscloud_lan.test_lan_2.id}"
  dhcp = false
  dhcpv6 = false
  firewall_active = false
  firewall_type = "BIDIRECTIONAL"
  ips = [ ionoscloud_ipblock.test_server.ips[0], ionoscloud_ipblock.test_server.ips[1] ]
  name = "updated"
  ipv6_cidr_block = cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,16)
  ipv6_ips = [ 
                cidrhost(cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,16),10),
                cidrhost(cidrsubnet(ionoscloud_lan.test_lan_2.ipv6_cidr_block,16,16),30)
             ]
}

data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  id = ` + constant.FullNicResourceName + `.id
}
`

const dataSourceNicById = constant.NicResource + ".test_nic_data"

const testAccDataSourceNicMatchId = testAccCheckNicConfigBasic + `
data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  id = ` + constant.FullNicResourceName + `.id
}
`

const testAccDataSourceNicMatchName = testAccCheckNicConfigBasic + `
data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name = ` + constant.FullNicResourceName + `.name 
}`

const testAccDataSourceNicMatchNameError = testAccCheckNicConfigBasic + `
data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name = "DoesNotExist"
}`

const testAccDataSourceNicMatchIdAndNameError = testAccCheckNicConfigBasic + `
data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  id = ` + constant.FullNicResourceName + `.id
  name = "doesNotExist"
}`

const testAccDataSourceNicMultipleResultsError = testAccCheckNicConfigBasic + `
resource "ionoscloud_ipblock" "test_server_multiple_results" {
  location = ionoscloud_datacenter.test_datacenter.location
  size = 2
  name = "test_server_ipblock"
}

resource ` + constant.NicResource + ` "database_nic_multiple_results" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  lan = "${ionoscloud_lan.test_lan_2.id}"
  firewall_active = true
  firewall_type = "INGRESS"
  ips = [ ionoscloud_ipblock.test_server_multiple_results.ips[0], ionoscloud_ipblock.test_server_multiple_results.ips[1] ]
  name = "%s"
}

data ` + constant.NicResource + ` test_nic_data {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name = ` + constant.FullNicResourceName + `.name 
}`
