package autoscaling

import (
	"context"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
)

type GroupService interface {
	GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error)
	ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error)
	CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error)
	UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupUpdate) (autoscaling.Group, *autoscaling.APIResponse, error)
	DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error)
}

func (c *Client) GetGroup(ctx context.Context, groupId string) (autoscaling.Group, *autoscaling.APIResponse, error) {
	group, apiResponse, err := c.GroupsApi.AutoscalingGroupsFindById(ctx, groupId).Execute()
	if apiResponse != nil {
		return group, apiResponse, err

	}
	return group, nil, err
}

func (c *Client) ListGroups(ctx context.Context) (autoscaling.GroupCollection, *autoscaling.APIResponse, error) {
	groups, apiResponse, err := c.GroupsApi.AutoscalingGroupsGet(ctx).Execute()
	if apiResponse != nil {
		return groups, apiResponse, err
	}
	return groups, nil, err
}

func (c *Client) CreateGroup(ctx context.Context, group autoscaling.Group) (autoscaling.GroupPostResponse, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.GroupsApi.AutoscalingGroupsPost(ctx).Group(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) UpdateGroup(ctx context.Context, groupId string, group autoscaling.GroupUpdate) (autoscaling.Group, *autoscaling.APIResponse, error) {
	groupResponse, apiResponse, err := c.GroupsApi.AutoscalingGroupsPut(ctx, groupId).Group(group).Execute()
	if apiResponse != nil {
		return groupResponse, apiResponse, err
	}
	return groupResponse, nil, err
}

func (c *Client) DeleteGroup(ctx context.Context, groupId string) (*autoscaling.APIResponse, error) {
	apiResponse, err := c.GroupsApi.AutoscalingGroupsDelete(ctx, groupId).Execute()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}
