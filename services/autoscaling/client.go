package autoscaling

import (
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"os"
)

type Client struct {
	autoscaling.APIClient
}

type ClientConfig struct {
	autoscaling.Configuration
}

// ClientService is a wrapper around autoscaling.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *autoscaling.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url string) ClientService {
	newConfigAutoscaling := autoscaling.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigAutoscaling.Debug = true
	}

	return &clientService{
		client: autoscaling.NewAPIClient(newConfigAutoscaling),
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
