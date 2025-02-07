package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	objstorage "github.com/ionos-cloud/sdk-go-object-storage"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	apiGatewayService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/apigateway"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var Version = "DEV"

type ClientOptions struct {
	Username         string
	Password         string
	Token            string
	Url              string
	Version          string
	TerraformVersion string
	Insecure         bool
}

// Provider returns a schema.Provider for ionoscloud
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(shared.IonosUsernameEnvVar, nil),
				Description: "IonosCloud username for API operations. If token is provided, token is preferred",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(shared.IonosPasswordEnvVar, nil),
				Description: "IonosCloud password for API operations. If token is provided, token is preferred",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(shared.IonosTokenEnvVar, nil),
				Description: "IonosCloud bearer token for API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(shared.IonosApiUrlEnvVar, ""),
				Description: "IonosCloud REST API URL. Usually not necessary to be set, SDKs know internally how to route requests to the API.",
			},
			"retries": {
				Type:       schema.TypeInt,
				Optional:   true,
				Default:    50,
				Deprecated: "Timeout is used instead of this functionality",
			},
			"contract_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "To be set only for reseller accounts. Allows to run terraform on a contract number under a reseller account.",
			},
			"s3_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_S3_ACCESS_KEY", nil),
				Description: "Access key for IONOS Object Storage operations.",
			},
			"s3_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_S3_SECRET_KEY", nil),
				Description: "Secret key for IONOS Object Storage operations.",
			},
			"s3_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "eu-central-3",
				DefaultFunc: schema.EnvDefaultFunc("IONOS_S3_REGION", nil),
				Description: "Region for IONOS Object Storage operations.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_ALLOW_INSECURE", nil),
				Description: "This field is to be set only for testing purposes. It is not recommended to set this field in production environments.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			constant.DatacenterResource:                        resourceDatacenter(),
			constant.IpBlockResource:                           resourceIPBlock(),
			constant.FirewallResource:                          resourceFirewall(),
			constant.LanResource:                               resourceLan(),
			constant.LoadBalancerResource:                      resourceLoadbalancer(),
			constant.NicResource:                               resourceNic(),
			constant.ServerResource:                            resourceServer(),
			constant.ServerCubeResource:                        resourceCubeServer(),
			constant.ServerVCPUResource:                        resourceVCPUServer(),
			constant.VolumeResource:                            resourceVolume(),
			constant.GroupResource:                             resourceGroup(),
			constant.ShareResource:                             resourceShare(),
			constant.UserResource:                              resourceUser(),
			constant.SnapshotResource:                          resourceSnapshot(),
			constant.ResourceIpFailover:                        resourceLanIPFailover(),
			constant.K8sClusterResource:                        resourcek8sCluster(),
			constant.K8sNodePoolResource:                       resourceK8sNodePool(),
			constant.PCCResource:                               resourcePrivateCrossConnect(),
			constant.BackupUnitResource:                        resourceBackupUnit(),
			constant.S3KeyResource:                             resourceS3Key(),
			constant.NatGatewayResource:                        resourceNatGateway(),
			constant.NatGatewayRuleResource:                    resourceNatGatewayRule(),
			constant.NetworkLoadBalancerResource:               resourceNetworkLoadBalancer(),
			constant.NetworkLoadBalancerForwardingRuleResource: resourceNetworkLoadBalancerForwardingRule(),
			constant.NSGResource:                               resourceNSG(),
			constant.NSGSelectionResource:                      resourceDatacenterNSGSelection(),
			constant.NSGFirewallRuleResource:                   resourceNSGFirewallRule(),
			constant.NFSClusterResource:                        resourceNFSCluster(),
			constant.NFSShareResource:                          resourceNFSShare(),
			constant.PsqlClusterResource:                       resourceDbaasPgSqlCluster(),
			constant.PsqlUserResource:                          resourceDbaasPgSqlUser(),
			constant.PsqlDatabaseResource:                      resourceDbaasPgSqlDatabase(),
			constant.DBaaSMariaDBClusterResource:               resourceDBaaSMariaDBCluster(),
			constant.DBaasMongoClusterResource:                 resourceDbaasMongoDBCluster(),
			constant.DBaasMongoUserResource:                    resourceDbaasMongoUser(),
			constant.DBaaSInMemoryDBReplicaSetResource:         resourceDBaaSInMemoryDBReplicaSet(),
			constant.ALBResource:                               resourceApplicationLoadBalancer(),
			constant.ALBForwardingRuleResource:                 resourceApplicationLoadBalancerForwardingRule(),
			constant.TargetGroupResource:                       resourceTargetGroup(),
			constant.CertificateResource:                       resourceCertificateManager(),
			constant.AutoCertificateProviderResource:           resourceCertificateManagerProvider(),
			constant.AutoCertificateResource:                   resourceCertificateManagerAutoCertificate(),
			constant.ContainerRegistryResource:                 resourceContainerRegistry(),
			constant.ContainerRegistryTokenResource:            resourceContainerRegistryToken(),
			constant.DataplatformClusterResource:               resourceDataplatformCluster(),
			constant.DataplatformNodePoolResource:              resourceDataplatformNodePool(),
			constant.DNSZoneResource:                           resourceDNSZone(),
			constant.DNSRecordResource:                         resourceDNSRecord(),
			constant.LoggingPipelineResource:                   resourceLoggingPipeline(),
			constant.AutoscalingGroupResource:                  ResourceAutoscalingGroup(),
			constant.ServerBootDeviceSelectionResource:         resourceServerBootDeviceSelection(),
			constant.KafkaClusterResource:                      resourceKafkaCluster(),
			constant.KafkaClusterTopicResource:                 resourceKafkaTopic(),
			constant.CDNDistributionResource:                   resourceCDNDistribution(),
			constant.APIGatewayResource:                        resourceAPIGateway(),
			constant.APIGatewayRouteResource:                   resourceAPIGatewayRoute(),
			constant.WireGuardGatewayResource:                  resourceVpnWireguardGateway(),
			constant.WireGuardPeerResource:                     resourceVpnWireguardPeer(),
			constant.IPSecGatewayResource:                      resourceVpnIPSecGateway(),
			constant.IPSecTunnelResource:                       resourceVpnIPSecTunnel(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			constant.DatacenterResource:                        dataSourceDataCenter(),
			constant.LocationResource:                          dataSourceLocation(),
			constant.ImageResource:                             dataSourceImage(),
			constant.ResourceResource:                          dataSourceResource(),
			constant.SnapshotResource:                          dataSourceSnapshot(),
			constant.LanResource:                               dataSourceLan(),
			constant.PCCResource:                               dataSourcePcc(),
			constant.ServerResource:                            dataSourceServer(),
			constant.ServerCubeResource:                        dataSourceCubeServer(),
			constant.ServerVCPUResource:                        dataSourceVCPUServer(),
			constant.ServersDataSource:                         dataSourceServers(),
			constant.K8sClusterResource:                        dataSourceK8sCluster(),
			constant.K8sClustersDataSource:                     dataSourceK8sClusters(),
			constant.K8sNodePoolResource:                       dataSourceK8sNodePool(),
			constant.K8sNodePoolNodesResource:                  dataSourceK8sNodePoolNodes(),
			constant.DBaaSMariaDBClusterResource:               dataSourceDBaaSMariaDBCluster(),
			constant.DBaaSMariaDBBackupsDataSource:             dataSourceDBaaSMariaDBBackups(),
			constant.NatGatewayResource:                        dataSourceNatGateway(),
			constant.NatGatewayRuleResource:                    dataSourceNatGatewayRule(),
			constant.NetworkLoadBalancerResource:               dataSourceNetworkLoadBalancer(),
			constant.NetworkLoadBalancerForwardingRuleResource: dataSourceNetworkLoadBalancerForwardingRule(),
			constant.NSGResource:                               dataSourceNSG(),
			constant.NFSClusterResource:                        dataSourceNFSCluster(),
			constant.NFSShareResource:                          dataSourceNFSShare(),
			constant.TemplateResource:                          dataSourceTemplate(),
			constant.BackupUnitResource:                        dataSourceBackupUnit(),
			constant.FirewallResource:                          dataSourceFirewall(),
			constant.S3KeyResource:                             dataSourceObjectStorageKey(),
			constant.GroupResource:                             dataSourceGroup(),
			constant.UserResource:                              dataSourceUser(),
			constant.IpBlockResource:                           dataSourceIpBlock(),
			constant.VolumeResource:                            dataSourceVolume(),
			constant.NicResource:                               dataSourceNIC(),
			constant.ShareResource:                             dataSourceShare(),
			constant.ResourceIpFailover:                        dataSourceIpFailover(),
			constant.PsqlClusterResource:                       dataSourceDbaasPgSqlCluster(),
			constant.PsqlUserResource:                          dataSourceDbaasPgSqlUser(),
			constant.PsqlDatabaseResource:                      dataSourceDbaasPgSqlDatabase(),
			constant.PsqlDatabasesResource:                     dataSourceDbaasPgSqlDatabases(),
			constant.DBaasMongoClusterResource:                 dataSourceDbaasMongoCluster(),
			constant.DBaaSMongoTemplateResource:                dataSourceDbassMongoTemplate(),
			constant.PsqlVersionsResource:                      dataSourceDbaasPgSqlVersions(),
			constant.PsqlBackupsResource:                       dataSourceDbaasPgSqlBackups(),
			constant.ALBResource:                               dataSourceApplicationLoadBalancer(),
			constant.ALBForwardingRuleResource:                 dataSourceApplicationLoadBalancerForwardingRule(),
			constant.TargetGroupResource:                       dataSourceTargetGroup(),
			constant.DBaasMongoUserResource:                    dataSourceDbaasMongoUser(),
			constant.DBaaSInMemoryDBReplicaSetResource:         dataSourceDBaaSInMemoryDBReplicaSet(),
			constant.DBaaSInMemoryDBSnapshotResource:           dataSourceDBaaSInMemoryDBSnapshot(),
			constant.CertificateResource:                       dataSourceCertificate(),
			constant.AutoCertificateProviderResource:           dataSourceCertificateManagerProvider(),
			constant.AutoCertificateResource:                   dataSourceCertificateManagerAutoCertificate(),
			constant.ContainerRegistryResource:                 dataSourceContainerRegistry(),
			constant.ContainerRegistryTokenResource:            dataSourceContainerRegistryToken(),
			constant.ContainerRegistryLocationsResource:        dataSourceContainerRegistryLocations(),
			constant.DataplatformClusterResource:               dataSourceDataplatformCluster(),
			constant.DataplatformNodePoolResource:              dataSourceDataplatformNodePool(),
			constant.DataplatformNodePoolsDataSource:           dataSourceDataplatformNodePools(),
			constant.DataplatformVersionsDataSource:            dataSourceDataplatformVersions(),
			constant.DNSZoneDataSource:                         dataSourceDNSZone(),
			constant.DNSRecordDataSource:                       dataSourceDNSRecord(),
			constant.LoggingPipelineDataSource:                 dataSourceLoggingPipeline(),
			constant.AutoscalingGroupResource:                  DataSourceAutoscalingGroup(),
			constant.AutoscalingGroupServersResource:           DataSourceAutoscalingGroupServers(),
			constant.KafkaClusterResource:                      dataSourceKafkaCluster(),
			constant.KafkaClusterTopicResource:                 dataSourceKafkaTopic(),
			constant.CDNDistributionResource:                   dataSourceCDNDistribution(),
			constant.APIGatewayResource:                        dataSourceAPIGateway(),
			constant.APIGatewayRouteResource:                   dataSourceAPIGatewayRoute(),
			constant.WireGuardGatewayResource:                  dataSourceVpnWireguardGateway(),
			constant.WireGuardPeerResource:                     dataSourceVpnWireguardPeer(),
			constant.IPSecGatewayResource:                      dataSourceVpnIPSecGateway(),
			constant.IPSecTunnelResource:                       dataSourceVpnIPSecTunnel(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it is 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		log.Printf("[DEBUG] Setting terraformVersion to %s", terraformVersion)

		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {

	var clientOpts ClientOptions
	username, usernameOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")
	token, tokenOk := d.GetOk("token")
	// for some reason, ENVDEFAULTFUNC does not work for this(boolean?) field
	if insecure := os.Getenv("IONOS_ALLOW_INSECURE"); insecure != "" {
		_ = d.Set("insecure", true)
	}
	insecure, insecureSet := d.GetOk("insecure")

	if !tokenOk {
		if !usernameOk {
			diags := diag.FromErr(fmt.Errorf("neither IonosCloud token, nor IonosCloud username has been provided"))
			return nil, diags
		}

		if !passwordOk {
			diags := diag.FromErr(fmt.Errorf("neither IonosCloud token, nor IonosCloud password has been provided"))
			return nil, diags
		}
	}

	cleanedURL := utils.CleanURL(d.Get("endpoint").(string))

	if contractNumber, contractOk := d.GetOk("contract_number"); contractOk {
		// will inject x-contract-number to sdks
		// TODO - Use constant here
		if err := os.Setenv("IONOS_CONTRACT_NUMBER", contractNumber.(string)); err != nil {
			return nil, diag.FromErr(err)
		}
	}

	// Standard client configuration
	clientOpts.Username = username.(string)
	clientOpts.Password = password.(string)
	clientOpts.Token = token.(string)
	clientOpts.Url = cleanedURL
	clientOpts.TerraformVersion = terraformVersion
	if insecureSet {
		clientOpts.Insecure = insecure.(bool)
	}

	return NewSDKBundleClient(clientOpts), nil
}

// NewSDKBundleClient returns a new SDK bundle client
func NewSDKBundleClient(clientOpts ClientOptions) interface{} {
	clientOpts.Version = ionoscloud.Version
	return services.SdkBundle{
		CDNClient:          NewClientByType(clientOpts, cdnClient).(*cdnService.Client),
		AutoscalingClient:  NewClientByType(clientOpts, autoscalingClient).(*autoscalingService.Client),
		CertManagerClient:  NewClientByType(clientOpts, certManagerClient).(*cert.Client),
		CloudApiClient:     NewClientByType(clientOpts, ionosClient).(*ionoscloud.APIClient),
		ContainerClient:    NewClientByType(clientOpts, containerRegistryClient).(*crService.Client),
		DataplatformClient: NewClientByType(clientOpts, dataplatformClient).(*dataplatformService.Client),
		DNSClient:          NewClientByType(clientOpts, dnsClient).(*dnsService.Client),
		LoggingClient:      NewClientByType(clientOpts, loggingClient).(*loggingService.Client),
		MariaDBClient:      NewClientByType(clientOpts, mariaDBClient).(*mariadb.MariaDBClient),
		MongoClient:        NewClientByType(clientOpts, mongoClient).(*dbaasService.MongoClient),
		NFSClient:          NewClientByType(clientOpts, nfsClient).(*nfsService.Client),
		PsqlClient:         NewClientByType(clientOpts, psqlClient).(*dbaasService.PsqlClient),
		KafkaClient:        NewClientByType(clientOpts, kafkaClient).(*kafkaService.Client),
		APIGatewayClient:   NewClientByType(clientOpts, apiGatewayClient).(*apiGatewayService.Client),
		VPNClient:          NewClientByType(clientOpts, vpnClient).(*vpn.Client),
		InMemoryDBClient:   NewClientByType(clientOpts, inMemoryDBClient).(*inmemorydb.InMemoryDBClient),
	}
}

type clientType int

const (
	ionosClient clientType = iota
	cdnClient
	autoscalingClient
	certManagerClient
	containerRegistryClient
	dataplatformClient
	dnsClient
	loggingClient
	mariaDBClient
	mongoClient
	nfsClient
	psqlClient
	s3Client
	kafkaClient
	apiGatewayClient
	vpnClient
	inMemoryDBClient
)

// NewClientByType returns a new client based on the client type
func NewClientByType(clientOpts ClientOptions, clientType clientType) interface{} {
	switch clientType {
	case ionosClient:
		{
			newConfig := shared.NewConfiguration(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url)
			newConfig.UserAgent = fmt.Sprintf(
				"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
				Version, ionoscloud.Version, clientOpts.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
			)
			newConfig.MaxRetries = constant.MaxRetries
			newConfig.WaitTime = constant.MaxWaitTime
			client := ionoscloud.NewAPIClient(newConfig)
			client.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOpts.Insecure)}
			return client
		}
	case cdnClient:
		return cdnService.NewCDNClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case autoscalingClient:
		return autoscalingService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case certManagerClient:
		return cert.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case containerRegistryClient:
		return crService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case dataplatformClient:
		return dataplatformService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case dnsClient:
		return dnsService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case loggingClient:
		return loggingService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.TerraformVersion, clientOpts.Insecure)
	case mariaDBClient:
		return mariadb.NewMariaDBClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case mongoClient:
		return dbaasService.NewMongoClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case nfsClient:
		return nfsService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case psqlClient:
		return dbaasService.NewPsqlClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case s3Client:
		{
			config := objstorage.NewConfiguration(clientOpts.Url)
			config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOpts.Insecure)}
			myS3Client := objstorage.NewAPIClient(config)
			return myS3Client
		}
	case kafkaClient:
		return kafkaService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.Username, clientOpts.Insecure)
	case apiGatewayClient:
		return apiGatewayService.NewClient(
			clientOpts.Username, clientOpts.Password, clientOpts.Token,
			clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion, clientOpts.Insecure)
	case vpnClient:
		return vpn.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Username, clientOpts.Insecure)
	case inMemoryDBClient:
		return inmemorydb.NewInMemoryDBClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.Username, clientOpts.Insecure)
	default:
		log.Fatalf("[ERROR] unknown client type %d", clientType)
	}
	return nil
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(utils.DefaultTimeout),
	Update:  schema.DefaultTimeout(utils.DefaultTimeout),
	Delete:  schema.DefaultTimeout(utils.DefaultTimeout),
	Default: schema.DefaultTimeout(utils.DefaultTimeout),
}
