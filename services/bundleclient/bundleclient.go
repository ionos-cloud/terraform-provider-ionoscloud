package bundleclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/failover"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

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
		DNSClient:                     dnsService.NewClient(clientOptions, fileConfig),
		LoggingClient:                 loggingService.NewClient(clientOptions, fileConfig),
		MariaDBClient:                 mariadb.NewClient(clientOptions, fileConfig),
		NFSClient:                     nfsService.NewClient(clientOptions, fileConfig),
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

// shouldApplyOverrides handles the early-exit checks common to all client constructors.
// It returns true when the caller should proceed with custom location or failover configuration,
// or false if the client should be returned immediately using the provided config (e.g. env var or default).
func (c SdkBundle) shouldApplyOverrides(product string) bool {
	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint from IONOS_API_URL env variable")
		return false
	}
	if c.fileConfig == nil {
		return false
	}
	if c.fileConfig.GetProductOverrides(product) == nil {
		log.Printf("[DEBUG] Missing config for %s product in file config, using SDK defaults", product)
		return false
	}
	return true
}

// NewContainerRegistryClient creates a new Container Registry client for a specific location
func (c SdkBundle) NewContainerRegistryClient(location string) (*crService.Client, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, cr.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(fileconfiguration.ContainerRegistry) {
		return crService.NewClientFromConfig(config), nil
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.ContainerRegistry, location)
	if endpoint == nil {
		if location == "" {
			return nil, fmt.Errorf(
				"could not instantiate Container Registry client: invalid config found for %q product in file config: "+
					"no global endpoints defined", fileconfiguration.ContainerRegistry,
			)
		}

		return nil, fmt.Errorf(
			"could not instantiate Container Registry client: invalid config found for %q product in file config: "+
				"missing endpoint in location %q and no global endpoints defined for fallback",
			fileconfiguration.ContainerRegistry, location,
		)
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return crService.NewClientFromConfig(config), nil
}

// NewMongoClient creates a new MongoDB client for a specific location
func (c SdkBundle) NewMongoClient(location string) (*dbaasService.MongoClient, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, mongo.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(fileconfiguration.Mongo) {
		return dbaasService.NewMongoClientFromConfig(config), nil
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.Mongo, location)
	if endpoint == nil {
		if location == "" {
			return nil, fmt.Errorf(
				"could not instantiate Mongo client: invalid config found for %q product in file config: "+
					"no global endpoints defined", fileconfiguration.Mongo,
			)
		}

		return nil, fmt.Errorf(
			"could not instantiate Mongo client: invalid config found for %q product in file config: "+
				"missing endpoint in location %q and no global endpoints defined for fallback",
			fileconfiguration.Mongo, location,
		)
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return dbaasService.NewMongoClientFromConfig(config), nil
}

// NewPsqlClient creates a new PostgreSQL client for a specific location
func (c SdkBundle) NewPsqlClient(location string) (*dbaasService.PsqlClient, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, psql.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(fileconfiguration.PSQL) {
		return dbaasService.NewPsqlClientFromConfig(config), nil
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.PSQL, location)
	if endpoint == nil {
		if location == "" {
			return nil, fmt.Errorf(
				"could not instantiate PostgreSQL client: invalid config found for %q product in file config: "+
					"no global endpoints defined", fileconfiguration.PSQL,
			)
		}

		return nil, fmt.Errorf(
			"could not instantiate PostgreSQL client: invalid config found for %q product in file config: "+
				"missing endpoint in location %q and no global endpoints defined for fallback",
			fileconfiguration.PSQL, location,
		)
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	return dbaasService.NewPsqlClientFromConfig(config), nil
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

// NewCloudAPIClient creates a new *ionoscloud.APIClient for the given location.
// The endpoint is determined in the following order of precedence:
//  1. IONOS_API_URL environment variable
//  2. If no file config is provided, use the default endpoint from the SDK
//  3. If file config is provided but no product overrides exist for the cloud product, use the default endpoint from the SDK
//  4. If file config is provided and product overrides for cloud are found:
//     a. If a location override is found for the provided location, use that endpoint
//     b. If no location override is found but a global override exists, use the global endpoint as fallback
//     c. If neither is found, return an error
func (c SdkBundle) NewCloudAPIClient(location string) (*ionoscloud.APIClient, error) {
	config := c.newCloudAPIClientConfig()
	if !c.shouldApplyOverrides(fileconfiguration.Cloud) {
		return ionoscloud.NewAPIClient(config), nil
	}

	endpoint := c.fileConfig.GetLocationOverridesWithGlobalFallback(fileconfiguration.Cloud, location)
	if endpoint == nil {
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
	config.HTTPClient = &http.Client{
		Transport: shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData),
	}
	return ionoscloud.NewAPIClient(config), nil
}

// NewCloudAPIClientWithFailover creates a new *ionoscloud.APIClient that distributes requests
// across all global endpoints configured for the cloud product using the failover strategy
// defined in the file config. It is intended for resources that do not have a location attribute.
// The endpoint is determined in the following order of precedence:
//  1. IONOS_API_URL environment variable
//  2. If no file config is provided, use the default endpoint from the SDK
//  3. If file config is provided but no product overrides exist for the cloud product, use the default endpoint from the SDK
//  4. If file config is provided and product overrides for cloud are found:
//     a. If no failover block is defined, or the strategy is "none", use the default endpoint from the SDK
//     b. If the strategy is "roundRobin", configure failover across all global endpoints
//     c. If no global endpoints are found, return an error
//     d. Any other strategy value is an error
func (c SdkBundle) NewCloudAPIClientWithFailover() (*ionoscloud.APIClient, error) {
	config := c.newCloudAPIClientConfig()
	if !c.shouldApplyOverrides(fileconfiguration.Cloud) {
		return ionoscloud.NewAPIClient(config), nil
	}

	failoverOptions := c.fileConfig.GetFailoverOptions()
	if failoverOptions == nil {
		return ionoscloud.NewAPIClient(config), nil
	}

	switch failover.NormalizeStrategy(failoverOptions.Strategy) {
	case failover.NormalizeStrategy(failover.RoundRobin):
		// handled below
	case failover.NormalizeStrategy(failover.None), "":
		return ionoscloud.NewAPIClient(config), nil
	default:
		return nil, fmt.Errorf("invalid failover strategy %q defined in file config, only %q is supported", failoverOptions.Strategy, failover.RoundRobin)
	}

	product := c.fileConfig.GetProductOverrides(fileconfiguration.Cloud)
	var failoverEndpoints []failover.Endpoint
	var servers ionoscloud.ServerConfigurations
	for _, ep := range product.Endpoints {
		if ep.Location != "" {
			continue
		}
		failoverEndpoints = append(failoverEndpoints, failover.Endpoint{
			URL:                 ep.Name,
			SkipTLSVerify:       ep.SkipTLSVerify,
			CertificateAuthData: ep.CertificateAuthData,
		})
		servers = append(servers, ionoscloud.ServerConfiguration{
			URL:         ep.Name,
			Description: shared.EndpointOverridden + "global",
		})
		log.Printf("[DEBUG] Adding global override endpoint %+v for %s product from file config", ep, fileconfiguration.Cloud)
	}
	if len(failoverEndpoints) == 0 {
		return nil, fmt.Errorf("no global failoverEndpoints configured for %q", fileconfiguration.Cloud)
	}

	config.Servers = servers
	config.HTTPClient.Transport = failover.NewRoundTripper(failoverEndpoints, *failoverOptions, config.HTTPClient.Transport)
	return ionoscloud.NewAPIClient(config), nil
}
