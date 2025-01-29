package dbaas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func (c *PsqlClient) GetCluster(ctx context.Context, clusterId string) (psql.ClusterResponse, *psql.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// GetCluster Retrieves a Mongo cluster
func (c *MongoClient) GetCluster(ctx context.Context, clusterID string) (mongo.ClusterResponse, *shared.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

func (c *PsqlClient) ListClusters(ctx context.Context, filterName string) (psql.ClusterList, *psql.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGetExecute(request)
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// ListClusters Lists Mongo clusters
func (c *MongoClient) ListClusters(ctx context.Context, filterName string) (mongo.ClusterList, *shared.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGetExecute(request)
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// GetTemplates Lists Mongo templates
func (c *MongoClient) GetTemplates(ctx context.Context) (mongo.TemplateList, *shared.APIResponse, error) {
	templates, apiResponse, err := c.sdkClient.TemplatesApi.TemplatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return templates, apiResponse, err
}

func (c *PsqlClient) CreateCluster(ctx context.Context, cluster psql.CreateClusterRequest) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// CreateCluster Creates a Mongo cluster
func (c *MongoClient) CreateCluster(ctx context.Context, cluster mongo.CreateClusterRequest) (mongo.ClusterResponse, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// UpdateCluster Updates a Mongo cluster
func (c *MongoClient) UpdateCluster(ctx context.Context, clusterID string, cluster mongo.PatchClusterRequest) (mongo.ClusterResponse, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPatch(ctx, clusterID).PatchClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *PsqlClient) UpdateCluster(ctx context.Context, clusterId string, cluster psql.PatchClusterRequest) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPatch(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *PsqlClient) DeleteCluster(ctx context.Context, clusterId string) (psql.ClusterResponse, *psql.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// DeleteCluster Deletes a Mongo cluster
func (c *MongoClient) DeleteCluster(ctx context.Context, clusterID string) (mongo.ClusterResponse, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *PsqlClient) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	var clusterId string
	// This can be called from the users service, in which case the cluster_id is defined inside
	// the user resource, hence this if since the ID can be obtained in different ways depending on
	// the caller.
	if clusterIdIntf, ok := d.GetOk("cluster_id"); ok {
		clusterId = clusterIdIntf.(string)
	} else {
		clusterId = d.Id()
	}
	cluster, _, err := c.GetCluster(ctx, clusterId)
	if err != nil {
		return true, fmt.Errorf("check failed for cluster status: %w", err)
	}

	if cluster.Metadata == nil || cluster.Metadata.State == nil {
		return false, fmt.Errorf("cluster metadata or state is empty for id %s", clusterId)
	}

	log.Printf("[INFO] state of the cluster %s ", string(*cluster.Metadata.State))
	return strings.EqualFold(string(*cluster.Metadata.State), constant.Available), nil
}

func (c *PsqlClient) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.GetCluster(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for psql cluster deletion status: %w", err)
	}
	return false, nil
}

func (c *MongoClient) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cluster, _, err := c.GetCluster(ctx, d.Id())
	if err != nil {
		return true, fmt.Errorf("check failed for cluster status: %w", err)
	}

	if cluster.Metadata == nil || cluster.Metadata.State == nil {
		return false, fmt.Errorf("cluster metadata or state is empty for id %s", d.Id())
	}

	log.Printf("[INFO] state of the cluster %s ", string(*cluster.Metadata.State))
	return strings.EqualFold(string(*cluster.Metadata.State), constant.Available), nil
}

func (c *MongoClient) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.GetCluster(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for cluster deletion status: %w", err)
	}
	return false, nil
}

func GetPgSqlClusterDataCreate(d *schema.ResourceData) (*psql.CreateClusterRequest, error) {

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

	if _, ok := d.GetOk("connection_pooler"); ok {
		dbaasCluster.Properties.ConnectionPooler = getConnectionPoolerData(d)
	}

	if _, ok := d.GetOk("connections"); ok {
		dbaasCluster.Properties.Connections = GetPsqlClusterConnectionsData(d)
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
		dbaasCluster.Properties.MaintenanceWindow = GetPsqlClusterMaintenanceWindowData(d)
	}

	dbaasCluster.Properties.Credentials = GetPsqlClusterCredentialsData(d)

	if synchronizationMode, ok := d.GetOk("synchronization_mode"); ok {
		synchronizationMode := psql.SynchronizationMode(synchronizationMode.(string))
		dbaasCluster.Properties.SynchronizationMode = &synchronizationMode
	}

	if _, ok := d.GetOk("from_backup"); ok {
		if fromBackup, err := GetPsqlClusterFromBackupData(d); err != nil {
			return nil, err
		} else {
			dbaasCluster.Properties.FromBackup = fromBackup
		}
	}

	return &dbaasCluster, nil
}

func SetMongoClusterCreateProperties(d *schema.ResourceData) (*mongo.CreateClusterRequest, error) {

	mongoCluster := mongo.CreateClusterRequest{
		Properties: &mongo.CreateClusterProperties{},
	}

	if templateId, ok := d.GetOk("template_id"); ok {
		templateId := templateId.(string)
		mongoCluster.Properties.TemplateID = &templateId
	}

	if mongoVersion, ok := d.GetOk("mongodb_version"); ok {
		mongoVersion := mongoVersion.(string)
		mongoCluster.Properties.MongoDBVersion = &mongoVersion
	}

	if instances, ok := d.GetOk("instances"); ok {
		instances := instances.(int)
		mongoInstances := int32(instances)
		mongoCluster.Properties.Instances = mongoInstances
	}

	connections, err := GetMongoClusterConnectionsData(d)
	if err != nil {
		return nil, err
	}
	mongoCluster.Properties.Connections = connections

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		mongoCluster.Properties.Location = location
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		mongoCluster.Properties.DisplayName = displayName
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		mongoCluster.Properties.MaintenanceWindow = GetMongoClusterMaintenanceWindowData(d)
	}

	// enterprise settings below
	if clusterType, ok := d.GetOk("type"); ok {
		clusterType := clusterType.(string)
		mongoCluster.Properties.Type = &clusterType
	}

	if shards, ok := d.GetOk("shards"); ok {
		shards := (int32)(shards.(int))
		mongoCluster.Properties.Shards = &shards
	}

	if _, ok := d.GetOk("bi_connector"); ok {
		mongoCluster.Properties.BiConnector = GetMongoBiConnectorData(d)
	}

	if ram, ok := d.GetOk("ram"); ok {
		val := ram.(int)
		ram := (int32)(val)
		mongoCluster.Properties.Ram = &ram
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		val := storageSize.(int)
		storageSize := (int32)(val)
		mongoCluster.Properties.StorageSize = &storageSize
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageType := mongo.StorageType(storageType.(string))
		mongoCluster.Properties.StorageType = &storageType
	}

	if cores, ok := d.GetOk("cores"); ok {
		val := cores.(int)
		cores := (int32)(val)
		mongoCluster.Properties.Cores = &cores
	}

	if edition, ok := d.GetOk("edition"); ok {
		edition := edition.(string)
		mongoCluster.Properties.Edition = &edition
	}
	// to be added when there is api support
	// if _, ok := d.GetOk("from_backup"); ok {
	//	var fromBackup *mongo.CreateRestoreRequest
	//	fromBackup, err := GetMongoClusterFromBackupData(d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	mongoCluster.Properties.FromBackup = fromBackup
	//}

	if _, ok := d.GetOk("backup"); ok {
		var backup *mongo.BackupProperties
		backup = GetMongoClusterBackupData(d)
		mongoCluster.Properties.Backup = backup
	}

	return &mongoCluster, nil
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
		patchRequest.Properties.Connections, _ = GetMongoClusterConnectionsData(d)
	}

	if d.HasChange("maintenance_window") {
		_, mWin := d.GetChange("maintenance_window")
		if mWin != nil {
			mWinVal := GetMongoClusterMaintenanceWindowData(d)
			patchRequest.Properties.MaintenanceWindow = mWinVal
		}
	}

	// enterprise settings below
	if d.HasChange("type") {
		_, val := d.GetChange("type")
		clusterStr := val.(string)
		patchRequest.Properties.Type = &clusterStr
	}

	if d.HasChange("shards") {
		_, val := d.GetChange("shards")
		shards := int32(val.(int))
		patchRequest.Properties.Shards = &shards
	}

	if d.HasChange("ram") {
		_, val := d.GetChange("ram")
		ram := int32(val.(int))
		patchRequest.Properties.Ram = &ram
	}

	if d.HasChange("storage_size") {
		_, val := d.GetChange("storage_size")
		storageSize := int32(val.(int))
		patchRequest.Properties.StorageSize = &storageSize
	}

	if d.HasChange("storage_type") {
		_, val := d.GetChange("storage_type")
		storageType := mongo.StorageType(val.(string))
		patchRequest.Properties.StorageType = &storageType
	}

	if d.HasChange("cores") {
		_, val := d.GetChange("cores")
		cores := int32(val.(int))
		patchRequest.Properties.Cores = &cores
	}

	// must always be sent for enterprise, will be taken from template_id if playground or business
	_, val := d.GetChange("edition")
	if val.(string) == "enterprise" {
		edition := val.(string)
		patchRequest.Properties.Edition = &edition
	}

	if d.HasChange("bi_connector") {
		_, val := d.GetChange("bi_connector")
		if val != nil {
			patchRequest.Properties.BiConnector = GetMongoBiConnectorData(d)
		}
	}

	if d.HasChange("backup") {
		_, val := d.GetChange("backup")
		if val != nil {
			patchRequest.Properties.Backup = GetMongoClusterBackupData(d)
		}
	}

	return &patchRequest
}

func GetPgSqlClusterDataUpdate(d *schema.ResourceData) (*psql.PatchClusterRequest, diag.Diagnostics) {

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

	if _, ok := d.GetOk("connection_pooler"); ok {
		dbaasCluster.Properties.ConnectionPooler = getConnectionPoolerData(d)
	}

	dbaasCluster.Properties.Connections = GetPsqlClusterConnectionsData(d)

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.Properties.DisplayName = &displayName
	}

	dbaasCluster.Properties.MaintenanceWindow = GetPsqlClusterMaintenanceWindowData(d)

	return &dbaasCluster, nil
}

func GetPsqlClusterConnectionsData(d *schema.ResourceData) *[]psql.Connection {
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

// GetMongoClusterConnectionsData creates an sdk object for the Mongo connection list from the plan
func GetMongoClusterConnectionsData(d *schema.ResourceData) ([]mongo.Connection, error) {
	connections := make([]mongo.Connection, 0)

	if vdcValue, ok := d.GetOk("connections"); ok {
		vdcValue := vdcValue.([]interface{})
		if vdcValue != nil {
			for vdcIndex := range vdcValue {

				connection := mongo.Connection{}
				if datacenterId, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", vdcIndex)); ok {
					datacenterId := datacenterId.(string)
					connection.DatacenterId = datacenterId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = lanId
				}

				if cidrList, ok := d.GetOk(fmt.Sprintf("connections.%d.cidr_list", vdcIndex)); ok {
					cidrList := cidrList.([]interface{})
					var list []string
					for _, cidr := range cidrList {
						list = append(list, cidr.(string))
					}
					connection.CidrList = list
				}

				// if val, ok := d.GetOk(fmt.Sprintf("connections.%d.whitelist", vdcIndex)); ok {
				//	whitelist := val.([]interface{})
				//	if len(whitelist) > 0 {
				//
				//		list := make([]string, len(whitelist))
				//		err := utils.DecodeInterfaceToStruct(whitelist, list)
				//		if err != nil {
				//			return nil, fmt.Errorf("could not decode whitelist from %+v (%w)", whitelist, err)
				//		}
				//		if len(list) > 0 {
				//			connection.Whitelist = &list
				//		}
				//	}
				//}
				connections = append(connections, connection)
			}
		}
	}

	return connections, nil
}

func GetPsqlClusterMaintenanceWindowData(d *schema.ResourceData) *psql.MaintenanceWindow {
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

func GetMongoClusterMaintenanceWindowData(d *schema.ResourceData) *mongo.MaintenanceWindow {
	var maintenanceWindow mongo.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := mongo.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = dayOfTheWeek
	}

	return &maintenanceWindow
}

func GetMongoBiConnectorData(d *schema.ResourceData) *mongo.BiConnectorProperties {
	var biConnector mongo.BiConnectorProperties

	if enabled, ok := d.GetOk("bi_connector.0.enabled"); ok {
		timeV := enabled.(bool)
		biConnector.Enabled = &timeV
	}

	if port, ok := d.GetOk("bi_connector.0.host"); ok {
		port := port.(string)
		biConnector.Port = &port
	}

	if host, ok := d.GetOk("bi_connector.0.host"); ok {
		host := host.(string)
		biConnector.Host = &host
	}

	return &biConnector
}

func GetPsqlClusterCredentialsData(d *schema.ResourceData) *psql.DBUser {
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

func GetPsqlClusterFromBackupData(d *schema.ResourceData) (*psql.CreateRestoreRequest, error) {
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
			return nil, fmt.Errorf("an error occurred while converting recovery_target_time to time.Time: %w", err)

		}
		ionosTime.Time = convertedTime
		restore.RecoveryTargetTime = &ionosTime
	}

	return &restore, nil
}

func GetMongoClusterFromBackupData(d *schema.ResourceData) (*mongo.CreateRestoreRequest, error) {
	var restore mongo.CreateRestoreRequest

	if val, ok := d.GetOk("from_backup.0.snapshot_id"); ok {
		snapshotId := val.(string)
		restore.SnapshotId = &snapshotId
	}

	if targetTime, ok := d.GetOk("from_backup.0.recovery_target_time"); ok {
		var ionosTime mongo.IonosTime
		targetTime := targetTime.(string)
		layout := "2006-01-02T15:04:05Z"
		convertedTime, err := time.Parse(layout, targetTime)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while converting recovery_target_time to time.Time: %w", err)

		}
		ionosTime.Time = convertedTime
		restore.RecoveryTargetTime = &ionosTime
	}

	return &restore, nil
}

func GetMongoClusterBackupData(d *schema.ResourceData) *mongo.BackupProperties {
	var backup mongo.BackupProperties
	// to be added when backend supports the fields
	// if val, ok := d.GetOk("backup.0.snapshot_interval_hours"); ok {
	//	interval := int32(val.(int))
	//	backup.SnapshotIntervalHours = &interval
	//}
	//
	//if val, ok := d.GetOk("backup.0.point_in_time_window_hours"); ok {
	//	pointInTime := int32(val.(int))
	//	backup.PointInTimeWindowHours = &pointInTime
	//}

	if val, ok := d.GetOk("backup.0.location"); ok {
		location := val.(string)
		backup.Location = &location
	}
	// to be added at a later date
	//if _, ok := d.GetOk("backup.0.backup_retention"); ok {
	//	retention := GetMongoClusterBackupRetentionData(d)
	//	backup.BackupRetention = retention
	//}

	return &backup
}

// GetMongoClusterBackupRetentionData will be when we have support in backend
// func GetMongoClusterBackupRetentionData(d *schema.ResourceData) *mongo.BackupRetentionProperties {
//	var backup mongo.BackupRetentionProperties
//	path := "backup.0.backup_retention.0."
//	if val, ok := d.GetOk(path + "snapshot_retention_days"); ok {
//		days := int32(val.(int))
//		backup.SnapshotRetentionDays = &days
//	}
//
//	if val, ok := d.GetOk(path + "daily_snapshot_retention_days"); ok {
//		days := int32(val.(int))
//		backup.SnapshotRetentionDays = &days
//	}
//
//	if val, ok := d.GetOk(path + "weekly_snapshot_retention_weeks"); ok {
//		weeks := int32(val.(int))
//		backup.SnapshotRetentionDays = &weeks
//	}
//
//	if val, ok := d.GetOk(path + "monthly_snapshot_retention_months"); ok {
//		months := int32(val.(int))
//		backup.SnapshotRetentionDays = &months
//	}
//
//	return &backup
//}

func SetPgSqlClusterData(d *schema.ResourceData, cluster psql.ClusterResponse) error {

	resourceName := "psql cluster"

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties == nil {
		return fmt.Errorf("cluster properties should not be empty for id %s", *cluster.Id)
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

	if cluster.Properties.ConnectionPooler != nil {
		var connectionPooler []interface{}
		connectionPoolerEntry := setConnectionPoolerProperties(*cluster.Properties.ConnectionPooler)
		connectionPooler = append(connectionPooler, connectionPoolerEntry)
		if err := d.Set("connection_pooler", connectionPooler); err != nil {
			return utils.GenerateSetError(resourceName, "connection_pooler", err)
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

	if cluster.Properties.DnsName != nil {
		if err := d.Set("dns_name", *cluster.Properties.DnsName); err != nil {
			return utils.GenerateSetError(resourceName, "dns_name", err)
		}
	}

	return nil
}

func SetMongoDBClusterData(d *schema.ResourceData, cluster mongo.ClusterResponse) error {

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
		if len(cluster.Properties.Connections) > 0 {
			var connections []interface{}
			for _, connection := range cluster.Properties.Connections {
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

		// enterprise edition below
		if cluster.Properties.Type != nil {
			if err := d.Set("type", *cluster.Properties.Type); err != nil {
				return utils.GenerateSetError(resourceName, "type", err)
			}
		}

		if cluster.Properties.Shards != nil {
			if err := d.Set("shards", *cluster.Properties.Shards); err != nil {
				return utils.GenerateSetError(resourceName, "shards", err)
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

		if cluster.Properties.Cores != nil {
			if err := d.Set("cores", *cluster.Properties.Cores); err != nil {
				return utils.GenerateSetError(resourceName, "cores", err)
			}
		}
		if cluster.Properties.Edition != nil {
			if err := d.Set("edition", *cluster.Properties.Edition); err != nil {
				return utils.GenerateSetError(resourceName, "edition", err)
			}
		}

		if cluster.Properties.BiConnector != nil {
			var biConnector []interface{}
			conEntry := SetMongoBiConnectorProperties(*cluster.Properties.BiConnector)
			biConnector = append(biConnector, conEntry)
			if err := d.Set("bi_connector", biConnector); err != nil {
				return utils.GenerateSetError(resourceName, "bi_connector", err)
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
	// to ba added when there is backend support
	//if vdcConnection.Whitelist != nil {
	//	utils.SetPropWithNilCheck(connection, "whitelist", vdcConnection.Whitelist)
	//}

	return connection
}

func SetMongoMaintenanceWindowProperties(maintenanceWindow mongo.MaintenanceWindow) map[string]interface{} {
	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}

func SetMongoBiConnectorProperties(biConnectorProperties mongo.BiConnectorProperties) map[string]interface{} {
	biCon := map[string]interface{}{}
	if biConnectorProperties.Enabled != nil {
		utils.SetPropWithNilCheck(biCon, "enabled", *biConnectorProperties.Enabled)
	}
	utils.SetPropWithNilCheck(biCon, "port", biConnectorProperties.Port)
	utils.SetPropWithNilCheck(biCon, "host", biConnectorProperties.Host)

	return biCon
}

func SetMongoDBTemplateData(d *schema.ResourceData, template mongo.TemplateResponse) error {
	resourceName := "dbaas mongo template"

	if template.Id != nil {
		d.SetId(*template.Id)
	}
	if template.Properties != nil {
		if template.Properties.Name != nil {
			field := "name"
			if err := d.Set(field, *template.Properties.Name); err != nil {
				return utils.GenerateSetError(resourceName, field, err)
			}
		}
		if template.Properties.Edition != nil {
			field := "edition"
			if err := d.Set(field, *template.Properties.Edition); err != nil {
				return utils.GenerateSetError(resourceName, field, err)
			}
		}
		if template.Properties.Cores != nil {
			field := "cores"
			if err := d.Set(field, *template.Properties.Cores); err != nil {
				return utils.GenerateSetError(resourceName, field, err)
			}
		}
		if template.Properties.Ram != nil {
			field := "ram"
			if err := d.Set(field, *template.Properties.Ram); err != nil {
				return utils.GenerateSetError(resourceName, field, err)
			}
		}
		if template.Properties.StorageSize != nil {
			field := "storage_size"
			if err := d.Set(field, *template.Properties.StorageSize); err != nil {
				return utils.GenerateSetError(resourceName, field, err)
			}
		}
	}
	return nil
}

// MongoClusterCheckRequiredFieldsSet Checks if required fields are set in the cluster resource
func MongoClusterCheckRequiredFieldsSet(d *schema.ResourceData) error {

	clusterType := d.Get("edition").(string)
	requiredNotSet := "%s argument must be set for %s edition of mongo cluster"
	// if clusterTYpe != "" {
	//	server.Properties.Type = &serverType
	//}
	switch strings.ToLower(clusterType) {
	case "enterprise":

		if _, ok := d.GetOk("cores"); !ok {
			return fmt.Errorf(requiredNotSet, "cores", clusterType)
		}

		if _, ok := d.GetOk("ram"); !ok {
			return fmt.Errorf(requiredNotSet, "ram", clusterType)
		}

		if _, ok := d.GetOk("storage_size"); !ok {
			return fmt.Errorf(requiredNotSet, "storage_size", clusterType)
		}

		if _, ok := d.GetOk("storage_type"); !ok {
			return fmt.Errorf(requiredNotSet, "storage_type", clusterType)
		}

		if _, ok := d.GetOk("template_id"); ok {
			return fmt.Errorf("%s argument must NOT be set for %s edition of mongo cluster", "template_id", clusterType)
		}

	default: // playground or business
		if _, ok := d.GetOk("template_id"); !ok {
			return fmt.Errorf(requiredNotSet, "template_id", clusterType)
		}
	}
	return nil
}

func getConnectionPoolerData(d *schema.ResourceData) *psql.ConnectionPooler {
	var connectionPooler psql.ConnectionPooler

	enabledIntf := d.Get("connection_pooler.0.enabled")
	enabledValue := enabledIntf.(bool)
	connectionPooler.Enabled = &enabledValue

	poolModeIntf := d.Get("connection_pooler.0.pool_mode")
	poolModeValue := psql.PoolMode(poolModeIntf.(string))
	connectionPooler.PoolMode = &poolModeValue

	return &connectionPooler
}

func setConnectionPoolerProperties(connectionPoolerProperties psql.ConnectionPooler) map[string]interface{} {
	connectionPoolerMap := map[string]interface{}{}

	utils.SetPropWithNilCheck(connectionPoolerMap, "enabled", *connectionPoolerProperties.Enabled)
	utils.SetPropWithNilCheck(connectionPoolerMap, "pool_mode", *connectionPoolerProperties.PoolMode)

	return connectionPoolerMap
}
