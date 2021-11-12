package autoscaling

import (
	"context"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ActionsService interface {
	GetAction(ctx context.Context, groupId string, actionId string) (autoscaling.Action, *autoscaling.APIResponse, error)
	GetAllActions(ctx context.Context, groupId string) (autoscaling.ActionCollection, *autoscaling.APIResponse, error)
}

func (c *Client) GetAction(ctx context.Context, groupId string, actionId string) (autoscaling.Action, *autoscaling.APIResponse, error) {
	action, apiResponse, err := c.GroupsApi.AutoscalingGroupsActionsFindById(ctx, actionId, groupId).Execute()
	if apiResponse != nil {
		return action, apiResponse, err

	}
	return action, nil, err
}

func (c *Client) GetAllActions(ctx context.Context, groupId string) (autoscaling.ActionCollection, *autoscaling.APIResponse, error) {
	actions, apiResponse, err := c.GroupsApi.AutoscalingGroupsActionsGet(ctx, groupId).Execute()
	if apiResponse != nil {
		return actions, apiResponse, err
	}
	return actions, nil, err
}
