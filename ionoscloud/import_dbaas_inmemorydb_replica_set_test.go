//go:build all || dbaas || inMemoryDB
// +build all dbaas inMemoryDB

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDBaaSInMemoryDBReplicaSetImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDBaaSInMemoryDBReplicaSetDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: inMemoryDBReplicaSetConfigHashedPassword,
			},
			{
				ResourceName:            constant.DBaaSInMemoryDBReplicaSetResource + "." + constant.DBaaSReplicaSetTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}
