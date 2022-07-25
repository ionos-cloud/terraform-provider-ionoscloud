package dataplatform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dataplatform "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Client struct {
	dataplatform.APIClient
}

type ClientConfig struct {
	dataplatform.Configuration
}

// ClientService is a wrapper around dataplatform.APIClient
type ClientService interface {
	Get() *Client
	GetConfig() *ClientConfig
}

type clientService struct {
	client *dataplatform.APIClient
}

var _ ClientService = &clientService{}

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	newConfigDataplatform := dataplatform.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDataplatform.Debug = true
	}
	newConfigDataplatform.MaxRetries = 999
	newConfigDataplatform.MaxWaitTime = 4 * time.Second

	newConfigDataplatform.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDataplatform.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dataplatform.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &clientService{
		client: dataplatform.NewAPIClient(newConfigDataplatform),
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
