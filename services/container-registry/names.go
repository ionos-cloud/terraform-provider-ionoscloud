package container_registry

import (
	"context"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
)

type NamesService interface {
	GetNameAvailability(ctx context.Context, name string) (bool, *cr.APIResponse, error)
}

func (c *Client) GetNameAvailability(ctx context.Context, name string) (bool, *cr.APIResponse, error) {
	apiResponse, err := c.NamesApi.NamesCheckUsage(ctx, name).Execute()
	if err != nil {
		return false, apiResponse, nil
	}
	return true, nil, err
}