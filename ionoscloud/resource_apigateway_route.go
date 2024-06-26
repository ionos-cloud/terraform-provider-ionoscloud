package ionoscloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGatewayRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiGatewayRouteCreate,
		ReadContext:   resourceApiGatewayRouteRead,
		UpdateContext: resourceApiGatewayRouteUpdate,
		DeleteContext: resourceApiGatewayRouteDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway route",
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
						"destination": {},
					},
				},
			},
		},
	}
}

func resourceApiGatewayRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
