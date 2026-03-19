//go:build compute || all || server || cube
// +build compute all server cube

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCubeServerImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerConfigBasic,
			},
			{
				ResourceName:            constant.ServerCubeResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccCubeServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids", "allow_replace", "location", "boot_volume"},
			},
		},
	})
}

func TestAccCubeServerImportWithIPv6(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCubeServerEnableIpv6,
			},
			{
				ResourceName:            constant.ServerCubeResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccCubeServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids", "allow_replace", "location", "boot_volume"},
			},
		},
	})
}

func testAccCubeServerImportStateId(s *terraform.State) (string, error) {
	var importID string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerCubeResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
