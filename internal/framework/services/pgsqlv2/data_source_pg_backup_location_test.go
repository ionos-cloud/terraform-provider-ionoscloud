//go:build all || dbaas || psqlv2

package pgsqlv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

// TestAccPgBackupLocationsV2DataSource tests the backup locations data source.
// This is a standalone test that does not require a cluster.
func TestAccPgBackupLocationsV2DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: backupLocationDSConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(backupLocationDSAddr, "backup_locations.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet(backupLocationDSAddr, "backup_locations.0.id"),
					resource.TestCheckResourceAttrSet(backupLocationDSAddr, "backup_locations.0.location"),
				),
			},
		},
	})
}

// --- Configs ---

var backupLocationDSConfig = fmt.Sprintf(`
data "ionoscloud_pg_backup_location_v2" "test" {
  location = "%s"
}
`, testLocation)
