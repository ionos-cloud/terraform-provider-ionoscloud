package logging

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var pipelineResourceName = "Logging Pipeline"

// setOverrides sets the overrides for the client. searches in loaded config, or in location from plan
// if location is empty, it will search in the environment variable ionosAPIURLLogging
func (c *Client) setOverrides(location string) {
	if location == "" && os.Getenv(ionosAPIURLLogging) != "" {
		c.GetConfig().Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLLogging)),
			},
		}
	} else {
		loadedconfig.OverrideClientEndpoint(c, shared.Logging, location)
	}
}

// CreatePipeline creates a new pipeline
func (c *Client) CreatePipeline(ctx context.Context, d *schema.ResourceData) (logging.ProvisioningPipeline, utils.ApiResponseInfo, error) {
	location := d.Get("location").(string)
	loadedconfig.OverrideClientEndpoint(c, shared.Logging, location)
	request := setPipelinePostRequest(d)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPost(ctx).Pipeline(*request).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// IsPipelineAvailable checks if the pipeline is available
func (c *Client) IsPipelineAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	pipelineID := d.Id()
	location := d.Get("location").(string)
	pipeline, _, err := c.GetPipelineByID(ctx, location, pipelineID)
	if err != nil {
		return false, err
	}
	if pipeline.Metadata == nil || pipeline.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for pipeline with ID: %s", pipelineID)
	}
	log.Printf("[DEBUG] pipeline status: %s", *pipeline.Metadata.State)
	return strings.EqualFold(*pipeline.Metadata.State, constant.Available), nil
}

// UpdatePipeline updates a pipeline
func (c *Client) UpdatePipeline(ctx context.Context, id string, d *schema.ResourceData) (logging.Pipeline, utils.ApiResponseInfo, error) {
	loadedconfig.OverrideClientEndpoint(c, shared.Logging, d.Get("location").(string))
	request := setPipelinePatchRequest(d)
	pipelineResponse, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPatch(ctx, id).Pipeline(*request).Execute()
	apiResponse.LogInfo()
	return pipelineResponse, apiResponse, err
}

// DeletePipeline deletes a pipeline
func (c *Client) DeletePipeline(ctx context.Context, location, id string) (utils.ApiResponseInfo, error) {
	c.setOverrides(location)
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsPipelineDeleted checks if the pipeline is deleted
func (c *Client) IsPipelineDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	loadedconfig.OverrideClientEndpoint(c, shared.Logging, d.Get("location").(string))
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetPipelineByID returns a pipeline by its ID
func (c *Client) GetPipelineByID(ctx context.Context, location, id string) (logging.Pipeline, *shared.APIResponse, error) {
	loadedconfig.OverrideClientEndpoint(c, shared.Logging, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// ListPipelines returns a list of all pipelines
func (c *Client) ListPipelines(ctx context.Context, location string) (logging.PipelineListResponse, *shared.APIResponse, error) {
	loadedconfig.OverrideClientEndpoint(c, shared.Logging, location)
	pipelines, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines, apiResponse, err
}

func setPipelinePostRequest(d *schema.ResourceData) *logging.PipelineCreate {
	request := logging.PipelineCreate{Properties: logging.PipelineCreateProperties{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = name
	}

	var logs []logging.PipelineCreatePropertiesLogs
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if logElem, ok := logData.(map[string]interface{}); ok {
				// Populate the logElem entry.
				logSource := logElem["source"].(string)
				logTag := logElem["tag"].(string)
				logProtocol := logElem["protocol"].(string)
				newLog := *logging.NewPipelineCreatePropertiesLogs()
				newLog.Source = &logSource
				newLog.Tag = &logTag
				newLog.Protocol = &logProtocol

				// Logic for destinations.
				var destinations []logging.Destination
				for _, destinationData := range logElem["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewDestination()
						newDestination.Type = &destinationType
						newDestination.RetentionInDays = &retentionInDays
						destinations = append(destinations, newDestination)
					}
				}
				newLog.Destinations = destinations
				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = logs

	return &request
}

func setPipelinePatchRequest(d *schema.ResourceData) *logging.PipelinePatch {
	request := logging.PipelinePatch{Properties: logging.PipelinePatchProperties{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = &name
	}

	var logs []logging.PipelineCreatePropertiesLogs
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if logElem, ok := logData.(map[string]interface{}); ok {
				// Populate the logElem entry.
				logSource := logElem["source"].(string)
				logTag := logElem["tag"].(string)
				logProtocol := logElem["protocol"].(string)
				newLog := *logging.NewPipelineCreatePropertiesLogs()
				newLog.Source = &logSource
				newLog.Tag = &logTag
				newLog.Protocol = &logProtocol

				// Logic for destinations.
				var destinations []logging.Destination
				for _, destinationData := range logElem["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewDestination()
						newDestination.Type = &destinationType
						newDestination.RetentionInDays = &retentionInDays
						destinations = append(destinations, newDestination)
					}
				}
				newLog.Destinations = destinations
				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = logs

	return &request
}

// SetPipelineData sets the pipeline data
func (c *Client) SetPipelineData(d *schema.ResourceData, pipeline logging.Pipeline) error {
	d.SetId(*pipeline.Id)

	if pipeline.Properties == nil {
		return fmt.Errorf("expected properties in the response for the Logging pipeline with ID: %s, but received 'nil' instead", *pipeline.Id)
	}

	if pipeline.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the Logging pipeline with ID: %s, but received 'nil' instead", *pipeline.Id)
	}

	if pipeline.Properties.Name != nil {
		if err := d.Set("name", *pipeline.Properties.Name); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "name", err)
		}
	}

	if pipeline.Properties.GrafanaAddress != nil {
		if err := d.Set("grafana_address", *pipeline.Properties.GrafanaAddress); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "grafana_address", err)
		}
	}

	if pipeline.Properties.Logs != nil {
		logs := make([]interface{}, len(pipeline.Properties.Logs))
		for i, logElem := range pipeline.Properties.Logs {
			// Populate the logElem entry.
			logEntry := make(map[string]interface{})
			logEntry["source"] = *logElem.Source
			logEntry["tag"] = *logElem.Tag
			logEntry["protocol"] = *logElem.Protocol
			logEntry["public"] = *logElem.Public

			// Logic for destinations
			destinations := make([]interface{}, len(logElem.Destinations))
			for i, destination := range logElem.Destinations {
				destinationEntry := make(map[string]interface{})
				destinationEntry["type"] = *destination.Type
				destinationEntry["retention_in_days"] = *destination.RetentionInDays
				destinations[i] = destinationEntry
			}
			logEntry["destinations"] = destinations
			logs[i] = logEntry
		}
		if err := d.Set("log", logs); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "log", err)
		}
	}

	return nil
}
