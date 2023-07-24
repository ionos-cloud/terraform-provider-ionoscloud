package logging

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	logging "github.com/ionos-cloud/sdk-go-logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"log"
	"strings"
)

var pipelineResourceName = "Logging Pipeline"

func (c *Client) CreatePipeline(ctx context.Context, d *schema.ResourceData) (logging.Pipeline, utils.ApiResponseInfo, error) {
	request := setPipelinePostRequest(d)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPost(ctx).Pipeline(*request).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

func (c *Client) IsPipelineAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	pipelineId := d.Id()
	pipeline, _, err := c.GetPipelineById(ctx, pipelineId)
	if err != nil {
		return false, err
	}
	if pipeline.Metadata == nil || pipeline.Metadata.Status == nil {
		return false, fmt.Errorf("expected metadata, got empty for pipeline with ID: %s", pipelineId)
	}
	log.Printf("[DEBUG] pipeline status: %s", *pipeline.Metadata.Status)
	return strings.EqualFold(*pipeline.Metadata.Status, constant.Available), nil
}

func (c *Client) UpdatePipeline(ctx context.Context, id string, d *schema.ResourceData) (logging.Pipeline, utils.ApiResponseInfo, error) {
	request := setPipelinePatchRequest(d)
	pipelineResponse, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPatch(ctx, id).Pipeline(*request).Execute()
	apiResponse.LogInfo()
	return pipelineResponse, apiResponse, err
}

func (c *Client) DeletePipeline(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsPipelineDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) GetPipelineById(ctx context.Context, id string) (logging.Pipeline, *logging.APIResponse, error) {
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

func (c *Client) ListPipelines(ctx context.Context) (logging.PipelineListResponse, *logging.APIResponse, error) {
	pipelines, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines, apiResponse, err
}

func setPipelinePostRequest(d *schema.ResourceData) *logging.CreateRequest {
	request := logging.CreateRequest{Properties: &logging.CreateRequestProperties{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = &name
	}

	var logs []logging.CreateRequestPipeline
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if log, ok := logData.(map[string]interface{}); ok {
				// Populate the log entry.
				logSource := log["source"].(string)
				logTag := log["tag"].(string)
				logProtocol := log["protocol"].(string)
				newLog := *logging.NewCreateRequestPipeline()
				newLog.Source = &logSource
				newLog.Tag = &logTag
				newLog.Protocol = &logProtocol

				// Logic for destinations.
				var destinations []logging.Destination
				for _, destinationData := range log["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewDestination()
						newDestination.Type = &destinationType
						newDestination.RetentionInDays = &retentionInDays
						destinations = append(destinations, newDestination)
					}
				}
				newLog.Destinations = &destinations
				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = &logs

	return &request
}

func setPipelinePatchRequest(d *schema.ResourceData) *logging.PatchRequest {
	request := logging.PatchRequest{Properties: &logging.PatchRequestProperties{}}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = &name
	}

	var logs []logging.PatchRequestPipeline
	if logsValue, ok := d.GetOk("log"); ok {
		for _, logData := range logsValue.([]interface{}) {
			if log, ok := logData.(map[string]interface{}); ok {
				// Populate the log entry.
				logSource := log["source"].(string)
				logTag := log["tag"].(string)
				logProtocol := log["protocol"].(string)
				newLog := *logging.NewPatchRequestPipeline()
				newLog.Source = &logSource
				newLog.Tag = &logTag
				newLog.Protocol = &logProtocol

				// Logic for destinations.
				var destinations []logging.Destination
				for _, destinationData := range log["destinations"].([]interface{}) {
					if destination, ok := destinationData.(map[string]interface{}); ok {
						destinationType := destination["type"].(string)
						retentionInDays := int32(destination["retention_in_days"].(int))
						newDestination := *logging.NewDestination()
						newDestination.Type = &destinationType
						newDestination.RetentionInDays = &retentionInDays
						destinations = append(destinations, newDestination)
					}
				}
				newLog.Destinations = &destinations
				logs = append(logs, newLog)
			}
		}
	}

	request.Properties.Logs = &logs

	return &request
}

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

	if pipeline.Properties.Logs != nil {
		logs := make([]interface{}, len(*pipeline.Properties.Logs))
		for i, log := range *pipeline.Properties.Logs {
			// Populate the log entry.
			logEntry := make(map[string]interface{})
			logEntry["source"] = *log.Source
			logEntry["tag"] = *log.Tag
			logEntry["protocol"] = *log.Protocol
			logEntry["public"] = *log.Public

			// Logic for destinations
			destinations := make([]interface{}, len(*log.Destinations))
			for i, destination := range *log.Destinations {
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
