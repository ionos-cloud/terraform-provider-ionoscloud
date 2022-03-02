package dbaas

import (
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"net"
	"net/http"
	"os"
	"time"
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
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 2 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: createTransport()}

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

func createTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   3,
		MaxConnsPerHost:       3,
	}
}
