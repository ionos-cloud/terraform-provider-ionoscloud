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
			{
				Config: redisDBReplicaSetConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestResource, replicaSetVersionAttribute, replicaSetVersionValue),
				),
			},
			{
				Config: redisDBReplicaSetDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetRAMAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue),
				),
			},
			{
				Config: redisDBReplicaSetDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBaaSRedisDBReplicaSetExists(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, &replicaSet),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetVersionAttribute, replicaSetVersionValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetReplicasAttribute, replicaSetReplicasValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetCoresAttribute, replicaSetCoresValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetRAMAttribute, replicaSetRAMValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, replicaSetDisplayNameAttribute, replicaSetDisplayNameValue),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsDatacenterIDAttribute, constant.DatacenterResource+"."+datacenterResourceName, "id"),
					resource.TestCheckResourceAttrPair(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterConnectionsAttribute+".0."+clusterConnectionsLanIDAttribute, constant.LanResource+"."+lanResourceName, "id"),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowDayOfTheWeekAttribute, clusterMaintenanceWindowDayOfTheWeekValue),
					resource.TestCheckResourceAttr(constant.DBaaSRedisDBReplicaSetResource+"."+constant.DBaaSReplicaSetTestDataSourceByName, clusterMaintenanceWindowAttribute+".0."+clusterMaintenanceWindowTimeAttribute, clusterMaintenanceWindowTimeValue)),
			},
			// other tests...
			{},
			// These should be at the very end
			{
				Config:      redisDBReplicaSetDataSourceWrongName,
				ExpectError: regexp.MustCompile("no Redis cluster found with the specified display name"),
			},
			{
				Config:      redisDBReplicaSetDataSourceWrongId,
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

		// TODO -- Create a constant for 'location'
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

const redisDBReplicaSetConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + datacenterResourceName + ` {
  name        = "redisdb_datacenter_example"
  location    = "es/vit"
  description = "Datacenter for testing RedisDB replica set"
}

resource ` + constant.LanResource + ` ` + lanResourceName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "redisdb_lan_example"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name                    = "example"
  datacenter_id           = ionoscloud_datacenter.datacenter_example.id
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

resource ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestResource + ` {
  ` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"  
  ` + replicaSetVersionAttribute + ` = "` + replicaSetVersionValue + `"
  ` + replicaSetDisplayNameAttribute + ` = "` + replicaSetDisplayNameValue + `"
  ` + replicaSetReplicasAttribute + ` = "` + replicaSetReplicasValue + `"
  ` + replicaSetPersistenceModeAttribute + ` = "` + replicaSetPersistenceModeValue + `"
  ` + replicaSetEvictionPolicyAttribute + ` = "` + replicaSetEvictionPolicyValue + `"
  ` + resources + `
  ` + connections + `
  ` + maintenanceWindow + `
  ` + credentials + `
}

resource ` + constant.RandomPassword + ` "replicaset_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
` + ServerImagePassword

// Internal resources
const resources = replicaSetResourcesAttribute + `{
	` + replicaSetCoresAttribute + ` = "` + replicaSetCoresValue + `"
	` + replicaSetRAMAttribute + ` = "` + replicaSetRAMValue + `"
}`

// For testing data source match by ID
const redisDBReplicaSetDataSourceMatchId = redisDBReplicaSetConfigBasic + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceById + ` {
	id = ` + constant.DBaaSRedisDBReplicaSetResource + `.` + constant.DBaaSReplicaSetTestResource + `.id
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For testing data source match by name
const redisDBReplicaSetDataSourceMatchName = redisDBReplicaSetConfigBasic + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "` + replicaSetDisplayNameValue + `"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by name
const redisDBReplicaSetDataSourceWrongName = redisDBReplicaSetConfigBasic + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceByName + ` {
	display_name	= "wrong_name"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// For negative case of data source match by ID - use a 0000 uuidv4
const redisDBReplicaSetDataSourceWrongId = redisDBReplicaSetConfigBasic + `
data ` + constant.DBaaSRedisDBReplicaSetResource + ` ` + constant.DBaaSReplicaSetTestDataSourceById + ` {
	id = "00000000-0000-0000-0000-000000000000"
	` + clusterLocationAttribute + ` = "` + replicaSetLocationValue + `"
}`

// TODO -- Rename anything related to clusters, for the moment just reuse them from MariaDB since
// they are the same
const connections = clusterConnectionsAttribute + `{
	` + clusterConnectionsDatacenterIDAttribute + ` = ` + constant.DatacenterResource + `.` + datacenterResourceName + `.id
    ` + clusterConnectionsLanIDAttribute + ` = ` + constant.LanResource + `.` + lanResourceName + `.id
	` + clusterConnectionsCidrAttribute + ` = ` + clusterConnectionsCidrValue + `
}`

const maintenanceWindow = clusterMaintenanceWindowAttribute + `{
	` + clusterMaintenanceWindowDayOfTheWeekAttribute + ` = "` + clusterMaintenanceWindowDayOfTheWeekValue + `"
	` + clusterMaintenanceWindowTimeAttribute + ` = "` + clusterMaintenanceWindowTimeValue + `"
}`

const credentials = clusterCredentialsAttribute + `{
	` + clusterCredentialsUsernameAttribute + ` = "` + clusterCredentialsUsernameValue + `"
	` + replicaSetPlainTextPasswordAttribute + ` = ` + constant.RandomPassword + `.replicaset_password.result
}`

// Attributes
const (
	replicaSetVersionAttribute           = "redis_version"
	replicaSetDisplayNameAttribute       = "display_name"
	replicaSetReplicasAttribute          = "replicas"
	replicaSetPersistenceModeAttribute   = "persistence_mode"
	replicaSetEvictionPolicyAttribute    = "eviction_policy"
	replicaSetResourcesAttribute         = "resources"
	replicaSetCoresAttribute             = "cores"
	replicaSetRAMAttribute               = "ram"
	replicaSetPlainTextPasswordAttribute = "plain_text_password"
)

// Values
const (
	replicaSetLocationValue         = "es/vit"
	replicaSetVersionValue          = "7.2"
	replicaSetDisplayNameValue      = "MyReplicaSet"
	replicaSetReplicasValue         = "4"
	replicaSetPersistenceModeValue  = "RDB"
	replicaSetEvictionPolicyValue   = "noeviction"
	replicaSetCoresValue            = "1"
	replicaSetRAMValue              = "6"
	clusterConnectionsCidrValue     = "local.database_ip_cidr"
	clusterCredentialsUsernameValue = "testuser"
	datacenterResourceName          = "datacenter_example"
	lanResourceName                 = "lan_example"
)
