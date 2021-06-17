package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	client.GetConfig().UserAgent = fmt.Sprintf("HashiCorp Terraform/%s Terraform Plugin SDK/%s Terraform Provider Ionoscloud/%s", terraformVersion, meta.SDKVersionString(), Version)

	log.Printf("[DEBUG] Terraform client UA set to %s", client.GetConfig().UserAgent)

	client.GetConfig().AddDefaultQueryParam("depth", "5")

	if len(c.Endpoint) > 0 {
		client.GetConfig().Host = c.Endpoint
	}
	return client, nil
}
