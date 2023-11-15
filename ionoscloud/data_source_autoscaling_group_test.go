//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const dataSourceAutoscalingGroupId = constant.DataSource + "." + constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupDataSourceById
const dataSourceAutoscalingGroupName = constant.DataSource + "." + constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupDataSourceByName

func TestAccDataSourceAutoscalingGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAG_ConfigBasic,
			},
			{
				Config: testAccDataSourceAutoscalingGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "name", resourceAGName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "datacenter_id", resourceAGName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "max_replica_count", resourceAGName, "max_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "min_replica_count", resourceAGName, "min_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "target_replica_count", resourceAGName, "target_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.metric", resourceAGName, "policy.0.metric"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.range", resourceAGName, "policy.0.range"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.amount", resourceAGName, "policy.0.scale_in_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.amount_type", resourceAGName, "policy.0.scale_in_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.termination_policy_type", resourceAGName, "policy.0.scale_in_action.0.termination_policy_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.cooldown_period", resourceAGName, "policy.0.scale_in_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_threshold", resourceAGName, "policy.0.scale_in_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.amount", resourceAGName, "policy.0.scale_out_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.amount_type", resourceAGName, "policy.0.scale_out_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.cooldown_period", resourceAGName, "policy.0.scale_out_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_threshold", resourceAGName, "policy.0.scale_out_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.unit", resourceAGName, "policy.0.unit"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.availability_zone", resourceAGName, "replica_configuration.0.availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cores", resourceAGName, "replica_configuration.0.cores"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cpu_family", resourceAGName, "replica_configuration.0.cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.lan", resourceAGName, "replica_configuration.0.nics.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.name", resourceAGName, "replica_configuration.0.nics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.dhcp", resourceAGName, "replica_configuration.0.nics.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.ram", resourceAGName, "replica_configuration.0.ram"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.image", resourceAGName, "replica_configuration.0.volumes.0.image"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.name", resourceAGName, "replica_configuration.0.volumes.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.size", resourceAGName, "replica_configuration.0.volumes.0.size"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.ssh_keys", resourceAGName, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.type", resourceAGName, "replica_configuration.0.volumes.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.user_data", resourceAGName, "replica_configuration.0.volumes.0.user_data"),
				),
			},
			{
				Config: testAccDataSourceAutoscalingGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "name", resourceAGName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "datacenter_id", resourceAGName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "max_replica_count", resourceAGName, "max_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "min_replica_count", resourceAGName, "min_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "target_replica_count", resourceAGName, "target_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.metric", resourceAGName, "policy.0.metric"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.range", resourceAGName, "policy.0.range"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", resourceAGName, "policy.0.scale_in_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", resourceAGName, "policy.0.scale_in_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", resourceAGName, "policy.0.scale_in_action.0.termination_policy_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", resourceAGName, "policy.0.scale_in_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_threshold", resourceAGName, "policy.0.scale_in_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", resourceAGName, "policy.0.scale_out_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", resourceAGName, "policy.0.scale_out_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", resourceAGName, "policy.0.scale_out_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_threshold", resourceAGName, "policy.0.scale_out_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.unit", resourceAGName, "policy.0.unit"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.availability_zone", resourceAGName, "replica_configuration.0.availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.cores", resourceAGName, "replica_configuration.0.cores"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.cpu_family", resourceAGName, "replica_configuration.0.cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", resourceAGName, "replica_configuration.0.nics.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", resourceAGName, "replica_configuration.0.nics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", resourceAGName, "replica_configuration.0.nics.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.ram", resourceAGName, "replica_configuration.0.ram"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image", resourceAGName, "replica_configuration.0.volumes.0.image"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", resourceAGName, "replica_configuration.0.volumes.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", resourceAGName, "replica_configuration.0.volumes.0.size"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.ssh_keys", resourceAGName, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", resourceAGName, "replica_configuration.0.volumes.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data", resourceAGName, "replica_configuration.0.volumes.0.user_data"),
				),
			},
		},
	})
}

const testAccDataSourceAutoscalingGroupMatchId = testAG_ConfigBasic + `
data ` + constant.AutoscalingGroupResource + ` ` + constant.AutoscalingGroupDataSourceById + ` {
  id = ` + resourceAGName + `.id
}
`

const testAccDataSourceAutoscalingGroupMatchName = testAG_ConfigBasic + `
data ` + constant.AutoscalingGroupResource + ` ` + constant.AutoscalingGroupDataSourceByName + ` {
  name = ` + resourceAGName + `.name
}
`
