//go:build compute || all || server

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServerImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBasic,
			},
			{
				ResourceName:            ServerResource + "." + ServerTestResource,
				ImportStateIdFunc:       testAccServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic"},
			},
		},
	})
}

func TestAccServerWithLabelsImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: randomProviderVersion343(),
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerCreationWithLabels,
			},
			{
				ResourceName:            ServerResource + "." + ServerTestResource,
				ImportStateIdFunc:       testAccServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic"},
			},
		},
	})
}

func testAccServerImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ServerResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
