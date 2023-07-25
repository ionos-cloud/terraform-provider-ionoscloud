//go:build compute || all || group

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroupImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccImportGroupConfigBasic,
			},

			{
				ResourceName:            constant.GroupResource + "." + constant.GroupTestResource,
				ImportStateIdFunc:       testAccGroupImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_id"},
			},
		},
	})
}

func testAccGroupImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.GroupResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}

var testAccImportGroupConfigBasic = `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id]
}
`
