package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
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
		client = ionoscloud.NewAPIClient(ionoscloud.NewConfiguration("", "", c.Token))
	} else {
		client = ionoscloud.NewAPIClient(ionoscloud.NewConfiguration(c.Username, c.Password, ""))
	}
	client.GetConfig().UserAgent = fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", terraformVersion, meta.SDKVersionString())

	log.Printf("[DEBUG] Terraform client UA set to %s", client.GetConfig().UserAgent)

	client.GetConfig().AddDefaultQueryParam("depth", "5")

	if len(c.Endpoint) > 0 {
		client.GetConfig().Host = c.Endpoint
	}
	return client, nil
}
