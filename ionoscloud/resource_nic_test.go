//go:build compute || all || nic

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"regexp"
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
		ExternalProviders: randomProviderVersion343(),
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
				Config: testAccDataSourceNicMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_type", fullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "pci_slot", fullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
				),
			},
			{
				Config: testAccDataSourceNicMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_type", fullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "pci_slot", fullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
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
		if rs.Type != NicResource {
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
			return fmt.Errorf("error occured while fetching nic: %s %w", rs.Primary.ID, err)
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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
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
` + ServerImagePassword

const testAccCheckNicConfigBasic = testCreateDataCenterAndServer + `
resource ` + NicResource + ` "database_nic" {
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
resource ` + NicResource + ` "database_nic" {
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
const dataSourceNicById = NicResource + ".test_nic_data"

const testAccDataSourceNicMatchId = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  id = ` + fullNicResourceName + `.id
}
`

const testAccDataSourceNicMatchName = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = ` + fullNicResourceName + `.name 
}`

const testAccDataSourceNicMatchNameError = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = "DoesNotExist"
}`

const testAccDataSourceNicMatchIdAndNameError = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	id = ` + fullNicResourceName + `.id
	name = "doesNotExist"
}`

const testAccDataSourceNicMultipleResultsError = testAccCheckNicConfigBasic + `
resource "ionoscloud_ipblock" "test_server_multiple_results" {
  location = ionoscloud_datacenter.test_datacenter.location
  size = 2
  name = "test_server_ipblock"
}

resource ` + NicResource + ` "database_nic_multiple_results" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  lan = 2
  firewall_active = true
  firewall_type = "INGRESS"
  ips = [ ionoscloud_ipblock.test_server_multiple_results.ips[0], ionoscloud_ipblock.test_server_multiple_results.ips[1] ]
  name = "%s"
}

data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = ` + fullNicResourceName + `.name 
}`
