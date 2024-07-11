//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDataSourceAutoscalingGroupServers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		ExternalProviders:        randomProviderVersion343(),
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_basic(constant.AutoscalingGroupTestResource),
			},
			{
				Config: testAccDataSourceAutoscalingGroupServers,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.AutoscalingGroupServersResource+"."+constant.AutoscalingGroupServersTestDataSource, "servers.#"),
				),
			},
		},
	})

}

var testAccDataSourceAutoscalingGroupServers = testAGConfig_basic(constant.AutoscalingGroupTestResource) + `
data ` + constant.AutoscalingGroupServersResource + ` ` + constant.AutoscalingGroupServersTestDataSource + ` {
  group_id = ` + constant.AutoscalingGroupResource + `.` + constant.AutoscalingGroupTestResource + `.id
}
`
