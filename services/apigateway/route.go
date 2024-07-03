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

// CreateRoute sends a POST request with the given data to the API to create a route.
func (c *Client) CreateRoute(ctx context.Context, d *schema.ResourceData) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	request := setRoutePostRequest(d)
	gatewayID := d.Get("gateway_id").(string)

	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesPost(ctx, gatewayID).RouteCreate(*request).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

// UpdateRoute sends a PUT request with the given data to the API to update the route.
func (c *Client) UpdateRoute(ctx context.Context, d *schema.ResourceData) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	request := setRoutePutRequest(d)
	gatewayID := d.Get("gateway_id").(string)
	routeID := d.Id()

	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesPut(ctx, gatewayID, routeID).RouteEnsure(*request).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

// GetRouteByID sends a GET request to the API to retrieve a route by ID from a given gateway.
func (c *Client) GetRouteByID(ctx context.Context, gatewayID string, routeID string) (apigateway.RouteRead, *apigateway.APIResponse, error) {
	route, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesFindById(ctx, gatewayID, routeID).Execute()
	apiResponse.LogInfo()
	return route, apiResponse, err
}

// ListRoutes sends a GET request to the API to retrieve all routes from a given gateway.
func (c *Client) ListRoutes(ctx context.Context, gatewayID string) (apigateway.RouteReadList, *apigateway.APIResponse, error) {
	routes, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesGet(ctx, gatewayID).Execute()
	apiResponse.LogInfo()
	return routes, apiResponse, err
}

// DeleteRoute sends a DELETE request to the API to delete a route by ID from a given gateway.
func (c *Client) DeleteRoute(ctx context.Context, gatewayID string, routeID string) (*apigateway.APIResponse, error) {
	apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesDelete(ctx, gatewayID, routeID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsAPIGatewayRouteAvailable checks if the API Gateway Route is available.
func (c *Client) IsAPIGatewayRouteAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	routeID := d.Id()
	gatewayID := d.Get("gateway_id").(string)

	route, _, err := c.GetRouteByID(ctx, gatewayID, routeID)
	if err != nil {
		return false, err
	}

	if route.Metadata == nil || route.Metadata.Status == nil {
		return false, fmt.Errorf("expected metadata, got empty for API Gateway Route with ID: %s", routeID)
	}
	log.Printf("[DEBUG] API Gateway Route status: %s", *route.Metadata.Status)
	return strings.EqualFold(*route.Metadata.Status, constant.Available), nil
}

// IsAPIGatewayRouteDeleted checks if the API Gateway Route has been deleted.
func (c *Client) IsAPIGatewayRouteDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	routeID := d.Id()
	gatewayID := d.Get("gateway_id").(string)

	_, apiResponse, err := c.sdkClient.RoutesApi.ApigatewaysRoutesFindById(ctx, gatewayID, routeID).Execute()
	return apiResponse.HttpNotFound(), err
}

// SetAPIGatewayRouteData sets the data for the API Gateway Route.
func (c *Client) SetAPIGatewayRouteData(d *schema.ResourceData, route apigateway.RouteRead) error {
	d.SetId(*route.Id)

	if route.Properties == nil {
		return fmt.Errorf("expected properties in the response for the API Gateway Route with ID: %s, but received 'nil' instead", *route.Id)
	}

	if route.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the API Gateway Route with ID: %s, but received 'nil' instead", *route.Id)

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
			return fmt.Errorf("error setting upstreams for the API Gateway Route with ID: %s", *route.Id)
		}
	}

	return nil
}

func setRoutePostRequest(d *schema.ResourceData) *apigateway.RouteCreate {
	return apigateway.NewRouteCreate(*setRouteConfig(d))
}

func setRoutePutRequest(d *schema.ResourceData) *apigateway.RouteEnsure {
	routeID := d.Id()
	route := setRouteConfig(d)

	return apigateway.NewRouteEnsure(routeID, *route)
}

func setRouteConfig(d *schema.ResourceData) *apigateway.Route {
	routeName := d.Get("name").(string)
	routeType := d.Get("type").(string)
	pathsRaw := d.Get("paths").([]interface{})
	methodsRaw := d.Get("methods").([]interface{})
	websocket := d.Get("websocket").(bool)
	upstreamsRaw := d.Get("upstreams").([]interface{})

	paths := make([]string, 0)
	for _, path := range pathsRaw {
		paths = append(paths, path.(string))
	}

	methods := make([]string, 0)
	for _, method := range methodsRaw {
		methods = append(methods, method.(string))
	}

	upstreams := make([]apigateway.RouteUpstreams, 0)
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
