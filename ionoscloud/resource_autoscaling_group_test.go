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

const resourceAutoscalingGroupName = constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupTestResource

func TestAccAutoscalingGroupBasic(t *testing.T) {
	var autoscalingGroup autoscaling.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigBasic, constant.AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", constant.AutoscalingGroupTestResource),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "5"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.metric", "INSTANCE_CPU_UTILIZATION_AVERAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", "ABSOLUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", "OLDEST_SERVER_FIRST"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", "ABSOLUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_threshold", "77"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.unit", "PER_HOUR"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cores", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_SKYLAKE"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", "nic_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", "30"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_password", "passw0rd"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.boot_order", "AUTO"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigUpdate, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", "ABSOLUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", "OLDEST_SERVER_FIRST"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_threshold", "33"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", "ABSOLUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", "PT5M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_threshold", "86"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_SKYLAKE"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", "nic_1_updated"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.1.name", "nic_2_updated"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.1.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", "30"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_password", "passw0rd"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.boot_order", "AUTO")),
			},
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigUpdateRemoveOptionalFields, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", constant.DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.range", "PT24H"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.delete_volumes", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_threshold", "35"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_threshold", "80"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_SKYLAKE"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", constant.LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", "nic_1_updated"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", "false"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_alias", "ubuntu:latest"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", "volume_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", "10"),
					utils.TestNotEmptySlice(constant.AutoscalingGroupResource, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_password", "passw0rd"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.boot_order", "AUTO")),
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
		_, apiResponse, err := client.GetGroup(ctx, rs.Primary.ID)
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

func testAccCheckAutoscalingGroupExists(n string, autoscalingGroup *autoscaling.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).AutoscalingClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundGroup, _, err := client.GetGroup(ctx, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error occured while fetching backup unit: %s", rs.Primary.ID)
		}
		if *foundGroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		autoscalingGroup = &foundGroup

		return nil
	}
}

const testAccCheckAutoscalingGroupConfigBasic = `
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
  max_replica_count      = 5
  min_replica_count      = 1
  //target_replica_count   = 2
  name           = "%s"
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
    nics {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1"
      dhcp      = true
    }
    ram          = 2048
    volumes {
      image_alias     = "ubuntu:latest"
      name      = "volume_1"
      size      = 30
      ssh_keys  = ["` + sshKey + `"]
      type      = "HDD"
      user_data    = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
      image_password= "passw0rd"
      boot_order = "AUTO"
    }
  }
}
`

const testAccCheckAutoscalingGroupConfigUpdate = `
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
    nics {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1_updated"
      dhcp      = true
    }
    nics {
      lan       = ` + constant.LanResource + `.autoscaling_lan_2.id
      name      = "nic_2_updated"
      dhcp      = true
    }
    ram          = 2048
    volumes {
      image_alias     = "ubuntu:latest"
      name      = "volume_1"
      size      = 30
      ssh_keys  = ["` + sshKey + `"]
      type      = "HDD"
      image_password= "passw0rd"
      boot_order = "AUTO"
    }
  }
}
`

const testAccCheckAutoscalingGroupConfigUpdateRemoveOptionalFields = `
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
    nics {
      lan       = ` + constant.LanResource + `.autoscaling_lan_1.id
      name      = "nic_1_updated"
      dhcp      = false
    }
    ram          = 2048
    volumes {
      image_alias     = "ubuntu:latest"
      name      = "volume_1"
      size      = 10
      ssh_keys  = ["` + sshKey + `"]
      type      = "HDD"
      image_password= "passw0rd"
      boot_order = "AUTO"
    }
  }
}
`
