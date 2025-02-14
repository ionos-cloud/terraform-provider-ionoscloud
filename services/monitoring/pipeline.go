package monitoring

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"log"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	monitoringSDK "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
)

// CreatePipeline creates a new pipeline.
func (c *Client) CreatePipeline(ctx context.Context, createReq monitoringSDK.PipelineCreate, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Monitoring, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPost(ctx).PipelineCreate(createReq).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// DeletePipeline deletes a pipeline using its ID.
func (c *Client) DeletePipeline(ctx context.Context, pipelineID, location string) (*shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Monitoring, location)
	apiResponse, err := c.sdkClient.PipelinesApi.PipelinesDelete(ctx, pipelineID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdatePipeline updates a pipeline using its ID.
func (c *Client) UpdatePipeline(ctx context.Context, updateReq monitoringSDK.PipelineEnsure, pipelineID, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Monitoring, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPut(ctx, pipelineID).PipelineEnsure(updateReq).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// GetPipelineByID retrieves a pipeline using its ID.
func (c *Client) GetPipelineByID(ctx context.Context, pipelineID, location string) (monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Monitoring, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, pipelineID).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// GetPipelines retrieves all pipelines from a location.
func (c *Client) GetPipelines(ctx context.Context, location string) ([]monitoringSDK.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Monitoring, location)
	pipelines, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines.Items, apiResponse, err
}

// IsPipelineReady checks if the pipeline is ready.
// backoff.Permanent is used to stop the retry.
func (c *Client) IsPipelineReady(ctx context.Context, pipelineID, location string) error {
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
func (c *Client) IsPipelineDeleted(ctx context.Context, pipelineID, location string) error {
	_, apiResponse, err := c.GetPipelineByID(ctx, pipelineID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return backoff.Permanent(fmt.Errorf("check failed for Monitoring pipeline with ID: %v, error: %w", pipelineID, err))
	}
	return fmt.Errorf("monitoring pipeline with ID: %s is not deleted yet, pipeline location: %s", pipelineID, location)
}
