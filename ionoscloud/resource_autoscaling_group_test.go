package ionoscloud

import (
	"context"
	"fmt"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAutoscalingGroup_Basic(t *testing.T) {
	var autoscalingGroup autoscaling.Group
	autoscalingGroupName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigBasic, autoscalingGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists("ionoscloud_autoscaling_group.autoscaling_group", &autoscalingGroup),
					resource.TestCheckResourceAttr("ionoscloud_autoscaling_group.autoscaling_group", "name", autoscalingGroupName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists("ionoscloud_autoscaling_group.autoscaling_group", &autoscalingGroup),
					resource.TestCheckResourceAttr("ionoscloud_autoscaling_group.autoscaling_group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckAutoscalingGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).AutoscalingClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "ionoscloud_autoscaling_group" {
			continue
		}

		_, apiResponse, err := client.GroupsApi.AutoscalingGroupsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking for the destruction of autoscaling group %s: %s",
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
		client := testAccProvider.Meta().(SdkBundle).AutoscalingClient

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

		foundGroup, _, err := client.GroupsApi.AutoscalingGroupsFindById(ctx, rs.Primary.ID).Execute()

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
resource "ionoscloud_datacenter" "autoscaling_group" {
   name     = "test_autoscaling_group"
   location = "de/txl"
}
resource "ionoscloud_lan" "autoscaling_group" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_group.id
    public           = false
    name             = "test_autoscaling_group"
}
resource "ionoscloud_autoscaling_template" "autoscaling_group" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "test_autoscaling_group"
    nics    {
		lan              = ionoscloud_lan.autoscaling_group.id
        name             = "test_autoscaling_group"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_group"
		size             = 50
    	type             = "HDD"
	}
}
resource "ionoscloud_autoscaling_group" "autoscaling_group" {
	datacenter {
       id                  = ionoscloud_datacenter.autoscaling_group.id
    }
	max_replica_count      = 5
	min_replica_count      = 1
	name				   = "%s"
	policy  {
    	metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		range              = "PT24H"
        scale_in_action {
			amount         =  1
			amount_type    = "ABSOLUTE"
			cooldown_period= "PT5M"
        }
		scale_in_threshold = 33
    	scale_out_action {
			amount         =  1
			amount_type    = "ABSOLUTE"
			cooldown_period= "PT5M"
        }
		scale_out_threshold = 77
        unit                = "PER_HOUR"
	}
    target_replica_count    = 1
	template {
		id = ionoscloud_autoscaling_template.autoscaling_group.id
    }
}
`

const testAccCheckAutoscalingGroupConfigUpdate = `
resource "ionoscloud_datacenter" "autoscaling_group" {
   name     = "test_autoscaling_group"
   location = "de/txl"
}
resource "ionoscloud_lan" "autoscaling_group" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_group.id
    public           = false
    name             = "test_autoscaling_group"
}
resource "ionoscloud_autoscaling_template" "autoscaling_group" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "test_autoscaling_group"
    nics    {
		lan              = ionoscloud_lan.autoscaling_group.id
        name             = "test_autoscaling_group"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_group"
		size             = 50
    	type             = "HDD"
	}
}
resource "ionoscloud_autoscaling_group" "autoscaling_group" {
	datacenter {
       id                  = ionoscloud_datacenter.autoscaling_group.id
    }
	max_replica_count      = 0
	min_replica_count      = 0
	name				   = "updated"
	policy  {
    	metric             = "INSTANCE_NETWORK_IN_BYTES"
		range              = "PT24H"
        scale_in_action {
			amount         =  1
			amount_type    = "ABSOLUTE"
			cooldown_period= "PT5M"
        }
		scale_in_threshold = 33
    	scale_out_action {
			amount         =  1
			amount_type    = "ABSOLUTE"
			cooldown_period= "PT5M"
        }
		scale_out_threshold = 86
        unit                = "PER_MINUTE"
	}
    target_replica_count    = 1
	template {
		id = ionoscloud_autoscaling_template.autoscaling_group.id
    }
}
`
