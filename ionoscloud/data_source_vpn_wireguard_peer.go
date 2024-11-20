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

func dataSourceVpnWireguardPeer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnWireguardPeerRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the WireGuard Peer. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Optional:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"allowed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnWireguardPeerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(services.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	location := d.Get("location").(string)
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the WireGuard Peer ID or name"))
	}

	var peer vpnSdk.WireguardPeerRead
	var err error
	if idOk {
		peer, _, err = client.GetWireguardPeerByID(ctx, gatewayID, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the WireGuard Peer with ID: %s, error: %w", id, err))
		}
	} else {
		var results []vpnSdk.WireguardPeerRead
		peers, _, err := client.ListWireguardPeers(ctx, gatewayID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching WireGuard Peers: %w", err))
		}
		for _, recordItem := range peers.Items {
			if len(results) == 1 {
				break
			}
			if strings.EqualFold(recordItem.Properties.Name, name) {
				results = append(results, recordItem)
			}
		}
		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no VPN WireGuard Peer found with the specified name = %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one VPN WireGuard Peer found with the specified name = %s", name))
		default:
			peer = results[0]
		}
	}
	if err := d.Set("id", peer.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
