//go:build all || dbaas || psql
// +build all dbaas psql

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"regexp"
	"testing"
)

func TestAccPgSqlDatabase(t *testing.T) {
	var database pgsql.DatabaseResource

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      pgSqlDatabaseDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: PgSqlDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					pgSqlDatabaseExistsCheck(PsqlDatabaseResource+"."+PsqlDatabaseTestResource, &database),
					resource.TestCheckResourceAttr(PsqlDatabaseResource+"."+PsqlDatabaseTestResource, databaseNameAttribute, databaseNameValue),
					resource.TestCheckResourceAttr(PsqlDatabaseResource+"."+PsqlDatabaseTestResource, databaseOwnerAttribute, databaseOwnerValue),
				),
			},
			{
				Config: PgSqlDatabaseDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlDatabaseResource+"."+PsqlDatabaseDataSourceByName, databaseNameAttribute, PsqlDatabaseResource+"."+PsqlDatabaseTestResource, databaseNameAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlDatabaseResource+"."+PsqlDatabaseDataSourceByName, databaseOwnerAttribute, PsqlDatabaseResource+"."+PsqlDatabaseTestResource, databaseOwnerAttribute),
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
					resource.TestCheckResourceAttr(DataSource+"."+PsqlDatabasesResource+"."+PsqlDatabasesDataSource, databasesAttribute+".#", "4"),
				),
			},
			{
				Config: PgSqlAllDatabasesFilterByOwnerDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+PsqlDatabasesResource+"."+PsqlDatabasesDataSource, databasesAttribute+".#", "1"),
					resource.TestCheckResourceAttr(DataSource+"."+PsqlDatabasesResource+"."+PsqlDatabasesDataSource, databasesAttribute+".0.name", databaseNameValue),
					resource.TestCheckResourceAttr(DataSource+"."+PsqlDatabasesResource+"."+PsqlDatabasesDataSource, databasesAttribute+".0.owner", databaseOwnerValue),
				),
			},
		},
	})
}

func pgSqlDatabaseExistsCheck(path string, database *pgsql.DatabaseResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).PsqlClient
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
			return fmt.Errorf("error occured while fetching the PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err)
		}
		database = &foundDatabase
		return nil
	}
}

func pgSqlDatabaseDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).PsqlClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != PsqlDatabaseResource {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		_, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occured while checking the deletion of PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err)
			}
		} else {
			return fmt.Errorf("PgSql database %s still exists in the cluster with ID: %s", name, clusterId)
		}
	}
	return nil
}

// Configurations

const PgSqlDatabaseConfig = PgSqlUserConfig + `
resource ` + PsqlDatabaseResource + ` ` + PsqlDatabaseTestResource + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id  
  ` + databaseNameAttribute + ` = "` + databaseNameValue + `"
  ` + databaseOwnerAttribute + ` = ` + PsqlUserResource + `.` + UserTestResource + `.username
}
`

const PgSqlDatabaseDataSource = PgSqlDatabaseConfig + `
data ` + PsqlDatabaseResource + ` ` + PsqlDatabaseDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id   
  ` + databaseNameAttribute + ` = ` + PsqlDatabaseResource + `.` + PsqlDatabaseTestResource + `.name
}
`

const PgSqlDatabaseDataSourceWrongName = PgSqlDatabaseConfig + `
data ` + PsqlDatabaseResource + ` ` + PsqlDatabaseDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id   
  ` + databaseNameAttribute + ` = "nonexistent"
}
`

const PgSqlAllDatabasesDataSource = PgSqlDatabaseConfig + `
data ` + PsqlDatabasesResource + ` ` + PsqlDatabasesDataSource + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id   
}
`

const PgSqlAllDatabasesFilterByOwnerDataSource = PgSqlDatabaseConfig + `
data ` + PsqlDatabasesResource + ` ` + PsqlDatabasesDataSource + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id 
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
