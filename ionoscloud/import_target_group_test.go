//go:build all || alb

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTargetGroupImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckTargetGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTargetGroupConfigBasic,
			},

			{
				ResourceName:      resourceNameTargetGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
