package lanSvc

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// Currently, this is not a service per se, but in the future, when the LAN service will be created,
// this function can be included in the service. Right now, it is just an utility function in other
// to reuse the code.
func FindLanById(client ionoscloud.APIClient, ctx context.Context, datacenterId string, lanId string) (ionoscloud.Lan, *ionoscloud.APIResponse, error) {
	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, datacenterId, lanId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return lan, apiResponse, fmt.Errorf("an error occured while retrieving existing IP failover groups for LAN with ID: %s, datacenter ID: %s, error: %w", lanId, datacenterId, err)
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return lan, apiResponse, fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead")
	}
	return lan, apiResponse, nil
}
