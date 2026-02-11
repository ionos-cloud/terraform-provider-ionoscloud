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
	clientConfig *ionoscloud.Configuration
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
		clientConfig: config,
	}
}

// CopyClientConfig creates a deep copy of the client configuration to ensure that modifications to the returned configuration
// do not affect the original configuration stored in Config.
func (c *Config) CopyClientConfig() *ionoscloud.Configuration {
	if c == nil {
		return nil
	}

	if c.clientConfig == nil {
		return nil
	}

	// Copy directly what can be copied with a simple assignment
	newConfig := *c.clientConfig

	// If deepcopy fails, log as a warning and default to nil for that field to avoid using shallow-copied data.
	if err := utils.Deepcopy(c.clientConfig.DefaultHeader, &newConfig.DefaultHeader); err != nil {
		log.Printf("[WARN] Failed to deepcopy DefaultHeader configuration, will default to nil: %v", err)
		newConfig.DefaultHeader = nil
	}

	if err := utils.Deepcopy(c.clientConfig.DefaultQueryParams, &newConfig.DefaultQueryParams); err != nil {
		log.Printf("[WARN] Failed to deepcopy DefaultQueryParams configuration, will default to nil: %v", err)
		newConfig.DefaultQueryParams = nil
	}

	if err := utils.Deepcopy(c.clientConfig.Servers, &newConfig.Servers); err != nil {
		log.Printf("[WARN] Failed to deepcopy Servers configuration, will default to nil: %v", err)
		newConfig.Servers = nil
	}

	if err := utils.Deepcopy(c.clientConfig.OperationServers, &newConfig.OperationServers); err != nil {
		log.Printf("[WARN] Failed to deepcopy OperationServers configuration, will default to nil: %v", err)
		newConfig.OperationServers = nil
	}

	if c.clientConfig.HTTPClient != nil {
		// Since http.Client contains either primitive, non-pointer types or interfaces, we can create a new
		// instance and copy the values directly.
		newHTTPClient := *c.clientConfig.HTTPClient
		newConfig.HTTPClient = &newHTTPClient
	}

	return &newConfig
}

// NewAPIClient create a new API client with server overrides based on the provided location. The server overrides are
// retrieved from the configuration file and applied to a copy of the client configuration, ensuring that the original
// configuration stored in Config remains unchanged. If the IONOS_API_URL environment variable is set, it will take
// precedence over any configuration file overrides.
func (c *Config) NewAPIClient(location string) *ionoscloud.APIClient {
	if c == nil {
		return nil
	}

	clientCfg := c.CopyClientConfig()
	if clientCfg == nil {
		return nil
	}

	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint %s from IONOS_API_URL env variable\n", os.Getenv(shared.IonosApiUrlEnvVar))
		return ionoscloud.NewAPIClient(clientCfg)
	}

	if c.fileConfig == nil {
		return ionoscloud.NewAPIClient(clientCfg)
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.Cloud, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s in location %s and no global endpoints defined", fileconfiguration.Cloud, location)
		return ionoscloud.NewAPIClient(clientCfg)
	}
	clientCfg.Servers = ionoscloud.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	clientCfg.HTTPClient = &http.Client{}
	clientCfg.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return ionoscloud.NewAPIClient(clientCfg)
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
