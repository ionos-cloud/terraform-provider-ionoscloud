//go:build nfs || all || nfs_cluster

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNFSCluster_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNFSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNFSClusterConfig_basic,
			},
			{
				ResourceName:      "ionoscloud_nfs_cluster.example",
				ImportStateIdFunc: testAccNFSClusterImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNFSClusterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_nfs_cluster" {
			continue
		}
		fmt.Println(rs.Type)

		importID = rs.Primary.Attributes["id"]
	}

	if importID == "" {
		return "", fmt.Errorf("Could not find nfs cluster")
	}

	return importID, nil
}
