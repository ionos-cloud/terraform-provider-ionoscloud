package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceApiGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiGatewayCreate,
		ReadContext:   resourceApiGatewayRead,
		UpdateContext: resourceApiGatewayUpdate,
		DeleteContext: resourceApiGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApiGatewayImport,
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
							Type:        schema.TypeString,
							Description: "The certificate ID for the domain.",
							Required:    true,
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

func resourceApiGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	logClient := meta.(services.SdkBundle).LoggingClient

	logs, ok := d.GetOk("logs")
	if ok && logs.(bool) {
		central, _, err := logClient.GetCentralLogging(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting Central Logging: %w", err))
		}

		if central.Properties == nil || central.Properties.Enabled == nil || !*central.Properties.Enabled {
			return diag.FromErr(fmt.Errorf("cannot create API Gateway with logs enabled, please use Logging API to enable Central Logging"))
		}
	}

	response, _, err := client.CreateAPIGateway(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating API Gateway: %w", err))
	}
	gatewayID := *response.Id
	d.SetId(gatewayID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for API Gateway with ID %v: %w", gatewayID, err))
	}
	if err := client.SetApiGatewayData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApiGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient

	response, _, err := client.UpdateAPIGateway(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating API Gateway: %w", err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for API Gateway %w", err))
	}
	if err := client.SetApiGatewayData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApiGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Id()
	_, err := client.DeleteApiGateway(ctx, gatewayID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsGatewayDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	return nil
}

func resourceApiGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := resourceApiGatewayImport(ctx, d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading API Gateway with ID: %v, error: %w", d.Id(), err))
	}
	return nil
}

func resourceApiGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Id()
	gateway, resp, err := client.GetApiGatewayById(ctx, gatewayID)
	if err != nil {
		if resp.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("API Gateway does not exist, error: %w", err)
		}
		return nil, fmt.Errorf("error importing API Gateway with ID: %v, error: %w", gatewayID, err)
	}
	log.Printf("[INFO] Gateway found: %+v", gateway)

	if err := client.SetApiGatewayData(d, gateway); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
