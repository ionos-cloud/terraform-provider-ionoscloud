package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAutoscalingGroupServers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAutoscalingGroupConfigBasic,
			},
			{
				Config: testAccDataSourceAutoscalingGroupServers,
				Check: resource.ComposeTestCheckFunc(
					testNotEmptySlice(AutoscalingGroupServersResource+"."+AutoscalingGroupServersTestDataSource, "servers.#"),
				),
			},
		},
	})

}

const testAccDataSourceAutoscalingGroupServers = testAccCheckAutoscalingGroupConfigBasic + `
data ` + AutoscalingGroupServersResource + ` ` + AutoscalingGroupServersTestDataSource + ` {
	group_id = ` + AutoscalingGroupResource + `.` + AutoscalingGroupTestResource + `.id
}
`
