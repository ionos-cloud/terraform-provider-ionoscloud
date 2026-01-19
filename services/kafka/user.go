package kafka

import (
	"context"
	"fmt"

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

// GetUserCredentials returns the access credentials using the cluster ID, user ID and location.
func (c *Client) GetUserCredentialsByID(ctx context.Context, clusterID, userID, location string) (kafka.UserReadAccess, utils.ApiResponseInfo, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Kafka, location)
	userCredentials, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersAccessGet(ctx, clusterID, userID).Execute()
	apiResponse.LogInfo()
	return userCredentials, apiResponse, err
}

// TODO -- Test method
// GetUserCredentialsByName returns the access credentials using the cluster ID, username and location.
func (c *Client) GetUserCredentialsByName(ctx context.Context, clusterID, username, location string) (kafka.UserReadAccess, utils.ApiResponseInfo, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Kafka, location)
	var result kafka.UserReadAccess
	var intermediate kafka.UserRead

	// Fetch all users.
	usersResp, _, err := c.GetUsers(ctx, clusterID, location)
	if err != nil {
		return result, nil, err
	}

	// Search for the appropriate user using the name.
	nResults := 0
	for _, user := range usersResp.Items {
		if user.Properties.Name == username {
			intermediate = user
			nResults++
		}
	}

	if nResults == 0 {
		return result, nil, fmt.Errorf("no Kafka user was found using the specified username: %s", username)
	}
	// TODO -- Check if the username is unique
	if nResults > 1 {
		return result, nil, fmt.Errorf("more than one Kafka user was found using the specified username: %s", username)
	}

	// Fetch access data using the user ID.
	result, apiResponse, err := c.GetUserCredentialsByID(ctx, clusterID, intermediate.Id, location)
	return result, apiResponse, err
}
