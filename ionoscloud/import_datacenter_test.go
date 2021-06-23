package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataCenter_ImportBasic(t *testing.T) {
	resourceName := "datacenter-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDatacenterConfig_basic, resourceName),
			},

			{
				ResourceName:      fmt.Sprintf("ionoscloud_datacenter.foobar"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
