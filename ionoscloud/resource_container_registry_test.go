//go:build all || cr
// +build all cr

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
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
					testAccCheckContainerRegistryExists(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, &containerRegistry),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time", "05:19:00+00:00"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1", "Tuesday"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location", "de/fra"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name", constant.ContainerRegistryTestResource),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning", "false"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.time", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.days.0", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "garbage_collection_schedule.0.days.1", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "location", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "name", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceById, "features.0.vulnerability_scanning", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "location", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "name", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "features.0.vulnerability_scanning", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryMatchNameAndLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "location", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "name", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "features.0.vulnerability_scanning", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryWrongNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria: name ="),
			},
			{
				Config:      testAccDataSourceContainerRegistryWrongLocationErr,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria: location ="),
			},
			{
				Config: testAccDataSourceContainerRegistryPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.time", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.0", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "garbage_collection_schedule.0.days.1", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "location", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "name", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestDataSourceByName, "features.0.vulnerability_scanning", constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning"),
				),
			},
			{
				Config:      testAccDataSourceContainerRegistryWrongPartialNameError,
				ExpectError: regexp.MustCompile("no registry found with the specified criteria: name ="),
			},
			{
				Config: testAccCheckContainerRegistryConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryExists(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, &containerRegistry),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.time", "01:23:00+00:00"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "garbage_collection_schedule.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "location", "de/fra"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "name", constant.ContainerRegistryTestResource),
					resource.TestCheckResourceAttr(constant.ContainerRegistryResource+"."+constant.ContainerRegistryTestResource, "features.0.vulnerability_scanning", "true"),
				),
			},
			{
				Config:      testAccDataSourceCRTokenNameMultipleRegsFound,
				ExpectError: regexp.MustCompile("more than one registry found with the specified criteria: name ="),
			},
			{
				Config:      testAccDataSourceCRTokenLocationMultipleRegsFound,
				ExpectError: regexp.MustCompile("more than one registry found with the specified criteria: name =  location = de/fra"),
			},
		},
	})
}

func testAccCheckContainerRegistryDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).ContainerClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ContainerRegistryResource {
			continue
		}

		_, apiResponse, err := client.GetRegistry(ctx, rs.Primary.ID)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of the container registry %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("container registry %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckContainerRegistryExists(n string, registry *cr.RegistryResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).ContainerClient

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
			return fmt.Errorf("an error occured while fetching container registry %s: %w", rs.Primary.ID, err)
		}
		if *foundRegistry.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		registry = &foundRegistry

		return nil
	}
}

const testAccCheckContainerRegistryConfigBasic = `
resource ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestResource + ` {
   garbage_collection_schedule {
    days			 = ["Monday", "Tuesday"]
    time             = "05:19:00+00:00"
  }
  features {
    vulnerability_scanning = false
  }
  location           = "de/fra"
  name		         = "` + constant.ContainerRegistryTestResource + `"
}
`

const testAccCheckContainerRegistryConfigUpdate = `
resource ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestResource + ` {
   garbage_collection_schedule {
    days			 = ["Monday"]
    time             = "01:23:00+00:00"
  }
  features {
    vulnerability_scanning = true
  }
  location           = "de/fra"
  name		         = "` + constant.ContainerRegistryTestResource + `"
}
`

const testAccDataSourceContainerRegistryMatchId = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceById + ` {
  id	= ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryMatchName = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  name	= "` + constant.ContainerRegistryTestResource + `"
}
`

const testAccDataSourceContainerRegistryMatchNameAndLocation = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  name	   = "` + constant.ContainerRegistryTestResource + `"
  location = "de/fra" 
}
`
const testAccDataSourceContainerRegistryWrongIdError = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  id	= "wrong_id"
}
`
const testAccDataSourceContainerRegistryWrongNameError = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  name	= "wrong_name"
}
`
const testAccDataSourceContainerRegistryWrongLocationErr = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  location	= "de/txl"
}
`
const testAccDataSourceContainerRegistryPartialMatchName = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  name	= "test"
  partial_match = true
}
`

const testAccDataSourceContainerRegistryWrongPartialNameError = testAccCheckContainerRegistryConfigBasic + `
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
  name	= "wrong_name"
  partial_match = true
}
`
const testAccDataSourceCRTokenNameMultipleRegsFound = testAccCheckContainerRegistryConfigUpdate + `
resource ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestResource + `1 {
   garbage_collection_schedule {
    days			 = ["Monday", "Tuesday"]
    time             = "05:19:00+00:00"
  }
  location           = "de/fra"
  name		         = "` + constant.ContainerRegistryTestResource + `1"
}
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
depends_on = [ ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `]
  partial_match = true
  name	= "` + constant.ContainerRegistryTestResource + `"
}
`

const testAccDataSourceCRTokenLocationMultipleRegsFound = testAccCheckContainerRegistryConfigUpdate + `
resource ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestResource + `1 {
   garbage_collection_schedule {
    days			 = ["Monday", "Tuesday"]
    time             = "05:19:00+00:00"
  }
  location           = "de/fra"
  name		         = "` + constant.ContainerRegistryTestResource + `1"
}
data ` + constant.ContainerRegistryResource + ` ` + constant.ContainerRegistryTestDataSourceByName + ` {
depends_on = [ ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `]
  location	= "de/fra"
}
`
