package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServer_ImportBasic(t *testing.T) {
	resourceName := "server-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckServerConfigBasic, resourceName),
			},

			{
				ResourceName:            "ionoscloud_server.webserver",
				ImportStateIdFunc:       testAccServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image"},
			},
		},
	})
}

func testAccServerImportStateId(s *terraform.State) (string, error) {
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_server" {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"], rs.Primary.Attributes["primary_nic"], rs.Primary.Attributes["firewallrule_id"])
	}

	return importID, nil
}
