package dns

import (
	"context"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var recordResourceName = "DNS Record"

// CreateRecord creates a new record
func (c *Client) CreateRecord(ctx context.Context, zoneId string, d *schema.ResourceData) (recordResponse dns.RecordRead, responseInfo *shared.APIResponse, err error) {
	recordUUID := uuidgen.ResourceUuid()
	request := setRecordPutRequest(d)
	recordResponse, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsPut(ctx, zoneId, recordUUID.String()).RecordEnsure(request).Execute()
	apiResponse.LogInfo()
	return recordResponse, apiResponse, err
}

// IsRecordCreated checks if record is created
func (c *Client) IsRecordCreated(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneId := d.Get("zone_id").(string)
	recordID := d.Id()
	record, _, err := c.GetRecordById(ctx, zoneId, recordID)
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] record state: %s", record.Metadata.State)

	return strings.EqualFold((string)(record.Metadata.State), (string)(dns.PROVISIONINGSTATE_AVAILABLE)), nil
}

// GetRecordById gets a record by ID
func (c *Client) GetRecordById(ctx context.Context, zoneId, recordID string) (dns.RecordRead, *shared.APIResponse, error) {
	record, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsFindById(ctx, zoneId, recordID).Execute()
	apiResponse.LogInfo()
	return record, apiResponse, err
}

// ListRecords lists records
func (c *Client) ListRecords(ctx context.Context, recordName string) (dns.RecordReadList, *shared.APIResponse, error) {
	request := c.sdkClient.RecordsApi.RecordsGet(ctx)
	if recordName != "" {
		request = request.FilterName(recordName)
	}
	records, apiResponse, err := c.sdkClient.RecordsApi.RecordsGetExecute(request)
	apiResponse.LogInfo()
	return records, apiResponse, err
}

func (c *Client) SetRecordData(d *schema.ResourceData, record dns.RecordRead) error {
	d.SetId(record.Id)

	if err := d.Set("name", record.Properties.Name); err != nil {
		return utils.GenerateSetError(recordResourceName, "name", err)
	}

	if err := d.Set("type", record.Properties.Type); err != nil {
		return utils.GenerateSetError(recordResourceName, "type", err)
	}

	if err := d.Set("content", record.Properties.Content); err != nil {
		return utils.GenerateSetError(recordResourceName, "content", err)
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

	if record.Properties.Priority != nil {
		if err := d.Set("priority", *record.Properties.Priority); err != nil {
			return utils.GenerateSetError(recordResourceName, "priority", err)
		}
	}

	if err := d.Set("fqdn", record.Metadata.Fqdn); err != nil {
		return utils.GenerateSetError(recordResourceName, "fqdn", err)
	}

	return nil
}

func (c *Client) DeleteRecord(ctx context.Context, zoneId, recordID string) (*shared.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsDelete(ctx, zoneId, recordID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsRecordDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneId := d.Get("zone_id").(string)
	recordID := d.Id()
	_, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsFindById(ctx, zoneId, recordID).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) UpdateRecord(ctx context.Context, zoneId, recordID string, d *schema.ResourceData) (dns.RecordRead, *shared.APIResponse, error) {
	request := setRecordPutRequest(d)
	recordResponse, apiResponse, err := c.sdkClient.RecordsApi.ZonesRecordsPut(ctx, zoneId, recordID).RecordEnsure(request).Execute()
	apiResponse.LogInfo()
	return recordResponse, apiResponse, err
}

func setRecordPutRequest(d *schema.ResourceData) dns.RecordEnsure {
	request := dns.RecordEnsure{
		Properties: dns.Record{},
	}

	// tread carefully, workaround for setting empty name
	isNull := d.GetRawConfig().AsValueMap()["name"].IsNull()
	if nameValue, _ := d.GetOk("name"); !isNull {
		name := nameValue.(string)
		request.Properties.Name = name
	}

	if typeValue, ok := d.GetOk("type"); ok {
		typeString := typeValue.(string)
		request.Properties.Type = dns.RecordType(typeString)
	}

	if contentValue, ok := d.GetOk("content"); ok {
		content := contentValue.(string)
		request.Properties.Content = content
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

	if enabledValue, ok := d.GetOkExists("enabled"); ok { //nolint:staticcheck
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return request
}
