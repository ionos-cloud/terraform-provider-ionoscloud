package dbaas

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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

//
//func NewMongoClient(username, password, token, url, version, terraformVersion string, insecure bool) *MongoClient {
//	newConfigDbaas := mongo.NewConfiguration(username, password, token, url)
//
//	if os.Getenv(constant.IonosDebug) != "" {
//		newConfigDbaas.Debug = true
//	}
//	newConfigDbaas.MaxRetries = constant.MaxRetries
//	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime
//
//	newConfigDbaas.HTTPClient = &http.sdkClient{Transport: utils.CreateTransport(insecure)}
//	newConfigDbaas.UserAgent = fmt.Sprintf(
//		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
//		version, mongo.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
//
//	return &MongoClient{
//		sdkClient: mongo.NewAPIClient(newConfigDbaas),
//	}
//}

func NewMongoClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *MongoClient {
	loadedconfig.SetClientOptionsFromLoadedConfig(&clientOptions, sharedLoadedConfig, shared.Cloud)
	newConfig := mongo.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)
	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-dbaas-mongo/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		mongo.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfig.Debug = true
	}
	newConfig.MaxRetries = constant.MaxRetries
	newConfig.WaitTime = constant.MaxWaitTime
	newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	client := MongoClient{
		sdkClient: mongo.NewAPIClient(newConfig),
	}
	return &client
}

//func NewPsqlClient(username, password, token, url, version, terraformVersion string, insecure bool) *PsqlClient {
//	newConfigDbaas := psql.NewConfiguration(username, password, token, url)
//
//	if os.Getenv(constant.IonosDebug) != "" {
//		newConfigDbaas.Debug = true
//	}
//	newConfigDbaas.MaxRetries = constant.MaxRetries
//	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime
//
//	newConfigDbaas.HTTPClient = &http.sdkClient{Transport: utils.CreateTransport(insecure)}
//	newConfigDbaas.UserAgent = fmt.Sprintf(
//		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
//		version, psql.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
//
//	return &PsqlClient{
//		sdkClient: psql.NewAPIClient(newConfigDbaas),
//	}
//}

func NewPSQLClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *PsqlClient {
	loadedconfig.SetClientOptionsFromLoadedConfig(&clientOptions, sharedLoadedConfig, shared.Cloud)
	newConfig := psql.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-postgres/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, psql.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
	if os.Getenv(constant.IonosDebug) != "" {
		newConfig.Debug = true
	}
	newConfig.MaxRetries = constant.MaxRetries
	newConfig.WaitTime = constant.MaxWaitTime
	newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	client := PsqlClient{
		sdkClient: psql.NewAPIClient(newConfig),
	}
	return &client
}
