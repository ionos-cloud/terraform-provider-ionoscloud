# Go API client for dns

Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.


## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/dns.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/dns.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/dns@latest
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

- *https://dns.de-fra.ionos.com* - Frankfurt

By default, *https://dns.de-fra.ionos.com* is used, however this can be overriden at authentication, either
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
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := dns.NewAPIClient(cfg)
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
        dns "github.com/ionos-cloud/sdk-go-bundle/products/dns"
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
        apiClient := dns.NewAPIClient(cfg)
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
         dns "github.com/ionos-cloud/sdk-go-bundle/products/dns"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := dns.NewAPIClient(cfg)
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
        dns "github.com/ionos-cloud/sdk-go-bundle/products/dns"
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
    apiClient := dns.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://dns.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
DNSSECApi | [**ZonesKeysDelete**](docs/api/DNSSECApi.md#zoneskeysdelete) | **Delete** /zones/{zoneId}/keys | Delete a DNSSEC key
DNSSECApi | [**ZonesKeysGet**](docs/api/DNSSECApi.md#zoneskeysget) | **Get** /zones/{zoneId}/keys | Retrieve a DNSSEC key
DNSSECApi | [**ZonesKeysPost**](docs/api/DNSSECApi.md#zoneskeyspost) | **Post** /zones/{zoneId}/keys | Create a DNSSEC key
QuotaApi | [**QuotaGet**](docs/api/QuotaApi.md#quotaget) | **Get** /quota | Retrieve resources quota
RecordsApi | [**RecordsGet**](docs/api/RecordsApi.md#recordsget) | **Get** /records | Retrieve all records from primary zones
RecordsApi | [**SecondaryzonesRecordsGet**](docs/api/RecordsApi.md#secondaryzonesrecordsget) | **Get** /secondaryzones/{secondaryZoneId}/records | Retrieve records for a secondary zone
RecordsApi | [**ZonesRecordsDelete**](docs/api/RecordsApi.md#zonesrecordsdelete) | **Delete** /zones/{zoneId}/records/{recordId} | Delete a record
RecordsApi | [**ZonesRecordsFindById**](docs/api/RecordsApi.md#zonesrecordsfindbyid) | **Get** /zones/{zoneId}/records/{recordId} | Retrieve a record
RecordsApi | [**ZonesRecordsGet**](docs/api/RecordsApi.md#zonesrecordsget) | **Get** /zones/{zoneId}/records | Retrieve records
RecordsApi | [**ZonesRecordsPost**](docs/api/RecordsApi.md#zonesrecordspost) | **Post** /zones/{zoneId}/records | Create a record
RecordsApi | [**ZonesRecordsPut**](docs/api/RecordsApi.md#zonesrecordsput) | **Put** /zones/{zoneId}/records/{recordId} | Update a record
ReverseRecordsApi | [**ReverserecordsDelete**](docs/api/ReverseRecordsApi.md#reverserecordsdelete) | **Delete** /reverserecords/{reverserecordId} | Delete a reverse DNS record
ReverseRecordsApi | [**ReverserecordsFindById**](docs/api/ReverseRecordsApi.md#reverserecordsfindbyid) | **Get** /reverserecords/{reverserecordId} | Retrieve a reverse DNS record
ReverseRecordsApi | [**ReverserecordsGet**](docs/api/ReverseRecordsApi.md#reverserecordsget) | **Get** /reverserecords | Retrieves existing reverse DNS records
ReverseRecordsApi | [**ReverserecordsPost**](docs/api/ReverseRecordsApi.md#reverserecordspost) | **Post** /reverserecords | Create a reverse DNS record
ReverseRecordsApi | [**ReverserecordsPut**](docs/api/ReverseRecordsApi.md#reverserecordsput) | **Put** /reverserecords/{reverserecordId} | Update a reverse DNS record
SecondaryZonesApi | [**SecondaryzonesAxfrGet**](docs/api/SecondaryZonesApi.md#secondaryzonesaxfrget) | **Get** /secondaryzones/{secondaryZoneId}/axfr | Get status of zone transfer
SecondaryZonesApi | [**SecondaryzonesAxfrPut**](docs/api/SecondaryZonesApi.md#secondaryzonesaxfrput) | **Put** /secondaryzones/{secondaryZoneId}/axfr | Start zone transfer
SecondaryZonesApi | [**SecondaryzonesDelete**](docs/api/SecondaryZonesApi.md#secondaryzonesdelete) | **Delete** /secondaryzones/{secondaryZoneId} | Delete a secondary zone
SecondaryZonesApi | [**SecondaryzonesFindById**](docs/api/SecondaryZonesApi.md#secondaryzonesfindbyid) | **Get** /secondaryzones/{secondaryZoneId} | Retrieve a secondary zone
SecondaryZonesApi | [**SecondaryzonesGet**](docs/api/SecondaryZonesApi.md#secondaryzonesget) | **Get** /secondaryzones | Retrieve secondary zones
SecondaryZonesApi | [**SecondaryzonesPost**](docs/api/SecondaryZonesApi.md#secondaryzonespost) | **Post** /secondaryzones | Create a secondary zone
SecondaryZonesApi | [**SecondaryzonesPut**](docs/api/SecondaryZonesApi.md#secondaryzonesput) | **Put** /secondaryzones/{secondaryZoneId} | Update a secondary zone
ZoneFilesApi | [**ZonesZonefileGet**](docs/api/ZoneFilesApi.md#zoneszonefileget) | **Get** /zones/{zoneId}/zonefile | Retrieve a zone file
ZoneFilesApi | [**ZonesZonefilePut**](docs/api/ZoneFilesApi.md#zoneszonefileput) | **Put** /zones/{zoneId}/zonefile | Updates a zone with a file
ZonesApi | [**ZonesDelete**](docs/api/ZonesApi.md#zonesdelete) | **Delete** /zones/{zoneId} | Delete a zone
ZonesApi | [**ZonesFindById**](docs/api/ZonesApi.md#zonesfindbyid) | **Get** /zones/{zoneId} | Retrieve a zone
ZonesApi | [**ZonesGet**](docs/api/ZonesApi.md#zonesget) | **Get** /zones | Retrieve zones
ZonesApi | [**ZonesPost**](docs/api/ZonesApi.md#zonespost) | **Post** /zones | Create a zone
ZonesApi | [**ZonesPut**](docs/api/ZonesApi.md#zonesput) | **Put** /zones/{zoneId} | Update a zone

</details>

## Documentation For Models

All URIs are relative to *https://dns.de-fra.ionos.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [Algorithm](docs/models/Algorithm)
 - [CommonZone](docs/models/CommonZone)
 - [CommonZoneRead](docs/models/CommonZoneRead)
 - [CommonZoneReadList](docs/models/CommonZoneReadList)
 - [DnssecKey](docs/models/DnssecKey)
 - [DnssecKeyCreate](docs/models/DnssecKeyCreate)
 - [DnssecKeyParameters](docs/models/DnssecKeyParameters)
 - [DnssecKeyReadCreation](docs/models/DnssecKeyReadCreation)
 - [DnssecKeyReadList](docs/models/DnssecKeyReadList)
 - [DnssecKeyReadListMetadata](docs/models/DnssecKeyReadListMetadata)
 - [DnssecKeyReadListProperties](docs/models/DnssecKeyReadListProperties)
 - [DnssecKeyReadListPropertiesKeyParameters](docs/models/DnssecKeyReadListPropertiesKeyParameters)
 - [DnssecKeyReadListPropertiesNsecParameters](docs/models/DnssecKeyReadListPropertiesNsecParameters)
 - [Error](docs/models/Error)
 - [ErrorMessages](docs/models/ErrorMessages)
 - [KeyData](docs/models/KeyData)
 - [KeyParameters](docs/models/KeyParameters)
 - [KskBits](docs/models/KskBits)
 - [Links](docs/models/Links)
 - [Metadata](docs/models/Metadata)
 - [MetadataForSecondaryZoneRecords](docs/models/MetadataForSecondaryZoneRecords)
 - [MetadataWithStateFqdnZoneId](docs/models/MetadataWithStateFqdnZoneId)
 - [MetadataWithStateFqdnZoneIdAllOf](docs/models/MetadataWithStateFqdnZoneIdAllOf)
 - [MetadataWithStateNameservers](docs/models/MetadataWithStateNameservers)
 - [MetadataWithStateNameserversAllOf](docs/models/MetadataWithStateNameserversAllOf)
 - [NsecMode](docs/models/NsecMode)
 - [NsecParameters](docs/models/NsecParameters)
 - [ProvisioningState](docs/models/ProvisioningState)
 - [Quota](docs/models/Quota)
 - [QuotaDetail](docs/models/QuotaDetail)
 - [Record](docs/models/Record)
 - [RecordCreate](docs/models/RecordCreate)
 - [RecordEnsure](docs/models/RecordEnsure)
 - [RecordRead](docs/models/RecordRead)
 - [RecordReadList](docs/models/RecordReadList)
 - [ReverseRecord](docs/models/ReverseRecord)
 - [ReverseRecordCreate](docs/models/ReverseRecordCreate)
 - [ReverseRecordEnsure](docs/models/ReverseRecordEnsure)
 - [ReverseRecordRead](docs/models/ReverseRecordRead)
 - [ReverseRecordsReadList](docs/models/ReverseRecordsReadList)
 - [SecondaryZone](docs/models/SecondaryZone)
 - [SecondaryZoneAllOf](docs/models/SecondaryZoneAllOf)
 - [SecondaryZoneCreate](docs/models/SecondaryZoneCreate)
 - [SecondaryZoneEnsure](docs/models/SecondaryZoneEnsure)
 - [SecondaryZoneRead](docs/models/SecondaryZoneRead)
 - [SecondaryZoneReadAllOf](docs/models/SecondaryZoneReadAllOf)
 - [SecondaryZoneReadList](docs/models/SecondaryZoneReadList)
 - [SecondaryZoneReadListAllOf](docs/models/SecondaryZoneReadListAllOf)
 - [SecondaryZoneRecordRead](docs/models/SecondaryZoneRecordRead)
 - [SecondaryZoneRecordReadList](docs/models/SecondaryZoneRecordReadList)
 - [SecondaryZoneRecordReadListMetadata](docs/models/SecondaryZoneRecordReadListMetadata)
 - [Zone](docs/models/Zone)
 - [ZoneAllOf](docs/models/ZoneAllOf)
 - [ZoneCreate](docs/models/ZoneCreate)
 - [ZoneEnsure](docs/models/ZoneEnsure)
 - [ZoneRead](docs/models/ZoneRead)
 - [ZoneReadAllOf](docs/models/ZoneReadAllOf)
 - [ZoneReadList](docs/models/ZoneReadList)
 - [ZoneReadListAllOf](docs/models/ZoneReadListAllOf)
 - [ZoneTransferPrimaryIpStatus](docs/models/ZoneTransferPrimaryIpStatus)
 - [ZoneTransferPrimaryIpsStatus](docs/models/ZoneTransferPrimaryIpsStatus)
 - [ZskBits](docs/models/ZskBits)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
