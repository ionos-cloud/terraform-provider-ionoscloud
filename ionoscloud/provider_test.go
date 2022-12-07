package ionoscloud

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

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
		t.Fatalf("err: %w", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	pbUsername := os.Getenv(ionoscloud.IonosUsernameEnvVar)
	pbPassword := os.Getenv(ionoscloud.IonosPasswordEnvVar)
	pbToken := os.Getenv(ionoscloud.IonosTokenEnvVar)
	if pbToken == "" {
		if pbUsername == "" || pbPassword == "" {
			t.Fatalf("%s/%s or %s must be set for acceptance tests", ionoscloud.IonosUsernameEnvVar, ionoscloud.IonosPasswordEnvVar, ionoscloud.IonosTokenEnvVar)
		}
	}

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
