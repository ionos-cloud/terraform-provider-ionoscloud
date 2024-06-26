package apigateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
