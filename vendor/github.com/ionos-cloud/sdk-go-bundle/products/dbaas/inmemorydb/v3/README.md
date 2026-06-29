# Go API client for inmemorydb

API description for the IONOS In-Memory DB

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb@latest
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

- *https://in-memory-db.de-fra.ionos.com* - Production de-fra
- *https://in-memory-db.de-txl.ionos.com* - Production de-txl
- *https://in-memory-db.es-vit.ionos.com* - Production es-vit
- *https://in-memory-db.gb-bhx.ionos.com* - Production gb-bhx
- *https://in-memory-db.gb-lhr.ionos.com* - Production gb-lhr
- *https://in-memory-db.us-ewr.ionos.com* - Production us-ewr
- *https://in-memory-db.us-las.ionos.com* - Production us-las
- *https://in-memory-db.us-mci.ionos.com* - Production us-mci
- *https://in-memory-db.fr-par.ionos.com* - Production fr-par

By default, *https://in-memory-db.de-fra.ionos.com* is used, however this can be overriden at authentication, either
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
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := inmemorydb.NewAPIClient(cfg)
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
        inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb"
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
        apiClient := inmemorydb.NewAPIClient(cfg)
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
         inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := inmemorydb.NewAPIClient(cfg)
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
        inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/inmemorydb"
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
    apiClient := inmemorydb.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://in-memory-db.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
ReplicaSetApi | [**ReplicasetsDelete**](docs/api/ReplicaSetApi.md#replicasetsdelete) | **Delete** /replicasets/{replicasetId} | Delete ReplicaSet
ReplicaSetApi | [**ReplicasetsFindById**](docs/api/ReplicaSetApi.md#replicasetsfindbyid) | **Get** /replicasets/{replicasetId} | Retrieve ReplicaSet
ReplicaSetApi | [**ReplicasetsGet**](docs/api/ReplicaSetApi.md#replicasetsget) | **Get** /replicasets | Retrieve all ReplicaSet
ReplicaSetApi | [**ReplicasetsPost**](docs/api/ReplicaSetApi.md#replicasetspost) | **Post** /replicasets | Create ReplicaSet
ReplicaSetApi | [**ReplicasetsPut**](docs/api/ReplicaSetApi.md#replicasetsput) | **Put** /replicasets/{replicasetId} | Ensure ReplicaSet
RestoreApi | [**SnapshotsRestoresFindById**](docs/api/RestoreApi.md#snapshotsrestoresfindbyid) | **Get** /snapshots/{snapshotId}/restores/{restoreId} | Retrieve Restore
RestoreApi | [**SnapshotsRestoresGet**](docs/api/RestoreApi.md#snapshotsrestoresget) | **Get** /snapshots/{snapshotId}/restores | Retrieve all Restore
RestoreApi | [**SnapshotsRestoresPost**](docs/api/RestoreApi.md#snapshotsrestorespost) | **Post** /snapshots/{snapshotId}/restores | Create Restore
SnapshotApi | [**SnapshotsFindById**](docs/api/SnapshotApi.md#snapshotsfindbyid) | **Get** /snapshots/{snapshotId} | Retrieve Snapshot
SnapshotApi | [**SnapshotsGet**](docs/api/SnapshotApi.md#snapshotsget) | **Get** /snapshots | Retrieve all Snapshot

</details>

## Documentation For Models

All URIs are relative to *https://in-memory-db.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [BackupProperties](docs/models/BackupProperties)
 - [Connection](docs/models/Connection)
 - [DayOfTheWeek](docs/models/DayOfTheWeek)
 - [Error](docs/models/Error)
 - [ErrorMessages](docs/models/ErrorMessages)
 - [EvictionPolicy](docs/models/EvictionPolicy)
 - [HashedPassword](docs/models/HashedPassword)
 - [Links](docs/models/Links)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [Metadata](docs/models/Metadata)
 - [Pagination](docs/models/Pagination)
 - [PersistenceMode](docs/models/PersistenceMode)
 - [ReplicaSet](docs/models/ReplicaSet)
 - [ReplicaSetCreate](docs/models/ReplicaSetCreate)
 - [ReplicaSetEnsure](docs/models/ReplicaSetEnsure)
 - [ReplicaSetMetadata](docs/models/ReplicaSetMetadata)
 - [ReplicaSetMetadataAllOf](docs/models/ReplicaSetMetadataAllOf)
 - [ReplicaSetRead](docs/models/ReplicaSetRead)
 - [ReplicaSetReadList](docs/models/ReplicaSetReadList)
 - [ReplicaSetReadListAllOf](docs/models/ReplicaSetReadListAllOf)
 - [ResourceState](docs/models/ResourceState)
 - [Resources](docs/models/Resources)
 - [Restore](docs/models/Restore)
 - [RestoreCreate](docs/models/RestoreCreate)
 - [RestoreMetadata](docs/models/RestoreMetadata)
 - [RestoreMetadataAllOf](docs/models/RestoreMetadataAllOf)
 - [RestoreRead](docs/models/RestoreRead)
 - [RestoreReadList](docs/models/RestoreReadList)
 - [RestoreReadListAllOf](docs/models/RestoreReadListAllOf)
 - [SnapshotCreate](docs/models/SnapshotCreate)
 - [SnapshotEnsure](docs/models/SnapshotEnsure)
 - [SnapshotMetadata](docs/models/SnapshotMetadata)
 - [SnapshotMetadataAllOf](docs/models/SnapshotMetadataAllOf)
 - [SnapshotRead](docs/models/SnapshotRead)
 - [SnapshotReadList](docs/models/SnapshotReadList)
 - [SnapshotReadListAllOf](docs/models/SnapshotReadListAllOf)
 - [User](docs/models/User)
 - [UserPassword](docs/models/UserPassword)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
