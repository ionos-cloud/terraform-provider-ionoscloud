package container_registry

import (
	"context"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
)

type RegistriesService interface {
	ListRegistries(ctx context.Context) (cr.RegistriesResponse, *cr.APIResponse, error)
	CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *cr.APIResponse, error)
	DeleteRegistry(ctx context.Context, registryId string) (*cr.APIResponse, error)
	GetRegistry(ctx context.Context, registryId string) (cr.RegistryResponse, *cr.APIResponse, error)
	PatchRegistry(ctx context.Context, registryId string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *cr.APIResponse, error)
	PutRegistry(ctx context.Context, registryId string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *cr.APIResponse, error)
	DeleteRepositories(ctx context.Context, registryId, repositoryId string) (*cr.APIResponse, error)
	ListTokens(ctx context.Context, registryId string) (cr.TokensResponse, *cr.APIResponse, error)
	CreateTokens(ctx context.Context, registryId string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *cr.APIResponse, error)
	DeleteToken(ctx context.Context, registryId, tokenId string) (*cr.APIResponse, error)
	GetToken(ctx context.Context, registryId, tokenId string) (cr.TokenResponse, *cr.APIResponse, error)
	PatchToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *cr.APIResponse, error)
	PutToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *cr.APIResponse, error)
}

func (c *Client) ListRegistries(ctx context.Context) (cr.RegistriesResponse, *cr.APIResponse, error) {
	registry, apiResponse, err := c.RegistriesApi.RegistriesGet(ctx).Execute()
	if apiResponse != nil {
		return registry, apiResponse, err

	}
	return registry, nil, err
}

func (c *Client) CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *cr.APIResponse, error) {
	registry, apiResponse, err := c.RegistriesApi.RegistriesPost(ctx).PostRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registry, apiResponse, err

	}
	return registry, nil, err
}

func (c *Client) DeleteRegistry(ctx context.Context, registryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.RegistriesApi.RegistriesDelete(ctx, registryId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) GetRegistry(ctx context.Context, registryId string) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesFindById(ctx, registryId).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) PatchRegistry(ctx context.Context, registryId string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesPatch(ctx, registryId).PatchRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) PutRegistry(ctx context.Context, registryId string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesPut(ctx, registryId).PutRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) DeleteRepositories(ctx context.Context, registryId, repositoryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.RepositoriesApi.RegistriesRepositoriesDelete(ctx, registryId, repositoryId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) ListTokens(ctx context.Context, registryId string) (cr.TokensResponse, *cr.APIResponse, error) {
	tokens, apiResponse, err := c.TokensApi.RegistriesTokensGet(ctx, registryId).Execute()
	if apiResponse != nil {
		return tokens, apiResponse, err

	}
	return tokens, nil, err
}

func (c *Client) CreateTokens(ctx context.Context, registryId string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPost(ctx, registryId).PostTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) DeleteToken(ctx context.Context, registryId, tokenId string) (*cr.APIResponse, error) {
	apiResponse, err := c.TokensApi.RegistriesTokensDelete(ctx, registryId, tokenId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) GetToken(ctx context.Context, registryId, tokenId string) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensFindById(ctx, registryId, tokenId).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) PatchToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPatch(ctx, registryId, tokenId).PatchTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) PutToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPut(ctx, registryId, tokenId).PutTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}
