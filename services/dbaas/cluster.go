package dbaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"time"
)

func (c *PsqlClient) GetCluster(ctx context.Context, clusterId string) (psql.ClusterResponse, *psql.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *MongoClient) GetCluster(ctx context.Context, clusterId string) (mongo.ClusterResponse, *mongo.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *PsqlClient) ListClusters(ctx context.Context, filterName string) (psql.ClusterList, *psql.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGetExecute(request)
	if apiResponse != nil {
		return clusters, apiResponse, err
	}
	return clusters, nil, err
}

func (c *MongoClient) ListClusters(ctx context.Context, filterName string) (mongo.ClusterList, *mongo.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGetExecute(request)
	if apiResponse != nil {
		return clusters, apiResponse, err
	}
	return clusters, nil, err
}

func (c *MongoClient) GetTemplates(ctx context.Context) (mongo.TemplateList, *mongo.APIResponse, error) {
	templates, apiResponse, err := c.sdkClient.TemplatesApi.TemplatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return templates, apiResponse, err
}

func (c *PsqlClient) CreateCluster(ctx context.Context, cluster psql.CreateClusterRequest) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *MongoClient) CreateCluster(ctx context.Context, cluster mongo.CreateClusterRequest) (mongo.ClusterResponse, *mongo.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *MongoClient) UpdateCluster(ctx context.Context, clusterId string, cluster mongo.PatchClusterRequest) (mongo.ClusterResponse, *mongo.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPatch(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *PsqlClient) UpdateCluster(ctx context.Context, clusterId string, cluster psql.PatchClusterRequest) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPatch(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *PsqlClient) DeleteCluster(ctx context.Context, clusterId string) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *MongoClient) DeleteCluster(ctx context.Context, clusterId string) (mongo.ClusterResponse, *mongo.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func GetDbaasPgSqlClusterDataCreate(d *schema.ResourceData) (*psql.CreateClusterRequest, error) {

	dbaasCluster := psql.CreateClusterRequest{
		Properties: &psql.CreateClusterProperties{},
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
		storageType := psql.StorageType(storageType.(string))
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
		synchronizationMode := psql.SynchronizationMode(synchronizationMode.(string))
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

func SetMongoClusterCreateProperties(d *schema.ResourceData) *mongo.CreateClusterRequest {

	dbaasCluster := mongo.CreateClusterRequest{
		Properties: &mongo.CreateClusterProperties{},
	}

	if templateId, ok := d.GetOk("template_id"); ok {
		templateId := templateId.(string)
		dbaasCluster.Properties.TemplateID = &templateId
	}

	if mongoVersion, ok := d.GetOk("mongodb_version"); ok {
		mongoVersion := mongoVersion.(string)
		dbaasCluster.Properties.MongoDBVersion = &mongoVersion
	}

	if instances, ok := d.GetOk("instances"); ok {
		instances := instances.(int)
		mongoInstances := int32(instances)
		dbaasCluster.Properties.Instances = &mongoInstances
	}

	dbaasCluster.Properties.Connections = GetDbaasMongoClusterConnectionsData(d)

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		dbaasCluster.Properties.Location = &location
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.Properties.DisplayName = &displayName
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dbaasCluster.Properties.MaintenanceWindow = GetDbaasMongoClusterMaintenanceWindowData(d)
	}

	return &dbaasCluster
}

func SetMongoClusterPatchProperties(d *schema.ResourceData) *mongo.PatchClusterRequest {

	patchRequest := mongo.PatchClusterRequest{
		Properties: mongo.NewPatchClusterProperties(),
	}

	if d.HasChange("display_name") {
		_, name := d.GetChange("display_name")
		nameStr := name.(string)
		patchRequest.Properties.DisplayName = &nameStr
	}

	if d.HasChange("instances") {
		_, instances := d.GetChange("instances")
		instancesInt := int32(instances.(int))
		patchRequest.Properties.Instances = &instancesInt
	}

	if d.HasChange("template_id") {
		_, template := d.GetChange("template_id")
		templateStr := template.(string)
		patchRequest.Properties.TemplateID = &templateStr
	}
	if d.HasChange("connections") {
		patchRequest.Properties.Connections = GetDbaasMongoClusterConnectionsData(d)
	}
	if d.HasChange("maintenance_window") {
		_, mWin := d.GetChange("maintenance_window")
		if mWin != nil {
			mWinVal := GetDbaasMongoClusterMaintenanceWindowData(d)
			patchRequest.Properties.MaintenanceWindow = mWinVal
		}
	}

	return &patchRequest
}

func GetDbaasPgSqlClusterDataUpdate(d *schema.ResourceData) (*psql.PatchClusterRequest, diag.Diagnostics) {

	dbaasCluster := psql.PatchClusterRequest{
		Properties: &psql.PatchClusterProperties{},
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

func GetDbaasClusterConnectionsData(d *schema.ResourceData) *[]psql.Connection {
	connections := make([]psql.Connection, 0)

	if vdcValue, ok := d.GetOk("connections"); ok {
		vdcValue := vdcValue.([]interface{})
		if vdcValue != nil {
			for vdcIndex := range vdcValue {

				connection := psql.Connection{}

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

func GetDbaasMongoClusterConnectionsData(d *schema.ResourceData) *[]mongo.Connection {
	connections := make([]mongo.Connection, 0)

	if vdcValue, ok := d.GetOk("connections"); ok {
		vdcValue := vdcValue.([]interface{})
		if vdcValue != nil {
			for vdcIndex := range vdcValue {

				connection := mongo.Connection{}

				if datacenterId, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", vdcIndex)); ok {
					datacenterId := datacenterId.(string)
					connection.DatacenterId = &datacenterId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = &lanId
				}

				if cidrList, ok := d.GetOk(fmt.Sprintf("connections.%d.cidr_list", vdcIndex)); ok {
					cidrList := cidrList.([]interface{})
					var list []string
					for _, cidr := range cidrList {
						list = append(list, cidr.(string))
					}
					connection.CidrList = &list
				}

				connections = append(connections, connection)
			}
		}

	}

	return &connections
}

func GetDbaasClusterMaintenanceWindowData(d *schema.ResourceData) *psql.MaintenanceWindow {
	var maintenanceWindow psql.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = &timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := psql.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}

	return &maintenanceWindow
}

func GetDbaasMongoClusterMaintenanceWindowData(d *schema.ResourceData) *mongo.MaintenanceWindow {
	var maintenanceWindow mongo.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = &timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := mongo.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}

	return &maintenanceWindow
}

func GetDbaasClusterCredentialsData(d *schema.ResourceData) *psql.DBUser {
	var user psql.DBUser

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

func GetDbaasClusterFromBackupData(d *schema.ResourceData) (*psql.CreateRestoreRequest, error) {
	var restore psql.CreateRestoreRequest

	if backupId, ok := d.GetOk("from_backup.0.backup_id"); ok {
		backupId := backupId.(string)
		restore.BackupId = &backupId
	}

	if targetTime, ok := d.GetOk("from_backup.0.recovery_target_time"); ok {
		var ionosTime psql.IonosTime
		targetTime := targetTime.(string)
		layout := "2006-01-02T15:04:05Z"
		convertedTime, err := time.Parse(layout, targetTime)
		if err != nil {
			return nil, fmt.Errorf("an error occured while converting recovery_target_time to time.Time: %w", err)

		}
		ionosTime.Time = convertedTime
		restore.RecoveryTargetTime = &ionosTime
	}

	return &restore, nil
}

func SetDbaasPgSqlClusterData(d *schema.ResourceData, cluster psql.ClusterResponse) error {

	resourceName := "psql cluster"

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
			return fmt.Errorf("error while setting location property for psql cluster %s: %w", d.Id(), err)
		}
	}

	if cluster.Properties.BackupLocation != nil {
		if err := d.Set("backup_location", *cluster.Properties.BackupLocation); err != nil {
			return fmt.Errorf("error while setting backup_location property for psql cluster %s: %w", d.Id(), err)
		}
	}

	if cluster.Properties.DisplayName != nil {
		if err := d.Set("display_name", *cluster.Properties.DisplayName); err != nil {
			return fmt.Errorf("error while setting display_name property for psql cluster %s: %w", d.Id(), err)
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
			return fmt.Errorf("error while setting SynchronizationMode property for psql cluster %s: %w", d.Id(), err)
		}
	}

	return nil
}

func SetDbaasMongoDBClusterData(d *schema.ResourceData, cluster mongo.ClusterResponse) error {

	resourceName := "dbaas mongo cluster"

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}
	if cluster.Properties != nil {

		if cluster.Properties.TemplateID != nil {
			if err := d.Set("template_id", *cluster.Properties.TemplateID); err != nil {
				return utils.GenerateSetError(resourceName, "template_id", err)
			}
		}
		if cluster.Properties.MongoDBVersion != nil {
			if err := d.Set("mongodb_version", *cluster.Properties.MongoDBVersion); err != nil {
				return utils.GenerateSetError(resourceName, "mongodb_version", err)
			}
		}
		if cluster.Properties.Instances != nil {
			if err := d.Set("instances", *cluster.Properties.Instances); err != nil {
				return utils.GenerateSetError(resourceName, "instances", err)
			}
		}
		if cluster.Properties.Connections != nil && len(*cluster.Properties.Connections) > 0 {
			var connections []interface{}
			for _, connection := range *cluster.Properties.Connections {
				connectionEntry := SetMongoConnectionProperties(connection)
				connections = append(connections, connectionEntry)
			}
			if err := d.Set("connections", connections); err != nil {
				return utils.GenerateSetError(resourceName, "connections", err)
			}
		}

		if cluster.Properties.Location != nil {
			if err := d.Set("location", *cluster.Properties.Location); err != nil {
				return fmt.Errorf("error while setting location property for psql cluster %s: %w", d.Id(), err)
			}
		}
		if cluster.Properties.DisplayName != nil {
			if err := d.Set("display_name", *cluster.Properties.DisplayName); err != nil {
				return fmt.Errorf("error while setting display_name property for psql cluster %s: %w", d.Id(), err)
			}
		}

		if cluster.Properties.ConnectionString != nil {
			if err := d.Set("connection_string", *cluster.Properties.ConnectionString); err != nil {
				return utils.GenerateSetError(resourceName, "connection_string", err)
			}
		}

		if cluster.Properties.MaintenanceWindow != nil {
			var maintenanceWindow []interface{}
			maintenanceWindowEntry := SetMongoMaintenanceWindowProperties(*cluster.Properties.MaintenanceWindow)
			maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
			if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
				return utils.GenerateSetError(resourceName, "maintenance_window", err)
			}
		}
	}

	return nil
}

func SetConnectionProperties(vdcConnection psql.Connection) map[string]interface{} {

	connection := map[string]interface{}{}

	utils.SetPropWithNilCheck(connection, "datacenter_id", vdcConnection.DatacenterId)
	utils.SetPropWithNilCheck(connection, "lan_id", vdcConnection.LanId)
	utils.SetPropWithNilCheck(connection, "cidr", vdcConnection.Cidr)

	return connection
}

func SetMaintenanceWindowProperties(maintenanceWindow psql.MaintenanceWindow) map[string]interface{} {

	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}

func SetMongoConnectionProperties(vdcConnection mongo.Connection) map[string]interface{} {
	connection := map[string]interface{}{}

	utils.SetPropWithNilCheck(connection, "datacenter_id", vdcConnection.DatacenterId)
	utils.SetPropWithNilCheck(connection, "lan_id", vdcConnection.LanId)
	utils.SetPropWithNilCheck(connection, "cidr_list", vdcConnection.CidrList)

	return connection
}

func SetMongoMaintenanceWindowProperties(maintenanceWindow mongo.MaintenanceWindow) map[string]interface{} {
	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}

func SetMongoDBTemplateData(d *schema.ResourceData, template mongo.TemplateResponse) error {
	resourceName := "dbaas mongo template"

	if template.Id != nil {
		d.SetId(*template.Id)
	}
	if template.Name != nil {
		field := "name"
		if err := d.Set(field, *template.Name); err != nil {
			return utils.GenerateSetError(resourceName, field, err)
		}
	}
	if template.Edition != nil {
		field := "edition"
		if err := d.Set(field, *template.Edition); err != nil {
			return utils.GenerateSetError(resourceName, field, err)
		}
	}
	if template.Cores != nil {
		field := "cores"
		if err := d.Set(field, *template.Cores); err != nil {
			return utils.GenerateSetError(resourceName, field, err)
		}
	}
	if template.Ram != nil {
		field := "ram"
		if err := d.Set(field, *template.Ram); err != nil {
			return utils.GenerateSetError(resourceName, field, err)
		}
	}
	if template.StorageSize != nil {
		field := "storage_size"
		if err := d.Set(field, *template.StorageSize); err != nil {
			return utils.GenerateSetError(resourceName, field, err)
		}
	}
	return nil
}

// todo: remove once mongo removes this field
const DefaultMongoDatabase = "admin"
