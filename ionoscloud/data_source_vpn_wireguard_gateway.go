package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
)

func dataSourceVpnWireguardGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnWireguardGatewayRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the WireGuard Gateway. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Required:    true,
			},
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the WireGuard Gateway",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_ipv4_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_ipv6_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"listen_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the WireGuard Gateway",
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

func dataSourceVpnWireguardGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	location := d.Get("location").(string)
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the WireGuard Gateway ID or name"))
	}

	var wireguardGw vpnSdk.WireguardGatewayRead
	var err error
	if idOk {
		wireguardGw, _, err = client.GetWireguardGatewayByID(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the WireGuard Gateway with ID: %s, error: %w", id, err))
		}
	} else {
		var results []vpnSdk.WireguardGatewayRead
		gateways, _, err := client.ListWireguardGateways(ctx, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching WireGuard Gateways: %w", err))
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
			return diag.FromErr(fmt.Errorf("no VPN WireGuard Gateway found with the specified name = %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one VPN WireGuard Gateway found with the specified name = %s", name))
		default:
			wireguardGw = results[0]
		}
	}
	if err := d.Set("id", wireguardGw.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := vpn.SetWireguardGWData(d, wireguardGw); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
