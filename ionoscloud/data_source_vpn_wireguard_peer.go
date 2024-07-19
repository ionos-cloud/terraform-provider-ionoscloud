package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
)

func dataSourceVpnWireguardPeer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnWireguardPeerRead,

		Schema: map[string]*schema.Schema{
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
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the wireguard peer ID or name"))
	}

	var peer vpnSdk.WireguardPeerRead
	var err error
	if idOk {
		peer, _, err = client.GetWireguardPeerByID(ctx, gatewayID, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the wireguard peer with ID: %s, error: %w", id, err))
		}
	} else {
		var results []vpnSdk.WireguardPeerRead
		peers, _, err := client.ListWireguardPeers(ctx, gatewayID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching wireguard peers: %w", err))
		}
		for _, recordItem := range peers.Items {
			if len(results) == 1 {
				break
			}
			if strings.EqualFold(recordItem.Properties.Name, name) {
				results = append(results, recordItem)
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no vpn wireguard peer found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one vpn wireguard [eer] found with the specified name = %s", name))
		} else {
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
