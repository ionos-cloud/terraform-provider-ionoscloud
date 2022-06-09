package dsaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var nodePoolResourceName = "DSaaS Node Pool"

type NodePoolService interface {
	GetNodePool(ctx context.Context, clusterId, nodePoolId string) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	ListNodePools(ctx context.Context, clusterId string) ([]dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
	CreateNodePool(ctx context.Context, clusterId string, cluster dsaas.CreateNodePoolRequest) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error)
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

func (c *Client) ListNodePools(ctx context.Context, clusterId string) (dsaas.NodePoolListResponseData, *dsaas.APIResponse, error) {
	nodePool, apiResponse, err := c.DataPlatformNodePoolApi.GetClusterNodepools(ctx, clusterId).Execute()
	if apiResponse != nil {
		return nodePool, apiResponse, err
	}
	return nodePool, nil, err
}

func (c *Client) CreateNodePool(ctx context.Context, clusterId string, cluster dsaas.CreateNodePoolRequest) (dsaas.NodePoolResponseData, *dsaas.APIResponse, error) {
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

func GetDSaaSNodePoolDataCreate(d *schema.ResourceData) *dsaas.CreateNodePoolRequest {

	dsaasNodePool := dsaas.CreateNodePoolRequest{
		Properties: &dsaas.CreateNodePoolProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dsaasNodePool.Properties.Name = &name
	}

	if nodeCountValue, ok := d.GetOk("node_count"); ok {
		nodeCount := int32(nodeCountValue.(int))
		dsaasNodePool.Properties.NodeCount = &nodeCount
	}

	if cpuFamilyValue, ok := d.GetOk("cpu_family"); ok {
		cpuFamily := cpuFamilyValue.(string)
		dsaasNodePool.Properties.CpuFamily = &cpuFamily
	}

	if coresCountValue, ok := d.GetOk("cores_count"); ok {
		coresCount := int32(coresCountValue.(int))
		dsaasNodePool.Properties.CoresCount = &coresCount
	}

	if ramSizeValue, ok := d.GetOk("ram_size"); ok {
		ramSize := int32(ramSizeValue.(int))
		dsaasNodePool.Properties.RamSize = &ramSize
	}

	if availabilityZoneValue, ok := d.GetOk("availability_zone"); ok {
		availabilityZone := dsaas.AvailabilityZone(availabilityZoneValue.(string))
		dsaasNodePool.Properties.AvailabilityZone = &availabilityZone
	}

	if storageTypeValue, ok := d.GetOk("availability_zone"); ok {
		storageType := dsaas.StorageType(storageTypeValue.(string))
		dsaasNodePool.Properties.StorageType = &storageType
	}

	if storageSizeValue, ok := d.GetOk("storage_size"); ok {
		storageSize := int32(storageSizeValue.(int))
		dsaasNodePool.Properties.StorageSize = &storageSize
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dsaasNodePool.Properties.MaintenanceWindow = GetDSaaSMaintenanceWindowData(d)
	}

	if labelsValue, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, v := range labelsValue.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		dsaasNodePool.Properties.Labels = &labels
	}

	if annotationsValue, ok := d.GetOk("annotations"); ok {
		annotations := make(map[string]interface{})
		for k, v := range annotationsValue.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
		dsaasNodePool.Properties.Annotations = &annotations
	}

	return &dsaasNodePool
}

func GetDSaaSNodePoolDataUpdate(d *schema.ResourceData) (*dsaas.PatchNodePoolRequest, diag.Diagnostics) {

	dsaasNodePool := dsaas.PatchNodePoolRequest{
		Properties: &dsaas.PatchNodePoolProperties{},
	}

	if _, ok := d.GetOk("name"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "name"))

	}

	if nodeCountValue, ok := d.GetOk("node_count"); ok {
		nodeCount := int32(nodeCountValue.(int))
		dsaasNodePool.Properties.NodeCount = &nodeCount
	}

	if _, ok := d.GetOk("cpu_family"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "cpu_family"))
	}

	if _, ok := d.GetOk("cores_count"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "cores_count"))

	}

	if _, ok := d.GetOk("ram_size"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "ram_size"))
	}

	if _, ok := d.GetOk("availability_zone"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "availability_zone"))
	}

	if _, ok := d.GetOk("availability_zone"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "availability_zone"))
	}

	if _, ok := d.GetOk("storage_size"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(nodePoolResourceName, "storage_size"))
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dsaasNodePool.Properties.MaintenanceWindow = GetDSaaSMaintenanceWindowData(d)
	}

	if labelsValue, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, v := range labelsValue.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		dsaasNodePool.Properties.Labels = &labels
	}

	if annotationsValue, ok := d.GetOk("annotations"); ok {
		annotations := make(map[string]interface{})
		for k, v := range annotationsValue.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
		dsaasNodePool.Properties.Annotations = &annotations
	}

	return &dsaasNodePool, nil
}

func SetDSaaSNodePoolData(d *schema.ResourceData, nodePool dsaas.NodePoolResponseData) error {

	if nodePool.Id != nil {
		d.SetId(*nodePool.Id)
	}

	if nodePool.Properties.Name != nil {
		if err := d.Set("name", *nodePool.Properties.Name); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "name", err)
		}
	}

	if nodePool.Properties.DataPlatformVersion != nil {
		if err := d.Set("data_platform_version", *nodePool.Properties.DataPlatformVersion); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "data_platform_version", err)
		}
	}

	if nodePool.Properties.DatacenterId != nil {
		if err := d.Set("datacenter_id", *nodePool.Properties.DatacenterId); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "datacenter_id", err)
		}
	}

	if nodePool.Properties.NodeCount != nil {
		if err := d.Set("node_count", *nodePool.Properties.NodeCount); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "node_count", err)
		}
	}

	if nodePool.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *nodePool.Properties.CpuFamily); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "cpu_family", err)
		}
	}

	if nodePool.Properties.CoresCount != nil {
		if err := d.Set("cores_count", *nodePool.Properties.CoresCount); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "cores_count", err)
		}
	}

	if nodePool.Properties.RamSize != nil {
		if err := d.Set("ram_size", *nodePool.Properties.RamSize); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "ram_size", err)
		}
	}

	if nodePool.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *nodePool.Properties.AvailabilityZone); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "availability_zone", err)
		}
	}

	if nodePool.Properties.StorageType != nil {
		if err := d.Set("storage_type", *nodePool.Properties.StorageType); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "storage_type", err)
		}
	}

	if nodePool.Properties.StorageSize != nil {
		if err := d.Set("storage_size", *nodePool.Properties.StorageSize); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "storage_size", err)
		}
	}

	if nodePool.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := SetMaintenanceWindowProperties(*nodePool.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "maintenance_window", err)
		}
	}

	if nodePool.Properties.Labels != nil {
		if err := d.Set("labels", *nodePool.Properties.Labels); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "storage_size", err)
		}
	}

	if nodePool.Properties.Annotations != nil {
		if err := d.Set("annotations", *nodePool.Properties.Annotations); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "annotations", err)
		}
	}
	return nil
}
