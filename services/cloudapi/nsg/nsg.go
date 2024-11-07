package nsg

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
)

// Service implements utility methods for the Network Security Group
type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

// PutServerNSG updates the security groups of a server
func (nsg *Service) PutServerNSG(ctx context.Context, dcID, serverID string, rawIDs []any) diag.Diagnostics {
	if dcID == "" || serverID == "" {
		return diag.Errorf("dcID and serverID must be set")
	}
	if len(rawIDs) == 0 {
		return nil
	}
	ids := make([]string, 0)
	for _, rawID := range rawIDs {
		if rawID != nil {
			id := rawID.(string)
			ids = append(ids, id)
		}
	}
	_, apiResponse, err := nsg.Client.SecurityGroupsApi.DatacentersServersSecuritygroupsPut(ctx, dcID, serverID).Securitygroups(*ionoscloud.NewListOfIds(ids)).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, nsg.Meta, nsg.D, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			nsg.D.SetId("")
		}
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for securitygroup state change on put. dcID: %s, server_id: %s, Response: (%w)", dcID, serverID, errState))
	}
	return nil
}

// PutNICNSG updates the security groups of a NIC
func (nsg *Service) PutNICNSG(ctx context.Context, dcID, serverID, nicID string, rawIDs []any) diag.Diagnostics {
	if dcID == "" || serverID == "" || nicID == "" {
		return diag.Errorf("dcID and serverID must be set")
	}
	if len(rawIDs) == 0 {
		return nil
	}
	ids := make([]string, 0)
	for _, rawId := range rawIDs {
		if rawId != nil {
			id := rawId.(string)
			ids = append(ids, id)
		}
	}
	_, apiResponse, err := nsg.Client.SecurityGroupsApi.DatacentersServersNicsSecuritygroupsPut(ctx, dcID, serverID, nicID).Securitygroups(*ionoscloud.NewListOfIds(ids)).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, nsg.Meta, nsg.D, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			nsg.D.SetId("")
		}
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for securitygroup state change on put. dcID: %s, server_id: %s, nic_id %s, Response: (%w)", dcID, serverID, nicID, errState))
	}
	return nil
}

// SetNSGInResourceData sets the security groups in the schema
func SetNSGInResourceData(d *schema.ResourceData, items *[]ionoscloud.SecurityGroup) error {
	if items == nil {
		return nil
	}
	nsgIDs := make([]string, 0)
	for _, group := range *items {
		if group.Id != nil {
			id := *group.Id
			nsgIDs = append(nsgIDs, id)
		}
	}
	if err := d.Set("security_groups_ids", nsgIDs); err != nil {
		return fmt.Errorf("error setting security_groups_ids %w", err)
	}
	return nil
}
