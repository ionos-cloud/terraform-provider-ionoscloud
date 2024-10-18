package vpn

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const ipsecTunnelResourceName = "VPN IPSec Tunnel"

var (
	// IPSecTunnelDiffieHellmanGroups is a list of supported Diffie-Hellman groups
	IPSecTunnelDiffieHellmanGroups = []string{"15-MODP3072", "16-MODP4096", "19-ECP256", "20-ECP384", "21-ECP521", "28-ECP256BP", "29-ECP384BP", "30-ECP512BP"}

	// IPSecTunnelEncryptionAlgorithms is a list of supported encryption algorithms
	IPSecTunnelEncryptionAlgorithms = []string{
		"AES128", "AES256", "AES128-CTR", "AES256-CTR", "AES128-GCM-16", "AES256-GCM-16", "AES128-GCM-12", "AES256-GCM-12", "AES128-CCM-12", "AES256-CCM-12",
	}

	// IPSecTunnelIntegrityAlgorithms is a list of supported integrity algorithms
	IPSecTunnelIntegrityAlgorithms = []string{"SHA256", "SHA384", "SHA512", "AES-XCBC"}
)

// CreateIPSecTunnel creates a new VPN IPSec Tunnel
func (c *Client) CreateIPSecTunnel(ctx context.Context, d *schema.ResourceData) (vpn.IPSecTunnelRead, *shared.APIResponse, error) {
	c.changeConfigURL(d.Get("location").(string))
	gatewayID := d.Get("gateway_id").(string)

	request := setIPSecTunnelCreateRequest(d)
	tunnel, apiResponse, err := c.sdkClient.IPSecTunnelsApi.IpsecgatewaysTunnelsPost(ctx, gatewayID).IPSecTunnelCreate(request).Execute()
	apiResponse.LogInfo()
	return tunnel, apiResponse, err
}

// GetIPSecTunnelByID retrieves a VPN IPSec Tunnel by its ID and location
func (c *Client) GetIPSecTunnelByID(ctx context.Context, id string, gatewayID string, location string) (vpn.IPSecTunnelRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)

	tunnel, apiResponse, err := c.sdkClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return tunnel, apiResponse, err
}

// ListIPSecTunnel retrieves all VPN IPSec Tunnels from a given gateway and location
func (c *Client) ListIPSecTunnel(ctx context.Context, gatewayID string, location string) (vpn.IPSecTunnelReadList, *shared.APIResponse, error) {
	c.changeConfigURL(location)

	gateways, apiResponse, err := c.sdkClient.IPSecTunnelsApi.IpsecgatewaysTunnelsGet(ctx, gatewayID).Execute()
	apiResponse.LogInfo()
	return gateways, apiResponse, err
}

// DeleteIPSecTunnel deletes a VPN IPSec Tunnel using its ID and location
func (c *Client) DeleteIPSecTunnel(ctx context.Context, id string, gatewayID string, location string) (utils.ApiResponseInfo, error) {
	c.changeConfigURL(location)

	apiResponse, err := c.sdkClient.IPSecTunnelsApi.IpsecgatewaysTunnelsDelete(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateIPSecTunnel updates a VPN IPSec Tunnel
func (c *Client) UpdateIPSecTunnel(ctx context.Context, d *schema.ResourceData) (vpn.IPSecTunnelRead, *shared.APIResponse, error) {
	c.changeConfigURL(d.Get("location").(string))
	gatewayID := d.Get("gateway_id").(string)

	request := setIPSecTunnelPutRequest(d)
	tunnel, apiResponse, err := c.sdkClient.IPSecTunnelsApi.IpsecgatewaysTunnelsPut(ctx, gatewayID, d.Id()).IPSecTunnelEnsure(request).Execute()
	apiResponse.LogInfo()
	return tunnel, apiResponse, err
}

// IsIPSecTunnelReady checks if a VPN IPSec Tunnel is ready to use
func (c *Client) IsIPSecTunnelReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	id := d.Id()
	location := d.Get("location").(string)
	gatewayID := d.Get("gateway_id").(string)

	tunnel, _, err := c.GetIPSecTunnelByID(ctx, id, gatewayID, location)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] VPN IPSec Gateway Tunnel state %s", tunnel.Metadata.Status)

	return strings.EqualFold(tunnel.Metadata.Status, constant.Available), nil
}

// IsIPSecTunnelDeleted checks if a VPN IPSec Tunnel has been deleted
func (c *Client) IsIPSecTunnelDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	id := d.Id()
	location := d.Get("location").(string)
	gatewayID := d.Get("gateway_id").(string)

	_, apiResponse, err := c.GetIPSecTunnelByID(ctx, id, gatewayID, location)
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// SetIPSecTunnelData sets the VPN IPSec Tunnel data to the terraform schema
func SetIPSecTunnelData(d *schema.ResourceData, tunnel vpn.IPSecTunnelRead) error {
	d.SetId(tunnel.Id)

	if err := d.Set("name", tunnel.Properties.Name); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "name", err)
	}

	if tunnel.Properties.Description != nil {
		if err := d.Set("description", *tunnel.Properties.Description); err != nil {
			return utils.GenerateSetError(ipsecGatewayResourceName, "description", err)
		}
	}

	if err := d.Set("remote_host", tunnel.Properties.RemoteHost); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "remote_host", err)
	}

	if err := d.Set("cloud_network_cidrs", tunnel.Properties.CloudNetworkCIDRs); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "cloud_network_cidrs", err)
	}

	if err := d.Set("peer_network_cidrs", tunnel.Properties.PeerNetworkCIDRs); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "peer_network_cidrs", err)
	}

	auth := map[string]interface{}{
		"method": tunnel.Properties.Auth.Method,
	}

	if err := d.Set("auth", []interface{}{auth}); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "auth", err)
	}

	ike := map[string]interface{}{
		"diffie_hellman_group": tunnel.Properties.Ike.DiffieHellmanGroup,
		"encryption_algorithm": tunnel.Properties.Ike.EncryptionAlgorithm,
		"integrity_algorithm":  tunnel.Properties.Ike.IntegrityAlgorithm,
		"lifetime":             tunnel.Properties.Ike.Lifetime,
	}
	if err := d.Set("ike", []interface{}{ike}); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "ike", err)
	}

	esp := map[string]interface{}{
		"diffie_hellman_group": tunnel.Properties.Esp.DiffieHellmanGroup,
		"encryption_algorithm": tunnel.Properties.Esp.EncryptionAlgorithm,
		"integrity_algorithm":  tunnel.Properties.Esp.IntegrityAlgorithm,
		"lifetime":             tunnel.Properties.Esp.Lifetime,
	}
	if err := d.Set("esp", []interface{}{esp}); err != nil {
		return utils.GenerateSetError(ipsecTunnelResourceName, "esp", err)
	}

	return nil
}

func setIPSecTunnelCreateRequest(d *schema.ResourceData) vpn.IPSecTunnelCreate {
	properties := setIPSecTunnelProperties(d)

	return vpn.IPSecTunnelCreate{Properties: properties}
}

func setIPSecTunnelPutRequest(d *schema.ResourceData) vpn.IPSecTunnelEnsure {
	properties := setIPSecTunnelProperties(d)

	return vpn.IPSecTunnelEnsure{Id: d.Id(), Properties: properties}
}

func setIPSecTunnelProperties(d *schema.ResourceData) vpn.IPSecTunnel {
	properties := vpn.IPSecTunnel{}

	properties.Name = d.Get("name").(string)
	properties.RemoteHost = d.Get("remote_host").(string)

	for _, v := range d.Get("cloud_network_cidrs").([]interface{}) {
		properties.CloudNetworkCIDRs = append(properties.CloudNetworkCIDRs, v.(string))
	}

	for _, v := range d.Get("peer_network_cidrs").([]interface{}) {
		properties.PeerNetworkCIDRs = append(properties.PeerNetworkCIDRs, v.(string))
	}

	properties.Auth = vpn.IPSecTunnelAuth{
		Method: d.Get("auth.0.method").(string),
	}
	if v, ok := d.GetOk("auth.0.psk_key"); ok {
		properties.Auth.Psk = &vpn.IPSecPSK{
			Key: v.(string),
		}
	}

	properties.Ike = vpn.IKEEncryption{
		DiffieHellmanGroup:  shared.ToPtr(d.Get("ike.0.diffie_hellman_group").(string)),
		EncryptionAlgorithm: shared.ToPtr(d.Get("ike.0.encryption_algorithm").(string)),
		IntegrityAlgorithm:  shared.ToPtr(d.Get("ike.0.integrity_algorithm").(string)),
		Lifetime:            shared.ToPtr(int32(d.Get("ike.0.lifetime").(int))),
	}

	properties.Esp = vpn.ESPEncryption{
		DiffieHellmanGroup:  shared.ToPtr(d.Get("esp.0.diffie_hellman_group").(string)),
		EncryptionAlgorithm: shared.ToPtr(d.Get("esp.0.encryption_algorithm").(string)),
		IntegrityAlgorithm:  shared.ToPtr(d.Get("esp.0.integrity_algorithm").(string)),
		Lifetime:            shared.ToPtr(int32(d.Get("esp.0.lifetime").(int))),
	}

	if v, ok := d.GetOk("description"); ok {
		properties.Description = shared.ToPtr(v.(string))
	}

	return properties
}
