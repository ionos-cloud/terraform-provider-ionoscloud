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

const ipsecGatewayResourceName = "VPN IPSec Gateway"

// CreateIPSecGateway creates a new VPN IPSec Gateway
func (c *Client) CreateIPSecGateway(ctx context.Context, d *schema.ResourceData) (vpn.IPSecGatewayRead, *shared.APIResponse, error) {
	c.changeConfigURL(d.Get("location").(string))

	request := setIPSecGatewayCreateRequest(d)
	gateway, apiResponse, err := c.sdkClient.IPSecGatewaysApi.IpsecgatewaysPost(ctx).IPSecGatewayCreate(request).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

// GetIPSecGatewayByID retrieves a VPN IPSec Gateway by its ID and location
func (c *Client) GetIPSecGatewayByID(ctx context.Context, id string, location string) (vpn.IPSecGatewayRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)

	gateway, apiResponse, err := c.sdkClient.IPSecGatewaysApi.IpsecgatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

// ListIPSecGateway retrieves all VPN IPSec Gateways from a given location
func (c *Client) ListIPSecGateway(ctx context.Context, location string) (vpn.IPSecGatewayReadList, *shared.APIResponse, error) {
	c.changeConfigURL(location)

	gateways, apiResponse, err := c.sdkClient.IPSecGatewaysApi.IpsecgatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return gateways, apiResponse, err
}

// DeleteIPSecGateway deletes a VPN IPSec Gateway by its ID and location
func (c *Client) DeleteIPSecGateway(ctx context.Context, id string, location string) (utils.ApiResponseInfo, error) {
	c.changeConfigURL(location)

	apiResponse, err := c.sdkClient.IPSecGatewaysApi.IpsecgatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateIPSecGateway updates a VPN IPSec Gateway
func (c *Client) UpdateIPSecGateway(ctx context.Context, d *schema.ResourceData) (vpn.IPSecGatewayRead, *shared.APIResponse, error) {
	c.changeConfigURL(d.Get("location").(string))

	request := setIPSecGatewayPutRequest(d)
	gateway, apiResponse, err := c.sdkClient.IPSecGatewaysApi.IpsecgatewaysPut(ctx, d.Id()).IPSecGatewayEnsure(request).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

// IsIPSecGatewayReady checks if a VPN IPSec Gateway is ready to use
func (c *Client) IsIPSecGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	id := d.Id()
	location := d.Get("location").(string)

	gateway, _, err := c.GetIPSecGatewayByID(ctx, id, location)
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] VPN IPSec Gateway state %s", gateway.Metadata.Status)

	return strings.EqualFold(gateway.Metadata.Status, constant.Available), nil
}

// IsIPSecGatewayDeleted checks if a VPN IPSec Gateway has been deleted
func (c *Client) IsIPSecGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	id := d.Id()
	location := d.Get("location").(string)

	_, apiResponse, err := c.GetIPSecGatewayByID(ctx, id, location)
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// SetIPSecGatewayData sets the VPN IPSec Gateway data to the Terraform schema
func SetIPSecGatewayData(d *schema.ResourceData, gateway vpn.IPSecGatewayRead) error {
	d.SetId(gateway.Id)

	if err := d.Set("name", gateway.Properties.Name); err != nil {
		return utils.GenerateSetError(ipsecGatewayResourceName, "name", err)
	}

	if gateway.Properties.Description != nil {
		if err := d.Set("description", *gateway.Properties.Description); err != nil {
			return utils.GenerateSetError(ipsecGatewayResourceName, "description", err)
		}
	}

	if err := d.Set("version", gateway.Properties.Version); err != nil {
		return utils.GenerateSetError(ipsecGatewayResourceName, "version", err)
	}
	if err := d.Set("gateway_ip", gateway.Properties.GatewayIP); err != nil {
		return utils.GenerateSetError(ipsecGatewayResourceName, "gateway_ip", err)
	}

	connections := make([]map[string]interface{}, len(gateway.Properties.Connections))
	for i, connection := range gateway.Properties.Connections {
		connectionData := map[string]interface{}{
			"datacenter_id": connection.DatacenterId,
			"lan_id":        connection.LanId,
			"ipv4_cidr":     connection.Ipv4CIDR,
			"ipv6_cidr":     connection.Ipv6CIDR,
		}
		connections[i] = connectionData
	}

	if err := d.Set("connections", connections); err != nil {
		return utils.GenerateSetError(ipsecGatewayResourceName, "connections", err)
	}

	if gateway.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", setIPSecMaintenanceWindowData(gateway.Properties)); err != nil {
			return utils.GenerateSetError(ipsecGatewayResourceName, "maintenance_window", err)
		}
	}
	if gateway.Properties.Tier != nil {
		if err := d.Set("tier", gateway.Properties.Tier); err != nil {
			return utils.GenerateSetError(ipsecGatewayResourceName, "tier", err)
		}
	}

	return nil
}

func setIPSecGatewayCreateRequest(d *schema.ResourceData) vpn.IPSecGatewayCreate {
	properties := setIPSecGatewayProperties(d)

	return vpn.IPSecGatewayCreate{Properties: properties}
}

func setIPSecGatewayPutRequest(d *schema.ResourceData) vpn.IPSecGatewayEnsure {
	properties := setIPSecGatewayProperties(d)

	return vpn.IPSecGatewayEnsure{Id: d.Id(), Properties: properties}
}

func setIPSecGatewayProperties(d *schema.ResourceData) vpn.IPSecGateway {
	properties := vpn.IPSecGateway{}

	properties.Name = d.Get("name").(string)
	properties.GatewayIP = d.Get("gateway_ip").(string)

	if v, ok := d.GetOk("description"); ok {
		properties.Description = shared.ToPtr(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		properties.Version = shared.ToPtr(v.(string))
	}

	connections := make([]vpn.Connection, len(d.Get("connections").([]interface{})))
	for i := range d.Get("connections").([]interface{}) {
		connections[i] = setConnectionData(d, i)
	}
	properties.Connections = connections

	if _, ok := d.GetOk("maintenance_window"); ok {
		properties.MaintenanceWindow = GetMaintenanceWindowData(d)
	}
	if v, ok := d.GetOk("tier"); ok {
		properties.Tier = shared.ToPtr(v.(string))
	}

	return properties
}

func setConnectionData(d *schema.ResourceData, index int) vpn.Connection {
	conn := vpn.Connection{}

	conn.DatacenterId = d.Get(fmt.Sprintf("connections.%d.datacenter_id", index)).(string)
	conn.LanId = d.Get(fmt.Sprintf("connections.%d.lan_id", index)).(string)
	conn.Ipv4CIDR = d.Get(fmt.Sprintf("connections.%d.ipv4_cidr", index)).(string)
	if v, ok := d.GetOk(fmt.Sprintf("connections.%d.ipv6_cidr", index)); ok {
		conn.Ipv6CIDR = shared.ToPtr(v.(string))
	}

	return conn
}

func setIPSecMaintenanceWindowData(ipSecGateway vpn.IPSecGateway) []interface{} {
	var maintenanceWindows []interface{}
	maintenanceWindow := map[string]interface{}{}
	utils.SetPropWithNilCheck(maintenanceWindow, "time", ipSecGateway.MaintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenanceWindow, "day_of_the_week", ipSecGateway.MaintenanceWindow.DayOfTheWeek)

	maintenanceWindows = append(maintenanceWindows, maintenanceWindow)
	return maintenanceWindows
}
