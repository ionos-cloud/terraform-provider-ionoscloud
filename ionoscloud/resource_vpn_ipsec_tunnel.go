package ionoscloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceVpnIPSecTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnIPSecTunnelCreate,
		ReadContext:   resourceVpnIPSecTunnelRead,
		UpdateContext: resourceVpnIPSecTunnelUpdate,
		DeleteContext: resourceVpnIPSecTunnelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of your IPSec Gateway Tunnel.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The human-readable description of your IPSec Gateway Tunnel.",
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the IPSec Gateway Tunnel. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Optional:    true,
				ForceNew:    true,
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPSec Gateway that the tunnel belongs to.",
				Required:    true,
			},
			"remote_host": {
				Type:        schema.TypeString,
				Description: "The remote peer host fully qualified domain name or public IPV4 IP to connect to.",
				Required:    true,
			},
			"auth": {
				Type:        schema.TypeList,
				Description: "Properties with all data needed to define IPSec Authentication.",
				Required:    true,
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:             schema.TypeString,
							Description:      "The Authentication Method to use for IPSec Authentication.",
							Optional:         true,
							Default:          "PSK",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PSK"}, false)),
						},
						"psk_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The Pre-Shared Key to use for IPSec Authentication. Note: Required if method is PSK.",
						},
					},
				},
			},
			"ike": {
				Type:        schema.TypeList,
				Description: "Settings for the initial security exchange phase.",
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diffie_hellman_group": {
							Type:             schema.TypeString,
							Description:      "The Diffie-Hellman Group to use for IPSec Encryption.",
							Optional:         true,
							Default:          "16-MODP4096",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelDiffieHellmanGroups, false)),
						},
						"encryption_algorithm": {
							Type:             schema.TypeString,
							Description:      "The encryption algorithm to use for IPSec Encryption.",
							Optional:         true,
							Default:          "AES256",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelEncryptionAlgorithms, false)),
						},
						"integrity_algorithm": {
							Type:             schema.TypeString,
							Description:      "The integrity algorithm to use for IPSec Encryption.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelIntegrityAlgorithms, false)),
							Optional:         true,
							Default:          "SHA256",
						},
						"lifetime": {
							Type:             schema.TypeInt,
							Description:      "The phase lifetime in seconds.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(3600, 86400)),
							Optional:         true,
							Default:          86400,
						},
					},
				},
			},
			"esp": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Settings for the IPSec SA (ESP) phase.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diffie_hellman_group": {
							Type:             schema.TypeString,
							Description:      "The Diffie-Hellman Group to use for IPSec Encryption.",
							Optional:         true,
							Default:          "16-MODP4096",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelDiffieHellmanGroups, false)),
						},
						"encryption_algorithm": {
							Type:             schema.TypeString,
							Description:      "The encryption algorithm to use for IPSec Encryption.",
							Optional:         true,
							Default:          "AES256",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelEncryptionAlgorithms, false)),
						},
						"integrity_algorithm": {
							Type:             schema.TypeString,
							Description:      "The integrity algorithm to use for IPSec Encryption.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(vpn.IPSecTunnelIntegrityAlgorithms, false)),
							Optional:         true,
							Default:          "SHA256",
						},
						"lifetime": {
							Type:             schema.TypeInt,
							Description:      "The phase lifetime in seconds.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(3600, 86400)),
							Optional:         true,
							Default:          86400,
						},
					},
				},
			},
			"cloud_network_cidrs": {
				Type:        schema.TypeList,
				Description: `The network CIDRs on the "Left" side that are allowed to connect to the IPSec tunnel, i.e. the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses.`,
				Required:    true,
				MinItems:    1,
				MaxItems:    20,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
				},
			},
			"peer_network_cidrs": {
				Type:        schema.TypeList,
				Description: `The network CIDRs on the "Right" side that are allowed to connect to the IPSec tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses.`,
				Required:    true,
				MinItems:    1,
				MaxItems:    20,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpnIPSecTunnelImport,
		},
	}
}

func resourceVpnIPSecTunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	pskKey := ""

	if d.Get("auth.0.method").(string) == "PSK" && d.Get("auth.0.psk_key").(string) == "" {
		return diag.FromErr(fmt.Errorf("psk_key is required when auth method is PSK"))
	}

	if d.Get("auth.0.method").(string) != "PSK" && d.Get("auth.0.psk_key").(string) != "" {
		return diag.FromErr(fmt.Errorf("psk_key is only required when auth method is PSK"))
	}

	if v, ok := d.GetOk("auth.0.psk_key"); ok {
		pskKey = v.(string)
	}

	tunnel, _, err := client.CreateIPSecTunnel(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tunnel.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsIPSecTunnelReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating %w ", err))
	}

	auth := d.Get("auth").([]interface{})
	auth[0].(map[string]interface{})["psk_key"] = pskKey
	if err = d.Set("auth", auth); err != nil {
		return diag.FromErr(err)
	}

	diags := resourceVpnIPSecTunnelRead(ctx, d, meta)
	if diags.HasError() {
		return diags
	}

	return nil
}

func resourceVpnIPSecTunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	id := d.Id()
	location := d.Get("location").(string)
	gatewayID := d.Get("gateway_id").(string)
	pskKey := d.Get("auth.0.psk_key").(string)

	tunnel, apiResponse, err := client.GetIPSecTunnelByID(ctx, id, gatewayID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		diags := diag.FromErr(fmt.Errorf("error while fetching IPSec Gateway Tunnel %s: %w", d.Id(), err))
		return diags
	}

	err = vpn.SetIPSecTunnelData(d, tunnel)
	if err != nil {
		return diag.FromErr(err)
	}

	auth := d.Get("auth").([]interface{})
	auth[0].(map[string]interface{})["psk_key"] = pskKey
	if err = d.Set("auth", auth); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
func resourceVpnIPSecTunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	pskKey := ""

	if d.Get("auth.0.method").(string) == "PSK" && d.Get("auth.0.psk_key").(string) == "" {
		return diag.FromErr(fmt.Errorf("psk_key is required when auth method is PSK"))
	}

	if d.Get("auth.0.method").(string) != "PSK" && d.Get("auth.0.psk_key").(string) != "" {
		return diag.FromErr(fmt.Errorf("psk_key is only required when auth method is PSK"))
	}

	if v, ok := d.GetOk("auth.0.psk_key"); ok {
		pskKey = v.(string)
	}

	tunnel, _, err := client.UpdateIPSecTunnel(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating IPSec Gateway Tunnel %s: %w", d.Id(), err))
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsIPSecTunnelReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("while waiting for IPSec Gateway Tunnel to be ready: %w", err))
	}

	err = vpn.SetIPSecTunnelData(d, tunnel)
	if err != nil {
		return diag.FromErr(err)
	}

	auth := d.Get("auth").([]interface{})
	auth[0].(map[string]interface{})["psk_key"] = pskKey
	if err = d.Set("auth", auth); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpnIPSecTunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	id := d.Id()
	location := d.Get("location").(string)
	gatewayID := d.Get("gateway_id").(string)

	apiResponse, err := client.DeleteIPSecTunnel(ctx, id, gatewayID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		diags := diag.FromErr(fmt.Errorf("error while deleting IPSec Gateway Tunnel %s: %w", d.Id(), err))
		return diags
	}

	time.Sleep(5 * time.Second)
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsIPSecTunnelDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("while deleting IPSec Gateway Tunnel %s : %w", d.Id(), err))
	}

	return nil
}

func resourceVpnIPSecTunnelImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	location := parts[0]
	gatewayID := parts[1]
	id := parts[2]

	if err := d.Set("location", location); err != nil {
		return nil, err
	}
	if err := d.Set("gateway_id", gatewayID); err != nil {
		return nil, err
	}
	d.SetId(id)

	diags := resourceVpnIPSecTunnelRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}
