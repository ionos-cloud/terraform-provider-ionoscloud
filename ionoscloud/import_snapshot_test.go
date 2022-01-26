//go:build compute || all || snapshot

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSnapshotImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSnapshotConfigBasic,
			},

			{
				ResourceName:            SnapshotResource + "." + SnapshotTestResource,
				ImportStateIdFunc:       testAccSnapshotImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"datacenter_id", "volume_id"},
			},
		},
	})
}

func testAccSnapshotImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != SnapshotResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
