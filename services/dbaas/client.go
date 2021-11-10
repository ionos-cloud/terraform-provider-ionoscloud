package dbaas

import (
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"os"
)

type Client struct {
	dbaas.APIClient
}

type ClientConfig struct {
	dbaas.Configuration
}

// ClientService is a wrapper around dbaas.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *dbaas.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url string) ClientService {
	newConfigDbaas := dbaas.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDbaas.Debug = true
	}

	return &clientService{
		client: dbaas.NewAPIClient(newConfigDbaas),
	}
}

func (c clientService) Get() *Client {
	return &Client{
		APIClient: *c.client,
	}
}

func (c clientService) GetConfig() *ClientConfig {
	return &ClientConfig{
		Configuration: *c.client.GetConfig(),
	}
}
