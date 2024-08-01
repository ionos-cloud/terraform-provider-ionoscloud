package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	s3 "github.com/ionos-cloud/sdk-go-s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	apiGatewayService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/apigateway"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
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
}

// Provider returns a schema.Provider for ionoscloud
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ionoscloud.IonosUsernameEnvVar, nil),
				Description: "IonosCloud username for API operations. If token is provided, token is preferred",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ionoscloud.IonosPasswordEnvVar, nil),
				Description: "IonosCloud password for API operations. If token is provided, token is preferred",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ionoscloud.IonosTokenEnvVar, nil),
				Description: "IonosCloud bearer token for API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ionoscloud.IonosApiUrlEnvVar, ""),
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
				Description: "Access key for IONOS S3 operations.",
			},
			"s3_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_S3_SECRET_KEY", nil),
				Description: "Secret key for IONOS S3 operations.",
			},
			"s3_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "eu-central-3",
				DefaultFunc: schema.EnvDefaultFunc("IONOS_S3_REGION", nil),
				Description: "Region for IONOS S3 operations.",
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
			constant.NFSClusterResource:                        resourceNFSCluster(),
			constant.NFSShareResource:                          resourceNFSShare(),
			constant.PsqlClusterResource:                       resourceDbaasPgSqlCluster(),
			constant.PsqlUserResource:                          resourceDbaasPgSqlUser(),
			constant.PsqlDatabaseResource:                      resourceDbaasPgSqlDatabase(),
			constant.DBaaSMariaDBClusterResource:               resourceDBaaSMariaDBCluster(),
			constant.DBaasMongoClusterResource:                 resourceDbaasMongoDBCluster(),
			constant.DBaasMongoUserResource:                    resourceDbaasMongoUser(),
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
			constant.NFSClusterResource:                        dataSourceNFSCluster(),
			constant.NFSShareResource:                          dataSourceNFSShare(),
			constant.TemplateResource:                          dataSourceTemplate(),
			constant.BackupUnitResource:                        dataSourceBackupUnit(),
			constant.FirewallResource:                          dataSourceFirewall(),
			constant.S3KeyResource:                             dataSourceS3Key(),
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

	cleanedUrl := cleanURL(d.Get("endpoint").(string))

	if contractNumber, contractOk := d.GetOk("contract_number"); contractOk {
		// will inject x-contract-number to sdks
		if err := os.Setenv(ionoscloud.IonosContractNumber, contractNumber.(string)); err != nil {
			return nil, diag.FromErr(err)
		}
	}

	// Standard client configuration
	clientOpts.Username = username.(string)
	clientOpts.Password = password.(string)
	clientOpts.Token = token.(string)
	clientOpts.Url = cleanedUrl
	clientOpts.Version = ionoscloud.Version
	clientOpts.TerraformVersion = terraformVersion

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
		APIGatewayClient:   NewClientByType(clientOpts, apiGatewayClient).(*apiGatewayService.Client),
		VPNClient:          NewClientByType(clientOpts, vpnClient).(*vpn.Client),
	}, nil
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
	apiGatewayClient
	vpnClient
)

func NewClientByType(clientOpts ClientOptions, clientType clientType) interface{} {
	switch clientType {
	case ionosClient:
		{
			newConfig := ionoscloud.NewConfiguration(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url)
			newConfig.UserAgent = fmt.Sprintf(
				"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
				Version, ionoscloud.Version, clientOpts.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
			)
			if os.Getenv(constant.IonosDebug) != "" {
				newConfig.Debug = true
			}
			newConfig.MaxRetries = constant.MaxRetries
			newConfig.WaitTime = constant.MaxWaitTime
			newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
			return ionoscloud.NewAPIClient(newConfig)
		}
	case cdnClient:
		return cdnService.NewCDNClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case autoscalingClient:
		return autoscalingService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case certManagerClient:
		return cert.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case containerRegistryClient:
		return crService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case dataplatformClient:
		return dataplatformService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case dnsClient:
		return dnsService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case loggingClient:
		return loggingService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.TerraformVersion)
	case mariaDBClient:
		return mariadb.NewMariaDBClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case mongoClient:
		return dbaasService.NewMongoClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case nfsClient:
		return nfsService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case psqlClient:
		return dbaasService.NewPsqlClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.Username)
	case s3Client:
		return s3.NewAPIClient(s3.NewConfiguration())
	case apiGatewayClient:
		return apiGatewayService.NewClient(
			clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion,
		)
	case vpnClient:
		return vpn.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Username)
	default:
		log.Fatalf("[ERROR] unknown client type %d", clientType)
	}
	return nil
}

// cleanURL makes sure trailing slash does not corrupt the state
func cleanURL(url string) string {
	length := len(url)
	if length > 1 && url[length-1] == '/' {
		url = url[:length-1]
	}

	return url
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(utils.DefaultTimeout),
	Update:  schema.DefaultTimeout(utils.DefaultTimeout),
	Delete:  schema.DefaultTimeout(utils.DefaultTimeout),
	Default: schema.DefaultTimeout(utils.DefaultTimeout),
}
