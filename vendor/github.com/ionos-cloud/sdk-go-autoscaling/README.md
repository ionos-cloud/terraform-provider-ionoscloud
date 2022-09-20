![CI](https://github.com/ionos-cloud/sdk-resources/workflows/%5B%20CI%20%5D%20CloudApi%20V6%20/%20Go/badge.svg)
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=alert_status)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=bugs)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=security_rating)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ionos-cloud_sdk-go&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=ionos-cloud_sdk-go)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/sdk-go.svg)](https://github.com/ionos-cloud/sdk-go/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/sdk-go.svg)](https://github.com/ionos-cloud/sdk-go/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/sdk-go.svg)](https://github.com/ionos-cloud/sdk-go)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")


# Go API client for ionoscloud

Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their manage Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls.

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go/v6
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go/v6
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.
```bash
go get github.com/ionos-cloud/sdk-go/v6@v6.0.0
```
To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go/v6@latest
```

## Environment Variables

| Environment Variable | Description                                                                                                                                                                                                                    |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IONOS_USERNAME`     | Specify the username used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_PASSWORD`     | Specify the password used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_TOKEN`        | Specify the token used to login, if a token is being used instead of username and password                                                                                                                                     |
| `IONOS_API_URL`      | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v6` path, so it should _not_ be included in the `IONOS_API_URL` environment variable |
| `IONOS_LOGLEVEL`     | Specify the Log Level used to log messages. Possible values: Off, Debug, Trace |

⚠️ **_Note: To overwrite the api endpoint - `api.ionos.com`, the environment variable `$IONOS_API_URL` can be set, and used with `NewConfigurationFromEnv()` function._**

## Examples

Examples for creating resources using the Go SDK can be found [here](examples/)

## Authentication

### Basic Authentication

- **Type**: HTTP basic authentication

Example

```golang
import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func basicAuthExample() error {
	cfg := ionoscloud.NewConfiguration("username_here", "pwd_here", "", "")
	cfg.Debug = true
	apiClient := ionoscloud.NewAPIClient(cfg)
	datacenters, _, err := apiClient.DataCentersApi.DatacentersGet(context.Background()).Depth(1).Execute()
	if err != nil {
		return fmt.Errorf("error retrieving datacenters %w", err)
	}
	if datacenters.HasItems() {
		for _, dc := range *datacenters.GetItems() {
			if dc.HasProperties() && dc.GetProperties().HasName() {
				fmt.Println(*dc.GetProperties().GetName())
			}
		}
	}
	return nil
}
```
### Token Authentication
There are 2 ways to generate your token:

 ### Generate token using [sdk-go-auth](https://github.com/ionos-cloud/sdk-go-auth):
```golang
    import (
        "context"
        "fmt"
        authApi "github.com/ionos-cloud/sdk-go-auth"
        "github.com/ionos-cloud/sdk-go/v6"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_USERNAME and IONOS_PASSWORD as env variables
        authClient := authApi.NewAPIClient(authApi.NewConfigurationFromEnv())
        jwt, _, err := authClient.TokensApi.TokensGenerate(context.Background()).Execute()
        if err != nil {
            return fmt.Errorf("error occurred while generating token (%w)", err)
        }
        if !jwt.HasToken() {
            return fmt.Errorf("could not generate token")
        }
        cfg := ionoscloud.NewConfiguration("", "", *jwt.GetToken(), "")
        cfg.Debug = true
        apiClient := ionoscloud.NewAPIClient(cfg)
        datacenters, _, err := apiClient.DataCentersApi.DatacenterGet(context.Background()).Depth(1).Execute()
        if err != nil {
            return fmt.Errorf("error retrieving datacenters (%w)", err)
        }
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
        "github.com/ionos-cloud/sdk-go/v6"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := authApi.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.Debug = true
        apiClient := ionoscloud.NewAPIClient(cfg)
        datacenters, _, err := apiClient.DataCenter6Api.DatacentersGet(context.Background()).Depth(1).Execute()
        if err != nil {
            return fmt.Errorf("error retrieving datacenters (%w)", err)
        }
        return nil
    }
```

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


#### How to set Depth parameter:

⚠️ **_Please use this parameter with caution. We recommend using the default value and raising its value only if it is needed._**

* On the configuration level:
```go
configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "URL")
configuration.SetDepth(5)
```
Using this method, the depth parameter will be set **on all the API calls**.

*  When calling a method:
```go
request := apiClient.DataCenterApi.DatacentersGet(context.Background()).Depth(1)
```
Using this method, the depth parameter will be set **on the current API call**.

* Using the default value:

If the depth parameter is not set, it will have the default value from the API that can be found [here](https://api.ionos.com/cloudapi/v6/swagger.json).

> Note: The priority for setting the depth parameter is: *set on function call > set on configuration level > set using the default value from the API*

### Pretty

The operations will also accept an optional _pretty_ argument. Setting this to a value of `true` or `false` controls whether the response is pretty-printed \(with indentation and new lines\). By default, the SDK sets the _pretty_ argument to `true`.

### Changing the base URL

Base URL for the HTTP operation can be changed by using the following function:

```go
requestProperties.SetURL("https://api.ionos.com/cloudapi/v6")
```

## Debugging

You can now inject any logger that implements Printf as a logger
instead of using the default sdk logger.
There are now Loglevels that you can set: `Off`, `Debug` and `Trace`.
`Off` - does not show any logs
`Debug` - regular logs, no sensitive information
`Trace` - we recommend you only set this field for debugging purposes. Disable it in your production environments because it can log sensitive data.
          It logs the full request and response without encryption, even for an HTTPS call. Verbose request and response logging can also significantly impact your application's performance.


```golang
package main
import "github.com/ionos-cloud/sdk-go/v6"
import "github.com/sirupsen/logrus"
func main() {
    // create your configuration. replace username, password, token and url with correct values, or use NewConfigurationFromEnv()
    // if you have set your env variables as explained above
    cfg := ionoscloud.NewConfiguration("username", "password", "token", "hostUrl")
    // enable request and response logging. this is the most verbose loglevel
    cfg.LogLevel = Trace
    // inject your own logger that implements Printf
    cfg.Logger = logrus.New()
    // create you api client with the configuration
    apiClient := ionoscloud.NewAPIClient(cfg)
}
```

If you want to see the API call request and response messages, you need to set the Debug field in the Configuration struct:

⚠️ **_Note: the field `Debug` is now deprecated and will be replaced with `LogLevel` in the future.

```golang
package main
import "github.com/ionos-cloud/sdk-go/v6"
func main() {
    // create your configuration. replace username, password, token and url with correct values, or use NewConfigurationFromEnv()
    // if you have set your env variables as explained above
    cfg := ionoscloud.NewConfiguration("username", "password", "token", "hostUrl")
    // enable request and response logging
    cfg.Debug = true
    // create you api client with the configuration
    apiClient := ionoscloud.NewAPIClient(cfg)
}
```

⚠️ **_Note: We recommend you only set this field for debugging purposes.
Disable it in your production environments because it can log sensitive data.
It logs the full request and response without encryption, even for an HTTPS call.
Verbose request and response logging can also significantly impact your application's performance._**


## Documentation for API Endpoints

All URIs are relative to *https://api.paas-public.k8s.stg.profitbricks.net/cloudapi/v6/containerregistries*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
LocationsApi | [**LocationsGet**](docs/api/LocationsApi.md#locationsget) | **Get** /locations | Get container registry locations
NamesApi | [**NamesFindByName**](docs/api/NamesApi.md#namesfindbyname) | **Head** /names/{name} | Get container registry name availability
RegistriesApi | [**RegistriesDelete**](docs/api/RegistriesApi.md#registriesdelete) | **Delete** /registries/{registryId} | Delete registry
RegistriesApi | [**RegistriesFindById**](docs/api/RegistriesApi.md#registriesfindbyid) | **Get** /registries/{registryId} | Get registry
RegistriesApi | [**RegistriesGet**](docs/api/RegistriesApi.md#registriesget) | **Get** /registries | List all container registries
RegistriesApi | [**RegistriesPatch**](docs/api/RegistriesApi.md#registriespatch) | **Patch** /registries/{registryId} | Update the properties of a registry
RegistriesApi | [**RegistriesPost**](docs/api/RegistriesApi.md#registriespost) | **Post** /registries | Create container registry
RegistriesApi | [**RegistriesPut**](docs/api/RegistriesApi.md#registriesput) | **Put** /registries/{registryId} | Create or replace container registry
RepositoriesApi | [**RegistriesRepositoriesDelete**](docs/api/RepositoriesApi.md#registriesrepositoriesdelete) | **Delete** /registries/{registryId}/repositories/{name} | Delete repository
TokensApi | [**RegistriesTokensDelete**](docs/api/TokensApi.md#registriestokensdelete) | **Delete** /registries/{registryId}/tokens/{tokenId} | Delete token
TokensApi | [**RegistriesTokensFindById**](docs/api/TokensApi.md#registriestokensfindbyid) | **Get** /registries/{registryId}/tokens/{tokenId} | Get Token Information
TokensApi | [**RegistriesTokensGet**](docs/api/TokensApi.md#registriestokensget) | **Get** /registries/{registryId}/tokens | List all tokens for the container registry
TokensApi | [**RegistriesTokensPatch**](docs/api/TokensApi.md#registriestokenspatch) | **Patch** /registries/{registryId}/tokens/{tokenId} | Update token
TokensApi | [**RegistriesTokensPost**](docs/api/TokensApi.md#registriestokenspost) | **Post** /registries/{registryId}/tokens | Create token
TokensApi | [**RegistriesTokensPut**](docs/api/TokensApi.md#registriestokensput) | **Put** /registries/{registryId}/tokens/{tokenId} | Create or replace token

</details>

## Documentation For Models

All URIs are relative to *https://api.paas-public.k8s.stg.profitbricks.net/cloudapi/v6/containerregistries*
<details >
<summary title="Click to toggle">API models list</summary>

 - [ApiErrorMessage](docs/models/ApiErrorMessage)
 - [ApiErrorResponse](docs/models/ApiErrorResponse)
 - [ApiResourceMetadata](docs/models/ApiResourceMetadata)
 - [Credentials](docs/models/Credentials)
 - [Day](docs/models/Day)
 - [Location](docs/models/Location)
 - [LocationsResponse](docs/models/LocationsResponse)
 - [PaginationLinks](docs/models/PaginationLinks)
 - [PatchRegistryInput](docs/models/PatchRegistryInput)
 - [PatchTokenInput](docs/models/PatchTokenInput)
 - [PostRegistryInput](docs/models/PostRegistryInput)
 - [PostRegistryOutput](docs/models/PostRegistryOutput)
 - [PostRegistryProperties](docs/models/PostRegistryProperties)
 - [PostTokenInput](docs/models/PostTokenInput)
 - [PostTokenOutput](docs/models/PostTokenOutput)
 - [PostTokenProperties](docs/models/PostTokenProperties)
 - [PutRegistryInput](docs/models/PutRegistryInput)
 - [PutRegistryOutput](docs/models/PutRegistryOutput)
 - [PutTokenInput](docs/models/PutTokenInput)
 - [PutTokenOutput](docs/models/PutTokenOutput)
 - [RegistriesResponse](docs/models/RegistriesResponse)
 - [RegistryProperties](docs/models/RegistryProperties)
 - [RegistryResponse](docs/models/RegistryResponse)
 - [Scope](docs/models/Scope)
 - [StorageUsage](docs/models/StorageUsage)
 - [TokenProperties](docs/models/TokenProperties)
 - [TokenResponse](docs/models/TokenResponse)
 - [TokensResponse](docs/models/TokensResponse)
 - [WeeklySchedule](docs/models/WeeklySchedule)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>



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