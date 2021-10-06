package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/dbaas"
)

func dataSourceDbaasPgSqlCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadCluster,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Optional:    true,
			},
			"postgres_version": {
				Type:        schema.TypeString,
				Description: "The PostgreSQL version of your cluster.",
				Computed:    true,
			},
			"replicas": {
				Type:        schema.TypeInt,
				Description: "The number of replicas in your cluster.",
				Computed:    true,
			},
			"cpu_core_count": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Computed:    true,
			},
			"ram_size": {
				Type:        schema.TypeString,
				Description: "The amount of memory per replica.",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeString,
				Description: "The amount of storage per replica.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "The storage type used in your cluster.",
				Computed:    true,
			},
			"vdc_connections": {
				Type:        schema.TypeList,
				Description: "The VDC to connect to your cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vdc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for the database.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24",
							Computed:    true,
						},
					},
				},
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				Computed:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "a weekly 4 hour-long window, during which maintenance might occur",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weekday": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				Description: "Credentials for the database user to be created.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Description: "the username for the initial postgres user. some system usernames\n          are restricted (e.g. \"postgres\", \"admin\", \"standby\")",
							Computed:    true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("display_name")

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and display_name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the dbaas cluster id or display_name"))
		return diags
	}

	var cluster dbaas.Cluster
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id.(string))
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the dbaas cluster with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		clusters, _, err := client.ListClusters(ctx, name.(string))

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %s", err.Error()))
			return diags
		}

		if clusters.Data != nil && len(*clusters.Data) > 0 {
			cluster = (*clusters.Data)[0]
		} else {
			return diag.FromErr(errors.New("dbaas cluster not found"))
		}
	}

	if cluster.Id != nil {
		if err := d.Set("id", *cluster.Id); err != nil {
			return nil
		}
	}

	if err := dbaasService.SetDbaasPgSqlClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
