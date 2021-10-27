package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewallImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigBasic,
			},

			{
				ResourceName:      FirewallResource + "." + FirewallTestResource,
				ImportStateIdFunc: testAccFirewallImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirewallImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != FirewallResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
