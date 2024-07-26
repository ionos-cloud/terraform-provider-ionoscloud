//go:build compute || all || cert

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"fmt"
	"testing"
)

func TestAccCMAutoCertificateImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCMAutoCertificateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: CMAutoCertificateConfig,
			},
			{
				ResourceName:      constant.AutoCertificateResource + "." + constant.TestCMAutoCertificateName,
				ImportStateIdFunc: testAccCMAutoCertificateImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCMAutoCertificateImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.AutoCertificateResource {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
