package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceAPIGatewayRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGatewayRouteCreate,
		ReadContext:   resourceAPIGatewayRouteRead,
		UpdateContext: resourceAPIGatewayRouteUpdate,
		DeleteContext: resourceAPIGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGatewayRouteImport,
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
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, false)),
				},
			},
			"upstreams": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:             schema.TypeString,
							Description:      "The target URL of the upstream.",
							Optional:         true,
							Default:          "http",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"http", "https", "grpc", "grpcs"}, false)),
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
							Type:             schema.TypeString,
							Description:      "The load balancer algorithm.",
							Optional:         true,
							Default:          "roundrobin",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"roundrobin", "least_connections"}, false)),
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
	client := meta.(bundleclient.SdkBundle).APIGatewayClient

	createdRoute, apiResponse, err := client.CreateRoute(ctx, d)
	if err != nil {
		d.SetId("")
		return utils.ToDiags(d, fmt.Sprintf("error creating API Gateway Route: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	d.SetId(createdRoute.Id)
	log.Printf("[INFO] Created API Gateway Route: %s", d.Id())

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsAPIGatewayRouteAvailable)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error waiting for API Gateway Route to be ready: %s", err), nil)
	}

	return resourceAPIGatewayRouteRead(ctx, d, meta)
}

func resourceAPIGatewayRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient

	routeID := d.Id()
	gatewayID := d.Get("gateway_id").(string)

	route, apiResponse, err := client.GetRouteByID(ctx, gatewayID, routeID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}

		return utils.ToDiags(d, fmt.Sprintf("error while fetching API Gateway Route: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retreived API Gateway Route %s: %+v", d.Id(), route)
	if err = client.SetAPIGatewayRouteData(d, route); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceAPIGatewayRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient

	updatedRoute, apiResponse, err := client.UpdateRoute(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error updating API Gateway Route: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsAPIGatewayRouteAvailable)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error waiting for API Gateway Route to be ready: %s", err), nil)
	}

	if err = client.SetAPIGatewayRouteData(d, updatedRoute); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceAPIGatewayRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).APIGatewayClient
	gatewayID := d.Get("gateway_id").(string)
	routeID := d.Id()

	apiResponse, err := client.DeleteRoute(ctx, gatewayID, routeID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error deleting API Gateway Route: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsAPIGatewayRouteDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error waiting for API Gateway Route to be deleted: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	return nil
}

func resourceAPIGatewayRouteImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "expected ID in the format gateway_id:route_id", nil)
	}

	if err := d.Set("gateway_id", parts[0]); err != nil {
		return nil, utils.GenerateSetError(constant.APIGatewayRouteResource, "gateway_id", err)
	}
	d.SetId(parts[1])

	diags := resourceAPIGatewayRouteRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, utils.ToError(d, diags[0].Summary, nil)
	}

	return []*schema.ResourceData{d}, nil
}
