package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"
	"time"
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
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.GroupShare{
		Properties: &ionoscloud.GroupShareProperties{},
	}

	tempSharePrivilege := d.Get("edit_privilege").(bool)
	request.Properties.SharePrivilege = &tempSharePrivilege
	tempEditPrivilege := d.Get("share_privilege").(bool)
	request.Properties.EditPrivilege = &tempEditPrivilege

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesPost(ctx,
		d.Get("group_id").(string), d.Get("resource_id").(string)).Resource(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a share: %s", err))
		return diags
	}

	if rsp.Id != nil {
		log.Printf("[DEBUG] SHARE ID: %s", *rsp.Id)
		d.SetId(*rsp.Id)
	}

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
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx,
		d.Get("group_id").(string), d.Get("resource_id").(string)).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Share ID %s %s", d.Id(), err))
		return diags
	}

	if rsp.Properties.EditPrivilege != nil {
		if err := d.Set("edit_privilege", *rsp.Properties.EditPrivilege); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.SharePrivilege != nil {
		if err := d.Set("share_privilege", *rsp.Properties.SharePrivilege); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	return nil
}

func resourceShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a share ID %s %s", d.Id(), err))
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
	client := meta.(*ionoscloud.APIClient)
	groupId := d.Get("group_id").(string)
	resourceId := d.Get("resource_id").(string)

	_, apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting share object %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for share %s to be deleted...", d.Id())

		sDeleted, dsErr := shareDeleted(ctx, client, d)
		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of share %s: %s", d.Id(), dsErr))
			return diags
		}

		if sDeleted {
			log.Printf("[INFO] Successfully deleted Share: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] share deletion timed out")
			diags := diag.FromErr(fmt.Errorf("share deletion timed out! WARNING: your share will still probably be deleted after some" +
				" time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}

	}

	d.SetId("")
	return nil
}
func shareDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	groupId := d.Get("group_id").(string)
	resourceId := d.Get("resource_id").(string)
	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, groupId, resourceId).Execute()
	logApiRequestTime(apiResponse)

	log.Printf("[INFO] Current deletion status for share %s: %+v", d.Id(), rsp)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, err

	}
	log.Printf("[INFO] share %s not deleted yet deleted : %+v", d.Id(), rsp)
	return false, nil
}

func resourceShareImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {group}/{resource}", d.Id())
	}

	grpId := parts[0]
	rscId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, rscId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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
