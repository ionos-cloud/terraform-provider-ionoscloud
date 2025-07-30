# Go API client for apigateway

API Gateway is an application that acts as a \"front door\" for backend services and APIs, handling client requests and routing them to the appropriate backend.


## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/apigateway.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/apigateway.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/apigateway@latest
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

- *https://apigateway.de-txl.ionos.com* - Production de-txl

By default, *https://apigateway.de-txl.ionos.com* is used, however this can be overriden at authentication, either
by setting the `IONOS_API_URL` environment variable or by specifying the `hostUrl` parameter when
initializing the sdk client.

### Basic Authentication

- **Type**: HTTP basic authentication

Example

```golang
import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	apigateway "github.com/ionos-cloud/sdk-go-bundle/products/apigateway"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := apigateway.NewAPIClient(cfg)
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
        apigateway "github.com/ionos-cloud/sdk-go-bundle/products/apigateway"
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
        apiClient := apigateway.NewAPIClient(cfg)
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
         apigateway "github.com/ionos-cloud/sdk-go-bundle/products/apigateway"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := apigateway.NewAPIClient(cfg)
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
        apigateway "github.com/ionos-cloud/sdk-go-bundle/products/apigateway"
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
    apiClient := apigateway.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://apigateway.de-txl.ionos.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
APIGatewaysApi | [**ApigatewaysDelete**](docs/api/APIGatewaysApi.md#apigatewaysdelete) | **Delete** /gateways/{apigatewayId} | Delete Gateway
APIGatewaysApi | [**ApigatewaysFindById**](docs/api/APIGatewaysApi.md#apigatewaysfindbyid) | **Get** /gateways/{apigatewayId} | Retrieve Gateway
APIGatewaysApi | [**ApigatewaysGet**](docs/api/APIGatewaysApi.md#apigatewaysget) | **Get** /gateways | Retrieve all Apigateways
APIGatewaysApi | [**ApigatewaysPost**](docs/api/APIGatewaysApi.md#apigatewayspost) | **Post** /gateways | Create Gateway
APIGatewaysApi | [**ApigatewaysPut**](docs/api/APIGatewaysApi.md#apigatewaysput) | **Put** /gateways/{apigatewayId} | Ensure Gateway
RoutesApi | [**ApigatewaysRoutesDelete**](docs/api/RoutesApi.md#apigatewaysroutesdelete) | **Delete** /gateways/{apigatewayId}/routes/{routeId} | Delete Route
RoutesApi | [**ApigatewaysRoutesFindById**](docs/api/RoutesApi.md#apigatewaysroutesfindbyid) | **Get** /gateways/{apigatewayId}/routes/{routeId} | Retrieve Route
RoutesApi | [**ApigatewaysRoutesGet**](docs/api/RoutesApi.md#apigatewaysroutesget) | **Get** /gateways/{apigatewayId}/routes | Retrieve all Routes
RoutesApi | [**ApigatewaysRoutesPost**](docs/api/RoutesApi.md#apigatewaysroutespost) | **Post** /gateways/{apigatewayId}/routes | Create Route
RoutesApi | [**ApigatewaysRoutesPut**](docs/api/RoutesApi.md#apigatewaysroutesput) | **Put** /gateways/{apigatewayId}/routes/{routeId} | Ensure Route

</details>

## Documentation For Models

All URIs are relative to *https://apigateway.de-txl.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [Error](docs/models/Error)
 - [ErrorMessages](docs/models/ErrorMessages)
 - [Gateway](docs/models/Gateway)
 - [GatewayCreate](docs/models/GatewayCreate)
 - [GatewayCustomDomains](docs/models/GatewayCustomDomains)
 - [GatewayEnsure](docs/models/GatewayEnsure)
 - [GatewayRead](docs/models/GatewayRead)
 - [GatewayReadList](docs/models/GatewayReadList)
 - [GatewayReadListAllOf](docs/models/GatewayReadListAllOf)
 - [Links](docs/models/Links)
 - [Metadata](docs/models/Metadata)
 - [MetadataWithEndpoint](docs/models/MetadataWithEndpoint)
 - [MetadataWithEndpointAllOf](docs/models/MetadataWithEndpointAllOf)
 - [MetadataWithStatus](docs/models/MetadataWithStatus)
 - [MetadataWithStatusAllOf](docs/models/MetadataWithStatusAllOf)
 - [Pagination](docs/models/Pagination)
 - [Route](docs/models/Route)
 - [RouteCreate](docs/models/RouteCreate)
 - [RouteEnsure](docs/models/RouteEnsure)
 - [RouteRead](docs/models/RouteRead)
 - [RouteReadList](docs/models/RouteReadList)
 - [RouteReadListAllOf](docs/models/RouteReadListAllOf)
 - [RouteUpstreams](docs/models/RouteUpstreams)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
