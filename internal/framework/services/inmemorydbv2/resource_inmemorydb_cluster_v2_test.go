//go:build all || dbaas || inmemorydbv2

package inmemorydbv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
			{
				Config: clusterCreateConfig + clusterDataSourcesConfig,
				Check: resource.ComposeTestCheckFunc(
					// by ID — all attributes
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "location", clusterResourceAddr, "location"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "persistence_mode", clusterResourceAddr, "persistence_mode"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "eviction_policy", clusterResourceAddr, "eviction_policy"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "snapshot.location", clusterResourceAddr, "snapshot.location"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "snapshot.retention_days", clusterResourceAddr, "snapshot.retention_days"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "snapshot.snapshot_hours.#", clusterResourceAddr, "snapshot.snapshot_hours.#"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "credentials.username", clusterResourceAddr, "credentials.username"),

					// by name — all attributes
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "location", clusterResourceAddr, "location"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "persistence_mode", clusterResourceAddr, "persistence_mode"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "eviction_policy", clusterResourceAddr, "eviction_policy"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "snapshot.location", clusterResourceAddr, "snapshot.location"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "snapshot.retention_days", clusterResourceAddr, "snapshot.retention_days"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "snapshot.snapshot_hours.#", clusterResourceAddr, "snapshot.snapshot_hours.#"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "credentials.username", clusterResourceAddr, "credentials.username"),

					// clusters list with name filter — verify item count and first item attributes
					resource.TestCheckResourceAttr(clustersDSAddr, "items.#", "1"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.persistence_mode", clusterResourceAddr, "persistence_mode"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.eviction_policy", clusterResourceAddr, "eviction_policy"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.credentials.username", clusterResourceAddr, "credentials.username"),

					// clusters list without filter — both clusters must appear
					resource.TestCheckResourceAttr(clustersAllDSAddr, "items.#", "2"),

					// snapshots filtered by cluster
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.#"),
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.0.id"),
					resource.TestCheckResourceAttrPair(snapshotsDSAddr, "items.0.cluster_id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.0.cluster_name"),
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.0.cluster_version"),
					resource.TestCheckResourceAttrSet(snapshotsDSAddr, "items.0.snapshot_location"),

					// snapshots without filter
					resource.TestCheckResourceAttrSet(snapshotsAllDSAddr, "items.#"),

					// snapshot locations list
					resource.TestCheckResourceAttrSet(snapshotLocationsDSAddr, "items.#"),
					resource.TestCheckResourceAttrSet(snapshotLocationsDSAddr, "items.0.id"),
					resource.TestCheckResourceAttrSet(snapshotLocationsDSAddr, "items.0.snapshot_region"),

					// versions list
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.#"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.id"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.version"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.status"),
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
			// include_resource = true: verify all resource attributes are populated correctly.
			{
				Query: true,
				Config: `list "ionoscloud_inmemorydb_cluster_v2" "test" {
			 provider         = ionoscloud
			 include_resource = true
			 config {
			   filters = [
			     { field_name = "name",     field_value = "tf-test-inmemorydbv2" },
			     { field_name = "location", field_value = "` + testLocation + `" },
			   ]
			 }
			}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceKnownValues(clusterResourceAddr,
						queryfilter.ByDisplayName(knownvalue.StringExact("tf-test-inmemorydbv2")),
						[]querycheck.KnownValueCheck{
							// Top-level
							{Path: tfjsonpath.New("name"), KnownValue: knownvalue.StringExact("tf-test-inmemorydbv2")},
							{Path: tfjsonpath.New("description"), KnownValue: knownvalue.StringExact("Test InMemoryDB v2 cluster")},
							{Path: tfjsonpath.New("version"), KnownValue: knownvalue.StringExact("9.0")},
							{Path: tfjsonpath.New("location"), KnownValue: knownvalue.StringExact(testLocation)},
							{Path: tfjsonpath.New("persistence_mode"), KnownValue: knownvalue.StringExact("None")},
							{Path: tfjsonpath.New("eviction_policy"), KnownValue: knownvalue.StringExact("allkeys-lru")},
							{Path: tfjsonpath.New("logs_enabled"), KnownValue: knownvalue.Bool(true)},
							{Path: tfjsonpath.New("metrics_enabled"), KnownValue: knownvalue.Bool(true)},
							{Path: tfjsonpath.New("dns_name"), KnownValue: knownvalue.NotNull()},
							// Instances
							{Path: tfjsonpath.New("instances").AtMapKey("count"), KnownValue: knownvalue.Int32Exact(1)},
							{Path: tfjsonpath.New("instances").AtMapKey("cores"), KnownValue: knownvalue.Int32Exact(1)},
							{Path: tfjsonpath.New("instances").AtMapKey("ram"), KnownValue: knownvalue.Int32Exact(4)},
							// Connections
							{Path: tfjsonpath.New("connections").AtMapKey("primary_instance_address"), KnownValue: knownvalue.StringExact("192.168.2.101/24")},
							// Snapshot
							{Path: tfjsonpath.New("snapshot").AtMapKey("location"), KnownValue: knownvalue.StringExact("eu-central-3")},
							{Path: tfjsonpath.New("snapshot").AtMapKey("retention_days"), KnownValue: knownvalue.Int32Exact(7)},
							// Maintenance window
							{Path: tfjsonpath.New("maintenance_window").AtMapKey("time"), KnownValue: knownvalue.StringExact("09:00:00")},
							{Path: tfjsonpath.New("maintenance_window").AtMapKey("day_of_the_week"), KnownValue: knownvalue.StringExact("Sunday")},
							// Credentials
							{Path: tfjsonpath.New("credentials").AtMapKey("username"), KnownValue: knownvalue.StringExact("cacheadmin")},
						},
					),
				},
			},
			{
				Config: clusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.retention_days", "14"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.snapshot_hours.#", "2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.snapshot_hours.0", "6"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.snapshot_hours.1", "18"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "12:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Wednesday"),
				),
			},
			{
				ResourceName:      clusterResourceAddr,
				ImportStateIdFunc: clusterImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"credentials.password.%",
					"credentials.password.hash",
					"credentials.password.algorithm",
					"timeouts",
				},
			},
		},
	})
}

var clusterDataSourcesConfig = fmt.Sprintf(`
resource "ionoscloud_inmemorydb_cluster_v2" "test2" {
  location         = "%[1]s"
  name             = "tf-test-inmemorydbv2-second"
  description      = "Second test InMemoryDB v2 cluster"
  version          = "9.0"
  persistence_mode = "None"
  eviction_policy  = "allkeys-lru"
  logs_enabled     = true
  metrics_enabled  = true

  instances = {
    count = 1
    cores = 1
    ram   = 4
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.2.102/24"
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 7
    snapshot_hours = [0, 12]
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "cacheadmin"
    password = {
      algorithm = "SHA-256"
      hash      = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
    }
  }
}

data "ionoscloud_inmemorydb_cluster_v2" "by_id" {
  id       = ionoscloud_inmemorydb_cluster_v2.test.id
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_cluster_v2" "by_name" {
  name     = ionoscloud_inmemorydb_cluster_v2.test.name
  location = "%[1]s"
}

data "ionoscloud_inmemorydb_clusters_v2" "test" {
  location   = "%[1]s"
  name       = ionoscloud_inmemorydb_cluster_v2.test.name
  depends_on = [ionoscloud_inmemorydb_cluster_v2.test2]
}

data "ionoscloud_inmemorydb_clusters_v2" "all" {
  location   = "%[1]s"
  depends_on = [ionoscloud_inmemorydb_cluster_v2.test2]
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
