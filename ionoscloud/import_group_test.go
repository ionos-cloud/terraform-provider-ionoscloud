package ionoscloud

import (
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
				Config: testAccCheckGroupConfigBasic,
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
