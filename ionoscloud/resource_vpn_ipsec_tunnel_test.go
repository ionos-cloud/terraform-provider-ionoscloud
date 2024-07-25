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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestVpnIPSecTunnelResource(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ProviderFactories: testAccProviderFactories,
			ExternalProviders: randomProviderVersion343(),
			CheckDestroy:      testCheckIPSecTunnelDestroy,
			Steps: []resource.TestStep{
				{
					Config: configIPSecTunnelBasic(tunnelResourceName, tunnelAttributeNameValue),
					Check:  checkIPSecTunnelResource(tunnelAttributeNameValue),
				},
				{
					Config: configIPSecTunnelBasic(tunnelResourceName, fmt.Sprintf("%v_updated", tunnelAttributeNameValue)),
					Check:  checkIPSecTunnelResource(fmt.Sprintf("%v_updated", tunnelAttributeNameValue)),
				},
				{
					Config: configIPSecTunnelDataSourceGetByID(tunnelResourceName, tunnelDataName, constant.IPSecTunnelResource+"."+tunnelResourceName+".id"),
					Check: checkIPSecTunnelResourceAttributesComparative(
						constant.DataSource+"."+constant.IPSecTunnelResource+"."+tunnelDataName, constant.IPSecTunnelResource+"."+tunnelResourceName,
					),
				},
				{
					Config: configIPSecTunnelDataSourceGetByName(
						tunnelResourceName, tunnelDataName, constant.IPSecTunnelResource+"."+tunnelResourceName+"."+tunnelAttributeName,
					),
					Check: checkIPSecTunnelResourceAttributesComparative(
						constant.DataSource+"."+constant.IPSecTunnelResource+"."+tunnelDataName, constant.IPSecTunnelResource+"."+tunnelResourceName,
					),
				},
				{
					Config:      configIPSecTunnelDataSourceGetByName(tunnelResourceName, tunnelDataName, `"willnotwork"`),
					ExpectError: regexp.MustCompile(`no VPN IPSec Gateway Tunnel found with the specified name =`),
				},
			},
		},
	)
}

func testCheckIPSecTunnelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).VPNClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.IPSecTunnelResource {
			continue
		}

		_, apiResponse, err := client.GetIPSecTunnelByID(ctx, rs.Primary.ID, rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["location"])
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

func testAccCheckIPSecTunnelExists(n string) resource.TestCheckFunc {
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

		foundTunnel, _, err := client.GetIPSecTunnelByID(ctx, rs.Primary.ID, rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching IPSec Gateway Tunnel with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if foundTunnel.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

func checkIPSecTunnelResource(attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckIPSecTunnelExists(constant.IPSecTunnelResource+"."+tunnelResourceName),
		checkIPSecTunnelResourceAttributes(constant.IPSecTunnelResource+"."+tunnelResourceName, attributeNameReferenceValue),
	)
}

func checkIPSecTunnelResourceAttributes(fullResourceName, attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeName, attributeNameReferenceValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeRemoteHost, tunnelAttributeRemoteHostValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeAuthMethod, tunnelAttributeAuthMethodValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeAuthPSKKey, tunnelAttributeAuthPSKKeyValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeIKE+tunnelAttributeDiffieHellmanGroup, tunnelAttributeDiffieHellmanGroupValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeIKE+tunnelAttributeEncryptionAlgorithm, tunnelAttributeEncryptionAlgorithmValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeIKE+tunnelAttributeIntegrityAlgorithm, tunnelAttributeIntegrityAlgorithmValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeIKE+tunnelAttributeLifetime, tunnelAttributeIKELifetimeValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeESP+tunnelAttributeDiffieHellmanGroup, tunnelAttributeDiffieHellmanGroupValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeESP+tunnelAttributeEncryptionAlgorithm, tunnelAttributeEncryptionAlgorithmValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeESP+tunnelAttributeIntegrityAlgorithm, tunnelAttributeIntegrityAlgorithmValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeESP+tunnelAttributeLifetime, tunnelAttributeESPLifetimeValue),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributeCloudNetworkCIDRs, "1"),
		resource.TestCheckResourceAttr(fullResourceName, tunnelAttributePeerNetworkCIDRs, "1"),
	)
}

func checkIPSecTunnelResourceAttributesComparative(fullResourceName, fullReferenceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributeName, fullReferenceResourceName, tunnelAttributeName),
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributeRemoteHost, fullReferenceResourceName, tunnelAttributeRemoteHost),
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributeAuthMethod, fullReferenceResourceName, tunnelAttributeAuthMethod),
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributeAuthPSKKey, fullReferenceResourceName, tunnelAttributeAuthPSKKey),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeIKE+tunnelAttributeDiffieHellmanGroup, fullReferenceResourceName, tunnelAttributeIKE+tunnelAttributeDiffieHellmanGroup,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeIKE+tunnelAttributeEncryptionAlgorithm, fullReferenceResourceName, tunnelAttributeIKE+tunnelAttributeEncryptionAlgorithm,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeIKE+tunnelAttributeIntegrityAlgorithm, fullReferenceResourceName, tunnelAttributeIKE+tunnelAttributeIntegrityAlgorithm,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeIKE+tunnelAttributeLifetime, fullReferenceResourceName, tunnelAttributeIKE+tunnelAttributeLifetime,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeESP+tunnelAttributeDiffieHellmanGroup, fullReferenceResourceName, tunnelAttributeESP+tunnelAttributeDiffieHellmanGroup,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeESP+tunnelAttributeEncryptionAlgorithm, fullReferenceResourceName, tunnelAttributeESP+tunnelAttributeEncryptionAlgorithm,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeESP+tunnelAttributeIntegrityAlgorithm, fullReferenceResourceName, tunnelAttributeESP+tunnelAttributeIntegrityAlgorithm,
		),
		resource.TestCheckResourceAttrPair(
			fullResourceName, tunnelAttributeESP+tunnelAttributeLifetime, fullReferenceResourceName, tunnelAttributeESP+tunnelAttributeLifetime,
		),
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributeCloudNetworkCIDRs, fullReferenceResourceName, tunnelAttributeCloudNetworkCIDRs),
		resource.TestCheckResourceAttrPair(fullResourceName, tunnelAttributePeerNetworkCIDRs, fullReferenceResourceName, tunnelAttributePeerNetworkCIDRs),
	)
}

func configIPSecTunnelBasic(resourceName, attributeName string) string {
	tunnelBasicConfig := fmt.Sprintf(
		templateIPSecTunnelConfig, tunnelResourceName, attributeName, tunnelAttributeRemoteHostValue, tunnelAttributeAuthMethodValue,
		tunnelAttributeAuthPSKKeyValue, tunnelAttributeDiffieHellmanGroupValue, tunnelAttributeEncryptionAlgorithmValue, tunnelAttributeIntegrityAlgorithmValue,
		tunnelAttributeIKELifetimeValue, tunnelAttributeDiffieHellmanGroupValue, tunnelAttributeEncryptionAlgorithmValue, tunnelAttributeIntegrityAlgorithmValue,
		tunnelAttributeESPLifetimeValue, tunnelAttributeCloudNetworkCIDRsValue, tunnelAttributePeerNetworkCIDRsValue,
	)

	return strings.Join([]string{defaultIPSecTunnelBaseConfig, tunnelBasicConfig}, "\n")
}

func configIPSecTunnelDataSourceGetByID(resourceName, dataSourceName, dataSourceAttributeID string) string {
	dataSourceBasicConfig := fmt.Sprintf(templateIPSecTunnelDataIDConfig, dataSourceName, dataSourceAttributeID)
	baseConfig := configIPSecTunnelBasic(resourceName, tunnelAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

func configIPSecTunnelDataSourceGetByName(resourceName, dataSourceName, dataSourceAttributeName string) string {
	dataSourceBasicConfig := fmt.Sprintf(templateIPSecTunnelDataNameConfig, dataSourceName, dataSourceAttributeName)
	baseConfig := configIPSecTunnelBasic(resourceName, tunnelAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

const (
	tunnelResourceName = "test_ipsec_tunnel"
	tunnelDataName     = "test_ipsec_tunnel_ds"

	tunnelAttributeName      = "name"
	tunnelAttributeNameValue = "ipsec_tunnel_test"

	tunnelAttributeRemoteHost      = "remote_host"
	tunnelAttributeRemoteHostValue = "vpn.example.com"

	tunnelAttributeAuthMethod      = "auth.0.method"
	tunnelAttributeAuthMethodValue = "PSK"

	tunnelAttributeAuthPSKKey      = "auth.0.psk_key"
	tunnelAttributeAuthPSKKeyValue = "X2wosbaw74M8hQGbK3jCCaEusR6CCFRa"

	tunnelAttributeIKE = "ike.0."
	tunnelAttributeESP = "esp.0."

	tunnelAttributeDiffieHellmanGroup      = "diffie_hellman_group"
	tunnelAttributeDiffieHellmanGroupValue = "16-MODP4096"

	tunnelAttributeEncryptionAlgorithm      = "encryption_algorithm"
	tunnelAttributeEncryptionAlgorithmValue = "AES256"

	tunnelAttributeIntegrityAlgorithm      = "integrity_algorithm"
	tunnelAttributeIntegrityAlgorithmValue = "SHA256"

	tunnelAttributeLifetime         = "lifetime"
	tunnelAttributeIKELifetimeValue = "86400"
	tunnelAttributeESPLifetimeValue = "3600"

	tunnelAttributeCloudNetworkCIDRs      = "cloud_network_cidrs.#"
	tunnelAttributeCloudNetworkCIDRsValue = `"0.0.0.0/0"`

	tunnelAttributePeerNetworkCIDRs      = "peer_network_cidrs.#"
	tunnelAttributePeerNetworkCIDRsValue = `"1.2.3.4/32"`

	tunnelAttributeLocation = "location"
)

const templateIPSecTunnelConfig = `
resource "ionoscloud_vpn_ipsec_tunnel" "%v" {
	name = "%v"
	location = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.location
	gateway_id = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.id
	remote_host = "%v"

	auth {
		method = "%v"
		psk_key = "%v"
	}

	ike {
      diffie_hellman_group = "%v"
      encryption_algorithm = "%v"
      integrity_algorithm = "%v"
      lifetime = %v
    }

    esp {
      diffie_hellman_group = "%v"
      encryption_algorithm = "%v"
      integrity_algorithm = "%v"
      lifetime = %v
    }

    cloud_network_cidrs = [
      %v
    ]

    peer_network_cidrs = [
      %v
    ]
}`

const templateIPSecTunnelDataIDConfig = `
data "ionoscloud_vpn_ipsec_tunnel" "%v" {
  	id = %v
	gateway_id = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.id
	location = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.location
}`

const templateIPSecTunnelDataNameConfig = `
data "ionoscloud_vpn_ipsec_tunnel" "%v" {
	name = %v
	gateway_id = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.id
	location = ionoscloud_vpn_ipsec_gateway.test_ipsec_gateway.location
}`

const defaultIPSecTunnelBaseConfig = `
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

resource "ionoscloud_vpn_ipsec_gateway" "test_ipsec_gateway" {
	name = "test_ipsec_gateway"
	version = "IKEv2"
	gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
	location = "de/fra"
	connections {
		datacenter_id = ionoscloud_datacenter.test_datacenter.id
		lan_id = ionoscloud_lan.test_lan.id
		ipv4_cidr = local.ipv4_cidr_block
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
