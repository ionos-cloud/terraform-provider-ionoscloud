package ionoscloud

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// Config represents
type Config struct {
	Username string
	Password string
	Endpoint string
	Retries  int
	Token    string
}

// Client returns a new client for accessing ionoscloud.
func (c *Config) Client(terraformVersion string) (*ionoscloud.APIClient, error) {
	var client *ionoscloud.APIClient
	if c.Token != "" {
		client = ionoscloud.NewAPIClient(ionoscloud.NewConfiguration("", "", c.Token, c.Endpoint))
	} else {
		client = ionoscloud.NewAPIClient(ionoscloud.NewConfiguration(c.Username, c.Password, "", c.Endpoint))
	}

	client.GetConfig().UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		Version, ionoscloud.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	tflog.Debug(context.Background(), "terraform client UA set", map[string]interface{}{"user_agent": client.GetConfig().UserAgent})

	client.GetConfig().AddDefaultQueryParam("depth", "5")

	if len(c.Endpoint) > 0 {
		client.GetConfig().Host = c.Endpoint
	}
	return client, nil
}
