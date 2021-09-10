package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
)

func dataSourcedbaasCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcedbaasReadCluster,
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
			"backup_enabled": {
				Type:       schema.TypeBool,
				Deprecated: "Deprecated: backup is always enabled.\n      Enables automatic backups of your cluster.",
				Computed:   true,
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

func dataSourcedbaasReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		cluster, _, err = client.ClustersApi.ClustersFindById(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the dbaas cluster with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		clusters, _, err := client.ClustersApi.ClustersGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %s", err.Error()))
			return diags
		}

		found := false
		if clusters.Data != nil {
			for _, c := range *clusters.Data {
				tmpCluster, _, err := client.ClustersApi.ClustersFindById(ctx, *c.Id).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas cluster with ID %s: %s", *c.Id, err.Error()))
					return diags
				}
				if tmpCluster.DisplayName != nil && *tmpCluster.DisplayName == name.(string) {
					/* lan found */
					cluster = tmpCluster
					found = true
					break
				}

			}
		}
		if !found {
			return diag.FromErr(errors.New("dbaas cluster not found"))
		}

	}

	if diags := setDbaasClusterData(d, &cluster); diags != nil {
		return diags
	}

	return nil

}

func setDbaasClusterData(d *schema.ResourceData, cluster *dbaas.Cluster) diag.Diagnostics {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
		if err := d.Set("id", *cluster.Id); err != nil {
			return diag.FromErr(err)
		}
	}

	if cluster.PostgresVersion != nil {
		if err := d.Set("postgres_version", *cluster.PostgresVersion); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting postgres_version property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.Replicas != nil {
		if err := d.Set("replicas", *cluster.Replicas); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting replicas property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.CpuCoreCount != nil {
		if err := d.Set("cpu_core_count", *cluster.CpuCoreCount); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting cpu_core_count for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.RamSize != nil {
		if err := d.Set("ram_size", *cluster.RamSize); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting ram_size property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.StorageSize != nil {
		if err := d.Set("storage_size", *cluster.StorageSize); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting storage_size property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.StorageType != nil {
		if err := d.Set("storage_type", *cluster.StorageType); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting storage_type property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.VdcConnections != nil && len(*cluster.VdcConnections) > 0 {
		var connections []interface{}
		for _, connection := range *cluster.VdcConnections {
			connectionEntry := make(map[string]interface{})

			if connection.VdcId != nil {
				connectionEntry["vdc_id"] = *connection.VdcId
			}

			if connection.LanId != nil {
				connectionEntry["lan_id"] = *connection.LanId
			}

			if connection.IpAddress != nil {
				connectionEntry["ip_address"] = *connection.IpAddress
			}

			connections = append(connections, connectionEntry)
		}
		if len(connections) > 0 {
			if err := d.Set("vdc_connections", connections); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting vdc_connections property for dbaas cluster  %s: %s", d.Id(), err))
				return diags
			}
		}
	}

	if cluster.Location != nil {
		if err := d.Set("location", *cluster.Location); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting location property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.DisplayName != nil {
		if err := d.Set("display_name", *cluster.DisplayName); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting display_name property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.BackupEnabled != nil {
		if err := d.Set("backup_enabled", *cluster.BackupEnabled); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting backup_enabled property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}

		maintenanceWindowEntry := make(map[string]interface{})

		if cluster.MaintenanceWindow.Time != nil {
			maintenanceWindowEntry["time"] = *cluster.MaintenanceWindow.Time
		}

		if cluster.MaintenanceWindow.Weekday != nil {
			maintenanceWindowEntry["weekday"] = *cluster.MaintenanceWindow.Weekday
		}

		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		err := d.Set("maintenance_window", maintenanceWindow)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting maintenance_window property for dbaas cluster %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}
