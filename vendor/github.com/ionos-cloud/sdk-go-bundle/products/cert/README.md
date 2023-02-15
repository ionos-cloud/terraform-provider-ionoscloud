# Go API client for cert

Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates with IONOS services and your internal connected resources. For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic.

The service provides the basic functions of uploading and deleting your certificates for this purpose.

## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/cert.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/cert.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/cert@latest
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
	cert "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "")
	cfg.LogLevel = Trace
	apiClient := cert.NewAPIClient(cfg)
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
        cert "github.com/ionos-cloud/sdk-go-bundle/products/cert"
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
        apiClient := cert.NewAPIClient(cfg)
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
         cert "github.com/ionos-cloud/sdk-go-bundle/products/cert"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := cert.NewAPIClient(cfg)
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
        cert "github.com/ionos-cloud/sdk-go-bundle/products/cert"
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
    apiClient := cert.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.ionos.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
CertificatesApi | [**CertificatesDelete**](docs/api/CertificatesApi.md#certificatesdelete) | **Delete** /certificatemanager/certificates/{certificateId} | Delete a Certificate by ID
CertificatesApi | [**CertificatesGet**](docs/api/CertificatesApi.md#certificatesget) | **Get** /certificatemanager/certificates | Get Certificates
CertificatesApi | [**CertificatesGetById**](docs/api/CertificatesApi.md#certificatesgetbyid) | **Get** /certificatemanager/certificates/{certificateId} | Get a Certificate by ID
CertificatesApi | [**CertificatesPatch**](docs/api/CertificatesApi.md#certificatespatch) | **Patch** /certificatemanager/certificates/{certificateId} | Update a Certificate Name by ID
CertificatesApi | [**CertificatesPost**](docs/api/CertificatesApi.md#certificatespost) | **Post** /certificatemanager/certificates | Add a New Certificate
InformationApi | [**GetInfo**](docs/api/InformationApi.md#getinfo) | **Get** /certificatemanager | Get the Service API Information
InformationApi | [**GetJsonOpenApiSpec**](docs/api/InformationApi.md#getjsonopenapispec) | **Get** /certificatemanager/openapi.json | Get the Open API Documentation JSON
InformationApi | [**GetYamlOpenApiSpec**](docs/api/InformationApi.md#getyamlopenapispec) | **Get** /certificatemanager/openapi.yaml | Get the Open API Documentation YAML

</details>

## Documentation For Models

All URIs are relative to *https://api.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [ApiInfoDto](docs/models/ApiInfoDto)
 - [CertificateCollectionDto](docs/models/CertificateCollectionDto)
 - [CertificateCollectionDtoLinks](docs/models/CertificateCollectionDtoLinks)
 - [CertificateDto](docs/models/CertificateDto)
 - [CertificatePatchDto](docs/models/CertificatePatchDto)
 - [CertificatePatchPropertiesDto](docs/models/CertificatePatchPropertiesDto)
 - [CertificatePostDto](docs/models/CertificatePostDto)
 - [CertificatePostPropertiesDto](docs/models/CertificatePostPropertiesDto)
 - [CertificatePropertiesDto](docs/models/CertificatePropertiesDto)
 - [ConfigPropertyDto](docs/models/ConfigPropertyDto)
 - [ResourceMetadataDto](docs/models/ResourceMetadataDto)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>