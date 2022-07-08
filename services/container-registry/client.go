package container_registry

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Client struct {
	cr.APIClient
}

type ClientConfig struct {
	cr.Configuration
}

// ClientService is a wrapper around cr.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *cr.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	newConfigregistry := cr.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigregistry.Debug = true
	}
	newConfigregistry.MaxRetries = 999
	newConfigregistry.MaxWaitTime = 4 * time.Second

	newConfigregistry.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigregistry.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cr.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &clientService{
		client: cr.NewAPIClient(newConfigregistry),
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
