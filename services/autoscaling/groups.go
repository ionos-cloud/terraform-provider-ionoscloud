package autoscaling

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"log"
)

type GroupService interface {
	GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error)
	ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error)
	CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error)
	UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupUpdate) (autoscaling.Group, *autoscaling.APIResponse, error)
	DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error)
}

func (c *Client) GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.GroupsApi.AutoscalingGroupsFindById(ctx, groupId).Execute()
	if apiResponse != nil {
		return group, apiResponse, err

	}
	return group, nil, err
}

func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error) {
	groups, apiResponse, err := c.GroupsApi.AutoscalingGroupsGet(ctx).Execute()
	if apiResponse != nil {
		return groups, apiResponse, err
	}
	return groups, nil, err
}

func (c *Client) CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.GroupsApi.AutoscalingGroupsPost(ctx).Group(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupUpdate) (autoscaling.Group, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.GroupsApi.AutoscalingGroupsPut(ctx, groupId).GroupUpdate(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error) {
	apiResponse, err := c.GroupsApi.AutoscalingGroupsDelete(ctx, groupId).Execute()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func GetAutoscalingGroupDataCreate(d *schema.ResourceData) (*autoscaling.Group, error) {

	group := autoscaling.Group{
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

	if value, ok := d.GetOk("target_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.TargetReplicaCount = &value
	}

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

	group.Properties.Datacenter = GetDatacenterData(d)

	return &group, nil
}

func GetAutoscalingGroupDataUpdate(d *schema.ResourceData) (*autoscaling.GroupUpdate, error) {

	group := autoscaling.GroupUpdate{
		Properties: &autoscaling.GroupUpdatableProperties{},
	}

	if value, ok := d.GetOk("max_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MaxReplicaCount = &value
	}

	if value, ok := d.GetOk("min_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.MinReplicaCount = &value
	}

	if value, ok := d.GetOk("target_replica_count"); ok {
		value := int64(value.(int))
		group.Properties.TargetReplicaCount = &value
	} else {
		group.Properties.TargetReplicaCount = nil
	}

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

	if d.HasChange("datacenter") {
		return nil, fmt.Errorf("datacenter property is immutable and can pe used only in create requests")
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

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.nics.%d.dhcp", index)); ok {
					value := value.(bool)
					nicEntry.Dhcp = &value
				} else {
					nicEntry.Dhcp = nil
				}

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
					sshKeyPaths := value.([]string)
					if len(sshKeyPaths) != 0 {
						for _, path := range sshKeyPaths {
							log.Printf("[DEBUG] Reading file %s", path)
							publicKey, err := readPublicKey(path)
							if err != nil {
								return nil, fmt.Errorf("error fetching sshkey from file (%s) (%s)", path, err.Error())
							}
							publicKeys = append(publicKeys, publicKey)
						}
					}
				}

				if value, ok := d.GetOk(fmt.Sprintf("replica_configuration.0.volumes.%d.ssh_keys", index)); ok {
					sshKeys := value.([]string)
					if len(sshKeys) != 0 {
						for _, key := range sshKeys {
							publicKeys = append(publicKeys, key)
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

func GetDatacenterData(d *schema.ResourceData) *autoscaling.Resource {
	datacenter := autoscaling.Resource{}

	if value, ok := d.GetOk("datacenter.0.href"); ok {
		value := value.(string)
		datacenter.Href = &value
	}

	if value, ok := d.GetOk("datacenter.0.type"); ok {
		value := value.(string)
		datacenter.Type = &value
	}

	if value, ok := d.GetOk("datacenter.0.id"); ok {
		value := value.(string)
		datacenter.Id = &value
	}

	return &datacenter
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

		if group.Properties.TargetReplicaCount != nil {
			if err := d.Set("target_replica_count", *group.Properties.TargetReplicaCount); err != nil {
				return generateSetError(resourceName, "target_replica_count", err)
			}
		}

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
			policy := setPolicyProperties(*group.Properties.Policy)
			if err := d.Set("policy", policy); err != nil {
				return generateSetError(resourceName, "policy", err)
			}
		}

		if group.Properties.ReplicaConfiguration != nil {
			replicaConfiguration := setReplicaConfiguration(*group.Properties.ReplicaConfiguration)
			if err := d.Set("replica_configuration", replicaConfiguration); err != nil {
				return generateSetError(resourceName, "replica_configuration", err)
			}
		}

		if group.Properties.Datacenter != nil {
			datacenter := setDatacenterProperties(*group.Properties.Datacenter)
			if err := d.Set("datacenter", datacenter); err != nil {
				return generateSetError(resourceName, "datacenter", err)
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
		scaleInAction := setScaleInActionProperties(*groupPolicy.ScaleInAction)
		policy["scale_in_action"] = scaleInAction
	}
	if groupPolicy.ScaleOutAction != nil {
		scaleOutAction := setScaleOutActionProperties(*groupPolicy.ScaleOutAction)
		policy["scale_out_action"] = scaleOutAction
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

func setReplicaConfiguration(replicaConfiguration autoscaling.ReplicaPropertiesPost) map[string]interface{} {

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
			volumeEntry := setVolumeProperties(volume)
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

func setVolumeProperties(replicaVolume autoscaling.ReplicaVolumePost) map[string]interface{} {
	volume := map[string]interface{}{}

	setPropWithNilCheck(volume, "image", replicaVolume.Image)
	setPropWithNilCheck(volume, "name", replicaVolume.Name)
	setPropWithNilCheck(volume, "size", replicaVolume.Size)
	setPropWithNilCheck(volume, "ssh_keys", replicaVolume.SshKeys)
	setPropWithNilCheck(volume, "type", replicaVolume.Type)
	setPropWithNilCheck(volume, "user_data", replicaVolume.UserData)
	setPropWithNilCheck(volume, "image_password", replicaVolume.ImagePassword)

	return volume
}

func setDatacenterProperties(datacenterResource autoscaling.Resource) map[string]interface{} {
	datacenter := map[string]interface{}{}

	setPropWithNilCheck(datacenter, "href", datacenterResource.Href)
	setPropWithNilCheck(datacenter, "type", datacenterResource.Type)
	setPropWithNilCheck(datacenter, "id", datacenterResource.Id)

	return datacenter
}
