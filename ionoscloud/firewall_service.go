package ionoscloud

import (
	"context"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

type FirewallService struct {
	client *ionoscloud.APIClient
}

func (fs *FirewallService) firewallsGet(ctx context.Context, datacenterId, serverId, nicId string) ([]ionoscloud.FirewallRule, error) {
	firewallRules, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Execute()
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

func (fs *FirewallService) firewallFindById(ctx context.Context, datacenterId, serverId, nicId, fwId string) (*ionoscloud.FirewallRule, *ionoscloud.APIResponse, error) {
	firewallRule, apiResponse, err := fs.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, fwId).Depth(1).Execute()
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
