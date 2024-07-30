//go:build all || dbaas || inMemoryDB
// +build all dbaas inMemoryDB

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDBaaSInMemoryDBReplicaSetImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDBaaSInMemoryDBReplicaSetDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: inMemoryDBReplicaSetConfigHashedPassword,
			},
			{
				ResourceName:            constant.DBaaSInMemoryDBReplicaSetResource + "." + constant.DBaaSReplicaSetTestResource,
				ImportState:             true,
				ImportStateIdFunc:       testAccDBaaSInMemoryDBReplicaSetImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}

func testAccDBaaSInMemoryDBReplicaSetImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaaSInMemoryDBReplicaSetResource {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
