package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"log"
)

func resourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutoscalingGroupCreate,
		ReadContext:   resourceAutoscalingGroupRead,
		UpdateContext: resourceAutoscalingGroupUpdate,
		DeleteContext: resourceAutoscalingGroupDelete,
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Description: "Absolute URL to the resource's representation",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of resource",
							Computed:    true,
						},
						"id": {
							Type:         schema.TypeString,
							Description:  "Unique identifier for the resource",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"location": {
				Type:        schema.TypeString,
				Description: "Location of the datacenter. This location is the same as the one from the selected template.",
				Computed:    true,
			},
			"max_replica_count": {
				Type:        schema.TypeInt,
				Description: "Maximum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 200 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 200 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"min_replica_count": {
				Type:        schema.TypeInt,
				Description: "Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 200 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 200 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "User-defined name for the Autoscaling Group.",
				Optional:     true,
				ValidateFunc: isValidName,
			},
			"policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Required:     true,
							Type:         schema.TypeString,
							Description:  "The Metric that should trigger Scaling Actions. The values of the Metric are checked in fixed intervals.",
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"range": {
							Required:     true,
							Type:         schema.TypeString,
							Description:  "Defines the range of time from which samples will be aggregated. Default is 120s.",
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"scale_in_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.",
									},
									"amount_type": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
									"cooldown_period": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).",
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
								},
							},
						},
						"scale_in_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"scale_out_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.",
									},
									"amount_type": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
									"cooldown_period": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).",
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
								},
							},
						},
						"scale_out_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"unit": {
							Required:     true,
							Type:         schema.TypeString,
							Description:  "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"target_replica_count": {
				Type:        schema.TypeInt,
				Description: "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. This value can be set only at Group creation time, subsequent change via update (PUT) request is not possible.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 200 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 200 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"template": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Description: "Absolute URL to the resource's representation",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of resource",
							Computed:    true,
						},
						"id": {
							Type:         schema.TypeString,
							Description:  "Unique identifier for the resource",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceAutoscalingGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).AutoscalingClient

	group := autoscaling.Group{
		Properties: &autoscaling.GroupProperties{},
	}

	if _, datacenterOk := d.GetOk("datacenter"); datacenterOk {
		if id, idOk := d.GetOk("datacenter.0.id"); idOk {
			id := id.(string)
			group.Properties.Datacenter = &autoscaling.GroupPropertiesDatacenter{}
			group.Properties.Datacenter.Id = &id
		}
	}

	maxReplicaCount := int64(d.Get("max_replica_count").(int))
	group.Properties.MaxReplicaCount = &maxReplicaCount

	minReplicaCount := int64(d.Get("min_replica_count").(int))
	group.Properties.MinReplicaCount = &minReplicaCount

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		group.Properties.Name = &name
	}

	if _, policyOk := d.GetOk("policy"); policyOk {
		groupPolicy := autoscaling.GroupPolicy{}

		if metric, metricOk := d.GetOk("policy.0.metric"); metricOk {
			metric := autoscaling.Metric(metric.(string))
			groupPolicy.Metric = &metric
		}

		if policyRange, policyRangeOk := d.GetOk("policy.0.range"); policyRangeOk {
			policyRange := policyRange.(string)
			groupPolicy.Range = &policyRange
		}

		if _, scaleInActionOk := d.GetOk("policy.0.scale_in_action"); scaleInActionOk {
			groupPolicyAction := autoscaling.GroupPolicyAction{}

			if amount, amountOk := d.GetOk("policy.0.scale_in_action.0.amount"); amountOk {
				amount := float32(amount.(int))
				groupPolicyAction.Amount = &amount
			}

			if amountType, amountTypeOk := d.GetOk("policy.0.scale_in_action.0.amount_type"); amountTypeOk {
				amountType := autoscaling.ActionAmount(amountType.(string))
				groupPolicyAction.AmountType = &amountType
			}

			if cooldownPeriod, cooldownPeriodOk := d.GetOk("policy.0.scale_in_action.0.cooldown_period"); cooldownPeriodOk {
				cooldownPeriod := cooldownPeriod.(string)
				groupPolicyAction.CooldownPeriod = &cooldownPeriod
			}

			groupPolicy.ScaleInAction = &groupPolicyAction
		}

		if scaleInThreshlod, scaleInThreshlodOk := d.GetOk("policy.0.scale_in_threshold"); scaleInThreshlodOk {
			scaleInThreshlod := float32(scaleInThreshlod.(int))
			groupPolicy.ScaleInThreshold = &scaleInThreshlod
		}

		if _, scaleOutActionOk := d.GetOk("policy.0.scale_out_action"); scaleOutActionOk {
			groupPolicyAction := autoscaling.GroupPolicyAction{}

			if amount, amountOk := d.GetOk("policy.0.scale_out_action.0.amount"); amountOk {
				amount := float32(amount.(int))
				groupPolicyAction.Amount = &amount
			}

			if amountType, amountTypeOk := d.GetOk("policy.0.scale_out_action.0.amount_type"); amountTypeOk {
				amountType := autoscaling.ActionAmount(amountType.(string))
				groupPolicyAction.AmountType = &amountType
			}

			if cooldownPeriod, cooldownPeriodOk := d.GetOk("policy.0.scale_out_action.0.cooldown_period"); cooldownPeriodOk {
				cooldownPeriod := cooldownPeriod.(string)
				groupPolicyAction.CooldownPeriod = &cooldownPeriod
			}

			groupPolicy.ScaleOutAction = &groupPolicyAction
		}

		if scaleOutThreshlod, scaleOutThreshlodOk := d.GetOk("policy.0.scale_out_threshold"); scaleOutThreshlodOk {
			scaleOutThreshlod := float32(scaleOutThreshlod.(int))
			groupPolicy.ScaleOutThreshold = &scaleOutThreshlod
		}

		if unit, unitOk := d.GetOk("policy.0.unit"); unitOk {
			unit := autoscaling.QueryUnit(unit.(string))
			groupPolicy.Unit = &unit
		}

		group.Properties.Policy = &groupPolicy
	}

	if targetReplicaCount, targetReplicaCountOk := d.GetOk("target_replica_count"); targetReplicaCountOk {
		targetReplicaCount := int64(targetReplicaCount.(int))
		group.Properties.TargetReplicaCount = &targetReplicaCount
	}

	if _, templateOk := d.GetOk("template"); templateOk {
		if id, idOk := d.GetOk("template.0.id"); idOk {
			id := id.(string)
			group.Properties.Template = &autoscaling.GroupPropertiesTemplate{}
			group.Properties.Template.Id = &id
		}
	}

	autoscalingGroup, _, err := client.GroupsApi.AutoscalingGroupsPost(ctx).Group(group).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating autoscaling group: %s", err))
		return diags
	}

	d.SetId(*autoscalingGroup.Id)

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).AutoscalingClient

	group, apiResponse, err := client.GroupsApi.AutoscalingGroupsFindById(ctx, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived autoscaling group %s: %+v", d.Id(), group)

	if group.Properties.Datacenter != nil {
		datacenter := make([]interface{}, 1)

		datacenterEntry := make(map[string]interface{})

		if group.Properties.Datacenter.Href != nil {
			datacenterEntry["href"] = *group.Properties.Datacenter.Href
		}

		if group.Properties.Datacenter.Type != nil {
			datacenterEntry["type"] = *group.Properties.Datacenter.Type
		}

		if group.Properties.Datacenter.Id != nil {
			datacenterEntry["id"] = *group.Properties.Datacenter.Id
		}

		datacenter[0] = datacenterEntry
		if err := d.Set("datacenter", datacenter); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting datacenter property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if group.Properties.Location != nil {
		if err := d.Set("location", *group.Properties.Location); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting location property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if group.Properties.MaxReplicaCount != nil {
		if err := d.Set("max_replica_count", *group.Properties.MaxReplicaCount); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting max_replica_count property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if group.Properties.MinReplicaCount != nil {
		if err := d.Set("min_replica_count", *group.Properties.MinReplicaCount); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting min_replica_count property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if group.Properties.Name != nil {
		if err := d.Set("name", *group.Properties.Name); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if group.Properties.Policy != nil {
		policy := make([]interface{}, 1)

		policyEntry := make(map[string]interface{})

		if group.Properties.Policy.Metric != nil {
			policyEntry["metric"] = *group.Properties.Policy.Metric
		}

		if group.Properties.Policy.Range != nil {
			policyEntry["range"] = *group.Properties.Policy.Range
		}

		if group.Properties.Policy.ScaleInAction != nil {
			scaleIn := make([]interface{}, 1)

			scaleInEntry := make(map[string]interface{})

			if group.Properties.Policy.ScaleInAction.Amount != nil {
				scaleInEntry["amount"] = *group.Properties.Policy.ScaleInAction.Amount
			}

			if group.Properties.Policy.ScaleInAction.AmountType != nil {
				scaleInEntry["amount_type"] = *group.Properties.Policy.ScaleInAction.AmountType
			}

			if group.Properties.Policy.ScaleInAction.CooldownPeriod != nil {
				scaleInEntry["cooldown_period"] = *group.Properties.Policy.ScaleInAction.CooldownPeriod
			}

			scaleIn[0] = scaleInEntry
			policyEntry["scale_in_action"] = scaleIn
		}

		if group.Properties.Policy.ScaleInThreshold != nil {
			policyEntry["scale_in_threshold"] = *group.Properties.Policy.ScaleInThreshold
		}

		if group.Properties.Policy.ScaleOutAction != nil {
			scaleOut := make([]interface{}, 1)

			scaleOutEntry := make(map[string]interface{})

			if group.Properties.Policy.ScaleOutAction.Amount != nil {
				scaleOutEntry["amount"] = *group.Properties.Policy.ScaleOutAction.Amount
			}

			if group.Properties.Policy.ScaleOutAction.AmountType != nil {
				scaleOutEntry["amount_type"] = *group.Properties.Policy.ScaleOutAction.AmountType
			}

			if group.Properties.Policy.ScaleOutAction.CooldownPeriod != nil {
				scaleOutEntry["cooldown_period"] = *group.Properties.Policy.ScaleOutAction.CooldownPeriod
			}

			scaleOut[0] = scaleOutEntry
			policyEntry["scale_out_action"] = scaleOut
		}

		if group.Properties.Policy.ScaleOutThreshold != nil {
			policyEntry["scale_out_threshold"] = *group.Properties.Policy.ScaleOutThreshold
		}

		if group.Properties.Policy.Unit != nil {
			policyEntry["unit"] = *group.Properties.Policy.Unit
		}

		policy[0] = policyEntry
		if err := d.Set("policy", policy); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting policy property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	if err := d.Set("target_replica_count", group.Properties.TargetReplicaCount); err != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting target_replica_count property for autoscaling group %s: %s", d.Id(), err))
		return diags
	}

	if group.Properties.Template != nil {
		template := make([]interface{}, 1)

		templateEntry := make(map[string]interface{})

		if group.Properties.Template.Href != nil {
			templateEntry["href"] = *group.Properties.Template.Href
		}

		if group.Properties.Template.Type != nil {
			templateEntry["type"] = *group.Properties.Template.Type
		}

		if group.Properties.Template.Id != nil {
			templateEntry["id"] = *group.Properties.Template.Id
		}

		template[0] = templateEntry
		if err := d.Set("template", template); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting template property for autoscaling group %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceAutoscalingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).AutoscalingClient

	group := autoscaling.Group{
		Properties: &autoscaling.GroupProperties{},
	}

	if _, datacenterOk := d.GetOk("datacenter"); datacenterOk {
		if id, idOk := d.GetOk("datacenter.0.id"); idOk {
			id := id.(string)
			group.Properties.Datacenter = &autoscaling.GroupPropertiesDatacenter{}
			group.Properties.Datacenter.Id = &id
		}
	}

	maxReplicaCount := int64(d.Get("max_replica_count").(int))
	group.Properties.MaxReplicaCount = &maxReplicaCount

	minReplicaCount := int64(d.Get("min_replica_count").(int))
	group.Properties.MinReplicaCount = &minReplicaCount

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		group.Properties.Name = &name
	}

	if _, policyOk := d.GetOk("policy"); policyOk {
		groupPolicy := autoscaling.GroupPolicy{}

		if metric, metricOk := d.GetOk("policy.0.metric"); metricOk {
			metric := autoscaling.Metric(metric.(string))
			groupPolicy.Metric = &metric
		}

		if policyRange, policyRangeOk := d.GetOk("policy.0.range"); policyRangeOk {
			policyRange := policyRange.(string)
			groupPolicy.Range = &policyRange
		}

		if _, scaleInActionOk := d.GetOk("policy.0.scale_in_action"); scaleInActionOk {
			groupPolicyAction := autoscaling.GroupPolicyAction{}

			if amount, amountOk := d.GetOk("policy.0.scale_in_action.0.amount"); amountOk {
				amount := float32(amount.(int))
				groupPolicyAction.Amount = &amount
			}

			if amountType, amountTypeOk := d.GetOk("policy.0.scale_in_action.0.amount_type"); amountTypeOk {
				amountType := autoscaling.ActionAmount(amountType.(string))
				groupPolicyAction.AmountType = &amountType
			}

			if cooldownPeriod, cooldownPeriodOk := d.GetOk("policy.0.scale_in_action.0.cooldown_period"); cooldownPeriodOk {
				cooldownPeriod := cooldownPeriod.(string)
				groupPolicyAction.CooldownPeriod = &cooldownPeriod
			}

			groupPolicy.ScaleInAction = &groupPolicyAction
		}

		if scaleInThreshlod, scaleInThreshlodOk := d.GetOk("policy.0.scale_in_threshold"); scaleInThreshlodOk {
			scaleInThreshlod := float32(scaleInThreshlod.(int))
			groupPolicy.ScaleInThreshold = &scaleInThreshlod
		}

		if _, scaleOutActionOk := d.GetOk("policy.0.scale_out_action"); scaleOutActionOk {
			groupPolicyAction := autoscaling.GroupPolicyAction{}

			if amount, amountOk := d.GetOk("policy.0.scale_out_action.0.amount"); amountOk {
				amount := float32(amount.(int))
				groupPolicyAction.Amount = &amount
			}

			if amountType, amountTypeOk := d.GetOk("policy.0.scale_out_action.0.amount_type"); amountTypeOk {
				amountType := autoscaling.ActionAmount(amountType.(string))
				groupPolicyAction.AmountType = &amountType
			}

			if cooldownPeriod, cooldownPeriodOk := d.GetOk("policy.0.scale_out_action.0.cooldown_period"); cooldownPeriodOk {
				cooldownPeriod := cooldownPeriod.(string)
				groupPolicyAction.CooldownPeriod = &cooldownPeriod
			}

			groupPolicy.ScaleOutAction = &groupPolicyAction
		}

		if scaleOutThreshlod, scaleOutThreshlodOk := d.GetOk("policy.0.scale_out_threshold"); scaleOutThreshlodOk {
			scaleOutThreshlod := float32(scaleOutThreshlod.(int))
			groupPolicy.ScaleOutThreshold = &scaleOutThreshlod
		}

		if unit, unitOk := d.GetOk("policy.0.unit"); unitOk {
			unit := autoscaling.QueryUnit(unit.(string))
			groupPolicy.Unit = &unit
		}

		group.Properties.Policy = &groupPolicy
	}

	if d.HasChange("target_replica_count") {
		diags := diag.FromErr(fmt.Errorf("target_replica_count can not pe used in update requests"))
		return diags
	} else {
		if targetReplicaCount, targetReplicaCountOk := d.GetOk("target_replica_count"); targetReplicaCountOk {
			targetReplicaCount := int64(targetReplicaCount.(int))
			group.Properties.TargetReplicaCount = &targetReplicaCount
		}
	}

	if _, templateOk := d.GetOk("template"); templateOk {
		if id, idOk := d.GetOk("template.0.id"); idOk {
			id := id.(string)
			group.Properties.Template = &autoscaling.GroupPropertiesTemplate{}
			group.Properties.Template.Id = &id
		}
	}

	_, _, err := client.GroupsApi.AutoscalingGroupsPut(ctx, d.Id()).Group(group).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating autoscaling group %s: %s", d.Id(), err))
		return diags
	}

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).AutoscalingClient

	_, err := client.GroupsApi.AutoscalingGroupsDelete(ctx, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting an autoscaling group %s %s", d.Id(), err))
		return diags
	}

	d.SetId("")

	return nil
}
