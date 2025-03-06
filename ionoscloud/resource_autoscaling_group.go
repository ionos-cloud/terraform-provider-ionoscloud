package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// ResourceAutoscalingGroup defines the schema for the Autoscaling Group resource
func ResourceAutoscalingGroup() *schema.Resource {
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
			// will be left until GA in case it is added again in swagger
			// "target_replica_count": {
			// 	Type:             schema.TypeInt,
			//	Description:      "The target number of VMs in this Group. Depending on the scaling policy, this number will be adjusted automatically. VMs will be created or destroyed automatically in order to adjust the actual number of VMs to this number. If targetReplicaCount is given in the request body then it must be >= minReplicaCount and <= maxReplicaCount.",
			//	Optional:         true,
			//	Computed:         true,
			//	ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 200)),
			// },
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
									"delete_volumes": {
										Required:    true,
										Type:        schema.TypeBool,
										Description: "If set to 'true', when deleting an replica during scale in, any attached volume will also be deleted. When set to 'false', all volumes remain in the datacenter and must be deleted manually. Note that every scale-out creates new volumes. When they are not deleted, they will eventually use all of your contracts resource limits. At this point, scaling out would not be possible anymore.",
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
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The minimum time that elapses after the start of this scaling action until the following scaling action is started. While a scaling action is in progress, no second action is initiated for the same VM Auto Scaling Group. Instead, the metric is re-evaluated after the current scaling action completes (either successfully or with errors). This is currently validated with a minimum value of 2 minutes and a maximum of 24 hours. The default value is 5 minutes if not specified.",
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
										Computed:    true,
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
						"range": {
							Optional:    true,
							Default:     "PT2M",
							Type:        schema.TypeString,
							Description: "Specifies the time range for which the samples are to be aggregated. Must be >= 2 minutes.",
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
						"ram": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
						},
						"nic": {
							Type:        schema.TypeList,
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
									"firewall_active": {
										Optional:    true,
										Type:        schema.TypeBool,
										Description: "Activate or deactivate the firewall. By default, an active firewall without any defined rules will block all incoming network traffic except for the firewall rules that explicitly allows certain protocols, IP addresses and ports.",
									},
									"firewall_type": {
										Optional:         true,
										Type:             schema.TypeString,
										Default:          "INGRESS",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INGRESS", "EGRESS", "BIDIRECTIONAL"}, true)),
										Description:      "The type of firewall rules that will be allowed on the NIC. If not specified, the default INGRESS value is used.",
									},
									"flow_log": {
										Optional:    true,
										Type:        schema.TypeList,
										Description: "List of all flow logs for the specified NIC.",
										Elem:        cloudapiflowlog.FlowlogSchemaResource,
									},
									"firewall_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of all firewall rules for the specified NIC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:             schema.TypeString,
													Required:         true,
													Description:      "The protocol for the rule. The property cannot be modified after its creation (not allowed in update requests).",
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"TCP", "UDP", "ICMP", "ANY"}, false)),
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the firewall rule.",
												},
												"source_mac": {
													Type:             schema.TypeString,
													Optional:         true,
													Description:      "Only traffic originating from the respective MAC address is permitted. Valid format: 'aa:bb:cc:dd:ee:ff'. The value 'null' allows traffic from any MAC address.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"), "Invalid MAC address format")),
												},
												"source_ip": {
													Type:             schema.TypeString,
													Optional:         true,
													Description:      "Only traffic originating from the respective IPv4 address is permitted. The value 'null' allows traffic from any IP address.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"), "Invalid IP address format")),
												},
												"target_ip": {
													Type:             schema.TypeString,
													Optional:         true,
													Description:      "If the target NIC has multiple IP addresses, only the traffic directed to the respective IP address of the NIC is allowed. The value 'null' allows traffic to any target IP address.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"), "Invalid IP address format")),
												},
												"icmp_code": {
													Type:             schema.TypeInt,
													Optional:         true,
													Description:      "Sets the allowed code (from 0 to 254) when ICMP protocol is selected. The value 'null' allows all codes.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 254)),
												},
												"icmp_type": {
													Type:             schema.TypeInt,
													Optional:         true,
													Description:      "Sets the allowed type (from 0 to 254) if the protocol ICMP is selected. The value 'null' allows all types.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 254)),
												},
												"port_range_start": {
													Type:             schema.TypeInt,
													Optional:         true,
													Description:      "Sets the initial range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value 'null' for 'port_range_start' and 'port_range_end' allows all ports.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
												},
												"port_range_end": {
													Type:             schema.TypeInt,
													Optional:         true,
													Description:      "Sets the end range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value 'null' for 'port_range_start' and 'port_range_end' allows all ports.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
												},
												"type": {
													Type:             schema.TypeString,
													Optional:         true,
													Computed:         true,
													Description:      "The firewall rule type. If not specified, the default value 'INGRESS' is used.",
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INGRESS", "EGRESS"}, false)),
												},
											},
										},
									},
									"target_group": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "In order to link VM to ALB, target group must be provided.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_group_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The ID of the target group.",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The port for the target group.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The weight for the target group.",
												},
											},
										},
									},
								}},
						},
						"volume": {
							Type:        schema.TypeSet,
							Description: "List of volumes associated with this Replica.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"type": {
										Required:         true,
										Type:             schema.TypeString,
										Description:      "Storage Type for this replica volume. Possible values: SSD, HDD, SSD_STANDARD or SSD_PREMIUM",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, true)),
									},
									"boot_order": {
										Required: true,
										Type:     schema.TypeString,
										Description: `Determines whether the volume will be used as a boot volume. Set to NONE, the volume will not be used as boot volume. 
Set to PRIMARY, the volume will be used as boot volume and set to AUTO will delegate the decision to the provisioning engine to decide whether to use the volume as boot volume.
Notice that exactly one volume can be set to PRIMARY or all of them set to AUTO.`,
									},
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
									"ssh_keys": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
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
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(utils.DefaultTimeout),
			Update:  schema.DefaultTimeout(utils.DefaultTimeout),
			Delete:  schema.DefaultTimeout(utils.DefaultTimeout),
			Default: schema.DefaultTimeout(utils.DefaultTimeout),
		},
	}
}

func resourceAutoscalingGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient

	var group autoscaling.GroupPost
	properties, err := expandProperties(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while expanding properties: %w", err))
	}
	group.Properties = *properties
	autoscalingGroup, _, err := client.CreateGroup(ctx, group)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Autoscaling Group: %w", err))
	}

	d.SetId(autoscalingGroup.Id)
	log.Printf("[INFO] Autoscaling Group created. Id set to %s", autoscalingGroup.Id)

	if err := checkAction(ctx, client, d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient
	group, apiResponse, err := client.GetGroup(ctx, d.Id(), 2)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Autoscaling Group with ID: %s not found, err: %+v", d.Id(), err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while retrieving Autoscaling Group with id %v, %w", d.Id(), err))
	}

	log.Printf("[INFO] successfully retrieved Autoscaling Group %s: %+v", d.Id(), group)
	if err := setAutoscalingGroupData(d, &group.Properties); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] successfully set Autoscaling Group data %s", d.Id())

	return nil
}

func resourceAutoscalingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient

	if d.HasChange("datacenter_id") {
		return diag.FromErr(fmt.Errorf("datacenter_id property is immutable and can be used only in create requests"))
	}

	replicaConfiguration, err := expandReplicaConfiguration(d.Get("replica_configuration").([]any))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while expanding replica configuration: %w", err))
	}

	group := autoscaling.GroupPut{
		Properties: autoscaling.GroupPutProperties{
			MaxReplicaCount:      int64(d.Get("max_replica_count").(int)),
			MinReplicaCount:      int64(d.Get("min_replica_count").(int)),
			Name:                 d.Get("name").(string),
			Policy:               *expandPolicy(d.Get("policy").([]any)),
			ReplicaConfiguration: *replicaConfiguration,
			Datacenter: autoscaling.GroupPutPropertiesDatacenter{
				Id: d.Get("datacenter_id").(string),
			},
		},
	}

	if _, _, err = client.UpdateGroup(ctx, d.Id(), group); err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating Autoscaling Group %s: %w", d.Id(), err))
	}

	log.Printf("[INFO] Autoscaling Group updated.")

	if err := checkAction(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceAutoscalingGroupRead(ctx, d, meta)
}

func resourceAutoscalingGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient
	if _, err := client.DeleteGroup(ctx, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while deleting an Autoscaling Group %s %w", d.Id(), err))
	}

	log.Printf("[INFO] Autoscaling Group deleted: %s.", d.Id())

	d.SetId("")

	return nil
}

func resourceAutoscalingGroupImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).AutoscalingClient

	groupID := d.Id()

	group, apiResponse, err := client.GetGroup(ctx, d.Id(), 0)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Autoscaling Group %q", groupID)
		}
		return nil, fmt.Errorf("an error occurred while retrieving Autoscaling Group %q, %w", groupID, err)
	}

	log.Printf("[INFO] Autoscaling Group found: %+v", group)

	if err := setAutoscalingGroupData(d, &group.Properties); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func expandProperties(d *schema.ResourceData) (*autoscaling.GroupProperties, error) {
	replicaConfiguration, err := expandReplicaConfiguration(d.Get("replica_configuration").([]any))
	if err != nil {
		return nil, err
	}

	return &autoscaling.GroupProperties{
		MaxReplicaCount:      shared.ToPtr(int64(d.Get("max_replica_count").(int))),
		MinReplicaCount:      shared.ToPtr(int64(d.Get("min_replica_count").(int))),
		Name:                 shared.ToPtr(d.Get("name").(string)),
		Policy:               expandPolicy(d.Get("policy").([]any)),
		ReplicaConfiguration: replicaConfiguration,
		Datacenter: autoscaling.GroupPropertiesDatacenter{
			Id: d.Get("datacenter_id").(string),
		},
	}, nil
}

func expandPolicy(l []any) *autoscaling.GroupPolicy {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	s := l[0].(map[string]interface{})

	// required fields
	metric := autoscaling.Metric(s["metric"].(string))
	scaleInAction := expandScaleInAction(s["scale_in_action"].([]any))
	scaleInThreshold := float32(s["scale_in_threshold"].(int))
	scaleOutAction := expandScaleOutAction(s["scale_out_action"].([]any))
	scaleOutThreshold := float32(s["scale_out_threshold"].(int))
	unit := autoscaling.QueryUnit(s["unit"].(string))
	out := &autoscaling.GroupPolicy{
		Metric:            metric,
		ScaleInAction:     *scaleInAction,
		ScaleInThreshold:  scaleInThreshold,
		ScaleOutAction:    *scaleOutAction,
		ScaleOutThreshold: scaleOutThreshold,
		Unit:              unit,
	}

	// optional fields
	if v, ok := s["range"]; ok {
		out.SetRange(v.(string))
	}

	return out
}

func expandReplicaConfiguration(l []any) (*autoscaling.ReplicaPropertiesPost, error) {
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	s := l[0].(map[string]interface{})

	// required fields
	availabilityZone := autoscaling.AvailabilityZone(s["availability_zone"].(string))
	cores := int32(s["cores"].(int))
	ram := int32(s["ram"].(int))
	out := &autoscaling.ReplicaPropertiesPost{
		AvailabilityZone: &availabilityZone,
		Cores:            cores,
		Ram:              ram,
	}

	// optional fields
	if v, ok := s["cpu_family"]; ok {
		out.SetCpuFamily(autoscaling.CpuFamily(v.(string)))
	}

	if v, ok := s["nic"]; ok {
		out.SetNics(expandNICs(v.([]any)))
	}

	if v, ok := s["volume"]; ok {
		volumes, err := expandVolumes(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
		out.SetVolumes(volumes)
	}

	return out, nil
}

func expandScaleInAction(l []any) *autoscaling.GroupPolicyScaleInAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	s := l[0].(map[string]interface{})

	// required fields
	amount := float32(s["amount"].(int))
	amountType := autoscaling.ActionAmount(s["amount_type"].(string))
	cooldownPeriod := s["cooldown_period"].(string)
	deleteVolumes := s["delete_volumes"].(bool)
	out := &autoscaling.GroupPolicyScaleInAction{
		Amount:         amount,
		AmountType:     amountType,
		CooldownPeriod: &cooldownPeriod,
		DeleteVolumes:  deleteVolumes,
	}

	// optional fields
	if v, ok := s["termination_policy_type"]; ok {
		out.SetTerminationPolicy(autoscaling.TerminationPolicyType(v.(string)))
	}

	return out
}

func expandScaleOutAction(l []any) *autoscaling.GroupPolicyScaleOutAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	s := l[0].(map[string]interface{})

	// required fields
	amount := float32(s["amount"].(int))
	amountType := autoscaling.ActionAmount(s["amount_type"].(string))
	cooldownPeriod := s["cooldown_period"].(string)

	return &autoscaling.GroupPolicyScaleOutAction{
		Amount:         amount,
		AmountType:     amountType,
		CooldownPeriod: &cooldownPeriod,
	}
}

func expandVolumes(l []any) ([]autoscaling.ReplicaVolumePost, error) {
	volumes := make([]autoscaling.ReplicaVolumePost, len(l))
	for i, entry := range l {
		s := entry.(map[string]interface{})

		// required fields
		name := s["name"].(string)
		size := int32(s["size"].(int))
		typ := autoscaling.VolumeHwType(s["type"].(string))
		bootOrder := s["boot_order"].(string)
		volumes[i] = autoscaling.ReplicaVolumePost{
			Name:      name,
			Size:      size,
			Type:      typ,
			BootOrder: bootOrder,
		}

		// optional fields
		if v, ok := s["image"]; ok {
			volumes[i].Image = shared.ToPtr(v.(string))
		}

		if v, ok := s["image_alias"]; ok {
			volumes[i].ImageAlias = shared.ToPtr(v.(string))
		}

		if volumes[i].Image == nil && volumes[i].ImageAlias == nil {
			return nil, fmt.Errorf("it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property")
		}

		if v, ok := s["ssh_keys"]; ok {
			keys, err := expandSSHKeys(v.([]any))
			if err != nil {
				return nil, err
			}
			volumes[i].SshKeys = keys
		}

		if v, ok := s["user_data"]; ok {
			volumes[i].UserData = shared.ToPtr(v.(string))
		}

		if v, ok := s["image_password"]; ok {
			volumes[i].ImagePassword = shared.ToPtr(v.(string))
		}

		if v, ok := s["bus"]; ok {
			volumes[i].SetBus(autoscaling.BusType(v.(string)))
		}

		if v, ok := s["backup_unit_id"]; ok {
			volumes[i].BackupunitId = shared.ToPtr(v.(string))
		}
	}

	return volumes, nil
}

func expandSSHKeys(l []any) ([]string, error) {
	sshKeys := make([]string, len(l))
	for i, entry := range l {
		pubKey, err := utils.ReadPublicKey(entry.(string))
		if err != nil {
			return nil, fmt.Errorf("error reading sshkey (%s) (%w)", entry, err)
		}

		sshKeys[i] = pubKey
	}

	return sshKeys, nil
}

func expandNICs(l []any) []autoscaling.ReplicaNic {
	nics := make([]autoscaling.ReplicaNic, len(l))
	for i, entry := range l {
		s := entry.(map[string]interface{})

		// required fields
		lan := int32(s["lan"].(int))
		name := s["name"].(string)
		nics[i] = autoscaling.ReplicaNic{
			Lan:  lan,
			Name: name,
		}

		// optional fields
		if v, ok := s["dhcp"]; ok {
			nics[i].SetDhcp(v.(bool))
		}

		if v, ok := s["firewall_active"]; ok {
			nics[i].SetFirewallActive(v.(bool))
		}

		if v, ok := s["firewall_type"]; ok {
			nics[i].SetFirewallType(v.(string))
		}

		if v, ok := s["flow_log"]; ok {
			nics[i].SetFlowLogs(expandFlowLogs(v.([]any)))
		}

		if v, ok := s["firewall_rule"]; ok {
			nics[i].SetFirewallRules(expandFirewallRules(v.([]any)))
		}

		if v, ok := s["target_group"]; ok {
			nics[i].TargetGroup = expandTargetGroup(v.([]any))
		}
	}

	return nics
}

func expandFlowLogs(l []any) []autoscaling.NicFlowLog {
	flowLogs := make([]autoscaling.NicFlowLog, len(l))
	for i, entry := range l {
		s := entry.(map[string]interface{})

		// all fields are required
		flowLogs[i] = autoscaling.NicFlowLog{
			Name:      s["name"].(string),
			Action:    s["action"].(string),
			Direction: s["direction"].(string),
			Bucket:    s["bucket"].(string),
		}
	}

	return flowLogs
}

func expandTargetGroup(l []any) *autoscaling.TargetGroup {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	s := l[0].(map[string]interface{})

	// required fields
	targetGroupID := s["target_group_id"].(string)
	port := int32(s["port"].(int))
	weight := int32(s["weight"].(int))
	return &autoscaling.TargetGroup{
		TargetGroupId: targetGroupID,
		Port:          port,
		Weight:        weight,
	}
}

func expandFirewallRules(l []any) []autoscaling.NicFirewallRule {
	rules := make([]autoscaling.NicFirewallRule, len(l))
	for i, entry := range l {
		s := entry.(map[string]interface{})

		// required fields
		rules[i] = autoscaling.NicFirewallRule{
			Protocol: s["protocol"].(string),
		}

		// optional fields
		if v, ok := s["name"]; ok {
			rules[i].SetName(v.(string))
		}

		if v, ok := s["source_mac"]; ok && v != "" {
			rules[i].SetSourceMac(v.(string))
		}

		if v, ok := s["source_ip"]; ok && v != "" {
			rules[i].SetSourceIp(v.(string))
		}

		if v, ok := s["target_ip"]; ok && v != "" {
			rules[i].SetTargetIp(v.(string))
		}

		if v, ok := s["icmp_code"]; ok && v != 0 {
			rules[i].SetIcmpCode(int32(v.(int)))
		}

		if v, ok := s["icmp_type"]; ok && v != 0 {
			rules[i].SetIcmpType(int32(v.(int)))
		}

		if v, ok := s["port_range_start"]; ok && v != 0 {
			rules[i].SetPortRangeStart(int32(v.(int)))
		}

		if v, ok := s["port_range_end"]; ok && v != 0 {
			rules[i].SetPortRangeEnd(int32(v.(int)))
		}

		if v, ok := s["type"]; ok {
			rules[i].SetType(v.(string))
		}
	}

	return rules
}

func setAutoscalingGroupData(d *schema.ResourceData, groupProperties *autoscaling.GroupProperties) error {

	resourceName := "autoscaling groupProperties"
	if groupProperties != nil {
		if groupProperties.MaxReplicaCount != nil {
			if err := d.Set("max_replica_count", *groupProperties.MaxReplicaCount); err != nil {
				return utils.GenerateSetError(resourceName, "max_replica_count", err)
			}
		}

		if groupProperties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *groupProperties.MinReplicaCount); err != nil {
				return utils.GenerateSetError(resourceName, "min_replica_count", err)
			}
		}

		// if groupProperties.TargetReplicaCount != nil {
		//	if err := d.Set("target_replica_count", *groupProperties.TargetReplicaCount); err != nil {
		//		return utils.GenerateSetError(resourceName, "target_replica_count", err)
		//	}
		//}

		if groupProperties.Name != nil {
			if err := d.Set("name", *groupProperties.Name); err != nil {
				return utils.GenerateSetError(resourceName, "name", err)
			}
		}

		if groupProperties.Policy != nil {
			if err := d.Set("policy", flattenPolicyProperties(groupProperties.Policy)); err != nil {
				return utils.GenerateSetError(resourceName, "policy", err)
			}
		}

		if groupProperties.ReplicaConfiguration != nil {
			if err := d.Set("replica_configuration", flattenReplicaConfiguration(d, groupProperties.ReplicaConfiguration)); err != nil {
				return utils.GenerateSetError(resourceName, "replica_configuration", err)
			}
		}

		if err := d.Set("datacenter_id", groupProperties.Datacenter.Id); err != nil {
			return utils.GenerateSetError(resourceName, "datacenter_id", err)
		}

		if err := d.Set("location", groupProperties.Location); err != nil {
			return utils.GenerateSetError(resourceName, "location", err)
		}
	}
	return nil
}

func flattenPolicyProperties(gp *autoscaling.GroupPolicy) []map[string]any {
	if gp == nil {
		groupPolicies := make([]map[string]any, 0)
		return groupPolicies
	}

	groupPolicies := make([]map[string]any, 1)
	policy := map[string]any{}
	utils.SetPropWithNilCheck(policy, "metric", gp.Metric)
	utils.SetPropWithNilCheck(policy, "range", gp.Range)
	utils.SetPropWithNilCheck(policy, "scale_in_threshold", gp.ScaleInThreshold)
	utils.SetPropWithNilCheck(policy, "scale_out_threshold", gp.ScaleOutThreshold)
	utils.SetPropWithNilCheck(policy, "unit", gp.Unit)
	utils.SetPropWithNilCheck(policy, "scale_in_action", flattenScaleInActionProperties(gp.ScaleInAction))
	utils.SetPropWithNilCheck(policy, "scale_out_action", flattenScaleOutActionProperties(gp.ScaleOutAction))

	groupPolicies[0] = policy
	return groupPolicies
}

func flattenScaleInActionProperties(scaleInAction autoscaling.GroupPolicyScaleInAction) []map[string]any {
	scaleInActions := make([]map[string]any, 1)
	scaleIn := map[string]any{}

	utils.SetPropWithNilCheck(scaleIn, "amount", scaleInAction.Amount)
	utils.SetPropWithNilCheck(scaleIn, "amount_type", scaleInAction.AmountType)
	utils.SetPropWithNilCheck(scaleIn, "termination_policy_type", scaleInAction.TerminationPolicy)
	utils.SetPropWithNilCheck(scaleIn, "cooldown_period", scaleInAction.CooldownPeriod)
	utils.SetPropWithNilCheck(scaleIn, "delete_volumes", scaleInAction.DeleteVolumes)

	scaleInActions[0] = scaleIn
	return scaleInActions
}

func flattenScaleOutActionProperties(scaleOutAction autoscaling.GroupPolicyScaleOutAction) []map[string]any {
	scaleOutActions := make([]map[string]any, 1)
	scaleOut := map[string]any{}

	utils.SetPropWithNilCheck(scaleOut, "amount", scaleOutAction.Amount)
	utils.SetPropWithNilCheck(scaleOut, "amount_type", scaleOutAction.AmountType)
	utils.SetPropWithNilCheck(scaleOut, "cooldown_period", scaleOutAction.CooldownPeriod)

	scaleOutActions[0] = scaleOut
	return scaleOutActions
}

func flattenReplicaConfiguration(d *schema.ResourceData, replicaConfiguration *autoscaling.ReplicaPropertiesPost) []map[string]any {
	if replicaConfiguration == nil {
		replicas := make([]map[string]any, 0)
		return replicas
	}
	replicas := make([]map[string]any, 1)
	replica := map[string]any{}

	utils.SetPropWithNilCheck(replica, "availability_zone", replicaConfiguration.AvailabilityZone)
	utils.SetPropWithNilCheck(replica, "cores", replicaConfiguration.Cores)
	utils.SetPropWithNilCheck(replica, "cpu_family", replicaConfiguration.CpuFamily)
	utils.SetPropWithNilCheck(replica, "ram", replicaConfiguration.Ram)
	utils.SetPropWithNilCheck(replica, "nic", flattenNIC(replicaConfiguration.Nics))
	utils.SetPropWithNilCheck(replica, "volume", flattenVolume(d, replicaConfiguration.Volumes))

	replicas[0] = replica
	return replicas
}

func flattenNIC(replicaNICs []autoscaling.ReplicaNic) []map[string]any {
	if replicaNICs == nil {
		nics := make([]map[string]any, 0)
		return nics
	}

	nics := make([]map[string]any, len(replicaNICs))
	for i, nic := range replicaNICs {
		trNIC := map[string]any{}
		utils.SetPropWithNilCheck(trNIC, "lan", nic.Lan)
		utils.SetPropWithNilCheck(trNIC, "name", nic.Name)
		utils.SetPropWithNilCheck(trNIC, "dhcp", nic.Dhcp)
		utils.SetPropWithNilCheck(trNIC, "firewall_active", nic.FirewallActive)
		utils.SetPropWithNilCheck(trNIC, "firewall_type", nic.FirewallType)
		utils.SetPropWithNilCheck(trNIC, "firewall_rule", flattenFirewallRules(nic.FirewallRules))
		utils.SetPropWithNilCheck(trNIC, "flow_log", flattenFlowLogs(nic.FlowLogs))
		utils.SetPropWithNilCheck(trNIC, "target_group", flattenTargetGroup(nic.TargetGroup))
		nics[i] = trNIC
	}
	return nics
}

func flattenFirewallRules(rules []autoscaling.NicFirewallRule) []map[string]any {
	if rules == nil {
		firewallRules := make([]map[string]any, 0)
		return firewallRules
	}

	firewallRules := make([]map[string]any, len(rules))
	for i, rule := range rules {
		trRule := map[string]any{}
		utils.SetPropWithNilCheck(trRule, "name", rule.Name)
		utils.SetPropWithNilCheck(trRule, "protocol", rule.Protocol)
		utils.SetPropWithNilCheck(trRule, "icmp_type", rule.IcmpType)
		utils.SetPropWithNilCheck(trRule, "icmp_code", rule.IcmpCode)
		utils.SetPropWithNilCheck(trRule, "port_range_start", rule.PortRangeStart)
		utils.SetPropWithNilCheck(trRule, "port_range_end", rule.PortRangeEnd)
		utils.SetPropWithNilCheck(trRule, "source_mac", rule.SourceMac)
		utils.SetPropWithNilCheck(trRule, "source_ip", rule.SourceIp)
		utils.SetPropWithNilCheck(trRule, "target_ip", rule.TargetIp)
		utils.SetPropWithNilCheck(trRule, "type", rule.Type)
		firewallRules[i] = trRule
	}
	return firewallRules
}

func flattenFlowLogs(logs []autoscaling.NicFlowLog) []map[string]any {
	if logs == nil {
		flowLogs := make([]map[string]any, 0)
		return flowLogs
	}

	flowLogs := make([]map[string]any, len(logs))
	for i, flog := range logs {
		trLog := map[string]any{}
		utils.SetPropWithNilCheck(trLog, "name", flog.Name)
		utils.SetPropWithNilCheck(trLog, "action", flog.Action)
		utils.SetPropWithNilCheck(trLog, "direction", flog.Direction)
		utils.SetPropWithNilCheck(trLog, "bucket", flog.Bucket)
		flowLogs[i] = trLog
	}
	return flowLogs
}

func flattenTargetGroup(tg *autoscaling.TargetGroup) []map[string]any {
	if tg == nil {
		targetGroups := make([]map[string]any, 0)
		return targetGroups
	}

	targetGroups := make([]map[string]any, 1)
	targetGroup := map[string]any{}
	utils.SetPropWithNilCheck(targetGroup, "target_group_id", tg.TargetGroupId)
	utils.SetPropWithNilCheck(targetGroup, "port", tg.Port)
	utils.SetPropWithNilCheck(targetGroup, "weight", tg.Weight)

	targetGroups[0] = targetGroup
	return targetGroups
}

func flattenVolume(d *schema.ResourceData, replicaVolumes []autoscaling.ReplicaVolumePost) []map[string]any {
	if replicaVolumes == nil {
		volumes := make([]map[string]any, 0)
		return volumes
	}

	volumes := make([]map[string]any, len(replicaVolumes))
	for i, volume := range replicaVolumes {
		trVolume := map[string]any{}
		utils.SetPropWithNilCheck(trVolume, "name", volume.Name)
		utils.SetPropWithNilCheck(trVolume, "image_alias", volume.ImageAlias)
		utils.SetPropWithNilCheck(trVolume, "image", volume.Image)
		utils.SetPropWithNilCheck(trVolume, "size", volume.Size)
		utils.SetPropWithNilCheck(trVolume, "type", volume.Type)
		utils.SetPropWithNilCheck(trVolume, "bus", volume.Bus)
		utils.SetPropWithNilCheck(trVolume, "boot_order", volume.BootOrder)
		// we need to take these from schema as they are not returned by API
		volumeMap, ok := d.GetOk("replica_configuration.0.volume")
		if ok {
			volumeMap := (volumeMap).(*schema.Set).List()[i].(map[string]any)
			trVolume["image_password"] = volumeMap["image_password"]
			trVolume["ssh_keys"] = volumeMap["ssh_keys"]
			trVolume["user_data"] = volumeMap["user_data"]
		}
		volumes[i] = trVolume
	}
	return volumes
}
func actionReady(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData, actionID string) (bool, error) {

	action, _, err := client.GetAction(ctx, d.Id(), actionID)
	if err != nil {
		return true, fmt.Errorf("error checking action status: %w", err)
	}

	if action.Properties.ActionStatus == autoscaling.ACTIONSTATUS_FAILED {
		return false, fmt.Errorf("action failed")
	}

	if action.Properties == nil {
		return false, fmt.Errorf("expected a value for the action status but received 'nil' instead")
	}
	return strings.EqualFold(string(action.Properties.ActionStatus), string(autoscaling.ACTIONSTATUS_SUCCESSFUL)), nil
}

// checkAction gets the triggered action and waits for it to be ready
func checkAction(ctx context.Context, client *autoscalingService.Client, d *schema.ResourceData) error {
	actions, _, err := client.GetAllActions(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("error fetching group actions: %w", err)
	}

	if actions.Items == nil {
		return fmt.Errorf("no action triggered for group: %s", d.Id())
	}

	if len(actions.Items) == 0 {
		return fmt.Errorf("no action triggered for group: %s", d.Id())
	}

	actionID := (actions.Items)[0].Id

	// wait for completion of triggered action
	for {
		log.Printf("[INFO] waiting for action %s to be ready...", actionID)

		actionSuccessful, rsErr := actionReady(ctx, client, d, actionID)
		if rsErr != nil {
			return fmt.Errorf("error while checking status of action %s: %w", actionID, rsErr)
		}

		if actionSuccessful {
			log.Printf("[INFO] action was ready: %s", actionID)
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
