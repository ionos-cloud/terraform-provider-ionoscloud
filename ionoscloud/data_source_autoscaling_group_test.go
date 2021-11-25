package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceAutoscalingGroupId = DataSource + "." + AutoscalingGroupResource + "." + AutoscalingGroupDataSourceById
const dataSourceAutoscalingGroupName = DataSource + "." + AutoscalingGroupResource + "." + AutoscalingGroupDataSourceByName

func TestAccDataSourceAutoscalingGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigBasic, AutoscalingGroupTestResource),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceAutoscalingGroupMatchId, AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "name", resourceAutoscalingGroupName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "datacenter_id", resourceAutoscalingGroupName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "max_replica_count", resourceAutoscalingGroupName, "max_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "min_replica_count", resourceAutoscalingGroupName, "min_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "target_replica_count", resourceAutoscalingGroupName, "target_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.metric", resourceAutoscalingGroupName, "policy.0.metric"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.range", resourceAutoscalingGroupName, "policy.0.range"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.amount", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.amount_type", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.termination_policy_type", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_action.0.cooldown_period", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_in_threshold", resourceAutoscalingGroupName, "policy.0.scale_in_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.amount", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.amount_type", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.cooldown_period", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_threshold", resourceAutoscalingGroupName, "policy.0.scale_out_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.unit", resourceAutoscalingGroupName, "policy.0.unit"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.availability_zone", resourceAutoscalingGroupName, "replica_configuration.0.availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cores", resourceAutoscalingGroupName, "replica_configuration.0.cores"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cpu_family", resourceAutoscalingGroupName, "replica_configuration.0.cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.lan", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.name", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nics.0.dhcp", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.ram", resourceAutoscalingGroupName, "replica_configuration.0.ram"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.image", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.name", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.size", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.ssh_keys", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.type", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.volumes.0.user_data", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceAutoscalingGroupMatchName, AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "name", resourceAutoscalingGroupName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "datacenter_id", resourceAutoscalingGroupName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "max_replica_count", resourceAutoscalingGroupName, "max_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "min_replica_count", resourceAutoscalingGroupName, "min_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "target_replica_count", resourceAutoscalingGroupName, "target_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.metric", resourceAutoscalingGroupName, "policy.0.metric"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.range", resourceAutoscalingGroupName, "policy.0.range"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_in_threshold", resourceAutoscalingGroupName, "policy.0.scale_in_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_threshold", resourceAutoscalingGroupName, "policy.0.scale_out_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.unit", resourceAutoscalingGroupName, "policy.0.unit"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.availability_zone", resourceAutoscalingGroupName, "replica_configuration.0.availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.cores", resourceAutoscalingGroupName, "replica_configuration.0.cores"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.cpu_family", resourceAutoscalingGroupName, "replica_configuration.0.cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.ram", resourceAutoscalingGroupName, "replica_configuration.0.ram"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.ssh_keys", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data", resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data"),
				),
			},
		},
	})
}

const testAccDataSourceAutoscalingGroupMatchId = testAccCheckAutoscalingGroupConfigBasic + `
data ` + AutoscalingGroupResource + ` ` + AutoscalingGroupDataSourceById + ` {
  id			= ` + resourceAutoscalingGroupName + `.id
}
`

const testAccDataSourceAutoscalingGroupMatchName = testAccCheckAutoscalingGroupConfigBasic + `
data ` + AutoscalingGroupResource + ` ` + AutoscalingGroupDataSourceByName + ` {
  name			= ` + resourceAutoscalingGroupName + `.name
}
`
