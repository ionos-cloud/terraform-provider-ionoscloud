package dbaas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
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

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	newConfigDbaas := dbaas.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 4 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dbaas.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

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
