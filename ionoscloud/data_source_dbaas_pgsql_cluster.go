package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"log"
	"strings"
)

func dataSourceDbaasPgSqlCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadCluster,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "The id of your cluster.",
				Optional:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
				Optional:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"postgres_version": {
				Type:        schema.TypeString,
				Description: "The PostgreSQL version of your cluster.",
				Optional:    true,
			},
			"instances": {
				Type:        schema.TypeInt,
				Description: "The total number of instances in the cluster (one master and n-1 standbys)",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Computed:    true,
			},
			"ram": {
				Type:        schema.TypeInt,
				Description: "The amount of memory per instance in megabytes. Has to be a multiple of 1024.",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The amount of storage per instance in megabytes.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "The storage type used in your cluster.",
				Computed:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "Details about the network connection for your cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your cluster to.",
							Computed:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The LAN to connect your cluster to.",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24",
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
			"backup_location": {
				Type:        schema.TypeString,
				Description: "The S3 location where the backups will be stored.",
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
						"day_of_the_week": {
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
							Description: "the username for the initial postgres user. some system usernames are restricted (e.g. \"postgres\", \"admin\", \"standby\")",
							Computed:    true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
			"synchronization_mode": {
				Type:        schema.TypeString,
				Description: "Represents different modes of replication.",
				Computed:    true,
			},
			"from_backup": {
				Type:        schema.TypeList,
				Description: "The PostgreSQL version of your cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the backup you want to restore.",
							Computed:    true,
						},
						"recovery_target_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.",
						},
					},
				},
			},
			"datacenter_name": {
				Type:        schema.TypeString,
				Description: "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudClient := meta.(SdkBundle).CloudApiClient
	client := meta.(SdkBundle).PsqlClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("display_name")
	locationValue, locationOk := d.GetOk("location")
	dcNameValue, dcNameOk := d.GetOk("datacenter_name")
	pgVersionValue, pgVersionOk := d.GetOk("postgres_version")

	id := idValue.(string)
	name := nameValue.(string)
	location := locationValue.(string)
	dcName := dcNameValue.(string)
	pgVersion := pgVersionValue.(string)

	if idOk && (nameOk || locationOk || dcNameOk || pgVersionOk) {
		diags := diag.FromErr(errors.New("id and display_name/location/datacenter_name/postgres_version cannot be both specified in the same time, choose between id or a combination of other parameters"))
		return diags
	}
	if !idOk && !nameOk && !locationOk && !dcNameOk && !pgVersionOk {
		diags := diag.FromErr(errors.New("please provide either the dbaas cluster id or other parameter like display_name, location or postgres_version"))
		return diags
	}

	var cluster dbaas.ClusterResponse
	var err error

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for DBaaS Postgres Cluster by id %s", id)

		cluster, _, err = client.GetCluster(ctx, id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the dbaas cluster with ID %s: %w", id, err))
			return diags
		}
	} else {
		var results []dbaas.ClusterResponse

		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for DBaaS Postgres Cluster by name with partial_match %t and name: %s", partialMatch, name)
			if partialMatch {
				clusters, _, err := client.ListClusters(ctx, name)
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %s", err.Error()))
					return diags
				}
				if len(*clusters.Items) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria name with partial match: %s", name))
				}
				results = *clusters.Items
			} else {
				clusters, _, err := client.ListClusters(ctx, "")
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %s", err.Error()))
					return diags
				}
				if clusters.Items != nil && len(*clusters.Items) > 0 && nameOk {
					var nameResults []dbaas.ClusterResponse
					for _, clusterItem := range *clusters.Items {
						if clusterItem.Properties != nil && clusterItem.Properties.DisplayName != nil && strings.EqualFold(*clusterItem.Properties.DisplayName, name) {
							nameResults = append(nameResults, clusterItem)
						}
					}
					if len(nameResults) == 0 {
						return diag.FromErr(fmt.Errorf("no result found with the specified criteria name: %s", name))
					}
					results = nameResults
				}
			}
		} else {
			clusters, _, err := client.ListClusters(ctx, "")
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %s", err.Error()))
				return diags
			}
			results = *clusters.Items
		}

		if locationOk && location != "" {
			var locationResults []dbaas.ClusterResponse
			for _, cluster := range results {
				if cluster.Properties != nil && cluster.Properties.Location != nil && strings.EqualFold(string(*cluster.Properties.Location), location) {
					locationResults = append(locationResults, cluster)
				}
			}
			if len(locationResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria location: %s", location))
			}
			results = locationResults
		}

		if dcNameOk && dcName != "" {
			var dcNameResults []dbaas.ClusterResponse
			for _, cluster := range results {
				if cluster.Properties != nil && cluster.Properties.Connections != nil && *(*cluster.Properties.Connections)[0].DatacenterId != "" {
					searchedDcId := *(*cluster.Properties.Connections)[0].DatacenterId
					datacenter, apiResponse, err := cloudClient.DataCentersApi.DatacentersFindById(ctx, searchedDcId).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching the datacenter while searching by id %s %w", id, err))
					}
					if datacenter.Properties != nil {
						actualDcName := datacenter.Properties.Name
						if strings.EqualFold(*actualDcName, dcName) {
							dcNameResults = append(dcNameResults, cluster)
						}
					}
				}
			}
			if len(dcNameResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria datacenter name: %s", dcName))
			}
			results = dcNameResults
		}

		if pgVersionOk && pgVersion != "" {
			var pgVersionResults []dbaas.ClusterResponse
			for _, cluster := range results {
				if cluster.Properties != nil && cluster.Properties.PostgresVersion != nil && strings.EqualFold(*cluster.Properties.PostgresVersion, location) {
					pgVersionResults = append(pgVersionResults, cluster)
				}
			}
			if len(pgVersionResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: postgres version: %s", pgVersion))
			}
			results = pgVersionResults
		}

		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DBaaS cluster found with the specified criteria"))
		} else {
			cluster = results[0]
		}

	}

	if err := dbaasService.SetDbaasPgSqlClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
