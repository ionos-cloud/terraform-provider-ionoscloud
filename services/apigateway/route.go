package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func (c *Client) CreateRoute(ctx context.Context, d *schema.ResourceData) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	request := setRoutePostRequest(d)
	gatewayId := d.Get("gateway_id").(string)

	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesPost(ctx, gatewayId).RouteCreate(*request).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

func (c *Client) UpdateRoute(ctx context.Context, d *schema.ResourceData) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	request := setRoutePutRequest(d)
	gatewayId := d.Get("gateway_id").(string)
	routeId := d.Id()

	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesPut(ctx, gatewayId, routeId).RouteEnsure(*request).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

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

func (c *Client) DeleteRoute(ctx context.Context, gatewayId string, routeId string) (*apigateway.APIResponse, error) {
	apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesDelete(ctx, gatewayId, routeId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsApiGatewayRouteAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	routeId := d.Id()
	gatewayId := d.Get("gateway_id").(string)

	route, _, err := c.GetRouteById(ctx, gatewayId, routeId)
	if err != nil {
		return false, err
	}

	if route.Metadata == nil || route.Metadata.Status == nil {
		return false, fmt.Errorf("expected metadata, got empty for ApiGateway Route with ID: %s", route)
	}
	log.Printf("[DEBUG] ApiGateway Route status: %s", *route.Metadata.Status)
	return strings.EqualFold(*route.Metadata.Status, constant.Available), nil
}

func (c *Client) IsApiGatewayRouteDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	routeId := d.Id()
	gatewayId := d.Get("gateway_id").(string)

	_, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesFindById(ctx, gatewayId, routeId).Execute()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) SetApiGatewayRouteData(d *schema.ResourceData, route apigateway.RouteRead) error {
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

func setRoutePostRequest(d *schema.ResourceData) *apigateway.RouteCreate {
	return apigateway.NewRouteCreate(*setRouteConfig(d))
}

func setRoutePutRequest(d *schema.ResourceData) *apigateway.RouteEnsure {
	routeId := d.Id()
	route := setRouteConfig(d)

	return apigateway.NewRouteEnsure(routeId, *route)
}

func setRouteConfig(d *schema.ResourceData) *apigateway.Route {
	routeName := d.Get("name").(string)
	routeType := d.Get("type").(string)
	pathsRaw := d.Get("paths").([]interface{})
	methodsRaw := d.Get("methods").([]interface{})
	websocket := d.Get("websocket").(bool)
	upstreamsRaw := d.Get("upstreams").([]interface{})

	var paths []string
	for _, path := range pathsRaw {
		paths = append(paths, path.(string))
	}

	var methods []string
	for _, method := range methodsRaw {
		methods = append(methods, method.(string))
	}

	var upstreams []apigateway.RouteUpstreams
	for _, upstream := range upstreamsRaw {
		upstreamData := upstream.(map[string]interface{})
		scheme := upstreamData["scheme"].(string)
		host := upstreamData["host"].(string)
		port := int32(upstreamData["port"].(int))
		loadbalancer := upstreamData["loadbalancer"].(string)
		weight := int32(upstreamData["weight"].(int))

		upstreamObj := apigateway.RouteUpstreams{
			Scheme:       &scheme,
			Host:         &host,
			Port:         &port,
			Loadbalancer: &loadbalancer,
			Weight:       &weight,
		}

		upstreams = append(upstreams, upstreamObj)
	}

	route := apigateway.NewRoute(routeName, routeType, paths, methods, upstreams)
	route.SetWebsocket(websocket)

	return route
}
