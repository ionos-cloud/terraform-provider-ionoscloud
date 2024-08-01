package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-api-gateway"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// GetAPIGatewayByID returns a gateway given an ID
func (c *Client) GetAPIGatewayByID(ctx context.Context, id string) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	apiGateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiGateway, apiResponse, err
}

// ListAPIGateways returns a list of all gateways
func (c *Client) ListAPIGateways(ctx context.Context) (apigateway.GatewayReadList, *apigateway.APIResponse, error) {
	apiGateways, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return apiGateways, apiResponse, err
}

// DeleteAPIGateway deletes a gateway given an ID
func (c *Client) DeleteAPIGateway(ctx context.Context, id string) (*apigateway.APIResponse, error) {
	apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateAPIGateway updates a gateway given an ID or creates a new one if it doesn't exist
func (c *Client) UpdateAPIGateway(ctx context.Context, d *schema.ResourceData) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	gateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysPut(ctx, d.Id()).
		GatewayEnsure(*setGatewayPutRequest(d)).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

// CreateAPIGateway creates a new gateway
func (c *Client) CreateAPIGateway(ctx context.Context, d *schema.ResourceData) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	gateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysPost(ctx).
		GatewayCreate(*setGatewayPostRequest(d)).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

// SetAPIGatewayData sets the data of the gateway in the terraform resource
func (c *Client) SetAPIGatewayData(d *schema.ResourceData, apiGateway apigateway.GatewayRead) error {
	d.SetId(*apiGateway.Id)

	if apiGateway.Properties == nil {
		return fmt.Errorf("expected properties in the response for the API Gateway with ID: %s, but received 'nil' instead", *apiGateway.Id)
	}

	if apiGateway.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the API Gateway with ID: %s, but received 'nil' instead", *apiGateway.Id)
	}

	if apiGateway.Properties.Name != nil {
		if err := d.Set("name", *apiGateway.Properties.Name); err != nil {
			return err
		}
	}

	if apiGateway.Properties.Logs != nil {
		if err := d.Set("logs", *apiGateway.Properties.Logs); err != nil {
			return err
		}
	}

	if apiGateway.Properties.Metrics != nil {
		if err := d.Set("metrics", *apiGateway.Properties.Metrics); err != nil {
			return err
		}
	}

	if apiGateway.Properties.CustomDomains != nil {
		var customDomains []map[string]interface{}
		for _, customDomain := range *apiGateway.Properties.CustomDomains {
			customDomainData := map[string]interface{}{}

			utils.SetPropWithNilCheck(customDomainData, "name", customDomain.Name)
			utils.SetPropWithNilCheck(customDomainData, "certificate_id", customDomain.CertificateId)

			customDomains = append(customDomains, customDomainData)
		}

		if err := d.Set("custom_domains", customDomains); err != nil {
			return fmt.Errorf("error setting custom_domains for the API Gateway with ID %s: %w", *apiGateway.Id, err)
		}
	}

	if apiGateway.Metadata.PublicEndpoint != nil {
		if err := d.Set("public_endpoint", *apiGateway.Metadata.PublicEndpoint); err != nil {
			return err
		}
	}

	return nil
}

// IsGatewayReady checks if the gateway is ready
func (c *Client) IsGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Id()
	gateway, _, err := c.GetAPIGatewayByID(ctx, gatewayID)
	if err != nil {
		return true, fmt.Errorf("status check failed for Gateway ID: %v, error: %w", gatewayID, err)
	}

	if gateway.Metadata == nil || gateway.Metadata.Status == nil {
		return false, fmt.Errorf("metadata or status is empty for API Gateway ID: %v", gatewayID)
	}

	log.Printf("[INFO] state of the gateway with ID %s is: %s ", gatewayID, *gateway.Metadata.Status)
	return strings.EqualFold(*gateway.Metadata.Status, constant.Available), nil
}

// IsGatewayDeleted checks if the gateway is deleted
func (c *Client) IsGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Id()
	_, apiResponse, err := c.GetAPIGatewayByID(ctx, gatewayID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for API Gateway deletion status, ID: %v, error: %w", gatewayID, err)
	}
	return false, nil
}

func setGatewayPostRequest(d *schema.ResourceData) *apigateway.GatewayCreate {
	return apigateway.NewGatewayCreate(setGatewayConfig(d))
}

func setGatewayPutRequest(d *schema.ResourceData) *apigateway.GatewayEnsure {
	gatewayID := d.Id()
	gateway := setGatewayConfig(d)

	return apigateway.NewGatewayEnsure(gatewayID, gateway)
}

func setGatewayConfig(d *schema.ResourceData) apigateway.Gateway {
	gatewayName := d.Get("name").(string)
	logs := d.Get("logs").(bool)
	metrics := d.Get("metrics").(bool)
	customDomainsRaw := d.Get("custom_domains").([]interface{})

	customDomains := make([]apigateway.GatewayCustomDomains, len(customDomainsRaw))
	for i, domain := range customDomainsRaw {
		domainData := domain.(map[string]interface{})
		name := domainData["name"].(string)
		certificateID := domainData["certificate_id"].(string)

		customDomainObj := apigateway.GatewayCustomDomains{
			Name:          &name,
			CertificateId: &certificateID,
		}

		customDomains[i] = customDomainObj
	}

	return apigateway.Gateway{
		Name:          apigateway.ToPtr(gatewayName),
		Logs:          apigateway.ToPtr(logs),
		Metrics:       apigateway.ToPtr(metrics),
		CustomDomains: &customDomains,
	}
}
