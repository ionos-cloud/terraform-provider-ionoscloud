[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)

# Go API client for ionoscloud

*Managed Stackable Data Platform* by IONOS Cloud provides a preconfigured Kubernetes cluster
with pre-installed and managed Stackable operators. After the provision of these Stackable operators,
the customer can interact with them directly
and build his desired application on top of the Stackable platform.

The Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API
in addition or as an alternative to the *Data Center Designer* (DCD).

## Getting Started

To get your DataPlatformCluster up and running, the following steps needs to be performed.

### IONOS Cloud Account

The first step is the creation of a IONOS Cloud account if not already existing.

To register a **new account**, visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).

### Virtual Data Center (VDC)

The Managed Stackable Data Platform needs a virtual data center (VDC) hosting the cluster.
This could either be a VDC that already exists, especially if you want to connect the managed data platform
to other services already running within your VDC.
Otherwise, if you want to place the Managed Stackable Data Platform in a new VDC or you have not yet created a VDC,
you need to do so.

A new VDC can be created via the IONOS Cloud API, the IONOS Cloud CLI (`ionosctl`), or the DCD Web interface.
For more information, see the
[official documentation](https://docs.ionos.com/cloud/getting-started/basic-tutorials/data-center-basics).

### Get a authentication token

To interact with this API a user specific authentication token is needed.
This token can be generated using the IONOS Cloud CLI the following way:

```
ionosctl token generate
```

For more information, [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate).

### Create a new DataPlatformCluster

Before using the Managed Stackable Data Platform, a new DataPlatformCluster must be created.

To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.

The provisioning of the cluster might take some time. To check the current provisioning status,
you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID
that was presented to you in the response of the create cluster call.

### Add a DataPlatformNodePool

To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool
that is provisioned along with the cluster is reserved for the Stackable operators.
You may create further node pools with resources tailored to your use case.

To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post)
endpoint.

### Receive Kubeconfig

Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API.
The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.

To protect the deployment of the Stackable distribution, the kubeconfig does not provide you with administration
rights for the cluster. What that means is that your actions and deployments are limited to the **default** namespace.

If you still want to group your deployments, you have the option to create subnamespaces within the default namespace.
This is made possible by the concept of *hierarchical namespaces* (HNS). You can find more details
[here](https://kubernetes.io/blog/2020/08/14/introducing-hierarchical-namespaces/).

The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint
using the cluster ID of the created DataPlatformCluster.

### Create Stackable Services

You can leverage the `kubeconfig.json` file to access the Managed Stackable Data Platform cluster and manage the
deployment of [Stackable data apps](https://stackable.tech/en/platform/).

With the Stackable operators, you can deploy the
[data apps](https://docs.stackable.tech/home/stable/getting_started.html#_deploying_stackable_services)
you want in your Data Platform cluster.

## Authorization

All endpoints are secured, so only an authenticated user can access them.
As Authentication mechanism the default IONOS Cloud authentication mechanism
is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).

### Basic Auth

The basic auth scheme uses the IONOS Cloud user credentials in form of a *Basic Authentication* header
accordingly to [RFC 7617](https://datatracker.ietf.org/doc/html/rfc7617).

### API Key as Bearer Token

The Bearer auth token used at the API Gateway is a user-related token created with the IONOS Cloud CLI
(For details, see the
[documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate)).
For every request to be authenticated, the token is passed as *Authorization Bearer* header along with the request.

### Permissions and Access Roles

Currently, an administrator can see and manipulate all resources in a contract.
Furthermore, users with the group privilege `Manage Dataplatform` can access the API.

## Components

The Managed Stackable Data Platform by IONOS Cloud consists of two components.
The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.

### DataPlatformCluster

A DataPlatformCluster is the virtual instance of the customer services and operations running the managed services
like Stackable operators.
A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer.
Therefore, it's possible to integrate the cluster with other resources as VLANs
e.g. to shape the data center in the customer's need
and integrate the cluster within the topology the customer wants to build.

In addition to the Kubernetes cluster, a small node pool is provided
which is exclusively used to run the Stackable operators.

### DataPlatformNodePool

A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top.
All nodes within a node pool are identical in setup.
The nodes of a pool are provisioned into virtual data centers at a location of your choice
and you can freely specify the properties of all the nodes at once before creation.

Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are read-only in the customer's VDC
and can only be modified or deleted via the API.

## References


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 1.2.0
- Package version: 1.1.1
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [https://www.ionos.com](https://www.ionos.com)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
go get github.com/antihax/optional
```

Put the package under your project folder and add the following in import:

```golang
import "./ionoscloud"
```

## Authentication

All available server URLs are:

- *https://api.ionos.com/dataplatform* - IONOS Cloud - Managed Stackable Data Platform API

By default, *https://api.ionos.com/dataplatform* is used, however this can be overriden at authentication, either
by setting the `IONOS_API_URL` environment variable or by specifying the `hostUrl` parameter when
initializing the sdk client.

The username and password or the authentication token can be manually specified when initializing
the sdk client:

```golang

client := ionoscloud.NewAPIClient(ionoscloud.NewConfiguration(username, password, token, hostUrl))

```

Environment variables can also be used. The sdk uses the following variables:
- IONOS_TOKEN    - login via token. This is the recommended way to authenticate.
- IONOS_USERNAME - to specify the username used to login
- IONOS_PASSWORD - to specify the password
- IONOS_API_URL  - to specify the API server URL

In this case, the client configuration needs to be initialized using `NewConfigurationFromEnv()`.

```golang

client := ionoscloud.NewAPIClient(ionoscloud.NewConfigurationFromEnv())

```


## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/dataplatform*
<details >
    <summary title="Click to toggle">API Endpoints table</summary>


| Class | Method | HTTP request | Description |
| ------------- | ------------- | ------------- | ------------- |
| DataPlatformClusterApi | [**ClustersDelete**](docs/api/DataPlatformClusterApi.md#ClustersDelete) | **Delete** /clusters/{clusterId} | Delete a DataPlatformCluster |
| DataPlatformClusterApi | [**ClustersFindById**](docs/api/DataPlatformClusterApi.md#ClustersFindById) | **Get** /clusters/{clusterId} | Retrieve a DataPlatformCluster |
| DataPlatformClusterApi | [**ClustersGet**](docs/api/DataPlatformClusterApi.md#ClustersGet) | **Get** /clusters | List the DataPlatformClusters |
| DataPlatformClusterApi | [**ClustersKubeconfigFindByClusterId**](docs/api/DataPlatformClusterApi.md#ClustersKubeconfigFindByClusterId) | **Get** /clusters/{clusterId}/kubeconfig | Read the Kubeconfig |
| DataPlatformClusterApi | [**ClustersPatch**](docs/api/DataPlatformClusterApi.md#ClustersPatch) | **Patch** /clusters/{clusterId} | Partially Modify a DataPlatformCluster |
| DataPlatformClusterApi | [**ClustersPost**](docs/api/DataPlatformClusterApi.md#ClustersPost) | **Post** /clusters | Create a DataPlatformCluster |
| DataPlatformMetaDataApi | [**VersionsGet**](docs/api/DataPlatformMetaDataApi.md#VersionsGet) | **Get** /versions | Managed Stackable Data Platform API Versions |
| DataPlatformNodePoolApi | [**ClustersNodepoolsDelete**](docs/api/DataPlatformNodePoolApi.md#ClustersNodepoolsDelete) | **Delete** /clusters/{clusterId}/nodepools/{nodepoolId} | Remove a DataPlatformNodePool from a DataPlatformCluster |
| DataPlatformNodePoolApi | [**ClustersNodepoolsFindById**](docs/api/DataPlatformNodePoolApi.md#ClustersNodepoolsFindById) | **Get** /clusters/{clusterId}/nodepools/{nodepoolId} | Retrieve a DataPlatformNodePool |
| DataPlatformNodePoolApi | [**ClustersNodepoolsGet**](docs/api/DataPlatformNodePoolApi.md#ClustersNodepoolsGet) | **Get** /clusters/{clusterId}/nodepools | List the DataPlatformNodePools of a DataPlatformCluster |
| DataPlatformNodePoolApi | [**ClustersNodepoolsPatch**](docs/api/DataPlatformNodePoolApi.md#ClustersNodepoolsPatch) | **Patch** /clusters/{clusterId}/nodepools/{nodepoolId} | Partially Modify a DataPlatformNodePool |
| DataPlatformNodePoolApi | [**ClustersNodepoolsPost**](docs/api/DataPlatformNodePoolApi.md#ClustersNodepoolsPost) | **Post** /clusters/{clusterId}/nodepools | Create a DataPlatformNodePool for a distinct DataPlatformCluster |

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com/dataplatform*
<details >
<summary title="Click to toggle">API models list</summary>

 - [AutoScaling](docs/models/AutoScaling)
 - [AutoScalingBase](docs/models/AutoScalingBase)
 - [AvailabilityZone](docs/models/AvailabilityZone)
 - [Cluster](docs/models/Cluster)
 - [ClusterListResponseData](docs/models/ClusterListResponseData)
 - [ClusterResponseData](docs/models/ClusterResponseData)
 - [CreateClusterProperties](docs/models/CreateClusterProperties)
 - [CreateClusterRequest](docs/models/CreateClusterRequest)
 - [CreateNodePoolProperties](docs/models/CreateNodePoolProperties)
 - [CreateNodePoolRequest](docs/models/CreateNodePoolRequest)
 - [ErrorMessage](docs/models/ErrorMessage)
 - [ErrorResponse](docs/models/ErrorResponse)
 - [Lan](docs/models/Lan)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [Metadata](docs/models/Metadata)
 - [NodePool](docs/models/NodePool)
 - [NodePoolListResponseData](docs/models/NodePoolListResponseData)
 - [NodePoolResponseData](docs/models/NodePoolResponseData)
 - [PatchClusterProperties](docs/models/PatchClusterProperties)
 - [PatchClusterRequest](docs/models/PatchClusterRequest)
 - [PatchNodePoolProperties](docs/models/PatchNodePoolProperties)
 - [PatchNodePoolRequest](docs/models/PatchNodePoolRequest)
 - [Route](docs/models/Route)
 - [StorageType](docs/models/StorageType)
 - [VersionsGet200Response](docs/models/VersionsGet200Response)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
