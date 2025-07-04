package logging

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

var pipelineResourceName = "Logging Pipeline"

// CreatePipeline creates a new pipeline
func (c *Client) CreatePipeline(ctx context.Context, d *schema.ResourceData) (logging.PipelineRead, utils.ApiResponseInfo, error) {
	location := d.Get("location").(string)
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, location)
	request := setPipelinePostRequest(d)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPost(ctx).PipelineCreate(*request).Execute()
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
	log.Printf("[DEBUG] pipeline status: %s", pipeline.Metadata.State)
	return strings.EqualFold(pipeline.Metadata.State, constant.Available), nil
}

// UpdatePipeline updates a pipeline
func (c *Client) UpdatePipeline(ctx context.Context, id string, d *schema.ResourceData) (logging.PipelineRead, utils.ApiResponseInfo, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, d.Get("location").(string))
	request := setPipelinePatchRequest(d)
	pipelineResponse, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPatch(ctx, id).PipelinePatch(*request).Execute()
	apiResponse.LogInfo()
	return pipelineResponse, apiResponse, err
}

// DeletePipeline deletes a pipeline
func (c *Client) DeletePipeline(ctx context.Context, location, id string) (utils.ApiResponseInfo, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, location)
	apiResponse, err := c.sdkClient.PipelinesApi.PipelinesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsPipelineDeleted checks if the pipeline is deleted
func (c *Client) IsPipelineDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, d.Get("location").(string))
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetPipelineByID returns a pipeline by its ID
func (c *Client) GetPipelineByID(ctx context.Context, location, id string) (logging.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// ListPipelines returns a list of all pipelines
func (c *Client) ListPipelines(ctx context.Context, location string) (logging.PipelineReadList, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Logging, location)
	pipelines, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines, apiResponse, err
}

func setPipelinePostRequest(d *schema.ResourceData) *logging.PipelineCreate {
	request := logging.PipelineCreate{Properties: logging.PipelineNoAddr{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = name
	}

	var logs []logging.PipelineNoAddrLogs
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if logElem, ok := logData.(map[string]interface{}); ok {
				// Populate the logElem entry.
				logSource := logElem["source"].(string)
				logTag := logElem["tag"].(string)
				logProtocol := logElem["protocol"].(string)

				// Logic for destinations.
				var destinations []logging.PipelineNoAddrLogsDestinations
				for _, destinationData := range logElem["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewPipelineNoAddrLogsDestinations(destinationType, retentionInDays)
						destinations = append(destinations, newDestination)
					}
				}
				newLog := *logging.NewPipelineNoAddrLogs(logSource, logTag, logProtocol, destinations)

				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = logs

	return &request
}

func setPipelinePatchRequest(d *schema.ResourceData) *logging.PipelinePatch {
	request := logging.PipelinePatch{Properties: logging.PipelineNoAddr{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = name
	}

	var logs []logging.PipelineNoAddrLogs
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if logElem, ok := logData.(map[string]interface{}); ok {
				// Populate the logElem entry.
				logSource := logElem["source"].(string)
				logTag := logElem["tag"].(string)
				logProtocol := logElem["protocol"].(string)

				// Logic for destinations.
				var destinations []logging.PipelineNoAddrLogsDestinations
				for _, destinationData := range logElem["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewPipelineNoAddrLogsDestinations(destinationType, retentionInDays)
						destinations = append(destinations, newDestination)
					}
				}
				newLog := *logging.NewPipelineNoAddrLogs(logSource, logTag, logProtocol, destinations)
				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = logs

	return &request
}

// SetPipelineData sets the pipeline data
func (c *Client) SetPipelineData(d *schema.ResourceData, pipeline logging.PipelineRead) error {
	d.SetId(pipeline.Id)

	if err := d.Set("name", pipeline.Properties.Name); err != nil {
		return utils.GenerateSetError(pipelineResourceName, "name", err)
	}

	if pipeline.Properties.GrafanaAddress != nil {
		if err := d.Set("grafana_address", *pipeline.Properties.GrafanaAddress); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "grafana_address", err)
		}
	}
	if pipeline.Properties.HttpAddress != nil && *pipeline.Properties.HttpAddress != "" {
		if err := d.Set("http_address", *pipeline.Properties.HttpAddress); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "http_address", err)
		}
	}
	if pipeline.Properties.TcpAddress != nil && *pipeline.Properties.TcpAddress != "" {
		if err := d.Set("tcp_address", *pipeline.Properties.TcpAddress); err != nil {
			return utils.GenerateSetError(pipelineResourceName, "tcp_address", err)
		}
	}

	if pipeline.Properties.Logs != nil {
		logs := make([]interface{}, len(pipeline.Properties.Logs))
		for i, logElem := range pipeline.Properties.Logs {
			// Populate the logElem entry.
			logEntry := make(map[string]interface{})
			utils.SetPropWithNilCheck(logEntry, "source", logElem.Source)
			utils.SetPropWithNilCheck(logEntry, "tag", logElem.Tag)
			utils.SetPropWithNilCheck(logEntry, "protocol", logElem.Protocol)
			// Logic for destinations
			destinations := make([]interface{}, len(logElem.Destinations))
			for i, destination := range logElem.Destinations {
				destinationEntry := make(map[string]interface{})
				utils.SetPropWithNilCheck(destinationEntry, "type", destination.Type)
				utils.SetPropWithNilCheck(destinationEntry, "retention_in_days", destination.RetentionInDays)
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
