package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAutoscalingGroup_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutoscalingGroupCreateResources,
			},
			{
				Config: testAccDataSourceAutoscalingGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_autoscaling_group.autoscaling_group", "name", "test_ds_autoscaling_group"),
				),
			},
		},
	})
}

func TestAccDataSourceAutoscalingGroup_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutoscalingGroupCreateResources,
			},
			{
				Config: testAccDataSourceAutoscalingGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_autoscaling_group.autoscaling_group", "name", "test_ds_autoscaling_group"),
				),
			},
		},
	})

}

const testAccDataSourceAutoscalingGroupCreateResources = `
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
	name				   = "test_ds_autoscaling_group"
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

const testAccDataSourceAutoscalingGroupMatchId = `
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
	name				   = "test_ds_autoscaling_group"
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

data "ionoscloud_autoscaling_group" "autoscaling_group" {
  id			= ionoscloud_autoscaling_group.autoscaling_group.id
}
`

const testAccDataSourceAutoscalingGroupMatchName = `
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
	name				   = "test_ds_autoscaling_group"
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

data "ionoscloud_autoscaling_group" "autoscaling_group" {
  name			= "test_ds"
}
`
