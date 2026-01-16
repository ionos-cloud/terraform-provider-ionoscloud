package kafka

import (
	"context"

	kafka "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

// GetUsers returns the list of users for a cluster using the cluster ID and the location in which the cluster resides
func (c *Client) GetUsers(ctx context.Context, clusterID, location string) (kafka.UserReadList, utils.ApiResponseInfo, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Kafka, location)
	users, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersGet(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return users, apiResponse, err
}
