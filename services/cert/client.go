package cert

import (
	certmanager "github.com/ionos-cloud/sdk-cert-go"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
)

type Client struct {
	certmanager.APIClient
}

type ClientConfig struct {
	certmanager.Configuration
}

// ClientService is a wrapper around dbaas.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *certmanager.APIClient
}

func NewClientService(username, password, token, url string) ClientService {
	certConfig := certmanager.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		certConfig.Debug = true
	}
	certConfig.MaxRetries = utils.MaxRetries
	certConfig.MaxWaitTime = utils.MaxWaitTime

	certConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}

	return &clientService{
		client: certmanager.NewAPIClient(certConfig),
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
