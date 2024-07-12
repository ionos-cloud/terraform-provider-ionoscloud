//go:build all || dbaas || mongo
// +build all dbaas mongo

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongoUserImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckMongoUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMongoUserConfigBasic,
			},
			{
				ResourceName:            constant.DBaasMongoUserResource + "." + constant.UserTestResource,
				ImportStateIdFunc:       MongoUserImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func MongoUserImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaasMongoUserResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["username"])
	}

	return importID, nil
}
