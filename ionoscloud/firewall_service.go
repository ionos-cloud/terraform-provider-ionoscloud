package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strconv"
)

type FirewallService struct {
	client     *ionoscloud.APIClient
	meta       interface{}
	schemaData *schema.ResourceData
}

func (fs *FirewallService) firewallsGet(ctx context.Context, datacenterId, serverId, nicId string, depth int32) ([]ionoscloud.FirewallRule, error) {
	firewallRules, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(depth).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.HttpNotFound() || firewallRules.Items == nil || len(*firewallRules.Items) == 0 {
		log.Printf("[DEBUG] no firewalls found for datacenter %s, server %s, nic %s", datacenterId, serverId, nicId)
		return nil, nil
	}
	return *firewallRules.Items, nil
}

func (fs *FirewallService) firewallFindById(ctx context.Context, datacenterId, serverId, nicId, firewallId string) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewallRule, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, firewallId).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return nil, apiResponse, err
	}
	if apiResponse.HttpNotFound() {
		log.Printf("[DEBUG] no firewall rule found for datacenter %s, server %s, nic %s", datacenterId, serverId, nicId)
		return nil, apiResponse, nil
	}
	return &firewallRule, apiResponse, nil
}

func (fs *FirewallService) deleteFirewallRule(ctx context.Context, datacenterId, serverId, nicId, firewallId string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDelete(ctx, datacenterId, serverId, nicId, firewallId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(fs.meta, fs.schemaData, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, ID: %s, Response: (%w)", datacenterId, serverId, nicId, firewallId, errState)
	}
	return apiResponse, nil
}

func (fs *FirewallService) createFirewallRule(ctx context.Context, datacenterId, serverId, nicId string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewall, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, datacenterId, serverId, nicId).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating firewall rule for dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(fs.meta, fs.schemaData, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, errState)
	}
	return &firewall, apiResponse, nil
}

func (fs *FirewallService) updateFirewallRule(ctx context.Context, datacenterId, serverId, nicId, id string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewall, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPut(ctx, datacenterId, serverId, nicId, id).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating firewall rule for dcId: %s, server_id: %s, nic_id: %s, id %s, Response: (%w)", datacenterId, serverId, nicId, id, err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(fs.meta, fs.schemaData, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, errState)
	}
	return &firewall, apiResponse, nil
}

func SetFirewallProperties(firewall ionoscloud.FirewallRule) map[string]interface{} {

	fw := map[string]interface{}{}
	if firewall.Properties != nil {
		utils.SetPropWithNilCheck(fw, "protocol", firewall.Properties.Protocol)
		utils.SetPropWithNilCheck(fw, "name", firewall.Properties.Name)
		utils.SetPropWithNilCheck(fw, "source_mac", firewall.Properties.SourceMac)
		utils.SetPropWithNilCheck(fw, "source_ip", firewall.Properties.SourceIp)
		utils.SetPropWithNilCheck(fw, "target_ip", firewall.Properties.TargetIp)
		utils.SetPropWithNilCheck(fw, "port_range_start", firewall.Properties.PortRangeStart)
		utils.SetPropWithNilCheck(fw, "port_range_end", firewall.Properties.PortRangeEnd)
		utils.SetPropWithNilCheck(fw, "type", firewall.Properties.Type)
		if firewall.Properties.IcmpType != nil {
			fw["icmp_type"] = strconv.Itoa(int(*firewall.Properties.IcmpType))
		}
		if firewall.Properties.IcmpCode != nil {
			fw["icmp_code"] = strconv.Itoa(int(*firewall.Properties.IcmpCode))
		}
	}
	return fw
}

func GetChangesInFirewallRules(oldValues, newValues []interface{}) ([]ionoscloud.FirewallruleProperties, []ionoscloud.FirewallruleProperties, error) {
	onlyOld := slice.Difference(oldValues, newValues)
	onlyNew := slice.Difference(newValues, oldValues)
	oldFwSlice := make([]ionoscloud.FirewallruleProperties, len(onlyOld))
	newFwSlice := make([]ionoscloud.FirewallruleProperties, len(onlyNew))
	err := utils.DecodeInterfaceToStruct(onlyNew, newFwSlice)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %s to new values of firewall rules %w", newValues, err)
	}
	err = utils.DecodeInterfaceToStruct(onlyOld, oldFwSlice)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %s to values of firewall rules %w", oldValues, err)
	}
	return oldFwSlice, newFwSlice, nil
}

// FwPropUnsetSetFieldIfNotSetInSchema will only set the in32 types if they actually exist in the schema
// mutates fwProp
func FwPropUnsetSetFieldIfNotSetInSchema(fwProp *ionoscloud.FirewallruleProperties, path string, d *schema.ResourceData) {
	if *fwProp.PortRangeStart == 0 {
		if _, isSet := d.GetOkExists(path + ".port_range_start"); isSet != true {
			fwProp.PortRangeStart = nil
		}
	}
	if *fwProp.PortRangeEnd == 0 {
		if _, isSet := d.GetOkExists(path + ".port_range_end"); isSet != true {
			fwProp.PortRangeEnd = nil
		}
	}
}
