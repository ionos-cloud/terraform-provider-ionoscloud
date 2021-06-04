// +build k8s

package ionoscloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAcck8sNodepool_ImportBasic(t *testing.T) {
	resourceName := "terraform_acctest"
	publicIp1 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_1")
	if publicIp1 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_1 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}
	publicIp2 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_2")
	if publicIp2 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_2 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigBasic, resourceName, publicIp1, publicIp2),
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
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_node_pool" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)
	}

	return importID, nil
}
