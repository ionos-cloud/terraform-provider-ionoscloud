package bundleclient

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/sdk-go/v6"

	apiGatewayService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/apigateway"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
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
)

func New(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *SdkBundle {
	return &SdkBundle{
		CDNClient:                     cdnService.NewCDNClient(clientOptions, fileConfig),
		AutoscalingClient:             autoscalingService.NewClient(clientOptions, fileConfig),
		CertManagerClient:             cert.NewClient(clientOptions, fileConfig),
		CloudApiClient:                cloudapi.NewClient(clientOptions, fileConfig),
		ContainerClient:               crService.NewClient(clientOptions, fileConfig),
		DataplatformClient:            dataplatformService.NewClient(clientOptions, fileConfig),
		DNSClient:                     dnsService.NewClient(clientOptions, fileConfig),
		LoggingClient:                 loggingService.NewClient(clientOptions, fileConfig),
		MariaDBClient:                 mariadb.NewClient(clientOptions, fileConfig),
		MongoClient:                   dbaasService.NewMongoClient(clientOptions, fileConfig),
		NFSClient:                     nfsService.NewClient(clientOptions, fileConfig),
		PsqlClient:                    dbaasService.NewPSQLClient(clientOptions, fileConfig),
		KafkaClient:                   kafkaService.NewClient(clientOptions, fileConfig),
		APIGatewayClient:              apiGatewayService.NewClient(clientOptions, fileConfig),
		VPNClient:                     vpn.NewClient(clientOptions, fileConfig),
		InMemoryDBClient:              inmemorydb.NewClient(clientOptions, fileConfig),
		S3Client:                      objectStorageService.NewClient(clientOptions, fileConfig),
		ObjectStorageManagementClient: objectStorageManagementService.NewClient(clientOptions, fileConfig),
		MonitoringClient:              monitoringService.NewClient(clientOptions, fileConfig),
	}
}

type SdkBundle struct {
	CloudApiClient                *ionoscloud.APIClient
	InMemoryDBClient              *inmemorydb.Client
	PsqlClient                    *dbaasService.PsqlClient
	MongoClient                   *dbaasService.MongoClient
	MariaDBClient                 *mariadb.Client
	NFSClient                     *nfsService.Client
	CertManagerClient             *cert.Client
	ContainerClient               *crService.Client
	DataplatformClient            *dataplatformService.Client
	DNSClient                     *dnsService.Client
	LoggingClient                 *loggingService.Client
	AutoscalingClient             *autoscalingService.Client
	KafkaClient                   *kafkaService.Client
	CDNClient                     *cdnService.Client
	APIGatewayClient              *apiGatewayService.Client
	VPNClient                     *vpn.Client
	S3Client                      *objectStorageService.Client
	ObjectStorageManagementClient *objectStorageManagementService.Client
	MonitoringClient              *monitoringService.Client
}
