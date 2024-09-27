//go:build all || dbaas || psql
// +build all dbaas psql

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccDBaaSPgSqlClusterBasic(t *testing.T) {
	var dbaasCluster psql.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.enabled", "false"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.pool_mode", "session"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".cluster_password", "result"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
					resource.TestCheckResourceAttrSet(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "dns_name"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "display_name", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "instances", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "cores", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "ram", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "storage_size", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "storage_type", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connection_pooler.0.enabled", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.enabled"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connection_pooler.0.pool_mode", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.pool_mode"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.datacenter_id", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.lan_id", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.cidr", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "location", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "backup_location", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "display_name", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.time", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "credentials.username", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "credentials.password", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.password"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "dns_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "dns_name", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "dns_name"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "display_name", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "instances", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "cores", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "ram", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "storage_size", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "storage_type", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connection_pooler.0.enabled", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.enabled"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connection_pooler.0.pool_mode", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.pool_mode"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.datacenter_id", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.lan_id", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.cidr", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "location", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "backup_location", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "display_name", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.time", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "credentials.username", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "credentials.password", constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config:      testAccDataSourceDBaaSPgSqlClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DBaaS cluster found with the specified name"),
			},
			{
				PreConfig: sleepUntilBackupIsReady,
				Config:    testAccDataSourceDbaasPgSqlClusterBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_backups.0.cluster_id", constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_id"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_backups.0.size"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_backups.0.location"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_backups.0.version"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.PsqlBackupsResource+"."+constant.PsqlBackupsTest, "cluster_backups.0.is_active"),
					utils.TestNotEmptySlice(constant.PsqlBackupsResource, "cluster_backups.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlVersionsByClusterId,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.PsqlVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlAllVersions,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.PsqlVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.enabled", "true"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.pool_mode", "transaction"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example_update", "id"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example_update", "id"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.101/24"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".cluster_password", "result"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.enabled", "true"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connection_pooler.0.pool_mode", "transaction"),
					resource.TestCheckNoResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.%"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".cluster_password", "result"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				// we need this as a separate test because the psql cluster needs to be deleted first
				// in order to be able to delete the associated lan after
				// otherwise we get 'Access Denied: Lan 1 is delete-protected by DBAAS'
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveDBaaS,
			},
		},
	})
}

// sleepUntilBackupIsReady waits 60s until backup is ready
func sleepUntilBackupIsReady() {
	time.Sleep(60 * time.Second)
}

func TestAccDBaaSPgSqlClusterAdditionalParameters(t *testing.T) {
	var dbaasCluster psql.ClusterResponse
	t.Skip()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccFromBackup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(constant.PsqlClusterResource+"."+constant.DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS")),
			},
		},
	})
}

func testAccCheckDbaasPgSqlClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).PsqlClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PsqlBackupsResource {
			continue
		}

		_, apiResponse, err := client.GetCluster(ctx, rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of psql cluster %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("PgSQL cluster with ID: %v still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDbaasPgSqlClusterExists(n string, cluster *psql.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).PsqlClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		defer cancel()

		foundCluster, _, err := client.GetCluster(ctx, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching PgSQL cluster with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccCheckDbaasPgSqlClusterConfigBasic = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 1
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  connection_pooler {
    enabled = false
    pool_mode = "session"
  }
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  backup_location = "de"
  display_name = "` + constant.DBaaSClusterTestResource + `"
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.cluster_password.result
  }
  synchronization_mode = "ASYNCHRONOUS"
}

resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccCheckDbaasPgSqlClusterConfigUpdate = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.LanResource + ` "lan_example_update" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}


resource ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 2
  cores              = 2
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  connection_pooler {
	enabled = true
	pool_mode = "transaction"
  }
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example_update.id 
    lan_id          =  ` + constant.LanResource + `.lan_example_update.id     
    cidr            =  "192.168.1.101/24"
  }
  location = ` + constant.DatacenterResource + `.datacenter_example_update.location
  backup_location = "de"
  display_name = "` + constant.UpdatedResources + `"
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "10:00:00"
  }
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.cluster_password.result
  }
  synchronization_mode = "ASYNCHRONOUS"
}

resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.LanResource + ` "lan_example_update" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}


resource ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 2
  cores              = 2
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  location = ` + constant.DatacenterResource + `.datacenter_example_update.location
  backup_location = "de"
  display_name = "` + constant.UpdatedResources + `"
  connection_pooler {
	enabled = true
	pool_mode = "transaction"
  } 
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "10:00:00"
  }
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.cluster_password.result
  }
  synchronization_mode = "ASYNCHRONOUS"
}

resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
const testAccCheckDbaasPgSqlClusterConfigUpdateRemoveDBaaS = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.LanResource + ` "lan_example_update" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}
`

const testAccFromBackup = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing psql cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 1
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  display_name = "` + constant.DBaaSClusterTestResource + `"
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.cluster_password.result
  }
  synchronization_mode = "ASYNCHRONOUS"
  from_backup {
	backup_id = "f767c6e5-747c-11ec-9bb6-4aa52b3d55f1-4oymiqu-12"
    recovery_target_time = "2022-01-13T16:27:42Z"
  }
}
resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id	= ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchName = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + constant.DBaaSClusterTestResource + `"
}
`

const testAccDataSourceDBaaSPgSqlClusterWrongNameError = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`

const testAccDataSourceDbaasPgSqlClusterBackups = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlBackupsResource + ` ` + constant.PsqlBackupsTest + ` {
	cluster_id = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDbaasPgSqlAllVersions = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlVersionsResource + ` ` + constant.PsqlVersionsTest + ` {
}
`

const testAccDataSourceDbaasPgSqlVersionsByClusterId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + constant.PsqlVersionsResource + ` ` + constant.PsqlVersionsTest + ` {
	cluster_id = ` + constant.PsqlClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`
