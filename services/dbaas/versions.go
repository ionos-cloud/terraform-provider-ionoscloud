package dbaas

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type VersionService interface {
	GetClusterVersions(ctx context.Context, clusterId string) (dbaas.PostgresVersionList, *dbaas.APIResponse, error)
	GetAllVersions(ctx context.Context) (dbaas.PostgresVersionList, *dbaas.APIResponse, error)
}

func (c *Client) GetClusterVersions(ctx context.Context, clusterId string) (dbaas.PostgresVersionList, *dbaas.APIResponse, error) {
	versions, apiResponse, err := c.ClustersApi.ClusterPostgresVersionsGet(ctx, clusterId).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err

	}
	return versions, nil, err
}

func (c *Client) GetAllVersions(ctx context.Context) (dbaas.PostgresVersionList, *dbaas.APIResponse, error) {
	versions, apiResponse, err := c.ClustersApi.PostgresVersionsGet(ctx).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err
	}
	return versions, nil, err
}

func SetPgSqlVersionsData(d *schema.ResourceData, postgresVersions dbaas.PostgresVersionList) diag.Diagnostics {

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
