package cloudapilan

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Currently, this is not a service per se, but in the future, when the LAN service will be created,
// this function can be included in the service. Right now, it is just an utility function in order
// to reuse the code.
func FindLanById(client ionoscloud.APIClient, ctx context.Context, datacenterId string, lanId string) (ionoscloud.Lan, *shared.APIResponse, error) {
	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, datacenterId, lanId).Execute()
	apiResponse.LogInfo()
	return lan, apiResponse, err
}
