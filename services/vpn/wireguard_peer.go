package vpn

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// CreateWireguardGatewayPeers creates a new wireguard peer
func (c *Client) CreateWireguardGatewayPeers(ctx context.Context, d *schema.ResourceData, gatewayID string) (vpn.WireguardPeerRead, utils.ApiResponseInfo, error) {
	c.changeConfigURL(d.Get("location").(string))
	request, err := setWireguardPeersPostRequest(d)
	if err != nil {
		return vpn.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPost(ctx, gatewayID).WireguardPeerCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// IsWireguardPeerAvailable checks if the wireguard peer is available
func (c *Client) IsWireguardPeerAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	c.changeConfigURL(d.Get("location").(string))
	wireguardID := d.Id()
	gatewayID := d.Get("gateway_id").(string)
	location := d.Get("location").(string)
	wireguard, _, err := c.GetWireguardPeerByID(ctx, gatewayID, wireguardID, location)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

// UpdateWireguardPeer updates a wireguard peer
func (c *Client) UpdateWireguardPeer(ctx context.Context, gatewayID, id string, d *schema.ResourceData) (vpn.WireguardPeerRead, utils.ApiResponseInfo, error) {
	c.changeConfigURL(d.Get("location").(string))
	request, err := setWireguardPeerPatchRequest(d)
	if err != nil {
		return vpn.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPut(ctx, gatewayID, id).WireguardPeerEnsure(request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

// DeleteWireguardPeer deletes a wireguard peer
func (c *Client) DeleteWireguardPeer(ctx context.Context, gatewayID, id, location string) (utils.ApiResponseInfo, error) {
	c.changeConfigURL(location)
	apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersDelete(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsWireguardPeerDeleted checks if the wireguard peer is deleted
func (c *Client) IsWireguardPeerDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	c.changeConfigURL(d.Get("location").(string))
	gatewayID := d.Get("gateway_id").(string)
	_, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetWireguardPeerByID returns a wireguard by its ID
func (c *Client) GetWireguardPeerByID(ctx context.Context, gatewayID, id, location string) (vpn.WireguardPeerRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// ListWireguardPeers returns a list of all wireguards
func (c *Client) ListWireguardPeers(ctx context.Context, gatewayID, location string) (vpn.WireguardPeerReadList, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	wireguards, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersGet(ctx, gatewayID).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

// IsWireguardPeerReady checks if the wireguard peer is ready
func (c *Client) IsWireguardPeerReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	location := d.Get("location").(string)
	c.changeConfigURL(location)
	gatewayID := d.Get("gateway_id").(string)
	cluster, _, err := c.GetWireguardPeerByID(ctx, gatewayID, d.Id(), location)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wierguard gateway state %s", cluster.Metadata.Status)
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

func setWireguardPeersPostRequest(d *schema.ResourceData) (*vpn.WireguardPeerCreate, error) {
	request := vpn.WireguardPeerCreate{Properties: vpn.WireguardPeer{}}
	name := d.Get("name").(string)

	request.Properties.Name = name

	if value, ok := d.GetOk("description"); ok {
		valueStr := value.(string)
		request.Properties.Description = &valueStr
	}

	if _, ok := d.GetOk("endpoint"); ok {
		request.Properties.Endpoint = getEndpointData(d)
	}
	if v, ok := d.GetOk("allowed_ips"); ok {
		raw := v.([]interface{})
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return nil, err
		}
		if len(ips) > 0 {
			request.Properties.AllowedIPs = ips
		}
	}
	if value, ok := d.GetOk("public_key"); ok {
		valueStr := value.(string)
		request.Properties.PublicKey = valueStr
	}

	return &request, nil
}
func getEndpointData(d *schema.ResourceData) *vpn.WireguardEndpoint {
	endpoint := vpn.NewWireguardEndpointWithDefaults()
	if endpointValues, ok := d.GetOk("endpoint"); ok {
		endpointMap := endpointValues.([]any)
		if endpointMap != nil {
			if host, ok := d.GetOk("endpoint.0.host"); ok {
				host := host.(string)
				endpoint.Host = host
			}
			if port, ok := d.GetOk("endpoint.0.port"); ok {
				port := port.(int)
				endpoint.Port = shared.ToPtr(int32(port))
			}
		}
	}
	return endpoint

}

func setWireguardPeerPatchRequest(d *schema.ResourceData) (vpn.WireguardPeerEnsure, error) {
	request := vpn.WireguardPeerEnsure{Properties: vpn.WireguardPeer{}}

	request.Id = d.Id()
	request.Properties.Name = d.Get("name").(string)
	request.Properties.PublicKey = d.Get("public_key").(string)

	if v, ok := d.GetOk("description"); ok {
		request.Properties.Description = shared.ToPtr(v.(string))
	}
	if _, ok := d.GetOk("endpoint"); ok {
		request.Properties.Endpoint = getEndpointData(d)
	}

	if v, ok := d.GetOk("allowed_ips"); ok {
		raw := v.([]interface{})
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return request, err
		}
		if len(ips) > 0 {
			request.Properties.AllowedIPs = ips
		}
	}

	return request, nil
}

var resPeerName = "vpnSdk wireguard peer"

// SetWireguardPeerData sets the wireguard peer data
func SetWireguardPeerData(d *schema.ResourceData, wireguard vpn.WireguardPeerRead) error {
	d.SetId(wireguard.Id)

	if err := d.Set("name", wireguard.Properties.Name); err != nil {
		return utils.GenerateSetError(resPeerName, "name", err)
	}
	if err := d.Set("description", wireguard.Properties.Description); err != nil {
		return utils.GenerateSetError(resPeerName, "description", err)
	}
	if wireguard.Properties.Endpoint != nil {
		endpointEntry, err := wireguard.Properties.Endpoint.ToMap()
		if err != nil {
			return err
		}

		if len(endpointEntry) > 0 {
			var endpoiontSlice []any
			endpoiontSlice = append(endpoiontSlice, endpointEntry)
			if err := d.Set("endpoint", endpoiontSlice); err != nil {
				return utils.GenerateSetError(resPeerName, "endpoint", err)
			}
		}

	}

	if err := d.Set("allowed_ips", wireguard.Properties.AllowedIPs); err != nil {
		return utils.GenerateSetError(resPeerName, "allowed_ips", err)
	}
	if err := d.Set("public_key", wireguard.Properties.PublicKey); err != nil {
		return utils.GenerateSetError(resPeerName, "public_key", err)
	}

	if err := d.Set("status", wireguard.Metadata.Status); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "status", err)
	}

	return nil
}
