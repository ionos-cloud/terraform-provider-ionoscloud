package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				Config: fmt.Sprintf(testAccCheckDatacenterConfigBasic, resourceName),
			},

			{
				ResourceName:      fmt.Sprintf("ionoscloud_datacenter.foobar"),
				ImportStateIdFunc: testAccDatacenterImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDatacenterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		importID = fmt.Sprintf("%s", rs.Primary.Attributes["id"])
	}

	return importID, nil
}
