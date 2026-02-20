package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceDbaasPgSqlUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadUser,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_system_user": {
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

func dataSourceDbaasPgSqlReadUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return utils.ToDiags(d, fmt.Sprintf("no PgSql user found with the specified username: %s and cluster ID: %s", username, clusterId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the PgSql user: %s, cluster ID: %s, err: %s", username, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}
