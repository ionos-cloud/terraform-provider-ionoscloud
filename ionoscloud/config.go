package ionoscloud

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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
		client = ionoscloud.NewAPIClient(shared.NewConfiguration("", "", c.Token, c.Endpoint))
	} else {
		client = ionoscloud.NewAPIClient(shared.NewConfiguration(c.Username, c.Password, "", c.Endpoint))
	}

	client.GetConfig().UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		Version, ionoscloud.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	log.Printf("[DEBUG] Terraform client UA set to %s", client.GetConfig().UserAgent)

	client.GetConfig().AddDefaultQueryParam("depth", "5")

	if len(c.Endpoint) > 0 {
		client.GetConfig().Host = c.Endpoint
	}
	return client, nil
}
