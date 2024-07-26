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

func dataSourceVpnIPSecTunnel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnIPSecTunnelRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPSec Gateway Tunnel.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human readable name of your IPSec Gateway Tunnel.",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The human readable description of your IPSec Gateway Tunnel.",
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the IPSec Gateway Tunnel. Supported locations: de/fra, de/txl, es/vit, gb/lhr, us/ewr, us/las, us/mci, fr/par",
				Required:    true,
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPSec Gateway that the tunnel belongs to.",
				Required:    true,
			},
			"remote_host": {
				Type:        schema.TypeString,
				Description: "The remote peer host fully qualified domain name or public IPV4 IP to connect to.",
				Computed:    true,
			},
			"auth": {
				Type:        schema.TypeList,
				Description: "Properties with all data needed to define IPSec Authentication.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:        schema.TypeString,
							Description: "The Authentication Method to use for IPSec Authentication.",
							Computed:    true,
						},
					},
				},
			},
			"ike": {
				Type:        schema.TypeList,
				Description: "Settings for the initial security exchange phase.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diffie_hellman_group": {
							Type:        schema.TypeString,
							Description: "The Diffie-Hellman Group to use for IPSec Encryption.",
							Computed:    true,
						},
						"encryption_algorithm": {
							Type:        schema.TypeString,
							Description: "The encryption algorithm to use for IPSec Encryption.",
							Computed:    true,
						},
						"integrity_algorithm": {
							Type:        schema.TypeString,
							Description: "The integrity algorithm to use for IPSec Encryption.",
							Computed:    true,
						},
						"lifetime": {
							Type:        schema.TypeInt,
							Description: "The phase lifetime in seconds.",
							Computed:    true,
						},
					},
				},
			},
			"esp": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Settings for the IPSec SA (ESP) phase.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diffie_hellman_group": {
							Type:        schema.TypeString,
							Description: "The Diffie-Hellman Group to use for IPSec Encryption.",
							Computed:    true,
						},
						"encryption_algorithm": {
							Type:        schema.TypeString,
							Description: "The encryption algorithm to use for IPSec Encryption.",
							Computed:    true,
						},
						"integrity_algorithm": {
							Type:        schema.TypeString,
							Description: "The integrity algorithm to use for IPSec Encryption.",
							Computed:    true,
						},
						"lifetime": {
							Type:        schema.TypeInt,
							Description: "The phase lifetime in seconds.",
							Computed:    true,
						},
					},
				},
			},
			"cloud_network_cidrs": {
				Type:        schema.TypeList,
				Description: `The network CIDRs on the "Left" side that are allowed to connect to the IPSec tunnel, i.e the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses.`,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"peer_network_cidrs": {
				Type:        schema.TypeList,
				Description: `The network CIDRs on the "Right" side that are allowed to connect to the IPSec tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses.`,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVpnIPSecTunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)
	gatewayID := d.Get("gateway_id").(string)
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the IPSec Gateway Tunnel ID or name"))
	}

	var tunnel vpnSdk.IPSecTunnelRead
	var err error
	if idOk {
		tunnel, _, err = client.GetIPSecTunnelByID(ctx, id, gatewayID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the IPSec Gateway Tunnel with ID: %s, error: %w", id, err))
		}
	} else {
		var results []vpnSdk.IPSecTunnelRead
		gateways, _, err := client.ListIPSecTunnel(ctx, gatewayID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching IPSec Gateway Tunnels: %w", err))
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
			return diag.FromErr(fmt.Errorf("no VPN IPSec Gateway Tunnel found with the specified name = %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one VPN IPSec Gateway Tunnel found with the specified name = %s", name))
		default:
			tunnel = results[0]
		}
	}
	if err := d.Set("id", tunnel.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := vpn.SetIPSecTunnelData(d, tunnel); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
