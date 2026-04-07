package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourcePrivateCrossConnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCrossConnectCreate,
		ReadContext:   resourcePrivateCrossConnectRead,
		UpdateContext: resourcePrivateCrossConnectUpdate,
		DeleteContext: resourcePrivateCrossConnectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivateCrossConnectImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The desired name",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The desired description",
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				ForceNew:    true,
			},
			"connectable_datacenters": {
				Type:        schema.TypeList,
				Description: "A list containing all the connectable datacenters",
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The UUID of the connectable datacenter",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the connectable datacenter",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The physical location of the connectable datacenter",
							Computed:    true,
						},
					},
				},
			},
			"peers": {
				Type:        schema.TypeList,
				Description: "A list containing the details of all cross-connected datacenters",
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The id of the cross-connected LAN",
							Computed:    true,
						},
						"lan_name": {
							Type:        schema.TypeString,
							Description: "The name of the cross-connected LAN",
							Computed:    true,
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The id of the cross-connected datacenter",
							Computed:    true,
						},
						"datacenter_name": {
							Type:        schema.TypeString,
							Description: "The name of the cross-connected datacenter",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location of the cross-connected datacenter",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcePrivateCrossConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	pcc := ionoscloud.PrivateCrossConnect{
		Properties: &ionoscloud.PrivateCrossConnectProperties{
			Name: &name,
		},
	}

	if descVal, descOk := d.GetOk("description"); descOk {
		log.Printf("[INFO] Setting PCC description to : %s", descVal.(string))
		description := descVal.(string)
		pcc.Properties.Description = &description
	}

	rsp, apiResponse, err := client.PrivateCrossConnectsApi.PccsPost(ctx).Pcc(pcc).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		return diagutil.ToDiags(d, fmt.Errorf("error creating cross connect: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*rsp.Id)
	log.Printf("[INFO] Created PCC: %s", d.Id())

	if diags := waitForPCCToBeReady(ctx, d, client); diags != nil {
		return diags
	}
	return resourcePrivateCrossConnectRead(ctx, d, meta)
}

func resourcePrivateCrossConnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	pcc, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching PCC: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] Successfully retrieved PCC %s: %+v", d.Id(), pcc)

	if err = setPccDataSource(d, &pcc); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func resourcePrivateCrossConnectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	request := ionoscloud.PrivateCrossConnect{}
	name := d.Get("name").(string)
	request.Properties = &ionoscloud.PrivateCrossConnectProperties{
		Name: &name,
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		log.Printf("[INFO] PCC name changed from %+v to %+v", oldName, newName)
		name := newName.(string)
		request.Properties.Name = &name
	}

	log.Printf("[INFO] Attempting update PCC %s", d.Id())

	if d.HasChange("description") {
		oldDesc, newDesc := d.GetChange("description")
		log.Printf("[INFO] PCC description changed from %+v to %+v", oldDesc, newDesc)
		descriprion := newDesc.(string)
		if newDesc != nil {
			request.Properties.Description = &descriprion
		}
	}

	_, apiResponse, err := client.PrivateCrossConnectsApi.PccsPatch(ctx, d.Id()).Pcc(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while updating PCC: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if diags := waitForPCCToBeReady(ctx, d, client); diags != nil {
		return diags
	}

	return resourcePrivateCrossConnectRead(ctx, d, meta)
}

func resourcePrivateCrossConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.PrivateCrossConnectsApi.PccsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting PCC: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	for {
		log.Printf("[INFO] Waiting for PCC %s to be deleted...", d.Id())

		pccDeleted, dsErr := privateCrossConnectDeleted(ctx, client, d)

		if dsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while deleting PCC: %w", err), nil)
		}

		if pccDeleted {
			log.Printf("[INFO] Successfully deleted PCC: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] delete timed out")
			return diagutil.ToDiags(d, fmt.Errorf("pcc removal timed out! WARNING: your pcc will still probably be removed after some time but the terraform state wont reflect that; check the updates in your Ionos Cloud account"), nil)
		}
	}

	return nil
}

func resourcePrivateCrossConnectImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<pcc-id> or <pcc-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	pccId := parts[0]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	pcc, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, pccId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find PCC %q", pccId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("unable to retrieve PCC, error: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] PCC found: %+v", pcc)

	d.SetId(*pcc.Id)
	if err := d.Set("name", *pcc.Properties.Name); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("description", *pcc.Properties.Description); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if err = setPccDataSource(d, &pcc); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	log.Printf("[INFO] Importing PCC %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func privateCrossConnectReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	rsp, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking PCC status: %w", err)
	}
	return strings.EqualFold(*rsp.Metadata.State, constant.Available), nil
}

func privateCrossConnectDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking PCC deletion status: %w", err)
	}
	return false, nil
}

func waitForPCCToBeReady(ctx context.Context, d *schema.ResourceData, client *ionoscloud.APIClient) diag.Diagnostics {
	for {
		log.Printf("[INFO] Waiting for PCC %s to be ready...", d.Id())

		pccReady, rsErr := privateCrossConnectReady(ctx, client, d)

		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of PCC: %w", rsErr), nil)
		}

		if pccReady {
			log.Printf("[INFO] PCC ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] update timed out")
			return diagutil.ToDiags(d, fmt.Errorf("pcc readiness check timed out! WARNING: your pcc will still probably be created/updated after some time "+
				"but the terraform state wont reflect that; check your Ionos Cloud account to see the updates"), nil)
		}
	}
	return nil
}
