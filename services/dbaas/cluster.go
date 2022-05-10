package dbaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"time"
)

type ClusterService interface {
	GetCluster(ctx context.Context, clusterId string) (dbaas.ClusterResponse, *dbaas.APIResponse, error)
	ListClusters(ctx context.Context, filterName string) (dbaas.ClusterList, *dbaas.APIResponse, error)
	CreateCluster(ctx context.Context, cluster dbaas.CreateClusterRequest) (dbaas.ClusterResponse, *dbaas.APIResponse, error)
	UpdateCluster(ctx context.Context, clusterId string, cluster dbaas.PatchClusterRequest) (dbaas.ClusterResponse, *dbaas.APIResponse, error)
	DeleteCluster(ctx context.Context, clusterId string) (dbaas.ClusterResponse, *dbaas.APIResponse, error)
}

func (c *Client) GetCluster(ctx context.Context, clusterId string) (dbaas.ClusterResponse, *dbaas.APIResponse, error) {
	cluster, apiResponse, err := c.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *Client) ListClusters(ctx context.Context, filterName string) (dbaas.ClusterList, *dbaas.APIResponse, error) {
	request := c.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.ClustersApi.ClustersGetExecute(request)
	if apiResponse != nil {
		return clusters, apiResponse, err
	}
	return clusters, nil, err
}

func (c *Client) CreateCluster(ctx context.Context, cluster dbaas.CreateClusterRequest) (dbaas.ClusterResponse, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) UpdateCluster(ctx context.Context, clusterId string, cluster dbaas.PatchClusterRequest) (dbaas.ClusterResponse, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersPatch(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) DeleteCluster(ctx context.Context, clusterId string) (dbaas.ClusterResponse, *dbaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.ClustersApi.ClustersDelete(ctx, clusterId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func GetDbaasPgSqlClusterDataCreate(d *schema.ResourceData) (*dbaas.CreateClusterRequest, error) {

	dbaasCluster := dbaas.CreateClusterRequest{
		Properties: &dbaas.CreateClusterProperties{},
	}

	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
		postgresVersion := postgresVersion.(string)
		dbaasCluster.Properties.PostgresVersion = &postgresVersion
	}

	if instances, ok := d.GetOk("instances"); ok {
		instances := int32(instances.(int))
		dbaasCluster.Properties.Instances = &instances
	}

	if cores, ok := d.GetOk("cores"); ok {
		cores := int32(cores.(int))
		dbaasCluster.Properties.Cores = &cores
	}

	if ram, ok := d.GetOk("ram"); ok {
		ram := int32(ram.(int))
		dbaasCluster.Properties.Ram = &ram
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		storageSize := int32(storageSize.(int))
		dbaasCluster.Properties.StorageSize = &storageSize
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageType := dbaas.StorageType(storageType.(string))
		dbaasCluster.Properties.StorageType = &storageType
	}

	if _, ok := d.GetOk("connections"); ok {
		dbaasCluster.Properties.Connections = GetDbaasClusterConnectionsData(d)
	} else {
		return nil, fmt.Errorf("connections parameter is required in create cluster requests")
	}

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		dbaasCluster.Properties.Location = &location
	}

	if backupLocation, ok := d.GetOk("backup_location"); ok {
		backupLocation := backupLocation.(string)
		dbaasCluster.Properties.BackupLocation = &backupLocation
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.Properties.DisplayName = &displayName
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dbaasCluster.Properties.MaintenanceWindow = GetDbaasClusterMaintenanceWindowData(d)
	}

	dbaasCluster.Properties.Credentials = GetDbaasClusterCredentialsData(d)

	if synchronizationMode, ok := d.GetOk("synchronization_mode"); ok {
		synchronizationMode := dbaas.SynchronizationMode(synchronizationMode.(string))
		dbaasCluster.Properties.SynchronizationMode = &synchronizationMode
	}

	if _, ok := d.GetOk("from_backup"); ok {
		if fromBackup, err := GetDbaasClusterFromBackupData(d); err != nil {
			return nil, err
		} else {
			dbaasCluster.Properties.FromBackup = fromBackup
		}
	}

	return &dbaasCluster, nil
}

func GetDbaasPgSqlClusterDataUpdate(d *schema.ResourceData) (*dbaas.PatchClusterRequest, diag.Diagnostics) {

	dbaasCluster := dbaas.PatchClusterRequest{
		Properties: &dbaas.PatchClusterProperties{},
	}

	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
		postgresVersion := postgresVersion.(string)
		dbaasCluster.Properties.PostgresVersion = &postgresVersion
	}

	if instances, ok := d.GetOk("instances"); ok {
		instances := int32(instances.(int))
		dbaasCluster.Properties.Instances = &instances
	}

	if cores, ok := d.GetOk("cores"); ok {
		cores := int32(cores.(int))
		dbaasCluster.Properties.Cores = &cores
	}

	if ram, ok := d.GetOk("ram"); ok {
		ram := int32(ram.(int))
		dbaasCluster.Properties.Ram = &ram
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		storageSize := int32(storageSize.(int))
		dbaasCluster.Properties.StorageSize = &storageSize
	}

	dbaasCluster.Properties.Connections = GetDbaasClusterConnectionsData(d)

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.Properties.DisplayName = &displayName
	}

	dbaasCluster.Properties.MaintenanceWindow = GetDbaasClusterMaintenanceWindowData(d)

	return &dbaasCluster, nil
}

func GetDbaasClusterConnectionsData(d *schema.ResourceData) *[]dbaas.Connection {
	connections := make([]dbaas.Connection, 0)

	if vdcValue, ok := d.GetOk("connections"); ok {
		vdcValue := vdcValue.([]interface{})
		if vdcValue != nil {
			for vdcIndex := range vdcValue {

				connection := dbaas.Connection{}

				if datacenterId, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", vdcIndex)); ok {
					datacenterId := datacenterId.(string)
					connection.DatacenterId = &datacenterId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = &lanId
				}

				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.cidr", vdcIndex)); ok {
					cidr := cidr.(string)
					connection.Cidr = &cidr
				}

				connections = append(connections, connection)
			}
		}

	}

	return &connections
}

func GetDbaasClusterMaintenanceWindowData(d *schema.ResourceData) *dbaas.MaintenanceWindow {
	var maintenanceWindow dbaas.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = &timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := dbaas.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}

	return &maintenanceWindow
}

func GetDbaasClusterCredentialsData(d *schema.ResourceData) *dbaas.DBUser {
	var user dbaas.DBUser

	if username, ok := d.GetOk("credentials.0.username"); ok {
		username := username.(string)
		user.Username = &username
	}

	if password, ok := d.GetOk("credentials.0.password"); ok {
		password := password.(string)
		user.Password = &password
	}

	return &user
}

func GetDbaasClusterFromBackupData(d *schema.ResourceData) (*dbaas.CreateRestoreRequest, error) {
	var restore dbaas.CreateRestoreRequest

	if backupId, ok := d.GetOk("from_backup.0.backup_id"); ok {
		backupId := backupId.(string)
		restore.BackupId = &backupId
	}

	if targetTime, ok := d.GetOk("from_backup.0.recovery_target_time"); ok {
		var ionosTime dbaas.IonosTime
		targetTime := targetTime.(string)
		layout := "2006-01-02T15:04:05Z"
		convertedTime, err := time.Parse(layout, targetTime)
		if err != nil {
			return nil, fmt.Errorf("an error occured while converting recovery_target_time to time.Time: %s", err)

		}
		ionosTime.Time = convertedTime
		restore.RecoveryTargetTime = &ionosTime
	}

	return &restore, nil
}

func SetDbaasPgSqlClusterData(d *schema.ResourceData, cluster dbaas.ClusterResponse) error {

	resourceName := "dbaas cluster"

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties.PostgresVersion != nil {
		if err := d.Set("postgres_version", *cluster.Properties.PostgresVersion); err != nil {
			return utils.GenerateSetError(resourceName, "postgres_version", err)
		}
	}

	if cluster.Properties.Instances != nil {
		if err := d.Set("instances", *cluster.Properties.Instances); err != nil {
			return utils.GenerateSetError(resourceName, "instances", err)
		}
	}

	if cluster.Properties.Cores != nil {
		if err := d.Set("cores", *cluster.Properties.Cores); err != nil {
			return utils.GenerateSetError(resourceName, "cores", err)
		}
	}

	if cluster.Properties.Ram != nil {
		if err := d.Set("ram", *cluster.Properties.Ram); err != nil {
			return utils.GenerateSetError(resourceName, "ram", err)
		}
	}

	if cluster.Properties.StorageSize != nil {
		if err := d.Set("storage_size", *cluster.Properties.StorageSize); err != nil {
			return utils.GenerateSetError(resourceName, "storage_size", err)
		}
	}

	if cluster.Properties.StorageType != nil {
		if err := d.Set("storage_type", *cluster.Properties.StorageType); err != nil {
			return utils.GenerateSetError(resourceName, "storage_type", err)
		}
	}

	if cluster.Properties.Connections != nil && len(*cluster.Properties.Connections) > 0 {
		var connections []interface{}
		for _, connection := range *cluster.Properties.Connections {
			connectionEntry := SetConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}
		if err := d.Set("connections", connections); err != nil {
			return utils.GenerateSetError(resourceName, "connections", err)
		}
	}

	if cluster.Properties.Location != nil {
		if err := d.Set("location", *cluster.Properties.Location); err != nil {
			return fmt.Errorf("error while setting location property for dbaas cluster %s: %s", d.Id(), err)
		}
	}

	if cluster.Properties.BackupLocation != nil {
		if err := d.Set("backup_location", *cluster.Properties.BackupLocation); err != nil {
			return fmt.Errorf("error while setting backup_location property for dbaas cluster %s: %s", d.Id(), err)
		}
	}

	if cluster.Properties.DisplayName != nil {
		if err := d.Set("display_name", *cluster.Properties.DisplayName); err != nil {
			return fmt.Errorf("error while setting display_name property for dbaas cluster %s: %s", d.Id(), err)
		}
	}

	if cluster.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := SetMaintenanceWindowProperties(*cluster.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(resourceName, "maintenance_window", err)
		}
	}

	if cluster.Properties.SynchronizationMode != nil {
		if err := d.Set("synchronization_mode", *cluster.Properties.SynchronizationMode); err != nil {
			return fmt.Errorf("error while setting SynchronizationMode property for dbaas cluster %s: %s", d.Id(), err)
		}
	}

	return nil
}

func SetConnectionProperties(vdcConnection dbaas.Connection) map[string]interface{} {

	connection := map[string]interface{}{}

	utils.SetPropWithNilCheck(connection, "datacenter_id", vdcConnection.DatacenterId)
	utils.SetPropWithNilCheck(connection, "lan_id", vdcConnection.LanId)
	utils.SetPropWithNilCheck(connection, "cidr", vdcConnection.Cidr)

	return connection
}

func SetMaintenanceWindowProperties(maintenanceWindow dbaas.MaintenanceWindow) map[string]interface{} {

	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}
