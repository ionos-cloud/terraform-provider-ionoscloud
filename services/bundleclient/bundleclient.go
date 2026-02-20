package bundleclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	objectStorageService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// New creates a new SdkBundle client
func New(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *SdkBundle {
	return &SdkBundle{
		CDNClient:                     cdnService.NewClient(clientOptions, fileConfig),
		AutoscalingClient:             autoscalingService.NewClient(clientOptions, fileConfig),
		CertManagerClient:             cert.NewClient(clientOptions, fileConfig),
		ContainerClient:               crService.NewClient(clientOptions, fileConfig),
		DNSClient:                     dnsService.NewClient(clientOptions, fileConfig),
		LoggingClient:                 loggingService.NewClient(clientOptions, fileConfig),
		MariaDBClient:                 mariadb.NewClient(clientOptions, fileConfig),
		MongoClient:                   dbaasService.NewMongoClient(clientOptions, fileConfig),
		NFSClient:                     nfsService.NewClient(clientOptions, fileConfig),
		PsqlClient:                    dbaasService.NewPSQLClient(clientOptions, fileConfig),
		KafkaClient:                   kafkaService.NewClient(clientOptions, fileConfig),
		VPNClient:                     vpn.NewClient(clientOptions, fileConfig),
		InMemoryDBClient:              inmemorydb.NewClient(clientOptions, fileConfig),
		S3Client:                      objectStorageService.NewClient(clientOptions, fileConfig),
		ObjectStorageManagementClient: objectStorageManagementService.NewClient(clientOptions, fileConfig),
		MonitoringClient:              monitoringService.NewClient(clientOptions, fileConfig),

		clientOptions: clientOptions,
		fileConfig:    fileConfig,
	}
}

// SdkBundle is a struct that defines the bundle client. It is used for both sdkv2 and plugin framework
type SdkBundle struct {
	InMemoryDBClient              *inmemorydb.Client
	PsqlClient                    *dbaasService.PsqlClient
	MongoClient                   *dbaasService.MongoClient
	MariaDBClient                 *mariadb.Client
	NFSClient                     *nfsService.Client
	CertManagerClient             *cert.Client
	ContainerClient               *crService.Client
	DNSClient                     *dnsService.Client
	LoggingClient                 *loggingService.Client
	AutoscalingClient             *autoscalingService.Client
	KafkaClient                   *kafkaService.Client
	CDNClient                     *cdnService.Client
	VPNClient                     *vpn.Client
	S3Client                      *objectStorageService.Client
	ObjectStorageManagementClient *objectStorageManagementService.Client
	MonitoringClient              *monitoringService.Client

	clientOptions clientoptions.TerraformClientOptions
	fileConfig    *fileconfiguration.FileConfig
}

func (c SdkBundle) newBundleClientConfig(userAgent string) *shared.Configuration {
	config := shared.NewConfiguration(c.clientOptions.Credentials.Username, c.clientOptions.Credentials.Password, c.clientOptions.Credentials.Token, c.clientOptions.Endpoint)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = userAgent
	config.HTTPClient = &http.Client{Transport: shared.CreateTransport(c.clientOptions.SkipTLSVerify, c.clientOptions.Certificate)}
	return config
}

// NewContainerRegistryClient creates a new Container Registry client for a specific location
func (c SdkBundle) NewContainerRegistryClient(location string) *crService.Client {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, cr.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint %s from IONOS_API_URL env variable\n", os.Getenv(shared.IonosApiUrlEnvVar))
		return crService.NewClientFromConfig(config)
	}

	if c.fileConfig == nil {
		return crService.NewClientFromConfig(config)
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.ContainerRegistry, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s product in location %s and no global endpoints defined", fileconfiguration.ContainerRegistry, location)
		return crService.NewClientFromConfig(config)
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return crService.NewClientFromConfig(config)
}

// newCloudAPIClientConfig creates a new *ionoscloud.Configuration using the client options defined in the SdkBundle struct.
func (c SdkBundle) newCloudAPIClientConfig() *ionoscloud.Configuration {
	config := ionoscloud.NewConfiguration(
		c.clientOptions.Credentials.Username, c.clientOptions.Credentials.Password, c.clientOptions.Credentials.Token, c.clientOptions.Endpoint,
	)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, ionoscloud.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(c.clientOptions.SkipTLSVerify, c.clientOptions.Certificate)
	return config
}

// NewCloudAPIClient creates a new *ionoscloud.APIClient using the client options and file config defined in the SdkBundle struct.
// The endpoint is determined in the following order of precedence:
//  1. IONOS_API_URL environment variable
//  2. If no file config is provided, use the default endpoint from the SDK
//  3. If file config is provided but no product overrides exist for cloud product, use the default endpoint from the SDK
//  4. If file config is provided and product overrides for cloud are found:
//     a. If location override is found for the provided location, use the endpoint from the location override
//     b. If no location override is found but a global override is found, use the endpoint from the global override
//     c. If no location or global override is found, log a warning for invalid configuration and return nil with an error
func (c SdkBundle) NewCloudAPIClient(location string) (*ionoscloud.APIClient, error) {
	config := c.newCloudAPIClientConfig()

	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint from IONOS_API_URL env variable")
		return ionoscloud.NewAPIClient(config), nil
	}

	if c.fileConfig == nil {
		return ionoscloud.NewAPIClient(config), nil
	}

	if c.fileConfig.GetProductOverrides(fileconfiguration.Cloud) == nil {
		log.Printf("[DEBUG] Missing config for %s product in file config, using SDK defaults", fileconfiguration.Cloud)
		return ionoscloud.NewAPIClient(config), nil
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.Cloud, location)
	if endpoint == nil {
		if location == "" {
			return nil, fmt.Errorf(
				"could not instantiate Cloud API client: invalid config found for %q product in file config: "+
					"no global endpoints defined", fileconfiguration.Cloud,
			)
		}

		return nil, fmt.Errorf(
			"could not instantiate Cloud API client: invalid config found for %q product in file config: "+
				"missing endpoint in location %q and no global endpoints defined for fallback",
			fileconfiguration.Cloud, location,
		)
	}
	config.Servers = ionoscloud.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return ionoscloud.NewAPIClient(config), nil
}
