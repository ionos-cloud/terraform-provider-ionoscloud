package autoscaling

import (
	"context"

	autoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func (c *Client) GetAction(ctx context.Context, groupId string, actionId string) (autoscaling.Action, *shared.APIResponse, error) {
	action, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsFindById(ctx, actionId, groupId).Execute()
	apiResponse.LogInfo()
	return action, apiResponse, err
}

func (c *Client) GetAllActions(ctx context.Context, groupId string) (autoscaling.ActionCollection, *shared.APIResponse, error) {
	actions, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsActionsGet(ctx, groupId).Execute()
	apiResponse.LogInfo()
	return actions, apiResponse, err
}
