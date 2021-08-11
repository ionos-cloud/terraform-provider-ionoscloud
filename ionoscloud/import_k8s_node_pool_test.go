// +build k8s

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAcck8sNodepool_ImportBasic(t *testing.T) {
	resourceName := "terraform_acctest"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigBasic, resourceName),
			},
			{
				ResourceName:            fmt.Sprintf("ionoscloud_k8s_node_pool.%s", resourceName),
				ImportStateIdFunc:       testAcck8sNodepoolImportStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"maintenance_window.0.time"},
			},
		},
	})
}

func testAcck8sNodepoolImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_node_pool" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)
	}

	return importID, nil
}
