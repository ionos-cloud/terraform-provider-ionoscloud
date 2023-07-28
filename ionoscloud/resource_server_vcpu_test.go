//go:build compute || all || server || vcpu
// +build compute all server vcpu

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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

//ToDo: add backup unit back in tests when stable

func TestAccServerVCPUBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckServerVCPUNoPwdOrSSH,
				ExpectError: regexp.MustCompile(`either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided`),
			},
			{
				Config: testAccCheckServerVCPUSshKeysDirectly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ssh_keys.0", sshKey)),
			},
			{
				Config: testAccCheckServerVCPUNoNic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
				),
			},
			{
				Config: testAccCheckServerVCPUNoNicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigEmptyNicIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "volume.0.boot_server", ServerVCPUResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.id", ServerVCPUResource+"."+ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "volume.0.boot_server", ServerVCPUResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.id", ServerVCPUResource+"."+ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccDataSourceServerVCPUMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "name", ServerVCPUResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "cores", ServerVCPUResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "ram", ServerVCPUResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "availability_zone", ServerVCPUResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrSet(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "type", ServerVCPUResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.name", ServerVCPUResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.size", ServerVCPUResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.type", ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.bus", ServerVCPUResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.availability_zone", ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "volumes.0.boot_server", ServerVCPUResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.lan", ServerVCPUResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.name", ServerVCPUResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.dhcp", ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_active", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_type", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.ips.0", ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.ips.1", ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.protocol", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.name", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.type", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerVCPUMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "name", ServerVCPUResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "cores", ServerVCPUResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "ram", ServerVCPUResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "availability_zone", ServerVCPUResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrSet(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "type", ServerVCPUResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.name", ServerVCPUResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.size", ServerVCPUResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.type", ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.bus", ServerVCPUResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.boot_server", ServerVCPUResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "volumes.0.availability_zone", ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.lan", ServerVCPUResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.name", ServerVCPUResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.dhcp", ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_active", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_type", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.ips.0", ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.ips.1", ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.name", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.type", ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceServerVCPUWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},
			{
				Config: testAccCheckServerVCPUConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.bus", "IDE"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "21"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:18"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
				),
			},
		},
	})
}

// issue #379
func TestAccServerVCPUNoBootVolumeBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUConfigNoBootVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "UNKNOWN"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.id", ServerVCPUResource+"."+ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigNoBootVolumeRemoveServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUAndVolumesDestroyed(DatacenterResource + "." + DatacenterTestResource),
				),
			},
		},
	})
}

// tests server with no cdromimage and with multiple firewall rules inline
func TestAccServerVCPUBootCdromNoImageAndInlineFwRules(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUConfigBootCdromNoImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "1"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfig2Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.name", ServerTestResource+"2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "2"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfig3Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.name", ServerTestResource+"2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.name", ServerTestResource+"3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "44"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigRemove2FwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource+"3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "44"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigRemoveAllFwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "0"),
				),
			},
		},
	})
}

// create and updates 5 inline rules after the server is created
func TestAccServerVCPUResolveImageNameAdd5FwRulesOnUpdate(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName5fwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "4"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "24"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "24"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_start", "26"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_end", "26"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName5fwRulesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "6"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_start", "24"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.port_range_end", "24"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_start", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.port_range_end", "25"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_start", "26"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.port_range_end", "26"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_start", "27"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.port_range_end", "27"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.1.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.2.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.3.type", "INGRESS"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.4.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
				),
			},
		},
	})
}

// also tests creating 5 fw rules inline
func TestAccServerVCPUWithSnapshotAnd5FwRulesInline(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUWithSnapshot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+"webserver", "nic.0.firewall.#", "5"),
				),
			},
		},
	})
}

func TestAccServerVCPUWithICMP(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUNoFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
				),
			},
			{
				Config: testAccCheckServerVCPUSeparateFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(ServerVCPUResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttrSet(ServerVCPUResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(FirewallResource+"."+FirewallTestResource, "type", "INGRESS"),
				),
			},
		},
	})
}

func TestAccServerVCPUWithLabels(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			// Clean server creation using labels in configuration.
			{
				Config: testAccCheckServerVCPUCreationWithLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.1.value", "labelvalue1"),
				),
			},
			// Check that labels are present in the server data source.
			{
				Config: testAccCheckDataSourceServerVCPUWithLabels,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "labels.#", "2"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(DataSource+"."+ServerVCPUResource+"."+ServerDataSourceById, "labels.1.value", "labelvalue1"),
				),
			},
			// Update server labels.
			{
				Config: testAccCheckServerVCPUUpdateLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.0.key", "updatedlabelkey0"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.0.value", "updatedlabelvalue0"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.1.key", "updatedlabelkey1"),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.1.value", "updatedlabelvalue1"),
				),
			},
			// Delete server labels.
			{
				Config: testAccCheckServerVCPUDeleteLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(ServerVCPUResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerVCPUResource+"."+ServerTestResource, "label.#", "0"),
				),
			},
		},
	})
}

func testAccCheckServerVCPUDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ServerVCPUResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("unable to fetch server %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("server still exists %s", rs.Primary.ID)

		}
	}

	return nil
}

func testAccCheckServerVCPUAndVolumesDestroyed(dcName string) resource.TestCheckFunc {
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

func testAccCheckServerVCPUExists(serverName string, server *ionoscloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[serverName]

		if !ok {
			return fmt.Errorf("testAccCheckServerVCPUExists: Not found: %s", serverName)
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

const bootCdromImageIdForVCPUServer = "0e4d57f9-cd78-11e9-b88c-525400f64d8d"

const testAccCheckServerVCPUConfigUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
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
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + UpdatedResources + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password_updated.result
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

const testAccDataSourceServerVCPUMatchId = testAccCheckServerVCPUConfigBasic + `
data ` + ServerVCPUResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerVCPUResource + `.` + ServerTestResource + `.id
}
`

const testAccDataSourceServerVCPUMatchName = testAccCheckServerVCPUConfigBasic + `
data ` + ServerVCPUResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + ServerTestResource + `"
}
`
const testAccDataSourceServerVCPUWrongNameError = testAccCheckServerVCPUConfigBasic + `
data ` + ServerVCPUResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckServerVCPUConfigBootCdromNoImage = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id        = ` + ServerVCPUResource + `.` + ServerTestResource + `.id
  nic_id           = ` + ServerVCPUResource + `.` + ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
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

const testAccCheckServerVCPUConfig2Fw = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + FirewallResource + ` ` + FirewallTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id        = ` + ServerVCPUResource + `.` + ServerTestResource + `.id
  nic_id           = ` + ServerVCPUResource + `.` + ServerTestResource + `.nic[0].id
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

resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
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

const testAccCheckServerVCPUConfig3Fw = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
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

const testAccCheckServerVCPUConfigRemove2FwRules = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
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

const testAccCheckServerVCPUConfigRemoveAllFwRules = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
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
const testAccCheckServerVCPUResolveImageNameNoNic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
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

const testAccCheckServerVCPUResolveImageName = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
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

const testAccCheckServerVCPUResolveImageName5fwRules = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
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

const testAccCheckServerVCPUResolveImageName5fwRulesUpdate = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
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

const testAccCheckServerVCPUWithSnapshot = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` "webserver" {
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
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
  server_id        = ` + ServerVCPUResource + `.webserver.id
  nic_id           = ` + ServerVCPUResource + `.webserver.nic[0].id
  protocol         = "TCP"
  name             = "external_rule"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  volume_id = ` + ServerVCPUResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  depends_on = [` + SnapshotResource + `.test_snapshot]
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
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

const testAccCheckServerVCPUNoFirewall = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
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
const testAccCheckServerVCPUSeparateFirewall = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
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
  server_id           = ` + ServerVCPUResource + `.` + ServerTestResource + `.id
  nic_id              = ` + ServerVCPUResource + `.` + ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUICMP = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
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
  server_id           = ` + ServerVCPUResource + `.` + ServerTestResource + `.id
  nic_id              = ` + ServerVCPUResource + `.` + ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckDataSourceServerVCPUWithLabels = testAccCheckServerVCPUCreationWithLabels + `
data ` + ServerVCPUResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerVCPUResource + `.` + ServerTestResource + `.id
}
`

const testAccCheckServerVCPUUpdateLabels = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
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

const testAccCheckServerVCPUDeleteLabels = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
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

const testAccCheckServerVCPUConfigNoBootVolume = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/txl"
}

resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + ServerVCPUResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  
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
  server_id               = ` + ServerVCPUResource + `.` + ServerTestResource + `.id
  name                    = "Another Volume Example"
  availability_zone       = "ZONE_1"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}
`
const testAccCheckServerVCPUConfigNoBootVolumeRemoveServer = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/txl"
}
`
