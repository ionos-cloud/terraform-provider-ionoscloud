package dbaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pgsql "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *PsqlClient) CreateDatabase(ctx context.Context, clusterId string, database pgsql.Database) (pgsql.DatabaseResource, utils.ApiResponseInfo, error) {
	databaseResponse, apiResponse, err := c.sdkClient.DatabasesApi.DatabasesPost(ctx, clusterId).Database(database).Execute()
	apiResponse.LogInfo()
	return databaseResponse, apiResponse, err
}

func (c *PsqlClient) DeleteDatabase(ctx context.Context, clusterId, name string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.DatabasesApi.DatabasesDelete(ctx, clusterId, name).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *PsqlClient) FindDatabaseByName(ctx context.Context, clusterId, name string) (pgsql.DatabaseResource, utils.ApiResponseInfo, error) {
	database, apiResponse, err := c.sdkClient.DatabasesApi.DatabasesGet(ctx, clusterId, name).Execute()
	apiResponse.LogInfo()
	return database, apiResponse, err
}

func (c *PsqlClient) GetDatabases(ctx context.Context, clusterId string) (pgsql.DatabaseList, utils.ApiResponseInfo, error) {
	databases, apiResponse, err := c.sdkClient.DatabasesApi.DatabasesList(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return databases, apiResponse, err
}

func SetDatabasePgSqlData(d *schema.ResourceData, database *pgsql.DatabaseResource) error {
	resourceName := "PgSQL database"
	d.SetId(database.Id)

	if database.Properties.Name != "" {
		if err := d.Set("name", database.Properties.Name); err != nil {
			return utils.GenerateSetError(resourceName, "name", err)
		}
	}
	if database.Properties.Owner != "" {
		if err := d.Set("owner", database.Properties.Owner); err != nil {
			return utils.GenerateSetError(resourceName, "owner", err)
		}
	}
	return nil
}
