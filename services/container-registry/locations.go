package container_registry

import (
	"context"
	registry "github.com/ionos-cloud/sdk-go-autoscaling"
)

type LocationsService interface {
	GetAllLocations(ctx context.Context) (registry.LocationsResponse, *registry.APIResponse, error)
}

func (c *Client) GetAllLocations(ctx context.Context) (registry.LocationsResponse, *registry.APIResponse, error) {
	versions, apiResponse, err := c.LocationsApi.LocationsGet(ctx).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err

	}
	return versions, nil, err
}
