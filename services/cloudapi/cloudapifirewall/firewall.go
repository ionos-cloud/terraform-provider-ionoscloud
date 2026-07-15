package cloudapifirewall

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Service struct {
	Client *ionoscloud.APIClient
	Meta   any
	D      *schema.ResourceData
}

func (fs *Service) Get(ctx context.Context, datacenterID, serverID, nicID string, depth int32) ([]ionoscloud.FirewallRule, error) {
	firewallRules, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterID, serverID, nicID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	if apiResponse.HttpNotFound() || len(firewallRules.Items) == 0 {
		tflog.Debug(ctx, "no firewalls found", map[string]any{"datacenter_id": datacenterID, "server_id": serverID, "nic_id": nicID})
		return nil, nil
	}
	return firewallRules.Items, nil
}

func (fs *Service) FindById(ctx context.Context, datacenterID, serverID, nicID, firewallID string) (*ionoscloud.FirewallRule, *shared.APIResponse, error) {
	firewallRule, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterID, serverID, nicID, firewallID).Depth(1).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, err
	}
	if apiResponse.HttpNotFound() {
		tflog.Debug(ctx, "no firewall rule found", map[string]any{"datacenter_id": datacenterID, "server_id": serverID, "nic_id": nicID})
		return nil, apiResponse, nil
	}
	return &firewallRule, apiResponse, nil
}

func (fs *Service) Delete(ctx context.Context, datacenterID, serverID, nicID, firewallID string) (*shared.APIResponse, error) {
	apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDelete(ctx, datacenterID, serverID, nicID, firewallID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	if errState := bundleclient.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutDelete); errState != nil {
		return apiResponse, fmt.Errorf("on delete an error occurred while waiting for state change dcID: %s, server_id: %s, nic_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, nicID, firewallID, errState)
	}
	return apiResponse, nil
}

func (fs *Service) Create(ctx context.Context, datacenterID, serverID, nicID string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *shared.APIResponse, error) {
	firewall, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, datacenterID, serverID, nicID).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occurred while creating firewall rule for dcID: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterID, serverID, nicID, err)
	}
	if errState := bundleclient.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutCreate); errState != nil {
		return nil, apiResponse, fmt.Errorf("on create an error occurred while waiting for state change dcID: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterID, serverID, nicID, errState)
	}
	return &firewall, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterID, serverID, nicID, id string, firewallrule ionoscloud.FirewallRule) (*ionoscloud.FirewallRule, *shared.APIResponse, error) {
	firewall, apiResponse, err := fs.Client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPut(ctx, datacenterID, serverID, nicID, id).Firewallrule(firewallrule).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occurred while updating firewall rule for dcID: %s, server_id: %s, nic_id: %s, id %s, Response: (%w)", datacenterID, serverID, nicID, id, err)
	}
	if errState := bundleclient.WaitForStateChange(ctx, fs.Meta, fs.D, apiResponse, schema.TimeoutUpdate); errState != nil {
		return nil, apiResponse, fmt.Errorf("on update an error occurred while waiting for state change dcID: %s, server_id: %s, nic_id: %s, Response: (%w)", datacenterID, serverID, nicID, errState)
	}
	return &firewall, apiResponse, nil
}

func SetProperties(firewall ionoscloud.FirewallRule) map[string]any {

	fw := map[string]any{}
	utils.SetPropWithNilCheck(fw, "protocol", firewall.Properties.Protocol)
	utils.SetPropWithNilCheck(fw, "name", firewall.Properties.Name)
	utils.SetPropWithNilCheck(fw, "source_mac", firewall.Properties.SourceMac.Get())
	utils.SetPropWithNilCheck(fw, "source_ip", firewall.Properties.SourceIp.Get())
	utils.SetPropWithNilCheck(fw, "target_ip", firewall.Properties.TargetIp.Get())
	utils.SetPropWithNilCheck(fw, "port_range_start", firewall.Properties.PortRangeStart)
	utils.SetPropWithNilCheck(fw, "port_range_end", firewall.Properties.PortRangeEnd)
	utils.SetPropWithNilCheck(fw, "type", firewall.Properties.Type)
	if firewall.Properties.IcmpType.IsSet() {
		fw["icmp_type"] = strconv.Itoa(int(*firewall.Properties.IcmpType.Get()))
	}
	if firewall.Properties.IcmpCode.IsSet() {
		fw["icmp_code"] = strconv.Itoa(int(*firewall.Properties.IcmpCode.Get()))
	}
	return fw
}

// DecodeTo - receives old and new values as slice of interfaces from schema, decodes and returns firewall properties
func DecodeTo(ctx context.Context, oldValues, newValues []any) ([]ionoscloud.FirewallruleProperties, []ionoscloud.FirewallruleProperties, error) {
	oldFirewallProperties := make([]ionoscloud.FirewallruleProperties, len(oldValues))
	newFirewallProperties := make([]ionoscloud.FirewallruleProperties, len(newValues))
	err := utils.DecodeInterfaceToStruct(ctx, newValues, newFirewallProperties)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode from %+v to new values of firewall rules %w", newValues, err)
	}
	err = utils.DecodeInterfaceToStruct(ctx, oldValues, oldFirewallProperties)
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
func (fs *Service) GetAndUpdateFirewalls(ctx context.Context, dcID, serverID, nicID, path string) (firewallRules []ionoscloud.FirewallRule, firewallRuleIDs []string, diags diag.Diagnostics) {
	firewallRuleIDs = []string{}
	if fs.D.HasChange(path) {
		oldValues, newValues := fs.D.GetChange(path)
		oldValuesIntf := oldValues.([]any)
		newValuesIntf := newValues.([]any)
		onlyOld := slice.Difference(oldValuesIntf, newValuesIntf)
		onlyNew := slice.Difference(newValuesIntf, oldValuesIntf)
		oldFirewalls, newFirewalls, err := DecodeTo(ctx, onlyOld, onlyNew)
		if err != nil {
			return firewallRules, firewallRuleIDs, diag.FromErr(fmt.Errorf("could not get changes for firewall rules %w", err))
		}

		firewallRuleIdsIntf := fs.D.Get("firewallrule_ids").([]any)
		firewallRuleIDs = slice.AnyToString(firewallRuleIdsIntf)

		if nicID != "" {
			// delete old rules
			for idx := range oldFirewalls {
				// we need the id, but we can't get it from oldFirewalls because it's only the property
				oldID := onlyOld[idx].(map[string]any)["id"].(string)

				if deleteRule := !utils.IsValueInSliceOfMap(onlyNew, "id", oldID); deleteRule {
					_, err = fs.Delete(ctx, dcID, serverID, nicID, oldID)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(fmt.Errorf("an error occurred while deleting firewall prop for dcID: %s, server_id: %s, "+
							"nic_id %s, ID: %s, (%w)", dcID, serverID, nicID, oldID, err))
					}

					if slice.Contains(firewallRuleIDs, oldID) {
						firewallRuleIDs = slice.DeleteFrom(firewallRuleIDs, oldID)
					}
				}
			}
		}

		// create updated rules
		for idx := range newFirewalls {
			PropUnsetSetFieldIfNotSetInSchema(&newFirewalls[idx], path, fs.D)
			prop := newFirewalls[idx]
			fwRule := ionoscloud.FirewallRule{
				Properties: prop,
			}
			var firewall *ionoscloud.FirewallRule
			if nicID != "" {
				if id, ok := onlyNew[idx].(map[string]any)["id"]; ok && id != "" {
					// do not send protocol, it's an update
					fwRule.Properties = SetNullableFields(fwRule.Properties)
					fwRule.Properties.Protocol = nil
					firewall, _, err = fs.Update(ctx, dcID, serverID, nicID, id.(string), fwRule)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(err)
					}
				} else {
					firewall, _, err = fs.Create(ctx, dcID, serverID, nicID, fwRule)
					if err != nil {
						return firewallRules, []string{}, diag.FromErr(err)
					}
					firewallRuleIDs = append(firewallRuleIDs, *firewall.Id)
				}
			}
			firewallRules = append(firewallRules, fwRule)
		}
	}
	return firewallRules, firewallRuleIDs, nil
}

func SetNullableFields(prop ionoscloud.FirewallruleProperties) ionoscloud.FirewallruleProperties {
	if !prop.SourceIp.IsSet() {
		prop.SetSourceIpNil()
	}
	if !prop.SourceMac.IsSet() {
		prop.SetSourceMacNil()
	}
	if !prop.IpVersion.IsSet() {
		prop.SetIpVersionNil()
	}
	if !prop.TargetIp.IsSet() {
		prop.SetTargetIpNil()
	}
	if !prop.IcmpCode.IsSet() {
		prop.SetIcmpCodeNil()
	}
	if !prop.IcmpType.IsSet() {
		prop.SetIcmpTypeNil()
	}
	return prop
}

func (fs *Service) AddToMapIfRuleExists(ctx context.Context, datacenterID, serverID, nicID, ruleID string) (map[string]any, error) {
	var firewallEntry map[string]any
	if datacenterID == "" || serverID == "" || nicID == "" || ruleID == "" {
		tflog.Debug(ctx, "cannot search for firewall rules: missing IDs", map[string]any{"datacenter_id": datacenterID, "server_id": serverID, "nic_id": nicID, "rule_id": ruleID})
		return firewallEntry, nil
	}

	firewall, apiResponse, err := fs.FindById(ctx, datacenterID, serverID, nicID, ruleID)
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return firewallEntry, fmt.Errorf("error, while searching for firewall rule with id %s (%w)", ruleID, err)
		}
	}
	if firewall == nil {
		return firewallEntry, nil
	}
	if firewall.Properties.Name != nil {
		tflog.Debug(ctx, "found firewall rule", map[string]any{"name": *firewall.Properties.Name})
	}
	firewallEntry = SetProperties(*firewall)
	firewallEntry["id"] = *firewall.Id

	return firewallEntry, nil
}

func SetIdsInSchema(d *schema.ResourceData, firewallRuleIDs []string) error {
	if len(firewallRuleIDs) == 0 {
		return nil
	}
	if err := d.Set("firewallrule_id", firewallRuleIDs[0]); err != nil {
		return utils.GenerateSetError(constant.ServerResource, "firewallrule_id", err)
	}
	if err := d.Set("firewallrule_ids", slice.ToAnyList(firewallRuleIDs)); err != nil {
		return utils.GenerateSetError(constant.ServerResource, "firewallrule_ids", err)
	}
	return nil
}

func ExtractOrderedFirewallIDs(foundRules, sentRules []ionoscloud.FirewallRule) []string {
	var ruleIDs = []string{}

	if len(sentRules) == 0 || len(foundRules) == 0 {
		return []string{}
	}

	// keep order of ruleIDs
	for _, rule := range sentRules {
		for _, foundRule := range foundRules {
			// computed, make equal for comparison
			if foundRule.Properties.IpVersion.IsSet() {
				rule.Properties.IpVersion = foundRule.Properties.IpVersion
			}
			// we need deepEqual here, because the structures contain pointers and cannot be compared using the stricter `==`
			if reflect.DeepEqual(rule.Properties, foundRule.Properties) {
				ruleIDs = append(ruleIDs, *foundRule.Id)
			}
		}
	}
	return ruleIDs
}

func SetFwRuleIdsInSchemaInCaseOfProviderUpdate(d *schema.ResourceData) error {
	if _, ok := d.GetOk("firewallrule_ids"); !ok {
		if fwRuleItf, ok := d.GetOk("firewallrule_id"); ok {
			firewallRule := fwRuleItf.(string)
			var firewallRuleIDs []string
			firewallRuleIDs = append(firewallRuleIDs, firewallRule)
			if err := d.Set("firewallrule_ids", firewallRuleIDs); err != nil {
				return utils.GenerateSetError("server", "firewallrule_ids", err)
			}
		}
	}
	return nil
}
