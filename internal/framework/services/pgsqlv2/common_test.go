//go:build all || dbaas || psqlv2

package pgsqlv2_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const (
	testLocation        = "de/txl"
	testLocationChanged = "de/fra"

	clusterResourceAddr  = constant.PsqlV2ClusterResource + ".test"
	clusterDSByIDAddr    = "data." + constant.PsqlV2ClusterDataSource + ".by_id"
	clusterDSByNameAddr  = "data." + constant.PsqlV2ClusterDataSource + ".by_name"
	clustersDSAddr       = "data." + constant.PsqlV2ClustersDataSource + ".test"
	backupsDSAddr        = "data." + constant.PsqlV2BackupsDataSource + ".test"
	versionsDSAddr       = "data." + constant.PsqlV2VersionsDataSource + ".test"
	backupLocationDSAddr = "data." + constant.PsqlV2BackupLocationDataSource + ".test"
)

// --- Shared Terraform Configs ---

// infraConfig creates the datacenter, LAN, and random password required by the PgSQL v2 cluster.
var infraConfig = fmt.Sprintf(`
resource "ionoscloud_datacenter" "test" {
  name     = "tf-test-pgsqlv2"
  location = "%[1]s"
}

resource "ionoscloud_lan" "test" {
  datacenter_id = ionoscloud_datacenter.test.id
  public        = false
  name          = "tf-test-pgsqlv2"
}

resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  min_special      = 1
  override_special = "@$!%%*?&"
}
`, testLocation)

// clusterCreateConfig creates a PgSQL v2 cluster with all attributes explicitly set.
var clusterCreateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_pg_cluster_v2" "test" {
  name              = "tf-test-pgsqlv2"
  description       = "Test PgSQL v2 cluster"
  version           = "17"
  location          = "%[1]s"
  backup_location   = "eu-central-3"
  replication_mode  = "ASYNCHRONOUS"
  connection_pooler = "DISABLED"
  logs_enabled      = true
  metrics_enabled   = true

  instances = {
    count        = 1
    cores        = 1
    ram          = 4
    storage_size = 10
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.1.100/24"
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "testuser"
    password = random_password.cluster_password.result
    database = "testdb"
  }
}
`, testLocation)

// clusterUpdateConfig updates all mutable attributes of the PgSQL v2 cluster.
var clusterUpdateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_pg_cluster_v2" "test" {
  name              = "tf-test-pgsqlv2-updated"
  description       = "Updated PgSQL v2 cluster"
  version           = "17"
  location          = "%[1]s"
  backup_location   = "eu-central-3"
  replication_mode  = "ASYNCHRONOUS"
  connection_pooler = "DISABLED"
  logs_enabled      = false
  metrics_enabled   = false

  instances = {
    count        = 2
    cores        = 2
    ram          = 4
    storage_size = 20
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.1.100/24"
  }

  maintenance_window = {
    time            = "12:00:00"
    day_of_the_week = "Wednesday"
  }

  credentials = {
    username = "testuser"
    password = random_password.cluster_password.result
    database = "testdb"
  }
}
`, testLocation)

// --- Shared Helper Functions ---

func checkClusterV2Exists(resourceAddr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAddr)
		}
		client, err := acctest.NewTestBundleClientFromEnv().NewPgSQLV2Client(context.Background(), testLocation)
		if err != nil {
			return fmt.Errorf("failed to create PgSQL v2 client: %w", err)
		}
		_, _, err = client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching PgSQL v2 cluster %s: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func checkClusterV2ExistsInLocation(resourceAddr, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAddr)
		}
		client, err := acctest.NewTestBundleClientFromEnv().NewPgSQLV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create PgSQL v2 client for location %s: %w", location, err)
		}
		_, _, err = client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching PgSQL v2 cluster %s in location %s: %w", rs.Primary.ID, location, err)
		}
		return nil
	}
}

func checkClusterV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PsqlV2ClusterResource {
			continue
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewPgSQLV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create PgSQL v2 client for location %s: %w", location, err)
		}
		_, apiResponse, err := client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || !apiResponse.HttpNotFound() {
				return fmt.Errorf("error checking PgSQL v2 cluster %s destruction: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("PgSQL v2 cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func pgClusterV2ImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PsqlV2ClusterResource {
			continue
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.ID), nil
	}
	return "", fmt.Errorf("no %s resource found in state", constant.PsqlV2ClusterResource)
}

// infraConfigNewLocation creates infra in a different location for the location-change test.
var infraConfigNewLocation = fmt.Sprintf(`
resource "ionoscloud_datacenter" "test" {
  name     = "tf-test-pgsqlv2"
  location = "%[1]s"
}

resource "ionoscloud_lan" "test" {
  datacenter_id = ionoscloud_datacenter.test.id
  public        = false
  name          = "tf-test-pgsqlv2"
}

resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  min_special      = 1
  override_special = "@$!%%*?&"
}
`, testLocationChanged)

// clusterLocationChangeConfig recreates the cluster in a different location.
// Because location has RequiresReplace, this should destroy the old cluster and create a new one.
var clusterLocationChangeConfig = infraConfigNewLocation + fmt.Sprintf(`
resource "ionoscloud_pg_cluster_v2" "test" {
  name              = "tf-test-pgsqlv2-updated"
  description       = "Updated PgSQL v2 cluster"
  version           = "17"
  location          = "%[1]s"
  backup_location   = "eu-central-3"
  replication_mode  = "ASYNCHRONOUS"
  connection_pooler = "DISABLED"
  logs_enabled      = false
  metrics_enabled   = false

  instances = {
    count        = 2
    cores        = 2
    ram          = 4
    storage_size = 20
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.1.100/24"
  }

  maintenance_window = {
    time            = "12:00:00"
    day_of_the_week = "Wednesday"
  }

  credentials = {
    username = "testuser"
    password = random_password.cluster_password.result
    database = "testdb"
  }
}
`, testLocationChanged)
