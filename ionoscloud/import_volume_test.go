//go:build compute || all || volume

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolumeImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeConfigBasic,
			},

			{
				ResourceName:            VolumeResource + "." + VolumeTestResource,
				ImportStateIdFunc:       testAccVolumeImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "user_data", "backup_unit_id"},
			},
		},
	})
}

func testAccVolumeImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != VolumeResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
