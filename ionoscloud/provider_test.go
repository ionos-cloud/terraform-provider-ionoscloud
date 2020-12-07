package ionoscloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ionoscloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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
}
