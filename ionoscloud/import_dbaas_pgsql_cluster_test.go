//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDbaasPgSqlClusterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
			},
			{
				ResourceName:            DBaaSClusterResource + "." + DBaaSClusterTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections,
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigCleanup,
			},
		},
	})
}
