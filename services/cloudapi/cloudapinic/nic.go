package cloudapinic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

func (fs *Service) Get(ctx context.Context, datacenterID, serverID string, depth int32) ([]ionoscloud.Nic, error) {
	nics, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsGet(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	if apiResponse.HttpNotFound() || nics.Items == nil || len(*nics.Items) == 0 {
		log.Printf("[DEBUG] no nics found for datacenter %s, server %s", datacenterID, serverID)
		return nil, nil
	}
	return *nics.Items, nil
}

func (fs *Service) FindById(ctx context.Context, datacenterId, serverId, ID string, depth int32) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
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
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return apiResponse, nil
}

// Create - creates the resource in the backend and waits until it is in an `AVAILABLE` state
func (fs *Service) Create(ctx context.Context, datacenterID, serverID string, nic ionoscloud.Nic) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
	val, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsPost(ctx, datacenterID, serverID).Nic(nic).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating val rule for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			fs.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &val, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterID, serverID, ID string, nicProperties ionoscloud.NicProperties) (*ionoscloud.Nic, *ionoscloud.APIResponse, error) {
	updatedNic, apiResponse, err := fs.Client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, datacenterID, serverID, ID).Nic(nicProperties).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating nic for dcId: %s, server_id: %s, id %s, Response: (%w)", datacenterID, serverID, ID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
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

	return nic, nil
}
