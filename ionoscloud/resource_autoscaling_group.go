package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutoscalingGroupCreate,
		ReadContext:   resourceAutoscalingGroupRead,
		UpdateContext: resourceAutoscalingGroupUpdate,
		DeleteContext: resourceAutoscalingGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAutoscalingGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"max_replica_count": {
				Type:             schema.TypeInt,
				Description:      "The maximum value for the number of replicas on a VM Auto Scaling Group. Must be >= 0 and <= 200. Will be enforced for both automatic and manual changes.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 200)),
			},
			"min_replica_count": {
				Type:             schema.TypeInt,
				Description:      "The minimum value for the number of replicas on a VM Auto Scaling Group. Must be >= 0 and <= 200. Will be enforced for both automatic and manual changes",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 200)),
			},
			//will be left until GA in case it is added again in swagger
			//"target_replica_count": {
			//	Type:             schema.TypeInt,
			//	Description:      "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.",
			//	Optional:         true,
			//	Computed:         true,
			//	ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 200)),
			//},
			"name": {
				Type:             schema.TypeString,
				Description:      "User-defined name for the Autoscaling Group.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 255)),
			},
			"policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Defines the behavior of this VM Auto Scaling Group. A policy consists of triggers and actions, where an action is an automated behavior, and the trigger defines the circumstances under which the action is triggered. Currently, two separate actions are supported, namely scaling inward and outward, triggered by the thresholds defined for a particular metric.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Required:         true,
							Type:             schema.TypeString,
							Description:      "The Metric that should trigger the scaling actions. Metric values are checked at fixed intervals.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(autoscaling.METRIC_CPU_UTILIZATION_AVERAGE), string(autoscaling.METRIC_NETWORK_IN_BYTES), string(autoscaling.METRIC_NETWORK_IN_PACKETS), string(autoscaling.METRIC_NETWORK_OUT_PACKETS), string(autoscaling.METRIC_NETWORK_OUT_BYTES)}, true)),
						},
						"range": {
							Optional:    true,
							Default:     "PT2M",
							Type:        schema.TypeString,
							Description: "Specifies the time range for which the samples are to be aggregated. Must be >= 2 minutes.",
						},
						"scale_in_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Defines the action to be taken when the 'scaleInThreshold' is exceeded. Here, scaling is always about removing VMs associated with this VM Auto Scaling Group. By default, the termination policy is 'OLDEST_SERVER_FIRST' is effective.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When 'amountType=ABSOLUTE' specifies the absolute number of VMs that are removed. The value must be between 1 to 10. 'amountType=PERCENTAGE' specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200. At least one VM is always removed. Note that for 'SCALE_IN' operations, volumes are not deleted after the server is deleted.",
									},
									"amount_type": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(autoscaling.ACTIONAMOUNT_ABSOLUTE), string(autoscaling.ACTIONAMOUNT_PERCENTAGE)}, true)),
									},
									"termination_policy_type": {
										Optional:         true,
										Computed:         true,
										Type:             schema.TypeString,
										Description:      "The type of termination policy for the VM Auto Scaling Group to follow a specific pattern for scaling-in replicas. The default termination policy is 'OLDEST_SERVER_FIRST'.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"OLDEST_SERVER_FIRST", "NEWEST_SERVER_FIRST", "RANDOM"}, true)),
									},
									"cooldown_period": {
										Optional:    true,
										Default:     "PT5M",
										Type:        schema.TypeString,
										Description: "The minimum time that elapses after the start of this scaling action until the following scaling action is started. While a scaling action is in progress, no second action is initiated for the same VM Auto Scaling Group. Instead, the metric is re-evaluated after the current scaling action completes (either successfully or with errors). This is currently validated with a minimum value of 2 minutes and a maximum of 24 hours. The default value is 5 minutes if not specified.",
									},
									"delete_volumes": {
										Optional:    true,
										Type:        schema.TypeBool,
										Description: "If set to `true`, when deleting an replica during scale in, any attached volume will also be deleted. When set to `false`, all volumes remain in the datacenter and must be deleted manually.  **Note**, that every scale-out creates new volumes. When they are not deleted, they will eventually use all of your contracts resource limits. At this point, scaling out would not be possible anymore.",
									},
								},
							},
						},
						"scale_in_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The upper threshold for the value of the 'metric'. Used with the 'greater than' (>) operator. A scale-out action is triggered when this value is exceeded, specified by the 'scale_out_action' property. The value must have a lower minimum delta to the 'scale_in_threshold', depending on the metric, to avoid competing for actions simultaneously. If 'properties.policy.unit=TOTAL', a value >= 40 must be chosen.",
						},
						"scale_out_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Defines the action to be performed when the 'scaleOutThreshold' is exceeded. Here, scaling is always about adding new VMs to this VM Auto Scaling Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When 'amountType=ABSOLUTE' specifies the absolute number of VMs that are added. The value must be between 1 to 10. 'amountType=PERCENTAGE' specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200. At least one VM is always added or removed.",
									},
									"amount_type": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(autoscaling.ACTIONAMOUNT_ABSOLUTE), string(autoscaling.ACTIONAMOUNT_PERCENTAGE)}, true)),
									},
									"cooldown_period": {
										Optional:    true,
										Default:     "PT5M",
										Type:        schema.TypeString,
										Description: "The minimum time that elapses after the start of this scaling action until the following scaling action is started. While a scaling action is in progress, no second action is initiated for the same VM Auto Scaling Group. Instead, the metric is re-evaluated after the current scaling action completes (either successfully or with errors). This is currently validated with a minimum value of 2 minutes and a maximum of 24 hours. The default value is 5 minutes if not specified.",
									},
								},
							},
						},
						"scale_out_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The upper threshold for the value of the 'metric'. Used with the 'greater than' (>) operator. A scale-out action is triggered when this value is exceeded, specified by the 'scaleOutAction' property. The value must have a lower minimum delta to the 'scaleInThreshold', depending on the metric, to avoid competing for actions simultaneously. If 'properties.policy.unit=TOTAL', a value >= 40 must be chosen.",
						},
						"unit": {
							Required:         true,
							Type:             schema.TypeString,
							Description:      "Units of the applied Metric. Possible values are: PER_HOUR, PER_MINUTE, PER_SECOND, TOTAL.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PER_HOUR", "PER_MINUTE", "PER_SECOND", "TOTAL"}, true)),
						},
					},
				},
			},
			"replica_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Required:         true,
							Type:             schema.TypeString,
							Description:      "The zone where the VMs are created using this configuration.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
						},
						"cores": {
							Required:         true,
							Type:             schema.TypeInt,
							Description:      "The total number of cores for the VMs.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						},
						"cpu_family": {
							Optional:         true,
							Type:             schema.TypeString,
							Description:      "The zone where the VMs are created using this configuration.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AMD_OPTERON", "INTEL_SKYLAKE", "INTEL_XEON"}, true)),
						},
						"nic": {
							Type:        schema.TypeSet,
							Description: "Set of NICs associated with this Replica.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lan": {
										Required:         true,
										Type:             schema.TypeInt,
										Description:      "Lan ID for this replica Nic.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
									},
									"name": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "Name for this replica NIC.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 255)),
									},
									"dhcp": {
										Optional:    true,
										Type:        schema.TypeBool,
										Default:     true,
										Description: "Dhcp flag for this replica Nic. This is an optional attribute with default value of 'true' if not given in the request payload or given as null.",
									},
								}},
						},
						"ram": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
						},
						//TODO: there might be a problem un update of multiple volumes because the order isn't guaranteed when we get the response from the API. We might have to move from TypeList to TypeSet, like for nics.
						"volume": {
							Type:        schema.TypeSet,
							Description: "List of volumes associated with this Replica.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Optional:         true,
										Type:             schema.TypeString,
										Description:      "The image installed on the disk. Currently, only the UUID of the image is supported. Note that either 'image' or 'imageAlias' must be specified, but not both.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
									},
									"image_alias": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "The image installed on the volume. Must be an 'imageAlias' as specified via the images API. Note that one of 'image' or 'imageAlias' must be set, but not both.",
									},
									"name": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "Name for this replica volume.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 255)),
									},
									"size": {
										Required:         true,
										Type:             schema.TypeInt,
										Description:      "User-defined size for this replica volume in GB.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
									},
									"ssh_keys": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"type": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "Storage Type for this replica volume. Possible values: SSD, HDD, SSD_STANDARD or SSD_PREMIUM",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, true)),
									},
									"user_data": {
										Optional:    true,
										Computed:    true,
										Type:        schema.TypeString,
										Description: "User-data (Cloud Init) for this replica volume.",
									},
									"image_password": {
										Optional:    true,
										Sensitive:   true,
										Type:        schema.TypeString,
										Description: "Image password for this replica volume.",
									},
									"boot_order": {
										Optional: true,
										Type:     schema.TypeString,
										Description: `Determines whether the volume will be used as a boot volume. Set to NONE, the volume will not be used as boot volume. 
Set to PRIMARY, the volume will be used as boot volume and set to AUTO will delegate the decision to the provisioning engine to decide whether to use the volume as boot volume.
Notice that exactly one volume can be set to PRIMARY or all of them set to AUTO.`,
									},
									"bus": {
										Optional:    true,
										Type:        schema.TypeString,
										Default:     autoscaling.BUSTYPE_VIRTIO,
										Description: `The bus type of the volume. Default setting is 'VIRTIO'. The bus type 'IDE' is also supported.`,
									},
									"backup_unit_id": {
										Type:        schema.TypeString,
										Description: "The uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property.",
										Optional:    true,
										Computed:    true,
									},
								}},
						},
					},
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "Unique identifier for the resource",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
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

func resourceAutoscalingGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient

	group, err := autoscalingService.GetAutoscalingGroupDataCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured getting data from the provided schema: %w", err))
	}
	if group.Properties != nil && group.Properties.Name != nil {
		log.Printf("[DEBUG] autoscaling group data extracted: %+v", *group.Properties.Name)
	}

	autoscalingGroup, _, err := client.CreateGroup(ctx, *group)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating autoscaling group: %w", err))
		return diags
	}

	d.SetId(*autoscalingGroup.Id)
	log.Printf("[INFO] autoscaling Group created. Id set to %s", *autoscalingGroup.Id)

	if err := checkAction(ctx, client, d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if err := autoscalingService.SetAutoscalingGroupData(d, autoscalingGroup.Properties); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return nil
}

func resourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	client := meta.(services.SdkBundle).AutoscalingClient

	group, apiResponse, err := client.GetGroup(ctx, d.Id(), 2)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] resource %s not found: %+v", d.Id(), err)
			d.SetId("")
			return nil
		} else {
			diags := diag.FromErr(fmt.Errorf("error while retrieving autoscaling group with id %v, %w", d.Id(), err))
			return diags
		}
	}

	log.Printf("[INFO] successfully retrieved autoscaling group %s: %+v", d.Id(), group)

	if err := autoscalingService.SetAutoscalingGroupData(d, group.Properties); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] successfully set autoscaling group data %s", d.Id())

	return nil
}

func resourceAutoscalingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	client := meta.(services.SdkBundle).AutoscalingClient

	group, err := autoscalingService.GetAutoscalingGroupDataUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if group.Properties != nil && group.Properties.Name != nil {
		log.Printf("[DEBUG] autoscaling group data extracted: %s", *group.Properties.Name)
	}
	updatedGroup, _, err := client.UpdateGroup(ctx, d.Id(), *group)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating autoscaling group %s: %w", d.Id(), err))
		return diags
	}
	log.Printf("[INFO] autoscaling Group updated.")

	if err := checkAction(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(autoscalingService.SetAutoscalingGroupData(d, updatedGroup.Properties))
}

func resourceAutoscalingGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient

	_, err := client.DeleteGroup(ctx, d.Id())
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting an autoscaling group %s %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] autoscaling Group deleted: %s.", d.Id())

	d.SetId("")

	return nil
}

func resourceAutoscalingGroupImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).AutoscalingClient

	groupId := d.Id()

	group, apiResponse, err := client.GetGroup(ctx, d.Id(), 0)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find autoscaling group %q", groupId)
		}
		return nil, fmt.Errorf("an error occured while retrieveing autoscaling grouo %q, %w", groupId, err)
	}

	log.Printf("[INFO] autoscaling group found: %+v", group)

	if err := autoscalingService.SetAutoscalingGroupData(d, group.Properties); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func actionReady(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData, actionId string) (bool, error) {

	action, _, err := client.GetAction(ctx, d.Id(), actionId)
	if err != nil {
		return true, fmt.Errorf("error checking action status: %w", err)
	}

	if *action.Properties.ActionStatus == autoscaling.ACTIONSTATUS_FAILED {
		return false, fmt.Errorf("action failed")
	}
	return *action.Properties.ActionStatus == autoscaling.ACTIONSTATUS_SUCCESSFUL, nil
}

// checkAction gets the triggered action and waits for it to be ready
func checkAction(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData) error {
	actions, _, err := client.GetAllActions(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("error fetching group actions: %w", err)
	}

	var actionId string

	if actions.Items != nil && len(*actions.Items) > 0 && (*actions.Items)[0].Id != nil {
		//latest action
		actionId = *(*actions.Items)[0].Id
	} else {
		return fmt.Errorf("no action triggered for group: %s", d.Id())
	}

	//wait for completion of triggered action
	for {
		log.Printf("[INFO] waiting for action %s to be ready...", actionId)

		actionSuccessful, rsErr := actionReady(ctx, client, d, actionId)
		if rsErr != nil {
			return fmt.Errorf("error while checking status of action %s: %w", actionId, rsErr)
		}

		if actionSuccessful {
			log.Printf("[INFO] action was ready: %s", actionId)
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			return fmt.Errorf("group creation/update timed out! WARNING: your group was created/updated but the action was not yet ready. " +
				"Check your Ionos Cloud account for updates")
		}
	}

	return nil
}
