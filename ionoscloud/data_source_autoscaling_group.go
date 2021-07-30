package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"strings"
)

func dataSourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "User-defined name for the Autoscaling Group.",
				Optional:    true,
			},
			"datacenter": {
				Type:        schema.TypeList,
				Description: "VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.",
				Computed:    true,
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
							Type:        schema.TypeString,
							Description: "Unique identifier for the resource",
							Computed:    true,
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
				Computed:    true,
			},
			"min_replica_count": {
				Type:        schema.TypeInt,
				Description: "Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.",
				Computed:    true,
			},
			"policy": {
				Type:        schema.TypeList,
				Description: "Specifies the behavior of this Autoscaling Group. A policy consists of Triggers and Actions, whereby an Action is some kind of automated behavior, and a Trigger is defined by the circumstances under which the Action is triggered. Currently, two separate Actions, namely Scaling In and Out are supported, triggered through Thresholds defined on a given Metric.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The Metric that should trigger Scaling Actions. The values of the Metric are checked in fixed intervals.",
						},
						"range": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Defines the range of time from which samples will be aggregated. Default is 120s.",
						},
						"scale_in_action": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.",
									},
									"amount_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
									},
									"cooldown_period": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).",
									},
								},
							},
						},
						"scale_in_threshold": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"scale_out_action": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the Autoscaling Group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed.",
									},
									"amount_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
									},
									"cooldown_period": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling Action has started, until the next Scaling Action will be started. Additionally, if a Scaling Action is currently in progress, no second Scaling Action will be started for the same Autoscaling Group. Instead, the Metric will be re-evaluated after the current Scaling Action completed (either successful or with failures).",
									},
								},
							},
						},
						"scale_out_threshold": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "A lower threshold on the value of `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In Action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"unit": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Specifies the Action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this Autoscaling Group.",
						},
					},
				},
			},
			"target_replica_count": {
				Type:        schema.TypeInt,
				Description: "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. This value can be set only at Group creation time, subsequent change via update (PUT) request is not possible.",
				Computed:    true,
			},
			"template": {
				Type:        schema.TypeList,
				Description: "VMs for this Autoscaling Group will be created using this Template.",
				Computed:    true,
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
							Type:        schema.TypeString,
							Description: "Unique identifier for the resource",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).AutoscalingClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		diags := diag.FromErr(fmt.Errorf("id and name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(fmt.Errorf("please provide either the group id or name"))
		return diags
	}

	var group autoscaling.Group
	var err error

	if idOk {
		/* search by ID */
		group, _, err = client.GroupsApi.AutoscalingGroupsFindById(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		var groups autoscaling.GroupCollection

		groups, _, err := client.GroupsApi.AutoscalingGroupsGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group: %s", err.Error()))
			return diags
		}

		found := false
		if groups.Items != nil {
			for _, g := range *groups.Items {
				tmpGroup, _, err := client.GroupsApi.AutoscalingGroupsFindById(ctx, *g.Id).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group %s: %s", *g.Id, err))
					return diags
				}
				if tmpGroup.Properties.Name != nil {
					if strings.Contains(*tmpGroup.Properties.Name, name.(string)) {
						group = tmpGroup
						found = true
						break
					}
				}
			}
		}

		if !found {
			diags := diag.FromErr(fmt.Errorf("group not found"))
			return diags
		}
	}

	if diags := setAutoscalingGroupData(d, &group); diags != nil {
		return diags
	}

	return nil
}

func setAutoscalingGroupData(d *schema.ResourceData, group *autoscaling.Group) diag.Diagnostics {
	d.SetId(*group.Id)
	if err := d.Set("id", *group.Id); err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	if group.Properties != nil {
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
	}

	return nil
}
