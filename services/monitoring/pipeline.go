package monitoring

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	monitoringSDK "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
)

var (
	ionosAPIURLMonitoring = "IONOS_API_URL_MONITORING"
	locationToURL         = map[string]string{
		"":       "https://monitoring.de-fra.ionos.com",
		"de/fra": "https://monitoring.de-fra.ionos.com",
		"de/txl": "https://monitoring.de-txl.ionos.com",
		"es/vit": "https://monitoring.es-vit.ionos.com",
		"gb/lhr": "https://monitoring.gb-lhr.ionos.com",
		"fr/par": "https://monitoring.fr-par.ionos.com",
	}
)

func (c *MonitoringClient) changeConfigURL(location string) {
	config := c.Client.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLMonitoring) != "" {
		config.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLMonitoring)),
			},
		}
		return
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

// CreatePipeline creates a new pipeline.
func (c *MonitoringClient) CreatePipeline(ctx context.Context, createReq monitoringSDK.PipelineCreate, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	pipeline, apiResponse, err := c.Client.PipelinesApi.PipelinesPost(ctx).PipelineCreate(createReq).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// DeletePipeline deletes a pipeline using its ID.
func (c *MonitoringClient) DeletePipeline(ctx context.Context, pipelineID, location string) (*shared.APIResponse, error) {
	c.changeConfigURL(location)
	apiResponse, err := c.Client.PipelinesApi.PipelinesDelete(ctx, pipelineID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdatePipeline updates a pipeline using its ID.
func (c *MonitoringClient) UpdatePipeline(ctx context.Context, updateReq monitoringSDK.PipelineEnsure, pipelineID, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	pipeline, apiResponse, err := c.Client.PipelinesApi.PipelinesPut(ctx, pipelineID).PipelineEnsure(updateReq).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// GetPipelineByID retrieves a pipeline using its ID.
func (c *MonitoringClient) GetPipelineByID(ctx context.Context, pipelineID, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	pipeline, apiResponse, err := c.Client.PipelinesApi.PipelinesFindById(ctx, pipelineID).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// GetPipelines retrieves all pipelines from a location.
func (c *MonitoringClient) GetPipelines(ctx context.Context, location string) ([]monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	c.changeConfigURL(location)
	pipelines, apiResponse, err := c.Client.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines.Items, apiResponse, err
}

// IsPipelineReady checks if the pipeline is ready.
// backoff.Permanent is used to stop the retry.
func (c *MonitoringClient) IsPipelineReady(ctx context.Context, pipelineID, location string) error {
	pipeline, _, err := c.GetPipelineByID(ctx, pipelineID, location)
	if err != nil {
		return backoff.Permanent(err)
	}
	log.Printf("[DEBUG] Monitoring pipeline state: %s", pipeline.Metadata.Status)

	if strings.EqualFold(pipeline.Metadata.Status, constant.Available) {
		return nil
	}
	return fmt.Errorf("pipeline is not ready, current state: %s", pipeline.Metadata.Status)
}

// IsPipelineDeleted checks if the pipeline is deleted.
// backoff.Permanent is used to stop the retry.
func (c *MonitoringClient) IsPipelineDeleted(ctx context.Context, pipelineID, location string) error {
	_, apiResponse, err := c.GetPipelineByID(ctx, pipelineID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return backoff.Permanent(fmt.Errorf("check failed for Monitoring pipeline with ID: %v, error: %w", pipelineID, err))
	}
	return fmt.Errorf("Monitoring pipeline with ID: %s is not deleted yet, pipeline location: %s", pipelineID, location)
}
