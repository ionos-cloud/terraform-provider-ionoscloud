package logging

import (
	"context"

	logging "github.com/ionos-cloud/sdk-go-logging"
)

func (c *Client) GetCentralLogging(ctx context.Context) (logging.CentralLoggingResponse, *logging.APIResponse, error) {
	centralLogging, apiResponse, err := c.sdkClient.CentralApi.CentralLoggingGet(ctx).Execute()
	apiResponse.LogInfo()
	return centralLogging, apiResponse, err
}
