//go:build all || dbaas || psqlv2

package pgsqlv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

// TestAccPgVersionsV2DataSource tests the versions data source.
// This is a standalone test that does not require a cluster.
func TestAccPgVersionsV2DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: versionsDSConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(versionsDSAddr, "versions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "versions.0.id"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "versions.0.version"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "versions.0.status"),
				),
			},
		},
	})
}

// --- Configs ---

var versionsDSConfig = fmt.Sprintf(`
data "ionoscloud_pg_versions_v2" "test" {
  location = "%s"
}
`, testLocation)
