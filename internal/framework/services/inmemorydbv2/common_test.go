//go:build all || dbaas || inmemorydbv2

package inmemorydbv2_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const (
	testLocation      = "de/txl"
	testLocationOther = "de/fra"

	clusterResourceAddr     = constant.InMemoryDBV2ClusterResource + ".test"
	clusterDSByIDAddr       = "data." + constant.InMemoryDBV2ClusterDataSource + ".by_id"
	clusterDSByNameAddr     = "data." + constant.InMemoryDBV2ClusterDataSource + ".by_name"
	clustersDSAddr          = "data." + constant.InMemoryDBV2ClustersDataSource + ".test"
	clustersAllDSAddr       = "data." + constant.InMemoryDBV2ClustersDataSource + ".all"
	snapshotsDSAddr         = "data." + constant.InMemoryDBV2SnapshotsDataSource + ".test"
	snapshotsAllDSAddr      = "data." + constant.InMemoryDBV2SnapshotsDataSource + ".all"
	snapshotLocationsDSAddr = "data." + constant.InMemoryDBV2SnapshotLocationsDataSource + ".test"
	versionsDSAddr          = "data." + constant.InMemoryDBV2VersionsDataSource + ".test"
)

// --- Shared Infrastructure ---

// infraConfig creates the datacenter and LAN needed by the cluster.
var infraConfig = fmt.Sprintf(`
resource "ionoscloud_datacenter" "test" {
  name     = "tf-test-inmemorydbv2"
  location = "%[1]s"
}

resource "ionoscloud_lan" "test" {
  datacenter_id = ionoscloud_datacenter.test.id
  public        = false
  name          = "tf-test-inmemorydbv2"
}
`, testLocation)

// --- Cluster configs ---

var clusterCreateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_inmemorydb_cluster_v2" "test" {
  location         = "%[1]s"
  name             = "tf-test-inmemorydbv2"
  description      = "Test InMemoryDB v2 cluster"
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
    primary_instance_address = "192.168.2.101/24"
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
`, testLocation)

var clusterUpdateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_inmemorydb_cluster_v2" "test" {
  location         = "%[1]s"
  name             = "tf-test-inmemorydbv2"
  description      = "Test InMemoryDB v2 cluster"
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
    primary_instance_address = "192.168.2.101/24"
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 14
    snapshot_hours = [6, 18]
  }

  maintenance_window = {
    time            = "12:00:00"
    day_of_the_week = "Wednesday"
  }

  credentials = {
    username = "cacheadmin"
    password = {
      algorithm = "SHA-256"
      hash      = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
    }
  }
}
`, testLocation)

// --- Helper Functions ---

func checkClusterExists(resourceAddr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAddr)
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewInMemoryDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create InMemoryDB v2 client: %w", err)
		}
		_, _, err = client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching InMemoryDB v2 cluster %s: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func checkClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.InMemoryDBV2ClusterResource {
			continue
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewInMemoryDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create InMemoryDB v2 client: %w", err)
		}
		_, apiResponse, err := client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || !apiResponse.HttpNotFound() {
				return fmt.Errorf("error checking cluster %s destruction: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("InMemoryDB v2 cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func clusterImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.InMemoryDBV2ClusterResource {
			continue
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.ID), nil
	}
	return "", fmt.Errorf("no %s resource found in state", constant.InMemoryDBV2ClusterResource)
}
