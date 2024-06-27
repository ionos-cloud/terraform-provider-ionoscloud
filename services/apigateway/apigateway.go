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

func (c *Client) GetApiGatewayById(ctx context.Context, id string) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	apiGateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiGateway, apiResponse, err
}

func (c *Client) ListApiGateways(ctx context.Context) (apigateway.GatewayReadList, *apigateway.APIResponse, error) {
	apiGateways, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysGet(ctx).Execute()
	apiResponse.LogInfo()
	return apiGateways, apiResponse, err
}

func (c *Client) DeleteApiGateway(ctx context.Context, id string) (*apigateway.APIResponse, error) {
	apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) UpdateApiGateway(ctx context.Context, id string, gw apigateway.GatewayEnsure) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	gateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysPut(ctx, id).GatewayEnsure(gw).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

func (c *Client) CreateApiGateway(ctx context.Context, gw apigateway.GatewayCreate) (apigateway.GatewayRead, *apigateway.APIResponse, error) {
	gateway, apiResponse, err := c.sdkClient.APIGatewaysApi.ApigatewaysPost(ctx).GatewayCreate(gw).Execute()
	apiResponse.LogInfo()
	return gateway, apiResponse, err
}

func (c *Client) SetApiGatewayData(d *schema.ResourceData, apiGateway apigateway.GatewayRead) error {
	d.SetId(*apiGateway.Id)

	if apiGateway.Properties == nil {
		return fmt.Errorf("expected properties in the response for the ApiGateway with ID: %s, but received 'nil' instead", *apiGateway.Id)
	}

	if apiGateway.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the ApiGateway with ID: %s, but received 'nil' instead", *apiGateway.Id)
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
			return fmt.Errorf("error setting custom_domains for the ApiGateway with ID: %s", *apiGateway.Id)
		}
	}

	if apiGateway.Metadata.PublicEndpoint != nil {
		if err := d.Set("public_endpoint", *apiGateway.Metadata.PublicEndpoint); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) IsGatewayReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Id()
	gateway, _, err := c.GetApiGatewayById(ctx, gatewayID)
	if err != nil {
		return true, fmt.Errorf("status check failed for Gateway ID: %v, error: %w", gatewayID, err)
	}

	if gateway.Metadata == nil || gateway.Metadata.Status == nil {
		return false, fmt.Errorf("metadata or status is empty for Gateway ID: %v", gatewayID)
	}

	log.Printf("[INFO] state of the MariaDB cluster with ID: %v is: %s ", gatewayID, *gateway.Metadata.Status)
	return strings.EqualFold(*gateway.Metadata.Status, constant.Available), nil
}

func (c *Client) IsGatewayDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	gatewayID := d.Id()
	_, apiResponse, err := c.GetApiGatewayById(ctx, gatewayID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for Gateway deletion status, ID: %v, error: %w", gatewayID, err)
	}
	return false, nil
}

func GetGatewayDataCreate(d *schema.ResourceData) (*apigateway.GatewayCreate, error) {
	gateway := apigateway.GatewayCreate{
		Properties: &apigateway.Gateway{},
	}

	if v, ok := d.GetOk("name"); ok {
		gateway.Properties.Name = apigateway.ToPtr(v.(string))
	}

	if v, ok := d.GetOk("logs"); ok {
		gateway.Properties.Logs = apigateway.ToPtr(v.(bool))
	}

	if v, ok := d.GetOk("metrics"); ok {
		gateway.Properties.Metrics = apigateway.ToPtr(v.(bool))
	}

	if v, ok := d.GetOk("custom_domains"); ok {
		domains := v.([]interface{})
		customDomains := make([]apigateway.GatewayCustomDomains, len(domains))
		for i, domain := range domains {
			domainMap := domain.(map[string]interface{})
			customDomains[i] = apigateway.GatewayCustomDomains{
				Name:          apigateway.ToPtr(domainMap["name"].(string)),
				CertificateId: apigateway.ToPtr(domainMap["certificate_id"].(string)),
			}
		}
		gateway.Properties.CustomDomains = &customDomains
	}

	return &gateway, nil
}

func GetGatewayDataEnsure(d *schema.ResourceData) (*apigateway.GatewayEnsure, error) {
	// TODO: This doesn't exist for MariaDB. Check if it's needed for API Gateway
	return &apigateway.GatewayEnsure{}, nil
}
