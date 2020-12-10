package profitbricks

import (
	"fmt"
	"net/http"
	"time"
)

const (
	// Kubernetes cluster/nodepool resource state is deploying
	K8sStateDeploying = "DEPLOYING"
	// Kubernetes cluster/nodepool resource state is active
	K8sStateActive = "ACTIVE"
	// Kubernetes cluster/nodepool resource state is failed
	K8sStateFailed = "FAILED"
	// Kubernetes cluster/nodepool resource state is updating
	K8sStateUpdating = "UPDATING"
	// Kubernetes cluster/nodepool resource state is failed_updating
	K8sStateFailedUpdating = "FAILED_UPDATING"
	// Kubernetes cluster/nodepool resource state is destroying
	K8sStateDestroying = "DESTROYING"
	// Kubernetes cluster/nodepool resource state is failed_destroying
	K8sStateFailedDestroying = "FAILED_DESTROYING"
	// Kubernetes cluster/nodepool resource state is terminated
	K8sStateTerminated = "TERMINATED"
)

const (
	// Kubernetes Node resource state is ready
	K8sNodeStateReady = "READY"
	// Kubernetes Node resource state is provisioning
	K8sNodeStateProvisioning = "PROVISIONING"
	// Kubernetes Node resource state is provisioned
	K8sNodeStateProvisioned = "PROVISIONED"
	// Kubernetes Node resource state is terminating
	K8sNodeStateTerminating = "TERMINATING"
	// Kubernetes Node resource state is rebuilding
	K8sNodeStateRebuilding = "REBUILDING"
)

type KubernetesClusters struct {
	// URL to the collection representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// Unique representation for Kubernetes Cluster as a collection on a resource.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// Slice of items in that collection
	// Read Only: true
	Items []KubernetesCluster `json:"items"`

	// The type of resource within a collection
	// Read Only: true
	// Enum: [collection]
	PBType string `json:"type,omitempty"`
}

type KubernetesCluster struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// metadata
	Metadata *Metadata `json:"metadata,omitempty"`

	// properties
	// Required: true
	Properties *KubernetesClusterProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [k8s]
	PBType string `json:"type,omitempty"`

	// Entities of a cluster
	Entities KubernetesClusterEntities `json:"entities,omitempty"`
}

type UpdatedKubernetesCluster struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// metadata
	Metadata *Metadata `json:"metadata,omitempty"`

	// properties
	// Required: true
	Properties *KubernetesClusterProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [k8s]
	PBType string `json:"type,omitempty"`

	// Entities of a cluster
	Entities KubernetesClusterEntities `json:"-"`
}

type KubernetesClusterEntities struct {
	// NodePools of a cluster
	NodePools *KubernetesNodePools `json:"nodepools,omitempty"`
}

type AutoScaling struct {
	// The minimum number of nodes this node pool can be scaled down to
	// Required: true
	MinNodeCount uint32 `json:"minNodeCount,omitempty"`
	// The maximum number of nodes this node pool can be scaled up to
	// Required: true
	MaxNodeCount uint32 `json:"maxNodeCount,omitempty"`
}

type KubernetesNodePoolLAN struct {
	// The ID of the lan the Kubernetes Node Pool should be part of
	// Required: true
	ID uint32 `json:"id"`
}

type MaintenanceWindow struct {
	// The english name of the day of the week
	// Required: false
	DayOfTheWeek string `json:"dayOfTheWeek,omitempty"`
	// A string of the following format: 08:00:00
	// Required: false
	Time string `json:"time,omitempty"`
}

type KubernetesClusterProperties struct {
	// A Kubernetes Cluster Name. Valid Kubernetes Cluster name must be 63 characters or less and must not be empty
	// and begin and end with an alphanumeric character ([a-z0-9]) with dashes (-), dots (.) and alphanumerics
	// between.
	// Required: true
	Name string `json:"name"`
	// The desired Kubernetes Version
	// Please consult the API documentation for supported versions
	// Required: false
	K8sVersion string `json:"k8sVersion,omitempty"`
	// The desired Maintanance Window
	// Required: false
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
}

type KubernetesConfig struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// properties
	// Required: true
	Properties KubernetesConfigProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [kubeconfig]
	PBType string `json:"type,omitempty"`
}

type KubernetesConfigProperties struct {
	// A Kubernetes Config file data
	KubeConfig string `json:"kubeconfig,omitempty"`
}

type KubernetesNodePool struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// metadata
	Metadata *Metadata `json:"metadata,omitempty"`

	// properties
	// Required: true
	Properties *KubernetesNodePoolProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [nodepool]
	PBType string `json:"type,omitempty"`
}

type KubernetesNodePoolProperties struct {
	// The availability zone in which the servers should exist
	// Required: true
	// Enum: [AUTO ZONE_1 ZONE_2]
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// Number of cores for node
	// Required: true
	CoresCount uint32 `json:"coresCount,omitempty"`

	// A valid cpu family name
	// Required: true
	CPUFamily string `json:"cpuFamily,omitempty"`

	// The unique identifier of the data center where the worker nodes of the node pool will be provisioned.
	// Required: true
	DatacenterID string `json:"datacenterId,omitempty"`

	// A Kubernetes Node Pool Name. Valid Kubernetes Node Pool name must be 63 characters or less and must not be
	// empty or begin and end with an alphanumeric character ([a-z0-9]) with dashes (-), dots (.) and alphanumerics
	// between.
	// Required: true
	Name string `json:"name,omitempty"`

	// Number of nodes part of the Node Pool
	// Required: true
	NodeCount uint32 `json:"nodeCount,omitempty"`

	// RAM size for node, minimum size 2048MB is recommended
	// Required: true
	RAMSize uint32 `json:"ramSize,omitempty"`

	// The size of the volume in GB. The size should be greater than 10GB.
	// Required: true
	StorageSize uint32 `json:"storageSize,omitempty"`

	// Hardware type of the volume
	// Required: true
	// Enum: [HDD SSD]
	StorageType string `json:"storageType,omitempty"`

	// The kubernetes version in which a nodepool is running.
	// Required: true
	K8sVersion string `json:"k8sVersion,omitempty"`

	// Whether this cluster should autoscale. This is comprised of a minimum and maximum number of nodes
	// Required: false
	AutoScaling *AutoScaling `json:"autoScaling,omitempty"`

	// The LANs the nodes in this pool should be part of
	// Required: true
	LANs *[]KubernetesNodePoolLAN `json:"lans,omitempty"`

	// The desired Maintanance Window
	// Required: false
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// Public Ips
	// Required: false
	PublicIPs *[]string `json:"publicIps,omitempty"`
}

type KubernetesNodePools struct {
	// URL to the collection representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// Unique representation for Kubernetes Nodes as a collection on a resource.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// Slice of items in that collection
	// Read Only: true
	Items []KubernetesNodePool `json:"items"`

	// The type of resource within a collection
	// Read Only: true
	// Enum: [nodepool]
	Type string `json:"type,omitempty"`
}

type KubernetesNodes struct {
	// URL to the collection representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// Unique representation for Kubernetes Node Pool as a collection on a resource.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// Slice of items in that collection
	// Read Only: true
	Items []KubernetesNode `json:"items"`

	// The type of resource within a collection
	// Read Only: true
	// Enum: [nodepool]
	Type string `json:"type,omitempty"`
}

type KubernetesNode struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// metadata
	Metadata *Metadata `json:"metadata,omitempty"`

	// The properties of the node
	Properties *KubernetesNodeProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [nodepool]
	PBType string `json:"type,omitempty"`
}

type KubernetesNodeProperties struct {
	// The generated unique name of the node.
	// Read Only: true
	Name string `json:"name,omitempty"`

	// The assigned public IP of the node.
	// Read Only: true
	PublicIP string `json:"publicIP,omitempty"`

	// The k8s version that the node has.
	// Read Only: false
	K8sVersion string `json:"k8sVersion,omitempty"`
}

// ListKubernetesClusters gets a list of all clusters
func (c *Client) ListKubernetesClusters() (*KubernetesClusters, error) {
	rsp := &KubernetesClusters{}
	return rsp, c.GetOK(kubernetesClustersPath(), rsp)
}

// GetKubernetesCluster gets cluster with given id
func (c *Client) GetKubernetesCluster(clusterID string) (*KubernetesCluster, error) {
	rsp := &KubernetesCluster{}
	return rsp, c.GetOK(kubernetesClusterPath(clusterID), rsp)
}

// CreateKubernetesCluster creates a cluster
func (c *Client) CreateKubernetesCluster(cluster KubernetesCluster) (*KubernetesCluster, error) {
	rsp := &KubernetesCluster{}
	return rsp, c.PostAcc(kubernetesClustersPath(), cluster, rsp)
}

// DeleteKubernetesCluster deletes cluster
func (c *Client) DeleteKubernetesCluster(clusterID string) (*http.Header, error) {
	h := &http.Header{}
	return h, c.Delete(kubernetesClusterPath(clusterID), h, http.StatusAccepted)
}

// UpdateKubernetesCluster updates cluster
func (c *Client) UpdateKubernetesCluster(clusterID string, cluster UpdatedKubernetesCluster) (*KubernetesCluster, error) {
	rsp := &KubernetesCluster{}
	return rsp, c.Put(kubernetesClusterPath(clusterID), cluster, rsp, http.StatusOK)
}

// GetKubeconfig returns the kubeconfig of cluster
func (c *Client) GetKubeconfig(clusterID string) (string, error) {
	rsp := &KubernetesConfig{}
	if err := c.GetOK(kubeConfigPath(clusterID), rsp); err != nil {
		return "", err
	}
	return rsp.Properties.KubeConfig, nil
}

// ListKubernetesNodePools gets a list of all node pools of a cluster
func (c *Client) ListKubernetesNodePools(clusterID string) (*KubernetesNodePools, error) {
	rsp := &KubernetesNodePools{}
	return rsp, c.GetOK(kubernetesNodePoolsPath(clusterID), rsp)
}

// CreateKubernetesNodePool creates a new node pool for cluster
func (c *Client) CreateKubernetesNodePool(clusterID string, nodePool KubernetesNodePool) (*KubernetesNodePool, error) {
	rsp := &KubernetesNodePool{}
	return rsp, c.PostAcc(kubernetesNodePoolsPath(clusterID), nodePool, rsp)
}

// DeleteKubernetesNodePool deletes node pool from cluster
func (c *Client) DeleteKubernetesNodePool(clusterID, nodePoolID string) (*http.Header, error) {
	return c.DeleteAcc(kubernetesNodePoolPath(clusterID, nodePoolID))
}

// GetKubernetesNodePool gets node pool of the cluster
func (c *Client) GetKubernetesNodePool(clusterID, nodePoolID string) (*KubernetesNodePool, error) {
	rsp := &KubernetesNodePool{}
	return rsp, c.GetOK(kubernetesNodePoolPath(clusterID, nodePoolID), rsp)
}

// Update KubernetesNodePool updates node pool
func (c *Client) UpdateKubernetesNodePool(clusterID, nodePoolID string, nodePool KubernetesNodePool) (*KubernetesNodePool, error) {
	rsp := &KubernetesNodePool{}
	return rsp, c.PutAcc(kubernetesNodePoolPath(clusterID, nodePoolID), nodePool, rsp)
}

// ListKubernetesNodes gets a list of all nodes of a node pool
func (c *Client) ListKubernetesNodes(clusterID, nodePoolID string) (*KubernetesNodes, error) {
	rsp := &KubernetesNodes{}
	return rsp, c.GetOK(kubernetesNodesPath(clusterID, nodePoolID), rsp)
}

// GetKubernetesNode gets node of a node pool
func (c *Client) GetKubernetesNode(clusterID, nodePoolID, nodeID string) (*KubernetesNode, error) {
	rsp := &KubernetesNode{}
	return rsp, c.GetOK(kubernetesNodePath(clusterID, nodePoolID, nodeID), rsp)
}

// DeleteKubernetesNode deletes a node from a node pool, decreasing its size by 1.
func (c *Client) DeleteKubernetesNode(clusterID, nodePoolID, nodeID string) (*http.Header, error) {
	return c.DeleteAcc(kubernetesNodePath(clusterID, nodePoolID, nodeID))
}

// ReplaceKubernetesNode replaces a node of a node pool.
func (c *Client) ReplaceKubernetesNode(clusterID, nodePoolID, nodeID string) (*http.Header, error) {
	url := kubernetesNodeReplacePath(clusterID, nodePoolID, nodeID)
	rsp, err := c.R().SetError(ApiError{}).Post(url)
	if err != nil {
		return nil, NewClientError(HttpClientError, fmt.Sprintf("[POST] %s: Client error: %s", url, err))
	}
	h := rsp.Header()
	return &h, validateResponse(rsp, http.StatusAccepted)
}

// Enabled returns true when max > 0.
func (a *AutoScaling) Enabled() bool {
	return a != nil && a.MaxNodeCount > 0
}

func (c *Client) WaitForKubernetesNodePoolState(
	clusterID, nodePoolID string,
	state string,
	timeout, interval time.Duration) error {
	if err := PollImmediate(interval, timeout, func() (bool, error) {
		np, err := c.GetKubernetesNodePool(clusterID, nodePoolID)
		if err != nil {
			return false, err
		}
		return np != nil && np.Metadata != nil && np.Metadata.State == state, err
	}); err != nil {
		return fmt.Errorf("error waiting for nodepool state %s: %w", state, err)
	}
	return nil
}

func (c *Client) WaitForKubernetesClusterState(
	clusterID string, state string, timeout, interval time.Duration) error {
	if err := PollImmediate(interval, timeout, func() (bool, error) {
		cl, err := c.GetKubernetesCluster(clusterID)
		if err != nil {
			return false, err
		}
		return cl != nil && cl.Metadata != nil && cl.Metadata.State == state, err
	}); err != nil {
		return fmt.Errorf("error waiting for cluster state %s: %w", state, err)
	}
	return nil
}
