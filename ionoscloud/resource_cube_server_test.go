//go:build compute || all || server || cube
// +build compute all server cube

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const bootCdromImageIdCube = "83f21679-3321-11eb-a681-1e659523cb7b"

func TestAccCubeServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(constant.ServerCubeResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "hostname", constant.ServerTestHostname),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.boot_server", constant.ServerCubeResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccDataSourceCubeServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "name", constant.ServerCubeResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "hostname", constant.ServerCubeResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "availability_zone", constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "type", constant.ServerCubeResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "volumes.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "volumes.0.type", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "volumes.0.bus", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "volumes.0.availability_zone", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "volumes.0.boot_server", constant.ServerCubeResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.lan", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.dhcp", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_active", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_type", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ips.0", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ips.1", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.protocol", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.firewall_rules.0.type", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceCubeServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "name", constant.ServerCubeResource+"."+constant.ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "hostname", constant.ServerCubeResource+"."+constant.ServerTestResource, "hostname"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "availability_zone", constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "type", constant.ServerCubeResource+"."+constant.ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "volumes.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "volumes.0.type", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "volumes.0.bus", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "volumes.0.boot_server", constant.ServerCubeResource+"."+constant.ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "volumes.0.availability_zone", constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.lan", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.dhcp", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_active", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_type", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.ips.0", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.ips.1", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.name", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceByName, "nics.0.firewall_rules.0.type", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceCubeServerWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},
			{
				Config: testAccCheckCubeServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(constant.ServerCubeResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "hostname", "updatedhostname"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckCubeServerEnableIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcpv6", "true"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.dhcpv6", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_cidr_block", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.0", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.1", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.2", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.2"),
				),
			},
			{
				Config: testAccCheckCubeServerUpdateIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcpv6", "false"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.dhcpv6", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcpv6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_cidr_block", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_cidr_block"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.0", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.1", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "nics.0.ipv6_ips.2", constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.ipv6_ips.2"),
				),
			},
			{
				Config: testAccCheckCubeServerSuspend,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "vm_state", constant.CubeVMStateStop),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "vm_state", constant.ServerCubeResource+"."+constant.ServerTestResource, "vm_state"),
				),
			},
			{
				Config:      testAccCheckCubeServerUpdateWhenSuspended,
				ExpectError: regexp.MustCompile(`cannot update a suspended Cube Server, must change the state to RUNNING first`),
			},
			{
				Config: testAccCheckCubeServerResume,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "vm_state", constant.VMStateStart),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServerCubeResource+"."+constant.ServerDataSourceById, "vm_state", constant.ServerCubeResource+"."+constant.ServerTestResource, "vm_state"),
				),
			},
		},
	})
}

// func TestAccCubeServerBootCdromNoImage(t *testing.T) { // todo is returning 500 interal, for the moment
//	var server ionoscloud.Server
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
//		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckCubeServerConfigBootCdromNoImage,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "AUTO"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.licence_type", "OTHER"),
//					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
//				),
//			},
//		},
//	})
// }

func TestAccCubeServerResolveImageName(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(constant.ServerCubeResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
				),
			},
		},
	})
}

// func TestAccCubeServerWithSnapshot(t *testing.T) { // todo for now is a vdc problem and the snapshot with a das volume when is deleting but the state remains procesing
//	var server ionoscloud.Server
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ExternalProviders: randomProviderVersion343(),
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
//		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(testAccCheckCubeServerWithSnapshot),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "AUTO"),
//					utils.TestImageNotNull(ServerCubeResource, "boot_image"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
//					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
//				),
//			},
//		},
//	})
// }

func TestAccCubeServerWithICMP(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerNoFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(constant.ServerCubeResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "availability_zone", "AUTO"),
					utils.TestImageNotNull(constant.ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttrPair(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.lan", constant.LanResource+"."+constant.LanTestResource, "id"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "10"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "1"),
				),
			},
			{
				Config: testAccCheckCubeServerICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(constant.ServerCubeResource+"."+constant.ServerTestResource, &server),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.name", constant.ServerTestResource),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(constant.ServerCubeResource+"."+constant.ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
				),
			},
		},
	})
}

func testAccCheckCubeServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerCubeResource {
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

func testAccCheckCubeServerExists(n string, server *ionoscloud.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckServerExists: Not found: %s", n)
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

const testAccCheckCubeServerConfigUpdate = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  hostname = "updatedhostname"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password_updated.result
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id

  volume {
    name            = "` + constant.ServerTestResource + `"
    licence_type    = "LINUX"
    disk_type = "DAS"
	}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.UpdatedResources + `"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
     firewall {
      protocol = "TCP"
      name = "` + constant.UpdatedResources + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}
` + ServerImagePasswordUpdated

const testAccDataSourceCubeServerMatchId = testAccCheckCubeServerConfigBasic + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id			= ` + constant.ServerCubeResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceCubeServerMatchName = testAccCheckCubeServerConfigBasic + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "` + constant.ServerTestResource + `"
}
`
const testAccDataSourceCubeServerWrongNameError = testAccCheckCubeServerConfigBasic + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckCubeServerConfigBootCdromNoImage = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  boot_cdrom = "` + bootCdromImageIdCube + `" 
  volume {
    name = "` + constant.ServerTestResource + `"
    disk_type = "DAS"
	licence_type = "LINUX"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + constant.ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckCubeServerResolveImageName = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/fra"
  description = "Test datacenter done by TF"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public        = true
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name              = "` + constant.ServerTestResource + `"
  datacenter_id     = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name        = "ubuntu:latest"
  image_password    = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name           = "` + constant.ServerTestResource + `"
    disk_type      = "DAS"
  }
  nic {
    lan             = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + constant.ServerTestResource + `"
      port_range_start = 22
      port_range_end   = 22
    }
  }
}
` + ServerImagePassword

const testAccCheckCubeServerWithSnapshot = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` "webserver" {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "webserver"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    disk_type = "DAS"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + constant.SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  volume_id = ` + constant.ServerCubeResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  depends_on = [` + constant.SnapshotResource + `.test_snapshot]
  name = "` + constant.ServerTestResource + `"
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name = "terraform_snapshot"
  volume {
    name = "` + constant.ServerTestResource + `"
    disk_type = "DAS"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
` + ServerImagePassword

const testAccCheckCubeServerNoFirewall = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
	disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
    firewall {
      protocol         = "ICMP"
      name             = "` + constant.ServerTestResource + `"
      icmp_type        = "10"
      icmp_code        = "1"
	  }
  }
}
` + ServerImagePassword

const testAccCheckCubeServerICMP = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}

resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  availability_zone = "AUTO"
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
	licence_type    = "LINUX"
    disk_type = "DAS"
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
` + ServerImagePassword
