package dnsaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnsaas "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
)

var recordResourceName = "DNS Record"

func (c *Client) CreateRecord(ctx context.Context, zoneId string, d *schema.ResourceData) (recordResponse dnsaas.RecordResponse, responseInfo utils.ApiResponseInfo, err error) {
	request := setRecordCreateRequest(d)
	recordResponse, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsPost(ctx, zoneId).RecordCreateRequest(*request).Execute()
	apiResponse.LogInfo()
	return recordResponse, apiResponse, err
}

func (c *Client) IsRecordCreated(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()
	record, _, err := c.GetRecordById(ctx, zoneId, recordId)
	if err != nil {
		return false, err
	}
	if record.Metadata == nil || record.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for record with ID: %s, zone ID: %s", recordId, zoneId)
	}
	log.Printf("[DEBUG] record state: %s", *record.Metadata.State)

	return strings.EqualFold((string)(*record.Metadata.State), (string)(dnsaas.CREATED)), nil
}

func (c *Client) GetRecordById(ctx context.Context, zoneId, recordId string) (dnsaas.RecordResponse, *dnsaas.APIResponse, error) {
	record, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsFindById(ctx, zoneId, recordId).Execute()
	apiResponse.LogInfo()
	return record, apiResponse, err
}

func (c *Client) ListRecords(ctx context.Context, zoneId, recordName string) (dnsaas.RecordsResponse, *dnsaas.APIResponse, error) {
	request := c.sdkClient.RecordsApi.RecordsGet(ctx)
	if recordName != "" {
		request = request.FilterName(recordName)
	}
	records, apiResponse, err := c.sdkClient.RecordsApi.RecordsGetExecute(request)
	apiResponse.LogInfo()
	return records, apiResponse, err
}

func (c *Client) SetRecordData(d *schema.ResourceData, record dnsaas.RecordResponse) error {
	if record.Id != nil {
		d.SetId(*record.Id)
	}

	if record.Properties.Name != nil {
		if err := d.Set("name", *record.Properties.Name); err != nil {
			return utils.GenerateSetError(recordResourceName, "name", err)
		}
	}

	if record.Properties.Type != nil {
		if err := d.Set("type", *record.Properties.Type); err != nil {
			return utils.GenerateSetError(recordResourceName, "type", err)
		}
	}

	if record.Properties.Content != nil {
		if err := d.Set("content", *record.Properties.Content); err != nil {
			return utils.GenerateSetError(recordResourceName, "content", err)
		}
	}

	if record.Properties.Ttl != nil {
		if err := d.Set("ttl", *record.Properties.Ttl); err != nil {
			return utils.GenerateSetError(recordResourceName, "ttl", err)
		}
	}

	if record.Properties.Enabled != nil {
		if err := d.Set("enabled", *record.Properties.Enabled); err != nil {
			return utils.GenerateSetError(recordResourceName, "enabled", err)
		}
	}

	return nil
}

func (c *Client) DeleteRecord(ctx context.Context, zoneId, recordId string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsDelete(ctx, zoneId, recordId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsRecordDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()
	_, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsFindById(ctx, zoneId, recordId).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) UpdateRecord(ctx context.Context, zoneId, recordId string, d *schema.ResourceData) (utils.ApiResponseInfo, error) {
	request := setRecordPutRequest(d)
	_, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsPut(ctx, zoneId, recordId).RecordUpdateRequest(*request).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func setRecordPutRequest(d *schema.ResourceData) *dnsaas.RecordUpdateRequest {
	request := dnsaas.RecordUpdateRequest{
		Properties: &dnsaas.RecordProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = &name
	}

	if typeValue, ok := d.GetOk("type"); ok {
		typeString := typeValue.(string)
		recordType := (dnsaas.RecordType)(typeString)
		request.Properties.Type = &recordType
	}

	if contentValue, ok := d.GetOk("content"); ok {
		content := contentValue.(string)
		request.Properties.Content = &content
	}

	if ttlValue, ok := d.GetOk("ttl"); ok {
		ttl := ttlValue.(int)
		castedTtl := (int32)(ttl)
		request.Properties.Ttl = &castedTtl
	}

	if priorityValue, ok := d.GetOk("priority"); ok {
		priority := priorityValue.(int)
		castedPriority := (int32)(priority)
		request.Properties.Priority = &castedPriority
	}

	if enabledValue, ok := d.GetOk("enabled"); ok {
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return &request
}

func setRecordCreateRequest(d *schema.ResourceData) *dnsaas.RecordCreateRequest {
	request := dnsaas.RecordCreateRequest{
		Properties: &dnsaas.RecordProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = &name
	}

	if typeValue, ok := d.GetOk("type"); ok {
		typeString := typeValue.(string)
		recordType := (dnsaas.RecordType)(typeString)
		request.Properties.Type = &recordType
	}

	if contentValue, ok := d.GetOk("content"); ok {
		content := contentValue.(string)
		request.Properties.Content = &content
	}

	if ttlValue, ok := d.GetOk("ttl"); ok {
		ttl := ttlValue.(int)
		castedTtl := (int32)(ttl)
		request.Properties.Ttl = &castedTtl
	}

	if priorityValue, ok := d.GetOk("priority"); ok {
		priority := priorityValue.(int)
		castedPriority := (int32)(priority)
		request.Properties.Priority = &castedPriority
	}

	if enabledValue, ok := d.GetOk("enabled"); ok {
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return &request
}
