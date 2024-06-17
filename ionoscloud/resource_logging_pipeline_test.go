//go:build all || logging

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
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
					testAccLoggingPipelineExistenceCheck(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
					resource.TestCheckResourceAttrSet(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineGrafanaAddressAttribute),
				),
			},
			{
				Config: LoggingPipelineDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineGrafanaAddressAttribute),
				),
			},
			{
				Config: LoggingPipelineDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineNameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, pipelineGrafanaAddressAttribute),
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
					testAccLoggingPipelineExistenceCheck(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineNameAttribute, pipelineNameUpdatedValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceUpdatedValue,
						pipelineLogTagAttribute:      pipelineLogTagUpdatedValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolUpdatedValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionUpdatedValue,
					}),
				),
			},
		},
	})
}

func testAccLoggingPipelineDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).LoggingClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.LoggingPipelineResource {
			continue
		}
		pipelineId := rs.Primary.ID
		_, apiResponse, err := client.GetPipelineById(ctx, pipelineId)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of Logging pipeline with ID: %s, error: %w", pipelineId, err)
			}
		} else {
			return fmt.Errorf("Logging pipeline with ID: %s still exists", pipelineId)
		}
	}
	return nil
}

func testAccLoggingPipelineExistenceCheck(path string, pipeline *logging.Pipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).LoggingClient
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
			return fmt.Errorf("an error occurred while fetching Logging pipeline with ID: %s, error: %w", pipelineId, err)
		}
		pipeline = &pipelineResponse
		return nil
	}
}

const LoggingPipelineDataSourceMatchById = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	id = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.id
}
`

const LoggingPipelineDataSourceMatchByName = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	name = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.name
}
`

const LoggingPipelineDataSourceInvalidBothIDAndName = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	id = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.id
	name = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.name
}
`

const LoggingPipelineDataSourceInvalidNoIDNoName = `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + ` {
}
`

const LoggingPipelineDataSourceWrongNameError = `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + ` {
	name = "nonexistent"
}
`

const LoggingPipelineConfigUpdate = `
resource ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestResourceName + ` {
	` + pipelineNameAttribute + ` = "` + pipelineNameUpdatedValue + `"
	` + pipelineLogUpdated + `
}
`
