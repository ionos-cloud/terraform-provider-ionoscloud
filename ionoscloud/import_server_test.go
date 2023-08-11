//go:build compute || all || server || enterprise
// +build compute all server enterprise

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServerImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBasic,
			},
			{
				ResourceName:            constant.ServerResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccServerImportStateIdWithNicAndFw,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids"},
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
				ResourceName:            constant.ServerResource + "." + constant.ServerTestResource,
				ImportStateIdFunc:       testAccServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name", "volume.0.user_data", "volume.0.backup_unit_id", "firewallrule_id", "primary_nic", "inline_volume_ids", "primary_ip"},
			},
		},
	})
}
func testAccServerImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
		if rs.Primary.Attributes["primary_nic"] != "" {
			importID = fmt.Sprintf("%s/%s", importID, rs.Primary.Attributes["primary_nic"])
		}
	}

	return importID, nil
}
func testAccServerImportStateIdWithNicAndFw(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
		//we might get the primary nic id and the primary firewall id here as import optionals
		if nicID, ok := rs.Primary.Attributes["primary_nic"]; ok {
			importID += "/" + nicID
			if primaryFwID, ok := rs.Primary.Attributes["firewallrule_id"]; ok {
				importID += "/" + primaryFwID
			}
		}

	}

	return importID, nil
}

const testAccCheckServerImport = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = random_password.image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = false
  }
  label {
    key = "labelkey0"
    value = "labelvalue0"
  }
  label {
    key = "labelkey1"
    value = "labelvalue1"
  }
}`
