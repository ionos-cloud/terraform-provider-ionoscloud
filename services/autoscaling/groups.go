package autoscaling

import (
	"context"

	autoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetGroup returns a group given an ID
func (c *Client) GetGroup(ctx context.Context, groupId string, depth float32) (autoscaling.Group, *shared.APIResponse, error) {
	group, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsFindById(ctx, groupId).Depth(depth).Execute()
	apiResponse.LogInfo()
	return group, apiResponse, err
}

// ListGroups returns a list of all groups
func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *shared.APIResponse, error) {
	groups, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsGet(ctx).Execute()
	apiResponse.LogInfo()
	return groups, apiResponse, err
}

// CreateGroup creates a new group
func (c *Client) CreateGroup(ctx context.Context, group autoscaling.GroupPost) (autoscaling.GroupPostResponse, *shared.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPost(ctx).GroupPost(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

// UpdateGroup updates a group given an ID
func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupPut) (autoscaling.Group, *shared.APIResponse, error) {
	groupResponse, apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsPut(ctx, groupId).GroupPut(group).Execute()
	apiResponse.LogInfo()
	return groupResponse, apiResponse, err
}

// DeleteGroup deletes a group given an ID
func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.AutoScalingGroupsApi.GroupsDelete(ctx, groupId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}
