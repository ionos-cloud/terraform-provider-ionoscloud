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
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway.",
				Required:    true,
			},
			"logs": {
				Type:        schema.TypeBool,
				Description: "Enable or disable logging.",
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceApiGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient

	gateway := services.GatewayCreate{
		Properties: &services.Gateway{
			Name:    utils.String(d.Get("name").(string)),
			Logs:    utils.Bool(d.Get("logs").(bool)),
			Metrics: utils.Bool(d.Get("metrics").(bool)),
		},
	}

	if v, ok := d.GetOk("custom_domains"); ok {
		domains := v.([]interface{})
		customDomains := make([]services.GatewayCustomDomains, len(domains))
		for i, domain := range domains {
			domainMap := domain.(map[string]interface{})
			customDomains[i] = services.GatewayCustomDomains{
				Name:          utils.String(domainMap["name"].(string)),
				CertificateId: utils.String(domainMap["certificate_id"].(string)),
			}
		}
		gateway.Properties.CustomDomains = &customDomains
	}

	response, _, err := client.CreateApiGateway(ctx, gateway)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating API Gateway: %w", err))
	}
	gatewayID := *response.Id
	d.SetId(gatewayID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsApiGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	if err := client.SetApiGatewayData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApiGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Id()
	gateway, _, err := client.GetApiGatewayById(ctx, gatewayID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	log.Printf("[INFO] Successfully retrieved API Gateway with ID: %v, gateway info: %+v", gatewayID, gateway)
	if err := client.SetApiGatewayData(d, gateway); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApiGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Id()

	gateway := services.GatewayUpdate{
		Properties: &services.Gateway{
			Name:    utils.String(d.Get("name").(string)),
			Logs:    utils.Bool(d.Get("logs").(bool)),
			Metrics: utils.Bool(d.Get("metrics").(bool)),
		},
	}

	if v, ok := d.GetOk("custom_domains"); ok {
		domains := v.([]interface{})
		customDomains := make([]services.GatewayCustomDomains, len(domains))
		for i, domain := range domains {
			domainMap := domain.(map[string]interface{})
			customDomains[i] = services.GatewayCustomDomains{
				Name:          utils.String(domainMap["name"].(string)),
				CertificateId: utils.String(domainMap["certificate_id"].(string)),
			}
		}
		gateway.Properties.CustomDomains = &customDomains
	}

	response, _, err := client.UpdateApiGateway(ctx, gatewayID, gateway)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsApiGatewayReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for API Gateway with ID: %v after update, error: %w", gatewayID, err))
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
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for API Gateway with ID: %v, error: %w", gatewayID, err))
	}
	return nil
}

func resourceApiGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Id()
	gateway, resp, err := client.GetApiGatewayById(ctx, gatewayID)
	if err != nil {
		if resp.HttpNotFound() {
			return nil, fmt.Errorf("API Gateway does not exist, error: %w", err)
		}
		return nil, fmt.Errorf("error importing API Gateway with ID: %v, error: %w", gatewayID, err)
	}
	if err := client.SetApiGatewayData(d, gateway); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
