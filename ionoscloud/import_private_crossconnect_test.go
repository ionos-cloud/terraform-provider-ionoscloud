package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccprivateCrossConnect_ImportBasic(t *testing.T) {
	resourceName := "example"
	resourceDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckprivateCrossConnectConfigBasic, resourceName, resourceDescription),
			},
			{
				ResourceName:      fmt.Sprintf("ionoscloud_private_crossconnect.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
