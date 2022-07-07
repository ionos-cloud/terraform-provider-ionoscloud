package container_registry

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	registry "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Client struct {
	registry.APIClient
}

type ClientConfig struct {
	registry.Configuration
}

// ClientService is a wrapper around registry.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *registry.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	newConfigregistry := registry.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigregistry.Debug = true
	}
	newConfigregistry.MaxRetries = 999
	newConfigregistry.MaxWaitTime = 4 * time.Second

	newConfigregistry.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigregistry.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-registry/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, registry.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &clientService{
		client: registry.NewAPIClient(newConfigregistry),
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
