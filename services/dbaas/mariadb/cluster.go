package mariadb

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var locationToURL = map[string]string{
	"":       "https://mariadb.de-txl.ionos.com",
	"de/fra": "https://mariadb.de-fra.ionos.com",
	"de/txl": "https://mariadb.de-txl.ionos.com",
	"es/vit": "https://mariadb.es-vit.ionos.com",
	"fr/par": "https://mariadb.fr-par.ionos.com",
	"gb/lhr": "https://mariadb.gb-lhr.ionos.com",
	"us/ewr": "https://mariadb.us-ewr.ionos.com",
	"us/las": "https://mariadb.us-las.ionos.com",
	"us/mci": "https://mariadb.us-mci.ionos.com",
}
var ionosAPIURLMariaDB = "IONOS_API_URL_MARIADB"

// modifyConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *MariaDBClient) modifyConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLMariaDB) != "" {
		clientConfig.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLMariaDB)),
			},
		}
		return
	}
	clientConfig.Servers = shared.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

// GetCluster retrieves a cluster by its ID and the location in which the cluster is created.
func (c *MariaDBClient) GetCluster(ctx context.Context, clusterID, location string) (mariadb.ClusterResponse, *shared.APIResponse, error) {
	c.modifyConfigURL(location)
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListClusters retrieves a list of clusters based on the location. Filters can be used.
func (c *MariaDBClient) ListClusters(ctx context.Context, filterName, location string) (mariadb.ClusterList, *shared.APIResponse, error) {
	c.modifyConfigURL(location)
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGetExecute(request)
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// CreateCluster creates a new cluster using the provided data in the request and the location.
func (c *MariaDBClient) CreateCluster(ctx context.Context, cluster mariadb.CreateClusterRequest, location string) (mariadb.ClusterResponse, *shared.APIResponse, error) {
	c.modifyConfigURL(location)
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).CreateClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// UpdateCluster updates a cluster by its ID and the location in which the cluster is created.
func (c *MariaDBClient) UpdateCluster(ctx context.Context, cluster mariadb.PatchClusterRequest, clusterID, location string) (mariadb.ClusterResponse, *shared.APIResponse, error) {
	c.modifyConfigURL(location)
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPatch(ctx, clusterID).PatchClusterRequest(cluster).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// DeleteCluster deletes a cluster by its ID and the location in which the cluster is created.
func (c *MariaDBClient) DeleteCluster(ctx context.Context, clusterID, location string) (mariadb.ClusterResponse, *shared.APIResponse, error) {
	c.modifyConfigURL(location)
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *MariaDBClient) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	location := d.Get("location").(string)
	cluster, _, err := c.GetCluster(ctx, clusterID, location)
	if err != nil {
		return true, fmt.Errorf("status check failed for MariaDB cluster with ID: %v, error: %w", clusterID, err)
	}

	if cluster.Metadata == nil || cluster.Metadata.State == nil {
		return false, fmt.Errorf("cluster metadata or state is empty for MariaDB cluster with ID: %v", clusterID)
	}

	log.Printf("[INFO] state of the MariaDB cluster with ID: %v is: %s ", clusterID, string(*cluster.Metadata.State))
	return strings.EqualFold(string(*cluster.Metadata.State), constant.Available), nil
}

func (c *MariaDBClient) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	_, apiResponse, err := c.GetCluster(ctx, clusterID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for MariaDB cluster deletion status, cluster ID: %v, error: %w", clusterID, err)
	}
	return false, nil
}

func GetMariaDBClusterDataCreate(d *schema.ResourceData) (*mariadb.CreateClusterRequest, error) {
	cluster := mariadb.CreateClusterRequest{
		Properties: &mariadb.CreateClusterProperties{},
	}

	if mariaDBVersion, ok := d.GetOk("mariadb_version"); ok {
		mariaDBVersion := mariaDBVersion.(string)
		cluster.Properties.MariadbVersion = (mariadb.MariadbVersion)(mariaDBVersion)
	}

	if instances, ok := d.GetOk("instances"); ok {
		instances := int32(instances.(int))
		cluster.Properties.Instances = instances
	}

	if cores, ok := d.GetOk("cores"); ok {
		cores := int32(cores.(int))
		cluster.Properties.Cores = cores
	}

	if ram, ok := d.GetOk("ram"); ok {
		ram := int32(ram.(int))
		cluster.Properties.Ram = ram
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		storageSize := int32(storageSize.(int))
		cluster.Properties.StorageSize = storageSize
	}

	if _, ok := d.GetOk("connections"); ok {
		cluster.Properties.Connections = GetMariaClusterConnectionsData(d)
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		cluster.Properties.DisplayName = displayName
	}

	cluster.Properties.Credentials = GetMariaClusterCredentialsData(d)

	if _, ok := d.GetOk("maintenance_window"); ok {
		cluster.Properties.MaintenanceWindow = GetMariaClusterMaintenanceWindowData(d)
	}

	if _, ok := d.GetOk("backup"); ok {
		cluster.Properties.Backup = GetMariaClusterBackupData(d)
	}

	return &cluster, nil
}

// GetMariaDBClusterDataUpdate retrieves the data from the terraform resource and sets it in the MariaDB cluster struct.
func GetMariaDBClusterDataUpdate(d *schema.ResourceData) (*mariadb.PatchClusterRequest, error) {
	cluster := mariadb.PatchClusterRequest{
		Properties: &mariadb.PatchClusterProperties{},
	}

	if d.HasChange("mariadb_version") {
		_, newValue := d.GetChange("mariadb_version")
		newVersion := newValue.(string)
		cluster.Properties.MariadbVersion = (*mariadb.MariadbVersion)(&newVersion)
	}

	if d.HasChange("instances") {
		_, n := d.GetChange("instances")
		nInt := int32(n.(int))
		cluster.Properties.Instances = &nInt
	}

	if d.HasChange("cores") {
		_, n := d.GetChange("cores")
		nInt := int32(n.(int))
		cluster.Properties.Cores = &nInt
	}

	if d.HasChange("ram") {
		_, n := d.GetChange("ram")
		nInt := int32(n.(int))
		cluster.Properties.Ram = &nInt
	}

	if d.HasChange("storage_size") {
		_, n := d.GetChange("storage_size")
		nInt := int32(n.(int))
		cluster.Properties.StorageSize = &nInt
	}

	if d.HasChange("display_name") {
		_, n := d.GetChange("display_name")
		nString := n.(string)
		cluster.Properties.DisplayName = &nString
	}

	if d.HasChange("maintenance_window") {
		cluster.Properties.MaintenanceWindow = GetMariaClusterMaintenanceWindowData(d)
	}

	return &cluster, nil
}

// GetMariaClusterConnectionsData retrieves the data from the terraform resource and sets it in the MariaDB connection struct.
func GetMariaClusterConnectionsData(d *schema.ResourceData) []mariadb.Connection {
	connections := make([]mariadb.Connection, 0)

	if connectionsIntf, ok := d.GetOk("connections"); ok {
		connectionsValues := connectionsIntf.([]interface{})
		for connectionIdx := range connectionsValues {
			connection := mariadb.Connection{}

			if datacenterID, ok := d.GetOk(fmt.Sprintf("connections.%d.datacenter_id", connectionIdx)); ok {
				datacenterID := datacenterID.(string)
				connection.DatacenterId = datacenterID
			}

			if lanID, ok := d.GetOk(fmt.Sprintf("connections.%d.lan_id", connectionIdx)); ok {
				lanID := lanID.(string)
				connection.LanId = lanID
			}

			if cidr, ok := d.GetOk(fmt.Sprintf("connections.%d.cidr", connectionIdx)); ok {
				cidr := cidr.(string)
				connection.Cidr = cidr
			}
			connections = append(connections, connection)
		}
	}
	return connections
}

// GetMariaClusterMaintenanceWindowData retrieves the data from the terraform resource and sets it in the MariaDB MaintenanceWindow struct.
func GetMariaClusterMaintenanceWindowData(d *schema.ResourceData) *mariadb.MaintenanceWindow {
	var maintenanceWindow mariadb.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := mariadb.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = dayOfTheWeek
	}

	return &maintenanceWindow
}

// GetMariaClusterBackupData retrieves the data from the terraform resource and sets it in the MariaDB Backup struct.
func GetMariaClusterBackupData(d *schema.ResourceData) *mariadb.BackupProperties {
	var backup mariadb.BackupProperties

	if loc, ok := d.GetOk("backup.0.location"); ok {
		loc := loc.(string)
		backup.Location = &loc
	}

	return &backup
}

// GetMariaClusterCredentialsData retrieves the data from the terraform resource and sets it in the MariaDB DBUser struct.
func GetMariaClusterCredentialsData(d *schema.ResourceData) mariadb.DBUser {
	var user mariadb.DBUser

	if username, ok := d.GetOk("credentials.0.username"); ok {
		username := username.(string)
		user.Username = username
	}

	if password, ok := d.GetOk("credentials.0.password"); ok {
		password := password.(string)
		user.Password = password
	}

	return user
}

func (c *MariaDBClient) SetMariaDBClusterData(d *schema.ResourceData, cluster mariadb.ClusterResponse) error {

	resourceName := "MariaDB cluster"

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties == nil {
		return fmt.Errorf("response properties should not be empty for MariaDB cluster with ID %v", *cluster.Id)
	}

	if cluster.Properties.MariadbVersion != nil {
		if err := d.Set("mariadb_version", *cluster.Properties.MariadbVersion); err != nil {
			return utils.GenerateSetError(resourceName, "mariadb_version", err)
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

	if len(cluster.Properties.Connections) > 0 {
		var connections []interface{}
		for _, connection := range cluster.Properties.Connections {
			connectionEntry := c.SetConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}
		if err := d.Set("connections", connections); err != nil {
			return utils.GenerateSetError(resourceName, "connections", err)
		}
	}

	if cluster.Properties.DisplayName != nil {
		if err := d.Set("display_name", *cluster.Properties.DisplayName); err != nil {
			return utils.GenerateSetError(resourceName, "display_name", err)
		}
	}

	if cluster.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := c.SetMaintenanceWindowProperties(*cluster.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(resourceName, "maintenance_window", err)
		}
	}

	if cluster.Properties.Backup != nil {
		var bac []interface{}
		backupEntry := c.SetBackupProperties(*cluster.Properties.Backup)
		bac = append(bac, backupEntry)
		if err := d.Set("backup", bac); err != nil {
			return utils.GenerateSetError(resourceName, "backup", err)
		}
	}

	if cluster.Properties.DnsName != nil {
		if err := d.Set("dns_name", *cluster.Properties.DnsName); err != nil {
			return utils.GenerateSetError(resourceName, "dns_name", err)
		}
	}

	return nil
}

func (c *MariaDBClient) SetConnectionProperties(connection mariadb.Connection) map[string]interface{} {
	connectionMap := map[string]interface{}{}

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "cidr", connection.Cidr)

	return connectionMap
}

func (c *MariaDBClient) SetMaintenanceWindowProperties(maintenanceWindow mariadb.MaintenanceWindow) map[string]interface{} {
	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}

func (c *MariaDBClient) SetBackupProperties(backup mariadb.BackupProperties) map[string]interface{} {
	bac := map[string]interface{}{}

	utils.SetPropWithNilCheck(bac, "location", backup.Location)

	return bac
}
