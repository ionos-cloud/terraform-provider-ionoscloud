//go:build compute || all || group

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroupImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccImportGroupConfigBasic,
			},

			{
				ResourceName:            GroupResource + "." + GroupTestResource,
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
		if rs.Type != GroupResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}

var testAccImportGroupConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  user_ids = [` + UserResource + `.` + UserTestResource + `.id]
}
`
