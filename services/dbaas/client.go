package dbaas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type PsqlClient struct {
	psql.APIClient
}

type MongoClient struct {
	mongo.APIClient
}

type PsqlClientConfig struct {
	psql.Configuration
}

type MongoClientConfig struct {
	mongo.Configuration
}

// PsqlClientService is a wrapper around psql.APIClient
type PsqlClientService interface {
	Get() *PsqlClient
	GetConfig() *PsqlClientConfig
}

type MongoClientService interface {
	Get() *MongoClient
	GetConfig() *MongoClientConfig
}

type clientService struct {
	client *psql.APIClient
}

type mongoClientService struct {
	client *mongo.APIClient
}

var _ PsqlClientService = &clientService{}
var _ MongoClientService = &mongoClientService{}

func NewPsqlClientService(username, password, token, url, version, terraformVersion string) PsqlClientService {
	newConfigDbaas := psql.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 4 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, psql.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &clientService{
		client: psql.NewAPIClient(newConfigDbaas),
	}
}

func NewMongoClientService(username, password, token, url, version, terraformVersion string) MongoClientService {
	newConfigDbaas := mongo.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 4 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, mongo.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &mongoClientService{
		client: mongo.NewAPIClient(newConfigDbaas),
	}
}

func (c clientService) Get() *PsqlClient {
	return &PsqlClient{
		APIClient: *c.client,
	}
}

func (c clientService) GetConfig() *PsqlClientConfig {
	return &PsqlClientConfig{
		Configuration: *c.client.GetConfig(),
	}
}

func (c mongoClientService) Get() *MongoClient {
	return &MongoClient{
		*c.client,
	}
}

func (c mongoClientService) GetConfig() *MongoClientConfig {
	return &MongoClientConfig{
		Configuration: *c.client.GetConfig(),
	}
}
