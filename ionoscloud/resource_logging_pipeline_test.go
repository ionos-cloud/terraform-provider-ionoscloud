//go:build all || logging

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
)

func TestAccLoggingPipeline(t *testing.T) {
	var Pipeline logging.PipelineRead
	resourceName := constant.LoggingPipelineResource + "." + constant.LoggingPipelineTestResourceName
	dataSourceName := constant.DataSource + "." + constant.LoggingPipelineResource + "." + constant.LoggingPipelineTestDataSourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccLoggingPipelineDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: LoggingPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, nameAttribute, pipelineNameValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceValue,
						pipelineLogTagAttribute:      pipelineLogTagValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionValue,
					}),
					resource.TestCheckResourceAttrSet(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, "http_address"),
					resource.TestCheckNoResourceAttr(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, "tcp_address"),
					resource.TestCheckResourceAttrSet(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineGrafanaAddressAttribute),
					resource.TestCheckResourceAttrSet(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, "key"),
					resource.TestCheckResourceAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue1),
					resource.TestCheckTypeSetElemAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return rs.Primary.Attributes["location"] + ":" + rs.Primary.ID, nil
				},
				// the key is only returned in the create response
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config: LoggingPipelineDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, nameAttribute, pipelineNameValue),
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
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, "http_address"),
					resource.TestCheckNoResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, "tcp_address"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, pipelineLogAttribute+".*."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue1),
					resource.TestCheckTypeSetElemAttr(dataSourceName, pipelineLogAttribute+".*."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue2),
					resource.TestCheckResourceAttrPair(dataSourceName, pipelineStatusMessageAttribute, resourceName, pipelineStatusMessageAttribute),
				),
			},
			{
				Config: LoggingPipelineDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, &Pipeline),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, nameAttribute, pipelineNameValue),
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
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, "http_address"),
					resource.TestCheckNoResourceAttr(constant.DataSource+"."+constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestDataSourceName, "tcp_address"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, pipelineLogAttribute+".*."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue1),
					resource.TestCheckTypeSetElemAttr(dataSourceName, pipelineLogAttribute+".*."+pipelineLogLabelsAttribute+".*", pipelineLogLabelValue2),
					resource.TestCheckResourceAttrPair(dataSourceName, pipelineStatusMessageAttribute, resourceName, pipelineStatusMessageAttribute),
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
					resource.TestCheckResourceAttr(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, nameAttribute, pipelineNameUpdatedValue),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".*", map[string]string{
						pipelineLogSourceAttribute:   pipelineLogSourceUpdatedValue,
						pipelineLogTagAttribute:      pipelineLogTagUpdatedValue,
						pipelineLogProtocolAttribute: pipelineLogProtocolUpdatedValue,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.LoggingPipelineResource+"."+constant.LoggingPipelineTestResourceName, pipelineLogAttribute+".0."+pipelineLogDestinationAttribute+".*", map[string]string{
						pipelineLogDestinationTypeAttribute:      pipelineLogDestinationTypeValue,
						pipelineLogDestinationRetentionAttribute: pipelineLogDestinationRetentionUpdatedValue,
					}),
					resource.TestCheckResourceAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".*", pipelineLogLabelUpdatedValue),
				),
			},
			{
				Config: LoggingPipelineConfigUpdateNoLabels,
				Check: resource.ComposeTestCheckFunc(
					testAccLoggingPipelineExistenceCheck(resourceName, &Pipeline),
					resource.TestCheckResourceAttr(resourceName, nameAttribute, pipelineNameUpdatedValue),
					resource.TestCheckResourceAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogTagAttribute, pipelineLogTagUpdatedValue),
					resource.TestCheckResourceAttr(resourceName, pipelineLogAttribute+".0."+pipelineLogLabelsAttribute+".#", "0"),
				),
			},
		},
	})
}

func testAccLoggingPipelineDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).LoggingClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.LoggingPipelineResource {
			continue
		}
		pipelineID := rs.Primary.ID
		_, apiResponse, err := client.GetPipelineByID(ctx, rs.Primary.Attributes["location"], pipelineID)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of Logging pipeline with ID: %s, error: %w", pipelineID, err)
			}
		} else {
			return fmt.Errorf("logging pipeline with ID: %s still exists", pipelineID)
		}
	}
	return nil
}

func testAccLoggingPipelineExistenceCheck(path string, pipeline *logging.PipelineRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).LoggingClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the Logging pipeline")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		pipelineID := rs.Primary.ID
		pipelineResponse, _, err := client.GetPipelineByID(ctx, rs.Primary.Attributes["location"], pipelineID)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching Logging pipeline with ID: %s location %s, error: %w", pipelineID, rs.Primary.Attributes["location"], err)
		}
		pipeline = &pipelineResponse
		return nil
	}
}

const LoggingPipelineDataSourceMatchById = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	id = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.id
	location = "es/vit"
}
`

const LoggingPipelineDataSourceMatchByName = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	name = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.name
	location = "es/vit"
}
`

const LoggingPipelineDataSourceInvalidBothIDAndName = LoggingPipelineConfig + `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + `{
	id = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.id
	name = ` + constant.LoggingPipelineResource + `.` + constant.LoggingPipelineTestResourceName + `.name
	location = "es/vit"
}
`

const LoggingPipelineDataSourceInvalidNoIDNoName = `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + ` {
	location = "es/vit"
}
`

const LoggingPipelineDataSourceWrongNameError = `
` + constant.DataSource + ` ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestDataSourceName + ` {
	name = "nonexistent"
	location = "es/vit"
}
`

const LoggingPipelineConfigUpdate = `
resource ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestResourceName + ` {
	` + nameAttribute + ` = "` + pipelineNameUpdatedValue + `"
	location = "es/vit"
	` + pipelineLogUpdated + `
}
`

const LoggingPipelineConfigUpdateNoLabels = `
resource ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestResourceName + ` {
	` + nameAttribute + ` = "` + pipelineNameUpdatedValue + `"
	location = "es/vit"
	` + pipelineLogUpdatedNoLabels + `
}
`
