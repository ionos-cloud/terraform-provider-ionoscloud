//go:build all || dbaas || redis
// +build all dbaas redis

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redisdb "github.com/ionos-cloud/sdk-go-dbaas-redis"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDBaaSRedisDBReplicaSetBasic(t *testing.T) {
	var replicaSet redisdb.ReplicaSetRead
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDBaaSRedisDBReplicaSetDestroyCheck,
		Steps: []resource.TestStep{
			// This step tests a configuration that uses a hashed password.
			{
				Config: redisDBReplicaSetConfigHashedPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValue,
						clusterRamAttribute:   replicaSetRAMValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetHashedPasswordAttribute+".*", map[string]string{
						replicaSetHashAttribute:      replicaSetHashValue,
						replicaSetAlgorithmAttribute: replicaSetAlgorithmValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekAttribute,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			// This step deletes the replica set that was previously created, in order to make place
			// for another replica set with different credentials.
			{
				Config: redisDBReplicaSetConfigSetup,
			},
			// This step tests a configuration that uses the plain text password, this configuration
			// will be also used for 'update' tests.
			{
				Config: redisDBReplicaSetConfigPlainTextPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValue,
						clusterRamAttribute:   replicaSetRAMValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsUsernameAttribute, clusterCredentialsUsernameValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetPlainTextPasswordAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekAttribute,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			{
				Config: redisDBReplicaSetDataSourceMatchID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetResourcesAttribute+".0."+clusterCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, replicaSetResourcesAttribute+".0."+clusterRamAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByID, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
				),
			},
			{
				Config: redisDBReplicaSetDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetResourcesAttribute+".0."+clusterCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetResourcesAttribute+".0."+clusterRamAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue)),
			},
			// This step tests for basic updates for different attributes.
			// TODO -- Check what fields can actually be updated since for some of them it seems
			// that we are receiving API errors.
			{
				Config: redisDBReplicaSetConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterDisplayNameAttribute, replicaSetDisplayNameUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetLocationAttribute, replicaSetLocationValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetReplicasAttribute, replicaSetReplicasUpdateValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".*", map[string]string{
						clusterCoresAttribute: replicaSetCoresValueUpdate,
						clusterRamAttribute:   replicaSetRAMValueUpdate,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetResourcesAttribute+".0.storage"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetPersistenceModeAttribute, replicaSetPersistenceModeUpdateValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetEvictionPolicyAttribute, replicaSetEvictionPolicyUpdateValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsDatacenterIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsLanIDAttribute),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, connections+".0."+clusterConnectionsCidrAttribute),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+clusterCredentialsUsernameAttribute, clusterCredentialsUsernameValue),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterCredentialsAttribute+".0."+replicaSetPlainTextPasswordAttribute),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, clusterMaintenanceWindowAttribute+".*", map[string]string{
						clusterMaintenanceWindowTimeAttribute:         clusterMaintenanceWindowTimeUpdateValue,
						clusterMaintenanceWindowDayOfTheWeekAttribute: clusterMaintenanceWindowDayOfTheWeekUpdateValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetDNSNameAttribute),
				),
			},
			{
				Config:      redisDBReplicaSetDataSourceWrongName,
				ExpectError: regexp.MustCompile("no Redis cluster found with the specified display name"),
			},
			{
				Config:      redisDBReplicaSetDataSourceWrongID,
				ExpectError: regexp.MustCompile("an error occurred while fetching the Redis cluster with ID"),
			},
		},
	})
}

func testAccCheckDBaaSRedisDBReplicaSetDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).RedisDBClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaaSRedisDBReplicaSetResource {
			continue
		}
		_, apiResponse, err := client.GetReplicaSet(ctx, rs.Primary.ID, rs.Primary.Attributes[clusterLocationAttribute])
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking the destruction of RedisDB replica set with ID: %v, error: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("RedisDB replica set with ID: %v still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDBaaSRedisDBReplicaSetExists(n string, replicaSet *redisdb.ReplicaSetRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).RedisDBClient
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
			return fmt.Errorf("an error occurred while fetching RedisDB replica set with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *foundReplicaSet.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}
		replicaSet = &foundReplicaSet
		return nil
	}
}

// This configuration contains the resources that need to be created before creating a RedisDB replica set
const redisDBReplicaSetConfigSetup = `
resource ` + constant.DatacenterResource + ` ` + datacenterResourceName + ` {
  name        = "redisdb_datacenter_example"
  location    = "` + replicaSetLocationValue + `"
  description = "Datacenter for testing RedisDB replica set"
}

resource ` + constant.LanResource + ` ` + lanResourceName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id 
  public        = false
  name          = "redisdb_lan_example"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name                    = "example"
  datacenter_id           = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
  cores                   = 2
  ram                     = 2048
  availability_zone       = "ZONE_1"
  cpu_family              = "INTEL_SKYLAKE"
  image_name              = "debian-10-genericcloud-amd64-20240114-1626"
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

const redisDBReplicaSetConfigHashedPassword = redisDBReplicaSetConfigSetup + `
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyValue + `"
  ` + resources + `
  ` + connections + `
  ` + maintenanceWindow + `
  ` + credentialsHashedPassword + `
}
`

const redisDBReplicaSetConfigPlainTextPassword = redisDBReplicaSetConfigSetup + `
resource ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyValue + `"
  ` + resources + `
  ` + connections + `
  ` + maintenanceWindow + `
  ` + credentialsPlainTextPassword + `
}
`

const redisDBReplicaSetConfigUpdate = redisDBReplicaSetConfigSetup + `
resource ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionUpdateValue + `"
  ` + clusterDisplayNameAttribute + ` = "` + replicaSetDisplayNameUpdateValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasUpdateValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeUpdateValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyUpdateValue + `"
  ` + resourcesUpdate + `
  ` + connections + `
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
const redisDBReplicaSetDataSourceMatchID = redisDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByID + ` {
	id = ` + constant.DBaaSRedisDBReplicaSetResource + `.` + constant.DBaaSReplicaSetTestResource + `.id
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For testing data source match by name
const redisDBReplicaSetDataSourceMatchName = redisDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "` + replicaSetDisplayNameValue + `"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by name
const redisDBReplicaSetDataSourceWrongName = redisDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "wrong_name"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by ID - use a 0000 uuidv4
const redisDBReplicaSetDataSourceWrongID = redisDBReplicaSetConfigPlainTextPassword + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByID + ` {
	id = "00000000-0000-0000-0000-000000000000"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

const connections = clusterConnectionsAttribute + `{
	` + clusterConnectionsDatacenterIDAttribute + ` = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
    ` + clusterConnectionsLanIDAttribute + ` = ` + constant.LanResource + `.` + lanResourceName + `.id
	` + clusterConnectionsCidrAttribute + ` = ` + clusterConnectionsCidrValue + `
}`

const maintenanceWindow = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeValue + `"
}`

const maintenanceWindowUpdate = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekUpdateValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeUpdateValue + `"
}`
const credentialsPlainTextPassword = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameValue + `"
	` + replicaSetPlainTextPasswordAttribute + ` = ` + constant.RandomPassword + `.replicaset_password.result
}`

const credentialsHashedPassword = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameValue + `"
	` + hashedPasswordObject + `
}`

const hashedPasswordObject = replicaSetHashedPasswordAttribute + `{
    ` + replicaSetHashAttribute + ` = "` + replicaSetHashValue + `"
	` + replicaSetAlgorithmAttribute + ` = "` + replicaSetAlgorithmValue + `"
`

// Attributes
const (
	replicaSetLocationAttribute          = "location"
	replicaSetVersionAttribute           = "redis_version"
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
	replicaSetLocationValue              = "es/vit"
	replicaSetLocationUpdateValue        = "de/txl"
	replicaSetVersionValue               = "7.2"
	replicaSetVersionUpdateValue         = "7.0"
	replicaSetDisplayNameValue           = "MyReplicaSet"
	replicaSetDisplayNameUpdateValue     = "UpdatedReplicaSet"
	replicaSetReplicasValue              = "4"
	replicaSetReplicasUpdateValue        = "6"
	replicaSetPersistenceModeValue       = "RDB"
	replicaSetPersistenceModeUpdateValue = "RDB_AOF"
	replicaSetEvictionPolicyValue        = "noeviction"
	replicaSetEvictionPolicyUpdateValue  = "allkeys-lru"
	replicaSetCoresValue                 = "1"
	replicaSetCoresValueUpdate           = "2"
	replicaSetRAMValue                   = "6"
	replicaSetRAMValueUpdate             = "8"
	clusterConnectionsCidrValue          = "local.database_ip_cidr"
	clusterCredentialsUsernameValue      = "testuser"
	datacenterResourceName               = "datacenter_example"
	lanResourceName                      = "lan_example"
	replicaSetHashValue                  = "492f3f38d6b5d3ca859514e250e25ba65935bcdd9f4f40c124b773fe536fee7d"
	replicaSetAlgorithmValue             = "SHA-256"
)
