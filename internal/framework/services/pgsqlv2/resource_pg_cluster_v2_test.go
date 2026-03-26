//go:build all || dbaas || psqlv2

package pgsqlv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

// TestAccPgClusterV2 tests the full lifecycle of the PgSQL v2 cluster resource
// and the cluster/clusters data sources.
func TestAccPgClusterV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: clusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterV2Exists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-pgsqlv2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Test PgSQL v2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "version", "17"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "location", testLocation),
					resource.TestCheckResourceAttr(clusterResourceAddr, "backup_location", "de"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "replication_mode", "ASYNCHRONOUS"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connection_pooler", "TRANSACTION"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "4"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.storage_size", "10"),
					resource.TestCheckResourceAttrPair(clusterResourceAddr, "connections.datacenter_id", "ionoscloud_datacenter.test", "id"),
					resource.TestCheckResourceAttrPair(clusterResourceAddr, "connections.lan_id", "ionoscloud_lan.test", "id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.primary_instance_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "09:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.username", "testuser"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.database", "testdb"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "dns_name"),
				),
			},
			{
				ResourceName:            clusterResourceAddr,
				ImportState:             true,
				ImportStateIdFunc:       pgClusterV2ImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "credentials.password", "restore_from_backup"},
			},
			{
				Config: clusterDSByIDConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "backup_location", clusterResourceAddr, "backup_location"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "replication_mode", clusterResourceAddr, "replication_mode"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connection_pooler", clusterResourceAddr, "connection_pooler"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.storage_size", clusterResourceAddr, "instances.storage_size"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
				),
			},
			{
				Config: clusterDSByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "backup_location", clusterResourceAddr, "backup_location"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "replication_mode", clusterResourceAddr, "replication_mode"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connection_pooler", clusterResourceAddr, "connection_pooler"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.storage_size", clusterResourceAddr, "instances.storage_size"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
				),
			},
			{
				Config: clustersDSFilteredByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(clustersDSAddr, "clusters.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.backup_location", clusterResourceAddr, "backup_location"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "clusters.0.replication_mode", clusterResourceAddr, "replication_mode"),
				),
			},
			{
				Config: clustersDSNoFilterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(clustersDSAddr, "clusters.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
			{
				Config: backupsDSByClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(backupsDSAddr, "location", testLocation),
					resource.TestCheckResourceAttrPair(backupsDSAddr, "cluster_id", clusterResourceAddr, "id"),
				),
			},
			{
				Config: backupsDSAllConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(backupsDSAddr, "location", testLocation),
				),
			},
			{
				Config:      clusterDSErrorBothIDAndName,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      clusterDSErrorNoIDNoName,
				ExpectError: regexp.MustCompile("Missing Attribute Configuration"),
			},
			{
				Config:      clusterDSErrorInvalidName,
				ExpectError: regexp.MustCompile("no PostgreSQL v2 cluster found"),
			},
			{
				Config:      clusterDSErrorInvalidID,
				ExpectError: regexp.MustCompile("failed to get PostgreSQL v2 cluster"),
			},
			{
				Config: clusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterV2Exists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-pgsqlv2-updated"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Updated PgSQL v2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "version", "17"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "location", testLocation),
					resource.TestCheckResourceAttr(clusterResourceAddr, "backup_location", "eu-central-2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "replication_mode", "ASYNCHRONOUS"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connection_pooler", "DISABLED"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "4"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.storage_size", "20"),
					resource.TestCheckResourceAttrPair(clusterResourceAddr, "connections.datacenter_id", "ionoscloud_datacenter.test", "id"),
					resource.TestCheckResourceAttrPair(clusterResourceAddr, "connections.lan_id", "ionoscloud_lan.test", "id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.primary_instance_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "12:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Wednesday"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.username", "testuser"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.database", "testdb"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "dns_name"),
				),
			},
			{
				Config: clusterLocationChangeConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterV2ExistsInLocation(clusterResourceAddr, testLocationChanged),
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-pgsqlv2-updated"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "location", testLocationChanged),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "4"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.storage_size", "20"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "dns_name"),
				),
			},
		},
	})
}

// --- Cluster data source configs ---

var clusterDSByIDConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "by_id" {
  id       = ionoscloud_pg_cluster_v2.test.id
  location = "%s"
}
`, testLocation)

var clusterDSByNameConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "by_name" {
  name     = ionoscloud_pg_cluster_v2.test.name
  location = "%s"
}
`, testLocation)

var clusterDSErrorBothIDAndName = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "error_test" {
  id       = ionoscloud_pg_cluster_v2.test.id
  name     = ionoscloud_pg_cluster_v2.test.name
  location = "%s"
}
`, testLocation)

var clusterDSErrorNoIDNoName = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "error_test" {
  location = "%s"
}
`, testLocation)

var clusterDSErrorInvalidName = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "error_test" {
  name     = "this-cluster-does-not-exist"
  location = "%s"
}
`, testLocation)

var clusterDSErrorInvalidID = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_cluster_v2" "error_test" {
  id       = "00000000-0000-0000-0000-000000000000"
  location = "%s"
}
`, testLocation)

// --- Clusters data source configs ---

var clustersDSFilteredByNameConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_clusters_v2" "test" {
  location = "%s"
  name     = ionoscloud_pg_cluster_v2.test.name
}
`, testLocation)

var clustersDSNoFilterConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_clusters_v2" "test" {
  location = "%s"
}
`, testLocation)

// --- Backups data source configs ---

var backupsDSByClusterConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_backups_v2" "test" {
  location   = "%s"
  cluster_id = ionoscloud_pg_cluster_v2.test.id
}
`, testLocation)

var backupsDSAllConfig = clusterCreateConfig + fmt.Sprintf(`
data "ionoscloud_pg_backups_v2" "test" {
  location = "%s"
}
`, testLocation)
