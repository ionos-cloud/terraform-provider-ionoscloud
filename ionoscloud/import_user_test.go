//go:build compute || all || user

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccImportUserConfigBasic,
			},

			{
				ResourceName:            UserResource + "." + UserTestResource,
				ImportStateIdFunc:       testAccUserImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != UserResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}

var testAccImportUserConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "` + UserTestResource + `"
  last_name = "` + UserTestResource + `"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = true
  force_sec_auth= true
  active  = true
}`
