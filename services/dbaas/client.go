package dbaas

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type MongoClient struct {
	sdkClient *mongo.APIClient
}

type PsqlClient struct {
	sdkClient *psql.APIClient
}

func NewMongoClient(username, password, token, url, version, terraformVersion string, insecure bool) *MongoClient {
	newConfigDbaas := mongo.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, mongo.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &MongoClient{
		sdkClient: mongo.NewAPIClient(newConfigDbaas),
	}
}

func NewPsqlClient(username, password, token, url, version, terraformVersion string, insecure bool) *PsqlClient {
	newConfigDbaas := shared.NewConfiguration(username, password, token, url)

	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, psql.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
	psqlClient := &PsqlClient{
		sdkClient: psql.NewAPIClient(newConfigDbaas),
	}
	psqlClient.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	return psqlClient
}
