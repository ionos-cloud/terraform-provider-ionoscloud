//go:build all || cr
// +build all cr

package ionoscloud

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
)

func TestAccContainerRegistryTokenBasic(t *testing.T) {
	var containerRegistryToken cr.TokenResponse
	defer removeTestFile()
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
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "expiry_date", "2023-01-13 16:27:42 +0000 UTC"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.0", "push"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.name", "Scope1"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.type", "repository"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "status", "enabled"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name", ContainerRegistryTokenTestResource),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryTokenMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "expiry_date", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "scopes.0.actions.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "scopes.0.name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "scopes.0.type", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "status", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceById, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name"),
				),
			},
			{
				Config: testAccDataSourceContainerRegistryTokenMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "expiry_date", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.actions.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.type", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "status", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name")),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongNameError,
				ExpectError: regexp.MustCompile("no token found with the specified name"),
			},
			{
				Config: testAccDataSourceContainerRegistryTokenPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "expiry_date", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.actions.0", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "scopes.0.type", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "status", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestDataSourceByName, "name", ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name")),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongPartialNameError,
				ExpectError: regexp.MustCompile("no token found with the specified name"),
			},
			{
				Config: testAccCheckContainerRegistryTokenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryTokenExists(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, &containerRegistryToken),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "expiry_date", "2023-01-23 16:27:42 +0000 UTC"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.0", "push"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.actions.1", "pull"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.name", "Scope1"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.0.type", "repository"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.1.actions.0", "push"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.1.actions.1", "pull"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.1.name", "Scope2"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "scopes.1.type", "backup"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "status", "disabled"),
					resource.TestCheckResourceAttr(ContainerRegistryTokenResource+"."+ContainerRegistryTokenTestResource, "name", ContainerRegistryTokenTestResource)),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenMultipleTokensFound,
				ExpectError: regexp.MustCompile("more than one registry found with the specified criteria"),
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
				return fmt.Errorf("an error occurred while checking the destruction of the container registry token %s: %w", rs.Primary.ID, err)
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
			return fmt.Errorf("an error occured while fetching container registry token %s: %w", rs.Primary.ID, err)
		}
		if *foundToken.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		registry = &foundToken

		return nil
	}
}

func removeTestFile() {
	os.Remove(testFileName)
}

const testFileName = "pass.txt"
const testAccCheckContainerRegistryTokenConfigBasic = testAccCheckContainerRegistryConfigBasic + `
resource ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestResource + ` {
  expiry_date        = "2023-01-13 16:27:42Z"
  name				 = "` + ContainerRegistryTokenTestResource + `"
  scopes  {
    actions			 = ["push"]
    name             = "Scope1"
    type             = "repository"
  }
  status	         = "enabled"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
  save_password_to_file = "` + testFileName + `"
}
`

const testAccCheckContainerRegistryTokenConfigUpdate = testAccCheckContainerRegistryConfigBasic + `
resource ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestResource + ` {
  expiry_date        = "2023-01-23 16:27:42Z"
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

const testAccDataSourceContainerRegistryTokenMatchNameAndLocation = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  location = "` + ContainerRegistryTokenResource + `.` + ContainerRegistryTestResource + `.location
  name	= "` + ContainerRegistryTokenTestResource + `"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenWrongNameError = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  name	= "wrong_name"
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
  name	= "wrong_name"
  partial_match = true
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMultipleTokensFound = testAccCheckContainerRegistryTokenConfigBasic + `
resource ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestResource + `1 {
  expiry_date        = "2023-01-13 16:27:42Z"
  name				 = "` + ContainerRegistryTokenTestResource + `1"
  scopes  {
    actions			 = ["push"]
    name             = "Scope1"
    type             = "repository"
  }
  status	         = "enabled"
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
  save_password_to_file = "` + testFileName + `"
}
data ` + ContainerRegistryTokenResource + ` ` + ContainerRegistryTokenTestDataSourceByName + ` {
  name	= "` + ContainerRegistryTokenTestResource + `"
  partial_match = true
  registry_id        = ` + ContainerRegistryResource + `.` + ContainerRegistryTestResource + `.id
}
`
