//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const resourceAGName = constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupTestResource

func TestAccAutoscalingGroupBasic(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: randomProviderVersion343(),
		CheckDestroy:      testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAG_ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", constant.AutoscalingGroupTestResource),
					resource.TestCheckResourceAttrPair(resourceAGName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.volume.0.image_password", constant.RandomPassword+".server_image_password", "result"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "5"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_CPU_UTILIZATION_AVERAGE"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.termination_policy_type", string(autoscaling.TERMINATIONPOLICYTYPE_OLDEST_SERVER_FIRST)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "77"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_HOUR"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.name", "nic_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.size", "30"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volume.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.boot_order", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.bus", "IDE"),
				),
			},
			{
				Config: fmt.Sprintf(testAGGroup_ConfigUpdate, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAGName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.termination_policy_type", string(autoscaling.TERMINATIONPOLICYTYPE_OLDEST_SERVER_FIRST)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "86"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.name", "nic_1_updated"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.1.name", "nic_2_updated"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.1.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.#", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.size", "30"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volume.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.boot_order", "AUTO")),
			},
			{
				Config: fmt.Sprintf(testAG_Update_RemoveNic_AddVolumes, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAGName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.termination_policy_type", string(autoscaling.TERMINATIONPOLICYTYPE_OLDEST_SERVER_FIRST)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_ABSOLUTE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "86"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.name", "nic_1_updated"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.#", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.#", "3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.size", "30"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volume.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.boot_order", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.name", "volume_2"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.size", "10"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.type", "SSD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.boot_order", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.name", "volume_3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.size", "5"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.boot_order", "AUTO")),
			},
			{
				Config: fmt.Sprintf(testAG_Update_RemoveOptionalFields, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAGName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAGName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAGName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "max_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.range", "PT2M"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_PERCENTAGE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_action.0.delete_volumes", "true"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_in_threshold", "35"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_action.0.amount_type", string(autoscaling.ACTIONAMOUNT_PERCENTAGE)),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.scale_out_threshold", "80"),
					resource.TestCheckResourceAttr(resourceAGName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.cpu_family", string(autoscaling.CPUFAMILY_INTEL_SKYLAKE)),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.#", "1"),
					resource.TestCheckResourceAttrPair(resourceAGName, "replica_configuration.0.nic.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.name", "nic_1_updated"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.#", "3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.size", "10"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volume.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.0.boot_order", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.name", "volume_2"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.size", "20"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.type", "SSD"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.1.boot_order", "AUTO"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.name", "volume_3"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.size", "5"),
					resource.TestCheckResourceAttr(resourceAGName, "replica_configuration.0.volume.2.type", "HDD")),
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
			return fmt.Errorf("error occured while fetching autoscaling group: %s, %w", rs.Primary.ID, err)
		}

		if *foundGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		autoscalingGroup = &foundGroup

		return nil
	}
}

const testAG_ConfigBasic = `
resource ` + constant.DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fra"
}
resource ` + constant.LanResource + ` "autoscaling_lan_1" {
  datacenter_id    = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  public           = false
name             = "test_autoscaling_group_1"
}

resource ` + constant.AutoscalingGroupResource + `  ` + constant.AutoscalingGroupTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  location = "de/fra"
  max_replica_count      = 5
  min_replica_count      = 1
  //target_replica_count   = 2
  name           = "` + constant.AutoscalingGroupTestResource + `"
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
    nic {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1"
      dhcp      = true
    }
    ram          = 2048
    volume {
      image_alias = "ubuntu:latest"
      name        = "volume_1"
      image_password = ` + constant.RandomPassword + `.server_image_password.result
      size        = 30
      ssh_keys    = ["` + sshKey + `"]
      type        = "HDD"
      user_data    = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
      boot_order = "AUTO"
      bus = "IDE"
    }
  }
}


` + ServerImagePassword

const testAGGroup_ConfigUpdate = `
resource ` + constant.DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fra"
}
resource ` + constant.LanResource + ` "autoscaling_lan_1" {
  datacenter_id    = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_1"
}

resource ` + constant.LanResource + ` "autoscaling_lan_2" {
  datacenter_id    = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_2"
}

resource ` + constant.AutoscalingGroupResource + `  ` + constant.AutoscalingGroupTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  location = "de/fra"
  max_replica_count      = 2
  min_replica_count      = 1
  //target_replica_count   = 2
  name           = "%s"
  policy {
    metric             = "INSTANCE_NETWORK_IN_BYTES"
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
    scale_out_threshold = 86
    unit                = "PER_MINUTE"
  }
  replica_configuration {
    availability_zone = "AUTO"
    cores             = "3"
    cpu_family        = "INTEL_SKYLAKE"
    nic {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1_updated"
      dhcp      = true
    }
    nic {
      lan       = ` + constant.LanResource + `.autoscaling_lan_2.id
      name      = "nic_2_updated"
      dhcp      = true
    }
    ram          = 2048
    volume {
      image_alias     = "ubuntu:latest"
      name      = "volume_1"
      size      = 30
      ssh_keys  = ["` + sshKey + `"]
      type      = "HDD"
      image_password= random_password.image_password.result
      boot_order = "AUTO"
    }
  }
}

resource "random_password" "image_password" {
  length = 16
  special = false
}
`
const testAG_Update_RemoveNic_AddVolumes = `
resource ` + constant.DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fra"
}
resource ` + constant.LanResource + ` "autoscaling_lan_1" {
  datacenter_id    = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_1"
}

resource ` + constant.LanResource + ` "autoscaling_lan_2" {
  datacenter_id    = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_2"
}

resource ` + constant.AutoscalingGroupResource + `  ` + constant.AutoscalingGroupTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  location = "de/fra"
  max_replica_count      = 2
  min_replica_count      = 1
  //target_replica_count   = 2
  name           = "%s"
  policy {
    metric             = "INSTANCE_NETWORK_IN_BYTES"
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
    scale_out_threshold = 86
    unit                = "PER_MINUTE"
  }
  replica_configuration {
    availability_zone = "AUTO"
    cores             = "3"
    cpu_family        = "INTEL_SKYLAKE"
    ram          = 2048
    nic {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1_updated"
      dhcp      = true
    }
    volume {
      image_alias    = "ubuntu:latest"
      name           = "volume_1"
      size           = 30
      ssh_keys       = ["` + sshKey + `"]
      type           = "HDD"
      image_password = random_password.image_password.result
      boot_order     = "AUTO"
    }
    volume {
      image_alias    = "ubuntu:latest"
      name           = "volume_2"
      size           = 10
      ssh_keys       = ["` + sshKey + `"]
      type           = "SSD"
      image_password = random_password.image_password.result
      boot_order     = "AUTO"
    }
    volume {
      image_alias    = "ubuntu:latest"
      name           = "volume_3"
      size           = 5
      ssh_keys       = ["` + sshKey + `"]
      type           = "HDD"
      image_password = random_password.image_password.result
      boot_order     = "AUTO"
    }
  }
}

resource "random_password" "image_password" {
  length = 16
  special = false
}
`

const testAG_Update_RemoveOptionalFields = `
resource ` + constant.DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fra"
}
resource ` + constant.LanResource + ` "autoscaling_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  public        = false
  name          = "test_autoscaling_group_1"
}

resource ` + constant.LanResource + ` "autoscaling_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
    public      = false
    name        = "test_autoscaling_group_2"
}

resource ` + constant.AutoscalingGroupResource + `  ` + constant.AutoscalingGroupTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.autoscaling_datacenter.id
  location = "de/fra"
  max_replica_count      = 2
  min_replica_count      = 1
  //target_replica_count   = 2
  name           = "%s"
  policy {
    metric = "INSTANCE_NETWORK_IN_BYTES"
    scale_in_action {
      amount      =  2
      amount_type = "PERCENTAGE"
      delete_volumes = true
    }
    scale_in_threshold = 35
    scale_out_action {
      amount      =  2
      amount_type = "PERCENTAGE"
    }
    scale_out_threshold = 80
    unit                = "PER_MINUTE"
  }
  replica_configuration {
    availability_zone = "AUTO"
    cores             = "3"
    cpu_family        = "INTEL_SKYLAKE"
    nic {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1_updated"
      dhcp      = false
    }
    ram          = 2048
    volume {
      image_alias = "ubuntu:latest"
      name        = "volume_1"
      size        = 10
      type        = "HDD"
      boot_order  = "AUTO"
      ssh_keys       = ["` + sshKey + `"]
      image_password = random_password.image_password.result
    }
    volume {
      image_alias    = "ubuntu:latest"
      name           = "volume_2"
      size           = 20
      ssh_keys       = ["` + sshKey + `"]
      type           = "SSD"
      image_password = random_password.image_password.result
      boot_order     = "AUTO"
    }
    volume {
      image_alias    = "ubuntu:latest"
      name           = "volume_3"
      size           = 5
      ssh_keys       = ["` + sshKey + `"]
      type           = "HDD"
      image_password = random_password.image_password.result
      boot_order     = "AUTO"
    }
  }
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
`
