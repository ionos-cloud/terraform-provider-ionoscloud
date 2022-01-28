//go:build compute || all || ipblock

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPBlockImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckIPBlockDestroyCheck,
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
