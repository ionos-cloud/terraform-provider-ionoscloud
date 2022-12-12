package dbaas

import (
	"context"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *MongoClient) CreateUser(ctx context.Context, clusterId string, user mongo.User) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPost(ctx, clusterId).User(user).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *MongoClient) UpdateUser(ctx context.Context, clusterId, username string, patchUserReq mongo.PatchUserRequest) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPatch(ctx, clusterId, DefaultMongoDatabase, username).PatchUserRequest(patchUserReq).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *MongoClient) GetUsers(ctx context.Context, clusterId string) (mongo.UsersList, utils.ApiResponseInfo, error) {
	users, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return users, apiResponse, err
}

func (c *MongoClient) FindUserByUsername(ctx context.Context, clusterId, username string) (mongo.User, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersFindById(ctx, clusterId, DefaultMongoDatabase, username).Execute()
	apiResponse.LogInfo()
	return user, apiResponse, err
}

func (c *MongoClient) DeleteUser(ctx context.Context, clusterId, username string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersDelete(ctx, clusterId, DefaultMongoDatabase, username).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}
