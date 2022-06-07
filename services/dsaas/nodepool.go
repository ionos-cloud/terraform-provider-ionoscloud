package dsaas

import (
	"context"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
)

type NodePoolService interface {
	GetNodePool(ctx context.Context, clusterId, nodePoolId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	ListNodePools(ctx context.Context, clusterId string) ([]dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	CreateNodePool(ctx context.Context, cluster dsaas.CreateNodePoolRequest, clusterId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	UpdateNodePool(ctx context.Context, clusterId, nodePoolId string, cluster dsaas.PatchNodePoolRequest) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	DeleteNodePool(ctx context.Context, clusterId, nodePoolId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
}

func (c *Client) GetNodePool(ctx context.Context, clusterId, nodePoolId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
	cluster, apiResponse, err := c.DataPlatformNodePoolApi.GetClusterNodepool(ctx, clusterId, nodePoolId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *Client) ListNodePools(ctx context.Context, clusterId string) ([]dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
	nodePool, apiResponse, err := c.DataPlatformNodePoolApi.GetClusterNodepools(ctx, clusterId).Execute()
	if apiResponse != nil {
		return nodePool, apiResponse, err
	}
	return nodePool, nil, err
}

func (c *Client) CreateNodePool(ctx context.Context, cluster dsaas.CreateNodePoolRequest, clusterId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformNodePoolApi.CreateClusterNodepool(ctx, clusterId).CreateNodePoolRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) UpdateNodePool(ctx context.Context, clusterId, nodePoolId string, cluster dsaas.PatchNodePoolRequest) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformNodePoolApi.PatchClusterNodepool(ctx, clusterId, nodePoolId).PatchNodePoolRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) DeleteNodePool(ctx context.Context, clusterId, nodePoolId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformNodePoolApi.DeleteClusterNodepool(ctx, clusterId, nodePoolId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

//func GetdsaasPgSqlClusterDataCreate(d *schema.ResourceData) (*dsaas.CreateClusterRequest, error) {
//
//	dsaasCluster := dsaas.CreateClusterRequest{
//		Properties: &dsaas.CreateClusterProperties{},
//	}
//
//	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
//		postgresVersion := postgresVersion.(string)
//		dsaasCluster.Properties.PostgresVersion = &postgresVersion
//	}
//
//	if instances, ok := d.GetOk("instances"); ok {
//		instances := int32(instances.(int))
//		dsaasCluster.Properties.Instances = &instances
//	}
//
//	if cores, ok := d.GetOk("cores"); ok {
//		cores := int32(cores.(int))
//		dsaasCluster.Properties.Cores = &cores
//	}
//
//	if ram, ok := d.GetOk("ram"); ok {
//		ram := int32(ram.(int))
//		dsaasCluster.Properties.Ram = &ram
//	}
//
//	if storageSize, ok := d.GetOk("storage_size"); ok {
//		storageSize := int32(storageSize.(int))
//		dsaasCluster.Properties.StorageSize = &storageSize
//	}
//
//	if storageType, ok := d.GetOk("storage_type"); ok {
//		storageType := dsaas.StorageType(storageType.(string))
//		dsaasCluster.Properties.StorageType = &storageType
//	}
//
//	if _, ok := d.GetOk("connections"); ok {
//		dsaasCluster.Properties.Connections = GetdsaasClusterConnectionsData(d)
//	} else {
//		return nil, fmt.Errorf("connections parameter is required in create cluster requests")
//	}
//
//	if location, ok := d.GetOk("location"); ok {
//		location := location.(string)
//		dsaasCluster.Properties.Location = &location
//	}
//
//	if backupLocation, ok := d.GetOk("backup_location"); ok {
//		backupLocation := backupLocation.(string)
//		dsaasCluster.Properties.BackupLocation = &backupLocation
//	}
//
//	if displayName, ok := d.GetOk("display_name"); ok {
//		displayName := displayName.(string)
//		dsaasCluster.Properties.DisplayName = &displayName
//	}
//
//	if _, ok := d.GetOk("maintenance_window"); ok {
//		dsaasCluster.Properties.MaintenanceWindow = GetdsaasClusterMaintenanceWindowData(d)
//	}
//
//	dsaasCluster.Properties.Credentials = GetdsaasClusterCredentialsData(d)
//
//	if synchronizationMode, ok := d.GetOk("synchronization_mode"); ok {
//		synchronizationMode := dsaas.SynchronizationMode(synchronizationMode.(string))
//		dsaasCluster.Properties.SynchronizationMode = &synchronizationMode
//	}
//
//	if _, ok := d.GetOk("from_backup"); ok {
//		if fromBackup, err := GetdsaasClusterFromBackupData(d); err != nil {
//			return nil, err
//		} else {
//			dsaasCluster.Properties.FromBackup = fromBackup
//		}
//	}
//
//	return &dsaasCluster, nil
//}
//
//func GetdsaasPgSqlClusterDataUpdate(d *schema.ResourceData) (*dsaas.PatchClusterRequest, diag.Diagnostics) {
//
//	dsaasCluster := dsaas.PatchClusterRequest{
//		Properties: &dsaas.PatchClusterProperties{},
//	}
//
//	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
//		postgresVersion := postgresVersion.(string)
//		dsaasCluster.Properties.PostgresVersion = &postgresVersion
//	}
//
//	if instances, ok := d.GetOk("instances"); ok {
//		instances := int32(instances.(int))
//		dsaasCluster.Properties.Instances = &instances
//	}
//
//	if cores, ok := d.GetOk("cores"); ok {
//		cores := int32(cores.(int))
//		dsaasCluster.Properties.Cores = &cores
//	}
//
//	if ram, ok := d.GetOk("ram"); ok {
//		ram := int32(ram.(int))
//		dsaasCluster.Properties.Ram = &ram
//	}
//
//	if storageSize, ok := d.GetOk("storage_size"); ok {
//		storageSize := int32(storageSize.(int))
//		dsaasCluster.Properties.StorageSize = &storageSize
//	}
//
//	dsaasCluster.Properties.Connections = GetdsaasClusterConnectionsData(d)
//
//	if displayName, ok := d.GetOk("display_name"); ok {
//		displayName := displayName.(string)
//		dsaasCluster.Properties.DisplayName = &displayName
//	}
//
//	dsaasCluster.Properties.MaintenanceWindow = GetdsaasClusterMaintenanceWindowData(d)
//
//	return &dsaasCluster, nil
//}
//
//func GetdsaasClusterConnectionsData(d *schema.ResourceData) *[]dsaas.Connection {
//	connections := make([]dsaas.Connection, 0)
//
//	if vdcValue, ok := d.GetOk("connections"); ok {
//		vdcValue := vdcValue.([]interface{})
//		if vdcValue != nil {
//			for vdcIndex := range vdcValue {
//
//				connection := dsaas.Connection{}
//
//				if datacenterId, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", vdcIndex)); ok {
//					datacenterId := datacenterId.(string)
//					connection.DatacenterId = &datacenterId
//				}
//
//				if lanId, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", vdcIndex)); ok {
//					lanId := lanId.(string)
//					connection.LanId = &lanId
//				}
//
//				if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.cidr", vdcIndex)); ok {
//					cidr := cidr.(string)
//					connection.Cidr = &cidr
//				}
//
//				connections = append(connections, connection)
//			}
//		}
//
//	}
//
//	return &connections
//}
//
//func GetdsaasClusterMaintenanceWindowData(d *schema.ResourceData) *dsaas.MaintenanceWindow {
//	var maintenanceWindow dsaas.MaintenanceWindow
//
//	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
//		timeV := timeV.(string)
//		maintenanceWindow.Time = &timeV
//	}
//
//	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
//		dayOfTheWeek := dsaas.DayOfTheWeek(dayOfTheWeek.(string))
//		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
//	}
//
//	return &maintenanceWindow
//}
//
//func GetdsaasClusterCredentialsData(d *schema.ResourceData) *dsaas.DBUser {
//	var user dsaas.DBUser
//
//	if username, ok := d.GetOk("credentials.0.username"); ok {
//		username := username.(string)
//		user.Username = &username
//	}
//
//	if password, ok := d.GetOk("credentials.0.password"); ok {
//		password := password.(string)
//		user.Password = &password
//	}
//
//	return &user
//}
//
//func GetdsaasClusterFromBackupData(d *schema.ResourceData) (*dsaas.CreateRestoreRequest, error) {
//	var restore dsaas.CreateRestoreRequest
//
//	if backupId, ok := d.GetOk("from_backup.0.backup_id"); ok {
//		backupId := backupId.(string)
//		restore.BackupId = &backupId
//	}
//
//	if targetTime, ok := d.GetOk("from_backup.0.recovery_target_time"); ok {
//		var ionosTime dsaas.IonosTime
//		targetTime := targetTime.(string)
//		layout := "2006-01-02T15:04:05Z"
//		convertedTime, err := time.Parse(layout, targetTime)
//		if err != nil {
//			return nil, fmt.Errorf("an error occured while converting recovery_target_time to time.Time: %s", err)
//
//		}
//		ionosTime.Time = convertedTime
//		restore.RecoveryTargetTime = &ionosTime
//	}
//
//	return &restore, nil
//}
//
//func SetdsaasPgSqlClusterData(d *schema.ResourceData, cluster dsaas.ClusterResponse) error {
//
//	resourceName := "dsaas cluster"
//
//	if cluster.Id != nil {
//		d.SetId(*cluster.Id)
//	}
//
//	if cluster.Properties.PostgresVersion != nil {
//		if err := d.Set("postgres_version", *cluster.Properties.PostgresVersion); err != nil {
//			return utils.GenerateSetError(resourceName, "postgres_version", err)
//		}
//	}
//
//	if cluster.Properties.Instances != nil {
//		if err := d.Set("instances", *cluster.Properties.Instances); err != nil {
//			return utils.GenerateSetError(resourceName, "instances", err)
//		}
//	}
//
//	if cluster.Properties.Cores != nil {
//		if err := d.Set("cores", *cluster.Properties.Cores); err != nil {
//			return utils.GenerateSetError(resourceName, "cores", err)
//		}
//	}
//
//	if cluster.Properties.Ram != nil {
//		if err := d.Set("ram", *cluster.Properties.Ram); err != nil {
//			return utils.GenerateSetError(resourceName, "ram", err)
//		}
//	}
//
//	if cluster.Properties.StorageSize != nil {
//		if err := d.Set("storage_size", *cluster.Properties.StorageSize); err != nil {
//			return utils.GenerateSetError(resourceName, "storage_size", err)
//		}
//	}
//
//	if cluster.Properties.StorageType != nil {
//		if err := d.Set("storage_type", *cluster.Properties.StorageType); err != nil {
//			return utils.GenerateSetError(resourceName, "storage_type", err)
//		}
//	}
//
//	if cluster.Properties.Connections != nil && len(*cluster.Properties.Connections) > 0 {
//		var connections []interface{}
//		for _, connection := range *cluster.Properties.Connections {
//			connectionEntry := SetConnectionProperties(connection)
//			connections = append(connections, connectionEntry)
//		}
//		if err := d.Set("connections", connections); err != nil {
//			return utils.GenerateSetError(resourceName, "connections", err)
//		}
//	}
//
//	if cluster.Properties.Location != nil {
//		if err := d.Set("location", *cluster.Properties.Location); err != nil {
//			return fmt.Errorf("error while setting location property for dsaas cluster %s: %s", d.Id(), err)
//		}
//	}
//
//	if cluster.Properties.BackupLocation != nil {
//		if err := d.Set("backup_location", *cluster.Properties.BackupLocation); err != nil {
//			return fmt.Errorf("error while setting backup_location property for dsaas cluster %s: %s", d.Id(), err)
//		}
//	}
//
//	if cluster.Properties.DisplayName != nil {
//		if err := d.Set("display_name", *cluster.Properties.DisplayName); err != nil {
//			return fmt.Errorf("error while setting display_name property for dsaas cluster %s: %s", d.Id(), err)
//		}
//	}
//
//	if cluster.Properties.MaintenanceWindow != nil {
//		var maintenanceWindow []interface{}
//		maintenanceWindowEntry := SetMaintenanceWindowProperties(*cluster.Properties.MaintenanceWindow)
//		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
//		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
//			return utils.GenerateSetError(resourceName, "maintenance_window", err)
//		}
//	}
//
//	if cluster.Properties.SynchronizationMode != nil {
//		if err := d.Set("synchronization_mode", *cluster.Properties.SynchronizationMode); err != nil {
//			return fmt.Errorf("error while setting SynchronizationMode property for dsaas cluster %s: %s", d.Id(), err)
//		}
//	}
//
//	return nil
//}
//
//func SetConnectionProperties(vdcConnection dsaas.Connection) map[string]interface{} {
//
//	connection := map[string]interface{}{}
//
//	utils.SetPropWithNilCheck(connection, "datacenter_id", vdcConnection.DatacenterId)
//	utils.SetPropWithNilCheck(connection, "lan_id", vdcConnection.LanId)
//	utils.SetPropWithNilCheck(connection, "cidr", vdcConnection.Cidr)
//
//	return connection
//}
//
//func SetMaintenanceWindowProperties(maintenanceWindow dsaas.MaintenanceWindow) map[string]interface{} {
//
//	maintenance := map[string]interface{}{}
//
//	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
//	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)
//
//	return maintenance
//}
