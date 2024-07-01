package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceApiGatewayRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiGatewayRouteRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the API Gateway route.",
				Optional:    true,
			},
			"apigateway_id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the API Gateway.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway.",
				Optional:    true,
			},
			"websocket": {
				Type:        schema.TypeBool,
				Description: "This field enables or disables websocket support.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "This field specifies the protocol used by the ingress to route traffic to the backend service.",
				Computed:    true,
			},
			"paths": {
				Type:        schema.TypeList,
				Description: "The paths that the route should match.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"methods": {
				Type:        schema.TypeList,
				Description: "The HTTP methods that the route should match.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"upstreams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:        schema.TypeString,
							Description: "The target URL of the upstream.",
							Computed:    true,
						},
						"loadbalancer": {
							Type:        schema.TypeString,
							Description: "The load balancer algorithm.",
							Computed:    true,
						},
						"host": {
							Type:        schema.TypeString,
							Description: "The host of the upstream.",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port of the upstream.",
							Computed:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true, // Alex: Added because of a panic, please remove if wrong
							Description: "Weight with which to split traffic to the upstream.",
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

func dataSourceApiGatewayRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)
	gatewayId := d.Get("apigateway_id").(string)

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the ApiGateway route ID or name"))
	}

	var route apigateway.RouteRead
	var err error
	if idOk {
		route, _, err = client.GetRouteById(ctx, gatewayId, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the ApiGateway route with ID: %s, error: %w", idValue, err))
		}
	} else {
		routes, _, err := client.ListRoutes(ctx, gatewayId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching ApiGateways route: %w", err))
		}

		var results []apigateway.RouteRead
		for _, r := range *routes.Items {
			if r.Properties != nil && r.Properties.Name != nil && utils.NameMatches(*r.Properties.Name, name, partialMatch) {
				results = append(results, r)
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no ApiGateway route found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one ApiGateway route found with the specified name: %s", name))
		} else {
			route = results[0]
		}
	}

	if err = client.SetApiGatewayRouteRead(d, route); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
