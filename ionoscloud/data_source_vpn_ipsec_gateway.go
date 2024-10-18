package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
)

func dataSourceVpnIPSecGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnIPSecGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPSec Gateway.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human readable name of your IPSec Gateway.",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The human-readable description of your IPSec Gateway.",
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the IPSec Gateway. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Required:    true,
			},
			"gateway_ip": {
				Type:        schema.TypeString,
				Description: "Public IP address to be assigned to the gateway.",
				Computed:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your gateway. Note: all connections must belong to the same datacenter.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your VPN Gateway to.",
							Computed:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your VPN Gateway to.",
							Computed:    true,
						},
						"ipv4_cidr": {
							Type:        schema.TypeString,
							Description: "Describes the private ipv4 subnet in your LAN that should be accessible by the VPN Gateway. Note: this should be the subnet already assigned to the LAN",
							Computed:    true,
						},
						"ipv6_cidr": {
							Type:        schema.TypeString,
							Description: "Describes the ipv6 subnet in your LAN that should be accessible by the VPN Gateway. Note: this should be the subnet already assigned to the LAN",
							Computed:    true,
						},
					},
				},
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IKE version that is permitted for the VPN tunnels.",
				Computed:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "a weekly 4 hour-long window, during which maintenance might occur",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"day_of_the_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tier": {
				Type:        schema.TypeString,
				Description: "Gateway performance options",
				Computed:    true,
			},
		},
	}
}

func dataSourceVpnIPSecGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the IPSec Gateway ID or name"))
	}

	var gateway vpnSdk.IPSecGatewayRead
	var err error
	if idOk {
		gateway, _, err = client.GetIPSecGatewayByID(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the IPSec Gateway with ID: %s, error: %w", id, err))
		}
	} else {
		var results []vpnSdk.IPSecGatewayRead
		gateways, _, err := client.ListIPSecGateway(ctx, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching IPSec Gateways: %w", err))
		}

		for _, recordItem := range gateways.Items {
			if len(results) == 1 {
				break
			}

			if strings.EqualFold(recordItem.Properties.Name, name) {
				results = append(results, recordItem)
			}
		}

		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no VPN IPSec Gateway found with the specified name = %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one VPN IPSec Gateway found with the specified name = %s", name))
		default:
			gateway = results[0]
		}
	}
	if err := d.Set("id", gateway.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := vpn.SetIPSecGatewayData(d, gateway); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
