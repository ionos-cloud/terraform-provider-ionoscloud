//go:build all || dbaas || psql
// +build all dbaas psql

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPgSqlDatabaseImport(t *testing.T) {
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
			},
			{
				ResourceName:      constant.PsqlDatabaseResource + "." + constant.PsqlDatabaseTestResource,
				ImportStateIdFunc: PgSqlDatabaseImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func PgSqlDatabaseImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PsqlDatabaseResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["name"])
	}

	return importID, nil
}
