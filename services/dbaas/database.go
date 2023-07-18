package dbaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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

func SetDatabasePgSqlData(d *schema.ResourceData, database *pgsql.DatabaseResource) error {
	d.SetId(*database.Id)
	if database.Properties == nil {
		return fmt.Errorf("expected properties in the response for the PgSql database with ID: %s, but received 'nil' instead", *database.Id)
	}
	if database.Properties.Name != nil {
		if err := d.Set("name", *database.Properties.Name); err != nil {
			return err
		}
	}
	if database.Properties.Owner != nil {
		if err := d.Set("owner", *database.Properties.Owner); err != nil {
			return err
		}
	}
	return nil
}
