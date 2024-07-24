package vpn

import (
	"context"
	"fmt"
	"log"
	"strings"

	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var wireguardResourceName = "vpnSdk wireguard gateway"

// CreateWireguardGateway creates a new wireguard gateway
func (c *Client) CreateWireguardGateway(ctx context.Context, d *schema.ResourceData) (vpnSdk.WireguardGatewayRead, utils.ApiResponseInfo, error) {
	request := setWireguardGWPostRequest(d)
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPost(ctx).WireguardGatewayCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// CreateWireguardGatewayPeers creates a new wireguard peer
func (c *Client) CreateWireguardGatewayPeers(ctx context.Context, d *schema.ResourceData, gatewayID string) (vpnSdk.WireguardPeerRead, utils.ApiResponseInfo, error) {
	request, err := setWireguardPeersPostRequest(d)
	if err != nil {
		return vpnSdk.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPost(ctx, gatewayID).WireguardPeerCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// IsWireguardAvailable checks if the wireguard is available
func (c *Client) IsWireguardAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	wireguardID := d.Id()
	wireguard, _, err := c.GetWireguardGatewayByID(ctx, wireguardID)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

// IsWireguardPeerAvailable checks if the wireguard peer is available
func (c *Client) IsWireguardPeerAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	wireguardID := d.Id()
	gatewayID := d.Get("gateway_id").(string)
	wireguard, _, err := c.GetWireguardPeerByID(ctx, gatewayID, wireguardID)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

// UpdateWireguardGateway updates a wireguard gateway
func (c *Client) UpdateWireguardGateway(ctx context.Context, id string, d *schema.ResourceData) (vpnSdk.WireguardGatewayRead, utils.ApiResponseInfo, error) {
	request := setWireguardGatewayPutRequest(d)
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPut(ctx, id).WireguardGatewayEnsure(*request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

// UpdateWireguardPeer updates a wireguard peer
func (c *Client) UpdateWireguardPeer(ctx context.Context, gatewayID, id string, d *schema.ResourceData) (vpnSdk.WireguardPeerRead, utils.ApiResponseInfo, error) {
	request, err := setWireguardPeerPatchRequest(d)
	if err != nil {
		return vpnSdk.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPut(ctx, gatewayID, id).WireguardPeerEnsure(request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

// DeleteWireguardGateway deletes a wireguard gateway
func (c *Client) DeleteWireguardGateway(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// DeleteWireguardPeer deletes a wireguard peer
func (c *Client) DeleteWireguardPeer(ctx context.Context, gatewayID, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersDelete(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsWireguardGatewayDeleted checks if the wireguard gateway is deleted
func (c *Client) IsWireguardGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// IsWireguardPeerDeleted checks if the wireguard peer is deleted
func (c *Client) IsWireguardPeerDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Get("gateway_id").(string)
	_, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetWireguardGatewayByID returns a wireguard by its ID
func (c *Client) GetWireguardGatewayByID(ctx context.Context, id string) (vpnSdk.WireguardGatewayRead, *shared.APIResponse, error) {
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// GetWireguardPeerByID returns a wireguard by its ID
func (c *Client) GetWireguardPeerByID(ctx context.Context, gatewayID, id string) (vpnSdk.WireguardPeerRead, *shared.APIResponse, error) {
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// ListWireguardGateways returns a list of all wireguards
func (c *Client) ListWireguardGateways(ctx context.Context) (vpnSdk.WireguardGatewayReadList, *shared.APIResponse, error) {
	wireguards, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

// ListWireguardPeers returns a list of all wireguards
func (c *Client) ListWireguardPeers(ctx context.Context, gatewayID string) (vpnSdk.WireguardPeerReadList, *shared.APIResponse, error) {
	wireguards, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersGet(ctx, gatewayID).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

// IsWireguardGatewayReady checks if the wireguard gateway is ready
func (c *Client) IsWireguardGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cluster, _, err := c.GetWireguardGatewayByID(ctx, d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wierguard gateway state %s", cluster.Metadata.Status)
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

// IsWireguardPeerReady checks if the wireguard peer is ready
func (c *Client) IsWireguardPeerReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Get("gateway_id").(string)
	cluster, _, err := c.GetWireguardPeerByID(ctx, gatewayID, d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wierguard gateway state %s", cluster.Metadata.Status)
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

func setWireguardGWPostRequest(d *schema.ResourceData) *vpnSdk.WireguardGatewayCreate {
	request := vpnSdk.WireguardGatewayCreate{Properties: vpnSdk.WireguardGateway{}}
	name := d.Get("name").(string)
	gatewayIP := d.Get("gateway_ip").(string)
	privateKey := d.Get("private_key").(string)

	request.Properties.Name = name
	request.Properties.GatewayIP = gatewayIP
	request.Properties.PrivateKey = privateKey

	if value, ok := d.GetOk("description"); ok {
		valueStr := value.(string)
		request.Properties.Description = &valueStr
	}
	if value, ok := d.GetOk("interface_ipv4_cidr"); ok {
		valueStr := value.(string)
		request.Properties.InterfaceIPv4CIDR = &valueStr
	}
	if value, ok := d.GetOk("interface_ipv6_cidr"); ok {
		valueStr := value.(string)
		request.Properties.InterfaceIPv6CIDR = &valueStr
	}
	if value, ok := d.GetOk("listenPort"); ok {
		valueStr := (int32)(value.(int))
		request.Properties.ListenPort = &valueStr
	}
	request.Properties.Connections = getWireguardGwConnectionsData(d)

	return &request
}

func setWireguardPeersPostRequest(d *schema.ResourceData) (*vpnSdk.WireguardPeerCreate, error) {
	request := vpnSdk.WireguardPeerCreate{Properties: vpnSdk.WireguardPeer{}}
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
func getEndpointData(d *schema.ResourceData) *vpnSdk.WireguardEndpoint {
	endpoint := vpnSdk.NewWireguardEndpointWithDefaults()
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

func getWireguardGwConnectionsData(d *schema.ResourceData) []vpnSdk.Connection {
	connections := make([]vpnSdk.Connection, 0)

	if connectionValues, ok := d.GetOk("connections"); ok {
		connectionsItf := connectionValues.([]any)
		if connectionsItf != nil {
			for idx := range connectionsItf {
				connection := vpnSdk.Connection{}
				if datacenterID, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", idx)); ok {
					datacenterID := datacenterID.(string)
					connection.DatacenterId = datacenterID
				}
				if lanID, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", idx)); ok {
					lanID := lanID.(string)
					connection.LanId = lanID
				}
				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.ipv4_cidr", idx)); ok {
					cidr := cidr.(string)
					connection.Ipv4CIDR = cidr
				}

				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.ipv6_cidr", idx)); ok {
					cidr := cidr.(string)
					connection.Ipv6CIDR = &cidr
				}

				connections = append(connections, connection)
			}
		}
	}

	return connections
}

func setWireguardGatewayPutRequest(d *schema.ResourceData) *vpnSdk.WireguardGatewayEnsure {
	request := vpnSdk.WireguardGatewayEnsure{Properties: vpnSdk.WireguardGateway{}}
	request.Id = d.Id()
	request.Properties.GatewayIP = d.Get("gateway_ip").(string)
	request.Properties.Name = d.Get("name").(string)
	request.Properties.PrivateKey = d.Get("private_key").(string)
	request.Properties.Connections = getWireguardGwConnectionsData(d)
	if val, ok := d.GetOk("interface_ipv4_cidr"); ok {
		request.Properties.InterfaceIPv4CIDR = shared.ToPtr(val.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		request.Properties.Description = shared.ToPtr(v.(string))
	}
	if v, ok := d.GetOk("interface_ipv6_cidr"); ok {
		request.Properties.InterfaceIPv6CIDR = shared.ToPtr(v.(string))
	}
	if v, ok := d.GetOk("listen_port"); ok {
		request.Properties.ListenPort = shared.ToPtr(int32(v.(int)))
	}
	return &request
}

func setWireguardPeerPatchRequest(d *schema.ResourceData) (vpnSdk.WireguardPeerEnsure, error) {
	request := vpnSdk.WireguardPeerEnsure{Properties: vpnSdk.WireguardPeer{}}

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

// SetWireguardGWData sets the wireguard gateway data
func SetWireguardGWData(d *schema.ResourceData, wireguard vpnSdk.WireguardGatewayRead) error {
	d.SetId(wireguard.Id)

	if err := d.Set("name", wireguard.Properties.Name); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "name", err)
	}
	if err := d.Set("description", wireguard.Properties.Description); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "description", err)
	}
	if err := d.Set("gateway_ip", wireguard.Properties.GatewayIP); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "gateway_ip", err)
	}
	if err := d.Set("public_key", wireguard.Metadata.PublicKey); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "public_key", err)
	}
	if err := d.Set("interface_ipv4_cidr", wireguard.Properties.InterfaceIPv4CIDR); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "interface_ipv4_cidr", err)
	}
	if err := d.Set("interface_ipv6_cidr", wireguard.Properties.InterfaceIPv6CIDR); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "interface_ipv6_cidr", err)
	}

	var connections []map[string]any
	for _, connection := range wireguard.Properties.Connections {
		connection, err := utils.DecodeStructToMap(connection)
		if err != nil {
			return err
		}
		connections = append(connections, connection)
	}
	if err := d.Set("connections", connections); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "connections", err)
	}

	if err := d.Set("listen_port", wireguard.Properties.ListenPort); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "listenPort", err)
	}
	if err := d.Set("status", wireguard.Metadata.Status); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "status", err)
	}

	return nil
}

var resPeerName = "vpnSdk wireguard peer"

// SetWireguardPeerData sets the wireguard peer data
func SetWireguardPeerData(d *schema.ResourceData, wireguard vpnSdk.WireguardPeerRead) error {
	d.SetId(wireguard.Id)

	if err := d.Set("name", wireguard.Properties.Name); err != nil {
		return utils.GenerateSetError(resPeerName, "name", err)
	}
	if err := d.Set("description", wireguard.Properties.Description); err != nil {
		return utils.GenerateSetError(resPeerName, "description", err)
	}
	if wireguard.Properties.Endpoint != nil {
		var endpoiontSlice []any
		endpointEntry, err := wireguard.Properties.Endpoint.ToMap()
		if err != nil {
			return err
		}
		endpoiontSlice = append(endpoiontSlice, endpointEntry)
		if err := d.Set("endpoint", endpoiontSlice); err != nil {
			return utils.GenerateSetError(resPeerName, "endpoint", err)
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
