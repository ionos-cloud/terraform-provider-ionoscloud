//go:build all || dbaas || inmemorydbv2

package inmemorydbv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccInMemoryDBV2Cluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "location", testLocation),
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-inmemorydbv2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Test InMemoryDB v2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "version", "9.0"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "persistence_mode", "None"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "eviction_policy", "allkeys-lru"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "true"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "dns_name"),

					// Instances
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "4"),

					// Connection
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.primary_instance_address", "192.168.2.101/24"),

					// Snapshot
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.location", "eu-central-3"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.retention_days", "7"),

					// Maintenance window
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "09:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Sunday"),

					// Credentials (username readable; password hash is sensitive)
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.username", "cacheadmin"),
				),
			},
			// TODO -- For this step we will add another cluster because we are also testing some filtering
			// TODO -- This step needs to be reviewed thoroughly
			{
				Config: clusterCreateConfig + clusterDataSourcesConfig,
				Check: resource.ComposeTestCheckFunc(
					// by ID
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "dns_name", clusterResourceAddr, "dns_name"),

					// by name
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "dns_name", clusterResourceAddr, "dns_name"),

					// list with name filter
					resource.TestCheckResourceAttr(clustersDSAddr, "items.#", "1"),
					// list without name filter
					resource.TestCheckResourceAttrSet(clustersAllDSAddr, "items.#"),

					// snapshots filtered by cluster
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.#"),
					// snapshots without filter
					resource.TestCheckResourceAttrSet(snapshotsAllDSAddr, "items.#"),

					// snapshot locations list
					resource.TestCheckResourceAttrSet(snapshotLocationsDSAddr, "items.#"),

					// versions list
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.#"),
				),
			},
			// List without filters — confirms the cluster appears in results.
			{
				Query: true,
				Config: `list "ionoscloud_inmemorydb_cluster_v2" "test" {
  provider = ionoscloud
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(clusterResourceAddr, map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"location": knownvalue.StringExact(testLocation),
					}),
				},
			},
			// Filter by name + correct location: unique name guarantees exactly 1 result.
			{
				Query: true,
				Config: `list "ionoscloud_inmemorydb_cluster_v2" "test" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = "tf-test-inmemorydbv2" },
      { field_name = "location", field_value = "` + testLocation + `" },
    ]
  }
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 1),
				},
			},
			// Filter by name + wrong location: must return 0, proving location filter is evaluated.
			{
				Query: true,
				Config: `list "ionoscloud_inmemorydb_cluster_v2" "test" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = "tf-test-inmemorydbv2" },
      { field_name = "location", field_value = "` + testLocationOther + `" },
    ]
  }
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 0),
				},
			},
			// TODO -- After manual testing, review this update configuration
			{
				Config: clusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-inmemorydbv2-updated"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Updated InMemoryDB v2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "eviction_policy", "allkeys-lfu"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "8"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.retention_days", "14"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "12:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Wednesday"),
				),
			},
			// TODO -- Review this step
			{
				ResourceName:      clusterResourceAddr,
				ImportStateIdFunc: clusterImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"credentials.password.hash",
					"credentials.password.algorithm",
					"timeouts",
				},
			},
		},
	})
}

var clusterDataSourcesConfig = fmt.Sprintf(`
data "ionoscloud_inmemorydb_cluster_v2" "by_id" {
  id       = ionoscloud_inmemorydb_cluster_v2.test.id
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_cluster_v2" "by_name" {
  name     = ionoscloud_inmemorydb_cluster_v2.test.name
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_clusters_v2" "test" {
  location = "%[1]s"
  name     = ionoscloud_inmemorydb_cluster_v2.test.name
}

data "ionoscloud_inmemorydb_clusters_v2" "all" {
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_snapshots_v2" "test" {
  location   = "%[1]s"
  cluster_id = ionoscloud_inmemorydb_cluster_v2.test.id
}

data "ionoscloud_inmemorydb_snapshots_v2" "all" {
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_snapshot_locations_v2" "test" {
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_versions_v2" "test" {
  location = "%[1]s"
}
`, testLocation)
