package autoscaling

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type GroupService interface {
	GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error)
	ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error)
	CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error)
	UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *autoscaling.APIResponse, error)
	DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error)
}

func (c *Client) GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsFindById(ctx, groupId).Depth(2).Execute()
	if apiResponse != nil {
		return group, apiResponse, err

	}
	return group, nil, err
}

func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error) {
	groups, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsGet(ctx).Execute()
	if apiResponse != nil {
		return groups, apiResponse, err
	}
	return groups, nil, err
}

func (c *Client) CreateGroup(ctx context.Context, group autoscaling.GroupPost) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPost(ctx).GroupPost(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPut(ctx, groupId).GroupPut(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error) {
	apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsDelete(ctx, groupId).Execute()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func GetAutoscalingGroupDataCreate(d *schema.ResourceData) (*autoscaling.GroupPost, error) {

	group := autoscaling.GroupPost{
		Properties: &autoscaling.GroupProperties{},
	}

	if value, ok := d.GetOk("max_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MaxReplicaCount = &value
	}

	if value, ok := d.GetOk("min_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MinReplicaCount = &value
	}

	//if value, ok := d.GetOk("target_replica_count"); ok {
	//	value := int64(value.(int))
	//	group.Properties.TargetReplicaCount = &value
	//}

	if value, ok := d.GetOk("name"); ok {
		value := value.(string)
		group.Properties.Name = &value
	}

	group.Properties.Policy = GetAutoscalingGroupPolicyData(d)

	if replicaConfiguration, err := GetReplicaConfigurationPostData(d); err != nil {
		return nil, err
	} else {
		group.Properties.ReplicaConfiguration = replicaConfiguration
	}

	if value, ok := d.GetOk("datacenter_id"); ok {
		value := value.(string)
		datacenter := autoscaling.GroupPropertiesDatacenter{
			Id: &value,
		}
		group.Properties.Datacenter = &datacenter
	}
	return &group, nil
}

func GetAutoscalingGroupDataUpdate(d *schema.ResourceData) (*autoscaling.GroupPut, error) {

	group := autoscaling.GroupPut{
		Properties: &autoscaling.GroupPutProperties{},
	}

	if value, ok := d.GetOk("max_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MaxReplicaCount = &value
	}

	if value, ok := d.GetOk("min_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MinReplicaCount = &value
	}

	//if value, ok := d.GetOk("target_replica_count"); ok {
	//	value := int64(value.(int))
	//	group.Properties.TargetReplicaCount = &value
	//} else {
	//	group.Properties.TargetReplicaCount = nil
	//}

	if value, ok := d.GetOk("name"); ok {
		value := value.(string)
		group.Properties.Name = &value
	}

	group.Properties.Policy = GetAutoscalingGroupPolicyData(d)

	if replicaConfiguration, err := GetReplicaConfigurationPostData(d); err != nil {
		return nil, err
	} else {
		group.Properties.ReplicaConfiguration = replicaConfiguration
	}

	if d.HasChange("datacenter_id") {
		return nil, fmt.Errorf("datacenter_id property is immutable and can be used only in create requests")
	} else {
		if value, ok := d.GetOk("datacenter_id"); ok {
			value := value.(string)
			datacenter := autoscaling.GroupPutPropertiesDatacenter{
				Id: &value,
			}
			group.Properties.Datacenter = &datacenter
		}
	}
	return &group, nil
}

func GetAutoscalingGroupPolicyData(d *schema.ResourceData) *autoscaling.GroupPolicy {

	groupPolicy := autoscaling.GroupPolicy{}

	if value, ok := d.GetOk("policy.0.metric"); ok {
		value := autoscaling.Metric(value.(string))
		groupPolicy.Metric = &value
	}

	if value, ok := d.GetOk("policy.0.range"); ok {
		value := value.(string)
		groupPolicy.Range = &value
	} else {
		groupPolicy.Range = nil
	}

	groupPolicy.ScaleInAction = GetScaleInActionData(d)

	if value, ok := d.GetOk("policy.0.scale_in_threshold"); ok {
		value := float32(value.(int))
		groupPolicy.ScaleInThreshold = &value
	}

	groupPolicy.ScaleOutAction = GetScaleOutActionData(d)

	if value, ok := d.GetOk("policy.0.scale_out_threshold"); ok {
		value := float32(value.(int))
		groupPolicy.ScaleOutThreshold = &value
	}

	if value, ok := d.GetOk("policy.0.unit"); ok {
		value := autoscaling.QueryUnit(value.(string))
		groupPolicy.Unit = &value
	}

	return &groupPolicy
}

func GetScaleInActionData(d *schema.ResourceData) *autoscaling.GroupPolicyScaleInAction {
	scaleInAction := autoscaling.GroupPolicyScaleInAction{}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.amount"); ok {
		value := float32(value.(int))
		scaleInAction.Amount = &value
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.amount_type"); ok {
		value := autoscaling.ActionAmount(value.(string))
		scaleInAction.AmountType = &value
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.termination_policy_type"); ok {
		value := autoscaling.TerminationPolicyType(value.(string))
		scaleInAction.TerminationPolicy = &value
	} else {
		scaleInAction.TerminationPolicy = nil
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.cooldown_period"); ok {
		value := value.(string)
		scaleInAction.CooldownPeriod = &value
	} else {
		scaleInAction.CooldownPeriod = nil
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.delete_volumes"); ok {
		value := value.(bool)
		scaleInAction.DeleteVolumes = &value
	} else {
		scaleInAction.DeleteVolumes = nil
	}

	return &scaleInAction
}

func GetScaleOutActionData(d *schema.ResourceData) *autoscaling.GroupPolicyScaleOutAction {
	scaleOutAction := autoscaling.GroupPolicyScaleOutAction{}

	if value, ok := d.GetOk("policy.0.scale_out_action.0.amount"); ok {
		value := float32(value.(int))
		scaleOutAction.Amount = &value
	}

	if value, ok := d.GetOk("policy.0.scale_out_action.0.amount_type"); ok {
		value := autoscaling.ActionAmount(value.(string))
		scaleOutAction.AmountType = &value
	}

	if value, ok := d.GetOk("policy.0.scale_out_action.0.cooldown_period"); ok {
		value := value.(string)
		scaleOutAction.CooldownPeriod = &value
	} else {
		scaleOutAction.CooldownPeriod = nil
	}

	return &scaleOutAction
}

func GetReplicaConfigurationPostData(d *schema.ResourceData) (*autoscaling.ReplicaPropertiesPost, error) {
	replica := autoscaling.ReplicaPropertiesPost{}

	if value, ok := d.GetOk("replica_configuration.0.availability_zone"); ok {
		value := autoscaling.AvailabilityZone(value.(string))
		replica.AvailabilityZone = &value
	}

	if value, ok := d.GetOk("replica_configuration.0.cores"); ok {
		value := int32(value.(int))
		replica.Cores = &value
	}

	if value, ok := d.GetOk("replica_configuration.0.cpu_family"); ok {
		value := autoscaling.CpuFamily(value.(string))
		replica.CpuFamily = &value
	} else {
		replica.CpuFamily = nil
	}

	replica.Nics = GetNicsData(d)
	if replica.Nics != nil && *replica.Nics == nil {
		*replica.Nics = make([]autoscaling.ReplicaNic, 0)
	}
	if value, ok := d.GetOk("replica_configuration.0.ram"); ok {
		value := int32(value.(int))
		replica.Ram = &value
	}

	if volumes, err := GetVolumesData(d); err != nil {
		return nil, err
	} else {
		replica.Volumes = volumes
	}

	return &replica, nil

}

func GetNicsData(d *schema.ResourceData) *[]autoscaling.ReplicaNic {
	var nics []autoscaling.ReplicaNic

	if nicsValue, ok := d.GetOk("replica_configuration.0.nic"); ok {
		nicsValue := nicsValue.(*schema.Set)
		if nicsValue != nil {
			for _, val := range nicsValue.List() {
				mmap := val.(map[string]any)
				var nicEntry = autoscaling.NewReplicaNic(int32(mmap["lan"].(int)), mmap["name"].(string))
				nicEntry.Dhcp = new(bool)
				*nicEntry.Dhcp = mmap["dhcp"].(bool)
				nics = append(nics, *nicEntry)
			}
		}

	}

	return &nics
}

func GetVolumesData(d *schema.ResourceData) (*[]autoscaling.ReplicaVolumePost, error) {
	var volumes []autoscaling.ReplicaVolumePost

	if volumesValue, ok := d.GetOk("replica_configuration.0.volumes"); ok {
		volumesValue := volumesValue.([]any)
		if volumesValue != nil {
			for index := range volumesValue {
				var volumeEntry autoscaling.ReplicaVolumePost

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image", index)); ok {
					value := value.(string)
					volumeEntry.Image = &value
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image_alias", index)); ok {
					value := value.(string)
					volumeEntry.ImageAlias = &value
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.name", index)); ok {
					value := value.(string)
					volumeEntry.Name = &value
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.size", index)); ok {
					value := int32(value.(int))
					volumeEntry.Size = &value
				}

				var publicKeys []string

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.ssh_keys", index)); ok {
					sshKeys := value.([]any)
					if len(sshKeys) != 0 {
						for _, keyOrPath := range sshKeys {
							//log.Printf("[DEBUG] Reading file %s", keyOrPath)
							publicKey, err := utils.ReadPublicKey(keyOrPath.(string))
							if err != nil {
								return nil, fmt.Errorf("error reading sshkey (%s) (%w)", keyOrPath, err)
							}
							publicKeys = append(publicKeys, publicKey)
						}
					}
				}

				volumeEntry.SshKeys = &publicKeys

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.type", index)); ok {
					value := autoscaling.VolumeHwType(value.(string))
					volumeEntry.Type = &value
				}

				//if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.user_data", index)); ok {
				//	value := value.(string)
				//	volumeEntry.UserData = &value
				//} else {
				//	volumeEntry.UserData = nil
				//}
				if userData, ok := d.GetOk("replica_configuration.0.volumes.%d.user_data"); ok {
					if *volumeEntry.Image == "" && *volumeEntry.ImageAlias == "" {
						return nil, fmt.Errorf("it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property ")
					} else {
						userData := userData.(string)
						volumeEntry.UserData = &userData
					}
				}
				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image_password", index)); ok {
					value := value.(string)
					volumeEntry.ImagePassword = &value
				} else {
					volumeEntry.ImagePassword = nil
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.boot_order", index)); ok {
					value := value.(string)
					volumeEntry.BootOrder = &value
				} else {
					volumeEntry.BootOrder = nil
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.bus", index)); ok {
					value := autoscaling.BusType(value.(string))
					volumeEntry.Bus = &value
				} else {
					volumeEntry.Bus = nil
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.backup_unit_id", index)); ok {
					value := value.(string)
					volumeEntry.BackupunitId = &value
				} else {
					volumeEntry.BackupunitId = nil
				}

				volumes = append(volumes, volumeEntry)
			}

		}

	}

	return &volumes, nil
}

func SetAutoscalingGroupData(d *schema.ResourceData, groupProperties *autoscaling.GroupProperties) error {

	resourceName := "autoscaling groupProperties"
	if groupProperties != nil {
		if groupProperties.MaxReplicaCount != nil {
			if err := d.Set("max_replica_count", *groupProperties.MaxReplicaCount); err != nil {
				return generateSetError(resourceName, "max_replica_count", err)
			}
		}

		if groupProperties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *groupProperties.MinReplicaCount); err != nil {
				return generateSetError(resourceName, "min_replica_count", err)
			}
		}

		//if groupProperties.TargetReplicaCount != nil {
		//	if err := d.Set("target_replica_count", *groupProperties.TargetReplicaCount); err != nil {
		//		return generateSetError(resourceName, "target_replica_count", err)
		//	}
		//}

		if groupProperties.Name != nil {
			if err := d.Set("name", *groupProperties.Name); err != nil {
				return generateSetError(resourceName, "name", err)
			}
		}

		if groupProperties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *groupProperties.MinReplicaCount); err != nil {
				return generateSetError(resourceName, "min_replica_count", err)
			}
		}

		if groupProperties.Policy != nil {
			var policies []any
			policy := setPolicyProperties(*groupProperties.Policy)
			policies = append(policies, policy)
			if err := d.Set("policy", policies); err != nil {
				return generateSetError(resourceName, "policy", err)
			}
		}

		if groupProperties.ReplicaConfiguration != nil {
			var replicaConfigurations []any
			replicaConfiguration := setReplicaConfiguration(d, *groupProperties.ReplicaConfiguration)
			replicaConfigurations = append(replicaConfigurations, replicaConfiguration)
			if err := d.Set("replica_configuration", replicaConfigurations); err != nil {
				return generateSetError(resourceName, "replica_configuration", err)
			}
		}

		if groupProperties.Datacenter != nil {
			if err := d.Set("datacenter_id", *groupProperties.Datacenter.Id); err != nil {
				return generateSetError(resourceName, "datacenter_id", err)
			}
		}

		if groupProperties.Location != nil {
			if err := d.Set("location", *groupProperties.Location); err != nil {
				return generateSetError(resourceName, "location", err)
			}
		}

	}
	return nil
}

func setPolicyProperties(groupPolicy autoscaling.GroupPolicy) map[string]any {

	policy := map[string]any{}

	setPropWithNilCheck(policy, "metric", groupPolicy.Metric)
	setPropWithNilCheck(policy, "range", groupPolicy.Range)
	setPropWithNilCheck(policy, "scale_in_threshold", groupPolicy.ScaleInThreshold)
	setPropWithNilCheck(policy, "scale_out_threshold", groupPolicy.ScaleOutThreshold)
	setPropWithNilCheck(policy, "unit", groupPolicy.Unit)

	if groupPolicy.ScaleInAction != nil {
		var scaleInActions []any
		scaleInAction := setScaleInActionProperties(*groupPolicy.ScaleInAction)
		scaleInActions = append(scaleInActions, scaleInAction)
		policy["scale_in_action"] = scaleInActions
	}
	if groupPolicy.ScaleOutAction != nil {
		var scaleOutActions []any
		scaleOutAction := setScaleOutActionProperties(*groupPolicy.ScaleOutAction)
		scaleOutActions = append(scaleOutActions, scaleOutAction)
		policy["scale_out_action"] = scaleOutActions
	}

	return policy
}

func setScaleInActionProperties(scaleInAction autoscaling.GroupPolicyScaleInAction) map[string]any {

	scaleIn := map[string]any{}

	setPropWithNilCheck(scaleIn, "amount", scaleInAction.Amount)
	setPropWithNilCheck(scaleIn, "amount_type", scaleInAction.AmountType)
	setPropWithNilCheck(scaleIn, "termination_policy_type", scaleInAction.TerminationPolicy)
	setPropWithNilCheck(scaleIn, "cooldown_period", scaleInAction.CooldownPeriod)
	setPropWithNilCheck(scaleIn, "delete_volumes", scaleInAction.DeleteVolumes)

	return scaleIn
}

func setScaleOutActionProperties(scaleOutAction autoscaling.GroupPolicyScaleOutAction) map[string]any {

	scaleOut := map[string]any{}

	setPropWithNilCheck(scaleOut, "amount", scaleOutAction.Amount)
	setPropWithNilCheck(scaleOut, "amount_type", scaleOutAction.AmountType)
	setPropWithNilCheck(scaleOut, "cooldown_period", scaleOutAction.CooldownPeriod)

	return scaleOut
}

func setReplicaConfiguration(d *schema.ResourceData, replicaConfiguration autoscaling.ReplicaPropertiesPost) map[string]any {

	replica := map[string]any{}

	setPropWithNilCheck(replica, "availability_zone", replicaConfiguration.AvailabilityZone)
	setPropWithNilCheck(replica, "cores", replicaConfiguration.Cores)
	setPropWithNilCheck(replica, "cpu_family", replicaConfiguration.CpuFamily)
	setPropWithNilCheck(replica, "ram", replicaConfiguration.Ram)

	if replicaConfiguration.Nics != nil {
		var nics []any
		for _, nic := range *replicaConfiguration.Nics {
			nicEntry := setNicProperties(nic)
			nics = append(nics, nicEntry)
		}
		replica["nic"] = nics
	}

	if replicaConfiguration.Volumes != nil {
		var volumes []any
		for idx, volume := range *replicaConfiguration.Volumes {
			volumeEntry := setVolumeProperties(d, idx, volume)
			volumes = append(volumes, volumeEntry)
		}
		replica["volume"] = volumes
	}

	return replica
}

func setNicProperties(replicaNic autoscaling.ReplicaNic) map[string]any {
	nic := map[string]any{}

	setPropWithNilCheck(nic, "lan", replicaNic.Lan)
	setPropWithNilCheck(nic, "name", replicaNic.Name)
	setPropWithNilCheck(nic, "dhcp", replicaNic.Dhcp)

	return nic
}

func setVolumeProperties(d *schema.ResourceData, index int, replicaVolume autoscaling.ReplicaVolumePost) map[string]any {
	volume := map[string]any{}

	setPropWithNilCheck(volume, "image", replicaVolume.Image)
	setPropWithNilCheck(volume, "image_alias", replicaVolume.ImageAlias)
	setPropWithNilCheck(volume, "name", replicaVolume.Name)
	setPropWithNilCheck(volume, "size", replicaVolume.Size)
	//setPropWithNilCheck(volume, "ssh_keys", replicaVolume.SshKeys)
	setPropWithNilCheck(volume, "type", replicaVolume.Type)
	setPropWithNilCheck(volume, "user_data", replicaVolume.UserData)
	//setPropWithNilCheck(volume, "image_password", replicaVolume.ImagePassword)
	setPropWithNilCheck(volume, "boot_order", replicaVolume.BootOrder)
	setPropWithNilCheck(volume, "bus", replicaVolume.Bus)
	//we need to take these from schema as they are not returned by API
	if password, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image_password", index)); ok {
		volume["image_password"] = password
	}
	if keys, ok := d.GetOk("replica_configuration.0.volumes.0.ssh_keys"); ok {
		volume["ssh_keys"] = keys
	}
	return volume
}
