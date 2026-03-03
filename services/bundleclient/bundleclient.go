package bundleclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
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

// prepareCloudAPIClient handles the early-exit checks common to all Cloud API client constructors.
// It returns nil when the caller should proceed with custom endpoint configuration,
// or client when a default client should be returned immediately.
// Exactly one of the two return values is always nil.
func (c SdkBundle) prepareCloudAPIClient(config *ionoscloud.Configuration) *ionoscloud.APIClient {

	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		log.Printf("[DEBUG] Using custom endpoint from IONOS_API_URL env variable")
		return ionoscloud.NewAPIClient(config)
	}
	if c.fileConfig == nil {
		return ionoscloud.NewAPIClient(config)
	}
	if c.fileConfig.GetProductOverrides(fileconfiguration.Cloud) == nil {
		log.Printf("[DEBUG] Missing config for %s product in file config, using SDK defaults", fileconfiguration.Cloud)
		return ionoscloud.NewAPIClient(config)
	}
	return nil
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
	client := c.prepareCloudAPIClient(config)
	if client != nil {
		return client, nil
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
	client := c.prepareCloudAPIClient(config)
	if client != nil {
		return client, nil
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
	var endpoints []failover.FailoverEndpoint
	var servers ionoscloud.ServerConfigurations
	for _, ep := range product.Endpoints {
		if ep.Location != "" {
			continue
		}
		endpoints = append(endpoints, failover.FailoverEndpoint{
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
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no global endpoints configured for %q", fileconfiguration.Cloud)
	}

	config.Servers = servers
	config.HTTPClient.Transport = failover.NewRoundTripper(endpoints, *failoverOptions, config.HTTPClient.Transport)
	return ionoscloud.NewAPIClient(config), nil
}
