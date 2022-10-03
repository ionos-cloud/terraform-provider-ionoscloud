package dbaas

import (
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
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

func NewClientService(username, password, token, url string) PsqlClientService {
	newConfigDbaas := psql.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 4 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}

	return &clientService{
		client: psql.NewAPIClient(newConfigDbaas),
	}
}

func NewMongoClientService(username, password, token, url string) MongoClientService {
	newConfigDbaas := mongo.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = 999
	newConfigDbaas.MaxWaitTime = 4 * time.Second

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}

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
