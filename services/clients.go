package services

import (
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
)

type SdkBundle struct {
	CloudApiClient     *ionoscloud.APIClient
	RedisClient        *dbaas.RedisClient
	PsqlClient         *dbaas.PsqlClient
	MongoClient        *dbaas.MongoClient
	MariaDBClient      *mariadb.MariaDBClient
	NFSClient          *nfs.Client
	CertManagerClient  *cert.Client
	ContainerClient    *containerregistry.Client
	DataplatformClient *dataplatform.Client
	DNSClient          *dns.Client
	LoggingClient      *logging.Client
	AutoscalingClient  *autoscaling.Client
	KafkaClient        *kafka.Client
	CDNClient          *cdn.Client
	APIGatewayClient   *apigateway.Client
	VPNClient          *vpn.Client
}
