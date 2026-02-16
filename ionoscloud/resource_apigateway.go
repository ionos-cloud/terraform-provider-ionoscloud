package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceAPIGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGatewayCreate,
		ReadContext:   resourceAPIGatewayRead,
		UpdateContext: resourceAPIGatewayUpdate,
		DeleteContext: resourceAPIGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGatewayImport,
		},
		Schema: map[string]*schema.Schema{
			// computed ID
			"id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the API Gateway.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway.",
				Required:    true,
			},
			"logs": {
				Type:        schema.TypeBool,
				Description: "Enable or disable logging. NOTE: Central Logging must be enabled through the Logging API to enable this feature.",
				Optional:    true,
				Default:     false,
			},
			"metrics": {
				Type:        schema.TypeBool,
				Description: "Enable or disable metrics.",
				Optional:    true,
				Default:     false,
			},
			"custom_domains": {
				Type:        schema.TypeList,
				Description: "Custom domains for the API Gateway.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The domain name.",
							Required:    true,
						},
						"certificate_id": {
							Type:             schema.TypeString,
							Description:      "The certificate ID for the domain.",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
						},
					},
				},
			},
			"public_endpoint": {
				Type:        schema.TypeString,
				Description: "The public endpoint of the API Gateway.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceAPIGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	logClient := meta.(bundleclient.SdkBundle).LoggingClient

	logs, ok := d.GetOk("logs")
	if ok && logs.(bool) {
		central, apiResponse, err := logClient.GetCentralLogging(ctx)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("error getting Central Logging: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if len(central.Items) == 0 {
			return utils.ToDiags(d, "central Logging is not enabled, please use Logging API to enable Central Logging", nil)
		}
		// will only be one item in the list, we just have to check if it is enabled
		if !central.Items[0].Properties.Enabled {
			return utils.ToDiags(d, "cannot create API Gateway with logs disabled, please use Logging API to enable Central Logging", nil)
		}
	}

	response, apiResponse, err := client.CreateAPIGateway(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error creating API Gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	gatewayID := response.Id
	d.SetId(gatewayID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsGatewayReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error checking status for API Gateway with ID %v: %s", gatewayID, err), nil)
	}

	return resourceAPIGatewayRead(ctx, d, meta)
}

func resourceAPIGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	logClient := meta.(bundleclient.SdkBundle).LoggingClient

	logs, ok := d.GetOk("logs")
	if ok && logs.(bool) {
		central, apiResponse, err := logClient.GetCentralLogging(ctx)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("error getting Central Logging: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if len(central.Items) == 0 {
			return utils.ToDiags(d, "central Logging is not enabled, please use Logging API to enable Central Logging", nil)
		}
		// will only be one item in the list, we just have to check if it is enabled
		if !central.Items[0].Properties.Enabled {
			return utils.ToDiags(d, "cannot create API Gateway with logs disabled, please use Logging API to enable Central Logging", nil)
		}
	}

	_, apiResponse, err := client.UpdateAPIGateway(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error updating API Gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsGatewayReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error checking status for API Gateway %s", err), nil)
	}

	return resourceAPIGatewayRead(ctx, d, meta)
}

func resourceAPIGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	gatewayID := d.Id()
	apiResponse, err := client.DeleteAPIGateway(ctx, gatewayID)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error deleting API Gateway with ID: %v, error: %s", gatewayID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsGatewayDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("deletion check failed for API Gateway with ID: %v, error: %s", gatewayID, err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func resourceAPIGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := resourceAPIGatewayImport(ctx, d, meta)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error reading API Gateway: %s", err), nil)
	}
	return nil
}

func resourceAPIGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	gatewayID := d.Id()
	gateway, resp, err := client.GetAPIGatewayByID(ctx, gatewayID)
	if err != nil {
		if resp.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("API Gateway does not exist, error: %s", err), &utils.DiagsOpts{StatusCode: resp.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("error importing API Gateway with ID: %v, error: %s", gatewayID, err), &utils.DiagsOpts{StatusCode: resp.StatusCode})
	}
	log.Printf("[INFO] Gateway found: %+v", gateway)

	if err := client.SetAPIGatewayData(d, gateway); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}
