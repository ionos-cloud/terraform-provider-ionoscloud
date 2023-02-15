# Go API client for containerregistry

Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls.

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/containerregistry.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/containerregistry.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/containerregistry@latest
```

## Environment Variables

| Environment Variable | Description                                                                                                                                                                                                                    |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IONOS_USERNAME`     | Specify the username used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_PASSWORD`     | Specify the password used to login, to authenticate against the IONOS Cloud API                                                                                                                                                |
| `IONOS_TOKEN`        | Specify the token used to login, if a token is being used instead of username and password                                                                                                                                     |
| `IONOS_API_URL`      | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v6` path, so it should _not_ be included in the `IONOS_API_URL` environment variable |
| `IONOS_LOGLEVEL`     | Specify the Log Level used to log messages. Possible values: Off, Debug, Trace |
| `IONOS_PINNED_CERT`  | Specify the SHA-256 public fingerprint here, enables certificate pinning                                                                                                                                                       |

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
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	containerregistry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "")
	cfg.LogLevel = Trace
	apiClient := containerregistry.NewAPIClient(cfg)
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
        containerregistry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
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
        cfg := shared.NewConfiguration("", "", *jwt.GetToken(), "")
        cfg.LogLevel = Trace
        apiClient := containerregistry.NewAPIClient(cfg)
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
         containerregistry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := containerregistry.NewAPIClient(cfg)
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

You can now inject any logger that implements Printf as a logger
instead of using the default sdk logger.
There are now Loglevels that you can set: `Off`, `Debug` and `Trace`.
`Off` - does not show any logs
`Debug` - regular logs, no sensitive information
`Trace` - we recommend you only set this field for debugging purposes. Disable it in your production environments because it can log sensitive data.
          It logs the full request and response without encryption, even for an HTTPS call. Verbose request and response logging can also significantly impact your application's performance.


```golang
package main

    import (
        containerregistry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
        "github.com/ionos-cloud/sdk-go-bundle/shared"
        "github.com/sirupsen/logrus"
    )

func main() {
    // create your configuration. replace username, password, token and url with correct values, or use NewConfigurationFromEnv()
    // if you have set your env variables as explained above
    cfg := shared.NewConfiguration("username", "password", "token", "hostUrl")
    // enable request and response logging. this is the most verbose loglevel
    cfg.LogLevel = Trace
    // inject your own logger that implements Printf
    cfg.Logger = logrus.New()
    // create you api client with the configuration
    apiClient := containerregistry.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/containerregistries*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
LocationsApi | [**LocationsGet**](docs/api/LocationsApi.md#locationsget) | **Get** /locations | Get container registry locations
NamesApi | [**NamesCheckUsage**](docs/api/NamesApi.md#namescheckusage) | **Head** /names/{name} | Get container registry name availability
RegistriesApi | [**RegistriesDelete**](docs/api/RegistriesApi.md#registriesdelete) | **Delete** /registries/{registryId} | Delete registry
RegistriesApi | [**RegistriesFindById**](docs/api/RegistriesApi.md#registriesfindbyid) | **Get** /registries/{registryId} | Get a registry
RegistriesApi | [**RegistriesGet**](docs/api/RegistriesApi.md#registriesget) | **Get** /registries | List all container registries
RegistriesApi | [**RegistriesPatch**](docs/api/RegistriesApi.md#registriespatch) | **Patch** /registries/{registryId} | Update the properties of a registry
RegistriesApi | [**RegistriesPost**](docs/api/RegistriesApi.md#registriespost) | **Post** /registries | Create container registry
RegistriesApi | [**RegistriesPut**](docs/api/RegistriesApi.md#registriesput) | **Put** /registries/{registryId} | Create or replace a container registry
RepositoriesApi | [**RegistriesRepositoriesDelete**](docs/api/RepositoriesApi.md#registriesrepositoriesdelete) | **Delete** /registries/{registryId}/repositories/{name} | Delete repository
TokensApi | [**RegistriesTokensDelete**](docs/api/TokensApi.md#registriestokensdelete) | **Delete** /registries/{registryId}/tokens/{tokenId} | Delete token
TokensApi | [**RegistriesTokensFindById**](docs/api/TokensApi.md#registriestokensfindbyid) | **Get** /registries/{registryId}/tokens/{tokenId} | Get token information
TokensApi | [**RegistriesTokensGet**](docs/api/TokensApi.md#registriestokensget) | **Get** /registries/{registryId}/tokens | List all tokens for the container registry
TokensApi | [**RegistriesTokensPatch**](docs/api/TokensApi.md#registriestokenspatch) | **Patch** /registries/{registryId}/tokens/{tokenId} | Update token
TokensApi | [**RegistriesTokensPost**](docs/api/TokensApi.md#registriestokenspost) | **Post** /registries/{registryId}/tokens | Create token
TokensApi | [**RegistriesTokensPut**](docs/api/TokensApi.md#registriestokensput) | **Put** /registries/{registryId}/tokens/{tokenId} | Create or replace token

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com/containerregistries*
<details >
<summary title="Click to toggle">API models list</summary>

 - [ApiErrorMessage](docs/models/ApiErrorMessage)
 - [ApiErrorResponse](docs/models/ApiErrorResponse)
 - [ApiResourceMetadata](docs/models/ApiResourceMetadata)
 - [Credentials](docs/models/Credentials)
 - [Day](docs/models/Day)
 - [Location](docs/models/Location)
 - [LocationsResponse](docs/models/LocationsResponse)
 - [Pagination](docs/models/Pagination)
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