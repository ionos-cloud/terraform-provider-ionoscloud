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
resource "ionoscloud_autoscaling_group" "autoscaling_group" {
	datacenter {
       id                  = ionoscloud_datacenter.autoscaling_group.id
    }
	max_replica_count      = 5
	min_replica_count      = 1
	name				   = "autoscaling_group"
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
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional)[string] The name of the Autoscaling Group.
- `datacenter` - (Required)[list] VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.
  - `href` - (Computed)[string] Absolute URL to the resource's representation
  - `type` - (Computed)[string]Type of resource
  - `id` - (Required)[string] Unique identifier for the resource
- `location` - (Computed)[string] Location of the datacenter. This location is the same as the one from the selected template.
- `max_replica_count` - (Required)[int] Maximum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
- `min_replica_count` - (Required)[int] Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
- `policy` - (Required)[list] Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.
  - `metric` - (Required)[string] The Metric that should trigger Scaling Actions. The values of the Metric are checked in fixed intervals.
  - `range` - (Required)[string] Defines the range of time from which samples will be aggregated. Default is 120s. *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
  - `scale_in_action` - (Required)[list] Specifies the Action to take when the `scaleInThreshold
    - `amount` - (Required)[int] When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.
    - `amount_type` - (Required)[string] The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].
    - `cooldown_period` - (Required)[string] Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures). *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
  - `scale_in_threshold` - (Required)[int] A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
  - `scale_out_action` - (Required)[list] Specifies the Action to take when the `scaleInThreshold
    - `amount` - (Required)[int] When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.
    - `amount_type` - (Required)[string] The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].
    - `cooldown_period` - (Required)[string] Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures). *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
  - `scale_out_threshold` - (Required)[int] A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
  - `unit` - (Required)[string] Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.
- `target_replica_count` - (Required)[int] The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. This value can be set only at Group creation time, subsequent change via update (PUT) request is not possible.
- `template` - (Required)[list] VMs for this Autoscaling Group will be created using this Template.
  - `href` - (Computed)[string] Absolute URL to the resource's representation
  - `type` - (Computed)[string] Type of resource
  - `id` - (Required)[string] Unique identifier for the resource