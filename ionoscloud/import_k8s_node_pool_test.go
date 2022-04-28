//go:build all || k8s
// +build all k8s

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccK8sNodePoolImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigBasic,
			},
			{
				ResourceName:            ResourceNameK8sNodePool,
				ImportStateIdFunc:       testAccK8sNodePoolImportStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"maintenance_window.0.time", "allow_replace"},
			},
		},
	})
}

func testAccK8sNodePoolImportStateID(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != K8sNodePoolResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)
	}
	return importID, nil
}
