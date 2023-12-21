//go:build all || k8s
// +build all k8s

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

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
				ResourceName:            constant.K8sClusterResource + "." + constant.K8sClusterTestResource,
				ImportStateIdFunc:       testAccK8sClusterImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_replace"},
			},
		},
	})
}

func testAccK8sClusterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.K8sClusterResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
