//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDataSourceAutoscalingGroupServers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: randomProviderVersion343(),
		Steps: []resource.TestStep{
			{
				Config: testAG_ConfigBasic,
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

const testAccDataSourceAutoscalingGroupServers = testAG_ConfigBasic + `
data ` + constant.AutoscalingGroupServersResource + ` ` + constant.AutoscalingGroupServersTestDataSource + ` {
  group_id = ` + constant.AutoscalingGroupResource + `.` + constant.AutoscalingGroupTestResource + `.id
}
`
