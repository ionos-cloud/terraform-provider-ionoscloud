package ionoscloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"ionoscloud": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"ionoscloud": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccProtoV5ProviderFactoriesInternal(t *testing.T, v **schema.Provider) map[string]func() (tfprotov5.ProviderServer, error) {
	providerServerFactory, p, err := ProtoV5ProviderServerFactory(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	providerServer := providerServerFactory()
	*v = p

	return map[string]func() (tfprotov5.ProviderServer, error){
		"ionoscloud": func() (tfprotov5.ProviderServer, error) {
			return providerServer, nil
		},
	}
}

func ProtoV5ProviderServerFactory(ctx context.Context) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	primary := Provider()
	servers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(provider.New()),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, primary, nil
}

func ProtoV6ProviderServerFactory(ctx context.Context) (func() tfprotov6.ProviderServer, *schema.Provider, error) {
	primary := Provider()
	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		primary.GRPCProvider, // Example terraform-plugin-sdk provider
	)
	if err != nil {
		return nil, nil, err
	}
	servers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(provider.New()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)
	if err != nil {
		return nil, nil, err
	}
	return muxServer.ProviderServer, primary, nil
}

func testAccProtoV6ProviderFactoriesInternal(t *testing.T, v **schema.Provider) map[string]func() (tfprotov6.ProviderServer, error) {
	providerServerFactory, p, err := ProtoV6ProviderServerFactory(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	providerServer := providerServerFactory()
	*v = p

	return map[string]func() (tfprotov6.ProviderServer, error){
		"ionoscloud": func() (tfprotov6.ProviderServer, error) {
			return providerServer, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	// pbUsername := os.Getenv(ionoscloud.IonosUsernameEnvVar)
	// pbPassword := os.Getenv(ionoscloud.IonosPasswordEnvVar)
	// pbToken := os.Getenv(ionoscloud.IonosTokenEnvVar)
	// if pbToken == "" {
	//	if pbUsername == "" || pbPassword == "" {
	//		t.Fatalf("%s/%s or %s must be set for acceptance tests", ionoscloud.IonosUsernameEnvVar, ionoscloud.IonosPasswordEnvVar, ionoscloud.IonosTokenEnvVar)
	//	}
	//}

	diags := testAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}

	return
}

func randomProviderVersion343() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			VersionConstraint: "3.4.3",
			Source:            "hashicorp/random",
		},
	}
}
