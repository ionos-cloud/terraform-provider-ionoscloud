package inmemorydb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	inMemoryDB "github.com/ionos-cloud/sdk-go-dbaas-in-memory-db"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var locationToURL = map[string]string{
	"de/fra": "https://in-memory-db.de-fra.ionos.com",
	"de/txl": "https://in-memory-db.de-txl.ionos.com",
	"es/vit": "https://in-memory-db.es-vit.ionos.com",
	"gb/lhr": "https://in-memory-db.gb-lhr.ionos.com",
	"us/ewr": "https://in-memory-db.us-ewr.ionos.com",
	"us/las": "https://in-memory-db.us-las.ionos.com",
	"us/mci": "https://in-memory-db.us-mci.ionos.com",
	"fr/par": "https://in-memory-db.fr-par.ionos.com",
}

// modifyConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *InMemoryDBClient) modifyConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	clientConfig.Servers = inMemoryDB.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

// CreateReplicaSet sends a 'POST' request to the API to create a replica set.
func (c *InMemoryDBClient) CreateReplicaSet(ctx context.Context, replicaSet inMemoryDB.ReplicaSetCreate, location string) (inMemoryDB.ReplicaSetRead, *inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	replicaSetResponse, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsPost(ctx).ReplicaSetCreate(replicaSet).Execute()
	apiResponse.LogInfo()
	return replicaSetResponse, apiResponse, err
}

//nolint:golint
func (c *InMemoryDBClient) IsReplicaSetReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	replicaSetID := d.Id()
	location := d.Get("location").(string)
	replicaSet, _, err := c.GetReplicaSet(ctx, replicaSetID, location)
	if err != nil {
		return false, fmt.Errorf("status check failed for InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err)
	}
	if replicaSet.Metadata == nil || replicaSet.Metadata.State == nil {
		return false, fmt.Errorf("metadata or state is empty for InMemoryDB replica set with ID: %v", replicaSetID)
	}
	log.Printf("[INFO] state of the InMemoryDB replica set with ID: %v is: %v", replicaSetID, *replicaSet.Metadata.State)
	return strings.EqualFold(*replicaSet.Metadata.State, constant.Available), nil
}

//nolint:golint
func (c *InMemoryDBClient) DeleteReplicaSet(ctx context.Context, replicaSetID, location string) (*inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsDelete(ctx, replicaSetID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

//nolint:golint
func (c *InMemoryDBClient) UpdateReplicaSet(ctx context.Context, replicaSetID, location string, replicaSet inMemoryDB.ReplicaSetEnsure) (inMemoryDB.ReplicaSetRead, *inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	replicaSetResponse, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsPut(ctx, replicaSetID).ReplicaSetEnsure(replicaSet).Execute()
	apiResponse.LogInfo()
	return replicaSetResponse, apiResponse, err
}

//nolint:golint
func (c *InMemoryDBClient) IsReplicaSetDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	replicaSetID := d.Id()
	_, apiResponse, err := c.GetReplicaSet(ctx, replicaSetID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err)
	}
	return false, nil
}

//nolint:golint
func (c *InMemoryDBClient) GetReplicaSet(ctx context.Context, replicaSetID, location string) (inMemoryDB.ReplicaSetRead, *inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	replicaSet, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsFindById(ctx, replicaSetID).Execute()
	apiResponse.LogInfo()
	return replicaSet, apiResponse, err
}

//nolint:golint
func (c *InMemoryDBClient) GetSnapshot(ctx context.Context, snapshotID, location string) (inMemoryDB.SnapshotRead, *inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	snapshot, apiResponse, err := c.sdkClient.SnapshotApi.SnapshotsFindById(ctx, snapshotID).Execute()
	apiResponse.LogInfo()
	return snapshot, apiResponse, err
}

//nolint:golint
func (c *InMemoryDBClient) ListReplicaSets(ctx context.Context, filterName, location string) (inMemoryDB.ReplicaSetReadList, *inMemoryDB.APIResponse, error) {
	c.modifyConfigURL(location)
	request := c.sdkClient.ReplicaSetApi.ReplicasetsGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	replicaSets, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsGetExecute(request)
	apiResponse.LogInfo()
	return replicaSets, apiResponse, err
}

// GetReplicaSetDataProperties reads all the properties from the configuration file and returns
// a structure that will be used to populate update/create requests.
func GetReplicaSetDataProperties(d *schema.ResourceData) *inMemoryDB.ReplicaSet {
	replicaSet := inMemoryDB.ReplicaSet{}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		replicaSet.DisplayName = &displayName
	}

	if inMemoryDBVersion, ok := d.GetOk("version"); ok {
		inMemoryDBVersion := inMemoryDBVersion.(string)
		replicaSet.Version = &inMemoryDBVersion
	}

	if replicas, ok := d.GetOk("replicas"); ok {
		replicas := int32(replicas.(int))
		replicaSet.Replicas = &replicas
	}

	if persistenceMode, ok := d.GetOk("persistence_mode"); ok {
		persistenceMode := inMemoryDB.PersistenceMode(persistenceMode.(string))
		replicaSet.PersistenceMode = &persistenceMode
	}

	if evictionPolicy, ok := d.GetOk("eviction_policy"); ok {
		evictionPolicy := inMemoryDB.EvictionPolicy(evictionPolicy.(string))
		replicaSet.EvictionPolicy = &evictionPolicy
	}

	if initialSnapshotID, ok := d.GetOk("initial_snapshot_id"); ok {
		initialSnapshotID := initialSnapshotID.(string)
		replicaSet.InitialSnapshotId = &initialSnapshotID
	}

	if _, ok := d.GetOk("resources"); ok {
		replicaSet.Resources = getResources(d)
	}

	if _, ok := d.GetOk("connections"); ok {
		replicaSet.Connections = getConnections(d)
	}

	if _, ok := d.GetOk("credentials"); ok {
		replicaSet.Credentials = getCredentials(d)
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		replicaSet.MaintenanceWindow = getMaintenanceWindow(d)
	}

	return &replicaSet
}

// GetReplicaSetDataCreate reads the data from the tf configuration files and populates a
// create request.
func GetReplicaSetDataCreate(d *schema.ResourceData) inMemoryDB.ReplicaSetCreate {
	return inMemoryDB.ReplicaSetCreate{
		Properties: GetReplicaSetDataProperties(d),
	}
}

// GetReplicaSetDataUpdate reads the data from the tf configuration files and populates an
// update request.
func GetReplicaSetDataUpdate(d *schema.ResourceData) inMemoryDB.ReplicaSetEnsure {
	replicaStateID := d.Id()
	return inMemoryDB.ReplicaSetEnsure{
		Id:         &replicaStateID,
		Properties: GetReplicaSetDataProperties(d),
	}
}

//nolint:all
func (c *InMemoryDBClient) SetReplicaSetData(d *schema.ResourceData, replicaSet inMemoryDB.ReplicaSetRead) error {
	resourceName := "InMemoryDB replica set"
	if replicaSet.Id != nil {
		d.SetId(*replicaSet.Id)
	}

	if replicaSet.Metadata == nil {
		return fmt.Errorf("response metadata should not be empty for InMemoryDB replica set with ID: %v", *replicaSet.Id)
	}

	if replicaSet.Properties == nil {
		return fmt.Errorf("response properties should not be empty for InMemoryDB replica set with ID: %v", *replicaSet.Id)
	}

	if replicaSet.Properties.DisplayName != nil {
		if err := d.Set("display_name", *replicaSet.Properties.DisplayName); err != nil {
			return utils.GenerateSetError(resourceName, "display_name", err)
		}
	}

	if replicaSet.Properties.Version != nil {
		if err := d.Set("version", *replicaSet.Properties.Version); err != nil {
			return utils.GenerateSetError(resourceName, "version", err)
		}
	}

	if replicaSet.Properties.Replicas != nil {
		if err := d.Set("replicas", *replicaSet.Properties.Replicas); err != nil {
			return utils.GenerateSetError(resourceName, "replicas", err)
		}
	}

	if replicaSet.Properties.PersistenceMode != nil {
		if err := d.Set("persistence_mode", *replicaSet.Properties.PersistenceMode); err != nil {
			return utils.GenerateSetError(resourceName, "persistence_mode", err)
		}
	}

	if replicaSet.Properties.EvictionPolicy != nil {
		if err := d.Set("eviction_policy", *replicaSet.Properties.EvictionPolicy); err != nil {
			return utils.GenerateSetError(resourceName, "eviction_policy", err)
		}
	}

	if replicaSet.Properties.InitialSnapshotId != nil {
		if err := d.Set("initial_snapshot_id", *replicaSet.Properties.InitialSnapshotId); err != nil {
			return utils.GenerateSetError(resourceName, "initial_snapshot_id", err)
		}
	}

	if replicaSet.Metadata.DnsName != nil {
		if err := d.Set("dns_name", *replicaSet.Metadata.DnsName); err != nil {
			return utils.GenerateSetError(resourceName, "dns_name", err)
		}
	}

	if replicaSet.Properties.Resources != nil {
		var resources []interface{}
		resourceEntry := setResourceProperties(*replicaSet.Properties.Resources)
		resources = append(resources, resourceEntry)
		if err := d.Set("resources", resources); err != nil {
			return utils.GenerateSetError(resourceName, "resources", err)
		}
	}

	if replicaSet.Properties.Connections != nil {
		var connections []interface{}
		for _, connection := range *replicaSet.Properties.Connections {
			connectionEntry := setConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}
		if err := d.Set("connections", connections); err != nil {
			return utils.GenerateSetError(resourceName, "connections", err)
		}
	}

	if replicaSet.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := setMaintenanceWindowProperties(*replicaSet.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(resourceName, "maintenance_window", err)
		}
	}

	if replicaSet.Metadata != nil && replicaSet.Metadata.DnsName != nil {
		if err := d.Set("dns_name", *replicaSet.Metadata.DnsName); err != nil {
			return utils.GenerateSetError(resourceName, "dns_name", err)
		}
	}

	return nil
}

//nolint:golint
func (c *InMemoryDBClient) SetSnapshotData(d *schema.ResourceData, snapshot inMemoryDB.SnapshotRead) error {
	if snapshot.Id == nil {
		return fmt.Errorf("expected a valid ID for InMemoryDB snapshot, but got 'nil' instead")
	}
	d.SetId(*snapshot.Id)
	if snapshot.Metadata == nil {
		return fmt.Errorf("response metadata should not be empty for InMemoryDB snapshot with ID: %v", *snapshot.Id)
	}
	var metadata []interface{}
	metadataEntry := make(map[string]interface{})
	if snapshot.Metadata.CreatedDate != nil {
		metadataEntry["created_date"] = (snapshot.Metadata.CreatedDate).Time.Format(constant.DatetimeZLayout)
	}
	if snapshot.Metadata.LastModifiedDate != nil {
		metadataEntry["last_modified_date"] = (snapshot.Metadata.LastModifiedDate).Time.Format(constant.DatetimeZLayout)
	}
	if snapshot.Metadata.ReplicasetId != nil {
		metadataEntry["replica_set_id"] = *snapshot.Metadata.ReplicasetId
	}
	if snapshot.Metadata.SnapshotTime != nil {
		metadataEntry["snapshot_time"] = (snapshot.Metadata.SnapshotTime).Time.Format(constant.DatetimeZLayout)
	}
	if snapshot.Metadata.DatacenterId != nil {
		metadataEntry["datacenter_id"] = *snapshot.Metadata.DatacenterId
	}
	metadata = append(metadata, metadataEntry)
	if err := d.Set("metadata", metadata); err != nil {
		return utils.GenerateSetError(constant.DBaaSInMemoryDBSnapshotResource, "metadata", err)
	}
	return nil
}

// getResources returns information about the 'resources' attribute defined in the tf configuration
// for the ReplicaSet resource, this information will be latter used to populate the request.
func getResources(d *schema.ResourceData) *inMemoryDB.Resources {
	var resources inMemoryDB.Resources
	if cores, ok := d.GetOk("resources.0.cores"); ok {
		cores := int32(cores.(int))
		resources.Cores = &cores
	}
	if ram, ok := d.GetOk("resources.0.ram"); ok {
		ram := int32(ram.(int))
		resources.Ram = &ram
	}
	return &resources
}

// getConnections returns information about the 'connections' attribute defined in the tf configuration.
func getConnections(d *schema.ResourceData) *[]inMemoryDB.Connection {
	var connections []inMemoryDB.Connection
	var connection inMemoryDB.Connection
	if datacenterID, ok := d.GetOk("connections.0.datacenter_id"); ok {
		datacenterID := datacenterID.(string)
		connection.DatacenterId = &datacenterID
	}
	if lanID, ok := d.GetOk("connections.0.lan_id"); ok {
		lanID := lanID.(string)
		connection.LanId = &lanID
	}
	if cidr, ok := d.GetOk("connections.0.cidr"); ok {
		cidr := cidr.(string)
		connection.Cidr = &cidr
	}
	connections = append(connections, connection)
	return &connections
}

// getCredentials returns information about the 'credentials' attribute defined in the tf configuration.
func getCredentials(d *schema.ResourceData) *inMemoryDB.User {
	var user inMemoryDB.User
	var password inMemoryDB.UserPassword
	if username, ok := d.GetOk("credentials.0.username"); ok {
		username := username.(string)
		user.Username = &username
	}
	if plainTextPassword, ok := d.GetOk("credentials.0.plain_text_password"); ok {
		plainTextPassword := plainTextPassword.(string)
		password.PlainTextPassword = &plainTextPassword
	}
	if _, ok := d.GetOk("credentials.0.hashed_password"); ok {
		password.HashedPassword = getHashPasswordInfo(d)
	}
	user.Password = &password
	return &user
}

// getHashPasswordInfo returns information about the 'hashed_password' attribute defined in the tf configuration.
func getHashPasswordInfo(d *schema.ResourceData) *inMemoryDB.HashedPassword {
	var hashedPassword inMemoryDB.HashedPassword
	if algorithm, ok := d.GetOk("credentials.0.hashed_password.0.algorithm"); ok {
		algorithm := algorithm.(string)
		hashedPassword.Algorithm = &algorithm
	}
	if hash, ok := d.GetOk("credentials.0.hashed_password.0.hash"); ok {
		hash := hash.(string)
		hashedPassword.Hash = &hash
	}
	return &hashedPassword
}

// getMaintenanceWindow returns information about the 'maintenance_window' attribute defined in the tf configuration.
func getMaintenanceWindow(d *schema.ResourceData) *inMemoryDB.MaintenanceWindow {
	var maintenanceWindow inMemoryDB.MaintenanceWindow
	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := inMemoryDB.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}
	if time, ok := d.GetOk("maintenance_window.0.time"); ok {
		time := time.(string)
		maintenanceWindow.Time = &time
	}
	return &maintenanceWindow
}

func setConnectionProperties(connection inMemoryDB.Connection) map[string]interface{} {
	connectionMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "cidr", connection.Cidr)

	return connectionMap
}

func setResourceProperties(resource inMemoryDB.Resources) map[string]interface{} {
	resourceMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(resourceMap, "cores", resource.Cores)
	utils.SetPropWithNilCheck(resourceMap, "ram", resource.Ram)
	utils.SetPropWithNilCheck(resourceMap, "storage", resource.Storage)

	return resourceMap
}

func setMaintenanceWindowProperties(maintenanceWindow inMemoryDB.MaintenanceWindow) map[string]interface{} {
	maintenanceWindowMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(maintenanceWindowMap, "day_of_the_week", maintenanceWindow.DayOfTheWeek)
	utils.SetPropWithNilCheck(maintenanceWindowMap, "time", maintenanceWindow.Time)

	return maintenanceWindowMap
}
