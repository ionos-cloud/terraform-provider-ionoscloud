//go:build all || objectstorage || objectstoragemanagement
// +build all objectstorage objectstoragemanagement

package objectstoragemanagement_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccAccesskeyResource(t *testing.T) {
	description := acctest.GenerateRandomResourceName("description")
	descriptionUpdated := acctest.GenerateRandomResourceName("description")
	name := "ionoscloud_object_storage_accesskey.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyConfigDescription(description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", description),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "canonical_user_id"),
					resource.TestCheckResourceAttrSet(name, "contract_user_id"),
				),
			},
			{
				Config: testAccAccesskeyConfigDescription(descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", descriptionUpdated),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "canonical_user_id"),
					resource.TestCheckResourceAttrSet(name, "contract_user_id"),
				),
			},
		},
	})
}

// TestAccAccesskeyFailoverCreatesOnSecondEndpoint verifies that when the first configured endpoint
// returns a retryable HTTP 503 error, the failover round-tripper retries the request against
// the second endpoint (the real IONOS API) and the user is successfully created there.
//
// The test starts a local HTTP server that always returns 503, writes a temporary file config
// that lists the mock server first and the real IONOS API second, and configures the provider
// via IONOS_CONFIG_FILE to use that config with the roundRobin failover strategy.
func TestAccAccesskeyFailoverCreatesOnSecondEndpoint(t *testing.T) {
	if os.Getenv("IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT") != "" {
		t.Skip("IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT is set; this test requires control over endpoint configuration")
	}

	// Count how many times the mock server is called so we can assert failover occurred.
	var mockCallCount int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&mockCallCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer mockServer.Close()

	description := acctest.GenerateRandomResourceName("failover-desc")
	name := "ionoscloud_object_storage_accesskey.test"

	// File config: mock server (returns 503) first, real IONOS API second.
	configContent := fmt.Sprintf(`version: 1.0
environments:
  - name: default
    products:
      - name: objectstoragemanagement
        endpoints:
          - name: %s
          - name: https://s3.ionos.com
failover:
  strategy: roundRobin
  failoverOnStatusCodes:
    - 503
  retryableMethods:
    - GET
    - POST
    - PUT
    - DELETE
    - HEAD
    - OPTIONS
  maxRetries: 1
`, mockServer.URL)

	tmpFile := createTempConfigFile(t, configContent)
	defer os.Remove(tmpFile)

	t.Setenv("IONOS_CONFIG_FILE", tmpFile)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyConfigDescription(description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", description),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "canonical_user_id"),
					resource.TestCheckResourceAttrSet(name, "contract_user_id"),
					// Verify the mock server was hit, proving requests went through the failover path.
					func(_ *terraform.State) error {
						if atomic.LoadInt32(&mockCallCount) == 0 {
							return fmt.Errorf("mock server was never called; failover round-tripper did not route through the first endpoint")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAccesskeyFailoverNetworkError verifies failover when a network-level error occurs (e.g., connection refused).
func TestAccAccesskeyFailoverNetworkError(t *testing.T) {
	if os.Getenv("IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT") != "" {
		t.Skip("IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT is set; this test requires control over endpoint configuration")
	}

	// Use an address that is likely to refuse connection
	badAddr := "http://127.0.0.1:1"

	description := acctest.GenerateRandomResourceName("failover-net")
	descriptionUpdated := acctest.GenerateRandomResourceName("failover-net-upd")
	name := "ionoscloud_object_storage_accesskey.test"

	configContent := fmt.Sprintf(`version: 1.0
environments:
  - name: default
    products:
      - name: objectstoragemanagement
        endpoints:
          - name: %s
          - name: https://s3.ionos.com
failover:
  strategy: roundRobin
  retryableMethods:
    - GET
    - POST
    - DELETE
    - PUT
  maxRetries: 1
`, badAddr)

	tmpFile := createTempConfigFile(t, configContent)
	defer os.Remove(tmpFile)
	t.Setenv("IONOS_CONFIG_FILE", tmpFile)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyConfigDescription(description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", description),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
				),
			},
			{
				Config: testAccAccesskeyConfigDescription(descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", descriptionUpdated),
					resource.TestCheckResourceAttrSet(name, "id"),
				),
			},
		},
	})
}

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "ionos-failover-test-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp config file: %s", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write config file: %s", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func testAccAccesskeyConfigDescription(description string) string {
	return fmt.Sprintf(`
resource "ionoscloud_object_storage_accesskey" "test" {
  description = %[1]q
}
`, description)
}
