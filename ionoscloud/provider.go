package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var Version = "DEV"

type SdkBundle struct {
	CloudApiClient     *ionoscloud.APIClient
	PsqlClient         *dbaasService.PsqlClient
	MongoClient        *dbaasService.MongoClient
	CertManagerClient  *cert.Client
	ContainerClient    *crService.Client
	DataplatformClient *dataplatformService.Client
	DNSClient          *dnsService.Client
	LoggingClient      *loggingService.Client
}

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
		},
		ResourcesMap: map[string]*schema.Resource{
			DatacenterResource:          resourceDatacenter(),
			IpBlockResource:             resourceIPBlock(),
			FirewallResource:            resourceFirewall(),
			LanResource:                 resourceLan(),
			"ionoscloud_loadbalancer":   resourceLoadbalancer(),
			NicResource:                 resourceNic(),
			ServerResource:              resourceServer(),
			ServerCubeResource:          resourceCubeServer(),
			VolumeResource:              resourceVolume(),
			GroupResource:               resourceGroup(),
			ShareResource:               resourceShare(),
			UserResource:                resourceUser(),
			SnapshotResource:            resourceSnapshot(),
			ResourceIpFailover:          resourceLanIPFailover(),
			K8sClusterResource:          resourcek8sCluster(),
			K8sNodePoolResource:         resourceK8sNodePool(),
			PCCResource:                 resourcePrivateCrossConnect(),
			BackupUnitResource:          resourceBackupUnit(),
			S3KeyResource:               resourceS3Key(),
			NatGatewayResource:          resourceNatGateway(),
			NatGatewayRuleResource:      resourceNatGatewayRule(),
			NetworkLoadBalancerResource: resourceNetworkLoadBalancer(),
			NetworkLoadBalancerForwardingRuleResource: resourceNetworkLoadBalancerForwardingRule(),
			PsqlClusterResource:                       resourceDbaasPgSqlCluster(),
			DBaasMongoClusterResource:                 resourceDbaasMongoDBCluster(),
			DBaasMongoUserResource:                    resourceDbaasMongoUser(),
			ALBResource:                               resourceApplicationLoadBalancer(),
			ALBForwardingRuleResource:                 resourceApplicationLoadBalancerForwardingRule(),
			TargetGroupResource:                       resourceTargetGroup(),
			CertificateResource:                       resourceCertificateManager(),
			ContainerRegistryResource:                 resourceContainerRegistry(),
			ContainerRegistryTokenResource:            resourceContainerRegistryToken(),
			DataplatformClusterResource:               resourceDataplatformCluster(),
			DataplatformNodePoolResource:              resourceDataplatformNodePool(),
			DNSZoneResource:                           resourceDNSZone(),
			DNSRecordResource:                         resourceDNSRecord(),
			LoggingPipelineResource:                   resourceLoggingPipeline(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			DatacenterResource:                        dataSourceDataCenter(),
			LocationResource:                          dataSourceLocation(),
			ImageResource:                             dataSourceImage(),
			ResourceResource:                          dataSourceResource(),
			SnapshotResource:                          dataSourceSnapshot(),
			LanResource:                               dataSourceLan(),
			PCCResource:                               dataSourcePcc(),
			ServerResource:                            dataSourceServer(),
			ServerCubeResource:                        dataSourceCubeServer(),
			ServersDataSource:                         dataSourceServers(),
			K8sClusterResource:                        dataSourceK8sCluster(),
			K8sNodePoolResource:                       dataSourceK8sNodePool(),
			K8sNodePoolNodesResource:                  dataSourceK8sNodePoolNodes(),
			NatGatewayResource:                        dataSourceNatGateway(),
			NatGatewayRuleResource:                    dataSourceNatGatewayRule(),
			NetworkLoadBalancerResource:               dataSourceNetworkLoadBalancer(),
			NetworkLoadBalancerForwardingRuleResource: dataSourceNetworkLoadBalancerForwardingRule(),
			TemplateResource:                          dataSourceTemplate(),
			BackupUnitResource:                        dataSourceBackupUnit(),
			FirewallResource:                          dataSourceFirewall(),
			S3KeyResource:                             dataSourceS3Key(),
			GroupResource:                             dataSourceGroup(),
			UserResource:                              dataSourceUser(),
			IpBlockResource:                           dataSourceIpBlock(),
			VolumeResource:                            dataSourceVolume(),
			NicResource:                               dataSourceNIC(),
			ShareResource:                             dataSourceShare(),
			ResourceIpFailover:                        dataSourceIpFailover(),
			PsqlClusterResource:                       dataSourceDbaasPgSqlCluster(),
			DBaasMongoClusterResource:                 dataSourceDbaasMongoCluster(),
			DBaaSMongoTemplateResource:                dataSourceDbassMongoTemplate(),
			PsqlVersionsResource:                      dataSourceDbaasPgSqlVersions(),
			PsqlBackupsResource:                       dataSourceDbaasPgSqlBackups(),
			ALBResource:                               dataSourceApplicationLoadBalancer(),
			ALBForwardingRuleResource:                 dataSourceApplicationLoadBalancerForwardingRule(),
			TargetGroupResource:                       dataSourceTargetGroup(),
			DBaasMongoUserResource:                    dataSourceDbaasMongoUser(),
			CertificateResource:                       dataSourceCertificate(),
			ContainerRegistryResource:                 dataSourceContainerRegistry(),
			ContainerRegistryTokenResource:            dataSourceContainerRegistryToken(),
			ContainerRegistryLocationsResource:        dataSourceContainerRegistryLocations(),
			DataplatformClusterResource:               dataSourceDataplatformCluster(),
			DataplatformNodePoolResource:              dataSourceDataplatformNodePool(),
			DataplatformNodePoolsDataSource:           dataSourceDataplatformNodePools(),
			DataplatformVersionsDataSource:            dataSourceDataplatformVersions(),
			DNSZoneDataSource:                         dataSourceDNSZone(),
			DNSRecordDataSource:                       dataSourceDNSRecord(),
			LoggingPipelineDataSource:                 dataSourceLoggingPipeline(),
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

	// Standard client configuration
	clientOpts.Username = username.(string)
	clientOpts.Password = password.(string)
	clientOpts.Token = token.(string)
	clientOpts.Url = cleanedUrl
	clientOpts.Version = ionoscloud.Version
	clientOpts.TerraformVersion = terraformVersion

	return SdkBundle{
		CloudApiClient:     NewClientByType(clientOpts, ionosClient).(*ionoscloud.APIClient),
		PsqlClient:         NewClientByType(clientOpts, psqlClient).(*dbaasService.PsqlClient),
		MongoClient:        NewClientByType(clientOpts, mongoClient).(*dbaasService.MongoClient),
		CertManagerClient:  NewClientByType(clientOpts, certManagerClient).(*cert.Client),
		ContainerClient:    NewClientByType(clientOpts, containerRegistryClient).(*crService.Client),
		DataplatformClient: NewClientByType(clientOpts, dataplatformClient).(*dataplatformService.Client),
		DNSClient:          NewClientByType(clientOpts, dnsClient).(*dnsService.Client),
		LoggingClient:      NewClientByType(clientOpts, loggingClient).(*loggingService.Client),
	}, nil
}

func NewClientByType(clientOpts ClientOptions, clientType clientType) interface{} {
	switch clientType {
	case ionosClient:
		{
			newConfig := ionoscloud.NewConfiguration(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url)
			newConfig.UserAgent = fmt.Sprintf(
				"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
				Version, ionoscloud.Version, clientOpts.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
			if os.Getenv(utils.IonosDebug) != "" {
				newConfig.Debug = true
			}
			newConfig.MaxRetries = utils.MaxRetries
			newConfig.WaitTime = utils.MaxWaitTime
			newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
			return ionoscloud.NewAPIClient(newConfig)
		}
	case psqlClient:
		return dbaasService.NewPsqlClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.Username)
	case mongoClient:
		return dbaasService.NewMongoClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case certManagerClient:
		return cert.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case containerRegistryClient:
		return crService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case dataplatformClient:
		return dataplatformService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case dnsClient:
		return dnsService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
	case loggingClient:
		return loggingService.NewClient(clientOpts.Username, clientOpts.Password, clientOpts.Token, clientOpts.Url, clientOpts.Version, clientOpts.TerraformVersion)
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

// getStateChangeConf gets the default configuration for tracking a request progress
func getStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *resource.StateChangeConf {
	stateConf := &resource.StateChangeConf{
		Pending:        resourcePendingStates,
		Target:         resourceTargetStates,
		Refresh:        resourceStateRefreshFunc(meta, location),
		Timeout:        d.Timeout(timeoutType),
		MinTimeout:     5 * time.Second,
		Delay:          0,   // Don't delay the start
		NotFoundChecks: 600, //Setting high number, to support long timeouts
	}

	return stateConf
}

type RequestFailedError struct {
	msg string
}

func (e RequestFailedError) Error() string {
	return e.msg
}

func IsRequestFailed(err error) bool {
	_, ok := err.(RequestFailedError)
	return ok
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(SdkBundle).CloudApiClient

		log.Printf("[INFO] Checking PATH %s\n", path)
		if path == "" {
			return nil, "", fmt.Errorf("can not check a state when path is empty")
		}

		request, apiResponse, err := client.GetRequestStatus(context.Background(), path)
		logApiRequestTime(apiResponse)
		if err != nil {
			return nil, "", fmt.Errorf("request failed with following error: %w", err)
		}
		if request != nil && request.Metadata != nil && request.Metadata.Status != nil {
			if *request.Metadata.Status == "FAILED" {
				var msg string
				if request.Metadata.Message != nil {
					msg = fmt.Sprintf("Request failed with following error: %s", *request.Metadata.Message)
				} else {
					msg = "Request failed with an unknown error"
				}
				return nil, "", RequestFailedError{msg}
			}

			if *request.Metadata.Status == "DONE" {
				return request, "DONE", nil
			}
		} else {
			if request == nil {
				log.Printf("[DEBUG] request is nil")
			} else if request.Metadata == nil {
				log.Printf("[DEBUG] request metadata is nil")
			}
			if request != nil && request.Metadata != nil && request.Metadata.Message != nil {
				log.Printf("[DEBUG] request failed with following error: %s", *request.Metadata.Message)
			}
			if apiResponse != nil {
				log.Printf("[DEBUG] response message %s", apiResponse.Message)
			}
			return nil, "", fmt.Errorf("request metadata status is nil for path %s", path)
		}

		return nil, *request.Metadata.Status, nil
	}
}

// resourcePendingStates defines states of working in progress
var resourcePendingStates = []string{
	"RUNNING",
	"QUEUED",
}

// resourceTargetStates defines states of completion
var resourceTargetStates = []string{
	"DONE",
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(utils.DefaultTimeout),
	Update:  schema.DefaultTimeout(utils.DefaultTimeout),
	Delete:  schema.DefaultTimeout(utils.DefaultTimeout),
	Default: schema.DefaultTimeout(utils.DefaultTimeout),
}
