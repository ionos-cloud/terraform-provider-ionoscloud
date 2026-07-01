//go:build all || dbaas || mariadbv2

package mariadbv2_test

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

func TestAccMariaDBV2Cluster(t *testing.T) {
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
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-mariadbv2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Test MariaDB V2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "version", "11.4"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "true"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "dns_name"),

					// Instances
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "4"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.storage_size", "10"),

					// Connections
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.primary_instance_address", "192.168.2.101/24"),

					// Backup
					resource.TestCheckResourceAttr(clusterResourceAddr, "backup.location", "eu-central-3"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "backup.retention_days", "7"),

					// Maintenance window
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "09:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Sunday"),

					// Credentials (username and database readable; password is write-only)
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.username", "dbadmin"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.database", "my_database"),
				),
			},
			{
				Config: clusterCreateConfig + clusterDataSourcesConfig,
				Check: resource.ComposeTestCheckFunc(
					// by ID
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "location", clusterResourceAddr, "location"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.storage_size", clusterResourceAddr, "instances.storage_size"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "backup.location", clusterResourceAddr, "backup.location"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "backup.retention_days", clusterResourceAddr, "backup.retention_days"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "credentials.username", clusterResourceAddr, "credentials.username"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "credentials.database", clusterResourceAddr, "credentials.database"),

					// by name
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "location", clusterResourceAddr, "location"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "instances.storage_size", clusterResourceAddr, "instances.storage_size"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "backup.location", clusterResourceAddr, "backup.location"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "backup.retention_days", clusterResourceAddr, "backup.retention_days"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "credentials.username", clusterResourceAddr, "credentials.username"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "credentials.database", clusterResourceAddr, "credentials.database"),

					// clusters list with name filter
					resource.TestCheckResourceAttr(clustersDSAddr, "items.#", "1"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.name", clusterResourceAddr, "name"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.description", clusterResourceAddr, "description"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.location", clusterResourceAddr, "location"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.logs_enabled", clusterResourceAddr, "logs_enabled"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.metrics_enabled", clusterResourceAddr, "metrics_enabled"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.dns_name", clusterResourceAddr, "dns_name"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.instances.count", clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.instances.ram", clusterResourceAddr, "instances.ram"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.instances.storage_size", clusterResourceAddr, "instances.storage_size"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.connections.lan_id", clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.connections.primary_instance_address", clusterResourceAddr, "connections.primary_instance_address"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.backup.location", clusterResourceAddr, "backup.location"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.backup.retention_days", clusterResourceAddr, "backup.retention_days"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.maintenance_window.time", clusterResourceAddr, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.maintenance_window.day_of_the_week", clusterResourceAddr, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.credentials.username", clusterResourceAddr, "credentials.username"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.credentials.database", clusterResourceAddr, "credentials.database"),

					// clusters list without filter — verify both test clusters appear
					resource.TestCheckResourceAttr(clustersAllDSAddr, "items.#", "2"),

					// backup locations
					resource.TestCheckResourceAttrSet(backupLocationsAddr, "items.#"),
					resource.TestCheckResourceAttrSet(backupLocationsAddr, "items.0.id"),
					resource.TestCheckResourceAttrSet(backupLocationsAddr, "items.0.location"),

					// versions
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.#"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.id"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.version"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.status"),
				),
			},
			// List filtered by location — avoids querying unavailable regional endpoints.
			{
				Query: true,
				Config: `list "ionoscloud_mariadb_cluster_v2" "test" {
			 provider = ionoscloud
			 config {
			   filters = [
			     { field_name = "location", field_value = "` + testLocation + `" },
			   ]
			 }
			}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(clusterResourceAddr, map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"location": knownvalue.StringExact(testLocation),
					}),
				},
			},
			// Filter by name + correct location
			{
				Query: true,
				Config: `list "ionoscloud_mariadb_cluster_v2" "test" {
			 provider = ionoscloud
			 config {
			   filters = [
			     { field_name = "name",     field_value = "tf-test-mariadbv2" },
			     { field_name = "location", field_value = "` + testLocation + `" },
			   ]
			 }
			}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 1),
				},
			},
			// Filter by name + wrong location: must return 0
			{
				Query: true,
				Config: `list "ionoscloud_mariadb_cluster_v2" "test" {
			 provider = ionoscloud
			 config {
			   filters = [
			     { field_name = "name",     field_value = "tf-test-mariadbv2" },
			     { field_name = "location", field_value = "` + testLocationOther + `" },
			   ]
			 }
			}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 0),
				},
			},
			// include_resource = true: verify all resource attributes
			{
				Query: true,
				Config: `list "ionoscloud_mariadb_cluster_v2" "test" {
			 provider         = ionoscloud
			 include_resource = true
			 config {
			   filters = [
			     { field_name = "name",     field_value = "tf-test-mariadbv2" },
			     { field_name = "location", field_value = "` + testLocation + `" },
			   ]
			 }
			}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceKnownValues(clusterResourceAddr,
						queryfilter.ByDisplayName(knownvalue.StringExact("tf-test-mariadbv2")),
						[]querycheck.KnownValueCheck{
							{Path: tfjsonpath.New("name"), KnownValue: knownvalue.StringExact("tf-test-mariadbv2")},
							{Path: tfjsonpath.New("description"), KnownValue: knownvalue.StringExact("Test MariaDB V2 cluster")},
							{Path: tfjsonpath.New("version"), KnownValue: knownvalue.StringExact("11.4")},
							{Path: tfjsonpath.New("location"), KnownValue: knownvalue.StringExact(testLocation)},
							{Path: tfjsonpath.New("logs_enabled"), KnownValue: knownvalue.Bool(true)},
							{Path: tfjsonpath.New("metrics_enabled"), KnownValue: knownvalue.Bool(true)},
							{Path: tfjsonpath.New("dns_name"), KnownValue: knownvalue.NotNull()},
							// Instances
							{Path: tfjsonpath.New("instances").AtMapKey("count"), KnownValue: knownvalue.Int32Exact(1)},
							{Path: tfjsonpath.New("instances").AtMapKey("cores"), KnownValue: knownvalue.Int32Exact(1)},
							{Path: tfjsonpath.New("instances").AtMapKey("ram"), KnownValue: knownvalue.Int32Exact(4)},
							{Path: tfjsonpath.New("instances").AtMapKey("storage_size"), KnownValue: knownvalue.Int32Exact(10)},
							// Connection
							{Path: tfjsonpath.New("connections").AtMapKey("primary_instance_address"), KnownValue: knownvalue.StringExact("192.168.2.101/24")},
							// Backup
							{Path: tfjsonpath.New("backup").AtMapKey("location"), KnownValue: knownvalue.StringExact("eu-central-3")},
							{Path: tfjsonpath.New("backup").AtMapKey("retention_days"), KnownValue: knownvalue.Int32Exact(7)},
							// Maintenance window
							{Path: tfjsonpath.New("maintenance_window").AtMapKey("time"), KnownValue: knownvalue.StringExact("09:00:00")},
							{Path: tfjsonpath.New("maintenance_window").AtMapKey("day_of_the_week"), KnownValue: knownvalue.StringExact("Sunday")},
							// Credentials
							{Path: tfjsonpath.New("credentials").AtMapKey("username"), KnownValue: knownvalue.StringExact("dbadmin")},
							{Path: tfjsonpath.New("credentials").AtMapKey("database"), KnownValue: knownvalue.StringExact("my_database")},
						},
					),
				},
			},
			{
				Config: clusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Updated MariaDB V2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "backup.retention_days", "14"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "12:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Wednesday"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "false"),
				),
			},
			{
				ResourceName:      clusterResourceAddr,
				ImportStateIdFunc: clusterImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"credentials.password",
					"timeouts",
				},
			},
		},
	})
}

var clusterDataSourcesConfig = fmt.Sprintf(`
resource "ionoscloud_mariadb_cluster_v2" "test2" {
  location    = "%[1]s"
  name        = "tf-test2-mariadbv2"
  description = "Second test MariaDB V2 cluster"
  version     = "11.4"

  logs_enabled    = true
  metrics_enabled = true

  instances = {
    count        = 1
    cores        = 1
    ram          = 4
    storage_size = 10
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.2.102/24"
  }

  backup = {
    location       = "eu-central-3"
    retention_days = 7
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "dbadmin"
    password = "P@ssw0rd123!"
    database = "my_database"
  }
}

data "ionoscloud_mariadb_cluster_v2" "by_id" {
  id       = ionoscloud_mariadb_cluster_v2.test.id
  location = "%[1]s"
}

data "ionoscloud_mariadb_cluster_v2" "by_name" {
  name     = ionoscloud_mariadb_cluster_v2.test.name
  location = "%[1]s"
}

data "ionoscloud_mariadb_clusters_v2" "test" {
  location   = "%[1]s"
  name       = ionoscloud_mariadb_cluster_v2.test.name
  depends_on = [ionoscloud_mariadb_cluster_v2.test2]
}

data "ionoscloud_mariadb_clusters_v2" "all" {
  location   = "%[1]s"
  depends_on = [ionoscloud_mariadb_cluster_v2.test2]
}

data "ionoscloud_mariadb_backups_v2" "test" {
  location   = "%[1]s"
  cluster_id = ionoscloud_mariadb_cluster_v2.test.id
}

data "ionoscloud_mariadb_backup_locations_v2" "test" {
  location = "%[1]s"
}

data "ionoscloud_mariadb_versions_v2" "test" {
  location = "%[1]s"
}
`, testLocation)
