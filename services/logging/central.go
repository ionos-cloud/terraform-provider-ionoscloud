package logging

import (
	"context"

	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func (c *Client) GetCentralLogging(ctx context.Context) (logging.CentralLoggingResponse, *shared.APIResponse, error) {
	centralLogging, apiResponse, err := c.sdkClient.CentralApi.CentralLoggingGet(ctx).Execute()
	apiResponse.LogInfo()
	return centralLogging, apiResponse, err
}
