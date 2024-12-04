# Go API client for vpn

The Managed VPN Gateway service provides secure and scalable connectivity, enabling encrypted communication between your IONOS cloud resources in a VDC and remote networks (on-premises, multi-cloud, private LANs in other VDCs etc).

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/vpn.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/vpn.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/vpn@latest
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

- *https://vpn.de-fra.ionos.com* - Production de-fra
- *https://vpn.de-txl.ionos.com* - Production de-txl
- *https://vpn.es-vit.ionos.com* - Production es-vit
- *https://vpn.gb-bhx.ionos.com* - Production gb-bhx
- *https://vpn.gb-lhr.ionos.com* - Production gb-lhr
- *https://vpn.us-ewr.ionos.com* - Production us-ewr
- *https://vpn.us-las.ionos.com* - Production us-las
- *https://vpn.us-mci.ionos.com* - Production us-mci
- *https://vpn.fr-par.ionos.com* - Production fr-par

By default, *https://vpn.de-fra.ionos.com* is used, however this can be overriden at authentication, either
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
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := vpn.NewAPIClient(cfg)
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
        vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn"
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
        apiClient := vpn.NewAPIClient(cfg)
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
         vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := vpn.NewAPIClient(cfg)
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
        vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn"
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
    apiClient := vpn.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://vpn.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
IPSecGatewaysApi | [**IpsecgatewaysDelete**](docs/api/IPSecGatewaysApi.md#ipsecgatewaysdelete) | **Delete** /ipsecgateways/{gatewayId} | Delete IPSecGateway
IPSecGatewaysApi | [**IpsecgatewaysFindById**](docs/api/IPSecGatewaysApi.md#ipsecgatewaysfindbyid) | **Get** /ipsecgateways/{gatewayId} | Retrieve IPSecGateway
IPSecGatewaysApi | [**IpsecgatewaysGet**](docs/api/IPSecGatewaysApi.md#ipsecgatewaysget) | **Get** /ipsecgateways | Retrieve all IPSecGateways
IPSecGatewaysApi | [**IpsecgatewaysPost**](docs/api/IPSecGatewaysApi.md#ipsecgatewayspost) | **Post** /ipsecgateways | Create IPSecGateway
IPSecGatewaysApi | [**IpsecgatewaysPut**](docs/api/IPSecGatewaysApi.md#ipsecgatewaysput) | **Put** /ipsecgateways/{gatewayId} | Ensure IPSecGateway
IPSecTunnelsApi | [**IpsecgatewaysTunnelsDelete**](docs/api/IPSecTunnelsApi.md#ipsecgatewaystunnelsdelete) | **Delete** /ipsecgateways/{gatewayId}/tunnels/{tunnelId} | Delete IPSecTunnel
IPSecTunnelsApi | [**IpsecgatewaysTunnelsFindById**](docs/api/IPSecTunnelsApi.md#ipsecgatewaystunnelsfindbyid) | **Get** /ipsecgateways/{gatewayId}/tunnels/{tunnelId} | Retrieve IPSecTunnel
IPSecTunnelsApi | [**IpsecgatewaysTunnelsGet**](docs/api/IPSecTunnelsApi.md#ipsecgatewaystunnelsget) | **Get** /ipsecgateways/{gatewayId}/tunnels | Retrieve all IPSecTunnels
IPSecTunnelsApi | [**IpsecgatewaysTunnelsPost**](docs/api/IPSecTunnelsApi.md#ipsecgatewaystunnelspost) | **Post** /ipsecgateways/{gatewayId}/tunnels | Create IPSecTunnel
IPSecTunnelsApi | [**IpsecgatewaysTunnelsPut**](docs/api/IPSecTunnelsApi.md#ipsecgatewaystunnelsput) | **Put** /ipsecgateways/{gatewayId}/tunnels/{tunnelId} | Ensure IPSecTunnel
WireguardGatewaysApi | [**WireguardgatewaysDelete**](docs/api/WireguardGatewaysApi.md#wireguardgatewaysdelete) | **Delete** /wireguardgateways/{gatewayId} | Delete WireguardGateway
WireguardGatewaysApi | [**WireguardgatewaysFindById**](docs/api/WireguardGatewaysApi.md#wireguardgatewaysfindbyid) | **Get** /wireguardgateways/{gatewayId} | Retrieve WireguardGateway
WireguardGatewaysApi | [**WireguardgatewaysGet**](docs/api/WireguardGatewaysApi.md#wireguardgatewaysget) | **Get** /wireguardgateways | Retrieve all WireguardGateways
WireguardGatewaysApi | [**WireguardgatewaysPost**](docs/api/WireguardGatewaysApi.md#wireguardgatewayspost) | **Post** /wireguardgateways | Create WireguardGateway
WireguardGatewaysApi | [**WireguardgatewaysPut**](docs/api/WireguardGatewaysApi.md#wireguardgatewaysput) | **Put** /wireguardgateways/{gatewayId} | Ensure WireguardGateway
WireguardPeersApi | [**WireguardgatewaysPeersDelete**](docs/api/WireguardPeersApi.md#wireguardgatewayspeersdelete) | **Delete** /wireguardgateways/{gatewayId}/peers/{peerId} | Delete WireguardPeer
WireguardPeersApi | [**WireguardgatewaysPeersFindById**](docs/api/WireguardPeersApi.md#wireguardgatewayspeersfindbyid) | **Get** /wireguardgateways/{gatewayId}/peers/{peerId} | Retrieve WireguardPeer
WireguardPeersApi | [**WireguardgatewaysPeersGet**](docs/api/WireguardPeersApi.md#wireguardgatewayspeersget) | **Get** /wireguardgateways/{gatewayId}/peers | Retrieve all WireguardPeers
WireguardPeersApi | [**WireguardgatewaysPeersPost**](docs/api/WireguardPeersApi.md#wireguardgatewayspeerspost) | **Post** /wireguardgateways/{gatewayId}/peers | Create WireguardPeer
WireguardPeersApi | [**WireguardgatewaysPeersPut**](docs/api/WireguardPeersApi.md#wireguardgatewayspeersput) | **Put** /wireguardgateways/{gatewayId}/peers/{peerId} | Ensure WireguardPeer

</details>

## Documentation For Models

All URIs are relative to *https://vpn.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [Connection](docs/models/Connection)
 - [DayOfTheWeek](docs/models/DayOfTheWeek)
 - [ESPEncryption](docs/models/ESPEncryption)
 - [Error](docs/models/Error)
 - [ErrorMessages](docs/models/ErrorMessages)
 - [IKEEncryption](docs/models/IKEEncryption)
 - [IPSecGateway](docs/models/IPSecGateway)
 - [IPSecGatewayCreate](docs/models/IPSecGatewayCreate)
 - [IPSecGatewayEnsure](docs/models/IPSecGatewayEnsure)
 - [IPSecGatewayMetadata](docs/models/IPSecGatewayMetadata)
 - [IPSecGatewayRead](docs/models/IPSecGatewayRead)
 - [IPSecGatewayReadList](docs/models/IPSecGatewayReadList)
 - [IPSecGatewayReadListAllOf](docs/models/IPSecGatewayReadListAllOf)
 - [IPSecPSK](docs/models/IPSecPSK)
 - [IPSecTunnel](docs/models/IPSecTunnel)
 - [IPSecTunnelAuth](docs/models/IPSecTunnelAuth)
 - [IPSecTunnelCreate](docs/models/IPSecTunnelCreate)
 - [IPSecTunnelEnsure](docs/models/IPSecTunnelEnsure)
 - [IPSecTunnelMetadata](docs/models/IPSecTunnelMetadata)
 - [IPSecTunnelRead](docs/models/IPSecTunnelRead)
 - [IPSecTunnelReadList](docs/models/IPSecTunnelReadList)
 - [IPSecTunnelReadListAllOf](docs/models/IPSecTunnelReadListAllOf)
 - [Links](docs/models/Links)
 - [MaintenanceWindow](docs/models/MaintenanceWindow)
 - [Metadata](docs/models/Metadata)
 - [Pagination](docs/models/Pagination)
 - [ResourceStatus](docs/models/ResourceStatus)
 - [WireguardEndpoint](docs/models/WireguardEndpoint)
 - [WireguardGateway](docs/models/WireguardGateway)
 - [WireguardGatewayCreate](docs/models/WireguardGatewayCreate)
 - [WireguardGatewayEnsure](docs/models/WireguardGatewayEnsure)
 - [WireguardGatewayMetadata](docs/models/WireguardGatewayMetadata)
 - [WireguardGatewayMetadataAllOf](docs/models/WireguardGatewayMetadataAllOf)
 - [WireguardGatewayRead](docs/models/WireguardGatewayRead)
 - [WireguardGatewayReadList](docs/models/WireguardGatewayReadList)
 - [WireguardGatewayReadListAllOf](docs/models/WireguardGatewayReadListAllOf)
 - [WireguardPeer](docs/models/WireguardPeer)
 - [WireguardPeerCreate](docs/models/WireguardPeerCreate)
 - [WireguardPeerEnsure](docs/models/WireguardPeerEnsure)
 - [WireguardPeerMetadata](docs/models/WireguardPeerMetadata)
 - [WireguardPeerRead](docs/models/WireguardPeerRead)
 - [WireguardPeerReadList](docs/models/WireguardPeerReadList)
 - [WireguardPeerReadListAllOf](docs/models/WireguardPeerReadListAllOf)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
