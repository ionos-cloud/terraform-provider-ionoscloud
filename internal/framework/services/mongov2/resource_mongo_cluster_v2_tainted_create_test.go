package mongov2_test

import (
	"context"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/stretchr/testify/assert"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/frameworktest"
)

// TestMongoClusterV2Create_PostSucceedsWaitFails verifies that when the POST succeeds but the cluster
// never becomes available, Create returns an error AND still writes the cluster id/location to state,
// so Terraform tracks and taints the cluster instead of orphaning it in the API.
func TestMongoClusterV2Create_PostSucceedsWaitFails(t *testing.T) {
	ctx := context.Background()

	const (
		typeName  = "ionoscloud_mongo_cluster_v2"
		location  = "de/fra"
		clusterID = "11111111-1111-1111-1111-111111111111"
	)

	frameworktest.FailedProvisioningStub{Product: fileconfiguration.MongoDBV2, Location: location, ClusterID: clusterID}.Serve(t)
	server, schema := frameworktest.NewConfiguredServer(ctx, t)
	resourceType := frameworktest.ResourceType(t, schema, typeName)

	// ENTERPRISE edition: instances block is required, instance_template_id/instance_count must be absent.
	plannedJSON := `{
		"name": "test",
		"version": "7.0",
		"location": "` + location + `",
		"edition": "ENTERPRISE",
		"connections": {"datacenter_id": "dc-id", "lan_id": "1", "instance_addresses": ["10.0.0.1"]},
		"credentials": {"username": "admin", "database": "testdb", "password": "Secret12345!"},
		"instances": {"cores": 1, "count": 1, "ram": 4, "storage_size": 10},
		"maintenance_window": {"day_of_the_week": "Sunday", "time": "09:00:00"},
		"snapshot": {"location": "eu-central-2", "retention_days": 7},
		"timeouts": {"create": "1s"}
	}`

	resp := frameworktest.Create(ctx, t, server, typeName, resourceType, plannedJSON)

	assert.Truef(t, frameworktest.HasErrorDiag(resp.Diagnostics, "become available"),
		"expected an error about the cluster not becoming available, got: %s", frameworktest.FormatDiags(resp.Diagnostics))

	state := frameworktest.StateStrings(t, resp.NewState, resourceType, "id", "location")
	assert.Equal(t, clusterID, state["id"], "id must be persisted to state after a successful POST")
	assert.Equal(t, location, state["location"], "location must be persisted to state after a successful POST")
}
