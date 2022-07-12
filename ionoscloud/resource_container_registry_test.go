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

func TestAccContainerRegistryBasic(t *testing.T) {
	var containerRegistry cr.RegistryResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckContainerRegistryDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckContainerRegistryConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryExists(ContainerRegistryResource+"."+ContainerRegistryTestResource, &containerRegistry),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "garbage_collection_schedule.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1", "Tuesday"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "maintenance_window.0.days.0", "Sunday"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "location", "de/txl"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "name", ContainerRegistryTestResource),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.days.1", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "maintenance_window.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "maintenance_window.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "location", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceById, "name", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "name"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "location", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "name", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "name"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryWrongNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria"),
			}, {
				Config: testAccDataSourceContainerRegistryPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.time", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.days.0", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "maintenance_window.0.days.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "location", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "name", ContainerRegistryResource+"."+ContainerRegistryTestDataSourceByName, "name"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryWrongPartialNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria"),
			},
			{
				Config: testAccCheckContainerRegistryConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryExists(ContainerRegistryResource+"."+ContainerRegistryTestResource, &containerRegistry),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "garbage_collection_schedule.0.time", "11:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "maintenance_window.0.days.0", "Saturday"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "location", "de/txl"),
					resource.TestCheckResourceAttr(ContainerRegistryResource+"."+ContainerRegistryTestResource, "name", ContainerRegistryTestResource),
				),
			},
		},
	})
}

func testAccCheckContainerRegistryDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).ContainerClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ContainerRegistryResource {
			continue
		}

		_, apiResponse, err := client.GetRegistry(ctx, rs.Primary.ID)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of the container registry %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("container registry %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckContainerRegistryExists(n string, registry *cr.RegistryResponse) resource.TestCheckFunc {
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

		foundRegistry, _, err := client.GetRegistry(ctx, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occured while fetching container registry %s: %s", rs.Primary.ID, err)
		}
		if *foundRegistry.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		registry = &foundRegistry

		return nil
	}
}

const testAccCheckContainerRegistryConfigBasic = `
resource ` + ContainerRegistryResource + ` ` + ContainerRegistryTestResource + ` {
   garbage_collection_schedule {
    days			 = ["Monday", "Tuesday"]
    time             = "10:00:00"
  }
  location           = "de/txl"
  maintenance_window {
    days			 = ["Sunday"]
    time             = "09:00:00"
  }
  name		         = "` + ContainerRegistryTestResource + `"
}
`

const testAccCheckContainerRegistryConfigUpdate = `
resource ` + ContainerRegistryResource + ` ` + ContainerRegistryTestResource + ` {
   garbage_collection_schedule {
    days			 = ["Monday"]
    time             = "11:00:00"
  }
  location           = "de/txl"
  maintenance_window {
    days			 = ["Saturday"]
    time             = "10:00:00"
  }
  name		         = "` + ContainerRegistryTestResource + `"
}
`

const testAccDataSourceContainerRegistryMatchId = testAccCheckContainerRegistryConfigBasic + `
data ` + ContainerRegistryResource + ` ` + ContainerRegistryTestDataSourceById + ` {
  id	= ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryMatchName = testAccCheckContainerRegistryConfigBasic + `
data ` + ContainerRegistryResource + ` ` + ContainerRegistryTestDataSourceByName + ` {
  name	= "` + ContainerRegistryTestResource + `"
}
`

const testAccDataSourceContainerRegistryWrongNameError = testAccCheckContainerRegistryConfigBasic + `
data ` + ContainerRegistryResource + ` ` + ContainerRegistryTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`

const testAccDataSourceContainerRegistryPartialMatchName = testAccCheckContainerRegistryConfigBasic + `
data ` + ContainerRegistryResource + ` ` + ContainerRegistryTestDataSourceByName + ` {
  name	= "test"
  partial_match = true
}
`

const testAccDataSourceContainerRegistryWrongPartialNameError = testAccCheckContainerRegistryConfigBasic + `
data ` + ContainerRegistryResource + ` ` + ContainerRegistryTestDataSourceByName + ` {
  display_name	= "wrong_name"
  partial_match = true
}
`
