//go:build all || dbaas || mariadbv2

package mariadbv2_test

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

	clusterResourceAddr = constant.MariaDBV2ClusterResource + ".test"
	clusterDSByIDAddr   = "data." + constant.MariaDBV2ClusterDataSource + ".by_id"
	clusterDSByNameAddr = "data." + constant.MariaDBV2ClusterDataSource + ".by_name"
	clustersDSAddr      = "data." + constant.MariaDBV2ClustersDataSource + ".test"
	clustersAllDSAddr   = "data." + constant.MariaDBV2ClustersDataSource + ".all"
	backupsDSAddr       = "data." + constant.MariaDBV2BackupsDataSource + ".test"
	backupLocationsAddr = "data." + constant.MariaDBV2BackupLocationsDataSource + ".test"
	versionsDSAddr      = "data." + constant.MariaDBV2VersionsDataSource + ".test"
)

// infraConfig creates the datacenter and LAN needed by the cluster.
var infraConfig = fmt.Sprintf(`
resource "ionoscloud_datacenter" "test" {
  name     = "tf-test-mariadbv2"
  location = "%[1]s"
}

resource "ionoscloud_lan" "test" {
  datacenter_id = ionoscloud_datacenter.test.id
  public        = false
  name          = "tf-test-mariadbv2"
}
`, testLocation)

var clusterCreateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_mariadb_cluster_v2" "test" {
  location    = "%[1]s"
  name        = "tf-test-mariadbv2"
  description = "Test MariaDB V2 cluster"
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
    primary_instance_address = "192.168.2.101/24"
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
`, testLocation)

var clusterUpdateConfig = infraConfig + fmt.Sprintf(`
resource "ionoscloud_mariadb_cluster_v2" "test" {
  location    = "%[1]s"
  name        = "tf-test-mariadbv2"
  description = "Updated MariaDB V2 cluster"
  version     = "11.4"

  logs_enabled    = false
  metrics_enabled = false

  instances = {
    count        = 1
    cores        = 1
    ram          = 4
    storage_size = 10
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.test.id
    lan_id                   = ionoscloud_lan.test.id
    primary_instance_address = "192.168.2.101/24"
  }

  backup = {
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

// checkClusterExists checks that the cluster resource exists in the API.
func checkClusterExists(resourceAddr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAddr)
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewMariaDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create MariaDB V2 client: %w", err)
		}
		_, _, err = client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching MariaDB V2 cluster %s: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

// checkClusterDestroy checks that all MariaDB V2 cluster resources have been destroyed.
func checkClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.MariaDBV2ClusterResource {
			continue
		}
		location := rs.Primary.Attributes["location"]
		client, err := acctest.NewTestBundleClientFromEnv().NewMariaDBV2Client(context.Background(), location)
		if err != nil {
			return fmt.Errorf("failed to create MariaDB V2 client: %w", err)
		}
		_, apiResponse, err := client.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || !apiResponse.HttpNotFound() {
				return fmt.Errorf("error checking cluster %s destruction: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("MariaDB V2 cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// clusterImportStateID returns the import ID in "location:cluster_id" format.
func clusterImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.MariaDBV2ClusterResource {
			continue
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.ID), nil
	}
	return "", fmt.Errorf("no %s resource found in state", constant.MariaDBV2ClusterResource)
}
