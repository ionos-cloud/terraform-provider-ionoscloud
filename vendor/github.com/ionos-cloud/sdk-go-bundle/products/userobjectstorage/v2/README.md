# Go API client for userobjectstorage

## Overview
The IONOS Object Storage API for user-owned buckets is a REST-based API that allows developers and applications to interact directly with
IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless
compatibility with existing tools and libraries tailored for S3 systems.

### API References
- [Object Storage Management API Reference](https://api.ionos.com/docs/s3-management/v1/) for managing Access Keys
- [Object Storage API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/)
- Object Storage API Reference for user-owned buckets - current document

### User documentation
[IONOS Object Storage User Guide](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage)
* [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets)
* [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility)
* [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)

## Endpoints for user-owned buckets
| Location | Region Name | Bucket Type | Endpoint |
| --- | --- | --- | --- |
| **Frankfurt, Germany** | **de** | User-owned | `https://s3.eu-central-1.ionoscloud.com`, <br/><br/>**s3 legacy endpoint:** `https://s3-de-central.profitbricks.com`  |
| **Berlin, Germany** | **eu-central-2** | User-owned | `https://s3.eu-central-2.ionoscloud.com` |
| **Logroño, Spain** | **eu-south-2** | User-owned | `https://s3.eu-south-2.ionoscloud.com` |

## Changelog
- **30.05.2024** Renaming to Storage Object API for user-owned buckets
- **25.09.2023** Storage object operation names are now used for headlines.
- **20.09.2023** Improved description for [HeadBucket](#tag/Basic-Operations/operation/HeadBucket) and [GetBucketLocation](#tag/Location/operation/GetBucketLocation).
- **13.09.2023** Improved description for [Bucket Policy-related operations](#tag/Policy/operation/PutBucketPolicy).
- **06.09.2023** Improved description for [Bucket ACL-related operations](#tag/ACL/operation/GetBucketAcl).
- **30.08.2023** Improved description for [Object Lock-related operations](#tag/Object-Lock/operation/GetObjectLockConfiguration).
- **24.07.2023** Improved description for [ListObjectsV2](#tag/Basic-Operations/operation/ListObjectsV2).
- **17.07.2023** Improved description for [ListBuckets](#tag/Basic-Operations/operation/ListBuckets).
- **07.07.2023** Improved description for [PutBucketReplication](#tag/Replication/operation/PutBucketReplication),
  [GetBucketReplication](#tag/Replication/operation/GetBucketReplication), [DeleteBucketReplication](#tag/Replication/operation/DeleteBucketReplication).
- **05.07.2023** Improved description for [PutBucketVersioning](#tag/Versioning/operation/PutBucketVersioning)
  and [GetBucketVersioning](#tag/Versioning/operation/GetBucketVersioning).
- **29.06.2023** Improved description for [PutBucketLifecycleConfiguration](#tag/Lifecycle/operation/PutBucketLifecycle).
- **19.04.2023** Improved description on how to use the encryption with IONOS Object Storage managed (SSE-S3) and customer managed keys (SSE-C)
  for [PutBucketEncryption](#tag/Encryption/operation/PutBucketEncryption) and [PutObject](#tag/Basic-Operations/operation/PutObject).


## Overview
The IONOS Cloud SDK for GO provides you with access to the IONOS Cloud API. The client library supports both simple and complex requests.
It is designed for developers who are building applications in GO . The SDK for GO wraps the IONOS Cloud API. All API operations are performed over SSL and authenticated using your IONOS Cloud portal credentials.
The API can be accessed within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Installing

### Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage.git
```
To update the SDK use go get -u to retrieve the latest version of the SDK.
```bash
go get -u github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage.git
```
### Go Modules

If you are using Go modules, your go get will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use @<tag> in your go get command.

To get the latest SDK repository, use @latest.
```bash
go get github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage@latest
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

- *https://s3.eu-central-1.ionoscloud.com* - URL for &#x60;de&#x60; (Frankfurt, Germany)
- *https://s3.de-central.profitbricks.com* - Legacy URL for &#x60;de&#x60; (Frankfurt, Germany)
- *https://s3.eu-central-2.ionoscloud.com* - URL for &#x60;eu-central-2&#x60; (Berlin, Germany)
- *https://s3.eu-south-2.ionoscloud.com* - URL for &#x60;eu-south-2&#x60; (Logroño, Spain)

By default, *https://s3.eu-central-1.ionoscloud.com* is used, however this can be overriden at authentication, either
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
	userobjectstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage"
	"log"
)

func basicAuthExample() error {
	cfg := shared.NewConfiguration("username_here", "pwd_here", "", "hostUrl_here")
	cfg.LogLevel = Trace
	apiClient := userobjectstorage.NewAPIClient(cfg)
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
        userobjectstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage"
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
        apiClient := userobjectstorage.NewAPIClient(cfg)
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
         userobjectstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage"
        "log"
    )

    func TokenAuthExample() error {
        //note: to use NewConfigurationFromEnv(), you need to previously set IONOS_TOKEN as env variables
        authClient := auth.NewAPIClient(authApi.NewConfigurationFromEnv())
        cfg.LogLevel = Trace
        apiClient := userobjectstorage.NewAPIClient(cfg)
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
        userobjectstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage"
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
    apiClient := userobjectstorage.NewAPIClient(cfg)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://s3.eu-central-1.ionoscloud.com*
<details >
<summary title="Click to toggle">API Endpoints table</summary>


Class | Method | HTTP request | Description
------------- | ------------- | ------------- | -------------
BucketsApi | [**CreateBucket**](docs/api/BucketsApi.md#createbucket) | **Put** /{Bucket} | CreateBucket
BucketsApi | [**DeleteBucket**](docs/api/BucketsApi.md#deletebucket) | **Delete** /{Bucket} | DeleteBucket
BucketsApi | [**GetBucketLocation**](docs/api/BucketsApi.md#getbucketlocation) | **Get** /{Bucket}?location | GetBucketLocation
BucketsApi | [**HeadBucket**](docs/api/BucketsApi.md#headbucket) | **Head** /{Bucket} | HeadBucket
BucketsApi | [**ListBuckets**](docs/api/BucketsApi.md#listbuckets) | **Get** / | ListBuckets
CORSApi | [**DeleteBucketCors**](docs/api/CORSApi.md#deletebucketcors) | **Delete** /{Bucket}?cors | DeleteBucketCors
CORSApi | [**GetBucketCors**](docs/api/CORSApi.md#getbucketcors) | **Get** /{Bucket}?cors | GetBucketCors
CORSApi | [**PutBucketCors**](docs/api/CORSApi.md#putbucketcors) | **Put** /{Bucket}?cors | PutBucketCors
EncryptionApi | [**DeleteBucketEncryption**](docs/api/EncryptionApi.md#deletebucketencryption) | **Delete** /{Bucket}?encryption | DeleteBucketEncryption
EncryptionApi | [**GetBucketEncryption**](docs/api/EncryptionApi.md#getbucketencryption) | **Get** /{Bucket}?encryption | GetBucketEncryption
EncryptionApi | [**PutBucketEncryption**](docs/api/EncryptionApi.md#putbucketencryption) | **Put** /{Bucket}?encryption | PutBucketEncryption
LifecycleApi | [**DeleteBucketLifecycle**](docs/api/LifecycleApi.md#deletebucketlifecycle) | **Delete** /{Bucket}?lifecycle | DeleteBucketLifecycle
LifecycleApi | [**GetBucketLifecycle**](docs/api/LifecycleApi.md#getbucketlifecycle) | **Get** /{Bucket}?lifecycle | GetBucketLifecycle
LifecycleApi | [**PutBucketLifecycle**](docs/api/LifecycleApi.md#putbucketlifecycle) | **Put** /{Bucket}?lifecycle | PutBucketLifecycle
LoggingApi | [**GetBucketLogging**](docs/api/LoggingApi.md#getbucketlogging) | **Get** /{Bucket}?logging | GetBucketLogging
LoggingApi | [**PutBucketLogging**](docs/api/LoggingApi.md#putbucketlogging) | **Put** /{Bucket}?logging | PutBucketLogging
ObjectLockApi | [**GetObjectLegalHold**](docs/api/ObjectLockApi.md#getobjectlegalhold) | **Get** /{Bucket}/{Key}?legal-hold | GetObjectLegalHold
ObjectLockApi | [**GetObjectLockConfiguration**](docs/api/ObjectLockApi.md#getobjectlockconfiguration) | **Get** /{Bucket}?object-lock | GetObjectLockConfiguration
ObjectLockApi | [**GetObjectRetention**](docs/api/ObjectLockApi.md#getobjectretention) | **Get** /{Bucket}/{Key}?retention | GetObjectRetention
ObjectLockApi | [**PutObjectLegalHold**](docs/api/ObjectLockApi.md#putobjectlegalhold) | **Put** /{Bucket}/{Key}?legal-hold | PutObjectLegalHold
ObjectLockApi | [**PutObjectLockConfiguration**](docs/api/ObjectLockApi.md#putobjectlockconfiguration) | **Put** /{Bucket}?object-lock | PutObjectLockConfiguration
ObjectLockApi | [**PutObjectRetention**](docs/api/ObjectLockApi.md#putobjectretention) | **Put** /{Bucket}/{Key}?retention | PutObjectRetention
ObjectsApi | [**CopyObject**](docs/api/ObjectsApi.md#copyobject) | **Put** /{Bucket}/{Key}?x-amz-copy-source | CopyObject
ObjectsApi | [**DeleteObject**](docs/api/ObjectsApi.md#deleteobject) | **Delete** /{Bucket}/{Key} | DeleteObject
ObjectsApi | [**DeleteObjects**](docs/api/ObjectsApi.md#deleteobjects) | **Post** /{Bucket}?delete | DeleteObjects
ObjectsApi | [**GetObject**](docs/api/ObjectsApi.md#getobject) | **Get** /{Bucket}/{Key} | GetObject
ObjectsApi | [**HeadObject**](docs/api/ObjectsApi.md#headobject) | **Head** /{Bucket}/{Key} | HeadObject
ObjectsApi | [**ListObjects**](docs/api/ObjectsApi.md#listobjects) | **Get** /{Bucket} | ListObjects
ObjectsApi | [**ListObjectsV2**](docs/api/ObjectsApi.md#listobjectsv2) | **Get** /{Bucket}?list-type&#x3D;2 | ListObjectsV2
ObjectsApi | [**OPTIONSObject**](docs/api/ObjectsApi.md#optionsobject) | **Options** /{Bucket} | OPTIONSObject
ObjectsApi | [**POSTObject**](docs/api/ObjectsApi.md#postobject) | **Post** /{Bucket}/{Key} | POSTObject
ObjectsApi | [**PutObject**](docs/api/ObjectsApi.md#putobject) | **Put** /{Bucket}/{Key} | PutObject
PolicyApi | [**DeleteBucketPolicy**](docs/api/PolicyApi.md#deletebucketpolicy) | **Delete** /{Bucket}?policy | DeleteBucketPolicy
PolicyApi | [**GetBucketPolicy**](docs/api/PolicyApi.md#getbucketpolicy) | **Get** /{Bucket}?policy | GetBucketPolicy
PolicyApi | [**GetBucketPolicyStatus**](docs/api/PolicyApi.md#getbucketpolicystatus) | **Get** /{Bucket}?policyStatus | GetBucketPolicyStatus
PolicyApi | [**PutBucketPolicy**](docs/api/PolicyApi.md#putbucketpolicy) | **Put** /{Bucket}?policy | PutBucketPolicy
PublicAccessBlockApi | [**DeletePublicAccessBlock**](docs/api/PublicAccessBlockApi.md#deletepublicaccessblock) | **Delete** /{Bucket}?publicAccessBlock | DeletePublicAccessBlock
PublicAccessBlockApi | [**GetPublicAccessBlock**](docs/api/PublicAccessBlockApi.md#getpublicaccessblock) | **Get** /{Bucket}?publicAccessBlock | GetPublicAccessBlock
PublicAccessBlockApi | [**PutPublicAccessBlock**](docs/api/PublicAccessBlockApi.md#putpublicaccessblock) | **Put** /{Bucket}?publicAccessBlock | PutPublicAccessBlock
ReplicationApi | [**GetBucketReplication**](docs/api/ReplicationApi.md#getbucketreplication) | **Get** /{Bucket}?replication | GetBucketReplication
TaggingApi | [**DeleteBucketTagging**](docs/api/TaggingApi.md#deletebuckettagging) | **Delete** /{Bucket}?tagging | DeleteBucketTagging
TaggingApi | [**DeleteObjectTagging**](docs/api/TaggingApi.md#deleteobjecttagging) | **Delete** /{Bucket}/{Key}?tagging | DeleteObjectTagging
TaggingApi | [**GetBucketTagging**](docs/api/TaggingApi.md#getbuckettagging) | **Get** /{Bucket}?tagging | GetBucketTagging
TaggingApi | [**GetObjectTagging**](docs/api/TaggingApi.md#getobjecttagging) | **Get** /{Bucket}/{Key}?tagging | GetObjectTagging
TaggingApi | [**PutBucketTagging**](docs/api/TaggingApi.md#putbuckettagging) | **Put** /{Bucket}?tagging | PutBucketTagging
TaggingApi | [**PutObjectTagging**](docs/api/TaggingApi.md#putobjecttagging) | **Put** /{Bucket}/{Key}?tagging | PutObjectTagging
UploadsApi | [**AbortMultipartUpload**](docs/api/UploadsApi.md#abortmultipartupload) | **Delete** /{Bucket}/{Key}?uploadId | AbortMultipartUpload
UploadsApi | [**CompleteMultipartUpload**](docs/api/UploadsApi.md#completemultipartupload) | **Post** /{Bucket}/{Key}?uploadId | CompleteMultipartUpload
UploadsApi | [**CreateMultipartUpload**](docs/api/UploadsApi.md#createmultipartupload) | **Post** /{Bucket}/{Key}?uploads | CreateMultipartUpload
UploadsApi | [**ListMultipartUploads**](docs/api/UploadsApi.md#listmultipartuploads) | **Get** /{Bucket}?uploads | ListMultipartUploads
UploadsApi | [**ListParts**](docs/api/UploadsApi.md#listparts) | **Get** /{Bucket}/{Key}?uploadId | ListParts
UploadsApi | [**UploadPart**](docs/api/UploadsApi.md#uploadpart) | **Put** /{Bucket}/{Key}?uploadId | UploadPart
UploadsApi | [**UploadPartCopy**](docs/api/UploadsApi.md#uploadpartcopy) | **Put** /{Bucket}/{Key}?x-amz-copy-source&amp;partNumber&amp;uploadId | UploadPartCopy
VersioningApi | [**GetBucketVersioning**](docs/api/VersioningApi.md#getbucketversioning) | **Get** /{Bucket}?versioning | GetBucketVersioning
VersioningApi | [**PutBucketVersioning**](docs/api/VersioningApi.md#putbucketversioning) | **Put** /{Bucket}?versioning | PutBucketVersioning
VersionsApi | [**ListObjectVersions**](docs/api/VersionsApi.md#listobjectversions) | **Get** /{Bucket}?versions | ListObjectVersions
WebsiteApi | [**DeleteBucketWebsite**](docs/api/WebsiteApi.md#deletebucketwebsite) | **Delete** /{Bucket}?website | DeleteBucketWebsite
WebsiteApi | [**GetBucketWebsite**](docs/api/WebsiteApi.md#getbucketwebsite) | **Get** /{Bucket}?website | GetBucketWebsite
WebsiteApi | [**PutBucketWebsite**](docs/api/WebsiteApi.md#putbucketwebsite) | **Put** /{Bucket}?website | PutBucketWebsite

</details>

## Documentation For Models

All URIs are relative to *https://s3.eu-central-1.ionoscloud.com*
<details >
<summary title="Click to toggle">API models list</summary>

 - [AbortIncompleteMultipartUpload](docs/models/AbortIncompleteMultipartUpload)
 - [BlockPublicAccessOutput](docs/models/BlockPublicAccessOutput)
 - [BlockPublicAccessPayload](docs/models/BlockPublicAccessPayload)
 - [Bucket](docs/models/Bucket)
 - [BucketPolicy](docs/models/BucketPolicy)
 - [BucketPolicyStatement](docs/models/BucketPolicyStatement)
 - [BucketPolicyStatementCondition](docs/models/BucketPolicyStatementCondition)
 - [BucketPolicyStatementConditionDateGreaterThan](docs/models/BucketPolicyStatementConditionDateGreaterThan)
 - [BucketPolicyStatementConditionDateGreaterThanOneOf](docs/models/BucketPolicyStatementConditionDateGreaterThanOneOf)
 - [BucketPolicyStatementConditionDateGreaterThanOneOf1](docs/models/BucketPolicyStatementConditionDateGreaterThanOneOf1)
 - [BucketPolicyStatementConditionDateLessThan](docs/models/BucketPolicyStatementConditionDateLessThan)
 - [BucketPolicyStatementConditionDateLessThanOneOf](docs/models/BucketPolicyStatementConditionDateLessThanOneOf)
 - [BucketPolicyStatementConditionIpAddress](docs/models/BucketPolicyStatementConditionIpAddress)
 - [BucketPolicyStatementPrincipal](docs/models/BucketPolicyStatementPrincipal)
 - [BucketPolicyStatementPrincipalAnyOf](docs/models/BucketPolicyStatementPrincipalAnyOf)
 - [BucketVersioningStatus](docs/models/BucketVersioningStatus)
 - [CORSRule](docs/models/CORSRule)
 - [CSVInput](docs/models/CSVInput)
 - [CSVOutput](docs/models/CSVOutput)
 - [CommonPrefix](docs/models/CommonPrefix)
 - [CompleteMultipartUploadOutput](docs/models/CompleteMultipartUploadOutput)
 - [CompletedPart](docs/models/CompletedPart)
 - [CopyObjectOutput](docs/models/CopyObjectOutput)
 - [CopyObjectRequest](docs/models/CopyObjectRequest)
 - [CopyObjectResult](docs/models/CopyObjectResult)
 - [CopyPartResult](docs/models/CopyPartResult)
 - [CreateBucketRequest](docs/models/CreateBucketRequest)
 - [CreateBucketRequestCreateBucketConfiguration](docs/models/CreateBucketRequestCreateBucketConfiguration)
 - [CreateMultipartUploadOutput](docs/models/CreateMultipartUploadOutput)
 - [DefaultRetention](docs/models/DefaultRetention)
 - [DeleteMarkerEntry](docs/models/DeleteMarkerEntry)
 - [DeleteObjectsOutput](docs/models/DeleteObjectsOutput)
 - [DeleteObjectsRequest](docs/models/DeleteObjectsRequest)
 - [DeleteObjectsRequestDelete](docs/models/DeleteObjectsRequestDelete)
 - [DeletedObject](docs/models/DeletedObject)
 - [DeletionError](docs/models/DeletionError)
 - [Destination](docs/models/Destination)
 - [EncodingType](docs/models/EncodingType)
 - [Encryption](docs/models/Encryption)
 - [Error](docs/models/Error)
 - [ErrorDocument](docs/models/ErrorDocument)
 - [ErrorError](docs/models/ErrorError)
 - [Example](docs/models/Example)
 - [ExampleCompleteMultipartUpload](docs/models/ExampleCompleteMultipartUpload)
 - [ExpirationStatus](docs/models/ExpirationStatus)
 - [ExpressionType](docs/models/ExpressionType)
 - [GetBucketCorsOutput](docs/models/GetBucketCorsOutput)
 - [GetBucketEncryptionOutput](docs/models/GetBucketEncryptionOutput)
 - [GetBucketLifecycleOutput](docs/models/GetBucketLifecycleOutput)
 - [GetBucketLocation200Response](docs/models/GetBucketLocation200Response)
 - [GetBucketLogging200Response](docs/models/GetBucketLogging200Response)
 - [GetBucketPolicyStatusOutput](docs/models/GetBucketPolicyStatusOutput)
 - [GetBucketReplicationOutput](docs/models/GetBucketReplicationOutput)
 - [GetBucketTaggingOutput](docs/models/GetBucketTaggingOutput)
 - [GetBucketVersioningOutput](docs/models/GetBucketVersioningOutput)
 - [GetBucketWebsiteOutput](docs/models/GetBucketWebsiteOutput)
 - [GetObjectLockConfigurationOutput](docs/models/GetObjectLockConfigurationOutput)
 - [GetObjectLockConfigurationOutputObjectLockConfiguration](docs/models/GetObjectLockConfigurationOutputObjectLockConfiguration)
 - [GetObjectOutput](docs/models/GetObjectOutput)
 - [GetObjectRetentionOutput](docs/models/GetObjectRetentionOutput)
 - [GetObjectTaggingOutput](docs/models/GetObjectTaggingOutput)
 - [HeadObjectOutput](docs/models/HeadObjectOutput)
 - [IndexDocument](docs/models/IndexDocument)
 - [Initiator](docs/models/Initiator)
 - [InputSerialization](docs/models/InputSerialization)
 - [InputSerializationJSON](docs/models/InputSerializationJSON)
 - [JSONOutput](docs/models/JSONOutput)
 - [LifecycleExpiration](docs/models/LifecycleExpiration)
 - [ListAllMyBucketsResult](docs/models/ListAllMyBucketsResult)
 - [ListMultipartUploadsOutput](docs/models/ListMultipartUploadsOutput)
 - [ListObjectVersionsOutput](docs/models/ListObjectVersionsOutput)
 - [ListObjectsOutput](docs/models/ListObjectsOutput)
 - [ListObjectsV2Output](docs/models/ListObjectsV2Output)
 - [ListObjectsV2OutputListBucketResult](docs/models/ListObjectsV2OutputListBucketResult)
 - [ListPartsOutput](docs/models/ListPartsOutput)
 - [Metadata1](docs/models/Metadata1)
 - [MetadataEntry](docs/models/MetadataEntry)
 - [MfaDeleteStatus](docs/models/MfaDeleteStatus)
 - [MultipartUpload](docs/models/MultipartUpload)
 - [NoncurrentVersionExpiration](docs/models/NoncurrentVersionExpiration)
 - [Object](docs/models/Object)
 - [ObjectIdentifier](docs/models/ObjectIdentifier)
 - [ObjectLegalHoldConfiguration](docs/models/ObjectLegalHoldConfiguration)
 - [ObjectLockRetention](docs/models/ObjectLockRetention)
 - [ObjectLockRule](docs/models/ObjectLockRule)
 - [ObjectStorageClass](docs/models/ObjectStorageClass)
 - [ObjectVersion](docs/models/ObjectVersion)
 - [ObjectVersionStorageClass](docs/models/ObjectVersionStorageClass)
 - [OutputSerialization](docs/models/OutputSerialization)
 - [Owner](docs/models/Owner)
 - [Part](docs/models/Part)
 - [PolicyStatus](docs/models/PolicyStatus)
 - [PutBucketCorsRequest](docs/models/PutBucketCorsRequest)
 - [PutBucketCorsRequestCORSConfiguration](docs/models/PutBucketCorsRequestCORSConfiguration)
 - [PutBucketEncryptionRequest](docs/models/PutBucketEncryptionRequest)
 - [PutBucketEncryptionRequestServerSideEncryptionConfiguration](docs/models/PutBucketEncryptionRequestServerSideEncryptionConfiguration)
 - [PutBucketLifecycleRequest](docs/models/PutBucketLifecycleRequest)
 - [PutBucketLifecycleRequestLifecycleConfiguration](docs/models/PutBucketLifecycleRequestLifecycleConfiguration)
 - [PutBucketLoggingRequest](docs/models/PutBucketLoggingRequest)
 - [PutBucketLoggingRequestBucketLoggingStatus](docs/models/PutBucketLoggingRequestBucketLoggingStatus)
 - [PutBucketTaggingRequest](docs/models/PutBucketTaggingRequest)
 - [PutBucketTaggingRequestTagging](docs/models/PutBucketTaggingRequestTagging)
 - [PutBucketVersioningRequest](docs/models/PutBucketVersioningRequest)
 - [PutBucketVersioningRequestVersioningConfiguration](docs/models/PutBucketVersioningRequestVersioningConfiguration)
 - [PutBucketWebsiteRequest](docs/models/PutBucketWebsiteRequest)
 - [PutBucketWebsiteRequestWebsiteConfiguration](docs/models/PutBucketWebsiteRequestWebsiteConfiguration)
 - [PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo](docs/models/PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo)
 - [PutObjectLegalHoldRequest](docs/models/PutObjectLegalHoldRequest)
 - [PutObjectLegalHoldRequestLegalHold](docs/models/PutObjectLegalHoldRequestLegalHold)
 - [PutObjectLockConfigurationRequest](docs/models/PutObjectLockConfigurationRequest)
 - [PutObjectLockConfigurationRequestObjectLockConfiguration](docs/models/PutObjectLockConfigurationRequestObjectLockConfiguration)
 - [PutObjectLockConfigurationRequestObjectLockConfigurationRule](docs/models/PutObjectLockConfigurationRequestObjectLockConfigurationRule)
 - [PutObjectRequest](docs/models/PutObjectRequest)
 - [PutObjectRetentionRequest](docs/models/PutObjectRetentionRequest)
 - [PutObjectRetentionRequestRetention](docs/models/PutObjectRetentionRequestRetention)
 - [Redirect](docs/models/Redirect)
 - [RedirectAllRequestsTo](docs/models/RedirectAllRequestsTo)
 - [ReplicationConfiguration](docs/models/ReplicationConfiguration)
 - [ReplicationRule](docs/models/ReplicationRule)
 - [RoutingRule](docs/models/RoutingRule)
 - [RoutingRuleCondition](docs/models/RoutingRuleCondition)
 - [Rule](docs/models/Rule)
 - [ServerSideEncryption](docs/models/ServerSideEncryption)
 - [ServerSideEncryptionByDefault](docs/models/ServerSideEncryptionByDefault)
 - [ServerSideEncryptionConfiguration](docs/models/ServerSideEncryptionConfiguration)
 - [ServerSideEncryptionRule](docs/models/ServerSideEncryptionRule)
 - [StorageClass](docs/models/StorageClass)
 - [Tag](docs/models/Tag)
 - [UploadPartCopyOutput](docs/models/UploadPartCopyOutput)
 - [UploadPartRequest](docs/models/UploadPartRequest)


[[Back to API list]](#documentation-for-api-endpoints) [[Back to Model list]](#documentation-for-models)

</details>
