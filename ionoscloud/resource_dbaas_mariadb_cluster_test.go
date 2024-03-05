package ionoscloud

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mariadb "github.com/ionos-cloud/sdk-go-dbaas-maria"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func TestAccDBaaSMariaDBClusterBasic(t *testing.T) {
	var cluster mariadb.ClusterResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDBaaSMariaDBClusterDestroyCheck,
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
					resource.TestCheckResourceAttr(constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute, clusterConnectionsCidrValue),
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
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceById, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute, clusterConnectionsCidrValue),
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
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute, clusterConnectionsCidrValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSMariaDBClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
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
		_, apiResponse, err := client.GetCluster(ctx, rs.Primary.ID)
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking the destruction of MariaDB cluster with ID: %v, error: %w", rs.Primary.ID, err)
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

		foundCluster, _, err := client.GetCluster(ctx, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("an error occured while fetching MariaDB cluster with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const mariaDBClusterConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + datacenterResourceName + ` {
  name        = "mariadb_datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing MariaDB cluster"
}

resource ` + constant.LanResource + ` ` + lanResourceName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "mariadb_lan_example"
}

resource ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  ` + clusterVersionAttribute + ` = "` + clusterVersionValue + `"
  ` + clusterInstancesAttribute + ` = "` + clusterInstancesValue + `"
  ` + clusterCoresAttribute + ` = "` + clusterCoresValue + `"
  ` + clusterRamAttribute + ` = "` + clusterRamValue + `"
  ` + clusterStorageSizeAttribute + ` = "` + clusterStorageSizeValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + clusterDisplayNameValue + `"
  ` + connections + `
  ` + maintenanceWindow + `
  ` + credentials + `
}

resource ` + constant.RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const mariaDBClusterDataSourceMatchId = mariaDBClusterConfigBasic + `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
	id = ` + constant.DBaaSMariaDBClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`

const mariaDBClusterDataSourceMatchName = mariaDBClusterConfigBasic + `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
	display_name	= "` + clusterDisplayNameValue + `"
}
`

const mariaDBClusterDataSourceWrongName = `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name = "wrong_name"
}
`

// Any valid UUID can be used here since there is no cluster created, so no cluster will be found
const mariaDBClusterDataSourceWrongId = `
data ` + constant.DBaaSMariaDBClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id = "178d6a7d-5ed4-44de-88f0-27f1182d8ae8"
}
`

// Internal resources
const connections = clusterConnectionsAttribute + `{
	` + clusterConnectionsDatacenterIDAttribute + ` = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
    ` + clusterConnectionsLanIDAttribute + ` = ` + constant.LanResource + `.` + lanResourceName + `.id
	` + clusterConnectionsCidrAttribute + ` = "` + clusterConnectionsCidrValue + `"
}`

const maintenanceWindow = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeValue + `"
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
	clusterInstancesValue           = "1"
	clusterCoresValue               = "4"
	clusterRamValue                 = "4"
	clusterStorageSizeValue         = "10"
	clusterConnectionsCidrValue     = "192.168.1.100/24"
	clusterDisplayNameValue         = constant.DBaaSClusterTestResource
	clusterCredentialsUsernameValue = "username"
	datacenterResourceName          = "datacenter_example"
	lanResourceName                 = "lan_example"
)
