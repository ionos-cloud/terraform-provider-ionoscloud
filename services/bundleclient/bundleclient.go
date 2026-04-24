package bundleclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	pgsqlv2sdk "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"
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
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	objectStorageService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/configlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// New creates a new SdkBundle client
func New(ctx context.Context, clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *SdkBundle {
	return &SdkBundle{
		CDNClient:         cdnService.NewClient(ctx, clientOptions, fileConfig),
		AutoscalingClient: autoscalingService.NewClient(ctx, clientOptions, fileConfig),
		CertManagerClient: cert.NewClient(ctx, clientOptions, fileConfig),
		DNSClient:         dnsService.NewClient(ctx, clientOptions, fileConfig),
		LoggingClient:     loggingService.NewClient(clientOptions, fileConfig),
		MariaDBClient:     mariadb.NewClient(clientOptions, fileConfig),
		NFSClient:         nfsService.NewClient(clientOptions, fileConfig),
		KafkaClient:       kafkaService.NewClient(clientOptions, fileConfig),
		VPNClient:         vpn.NewClient(clientOptions, fileConfig),
		InMemoryDBClient:  inmemorydb.NewClient(clientOptions, fileConfig),
		S3Client:          objectStorageService.NewClient(ctx, clientOptions, fileConfig),
		MonitoringClient:  monitoringService.NewClient(clientOptions, fileConfig),

		clientOptions: clientOptions,
		fileConfig:    fileConfig,
	}
}

// SdkBundle is a struct that defines the bundle client. It is used for both sdkv2 and plugin framework
type SdkBundle struct {
	InMemoryDBClient  *inmemorydb.Client
	PsqlClient        *dbaasService.PsqlClient
	MongoClient       *dbaasService.MongoClient
	MariaDBClient     *mariadb.Client
	NFSClient         *nfsService.Client
	CertManagerClient *cert.Client
	DNSClient         *dnsService.Client
	LoggingClient     *loggingService.Client
	AutoscalingClient *autoscalingService.Client
	KafkaClient       *kafkaService.Client
	CDNClient         *cdnService.Client
	VPNClient         *vpn.Client
	S3Client          *objectStorageService.Client
	MonitoringClient  *monitoringService.Client

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

// shouldApplyOverridesCustomEnv handles the early-exit checks common to all client constructors.
// It returns true when the caller should proceed with custom location or failover configuration,
// or false if the client should be returned immediately using the provided config (e.g. env var provided as parameter or default).
func (c SdkBundle) shouldApplyOverridesCustomEnv(ctx context.Context, product, productEnvVar string) bool {
	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		tflog.Debug(ctx, "using endpoint from env", map[string]interface{}{"product": product, "env": shared.IonosApiUrlEnvVar, "url": os.Getenv(shared.IonosApiUrlEnvVar)})
		return false
	}
	if productEnvVar != "" && os.Getenv(productEnvVar) != "" {
		tflog.Debug(ctx, "using endpoint from env", map[string]interface{}{"product": product, "env": productEnvVar, "url": os.Getenv(productEnvVar)})
		return false
	}
	if c.fileConfig == nil {
		tflog.Debug(ctx, "no file config available", map[string]interface{}{"product": product})
		return false
	}
	if c.fileConfig.GetProductOverrides(product) == nil {
		tflog.Debug(ctx, "no config in file config, using SDK defaults", map[string]interface{}{"product": product})
		return false
	}
	return true
}

// shouldApplyOverrides handles the early-exit checks common to all client constructors.
// It returns true when the caller should proceed with custom location or failover configuration,
// or false if the client should be returned immediately using the provided config (e.g. IONOS_API_URL env var or default).
func (c SdkBundle) shouldApplyOverrides(ctx context.Context, product string) bool {
	return c.shouldApplyOverridesCustomEnv(ctx, product, "")
}

// NewContainerRegistryClient creates a new Container Registry client for a specific location
func (c SdkBundle) NewContainerRegistryClient(ctx context.Context, location string) (*crService.Client, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, cr.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(ctx, fileconfiguration.ContainerRegistry) {
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
	tflog.Debug(ctx, "Container Registry: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": endpoint.Name})
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
func (c SdkBundle) NewMongoClient(ctx context.Context, location string) (*dbaasService.MongoClient, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, mongo.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(ctx, fileconfiguration.Mongo) {
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
	tflog.Debug(ctx, "Mongo: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": endpoint.Name})
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
func (c SdkBundle) NewPsqlClient(ctx context.Context, location string) (*dbaasService.PsqlClient, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, psql.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverrides(ctx, fileconfiguration.PSQL) {
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
	tflog.Debug(ctx, "PostgreSQL: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": endpoint.Name})
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

// NewPgSQLV2Client creates a new PostgreSQL v2 client for a specific location.
func (c SdkBundle) NewPgSQLV2Client(ctx context.Context, location string) (*pgsqlv2Service.Client, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-pgsqlv2/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, pgsqlv2sdk.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if c.fileConfig != nil {
		endpoint := c.fileConfig.GetProductLocationOverrides(fileconfiguration.PSQLV2, location)
		if endpoint == nil {
			tflog.Warn(ctx, "product missing from config file or location not defined, using internal locations map to configure endpoint", map[string]interface{}{"product": fileconfiguration.PSQLV2, "location": location})
		} else {
			config.Servers = shared.ServerConfigurations{
				{
					URL:         endpoint.Name,
					Description: shared.EndpointOverridden + location,
				},
			}
			config.HTTPClient = &http.Client{}
			config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
			return pgsqlv2Service.NewClientFromConfig(config), nil
		}
	}

	endpoint := pgsqlv2Service.LocationToURL[location]
	if endpoint == "" {
		return nil, fmt.Errorf("can't configure endpoint for location %q, available locations: %s", location, strings.Join(pgsqlv2Service.AvailableLocations(), ", "))
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint,
			Description: "endpoint from the internal locations map, location: " + location,
		},
	}

	return pgsqlv2Service.NewClientFromConfig(config), nil
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
func (c SdkBundle) NewCloudAPIClient(ctx context.Context, location string) (*ionoscloud.APIClient, error) {
	config := c.newCloudAPIClientConfig()
	if !c.shouldApplyOverrides(ctx, fileconfiguration.Cloud) {
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
	tflog.Debug(ctx, "Cloud API: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": endpoint.Name})
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
//     a. If no failover block is defined, or the strategy is "none", use first endpoint found in file configuration
//     b. If the strategy is "roundRobin", configure failover across all global endpoints in the file configuration
//     c. If no global endpoints are found, return an error
//     d. Any other strategy value is an error
func (c SdkBundle) NewCloudAPIClientWithFailover(ctx context.Context) (*ionoscloud.APIClient, error) {
	config := c.newCloudAPIClientConfig()
	if !c.shouldApplyOverrides(ctx, fileconfiguration.Cloud) {
		return ionoscloud.NewAPIClient(config), nil
	}

	failoverOptions := c.fileConfig.GetFailoverOptions()
	if failoverOptions == nil {
		failoverOptions = &failover.Options{Strategy: failover.None}
	}

	endpoints := c.fileConfig.FilterGlobalOverrides(fileconfiguration.Cloud)
	var failoverEndpoints []failover.Endpoint
	var servers ionoscloud.ServerConfigurations
	for _, ep := range endpoints {
		failoverEndpoints = append(failoverEndpoints, failover.Endpoint{
			URL:                 ep.Name,
			SkipTLSVerify:       ep.SkipTLSVerify,
			CertificateAuthData: ep.CertificateAuthData,
		})
		servers = append(servers, ionoscloud.ServerConfiguration{
			URL:         ep.Name,
			Description: shared.EndpointOverridden + "global",
		})
		tflog.Debug(ctx, "adding global override endpoint from file config", map[string]interface{}{"url": ep.Name, "skip_tls_verify": ep.SkipTLSVerify, "product": fileconfiguration.Cloud})
		if ep.CertificateAuthData != "" {
			tflog.Debug(ctx, "certificateAuthData present", map[string]interface{}{"product": fileconfiguration.Cloud, "cert_auth_data_len": len(ep.CertificateAuthData), "url": ep.Name})
		}
	}
	tflog.Debug(ctx, "failover config", map[string]interface{}{"product": fileconfiguration.Cloud, "strategy": failoverOptions.Strategy, "endpoints": len(failoverEndpoints)})
	if len(failoverEndpoints) == 0 {
		return nil, fmt.Errorf("no global failover endpoints configured for %q", fileconfiguration.Cloud)
	}
	//nolint:exhaustive
	switch failover.NormalizeStrategy(failoverOptions.Strategy) {
	case failover.NormalizeStrategy(failover.RoundRobin):
		config.HTTPClient.Transport = failover.NewRoundTripper(failoverEndpoints, *failoverOptions, config.HTTPClient.Transport)
	case failover.NormalizeStrategy(failover.None), "":
		servers = servers[0:1]
		ep := failoverEndpoints[0]
		if ep.SkipTLSVerify || ep.CertificateAuthData != "" {
			config.HTTPClient.Transport = shared.CreateTransport(ep.SkipTLSVerify, ep.CertificateAuthData)
		}
	default:
		return nil, fmt.Errorf("invalid failover strategy %q defined in file config, only %q, %q or an empty value are supported",
			failoverOptions.Strategy, failover.RoundRobin, failover.None)
	}
	config.Servers = servers
	return ionoscloud.NewAPIClient(config), nil
}

// NewObjectStorageManagementClient creates a new Object Storage Management client that distributes requests
// across all global endpoints configured for this product using the failover strategy
// defined in the file config.
func (c SdkBundle) NewObjectStorageManagementClient(ctx context.Context) (*objectStorageManagementService.Client, error) {
	config := c.newBundleClientConfig(fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-object-storage-management/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		c.clientOptions.Version, objectstoragemanagement.Version, c.clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	))

	if !c.shouldApplyOverridesCustomEnv(ctx, fileconfiguration.ObjectStorageManagement, objectStorageManagementService.IonosAPIURLObjectStorageManagement) {
		return objectStorageManagementService.NewClientFromConfig(config), nil
	}

	failoverOptions := c.fileConfig.GetFailoverOptions()
	if failoverOptions == nil {
		failoverOptions = &failover.Options{Strategy: failover.None}
	}

	endpoints := c.fileConfig.FilterGlobalOverrides(fileconfiguration.ObjectStorageManagement)
	var failoverEndpoints []failover.Endpoint
	var servers shared.ServerConfigurations

	for _, ep := range endpoints {
		failoverEndpoints = append(failoverEndpoints, failover.Endpoint{
			URL:                 ep.Name,
			SkipTLSVerify:       ep.SkipTLSVerify,
			CertificateAuthData: ep.CertificateAuthData,
		})
		servers = append(servers, shared.ServerConfiguration{
			URL:         ep.Name,
			Description: shared.EndpointOverridden + "global",
		})
		tflog.Debug(ctx, "adding global override endpoint from file config", map[string]interface{}{"url": ep.Name, "skip_tls_verify": ep.SkipTLSVerify, "product": fileconfiguration.ObjectStorageManagement})
		if ep.CertificateAuthData != "" {
			tflog.Debug(ctx, "certificateAuthData present", map[string]interface{}{"product": fileconfiguration.ObjectStorageManagement, "cert_auth_data_len": len(ep.CertificateAuthData), "url": ep.Name})
		}
	}
	tflog.Debug(ctx, "failover config", map[string]interface{}{"product": fileconfiguration.ObjectStorageManagement, "strategy": failoverOptions.Strategy, "endpoints": len(failoverEndpoints)})

	if len(failoverEndpoints) == 0 {
		return nil, fmt.Errorf("no global failover endpoints configured for %q", fileconfiguration.ObjectStorageManagement)
	}

	//nolint:exhaustive
	switch failover.NormalizeStrategy(failoverOptions.Strategy) {

	case failover.NormalizeStrategy(failover.RoundRobin):
		config.HTTPClient.Transport = failover.NewRoundTripper(failoverEndpoints, *failoverOptions, config.HTTPClient.Transport)

	case failover.NormalizeStrategy(failover.None), "":
		servers = servers[0:1]
		ep := failoverEndpoints[0]
		if ep.SkipTLSVerify || ep.CertificateAuthData != "" {
			config.HTTPClient.Transport = shared.CreateTransport(ep.SkipTLSVerify, ep.CertificateAuthData)
		}

	default:
		return nil, fmt.Errorf("invalid failover strategy %q defined in file config, only %q, %q or an empty value are supported",
			failoverOptions.Strategy, failover.RoundRobin, failover.None)
	}

	config.Servers = servers

	return objectStorageManagementService.NewClientFromConfig(config), nil
}
