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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ssh_keys.0", sshKey)),
			},
			{
				Config: testAccCheckServerVCPUNoNic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
				),
			},
			{
				Config: testAccCheckServerVCPUShutDown,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "vm_state", "SHUTOFF"),
				),
			},
			{
				Config: testAccCheckServerVCPUNoNicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "vm_state", "SHUTOFF"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
				),
			},
			{
				Config: testAccCheckServerVCPUPowerOn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "vm_state", "RUNNING"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigEmptyNicIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerVCPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerVCPUResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerVCPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerVCPUResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccDataSourceServerVCPUMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "cores", constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "ram", constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "availability_zone", constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "cpu_family"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.size", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.bus", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.availability_zone", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "volumes.0.boot_server", constant.ServerVCPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.lan", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.dhcp", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_active", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.ips.0", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.ips.1", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.protocol", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerVCPUMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "cores", constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "ram", constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "availability_zone", constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "cpu_family"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.size", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.bus", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.boot_server", constant.ServerVCPUResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "volumes.0.availability_zone", constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.lan", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.dhcp", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_active", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.ips.0", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.ips.1", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.name", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.type", constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceServerVCPUWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},

			{
				Config: testAccCheckServerVCPUConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus", "IDE"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "21"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:18"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "UNKNOWN"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerVCPUResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigNoBootVolumeRemoveServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUAndVolumesDestroyed(constant.DatacenterResource + "." + constant.DatacenterTestResource),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "1"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfig2Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", constant.ServerTestResource+"2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "2"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfig3Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", constant.ServerTestResource+"2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", constant.ServerTestResource+"3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "44"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigRemove2FwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource+"3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "44"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerVCPUConfigRemoveAllFwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "0"),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName5fwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "4"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "24"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "24"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_start", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_start", "26"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_end", "26"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageName5fwRulesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "6"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "24"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "24"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_start", "26"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_end", "26"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_start", "27"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_end", "27"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.1.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.2.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.3.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.4.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+"webserver", "nic.0.firewall.#", "5"),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
				),
			},
			{
				Config: testAccCheckServerVCPUSeparateFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					utils.TestImageNotNull(constant.ServerVCPUResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerVCPUICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttrSet(constant.ServerVCPUResource+"."+constant.ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "type", constant.VCPUType),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
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
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.1.value", "labelvalue1"),
				),
			},
			// Check that labels are present in the server data source.
			{
				Config: testAccCheckDataSourceServerVCPUWithLabels,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "labels.#", "2"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerVCPUResource+"."+constant.ServerDataSourceById, "labels.1.value", "labelvalue1"),
				),
			},
			// Update server labels.
			{
				Config: testAccCheckServerVCPUUpdateLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.0.key", "updatedlabelkey0"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.0.value", "updatedlabelvalue0"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.1.key", "updatedlabelkey1"),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.1.value", "updatedlabelvalue1"),
				),
			},
			// Delete server labels.
			{
				Config: testAccCheckServerVCPUDeleteLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerVCPUExists(constant.ServerVCPUResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerVCPUResource+"."+constant.ServerTestResource, "label.#", "0"),
				),
			},
		},
	})
}

func testAccCheckServerVCPUDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerVCPUResource {
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
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_ipblock" "webserver_ipblock_update" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password_updated.result
  volume {
    name = "` + constant.UpdatedResources + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.UpdatedResources + `"
    dhcp = false
    firewall_active = false
    ips            = [ ionoscloud_ipblock.webserver_ipblock_update.ips[0], ionoscloud_ipblock.webserver_ipblock_update.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + constant.UpdatedResources + `"
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
data ` + constant.ServerVCPUResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceServerVCPUMatchName = testAccCheckServerVCPUConfigBasic + `
data ` + constant.ServerVCPUResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.ServerTestResource + `"
}
`
const testAccDataSourceServerVCPUWrongNameError = testAccCheckServerVCPUConfigBasic + `
data ` + constant.ServerVCPUResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckServerVCPUConfigBootCdromNoImage = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id        = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
  nic_id           = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      type = "EGRESS"
      name = "` + constant.ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerVCPUConfig2Fw = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id        = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
  nic_id           = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
}

resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `"
	  type = "EGRESS"
      port_range_start = 25
      port_range_end = 25
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[0]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[1]
    }
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `2"
      port_range_start = 23
      port_range_end = 23
    }
  }
}`

const testAccCheckServerVCPUConfig3Fw = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `"
      type = "EGRESS"
      port_range_start = 25
      port_range_end = 25
    }
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `2"
      port_range_start = 23
      port_range_end = 23
    }
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `3"
      port_range_start = 44
      port_range_end = 44
    }
  }
}`

const testAccCheckServerVCPUConfigRemove2FwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `3"
      port_range_start = 44
      port_range_end = 44
    }
  }
}`

const testAccCheckServerVCPUConfigRemoveAllFwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdForVCPUServer + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
	licence_type = "OTHER"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}`
const testAccCheckServerVCPUResolveImageNameNoNic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  image_name        = "ubuntu:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name           = "` + constant.ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUResolveImageName = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  image_name        = "ubuntu:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name           = "` + constant.ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp            = true
    firewall_active = true
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUResolveImageName5fwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  image_name        = "ubuntu:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name           = "` + constant.ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "ICMP"
      name             = "` + constant.ServerTestResource + `"
      type             = "INGRESS"
      icmp_code        = 4
      icmp_type        = 5
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 23
      port_range_end   = 23
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
	firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `5"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUResolveImageName5fwRulesUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/txl"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  image_name        = "ubuntu:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name           = "` + constant.ServerTestResource + `"
    size              = 5
    disk_type      = "SSD Standard"
  }
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "ICMP"
      name             = "` + constant.ServerTestResource + `"
      type             = "INGRESS"
      icmp_code        = 5
      icmp_type        = 6
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
	firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `5"
      type             = "EGRESS"
      port_range_start = 27
      port_range_end   = 27
    }
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUWithSnapshot = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` "webserver" {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
	image_name = "ubuntu:latest"
	image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `"
      port_range_start = 22
      type             = "EGRESS"
      port_range_end   = 22
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `2"
      type             = "INGRESS"
      port_range_start = 23
      port_range_end   = 23
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 24
      port_range_end   = 24
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 25
      port_range_end   = 25
    }
	firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `5"
      type             = "INGRESS"
      port_range_start = 26
      port_range_end   = 26
    }
  }
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id        = ` + constant.ServerVCPUResource + `.webserver.id
  nic_id           = ` + constant.ServerVCPUResource + `.webserver.nic[0].id
  protocol         = "TCP"
  name             = "external_rule"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + constant.SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  volume_id = ` + constant.ServerVCPUResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  depends_on = [` + constant.SnapshotResource + `.test_snapshot]
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name = "terraform_snapshot"
  volume {
    name = "` + constant.ServerTestResource + `"
    size = 5
    disk_type = "SSD Standard"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUNoFirewall = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`
const testAccCheckServerVCPUSeparateFirewall = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name 			= "system"
    dhcp            = true
    firewall_active = true
    }
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id           = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
  nic_id              = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerVCPUICMP = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
	disk_type = "HDD"
}
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name 			= "system"
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "ICMP"
      name             = "` + constant.ServerTestResource + `"
      icmp_type        = "12"
      icmp_code        = "0"
	  }
    }
}
resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id           = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
  nic_id              = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckDataSourceServerVCPUWithLabels = testAccCheckServerVCPUCreationWithLabels + `
data ` + constant.ServerVCPUResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckServerVCPUUpdateLabels = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
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
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
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
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
}`

const testAccCheckServerVCPUConfigNoBootVolume = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/txl"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
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
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "INGRESS"
  }
}

resource "ionoscloud_volume" "exampleVol1" {
  datacenter_id           = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id               = ` + constant.ServerVCPUResource + `.` + constant.ServerTestResource + `.id
  name                    = "Another Volume Example"
  availability_zone       = "ZONE_1"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}
`
const testAccCheckServerVCPUConfigNoBootVolumeRemoveServer = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/txl"
}
`
