package redisdb

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	redis "github.com/ionos-cloud/sdk-go-dbaas-redis"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"net/http"
	"os"
	"runtime"
)

type RedisDBClient struct {
	sdkClient *redis.APIClient
}

func NewRedisDBClient(username, password, token, url, version, terraformVersion string) *RedisDBClient {
	newConfigDbaas := redis.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-redisdb/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, redis.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &RedisDBClient{
		sdkClient: redis.NewAPIClient(newConfigDbaas),
	}
}
