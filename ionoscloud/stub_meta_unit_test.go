package ionoscloud

import (
	"context"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// stubMeta returns a provider meta whose Cloud API client talks to the stub at url.
func stubMeta(t *testing.T, url string) bundleclient.SdkBundle {
	t.Helper()
	t.Setenv(shared.IonosApiUrlEnvVar, "")
	return *bundleclient.New(context.Background(), clientoptions.TerraformClientOptions{
		ClientOptions: shared.ClientOptions{
			Endpoint:    url,
			Credentials: shared.Credentials{Token: "token-for-the-stub"},
		},
	}, nil)
}
