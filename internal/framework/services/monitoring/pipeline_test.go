//go:build all || monitoring
// +build all monitoring

package monitoring_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccPipeline(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: pipelineBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					checkPipelineExists(context.Background(), "ionoscloud_monitoring_pipeline.test"),
					resource.TestCheckResourceAttr("ionoscloud_monitoring_pipeline.test", "name", "TFTestPipeline"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "grafana_endpoint"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "http_endpoint"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "key"),
				),
			},
			{
				ResourceName:            "ionoscloud_monitoring_pipeline.test",
				ImportState:             true,
				ImportStateIdFunc:       monitoringImportStateID,
				ImportStateVerifyIgnore: []string{"key", "location"},
				ImportStateVerify:       true,
			},
			{
				Config: dataSourceByName,
				Check: resource.ComposeTestCheckFunc(
					checkPipelineExists(context.Background(), "ionoscloud_monitoring_pipeline.test"),
					resource.TestCheckResourceAttr("data.ionoscloud_monitoring_pipeline.test", "name", "TFTestPipeline"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_monitoring_pipeline.test", "grafana_endpoint"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_monitoring_pipeline.test", "http_endpoint"),
				),
			},
			{
				Config: dataSourceByID,
				Check: resource.ComposeTestCheckFunc(
					checkPipelineExists(context.Background(), "ionoscloud_monitoring_pipeline.test"),
					resource.TestCheckResourceAttr("data.ionoscloud_monitoring_pipeline.test", "name", "TFTestPipeline"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_monitoring_pipeline.test", "grafana_endpoint"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_monitoring_pipeline.test", "http_endpoint"),
				),
			},
			{
				Config:      invalidDSConfigBothNameID,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      invalidDSConfigNoNameNoID,
				ExpectError: regexp.MustCompile("Missing Attribute Configuration"),
			},
			{
				Config:      invalidDSConfigMultipleOccurrences,
				ExpectError: regexp.MustCompile("multiple Monitoring pipelines found with the same name"),
			},
			{
				Config:      invalidDSConfigInvalidName,
				ExpectError: regexp.MustCompile("no Monitoring pipeline found with the specified name"),
			},
			{
				Config:      invalidDSConfigInvalidID,
				ExpectError: regexp.MustCompile("failed to get Monitoring pipeline"),
			},
			{
				Config: pipelineBasicUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkPipelineExists(context.Background(), "ionoscloud_monitoring_pipeline.test"),
					resource.TestCheckResourceAttr("ionoscloud_monitoring_pipeline.test", "name", "updatedTestName"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "grafana_endpoint"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "http_endpoint"),
					resource.TestCheckResourceAttrSet("ionoscloud_monitoring_pipeline.test", "key"),
				),
			},
		},
	})
}

const (
	pipelineBasicConfig = `
	resource "ionoscloud_monitoring_pipeline" "test" {
		name = "TFTestPipeline"
	}
`

	pipelineBasicUpdateConfig = `
	resource "ionoscloud_monitoring_pipeline" "test" {
		name = "updatedTestName"
	}
`

	dataSourceByName = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {
		name = ionoscloud_monitoring_pipeline.test.name
	}
`

	dataSourceByID = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {
		id = ionoscloud_monitoring_pipeline.test.id
	}
`

	invalidDSConfigBothNameID = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {
		id = ionoscloud_monitoring_pipeline.test.id
		name = ionoscloud_monitoring_pipeline.test.name
	}
`

	invalidDSConfigNoNameNoID = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {

	}
`

	// This works because we can create two pipelines with the same name
	invalidDSConfigMultipleOccurrences = pipelineBasicConfig + `
	resource "ionoscloud_monitoring_pipeline" "secondPipeline" {
		name = "TFTestPipeline"
	}

	data "ionoscloud_monitoring_pipeline" "test" {
		name = ionoscloud_monitoring_pipeline.test.name
	}
`

	invalidDSConfigInvalidName = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {
		name = "itdoesntexist"
	}
`

	invalidDSConfigInvalidID = pipelineBasicConfig + `
	data "ionoscloud_monitoring_pipeline" "test" {
		id = "itdoesntexist"
	}
`
)

func checkPipelineExists(ctx context.Context, accessPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[accessPath]
		if !ok {
			return fmt.Errorf("Not found: %s", accessPath)
		}
		client := acctest.MonitoringClient()
		_, _, err := client.GetPipelineByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching Monitoring pipeline with ID: %v, error: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func checkPipelineDestroy(s *terraform.State) error {
	client := acctest.MonitoringClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_monitoring_pipeline" {
			continue
		}
		_, apiResponse, err := client.GetPipelineByID(context.Background(), rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking the destruction of Monitoring pipeline with ID: %v, error: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("Monitoring pipeline with ID: %v still exists", rs.Primary.ID)
		}
	}
	return nil
}

func monitoringImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_monitoring_pipeline" {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
