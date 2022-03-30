//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDbaasPgSqlClusterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasic,
			},
			{
				ResourceName:            DBaaSClusterResource + "." + DBaaSClusterTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasicRemoveCluster,
			},
		},
	})
}

const testAccCheckDbaasPgSqlClusterConfigBasicRemoveCluster = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}
`
