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
	c.changeConfigURL(d.Get("location").(string))
	request := setWireguardGWPostRequest(d)
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPost(ctx).WireguardGatewayCreate(*request).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// IsWireguardAvailable checks if the wireguard is available
func (c *Client) IsWireguardAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	location := d.Get("location").(string)
	c.changeConfigURL(location)
	wireguardID := d.Id()
	wireguard, _, err := c.GetWireguardGatewayByID(ctx, wireguardID, location)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] wireguard status: %s", wireguard.Metadata.Status)
	return strings.EqualFold(wireguard.Metadata.Status, constant.Available), nil
}

// UpdateWireguardGateway updates a wireguard gateway
func (c *Client) UpdateWireguardGateway(ctx context.Context, id string, d *schema.ResourceData) (vpnSdk.WireguardGatewayRead, utils.ApiResponseInfo, error) {
	c.changeConfigURL(d.Get("location").(string))
	request := setWireguardGatewayPutRequest(d)
	wireguardResponse, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysPut(ctx, id).WireguardGatewayEnsure(*request).Execute()
	apiResponse.LogInfo()
	return wireguardResponse, apiResponse, err
}

// DeleteWireguardGateway deletes a wireguard gateway
func (c *Client) DeleteWireguardGateway(ctx context.Context, id, location string) (utils.ApiResponseInfo, error) {
	c.changeConfigURL(location)
	apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsWireguardGatewayDeleted checks if the wireguard gateway is deleted
func (c *Client) IsWireguardGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	c.changeConfigURL(d.Get("location").(string))
	_, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetWireguardGatewayByID returns a wireguard by its ID
func (c *Client) GetWireguardGatewayByID(ctx context.Context, id, location string) (vpnSdk.WireguardGatewayRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	wireguard, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return wireguard, apiResponse, err
}

// ListWireguardGateways returns a list of all wireguards
func (c *Client) ListWireguardGateways(ctx context.Context, location string) (vpnSdk.WireguardGatewayReadList, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	wireguards, apiResponse, err := c.sdkClient.WireguardGatewaysApi.WireguardgatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return wireguards, apiResponse, err
}

// IsWireguardGatewayReady checks if the wireguard gateway is ready
func (c *Client) IsWireguardGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	location := d.Get("location").(string)
	c.changeConfigURL(location)
	cluster, _, err := c.GetWireguardGatewayByID(ctx, d.Id(), location)
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
	if _, ok := d.GetOk("maintenance_window"); ok {
		request.Properties.MaintenanceWindow = GetMaintenanceWindowData(d)
	}
	if value, ok := d.GetOk("tier"); ok {
		valueStr := value.(string)
		request.Properties.Tier = &valueStr
	}

	request.Properties.Connections = getWireguardGwConnectionsData(d)

	return &request
}

func getWireguardGwConnectionsData(d *schema.ResourceData) []vpnSdk.Connection {
	connections := make([]vpnSdk.Connection, 0)

	if connectionValues, ok := d.GetOk("connections"); ok {
		connectionsItf := connectionValues.([]any)
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
	if _, ok := d.GetOk("maintenance_window"); ok {
		request.Properties.MaintenanceWindow = GetMaintenanceWindowData(d)
	}
	if v, ok := d.GetOk("tier"); ok {
		request.Properties.Tier = shared.ToPtr(v.(string))
	}
	return &request
}

// SetWireguardGWData sets the wireguard gateway data
func SetWireguardGWData(d *schema.ResourceData, wireguard vpnSdk.WireguardGatewayRead) error {
	d.SetId(wireguard.Id)

	if err := d.Set("name", wireguard.Properties.Name); err != nil {
		return utils.GenerateSetError(wireguardResourceName, "name", err)
	}
	// TODO -- Check if reading a GW with a nil description will lead to an error here.
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

	var connections []map[string]any // nolint: prealloc
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
	if wireguard.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", setWireguardMaintenanceWindowData(wireguard.Properties)); err != nil {
			return utils.GenerateSetError(wireguardResourceName, "maintenance_window", err)
		}
	}
	if wireguard.Properties.Tier != nil {
		if err := d.Set("tier", wireguard.Properties.Tier); err != nil {
			return utils.GenerateSetError(wireguardResourceName, "tier", err)
		}
	}

	return nil
}

func setWireguardMaintenanceWindowData(wireguardGateway vpnSdk.WireguardGateway) []interface{} {
	var maintenanceWindows []interface{}
	maintenanceWindow := map[string]interface{}{}
	utils.SetPropWithNilCheck(maintenanceWindow, "time", wireguardGateway.MaintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenanceWindow, "day_of_the_week", wireguardGateway.MaintenanceWindow.DayOfTheWeek)

	maintenanceWindows = append(maintenanceWindows, maintenanceWindow)
	return maintenanceWindows
}
