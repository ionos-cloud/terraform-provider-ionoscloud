package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling/v2"

	as "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
)

// DataSourceAutoscalingGroup defines the schema for the Autoscaling Group data source
func DataSourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "UUID of the Autoscaling Group.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "User-defined name for the Autoscaling Group.",
				Optional:    true,
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
										Description: "When 'amountType=ABSOLUTE' specifies the absolute number of VMs that are removed. The value must be between 1 to 10. 'amountType=PERCENTAGE' specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200. At least one VM is always removed. Note that for 'SCALE_IN' operations, volumes are not deleted after the server is deleted.",
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
									"delete_volumes": {
										Computed:    true,
										Type:        schema.TypeBool,
										Description: "If set to 'true', when deleting an replica during scale in, any attached volume will also be deleted. When set to 'false', all volumes remain in the datacenter and must be deleted manually. Note that every scale-out creates new volumes. When they are not deleted, they will eventually use all of your contracts resource limits. At this point, scaling out would not be possible anymore.",
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
										Description: "When 'amountType=ABSOLUTE' specifies the absolute number of VMs that are added. The value must be between 1 to 10. 'amountType=PERCENTAGE' specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200. At least one VM is always added. Note that for 'SCALE_IN' operations, volumes are not deleted after the server is deleted.",
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
				Description: "VMs for this Autoscaling Group will be created in this Virtual Datacenter. Please note, that it has to have the same `location` as the `template`.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The zone where the VMs are created using this configuration.",
						},
						"cores": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The total number of cores for the VMs.",
						},
						"cpu_family": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The zone where the VMs are created using this configuration.",
						},
						"nic": {
							Type:        schema.TypeSet,
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
									"firewall_active": {
										Computed:    true,
										Type:        schema.TypeBool,
										Description: "Activate or deactivate the firewall. By default, an active firewall without any defined rules will block all incoming network traffic except for the firewall rules that explicitly allows certain protocols, IP addresses and ports.",
									},
									"firewall_type": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The type of firewall rules that will be allowed on the NIC. If not specified, the default INGRESS value is used.",
									},
									"flow_log": {
										Computed:    true,
										Type:        schema.TypeList,
										Description: "Flow log configuration for the NIC. By default, the flow log is inactive. If you want to activate the flow log, you must specify the target resource and the type of traffic to log.",
										Elem:        cloudapiflowlog.FlowlogSchemaDatasource,
									},
									"firewall_rule": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "List of all firewall rules for the specified NIC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Computed:    true,
													Type:        schema.TypeString,
													Description: "The protocol for the rule. The property cannot be modified after its creation (not allowed in update requests).",
												},
												"name": {
													Computed:    true,
													Type:        schema.TypeString,
													Description: "The name of the firewall rule.",
												},
												"source_mac": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Only traffic originating from the respective MAC address is permitted. Valid format: 'aa:bb:cc:dd:ee:ff'. The value 'null' allows traffic from any MAC address.",
												},
												"source_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Only traffic originating from the respective IPv4 address is permitted. The value 'null' allows traffic from any IP address.",
												},
												"target_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "If the target NIC has multiple IP addresses, only the traffic directed to the respective IP address of the NIC is allowed. The value 'null' allows traffic to any target IP address.",
												},
												"icmp_code": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Sets the allowed code (from 0 to 254) when ICMP protocol is selected. The value 'null' allows all codes.",
												},
												"icmp_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Sets the allowed type (from 0 to 254) if the protocol ICMP is selected. The value 'null' allows all types.",
												},
												"port_range_start": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Sets the initial range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value 'null' for 'port_range_start' and 'port_range_end' allows all ports.",
												},
												"port_range_end": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Sets the end range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value 'null' for 'port_range_start' and 'port_range_end' allows all ports.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The firewall rule type. If not specified, the default value 'INGRESS' is used.",
												},
											},
										},
									},
									"target_group": {
										Computed:    true,
										Type:        schema.TypeList,
										Description: "In order to link VM to ALB, target group must be provided.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the target group.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The port for the target group.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The weight for the target group.",
												},
											},
										},
									},
								}},
						},
						"ram": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
						},
						"volume": {
							Type:        schema.TypeSet,
							Description: "List of volumes associated with this Replica. Only a single volume is currently supported.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The image installed on the volume. Only the UUID of the image is presently supported.",
									},
									"image_alias": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The image installed on the volume. Must be an 'imageAlias' as specified via the images API.",
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
									"boot_order": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: `Determines whether the volume will be used as a boot volume: NONE - the volume will not be used as boot volume, PRIMARY - the volume will be used as boot volume, AUTO - will delegate the decision to the provisioning engine to decide whether to use the volume as boot volume.`,
									},
									"bus": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: `The bus type of the volume. Default setting is 'VIRTIO'. The bus type 'IDE' is also supported.`,
									},
									"backup_unit_id": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "The uuid of the Backup Unit that user has access to.",
									},
								}},
						},
					},
				},
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "Unique identifier for the resource",
				Computed:    true,
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

func dataSourceAutoscalingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).AutoscalingClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if err := validateIDAndName(idOk, nameOk); err != nil {
		return diag.FromErr(err)
	}

	var (
		group *autoscaling.Group
		err   error
	)

	if idOk {
		group, err = fetchGroupByID(ctx, client, id.(string))
	} else {
		group, err = fetchGroupByName(ctx, client, name.(string))
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.Id)
	if err = setAutoscalingGroupData(d, &group.Properties); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func validateIDAndName(idOk, nameOk bool) error {
	if idOk && nameOk {
		return fmt.Errorf("id and name cannot be both specified at the same time")
	}
	if !idOk && !nameOk {
		return fmt.Errorf("please provide either the group id or name")
	}
	return nil
}

func fetchGroupByID(ctx context.Context, client *as.Client, id string) (*autoscaling.Group, error) {
	group, _, err := client.GetGroup(ctx, id, 2)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while fetching group with ID %s: %w", id, err)
	}
	return &group, nil
}

func fetchGroupByName(ctx context.Context, client *as.Client, name string) (*autoscaling.Group, error) {
	groups, _, err := client.ListGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting groups: %w", err)
	}

	if groups.Items == nil {
		return nil, fmt.Errorf("no group found")
	}

	for _, g := range groups.Items {
		tmpGroup, _, err := client.GetGroup(ctx, g.Id, 2)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching group %s: %w", g.Id, err)
		}

		if tmpGroup.Properties.Name != nil && strings.EqualFold(*tmpGroup.Properties.Name, name) {
			return &tmpGroup, nil
		}
	}

	return nil, fmt.Errorf("no group found with the specified criteria: name = %s", name)
}
