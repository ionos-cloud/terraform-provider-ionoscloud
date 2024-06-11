package redisdb

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redisdb "github.com/ionos-cloud/sdk-go-dbaas-redis"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"log"
	"strings"
)

var locationToURL = map[string]string{
	"de/fra": "https://redis.de-fra.ionos.com",
	"de/txl": "https://redis.de-txl.ionos.com",
	"es/vit": "https://redis.es-vit.ionos.com",
	"gb/lhr": "https://redis.gb-lhr.ionos.com",
	"us/ewr": "https://redis.us-ewr.ionos.com",
	"us/las": "https://redis.us-las.ionos.com",
	"us/mci": "https://redis.us-mci.ionos.com",
	"fr/par": "https://redis.fr-par.ionos.com",
}

// modifyConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *RedisDBClient) modifyConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	clientConfig.Servers = redisdb.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

// CreateRedisDPReplicaSet sends a 'POST' request to the API to create a replica set.
func (c *RedisDBClient) CreateRedisDPReplicaSet(ctx context.Context, replicaSet redisdb.ReplicaSetCreate, location string) (redisdb.ReplicaSetRead, *redisdb.APIResponse, error) {
	c.modifyConfigURL(location)
	replicaSetResponse, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsPost(ctx).ReplicaSetCreate(replicaSet).Execute()
	apiResponse.LogInfo()
	return replicaSetResponse, apiResponse, err
}

func (c *RedisDBClient) GetReplicaSet(ctx context.Context, replicaSetID, location string) (redisdb.ReplicaSetRead, *redisdb.APIResponse, error) {
	c.modifyConfigURL(location)
	replicaSet, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsFindById(ctx, replicaSetID).Execute()
	apiResponse.LogInfo()
	return replicaSet, apiResponse, err
}

func (c *RedisDBClient) ListReplicaSets(ctx context.Context, filterName, location string) (redisdb.ReplicaSetReadList, *redisdb.APIResponse, error) {
	c.modifyConfigURL(location)
	request := c.sdkClient.ReplicaSetApi.ReplicasetsGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	replicaSets, apiResponse, err := c.sdkClient.ReplicaSetApi.ReplicasetsGetExecute(request)
	apiResponse.LogInfo()
	return replicaSets, apiResponse, err
}

func (c *RedisDBClient) IsReplicaSetReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	replicaSetID := d.Id()
	location := d.Get("location").(string)
	replicaSet, _, err := c.GetReplicaSet(ctx, replicaSetID, location)
	if err != nil {
		return false, fmt.Errorf("status check failed for Redis DB replica set with ID: %v, error: %w", replicaSetID, err)
	}
	if replicaSet.Metadata == nil || replicaSet.Metadata.State == nil {
		return false, fmt.Errorf("metadata or state is empty for Redis DB replica set with ID: %v", replicaSetID)
	}
	log.Printf("[INFO] state of the RedisDB replica set with ID: %v is: %v", replicaSetID, *replicaSet.Metadata.State)
	return strings.EqualFold(string(*replicaSet.Metadata.State), constant.Available), nil
}

// GetRedisDBReplicaSetDataCreate reads the data from the tf configuration files and populates a
// request.
func GetRedisDBReplicaSetDataCreate(d *schema.ResourceData) redisdb.ReplicaSetCreate {
	replicaSet := redisdb.ReplicaSetCreate{
		Properties: &redisdb.ReplicaSet{},
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		replicaSet.Properties.DisplayName = &displayName
	}

	if redisVersion, ok := d.GetOk("redis_version"); ok {
		redisVersion := redisVersion.(string)
		replicaSet.Properties.RedisVersion = &redisVersion
	}

	if replicas, ok := d.GetOk("replicas"); ok {
		replicas := int32(replicas.(int))
		replicaSet.Properties.Replicas = &replicas
	}

	if persistenceMode, ok := d.GetOk("persistence_mode"); ok {
		persistenceMode := redisdb.PersistenceMode(persistenceMode.(string))
		replicaSet.Properties.PersistenceMode = &persistenceMode
	}

	if evictionPolicy, ok := d.GetOk("eviction_policy"); ok {
		evictionPolicy := redisdb.EvictionPolicy(evictionPolicy.(string))
		replicaSet.Properties.EvictionPolicy = &evictionPolicy
	}

	if initialSnapshotId, ok := d.GetOk("initial_snapshot_id"); ok {
		initialSnapshotId := initialSnapshotId.(string)
		replicaSet.Properties.InitialSnapshotId = &initialSnapshotId
	}

	if _, ok := d.GetOk("resources"); ok {
		replicaSet.Properties.Resources = getResources(d)
	}

	if _, ok := d.GetOk("connections"); ok {
		replicaSet.Properties.Connections = getConnections(d)
	}

	if _, ok := d.GetOk("credentials"); ok {
		replicaSet.Properties.Credentials = getCredentials(d)
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		replicaSet.Properties.MaintenanceWindow = getMaintenanceWindow(d)
	}

	return replicaSet
}

func (c *RedisDBClient) SetRedisDBReplicaSetData(d *schema.ResourceData, replicaSet redisdb.ReplicaSetRead) error {
	resourceName := "RedisDB replica set"
	if replicaSet.Id != nil {
		d.SetId(*replicaSet.Id)
	}
	if replicaSet.Properties == nil {
		return fmt.Errorf("response properties should not be empty for RedisDB replica set with ID: %v", *replicaSet.Id)
	}
	if replicaSet.Properties.DisplayName != nil {
		if err := d.Set("display_name", *replicaSet.Properties.DisplayName); err != nil {
			return utils.GenerateSetError(resourceName, "display_name", err)
		}
	}

	if replicaSet.Properties.RedisVersion != nil {
		if err := d.Set("redis_version", *replicaSet.Properties.RedisVersion); err != nil {
			return utils.GenerateSetError(resourceName, "redis_version", err)
		}
	}

	if replicaSet.Properties.Replicas != nil {
		if err := d.Set("replicas", *replicaSet.Properties.Replicas); err != nil {
			return utils.GenerateSetError(resourceName, "replicas", err)
		}
	}

	// TODO -- Check if we need to set anything for 'credentials'

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

	return nil
}

// getResources returns information about the 'resources' attribute defined in the tf configuration
// for the ReplicaSet resource, this information will be latter used to populate the request.
func getResources(d *schema.ResourceData) *redisdb.Resources {
	var resources redisdb.Resources
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
func getConnections(d *schema.ResourceData) *[]redisdb.Connection {
	var connections []redisdb.Connection
	var connection redisdb.Connection
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
func getCredentials(d *schema.ResourceData) *redisdb.User {
	var user redisdb.User
	var password redisdb.UserPassword
	if username, ok := d.GetOk("credentials.0.username"); ok {
		username := username.(string)
		user.Username = &username
	}
	if plainTextPassword, ok := d.GetOk("credentials.0.plain_text_password"); ok {
		plainTextPassword := plainTextPassword.(string)
		password.PlaintextPassword = &plainTextPassword
	}
	if _, ok := d.GetOk("credentials.0.hashed_password"); ok {
		password.HashedPassword = getHashPasswordInfo(d)
	}
	user.Password = &password
	return &user
}

// getHashPasswordInfo returns information about the 'hashed_password' attribute defined in the tf configuration.
func getHashPasswordInfo(d *schema.ResourceData) *redisdb.HashedPassword {
	var hashedPassword redisdb.HashedPassword
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
func getMaintenanceWindow(d *schema.ResourceData) *redisdb.MaintenanceWindow {
	var maintenanceWindow redisdb.MaintenanceWindow
	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := redisdb.DayOfTheWeek(dayOfTheWeek.(string))
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}
	if time, ok := d.GetOk("maintenance_window.0.time"); ok {
		time := time.(string)
		maintenanceWindow.Time = &time
	}
	return &maintenanceWindow
}

func setConnectionProperties(connection redisdb.Connection) map[string]interface{} {
	connectionMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "cidr", connection.Cidr)

	return connectionMap
}

func setResourceProperties(resource redisdb.Resources) map[string]interface{} {
	resourceMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(resourceMap, "cores", resource.Cores)
	utils.SetPropWithNilCheck(resourceMap, "ram", resource.Ram)
	utils.SetPropWithNilCheck(resourceMap, "storage", resource.Storage)

	return resourceMap
}

func setMaintenanceWindowProperties(maintenanceWindow redisdb.MaintenanceWindow) map[string]interface{} {
	maintenanceWindowMap := make(map[string]interface{})

	utils.SetPropWithNilCheck(maintenanceWindowMap, "day_of_the_week", maintenanceWindow.DayOfTheWeek)
	utils.SetPropWithNilCheck(maintenanceWindowMap, "time", maintenanceWindow.Time)

	return maintenanceWindowMap
}
