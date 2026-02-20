package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceVpnWireguardPeer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnWireguardPeerCreate,
		ReadContext:   resourceVpnWireguardPeerRead,
		UpdateContext: resourceVpnWireguardPeerUpdate,
		DeleteContext: resourceVpnWireguardPeerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpnWireguardPeerImport,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The ID of the WireGuard Peer that the peer will connect to.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the WireGuard Peer. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The human readable name of your WireGuard Gateway Peer.",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Human readable description of the WireGuard Gateway Peer.",
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"endpoint": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname or IPV4 address that the WireGuard Server will connect to.",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      51820,
							Description:  "Port that the WireGuard Server will connect to.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
				Description: "Endpoint configuration for the WireGuard Peer.",
			},
			"allowed_ips": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The subnet CIDRs that are allowed to connect to the WireGuard Gateway.",
				MinItems:    1,
				MaxItems:    20,
			},
			"public_key": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "WireGuard public key of the connecting peer",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the WireGuard Gateway",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVpnWireguardPeerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	peer, apiResponse, err := client.CreateWireguardGatewayPeers(ctx, d, gatewayID)
	if err != nil {
		d.SetId("")
		return utils.ToDiags(d, fmt.Sprintf("error creating WireGuard Peer: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		d.SetId("")
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func resourceVpnWireguardPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	location := d.Get("location").(string)
	peer, apiResponse, err := client.GetWireguardPeerByID(ctx, gatewayID, d.Id(), location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[DEBUG] cannot find peer by gatewayID %s and id %s", gatewayID, d.Id())
			d.SetId("")
			return nil
		}
	}
	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceVpnWireguardPeerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	_, apiResponse, err := client.UpdateWireguardPeer(ctx, gatewayID, d.Id(), d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error updating WireGuard Peer: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	return nil
}

func resourceVpnWireguardPeerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(bundleclient.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	location := d.Get("location").(string)
	apiResponse, err := client.DeleteWireguardPeer(ctx, gatewayID, d.Id(), location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error deleting WireGuard Peer: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsWireguardPeerDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("deleting %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	log.Printf("[INFO] Successfully deleted WireGuard Peer: %s", d.Id())

	d.SetId("")
	return nil
}

func resourceVpnWireguardPeerImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(bundleclient.SdkBundle).VPNClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return nil, utils.ToError(d, "invalid import format:, expecting the following format: location:gateway_id:id", nil)
	}
	location := parts[0]
	gatewayID := parts[1]
	peerID := parts[2]
	peer, apiResponse, err := client.GetWireguardPeerByID(ctx, gatewayID, peerID, location)
	if err != nil {
		return nil, utils.ToError(d, err.Error(), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := d.Set("gateway_id", gatewayID); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}
