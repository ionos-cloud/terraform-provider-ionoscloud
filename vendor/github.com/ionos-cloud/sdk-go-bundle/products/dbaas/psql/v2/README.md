# Go API client for psql

An enterprise-grade Database is provided as a Service (DBaaS) solution that
can be managed through a browser-based \"Data Center Designer\" (DCD) tool or
via an easy to use API.

The API allows you to create additional PostgreSQL database clusters or modify existing
ones. It is designed to allow users to leverage the same power and
flexibility found within the DCD visual tool. Both tools are consistent with
their concepts and lend well to making the experience smooth and intuitive.


## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/psql.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/psql.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/psql@latest
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

- *https://api.ionos.com/databases/postgresql* - Production

By default, *https://api.ionos.com/databases/postgresql* is used, however this can be overriden at authentication, either
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
	psql "github.com/ionos-cloud/sdk-go-bundle/products/psql"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := psql.NewAPIClient(cfg)
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
        psql "github.com/ionos-cloud/sdk-go-bundle/products/psql"
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
        apiClient := psql.NewAPIClient(cfg)
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
         psql "github.com/ionos-cloud/sdk-go-bundle/products/psql"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := psql.NewAPIClient(cfg)
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
        psql "github.com/ionos-cloud/sdk-go-bundle/products/psql"
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
    apiClient := psql.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/databases/postgresql*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
BackupsApi | [**ClusterBackupsGet**](docs/api/BackupsApi.md#clusterbackupsget) | **Get** /clusters/{clusterId}/backups | List backups of cluster
BackupsApi | [**ClustersBackupsFindById**](docs/api/BackupsApi.md#clustersbackupsfindbyid) | **Get** /clusters/backups/{backupId} | Fetch a cluster backup
BackupsApi | [**ClustersBackupsGet**](docs/api/BackupsApi.md#clustersbackupsget) | **Get** /clusters/backups | List cluster backups
ClustersApi | [**ClusterPostgresVersionsGet**](docs/api/ClustersApi.md#clusterpostgresversionsget) | **Get** /clusters/{clusterId}/postgresversions | List PostgreSQL versions
ClustersApi | [**ClustersDelete**](docs/api/ClustersApi.md#clustersdelete) | **Delete** /clusters/{clusterId} | Delete a cluster
ClustersApi | [**ClustersFindById**](docs/api/ClustersApi.md#clustersfindbyid) | **Get** /clusters/{clusterId} | Fetch a cluster
ClustersApi | [**ClustersGet**](docs/api/ClustersApi.md#clustersget) | **Get** /clusters | List clusters
ClustersApi | [**ClustersPatch**](docs/api/ClustersApi.md#clusterspatch) | **Patch** /clusters/{clusterId} | Patch a cluster
ClustersApi | [**ClustersPost**](docs/api/ClustersApi.md#clusterspost) | **Post** /clusters | Create a cluster
ClustersApi | [**ClustersVersionsGet**](docs/api/ClustersApi.md#clustersversionsget) | **Get** /clusters/{clusterId}/versions | Supported PostgreSQL versions of cluster
ClustersApi | [**PostgresVersionsGet**](docs/api/ClustersApi.md#postgresversionsget) | **Get** /clusters/postgresversions | List all PostgreSQL versions
DatabasesApi | [**DatabasesDelete**](docs/api/DatabasesApi.md#databasesdelete) | **Delete** /clusters/{clusterId}/databases/{databasename} | Delete database
DatabasesApi | [**DatabasesGet**](docs/api/DatabasesApi.md#databasesget) | **Get** /clusters/{clusterId}/databases/{databasename} | Get database
DatabasesApi | [**DatabasesList**](docs/api/DatabasesApi.md#databaseslist) | **Get** /clusters/{clusterId}/databases | List databases
DatabasesApi | [**DatabasesPost**](docs/api/DatabasesApi.md#databasespost) | **Post** /clusters/{clusterId}/databases | Create a database
LogsApi | [**ClusterLogsGet**](docs/api/LogsApi.md#clusterlogsget) | **Get** /clusters/{clusterId}/logs | Get logs of your cluster
MetadataApi | [**InfosVersionGet**](docs/api/MetadataApi.md#infosversionget) | **Get** /infos/version | Get the current API version
MetadataApi | [**InfosVersionsGet**](docs/api/MetadataApi.md#infosversionsget) | **Get** /infos/versions | Fetch all API versions
MetadataApi | [**VersionsGet**](docs/api/MetadataApi.md#versionsget) | **Get** /versions | PostgreSQL versions for new clusters
RestoresApi | [**ClusterRestorePost**](docs/api/RestoresApi.md#clusterrestorepost) | **Post** /clusters/{clusterId}/restore | In-place restore of a cluster
UsersApi | [**UsersDelete**](docs/api/UsersApi.md#usersdelete) | **Delete** /clusters/{clusterId}/users/{username} | Delete user
UsersApi | [**UsersGet**](docs/api/UsersApi.md#usersget) | **Get** /clusters/{clusterId}/users/{username} | Get user
UsersApi | [**UsersList**](docs/api/UsersApi.md#userslist) | **Get** /clusters/{clusterId}/users | List users
UsersApi | [**UsersPatch**](docs/api/UsersApi.md#userspatch) | **Patch** /clusters/{clusterId}/users/{username} | Patch user
UsersApi | [**UsersPost**](docs/api/UsersApi.md#userspost) | **Post** /clusters/{clusterId}/users | Create a user

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com/databases/postgresql*
<details >
<summary title="Click to toggle">API models list</summary>

 - [APIVersion](docs/models/APIVersion)
 - [BackupMetadata](docs/models/BackupMetadata)
 - [BackupResponse](docs/models/BackupResponse)
 - [ClusterBackup](docs/models/ClusterBackup)
 - [ClusterBackupList](docs/models/ClusterBackupList)
 - [ClusterBackupListAllOf](docs/models/ClusterBackupListAllOf)
 - [ClusterList](docs/models/ClusterList)
 - [ClusterListAllOf](docs/models/ClusterListAllOf)
 - [ClusterLogs](docs/models/ClusterLogs)
 - [ClusterLogsInstances](docs/models/ClusterLogsInstances)
 - [ClusterLogsInstancesMessages](docs/models/ClusterLogsInstancesMessages)
 - [ClusterMetadata](docs/models/ClusterMetadata)
 - [ClusterProperties](docs/models/ClusterProperties)
 - [ClusterResponse](docs/models/ClusterResponse)
 - [Connection](docs/models/Connection)
 - [ConnectionPooler](docs/models/ConnectionPooler)
 - [CreateClusterProperties](docs/models/CreateClusterProperties)
 - [CreateClusterRequest](docs/models/CreateClusterRequest)
 - [CreateRestoreRequest](docs/models/CreateRestoreRequest)
 - [DBUser](docs/models/DBUser)
 - [Database](docs/models/Database)
 - [DatabaseItems](docs/models/DatabaseItems)
 - [DatabaseList](docs/models/DatabaseList)
 - [DatabaseProperties](docs/models/DatabaseProperties)
 - [DatabaseResource](docs/models/DatabaseResource)
 - [DayOfTheWeek](docs/models/DayOfTheWeek)
 - [DeprecatedPagination](docs/models/DeprecatedPagination)
 - [ErrorMessage](docs/models/ErrorMessage)
 - [ErrorResponse](docs/models/ErrorResponse)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [Metadata](docs/models/Metadata)
 - [Pagination](docs/models/Pagination)
 - [PaginationLinks](docs/models/PaginationLinks)
 - [PatchClusterProperties](docs/models/PatchClusterProperties)
 - [PatchClusterRequest](docs/models/PatchClusterRequest)
 - [PatchUserProperties](docs/models/PatchUserProperties)
 - [PoolMode](docs/models/PoolMode)
 - [PostgresVersionList](docs/models/PostgresVersionList)
 - [PostgresVersionListData](docs/models/PostgresVersionListData)
 - [Resource](docs/models/Resource)
 - [ResourceMetadata](docs/models/ResourceMetadata)
 - [ResourceType](docs/models/ResourceType)
 - [State](docs/models/State)
 - [StorageType](docs/models/StorageType)
 - [SynchronizationMode](docs/models/SynchronizationMode)
 - [User](docs/models/User)
 - [UserItems](docs/models/UserItems)
 - [UserList](docs/models/UserList)
 - [UserProperties](docs/models/UserProperties)
 - [UserResource](docs/models/UserResource)
 - [UsersPatchRequest](docs/models/UsersPatchRequest)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
