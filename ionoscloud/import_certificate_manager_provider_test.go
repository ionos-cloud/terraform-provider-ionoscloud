//go:build compute || all || cert

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"fmt"
	"testing"
)

func TestAccCMProviderImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCMProviderDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: CMProviderConfig,
			},
			{
				ResourceName:            constant.AutoCertificateProviderResource + "." + constant.TestCMProviderName,
				ImportStateIdFunc:       testAccCMProviderImportStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_account_binding"},
			},
		},
	})
}

func testAccCMProviderImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.AutoCertificateProviderResource {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
