//go:build k8s
// +build k8s

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcck8sClusterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckK8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sClusterConfigBasic,
			},
			{
				ResourceName:      K8sClusterResource + "." + K8sClusterTestResource,
				ImportStateIdFunc: testAccK8sClusterImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccK8sClusterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != K8sClusterResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
