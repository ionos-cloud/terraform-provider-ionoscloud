package autoscaling

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *Client) GetGroup(ctx context.Context, groupId string, depth float32) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsFindById(ctx, groupId).Depth(depth).Execute()
	apiResponse.LogInfo()
	return group, apiResponse, err
}

func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error) {
	groups, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsGet(ctx).Execute()
	apiResponse.LogInfo()
	return groups, apiResponse, err
}

func (c *Client) CreateGroup(ctx context.Context, group autoscaling.GroupPost) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPost(ctx).GroupPost(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPut(ctx, groupId).GroupPut(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error) {
	apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsDelete(ctx, groupId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
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
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.cooldown_period"); ok {
		value := value.(string)
		scaleInAction.CooldownPeriod = &value
	}

	if value, ok := d.GetOk("policy.0.scale_in_action.0.delete_volumes"); ok {
		value := value.(bool)
		scaleInAction.DeleteVolumes = &value
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
	}

	replica.Nics = GetNicsData(d)

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
	nics := make([]autoscaling.ReplicaNic, 0)

	if value, ok := d.GetOk("replica_configuration.0.nic"); ok {
		nicsSet := value.(*schema.Set)
		if nicsSet != nil {
			for _, val := range nicsSet.List() {
				nicsMap := val.(map[string]any)
				var nicEntry = autoscaling.NewReplicaNic(int32(nicsMap["lan"].(int)), nicsMap["name"].(string))
				nicEntry.Dhcp = new(bool)
				*nicEntry.Dhcp = nicsMap["dhcp"].(bool)
				nics = append(nics, *nicEntry)
			}
		}
	}

	return &nics
}

func GetVolumesData(d *schema.ResourceData) (*[]autoscaling.ReplicaVolumePost, error) {
	var volumes []autoscaling.ReplicaVolumePost

	if value, ok := d.GetOk("replica_configuration.0.volume"); ok {
		volumesSet := value.(*schema.Set)
		if volumesSet != nil {
			volumesValue := volumesSet.List()
			if volumesValue != nil {
				for index := range volumesValue {
					volumeMap := volumesValue[index].(map[string]any)
					var volumeEntry = *autoscaling.NewReplicaVolumePostWithDefaults()
					if val, ok := volumeMap["name"]; ok {
						volumeEntry.Name = new(string)
						*volumeEntry.Name = val.(string)
					}
					if val, ok := volumeMap["image_alias"]; ok {
						volumeEntry.ImageAlias = new(string)
						*volumeEntry.ImageAlias = val.(string)
					}
					if val, ok := volumeMap["image"]; ok {
						volumeEntry.Image = new(string)
						*volumeEntry.Image = val.(string)
					}
					if val, ok := volumeMap["size"]; ok {
						volumeEntry.Size = new(int32)
						*volumeEntry.Size = int32(val.(int))
					}
					if val, ok := volumeMap["type"]; ok {
						volumeEntry.Type = new(autoscaling.VolumeHwType)
						*volumeEntry.Type = autoscaling.VolumeHwType(val.(string))
					}
					if *volumeEntry.Image == "" && *volumeEntry.ImageAlias == "" {
						return nil, fmt.Errorf("it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property")
					}
					if val, ok := volumeMap["user_data"]; ok {
						volumeEntry.UserData = new(string)
						*volumeEntry.UserData = val.(string)
					}
					var publicKeys []string
					sshKeys := volumeMap["ssh_keys"].([]any)
					for _, keyOrPath := range sshKeys {
						publicKey, err := utils.ReadPublicKey(keyOrPath.(string))
						if err != nil {
							return nil, fmt.Errorf("error reading sshkey (%s) (%w)", keyOrPath, err)
						}
						publicKeys = append(publicKeys, publicKey)
					}

					volumeEntry.SshKeys = &publicKeys
					if val, ok := volumeMap["image_password"]; ok {
						volumeEntry.ImagePassword = new(string)
						*volumeEntry.ImagePassword = val.(string)
					}

					if val, ok := volumeMap["bus"]; ok {
						volumeEntry.Bus = new(autoscaling.BusType)
						*volumeEntry.Bus = autoscaling.BusType(val.(string))
					}
					if val, ok := volumeMap["backup_unit_id"]; ok {
						volumeEntry.BackupunitId = new(string)
						*volumeEntry.BackupunitId = val.(string)
					}
					if val, ok := volumeMap["boot_order"]; ok {
						volumeEntry.BootOrder = new(string)
						*volumeEntry.BootOrder = val.(string)
					}
					volumes = append(volumes, volumeEntry)
				}
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
				return utils.GenerateSetError(resourceName, "max_replica_count", err)
			}
		}

		if groupProperties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *groupProperties.MinReplicaCount); err != nil {
				return utils.GenerateSetError(resourceName, "min_replica_count", err)
			}
		}

		//if groupProperties.TargetReplicaCount != nil {
		//	if err := d.Set("target_replica_count", *groupProperties.TargetReplicaCount); err != nil {
		//		return utils.GenerateSetError(resourceName, "target_replica_count", err)
		//	}
		//}

		if groupProperties.Name != nil {
			if err := d.Set("name", *groupProperties.Name); err != nil {
				return utils.GenerateSetError(resourceName, "name", err)
			}
		}

		if groupProperties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *groupProperties.MinReplicaCount); err != nil {
				return utils.GenerateSetError(resourceName, "min_replica_count", err)
			}
		}

		if groupProperties.Policy != nil {
			var policies []any
			policy := setPolicyProperties(*groupProperties.Policy)
			policies = append(policies, policy)
			if err := d.Set("policy", policies); err != nil {
				return utils.GenerateSetError(resourceName, "policy", err)
			}
		}

		if groupProperties.ReplicaConfiguration != nil {
			var replicaConfigurations []any
			replicaConfiguration := setReplicaConfiguration(d, *groupProperties.ReplicaConfiguration)
			replicaConfigurations = append(replicaConfigurations, replicaConfiguration)
			if err := d.Set("replica_configuration", replicaConfigurations); err != nil {
				return utils.GenerateSetError(resourceName, "replica_configuration", err)
			}
		}

		if groupProperties.Datacenter != nil {
			if err := d.Set("datacenter_id", *groupProperties.Datacenter.Id); err != nil {
				return utils.GenerateSetError(resourceName, "datacenter_id", err)
			}
		}

		if groupProperties.Location != nil {
			if err := d.Set("location", *groupProperties.Location); err != nil {
				return utils.GenerateSetError(resourceName, "location", err)
			}
		}

	}
	return nil
}

func setPolicyProperties(groupPolicy autoscaling.GroupPolicy) map[string]any {

	policy := map[string]any{}

	utils.SetPropWithNilCheck(policy, "metric", groupPolicy.Metric)
	utils.SetPropWithNilCheck(policy, "range", groupPolicy.Range)
	utils.SetPropWithNilCheck(policy, "scale_in_threshold", groupPolicy.ScaleInThreshold)
	utils.SetPropWithNilCheck(policy, "scale_out_threshold", groupPolicy.ScaleOutThreshold)
	utils.SetPropWithNilCheck(policy, "unit", groupPolicy.Unit)

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

	utils.SetPropWithNilCheck(scaleIn, "amount", scaleInAction.Amount)
	utils.SetPropWithNilCheck(scaleIn, "amount_type", scaleInAction.AmountType)
	utils.SetPropWithNilCheck(scaleIn, "termination_policy_type", scaleInAction.TerminationPolicy)
	utils.SetPropWithNilCheck(scaleIn, "cooldown_period", scaleInAction.CooldownPeriod)
	utils.SetPropWithNilCheck(scaleIn, "delete_volumes", scaleInAction.DeleteVolumes)

	return scaleIn
}

func setScaleOutActionProperties(scaleOutAction autoscaling.GroupPolicyScaleOutAction) map[string]any {

	scaleOut := map[string]any{}

	utils.SetPropWithNilCheck(scaleOut, "amount", scaleOutAction.Amount)
	utils.SetPropWithNilCheck(scaleOut, "amount_type", scaleOutAction.AmountType)
	utils.SetPropWithNilCheck(scaleOut, "cooldown_period", scaleOutAction.CooldownPeriod)

	return scaleOut
}

func setReplicaConfiguration(d *schema.ResourceData, replicaConfiguration autoscaling.ReplicaPropertiesPost) map[string]any {

	replica := map[string]any{}

	utils.SetPropWithNilCheck(replica, "availability_zone", replicaConfiguration.AvailabilityZone)
	utils.SetPropWithNilCheck(replica, "cores", replicaConfiguration.Cores)
	utils.SetPropWithNilCheck(replica, "cpu_family", replicaConfiguration.CpuFamily)
	utils.SetPropWithNilCheck(replica, "ram", replicaConfiguration.Ram)

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

	utils.SetPropWithNilCheck(nic, "lan", replicaNic.Lan)
	utils.SetPropWithNilCheck(nic, "name", replicaNic.Name)
	utils.SetPropWithNilCheck(nic, "dhcp", replicaNic.Dhcp)

	return nic
}

func setVolumeProperties(d *schema.ResourceData, index int, replicaVolume autoscaling.ReplicaVolumePost) map[string]any {
	volume := map[string]any{}

	utils.SetPropWithNilCheck(volume, "image", replicaVolume.Image)
	utils.SetPropWithNilCheck(volume, "image_alias", replicaVolume.ImageAlias)
	utils.SetPropWithNilCheck(volume, "name", replicaVolume.Name)
	utils.SetPropWithNilCheck(volume, "size", replicaVolume.Size)
	utils.SetPropWithNilCheck(volume, "type", replicaVolume.Type)
	utils.SetPropWithNilCheck(volume, "boot_order", replicaVolume.BootOrder)
	utils.SetPropWithNilCheck(volume, "bus", replicaVolume.Bus)
	//we need to take these from schema as they are not returned by API
	volumeMap, ok := d.GetOk("replica_configuration.0.volume")
	if ok {
		volumeMap := (volumeMap).(*schema.Set).List()[index].(map[string]any)
		volume["image_password"] = volumeMap["image_password"]
		volume["ssh_keys"] = volumeMap["ssh_keys"]
		volume["user_data"] = volumeMap["user_data"]
	}

	return volume
}
