package bundle

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

type ClientOptions struct {
	shared.ClientOverrideOptions
	Version          string
	TerraformVersion string
}

//func NewSdkBundleClient(clientOptions ClientOptions, loadedConfig *shared.LoadedConfig) *services.SdkBundle {
//	return &services.SdkBundle{
//		CDNClient:          cdn.NewCDNClient(clientOptions, loadedConfig),
//		AutoscalingClient:  autoscaling.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		CertManagerClient:  cert.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		CloudApiClient:     cloudapi.NewClient(clientOptions, loadedConfig),
//		ContainerClient:    containerregistry.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		DataplatformClient: dataplatform.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		DNSClient:          dns.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		LoggingClient:      logging.NewClient(username, password, token, cleanedEndpoint, terraformVersion, insecureBool),
//		MariaDBClient:      mariadb.NewMariaDBClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		MongoClient:        dbaas.NewMongoClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		NFSClient:          nfs.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		PsqlClient:         dbaas.NewPsqlClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		KafkaClient:        kafka.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		APIGatewayClient: apigateway.NewClient(
//			username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool,
//		),
//		VPNClient:                     vpn.NewClient(username, password, token, cleanedEndpoint, terraformVersion, insecureBool),
//		InMemoryDBClient:              inmemorydb.NewInMemoryDBClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		S3Client:                      objectstorage.NewClient(accessKey, secretKey, region, endpoint, insecureBool),
//		ObjectStorageManagementClient: objectstoragemanagement.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//		MonitoringClient:              monitoring.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
//	}
//}
