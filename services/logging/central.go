package logging

import (
	"context"

	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetCentralLogging will retrieve the central logging configuration
func (c *Client) GetCentralLogging(ctx context.Context) (logging.CentralLoggingReadList, *shared.APIResponse, error) {
	centralLogging, apiResponse, err := c.sdkClient.CentralApi.CentralGet(ctx).Execute()
	apiResponse.LogInfo()
	return centralLogging, apiResponse, err
}
