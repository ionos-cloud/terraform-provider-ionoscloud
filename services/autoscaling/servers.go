package autoscaling

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func (c *Client) GetGroupServer(ctx context.Context, groupId string, serverId string) (autoscaling.Server, *autoscaling.APIResponse, error) {
	server, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsServersFindById(ctx, serverId, groupId).Execute()
	return server, apiResponse, err
}

func (c *Client) GetAllGroupServers(ctx context.Context, groupId string) (autoscaling.ServerCollection, *autoscaling.APIResponse, error) {
	servers, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsServersGet(ctx, groupId).Execute()
	return servers, apiResponse, err
}

func SetAutoscalingServersData(d *schema.ResourceData, groupServers autoscaling.ServerCollection) diag.Diagnostics {

	if groupServers.Items != nil {
		var servers []interface{}
		for _, groupServer := range *groupServers.Items {
			serverEntry := make(map[string]interface{})
			if groupServer.Id != nil {
				serverEntry["id"] = *groupServer.Id
			} else {
				return diag.FromErr(fmt.Errorf("expected an UUID for server entry but received 'nil'"))
			}
			servers = append(servers, serverEntry)
		}
		err := d.Set("servers", servers)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while setting group servers data: %w", err))
		}
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	return nil
}
