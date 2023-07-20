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

func TestAccPgSqlUser(t *testing.T) {
	var user pgsql.UserResource

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      pgSqlUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: PgSqlUserConfig,
				Check: resource.ComposeTestCheckFunc(
					pgSqlUserExistsCheck(PsqlUserResource+"."+UserTestResource, &user),
					resource.TestCheckResourceAttr(PsqlUserResource+"."+UserTestResource, usernameAttribute, usernameValue),
					resource.TestCheckResourceAttrSet(PsqlUserResource+"."+UserTestResource, passwordAttribute),
					resource.TestCheckResourceAttr(PsqlUserResource+"."+UserTestResource, isSystemUserAttribute, isSystemUserValue),
				),
			},
			{
				Config: PgSqlUserDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlUserResource+"."+UserDataSourceByName, usernameAttribute, PsqlUserResource+"."+UserTestResource, usernameAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlUserResource+"."+UserDataSourceByName, isSystemUserAttribute, PsqlUserResource+"."+UserTestResource, isSystemUserAttribute),
				),
			},
			{
				Config:      PgSqlUserDataSourceWrongUsername,
				ExpectError: regexp.MustCompile(`no PgSql user found with the specified username`),
			},
		},
	})
}

func pgSqlUserExistsCheck(path string, user *pgsql.UserResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).PsqlClient
		rs, ok := s.RootModule().Resources[path]
		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the PgSql user")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		clusterId := rs.Primary.Attributes["cluster_id"]
		username := rs.Primary.Attributes["username"]
		foundUser, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
		apiResponse.LogInfo()
		if err != nil {
			return fmt.Errorf("error occured while fetching the PgSql user: %s, cluster ID: %s, error: %w", username, clusterId, err)
		}
		user = &foundUser
		return nil
	}
}

func pgSqlUserDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).PsqlClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != PsqlUserResource {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		username := rs.Primary.Attributes["username"]
		_, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occured while checking the deletion of PgSql username: %s, cluster ID: %s, error: %w", username, clusterId, err)
			}
		} else {
			return fmt.Errorf("PgSql user %s still exists in the cluster with ID: %s", username, clusterId)
		}
	}
	return nil
}

const PgSqlUserDataSource = PgSqlUserConfig + `
data ` + PsqlUserResource + ` ` + UserDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id  
  username = ` + PsqlUserResource + `.` + UserTestResource + `.username
}
`

const PgSqlUserDataSourceWrongUsername = PgSqlUserConfig + `
data ` + PsqlUserResource + ` ` + UserDataSourceByName + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id  
  username = "nonexistent"
}
`
