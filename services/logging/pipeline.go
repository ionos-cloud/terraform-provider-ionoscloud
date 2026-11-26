package logging

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

var pipelineResourceName = "Logging Pipeline"

// ErrPipelineFailed is returned while waiting for a pipeline that ended up in the FAILED state
var ErrPipelineFailed = errors.New("logging pipeline is in FAILED state")

// CreatePipeline creates a new pipeline
func (c *Client) CreatePipeline(ctx context.Context, d *schema.ResourceData) (logging.PipelineRead, *shared.APIResponse, error) {
	location := d.Get("location").(string)
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, location)
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
	tflog.Debug(ctx, "logging pipeline status", map[string]any{"state": pipeline.Metadata.State})
	if strings.EqualFold(pipeline.Metadata.State, constant.Failed) {
		return false, fmt.Errorf("%w, pipeline ID: %s, status message: %q", ErrPipelineFailed, pipelineID, pipeline.Metadata.GetStatusMessage())
	}
	return strings.EqualFold(pipeline.Metadata.State, constant.Available), nil
}

// UpdatePipeline updates a pipeline
func (c *Client) UpdatePipeline(ctx context.Context, id string, d *schema.ResourceData) (logging.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, d.Get("location").(string))
	request := setPipelinePatchRequest(d)
	pipelineResponse, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesPatch(ctx, id).PipelinePatch(*request).Execute()
	apiResponse.LogInfo()
	return pipelineResponse, apiResponse, err
}

// DeletePipeline deletes a pipeline
func (c *Client) DeletePipeline(ctx context.Context, location, id string) (*shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, location)
	apiResponse, err := c.sdkClient.PipelinesApi.PipelinesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsPipelineDeleted checks if the pipeline is deleted
func (c *Client) IsPipelineDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, d.Get("location").(string))
	_, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetPipelineByID returns a pipeline by its ID
func (c *Client) GetPipelineByID(ctx context.Context, location, id string) (logging.PipelineRead, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, location)
	pipeline, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return pipeline, apiResponse, err
}

// ListPipelines returns a list of all pipelines
func (c *Client) ListPipelines(ctx context.Context, location string) (logging.PipelineReadList, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(ctx, c, fileconfiguration.Logging, location)
	pipelines, apiResponse, err := c.sdkClient.PipelinesApi.PipelinesGet(ctx).Execute()
	apiResponse.LogInfo()
	return pipelines, apiResponse, err
}

func setPipelinePostRequest(d *schema.ResourceData) *logging.PipelineCreate {
	request := logging.PipelineCreate{Properties: logging.PipelineNoAddr{}}

	if nameValue, ok := d.GetOk("name"); ok {
		request.Properties.Name = nameValue.(string)
	}

	logsValue, _ := d.Get("log").([]any)
	request.Properties.Logs = expandPipelineLogs(logsValue, nil)

	return &request
}

func setPipelinePatchRequest(d *schema.ResourceData) *logging.PipelinePatch {
	request := logging.PipelinePatch{Properties: logging.PipelineNoAddr{}}

	if nameValue, ok := d.GetOk("name"); ok {
		request.Properties.Name = nameValue.(string)
	}

	oldLogsValue, newLogsValue := d.GetChange("log")
	oldLogs, _ := oldLogsValue.([]any)
	newLogs, _ := newLogsValue.([]any)
	request.Properties.Logs = expandPipelineLogs(newLogs, oldLogs)

	return &request
}

// expandPipelineLogs converts the 'log' list from the Terraform configuration into SDK log objects.
// oldLogs is the previous value of 'log' (nil on create). Labels are only sent when configured, except when
// the log with the same tag previously had labels and now has none: then an empty list is sent so the
// labels are removed in the API. Logs are matched by tag, which is unique within a pipeline, so removing or
// reordering other logs does not affect the match.
func expandPipelineLogs(logs, oldLogs []any) []logging.PipelineNoAddrLogs {
	oldTagsWithLabels := make(map[string]bool, len(oldLogs))
	for _, oldLogData := range oldLogs {
		if oldLogElem, ok := oldLogData.(map[string]any); ok && len(logLabels(oldLogElem)) > 0 {
			if tag, ok := oldLogElem["tag"].(string); ok {
				oldTagsWithLabels[tag] = true
			}
		}
	}

	var result []logging.PipelineNoAddrLogs
	for _, logData := range logs {
		logElem, ok := logData.(map[string]any)
		if !ok {
			continue
		}

		var destinations []logging.PipelineNoAddrLogsDestinations
		destinationsValue, _ := logElem["destinations"].([]any)
		for _, destinationData := range destinationsValue {
			if destination, ok := destinationData.(map[string]any); ok {
				destinationType := destination["type"].(string)
				retentionInDays := int32(destination["retention_in_days"].(int)) //nolint:gosec // the schema only allows 7, 14 or 30
				destinations = append(destinations, *logging.NewPipelineNoAddrLogsDestinations(destinationType, retentionInDays))
			}
		}

		newLog := *logging.NewPipelineNoAddrLogs(logElem["source"].(string), logElem["tag"].(string), logElem["protocol"].(string), destinations)
		if labels := logLabels(logElem); len(labels) > 0 {
			newLog.Labels = labels
		} else if oldTagsWithLabels[newLog.Tag] {
			newLog.Labels = []string{}
		}
		result = append(result, newLog)
	}
	return result
}

// logLabels returns the 'labels' of a single 'log' element as a sorted string slice.
func logLabels(logElem map[string]any) []string {
	labelsSet, ok := logElem["labels"].(*schema.Set)
	if !ok || labelsSet == nil {
		return nil
	}
	labels := make([]string, 0, labelsSet.Len())
	for _, label := range labelsSet.List() {
		labels = append(labels, label.(string))
	}
	slices.Sort(labels)
	return labels
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

	if err := d.Set("status_message", pipeline.Metadata.GetStatusMessage()); err != nil {
		return utils.GenerateSetError(pipelineResourceName, "status_message", err)
	}

	if pipeline.Properties.Logs != nil {
		logs := make([]any, len(pipeline.Properties.Logs))
		for i, logElem := range pipeline.Properties.Logs {
			// Populate the logElem entry.
			logEntry := make(map[string]any)
			utils.SetPropWithNilCheck(logEntry, "source", logElem.Source)
			utils.SetPropWithNilCheck(logEntry, "tag", logElem.Tag)
			utils.SetPropWithNilCheck(logEntry, "protocol", logElem.Protocol)
			logEntry["labels"] = utils.ToInterfaceSlice(logElem.Labels)
			// Logic for destinations
			destinations := make([]any, len(logElem.Destinations))
			for i, destination := range logElem.Destinations {
				destinationEntry := make(map[string]any)
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
