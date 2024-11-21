package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func dataSourceShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceShareRead,
		Schema: map[string]*schema.Schema{
			"edit_privilege": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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

	client := meta.(services.SdkBundle).CloudApiClient

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, d.Get("group_id").(string), d.Get("resource_id").(string)).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			return diag.FromErr(fmt.Errorf("share ID %s group_id %s resource_id %s not found", d.Id(), d.Get("group_id").(string), d.Get("resource_id").(string)))
		}
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching a Share ID %s %w", d.Id(), err))
		return diags
	}
	d.SetId(*rsp.Id)
	if err := d.Set("edit_privilege", *rsp.Properties.EditPrivilege); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("share_privilege", *rsp.Properties.SharePrivilege); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	return nil

	return nil
}
