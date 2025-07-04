//go:build compute || all || server || enterprise
// +build compute all server enterprise

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

const bootCdromImageId = "aa97f2f4-ca28-11ec-925c-46570854dba1"

// ToDo: add backup unit back in tests when stable

func TestAccServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckServerNoPwdOrSSH,
				ExpectError: regexp.MustCompile(`either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided`),
			},
			{
				// ssh_key_path now accepts the ssh key directly too
				Config: testAccCheckServerSshDirectly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ssh_key_path.0", sshKey)),
			},
			{
				Config: testAccCheckServerSshKeysDirectly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ssh_keys.0", sshKey)),
			},
			{
				Config:      testAccCheckServerSshKeysAndKeyPathErr,
				ExpectError: regexp.MustCompile(`"ssh_keys": conflicts with ssh_key_path`),
			},
			{
				Config: testAccCheckServerNoNic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "hostname", constant.ServerTestHostname),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
				),
			},
			{
				Config: testAccCheckServerNoNicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "hostname", constant.ServerTestHostname),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "2048"),
				),
			},
			{
				Config: testAccCheckServerConfigEmptyNicIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttrSet(constant.ServerResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
				),
			},
			{
				Config: testAccCheckServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccDataSourceServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "name", constant.ServerResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "hostname", constant.ServerResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "cores", constant.ServerResource+"."+constant.ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "ram", constant.ServerResource+"."+constant.ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "type", constant.ServerResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "volumes.0.name", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "volumes.0.size", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "volumes.0.type", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "volumes.0.bus", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "volumes.0.boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.lan", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.name", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.dhcp", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_active", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_type", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ips.0", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ips.1", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.protocol", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.name", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.type", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "name", constant.ServerResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "hostname", constant.ServerResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "cores", constant.ServerResource+"."+constant.ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "ram", constant.ServerResource+"."+constant.ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "type", constant.ServerResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "volumes.0.name", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "volumes.0.size", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "volumes.0.type", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "volumes.0.bus", constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "volumes.0.boot_server", constant.ServerResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.lan", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.name", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.dhcp", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_active", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_type", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.ips.0", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.ips.1", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.name", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.type", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceServerWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},
			{
				Config: testAccCheckServerConfigIpv6Enabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcpv6", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "hostname", "updated"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.dhcpv6", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_cidr_block", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.0", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.1", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.2", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.2"),
				),
			},
			{
				Config: testAccCheckServerConfigIpv6Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcpv6", "false"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.dhcpv6", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_cidr_block", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.0", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.1", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.2", constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.2"),
				),
			},
			{
				Config: testAccCheckServerConfigShutdown,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "vm_state", constant.VMStateStop),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "vm_state", constant.ServerResource+"."+constant.ServerTestResource, "vm_state"),
				),
			},
			{
				Config: testAccCheckServerConfigPowerOn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "vm_state", constant.VMStateStart),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "vm_state", constant.ServerResource+"."+constant.ServerTestResource, "vm_state"),
				),
			},
			{
				Config: testAccCheckServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "security_groups_ids.#", "1"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus", "IDE"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.security_groups_ids.#", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock_update", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock_update", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "21"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:18"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock_update", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
				),
			},
		},
	})
}

// issue #379 - also tests if a server can be created without a volume
func TestAccServerNoBootVolumeBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerNoInlineVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckNoResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.#"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.mac", constant.NicMac),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerConfigInlineVolumeUnknownLicense,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "6"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "UNKNOWN"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.expose_serial", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.mac", constant.NicMac),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerNoInlineVolumeConfig, // checks if removing the volume works
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "ENTERPRISE"),
					resource.TestCheckNoResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.#"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.mac", constant.NicMac),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.id", constant.ServerResource+"."+constant.ServerTestResource, "primary_nic"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerConfigRemoveServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerAndVolumesDestroyed(constant.DatacenterResource + "." + constant.DatacenterTestResource),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBootCdromNoImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.mac", constant.NicMac),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "1"),
				),
			},
			{
				Config: testAccCheckServerConfig2Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "25"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", constant.ServerTestResource+"2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "2"),
				),
			},
			{
				Config: testAccCheckServerConfig3Fw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", constant.ServerTestResource+"2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", constant.ServerTestResource+"3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "44"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerConfigRemove2FwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource+"3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "44"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "44"),
				),
			},
			{
				Config: testAccCheckServerConfigRemoveAllFwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "OTHER"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "0"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName5fwRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "4"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "23"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "24"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "24"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_start", "25"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_end", "25"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_start", "26"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_end", "26"),
				),
			},
			{
				Config: testAccCheckServerResolveImageName5fwRulesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.#", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "test_server"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.name", "test_server2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.name", "test_server3"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.name", "test_server4"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.name", "test_server5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "6"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_start", "204"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.port_range_end", "204"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_start", "205"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.port_range_end", "205"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_start", "206"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.port_range_end", "206"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_start", "207"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.port_range_end", "207"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.1.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.2.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.3.type", "INGRESS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.4.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckServerResolveImageNameNoNic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerWithSnapshot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "SSD Standard"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+"webserver", "nic.0.firewall.#", "5"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerAndServersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "cores", "data.ionoscloud_template."+constant.ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "ram", "data.ionoscloud_template."+constant.ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "template_uuid", "data.ionoscloud_template."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "type", "CUBE"),
					utils.TestImageNotNull("ionoscloud_server", "boot_image"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "data.ionoscloud_template."+constant.ServerTestResource, "storage_size"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.licence_type", "LINUX"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", "ionoscloud_lan.webserver_lan", "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.#", "1"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerNoFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
				),
			},
			{
				Config: testAccCheckSeparateFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "ram", "1024"),
					utils.TestImageNotNull(constant.ServerResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.size", "5"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "volume.0.disk_type", "HDD"),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckServerICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "name", "allow-icmp"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.FirewallResource+"."+constant.FirewallTestResource, "type", "INGRESS"),
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
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			// Clean server creation using labels in configuration.
			{
				Config: testAccCheckServerCreationWithLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.1.value", "labelvalue1"),
				),
			},
			// Check that labels are present in the server data source.
			{
				Config: testAccCheckDataSourceServerWithLabels,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "labels.#", "2"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServerResource+"."+constant.ServerDataSourceById, "labels.1.value", "labelvalue1"),
				),
			},
			// Update server labels.
			{
				Config: testAccCheckServerUpdateLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.#", "2"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.0.key", "updatedlabelkey0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.0.value", "updatedlabelvalue0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.1.key", "updatedlabelkey1"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.1.value", "updatedlabelvalue1"),
				),
			},
			// Delete server labels.
			{
				Config: testAccCheckServerDeleteLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "label.#", "0"),
				),
			},
		},
	})
}

func TestAccServerBootDeviceSelection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccMultipleBootDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "boot_volume", constant.ServerResource+"."+constant.ServerTestResource, "inline_volume_ids.0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", ""),
				)},
			{
				Config: testExternalVolumeSelection,
				Check: resource.ComposeTestCheckFunc(
					// The Server resource state is not updated after the Apply because the changes are triggered through the boot selection resource, so a second refresh is not performed for the Server state
					// The new boot_volume value will be available in the next Update.
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", ""),
				)},
			{
				Config: testImageCdromSelection,
				Check: resource.ComposeTestCheckFunc(
					// The Server resource state is not updated after the Apply because the changes are triggered through the selection resource, so a second refresh is not performed for the Server state
					// The new boot_cdrom value will be available in the next Update
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "boot_volume", constant.VolumeResource+"."+constant.VolumeTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", ""),
				)},
			{
				Config: testAccPxeShellBoot,
				Check: resource.ComposeTestCheckFunc(
					// The Server resource state is not updated after the Apply because the changes are triggered through the selection resource, so a second refresh is not performed for the Server state
					// The null value for boot_cdrom will be available in the next Update
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_volume", ""),
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", `data.`+constant.ImageResource+"."+constant.ImageTestResource, "id"),
				)},
			{
				Config: testAccCheckServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					// The Server resource state is not updated after the Apply because the changes are triggered through the selection resource, so a second refresh is not performed for the Server state
					// The new boot_volume value will be available in the next Update.
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_volume", ""),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", ""),
				)},
			{
				Config: testAccCheckServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.ServerResource+"."+constant.ServerTestResource, "boot_volume", constant.ServerResource+"."+constant.ServerTestResource, "inline_volume_ids.0"),
					resource.TestCheckResourceAttr(constant.ServerResource+"."+constant.ServerTestResource, "boot_cdrom", ""),
				)},
		},
	})
}

func testAccCheckServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerResource {
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

func testAccCheckServerAndVolumesDestroyed(dcName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

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
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching Server: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		server = &foundServer
		return nil
	}
}

const testAccCheckServerConfigUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
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
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password_updated.result
  type = "ENTERPRISE"
  security_groups_ids = [ionoscloud_nsg.example_1.id]
  volume {
    name = "` + constant.UpdatedResources + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.UpdatedResources + `"
    dhcp = false
    firewall_active = false
    security_groups_ids = [ionoscloud_nsg.example_2.id]
    ips = [ ionoscloud_ipblock.webserver_ipblock_update.ips[0], ionoscloud_ipblock.webserver_ipblock_update.ips[1] ]
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
` + ServerImagePasswordUpdated + testSecurityGroups

const testAccDataSourceServerMatchId = testAccCheckServerConfigBasic + `
data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceServerMatchName = testAccCheckServerConfigBasic + `
data ` + constant.ServerResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.ServerTestResource + `"
}
`
const testAccDataSourceServerWrongNameError = testAccCheckServerConfigBasic + `
data ` + constant.ServerResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckServerConfigBootCdromNoImage = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id        = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id           = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol         = "TCP"
  name             = "SSH"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  boot_cdrom = "` + bootCdromImageId + `" 
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
    mac = "` + constant.NicMac + `"
	firewall {
      protocol = "TCP"
      type = "EGRESS"
      name = "` + constant.ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckServerConfig2Fw = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.FirewallResource + ` ` + constant.FirewallTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id        = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id           = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.nic[0].id
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

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  boot_cdrom = "` + bootCdromImageId + `" 
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

const testAccCheckServerConfig3Fw = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  boot_cdrom = "` + bootCdromImageId + `" 
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

const testAccCheckServerConfigRemove2FwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  boot_cdrom = "` + bootCdromImageId + `" 
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

const testAccCheckServerConfigRemoveAllFwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  boot_cdrom = "` + bootCdromImageId + `" 
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
const testAccCheckServerResolveImageNameNoNic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "us/las"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
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

const testAccCheckServerResolveImageName = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "us/las"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
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

const testAccCheckServerResolveImageName5fwRules = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "us/las"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
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

const testAccCheckServerResolveImageName5fwRulesUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "us/las"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores             = 1
  ram               = 1024
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
      port_range_start = 204
      port_range_end   = 204
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `3"
      type             = "INGRESS"
      port_range_start = 205
      port_range_end   = 205
    }
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `4"
      type             = "INGRESS"
      port_range_start = 206
      port_range_end   = 206
    }
	firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `5"
      type             = "EGRESS"
      port_range_start = 207
      port_range_end   = 207
    }
  }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerWithSnapshot = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` "webserver" {
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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
  server_id        = ` + constant.ServerResource + `.webserver.id
  nic_id           = ` + constant.ServerResource + `.webserver.nic[0].id
  protocol         = "TCP"
  name             = "external_rule"
  port_range_start = 28
  port_range_end   = 28
}

resource ` + constant.SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  volume_id = ` + constant.ServerResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  depends_on = [` + constant.SnapshotResource + `.test_snapshot]
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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

const testAccCheckCubeServerAndServersDataSource = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + " " + constant.DatacenterTestResource + `{
	name       = "volume-test"
	location   = "de/txl"
}

resource "ionoscloud_lan" "webserver_lan" {
 datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
 public = true
 name = "public"
}

resource "ionoscloud_server" ` + constant.ServerTestResource + ` {
 name              = "` + constant.ServerTestResource + `"
 image_name        = "ubuntu:latest"
 type              = "CUBE"
 template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
 image_password = ` + constant.RandomPassword + `.server_image_password.result
 datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
 volume {
   name            = "` + constant.ServerTestResource + `"
   licence_type    = "LINUX"
   disk_type = "DAS"
	}
 nic {
   lan             = ionoscloud_lan.webserver_lan.id
   name            = "` + constant.ServerTestResource + `"
   dhcp            = true
   firewall_active = true
 }
}
data ` + constant.ServersDataSource + ` ` + constant.ServerDataSourceByName + ` {
depends_on = [` + constant.ServerResource + `.` + constant.ServerTestResource + `]
 datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
 filter {
  name = "type"
  value = "CUBE"
 }
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerNoFirewall = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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
const testAccCheckSeparateFirewall = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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
  server_id           = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id              = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckServerICMP = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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
  server_id           = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  nic_id              = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.nic[0].id
  protocol            = "ICMP"
  name                = "allow-icmp"
  type                = "INGRESS"
}
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const testAccCheckDataSourceServerWithLabels = testAccCheckServerCreationWithLabels + `
data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckServerUpdateLabels = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
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

const testAccCheckServerDeleteLabels = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
}`
const testAccCheckServerNoInlineVolumeConfig = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  type = "ENTERPRISE"
  allow_replace = true
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "INGRESS"
    mac = "` + constant.NicMac + `"
  }
}

resource "ionoscloud_volume" "exampleVol1" {
  datacenter_id           = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id               = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name                    = "Another Volume Example"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}
`
const testAccCheckServerConfigInlineVolumeUnknownLicense = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  type = "ENTERPRISE"
  allow_replace = true
  volume {
    name = "system"
    size = 6
    licence_type = "UNKNOWN"
    disk_type = "SSD Standard"
    bus = "VIRTIO"
	expose_serial = true
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "INGRESS"
    mac = "` + constant.NicMac + `"
  }
}

resource "ionoscloud_volume" "exampleVol1" {
  datacenter_id           = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id               = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  name                    = "Another Volume Example"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}
`
const testAccCheckServerConfigRemoveServer = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
`

const testAccMultipleBootDevices = testAccCheckServerConfigBasic + `
resource ` + constant.VolumeResource + ` ` + constant.VolumeTestResource + ` {
  server_id         = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name              = "External Volume 1"
  size              = 10
  disk_type         = "SSD"
  image_name        = "debian:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
}
data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
  name     = "ubuntu-22.04"
  location = "us/las"
  type     = "CDROM"
}
`

const testExternalVolumeSelection = testAccMultipleBootDevices + `
resource ` + constant.ServerBootDeviceSelectionResource + ` ` + constant.TestServerBootDeviceSelectionResource + ` {
  datacenter_id  = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id      = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  boot_device_id = ` + constant.VolumeResource + `.` + constant.VolumeTestResource + `.id
}
`

const testImageCdromSelection = testAccMultipleBootDevices + `
resource ` + constant.ServerBootDeviceSelectionResource + ` ` + constant.TestServerBootDeviceSelectionResource + ` {
  datacenter_id  = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id      = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
  boot_device_id = data.` + constant.ImageResource + `.` + constant.ImageTestResource + `.id
}
`
const testAccPxeShellBoot = testAccMultipleBootDevices + `
resource ` + constant.ServerBootDeviceSelectionResource + ` ` + constant.TestServerBootDeviceSelectionResource + ` {
  datacenter_id  = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  server_id      = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`
