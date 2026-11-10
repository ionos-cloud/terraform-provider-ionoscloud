package ionoscloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// A resource deleted outside Terraform is removed from the state on refresh instead of failing the plan.
func TestReadRemovesResourceDeletedOutsideTerraform(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"httpStatus":404,"messages":[{"errorCode":"309","message":"Resource does not exist"}]}`))
	}))
	t.Cleanup(server.Close)

	t.Setenv(shared.IonosApiUrlEnvVar, "")
	t.Setenv("IONOS_API_URL_NFS", server.URL)
	t.Setenv("IONOS_API_URL_VPN", server.URL)

	ctx := context.Background()
	meta := *bundleclient.New(ctx, clientoptions.TerraformClientOptions{
		ClientOptions: shared.ClientOptions{
			Endpoint:    server.URL,
			Credentials: shared.Credentials{Token: "token-for-the-stub"},
		},
	}, nil)

	tests := []struct {
		name     string
		resource *schema.Resource
		read     schema.ReadContextFunc
		attrs    map[string]string
	}{
		{name: "nsg", resource: resourceNSG(), read: resourceNSGRead, attrs: map[string]string{"datacenter_id": "dc-1"}},
		{name: "nfs cluster", resource: resourceNFSCluster(), read: resourceNFSClusterRead},
		{name: "nfs share", resource: resourceNFSShare(), read: resourceNFSShareRead, attrs: map[string]string{"cluster_id": "cluster-1"}},
		{name: "wireguard gateway", resource: resourceVpnWireguardGateway(), read: resourceVpnWireguardGatewayRead},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.resource.Data(&terraform.InstanceState{ID: "gone", Attributes: tt.attrs})
			if diags := tt.read(ctx, d, meta); diags.HasError() {
				t.Fatalf("Read returned %v, want no error", diags)
			}
			if d.Id() != "" {
				t.Errorf("ID = %q, want it cleared", d.Id())
			}
		})
	}
}
