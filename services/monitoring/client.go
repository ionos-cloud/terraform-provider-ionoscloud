package monitoring

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	monitoring "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"net/http"
	"runtime"
)

// Client is a wrapper for the Monitoring SDK
type MonitoringClient struct {
	Client monitoring.APIClient
}

// NewClient returns a new Monitoring client
func NewClient(username, password, token, url, terraformVersion string, insecure bool) *MonitoringClient {
	config := shared.NewConfiguration(username, password, token, url)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-monitoring/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		monitoring.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) // nolint:staticcheck

	return &MonitoringClient{Client: *monitoring.NewAPIClient(config)}
}
