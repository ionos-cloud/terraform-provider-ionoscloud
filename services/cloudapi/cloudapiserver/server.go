package cloudapiserver

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

func (fs *Service) FindById(ctx context.Context, datacenterID, serverID string, depth int32) (*ionoscloud.Server, error) {
	server, apiResponse, err := fs.Client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (fs *Service) Delete(ctx context.Context, datacenterID, serverID, ID string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := fs.Client.ServersApi.DatacentersServersDelete(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for server state change on delete dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return apiResponse, nil
}

func (fs *Service) Create(ctx context.Context, datacenterID string, nic ionoscloud.Nic) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	server, apiResponse, err := fs.Client.ServersApi.DatacentersServersPost(ctx, datacenterID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating server for dcId: %s, Response: (%w)", datacenterID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			fs.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for nic state change on create dcId: %s, Response: (%w)", datacenterID, errState)
	}
	return &server, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterID, serverID, ID string, serverProperties ionoscloud.ServerProperties) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	updatedServer, apiResponse, err := fs.Client.ServersApi.DatacentersServersPatch(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating server for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on update dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &updatedServer, apiResponse, nil
}

func (fs *Service) Start(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case "ENTERPRISE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "RUNNING"))

	case "CUBE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersResumePost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "RUNNING"))
	}

	return fmt.Errorf("unknown server type: %s", serverType)

}

func (fs *Service) Stop(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case "ENTERPRISE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "SHUTOFF"))

	case "CUBE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersSuspendPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "SUSPENDED"))
	}

	return fmt.Errorf("unknown server type: %s", serverType)

}

func (fs *Service) checkExpectedVmStateFn(ctx context.Context, dcId, expectedState string) func(ctx context.Context, d *schema.ResourceData) (bool, error) {

	return func(ctx context.Context, d *schema.ResourceData) (bool, error) {
		ionoscloudServer, _, err := fs.Client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()
		serverType := *ionoscloudServer.Properties.Type
		if err != nil {
			return false, err
		}
		if *ionoscloudServer.Properties.VmState != expectedState {
			log.Printf("[INFO] Server (type: %s) vmState not yet changed to %s: %s", serverType, expectedState, d.Id())
			return false, nil
		}
		return true, nil
	}
}

// // DecodeTo - receives old and new values as slice of interfaces from schema, decodes and returns nic properties
// func DecodeTo(oldValues, newValues []interface{}) ([]ionoscloud.Nic, []ionoscloud.Nic, error) {
// 	oldNicProps := make([]ionoscloud.Nic, len(oldValues))
// 	newNicProps := make([]ionoscloud.Nic, len(newValues))
// 	err := utils.DecodeInterfaceToStruct(newValues, newNicProps)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("could not decode from %+v to new values of nic rules %w", newValues, err)
// 	}
// 	err = utils.DecodeInterfaceToStruct(oldValues, oldNicProps)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("could not decode from %+v to values of nic rules %w", oldValues, err)
// 	}
// 	return oldNicProps, newNicProps, nil
// }

// func SetNetworkProperties(nic ionoscloud.Nic) map[string]interface{} {
// 	network := map[string]interface{}{}
// 	if nic.Properties != nil {
// 		utils.SetPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
// 		utils.SetPropWithNilCheck(network, "dhcpv6", nic.Properties.Dhcpv6)
// 		utils.SetPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)
// 		utils.SetPropWithNilCheck(network, "firewall_type", nic.Properties.FirewallType)
// 		utils.SetPropWithNilCheck(network, "lan", nic.Properties.Lan)
// 		utils.SetPropWithNilCheck(network, "name", nic.Properties.Name)
// 		utils.SetPropWithNilCheck(network, "ips", nic.Properties.Ips)
// 		utils.SetPropWithNilCheck(network, "ipv6_ips", nic.Properties.Ipv6Ips)
// 		utils.SetPropWithNilCheck(network, "ipv6_cidr_block", nic.Properties.Ipv6CidrBlock)
// 		utils.SetPropWithNilCheck(network, "mac", nic.Properties.Mac)
// 		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
// 			network["ips"] = *nic.Properties.Ips
// 		}
// 	}
// 	return network
// }

// func GetNicFromSchema(d *schema.ResourceData, path string) (ionoscloud.Nic, error) {
// 	nic := ionoscloud.Nic{
// 		Properties: &ionoscloud.NicProperties{},
// 	}

// 	lanInt := int32(d.Get(path + "lan").(int))
// 	nic.Properties.Lan = &lanInt

// 	if v, ok := d.GetOk(path + "name"); ok {
// 		vStr := v.(string)
// 		nic.Properties.Name = &vStr
// 	}

// 	dhcp := d.Get(path + "dhcp").(bool)
// 	if dhcpv6, ok := d.GetOkExists(path + "dhcpv6"); ok {
// 		dhcpv6 := dhcpv6.(bool)
// 		nic.Properties.Dhcpv6 = &dhcpv6
// 	} else {
// 		nic.Properties.SetDhcpv6Nil()
// 	}

// 	fwActive := d.Get(path + "firewall_active").(bool)
// 	nic.Properties.Dhcp = &dhcp
// 	nic.Properties.FirewallActive = &fwActive

// 	if _, ok := d.GetOk(path + "firewall_type"); ok {
// 		raw := d.Get(path + "firewall_type").(string)
// 		nic.Properties.FirewallType = &raw
// 	}

// 	if v, ok := d.GetOk(path + "ips"); ok {
// 		raw := v.([]interface{})
// 		if raw != nil && len(raw) > 0 {
// 			ips := make([]string, 0)
// 			for _, rawIp := range raw {
// 				if rawIp != nil {
// 					ip := rawIp.(string)
// 					ips = append(ips, ip)
// 				}
// 			}
// 			if ips != nil && len(ips) > 0 {
// 				nic.Properties.Ips = &ips
// 			}
// 		}
// 	}

// 	if v, ok := d.GetOk(path + "ipv6_ips"); ok {
// 		raw := v.([]interface{})
// 		ipv6Ips := make([]string, len(raw))
// 		err := utils.DecodeInterfaceToStruct(raw, ipv6Ips)
// 		if err != nil {
// 			return nic, err
// 		}
// 		if len(ipv6Ips) > 0 {
// 			nic.Properties.Ipv6Ips = &ipv6Ips
// 		}
// 	}

// 	if v, ok := d.GetOk(path + "ipv6_cidr_block"); ok {
// 		ipv6Block := v.(string)
// 		nic.Properties.Ipv6CidrBlock = &ipv6Block
// 	}

// 	return nic, nil
// }

// func NicSetData(d *schema.ResourceData, nic *ionoscloud.Nic) error {
// 	if nic == nil {
// 		return fmt.Errorf("nic is empty")
// 	}

// 	if nic.Id != nil {
// 		d.SetId(*nic.Id)
// 	}

// 	if nic.Properties != nil {
// 		log.Printf("[INFO] LAN ON NIC: %d", nic.Properties.Lan)
// 		if nic.Properties.Dhcp != nil {
// 			if err := d.Set("dhcp", *nic.Properties.Dhcp); err != nil {
// 				return fmt.Errorf("error setting dhcp %w", err)
// 			}
// 		}

// 		if nic.Properties.Dhcpv6 != nil {
// 			if err := d.Set("dhcpv6", *nic.Properties.Dhcpv6); err != nil {
// 				return fmt.Errorf("error setting dhcpv6 %w", err)
// 			}
// 		}
// 		if nic.Properties.Lan != nil {
// 			if err := d.Set("lan", *nic.Properties.Lan); err != nil {
// 				return fmt.Errorf("error setting lan %w", err)
// 			}
// 		}
// 		if nic.Properties.Name != nil {
// 			if err := d.Set("name", *nic.Properties.Name); err != nil {
// 				return fmt.Errorf("error setting name %w", err)
// 			}
// 		}
// 		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
// 			if err := d.Set("ips", *nic.Properties.Ips); err != nil {
// 				return fmt.Errorf("error setting ips %w", err)
// 			}
// 		}
// 		//should not be checked for len, we want to set the empty slice anyway as the field is computed, and it will not be set by backend
// 		// if ipv6_cidr_block is not set on the lan
// 		if nic.Properties.Ipv6Ips != nil {
// 			if err := d.Set("ipv6_ips", *nic.Properties.Ipv6Ips); err != nil {
// 				return fmt.Errorf("error setting ipv6_ips %w", err)
// 			}
// 		}
// 		if nic.Properties.Ipv6CidrBlock != nil {
// 			if err := d.Set("ipv6_cidr_block", *nic.Properties.Ipv6CidrBlock); err != nil {
// 				return fmt.Errorf("error setting ipv6_cidr_block %w", err)
// 			}
// 		}
// 		if nic.Properties.FirewallActive != nil {
// 			if err := d.Set("firewall_active", *nic.Properties.FirewallActive); err != nil {
// 				return fmt.Errorf("error setting firewall_active %w", err)
// 			}
// 		}
// 		if nic.Properties.FirewallType != nil {
// 			if err := d.Set("firewall_type", *nic.Properties.FirewallType); err != nil {
// 				return fmt.Errorf("error setting firewall_type %w", err)
// 			}
// 		}
// 		if nic.Properties.Mac != nil {
// 			if err := d.Set("mac", *nic.Properties.Mac); err != nil {
// 				return fmt.Errorf("error setting mac %w", err)
// 			}
// 		}
// 		if nic.Properties.DeviceNumber != nil {
// 			if err := d.Set("device_number", *nic.Properties.DeviceNumber); err != nil {
// 				return fmt.Errorf("error setting device_number %w", err)
// 			}
// 		}
// 		if nic.Properties.PciSlot != nil {
// 			if err := d.Set("pci_slot", *nic.Properties.PciSlot); err != nil {
// 				return fmt.Errorf("error setting pci_slot %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }
