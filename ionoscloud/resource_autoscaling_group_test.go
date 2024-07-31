//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const resourceAGName = constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupTestResource

func TestAccAutoscalingGroup_basic(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_basic(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", constant.AutoscalingGroupTestResource),
					resource.TestCheckResourceAttrPair(resourceAGName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttrPair(resourceAGName, "location", constant.DatacenterResource+".autoscaling_datacenter", "location"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "5"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_CPU_UTILIZATION_AVERAGE"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "77"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_HOUR"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "2048"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_requiredUpdated(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_requiredUpdated(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", "test_autoscaling_group_new"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "4"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_PERCENTAGE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.delete_volumes", "false"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "34"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "10"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_PERCENTAGE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "78"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "1024"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_policyWithOptionals(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_policyWithOptionals(constant.AutoscalingGroupTestResource, "PT24H", "PT24H", "RANDOM", "PT24H"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.cooldown_period", "PT24H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.termination_policy_type", "RANDOM"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.cooldown_period", "PT24H"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAGConfig_policyWithOptionals(constant.AutoscalingGroupTestResource, "PT12H", "PT12H", "NEWEST_SERVER_FIRST", "PT12H"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT12H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.cooldown_period", "PT12H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.termination_policy_type", "NEWEST_SERVER_FIRST"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.cooldown_period", "PT12H"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_replicaNic(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_replicaNic(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.#", "1"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.name", "nic_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_active", "false"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.target_group.#", "0"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.#", "0"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.#", "0"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_nicWithTargetGroup(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
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
				Config: testAGConfig_nicWithTargetGroup(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.target_group.#", "1"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.target_group.0.target_group_id",
						constant.TargetGroupResource+".autoscaling_target_group", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.target_group.0.port", "80"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.target_group.0.weight", "1"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_nicWithFlowLog(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_nicWithFlowLog(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.0.name", "flow_log_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.0.bucket", "test-de-bucket"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.0.action", "ALL"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.flow_log.0.direction", "BIDIRECTIONAL"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_nicWithTcpFirewall(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_nicWithTCPFirewall(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_type", "INGRESS"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.name", "rule_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_start", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.port_range_end", "1000"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.type", "INGRESS"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_nicWithICMPFirewall(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_nicWithICMPFirewall(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_active", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_type", "INGRESS"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.name", "rule_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.protocol", "ICMP"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.icmp_code", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.icmp_type", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.firewall_rule.0.type", "INGRESS"),
				),
			},
			{
				ResourceName:      resourceAGName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAutoscalingGroup_replicaWithVolume(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		ExternalProviders:        randomProviderVersion343(),
		CheckDestroy:             testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAGConfig_replicaWithVolume(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.size", "30"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.bus", "IDE"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.volume.0.image_password", constant.RandomPassword+".image_password", "result"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.boot_order", "AUTO"),
					utils.TestNotEmptySlice(resourceAGName, "replica_configuration.0.volume.0.ssh_keys"),
				),
			},
			{
				Config: testAGConfig_replicaWithMultipleVolumes(constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.size", "20"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.type", "SSD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.bus", "VIRTIO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.volume.1.image_password", constant.RandomPassword+".image_password", "result"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.boot_order", "AUTO"),
					utils.TestNotEmptySlice(resourceAGName, "replica_configuration.0.volume.1.ssh_keys"),
				),
			},
			{
				ResourceName:            resourceAGName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replica_configuration.0.volume.0.ssh_keys", "replica_configuration.0.volume.0.image_password", "replica_configuration.0.volume.1.ssh_keys", "replica_configuration.0.volume.1.image_password"},
			},
		},
	})

}

func testAccCheckAutoscalingGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).AutoscalingClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != constant.AutoscalingGroupResource {
			continue
		}
		_, apiResponse, err := client.GetGroup(ctx, rs.Primary.ID, 0)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking for the destruction of autoscaling group %s: %w",
					rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("autoscaling group %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckAutoscalingGroupExists(name string, autoscalingGroup *autoscaling.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).AutoscalingClient

		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundGroup, _, err := client.GetGroup(ctx, rs.Primary.ID, 0)
		if err != nil {
			return fmt.Errorf("error occurred while fetching autoscaling group: %s, %w", rs.Primary.ID, err)
		}

		if *foundGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		autoscalingGroup = &foundGroup

		return nil
	}
}

func testAGConfig_base() string {
	return fmt.Sprintf(`
resource "ionoscloud_datacenter" "autoscaling_datacenter" {
	name     = "test_autoscaling_group"
	location = "de/fra"
}
`)
}

func testAGConfig_basic(name string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048
	}
}
`, name))
}

func testAGConfig_requiredUpdated(name string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group_new"
	max_replica_count = 4
	min_replica_count = 2

	policy {
		metric             = "INSTANCE_NETWORK_IN_BYTES"
		scale_in_action {
		  amount                  =  2
		  amount_type             = "PERCENTAGE"
		  delete_volumes = false
		}
		scale_in_threshold = 34
		scale_out_action  {
		  amount          =  10
		  amount_type     = "PERCENTAGE"
		}
		scale_out_threshold = 78
		unit                = "PER_MINUTE"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "1"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 1024
	}
}
`, name))
}

func testAGConfig_policyWithOptionals(rName string, rangeValue, scaleInCooldown, scaleInTermination, scaleOutCooldown string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
	      cooldown_period = %[3]q
		  termination_policy_type = %[4]q
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		  cooldown_period = %[5]q
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
		range               = %[2]q
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048
	}
}
`, rName, rangeValue, scaleInCooldown, scaleInTermination, scaleOutCooldown))
}

func testAGConfig_replicaNic(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource "ionoscloud_lan" "autoscaling_lan_1" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_datacenter.id
	public           = false
	name             = "test_autoscaling_group_1"
}

resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

		nic {
        	lan       =  ionoscloud_lan.autoscaling_lan_1.id
        	name      = "nic_1"
			dhcp      = true
			firewall_active = false
		}
	}
}
`, rName))
}

func testAGConfig_nicWithTargetGroup(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
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

resource "time_sleep" "wait_10_minutes" {
  depends_on = [ionoscloud_target_group.autoscaling_target_group]

  destroy_duration = "10m"
}

resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

		nic {
        	lan       =  ionoscloud_lan.autoscaling_lan_1.id
        	name      = "nic_1"
			dhcp      = true
			firewall_active = false
			
			target_group {
				 target_group_id = ionoscloud_target_group.autoscaling_target_group.id
        		 port            = 80
                 weight          = 1
			}
		}
	}

	depends_on = [time_sleep.wait_10_minutes]
}
`, rName))
}
func testAGConfig_nicWithFlowLog(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource "ionoscloud_lan" "autoscaling_lan_1" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_datacenter.id
	public           = false
	name             = "test_autoscaling_group_1"
}

resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

		nic {
        	lan       =  ionoscloud_lan.autoscaling_lan_1.id
        	name      = "nic_1"
			dhcp      = true
			firewall_active = false

			flow_log {
				name="flow_log_1"
	   			bucket="test-de-bucket"
				action="ALL"
				direction="BIDIRECTIONAL"
			}
		}
	}
}
`, rName))
}

func testAGConfig_nicWithTCPFirewall(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource "ionoscloud_lan" "autoscaling_lan_1" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_datacenter.id
	public           = false
	name             = "test_autoscaling_group_1"
}

resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

		nic {
        	lan       =  ionoscloud_lan.autoscaling_lan_1.id
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
		}
	}
}
`, rName))
}

func testAGConfig_nicWithICMPFirewall(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource "ionoscloud_lan" "autoscaling_lan_1" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_datacenter.id
	public           = false
	name             = "test_autoscaling_group_1"
}

resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

		nic {
        	lan       =  ionoscloud_lan.autoscaling_lan_1.id
        	name      = "nic_1"
			dhcp      = true
			firewall_active = true
			firewall_type = "INGRESS"
			firewall_rule {
				name = "rule_1"
				protocol = "ICMP"
				icmp_code = 1
				icmp_type = 1
				type = "INGRESS"
			}
		}
	}
}
`, rName))
}

func testAGConfig_replicaWithVolume(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

	 	volume {
      		image_alias     = "ubuntu:latest"
      		name      = "volume_1"
      		size      = 30
      		ssh_keys  = [%[2]q]
      		type      = "HDD"
      		image_password= random_password.image_password.result
      		boot_order = "AUTO"
			bus = "IDE"
    	}
	}
}

resource "random_password" "image_password" {
  length = 16
  special = false
}
`, rName, sshKey))
}

func testAGConfig_replicaWithMultipleVolumes(rName string) string {
	return utils.ConfigCompose(testAGConfig_base(), fmt.Sprintf(`
resource  "ionoscloud_autoscaling_group"  %[1]q {
	datacenter_id = ionoscloud_datacenter.autoscaling_datacenter.id
	name = "test_autoscaling_group"
	max_replica_count = 5
	min_replica_count = 1

	policy {
		metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		scale_in_action {
		  amount                  =  1
		  amount_type             = "ABSOLUTE"
		  delete_volumes = true
		}
		scale_in_threshold = 33
		scale_out_action  {
		  amount          =  1
		  amount_type     = "ABSOLUTE"
		}
		scale_out_threshold = 77
		unit                = "PER_HOUR"
	}

	replica_configuration {
		availability_zone = "AUTO"
		cores             = "2"
		cpu_family        = "INTEL_SKYLAKE"
	    ram               = 2048

	 	volume {
      		image_alias     = "ubuntu:latest"
      		name      = "volume_1"
      		size      = 30
      		ssh_keys  = [%[2]q]
      		type      = "HDD"
      		image_password= random_password.image_password.result
      		boot_order = "AUTO"
			bus = "IDE"
    	}

		volume {
    	  	image_alias    = "ubuntu:latest"
      		name           = "volume_2"
      		size           = 20
      		ssh_keys       = ["`+sshKey+`"]
      		type           = "SSD"
      		image_password = random_password.image_password.result
      		boot_order     = "AUTO"
    	}
	}
}

resource "random_password" "image_password" {
  length = 16
  special = false
}
`, rName, sshKey))
}
