package inmemorydb

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	inMemoryDB "github.com/ionos-cloud/sdk-go-dbaas-in-memory-db"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

//nolint:docs
type InMemoryDBClient struct {
	sdkClient *inMemoryDB.APIClient
}

//nolint:docs
func NewInMemoryDBClient(username, password, token, url, version, terraformVersion string) *InMemoryDBClient {
	newConfigDbaas := inMemoryDB.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-in-memory-db/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, inMemoryDB.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &InMemoryDBClient{
		sdkClient: inMemoryDB.NewAPIClient(newConfigDbaas),
	}
}
