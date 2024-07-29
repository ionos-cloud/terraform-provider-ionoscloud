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

func resourceAPIGatewayRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGatewayRouteCreate,
		ReadContext:   resourceAPIGatewayRouteRead,
		UpdateContext: resourceAPIGatewayRouteUpdate,
		DeleteContext: resourceAPIGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGatewayImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the API Gateway Route.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway Route.",
				Required:    true,
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Description: "The ID of the API Gateway that the route belongs to.",
				Required:    true,
			},
			"websocket": {
				Type:        schema.TypeBool,
				Description: "To enable websocket support.",
				Optional:    true,
				Default:     false,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "This field specifies the protocol used by the ingress to route traffic to the backend service.",
				Optional:    true,
				Default:     "http",
			},
			"paths": {
				Type:        schema.TypeList,
				Description: "The paths that the route should match.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"methods": {
				Type:        schema.TypeList,
				Description: "The HTTP methods that the route should match.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"upstreams": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:        schema.TypeString,
							Description: "The target URL of the upstream.",
							Optional:    true,
							Default:     "http",
						},
						"host": {
							Type:        schema.TypeString,
							Description: "The host of the upstream.",
							Required:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port of the upstream.",
							Optional:    true,
							Default:     80,
						},
						"loadbalancer": {
							Type:        schema.TypeString,
							Description: "The load balancer algorithm.",
							Optional:    true,
							Default:     "round_robin",
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Weight with which to split traffic to the upstream.",
							Optional:    true,
							Default:     100,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceAPIGatewayRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient

	createdRoute, _, err := client.CreateRoute(ctx, d)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating API Gateway Route: %w", err))
		return diags
	}

	d.SetId(*createdRoute.Id)
	log.Printf("[INFO] Created API Gateway Route: %s", d.Id())

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsAPIGatewayRouteAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error waiting for API Gateway Route to be ready: %w", err))
		return diags
	}

	return resourceAPIGatewayRouteRead(ctx, d, meta)
}

func resourceAPIGatewayRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient

	routeID := d.Id()
	gatewayID := d.Get("gateway_id").(string)

	route, apiResponse, err := client.GetRouteByID(ctx, gatewayID, routeID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		diags := diag.FromErr(fmt.Errorf("error while fetching API Gateway Route %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived API Gateway Route %s: %+v", d.Id(), route)
	if err = client.SetAPIGatewayRouteData(d, route); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAPIGatewayRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient

	updatedRoute, _, err := client.UpdateRoute(ctx, d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error updating API Gateway Route: %w", err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsAPIGatewayRouteAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error waiting for API Gateway Route to be ready: %w", err))
		return diags
	}

	if err = client.SetAPIGatewayRouteData(d, updatedRoute); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAPIGatewayRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ApiGatewayClient
	gatewayID := d.Get("gateway_id").(string)
	routeID := d.Id()

	apiResponse, err := client.DeleteRoute(ctx, gatewayID, routeID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error deleting API Gateway Route: %w", err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsAPIGatewayRouteDeleted)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error waiting for API Gateway Route to be deleted: %w", err))
		return diags
	}

	return nil
}

func resourceAPIGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diags := resourceApiGatewayRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}