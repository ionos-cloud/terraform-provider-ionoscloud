package ionoscloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceVpnIPSecGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnIPSecGatewayCreate,
		ReadContext:   resourceVpnIPSecGatewayRead,
		UpdateContext: resourceVpnIPSecGatewayUpdate,
		DeleteContext: resourceVpnIPSecGatewayDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human readable name of your IPSecGateway.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The human-readable description of your IPSec Gateway.",
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of the IPSec Gateway. Supported locations: %s", strings.Join(vpn.AvailableLocations, ", ")),
				Optional:    true,
				ForceNew:    true,
			},
			"gateway_ip": {
				Type:             schema.TypeString,
				Description:      "Public IP address to be assigned to the gateway. Note: This must be an IP address in the same datacenter as the connections.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your gateway. Note: all connections must belong to the same datacenter.",
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:             schema.TypeString,
							Description:      "The datacenter to connect your VPN Gateway to.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your VPN Gateway to.",
							Required:    true,
						},
						"ipv4_cidr": {
							Type:             schema.TypeString,
							Description:      "A LAN IPv4 address in CIDR notation that will be assigned to the VPN Gateway. This will be the private gateway address for LAN clients to route traffic over the VPN Gateway, this should be within the subnet already assigned to the LAN.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
							Required:         true,
						},
						"ipv6_cidr": {
							Type:             schema.TypeString,
							Description:      "A LAN IPv6 address in CIDR notation that will be assigned to the VPN Gateway. This will be the private gateway address for LAN clients to route traffic over the VPN Gateway, this should be within the subnet already assigned to the LAN.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
							Optional:         true,
						},
					},
				},
			},
			"version": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The IKE version that is permitted for the VPN tunnels.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IKEv2"}, false)),
				Default:          "IKEv2",
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A weekly 4 hour-long window, during which maintenance might occur",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "Start of the maintenance window in UTC time.",
							Required:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "The name of the week day",
							Required:    true,
						},
					},
				},
			},
			"tier": {
				Type:        schema.TypeString,
				Description: "Gateway performance options. See the documentation for the available options",
				Default:     constant.DefaultTier,
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpnIPSecGatewayImport,
		},
	}
}

func resourceVpnIPSecGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient

	gateway, apiResponse, err := client.CreateIPSecGateway(ctx, d)
	if err != nil {
		return utils.ToDiags(d, err.Error(), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	d.SetId(gateway.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsIPSecGatewayReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("creating %s ", err), nil)
	}

	return resourceVpnIPSecGatewayRead(ctx, d, meta)
}

func resourceVpnIPSecGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient
	id := d.Id()
	location := d.Get("location").(string)

	gateway, apiResponse, err := client.GetIPSecGatewayByID(ctx, id, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		return utils.ToDiags(d, fmt.Sprintf("error while fetching IPSec Gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	return utils.ToDiags(d, vpn.SetIPSecGatewayData(d, gateway).Error(), nil)
}
func resourceVpnIPSecGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient

	gateway, apiResponse, err := client.UpdateIPSecGateway(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error updating IPSec Gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsIPSecGatewayReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("while waiting for IPSec Gateway to be ready: %s", err), nil)
	}

	return utils.ToDiags(d, vpn.SetIPSecGatewayData(d, gateway).Error(), nil)
}

func resourceVpnIPSecGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).VPNClient
	id := d.Id()
	location := d.Get("location").(string)

	apiResponse, err := client.DeleteIPSecGateway(ctx, id, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		return utils.ToDiags(d, fmt.Sprintf("error while deleting IPSec Gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	time.Sleep(5 * time.Second)
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsIPSecGatewayDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("while deleting IPSec Gateway: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	return nil
}

func resourceVpnIPSecGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	location := parts[0]
	id := parts[1]

	if err := d.Set("location", location); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	d.SetId(id)

	diags := resourceVpnIPSecGatewayRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, utils.ToError(d, diags[0].Summary, nil)
	}
	return []*schema.ResourceData{d}, nil
}
