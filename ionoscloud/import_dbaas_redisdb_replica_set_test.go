//go:build all || dbaas || redis
// +build all dbaas redis

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"testing"
)

func TestAccDBaaSRedisDBReplicaSetImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDBaaSRedisDBReplicaSetDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: redisDBReplicaSetConfigHashedPassword,
			},
			{
				ResourceName:            constant.DBaaSRedisDBReplicaSetResource + "." + constant.DBaaSReplicaSetTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}
