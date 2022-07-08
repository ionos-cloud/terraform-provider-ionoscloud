package container_registry

import (
	"context"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
)

type NamesService interface {
	GetNameAvailability(ctx context.Context, name string) (bool, *cr.APIResponse, error)
}

func (c *Client) GetNameAvailability(ctx context.Context, name string) (bool, *cr.APIResponse, error) {
	apiResponse, err := c.NamesApi.NamesFindByName(ctx, name).Execute()
	if err == nil && apiResponse != nil && apiResponse.StatusCode == 200 {
		return true, apiResponse, nil
	}
	return false, nil, err
}
