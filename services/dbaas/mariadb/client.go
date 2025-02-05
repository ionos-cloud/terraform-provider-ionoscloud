package mariadb

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type MariaDBClient struct {
	sdkClient *mariadb.APIClient
}

func NewMariaDBClient(username, password, token, url, version, terraformVersion string, insecure bool) *MariaDBClient {
	newConfigDbaas := shared.NewConfiguration(username, password, token, url)

	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mariadb/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, mariadb.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
	client := &MariaDBClient{
		sdkClient: mariadb.NewAPIClient(newConfigDbaas),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}

	return client
}
