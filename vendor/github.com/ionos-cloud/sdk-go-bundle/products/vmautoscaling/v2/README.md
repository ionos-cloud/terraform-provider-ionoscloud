# Go API client for vmautoscaling

The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.

For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.


When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling@latest
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

- *https://api.ionos.com/autoscaling* - Production

By default, *https://api.ionos.com/autoscaling* is used, however this can be overriden at authentication, either
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
	vmautoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := vmautoscaling.NewAPIClient(cfg)
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
        vmautoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling"
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
        apiClient := vmautoscaling.NewAPIClient(cfg)
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
         vmautoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := vmautoscaling.NewAPIClient(cfg)
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
        vmautoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling"
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
    apiClient := vmautoscaling.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com/autoscaling*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
AutoScalingGroupsApi | [**GroupsActionsFindById**](docs/api/AutoScalingGroupsApi.md#groupsactionsfindbyid) | **Get** /groups/{groupId}/actions/{actionId} | Get Scaling Action Details by ID
AutoScalingGroupsApi | [**GroupsActionsGet**](docs/api/AutoScalingGroupsApi.md#groupsactionsget) | **Get** /groups/{groupId}/actions | Get Scaling Actions
AutoScalingGroupsApi | [**GroupsDelete**](docs/api/AutoScalingGroupsApi.md#groupsdelete) | **Delete** /groups/{groupId} | Delete a VM Auto Scaling Group by ID
AutoScalingGroupsApi | [**GroupsFindById**](docs/api/AutoScalingGroupsApi.md#groupsfindbyid) | **Get** /groups/{groupId} | Get an Auto Scaling by ID
AutoScalingGroupsApi | [**GroupsGet**](docs/api/AutoScalingGroupsApi.md#groupsget) | **Get** /groups | Get VM Auto Scaling Groups
AutoScalingGroupsApi | [**GroupsPost**](docs/api/AutoScalingGroupsApi.md#groupspost) | **Post** /groups | Create a VM Auto Scaling Group
AutoScalingGroupsApi | [**GroupsPut**](docs/api/AutoScalingGroupsApi.md#groupsput) | **Put** /groups/{groupId} | Update a VM Auto Scaling Group by ID
AutoScalingGroupsApi | [**GroupsServersFindById**](docs/api/AutoScalingGroupsApi.md#groupsserversfindbyid) | **Get** /groups/{groupId}/servers/{serverId} | Get VM Auto Scaling Group Server by ID
AutoScalingGroupsApi | [**GroupsServersGet**](docs/api/AutoScalingGroupsApi.md#groupsserversget) | **Get** /groups/{groupId}/servers | Get VM Auto Scaling Group Servers

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com/autoscaling*
<details >
<summary title="Click to toggle">API models list</summary>

 - [Action](docs/models/Action)
 - [ActionAmount](docs/models/ActionAmount)
 - [ActionCollection](docs/models/ActionCollection)
 - [ActionProperties](docs/models/ActionProperties)
 - [ActionResource](docs/models/ActionResource)
 - [ActionStatus](docs/models/ActionStatus)
 - [ActionType](docs/models/ActionType)
 - [ActionsLinkResource](docs/models/ActionsLinkResource)
 - [AvailabilityZone](docs/models/AvailabilityZone)
 - [BusType](docs/models/BusType)
 - [CpuFamily](docs/models/CpuFamily)
 - [DatacenterServer](docs/models/DatacenterServer)
 - [Error401](docs/models/Error401)
 - [Error401Message](docs/models/Error401Message)
 - [Error404](docs/models/Error404)
 - [Error404Message](docs/models/Error404Message)
 - [ErrorAuthorize](docs/models/ErrorAuthorize)
 - [ErrorAuthorizeMessage](docs/models/ErrorAuthorizeMessage)
 - [ErrorGroupValidate](docs/models/ErrorGroupValidate)
 - [ErrorGroupValidateMessage](docs/models/ErrorGroupValidateMessage)
 - [ErrorMessage](docs/models/ErrorMessage)
 - [ErrorMessageParse](docs/models/ErrorMessageParse)
 - [Group](docs/models/Group)
 - [GroupCollection](docs/models/GroupCollection)
 - [GroupEntities](docs/models/GroupEntities)
 - [GroupPolicy](docs/models/GroupPolicy)
 - [GroupPolicyScaleInAction](docs/models/GroupPolicyScaleInAction)
 - [GroupPolicyScaleOutAction](docs/models/GroupPolicyScaleOutAction)
 - [GroupPost](docs/models/GroupPost)
 - [GroupPostEntities](docs/models/GroupPostEntities)
 - [GroupPostResponse](docs/models/GroupPostResponse)
 - [GroupProperties](docs/models/GroupProperties)
 - [GroupPropertiesDatacenter](docs/models/GroupPropertiesDatacenter)
 - [GroupPut](docs/models/GroupPut)
 - [GroupPutProperties](docs/models/GroupPutProperties)
 - [GroupPutPropertiesDatacenter](docs/models/GroupPutPropertiesDatacenter)
 - [GroupResource](docs/models/GroupResource)
 - [Metadata](docs/models/Metadata)
 - [MetadataBasic](docs/models/MetadataBasic)
 - [MetadataState](docs/models/MetadataState)
 - [Metric](docs/models/Metric)
 - [NicFirewallRule](docs/models/NicFirewallRule)
 - [NicFlowLog](docs/models/NicFlowLog)
 - [ParseError](docs/models/ParseError)
 - [QueryUnit](docs/models/QueryUnit)
 - [ReplicaNic](docs/models/ReplicaNic)
 - [ReplicaPropertiesPost](docs/models/ReplicaPropertiesPost)
 - [ReplicaVolumePost](docs/models/ReplicaVolumePost)
 - [Server](docs/models/Server)
 - [ServerCollection](docs/models/ServerCollection)
 - [ServerProperties](docs/models/ServerProperties)
 - [ServersLinkResource](docs/models/ServersLinkResource)
 - [TargetGroup](docs/models/TargetGroup)
 - [TerminationPolicyType](docs/models/TerminationPolicyType)
 - [VolumeHwType](docs/models/VolumeHwType)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
