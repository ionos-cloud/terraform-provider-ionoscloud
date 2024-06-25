//go:build all || autoscaling
// +build all autoscaling

package asg

import (
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const dataSourceAutoscalingGroupId = constant.DataSource + "." + constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupDataSourceById
const dataSourceAutoscalingGroupName = constant.DataSource + "." + constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupDataSourceByName

func TestAccDataSourceAutoscalingGroup(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			ionoscloud.testAccPreCheck(t)
		},
		ProviderFactories: ionoscloud.testAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.11.1",
			},
		},
		CheckDestroy: testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfigDataSourceBasic(),
			},
			{
				Config: testAccDataSourceAutoscalingGroupMatchId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "name", resourceAGName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "datacenter_id", resourceAGName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "location", resourceAGName, "location"),
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
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_action.0.delete_volumes", resourceAGName, "policy.0.scale_out_action.0.delete_volumes"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.scale_out_threshold", resourceAGName, "policy.0.scale_out_threshold"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "policy.0.unit", resourceAGName, "policy.0.unit"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.availability_zone", resourceAGName, "replica_configuration.0.availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cores", resourceAGName, "replica_configuration.0.cores"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.cpu_family", resourceAGName, "replica_configuration.0.cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.lan", resourceAGName, "replica_configuration.0.nic.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.name", resourceAGName, "replica_configuration.0.nic.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.dhcp", resourceAGName, "replica_configuration.0.nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_active", resourceAGName, "replica_configuration.0.nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_type", resourceAGName, "replica_configuration.0.nic.0.firewall_type"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.name", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.name"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.bucket", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.bucket"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.action", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.action"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.direction", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.direction"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_rule.0.name", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_rule.0.protocol", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_rule.0.port_range_start", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_start"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_rule.0.port_range_end", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_end"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.firewall_rule.0.type", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.target_group.0.target_group_id", resourceAGName, "replica_configuration.0.nic.0.target_group.0.target_group_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.target_group.0.weight", resourceAGName, "replica_configuration.0.nic.0.target_group.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupId, "replica_configuration.0.nic.0.target_group.0.port", resourceAGName, "replica_configuration.0.nic.0.target_group.0.port"),
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
				Config: testAccDataSourceAutoscalingGroupMatchName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "name", resourceAGName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "datacenter_id", resourceAGName, "datacenter_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "max_replica_count", resourceAGName, "max_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "min_replica_count", resourceAGName, "min_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "target_replica_count", resourceAGName, "target_replica_count"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "location", resourceAGName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "policy.0.scale_out_action.0.delete_volumes", resourceAGName, "policy.0.scale_out_action.0.delete_volumes"),
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
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.lan", resourceAGName, "replica_configuration.0.nic.0.lan"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.name", resourceAGName, "replica_configuration.0.nic.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.dhcp", resourceAGName, "replica_configuration.0.nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_active", resourceAGName, "replica_configuration.0.nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_type", resourceAGName, "replica_configuration.0.nic.0.firewall_type"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.name", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.name"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.bucket", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.bucket"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.action", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.action"),
					//resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.flow_log.0.direction", resourceAGName, "replica_configuration.0.nic.0.flow_log.0.direction"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_rule.0.name", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_rule.0.protocol", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_start", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_start"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_end", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_end"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.firewall_rule.0.type", resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.target_group.0.target_group_id", resourceAGName, "replica_configuration.0.nic.0.target_group.0.target_group_id"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.target_group.0.weight", resourceAGName, "replica_configuration.0.nic.0.target_group.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceAutoscalingGroupName, "replica_configuration.0.nic.0.target_group.0.port", resourceAGName, "replica_configuration.0.nic.0.target_group.0.port"),
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

func testAGConfigDataSourceBasic() string {
	return fmt.Sprintf(`
resource "ionoscloud_datacenter" "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fra"
}
resource "ionoscloud_lan" "autoscaling_lan_1" {
  	datacenter_id    = ionoscloud_datacenter.autoscaling_datacenter.id
  	public           = false
	name             = "test_autoscaling_group_1"
}

resource "ionoscloud_target_group" "autoscaling_target_group" {
    name                      = "Target Group Example" 
    algorithm                 = "ROUND_ROBIN"
    protocol                  = "HTTP"
}

// This resource is used to wait for the target group to be destroyed after clearing all resources from autoscaling group
resource "time_sleep" "wait_5_minutes" {
	depends_on = [ionoscloud_target_group.autoscaling_target_group]
	destroy_duration = "5m"
}

resource "ionoscloud_autoscaling_group"  %[1]q {
  datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
  max_replica_count      = 5
  min_replica_count      = 1
  name           = %[1]q
  policy {
    metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
    range              = "PT24H"
    scale_in_action {
      amount                  =  1
      amount_type             = "ABSOLUTE"
      termination_policy_type = "OLDEST_SERVER_FIRST"
      cooldown_period         = "PT5M"
      delete_volumes = true
    }
    scale_in_threshold = 33
    scale_out_action  {
      amount          =  1
      amount_type     = "ABSOLUTE"
      cooldown_period = "PT5M"
    }
    scale_out_threshold = 77
    unit                = "PER_HOUR"
  }
  replica_configuration {
    availability_zone = "AUTO"
    cores             = "2"
    cpu_family        = "INTEL_SKYLAKE"
 	ram          = 2048
    nic {
      lan       = ionoscloud_lan.autoscaling_lan_1.id
      name      = "nic_1"
      dhcp      = true
      firewall_active = true
	  firewall_type = "INGRESS"
	  firewall_rule {
		name = "rule_1"
		protocol = "TCP"
		port_range_start = 1
		port_range_end = 1000
		type = "INGRESS"
      }

      //flow_log {
		//name="flow_log_1"
	  //  bucket="test-de-bucket"
		//action="ALL"
		//direction="BIDIRECTIONAL"
      //}

      target_group {
        target_group_id = ionoscloud_target_group.autoscaling_target_group.id
        port            = 80
        weight          = 50
      }
    }
   
    volume {
      image_alias = "ubuntu:latest"
      name        = "volume_1"
      image_password = random_password.image_password.result
      size        = 30
      ssh_keys    = ["`+ionoscloud.sshKey+`"]
      type        = "HDD"
      user_data    = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
      boot_order = "AUTO"
      bus = "IDE"
    }
  }

  depends_on = [time_sleep.wait_5_minutes]	
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
`, constant.AutoscalingGroupTestResource)
}

func testAccDataSourceAutoscalingGroupMatchId() string {
	return utils.ConfigCompose(testAGConfigDataSourceBasic(), fmt.Sprintf(`
data %[1]q %[2]q {
	  id = %[1]s.%[3]s.id
}
`, constant.AutoscalingGroupResource, constant.AutoscalingGroupDataSourceById, constant.AutoscalingGroupTestResource))
}

func testAccDataSourceAutoscalingGroupMatchName() string {
	return utils.ConfigCompose(testAGConfigDataSourceBasic(), fmt.Sprintf(`
data %[1]q %[2]q {
	  name = %[1]s.%[3]s.name
}
`, constant.AutoscalingGroupResource, constant.AutoscalingGroupDataSourceByName, constant.AutoscalingGroupTestResource))
}
