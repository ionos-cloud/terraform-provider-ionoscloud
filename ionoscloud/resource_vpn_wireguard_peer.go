package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
				//ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.AvailableLocations, false)),
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
	client := m.(services.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	peer, _, err := client.CreateWireguardGatewayPeers(ctx, d, gatewayID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("error creating WireGuard Peer: %w", err))
	}
	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return nil
}

func resourceVpnWireguardPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
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
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpnWireguardPeerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(services.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	_, _, err := client.UpdateWireguardPeer(ctx, gatewayID, d.Id(), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating WireGuard Peer: %w", err))
	}
	return nil
}

func resourceVpnWireguardPeerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(services.SdkBundle).VPNClient
	gatewayID := d.Get("gateway_id").(string)
	location := d.Get("location").(string)
	apiResponse, err := client.DeleteWireguardPeer(ctx, gatewayID, d.Id(), location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting WireGuard Peer: %w", err))
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsWireguardPeerDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleting %w", err))
	}

	log.Printf("[INFO] Successfully deleted WireGuard Peer: %s", d.Id())

	d.SetId("")
	return nil
}

func resourceVpnWireguardPeerImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(services.SdkBundle).VPNClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import format: %s, expecting the following format: location:gateway_id:id", d.Id())
	}
	location := parts[0]
	gatewayID := parts[1]
	peerID := parts[2]
	peer, _, err := client.GetWireguardPeerByID(ctx, gatewayID, peerID, location)
	if err != nil {
		return nil, err
	}
	if err := d.Set("gateway_id", gatewayID); err != nil {
		return nil, err
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}
	if err := vpn.SetWireguardPeerData(d, peer); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
