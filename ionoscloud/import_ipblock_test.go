//go:build compute || all || ipblock

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPBlockImportBasic(t *testing.T) {
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIPBlockConfigBasic, location),
			},

			{
				ResourceName:      IpBLockResource + ".webserver_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
