package autoscaling

import (
	"context"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ServersService interface {
	GetServer(ctx context.Context, groupId string, serverId string) (autoscaling.Server, *autoscaling.APIResponse, error)
	GetAllServers(ctx context.Context, groupId string) (autoscaling.ServerCollection, *autoscaling.APIResponse, error)
}

func (c *Client) GetServer(ctx context.Context, groupId string, serverId string) (autoscaling.Server, *autoscaling.APIResponse, error) {
	server, apiResponse, err := c.GroupsApi.AutoscalingGroupsServersFindById(ctx, serverId, groupId).Execute()
	if apiResponse != nil {
		return server, apiResponse, err

	}
	return server, nil, err
}

func (c *Client) GetAllServers(ctx context.Context, groupId string) (autoscaling.ServerCollection, *autoscaling.APIResponse, error) {
	servers, apiResponse, err := c.GroupsApi.AutoscalingGroupsServersGet(ctx, groupId).Execute()
	if apiResponse != nil {
		return servers, apiResponse, err
	}
	return servers, nil, err
}
