package ionoscloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVpnWireguardPeer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnWireguardPeerCreate,
		ReadContext:   resourceVpnWireguardPeerRead,
		//UpdateContext: resourceVpnWireguardPeerUpdate,
		//DeleteContext: resourceVpnWireguardPeerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The human readable name of your WireguardGateway Peer.",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Human readable description of the WireguardGateway Peer.",
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
				Description: "Endpoint configuration for the WireGuard peer.",
			},
			"allowed_ips": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The subnet CIDRs that are allowed to connect to the WireGuard Gateway.",
				MinItems:    1,
				MaxItems:    20,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.Any(validation.IsCIDR, validation.IsIPAddress),
				),
			},
			"public_key": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "WireGuard public key of the connecting peer",
				ValidateFunc: validation.StringIsNotEmpty, // Add more specific validation if needed
			},
		},
	}
}

func resourceVpnWireguardPeerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := m.(*ionoscloud.Client)
	//
	//// Extract fields from schema.ResourceData
	//name := d.Get("name").(string)
	//description := d.Get("description").(string)
	//publicKey := d.Get("public_key").(string)
	//// Assuming endpoint and allowed_ips need to be handled similarly
	//
	//// Construct request object (pseudo-code, adjust according to actual API client)
	//createRequest := ionoscloud.WireguardPeerCreateRequest{
	//	Name:        name,
	//	Description: description,
	//	PublicKey:   publicKey,
	//	// Add other fields as necessary
	//}
	//
	//// Send request to API
	//peer, _, err := client.WireguardPeers.Create(ctx, createRequest)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//// Set the resource ID
	//d.SetId(peer.ID)

	return resourceVpnWireguardPeerRead(ctx, d, m)
}

func resourceVpnWireguardPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(services.SdkBundle).VPNClient

	// Use the ID to read the resource
	//peer, _, err := client.WireguardPeers.GetByID(ctx, d.Id())
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	// Update the state
	return nil
}
