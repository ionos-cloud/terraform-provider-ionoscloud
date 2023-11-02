package autoscaling

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

type GroupService interface {
	GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error)
	ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error)
	CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error)
	UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *autoscaling.APIResponse, error)
	DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error)
}

func (c *Client) GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsFindById(ctx, groupId).Execute()
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

	if nicsValue, ok := d.GetOk("replica_configuration.0.nics"); ok {
		nicsValue := nicsValue.([]interface{})
		if nicsValue != nil {
			for index := range nicsValue {
				var nicEntry autoscaling.ReplicaNic

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.nics.%d.lan", index)); ok {
					value := int32(value.(int))
					nicEntry.Lan = &value
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.nics.%d.name", index)); ok {
					value := value.(string)
					nicEntry.Name = &value
				}

				value := d.Get(fmt.Sprintf("replica_configuration.0.nics.%d.dhcp", index)).(bool)
				nicEntry.Dhcp = &value

				nics = append(nics, nicEntry)
			}
		}

	}

	return &nics
}

func GetVolumesData(d *schema.ResourceData) (*[]autoscaling.ReplicaVolumePost, error) {
	var volumes []autoscaling.ReplicaVolumePost

	if volumesValue, ok := d.GetOk("replica_configuration.0.volumes"); ok {
		volumesValue := volumesValue.([]interface{})
		if volumesValue != nil {
			for index := range volumesValue {
				var volumeEntry autoscaling.ReplicaVolumePost

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image", index)); ok {
					value := value.(string)
					volumeEntry.Image = &value
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

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.ssh_key_paths", index)); ok {
					sshKeyPaths := value.([]interface{})
					if len(sshKeyPaths) != 0 {
						for _, path := range sshKeyPaths {
							log.Printf("[DEBUG] Reading file %s", path)
							publicKey, err := readPublicKey(path.(string))
							if err != nil {
								return nil, fmt.Errorf("error fetching sshkey from file (%s) (%s)", path, err.Error())
							}
							publicKeys = append(publicKeys, publicKey)
						}
					}
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.ssh_key_values", index)); ok {
					sshKeys := value.([]interface{})
					if len(sshKeys) != 0 {
						for _, key := range sshKeys {
							publicKeys = append(publicKeys, key.(string))
						}
					}
				}

				volumeEntry.SshKeys = &publicKeys

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.type", index)); ok {
					value := autoscaling.VolumeHwType(value.(string))
					volumeEntry.Type = &value
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.user_data", index)); ok {
					value := value.(string)
					volumeEntry.UserData = &value
				} else {
					volumeEntry.UserData = nil
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.image_password", index)); ok {
					value := value.(string)
					volumeEntry.ImagePassword = &value
				} else {
					volumeEntry.ImagePassword = nil
				}

				volumes = append(volumes, volumeEntry)
			}

		}

	}

	return &volumes, nil
}

func SetAutoscalingGroupData(d *schema.ResourceData, group autoscaling.Group) error {

	resourceName := "autoscaling group"

	if group.Id != nil {
		d.SetId(*group.Id)
	}

	if group.Properties != nil {
		if group.Properties.MaxReplicaCount != nil {
			if err := d.Set("max_replica_count", *group.Properties.MaxReplicaCount); err != nil {
				return generateSetError(resourceName, "max_replica_count", err)
			}
		}

		if group.Properties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *group.Properties.MinReplicaCount); err != nil {
				return generateSetError(resourceName, "min_replica_count", err)
			}
		}

		//if group.Properties.TargetReplicaCount != nil {
		//	if err := d.Set("target_replica_count", *group.Properties.TargetReplicaCount); err != nil {
		//		return generateSetError(resourceName, "target_replica_count", err)
		//	}
		//}

		if group.Properties.Name != nil {
			if err := d.Set("name", *group.Properties.Name); err != nil {
				return generateSetError(resourceName, "name", err)
			}
		}

		if group.Properties.MinReplicaCount != nil {
			if err := d.Set("min_replica_count", *group.Properties.MinReplicaCount); err != nil {
				return generateSetError(resourceName, "min_replica_count", err)
			}
		}

		if group.Properties.Policy != nil {
			var policies []interface{}
			policy := setPolicyProperties(*group.Properties.Policy)
			policies = append(policies, policy)
			if err := d.Set("policy", policies); err != nil {
				return generateSetError(resourceName, "policy", err)
			}
		}

		if group.Properties.ReplicaConfiguration != nil {
			var replicaConfigurations []interface{}
			replicaConfiguration := setReplicaConfiguration(d, *group.Properties.ReplicaConfiguration)
			replicaConfigurations = append(replicaConfigurations, replicaConfiguration)
			if err := d.Set("replica_configuration", replicaConfigurations); err != nil {
				return generateSetError(resourceName, "replica_configuration", err)
			}
		}

		if group.Properties.Datacenter != nil {
			if err := d.Set("datacenter_id", *group.Properties.Datacenter.Id); err != nil {
				return generateSetError(resourceName, "datacenter_id", err)
			}
		}

		if group.Properties.Location != nil {
			if err := d.Set("location", *group.Properties.Location); err != nil {
				return generateSetError(resourceName, "location", err)
			}
		}

	}
	return nil
}

func setPolicyProperties(groupPolicy autoscaling.GroupPolicy) map[string]interface{} {

	policy := map[string]interface{}{}

	setPropWithNilCheck(policy, "metric", groupPolicy.Metric)
	setPropWithNilCheck(policy, "range", groupPolicy.Range)
	setPropWithNilCheck(policy, "scale_in_threshold", groupPolicy.ScaleInThreshold)
	setPropWithNilCheck(policy, "scale_out_threshold", groupPolicy.ScaleOutThreshold)
	setPropWithNilCheck(policy, "unit", groupPolicy.Unit)

	if groupPolicy.ScaleInAction != nil {
		var scaleInActions []interface{}
		scaleInAction := setScaleInActionProperties(*groupPolicy.ScaleInAction)
		scaleInActions = append(scaleInActions, scaleInAction)
		policy["scale_in_action"] = scaleInActions
	}
	if groupPolicy.ScaleOutAction != nil {
		var scaleOutActions []interface{}
		scaleOutAction := setScaleOutActionProperties(*groupPolicy.ScaleOutAction)
		scaleOutActions = append(scaleOutActions, scaleOutAction)
		policy["scale_out_action"] = scaleOutActions
	}

	return policy
}

func setScaleInActionProperties(scaleInAction autoscaling.GroupPolicyScaleInAction) map[string]interface{} {

	scaleIn := map[string]interface{}{}

	setPropWithNilCheck(scaleIn, "amount", scaleInAction.Amount)
	setPropWithNilCheck(scaleIn, "amount_type", scaleInAction.AmountType)
	setPropWithNilCheck(scaleIn, "termination_policy_type", scaleInAction.TerminationPolicy)
	setPropWithNilCheck(scaleIn, "cooldown_period", scaleInAction.CooldownPeriod)

	return scaleIn
}

func setScaleOutActionProperties(scaleOutAction autoscaling.GroupPolicyScaleOutAction) map[string]interface{} {

	scaleOut := map[string]interface{}{}

	setPropWithNilCheck(scaleOut, "amount", scaleOutAction.Amount)
	setPropWithNilCheck(scaleOut, "amount_type", scaleOutAction.AmountType)
	setPropWithNilCheck(scaleOut, "cooldown_period", scaleOutAction.CooldownPeriod)

	return scaleOut
}

func setReplicaConfiguration(d *schema.ResourceData, replicaConfiguration autoscaling.ReplicaPropertiesPost) map[string]interface{} {

	replica := map[string]interface{}{}

	setPropWithNilCheck(replica, "availability_zone", replicaConfiguration.AvailabilityZone)
	setPropWithNilCheck(replica, "cores", replicaConfiguration.Cores)
	setPropWithNilCheck(replica, "cpu_family", replicaConfiguration.CpuFamily)
	setPropWithNilCheck(replica, "ram", replicaConfiguration.Ram)

	if replicaConfiguration.Nics != nil {
		var nics []interface{}
		for _, nic := range *replicaConfiguration.Nics {
			nicEntry := setNicProperties(nic)
			nics = append(nics, nicEntry)
		}
		replica["nics"] = nics
	}

	if replicaConfiguration.Volumes != nil {
		var volumes []interface{}
		for _, volume := range *replicaConfiguration.Volumes {
			volumeEntry := setVolumeProperties(d, volume)
			volumes = append(volumes, volumeEntry)
		}
		replica["volumes"] = volumes
	}

	return replica
}

func setNicProperties(replicaNic autoscaling.ReplicaNic) map[string]interface{} {
	nic := map[string]interface{}{}

	setPropWithNilCheck(nic, "lan", replicaNic.Lan)
	setPropWithNilCheck(nic, "name", replicaNic.Name)
	setPropWithNilCheck(nic, "dhcp", replicaNic.Dhcp)

	return nic
}

func setVolumeProperties(d *schema.ResourceData, replicaVolume autoscaling.ReplicaVolumePost) map[string]interface{} {
	volume := map[string]interface{}{}

	setPropWithNilCheck(volume, "image", replicaVolume.Image)
	setPropWithNilCheck(volume, "name", replicaVolume.Name)
	setPropWithNilCheck(volume, "size", replicaVolume.Size)
	setPropWithNilCheck(volume, "ssh_keys", replicaVolume.SshKeys)
	setPropWithNilCheck(volume, "type", replicaVolume.Type)
	setPropWithNilCheck(volume, "user_data", replicaVolume.UserData)

	if paths, ok := d.GetOk("replica_configuration.0.volumes.0.ssh_key_paths"); ok {
		volume["ssh_key_paths"] = paths
	}

	if paths, ok := d.GetOk("replica_configuration.0.volumes.0.ssh_key_values"); ok {
		volume["ssh_key_values"] = paths
	}

	if password, ok := d.GetOk("replica_configuration.0.volumes.0.image_password"); ok {
		volume["image_password"] = password
	}

	return volume
}
