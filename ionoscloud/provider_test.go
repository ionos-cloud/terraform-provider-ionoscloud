package ionoscloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ionoscloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
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
