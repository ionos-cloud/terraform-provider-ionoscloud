package cert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
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

func NewClientService(username, password, token, url, version, terraformVersion string) ClientService {
	certConfig := certmanager.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		certConfig.Debug = true
	}
	certConfig.MaxRetries = utils.MaxRetries
	certConfig.MaxWaitTime = utils.MaxWaitTime

	certConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	certConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, certmanager.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

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
