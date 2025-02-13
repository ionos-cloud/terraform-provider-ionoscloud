package ionoscloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"log"
	"os"

	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"

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
	usernameItf, usernameOk := d.GetOk("username")
	passwordItf, passwordOk := d.GetOk("password")
	tokenItf, tokenOk := d.GetOk("token")
	username := ""
	if usernameOk {
		username = usernameItf.(string)
	}
	password := ""
	if passwordOk {
		password = passwordItf.(string)
	}
	token := ""
	if tokenOk {
		token = tokenItf.(string)
	}
	// for some reason, ENVDEFAULTFUNC does not work for this(boolean?) field
	if insecure := os.Getenv("IONOS_ALLOW_INSECURE"); insecure != "" {
		_ = d.Set("insecure", true)
	}
	insecure, insecureSet := d.GetOk("insecure")
	insecureBool := false

	loadedConfig, readFileErr := shared.ReadConfigFromFile()
	if !tokenOk {
		if !usernameOk || !passwordOk {
			if readFileErr != nil {
				return nil, diag.Errorf("missing credentials, either token or username and password must be set, %s", readFileErr.Error())
			}
			profile := loadedConfig.GetCurrentProfile()
			if profile == nil {
				return nil, diag.Errorf("missing credentials, either token or username and password must be set")
			}
			token = profile.Credentials.Token
			username = profile.Credentials.Username
			password = profile.Credentials.Password
		}
		if token == "" && (username == "" || password == "") {
			return nil, diag.Errorf("missing credentials, either token or username and password must be set")
		}
	}

	endpoint := utils.CleanURL(d.Get("endpoint").(string))

	if contractNumber, contractOk := d.GetOk("contract_number"); contractOk {
		// will inject x-contract-number to sdks
		if err := os.Setenv(ionoscloud.IonosContractNumber, contractNumber.(string)); err != nil {
			return nil, diag.FromErr(err)
		}
	}

	if insecureSet {
		insecureBool = insecure.(bool)
	}
	clientOptions := bundle.ClientOptions{
		ClientOverrideOptions: shared.ClientOverrideOptions{
			Endpoint:      endpoint,
			SkipTLSVerify: insecureBool,
			//Certificate:   "",
			Credentials: shared.Credentials{
				Username: username,
				Password: password,
				Token:    token,
			},
		},
		Version:          "",
		TerraformVersion: terraformVersion,
	}

	client := services.SdkBundle{
		CDNClient:                     cdnService.NewCDNClient(clientOptions, loadedConfig),
		AutoscalingClient:             autoscalingService.NewClient(clientOptions, loadedConfig),
		CertManagerClient:             cert.NewClient(clientOptions, loadedConfig),
		CloudApiClient:                cloudapi.NewClient(clientOptions, loadedConfig),
		ContainerClient:               crService.NewClient(clientOptions, loadedConfig),
		DataplatformClient:            dataplatformService.NewClient(clientOptions, loadedConfig),
		DNSClient:                     dnsService.NewClient(clientOptions, loadedConfig),
		LoggingClient:                 loggingService.NewClient(clientOptions, loadedConfig),
		MariaDBClient:                 mariadb.NewMariaDBClient(clientOptions, loadedConfig),
		MongoClient:                   dbaasService.NewMongoClient(clientOptions, loadedConfig),
		NFSClient:                     nfsService.NewClient(clientOptions, loadedConfig),
		PsqlClient:                    dbaasService.NewPSQLClient(clientOptions, loadedConfig),
		KafkaClient:                   kafkaService.NewClient(clientOptions, loadedConfig),
		APIGatewayClient:              apiGatewayService.NewClient(clientOptions, loadedConfig),
		VPNClient:                     vpn.NewClient(clientOptions, loadedConfig),
		InMemoryDBClient:              inmemorydb.NewInMemoryDBClient(clientOptions, loadedConfig),
		ObjectStorageManagementClient: objectStorageManagementService.NewClient(clientOptions, loadedConfig),
		MonitoringClient:              monitoringService.NewClient(clientOptions, loadedConfig),
	}

	return client, nil
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(utils.DefaultTimeout),
	Update:  schema.DefaultTimeout(utils.DefaultTimeout),
	Delete:  schema.DefaultTimeout(utils.DefaultTimeout),
	Default: schema.DefaultTimeout(utils.DefaultTimeout),
}
