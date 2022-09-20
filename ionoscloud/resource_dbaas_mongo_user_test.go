//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"context"
	"fmt"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserMongoBasic(t *testing.T) {
	var user mongo.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMongoUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoUserExists(DBaasMongoUserResource+"."+UserTestResource, &user),
					resource.TestCheckResourceAttr(DBaasMongoUserResource+"."+UserTestResource, "username", UserTestResource),
					resource.TestCheckResourceAttr(DBaasMongoUserResource+"."+UserTestResource, "database", "admin"),
					resource.TestCheckResourceAttrSet(DBaasMongoUserResource+"."+UserTestResource, "password"),
					resource.TestCheckResourceAttr(DBaasMongoUserResource+"."+UserTestResource, "roles.0", "read"),
					resource.TestCheckResourceAttr(DBaasMongoUserResource+"."+UserTestResource, "roles.1", "readAnyDatabase"),
				),
			},
			//{
			//	Config: testAccDataSourceMongoUserMatchId,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "first_name", DBaasMongoUserResource+"."+UserTestResource, "first_name"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "last_name", DBaasMongoUserResource+"."+UserTestResource, "last_name"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "email", DBaasMongoUserResource+"."+UserTestResource, "email"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "administrator", DBaasMongoUserResource+"."+UserTestResource, "administrator"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "force_sec_auth", DBaasMongoUserResource+"."+UserTestResource, "force_sec_auth"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "sec_auth_active", DBaasMongoUserResource+"."+UserTestResource, "sec_auth_active"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "s3_canonical_user_id", DBaasMongoUserResource+"."+UserTestResource, "s3_canonical_user_id"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "active", DBaasMongoUserResource+"."+UserTestResource, "active"),
			//	),
			//},
			//{
			//	Config: testAccDataSourceMongoUserMatchUsername,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "first_name", DBaasMongoUserResource+"."+UserTestResource, "first_name"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "last_name", DBaasMongoUserResource+"."+UserTestResource, "last_name"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "email", DBaasMongoUserResource+"."+UserTestResource, "email"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "administrator", DBaasMongoUserResource+"."+UserTestResource, "administrator"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "force_sec_auth", DBaasMongoUserResource+"."+UserTestResource, "force_sec_auth"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "sec_auth_active", DBaasMongoUserResource+"."+UserTestResource, "sec_auth_active"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "s3_canonical_user_id", DBaasMongoUserResource+"."+UserTestResource, "s3_canonical_user_id"),
			//		resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoUserResource+"."+UserDataSourceById, "active", DBaasMongoUserResource+"."+UserTestResource, "active"),
			//	),
			//},
		},
	})
}

func testAccCheckMongoUserDestroyCheck(s *terraform.State) error {
	//client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	client := testAccProvider.Meta().(SdkBundle).MongoClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DBaasMongoUserResource {
			continue
		}
		//_, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		clusterId := rs.Primary.Attributes["cluster_id"]
		database := rs.Primary.Attributes["database"]
		username := rs.Primary.Attributes["username"]
		_, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
		//logApiRequestTime(apiResponse)

		if err != nil {
			if !httpMongoNotFound(apiResponse) {
				return fmt.Errorf("user still exists %s - an error occurred while checking it %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("user still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoUserExists(n string, user *mongo.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//client := testAccProvider.Meta().(SdkBundle).CloudApiClient
		client := testAccProvider.Meta().(SdkBundle).MongoClient
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
		database := rs.Primary.Attributes["database"]
		username := rs.Primary.Attributes["username"]
		//foundUser, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		foundUser, _, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
		//logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching User: %s %s", rs.Primary.ID, err)
		}
		//if *foundUser.Id != rs.Primary.ID { there is no id in the user: in swagger
		//	return fmt.Errorf("record not found")
		//}

		user = &foundUser

		return nil
	}
}

func httpMongoNotFound(resp *mongo.APIResponse) bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

var testAccCheckMongoUserConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + DBaasMongoClusterResource + ` ` + DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + DBaaSClusterTestResource + `"
  location = ` + DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
  
  credentials {
  	username = "username"
	password = "password"
  }
}



resource ` + DBaasMongoUserResource + ` ` + UserTestResource + ` {
  cluster_id = ` + DBaasMongoClusterResource + `.` + DBaaSClusterTestResource + `.id 
  username = "` + UserTestResource + `"
  database = "admin"
  password = "abc123-321CBA"
  roles = ["read", "readAnyDatabase"]
}`

var testAccDataSourceMongoUserMatchId = testAccCheckMongoUserConfigBasic + `
data ` + DBaasMongoUserResource + ` ` + UserDataSourceById + ` {
  id			= ` + DBaasMongoUserResource + `.` + UserTestResource + `.id
}
`

var testAccDataSourceMongoUserMatchUsername = testAccCheckMongoUserConfigBasic + `
data ` + DBaasMongoUserResource + ` ` + UserDataSourceById + ` {
  id			= ` + DBaasMongoUserResource + `.` + UserTestResource + `.id
}
`
