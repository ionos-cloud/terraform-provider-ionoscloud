package cloudapinic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

func (fs *Service) List(ctx context.Context, datacenterID, serverID string, depth int32) ([]ionoscloud.Nic, error) {
	emptyNicList := make([]ionoscloud.Nic, 0)
	nics, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsGet(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return emptyNicList, err
	}
	if nics.Items == nil {
		log.Printf("[DEBUG] empty nic list for datacenter %s, server %s", datacenterID, serverID)
		return emptyNicList, nil
	}
	return *nics.Items, nil
}

func (fs *Service) Get(ctx context.Context, datacenterId, serverId, ID string, depth int32) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
	nic, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, datacenterId, serverId, ID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[DEBUG] no nic found for datacenter %s, server %s", datacenterId, serverId)
		}
		return nil, apiResponse, err
	}
	return &nic, apiResponse, nil
}

func (fs *Service) Delete(ctx context.Context, datacenterID, serverID, ID string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsDelete(ctx, datacenterID, serverID, ID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	if errState := cloudapi.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutDelete); errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for nic state change on delete dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return apiResponse, nil
}

// Create - creates the resource in the backend and waits until it is in an `AVAILABLE` state
func (fs *Service) Create(ctx context.Context, datacenterID, serverID string, nic ionoscloud.Nic) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
	val, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsPost(ctx, datacenterID, serverID).Nic(nic).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating nic for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			fs.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for nic state change on create dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &val, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterID, serverID, ID string, nicProperties ionoscloud.NicProperties) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
	updatedNic, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, datacenterID, serverID, ID).Nic(nicProperties).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating nic for dcId: %s, server_id: %s, id %s, Response: (%w)", datacenterID, serverID, ID, err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutUpdate); errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for nic state change on update dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return &updatedNic, apiResponse, nil
}

// DecodeTo - receives old and new values as slice of interfaces from schema, decodes and returns nic properties
func DecodeTo(oldValues, newValues []interface{}) ([]ionoscloud.Nic, []ionoscloud.Nic, error) {
	oldNicProps := make([]ionoscloud.Nic, len(oldValues))
	newNicProps := make([]ionoscloud.Nic, len(newValues))
	err := utils.DecodeInterfaceToStruct(newValues, newNicProps)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %+v to new values of nic rules %w", newValues, err)
	}
	err = utils.DecodeInterfaceToStruct(oldValues, oldNicProps)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %+v to values of nic rules %w", oldValues, err)
	}
	return oldNicProps, newNicProps, nil
}

func SetNetworkProperties(nic ionoscloud.Nic) map[string]interface{} {
	network := map[string]interface{}{}
	if nic.Properties != nil {
		utils.SetPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
		utils.SetPropWithNilCheck(network, "dhcpv6", nic.Properties.Dhcpv6)
		utils.SetPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)
		utils.SetPropWithNilCheck(network, "firewall_type", nic.Properties.FirewallType)
		utils.SetPropWithNilCheck(network, "lan", nic.Properties.Lan)
		utils.SetPropWithNilCheck(network, "name", nic.Properties.Name)
		utils.SetPropWithNilCheck(network, "ips", nic.Properties.Ips)
		utils.SetPropWithNilCheck(network, "ipv6_ips", nic.Properties.Ipv6Ips)
		utils.SetPropWithNilCheck(network, "ipv6_cidr_block", nic.Properties.Ipv6CidrBlock)
		utils.SetPropWithNilCheck(network, "mac", nic.Properties.Mac)
		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}
	}
	nsgIDs := make([]string, 0)
	if nic.Entities != nil && nic.Entities.Securitygroups != nil && nic.Entities.Securitygroups.Items != nil {
		for _, group := range *nic.Entities.Securitygroups.Items {
			if group.Id != nil {
				id := *group.Id
				nsgIDs = append(nsgIDs, id)
			}
		}
	}
	utils.SetPropWithNilCheck(network, "security_groups_ids", nsgIDs)

	return network
}

func GetNicFromSchema(d *schema.ResourceData, path string) (ionoscloud.Nic, error) {
	nic := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{},
	}

	lanInt := int32(d.Get(path + "lan").(int))
	nic.Properties.Lan = &lanInt

	if v, ok := d.GetOk(path + "name"); ok {
		vStr := v.(string)
		nic.Properties.Name = &vStr
	}

	dhcp := d.Get(path + "dhcp").(bool)
	if dhcpv6, ok := d.GetOkExists(path + "dhcpv6"); ok {
		dhcpv6 := dhcpv6.(bool)
		nic.Properties.Dhcpv6 = &dhcpv6
	} else {
		nic.Properties.SetDhcpv6Nil()
	}

	fwActive := d.Get(path + "firewall_active").(bool)
	nic.Properties.Dhcp = &dhcp
	nic.Properties.FirewallActive = &fwActive

	if _, ok := d.GetOk(path + "firewall_type"); ok {
		raw := d.Get(path + "firewall_type").(string)
		nic.Properties.FirewallType = &raw
	}

	if v, ok := d.GetOk(path + "ips"); ok {
		raw := v.([]interface{})
		if raw != nil && len(raw) > 0 {
			ips := make([]string, 0)
			for _, rawIp := range raw {
				if rawIp != nil {
					ip := rawIp.(string)
					ips = append(ips, ip)
				}
			}
			if ips != nil && len(ips) > 0 {
				nic.Properties.Ips = &ips
			}
		}
	}

	if v, ok := d.GetOk(path + "ipv6_ips"); ok {
		raw := v.([]interface{})
		ipv6Ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ipv6Ips)
		if err != nil {
			return nic, err
		}
		if len(ipv6Ips) > 0 {
			nic.Properties.Ipv6Ips = &ipv6Ips
		}
	}

	if v, ok := d.GetOk(path + "ipv6_cidr_block"); ok {
		ipv6Block := v.(string)
		nic.Properties.Ipv6CidrBlock = &ipv6Block
	}
	if flowLogs, ok := d.GetOk("flowlog"); ok {
		nic.Entities = &ionoscloud.NicEntities{
			Flowlogs: &ionoscloud.FlowLogs{
				Items: &[]ionoscloud.FlowLog{},
			},
		}
		if flowLogList, ok := flowLogs.([]any); ok {
			for _, flowLogData := range flowLogList {
				if flowLog, ok := flowLogData.(map[string]any); ok {
					*nic.Entities.Flowlogs.Items = append(*nic.Entities.Flowlogs.Items, cloudapiflowlog.GetFlowlogFromMap(flowLog))
				}
			}
		}
	}
	return nic, nil
}

func NicSetData(d *schema.ResourceData, nic *ionoscloud.Nic) error {
	if nic == nil {
		return fmt.Errorf("nic is empty")
	}

	if nic.Id != nil {
		d.SetId(*nic.Id)
	}

	if nic.Properties != nil {
		log.Printf("[INFO] LAN ON NIC: %d", nic.Properties.Lan)
		if nic.Properties.Dhcp != nil {
			if err := d.Set("dhcp", *nic.Properties.Dhcp); err != nil {
				return fmt.Errorf("error setting dhcp %w", err)
			}
		}

		if nic.Properties.Dhcpv6 != nil {
			if err := d.Set("dhcpv6", *nic.Properties.Dhcpv6); err != nil {
				return fmt.Errorf("error setting dhcpv6 %w", err)
			}
		}
		if nic.Properties.Lan != nil {
			if err := d.Set("lan", *nic.Properties.Lan); err != nil {
				return fmt.Errorf("error setting lan %w", err)
			}
		}
		if nic.Properties.Name != nil {
			if err := d.Set("name", *nic.Properties.Name); err != nil {
				return fmt.Errorf("error setting name %w", err)
			}
		}
		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			if err := d.Set("ips", *nic.Properties.Ips); err != nil {
				return fmt.Errorf("error setting ips %w", err)
			}
		}
		// should not be checked for len, we want to set the empty slice anyway as the field is computed, and it will not be set by backend
		// if ipv6_cidr_block is not set on the lan
		if nic.Properties.Ipv6Ips != nil {
			if err := d.Set("ipv6_ips", *nic.Properties.Ipv6Ips); err != nil {
				return fmt.Errorf("error setting ipv6_ips %w", err)
			}
		}
		if nic.Properties.Ipv6CidrBlock != nil {
			if err := d.Set("ipv6_cidr_block", *nic.Properties.Ipv6CidrBlock); err != nil {
				return fmt.Errorf("error setting ipv6_cidr_block %w", err)
			}
		}
		if nic.Properties.FirewallActive != nil {
			if err := d.Set("firewall_active", *nic.Properties.FirewallActive); err != nil {
				return fmt.Errorf("error setting firewall_active %w", err)
			}
		}
		if nic.Properties.FirewallType != nil {
			if err := d.Set("firewall_type", *nic.Properties.FirewallType); err != nil {
				return fmt.Errorf("error setting firewall_type %w", err)
			}
		}
		if nic.Properties.Mac != nil {
			if err := d.Set("mac", *nic.Properties.Mac); err != nil {
				return fmt.Errorf("error setting mac %w", err)
			}
		}
		if nic.Properties.DeviceNumber != nil {
			if err := d.Set("device_number", *nic.Properties.DeviceNumber); err != nil {
				return fmt.Errorf("error setting device_number %w", err)
			}
		}
		if nic.Properties.PciSlot != nil {
			if err := d.Set("pci_slot", *nic.Properties.PciSlot); err != nil {
				return fmt.Errorf("error setting pci_slot %w", err)
			}
		}
		nsgIDs := make([]string, 0)
		if nic.Entities != nil && nic.Entities.Securitygroups != nil && nic.Entities.Securitygroups.Items != nil {
			for _, group := range *nic.Entities.Securitygroups.Items {
				if group.Id != nil {
					id := *group.Id
					nsgIDs = append(nsgIDs, id)
				}
			}
		}
		if err := d.Set("security_groups_ids", nsgIDs); err != nil {
			return fmt.Errorf("error setting security_groups_ids %w", err)
		}
	}

	if nic.Entities != nil && nic.Entities.Flowlogs != nil && nic.Entities.Flowlogs.Items != nil && len(*nic.Entities.Flowlogs.Items) > 0 {
		var flowlogs []map[string]any
		for _, flowLog := range *nic.Entities.Flowlogs.Items {
			result := map[string]any{}
			result, err := utils.DecodeStructToMap(flowLog.Properties)
			if err != nil {
				return err
			}
			result["id"] = *flowLog.Id
			flowlogs = append(flowlogs, result)
		}
		if err := d.Set("flowlog", flowlogs); err != nil {
			return fmt.Errorf("error setting flowlog %w", err)
		}
	}

	return nil
}
