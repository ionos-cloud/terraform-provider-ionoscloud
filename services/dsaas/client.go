package dsaas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Client struct {
	dsaas.APIClient
}

type ClientConfig struct {
	dsaas.Configuration
}

// ClientService is a wrapper around dsaas.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *dsaas.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	newConfigDsaas := dsaas.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDsaas.Debug = true
	}
	newConfigDsaas.MaxRetries = 999
	newConfigDsaas.MaxWaitTime = 4 * time.Second

	newConfigDsaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDsaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dsaas/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dsaas.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &clientService{
		client: dsaas.NewAPIClient(newConfigDsaas),
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
