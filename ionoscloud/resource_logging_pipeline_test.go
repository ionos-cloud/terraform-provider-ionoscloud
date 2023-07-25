//go:build all || logging

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	logging "github.com/ionos-cloud/sdk-go-logging"
)

func TestAccLoggingPipeline(t *testing.T) {
	var Pipeline logging.Pipeline

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccLoggingPipelineDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: LoggingPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
				),
			},
			{
				Config: LoggingPipelineDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
				),
			},
			{
				Config: LoggingPipelineDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
				),
			},
			{
				Config:      LoggingPipelineDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      LoggingPipelineDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the Logging pipeline ID or name"),
			},
			{
				Config:      LoggingPipelineDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no Logging pipelines found with the specified name"),
			},
			{
				Config: LoggingPipelineConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameUpdatedValue),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceUpdatedValue,
						pipelineLogTagAttribute:      pipelineLogTagUpdatedValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolUpdatedValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(LoggingPipelineResource+"."+LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionUpdatedValue,
					}),
				),
			},
		},
	})
}

func testAccLoggingPipelineDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).LoggingClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != LoggingPipelineResource {
			continue
		}
		pipelineId := rs.Primary.ID
		_, apiResponse, err := client.GetPipelineById(ctx, pipelineId)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occured while checking the destruction of Logging pipeline with ID: %s, error: %w", pipelineId, err)
			}
		} else {
			return fmt.Errorf("Logging pipeline with ID: %s still exists", pipelineId)
		}
	}
	return nil
}

func testAccLoggingPipelineExistenceCheck(path string, pipeline *logging.Pipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).LoggingClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the Logging pipeline")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		pipelineId := rs.Primary.ID
		pipelineResponse, _, err := client.GetPipelineById(ctx, pipelineId)
		if err != nil {
			return fmt.Errorf("an error occured while fetching Logging pipeline with ID: %s, error: %w", pipelineId, err)
		}
		pipeline = &pipelineResponse
		return nil
	}
}

const LoggingPipelineDataSourceMatchById = LoggingPipelineConfig + `
` + DataSource + ` ` + LoggingPipelineResource + ` ` + LoggingPipelineTestDataSourceName + `{
	id = ` + LoggingPipelineResource + `.` + LoggingPipelineTestResourceName + `.id
}
`

const LoggingPipelineDataSourceMatchByName = LoggingPipelineConfig + `
` + DataSource + ` ` + LoggingPipelineResource + ` ` + LoggingPipelineTestDataSourceName + `{
	name = ` + LoggingPipelineResource + `.` + LoggingPipelineTestResourceName + `.name
}
`

const LoggingPipelineDataSourceInvalidBothIDAndName = LoggingPipelineConfig + `
` + DataSource + ` ` + LoggingPipelineResource + ` ` + LoggingPipelineTestDataSourceName + `{
	id = ` + LoggingPipelineResource + `.` + LoggingPipelineTestResourceName + `.id
	name = ` + LoggingPipelineResource + `.` + LoggingPipelineTestResourceName + `.name
}
`

const LoggingPipelineDataSourceInvalidNoIDNoName = `
` + DataSource + ` ` + LoggingPipelineResource + ` ` + LoggingPipelineTestDataSourceName + ` {
}
`

const LoggingPipelineDataSourceWrongNameError = `
` + DataSource + ` ` + LoggingPipelineResource + ` ` + LoggingPipelineTestDataSourceName + ` {
	name = "nonexistent"
}
`

const LoggingPipelineConfigUpdate = `
resource ` + LoggingPipelineResource + ` ` + LoggingPipelineTestResourceName + ` {
	` + pipelineNameAttribute + ` = "` + pipelineNameUpdatedValue + `"
	` + pipelineLogUpdated + `
}
`
