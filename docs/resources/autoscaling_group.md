---
layout: "ionoscloud"
page_title: "IonosCloud: autoscaling_group"
sidebar_current: "docs-resource-autoscaling_group"
description: |-
Creates and manages IonosCloud Autoscaling Group.
---

# ionoscloud_autoscaling_group

Manages an Autoscaling Group on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_datacenter" "datacenter_example" {
    name     = "datacenter_example"
    location = "de/fkb"
}

resource "ionoscloud_lan" "lan_example_1" {
    datacenter_id    = ionoscloud_datacenter.datacenter_example.id
    public           = false
    name             = "lan_example_1"
}

resource "ionoscloud_lan" "lan_example_2" {
    datacenter_id    = ionoscloud_datacenter.datacenter_example.id
    public           = false
    name             = "lan_example_2"
}

resource "ionoscloud_autoscaling_group" "autoscaling_group_example" {
    datacenter_id = ionoscloud_datacenter.datacenter_example.id
    max_replica_count      = 5
    min_replica_count      = 1
    target_replica_count   = 2
    name				   = "autoscaling_group_example"
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
        ram				  = 2048
        nics {
            lan  		  = ionoscloud_lan.lan_example_1.id
            name		  = "nic_example_1"
            dhcp 		  = true
        }
        nics {
            lan  		  = ionoscloud_lan.lan_example_2.id
            name		  = "nic_example_2"
            dhcp 		  = true
        }
        volumes	{
            image  		  = "ee89912b-2290-11eb-af9f-1ee452559185"
            name		  = "volume_example"
            size 		  = 30
            ssh_key_paths = [ "PATH_TO/.ssh/id_rsa.pub"]
            ssh_key_values= [ "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\nNrRFi9wrf+M7Q== user@domain.local"]
            type		  = "HDD"
            user_data	  = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
            image_password= "passw0rd"
		}
	}
}
```

## Argument Reference

The following arguments are supported:

- `max_replica_count` - (Required)[int] Maximum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
- `min_replica_count` - (Required)[int] Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
- `target_replica_count` - (Optional)[int] The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.
- `name` - (Required)[string] User-defined name for the Autoscaling Group.
- `policy` - (Required) Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.
    - `metric` - (Required)[string] The Metric that should trigger the scaling actions. Metric values are checked at fixed intervals. Possible values: `INSTANCE_CPU_UTILIZATION_AVERAGE`, `INSTANCE_NETWORK_IN_BYTES`, `INSTANCE_NETWORK_IN_PACKETS`, `INSTANCE_NETWORK_OUT_BYTES`, `INSTANCE_NETWORK_OUT_PACKETS`
    - `range` - (Optional)[string] Defines the time range, for which the samples will be aggregated. Default is 120s. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    - `scale_in_action` - (Required) Specifies the action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this autoscaling group. Default termination policy is OLDEST_SERVER_FIRST.
        - `amount` - (Required)[int] When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.
        - `amount_type` - (Required)[string] The type for the given amount. Possible values are: `ABSOLUTE`, `PERCENTAGE`.
        - `termination_policy_type` - (Optional)[string] The type of the termination policy for the autoscaling group so that a specific pattern is followed for Scaling-In instances. Default termination policy is `OLDEST_SERVER_FIRST`. Possible values are: `OLDEST_SERVER_FIRST`, `NEWEST_SERVER_FIRST`, `RANDOM`
        - `cooldown_period` - (Optional)[string] Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    - `scale_in_threshold` - (Required)[int] A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
    - `scale_out_action` - (Required)[list] Specifies the Action to take when the `scaleInThreshold
        - `amount` - (Required)[int] When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.
        - `amount_type` - (Required)[string] The type for the given amount. Possible values are: `ABSOLUTE`, `PERCENTAGE`.
        - `cooldown_period` - (Optional)[string] Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    - `scale_out_threshold` - (Required)[int] A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
    - `unit` - (Required)[string] Units of the applied Metric. Possible values are: `PER_HOUR`, `PER_MINUTE`, `PER_SECOND`, `TOTAL`.
- `replica_configuration` - (Required)
    - `availability_zone` - (Required)[string] The zone where the VMs are created using this configuration. Possible values are: `AUTO`, `ZONE_1`, `ZONE_2`.
    - `cores` - (Required)[string] The total number of cores for the VMs.
    - `cpu_family` - (Optional)[string] PU family for the VMs created using this configuration. If null, the VM will be created with the default CPU family for the assigned location. Possible values are: `AMD_OPTERON`, `INTEL_SKYLAKE`, `INTEL_XEON`.
    - `nics` - (Optional)[list] List of NICs associated with this Replica.
        - `lan` - (Required)[int] Lan ID for this replica Nic.
        - `name` - (Required)[string] Name for this replica NIC.
        - `dhcp` - (Optional)[bool] Dhcp flag for this replica Nic. This is an optional attribute with default value of 'true' if not given in the request payload or given as null.
    - `ram` - (Required)[int] The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.
    - `volumes` - (Optional)[list] List of volumes associated with this Replica. Only a single volume is currently supported.
        `image` - (Required)[string] The image installed on the volume. Only the UUID of the image is presently supported.
        `name` - (Required)[string] Name for this replica volume.
        `size` - (Required)[int] Name for this replica volume.
        `ssh_key_paths` - (Optional)[list] List of ssh key paths.
        `ssh_key_values` - (Optional)[list] List of ssh key values.
        `ssh_keys` - (Computed) List of ssh key computed from what it was provided in `ssk_key_paths` and `ssh_key_values` fields.
        `type` - (Required)[string] Storage Type for this replica volume. Possible values: `SSD`, `HDD`, `SSD_STANDARD` or `SSD_PREMIUM`.
        `user_data` - (Optional)[string] User-data (Cloud Init) for this replica volume. Make sure you provide a Cloud Init compatible image in conjunction with this parameter.
        `image_password` - (Optional)[string] Image password for this replica volume.
- `datacenter_id` - (Required)[string] Unique identifier for the resource
- `location` - (Computed) Location of the data center.
