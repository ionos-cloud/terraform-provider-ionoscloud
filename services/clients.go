package services

import (
	"github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
)

type SdkBundle struct {
	CloudApiClient     *ionoscloud.APIClient
	PsqlClient         *dbaas.PsqlClient
	MongoClient        *dbaas.MongoClient
	CertManagerClient  *cert.Client
	ContainerClient    *containerregistry.Client
	DataplatformClient *dataplatform.Client
	DNSClient          *dns.Client
	LoggingClient      *logging.Client
	AutoscalingClient  *autoscaling.Client
}
