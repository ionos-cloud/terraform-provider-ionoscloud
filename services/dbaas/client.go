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
	sdkClient *psql.APIClient
}

type MongoClient struct {
	sdkClient *mongo.APIClient
}

func NewPsqlClient(username, password, token, url, version, terraformVersion string) *PsqlClient {
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

	return &PsqlClient{
		sdkClient: psql.NewAPIClient(newConfigDbaas),
	}
}

func NewMongoClient(username, password, token, url, version, terraformVersion string) *MongoClient {
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

	return &MongoClient{
		sdkClient: mongo.NewAPIClient(newConfigDbaas),
	}
}
