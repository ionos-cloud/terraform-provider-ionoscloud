# Go API client for dataplatform

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
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/dataplatform.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/dataplatform.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/dataplatform@latest
```

## Environment Variables

| Environment Variable | Description                                                                                                                                                                                                                    |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IONOS_USERNAME`     | Specify the username used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_PASSWORD`     | Specify the password used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_TOKEN`        | Specify the token used to login, if a token is being used instead of username and password                                                                                                                                     |
| `IONOS_API_URL`      | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v6` path, so it should _not_ be included in the `IONOS_API_URL` environment variable |
| `IONOS_LOG_LEVEL`    | Specify the Log Level used to log messages. Possible values: Off, Debug, Trace |
| `IONOS_PINNED_CERT`  | Specify the SHA-256 public fingerprint here, enables certificate pinning                                                                                                                                                       |

⚠️ **_Note: To overwrite the api endpoint - `api.ionos.com`, the environment variable `IONOS_API_URL` can be set, and used with `NewConfigurationFromEnv()` function._**

## Examples

Examples for creating resources using the Go SDK can be found [here](examples/)

## Authentication

All available server URLs are:

- *https://api.ionos.com/dataplatform* - IONOS Cloud - Managed Stackable Data Platform API

By default, *https://api.ionos.com/dataplatform* is used, however this can be overriden at authentication, either
by setting the `IONOS_API_URL` environment variable or by specifying the `hostUrl` parameter when
initializing the sdk client.

**NOTE**: We recommend passing the URL without the `https://` or `http://` prefix. The SDK
checks and adds it if necessary when configurations are created using `NewConfiguration` or
`NewConfigurationFromEnv`. This is to avoid issues caused by typos in the prefix that cannot
 be easily detected and debugged.

### Basic Authentication

- **Type**: HTTP basic authentication

Example

```golang
import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := dataplatform.NewAPIClient(cfg)
	return nil
}
```
### Token Authentication
There are 2 ways to generate your token:

 ### Generate token using sdk for [auth](https://github.com/ionos-cloud/products/auth):
```golang
    import (
        "context"
        "fmt"
        "github.com/ionos-cloud/sdk-go-bundle/products/auth"
        "github.com/ionos-cloud/sdk-go-bundle/shared"
        dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_USERNAME and IONOS_PASSWORD as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        jwt, _, err := auth.TokensApi.TokensGenerate(context.Background()).Execute()
        if err != nil {
            return fmt.Errorf("error occurred while generating token (%w)", err)
        }
        if !jwt.HasToken() {
            return fmt.Errorf("could not generate token")
        }
        cfg := shared.NewConfiguration("", "", *jwt.GetToken(), "hostUrl_here")
        cfg.LogLevel = Trace
        apiClient := dataplatform.NewAPIClient(cfg)
        return nil
    }
```
 ### Generate token using ionosctl:
  Install ionosctl as explained [here](https://github.com/ionos-cloud/ionosctl)
  Run commands to login and generate your token.
```golang
    ionosctl login
    ionosctl token generate
    export IONOS_TOKEN="insert_here_token_saved_from_generate_command"
```
 Save the generated token and use it to authenticate:
```golang
    import (
        "context"
        "fmt"
        "github.com/ionos-cloud/sdk-go-bundle/products/auth"
         dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := dataplatform.NewAPIClient(cfg)
        return nil
    }
```

## Certificate pinning:

You can enable certificate pinning if you want to bypass the normal certificate checking procedure,
by doing the following:

Set env variable IONOS_PINNED_CERT=<insert_sha256_public_fingerprint_here>

You can get the sha256 fingerprint most easily from the browser by inspecting the certificate.

### Depth

Many of the _List_ or _Get_ operations will accept an optional _depth_ argument. Setting this to a value between 0 and 5 affects the amount of data that is returned. The details returned vary depending on the resource being queried, but it generally follows this pattern. By default, the SDK sets the _depth_ argument to the maximum value.

| Depth | Description |
| :--- | :--- |
| 0 | Only direct properties are included. Children are not included. |
| 1 | Direct properties and children's references are returned. |
| 2 | Direct properties and children's properties are returned. |
| 3 | Direct properties, children's properties, and descendants' references are returned. |
| 4 | Direct properties, children's properties, and descendants' properties are returned. |
| 5 | Returns all available properties. |

### Changing the base URL

Base URL for the HTTP operation can be changed by using the following function:

```go
requestProperties.SetURL("https://api.ionos.com/cloudapi/v6")
```

## Debugging

You can inject any logger that implements Printf as a logger
instead of using the default sdk logger.
There are log levels that you can set: `Off`, `Debug` and `Trace`.
`Off` - does not show any logs
`Debug` - regular logs, no sensitive information
`Trace` - we recommend you only set this field for debugging purposes. Disable it in your production environments because it can log sensitive data.
          It logs the full request and response without encryption, even for an HTTPS call. Verbose request and response logging can also significantly impact your application's performance.


```golang
package main

    import (
        dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
        "github.com/ionos-cloud/sdk-go-bundle/shared"
        "github.com/sirupsen/logrus"
    )

func main() {
    // create your configuration. replace username, password, token and url with correct values, or use NewConfigurationFromEnv()
    // if you have set your env variables as explained above
    cfg := shared.NewConfiguration("username", "password", "token", "hostUrl")
    // enable request and response logging. this is the most verbose loglevel
    shared.SdkLogLevel = Trace
    // inject your own logger that implements Printf
    shared.SdkLogger = logrus.New()
    // create you api client with the configuration
    apiClient := dataplatform.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/dataplatform*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
DataPlatformClusterApi | [**ClustersDelete**](docs/api/DataPlatformClusterApi.md#clustersdelete) | **Delete** /clusters/{clusterId} | Delete a DataPlatformCluster
DataPlatformClusterApi | [**ClustersFindById**](docs/api/DataPlatformClusterApi.md#clustersfindbyid) | **Get** /clusters/{clusterId} | Retrieve a DataPlatformCluster
DataPlatformClusterApi | [**ClustersGet**](docs/api/DataPlatformClusterApi.md#clustersget) | **Get** /clusters | List the DataPlatformClusters
DataPlatformClusterApi | [**ClustersKubeconfigFindByClusterId**](docs/api/DataPlatformClusterApi.md#clusterskubeconfigfindbyclusterid) | **Get** /clusters/{clusterId}/kubeconfig | Read the Kubeconfig
DataPlatformClusterApi | [**ClustersPatch**](docs/api/DataPlatformClusterApi.md#clusterspatch) | **Patch** /clusters/{clusterId} | Partially Modify a DataPlatformCluster
DataPlatformClusterApi | [**ClustersPost**](docs/api/DataPlatformClusterApi.md#clusterspost) | **Post** /clusters | Create a DataPlatformCluster
DataPlatformMetaDataApi | [**VersionsGet**](docs/api/DataPlatformMetaDataApi.md#versionsget) | **Get** /versions | Managed Stackable Data Platform API Versions
DataPlatformNodePoolApi | [**ClustersNodepoolsDelete**](docs/api/DataPlatformNodePoolApi.md#clustersnodepoolsdelete) | **Delete** /clusters/{clusterId}/nodepools/{nodepoolId} | Remove a DataPlatformNodePool from a DataPlatformCluster
DataPlatformNodePoolApi | [**ClustersNodepoolsFindById**](docs/api/DataPlatformNodePoolApi.md#clustersnodepoolsfindbyid) | **Get** /clusters/{clusterId}/nodepools/{nodepoolId} | Retrieve a DataPlatformNodePool
DataPlatformNodePoolApi | [**ClustersNodepoolsGet**](docs/api/DataPlatformNodePoolApi.md#clustersnodepoolsget) | **Get** /clusters/{clusterId}/nodepools | List the DataPlatformNodePools of a DataPlatformCluster
DataPlatformNodePoolApi | [**ClustersNodepoolsPatch**](docs/api/DataPlatformNodePoolApi.md#clustersnodepoolspatch) | **Patch** /clusters/{clusterId}/nodepools/{nodepoolId} | Partially Modify a DataPlatformNodePool
DataPlatformNodePoolApi | [**ClustersNodepoolsPost**](docs/api/DataPlatformNodePoolApi.md#clustersnodepoolspost) | **Post** /clusters/{clusterId}/nodepools | Create a DataPlatformNodePool for a distinct DataPlatformCluster

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
