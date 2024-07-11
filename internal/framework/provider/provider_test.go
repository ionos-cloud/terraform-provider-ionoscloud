package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
)

func TestMuxServer(t *testing.T) {
	ctx := context.Background()
	var pv *schema.Provider
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(ctx, t, &pv),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "ionoscloud" {
					}`,
			},
		},
	})
}

func testAccProtoV5ProviderFactoriesInternal(ctx context.Context, t *testing.T, v **schema.Provider) map[string]func() (tfprotov5.ProviderServer, error) {
	providerServerFactory, p, err := ProtoV5ProviderServerFactory(ctx)

	if err != nil {
		t.Fatal(err)
	}

	providerServer := providerServerFactory()
	*v = p

	return map[string]func() (tfprotov5.ProviderServer, error){
		acctest.ProviderName: func() (tfprotov5.ProviderServer, error) {
			return providerServer, nil
		},
	}
}

func ProtoV5ProviderServerFactory(ctx context.Context) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	primary := ionoscloud.Provider()
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
