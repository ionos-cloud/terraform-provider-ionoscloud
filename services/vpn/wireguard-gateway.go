package vpn

import (
	"context"
	"fmt"
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var wireguardResourceName = "vpn wireguard gateway"

func (c *Client) CreateWireguardGateway(ctx context.Context, d *schema.ResourceData) (vpn.WireguardGatewayRead, utils.ApiResponseInfo, error) {
	request := setWireguardGWPostRequest(d)
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPost(ctx).WireguardGatewayCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}
func (c *Client) CreateWireguardGatewayPeers(ctx context.Context, d *schema.ResourceData, gatewayID string) (vpn.WireguardPeerRead, utils.ApiResponseInfo, error) {
	request, err := setWireguardPeersPostRequest(d)
	if err != nil {
		return vpn.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPost(ctx, gatewayID).WireguardPeerCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

func (c *Client) IsWireguardAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	wireguardId := d.Id()
	wireguard, _, err := c.GetWireguardGatewayByID(ctx, wireguardId)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

func (c *Client) IsWireguardPeerAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	wireguardId := d.Id()
	gatewayID := d.Get("gateway_id").(string)
	wireguard, _, err := c.GetWireguardPeerByID(ctx, gatewayID, wireguardId)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

func (c *Client) UpdateWireguardGateway(ctx context.Context, id string, d *schema.ResourceData) (vpn.WireguardGatewayRead, utils.ApiResponseInfo, error) {
	request := setWireguardGatewayPatchRequest(d)
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPut(ctx, id).WireguardGatewayEnsure(*request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

func (c *Client) UpdateWireguardPeer(ctx context.Context, gatewayID, id string, d *schema.ResourceData) (vpn.WireguardPeerRead, utils.ApiResponseInfo, error) {
	request, err := setWireguardPeerPatchRequest(d)
	if err != nil {
		return vpn.WireguardPeerRead{}, nil, fmt.Errorf("error decoding endpoint: %w", err)
	}
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersPut(ctx, gatewayID, id).WireguardPeerEnsure(request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

func (c *Client) DeleteWireguardGateway(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}
func (c *Client) DeleteWireguardPeer(ctx context.Context, gatewayID, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersDelete(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsWireguardGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) IsWireguardPeerDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Get("gateway_id").(string)
	_, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetWireguardGatewayByID returns a wireguard by its ID
func (c *Client) GetWireguardGatewayByID(ctx context.Context, id string) (vpn.WireguardGatewayRead, *shared.APIResponse, error) {
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// GetWireguardPeerByID returns a wireguard by its ID
func (c *Client) GetWireguardPeerByID(ctx context.Context, gatewayID, id string) (vpn.WireguardPeerRead, *shared.APIResponse, error) {
	wireguard, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersFindById(ctx, gatewayID, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// ListWireguardGateways returns a list of all wireguards
func (c *Client) ListWireguardGateways(ctx context.Context) (vpn.WireguardGatewayReadList, *shared.APIResponse, error) {
	wireguards, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

// ListWireguardPeers returns a list of all wireguards
func (c *Client) ListWireguardPeers(ctx context.Context, gatewayID string) (vpn.WireguardPeerReadList, *shared.APIResponse, error) {
	wireguards, apiResponse, err := c.sdkClient.WireguardPeersApi.WireguardgatewaysPeersGet(ctx, gatewayID).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

func (c *Client) IsWireguardGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cluster, _, err := c.GetWireguardGatewayByID(ctx, d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wierguard gateway state %s", cluster.Metadata.Status)
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

func (c *Client) IsWireguardPeerReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Get("gateway_id").(string)
	cluster, _, err := c.GetWireguardPeerByID(ctx, gatewayID, d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wierguard gateway state %s", cluster.Metadata.Status)
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

func setWireguardGWPostRequest(d *schema.ResourceData) *vpn.WireguardGatewayCreate {
	request := vpn.WireguardGatewayCreate{Properties: vpn.WireguardGateway{}}
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
	request.Properties.Connections = GetWireguardGwConnectionsData(d)

	return &request
}

func setWireguardPeersPostRequest(d *schema.ResourceData) (*vpn.WireguardPeerCreate, error) {
	request := vpn.WireguardPeerCreate{Properties: vpn.WireguardPeer{}}
	name := d.Get("name").(string)

	request.Properties.Name = name

	if value, ok := d.GetOk("description"); ok {
		valueStr := value.(string)
		request.Properties.Description = &valueStr
	}

	if value, ok := d.GetOk("endpoint"); ok {
		raw := value.(interface{})
		//vpn.WireguardEndpoint
		endpoint := vpn.NewWireguardEndpointWithDefaults()
		err := utils.DecodeInterfaceToStruct(raw, endpoint)
		if err != nil {
			return &request, err
		}
		request.Properties.Endpoint = endpoint
	}

	if value, ok := d.GetOk("allowed_ips"); ok {
		valueStr := value.([]string)
		request.Properties.AllowedIPs = valueStr
	}
	if value, ok := d.GetOk("public_key"); ok {
		valueStr := value.(string)
		request.Properties.PublicKey = valueStr
	}

	return &request, nil
}

func GetWireguardGwConnectionsData(d *schema.ResourceData) []vpn.Connection {
	connections := make([]vpn.Connection, 0)

	if connectionValues, ok := d.GetOk("connections"); ok {
		connectionsItf := connectionValues.([]any)
		if connectionsItf != nil {
			for vdcIndex := range connectionsItf {

				connection := vpn.Connection{}

				if datacenterId, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", vdcIndex)); ok {
					datacenterId := datacenterId.(string)
					connection.DatacenterId = datacenterId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = lanId
				}

				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.ipv4_cidr", vdcIndex)); ok {
					cidr := cidr.(string)
					connection.Ipv4CIDR = cidr
				}

				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.ipv6_cidr", vdcIndex)); ok {
					cidr := cidr.(string)
					connection.Ipv6CIDR = &cidr
				}

				connections = append(connections, connection)
			}
		}
	}

	return connections
}

func setWireguardGatewayPatchRequest(d *schema.ResourceData) *vpn.WireguardGatewayEnsure {
	request := vpn.WireguardGatewayEnsure{Properties: vpn.WireguardGateway{}}

	if d.HasChange("name") {
		_, val := d.GetChange("name")
		request.Properties.Name = val.(string)
	}
	if d.HasChange("gateway_ip") {
		_, val := d.GetChange("gateway_ip")
		request.Properties.GatewayIP = val.(string)
	}

	if d.HasChange("description") {
		_, val := d.GetChange("description")
		request.Properties.Description = ionoscloud.ToPtr(val.(string))
	}
	if d.HasChange("interface_ipv4_cidr") {
		_, val := d.GetChange("interface_ipv4_cidr")
		request.Properties.InterfaceIPv4CIDR = ionoscloud.ToPtr(val.(string))
	}
	if d.HasChange("interface_ipv6_cidr") {
		_, val := d.GetChange("interface_ipv6_cidr")
		request.Properties.InterfaceIPv6CIDR = ionoscloud.ToPtr(val.(string))
	}
	if d.HasChange("listen_port") {
		_, val := d.GetChange("listen_port")
		request.Properties.ListenPort = ionoscloud.ToPtr(int32(val.(int)))
	}
	if d.HasChange("connections") {
		request.Properties.Connections = GetWireguardGwConnectionsData(d)
	}

	return &request
}

func setWireguardPeerPatchRequest(d *schema.ResourceData) (vpn.WireguardPeerEnsure, error) {
	request := vpn.WireguardPeerEnsure{Properties: vpn.WireguardPeer{}}

	if d.HasChange("name") {
		_, val := d.GetChange("name")
		request.Properties.Name = val.(string)
	}

	if d.HasChange("description") {
		_, val := d.GetChange("description")
		request.Properties.Description = ionoscloud.ToPtr(val.(string))
	}
	if d.HasChange("endpoint") {
		_, val := d.GetChange("endpoint")
		raw := val.(interface{})
		endpoint := vpn.NewWireguardEndpointWithDefaults()
		err := utils.DecodeInterfaceToStruct(raw, endpoint)
		if err != nil {
			return request, err
		}
		request.Properties.Endpoint = endpoint
	}
	if d.HasChange("allowed_ips") {
		_, val := d.GetChange("allowed_ips")
		request.Properties.AllowedIPs = val.([]string)
	}
	if d.HasChange("public_key") {
		_, val := d.GetChange("public_key")
		request.Properties.PublicKey = val.(string)
	}

	return request, nil
}

func SetWireguardGWData(d *schema.ResourceData, wireguard vpn.WireguardGatewayRead) error {
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

	if err := d.Set("listen_port", wireguard.Properties.ListenPort); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "listenPort", err)
	}
	if err := d.Set("status", wireguard.Metadata.Status); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "status", err)
	}

	return nil
}

var resPeerName = "vpn wireguard peer"

func SetWireguardPeerData(d *schema.ResourceData, wireguard vpn.WireguardPeerRead) error {
	d.SetId(wireguard.Id)

	if err := d.Set("name", wireguard.Properties.Name); err != nil {
		return utils.GenerateSetError(resPeerName, "name", err)
	}
	if err := d.Set("description", wireguard.Properties.Description); err != nil {
		return utils.GenerateSetError(resPeerName, "description", err)
	}
	if err := d.Set("endpoint", wireguard.Properties.Endpoint); err != nil {
		return utils.GenerateSetError(resPeerName, "endpoint", err)
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
