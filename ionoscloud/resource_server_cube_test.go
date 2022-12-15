//go:build compute || all || server

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const bootCdromImageIdCube = "83f21679-3321-11eb-a681-1e659523cb7b"

func TestAccCubeServerBasic(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					utils.TestImageNotNull(ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "volume.0.boot_server", ServerCubeResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_type", "BIDIRECTIONAL"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name", "SSH"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
			{
				Config: testAccDataSourceCubeServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "name", ServerCubeResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "availability_zone", ServerCubeResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "type", ServerCubeResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "volumes.0.name", ServerCubeResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "volumes.0.type", ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "volumes.0.bus", ServerCubeResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "volumes.0.availability_zone", ServerCubeResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "volumes.0.boot_server", ServerCubeResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.lan", ServerCubeResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.name", ServerCubeResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.dhcp", ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_active", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_type", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.ips.0", ServerCubeResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.ips.1", ServerCubeResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.protocol", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.name", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.type", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceCubeServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "name", ServerCubeResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "availability_zone", ServerCubeResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "type", ServerCubeResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "volumes.0.name", ServerCubeResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "volumes.0.type", ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "volumes.0.bus", ServerCubeResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "volumes.0.boot_server", ServerCubeResource+"."+ServerTestResource, "id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "volumes.0.availability_zone", ServerCubeResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.lan", ServerCubeResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.name", ServerCubeResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.dhcp", ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_active", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_type", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.ips.0", ServerCubeResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.ips.1", ServerCubeResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.name", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerCubeResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.type", ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config:      testAccDataSourceCubeServerWrongNameError,
				ExpectError: regexp.MustCompile(`no server found with the specified criteria: name`),
			},
			{
				Config: testAccCheckCubeServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					utils.TestImageNotNull(ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password_updated", "result"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.ips.0", "ionoscloud_ipblock.webserver_ipblock", "ips.0"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.ips.1", "ionoscloud_ipblock.webserver_ipblock", "ips.1"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac", "00:0a:95:9d:68:17"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.2"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip", "ionoscloud_ipblock.webserver_ipblock", "ips.3"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.type", "EGRESS"),
				),
			},
		},
	})
}

//func TestAccCubeServerBootCdromNoImage(t *testing.T) { // todo is returning 500 interal, for the moment
//	var server ionoscloud.Server
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckCubeServerConfigBootCdromNoImage,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
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
//}

func TestAccCubeServerResolveImageName(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerResolveImageName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(ServerCubeResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					utils.TestImageNotNull(ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start", "22"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end", "22"),
				),
			},
		},
	})
}

//func TestAccCubeServerWithSnapshot(t *testing.T) { // todo for now is a vdc problem and the snapshot with a das volume when is deleting but the state remains procesing
//	var server ionoscloud.Server
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ExternalProviders: randomProviderVersion343(),
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(testAccCheckCubeServerWithSnapshot),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
//					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
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
//}

func TestAccCubeServerWithICMP(t *testing.T) {
	var server ionoscloud.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerNoFirewall,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "availability_zone", "ZONE_1"),
					utils.TestImageNotNull(ServerCubeResource, "boot_image"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "image_password", RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.name", "system"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "volume.0.disk_type", "DAS"),
					resource.TestCheckResourceAttrPair(ServerCubeResource+"."+ServerTestResource, "nic.0.lan", LanResource+"."+LanTestResource, "id"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.name", "system"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "10"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "1"),
				),
			},
			{
				Config: testAccCheckCubeServerICMP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCubeServerExists(ServerCubeResource+"."+ServerTestResource, &server),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.name", ServerTestResource),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_type", "12"),
					resource.TestCheckResourceAttr(ServerCubeResource+"."+ServerTestResource, "nic.0.firewall.0.icmp_code", "0"),
				),
			},
		},
	})
}

func testAccCheckCubeServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ServerCubeResource {
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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occured while fetching Server: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		server = &foundServer

		return nil
	}
}

const testAccCheckCubeServerConfigUpdate = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  name = "` + UpdatedResources + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password_updated.result
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id

  volume {
    name            = "` + ServerTestResource + `"
    licence_type    = "LINUX"
    disk_type = "DAS"
	}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "` + UpdatedResources + `"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
     firewall {
      protocol = "TCP"
      name = "` + UpdatedResources + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}
resource ` + RandomPassword + ` "server_image_password_updated" {
  length           = 16
  special          = false
}
`

const testAccDataSourceCubeServerMatchId = testAccCheckCubeServerConfigBasic + `
data ` + ServerCubeResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerCubeResource + `.` + ServerTestResource + `.id
}
`

const testAccDataSourceCubeServerMatchName = testAccCheckCubeServerConfigBasic + `
data ` + ServerCubeResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + ServerTestResource + `"
}
`
const testAccDataSourceCubeServerWrongNameError = testAccCheckCubeServerConfigBasic + `
data ` + ServerCubeResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "wrong_name"
}
`

const testAccCheckCubeServerConfigBootCdromNoImage = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  boot_cdrom = "` + bootCdromImageIdCube + `" 
  volume {
    name = "` + ServerTestResource + `"
    disk_type = "DAS"
	licence_type = "LINUX"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true 
	firewall {
      protocol = "TCP"
      name = "` + ServerTestResource + `"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckCubeServerResolveImageName = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
  name        = "test_server"
  location    = "de/fra"
  description = "Test datacenter done by TF"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public        = true
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name              = "` + ServerTestResource + `"
  datacenter_id     = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name        = "ubuntu:latest"
  image_password    = ` + RandomPassword + `.server_image_password.result
  volume {
    name           = "` + ServerTestResource + `"
    disk_type      = "DAS"
  }
  nic {
    lan             = ` + LanResource + `.` + LanTestResource + `.id
    dhcp            = true
    firewall_active = true
    firewall {
      protocol         = "TCP"
      name             = "` + ServerTestResource + `"
      port_range_start = 22
      port_range_end   = 22
    }
  }
}
` + ServerImagePassword

const testAccCheckCubeServerWithSnapshot = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "volume-test"
	location   = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerCubeResource + ` "webserver" {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name = "webserver"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name = "ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    disk_type = "DAS"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
resource ` + SnapshotResource + ` "test_snapshot" {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  volume_id = ` + ServerCubeResource + `.webserver.boot_volume
  name = "terraform_snapshot"
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  depends_on = [` + SnapshotResource + `.test_snapshot]
  name = "` + ServerTestResource + `"
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name = "terraform_snapshot"
  volume {
    name = "` + ServerTestResource + `"
    disk_type = "DAS"
  }
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    dhcp = true
    firewall_active = true
  }
}
` + ServerImagePassword

const testAccCheckCubeServerNoFirewall = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
	disk_type = "DAS"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
    firewall {
      protocol         = "ICMP"
      name             = "` + ServerTestResource + `"
      icmp_type        = "10"
      icmp_code        = "1"
	  }
  }
}
` + ServerImagePassword

const testAccCheckCubeServerICMP = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}

resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
	licence_type    = "LINUX"
    disk_type = "DAS"
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
` + ServerImagePassword
