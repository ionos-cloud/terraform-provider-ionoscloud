package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/autoscaling"
	"log"
	"time"
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
				Type:         schema.TypeInt,
				Description:  "Maximum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.",
				Required:     true,
				ValidateFunc: validation.All(validation.IntBetween(0, 200)),
			},
			"min_replica_count": {
				Type:         schema.TypeInt,
				Description:  "Minimum replica count value for `targetReplicaCount`. Will be enforced for both automatic and manual changes.",
				Required:     true,
				ValidateFunc: validation.All(validation.IntBetween(0, 200)),
			},
			"target_replica_count": {
				Type:         schema.TypeInt,
				Description:  "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntBetween(0, 200)),
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "User-defined name for the Autoscaling Group.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(0, 255)),
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
							Description:  "The Metric that should trigger the scaling actions. Metric values are checked at fixed intervals.",
							ValidateFunc: validation.All(validation.StringInSlice([]string{"INSTANCE_CPU_UTILIZATION_AVERAGE", "INSTANCE_NETWORK_IN_BYTES", "INSTANCE_NETWORK_IN_PACKETS", "INSTANCE_NETWORK_OUT_BYTES", "INSTANCE_NETWORK_OUT_PACKETS"}, true)),
						},
						"range": {
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Defines the time range, for which the samples will be aggregated. Default is 120s.",
						},
						"scale_in_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Specifies the action to take when the `scaleInThreshold` is exceeded. Hereby, scaling in is always about removing VMs that are currently associated with this autoscaling group. Default termination policy is OLDEST_SERVER_FIRST.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.",
									},
									"amount_type": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateFunc: validation.All(validation.StringInSlice([]string{"ABSOLUTE", "PERCENTAGE"}, true)),
									},
									"termination_policy_type": {
										Optional:     true,
										Computed:     true,
										Type:         schema.TypeString,
										Description:  "The type of the termination policy for the autoscaling group so that a specific pattern is followed for Scaling-In instances. Default termination policy is OLDEST_SERVER_FIRST.",
										ValidateFunc: validation.All(validation.StringInSlice([]string{"OLDEST_SERVER_FIRST", "NEWEST_SERVER_FIRST", "RANDOM"}, true)),
									},
									"cooldown_period": {
										Optional:    true,
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given.",
									},
								},
							},
						},
						"scale_in_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The lower threshold for the value of the `metric`. Will be used with `less than` (<) operator. Exceeding this will start a Scale-In action as specified by the `scaleInAction` property. The value must have a higher minimum delta to the `scaleOutThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"scale_out_action": {
							Required:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Specifies the action to take when the `scaleOutThreshold` is exceeded. Hereby, scaling out is always about adding new VMs to this autoscaling group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Required:    true,
										Type:        schema.TypeInt,
										Description: "When `amountType == ABSOLUTE`, this is the number of VMs added or removed in one step. When `amountType == PERCENTAGE`, this is a percentage value, which will be applied to the autoscaling group's current `targetReplicaCount` in order to derive the number of VMs that will be added or removed in one step. There will always be at least one VM added or removed. For SCALE_IN operation now volumes are NOT deleted after the server deletion.",
									},
									"amount_type": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "The type for the given amount. Possible values are: [ABSOLUTE, PERCENTAGE].",
										ValidateFunc: validation.All(validation.StringInSlice([]string{"ABSOLUTE", "PERCENTAGE"}, true)),
									},
									"cooldown_period": {
										Optional:    true,
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Minimum time to pass after this Scaling action has started, until the next Scaling action will be started. Additionally, if a Scaling action is currently in progress, no second Scaling action will be started for the same autoscaling group. Instead, the Metric will be re-evaluated after the current Scaling action is completed (either successfully or with failures). This is validated with a minimum value of 2 minutes and a maximum of 24 hours currently. Default value is 5 minutes if not given.",
									},
								},
							},
						},
						"scale_out_threshold": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The upper threshold for the value of the `metric`. Will be used with `greater than` (>) operator. Exceeding this will start a Scale-Out action as specified by the `scaleOutAction` property. The value must have a lower minimum delta to the `scaleInThreshold` depending on the `metric` to avoid competitive actions at the same time.",
						},
						"unit": {
							Required:     true,
							Type:         schema.TypeString,
							Description:  "Units of the applied Metric. Possible values are: PER_HOUR, PER_MINUTE, PER_SECOND, TOTAL.",
							ValidateFunc: validation.All(validation.StringInSlice([]string{"PER_HOUR", "PER_MINUTE", "PER_SECOND", "TOTAL"}, true)),
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
							Required:     true,
							Type:         schema.TypeString,
							Description:  "The zone where the VMs are created using this configuration.",
							ValidateFunc: validation.All(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
						},
						"cores": {
							Required:     true,
							Type:         schema.TypeInt,
							Description:  "The total number of cores for the VMs.",
							ValidateFunc: validation.All(validation.IntAtLeast(1)),
						},
						"cpu_family": {
							Optional:     true,
							Type:         schema.TypeString,
							Description:  "The zone where the VMs are created using this configuration.",
							ValidateFunc: validation.All(validation.StringInSlice([]string{"AMD_OPTERON", "INTEL_SKYLAKE", "INTEL_XEON"}, true)),
						},
						"nics": {
							Type:        schema.TypeList,
							Description: "List of NICs associated with this Replica.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lan": {
										Required:     true,
										Type:         schema.TypeInt,
										Description:  "Lan ID for this replica Nic.",
										ValidateFunc: validation.All(validation.IntAtLeast(1)),
									},
									"name": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "Name for this replica NIC.",
										ValidateFunc: validation.All(validation.StringLenBetween(0, 255)),
									},
									"dhcp": {
										Optional:    true,
										Type:        schema.TypeBool,
										Description: "Dhcp flag for this replica Nic. This is an optional attribute with default value of 'true' if not given in the request payload or given as null.",
									},
								}},
						},
						"ram": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
						},
						"volumes": {
							Type:        schema.TypeList,
							Description: "List of volumes associated with this Replica. Only a single volume is currently supported.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "The image installed on the volume. Only the UUID of the image is presently supported.",
										ValidateFunc: validation.All(validation.IsUUID),
									},
									"name": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "Name for this replica volume.",
										ValidateFunc: validation.All(validation.StringLenBetween(0, 255)),
									},
									"size": {
										Required:     true,
										Type:         schema.TypeInt,
										Description:  "User-defined size for this replica volume in GB.",
										ValidateFunc: validation.All(validation.IntAtLeast(1)),
									},
									"ssh_key_paths": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"ssh_key_values": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"ssh_keys": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
									"type": {
										Required:     true,
										Type:         schema.TypeString,
										Description:  "Storage Type for this replica volume. Possible values: SSD, HDD, SSD_STANDARD or SSD_PREMIUM",
										ValidateFunc: validation.All(validation.StringInSlice([]string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, true)),
									},
									"user_data": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "User-data (Cloud Init) for this replica volume.",
									},
									"image_password": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "Image password for this replica volume.",
									},
								}},
						},
					},
				},
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Description:  "Unique identifier for the resource",
				Required:     true,
				ValidateFunc: validation.All(validation.IsUUID),
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

func resourceAutoscalingGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).AutoscalingClient

	group, err := autoscalingService.GetAutoscalingGroupDataCreate(d)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured at getting data from the provided schema: %q", err))
		return diags
	}

	log.Printf("[DEBUG] autoscaling group data extracted: %+v", *group.Properties)

	autoscalingGroup, _, err := client.CreateGroup(ctx, *group)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating autoscaling group: %s", err))
		return diags
	}

	d.SetId(*autoscalingGroup.Id)
	log.Printf("[INFO] autoscaling Group created. Id set to %s", *autoscalingGroup.Id)

	if err := checkAction(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).AutoscalingClient

	group, apiResponse, err := client.GetGroup(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			log.Printf("[INFO] resource %s not found: %+v", d.Id(), err)
			d.SetId("")
			return nil
		} else {
			diags := diag.FromErr(fmt.Errorf("error while retrieving autoscaling group with id %v, %+v", d.Id(), err))
			return diags
		}
	}

	log.Printf("[INFO] successfully retreived autoscaling group %s: %+v", d.Id(), group)

	if err := autoscalingService.SetAutoscalingGroupData(d, group); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] successfully set autoscaling group data %s", d.Id())

	return nil
}

func resourceAutoscalingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).AutoscalingClient

	group, err := autoscalingService.GetAutoscalingGroupDataUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] autoscaling group data extracted: %+v", *group.Properties)

	_, _, err = client.UpdateGroup(ctx, d.Id(), *group)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating autoscaling group %s: %s", d.Id(), err))
		return diags
	}
	log.Printf("[INFO] autoscaling Group updated.")

	time.Sleep(SleepInterval * 20)

	if err := checkAction(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).AutoscalingClient

	_, err := client.DeleteGroup(ctx, d.Id())

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting an autoscaling group %s %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] autoscaling Group deleted: %s.", d.Id())

	d.SetId("")

	return nil
}

func resourceAutoscalingGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).AutoscalingClient

	groupId := d.Id()

	group, apiResponse, err := client.GetGroup(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find autoscaling group %q", groupId)
		}
		return nil, fmt.Errorf("an error occured while retrieveing autoscaling grouo %q, %q", groupId, err)
	}

	log.Printf("[INFO] autoscaling group found: %+v", group)

	if err := autoscalingService.SetAutoscalingGroupData(d, group); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func actionReady(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData, actionId string) (bool, error) {

	action, _, err := client.GetAction(ctx, d.Id(), actionId)

	if err != nil {
		return true, fmt.Errorf("error checking action status: %s", err)
	}

	if *action.Properties.ActionStatus == "FAILED" {
		return false, fmt.Errorf("action failed")
	}

	return *action.Properties.ActionStatus == "SUCCESSFUL", nil
}

// checkAction gets the triggered action and waits for it to be ready
func checkAction(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData) error {
	actions, _, err := client.GetAllActions(ctx, d.Id())

	if err != nil {
		return fmt.Errorf("error fetching group actions: %s", err)
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
			return fmt.Errorf("error while checking status of action %s: %s", actionId, rsErr)
		}

		if actionSuccessful {
			log.Printf("[INFO] action was ready: %s", actionId)
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			return fmt.Errorf("group creation/update timed out! WARNING: your group was created/updated but the action was not yet ready. " +
				"Check your Ionos Cloud account for updates")
		}
	}

	return nil
}
