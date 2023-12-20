package autoscaling

import (
	"context"

	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func (c *Client) GetAction(ctx context.Context, groupId string, actionId string) (autoscaling.Action, *autoscaling.APIResponse, error) {
	action, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsFindById(ctx, actionId, groupId).Execute()
	apiResponse.LogInfo()
	return action, apiResponse, err
}

func (c *Client) GetAllActions(ctx context.Context, groupId string) (autoscaling.ActionCollection, *autoscaling.APIResponse, error) {
	actions, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsGet(ctx, groupId).Execute()
	apiResponse.LogInfo()
	return actions, apiResponse, err
}
