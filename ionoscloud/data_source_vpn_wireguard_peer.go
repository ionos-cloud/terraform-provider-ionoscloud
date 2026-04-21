package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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

func dataSourceVpnWireguardPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("ID and name cannot be both specified at the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the WireGuard Peer ID or name"), nil)
	}

	var peer vpnSdk.WireguardPeerRead
	var apiResponse *shared.APIResponse
	var err error
	if idOk {
		id := idValue.(string)
		peer, apiResponse, err = client.GetWireguardPeerByID(ctx, gatewayID, id, location)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the WireGuard Peer with ID: %s, error: %w", id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
	} else {
		name := nameValue.(string)
		var results []vpnSdk.WireguardPeerRead
		peers, apiResponse, err := client.ListWireguardPeers(ctx, gatewayID, location)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching WireGuard Peers: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
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
			return diagutil.ToDiags(d, fmt.Errorf("no VPN WireGuard Peer found with the specified name = %s", name), nil)
		case len(results) > 1:
			return diagutil.ToDiags(d, fmt.Errorf("more than one VPN WireGuard Peer found with the specified name = %s", name), nil)
		default:
			peer = results[0]
		}
	}
	if err := d.Set("id", peer.Id); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
