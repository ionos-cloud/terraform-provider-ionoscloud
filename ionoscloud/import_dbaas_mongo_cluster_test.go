//go:build all || dbaas || mongo
// +build all dbaas mongo

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbaasMongoClusterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDbaasMongoClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasMongoClusterConfigBasic,
			},
			{
				ResourceName:            constant.DBaasMongoClusterResource + "." + constant.DBaaSClusterTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}
