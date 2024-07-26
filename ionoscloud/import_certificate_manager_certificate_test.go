//go:build compute || all || cert

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"
)

func TestAccCertificateImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCertificateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCertConfigBasic,
			},
			{
				ResourceName:            constant.CertificateResource + "." + constant.TestCertName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}
