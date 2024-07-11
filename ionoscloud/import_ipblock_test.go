//go:build compute || all || ipblock

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPBlockImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPBlockConfigBasic,
			},
			{
				ResourceName:      fullIpBlockResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
