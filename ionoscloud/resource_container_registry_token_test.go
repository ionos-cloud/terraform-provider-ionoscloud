//go:build all || cr
// +build all cr

package ionoscloud

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
)

func TestAccContainerRegistryTokenBasic(t *testing.T) {
	var containerRegistryToken cr.TokenResponse

	expiryDate := time.Now().Add(13 * time.Hour).UTC()
	expiryDateTZOffsetLayout := expiryDate.Format(constant.DatetimeTZOffsetLayout)

	templateData := struct{ ExpiryDate string }{ExpiryDate: expiryDate.Format(constant.DatetimeZLayout)}

	defer removeTestFile()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckContainerRegistryTokenDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: getConfigurationFromTemplate(testAccCheckContainerRegistryTokenConfigBasic, templateData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryTokenExists(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, &containerRegistryToken),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "expiry_date", expiryDateTZOffsetLayout),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.0", "push"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.name", "Scope1"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.type", "repository"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "status", "enabled"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "name", constant.ContainerRegistryTokenTestResource),
				),
			},
			{
				Config: getConfigurationFromTemplate(testAccDataSourceContainerRegistryTokenMatchId, templateData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "expiry_date", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "scopes.0.actions.0", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "scopes.0.name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "scopes.0.type", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "status", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceById, "name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "name"),
				),
			},
			{
				Config: getConfigurationFromTemplate(testAccDataSourceContainerRegistryTokenMatchName, templateData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "expiry_date", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.actions.0", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.type", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "status", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "name")),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongNameError,
				ExpectError: regexp.MustCompile("no token found with the specified name"),
			},
			{
				Config: getConfigurationFromTemplate(testAccDataSourceContainerRegistryTokenPartialMatchName, templateData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "expiry_date", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "expiry_date"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.actions.0", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "scopes.0.type", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "status", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "status"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestDataSourceByName, "name", constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "name")),
			},
			{
				Config:      testAccDataSourceContainerRegistryTokenWrongPartialNameError,
				ExpectError: regexp.MustCompile("no token found with the specified name"),
			},
			{
				Config: getConfigurationFromTemplate(testAccCheckContainerRegistryTokenConfigUpdate, templateData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerRegistryTokenExists(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, &containerRegistryToken),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "expiry_date", expiryDateTZOffsetLayout),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.0", "push"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.actions.1", "pull"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.name", "Scope1"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.0.type", "repository"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.1.actions.0", "push"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.1.actions.1", "pull"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.1.name", "Scope2"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "scopes.1.type", "backup"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "status", "disabled"),
					resource.TestCheckResourceAttr(constant.ContainerRegistryTokenResource+"."+constant.ContainerRegistryTokenTestResource, "name", constant.ContainerRegistryTokenTestResource)),
			},
			{
				Config:      getConfigurationFromTemplate(testAccDataSourceContainerRegistryTokenMultipleTokensFound, templateData),
				ExpectError: regexp.MustCompile("more than one token found with the specified criteria: name = test-container-registry-token"),
			},
		},
	})
}

func testAccCheckContainerRegistryTokenDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).ContainerClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ContainerRegistryTokenResource {
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

		foundToken, _, err := client.GetToken(ctx, rs.Primary.Attributes["registry_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching container registry token %s: %w", rs.Primary.ID, err)
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

// Use templating for 'expiry_date' to ensure the date is always in the future
const testAccCheckContainerRegistryTokenConfigBasic = testAccCheckContainerRegistryConfigBasic + `
resource ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestResource + ` {
  expiry_date    = "{{.ExpiryDate}}"
  name           = "` + constant.ContainerRegistryTokenTestResource + `"
  scopes {
    actions  = ["push"]
    name     = "Scope1"
    type     = "repository"
  }
  status                 = "enabled"
  registry_id            = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
  save_password_to_file  = "` + testFileName + `"
}
`

const testAccCheckContainerRegistryTokenConfigUpdate = testAccCheckContainerRegistryConfigBasic + `
resource ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestResource + ` {
  expiry_date    = "{{.ExpiryDate}}"
  name           = "` + constant.ContainerRegistryTokenTestResource + `"
  scopes {
    actions    = ["push", "pull"]
    name       = "Scope1"
    type       = "repository"
  }
  scopes {
    actions    = ["push", "pull"]
    name       = "Scope2"
    type       = "backup"
  }
  status         = "disabled"
  registry_id    = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMatchId = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceById + ` {
  id = ` + constant.ContainerRegistryTokenResource + `.` + constant.ContainerRegistryTokenTestResource + `.id
  registry_id = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMatchName = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  name = "` + constant.ContainerRegistryTokenTestResource + `"
  registry_id = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMatchNameAndLocation = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  location       = "` + constant.ContainerRegistryTokenResource + `.` + constant.ContainerRegistryTestResource + `.location
  name	         = "` + constant.ContainerRegistryTokenTestResource + `"
  registry_id    = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenWrongNameError = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  name = "wrong_name"
  registry_id = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenPartialMatchName = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  name             = "test"
  partial_match    = true
  registry_id      = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenWrongPartialNameError = testAccCheckContainerRegistryTokenConfigBasic + `
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  name             = "wrong_name"
  partial_match    = true
  registry_id      = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`

const testAccDataSourceContainerRegistryTokenMultipleTokensFound = testAccCheckContainerRegistryTokenConfigBasic + `
resource ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestResource + `1 {
  expiry_date    = "{{.ExpiryDate}}"
  name           = "` + constant.ContainerRegistryTokenTestResource + `1"
  scopes {
    actions    = ["push"]
    name       = "Scope1"
    type       = "repository"
  }
  status                   = "enabled"
  registry_id              = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
  save_password_to_file    = "` + testFileName + `"
}
data ` + constant.ContainerRegistryTokenResource + ` ` + constant.ContainerRegistryTokenTestDataSourceByName + ` {
  name             = "` + constant.ContainerRegistryTokenTestResource + `"
  partial_match    = true
  registry_id      = ` + constant.ContainerRegistryResource + `.` + constant.ContainerRegistryTestResource + `.id
}
`
