package autoscaling

import (
	"context"

	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func (c *Client) GetAction(ctx context.Context, groupId string, actionId string) (autoscaling.Action, *autoscaling.APIResponse, error) {
	action, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsFindById(ctx, actionId, groupId).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return action, apiResponse, err

	}
	return action, nil, err
}

func (c *Client) GetAllActions(ctx context.Context, groupId string) (autoscaling.ActionCollection, *autoscaling.APIResponse, error) {
	actions, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsGet(ctx, groupId).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return actions, apiResponse, err
	}
	return actions, nil, err
}
