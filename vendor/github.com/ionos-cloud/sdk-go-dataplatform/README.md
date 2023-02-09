# Go API client for ionoscloud

Managed Stackable Data Platform by IONOS Cloud provides a preconfigured Kubernetes cluster
with pre-installed and managed Stackable operators. After the provision of these Stackable operators,
the customer can interact with them directly
and build his desired application on top of the Stackable Platform.

Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API
in addition or as an alternative to the \"Data Center Designer\" (DCD).

## Getting Started

To get your DataPlatformCluster up and running, the following steps needs to be performed.

### IONOS Cloud Account

The first step is the creation of a IONOS Cloud account if not already existing.

To register a **new account** visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).

### Virtual Datacenter (VDC)

The Managed Data Stack needs a virtual datacenter (VDC) hosting the cluster.
This could either be a VDC that already exists, especially if you want to connect the managed DataPlatform
to other services already running within your VDC. Otherwise, if you want to place the Managed Data Stack in
a new VDC or you have not yet created a VDC, you need to do so.

A new VDC can be created via the IONOS Cloud API, the IONOS-CLI or the DCD Web interface.
For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/tutorials/data-center-basics)

### Get a authentication token

To interact with this API a user specific authentication token is needed.
This token can be generated using the IONOS-CLI the following way:

```
ionosctl token generate
```

For more information [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate)

### Create a new DataPlatformCluster

Before using Managed Stackable Data Platform, a new DataPlatformCluster must be created.

To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.

The provisioning of the cluster might take some time. To check the current provisioning status,
you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID
that was presented to you in the response of the create cluster call.

### Add a DataPlatformNodePool

To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool
that is provisioned along with the cluster is reserved for the Stackable operators.
You may create further node pools with resources tailored to your use-case.

To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post)
endpoint.

### Receive Kubeconfig

Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API.
The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.

The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint
using the cluster ID of the created DataPlatformCluster.

### Create Stackable Service

To create the desired application, the Stackable service needs to be provided,
using the received kubeconfig and
[deploy a Stackable service](https://docs.stackable.tech/home/getting_started.html#_deploying_stackable_services)

## Authorization

All endpoints are secured, so only an authenticated user can access them.
As Authentication mechanism the default IONOS Cloud authentication mechanism
is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).

### Basic-Auth

The basic auth scheme uses the IONOS Cloud user credentials in form of a Basic Authentication Header
accordingly to [RFC7617](https://datatracker.ietf.org/doc/html/rfc7617)

### API-Key as Bearer Token

The Bearer auth token used at the API-Gateway is a user related token created with the IONOS-CLI.
(See the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate) for details)
For every request to be authenticated, the token is passed as 'Authorization Bearer' header along with the request.

### Permissions and access roles

Currently, an admin can see and manipulate all resources in a contract.
A normal authenticated user can only see and manipulate resources he created.


## Components

The Managed Stackable Data Platform by IONOS Cloud consists of two components.
The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.

### DataPlatformCluster

A DataPlatformCluster is the virtual instance of the customer services and operations running the managed Services like Stackable operators.
A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer.
Therefore, it's possible to integrate the cluster with other resources as vLANs e.G.
to shape the datacenter in the customer's need and integrate the cluster within the topology the customer wants to build.

In addition to the Kubernetes cluster a small node pool is provided which is exclusively used to run the Stackable operators.

### DataPlatformNodePool

A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top.
All nodes within a node pool are identical in setup.
The nodes of a pool are provisioned into virtual data centers at a location of your choice
and you can freely specify the properties of all the nodes at once before creation.

Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are readonly in the customer's VDC
and can only be modified or deleted via the API.

### References


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.0.7
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [https://ionos.com](https://ionos.com)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import sw "./ionoscloud"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identifield by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```
ctx := context.WithValue(context.Background(), sw.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), sw.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/dataplatform*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DataPlatformClusterApi* | [**CreateCluster**](docs/api/DataPlatformClusterApi.md#createcluster) | **Post** /clusters | Create a DataPlatformCluster
*DataPlatformClusterApi* | [**DeleteCluster**](docs/api/DataPlatformClusterApi.md#deletecluster) | **Delete** /clusters/{clusterId} | Delete DataPlatformCluster
*DataPlatformClusterApi* | [**GetCluster**](docs/api/DataPlatformClusterApi.md#getcluster) | **Get** /clusters/{clusterId} | Retrieve a DataPlatformCluster
*DataPlatformClusterApi* | [**GetClusterKubeconfig**](docs/api/DataPlatformClusterApi.md#getclusterkubeconfig) | **Get** /clusters/{clusterId}/kubeconfig | Read the kubeconfig
*DataPlatformClusterApi* | [**GetClusters**](docs/api/DataPlatformClusterApi.md#getclusters) | **Get** /clusters | List DataPlatformCluster
*DataPlatformClusterApi* | [**PatchCluster**](docs/api/DataPlatformClusterApi.md#patchcluster) | **Patch** /clusters/{clusterId} | Partially modify a DataPlatformCluster
*DataPlatformMetaDataApi* | [**VersionsGet**](docs/api/DataPlatformMetaDataApi.md#versionsget) | **Get** /versions | Managed Data Stack API version
*DataPlatformNodePoolApi* | [**CreateClusterNodepool**](docs/api/DataPlatformNodePoolApi.md#createclusternodepool) | **Post** /clusters/{clusterId}/nodepools | Create a DataPlatformNodePool for a distinct DataPlatformCluster
*DataPlatformNodePoolApi* | [**DeleteClusterNodepool**](docs/api/DataPlatformNodePoolApi.md#deleteclusternodepool) | **Delete** /clusters/{clusterId}/nodepools/{nodepoolId} | Remove node pool from DataPlatformCluster.
*DataPlatformNodePoolApi* | [**GetClusterNodepool**](docs/api/DataPlatformNodePoolApi.md#getclusternodepool) | **Get** /clusters/{clusterId}/nodepools/{nodepoolId} | Retrieve a DataPlatformNodePool
*DataPlatformNodePoolApi* | [**GetClusterNodepools**](docs/api/DataPlatformNodePoolApi.md#getclusternodepools) | **Get** /clusters/{clusterId}/nodepools | List the DataPlatformNodePools of a  DataPlatformCluster
*DataPlatformNodePoolApi* | [**PatchClusterNodepool**](docs/api/DataPlatformNodePoolApi.md#patchclusternodepool) | **Patch** /clusters/{clusterId}/nodepools/{nodepoolId} | Partially modify a DataPlatformNodePool


## Documentation For Models

 - [AvailabilityZone](docs/models/AvailabilityZone.md)
 - [Cluster](docs/models/Cluster.md)
 - [ClusterListResponseData](docs/models/ClusterListResponseData.md)
 - [ClusterResponseData](docs/models/ClusterResponseData.md)
 - [CreateClusterProperties](docs/models/CreateClusterProperties.md)
 - [CreateClusterRequest](docs/models/CreateClusterRequest.md)
 - [CreateNodePoolProperties](docs/models/CreateNodePoolProperties.md)
 - [CreateNodePoolRequest](docs/models/CreateNodePoolRequest.md)
 - [ErrorMessage](docs/models/ErrorMessage.md)
 - [ErrorResponse](docs/models/ErrorResponse.md)
 - [MaintenanceWindow](docs/models/MaintenanceWindow.md)
 - [Metadata](docs/models/Metadata.md)
 - [NodePool](docs/models/NodePool.md)
 - [NodePoolListResponseData](docs/models/NodePoolListResponseData.md)
 - [NodePoolResponseData](docs/models/NodePoolResponseData.md)
 - [PatchClusterProperties](docs/models/PatchClusterProperties.md)
 - [PatchClusterRequest](docs/models/PatchClusterRequest.md)
 - [PatchNodePoolProperties](docs/models/PatchNodePoolProperties.md)
 - [PatchNodePoolRequest](docs/models/PatchNodePoolRequest.md)
 - [StorageType](docs/models/StorageType.md)


## Documentation For Authorization



### basicAuth

- **Type**: HTTP basic authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
    UserName: "username",
    Password: "password",
})
r, err := client.Service.Operation(auth, args)
```


### tokenAuth

- **Type**: HTTP Bearer token authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "BEARERTOKENSTRING")
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author



