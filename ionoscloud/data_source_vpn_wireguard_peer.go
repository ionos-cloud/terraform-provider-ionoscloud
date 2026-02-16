package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Computed: true,
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
	client := m.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	location := d.Get("location").(string)
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be both specified at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the WireGuard Peer ID or name", nil)
	}

	var peer vpnSdk.WireguardPeerRead
	var apiResponse *shared.APIResponse
	var err error
	if idOk {
		peer, apiResponse, err = client.GetWireguardPeerByID(ctx, gatewayID, id, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the WireGuard Peer with ID: %s, error: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []vpnSdk.WireguardPeerRead
		peers, apiResponse, err := client.ListWireguardPeers(ctx, gatewayID, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching WireGuard Peers: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
			return utils.ToDiags(d, fmt.Sprintf("no VPN WireGuard Peer found with the specified name = %s", name), nil)
		case len(results) > 1:
			return utils.ToDiags(d, fmt.Sprintf("more than one VPN WireGuard Peer found with the specified name = %s", name), nil)
		default:
			peer = results[0]
		}
	}
	if err := d.Set("id", peer.Id); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
