---
subcategory: "Autoscaling"
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
    location = "de/fra"
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

resource "ionoscloud_target_group" "autoscaling_target_group" {
  name                      = "Target Group Example"
  algorithm                 = "ROUND_ROBIN"
  protocol                  = "HTTP"
}

resource "ionoscloud_autoscaling_group" "autoscaling_group_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  max_replica_count      = 2
  min_replica_count      = 1
  name                   = "autoscaling_group_example"
  policy {
    metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
    range              = "PT24H"
    scale_in_action {
      amount                  =  1
      amount_type             = "ABSOLUTE"
      termination_policy_type = "OLDEST_SERVER_FIRST"
      cooldown_period         = "PT5M"
      delete_volumes          = true
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
    cores               = "2"
    cpu_family           = "INTEL_SKYLAKE"
    ram                  = 2048
    nic {
      lan   = ionoscloud_lan.lan_example_1.id
      name  = "nic_example_1"
      dhcp  = true
    }
    nic {
      lan   = ionoscloud_lan.lan_example_2.id
      name  = "nic_example_2"
      dhcp  = true
      firewall_active = true
      firewall_type = "INGRESS"
      firewall_rule {
        name = "rule_1"
        protocol = "TCP"
        port_range_start = 1
        port_range_end = 1000
        type = "INGRESS"
      }

      flow_log {
        name="flow_log_1"
        bucket="test-de-bucket"
        action="ALL"
        direction="BIDIRECTIONAL"
      }

      target_group {
        target_group_id = ionoscloud_target_group.autoscaling_target_group.id
        port            = 80
        weight          = 50
      }
    }
    volume    {
      image_alias    = "ubuntu:latest"
      name           = "volume_example"
      size           = 10
      type           = "HDD"
      user_data      = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
      image_password = random_password.server_image_password.result
      boot_order     = "AUTO"
    }
  }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```

## Argument Reference

The following arguments are supported:

- `max_replica_count` - (Required)[int] The maximum value for the number of replicas on a VM Auto Scaling Group. Must be >= 0 and <= 200. Will be enforced for both automatic and manual changes.
- `min_replica_count` - (Required)[int] The minimum value for the number of replicas on a VM Auto Scaling Group. Must be >= 0 and <= 200. Will be enforced for both automatic and manual changes.
- `target_replica_count` - (Optional)[int] The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.
- `name` - (Required)[string] User-defined name for the Autoscaling Group.
- `policy` - (Required)[List] Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.
    - `metric` - (Required)[string] The Metric that should trigger the scaling actions. Metric values are checked at fixed intervals. Possible values: `INSTANCE_CPU_UTILIZATION_AVERAGE`, `INSTANCE_NETWORK_IN_BYTES`, `INSTANCE_NETWORK_IN_PACKETS`, `INSTANCE_NETWORK_OUT_BYTES`, `INSTANCE_NETWORK_OUT_PACKETS`
    - `range` - (Optional)[string] Defines the time range, for which the samples will be aggregated. Default is 120s. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    - `scale_in_action` - (Required)[list] Specifies the action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this autoscaling group. Default termination policy is OLDEST_SERVER_FIRST.
        - `amount` - (Required)[int] When `amountType == ABSOLUTE`, this is the number of VMs removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be removed in one step. There will always be at least one VM removed. For SCALE_IN operation new volumes are NOT deleted after the server deletion.
        - `amount_type` - (Required)[string] The type for the given amount. Possible values are: `ABSOLUTE`, `PERCENTAGE`.
        - `termination_policy_type` - (Optional)[string] The type of the termination policy for the autoscaling group so that a specific pattern is followed for Scaling-In replicas. Default termination policy is `OLDEST_SERVER_FIRST`. Possible values are: `OLDEST_SERVER_FIRST`, `NEWEST_SERVER_FIRST`, `RANDOM`
        - `cooldown_period` - (Optional)[string] Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
        - `delete_volumes` - (Required)[bool] If set to `true`, when deleting a replica during scale in, any attached volume will also be deleted. When set to `false`, all volumes remain in the datacenter and must be deleted manually. Note that every scale-out creates new volumes. When they are not deleted, they will eventually use all of your contracts resource limits. At this point, scaling out would not be possible anymore.
    - `scale_in_threshold` - (Required)[int] A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
    - `scale_out_action` - (Required)[list] Specifies the action to take when the `scaleOutThreshold` is exceeded. Hereby, scaling out is always about adding new VMs to this autoscaling group.
        - `amount` - (Required)[int] When `amountType=ABSOLUTE` specifies the absolute number of VMs that are added. The value must be between 1 to 10. `amountType=PERCENTAGE` specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200. At least one VM is always added.
        - `amount_type` - (Required)[string] The type for the given amount. Possible values are: `ABSOLUTE`, `PERCENTAGE`.
        - `cooldown_period` - (Optional)[string] Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    - `scale_out_threshold` - (Required)[int] The upper threshold for the value of the `metric`. Used with the `greater than` (>) operator. A scale-out action is triggered when this value is exceeded, specified by the `scaleOutAction` property. The value must have a lower minimum delta to the `scaleInThreshold`, depending on the metric, to avoid competing for actions simultaneously. If `properties.policy.unit=TOTAL`, a value >= 40 must be chosen.
    - `unit` - (Required)[string] Units of the applied Metric. Possible values are: `PER_HOUR`, `PER_MINUTE`, `PER_SECOND`, `TOTAL`.
- `replica_configuration` - (Required)[List]  
    - `availability_zone` - (Required)[string] The zone where the VMs are created using this configuration. Possible values are: `AUTO`, `ZONE_1`, `ZONE_2`.
    - `cores` - (Required)[int] The total number of cores for the VMs.
    - `cpu_family` - (Optional)[string] PU family for the VMs created using this configuration. If null, the VM will be created with the default CPU family for the assigned location. Possible values are: `INTEL_SKYLAKE`, `INTEL_XEON`.
    - `nics` - (Optional)[set] List of NICs associated with this Replica.
        - `lan` - (Required)[int] Lan ID for this replica Nic.
        - `name` - (Required)[string] Name for this replica NIC.
        - `dhcp` - (Optional)[bool] Dhcp flag for this replica Nic. This is an optional attribute with default value of `true` if not given in the request payload or given as null.
        - `firewall_active` - (Optional)[bool] Firewall active flag.
        - `firewall_type` - (Optional)[string] The type of firewall rules that will be allowed on the NIC. Valid values: INGRESS EGRESS BIDIRECTIONAL. If not specified, the default INGRESS value is used.
        - `firewall_rules` - (Optional)[list] List of firewall rules associated with this NIC.
            - `protocol` - (Required)[string] The protocol for the rule: TCP, UDP, ICMP, ANY. Property cannot be modified after creation (disallowed in update requests).
            - `name` - (Optional)[string] The name of the firewall rule.
            - `source_mac` - (Optional)[string] Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Value null allows all source MAC address. Valid format: aa:bb:cc:dd:ee:ff.
            - `source_ip` -  (Optional)[string] Only traffic originating from the respective IPv4 address is allowed. Value null allows all source IPs.
            - `target_ip` - (Optional)[string] In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Value null allows all target IPs.
            - `port_range_start` - (Optional)[int] Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
            - `port_range_end` - (Optional)[int] Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
            - `icmp_type` - (Optional)[string] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen. Value null allows all codes.
            - `icmp_code` - (Optional)[int] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.
            - `type` - (Optional)[string] The type of firewall rule. If is not specified, it will take the default value INGRESS.
        - `flow_log` - (Optional)[list] Only 1 flow log can be configured. Only the name field can change as part of an update. Flow logs holistically capture network information such as source and destination IP addresses, source and destination ports, number of packets, amount of bytes, the start and end time of the recording, and the type of protocol – and log the extent to which your instances are being accessed.
            - `action` - (Required)[string] Specifies the action to be taken when the rule is matched. Possible values: ACCEPTED, REJECTED, ALL. Immutable, forces re-creation.
            - `bucket` - (Required)[string] Specifies the S3 IONOS bucket where the flow log data will be stored. The bucket must exist. Immutable, forces re-creation.
            - `direction` - (Required)[string] Specifies the traffic direction pattern. Valid values: INGRESS, EGRESS, BIDIRECTIONAL. Immutable, forces re-creation.
            - `name` - (Required)[string] Specifies the name of the flow log.
        - `target_group` - (Optional)[list] In order to link VM to ALB, target group must be provided
            - `target_group_id` (Required)[string] The ID of the target group.
            - `weight` (Optional)[int] The weight of the target group.
            - `port` (Optional)[int] The port of the target group.
    - `ram` - (Required)[int] The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.
    - `volume` - (Optional)[list] List of volumes associated with this Replica.
        - `image` - (Optional)[string] The image installed on the volume. Only the UUID of the image is presently supported.
        - `image_alias` - (Optional)[string] The image installed on the volume. Must be an `imageAlias` as specified via the images API. Note that one of `image` or `imageAlias` must be set, but not both.
        - `name` - (Required)[string] Name for this replica volume.
        - `size` - (Required)[int] Name for this replica volume.
        - `ssh_keys` - (Optional) List of ssh keys, supports values or paths to files. Cannot be changed at update.
        - `type` - (Required)[string] Storage Type for this replica volume. Possible values: `SSD`, `HDD`, `SSD_STANDARD` or `SSD_PREMIUM`.
        - `user_data` - (Optional)[string] User-data (Cloud Init) for this replica volume. Make sure you provide a Cloud Init compatible image in conjunction with this parameter.
        - `image_password` - (Optional)[string] Image password for this replica volume.
        - `bus` - (Optional)[string] The bus type of the volume. Default setting is `VIRTIO`. The bus type `IDE` is also supported.
        - `backup_unit_id` - (Optional)[string] The uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either `public image` or `imageAlias` in conjunction with this property.
        - `boot_order` - (Optional)[string] Determines whether the volume will be used as a boot volume. Set to NONE, the volume will not be used as boot volume. Set to PRIMARY, the volume will be used as boot volume and set to AUTO will delegate the decision to the provisioning engine to decide whether to use the volume as boot volume.
      Notice that exactly one volume can be set to PRIMARY or all of them set to AUTO.
- `datacenter_id` - (Required)[string] Unique identifier for the resource
- `location` - (Computed) Location of the data center.
