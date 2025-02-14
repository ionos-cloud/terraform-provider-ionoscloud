package dbaas

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type MongoClient struct {
	sdkClient *mongo.APIClient
}

type PsqlClient struct {
	sdkClient *psql.APIClient
}

func NewMongoClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *MongoClient {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cloud)
	config := mongo.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		mongo.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	fileconfiguration.AddCertsToClient(config.HTTPClient, clientOptions.Certificate)
	client := MongoClient{
		sdkClient: mongo.NewAPIClient(config),
	}
	return &client
}

func NewPSQLClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *PsqlClient {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cloud)
	config := psql.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, psql.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	fileconfiguration.AddCertsToClient(config.HTTPClient, clientOptions.Certificate)
	client := PsqlClient{
		sdkClient: psql.NewAPIClient(config),
	}
	return &client
}
