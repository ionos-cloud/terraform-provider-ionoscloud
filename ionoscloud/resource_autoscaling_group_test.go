package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"
)

var resourceAutoscalingGroupName = AutoscalingGroupResource + "." + AutoscalingGroupTestResource

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
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigBasic, AutoscalingGroupTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", AutoscalingGroupTestResource),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "5"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "target_replica_count", "2"),
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
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", "LAN NIC 1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "2048"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image", "ee89912b-2290-11eb-af9f-1ee452559185"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", "Volume 1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", "30"),
					testNotEmptySlice(AutoscalingGroupResource, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data", "ZWNobyAiSGVsbG8sIFdvcmxkIgo="),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_password", "passw0rd"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigUpdate, UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "6"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "target_replica_count", "4"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.range", "PT12H"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", "NEWEST_SERVER_FIRST"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", "PT10M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_threshold", "35"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", "PT10M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_threshold", "80"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.lan", LanResource+".autoscaling_lan_1", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.name", UpdatedResources),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.0.dhcp", "false"),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "replica_configuration.0.nics.1.lan", LanResource+".autoscaling_lan_2", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.1.name", "LAN NIC 2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics.1.dhcp", "true"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "1024"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image", "129db64f-2291-11eb-af9f-1ee452559185"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.name", "Volume 2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.size", "40"),
					testNotEmptySlice(AutoscalingGroupResource, "replica_configuration.0.volumes.0.ssh_keys"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.type", "HDD"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.user_data", ""),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes.0.image_password", "passw0rdupdated")),
			},
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingGroupConfigUpdateRemoveNicsVolumes, UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingGroupExists(resourceAutoscalingGroupName, &autoscalingGroup),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "name", UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceAutoscalingGroupName, "datacenter_id", DatacenterResource+".autoscaling_datacenter", "id"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "max_replica_count", "6"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "min_replica_count", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "target_replica_count", "4"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.metric", "INSTANCE_NETWORK_IN_BYTES"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.range", "PT12H"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.termination_policy_type", "NEWEST_SERVER_FIRST"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_action.0.cooldown_period", "PT10M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_in_threshold", "35"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount", "2"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.amount_type", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_action.0.cooldown_period", "PT10M"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.scale_out_threshold", "80"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "policy.0.unit", "PER_MINUTE"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.availability_zone", "ZONE_1"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cores", "3"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.ram", "1024"),
					resource.TestCheckNoResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.nics"),
					resource.TestCheckNoResourceAttr(resourceAutoscalingGroupName, "replica_configuration.0.volumes")),
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

		if rs.Type != AutoscalingGroupResource {
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
resource ` + DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fkb"
}
resource ` + LanResource + ` "autoscaling_lan_1" {
	datacenter_id    = ` + DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_1"
}

resource ` + AutoscalingGroupResource + `  ` + AutoscalingGroupTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.autoscaling_datacenter.id
	max_replica_count      = 5
	min_replica_count      = 1
	target_replica_count   = 2
	name				   = "%s"
	policy {
    	metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
		range              = "PT24H"
        scale_in_action {
			amount        		    =  1
			amount_type    			= "ABSOLUTE"
			termination_policy_type = "OLDEST_SERVER_FIRST"
			cooldown_period			= "PT5M"
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
		cores 			  = "2"
		cpu_family 		  = "INTEL_XEON"
		nics {
			lan  		  = ` + LanResource + `.autoscaling_lan_1.id
			name		  = "LAN NIC 1"
			dhcp 		  = true
		}
		ram				  = 2048
		volumes	{
			image  		  = "ee89912b-2290-11eb-af9f-1ee452559185"
			name		  = "Volume 1"
			size 		  = 30
			ssh_key_paths = [ "/home/iulia/.ssh/id_rsa.pub"]
			ssh_key_values= [ "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\nNrRFi9wrf+M7Q== user@domain.local"]
			type		  = "HDD"
			user_data	  = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
			image_password= "passw0rd"
		}
	}
}
`

const testAccCheckAutoscalingGroupConfigUpdate = `
resource ` + DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fkb"
}
resource ` + LanResource + ` "autoscaling_lan_1" {
	datacenter_id    = ` + DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_1"
}

resource ` + LanResource + ` "autoscaling_lan_2" {
	datacenter_id    = ` + DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_2"
}

resource ` + AutoscalingGroupResource + `  ` + AutoscalingGroupTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.autoscaling_datacenter.id
	max_replica_count      = 6
	min_replica_count      = 2
	target_replica_count   = 4
	name				   = "%s"
	policy  {
    	metric             = "INSTANCE_NETWORK_IN_BYTES"
		range              = "PT12H"
        scale_in_action {
			amount        		    =  2
			amount_type    			= "PERCENTAGE"
			termination_policy_type = "NEWEST_SERVER_FIRST"
			cooldown_period			= "PT10M"
        }
		scale_in_threshold = 35
    	scale_out_action {
			amount         =  2
			amount_type    = "PERCENTAGE"
			cooldown_period= "PT10M"
        }
		scale_out_threshold = 80
        unit                = "PER_MINUTE"
	}
    replica_configuration {
		availability_zone = "ZONE_1"
		cores 			  = "3"
		cpu_family 		  = "INTEL_XEON"
		nics {
			lan  		  = ` + LanResource + `.autoscaling_lan_1.id
			name		  = "` + UpdatedResources + `"
			dhcp 		  = false
		}
		nics {
			lan  		  = ` + LanResource + `.autoscaling_lan_2.id
			name		  = "LAN NIC 2"
			dhcp 		  = true
		}
		ram				  = 1024
		volumes	{
			image  		  = "129db64f-2291-11eb-af9f-1ee452559185"
			name		  = "Volume 2"
			size 		  = 40
			ssh_key_paths = []
			ssh_key_values= [ "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\nNrRFi9wrf+M7Q== user@domain.local"]
			type		  = "HDD"
			image_password= "passw0rdupdated"
		}
	}
}
`

const testAccCheckAutoscalingGroupConfigUpdateRemoveNicsVolumes = `
resource ` + DatacenterResource + ` "autoscaling_datacenter" {
   name     = "test_autoscaling_group"
   location = "de/fkb"
}
resource ` + LanResource + ` "autoscaling_lan_1" {
	datacenter_id    = ` + DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_1"
}

resource ` + LanResource + ` "autoscaling_lan_2" {
	datacenter_id    = ` + DatacenterResource + `.autoscaling_datacenter.id
    public           = false
    name             = "test_autoscaling_group_2"
}

resource ` + AutoscalingGroupResource + `  ` + AutoscalingGroupTestResource + ` {
	datacenter_id = ` + DatacenterResource + `.autoscaling_datacenter.id
	max_replica_count      = 6
	min_replica_count      = 2
	target_replica_count   = 4
	name				   = "%s"
	policy  {
    	metric             = "INSTANCE_NETWORK_IN_BYTES"
		range              = "PT12H"
        scale_in_action {
			amount        		    =  2
			amount_type    			= "PERCENTAGE"
			termination_policy_type = "NEWEST_SERVER_FIRST"
			cooldown_period			= "PT10M"
        }
		scale_in_threshold = 35
    	scale_out_action {
			amount         =  2
			amount_type    = "PERCENTAGE"
			cooldown_period= "PT10M"
        }
		scale_out_threshold = 80
        unit                = "PER_MINUTE"
	}
    replica_configuration {
		availability_zone = "ZONE_1"
		cores 			  = "3"
		cpu_family 		  = "INTEL_XEON"
		ram				  = 1024
	}
}
`
