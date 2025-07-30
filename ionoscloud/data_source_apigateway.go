package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceAPIGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the API Gateway.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway.",
				Optional:    true,
				Computed:    true,
			},
			"logs": {
				Type:        schema.TypeBool,
				Description: "This field enables or disables the collection and reporting of logs for observability of this instance.",
				Computed:    true,
			},
			"metrics": {
				Type:        schema.TypeBool,
				Description: "This field enables or disables the collection and reporting of metrics for observability of this instance.",
				Computed:    true,
			},
			"custom_domains": {
				Type:        schema.TypeList,
				Description: "The custom domain that the API Gateway instance should listen on.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The domain name of the distribution.",
							Computed:    true,
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Description: "The ID of the certificate to use for the distribution.",
							Computed:    true,
						},
					},
				},
			},
			"public_endpoint": {
				Type:        schema.TypeString,
				Description: "The public endpoint of the API Gateway.",
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using the name filter.",
				Default:     false,
				Optional:    true,
			},
		},
	}
}

func dataSourceAPIGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the API Gateway ID or name"))
	}

	var gateway apigateway.GatewayRead
	var err error
	if idOk {
		gateway, _, err = client.GetAPIGatewayByID(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the API Gateway with ID: %s, error: %w", idValue, err))
		}
	} else {
		gateways, _, err := client.ListAPIGateways(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching API Gateways: %w", err))
		}

		var results []apigateway.GatewayRead
		for _, gw := range gateways.Items {
			if utils.NameMatches(gw.Properties.Name, name, partialMatch) {
				results = append(results, gw)
			}
		}

		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no API Gateway found with the specified name: %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one API Gateway found with the specified name: %s", name))
		default:
			gateway = results[0]
		}
	}

	if err = client.SetAPIGatewayData(d, gateway); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
