package dbaas

import (
	"context"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *MongoClient) CreateUser(ctx context.Context, clusterId string, user mongo.User) (*mongo.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPost(ctx, clusterId).User(user).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func (c *MongoClient) UpdateUser(ctx context.Context, clusterId, username string, patchUserReq mongo.PatchUserRequest) (*mongo.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPatch(ctx, clusterId, "admin", username).PatchUserRequest(patchUserReq).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func (c *MongoClient) GetUsers(ctx context.Context, clusterId string) (mongo.UsersList, *mongo.APIResponse, error) {
	users, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return users, apiResponse, err
	}
	return users, nil, err
}

func (c *MongoClient) FindUserByUsername(ctx context.Context, clusterId, username string) (mongo.User, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersFindById(ctx, clusterId, "admin", username).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return user, apiResponse, err
	}
	return user, nil, err
}

func (c *MongoClient) DeleteUser(ctx context.Context, clusterId, username string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersDelete(ctx, clusterId, "admin", username).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}
