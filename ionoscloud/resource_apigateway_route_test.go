//go:build all || apigateway || route

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccAPIGatewayRouteBasic(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					VersionConstraint: "3.4.3",
					Source:            "hashicorp/random",
				},
				"time": {
					Source:            "hashicorp/time",
					VersionConstraint: "0.11.1",
				},
			},
			ProviderFactories: testAccProviderFactories,
			CheckDestroy:      testAccCheckAPIGatewayRouteDestroyCheck,
			Steps: []resource.TestStep{
				{
					Config: configAPIGatewayRouteBasic(routeResourceName, routeAttributeNameValue),
					Check:  checkAPIGatewayRouteResource(routeAttributeNameValue),
				},
				{
					Config: configAPIGatewayRouteBasic(routeResourceName, routeAttributeNameValueUpdated),
					Check:  checkAPIGatewayRouteResource(routeAttributeNameValueUpdated),
				},
			},
		},
	)
}

func TestAccAPIGatewayRouteDataSourceGetByID(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					VersionConstraint: "3.4.3",
					Source:            "hashicorp/random",
				},
				"time": {
					Source:            "hashicorp/time",
					VersionConstraint: "0.11.1",
				},
			},
			ProviderFactories: testAccProviderFactories,
			CheckDestroy:      testAccCheckAPIGatewayRouteDestroyCheck,
			Steps: []resource.TestStep{
				{
					Config: configAPIGatewayRouteDataSourceGetByID(routeResourceName, routeDataName, constant.ApiGatewayRouteResource+"."+routeResourceName+".id"),
					Check:  checkAPIGatewayRouteAttributes(constant.DataSource+"."+constant.ApiGatewayRouteResource+"."+routeDataName, routeAttributeNameValue),
				},
				{
					Config:      configAPIGatewayRouteDataSourceGetByID(routeResourceName, routeDataName, `"00000000-0000-0000-0000-000000000000"`),
					ExpectError: regexp.MustCompile("an error occurred while fetching the API Gateway Route with ID"),
				},
			},
		},
	)
}

func TestAccAPIGatewayRouteDataSourceGetByName(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					VersionConstraint: "3.4.3",
					Source:            "hashicorp/random",
				},
				"time": {
					Source:            "hashicorp/time",
					VersionConstraint: "0.11.1",
				},
			},
			ProviderFactories: testAccProviderFactories,
			CheckDestroy:      testAccCheckAPIGatewayRouteDestroyCheck,
			Steps: []resource.TestStep{
				{
					Config: configAPIGatewayRouteDataSourceGetByName(routeResourceName, routeDataName, constant.ApiGatewayRouteResource+"."+routeResourceName+".name"),
					Check:  checkAPIGatewayRouteAttributes(constant.DataSource+"."+constant.ApiGatewayRouteResource+"."+routeDataName, routeAttributeNameValue),
				},
				{
					Config:      configAPIGatewayRouteDataSourceGetByName(routeResourceName, routeDataName, `"wrongname"`),
					ExpectError: regexp.MustCompile("no API Gateway Route found with the specified name"),
				},
			},
		},
	)
}

func testAccCheckAPIGatewayRouteExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundRoute, _, err := client.GetRouteByID(ctx, rs.Primary.Attributes["gateway_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching API Gateway Route with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *foundRoute.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

func testAccCheckAPIGatewayRouteDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ApiGatewayRouteResource {
			continue
		}

		_, apiResponse, err := client.GetRouteByID(ctx, rs.Primary.Attributes["gateway_id"], rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of API Gateway Route with ID: %v, error: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("API Gateway Route with ID: %v still exists", rs.Primary.ID)
		}
	}

	return nil
}

func checkAPIGatewayRouteResource(attributeNameValueReference string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckAPIGatewayRouteExists(constant.ApiGatewayRouteResource+"."+routeResourceName),
		checkAPIGatewayRouteAttributes(constant.ApiGatewayRouteResource+"."+routeResourceName, attributeNameValueReference),
	)
}

func configAPIGatewayRouteBasic(resourceName, attributeName string) string {
	pathsValue := fmt.Sprintf(`"%s"`, routeAttributePathsValue)
	methodsValue := fmt.Sprintf(`"%s"`, routeAttributeMethodsValue)

	routeBasicConfig := fmt.Sprintf(
		templateAPIGatewayRouteConfig, resourceName, attributeName, routeAttributeTypeValue, pathsValue, methodsValue, routeAttributeWebsocketValue,
		routeAttributeUpstreamsSchemeValue, routeAttributeUpstreamsLoadbalancerValue, routeAttributeUpstreamsHostValue, routeAttributeUpstreamsPortValue,
		routeAttributeUpstreamsWeightValue,
	)
	return strings.Join([]string{defaultAPIGatewayConfig, routeBasicConfig}, "\n")
}

func checkAPIGatewayRouteAttributes(fullResourceName, attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeName, attributeNameReferenceValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeType, routeAttributeTypeValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributePaths, routeAttributePathsValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeMethods, routeAttributeMethodsValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeWebsocket, routeAttributeWebsocketValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeUpstreamsScheme, routeAttributeUpstreamsSchemeValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeUpstreamsLoadbalancer, routeAttributeUpstreamsLoadbalancerValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeUpstreamsHost, routeAttributeUpstreamsHostValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeUpstreamsPort, routeAttributeUpstreamsPortValue),
		resource.TestCheckResourceAttr(fullResourceName, routeAttributeUpstreamsWeight, routeAttributeUpstreamsWeightValue),
	)
}

func configAPIGatewayRouteDataSourceGetByID(resourceName, dataSourceName, attributeId string) string {
	dataSourceBasicConfig := fmt.Sprintf(templateAPIGatewayRouteDataIDConfig, dataSourceName, attributeId)
	baseConfig := configAPIGatewayRouteBasic(resourceName, routeAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

func configAPIGatewayRouteDataSourceGetByName(resourceName, dataSourceName, attributeName string) string {
	dataSourceBasicConfig := fmt.Sprintf(templateAPIGatewayRouteDataNameConfig, dataSourceName, attributeName)
	baseConfig := configAPIGatewayRouteBasic(resourceName, routeAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

var (
	routeResourceName = "example_route"
	routeDataName     = "example_route_data"

	routeAttributeName             = "name"
	routeAttributeNameValue        = "exampleroute"
	routeAttributeNameValueUpdated = "examplerouteupdated"

	routeAttributeType      = "type"
	routeAttributeTypeValue = "http"

	routeAttributePaths      = "paths.0"
	routeAttributePathsValue = "/foo/*"

	routeAttributeMethods      = "methods.0"
	routeAttributeMethodsValue = "GET"

	routeAttributeWebsocket      = "websocket"
	routeAttributeWebsocketValue = "false"

	routeAttributeUpstreamsScheme      = "upstreams.0.scheme"
	routeAttributeUpstreamsSchemeValue = "http"

	routeAttributeUpstreamsHost      = "upstreams.0.host"
	routeAttributeUpstreamsHostValue = "example.com"

	routeAttributeUpstreamsPort      = "upstreams.0.port"
	routeAttributeUpstreamsPortValue = "80"

	routeAttributeUpstreamsWeight      = "upstreams.0.weight"
	routeAttributeUpstreamsWeightValue = "100"

	routeAttributeUpstreamsLoadbalancer      = "upstreams.0.loadbalancer"
	routeAttributeUpstreamsLoadbalancerValue = "round-robin"
)

var templateAPIGatewayRouteConfig = `resource "ionoscloud_apigateway_route" "%s" {
	name = "%s"
	type = "%s"
	paths = [
		%s
	]
	methods = [
		%s
	]
	websocket = %s
	upstreams {
		scheme       = "%s"
		loadbalancer = "%s"
		host         = "%s"
		port         = %s
		weight       = %s
	}
	gateway_id = ionoscloud_apigateway.example_gateway.id
}`

var templateAPIGatewayRouteDataIDConfig = `data "ionoscloud_apigateway_route" "%s" {
	gateway_id = ionoscloud_apigateway.example_gateway.id
	id = %s
}`

var templateAPIGatewayRouteDataNameConfig = `data "ionoscloud_apigateway_route" "%s" {
	gateway_id = ionoscloud_apigateway.example_gateway.id
	name = %s
}`

var defaultAPIGatewayConfig = `resource "ionoscloud_apigateway" "example_gateway" {
	name              = "examplegateway"
    logs              = false
    metrics           = false

    custom_domains {
        name           = "example.com"
        certificate_id = "00000000-0000-0000-0000-000000000000"
    }

    custom_domains {
        name           = "example.org"
        certificate_id = "00000000-0000-0000-0000-000000000000"
    }
}`
