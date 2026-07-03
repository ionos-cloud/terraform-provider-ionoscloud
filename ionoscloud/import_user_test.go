//go:build compute || all || user

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccImportUserConfigBasic,
			},

			{
				ResourceName:            constant.UserResource + "." + constant.UserTestResource,
				ImportStateIdFunc:       testAccUserImportStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserImportStateID(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.UserResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}

var testAccImportUserConfigBasic = `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "` + constant.UserTestResource + `"
  last_name = "` + constant.UserTestResource + `"
  email = "` + utils.GenerateEmail() + `"
  password = random_password.user_password.result
  force_sec_auth= true
  active  = true
  group_ids 		= [ ionoscloud_group.group1.id]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}

resource "random_password" "user_password" {
  length  = 16
  special = true
}
`
