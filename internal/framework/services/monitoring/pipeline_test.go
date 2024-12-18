//go:build all || monitoring
// +build all monitoring

package monitoring_test

import (
	"context"
	"fmt"
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
				),
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
		},
	})
}

const (
	pipelineBasicConfig = `
	resource "ionoscloud_monitoring_pipeline" "test" {
		name = "TFTestPipeline"
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
