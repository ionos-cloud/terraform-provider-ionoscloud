package dns

import (
	"context"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var reverseRecordResourceName = "DNS Reverse Record"

// CreateReverseRecord creates a new reverse record
func (c *Client) CreateReverseRecord(ctx context.Context, d *schema.ResourceData) (recordResponse dns.ReverseRecordRead, responseInfo *shared.APIResponse, err error) {
	recordUUID := uuidgen.ResourceUuid()
	request := setReverseRecordPutRequest(d)
	reverseRecordResponse, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsPut(ctx, recordUUID.String()).ReverseRecordEnsure(request).Execute()
	apiResponse.LogInfo()
	return reverseRecordResponse, apiResponse, err
}

// GetReverseRecordById gets a reverse record by ID
func (c *Client) GetReverseRecordById(ctx context.Context, recordID string) (dns.ReverseRecordRead, *shared.APIResponse, error) {
	record, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsFindById(ctx, recordID).Execute()
	apiResponse.LogInfo()
	return record, apiResponse, err
}

// ListReverseRecords lists reverse records
func (c *Client) ListReverseRecords(ctx context.Context, ips []string) (dns.ReverseRecordsReadList, *shared.APIResponse, error) {
	request := c.sdkClient.ReverseRecordsApi.ReverserecordsGet(ctx)
	if ips != nil {
		request = request.FilterRecordIp(ips)
	}
	records, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsGetExecute(request)
	apiResponse.LogInfo()
	return records, apiResponse, err
}

func (c *Client) SetReverseRecordData(d *schema.ResourceData, record dns.ReverseRecordRead) error {
	d.SetId(record.Id)

	if err := d.Set("name", record.Properties.Name); err != nil {
		return utils.GenerateSetError(reverseRecordResourceName, "name", err)
	}

	if err := d.Set("description", record.Properties.Description); err != nil {
		return utils.GenerateSetError(reverseRecordResourceName, "description", err)
	}

	if err := d.Set("ip", record.Properties.Ip); err != nil {
		return utils.GenerateSetError(reverseRecordResourceName, "ip", err)
	}

	return nil
}

func (c *Client) DeleteReverseRecord(ctx context.Context, recordID string) (*shared.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsDelete(ctx, recordID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsReverseRecordDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	recordID := d.Id()
	_, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsFindById(ctx, recordID).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) UpdateReverseRecord(ctx context.Context, recordID string, d *schema.ResourceData) (dns.ReverseRecordRead, *shared.APIResponse, error) {
	request := setReverseRecordPutRequest(d)
	recordResponse, apiResponse, err := c.sdkClient.ReverseRecordsApi.ReverserecordsPut(ctx, recordID).ReverseRecordEnsure(request).Execute()
	apiResponse.LogInfo()
	return recordResponse, apiResponse, err
}

func setReverseRecordPutRequest(d *schema.ResourceData) dns.ReverseRecordEnsure {
	request := dns.ReverseRecordEnsure{
		Properties: dns.ReverseRecord{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.Name = name
	}

	if descriptionValue, ok := d.GetOk("description"); ok {
		descriptionString := descriptionValue.(string)
		request.Properties.Description = &descriptionString
	}

	if ipValue, ok := d.GetOk("ip"); ok {
		ip := ipValue.(string)
		request.Properties.Ip = ip
	}
	return request
}
