package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a share: %w", err))
		return diags
	}
	d.SetId(*rsp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceShareRead(ctx, d, meta)
}

func resourceShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, d.Get("group_id").(string), d.Get("resource_id").(string)).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Share ID %s %w", d.Id(), err))
		return diags
	}

	if err := d.Set("edit_privilege", *rsp.Properties.EditPrivilege); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("share_privilege", *rsp.Properties.SharePrivilege); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	return nil
}

func resourceShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a share ID %s %w", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceShareRead(ctx, d, meta)
}

func resourceShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	groupId := d.Get("group_id").(string)
	resourceId := d.Get("resource_id").(string)

	apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			diags := diag.FromErr(err)
			return diags
		}
		//try again in 20 seconds
		// todo: get rid of this retry
		time.Sleep(20 * time.Second)
		apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			if !httpNotFound(apiResponse) {
				diags := diag.FromErr(fmt.Errorf("an error occured while deleting a share %s %w", d.Id(), err))
				return diags
			}
		}
	}

	// Wait, catching any errors
	if apiResponse != nil && apiResponse.Header.Get("Location") != "" {
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}

	d.SetId("")
	return nil
}

func resourceShareImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {group}/{resource}", d.Id())
	}

	grpId := parts[0]
	rscId := parts[1]

	client := meta.(SdkBundle).CloudApiClient

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, rscId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the share of resource %q for group %q", rscId, grpId)
		}
		return nil, fmt.Errorf("share does not exist of resource %q for group %q", rscId, grpId)
	}

	log.Printf("[INFO] share found: %+v", share)

	d.SetId(*share.Id)

	if err := d.Set("group_id", grpId); err != nil {
		return nil, err
	}

	if err := d.Set("resource_id", rscId); err != nil {
		return nil, err
	}

	if share.Properties.EditPrivilege != nil {
		if err := d.Set("edit_privilege", *share.Properties.EditPrivilege); err != nil {
			return nil, err
		}
	}

	if share.Properties.SharePrivilege != nil {
		if err := d.Set("share_privilege", *share.Properties.SharePrivilege); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}
