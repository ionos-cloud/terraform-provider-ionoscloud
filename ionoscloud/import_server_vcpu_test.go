//go:build compute || all || server || vcpu
// +build compute all server vcpu

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccServerVCPUImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServerVCPUDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUConfigBasic,
			},
			{
				ResourceName:            constant.ServerVCPUResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccServerVCPUImportStateIdWithNicAndFw,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids"},
			},
		},
	})
}

func TestAccServerVCPUWithLabelsImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		ExternalProviders:        randomProviderVersion343(),
		CheckDestroy:             testAccCheckServerVCPUDestroyCheck,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerVCPUCreationWithLabels,
			},
			{
				ResourceName:            constant.ServerVCPUResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccServerVCPUImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids", "primary_ip"},
			},
		},
	})
}
func testAccServerVCPUImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerVCPUResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
		if rs.Primary.Attributes["primary_nic"] != "" {
			importID = fmt.Sprintf("%s/%s", importID, rs.Primary.Attributes["primary_nic"])
		}
	}

	return importID, nil
}
func testAccServerVCPUImportStateIdWithNicAndFw(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerVCPUResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
		// we might get the primary nic id and the primary firewall id here as import optionals
		if nicID, ok := rs.Primary.Attributes["primary_nic"]; ok {
			importID += "/" + nicID
			if primaryFwID, ok := rs.Primary.Attributes["firewallrule_id"]; ok {
				importID += "/" + primaryFwID
			}
		}

	}

	return importID, nil
}
