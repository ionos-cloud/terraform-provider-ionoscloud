package inmemorydb

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// overrideClientEndpoint todo - after move to bundle, replace with generic function from fileConfig
func (c *Client) overrideClientEndpoint(productName, location string) {
	// whatever is set, at the end we need to check if the IONOS_API_URL_productname is set and use override the endpoint if yes
	defer c.changeConfigURL(location)
	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		fmt.Printf("[DEBUG] Using custom endpoint %s\n", os.Getenv(shared.IonosApiUrlEnvVar))
		return
	}
	fileConfig := c.GetFileConfig()
	if fileConfig == nil {
		return
	}
	config := c.GetConfig()
	if config == nil {
		return
	}
	endpoint := fileConfig.GetProductLocationOverrides(productName, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s in location %s", productName, location)
		return
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
}

// CreateReplicaSet sends a 'POST' request to the API to create a replica set.
func (c *Client) CreateReplicaSet(ctx context.Context, replicaSet inmemorydb.ReplicaSetCreate, location string) (inmemorydb.ReplicaSetRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
	replicaSetResponse, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsPost(ctx).ReplicaSetCreate(replicaSet).Execute()
	apiResponse.LogInfo()
	return replicaSetResponse, apiResponse, err
}

// IsReplicaSetReady checks if the replica set is in the 'AVAILABLE' state.
func (c *Client) IsReplicaSetReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	replicaSetID := d.Id()
	location := d.Get("location").(string)
	replicaSet, _, err := c.GetReplicaSet(ctx, replicaSetID, location)
	if err != nil {
		return false, fmt.Errorf("status check failed for InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err)
	}

	log.Printf("[INFO] state of the InMemoryDB replica set with ID: %v is: %v", replicaSetID, replicaSet.Metadata.State)
	if utils.IsStateFailed(replicaSet.Metadata.State) {
		return false, fmt.Errorf("replica set with ID: %v is in FAILED state", replicaSetID)
	}

	return strings.EqualFold(replicaSet.Metadata.State, constant.Available), nil
}

// DeleteReplicaSet sends a 'DELETE' request to the API to delete a replica set.
func (c *Client) DeleteReplicaSet(ctx context.Context, replicaSetID, location string) (*shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
	apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsDelete(ctx, replicaSetID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateReplicaSet sends a 'PUT' request to the API to update a replica set.
func (c *Client) UpdateReplicaSet(ctx context.Context, replicaSetID, location string, replicaSet inmemorydb.ReplicaSetEnsure) (inmemorydb.ReplicaSetRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
	replicaSetResponse, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsPut(ctx, replicaSetID).ReplicaSetEnsure(replicaSet).Execute()
	apiResponse.LogInfo()
	return replicaSetResponse, apiResponse, err
}

// IsReplicaSetDeleted checks if the replica set is deleted.
func (c *Client) IsReplicaSetDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
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

// GetReplicaSet sends a 'GET' request to the API to get a replica set.
func (c *Client) GetReplicaSet(ctx context.Context, replicaSetID, location string) (inmemorydb.ReplicaSetRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
	replicaSet, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsFindById(ctx, replicaSetID).Execute()
	apiResponse.LogInfo()
	return replicaSet, apiResponse, err
}

// GetSnapshot sends a 'GET' request to the API to get a snapshot.
func (c *Client) GetSnapshot(ctx context.Context, snapshotID, location string) (inmemorydb.SnapshotRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
	snapshot, apiResponse, err := c.sdkClient.SnapshotApi.SnapshotsFindById(ctx, snapshotID).Execute()
	apiResponse.LogInfo()
	return snapshot, apiResponse, err
}

// ListReplicaSets sends a 'GET' request to the API to get a list of replica sets.
func (c *Client) ListReplicaSets(ctx context.Context, filterName, location string) (inmemorydb.ReplicaSetReadList, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.InMemoryDB, location)
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
func GetReplicaSetDataProperties(d *schema.ResourceData) inmemorydb.ReplicaSet {
	replicaSet := inmemorydb.ReplicaSet{}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		replicaSet.DisplayName = displayName
	}

	if inMemoryDBVersion, ok := d.GetOk("version"); ok {
		inMemoryDBVersion := inMemoryDBVersion.(string)
		replicaSet.Version = inMemoryDBVersion
	}

	if replicas, ok := d.GetOk("replicas"); ok {
		replicas := int32(replicas.(int))
		replicaSet.Replicas = replicas
	}

	if persistenceMode, ok := d.GetOk("persistence_mode"); ok {
		persistenceMode := inmemorydb.PersistenceMode(persistenceMode.(string))
		replicaSet.PersistenceMode = persistenceMode
	}

	if evictionPolicy, ok := d.GetOk("eviction_policy"); ok {
		evictionPolicy := inmemorydb.EvictionPolicy(evictionPolicy.(string))
		replicaSet.EvictionPolicy = evictionPolicy
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

	return replicaSet
}

// GetReplicaSetDataCreate reads the data from the tf configuration files and populates a
// create request.
func GetReplicaSetDataCreate(d *schema.ResourceData) inmemorydb.ReplicaSetCreate {
	return inmemorydb.ReplicaSetCreate{
		Properties: GetReplicaSetDataProperties(d),
	}
}

// GetReplicaSetDataUpdate reads the data from the tf configuration files and populates an
// update request.
func GetReplicaSetDataUpdate(d *schema.ResourceData) inmemorydb.ReplicaSetEnsure {
	replicaStateID := d.Id()
	return inmemorydb.ReplicaSetEnsure{
		Id:         replicaStateID,
		Properties: GetReplicaSetDataProperties(d),
	}
}

// SetReplicaSetData populates the tf resource data with the response from the API.
func (c *Client) SetReplicaSetData(d *schema.ResourceData, replicaSet inmemorydb.ReplicaSetRead) error {
	resourceName := "InMemoryDB replica set"

	d.SetId(replicaSet.Id)

	if err := d.Set("display_name", replicaSet.Properties.DisplayName); err != nil {
		return utils.GenerateSetError(resourceName, "display_name", err)
	}

	if err := d.Set("version", replicaSet.Properties.Version); err != nil {
		return utils.GenerateSetError(resourceName, "version", err)
	}

	if err := d.Set("replicas", replicaSet.Properties.Replicas); err != nil {
		return utils.GenerateSetError(resourceName, "replicas", err)
	}

	if err := d.Set("persistence_mode", replicaSet.Properties.PersistenceMode); err != nil {
		return utils.GenerateSetError(resourceName, "persistence_mode", err)
	}

	if err := d.Set("eviction_policy", replicaSet.Properties.EvictionPolicy); err != nil {
		return utils.GenerateSetError(resourceName, "eviction_policy", err)
	}

	if replicaSet.Properties.InitialSnapshotId != nil {
		if err := d.Set("initial_snapshot_id", *replicaSet.Properties.InitialSnapshotId); err != nil {
			return utils.GenerateSetError(resourceName, "initial_snapshot_id", err)
		}
	}

	if err := d.Set("dns_name", replicaSet.Metadata.DnsName); err != nil {
		return utils.GenerateSetError(resourceName, "dns_name", err)
	}

	var resources []interface{}
	resourceEntry := setResourceProperties(replicaSet.Properties.Resources)
	resources = append(resources, resourceEntry)
	if err := d.Set("resources", resources); err != nil {
		return utils.GenerateSetError(resourceName, "resources", err)
	}

	if replicaSet.Properties.Connections != nil {
		var connections []interface{}
		for _, connection := range replicaSet.Properties.Connections {
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

	if err := d.Set("dns_name", replicaSet.Metadata.DnsName); err != nil {
		return utils.GenerateSetError(resourceName, "dns_name", err)
	}

	return nil
}

// SetSnapshotData populates the tf resource data with the response from the API.
func (c *Client) SetSnapshotData(d *schema.ResourceData, snapshot inmemorydb.SnapshotRead) error {

	d.SetId(snapshot.Id)

	var metadata []interface{}
	metadataEntry := make(map[string]interface{})
	if snapshot.Metadata.CreatedDate != nil {
		metadataEntry["created_date"] = (snapshot.Metadata.CreatedDate).Time.Format(constant.DatetimeZLayout)
	}
	if snapshot.Metadata.LastModifiedDate != nil {
		metadataEntry["last_modified_date"] = (snapshot.Metadata.LastModifiedDate).Time.Format(constant.DatetimeZLayout)
	}

	metadataEntry["replica_set_id"] = snapshot.Metadata.ReplicasetId

	if snapshot.Metadata.SnapshotTime != nil {
		metadataEntry["snapshot_time"] = (snapshot.Metadata.SnapshotTime).Time.Format(constant.DatetimeZLayout)
	}

	metadataEntry["datacenter_id"] = snapshot.Metadata.DatacenterId

	metadata = append(metadata, metadataEntry)
	if err := d.Set("metadata", metadata); err != nil {
		return utils.GenerateSetError(constant.DBaaSInMemoryDBSnapshotResource, "metadata", err)
	}
	return nil
}

// getResources returns information about the 'resources' attribute defined in the tf configuration
// for the ReplicaSet resource, this information will be latter used to populate the request.
func getResources(d *schema.ResourceData) inmemorydb.Resources {
	var resources inmemorydb.Resources
	if cores, ok := d.GetOk("resources.0.cores"); ok {
		cores := int32(cores.(int))
		resources.Cores = cores
	}
	if ram, ok := d.GetOk("resources.0.ram"); ok {
		ram := int32(ram.(int))
		resources.Ram = ram
	}
	return resources
}

// getConnections returns information about the 'connections' attribute defined in the tf configuration.
func getConnections(d *schema.ResourceData) []inmemorydb.Connection {
	var connections []inmemorydb.Connection
	var connection inmemorydb.Connection
	if datacenterID, ok := d.GetOk("connections.0.datacenter_id"); ok {
		datacenterID := datacenterID.(string)
		connection.DatacenterId = datacenterID
	}
	if lanID, ok := d.GetOk("connections.0.lan_id"); ok {
		lanID := lanID.(string)
		connection.LanId = lanID
	}
	if cidr, ok := d.GetOk("connections.0.cidr"); ok {
		cidr := cidr.(string)
		connection.Cidr = cidr
	}
	connections = append(connections, connection)
	return connections
}

// getCredentials returns information about the 'credentials' attribute defined in the tf configuration.
func getCredentials(d *schema.ResourceData) inmemorydb.User {
	var user inmemorydb.User
	var password inmemorydb.UserPassword
	if username, ok := d.GetOk("credentials.0.username"); ok {
		username := username.(string)
		user.Username = username
	}
	if plainTextPassword, ok := d.GetOk("credentials.0.plain_text_password"); ok {
		plainTextPassword := plainTextPassword.(string)
		password.PlainTextPassword = &plainTextPassword
	}
	if _, ok := d.GetOk("credentials.0.hashed_password"); ok {
		password.HashedPassword = getHashPasswordInfo(d)
	}
	user.Password = &password
	return user
}

// getHashPasswordInfo returns information about the 'hashed_password' attribute defined in the tf configuration.
func getHashPasswordInfo(d *schema.ResourceData) *inmemorydb.HashedPassword {
	var hashedPassword inmemorydb.HashedPassword
	if algorithm, ok := d.GetOk("credentials.0.hashed_password.0.algorithm"); ok {
		algorithm := algorithm.(string)
		hashedPassword.Algorithm = algorithm
	}
	if hash, ok := d.GetOk("credentials.0.hashed_password.0.hash"); ok {
		hash := hash.(string)
		hashedPassword.Hash = hash
	}
	return &hashedPassword
}

// getMaintenanceWindow returns information about the 'maintenance_window' attribute defined in the tf configuration.
func getMaintenanceWindow(d *schema.ResourceData) *inmemorydb.MaintenanceWindow {
	var maintenanceWindow inmemorydb.MaintenanceWindow
	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := inmemorydb.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = dayOfTheWeek
	}
	if time, ok := d.GetOk("maintenance_window.0.time"); ok {
		time := time.(string)
		maintenanceWindow.Time = time
	}
	return &maintenanceWindow
}

func setConnectionProperties(connection inmemorydb.Connection) map[string]interface{} {
	connectionMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "cidr", connection.Cidr)

	return connectionMap
}

func setResourceProperties(resource inmemorydb.Resources) map[string]interface{} {
	resourceMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(resourceMap, "cores", resource.Cores)
	utils.SetPropWithNilCheck(resourceMap, "ram", resource.Ram)
	utils.SetPropWithNilCheck(resourceMap, "storage", resource.Storage)

	return resourceMap
}

func setMaintenanceWindowProperties(maintenanceWindow inmemorydb.MaintenanceWindow) map[string]interface{} {
	maintenanceWindowMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(maintenanceWindowMap, "day_of_the_week", maintenanceWindow.DayOfTheWeek)
	utils.SetPropWithNilCheck(maintenanceWindowMap, "time", maintenanceWindow.Time)

	return maintenanceWindowMap
}
