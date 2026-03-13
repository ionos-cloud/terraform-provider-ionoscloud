package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceShareCreate,
		ReadContext:   resourceShareRead,
		UpdateContext: resourceShareUpdate,
		DeleteContext: resourceShareDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceShareImporter,
		},
		Schema: map[string]*schema.Schema{
			"edit_privilege": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return diag.FromErr(err)
	}

	request := ionoscloud.GroupShare{
		Properties: &ionoscloud.GroupShareProperties{},
	}

	tempSharePrivilege := d.Get("share_privilege").(bool)
	request.Properties.SharePrivilege = &tempSharePrivilege
	tempEditPrivilege := d.Get("edit_privilege").(bool)
	request.Properties.EditPrivilege = &tempEditPrivilege

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesPost(ctx, d.Get("group_id").(string), d.Get("resource_id").(string)).Resource(request).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a share: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.StatusCode})
	}
	d.SetId(*rsp.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: schema.TimeoutCreate, RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceShareRead(ctx, d, meta)
}

func resourceShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return diag.FromErr(err)
	}

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, d.Get("group_id").(string), d.Get("resource_id").(string)).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching a Share: %w", err), nil)
	}

	if err := d.Set("edit_privilege", *rsp.Properties.EditPrivilege); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	if err := d.Set("share_privilege", *rsp.Properties.SharePrivilege); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func resourceShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return diag.FromErr(err)
	}

	tempSharePrivilege := d.Get("share_privilege").(bool)
	tempEditPrivilege := d.Get("edit_privilege").(bool)

	shareReq := ionoscloud.GroupShare{
		Properties: &ionoscloud.GroupShareProperties{
			EditPrivilege:  &tempEditPrivilege,
			SharePrivilege: &tempSharePrivilege,
		},
	}

	_, apiResponse, err := client.UserManagementApi.UmGroupsSharesPut(ctx,
		d.Get("group_id").(string), d.Get("resource_id").(string)).Resource(shareReq).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while patching a share: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: schema.TimeoutUpdate, RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceShareRead(ctx, d, meta)
}

func resourceShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := d.Get("group_id").(string)
	resourceId := d.Get("resource_id").(string)

	apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if !httpNotFound(apiResponse) {
			return diagutil.ToDiags(d, err, nil)
		}
	}
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: schema.TimeoutDelete, RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceShareImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, diagutil.ToError(d, fmt.Errorf("invalid import. Expecting {group}/{resource}"), nil)
	}

	grpId := parts[0]
	rscId := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return nil, err
	}

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, rscId).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to fetch the share of resource %q for group %q", rscId, grpId), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("share does not exist of resource %q for group %q", rscId, grpId), nil)
	}

	log.Printf("[INFO] share found: %+v", share)

	d.SetId(*share.Id)

	if err := d.Set("group_id", grpId); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if err := d.Set("resource_id", rscId); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if share.Properties.EditPrivilege != nil {
		if err := d.Set("edit_privilege", *share.Properties.EditPrivilege); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	if share.Properties.SharePrivilege != nil {
		if err := d.Set("share_privilege", *share.Properties.SharePrivilege); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	return []*schema.ResourceData{d}, nil
}
