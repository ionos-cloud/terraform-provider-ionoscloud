package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbaasMongoUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasMongoReadUser,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Id of the backup unit.",
				Optional:    true,
			},
			"cluster_id": {
				Type:             schema.TypeString,
				Description:      "The id of your cluster.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"username": {
				Type:             schema.TypeString,
				Description:      "The username to search for",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"database": {
				Type:             schema.TypeString,
				Description:      "The database",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A list of mongodb user roles. Examples: read, readWrite, readAnyDatabase",
							Computed:    true,
						},
						"database": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasMongoReadUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	clusterIdIf, idOk := d.GetOk("cluster_id")
	usernameIf, nameOk := d.GetOk("username")

	if !idOk || !nameOk {
		diags := diag.FromErr(errors.New("please provide cluster_id and username"))
		return diags
	}

	username := usernameIf.(string)
	clusterId := clusterIdIf.(string)
	var user mongo.User
	var err error

	users, _, err := client.GetUsers(ctx, clusterId)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas mongo users: %w", err))
		return diags
	}

	var results []mongo.User

	if users.Items != nil && len(*users.Items) > 0 {
		for _, userItem := range *users.Items {
			if userItem.Properties != nil && userItem.Properties.Username != nil && strings.EqualFold(*userItem.Properties.Username, username) {
				results = append(results, userItem)
			}
		}
	}

	if len(results) == 0 {
		return diag.FromErr(fmt.Errorf("no DBaaS mongo user found with the specified username = %s and cluster_id = %s", username, clusterId))
	} else if len(results) > 1 {
		return diag.FromErr(fmt.Errorf("more than one DBaaS mongo user found with the specified criteria username = %s and cluster_id = %s", username, clusterId))
	} else {
		user = results[0]
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return diag.FromErr(err)
	}
	if user.Properties != nil && user.Properties.Username != nil {
		d.SetId(clusterId + *user.Properties.Username)
	}

	return nil

}
