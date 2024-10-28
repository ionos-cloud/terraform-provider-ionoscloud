package autoscaling

import (
	"context"

	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func (c *Client) GetGroup(ctx context.Context, groupId string, depth float32) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsFindById(ctx, groupId).Depth(depth).Execute()
	apiResponse.LogInfo()
	return group, apiResponse, err
}

func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error) {
	groups, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsGet(ctx).Execute()
	apiResponse.LogInfo()
	return groups, apiResponse, err
}

func (c *Client) CreateGroup(ctx context.Context, group autoscaling.GroupPost) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPost(ctx).GroupPost(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPut(ctx, groupId).GroupPut(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error) {
	apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsDelete(ctx, groupId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}
