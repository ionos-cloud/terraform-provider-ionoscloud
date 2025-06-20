//go:build all || dbaas || inmemorydb
// +build all dbaas inmemorydb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDBaaSInMemoryDBReplicaSetBasic(t *testing.T) {
	var replicaSet inmemorydb.ReplicaSetRead
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDBaaSInMemoryDBReplicaSetDestroyCheck,
		// The tests contain multiple constants that are reused in other DBaaS tests, especially attributes like 'lan_id', 'datacenter_id' for
		// which there is no need to create new constants (there is a high probability that these attributes will remain the same in the future).
		Steps: []resource.TestStep{
			// This step tests a configuration that uses a hashed password.
			{
				Config: inMemoryDBReplicaSetConfigHashedPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSInMemoryDBReplicaSetExists(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValue,
						clusterRamAttribute:   replicaSetRAMValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetHashedPasswordAttribute+".*", map[string]string{
						replicaSetHashAttribute:      replicaSetHashValue,
						replicaSetAlgorithmAttribute: replicaSetAlgorithmValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			{
				Config: inMemoryDBReplicaSetConfigSetup,
			},
			// This step tests a configuration that uses the plain text password, this configuration
			// will be also used for 'update' tests.
			{
				Config: inMemoryDBReplicaSetConfigPlainTextPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSInMemoryDBReplicaSetExists(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValue,
						clusterRamAttribute:   replicaSetRAMValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsUsernameAttribute, clusterCredentialsUsernameAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetPlainTextPasswordAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			{
				Config: inMemoryDBReplicaSetDataSourceMatchID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSInMemoryDBReplicaSetExists(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetResourcesAttribute+".0."+clusterCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetResourcesAttribute+".0."+clusterRamAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+`.`+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
				),
			},
			{
				Config: inMemoryDBReplicaSetDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSInMemoryDBReplicaSetExists(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetResourcesAttribute+".0."+clusterCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetResourcesAttribute+".0."+clusterRamAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+`.`+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue)),
			},
			// This step tests for basic updates for different attributes.
			// TODO -- Check what fields can actually be updated since for some of them it seems
			// that we are receiving API errors.
			{
				Config: inMemoryDBReplicaSetConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSInMemoryDBReplicaSetExists(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasUpdateValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValueUpdate,
						clusterRamAttribute:   replicaSetRAMValueUpdate,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyUpdateValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterConnectionsAttribute+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckResourceAttr(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsUsernameAttribute, clusterCredentialsUsernameAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetPlainTextPasswordAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeUpdateValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekUpdateValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSInMemoryDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			{
				Config:      inMemoryDBReplicaSetDataSourceWrongName,
				ExpectError: regexp.MustCompile("no InMemoryDB replica set found with the specified display name"),
			},
			{
				Config:      inMemoryDBReplicaSetDataSourceWrongID,
				ExpectError: regexp.MustCompile("an error occurred while fetching the InMemoryDB replica set with ID"),
			},
		},
	})
}

func testAccCheckDBaaSInMemoryDBReplicaSetDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).InMemoryDBClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaaSInMemoryDBReplicaSetResource {
			continue
		}
		_, apiResponse, err := client.GetReplicaSet(ctx, rs.Primary.ID, rs.Primary.Attributes[clusterLocationAttribute])
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking the destruction of InMemoryDB replica set with ID: %v, error: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("InMemoryDB replica set with ID: %v still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDBaaSInMemoryDBReplicaSetExists(n string, replicaSet *inmemorydb.ReplicaSetRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).InMemoryDBClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundReplicaSet, _, err := client.GetReplicaSet(ctx, rs.Primary.ID, rs.Primary.Attributes[clusterLocationAttribute])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching InMemoryDB replica set with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if foundReplicaSet.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}
		replicaSet = &foundReplicaSet
		return nil
	}
}

// This configuration contains the resources that need to be created before creating a InMemoryDB replica set
// This configuration will be used when the API will be fixed.
const inMemoryDBReplicaSetConfigSetup = `
resource ` + constant.DatacenterResource + ` ` + datacenterResourceName + ` {
  name        = "in_memory_db_datacenter_example"
  location    = "` + replicaSetLocationValue + `"
  description = "Datacenter for testing InMemoryDB replica set"
}

resource ` + constant.LanResource + ` ` + lanResourceName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id 
  public        = false
  name          = "inMemorydb_lan_example"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name                    = "example"
  datacenter_id           = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
  cores                   = 2
  ram                     = 2048
  image_name             = "rocky:latest"
  image_password          = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name                  = "example"
    size                  = 20
    disk_type             = "SSD Standard"
  }
  nic {
    lan                   = ` + constant.LanResource + `.` + lanResourceName + `.id
    name                  = "example"
    dhcp                  = true
  }
}

locals {
 prefix                   = format("%s/%s", ` + constant.ServerResource + `.` + constant.ServerTestResource + `.nic[0].ips[0], "24")
 database_ip              = cidrhost(local.prefix, 1)
 database_ip_cidr         = format("%s/%s", local.database_ip, "24")
}

resource ` + constant.RandomPassword + ` "replicaset_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
` + ServerImagePassword

const inMemoryDBReplicaSetConfigHashedPassword = inMemoryDBReplicaSetConfigSetup + `
resource ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyValue + `"
  ` + resources + `
  ` + replicaSetConnections + `
  ` + replicaSetMaintenanceWindow + `
  ` + credentialsHashedPassword + `
}
`

const inMemoryDBReplicaSetConfigPlainTextPassword = inMemoryDBReplicaSetConfigSetup + `
resource ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyValue + `"
  ` + resources + `
  ` + replicaSetConnections + `
  ` + replicaSetMaintenanceWindow + `
  ` + credentialsPlainTextPassword + `
}
`

const inMemoryDBReplicaSetConfigUpdate = inMemoryDBReplicaSetConfigSetup + `
resource ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameUpdateValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasUpdateValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeUpdateValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyUpdateValue + `"
  ` + resourcesUpdate + `
  ` + replicaSetConnections + `
  ` + maintenanceWindowUpdate + `
  ` + credentialsPlainTextPassword + `
}
`

// Internal resources
const resources = replicaSetResourcesAttribute + `{
	` + clusterCoresAttribute + ` = "` + replicaSetCoresValue + `"
	` + clusterRamAttribute + ` = "` + replicaSetRAMValue + `"
}`

const resourcesUpdate = replicaSetResourcesAttribute + `{
	` + clusterCoresAttribute + ` = "` + replicaSetCoresValueUpdate + `"
	` + clusterRamAttribute + ` = "` + replicaSetRAMValueUpdate + `"
}`

// For testing data source match by ID
const inMemoryDBReplicaSetDataSourceMatchID = inMemoryDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByID + ` {
	id = ` + constant.DBaaSInMemoryDBReplicaSetResource + `.` + constant.DBaaSReplicaSetTestResource + `.id
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For testing data source match by name
const inMemoryDBReplicaSetDataSourceMatchName = inMemoryDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "` + replicaSetDisplayNameValue + `"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by name
const inMemoryDBReplicaSetDataSourceWrongName = inMemoryDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "wrong_name"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by ID - use a 0000 uuidv4
const inMemoryDBReplicaSetDataSourceWrongID = inMemoryDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSInMemoryDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByID + ` {
	id = "00000000-0000-0000-0000-000000000000"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

const replicaSetConnections = clusterConnectionsAttribute + `{
	` + clusterConnectionsDatacenterIDAttribute + ` = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
    ` + clusterConnectionsLanIDAttribute + ` = ` + constant.LanResource + `.` + lanResourceName + `.id
	` + clusterConnectionsCidrAttribute + ` = ` + replicaSetConnectionsCidrValue + `
}`

const replicaSetMaintenanceWindow = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeValue + `"
}`

const maintenanceWindowUpdate = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekUpdateValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeUpdateValue + `"
}`
const credentialsPlainTextPassword = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameAttribute + `"
	` + replicaSetPlainTextPasswordAttribute + ` = ` + constant.RandomPassword + `.replicaset_password.result
}`

const credentialsHashedPassword = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameAttribute + `"
	` + hashedPasswordObject + `
}`

const hashedPasswordObject = replicaSetHashedPasswordAttribute + `{
    ` + replicaSetHashAttribute + ` = "` + replicaSetHashValue + `"
	` + replicaSetAlgorithmAttribute + ` = "` + replicaSetAlgorithmValue + `"
}
`

// Attributes
const (
	replicaSetLocationAttribute          = "location"
	replicaSetVersionAttribute           = "version"
	replicaSetDNSNameAttribute           = "dns_name"
	replicaSetReplicasAttribute          = "replicas"
	replicaSetPersistenceModeAttribute   = "persistence_mode"
	replicaSetEvictionPolicyAttribute    = "eviction_policy"
	replicaSetResourcesAttribute         = "resources"
	replicaSetPlainTextPasswordAttribute = "plain_text_password"
	replicaSetHashedPasswordAttribute    = "hashed_password"
	replicaSetAlgorithmAttribute         = "algorithm"
	replicaSetHashAttribute              = "hash"
)

// Values
const (
	replicaSetLocationValue              = "gb/lhr"
	replicaSetLocationUpdateValue        = "de/txl"
	replicaSetVersionValue               = "7.2"
	replicaSetDisplayNameValue           = "test"
	replicaSetDisplayNameUpdateValue     = "UpdatedReplicaSet"
	replicaSetReplicasValue              = "3"
	replicaSetReplicasUpdateValue        = "5"
	replicaSetPersistenceModeValue       = "RDB"
	replicaSetPersistenceModeUpdateValue = "AOF"
	replicaSetEvictionPolicyValue        = "allkeys-lru"
	replicaSetEvictionPolicyUpdateValue  = "allkeys-lru"
	replicaSetCoresValue                 = "4"
	replicaSetCoresValueUpdate           = "6"
	replicaSetRAMValue                   = "4"
	replicaSetRAMValueUpdate             = "6"
	replicaSetConnectionsCidrValue       = "local.database_ip_cidr"
	replicaSetHashValue                  = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
	replicaSetAlgorithmValue             = "SHA-256"
)
