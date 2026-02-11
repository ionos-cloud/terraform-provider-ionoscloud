package cloudapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Config struct {
	fileConfig   *fileconfiguration.FileConfig
	clientConfig ionoscloud.Configuration
}

// NewConfig returns a new Config struct based on the provided client options and configuration file.
func NewConfig(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Config {
	config := ionoscloud.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, ionoscloud.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)

	return &Config{
		fileConfig:   fileConfig,
		clientConfig: *config,
	}
}

// CopyClientConfig creates a deep copy of the client configuration to ensure that modifications to the returned configuration
// do not affect the original configuration stored in Config.
func (c Config) CopyClientConfig() ionoscloud.Configuration {
	// Copy directly what can be copied with a simple assignment
	// Also allows for unmarshalling into the Logger interface without issues
	newConfig := c.clientConfig

	// Create a copy of the HTTPClient if it exists
	// http.DefaultClient is the control value
	httpClientCopy := http.DefaultClient
	if c.clientConfig.HTTPClient != nil {
		httpClientCopy = c.clientConfig.HTTPClient
		// Temporarily set to nil to avoid deepcopy issues with http.Client
		c.clientConfig.HTTPClient = nil
	}

	err := utils.Deepcopy(c.clientConfig, &newConfig)
	if err != nil {
		log.Printf("[ERROR] Failed to deepcopy client configuration, using a shallow copy: %v", err)
		if httpClientCopy != http.DefaultClient {
			c.clientConfig.HTTPClient = httpClientCopy
		}
		return c.clientConfig
	}

	if httpClientCopy != http.DefaultClient {
		// Since http.Client contains either primitive, non-pointer types or interfaces, we can create a new
		// instance and copy the values directly.
		newHTTPClient := *httpClientCopy
		newConfig.HTTPClient = &newHTTPClient
	}

	return newConfig
}

// NewAPIClient creates a new API client with server overrides based on the provided location. The server overrides are
// retrieved from the configuration file and applied to a copy of the client configuration, ensuring that the original
// configuration stored in Config remains unchanged. If the IONOS_API_URL environment variable is set, it will take
// precedence over any configuration file overrides.
func (c Config) NewAPIClient(location string) *ionoscloud.APIClient {
	clientCfg := c.CopyClientConfig()

	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint %s from IONOS_API_URL env variable\n", os.Getenv(shared.IonosApiUrlEnvVar))
		return ionoscloud.NewAPIClient(&clientCfg)
	}

	if c.fileConfig == nil {
		return ionoscloud.NewAPIClient(&clientCfg)
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.Cloud, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s product in location %s and no global endpoints defined", fileconfiguration.Cloud, location)
		return ionoscloud.NewAPIClient(&clientCfg)
	}
	clientCfg.Servers = ionoscloud.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	clientCfg.HTTPClient = &http.Client{}
	clientCfg.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return ionoscloud.NewAPIClient(&clientCfg)
}

func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *ionoscloud.APIClient {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cloud)
	config := ionoscloud.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, ionoscloud.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)
	client := ionoscloud.NewAPIClient(config)
	return client
}
