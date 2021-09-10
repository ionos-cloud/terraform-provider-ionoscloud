package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDbaasCluster_ImportBasic(t *testing.T) {
	resourceName := "ionoscloud_dbaas_cluster.test_dbaas_cluster"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDbaasClusterConfigBasic),
			},

			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDbaasClusterImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDbaasClusterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_dbaas_cluster" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
