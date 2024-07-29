//go:build compute || all || snapshot

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSnapshotConfigBasic,
			},

			{
				ResourceName:            constant.SnapshotResource + "." + constant.SnapshotTestResource,
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
		if rs.Type != constant.SnapshotResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
