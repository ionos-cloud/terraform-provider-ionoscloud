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

var zoneResourceName = "DNS Zone"

// CreateZone creates a new DNS Zone
func (c *Client) CreateZone(ctx context.Context, d *schema.ResourceData) (zoneResponse dns.ZoneRead, responseInfo *shared.APIResponse, err error) {
	zoneUuid := uuidgen.ResourceUuid()
	request := setZonePutRequest(d)
	responseData, apiResponse, err := c.sdkClient.ZonesApi.ZonesPut(ctx, zoneUuid.String()).ZoneEnsure(*request).Execute()
	apiResponse.LogInfo()
	return responseData, apiResponse, err
}

// IsZoneCreated checks if a zone is created
func (c *Client) IsZoneCreated(ctx context.Context, d *schema.ResourceData) (bool, error) {
	zoneID := d.Id()
	zone, _, err := c.GetZoneById(ctx, zoneID)
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] zone state: %s", zone.Metadata.State)

	return strings.EqualFold((string)(zone.Metadata.State), (string)(dns.PROVISIONINGSTATE_AVAILABLE)), nil
}

// GetZoneById gets a zone by ID
func (c *Client) GetZoneById(ctx context.Context, id string) (dns.ZoneRead, *shared.APIResponse, error) {
	zone, apiResponse, err := c.sdkClient.ZonesApi.ZonesFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return zone, apiResponse, err
}

// ListZones lists all zones
func (c *Client) ListZones(ctx context.Context, filterName string) (dns.ZoneReadList, *shared.APIResponse, error) {
	request := c.sdkClient.ZonesApi.ZonesGet(ctx)
	if filterName != "" {
		request = request.FilterZoneName(filterName)
	}
	zones, apiResponse, err := c.sdkClient.ZonesApi.ZonesGetExecute(request)
	apiResponse.LogInfo()
	return zones, apiResponse, err
}

// SetZoneData sets the data of a zone
func (c *Client) SetZoneData(d *schema.ResourceData, zone dns.ZoneRead) error {
	d.SetId(zone.Id)

	if err := d.Set("name", zone.Properties.ZoneName); err != nil {
		return utils.GenerateSetError(zoneResourceName, "name", err)
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

	if err := d.Set("nameservers", zone.Metadata.Nameservers); err != nil {
		return utils.GenerateSetError(zoneResourceName, "nameservers", err)
	}

	return nil
}

// UpdateZone updates a zone
func (c *Client) UpdateZone(ctx context.Context, id string, d *schema.ResourceData) (dns.ZoneRead, *shared.APIResponse, error) {
	request := setZonePutRequest(d)
	zoneResponse, apiResponse, err := c.sdkClient.ZonesApi.ZonesPut(ctx, id).ZoneEnsure(*request).Execute()
	apiResponse.LogInfo()
	return zoneResponse, apiResponse, err
}

// DeleteZone deletes a zone
func (c *Client) DeleteZone(ctx context.Context, id string) (*shared.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.ZonesApi.ZonesDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsZoneDeleted checks if a zone is deleted
func (c *Client) IsZoneDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.ZonesApi.ZonesFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func setZonePutRequest(d *schema.ResourceData) *dns.ZoneEnsure {
	request := dns.ZoneEnsure{
		Properties: dns.Zone{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		request.Properties.ZoneName = name
	}

	if descriptionValue, ok := d.GetOk("description"); ok {
		description := descriptionValue.(string)
		request.Properties.Description = &description
	}

	if enabledValue, ok := d.GetOkExists("enabled"); ok { //nolint:staticcheck
		enabled := enabledValue.(bool)
		request.Properties.Enabled = &enabled
	}
	return &request
}
