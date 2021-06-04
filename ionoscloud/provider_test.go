package ionoscloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"kubernetes": func() (*schema.Provider, error) {
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
	pbUsername := os.Getenv("IONOS_USERNAME")
	pbPassword := os.Getenv("IONOS_PASSWORD")
	pbToken := os.Getenv("IONOS_TOKEN")
	if pbToken == "" {
		if pbUsername == "" || pbPassword == "" {
			t.Fatal("IONOS_USERNAME/IONOS_PASSWORD or IONOS_TOKEN must be set for acceptance tests")
		}
	} else {
		if pbUsername != "" || pbPassword != "" {
			t.Fatal("IONOS_USERNAME/IONOS_PASSWORD or IONOS_TOKEN must be set for acceptance tests")
		}

	}
	diags := testAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}
	return
}
