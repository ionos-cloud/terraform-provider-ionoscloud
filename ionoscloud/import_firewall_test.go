//go:build compute || all || firewall

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFirewallImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigBasic,
			},

			{
				ResourceName:      constant.FirewallResource + "." + constant.FirewallTestResource,
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
		if rs.Type != constant.FirewallResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
