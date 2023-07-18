package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
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
	client := meta.(SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return diag.FromErr(fmt.Errorf("no PgSql user found with the specified username: %s and cluster ID: %s", username, clusterId))
		}
		return diag.FromErr(fmt.Errorf("an error occured while fetching the PgSql user: %s, cluster ID: %s, err: %w", username, clusterId, err))
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
