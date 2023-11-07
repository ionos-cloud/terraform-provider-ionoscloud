---
layout: "ionoscloud"
page_title: "IonosCloud : autoscaling_group"
sidebar_current: "docs-datasource-autoscaling_group"
description: |-
  Get information on a IonosCloud Autoscaling Group
---

# ionoscloud\_autoscaling_group

The autoscaling group data source can be used to search for and return an existing Autoscaling Group. You can provide a string for the name or id parameters which will be compared with provisioned Autoscaling Groups. If a single match is found, it will be returned. If your search results in multiple matches, an error will be generated. When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_autoscaling_group" "autoscaling_group" {
  name = "test_ds"
}
```

## Argument Reference

* `id` - (Optional) Id of an existing Autoscaling Group that you want to search for.
* `name` - (Optional) Name of an existing Autoscaling Group that you want to search for.

Either `name` or `id` must be provided. If none or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:


* `id` - UUID of the Autoscaling Group.
* `name` - The name of the Autoscaling Group.
* `datacenter` - VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.
    * `href` - Absolute URL to the resource's representation
    * `type` - Type of resource
    * `id` - Unique identifier for the resource
* `location` - Location of the datacenter. This location is the same as the one from the selected template.
* `max_replica_count` - Maximum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
* `min_replica_count` - Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.
* `policy` - Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.
    * `metric` - The Metric that should trigger Scaling Actions. The values of the Metric are checked in fixed intervals.
    * `range` - Defines the range of time from which samples will be aggregated. Default is 120s.
      *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    * `scale_in_action` - Specifies the Action to take when the `scaleInThreshold
        * `amount` - When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.
        * `amount_type` - The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].
        * `cooldown_period` - Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).
          *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    * `scale_in_threshold` - A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
    * `scale_out_action` - Specifies the Action to take when the `scaleInThreshold
        * `amount` - When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.
        * `amount_type` - The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].
        * `cooldown_period` - Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).
          *Note that when you set it to values like 5m the API will automatically transform it in PT5M, so the plan will show you a diff in state that should be ignored.*
    * `scale_out_threshold` - A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.
    * `unit` - Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.
* `target_replica_count` - The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. This value can be set only at Group creation time, subsequent change via update (PUT) request is not possible.
* `template` - VMs for this Autoscaling Group will be created using this Template.
    * `href` - Absolute URL to the resource's representation
    * `type` - Type of resource
    * `id` - Unique identifier for the resource