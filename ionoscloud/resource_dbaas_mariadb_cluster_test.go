//go:build all || dbaas || mariadb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func TestAccDBaaSMariaDBClusterBasic(t *testing.T) {
	var cluster mariadb.ClusterResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.11.1",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDBaaSMariaDBClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: mariaDBClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSMariaDBClusterExists(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, &cluster),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterVersionAttribute, clusterVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterInstancesAttribute, clusterInstancesValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterCoresAttribute, clusterCoresValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterRamAttribute, clusterRamValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterStorageSizeAttribute, clusterStorageSizeValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterDisplayNameAttribute, clusterDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsUsernameAttribute, clusterCredentialsUsernameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsPasswordAttribute, constant.RandomPassword+".cluster_password", "result"),
				),
			},
			{
				Config: mariaDBClusterDataSourceMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterVersionAttribute, clusterVersionValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterInstancesAttribute, clusterInstancesValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterCoresAttribute, clusterCoresValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterRamAttribute, clusterRamValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterStorageSizeAttribute, clusterStorageSizeValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterDisplayNameAttribute, clusterDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
				),
			},
			{
				Config: mariaDBClusterDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterVersionAttribute, clusterVersionValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterInstancesAttribute, clusterInstancesValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterCoresAttribute, clusterCoresValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterRamAttribute, clusterRamValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterStorageSizeAttribute, clusterStorageSizeValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterDisplayNameAttribute, clusterDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
				),
			},
			{
				Config: mariaDBBackupsDataSourceMatchClusterID,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.DBaaSMariaDBBackupsDataSource, "backups.#"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSMariaDBBackupsDataSource+"."+constant.DBaasMariaDBBackupsDataSourceName, "backups.0.cluster_id", constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, "id"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DBaaSMariaDBBackupsDataSource+"."+constant.DBaasMariaDBBackupsDataSourceName, "backups.0.earliest_recovery_target_time"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DBaaSMariaDBBackupsDataSource+"."+constant.DBaasMariaDBBackupsDataSourceName, "backups.0.size"),
					utils.TestNotEmptySlice(constant.DBaaSMariaDBBackupsDataSource, "backups.0.base_backups.#"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DBaaSMariaDBBackupsDataSource+"."+constant.DBaasMariaDBBackupsDataSourceName, "backups.0.base_backups.0.size"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DBaaSMariaDBBackupsDataSource+"."+constant.DBaasMariaDBBackupsDataSourceName, "backups.0.base_backups.0.created"),
				),
			},
			{
				Config:      mariaDBClusterDataSourceWrongName,
				ExpectError: regexp.MustCompile("no MariaDB cluster found with the specified display name"),
			},
			{
				Config:      mariaDBClusterDataSourceWrongId,
				ExpectError: regexp.MustCompile("an error occurred while fetching the MariaDB cluster with ID"),
			},
			{
				Config: mariaDBClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSMariaDBClusterExists(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, &cluster),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterVersionAttribute, clusterVersionUpdatedValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterInstancesAttribute, clusterInstancesUpdatedValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterCoresAttribute, clusterCoresUpdatedValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterRamAttribute, clusterRamUpdatedValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterStorageSizeAttribute, clusterStorageSizeUpdatedValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeUpdateValue),
				),
			},
			{
				Config:      mariaDBClusterInvalidUpdateConfig,
				ExpectError: regexp.MustCompile("downgrade is not supported from"),
			},
		},
	})
}

func testAccCheckDBaaSMariaDBClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).MariaDBClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaaSMariaDBClusterResource {
			continue
		}
		_, apiResponse, err := client.GetCluster(ctx, rs.Primary.ID, rs.Primary.Attributes[clusterLocationAttribute])
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of MariaDB cluster with ID: %v, error: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("MariaDB cluster with ID: %v still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDBaaSMariaDBClusterExists(n string, cluster *mariadb.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).MariaDBClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundCluster, _, err := client.GetCluster(ctx, rs.Primary.ID, rs.Primary.Attributes[clusterLocationAttribute])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching MariaDB cluster with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const mariaDBTestInfraConfig = `
resource ` + constant.DatacenterResource + ` ` + datacenterResourceName + ` {
  name        = "mariadb_datacenter_example"
  location    = "fr/par"
  description = "Datacenter for testing MariaDB cluster"
}

resource ` + constant.LanResource + ` ` + lanResourceName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "mariadb_lan_example"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name                    = "example"
  datacenter_id           = ionoscloud_datacenter.datacenter_example.id
  cores                   = 2
  ram                     = 2048
  availability_zone       = "ZONE_1"
  cpu_family              = "INTEL_SKYLAKE"
  image_name              = "rockylinux-8-GenericCloud-20240528"
  image_password          = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name                  = "example"
    size                  = 20
    disk_type             = "SSD Standard"
  }
  nic {
    lan                   = ionoscloud_lan.lan_example.id
    name                  = "example"
    dhcp                  = true
  }
}

locals {
 prefix                   = format("%s/%s", ionoscloud_server.test_server.nic[0].ips[0], "24")
 database_ip              = cidrhost(local.prefix, 1)
 database_ip_cidr         = format("%s/%s", local.database_ip, "24")
}

resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
` + ServerImagePassword

const mariaDBClusterConfigBasic = mariaDBTestInfraConfig + `
resource ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  ` + clusterVersionAttribute + ` = "` + clusterVersionValue + `"
  ` + clusterInstancesAttribute + ` = "` + clusterInstancesValue + `"
  ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
  ` + clusterCoresAttribute + ` = "` + clusterCoresValue + `"
  ` + clusterRamAttribute + ` = "` + clusterRamValue + `"
  ` + clusterStorageSizeAttribute + ` = "` + clusterStorageSizeValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + clusterDisplayNameValue + `"
  ` + connections + `
  ` + maintenanceWindow + `
  ` + credentials + `
}

# Wait few seconds after cluster creation so the backups can be properly retrieved
resource "time_sleep" "wait_10_minutes" {
  depends_on = [` + constant.DBaaSMariaDBClusterResource + `.` + constant.DBaaSClusterTestResource + `]
  create_duration = "600s"
}
`

const mariaDBClusterUpdateConfig = mariaDBTestInfraConfig + `
resource ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  ` + clusterVersionAttribute + ` = "` + clusterVersionUpdatedValue + `"
  ` + clusterInstancesAttribute + ` = "` + clusterInstancesUpdatedValue + `"
  ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
  ` + clusterCoresAttribute + ` = "` + clusterCoresUpdatedValue + `"
  ` + clusterRamAttribute + ` = "` + clusterRamUpdatedValue + `"
  ` + clusterStorageSizeAttribute + ` = "` + clusterStorageSizeUpdatedValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + clusterDisplayNameValue + `"
  ` + connections + `
  ` + maintenanceWindowUpdated + `
  ` + credentials + `
}
`

// This update is invalid because version downgrade is not allowed.
const mariaDBClusterInvalidUpdateConfig = mariaDBTestInfraConfig + `
resource ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  
  ` + clusterVersionAttribute + ` = "10.1"
  ` + clusterInstancesAttribute + ` = "` + clusterInstancesUpdatedValue + `"
  ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
  ` + clusterCoresAttribute + ` = "` + clusterCoresUpdatedValue + `"
  ` + clusterRamAttribute + ` = "` + clusterRamUpdatedValue + `"
  ` + clusterStorageSizeAttribute + ` = "` + clusterStorageSizeUpdatedValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + clusterDisplayNameValue + `"
  ` + connections + `
  ` + maintenanceWindowUpdated + `
  ` + credentials + `
}
`

const mariaDBClusterDataSourceMatchId = mariaDBClusterConfigBasic + `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
	id = ` + constant.DBaaSMariaDBClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
    ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
}

`
const mariaDBClusterDataSourceMatchName = mariaDBClusterConfigBasic + `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
	display_name	= "` + clusterDisplayNameValue + `"
    ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
}
`

const mariaDBBackupsDataSourceMatchClusterID = mariaDBClusterConfigBasic + `
data ` + constant.DBaaSMariaDBBackupsDataSource + ` ` + constant.DBaasMariaDBBackupsDataSourceName + ` {
	cluster_id = ` + constant.DBaaSMariaDBClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
	` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
    # Use the previously created 'time' resource to delay information retrieval for the data source
	depends_on = [time_sleep.wait_10_minutes]
}
`
const mariaDBClusterDataSourceWrongName = `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name = "wrong_name"
  ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
}
`

// Any valid UUID can be used here since there is no cluster created, so no cluster will be found
const mariaDBClusterDataSourceWrongId = `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id = "178d6a7d-5ed4-44de-88f0-27f1182d8ae8"
  ` + clusterLocationAttribute + ` = "` + clusterLocationValue + `"
}
`

// Internal resources
const connections = clusterConnectionsAttribute + `{
	` + clusterConnectionsDatacenterIDAttribute + ` = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
    ` + clusterConnectionsLanIDAttribute + ` = ` + constant.LanResource + `.` + lanResourceName + `.id
	` + clusterConnectionsCidrAttribute + ` = ` + clusterConnectionsCidrValue + `
}`

const maintenanceWindow = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeValue + `"
}`

const maintenanceWindowUpdated = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekUpdateValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeUpdateValue + `"
}`

const credentials = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameValue + `"
	` + clusterCredentialsPasswordAttribute + ` = ` + constant.RandomPassword + `.cluster_password.result
}`

// Attributes
const clusterVersionAttribute = "mariadb_version"

// Values
const (
	clusterVersionValue             = "10.6"
	clusterVersionUpdatedValue      = "10.11"
	clusterInstancesValue           = "1"
	clusterInstancesUpdatedValue    = "2"
	clusterLocationValue            = "fr/par"
	clusterCoresValue               = "4"
	clusterCoresUpdatedValue        = "5"
	clusterRamValue                 = "4"
	clusterRamUpdatedValue          = "5"
	clusterStorageSizeValue         = "10"
	clusterStorageSizeUpdatedValue  = "11"
	clusterConnectionsCidrValue     = "local.database_ip_cidr"
	clusterDisplayNameValue         = constant.DBaaSClusterTestResource
	clusterCredentialsUsernameValue = "username"
	datacenterResourceName          = "datacenter_example"
	lanResourceName                 = "lan_example"
)
