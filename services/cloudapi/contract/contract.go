package contract

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

// GetContractNumber fetches the contract number from the API.
// Returns the contract number if exactly one contract is found, empty string otherwise.
func GetContractNumber(ctx context.Context, client *bundleclient.SdkBundle) string {
	apiClient, err := client.NewCloudAPIClient(ctx, "")
	if err != nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	contracts, _, err := apiClient.ContractResourcesApi.ContractsGet(ctx).Execute()
	if err != nil {
		return ""
	}
	if contracts.Items == nil || len(*contracts.Items) != 1 {
		return ""
	}
	c := (*contracts.Items)[0]
	if c.Properties != nil && c.Properties.ContractNumber != nil {
		return fmt.Sprintf("%d", *c.Properties.ContractNumber)
	}
	return ""
}
