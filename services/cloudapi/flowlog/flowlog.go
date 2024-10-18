package cloudapiflowlog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var FlowlogSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The resource's unique identifier.",
		},
		"action": {
			Type:             schema.TypeString,
			Description:      "Specifies the traffic direction pattern. Valid values: ACCEPTED, REJECTED, ALL. Immutable, forces re-recreation of the nic resource.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ACCEPTED", "REJECTED", "ALL"}, true)),
			DiffSuppressFunc: utils.DiffToLower,
			Required:         true,
		},
		"bucket": {
			Type:        schema.TypeString,
			Description: "The bucket name of an existing IONOS Object Storage bucket. Immutable, forces re-recreation of the nic resource.",
			Required:    true,
		},
		"direction": {
			Type:             schema.TypeString,
			Description:      "Specifies the traffic direction pattern. Valid values: INGRESS, EGRESS, BIDIRECTIONAL. Immutable, forces re-recreation of the nic resource.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INGRESS", "EGRESS", "BIDIRECTIONAL"}, true)),
			DiffSuppressFunc: utils.DiffToLower,
			Required:         true,
		},
		"name": {
			Type:             schema.TypeString,
			Description:      "The resource name.",
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
	},
}

var FlowlogSchemaDatasource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The resource's unique identifier.",
		},
		"action": {
			Type:        schema.TypeString,
			Description: "Specifies the traffic direction pattern. Valid values: ACCEPTED, REJECTED, ALL.",
			Computed:    true,
		},
		"bucket": {
			Type:        schema.TypeString,
			Description: "The bucket name of an existing IONOS Object Storage bucket.",
			Computed:    true,
		},
		"direction": {
			Type:        schema.TypeString,
			Description: "Specifies the traffic direction pattern. Valid values: INGRESS, EGRESS, BIDIRECTIONAL.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The resource name.",
			Computed:    true,
		},
	},
}

type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

func (fw *Service) CreateOrPatchForServer(ctx context.Context, dcId, srvID, nicID, ID string, flowLog ionoscloud.FlowLog) error {
	if ID == "" {
		_, _, err := fw.Client.FlowLogsApi.DatacentersServersNicsFlowlogsPost(ctx, dcId, srvID, nicID).Flowlog(flowLog).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while creating flowlog in datacenter %s, server %s nic %s : %w", dcId, srvID, nicID, err)
		}
	} else {
		_, _, err := fw.Client.FlowLogsApi.DatacentersServersNicsFlowlogsPatch(ctx, dcId, srvID, nicID, ID).Flowlog(*flowLog.Properties).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while updating flowlog %s datacenter %s, server %s nic %s : %w", ID, dcId, srvID, nicID, err)
		}
	}
	return nil
}

func (fw *Service) CreateOrPatchForNLB(ctx context.Context, dcId, nlbID, ID string, flowLog ionoscloud.FlowLog) error {
	if ID == "" {
		_, _, err := fw.Client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPost(ctx, dcId, nlbID).NetworkLoadBalancerFlowLog(flowLog).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while creating flowlog in datacenter %s, nlb %s : %w", dcId, nlbID, err)
		}
	} else {
		_, _, err := fw.Client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPatch(ctx, dcId, nlbID, ID).NetworkLoadBalancerFlowLogProperties(*flowLog.Properties).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while updating flowlog %s datacenter %s, nlb %s : %w", ID, dcId, nlbID, err)
		}
	}
	return nil
}

func (fw *Service) CreateOrPatchForALB(ctx context.Context, dcId, albID, ID string, flowLog ionoscloud.FlowLog) error {
	if ID == "" {
		_, _, err := fw.Client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPost(ctx, dcId, albID).ApplicationLoadBalancerFlowLog(flowLog).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while creating flowlog in datacenter %s, alb %s : %w", dcId, albID, err)
		}
	} else {
		_, _, err := fw.Client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPatch(ctx, dcId, albID, ID).ApplicationLoadBalancerFlowLogProperties(*flowLog.Properties).Execute()
		if err != nil {
			return fmt.Errorf("error occurred while updating flowlog %s, datacenter %s, alb %s : %w", ID, dcId, albID, err)
		}
	}
	return nil
}

// GetFlowLogForALB - there can be only one flowlog per alb
func (fw *Service) GetFlowLogForALB(ctx context.Context, dcId, albID string, depth int32) (*ionoscloud.FlowLog, *ionoscloud.APIResponse, error) {
	flowLogs, apiResponse, err := fw.Client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsGet(ctx, dcId, albID).Depth(depth).Execute()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("error occurred while finding datacenter %s, alb %s : %w", dcId, albID, err)
	}
	if flowLogs.Items != nil && len(*flowLogs.Items) > 0 {
		return &(*flowLogs.Items)[0], apiResponse, nil
	}
	return nil, apiResponse, nil
}

// Delete - this method actually does not work for now
func (fw *Service) Delete(ctx context.Context, dcId string, srvID string, nicID, ID string) error {
	_, err := fw.Client.FlowLogsApi.DatacentersServersNicsFlowlogsDelete(ctx, dcId, srvID, nicID, ID).Execute()
	if err != nil {
		return fmt.Errorf("error occurred while deleting flowlog %s datacenter %s, server %s nic %s : %w", ID, dcId, srvID, nicID, err)
	}
	return nil
}
func GetFlowlogFromMap(flowLogMap map[string]any) ionoscloud.FlowLog {
	flowlog := ionoscloud.NewFlowLog(*ionoscloud.NewFlowLogProperties("", "", "", ""))
	*flowlog.Properties.Action = flowLogMap["action"].(string)
	*flowlog.Properties.Bucket = flowLogMap["bucket"].(string)
	*flowlog.Properties.Direction = flowLogMap["direction"].(string)
	*flowlog.Properties.Name = flowLogMap["name"].(string)
	return *flowlog
}
