//go:build all || dbaas || mongo
// +build all dbaas mongo

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserMongoBasic(t *testing.T) {
	var user mongo.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckMongoUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMongoUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoUserExists(constant.DBaasMongoUserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "username", constant.UserTestResource),
					resource.TestCheckResourceAttrSet(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "password"),
					resource.TestCheckResourceAttr(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "roles.#", "2"),
				),
			},
			{
				Config: testAccDataSourceMongoUserMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoUserResource+"."+constant.UserDataSourceById, "username", constant.DBaasMongoUserResource+"."+constant.UserTestResource, "username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoUserResource+"."+constant.UserDataSourceById, "roles.#", constant.DBaasMongoUserResource+"."+constant.UserTestResource, "roles.#"),
				),
			},
			{
				Config: testAccDataSourceMongoUserMatchUsername,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoUserResource+"."+constant.UserDataSourceById, "username", constant.DBaasMongoUserResource+"."+constant.UserTestResource, "username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoUserResource+"."+constant.UserDataSourceById, "roles.#", constant.DBaasMongoUserResource+"."+constant.UserTestResource, "roles.#"),
				),
			},
			{
				Config:      testAccDataSourceMongoUserWrongUsername,
				ExpectError: regexp.MustCompile(`no DBaaS mongo user found with the specified username =`),
			},
			{
				Config: testAccCheckMongoUserConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoUserExists(constant.DBaasMongoUserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "username", constant.UserTestResource),
					resource.TestCheckResourceAttrSet(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "password"),
					resource.TestCheckResourceAttr(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "roles.#", "1"),
					resource.TestCheckResourceAttr(constant.DBaasMongoUserResource+"."+constant.UserTestResource, "roles.0.role", "readWrite"),
				),
			},
		},
	})
}

func testAccCheckMongoUserDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).MongoClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DBaasMongoUserResource {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		username := rs.Primary.Attributes["username"]
		_, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("user still exists %s - an error occurred while checking it %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("user still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoUserExists(n string, user *mongo.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).MongoClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckUserExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		username := rs.Primary.Attributes["username"]
		foundUser, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
		apiResponse.LogInfo()
		if err != nil {
			return fmt.Errorf("error occurred while fetching User: %s %w", rs.Primary.ID, err)
		}

		user = &foundUser

		return nil
	}
}

var testAccCheckMongoUserConfigBasic = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + constant.DBaaSClusterTestResource + `"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
}

resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.DBaasMongoUserResource + ` ` + constant.UserTestResource + ` {
  cluster_id = ` + constant.DBaasMongoClusterResource + `.` + constant.DBaaSClusterTestResource + `.id 
  username = "` + constant.UserTestResource + `"
  password = ` + constant.RandomPassword + `.user_password.result
  roles {
    role = "read"
    database = "db1"
  }
  roles {
	role = "readWrite"
	database = "db2"
  }
}`

var testAccCheckMongoUserConfigUpdated = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + constant.DBaaSClusterTestResource + `"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
}

resource ` + constant.DBaasMongoUserResource + ` ` + constant.UserTestResource + ` {
  cluster_id = ` + constant.DBaasMongoClusterResource + `.` + constant.DBaaSClusterTestResource + `.id 
  username = "` + constant.UserTestResource + `"
  password = ` + constant.RandomPassword + `.user_password_updated.result
  roles {
    role = "readWrite"
    database = "db1"
  }
}
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

var testAccDataSourceMongoUserMatchId = testAccCheckMongoUserConfigBasic + `
data ` + constant.DBaasMongoUserResource + ` ` + constant.UserDataSourceById + ` {
  cluster_id = ` + constant.DBaasMongoUserResource + `.` + constant.UserTestResource + `.cluster_id
  username = ` + constant.DBaasMongoUserResource + `.` + constant.UserTestResource + `.username
}
`

var testAccDataSourceMongoUserMatchUsername = testAccCheckMongoUserConfigBasic + `
data ` + constant.DBaasMongoUserResource + ` ` + constant.UserDataSourceById + ` {
  cluster_id = ` + constant.DBaasMongoUserResource + `.` + constant.UserTestResource + `.cluster_id
  username = ` + constant.DBaasMongoUserResource + `.` + constant.UserTestResource + `.username
}
`

var testAccDataSourceMongoUserWrongUsername = testAccCheckMongoUserConfigBasic + `
data ` + constant.DBaasMongoUserResource + ` ` + constant.UserDataSourceById + ` {
  cluster_id = ` + constant.DBaasMongoUserResource + `.` + constant.UserTestResource + `.cluster_id
  username = "willnotwork"
}
`
