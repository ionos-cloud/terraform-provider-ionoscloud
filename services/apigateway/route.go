package apigateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *Client) GetRouteById(ctx context.Context, gatewayId string, routeId string) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesFindById(ctx, gatewayId, routeId).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

func (c *Client) ListRoutes(ctx context.Context, gatewayId string) (apigateway.RouteReadList, *apigateway.APIResponse, error) {
	routes, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesGet(ctx, gatewayId).Execute()
	apiResponse.LogInfo()
	return routes, apiResponse, err
}

func (c *Client) SetApiGatewayRouteRead(d *schema.ResourceData, route apigateway.RouteRead) error {
	d.SetId(*route.Id)

	if route.Properties == nil {
		return fmt.Errorf("expected properties in the response for the ApiGateway route with ID: %s, but received 'nil' instead", *route.Id)
	}

	if route.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the ApiGateway route with ID: %s, but received 'nil' instead", *route.Id)

	}

	if route.Properties.Name != nil {
		if err := d.Set("name", *route.Properties.Name); err != nil {
			return err
		}
	}

	if route.Properties.Websocket != nil {
		if err := d.Set("websocket", *route.Properties.Websocket); err != nil {
			return err
		}
	}

	if route.Properties.Type != nil {
		if err := d.Set("type", *route.Properties.Type); err != nil {
			return err
		}
	}

	if route.Properties.Paths != nil {
		if err := d.Set("paths", *route.Properties.Paths); err != nil {
			return err
		}
	}

	if route.Properties.Methods != nil {
		if err := d.Set("methods", *route.Properties.Methods); err != nil {
			return err
		}
	}

	if route.Properties.Upstreams != nil {
		var upstreams []map[string]interface{}
		for _, upstream := range *route.Properties.Upstreams {
			upstreamData := map[string]interface{}{}

			utils.SetPropWithNilCheck(upstreamData, "scheme", upstream.Scheme)
			utils.SetPropWithNilCheck(upstreamData, "host", upstream.Host)
			utils.SetPropWithNilCheck(upstreamData, "port", upstream.Port)
			utils.SetPropWithNilCheck(upstreamData, "loadbalancer", upstream.Loadbalancer)
			utils.SetPropWithNilCheck(upstreamData, "weight", upstream.Weight)

			upstreams = append(upstreams, upstreamData)
		}

		if err := d.Set("upstreams", upstreams); err != nil {
			return fmt.Errorf("error setting upstreams for the ApiGateway route with ID: %s", *route.Id)
		}
	}

	return nil
}
