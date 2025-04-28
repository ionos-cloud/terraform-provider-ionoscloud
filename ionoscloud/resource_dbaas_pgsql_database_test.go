//go:build all || dbaas || psql
// +build all dbaas psql

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	pgsql "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccPgSqlDatabase(t *testing.T) {
	var database pgsql.DatabaseResource

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             pgSqlDatabaseDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: PgSqlDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					pgSqlDatabaseExistsCheck(constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseTestResource, &database),
					resource.TestCheckResourceAttr(constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseTestResource, databaseNameAttribute, databaseNameValue),
					resource.TestCheckResourceAttr(constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseTestResource, databaseOwnerAttribute, databaseOwnerValue),
				),
			},
			{
				Config: PgSqlDatabaseDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseDataSourceByName, databaseNameAttribute, constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseTestResource, databaseNameAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseDataSourceByName, databaseOwnerAttribute, constant.PsqlDatabaseResource+"."+constant.PsqlDatabaseTestResource, databaseOwnerAttribute),
				),
			},
			{
				Config:      PgSqlDatabaseDataSourceWrongName,
				ExpectError: regexp.MustCompile(`no PgSql database found with the specified name`),
			},
			{
				Config: PgSqlAllDatabasesDataSource,
				// Check only the length since there are some databases that already exist in the
				// cluster.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.PsqlDatabasesResource+"."+constant.PsqlDatabasesDataSource, databasesAttribute+".#", "4"),
				),
			},
			{
				Config: PgSqlAllDatabasesFilterByOwnerDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.PsqlDatabasesResource+"."+constant.PsqlDatabasesDataSource, databasesAttribute+".#", "1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.PsqlDatabasesResource+"."+constant.PsqlDatabasesDataSource, databasesAttribute+".0.name", databaseNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.PsqlDatabasesResource+"."+constant.PsqlDatabasesDataSource, databasesAttribute+".0.owner", databaseOwnerValue),
				),
			},
		},
	})
}

func pgSqlDatabaseExistsCheck(path string, database *pgsql.DatabaseResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).PsqlClient
		rs, ok := s.RootModule().Resources[path]
		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the PgSql database")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		foundDatabase, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
		apiResponse.LogInfo()
		if err != nil {
			return fmt.Errorf("error occurred while fetching the PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err)
		}
		database = &foundDatabase
		return nil
	}
}

func pgSqlDatabaseDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).PsqlClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PsqlDatabaseResource {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		_, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the deletion of PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err)
			}
		} else {
			return fmt.Errorf("PgSql database %s still exists in the cluster with ID: %s", name, clusterId)
		}
	}
	return nil
}

// Configurations

const PgSqlDatabaseConfig = PgSqlUserConfig + `
resource ` + constant.PsqlDatabaseResource + ` ` + constant.PsqlDatabaseTestResource + ` {
  ` + clusterIdAttribute + ` = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id  
  ` + databaseNameAttribute + ` = "` + databaseNameValue + `"
  ` + databaseOwnerAttribute + ` = ` + constant.PsqlUserResource + `.` + constant.UserTestResource + `.username
}
`

const PgSqlDatabaseDataSource = PgSqlDatabaseConfig + `
data ` + constant.PsqlDatabaseResource + ` ` + constant.PsqlDatabaseDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id   
  ` + databaseNameAttribute + ` = ` + constant.PsqlDatabaseResource + `.` + constant.PsqlDatabaseTestResource + `.name
}
`

const PgSqlDatabaseDataSourceWrongName = PgSqlDatabaseConfig + `
data ` + constant.PsqlDatabaseResource + ` ` + constant.PsqlDatabaseDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id   
  ` + databaseNameAttribute + ` = "nonexistent"
}
`

const PgSqlAllDatabasesDataSource = PgSqlDatabaseConfig + `
data ` + constant.PsqlDatabasesResource + ` ` + constant.PsqlDatabasesDataSource + ` {
  ` + clusterIdAttribute + ` = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id   
}
`

const PgSqlAllDatabasesFilterByOwnerDataSource = PgSqlDatabaseConfig + `
data ` + constant.PsqlDatabasesResource + ` ` + constant.PsqlDatabasesDataSource + ` {
  ` + clusterIdAttribute + ` = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id 
  ` + databaseOwnerAttribute + ` = "` + databaseOwnerValue + `"
}
`

// Attributes
const databasesAttribute = "databases"
const databaseNameAttribute = "name"
const databaseOwnerAttribute = "owner"

// Values
const databaseNameValue = "testdatabase"
const databaseOwnerValue = usernameValue
