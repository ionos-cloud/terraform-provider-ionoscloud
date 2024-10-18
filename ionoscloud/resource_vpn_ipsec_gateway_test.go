//go:build all || vpn || vpn_ipsec
// +build all vpn vpn_ipsec

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestVpnIPSecGatewayResource(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ProviderFactories: testAccProviderFactories,
			ExternalProviders: randomProviderVersion343(),
			CheckDestroy:      testCheckIPSecGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: configIPSecGatewayBasic(gatewayResourceName, gatewayAttributeNameValue),
					Check:  checkIPSecGatewayResource(gatewayAttributeNameValue),
				},
				{
					Config: configIPSecGatewayBasic(gatewayResourceName, fmt.Sprintf("%v_updated", gatewayAttributeNameValue)),
					Check:  checkIPSecGatewayResource(fmt.Sprintf("%v_updated", gatewayAttributeNameValue)),
				},
				{
					Config: configIPSecGatewayDataSourceGetByID(gatewayResourceName, gatewayDataName, constant.IPSecGatewayResource+"."+gatewayResourceName+".id"),
					Check: checkIPSecGatewayResourceAttributesComparative(
						constant.DataSource+"."+constant.IPSecGatewayResource+"."+gatewayDataName, constant.IPSecGatewayResource+"."+gatewayResourceName,
					),
				},
				{
					Config: configIPSecGatewayDataSourceGetByName(
						gatewayResourceName, gatewayDataName, constant.IPSecGatewayResource+"."+gatewayResourceName+"."+gatewayAttributeName,
					),
					Check: checkIPSecGatewayResourceAttributesComparative(
						constant.DataSource+"."+constant.IPSecGatewayResource+"."+gatewayDataName, constant.IPSecGatewayResource+"."+gatewayResourceName,
					),
				},
				{
					Config:      configIPSecGatewayDataSourceGetByName(gatewayResourceName, gatewayDataName, `"willnotwork"`),
					ExpectError: regexp.MustCompile(`no VPN IPSec Gateway found with the specified name =`),
				},
			},
		},
	)
}

func testCheckIPSecGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).VPNClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.IPSecGatewayResource {
			continue
		}

		_, apiResponse, err := client.GetIPSecGatewayByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of resource %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("resource %s in %s still exists", rs.Primary.ID, rs.Primary.Attributes["location"])
		}
	}
	return nil
}

func testAccCheckIPSecGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).VPNClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundGateway, _, err := client.GetIPSecGatewayByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching IPSec Gateway with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if foundGateway.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

func checkIPSecGatewayResource(attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckIPSecGatewayExists(constant.IPSecGatewayResource+"."+gatewayResourceName),
		checkIPSecGatewayResourceAttributes(constant.IPSecGatewayResource+"."+gatewayResourceName, attributeNameReferenceValue),
	)
}

func checkIPSecGatewayResourceAttributes(fullResourceName, attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fullResourceName, gatewayAttributeName, attributeNameReferenceValue),
		resource.TestCheckResourceAttr(fullResourceName, gatewayAttributeVersion, gatewayAttributeVersionValue),
		resource.TestCheckResourceAttrPair(fullResourceName, gatewayAttributeIP, "ionoscloud_ipblock.test_ipblock", "ips.0"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.datacenter_id", "ionoscloud_datacenter.test_datacenter", "id"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.lan_id", "ionoscloud_lan.test_lan", "id"),
		resource.TestCheckResourceAttrSet(fullResourceName, "connections.0.ipv4_cidr"),
		resource.TestCheckResourceAttr(fullResourceName, "maintenance_window.0.day_of_the_week", "Monday"),
		resource.TestCheckResourceAttr(fullResourceName, "maintenance_window.0.time", "09:00:00"),
		resource.TestCheckResourceAttr(fullResourceName, "tier", "STANDARD"),
		//resource.TestCheckResourceAttrSet(fullResourceName, "connections.0.ipv6_cidr"),
	)
}

func checkIPSecGatewayResourceAttributesComparative(fullResourceName, fullReferenceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(fullResourceName, gatewayAttributeName, fullReferenceResourceName, gatewayAttributeName),
		resource.TestCheckResourceAttrPair(fullResourceName, gatewayAttributeVersion, fullReferenceResourceName, gatewayAttributeVersion),
		resource.TestCheckResourceAttrPair(fullResourceName, gatewayAttributeIP, fullReferenceResourceName, gatewayAttributeIP),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.datacenter_id", fullReferenceResourceName, "connections.0.datacenter_id"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.lan_id", fullReferenceResourceName, "connections.0.lan_id"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.ipv4_cidr", fullReferenceResourceName, "connections.0.ipv4_cidr"),
		resource.TestCheckResourceAttrPair(fullResourceName, "maintenance_window.0.day_of_the_week", fullReferenceResourceName, "maintenance_window.0.day_of_the_week"),
		resource.TestCheckResourceAttrPair(fullResourceName, "maintenance_window.0.time", fullReferenceResourceName, "maintenance_window.0.time"),
		resource.TestCheckResourceAttrPair(fullResourceName, "tier", fullReferenceResourceName, "tier"),
		//resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.ipv6_cidr", fullReferenceResourceName, "connections.0.ipv6_cidr"),
	)
}

func configIPSecGatewayBasic(resourceName, attributeName string) string {
	gatewayBasicConfig := fmt.Sprintf(templateIPSecGatewayConfig, gatewayResourceName, attributeName, gatewayAttributeVersionValue, gatewayAttributeLocationValue)

	return strings.Join([]string{defaultIPSecGatewayBaseConfig, gatewayBasicConfig}, "\n")
}

func configIPSecGatewayDataSourceGetByID(resourceName, dataSourceName, dataSourceAttributeID string) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateIPSecGatewayDataIDConfig, dataSourceName, dataSourceAttributeID, constant.IPSecGatewayResource+"."+resourceName+"."+gatewayAttributeLocation,
	)
	baseConfig := configIPSecGatewayBasic(resourceName, gatewayAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

func configIPSecGatewayDataSourceGetByName(resourceName, dataSourceName, dataSourceAttributeName string) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateIPSecGatewayDataNameConfig, dataSourceName, dataSourceAttributeName, constant.IPSecGatewayResource+"."+resourceName+"."+gatewayAttributeLocation,
	)
	baseConfig := configIPSecGatewayBasic(resourceName, gatewayAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

const (
	gatewayResourceName = "test_ipsec_gateway"
	gatewayDataName     = "test_ipsec_gateway_ds"

	gatewayAttributeName      = "name"
	gatewayAttributeNameValue = "ipsec_gateway_test"

	gatewayAttributeVersion      = "version"
	gatewayAttributeVersionValue = "IKEv2"

	gatewayAttributeIP = "gateway_ip"

	gatewayAttributeLocation      = "location"
	gatewayAttributeLocationValue = "de/fra"
)

const templateIPSecGatewayConfig = `
resource "ionoscloud_vpn_ipsec_gateway" "%v" {
	name = "%v"
	version = "%v"
	gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
	location = "%v"
	connections {
		datacenter_id = ionoscloud_datacenter.test_datacenter.id
		lan_id = ionoscloud_lan.test_lan.id
		ipv4_cidr = local.ipv4_cidr_block
	}
	maintenance_window {
    	day_of_the_week       = "Monday"
    	time                  = "09:00:00"
	}
  	tier = "STANDARD"
}`

const templateIPSecGatewayDataIDConfig = `
data "ionoscloud_vpn_ipsec_gateway" "%v" {
  	id = %v
	location = %v
}`

const templateIPSecGatewayDataNameConfig = `
data "ionoscloud_vpn_ipsec_gateway" "%v" {
	name = %v
	location = %v
}`

const defaultIPSecGatewayBaseConfig = `
resource "ionoscloud_datacenter" "test_datacenter" {
	name = "test_datacenter"
	location = "de/fra"
	sec_auth_protection = false
}

resource "ionoscloud_lan" "test_lan" {
	name = "test_lan"
	public = false
	datacenter_id = ionoscloud_datacenter.test_datacenter.id
	ipv6_cidr_block = local.lan_ipv6_cidr_block
}

resource "ionoscloud_ipblock" "test_ipblock" {
	name = "test_ipblock"
	location = "de/fra"
	size = 1
}

resource "ionoscloud_server" "test_server" {
	name = "test_server"
	datacenter_id = ionoscloud_datacenter.test_datacenter.id
	cores = 1
	ram = 2048
	image_name = "ubuntu:latest"
	image_password = random_password.server_image_password.result

	nic {
        lan = ionoscloud_lan.test_lan.id
        name = "test_nic"
        dhcp = true
		dhcpv6 = false
		ipv6_cidr_block = local.ipv6_cidr_block
        firewall_active   = false
    }

	volume {
		name         = "test_volume"
		disk_type    = "HDD"
		size         = 10
		licence_type = "OTHER"
		ssh_key_path = null
	}
}

resource "random_password" "server_image_password" {
	length           = 16
	special          = false
}

locals {
	lan_ipv6_cidr_block_parts = split("/", ionoscloud_datacenter.test_datacenter.ipv6_cidr_block)
	lan_ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "64")

	ipv4_cidr_block = format("%s/%s", ionoscloud_server.test_server.nic[0].ips[0], "24")
	ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "80")
}
`
