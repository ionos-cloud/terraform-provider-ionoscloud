package dbaas

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
)

type MongoClient struct {
	sdkClient *mongo.APIClient
}

type PsqlClient struct {
	sdkClient *psql.APIClient
}

func NewMongoClientFromConfig(config *shared.Configuration) *MongoClient {
	return &MongoClient{sdkClient: mongo.NewAPIClient(config)}
}

func NewPsqlClientFromConfig(config *shared.Configuration) *PsqlClient {
	return &PsqlClient{sdkClient: psql.NewAPIClient(config)}
}
