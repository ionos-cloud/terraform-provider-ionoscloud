package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceVpnWireguardGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnWireguardGatewayCreate,
		ReadContext:   resourceVpnWireguardGatewayRead,
		UpdateContext: resourceVpnWireguardGatewayUpdate,
		DeleteContext: resourceVpnWireguardGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpnWireguardGatewayImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the IPSec Gateway. Supported locations: de/fra, de/txl, es/vit, gb/lhr, us/ewr, us/las, us/mci, fr/par",
				Required:    true,
				ForceNew:    true,
			},
			"connections": {
				MinItems: 1,
				MaxItems: 10,
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
						},
						"lan_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipv4_cidr": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
							Optional:         true,
						},
						"ipv6_cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"private_key": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "PrivateKey used for WireGuard Serve",
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PublicKey used for WireGuard Server. Received in response from API",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interface_ipv4_cidr": {
				Type: schema.TypeString,
				Description: `The IPV4 address (with CIDR mask) to be assigned to the WireGuard interface. 
							 __Note__: either interfaceIPv4CIDR or interfaceIPv6CIDR is __required__.`,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"interface_ipv6_cidr": {
				Type: schema.TypeString,
				Description: `The IPV6 address (with CIDR mask) to be assigned to the WireGuard interface.
							 __Note__: either interfaceIPv6CIDR or interfaceIPv4CIDR is __required__.`,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"listen_port": {
				Type:     schema.TypeInt,
				Default:  51820,
				Optional: true,
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

func resourceVpnWireguardGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient

	gateway, _, err := client.CreateWireguardGateway(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(gateway.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsWireguardGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating %w ", err))
	}
	return resourceVpnWireguardGatewayRead(ctx, d, meta)
}

func resourceVpnWireguardGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	location := d.Get("location").(string)
	wireguard, _, err := client.GetWireguardGatewayByID(ctx, d.Id(), location)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(vpn.SetWireguardGWData(d, wireguard))
}
func resourceVpnWireguardGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient

	wireguard, _, err := client.UpdateWireguardGateway(ctx, d.Id(), d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsWireguardGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating %w ", err))
	}

	return diag.FromErr(vpn.SetWireguardGWData(d, wireguard))
}

func resourceVpnWireguardGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).VPNClient
	location := d.Get("location").(string)
	apiResponse, err := client.DeleteWireguardGateway(ctx, d.Id(), location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting wireguard gateway %s: %w", d.Id(), err))
		return diags
	}
	//todo: for now we need to keep this because otherwise we get an internal server error on the first find after the delete
	// remove when no longer necessary
	time.Sleep(5 * time.Second)
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsWireguardGatewayDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("while waiting for the wireguard gateway to be deleted %s : %w", d.Id(), err))
	}

	log.Printf("[INFO] Successfully deleted Wireguard Gateway %s", d.Id())
	d.SetId("")

	return nil
}

func resourceVpnWireguardGatewayImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(services.SdkBundle).VPNClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format: %s, expecting the following format: location:id", d.Id())
	}
	location := parts[0]
	ID := parts[1]
	gateway, _, err := client.GetWireguardGatewayByID(ctx, ID, location)
	if err != nil {
		return nil, err
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}
	if err := vpn.SetWireguardGWData(d, gateway); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
