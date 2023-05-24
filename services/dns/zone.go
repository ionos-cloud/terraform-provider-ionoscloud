package dns

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/uuidgen"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var zoneResourceName = "DNS Zone"

func (c *Client) CreateZone(ctx context.Context, d *schema.ResourceData) (zoneResponse dns.ZoneResponse, responseInfo utils.ApiResponseInfo, err error) {
	zoneUuid := uuidgen.ResourceUuid()
	request := setZonePutRequest(d)
	responseData, apiResponse, err := c.sdkClient.ZonesApi.ZonesPut(ctx, zoneUuid.String()).ZoneUpdateRequest(*request).Execute()
	apiResponse.LogInfo()
	return responseData, apiResponse, err
}

func (c *Client) IsZoneCreated(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneId := d.Id()
	zone, _, err := c.GetZoneById(ctx, zoneId)
	if err != nil {
		return false, err
	}
	if zone.Metadata == nil || zone.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for zone with ID: %s", zoneId)
	}
	log.Printf("[DEBUG] zone state: %s", *zone.Metadata.State)

	return strings.EqualFold((string)(*zone.Metadata.State), (string)(dns.CREATED)), nil
}

func (c *Client) GetZoneById(ctx context.Context, id string) (dns.ZoneResponse, *dns.APIResponse, error) {
	zone, apiResponse, err := c.sdkClient.ZonesApi.ZonesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return zone, apiResponse, err
}

func (c *Client) ListZones(ctx context.Context, filterName string) (dns.ZonesResponse, *dns.APIResponse, error) {
	request := c.sdkClient.ZonesApi.ZonesGet(ctx)
	if filterName != "" {
		request = request.FilterZoneName(filterName)
	}
	zones, apiResponse, err := c.sdkClient.ZonesApi.ZonesGetExecute(request)
	apiResponse.LogInfo()
	return zones, apiResponse, err
}

func (c *Client) SetZoneData(d *schema.ResourceData, zone dns.ZoneResponse) error {
	if zone.Id != nil {
		d.SetId(*zone.Id)
	}

	if zone.Properties == nil {
		return fmt.Errorf("expected properties in the zone response for the zone with ID: %s, but received 'nil' instead", *zone.Id)
	}

	if zone.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the zone with ID: %s, but received 'nil' instead", *zone.Id)
	}

	if zone.Properties.ZoneName != nil {
		if err := d.Set("name", *zone.Properties.ZoneName); err != nil {
			return utils.GenerateSetError(zoneResourceName, "name", err)
		}
	}

	if zone.Properties.Description != nil {
		if err := d.Set("description", *zone.Properties.Description); err != nil {
			return utils.GenerateSetError(zoneResourceName, "description", err)
		}
	}

	if zone.Properties.Enabled != nil {
		if err := d.Set("enabled", *zone.Properties.Enabled); err != nil {
			return utils.GenerateSetError(zoneResourceName, "enabled", err)
		}
	}

	if zone.Metadata.Nameservers != nil {
		if err := d.Set("nameservers", *zone.Metadata.Nameservers); err != nil {
			return utils.GenerateSetError(zoneResourceName, "nameservers", err)
		}
	}

	return nil
}

func (c *Client) UpdateZone(ctx context.Context, id string, d *schema.ResourceData) (dns.ZoneResponse, utils.ApiResponseInfo, error) {
	request := setZonePutRequest(d)
	zoneResponse, apiResponse, err := c.sdkClient.ZonesApi.ZonesPut(ctx, id).ZoneUpdateRequest(*request).Execute()
	apiResponse.LogInfo()
	return zoneResponse, apiResponse, err
}

func (c *Client) DeleteZone(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.ZonesApi.ZonesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsZoneDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.ZonesApi.ZonesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func setZoneCreateRequest(d *schema.ResourceData) *dns.ZoneCreateRequest {
	request := dns.ZoneCreateRequest{
		Properties: &dns.ZoneCreateRequestProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.ZoneName = &name
	}

	if descriptionValue, ok := d.GetOk("description"); ok {
		description := descriptionValue.(string)
		request.Properties.Description = &description
	}

	if enabledValue, ok := d.GetOkExists("enabled"); ok {
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return &request
}

func setZonePutRequest(d *schema.ResourceData) *dns.ZoneUpdateRequest {
	request := dns.ZoneUpdateRequest{
		Properties: &dns.ZoneUpdateRequestProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.ZoneName = &name
	}

	if descriptionValue, ok := d.GetOk("description"); ok {
		description := descriptionValue.(string)
		request.Properties.Description = &description
	}

	if enabledValue, ok := d.GetOkExists("enabled"); ok {
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return &request
}
