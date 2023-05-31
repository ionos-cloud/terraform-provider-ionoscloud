//go:build compute || all || server

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

const bootCdromImageId = "83f21679-3321-11eb-a681-1e659523cb7b"

//ToDo: add backup unit back in tests when stable

func TestAccServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckServerNoPwdOrSSH,
				ExpectError: regexp.MustCompile(`either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided`),
			},
			{
				//ssh_key_path now accepts the ssh key directly too
				Config: testAccCheckServerSshDirectly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ssh_key_path.0", sshKey)),
			},
			{
				Config: testAccCheckServerSshKeysDirectly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ssh_keys.0", sshKey)),
			},
			{
				Config:      testAccCheckServerSshKeysAndKeyPathErr,
				ExpectError: regexp.MustCompile(`"ssh_keys": conflicts with ssh_key_path`),
			},
			{
				Config: testAccCheckServerNoNic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
				),
			},
			{
				Config: testAccCheckServerNoNicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
				),
			},
			{
				Config: testAccCheckServerConfigEmptyNicIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.id", ServerResource+"."+ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
				),
			},
			{
				Config: testAccCheckServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.id", ServerResource+"."+ServerTestResource, "primary_nic"),
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
				Config: testAccDataSourceServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "name", ServerResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "cores", ServerResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "ram", ServerResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "availability_zone", ServerResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "cpu_family", ServerResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "type", ServerResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.name", ServerResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.size", ServerResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.type", ServerResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.bus", ServerResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.availability_zone", ServerResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.lan", ServerResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.name", ServerResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.dhcp", ServerResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_active", ServerResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_type", ServerResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.ips.0", ServerResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.ips.1", ServerResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.protocol", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.name", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.type", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "name", ServerResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "cores", ServerResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "ram", ServerResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "availability_zone", ServerResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "cpu_family", ServerResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "type", ServerResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.name", ServerResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.size", ServerResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.type", ServerResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.bus", ServerResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.boot_server", ServerResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.availability_zone", ServerResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.lan", ServerResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.name", ServerResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.dhcp", ServerResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_active", ServerResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_type", ServerResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.ips.0", ServerResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.ips.1", ServerResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.name", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.type", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceServerWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
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
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
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

// issue #379
func TestAccServerNoBootVolumeBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigNoBootVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.licence_type", "UNKNOWN"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.id", ServerResource+"."+ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerConfigNoBootVolumeRemoveServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerAndVolumesDestroyed(DatacenterResource + "." + DatacenterTestResource),
				),
			},
		},
	})
}

// tests server with no cdromimage and with multiple firewall rules inline
func TestAccServerBootCdromNoImageAndInlineFwRules(t *testing.T) {
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
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "1"),
				),
			},
			{
				Config: testAccCheckServerConfig2Fw,
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
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.name", ServerTestResource+"2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "2"),
				),
			},
			{
				Config: testAccCheckServerConfig3Fw,
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
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.name", ServerTestResource+"2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.name", ServerTestResource+"3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "44"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerConfigRemove2FwRules,
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
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource+"3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "44"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerConfigRemoveAllFwRules,
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
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "0"),
				),
			},
		},
	})
}

// create and updates 5 inline rules after the server is created
func TestAccServerResolveImageNameAdd5FwRulesOnUpdate(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName5fwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "4"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "24"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "24"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_start", "26"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_end", "26"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName5fwRulesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "6"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "24"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "24"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_start", "26"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_end", "26"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_start", "27"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_end", "27"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.1.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.2.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.3.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.4.type", "EGRESS"),
				),
			},
		},
	})
}

// also tests creating 5 fw rules inline
func TestAccServerWithSnapshotAnd5FwRulesInline(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerWithSnapshot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "INTEL_SKYLAKE"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+"webserver", "nic.0.firewall.#", "5"),
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
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerAndServersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "cores", "data.ionoscloud_template."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "ram", "data.ionoscloud_template."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "template_uuid", "data.ionoscloud_template."+ServerTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "type", "CUBE"),
					utils.TestImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "volume.0.size", "data.ionoscloud_template."+ServerTestResource, "storage_size"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.licence_type", "LINUX"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.#", "1"),
				),
			},
		},
	})
}

func TestAccServerWithICMP(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerNoFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
				),
			},
			{
				Config: testAccCheckSeparateFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "cpu_family", "AMD_OPTERON"),
					utils.TestImageNotNull(ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(ServerResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
		},
	})
}

func TestAccServerWithLabels(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			// Clean server creation using labels in configuration.
			{
				Config: testAccCheckServerCreationWithLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.1.value", "labelvalue1"),
				),
			},
			// Check that labels are present in the server data source.
			{
				Config: testAccCheckDataSourceServerWithLabels,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServerResource+"."+ServerDataSourceById, "labels.#", "2"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerResource+"."+ServerDataSourceById, "labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerResource+"."+ServerDataSourceById, "labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerResource+"."+ServerDataSourceById, "labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerResource+"."+ServerDataSourceById, "labels.1.value", "labelvalue1"),
				),
			},
			// Update server labels.
			{
				Config: testAccCheckServerUpdateLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.0.key", "updatedlabelkey0"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.0.value", "updatedlabelvalue0"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.1.key", "updatedlabelkey1"),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.1.value", "updatedlabelvalue1"),
				),
			},
			// Delete server labels.
			{
				Config: testAccCheckServerDeleteLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerResource+"."+ServerTestResource, "label.#", "0"),
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
		if rs.Type != ServerResource {
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

func testAccCheckServerAndVolumesDestroyed(dcName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		datacenterResourceState, ok := s.RootModule().Resources[dcName]
		if !ok {
			return fmt.Errorf("can not get data center resource named: %s", dcName)
		}

		dcId := datacenterResourceState.Primary.ID

		// Since we are creating only ONE server in the data center, we can use
		// DatacentersServersGet to check if the server was deleted properly.
		servers, apiResponse, err := client.ServersApi.DatacentersServersGet(ctx, dcId).Execute()
		logApiRequestTime(apiResponse)
		if err == nil {
			if serverItems, ok := servers.GetItemsOk(); ok {
				if len(*serverItems) > 0 {
					return fmt.Errorf("server still exists for data center with ID: %s", dcId)
				}
			}
		} else {
			return err
		}

		volumes, apiResponse, err := client.VolumesApi.DatacentersVolumesGet(ctx, dcId).Execute()
		logApiRequestTime(apiResponse)
		if err == nil {
			if volItems, ok := volumes.GetItemsOk(); ok {
				if len(*volItems) > 0 {
					return fmt.Errorf("volumes still exists for data center with ID: %s", dcId)
				}
			}
		} else {
			return err
		}
		return nil
	}
}

func testAccCheckServerExists(serverName string, server *ionoscloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[serverName]

		if !ok {
			return fmt.Errorf("testAccCheckServerExists: Not found: %s", serverName)
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
			return fmt.Errorf("error occured while fetching Server: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		server = &foundServer

		return nil
	}
}

const testAccCheckServerConfigUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password_updated.result
  type = "ENTERPRISE"
  volume {
    name = "` + UpdatedResources + `"
    size = 6
    disk_type = "SSD Standard"
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
}
` + ServerImagePasswordUpdated

const testAccDataSourceServerMatchId = testAccCheckServerConfigBasic + `
data ` + ServerResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerResource + `.` + ServerTestResource + `.id
}
`

const testAccDataSourceServerMatchName = testAccCheckServerConfigBasic + `
data ` + ServerResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + ServerTestResource + `"
}
`
const testAccDataSourceServerWrongNameError = testAccCheckServerConfigBasic + `
data ` + ServerResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

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

resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id        = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id           = ` + ServerResource + `.` + ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
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
      type = "EGRESS"
      name = "` + ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerConfig2Fw = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id        = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id           = ` + ServerResource + `.` + ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
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
	  type = "EGRESS"
      port_range_start = 25
      port_range_end = 25
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[0]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[1]
    }
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResource + `2"
      port_range_start = 23
      port_range_end = 23
    }
  }
}`

const testAccCheckServerConfig3Fw = `
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
      type = "EGRESS"
      port_range_start = 25
      port_range_end = 25
    }
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResource + `2"
      port_range_start = 23
      port_range_end = 23
    }
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResource + `3"
      port_range_start = 44
      port_range_end = 44
    }
  }
}`

const testAccCheckServerConfigRemove2FwRules = `
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
      name = "` + ServerTestResource + `3"
      port_range_start = 44
      port_range_end = 44
    }
  }
}`

const testAccCheckServerConfigRemoveAllFwRules = `
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
  }
}`
const testAccCheckServerResolveImageNameNoNic = `
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
  image_name        = "ubuntu:latest"
  image_password    = ` + RandomPassword + `.server_image_password.result
  volume {
    name           = "` + ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

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
  image_name        = "ubuntu:latest"
  image_password    = ` + RandomPassword + `.server_image_password.result
  volume {
    name           = "` + ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ` + LanResource + `.` + LanTestResource + `.id
    dhcp            = true
    firewall_active = true
  }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerResolveImageName5fwRules = `
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
  image_name        = "ubuntu:latest"
  image_password    = ` + RandomPassword + `.server_image_password.result
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
      protocol         = "ICMP"
      name             = "` + ServerTestResource + `"
      type             = "INGRESS"
      icmp_code        = 4
      icmp_type        = 5
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 23
      port_range_end   = 23
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
	firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `5"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
  }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerResolveImageName5fwRulesUpdate = `
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
  image_name        = "ubuntu:latest"
  image_password    = ` + RandomPassword + `.server_image_password.result
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
      protocol         = "ICMP"
      name             = "` + ServerTestResource + `"
      type             = "INGRESS"
      icmp_code        = 5
      icmp_type        = 6
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
	firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `5"
      type             = "EGRESS"
      port_range_start = 27
      port_range_end   = 27
    }
  }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

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
	image_name = "ubuntu:latest"
	image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `"
      port_range_start = 22
      type             = "EGRESS"
      port_range_end   = 22
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 23
      port_range_end   = 23
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
	firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `5"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
  }
}
resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id        = ` + ServerResource + `.webserver.id
  nic_id           = ` + ServerResource + `.webserver.nic[0].id
  protocol         = "TCP"
  name             = "external_rule"
  port_range_start = 28
  port_range_end   = 28
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
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckCubeServerAndServersDataSource = `
data "ionoscloud_template" ` + ServerTestResource + ` {
   name = "CUBES XS"
   cores = 1
   ram   = 1024
   storage_size = 30
}

resource ` + DatacenterResource + " " + DatacenterTestResource + `{
	name       = "volume-test"
	location   = "de/txl"
}

resource "ionoscloud_lan" "webserver_lan" {
 datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
 public = true
 name = "public"
}

resource "ionoscloud_server" ` + ServerTestResource + ` {
 name              = "` + ServerTestResource + `"
 availability_zone = "ZONE_2"
 image_name        = "ubuntu:latest"
 type              = "CUBE"
 template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
 image_password = ` + RandomPassword + `.server_image_password.result
 datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
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
}
data ` + ServersDataSource + ` ` + ServerDataSourceByName + ` {
depends_on = [` + ServerResource + `.` + ServerTestResource + `]
 datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
 filter {
  name = "type"
  value = "CUBE"
 }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerNoFirewall = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`
const testAccCheckSeparateFirewall = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan             = ` + LanResource + `.` + LanTestResource + `.id
    name 			= "system"
    dhcp            = true
    firewall_active = true
    }
}
resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id           = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id              = ` + ServerResource + `.` + ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerICMP = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan             = ` + LanResource + `.` + LanTestResource + `.id
    name 			= "system"
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "ICMP"
      name             = "` + ServerTestResource + `"
      icmp_type        = "12"
      icmp_code        = "0"
	  }
    }
}
resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id           = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id              = ` + ServerResource + `.` + ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckDataSourceServerWithLabels = testAccCheckServerCreationWithLabels + `
data ` + ServerResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerResource + `.` + ServerTestResource + `.id
}
`

const testAccCheckServerUpdateLabels = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
  image_name ="ubuntu:latest"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
  label {
    key = "updatedlabelkey0"
    value = "updatedlabelvalue0"
  }
  label {
    key = "updatedlabelkey1"
    value = "updatedlabelvalue1"
  }
}`

const testAccCheckServerDeleteLabels = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
  image_name ="ubuntu:latest"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
}`

const testAccCheckServerConfigNoBootVolume = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  type = "ENTERPRISE"
  
  volume {
    name = "system"
    size = 6
    licence_type = "UNKNOWN"
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "INGRESS"
  }
}

resource "ionoscloud_volume" "exampleVol1" {
  datacenter_id           = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id               = ` + ServerResource + `.` + ServerTestResource + `.id
  name                    = "Another Volume Example"
  availability_zone       = "ZONE_1"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}
`
const testAccCheckServerConfigNoBootVolumeRemoveServer = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
`
