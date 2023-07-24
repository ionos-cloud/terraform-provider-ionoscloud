package firewallSvc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"log"
	"reflect"
	"strconv"
)

type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

func (fs *Service) Get(ctx context.Context, datacenterId, serverId, nicId string, depth int32) ([]ionoscloud.FirewallRule, error) {
	firewallRules, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	if apiResponse.HttpNotFound() || firewallRules.Items == nil || len(*firewallRules.Items) == 0 {
		log.Printf("[DEBUG] no firewalls found for datacenter %s, server %s, nic %s", datacenterId, serverId, nicId)
		return nil, nil
	}
	return *firewallRules.Items, nil
}

func (fs *Service) FindById(ctx context.Context, datacenterId, serverId, nicId, firewallId string) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewallRule, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, firewallId).Depth(1).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, err
	}
	if apiResponse.HttpNotFound() {
		log.Printf("[DEBUG] no firewall rule found for datacenter %s, server %s, nic %s", datacenterId, serverId, nicId)
		return nil, apiResponse, nil
	}
	return &firewallRule, apiResponse, nil
}

func (fs *Service) Delete(ctx context.Context, datacenterId, serverId, nicId, firewallId string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDelete(ctx, datacenterId, serverId, nicId, firewallId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, ID: %s, Response: (%w)", datacenterId, serverId, nicId, firewallId, errState)
	}
	return apiResponse, nil
}

func (fs *Service) Create(ctx context.Context, datacenterId, serverId, nicId string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewall, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, datacenterId, serverId, nicId).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating firewall rule for dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, errState)
	}
	return &firewall, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterId, serverId, nicId, id string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewall, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPut(ctx, datacenterId, serverId, nicId, id).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating firewall rule for dcId: %s, server_id: %s, nic_id: %s, id %s, Response: (%w)", datacenterId, serverId, nicId, id, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for state change dcId: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterId, serverId, nicId, errState)
	}
	return &firewall, apiResponse, nil
}

func SetProperties(firewall ionoscloud.FirewallRule) map[string]interface{} {

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

// DecodeTo - receives old and new values as slice of interfaces from schema, decodes and returns firewall properties
func DecodeTo(oldValues, newValues []interface{}) ([]ionoscloud.FirewallruleProperties, []ionoscloud.FirewallruleProperties, error) {
	oldFirewallProperties := make([]ionoscloud.FirewallruleProperties, len(oldValues))
	newFirewallProperties := make([]ionoscloud.FirewallruleProperties, len(newValues))
	err := utils.DecodeInterfaceToStruct(newValues, newFirewallProperties)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %+v to new values of firewall rules %w", newValues, err)
	}
	err = utils.DecodeInterfaceToStruct(oldValues, oldFirewallProperties)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %+v to values of firewall rules %w", oldValues, err)
	}
	return oldFirewallProperties, newFirewallProperties, nil
}

// PropUnsetSetFieldIfNotSetInSchema will only set the in32 types if they actually exist in the schema
// used to unset the fields that are set during the decode process
// mutates fwProp
func PropUnsetSetFieldIfNotSetInSchema(fwProp *ionoscloud.FirewallruleProperties, path string, d *schema.ResourceData) {
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

// GetAndUpdateFirewalls - checks in schema and returns modified firewall rules as a slice of ionoscloud.FirewallRule and also returns a slice of firewall rule ids
func (fs *Service) GetAndUpdateFirewalls(ctx context.Context, dcId, serverId, nicId, path string) (firewallRules []ionoscloud.FirewallRule, firewallRuleIds []string, diags diag.Diagnostics) {
	firewallRuleIds = []string{}
	if fs.D.HasChange(path) {
		oldValues, newValues := fs.D.GetChange(path)
		oldValuesIntf := oldValues.([]interface{})
		newValuesIntf := newValues.([]interface{})
		onlyOld := slice.Difference(oldValuesIntf, newValuesIntf)
		onlyNew := slice.Difference(newValuesIntf, oldValuesIntf)
		oldFirewalls, newFirewalls, err := DecodeTo(onlyOld, onlyNew)
		if err != nil {
			return firewallRules, firewallRuleIds, diag.FromErr(fmt.Errorf("could not get changes for firewall rules %w", err))
		}

		firewallRuleIdsIntf := fs.D.Get("firewallrule_ids").([]interface{})
		firewallRuleIds = slice.AnyToString(firewallRuleIdsIntf)

		if nicId != "" {
			//delete old rules
			for idx := range oldFirewalls {
				//we need the id, but we can't get it from oldFirewalls because it's only the property
				oldId := onlyOld[idx].(map[string]interface{})["id"].(string)

				if deleteRule := !utils.IsValueInSliceOfMap(onlyNew, "id", oldId); deleteRule {
					_, err = fs.Delete(ctx, dcId, serverId, nicId, oldId)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(fmt.Errorf("an error occured while deleting firewall prop for dcId: %s, server_id: %s, "+
							"nic_id %s, ID: %s, (%w)", dcId, serverId, nicId, oldId, err))
					}

					if slice.Contains(firewallRuleIds, oldId) {
						firewallRuleIds = slice.DeleteFrom(firewallRuleIds, oldId)
					}
				}
			}
		}

		// create updated rules
		for idx := range newFirewalls {
			PropUnsetSetFieldIfNotSetInSchema(&newFirewalls[idx], path, fs.D)
			prop := newFirewalls[idx]
			fwRule := ionoscloud.FirewallRule{
				Properties: &prop,
			}
			var firewall *ionoscloud.FirewallRule
			if nicId != "" {
				if id, ok := onlyNew[idx].(map[string]interface{})["id"]; ok && id != "" {
					//do not send protocol, it's an update
					*fwRule.Properties = SetNullableFields(*fwRule.Properties)
					fwRule.Properties.Protocol = nil
					firewall, _, err = fs.Update(ctx, dcId, serverId, nicId, id.(string), fwRule)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(err)
					}
				} else {
					firewall, _, err = fs.Create(ctx, dcId, serverId, nicId, fwRule)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(err)
					}
					firewallRuleIds = append(firewallRuleIds, *firewall.Id)
				}
			}
			firewallRules = append(firewallRules, fwRule)
		}
	}
	return firewallRules, firewallRuleIds, nil
}

func SetNullableFields(prop ionoscloud.FirewallruleProperties) ionoscloud.FirewallruleProperties {
	if prop.SourceIp == nil {
		prop.SetSourceIpNil()
	}
	if prop.SourceMac == nil {
		prop.SetSourceMacNil()
	}
	if prop.IpVersion == nil {
		prop.SetIpVersionNil()
	}
	if prop.TargetIp == nil {
		prop.SetTargetIpNil()
	}
	if prop.IcmpCode == nil {
		prop.SetIcmpCodeNil()
	}
	if prop.IcmpType == nil {
		prop.SetIcmpTypeNil()
	}
	return prop
}

func (fs *Service) AddToMapIfRuleExists(ctx context.Context, datacenterId, serverId, nicId, ruleId string) (map[string]interface{}, error) {
	var firewallEntry map[string]interface{}
	if datacenterId == "" || serverId == "" || nicId == "" || ruleId == "" {
		log.Printf("[DEBUG] Cannot search for firewall rules without dcId %s, serverId %s, nicId %s, ruleId %s", datacenterId, serverId, nicId, ruleId)
		return firewallEntry, nil
	}

	firewall, apiResponse, err := fs.FindById(ctx, datacenterId, serverId, nicId, ruleId)
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return firewallEntry, fmt.Errorf("error, while searching for firewall rule with id %s (%w)", ruleId, err)
		}
	}
	if firewall == nil {
		return firewallEntry, nil
	}
	if firewall.Properties != nil && firewall.Properties.Name != nil {
		log.Printf("[DEBUG] found firewall rule with name %s", *firewall.Properties.Name)
	}
	firewallEntry = SetProperties(*firewall)
	firewallEntry["id"] = *firewall.Id

	return firewallEntry, nil
}

func SetIdsInSchema(d *schema.ResourceData, firewallRuleIds []string) error {
	if len(firewallRuleIds) == 0 {
		return nil
	}
	if err := d.Set("firewallrule_id", firewallRuleIds[0]); err != nil {
		return utils.GenerateSetError(constant.ServerResource, "firewallrule_id", err)
	}
	if err := d.Set("firewallrule_ids", slice.ToAnyList(firewallRuleIds)); err != nil {
		return utils.GenerateSetError(constant.ServerResource, "firewallrule_ids", err)
	}
	return nil
}

func ExtractOrderedFirewallIds(foundRules, sentRules []ionoscloud.FirewallRule) []string {
	var ruleIds = []string{}

	if len(sentRules) == 0 || len(foundRules) == 0 {
		return []string{}
	}

	//keep order of ruleIds
	for _, rule := range sentRules {
		for _, foundRule := range foundRules {
			// computed, make equal for comparison
			if rule.Properties != nil &&
				foundRule.Properties != nil && foundRule.Properties.IpVersion != nil {
				rule.Properties.IpVersion = foundRule.Properties.IpVersion
			}
			//we need deepEqual here, because the structures contain pointers and cannot be compared using the stricter `==`
			if reflect.DeepEqual(rule.Properties, foundRule.Properties) {
				ruleIds = append(ruleIds, *foundRule.Id)
			}
		}
	}
	return ruleIds
}

func SetFwRuleIdsInSchemaInCaseOfProviderUpdate(d *schema.ResourceData) error {
	if _, ok := d.GetOk("firewallrule_ids"); !ok {
		if fwRuleItf, ok := d.GetOk("firewallrule_id"); ok {
			firewallRule := fwRuleItf.(string)
			var firewallRuleIds []string
			firewallRuleIds = append(firewallRuleIds, firewallRule)
			if err := d.Set("firewallrule_ids", firewallRuleIds); err != nil {
				return utils.GenerateSetError("server", "firewallrule_ids", err)
			}
		}
	}
	return nil
}
