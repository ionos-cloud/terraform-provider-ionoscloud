package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/autoscaling"
)

func dataSourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "UUID of the Autoscaling Group.",
				Optional:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "User-defined name for the Autoscaling Group.",
				Optional:    true,
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
			"target_replica_count": {
				Type:        schema.TypeInt,
				Description: "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.",
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
							Description: "The Metric that should trigger the scaling actions. Metric values are checked at fixed intervals.",
						},
						"range": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Defines the time range, for which the samples will be aggregated. Default is 120s.",
						},
						"scale_in_action": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Specifies the action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this autoscaling group. Default termination policy is OLDEST_SERVER_FIRST.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.",
									},
									"amount_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
									},
									"termination_policy_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type of the termination policy for the autoscaling group so that a specific pattern is followed for Scaling-In instances. Default termination policy is OLDEST_SERVER_FIRST.",
									},
									"cooldown_period": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given.",
									},
								},
							},
						},
						"scale_in_threshold": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The lower threshold for the value of the `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"scale_out_action": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Specifies the action to take when the `scaleOutThreshold` is exceeded. Hereby, scaling out is always about adding new VMs to this autoscaling group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.",
									},
									"amount_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
									},
									"cooldown_period": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given.",
									},
								},
							},
						},
						"scale_out_threshold": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The upper threshold for the value of the `metric`. Will be used with `greater than` (>) operator. Exceeding this will start a Scale-Out action as specified by the `scaleOutAction` property. The value must have a lower minimum delta to the `scaleInThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"unit": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Units of the applied Metric. Possible values are: PER_HOUR, PER_MINUTE, PER_SECOND, TOTAL.",
						},
					},
				},
			},
			"replica_configuration": {
				Type:        schema.TypeList,
				Description: "VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it have the same `location` as the `template`.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The zone where the VMs are created using this configuration.",
						},
						"cores": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The total number of cores for the VMs.",
						},
						"cpu_family": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The zone where the VMs are created using this configuration.",
						},
						"nics": {
							Type:        schema.TypeList,
							Description: "List of NICs associated with this Replica.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lan": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Lan ID for this replica Nic.",
									},
									"name": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Name for this replica NIC.",
									},
									"dhcp": {
										Computed:    true,
										Type:        schema.TypeBool,
										Description: "Dhcp flag for this replica Nic. This is an optional attribute with default value of 'true' if not given in the request payload or given as null.",
									},
								}},
						},
						"ram": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
						},
						"volumes": {
							Type:        schema.TypeList,
							Description: "List of volumes associated with this Replica. Only a single volume is currently supported.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The image installed on the volume. Only the UUID of the image is presently supported.",
									},
									"name": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Name for this replica volume.",
									},
									"size": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "User-defined size for this replica volume in GB.",
									},
									"ssh_keys": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
									"type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Storage Type for this replica volume. Possible values: SSD, HDD, SSD_STANDARD or SSD_PREMIUM",
									},
									"user_data": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "user-data (Cloud Init) for this replica volume.",
									},
									"image_password": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Image password for this replica volume.",
									},
								}},
						},
					},
				},
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
				Description: "Location of the data center.",
				Computed:    true,
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
		groups, _, err := client.GroupsApi.AutoscalingGroupsGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group: %s", err.Error()))
			return diags
		}

		var results []autoscaling.Group

		if groups.Items != nil {
			for _, g := range *groups.Items {
				tmpGroup, _, err := client.GroupsApi.AutoscalingGroupsFindById(ctx, *g.Id).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group %s: %s", *g.Id, err))
					return diags
				}

				if tmpGroup.Properties.Name != nil && *tmpGroup.Properties.Name == name.(string) {
					results = append(results, tmpGroup)
					break
				}
			}
		}

		if results != nil && len(results) > 0 {
			group = results[0]
		} else {
			diags := diag.FromErr(fmt.Errorf("group not found"))
			return diags
		}
	}

	if err := autoscalingService.SetAutoscalingGroupData(d, group); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
