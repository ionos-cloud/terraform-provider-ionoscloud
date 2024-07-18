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

func (c *Client) IsWireguardAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	wireguardId := d.Id()
	wireguard, _, err := c.GetWireguardGatewayByID(ctx, wireguardId)
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

func (c *Client) DeleteWireguardGateway(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsWireguardGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetWireguardGatewayByID returns a wireguard by its ID
func (c *Client) GetWireguardGatewayByID(ctx context.Context, id string) (vpn.WireguardGatewayRead, *shared.APIResponse, error) {
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// ListWireguardGateways returns a list of all wireguards
func (c *Client) ListWireguardGateways(ctx context.Context) (vpn.WireguardGatewayReadList, *shared.APIResponse, error) {
	wireguards, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

func (c *Client) IsWireguardGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cluster, _, err := c.GetWireguardGatewayByID(ctx, d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] dataplatform cluster state %s", cluster.Metadata.Status)
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
	if d.HasChange("private_key") {

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

	//if wireguard.Properties.GrafanaAddress != nil {
	//	if err := d.Set("grafana_address", *wireguard.Properties.GrafanaAddress); err != nil {
	//		return utils.GenerateSetError(wireguardResourceName, "grafana_address", err)
	//	}
	//}
	//
	//if wireguard.Properties.Logs != nil {
	//	logs := make([]interface{}, len(wireguard.Properties.Logs))
	//	for i, logElem := range wireguard.Properties.Logs {
	//		// Populate the logElem entry.
	//		logEntry := make(map[string]interface{})
	//		logEntry["source"] = *logElem.Source
	//		logEntry["tag"] = *logElem.Tag
	//		logEntry["protocol"] = *logElem.Protocol
	//		logEntry["public"] = *logElem.Public
	//
	//		// Logic for destinations
	//		destinations := make([]interface{}, len(logElem.Destinations))
	//		for i, destination := range logElem.Destinations {
	//			destinationEntry := make(map[string]interface{})
	//			destinationEntry["type"] = *destination.Type
	//			destinationEntry["retention_in_days"] = *destination.RetentionInDays
	//			destinations[i] = destinationEntry
	//		}
	//		logEntry["destinations"] = destinations
	//		logs[i] = logEntry
	//	}
	//	if err := d.Set("log", logs); err != nil {
	//		return utils.GenerateSetError(wireguardResourceName, "log", err)
	//	}
	//}

	return nil
}
