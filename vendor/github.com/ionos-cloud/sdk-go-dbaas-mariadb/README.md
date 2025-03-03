[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)

# Go API client for ionoscloud

An enterprise-grade Database is provided as a Service (DBaaS) solution that
can be managed through a browser-based \"Data Center Designer\" (DCD) tool or
via an easy to use API.

The API allows you to create additional MariaDB database clusters or modify existing
ones. It is designed to allow users to leverage the same power and
flexibility found within the DCD visual tool. Both tools are consistent with
their concepts and lend well to making the experience smooth and intuitive.


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.1.0
- Package version: 1.1.3
- Build package: org.openapitools.codegen.languages.GoClientCodegen

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

- *https://mariadb.de-txl.ionos.com* - Berlin, Germany
- *https://mariadb.de-fra.ionos.com* - Frankfurt, Germany
- *https://mariadb.es-vit.ionos.com* - Logroño, Spain
- *https://mariadb.fr-par.ionos.com* - Paris, France
- *https://mariadb.gb-lhr.ionos.com* - London, Great Britain
- *https://mariadb.us-ewr.ionos.com* - Newark, USA
- *https://mariadb.us-las.ionos.com* - Las Vegas, USA
- *https://mariadb.us-mci.ionos.com* - Lenexa, USA

By default, *https://mariadb.de-txl.ionos.com* is used, however this can be overriden at authentication, either
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

All URIs are relative to *https://mariadb.de-txl.ionos.com*
<details >
    <summary title="Click to toggle">API Endpoints table</summary>


| Class | Method | HTTP request | Description |
| ------------- | ------------- | ------------- | ------------- |
| BackupsApi | [**BackupsFindById**](docs/api/BackupsApi.md#BackupsFindById) | **Get** /backups/{backupId} | Fetch backups |
| BackupsApi | [**BackupsGet**](docs/api/BackupsApi.md#BackupsGet) | **Get** /backups | List of backups. |
| BackupsApi | [**ClusterBackupsGet**](docs/api/BackupsApi.md#ClusterBackupsGet) | **Get** /clusters/{clusterId}/backups | List backups of cluster |
| ClustersApi | [**ClustersDelete**](docs/api/ClustersApi.md#ClustersDelete) | **Delete** /clusters/{clusterId} | Delete a cluster |
| ClustersApi | [**ClustersFindById**](docs/api/ClustersApi.md#ClustersFindById) | **Get** /clusters/{clusterId} | Fetch a cluster |
| ClustersApi | [**ClustersGet**](docs/api/ClustersApi.md#ClustersGet) | **Get** /clusters | List clusters |
| ClustersApi | [**ClustersPatch**](docs/api/ClustersApi.md#ClustersPatch) | **Patch** /clusters/{clusterId} | Update a cluster |
| ClustersApi | [**ClustersPost**](docs/api/ClustersApi.md#ClustersPost) | **Post** /clusters | Create a cluster |
| RestoreApi | [**ClustersRestore**](docs/api/RestoreApi.md#ClustersRestore) | **Post** /clusters/{clusterId}/restore | In-place restore of a cluster. |

</details>

## Documentation For Models

All URIs are relative to *https://mariadb.de-txl.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [Backup](docs/models/Backup)
 - [BackupList](docs/models/BackupList)
 - [BackupListAllOf](docs/models/BackupListAllOf)
 - [BackupResponse](docs/models/BackupResponse)
 - [BaseBackup](docs/models/BaseBackup)
 - [ClusterList](docs/models/ClusterList)
 - [ClusterListAllOf](docs/models/ClusterListAllOf)
 - [ClusterMetadata](docs/models/ClusterMetadata)
 - [ClusterProperties](docs/models/ClusterProperties)
 - [ClusterResponse](docs/models/ClusterResponse)
 - [ClustersGet400Response](docs/models/ClustersGet400Response)
 - [ClustersGet401Response](docs/models/ClustersGet401Response)
 - [ClustersGet403Response](docs/models/ClustersGet403Response)
 - [ClustersGet404Response](docs/models/ClustersGet404Response)
 - [ClustersGet405Response](docs/models/ClustersGet405Response)
 - [ClustersGet415Response](docs/models/ClustersGet415Response)
 - [ClustersGet422Response](docs/models/ClustersGet422Response)
 - [ClustersGet429Response](docs/models/ClustersGet429Response)
 - [ClustersGet500Response](docs/models/ClustersGet500Response)
 - [ClustersGet503Response](docs/models/ClustersGet503Response)
 - [Connection](docs/models/Connection)
 - [CreateClusterProperties](docs/models/CreateClusterProperties)
 - [CreateClusterRequest](docs/models/CreateClusterRequest)
 - [DBUser](docs/models/DBUser)
 - [DayOfTheWeek](docs/models/DayOfTheWeek)
 - [ErrorMessage](docs/models/ErrorMessage)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [MariadbVersion](docs/models/MariadbVersion)
 - [Pagination](docs/models/Pagination)
 - [PaginationLinks](docs/models/PaginationLinks)
 - [PatchClusterProperties](docs/models/PatchClusterProperties)
 - [PatchClusterRequest](docs/models/PatchClusterRequest)
 - [RestoreRequest](docs/models/RestoreRequest)
 - [State](docs/models/State)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
