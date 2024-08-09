package dbaas

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *MongoClient) CreateUser(ctx context.Context, clusterId string, user mongo.User) (mongo.User, utils.ApiResponseInfo, error) {
	userResp, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPost(ctx, clusterId).User(user).Execute()
	apiResponse.LogInfo()
	return userResp, apiResponse, err
}

func (c *PsqlClient) CreateUser(ctx context.Context, clusterId string, user pgsql.User) (pgsql.UserResource, utils.ApiResponseInfo, error) {
	userResp, apiResponse, err := c.sdkClient.UsersApi.UsersPost(ctx, clusterId).User(user).Execute()
	apiResponse.LogInfo()
	return userResp, apiResponse, err
}

func (c *MongoClient) UpdateUser(ctx context.Context, clusterId, username string, patchUserReq mongo.PatchUserRequest) (mongo.User, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersPatch(ctx, clusterId, username).PatchUserRequest(patchUserReq).Execute()
	apiResponse.LogInfo()
	return user, apiResponse, err
}

func (c *PsqlClient) UpdateUser(ctx context.Context, clusterId, username string, patchUserReq pgsql.UsersPatchRequest) (pgsql.UserResource, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.UsersPatch(ctx, clusterId, username).UsersPatchRequest(patchUserReq).Execute()
	apiResponse.LogInfo()
	return user, apiResponse, err
}

func (c *MongoClient) GetUsers(ctx context.Context, clusterId string) (mongo.UsersList, utils.ApiResponseInfo, error) {
	users, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return users, apiResponse, err
}

func (c *MongoClient) FindUserByUsername(ctx context.Context, clusterId, username string) (mongo.User, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersFindById(ctx, clusterId, username).Execute()
	apiResponse.LogInfo()
	return user, apiResponse, err
}

func (c *PsqlClient) FindUserByUsername(ctx context.Context, clusterId, username string) (pgsql.UserResource, utils.ApiResponseInfo, error) {
	user, apiResponse, err := c.sdkClient.UsersApi.UsersGet(ctx, clusterId, username).Execute()
	apiResponse.LogInfo()
	return user, apiResponse, err
}

func (c *MongoClient) DeleteUser(ctx context.Context, clusterId, username string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersDelete(ctx, clusterId, username).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *PsqlClient) DeleteUser(ctx context.Context, clusterId, username string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.UsersApi.UsersDelete(ctx, clusterId, username).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsUserReady - checks the cluster, as it will move to busy while the user is created or updated
// There is no metadata state on the user
func (c *MongoClient) IsUserReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterIdIf, idOk := d.GetOk("cluster_id")
	usernameIf, _ := d.GetOk("username")
	username := usernameIf.(string)
	clusterId := clusterIdIf.(string)
	if !idOk {
		return false, fmt.Errorf("id missing from schema for cluster with id %s", d.Id())
	}
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return true, fmt.Errorf("error checking cluster status: %w", err)
	}
	if cluster.Metadata == nil || cluster.Metadata.State == nil {
		log.Printf("cluster metadata or state is empty for cluster %s in cluster %s ", username, clusterId)
		return false, fmt.Errorf("cluster metadata or state is empty for id %s", d.Id())
	}
	log.Printf("[INFO] state of the cluster %s ", string(*cluster.Metadata.State))
	return strings.EqualFold(string(*cluster.Metadata.State), constant.Available), nil
}

func (c *MongoClient) IsUserDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterIdIf, idOk := d.GetOk("cluster_id")
	usernameIf, nameOk := d.GetOk("username")
	if !nameOk {
		return false, fmt.Errorf("username missing from schema for user with id %s", d.Id())
	}
	if !idOk {
		return false, fmt.Errorf("id missing from schema for user with id %s", d.Id())
	}
	_, apiResponse, err := c.sdkClient.UsersApi.ClustersUsersFindById(ctx, clusterIdIf.(string), usernameIf.(string)).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("error checking user deletion status: %w", err)
	}
	return false, nil
}

func (c *PsqlClient) IsUserDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	_, apiResponse, err := c.sdkClient.UsersApi.UsersGet(ctx, clusterId, username).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("error checking user deletion status: %w", err)
	}
	return false, nil
}

func SetUserMongoData(d *schema.ResourceData, user *mongo.User) error {
	if user.Properties != nil {
		if user.Properties.Username != nil {
			if err := d.Set("username", *user.Properties.Username); err != nil {
				return err
			}
		}

		if user.Properties.Roles != nil && len(*user.Properties.Roles) > 0 {
			userRoles := make([]interface{}, len(*user.Properties.Roles))
			for index, user := range *user.Properties.Roles {
				userEntry := make(map[string]interface{})

				if user.Role != nil {
					userEntry["role"] = *user.Role
				}

				if user.Database != nil {
					userEntry["database"] = user.Database
				}
				userRoles[index] = userEntry
			}

			if len(userRoles) > 0 {
				if err := d.Set("roles", userRoles); err != nil {
					return fmt.Errorf("error setting user roles for user (%w)", err)
				}
			}
		}
	}
	return nil
}

func SetUserPgSqlData(d *schema.ResourceData, user *pgsql.UserResource) error {
	resourceName := "PgSQL user"
	d.SetId(*user.Id)
	if user.Properties == nil {
		return fmt.Errorf("expected properties in the response for the PgSql user with ID: %s, but received 'nil' instead", *user.Id)
	}
	if user.Properties.Username != nil {
		if err := d.Set("username", *user.Properties.Username); err != nil {
			return utils.GenerateSetError(resourceName, "username", err)
		}
	}
	if user.Properties.System != nil {
		if err := d.Set("is_system_user", *user.Properties.System); err != nil {
			return utils.GenerateSetError(resourceName, "is_system_user", err)
		}
	}
	return nil
}
