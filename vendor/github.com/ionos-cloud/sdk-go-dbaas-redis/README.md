# Go API client for ionoscloud

Redis Database API

## Overview

You can use this SDK to manage your Redis Database resources.

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-redis"
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
ctx := context.WithValue(context.Background(), ionoscloud.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), ionoscloud.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

## Documentation for API Endpoints

All URIs are relative to *https://redis.de-fra.ionos.com*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ReplicaSetApi* | [**ReplicasetsDelete**](docs/api/ReplicaSetApi.md#replicasetsdelete) | **Delete** /replicasets/{replicasetId} | Delete ReplicaSet
*ReplicaSetApi* | [**ReplicasetsFindById**](docs/api/ReplicaSetApi.md#replicasetsfindbyid) | **Get** /replicasets/{replicasetId} | Retrieve ReplicaSet
*ReplicaSetApi* | [**ReplicasetsGet**](docs/api/ReplicaSetApi.md#replicasetsget) | **Get** /replicasets | Retrieve all Replicaset
*ReplicaSetApi* | [**ReplicasetsPost**](docs/api/ReplicaSetApi.md#replicasetspost) | **Post** /replicasets | Create ReplicaSet
*ReplicaSetApi* | [**ReplicasetsPut**](docs/api/ReplicaSetApi.md#replicasetsput) | **Put** /replicasets/{replicasetId} | Ensure ReplicaSet
*RestoreApi* | [**SnapshotsRestoresFindById**](docs/api/RestoreApi.md#snapshotsrestoresfindbyid) | **Get** /snapshots/{snapshotId}/restores/{restoreId} | Retrieve Restore
*RestoreApi* | [**SnapshotsRestoresGet**](docs/api/RestoreApi.md#snapshotsrestoresget) | **Get** /snapshots/{snapshotId}/restores | Retrieve all Restore
*RestoreApi* | [**SnapshotsRestoresPost**](docs/api/RestoreApi.md#snapshotsrestorespost) | **Post** /snapshots/{snapshotId}/restores | Create Restore
*SnapshotApi* | [**SnapshotsFindById**](docs/api/SnapshotApi.md#snapshotsfindbyid) | **Get** /snapshots/{snapshotId} | Retrieve Snapshot
*SnapshotApi* | [**SnapshotsGet**](docs/api/SnapshotApi.md#snapshotsget) | **Get** /snapshots | Retrieve all Snapshot


## Documentation For Models

 - [Connection](docs/models/Connection.md)
 - [DayOfTheWeek](docs/models/DayOfTheWeek.md)
 - [Error](docs/models/Error.md)
 - [ErrorMessages](docs/models/ErrorMessages.md)
 - [EvictionPolicy](docs/models/EvictionPolicy.md)
 - [HashedPassword](docs/models/HashedPassword.md)
 - [Links](docs/models/Links.md)
 - [MaintenanceWindow](docs/models/MaintenanceWindow.md)
 - [Metadata](docs/models/Metadata.md)
 - [Pagination](docs/models/Pagination.md)
 - [PersistenceMode](docs/models/PersistenceMode.md)
 - [ReplicaSet](docs/models/ReplicaSet.md)
 - [ReplicaSetCreate](docs/models/ReplicaSetCreate.md)
 - [ReplicaSetEnsure](docs/models/ReplicaSetEnsure.md)
 - [ReplicaSetMetadata](docs/models/ReplicaSetMetadata.md)
 - [ReplicaSetMetadataAllOf](docs/models/ReplicaSetMetadataAllOf.md)
 - [ReplicaSetRead](docs/models/ReplicaSetRead.md)
 - [ReplicaSetReadList](docs/models/ReplicaSetReadList.md)
 - [ReplicaSetReadListAllOf](docs/models/ReplicaSetReadListAllOf.md)
 - [ResourceState](docs/models/ResourceState.md)
 - [Resources](docs/models/Resources.md)
 - [Restore](docs/models/Restore.md)
 - [RestoreCreate](docs/models/RestoreCreate.md)
 - [RestoreMetadata](docs/models/RestoreMetadata.md)
 - [RestoreMetadataAllOf](docs/models/RestoreMetadataAllOf.md)
 - [RestoreRead](docs/models/RestoreRead.md)
 - [RestoreReadList](docs/models/RestoreReadList.md)
 - [RestoreReadListAllOf](docs/models/RestoreReadListAllOf.md)
 - [SnapshotCreate](docs/models/SnapshotCreate.md)
 - [SnapshotEnsure](docs/models/SnapshotEnsure.md)
 - [SnapshotMetadata](docs/models/SnapshotMetadata.md)
 - [SnapshotMetadataAllOf](docs/models/SnapshotMetadataAllOf.md)
 - [SnapshotRead](docs/models/SnapshotRead.md)
 - [SnapshotReadList](docs/models/SnapshotReadList.md)
 - [SnapshotReadListAllOf](docs/models/SnapshotReadListAllOf.md)
 - [User](docs/models/User.md)
 - [UserPassword](docs/models/UserPassword.md)


## Documentation For Authorization


Authentication schemes defined for the API:
### tokenAuth

- **Type**: HTTP Bearer token authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "BEARER_TOKEN_STRING")
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

support@cloud.ionos.com

