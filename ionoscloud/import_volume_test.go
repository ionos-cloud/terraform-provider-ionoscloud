package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolume_ImportBasic(t *testing.T) {
	resourceName := "database_volume"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVolumeDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckVolumeConfigBasic, resourceName),
			},

			{
				ResourceName:            fmt.Sprintf("ionoscloud_volume.%s", resourceName),
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
		if rs.Type != "ionoscloud_volume" {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
