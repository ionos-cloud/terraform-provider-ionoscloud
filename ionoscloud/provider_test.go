package ionoscloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"os"
	"testing"

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

func testAccPreCheck(t *testing.T) {
	pbUsername := os.Getenv(ionoscloud.IonosUsernameEnvVar)
	pbPassword := os.Getenv(ionoscloud.IonosPasswordEnvVar)
	pbToken := os.Getenv(ionoscloud.IonosTokenEnvVar)
	if pbToken == "" {
		if pbUsername == "" || pbPassword == "" {
			t.Fatalf("%s/%s or %s must be set for acceptance tests", ionoscloud.IonosUsernameEnvVar, ionoscloud.IonosPasswordEnvVar, ionoscloud.IonosTokenEnvVar)
		}
	} else {
		if pbUsername != "" || pbPassword != "" {
			t.Fatalf("%s/%s can not be set together with %s", ionoscloud.IonosUsernameEnvVar, ionoscloud.IonosPasswordEnvVar, ionoscloud.IonosTokenEnvVar)
		}

	}

	diags := testAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}

	return
}
