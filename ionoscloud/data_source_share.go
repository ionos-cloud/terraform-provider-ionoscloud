package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceShareRead,
		Schema: map[string]*schema.Schema{
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
			"edit_privilege": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	groupID := d.Get("group_id").(string)
	resourceID := d.Get("resource_id").(string)
	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, groupID, resourceID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			return utils.ToDiags(d, fmt.Sprintf("group_id %s resource_id %s not found", groupID, resourceID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching a share with group_id %s resource_id %s %s", groupID, resourceID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if rsp.Properties == nil {
		return utils.ToDiags(d, "no properties found in the response", nil)
	}
	d.SetId(*rsp.Id)
	if err := d.Set("edit_privilege", *rsp.Properties.EditPrivilege); err != nil {
		return utils.ToDiags(d, utils.GenerateSetError("share", "edit_privilege", err).Error(), nil)
	}
	if err := d.Set("share_privilege", *rsp.Properties.SharePrivilege); err != nil {
		return utils.ToDiags(d, utils.GenerateSetError("share", "share_privilege", err).Error(), nil)
	}
	return nil
}
