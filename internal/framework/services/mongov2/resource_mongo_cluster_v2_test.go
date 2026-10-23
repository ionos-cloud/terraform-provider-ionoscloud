//go:build all || dbaas || mongov2

package mongov2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccMongoV2Cluster(t *testing.T) {
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
					resource.TestCheckResourceAttr(clusterResourceAddr, "name", "tf-test-mongov2"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Test MongoDB V2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "edition", "ENTERPRISE"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "version", "7.0"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "true"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connection_uri"),

					// Instances (ENTERPRISE)
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.count", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.cores", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.ram", "8"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instances.storage_size", "10"),

					// Connections
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connections.lan_id"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.instance_addresses.#", "1"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "connections.instance_addresses.0", "192.168.1.100/24"),

					// Snapshot config
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.location", "eu-central-3"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.retention_days", "7"),

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
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "edition", clusterResourceAddr, "edition"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "version", clusterResourceAddr, "version"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "instances.cores", clusterResourceAddr, "instances.cores"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "connections.datacenter_id", clusterResourceAddr, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "snapshot.retention_days", clusterResourceAddr, "snapshot.retention_days"),
					resource.TestCheckResourceAttrPair(clusterDSByIDAddr, "credentials.username", clusterResourceAddr, "credentials.username"),

					// by name
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clusterDSByNameAddr, "name", clusterResourceAddr, "name"),

					// clusters list with name filter
					resource.TestCheckResourceAttr(clustersDSAddr, "items.#", "1"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.id", clusterResourceAddr, "id"),
					resource.TestCheckResourceAttrPair(clustersDSAddr, "items.0.edition", clusterResourceAddr, "edition"),

					// clusters list without filter — both test clusters appear
					resource.TestCheckResourceAttr(clustersAllDSAddr, "items.#", "2"),

					// snapshot locations
					resource.TestCheckResourceAttrSet(snapshotLocationsAddr, "items.#"),

					// templates filtered by edition
					resource.TestCheckResourceAttrSet(templatesDSAddr, "items.#"),

					// versions
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.#"),
					resource.TestCheckResourceAttrSet(versionsDSAddr, "items.0.version"),
				),
			},
			// List filtered by location — avoids querying other regional endpoints.
			{
				Query: true,
				Config: `list "ionoscloud_mongo_cluster_v2" "test" {
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
				Config: `list "ionoscloud_mongo_cluster_v2" "test" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = "tf-test-mongov2" },
      { field_name = "location", field_value = "` + testLocation + `" },
    ]
  }
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 1),
				},
			},
			// Filter by name + wrong (but reachable) location: must return 0
			{
				Query: true,
				Config: `list "ionoscloud_mongo_cluster_v2" "test" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = "tf-test-mongov2" },
      { field_name = "location", field_value = "` + testLocationOther + `" },
    ]
  }
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(clusterResourceAddr, 0),
				},
			},
			// Import
			{
				ResourceName:      clusterResourceAddr,
				ImportState:       true,
				ImportStateVerify: true,
				// password is write-only and never returned by the API.
				ImportStateVerifyIgnore: []string{"credentials.password", "timeouts"},
				ImportStateIdFunc:       clusterImportStateID,
			},
			// Update
			{
				Config: clusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "description", "Updated MongoDB V2 cluster"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "snapshot.retention_days", "14"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.time", "12:00:00"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "maintenance_window.day_of_the_week", "Wednesday"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "logs_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "metrics_enabled", "false"),
				),
			},
		},
	})
}

func TestAccMongoV2ClusterBusiness(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: businessCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkClusterExists(clusterResourceAddr),
					resource.TestCheckResourceAttr(clusterResourceAddr, "edition", "BUSINESS"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "instance_count", "1"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "instance_template_id"),
					// ENTERPRISE-only attrs must be absent/known for BUSINESS.
					resource.TestCheckNoResourceAttr(clusterResourceAddr, "bi_connector_enabled"),
					resource.TestCheckNoResourceAttr(clusterResourceAddr, "instances.count"),
					resource.TestCheckResourceAttrSet(clusterResourceAddr, "connection_uri"),
					resource.TestCheckResourceAttr(clusterResourceAddr, "credentials.username", "dbadmin"),
				),
			},
			// instance_template_id is immutable → changing it must produce a (replace) plan, not update.
			{
				Config:             businessCreateConfigOtherTemplate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Import round-trip.
			{
				ResourceName:            clusterResourceAddr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials.password", "timeouts"},
				ImportStateIdFunc:       clusterImportStateID,
			},
		},
	})
}

const (
	testLocation      = "gb/bhx"
	testLocationOther = "us/ewr"

	clusterResourceAddr   = constant.MongoV2ClusterResource + ".test"
	clusterDSByIDAddr     = "data." + constant.MongoV2ClusterDataSource + ".by_id"
	clusterDSByNameAddr   = "data." + constant.MongoV2ClusterDataSource + ".by_name"
	clustersDSAddr        = "data." + constant.MongoV2ClustersDataSource + ".test"
	clustersAllDSAddr     = "data." + constant.MongoV2ClustersDataSource + ".all"
	templatesDSAddr       = "data." + constant.MongoV2TemplatesDataSource + ".test"
	snapshotLocationsAddr = "data." + constant.MongoV2SnapshotLocationsDataSource + ".test"
	versionsDSAddr        = "data." + constant.MongoV2VersionsDataSource + ".test"
)

// infraConfig creates the datacenter and LAN needed by the cluster.
var infraConfig = fmt.Sprintf(`
resource "ionoscloud_datacenter" "test" {
  name     = "tf-test-mongov2"
  location = "%[1]s"
}

resource "ionoscloud_lan" "test" {
  datacenter_id = ionoscloud_datacenter.test.id
  public        = false
  name          = "tf-test-mongov2"
}
`, testLocation)

var clusterCreateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_mongo_cluster_v2" "test" {
  location    = "%[1]s"
  name        = "tf-test-mongov2"
  description = "Test MongoDB V2 cluster"
  edition     = "ENTERPRISE"
  version     = "7.0"

  logs_enabled    = true
  metrics_enabled = true

  instances = {
    count        = 1
    cores        = 1
    ram          = 8
    storage_size = 10
  }

  connections = {
    datacenter_id      = ionoscloud_datacenter.test.id
    lan_id             = ionoscloud_lan.test.id
    instance_addresses = ["192.168.1.100/24"]
  }

  snapshot = {
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
`, testLocation)

var clusterUpdateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_mongo_cluster_v2" "test" {
  location    = "%[1]s"
  name        = "tf-test-mongov2"
  description = "Updated MongoDB V2 cluster"
  edition     = "ENTERPRISE"
  version     = "7.0"

  logs_enabled    = false
  metrics_enabled = false

  instances = {
    count        = 1
    cores        = 1
    ram          = 8
    storage_size = 10
  }

  connections = {
    datacenter_id      = ionoscloud_datacenter.test.id
    lan_id             = ionoscloud_lan.test.id
    instance_addresses = ["192.168.1.100/24"]
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 14
  }

  maintenance_window = {
    time            = "12:00:00"
    day_of_the_week = "Wednesday"
  }

  credentials = {
    username = "dbadmin"
    password = "P@ssw0rd123!"
    database = "my_database"
  }
}
`, testLocation)

var clusterDataSourcesConfig = fmt.Sprintf(`
resource "ionoscloud_mongo_cluster_v2" "test2" {
  location    = "%[1]s"
  name        = "tf-test2-mongov2"
  description = "Second test MongoDB V2 cluster"
  edition     = "ENTERPRISE"
  version     = "7.0"

  instances = {
    count        = 1
    cores        = 1
    ram          = 8
    storage_size = 10
  }

  connections = {
    datacenter_id      = ionoscloud_datacenter.test.id
    lan_id             = ionoscloud_lan.test.id
    instance_addresses = ["192.168.1.101/24"]
  }

  snapshot = {
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

data "ionoscloud_mongo_cluster_v2" "by_id" {
  id       = ionoscloud_mongo_cluster_v2.test.id
  location = "%[1]s"
}

data "ionoscloud_mongo_cluster_v2" "by_name" {
  name     = ionoscloud_mongo_cluster_v2.test.name
  location = "%[1]s"
}

data "ionoscloud_mongo_clusters_v2" "test" {
  location   = "%[1]s"
  name       = ionoscloud_mongo_cluster_v2.test.name
  depends_on = [ionoscloud_mongo_cluster_v2.test2]
}

data "ionoscloud_mongo_clusters_v2" "all" {
  location   = "%[1]s"
  depends_on = [ionoscloud_mongo_cluster_v2.test2]
}

data "ionoscloud_mongo_templates_v2" "test" {
  location = "%[1]s"
  edition  = "ENTERPRISE"
}

data "ionoscloud_mongo_snapshot_locations_v2" "test" {
  location = "%[1]s"
}

data "ionoscloud_mongo_versions_v2" "test" {
  location = "%[1]s"
}
`, testLocation)

// businessTemplate is a BUSINESS instance-templates data source used by the BUSINESS cluster configs.
var businessTemplate = fmt.Sprintf(`
data "ionoscloud_mongo_templates_v2" "business" {
  location = "%[1]s"
  edition  = "BUSINESS"
}
`, testLocation)

// mongoBusinessCluster renders a BUSINESS cluster using the template at the given index.
// NOTE: the "other template" step assumes at least two BUSINESS templates exist.
func mongoBusinessCluster(templateIndex int) string {
	return businessTemplate + fmt.Sprintf(`
resource "ionoscloud_mongo_cluster_v2" "test" {
  location             = "%[1]s"
  name                 = "tf-test-mongov2-business"
  description          = "Business MongoDB V2 cluster"
  edition              = "BUSINESS"
  version              = "7.0"
  instance_template_id = data.ionoscloud_mongo_templates_v2.business.items[%[2]d].id
  instance_count       = 1

  connections = {
    datacenter_id      = ionoscloud_datacenter.test.id
    lan_id             = ionoscloud_lan.test.id
    instance_addresses = ["192.168.1.120/24"]
  }

  snapshot = {
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
`, testLocation, templateIndex)
}

var businessCreateConfig = infraConfig + mongoBusinessCluster(0)
var businessCreateConfigOtherTemplate = infraConfig + mongoBusinessCluster(1)

// checkClusterExists checks that the cluster resource exists in the API.
func checkClusterExists(resourceAddr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAddr)
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewMongoDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create MongoDB V2 client: %w", err)
		}
		_, _, err = client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching MongoDB V2 cluster %s: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

// checkClusterDestroy checks that all MongoDB V2 cluster resources have been destroyed.
func checkClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.MongoV2ClusterResource {
			continue
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewMongoDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create MongoDB V2 client: %w", err)
		}
		_, apiResponse, err := client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || !apiResponse.HttpNotFound() {
				return fmt.Errorf("error checking cluster %s destruction: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("MongoDB V2 cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// clusterImportStateID returns the import ID in "location:cluster_id" format.
func clusterImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.MongoV2ClusterResource {
			continue
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.ID), nil
	}
	return "", fmt.Errorf("no %s resource found in state", constant.MongoV2ClusterResource)
}
