package dbaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ClusterService interface {
	GetCluster(ctx context.Context, clusterId string) (dbaas.Cluster, *dbaas.APIResponse, error)
	ListClusters(ctx context.Context) (dbaas.ClusterList, *dbaas.APIResponse, error)
	CreateCluster(ctx context.Context, cluster dbaas.CreateClusterRequest) (dbaas.Cluster, *dbaas.APIResponse, error)
	CreateClusterFromBackup(ctx context.Context, cluster dbaas.CreateClusterRequest, backupId string) (dbaas.Cluster, *dbaas.APIResponse, error)
	UpdateCluster(ctx context.Context, clusterId string, cluster dbaas.PatchClusterRequest) (dbaas.Cluster, *dbaas.APIResponse, error)
	DeleteCluster(ctx context.Context, clusterId string) (dbaas.Cluster, *dbaas.APIResponse, error)
}

func (c *Client) GetCluster(ctx context.Context, clusterId string) (dbaas.Cluster, *dbaas.APIResponse, error) {
	cluster, apiResponse, err := c.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *Client) ListClusters(ctx context.Context) (dbaas.ClusterList, *dbaas.APIResponse, error) {
	clusters, apiResponse, err := c.ClustersApi.ClustersGet(ctx).Execute()
	if apiResponse != nil {
		return clusters, apiResponse, err
	}
	return clusters, nil, err
}

func (c *Client) CreateCluster(ctx context.Context, cluster dbaas.CreateClusterRequest) (dbaas.Cluster, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersPost(ctx).Cluster(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) CreateClusterFromBackup(ctx context.Context, cluster dbaas.CreateClusterRequest, backupId string) (dbaas.Cluster, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersPost(ctx).Cluster(cluster).FromBackup(backupId).Execute()
	if err != nil {
		fmt.Printf("error while creating from backup: %v", err)
	}
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) UpdateCluster(ctx context.Context, clusterId string, cluster dbaas.PatchClusterRequest) (dbaas.Cluster, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersPatch(ctx, clusterId).Cluster(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) DeleteCluster(ctx context.Context, clusterId string) (dbaas.Cluster, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersDelete(ctx, clusterId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func GetDbaasPgSqlClusterDataCreate(d *schema.ResourceData) *dbaas.CreateClusterRequest {

	dbaasCluster := dbaas.CreateClusterRequest{}

	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
		postgresVersion := postgresVersion.(string)
		dbaasCluster.PostgresVersion = &postgresVersion
	}

	if replicas, ok := d.GetOk("replicas"); ok {
		replicas := float32(replicas.(int))
		dbaasCluster.Replicas = &replicas
	}

	if cpuCoreCount, ok := d.GetOk("cpu_core_count"); ok {
		cpuCoreCount := float32(cpuCoreCount.(int))
		dbaasCluster.CpuCoreCount = &cpuCoreCount
	}

	if ramSize, ok := d.GetOk("ram_size"); ok {
		ramSize := ramSize.(string)
		dbaasCluster.RamSize = &ramSize
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		storageSize := storageSize.(string)
		dbaasCluster.StorageSize = &storageSize
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageType := dbaas.StorageType(storageType.(string))
		dbaasCluster.StorageType = &storageType
	}

	if vdcConnection, ok := d.GetOk("vdc_connections"); ok {
		if vdcConnection.([]interface{}) != nil {

			var vdcConnections []dbaas.VDCConnection

			for vdcIndex := range vdcConnection.([]interface{}) {

				connection := dbaas.VDCConnection{}

				if vdcId, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.vdc_id", vdcIndex)); ok {
					vdcId := vdcId.(string)
					connection.VdcId = &vdcId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = &lanId
				}

				if ipAddress, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.ip_address", vdcIndex)); ok {
					ipAddress := ipAddress.(string)
					connection.IpAddress = &ipAddress
				}

				vdcConnections = append(vdcConnections, connection)

			}

			if len(vdcConnections) > 0 {
				dbaasCluster.VdcConnections = &vdcConnections
			}
		}
	}

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		dbaasCluster.Location = &location
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.DisplayName = &displayName
	}

	backupEnabled := d.Get("backup_enabled").(bool)
	dbaasCluster.BackupEnabled = &backupEnabled

	if _, ok := d.GetOk("maintenance_window.0"); ok {
		dbaasCluster.MaintenanceWindow = &dbaas.MaintenanceWindow{}

		if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
			timeV := timeV.(string)
			dbaasCluster.MaintenanceWindow.Time = &timeV
		}

		if weekday, ok := d.GetOk("maintenance_window.0.weekday"); ok {
			weekday := weekday.(string)
			dbaasCluster.MaintenanceWindow.Weekday = &weekday
		}
	}

	if _, ok := d.GetOk("credentials.0"); ok {
		dbaasCluster.Credentials = &dbaas.DBUser{}
		if username, ok := d.GetOk("credentials.0.username"); ok {
			username := username.(string)
			dbaasCluster.Credentials.Username = &username
		}

		if password, ok := d.GetOk("credentials.0.password"); ok {
			password := password.(string)
			dbaasCluster.Credentials.Password = &password
		}
	}

	return &dbaasCluster
}

func GetDbaasPgSqlClusterDataUpdate(d *schema.ResourceData) (*dbaas.PatchClusterRequest, diag.Diagnostics) {

	cluster := dbaas.PatchClusterRequest{}

	if d.HasChange("postgres_version") {
		_, v := d.GetChange("postgres_version")
		vStr := v.(string)
		cluster.PostgresVersion = &vStr
	}

	if d.HasChange("replicas") {
		diags := diag.FromErr(fmt.Errorf("replicas parameter is immutable"))
		return nil, diags
	}

	if d.HasChange("cpu_core_count") {
		diags := diag.FromErr(fmt.Errorf("cpu_core_count parameter is immutable"))
		return nil, diags
	}

	if d.HasChange("ram_size") {
		diags := diag.FromErr(fmt.Errorf("ram_size parameter is immutable"))
		return nil, diags
	}

	if d.HasChange("storage_size") {
		_, v := d.GetChange("storage_size")
		vStr := v.(string)
		cluster.StorageSize = &vStr
	}

	if d.HasChange("vdc_connections") {
		diags := diag.FromErr(fmt.Errorf("vdc_connections parameter is immutable"))
		return nil, diags
	}

	if d.HasChange("display_name") {
		_, v := d.GetChange("display_name")
		vStr := v.(string)
		cluster.DisplayName = &vStr
	}

	if d.HasChange("backup_enabled") {
		vBool := d.Get("backup_enabled").(bool)
		cluster.BackupEnabled = &vBool
	}

	if d.HasChange("maintenance_window") {
		cluster.MaintenanceWindow = &dbaas.MaintenanceWindow{}

		if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
			timeV := timeV.(string)
			cluster.MaintenanceWindow.Time = &timeV
		}

		if weekday, ok := d.GetOk("maintenance_window.0.weekday"); ok {
			weekday := weekday.(string)
			cluster.MaintenanceWindow.Weekday = &weekday
		}
	}

	return &cluster, nil
}

func SetDbaasPgSqlClusterData(d *schema.ResourceData, cluster dbaas.Cluster) diag.Diagnostics {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
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
