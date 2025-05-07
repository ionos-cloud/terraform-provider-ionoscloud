# Go API client for mongo

With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.

MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.

The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.


## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/mongo.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/mongo.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/mongo@latest
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

- *https://api.ionos.com/databases/mongodb* - Production

By default, *https://api.ionos.com/databases/mongodb* is used, however this can be overriden at authentication, either
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
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/mongo"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := mongo.NewAPIClient(cfg)
	return nil
}
```
### Token Authentication
There are 2 ways to generate your token:

 ### Generate token using sdk for [auth](https://github.com/ionos-cloud/sdk-go-bundle/products/auth):
```golang
    import (
        "context"
        "fmt"
        "github.com/ionos-cloud/sdk-go-bundle/products/auth"
        "github.com/ionos-cloud/sdk-go-bundle/shared"
        mongo "github.com/ionos-cloud/sdk-go-bundle/products/mongo"
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
        apiClient := mongo.NewAPIClient(cfg)
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
         mongo "github.com/ionos-cloud/sdk-go-bundle/products/mongo"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := mongo.NewAPIClient(cfg)
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
        mongo "github.com/ionos-cloud/sdk-go-bundle/products/mongo"
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
    apiClient := mongo.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/databases/mongodb*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
ClustersApi | [**ClustersDelete**](docs/api/ClustersApi.md#clustersdelete) | **Delete** /clusters/{clusterId} | Delete a Cluster
ClustersApi | [**ClustersFindById**](docs/api/ClustersApi.md#clustersfindbyid) | **Get** /clusters/{clusterId} | Get a cluster by id
ClustersApi | [**ClustersGet**](docs/api/ClustersApi.md#clustersget) | **Get** /clusters | Get Clusters
ClustersApi | [**ClustersPatch**](docs/api/ClustersApi.md#clusterspatch) | **Patch** /clusters/{clusterId} | Patch a cluster
ClustersApi | [**ClustersPost**](docs/api/ClustersApi.md#clusterspost) | **Post** /clusters | Create a Cluster
ClustersApi | [**ClustersVersionsGet**](docs/api/ClustersApi.md#clustersversionsget) | **Get** /clusters/{clusterId}/versions | Get available MongoDB versions for this cluster
LogsApi | [**ClustersLogsGet**](docs/api/LogsApi.md#clusterslogsget) | **Get** /clusters/{clusterId}/logs | Get logs of your cluster
MetadataApi | [**InfosVersionGet**](docs/api/MetadataApi.md#infosversionget) | **Get** /infos/version | Get API Version
MetadataApi | [**InfosVersionsGet**](docs/api/MetadataApi.md#infosversionsget) | **Get** /infos/versions | Get All API Versions
MetadataApi | [**VersionsGet**](docs/api/MetadataApi.md#versionsget) | **Get** /versions | Get available MongoDB versions
RestoresApi | [**ClustersRestorePost**](docs/api/RestoresApi.md#clustersrestorepost) | **Post** /clusters/{clusterId}/restore | In-place restore of a cluster
SnapshotsApi | [**ClustersSnapshotsGet**](docs/api/SnapshotsApi.md#clusterssnapshotsget) | **Get** /clusters/{clusterId}/snapshots | Get the snapshots of your cluster
TemplatesApi | [**TemplatesGet**](docs/api/TemplatesApi.md#templatesget) | **Get** /templates | Get Templates
UsersApi | [**ClustersUsersDelete**](docs/api/UsersApi.md#clustersusersdelete) | **Delete** /clusters/{clusterId}/users/{username} | Delete a MongoDB User by ID
UsersApi | [**ClustersUsersFindById**](docs/api/UsersApi.md#clustersusersfindbyid) | **Get** /clusters/{clusterId}/users/{username} | Get a MongoDB User by ID
UsersApi | [**ClustersUsersGet**](docs/api/UsersApi.md#clustersusersget) | **Get** /clusters/{clusterId}/users | Get all Cluster Users
UsersApi | [**ClustersUsersPatch**](docs/api/UsersApi.md#clustersuserspatch) | **Patch** /clusters/{clusterId}/users/{username} | Patch a MongoDB User by ID
UsersApi | [**ClustersUsersPost**](docs/api/UsersApi.md#clustersuserspost) | **Post** /clusters/{clusterId}/users | Create MongoDB User

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com/databases/mongodb*
<details >
<summary title="Click to toggle">API models list</summary>

 - [APIVersion](docs/models/APIVersion)
 - [BackupProperties](docs/models/BackupProperties)
 - [BiConnectorProperties](docs/models/BiConnectorProperties)
 - [ClusterList](docs/models/ClusterList)
 - [ClusterListAllOf](docs/models/ClusterListAllOf)
 - [ClusterLogs](docs/models/ClusterLogs)
 - [ClusterLogsInstances](docs/models/ClusterLogsInstances)
 - [ClusterLogsInstancesMessages](docs/models/ClusterLogsInstancesMessages)
 - [ClusterProperties](docs/models/ClusterProperties)
 - [ClusterResponse](docs/models/ClusterResponse)
 - [Connection](docs/models/Connection)
 - [CreateClusterProperties](docs/models/CreateClusterProperties)
 - [CreateClusterRequest](docs/models/CreateClusterRequest)
 - [CreateRestoreRequest](docs/models/CreateRestoreRequest)
 - [DayOfTheWeek](docs/models/DayOfTheWeek)
 - [ErrorMessage](docs/models/ErrorMessage)
 - [ErrorResponse](docs/models/ErrorResponse)
 - [Health](docs/models/Health)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [Metadata](docs/models/Metadata)
 - [MongoDBVersionList](docs/models/MongoDBVersionList)
 - [MongoDBVersionListData](docs/models/MongoDBVersionListData)
 - [Pagination](docs/models/Pagination)
 - [PaginationLinks](docs/models/PaginationLinks)
 - [PatchClusterProperties](docs/models/PatchClusterProperties)
 - [PatchClusterRequest](docs/models/PatchClusterRequest)
 - [PatchUserProperties](docs/models/PatchUserProperties)
 - [PatchUserRequest](docs/models/PatchUserRequest)
 - [ResourceType](docs/models/ResourceType)
 - [SnapshotList](docs/models/SnapshotList)
 - [SnapshotListAllOf](docs/models/SnapshotListAllOf)
 - [SnapshotProperties](docs/models/SnapshotProperties)
 - [SnapshotResponse](docs/models/SnapshotResponse)
 - [State](docs/models/State)
 - [StorageType](docs/models/StorageType)
 - [TemplateList](docs/models/TemplateList)
 - [TemplateListAllOf](docs/models/TemplateListAllOf)
 - [TemplateProperties](docs/models/TemplateProperties)
 - [TemplateResponse](docs/models/TemplateResponse)
 - [User](docs/models/User)
 - [UserMetadata](docs/models/UserMetadata)
 - [UserProperties](docs/models/UserProperties)
 - [UserRoles](docs/models/UserRoles)
 - [UsersList](docs/models/UsersList)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
