package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPBlock_ImportBasic(t *testing.T) {
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIPBlockConfig_basic, location),
			},

			{
				ResourceName:      "ionoscloud_ipblock.webserver_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
