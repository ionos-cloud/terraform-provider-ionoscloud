package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
	"regexp"
	"testing"
)

func TestAccContainerRegistryTokenBasic(t *testing.T) {
	var containerRegistryToken cr.TokenResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckContainerRegistryTokenDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckContainerRegistryTokenConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryTokenExists(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, &containerRegistryToken),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "garbage_collection_schedule.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "garbage_collection_schedule.0.days.1", "Tuesday"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "maintenance_window.0.days.0", "Sunday"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "location", "de/txl"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name", ContainerRegistryTokenTestResource),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryTokenMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "garbage_collection_schedule.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "garbage_collection_schedule.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "garbage_collection_schedule.0.days.1", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "maintenance_window.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "maintenance_window.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "location", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryTokenMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.1", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "location", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria"),
			}, {
				Config: testAccDataSourceContainerRegistryTokenPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.1", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.time", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.days.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "location", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongPartialNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria"),
			},
			{
				Config: testAccCheckContainerRegistryTokenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryTokenExists(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, &containerRegistryToken),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "garbage_collection_schedule.0.time", "11:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "maintenance_window.0.days.0", "Saturday"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "location", "de/txl"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name", ContainerRegistryTokenTestResource),
				),
			},
		},
	})
}

func testAccCheckContainerRegistryTokenDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).ContainerClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ContainerRegistryTokenResource {
			continue
		}

		_, apiResponse, err := client.GetToken(ctx, rs.Primary.Attributes["registry_id"], rs.Primary.ID)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of the container registry token %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("container registry token %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckContainerRegistryTokenExists(n string, registry *cr.TokenResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).ContainerClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundToken, _, err := client.GetToken(ctx, rs.Primary.Attributes["registry_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occured while fetching container registry token %s: %s", rs.Primary.ID, err)
		}
		if *foundToken.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		registry = &foundToken

		return nil
	}
}

const testAccCheckContainerRegistryTokenConfigBasic = testAccCheckContainerRegistryConfigBasic + `
resource ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestResource + ` {
   credentials {
	username		 = "username"
    password         = "password"
  }
  expiry_date        = "2023-01-13T16:27:42Z"
  name				 = "` + ContainerRegistryTokenTestResource + `"
  scopes  {
    actions			 = ["push"]
    name             = "Scope1"
    type             = "repository"
  }
  status	         = "enabled"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccCheckContainerRegistryTokenConfigUpdate = testAccCheckContainerRegistryConfigBasic + `
resource ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestResource + ` {
   credentials {
	username		 = "usernameUpdated"
    password         = "passwordUpdated"
  }
  expiry_date        = "2023-01-23T16:27:42Z"
  name				 = "` + ContainerRegistryTokenTestResource + `"
  scopes  {
    actions			 = ["push", "pull"]
    name             = "Scope1"
    type             = "repository"
  }
  scopes  {
    actions			 = ["push", "pull"]
    name             = "Scope2"
    type             = "backup"
  }
  status	         = "disabled"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMatchId = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceById + ` {
  id	= ` + ContainerRegistryTokenResource + `.` + ContainerRegistryTokenTestResource + `.id
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMatchName = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  name	= "` + ContainerRegistryTokenTestResource + `"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenWrongNameError = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  display_name	= "wrong_name"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenPartialMatchName = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  name	= "test"
  partial_match = true
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenWrongPartialNameError = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  display_name	= "wrong_name"
  partial_match = true
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`
