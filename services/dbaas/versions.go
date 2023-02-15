package dbaas

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func (c *PsqlClient) GetClusterVersions(ctx context.Context, clusterId string) (psql.PostgresVersionList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.ClustersApi.ClusterPostgresVersionsGet(ctx, clusterId).Execute()
	return versions, apiResponse, err
}

func (c *PsqlClient) GetAllVersions(ctx context.Context) (psql.PostgresVersionList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.ClustersApi.PostgresVersionsGet(ctx).Execute()
	return versions, apiResponse, err
}

func SetPgSqlVersionsData(d *schema.ResourceData, postgresVersions psql.PostgresVersionList) diag.Diagnostics {

	if postgresVersions.Data != nil {
		var versions []string
		for _, version := range *postgresVersions.Data {
			if version.Name != nil {
				versions = append(versions, *version.Name)
			}
		}
		err := d.Set("postgres_versions", versions)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting postgres_versions: %w", err))
			return diags
		}
	}

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	return nil
}
