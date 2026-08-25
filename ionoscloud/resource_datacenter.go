package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	dcschema "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/compute/datacenter"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func ResourceDatacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterCreate,
		ReadContext:   resourceDatacenterRead,
		UpdateContext: resourceDatacenterUpdate,
		DeleteContext: resourceDatacenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDatacenterImport,
		},
		Identity: &schema.ResourceIdentity{
			SchemaFunc: dcschema.IdentitySchema,
		},
		Schema:   dcschema.Schema(),
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	datacenterName := d.Get("name").(string)
	datacenterLocation := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, datacenterLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenter := ionoscloud.DatacenterPost{
		Properties: &ionoscloud.DatacenterPropertiesPost{
			Name:     &datacenterName,
			Location: &datacenterLocation,
		},
	}
	if attr, ok := d.GetOk("description"); ok {
		attrStr := attr.(string)
		datacenter.Properties.Description = &attrStr
	}

	if attr, ok := d.GetOk("sec_auth_protection"); ok {
		attrStr := attr.(bool)
		datacenter.Properties.SecAuthProtection = &attrStr
	}

	createdDatacenter, apiResponse, err := client.DataCentersApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}
	d.SetId(*createdDatacenter.Id)

	tflog.Info(ctx, "datacenter created", map[string]any{"resource_id": d.Id()})

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	datacenterLocation := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, datacenterLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, err, nil)
	}

	if err := dcschema.PopulateResourceData(d, &datacenter); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceDatacenterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	obj := ionoscloud.DatacenterPropertiesPut{}

	datacenterLocation := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, datacenterLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		newNameStr := newName.(string)
		obj.Name = &newNameStr
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		newDescriptionStr := newDescription.(string)
		obj.Description = &newDescriptionStr
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")

		return diagutil.ToDiags(d, fmt.Errorf("data center is created in %s location. You can not change location of the data center to %s; it requires recreation of the data center", oldLocation, newLocation), nil)
	}

	if d.HasChange("sec_auth_protection") {
		_, newSecAuthProtection := d.GetChange("sec_auth_protection")
		newSecAuthProtectionStr := newSecAuthProtection.(bool)
		obj.SecAuthProtection = &newSecAuthProtectionStr
	}

	_, apiResponse, err := client.DataCentersApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	datacenterLocation := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, datacenterLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceDatacenterImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<datacenter-id> or <datacenter-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	dcID := parts[0]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find datacenter %q", dcID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while retrieving the datacenter: %w", err), nil)
	}

	tflog.Info(ctx, "datacenter imported", map[string]any{"resource_id": d.Id()})

	if err := dcschema.PopulateResourceData(d, &datacenter); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
