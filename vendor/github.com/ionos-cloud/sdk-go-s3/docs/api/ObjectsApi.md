# \ObjectsApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**CopyObject**](ObjectsApi.md#CopyObject) | **Put** /{Bucket}/{Key}?x-amz-copy-source | CopyObject|
|[**DeleteObject**](ObjectsApi.md#DeleteObject) | **Delete** /{Bucket}/{Key} | DeleteObject|
|[**DeleteObjects**](ObjectsApi.md#DeleteObjects) | **Post** /{Bucket}?delete | DeleteObjects|
|[**GetObject**](ObjectsApi.md#GetObject) | **Get** /{Bucket}/{Key} | GetObject|
|[**HeadObject**](ObjectsApi.md#HeadObject) | **Head** /{Bucket}/{Key} | HeadObject|
|[**ListObjects**](ObjectsApi.md#ListObjects) | **Get** /{Bucket} | ListObjects|
|[**ListObjectsV2**](ObjectsApi.md#ListObjectsV2) | **Get** /{Bucket}?list-type&#x3D;2 | ListObjectsV2|
|[**OPTIONSObject**](ObjectsApi.md#OPTIONSObject) | **Options** /{Bucket} | OPTIONSObject|
|[**POSTObject**](ObjectsApi.md#POSTObject) | **Post** /{Bucket}/{Key} | POSTObject|
|[**PutObject**](ObjectsApi.md#PutObject) | **Put** /{Bucket}/{Key} | PutObject|



## CopyObject

```go
var result CopyObjectResult = CopyObject(ctx, bucket, key)
                      .XAmzCopySource(xAmzCopySource)
                      .CacheControl(cacheControl)
                      .ContentDisposition(contentDisposition)
                      .ContentEncoding(contentEncoding)
                      .ContentLanguage(contentLanguage)
                      .ContentType(contentType)
                      .XAmzCopySourceIfMatch(xAmzCopySourceIfMatch)
                      .XAmzCopySourceIfModifiedSince(xAmzCopySourceIfModifiedSince)
                      .XAmzCopySourceIfNoneMatch(xAmzCopySourceIfNoneMatch)
                      .XAmzCopySourceIfUnmodifiedSince(xAmzCopySourceIfUnmodifiedSince)
                      .Expires(expires)
                      .XAmzMetadataDirective(xAmzMetadataDirective)
                      .XAmzTaggingDirective(xAmzTaggingDirective)
                      .XAmzServerSideEncryption(xAmzServerSideEncryption)
                      .XAmzStorageClass(xAmzStorageClass)
                      .XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .XAmzCopySourceServerSideEncryptionCustomerAlgorithm(xAmzCopySourceServerSideEncryptionCustomerAlgorithm)
                      .XAmzCopySourceServerSideEncryptionCustomerKey(xAmzCopySourceServerSideEncryptionCustomerKey)
                      .XAmzCopySourceServerSideEncryptionCustomerKeyMD5(xAmzCopySourceServerSideEncryptionCustomerKeyMD5)
                      .XAmzTagging(xAmzTagging)
                      .XAmzObjectLockMode(xAmzObjectLockMode)
                      .XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate)
                      .XAmzObjectLockLegalHold(xAmzObjectLockLegalHold)
                      .CopyObjectRequest(copyObjectRequest)
                      .Execute()
```

CopyObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    xAmzCopySource := "xAmzCopySource_example" // string | <p>Specifies the source object for the copy operation.
    key := "key_example" // string | The key of the destination object.
    cacheControl := "cacheControl_example" // string | Specifies caching behavior along the request/reply chain. (optional)
    contentDisposition := "contentDisposition_example" // string | Specifies presentational information for the object. (optional)
    contentEncoding := "contentEncoding_example" // string | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. (optional)
    contentLanguage := "contentLanguage_example" // string | The language the content is in. (optional)
    contentType := "contentType_example" // string | A standard MIME type describing the format of the object data. (optional)
    xAmzCopySourceIfMatch := "xAmzCopySourceIfMatch_example" // string | Copies the object if its entity tag (ETag) matches the specified tag. (optional)
    xAmzCopySourceIfModifiedSince := time.Now() // time.Time | Copies the object if it has been modified since the specified time. (optional)
    xAmzCopySourceIfNoneMatch := "xAmzCopySourceIfNoneMatch_example" // string | Copies the object if its entity tag (ETag) is different than the specified ETag. (optional)
    xAmzCopySourceIfUnmodifiedSince := time.Now() // time.Time | Copies the object if it hasn't been modified since the specified time. (optional)
    expires := time.Now() // time.Time | The date and time at which the object is no longer cacheable. (optional)
    xAmzMetadataDirective := "xAmzMetadataDirective_example" // string | Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. (optional)
    xAmzTaggingDirective := "xAmzTaggingDirective_example" // string | Specifies whether the object tag-set are copied from the source object or replaced with tag-set provided in the request. (optional)
    xAmzServerSideEncryption := "xAmzServerSideEncryption_example" // string | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). (optional)
    xAmzStorageClass := "xAmzStorageClass_example" // string | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. (optional)
    xAmzWebsiteRedirectLocation := "xAmzWebsiteRedirectLocation_example" // string | If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata. (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm` header. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    xAmzCopySourceServerSideEncryptionCustomerAlgorithm := "xAmzCopySourceServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use when decrypting the source object (AES256). (optional)
    xAmzCopySourceServerSideEncryptionCustomerKey := "xAmzCopySourceServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use to decrypt the source object. The encryption key provided in this header must be one that was used when the source object was created. (optional)
    xAmzCopySourceServerSideEncryptionCustomerKeyMD5 := "xAmzCopySourceServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    xAmzTagging := "xAmzTagging_example" // string | The tag-set for the object destination object this value must be used in conjunction with the `TaggingDirective`. The tag-set must be encoded as URL Query parameters. (optional)
    xAmzObjectLockMode := "xAmzObjectLockMode_example" // string | The Object Lock mode that you want to apply to the copied object. (optional)
    xAmzObjectLockRetainUntilDate := time.Now() // time.Time | The date and time when you want the copied object's Object Lock to expire. (optional)
    xAmzObjectLockLegalHold := "xAmzObjectLockLegalHold_example" // string | Specifies whether you want to apply a Legal Hold to the copied object. (optional)
    copyObjectRequest := *openapiclient.NewCopyObjectRequest() // CopyObjectRequest |  (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.CopyObject(context.Background(), bucket, key).XAmzCopySource(xAmzCopySource).CacheControl(cacheControl).ContentDisposition(contentDisposition).ContentEncoding(contentEncoding).ContentLanguage(contentLanguage).ContentType(contentType).XAmzCopySourceIfMatch(xAmzCopySourceIfMatch).XAmzCopySourceIfModifiedSince(xAmzCopySourceIfModifiedSince).XAmzCopySourceIfNoneMatch(xAmzCopySourceIfNoneMatch).XAmzCopySourceIfUnmodifiedSince(xAmzCopySourceIfUnmodifiedSince).Expires(expires).XAmzMetadataDirective(xAmzMetadataDirective).XAmzTaggingDirective(xAmzTaggingDirective).XAmzServerSideEncryption(xAmzServerSideEncryption).XAmzStorageClass(xAmzStorageClass).XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).XAmzCopySourceServerSideEncryptionCustomerAlgorithm(xAmzCopySourceServerSideEncryptionCustomerAlgorithm).XAmzCopySourceServerSideEncryptionCustomerKey(xAmzCopySourceServerSideEncryptionCustomerKey).XAmzCopySourceServerSideEncryptionCustomerKeyMD5(xAmzCopySourceServerSideEncryptionCustomerKeyMD5).XAmzTagging(xAmzTagging).XAmzObjectLockMode(xAmzObjectLockMode).XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate).XAmzObjectLockLegalHold(xAmzObjectLockLegalHold).CopyObjectRequest(copyObjectRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.CopyObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `CopyObject`: CopyObjectResult
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.CopyObject`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | The key of the destination object. | |

### Other Parameters

Other parameters are passed through a pointer to an apiCopyObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **xAmzCopySource** | **string** | &lt;p&gt;Specifies the source object for the copy operation. | |
| **cacheControl** | **string** | Specifies caching behavior along the request/reply chain. | |
| **contentDisposition** | **string** | Specifies presentational information for the object. | |
| **contentEncoding** | **string** | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. | |
| **contentLanguage** | **string** | The language the content is in. | |
| **contentType** | **string** | A standard MIME type describing the format of the object data. | |
| **xAmzCopySourceIfMatch** | **string** | Copies the object if its entity tag (ETag) matches the specified tag. | |
| **xAmzCopySourceIfModifiedSince** | **time.Time** | Copies the object if it has been modified since the specified time. | |
| **xAmzCopySourceIfNoneMatch** | **string** | Copies the object if its entity tag (ETag) is different than the specified ETag. | |
| **xAmzCopySourceIfUnmodifiedSince** | **time.Time** | Copies the object if it hasn&#39;t been modified since the specified time. | |
| **expires** | **time.Time** | The date and time at which the object is no longer cacheable. | |
| **xAmzMetadataDirective** | **string** | Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. | |
| **xAmzTaggingDirective** | **string** | Specifies whether the object tag-set are copied from the source object or replaced with tag-set provided in the request. | |
| **xAmzServerSideEncryption** | **string** | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). | |
| **xAmzStorageClass** | **string** | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. | |
| **xAmzWebsiteRedirectLocation** | **string** | If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata. | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm&#x60; header. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **xAmzCopySourceServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use when decrypting the source object (AES256). | |
| **xAmzCopySourceServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use to decrypt the source object. The encryption key provided in this header must be one that was used when the source object was created. | |
| **xAmzCopySourceServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **xAmzTagging** | **string** | The tag-set for the object destination object this value must be used in conjunction with the &#x60;TaggingDirective&#x60;. The tag-set must be encoded as URL Query parameters. | |
| **xAmzObjectLockMode** | **string** | The Object Lock mode that you want to apply to the copied object. | |
| **xAmzObjectLockRetainUntilDate** | **time.Time** | The date and time when you want the copied object&#39;s Object Lock to expire. | |
| **xAmzObjectLockLegalHold** | **string** | Specifies whether you want to apply a Legal Hold to the copied object. | |
| **copyObjectRequest** | [**CopyObjectRequest**](../models/CopyObjectRequest.md) |  | |

### Return type

[**CopyObjectResult**](../models/CopyObjectResult.md)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.CopyObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.CopyObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.CopyObject": {
    "port": "8443",
},
})
```


## DeleteObject

```go
var result map[string]interface{} = DeleteObject(ctx, bucket, key)
                      .XAmzMfa(xAmzMfa)
                      .VersionId(versionId)
                      .XAmzBypassGovernanceRetention(xAmzBypassGovernanceRetention)
                      .Execute()
```

DeleteObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    key := "key_example" // string | Key name of the object to delete.
    xAmzMfa := "xAmzMfa_example" // string | The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled. (optional)
    versionId := "versionId_example" // string | VersionId used to reference a specific version of the object. (optional)
    xAmzBypassGovernanceRetention := true // bool | Indicates whether S3 Object Lock should bypass Governance-mode restrictions to process this operation. To use this header, you must have the `PutBucketPublicAccessBlock` permission. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resp, err := apiClient.ObjectsApi.DeleteObject(context.Background(), bucket, key).XAmzMfa(xAmzMfa).VersionId(versionId).XAmzBypassGovernanceRetention(xAmzBypassGovernanceRetention).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.DeleteObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `DeleteObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.DeleteObject`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Key name of the object to delete. | |

### Other Parameters

Other parameters are passed through a pointer to an apiDeleteObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **xAmzMfa** | **string** | The concatenation of the authentication device&#39;s serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled. | |
| **versionId** | **string** | VersionId used to reference a specific version of the object. | |
| **xAmzBypassGovernanceRetention** | **bool** | Indicates whether S3 Object Lock should bypass Governance-mode restrictions to process this operation. To use this header, you must have the &#x60;PutBucketPublicAccessBlock&#x60; permission. | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.DeleteObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.DeleteObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.DeleteObject": {
    "port": "8443",
},
})
```


## DeleteObjects

```go
var result DeleteObjectsOutput = DeleteObjects(ctx, bucket)
                      .DeleteObjectsRequest(deleteObjectsRequest)
                      .XAmzMfa(xAmzMfa)
                      .XAmzBypassGovernanceRetention(xAmzBypassGovernanceRetention)
                      .Execute()
```

DeleteObjects



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    deleteObjectsRequest := *openapiclient.NewDeleteObjectsRequest() // DeleteObjectsRequest | 
    xAmzMfa := "xAmzMfa_example" // string | The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled. (optional)
    xAmzBypassGovernanceRetention := true // bool | Specifies whether you want to delete this object even if it has a Governance-type Object Lock in place. To use this header, you must have the `PutBucketPublicAccessBlock` permission. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.DeleteObjects(context.Background(), bucket).DeleteObjectsRequest(deleteObjectsRequest).XAmzMfa(xAmzMfa).XAmzBypassGovernanceRetention(xAmzBypassGovernanceRetention).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.DeleteObjects``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `DeleteObjects`: DeleteObjectsOutput
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.DeleteObjects`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiDeleteObjectsRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deleteObjectsRequest** | [**DeleteObjectsRequest**](../models/DeleteObjectsRequest.md) |  | |
| **xAmzMfa** | **string** | The concatenation of the authentication device&#39;s serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled. | |
| **xAmzBypassGovernanceRetention** | **bool** | Specifies whether you want to delete this object even if it has a Governance-type Object Lock in place. To use this header, you must have the &#x60;PutBucketPublicAccessBlock&#x60; permission. | |

### Return type

[**DeleteObjectsOutput**](../models/DeleteObjectsOutput.md)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.DeleteObjects"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.DeleteObjects": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.DeleteObjects": {
    "port": "8443",
},
})
```


## GetObject

```go
var result *os.File = GetObject(ctx, bucket, key)
                      .IfMatch(ifMatch)
                      .IfModifiedSince(ifModifiedSince)
                      .IfNoneMatch(ifNoneMatch)
                      .IfUnmodifiedSince(ifUnmodifiedSince)
                      .Range_(range_)
                      .ResponseCacheControl(responseCacheControl)
                      .ResponseContentDisposition(responseContentDisposition)
                      .ResponseContentEncoding(responseContentEncoding)
                      .ResponseContentLanguage(responseContentLanguage)
                      .ResponseContentType(responseContentType)
                      .ResponseExpires(responseExpires)
                      .VersionId(versionId)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .PartNumber(partNumber)
                      .Execute()
```

GetObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    key := "key_example" // string | <p> Key of the object to get. </p> <p> <b> Possible values:</b> length ≥ 1 </p>
    ifMatch := "ifMatch_example" // string | Return the object only if its entity tag (ETag) is the same as the one specified, otherwise return a 412 (precondition failed). (optional)
    ifModifiedSince := time.Now() // time.Time | Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified). (optional)
    ifNoneMatch := "ifNoneMatch_example" // string | Return the object only if its entity tag (ETag) is different from the one specified, otherwise return a 304 (not modified). (optional)
    ifUnmodifiedSince := time.Now() // time.Time | Return the object only if it has not been modified since the specified time, otherwise return a 412 (precondition failed). (optional)
    range_ := "range__example" // string | <p>Downloads the specified range bytes of an object. For more information about the HTTP Range header, see <a href=\"https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35\">Range</a>.</p> <note> <p>IONOS S3 Object Storage doesn't support retrieving multiple ranges of data per `GET` request.</p> </note> (optional)
    responseCacheControl := "responseCacheControl_example" // string | Sets the `Cache-Control` header of the response. (optional)
    responseContentDisposition := "responseContentDisposition_example" // string | Sets the `Content-Disposition` header of the response (optional)
    responseContentEncoding := "responseContentEncoding_example" // string | Sets the `Content-Encoding` header of the response. (optional)
    responseContentLanguage := "responseContentLanguage_example" // string | Sets the `Content-Language` header of the response. (optional)
    responseContentType := "responseContentType_example" // string | Sets the `Content-Type` header of the response. (optional)
    responseExpires := time.Now() // time.Time | Sets the `Expires` header of the response. (optional)
    versionId := "versionId_example" // string | VersionId used to reference a specific version of the object. (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when decrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage used to encrypt the data. This value is used to decrypt the object when recovering it and must match the one used when storing the data. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm` header. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    partNumber := int32(56) // int32 | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a 'ranged' GET request for the part specified. Useful for downloading just a part of an object. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.GetObject(context.Background(), bucket, key).IfMatch(ifMatch).IfModifiedSince(ifModifiedSince).IfNoneMatch(ifNoneMatch).IfUnmodifiedSince(ifUnmodifiedSince).Range_(range_).ResponseCacheControl(responseCacheControl).ResponseContentDisposition(responseContentDisposition).ResponseContentEncoding(responseContentEncoding).ResponseContentLanguage(responseContentLanguage).ResponseContentType(responseContentType).ResponseExpires(responseExpires).VersionId(versionId).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).PartNumber(partNumber).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.GetObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetObject`: *os.File
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.GetObject`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | &lt;p&gt; Key of the object to get. &lt;/p&gt; &lt;p&gt; &lt;b&gt; Possible values:&lt;/b&gt; length ≥ 1 &lt;/p&gt; | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **ifMatch** | **string** | Return the object only if its entity tag (ETag) is the same as the one specified, otherwise return a 412 (precondition failed). | |
| **ifModifiedSince** | **time.Time** | Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified). | |
| **ifNoneMatch** | **string** | Return the object only if its entity tag (ETag) is different from the one specified, otherwise return a 304 (not modified). | |
| **ifUnmodifiedSince** | **time.Time** | Return the object only if it has not been modified since the specified time, otherwise return a 412 (precondition failed). | |
| **range_** | **string** | &lt;p&gt;Downloads the specified range bytes of an object. For more information about the HTTP Range header, see &lt;a href&#x3D;\&quot;https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35\&quot;&gt;Range&lt;/a&gt;.&lt;/p&gt; &lt;note&gt; &lt;p&gt;IONOS S3 Object Storage doesn&#39;t support retrieving multiple ranges of data per &#x60;GET&#x60; request.&lt;/p&gt; &lt;/note&gt; | |
| **responseCacheControl** | **string** | Sets the &#x60;Cache-Control&#x60; header of the response. | |
| **responseContentDisposition** | **string** | Sets the &#x60;Content-Disposition&#x60; header of the response | |
| **responseContentEncoding** | **string** | Sets the &#x60;Content-Encoding&#x60; header of the response. | |
| **responseContentLanguage** | **string** | Sets the &#x60;Content-Language&#x60; header of the response. | |
| **responseContentType** | **string** | Sets the &#x60;Content-Type&#x60; header of the response. | |
| **responseExpires** | **time.Time** | Sets the &#x60;Expires&#x60; header of the response. | |
| **versionId** | **string** | VersionId used to reference a specific version of the object. | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when decrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage used to encrypt the data. This value is used to decrypt the object when recovering it and must match the one used when storing the data. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm&#x60; header. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **partNumber** | **int32** | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a &#39;ranged&#39; GET request for the part specified. Useful for downloading just a part of an object. | |

### Return type

[***os.File**](../models/*os.File.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.GetObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.GetObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.GetObject": {
    "port": "8443",
},
})
```


## HeadObject

```go
var result HeadObjectOutput = HeadObject(ctx, bucket, key)
                      .IfMatch(ifMatch)
                      .IfModifiedSince(ifModifiedSince)
                      .IfNoneMatch(ifNoneMatch)
                      .IfUnmodifiedSince(ifUnmodifiedSince)
                      .Range_(range_)
                      .VersionId(versionId)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .PartNumber(partNumber)
                      .Execute()
```

HeadObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    key := "key_example" // string | The object key.
    ifMatch := "ifMatch_example" // string | Return the object only if its entity tag (ETag) is the same as the one specified, otherwise return a 412 (precondition failed). (optional)
    ifModifiedSince := time.Now() // time.Time | Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified). (optional)
    ifNoneMatch := "ifNoneMatch_example" // string | Return the object only if its entity tag (ETag) is different from the one specified, otherwise return a 304 (not modified). (optional)
    ifUnmodifiedSince := time.Now() // time.Time | Return the object only if it has not been modified since the specified time, otherwise return a 412 (precondition failed). (optional)
    range_ := "range__example" // string | <p>Downloads the specified range bytes of an object. For more information about the HTTP Range header, see <a href=\"https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35\">Range</a>.</p> <note> <p>IONOS S3 Object Storage doesn't support retrieving multiple ranges of data per `GET` request.</p> </note> (optional)
    versionId := "versionId_example" // string | VersionId used to reference a specific version of the object. (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm` header. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    partNumber := int32(56) // int32 | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a 'ranged' HEAD request for the part specified. Useful querying about the size of the part and the number of parts in this object. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.HeadObject(context.Background(), bucket, key).IfMatch(ifMatch).IfModifiedSince(ifModifiedSince).IfNoneMatch(ifNoneMatch).IfUnmodifiedSince(ifUnmodifiedSince).Range_(range_).VersionId(versionId).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).PartNumber(partNumber).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.HeadObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `HeadObject`: HeadObjectOutput
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.HeadObject`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | The object key. | |

### Other Parameters

Other parameters are passed through a pointer to an apiHeadObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **ifMatch** | **string** | Return the object only if its entity tag (ETag) is the same as the one specified, otherwise return a 412 (precondition failed). | |
| **ifModifiedSince** | **time.Time** | Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified). | |
| **ifNoneMatch** | **string** | Return the object only if its entity tag (ETag) is different from the one specified, otherwise return a 304 (not modified). | |
| **ifUnmodifiedSince** | **time.Time** | Return the object only if it has not been modified since the specified time, otherwise return a 412 (precondition failed). | |
| **range_** | **string** | &lt;p&gt;Downloads the specified range bytes of an object. For more information about the HTTP Range header, see &lt;a href&#x3D;\&quot;https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35\&quot;&gt;Range&lt;/a&gt;.&lt;/p&gt; &lt;note&gt; &lt;p&gt;IONOS S3 Object Storage doesn&#39;t support retrieving multiple ranges of data per &#x60;GET&#x60; request.&lt;/p&gt; &lt;/note&gt; | |
| **versionId** | **string** | VersionId used to reference a specific version of the object. | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm&#x60; header. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **partNumber** | **int32** | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a &#39;ranged&#39; HEAD request for the part specified. Useful querying about the size of the part and the number of parts in this object. | |

### Return type

[**HeadObjectOutput**](../models/HeadObjectOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.HeadObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.HeadObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.HeadObject": {
    "port": "8443",
},
})
```


## ListObjects

```go
var result ListObjectsOutput = ListObjects(ctx, bucket)
                      .Delimiter(delimiter)
                      .EncodingType(encodingType)
                      .Marker(marker)
                      .MaxKeys(maxKeys)
                      .Prefix(prefix)
                      .XAmzRequestPayer(xAmzRequestPayer)
                      .MaxKeys2(maxKeys2)
                      .Marker2(marker2)
                      .Execute()
```

ListObjects



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    delimiter := "delimiter_example" // string | A delimiter is a character you use to group keys. (optional)
    encodingType := "encodingType_example" // string |  (optional)
    marker := "marker_example" // string | Marker is where you want IONOS S3 Object Storage to start listing from. IONOS S3 Object Storage starts listing after this specified key. Marker can be any key in the bucket. (optional)
    maxKeys := int32(56) // int32 | Sets the maximum number of keys returned in the response. By default the operation returns up to 1,000 key names. The response might contain fewer keys but will never contain more.  (optional)
    prefix := "prefix_example" // string | Limits the response to keys that begin with the specified prefix. (optional)
    xAmzRequestPayer := "xAmzRequestPayer_example" // string | Confirms that the requester knows that she or he will be charged for the list objects request. Bucket owners need not specify this parameter in their requests. (optional)
    maxKeys2 := "maxKeys_example" // string | Pagination limit (optional)
    marker2 := "marker_example" // string | Pagination token (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.ListObjects(context.Background(), bucket).Delimiter(delimiter).EncodingType(encodingType).Marker(marker).MaxKeys(maxKeys).Prefix(prefix).XAmzRequestPayer(xAmzRequestPayer).MaxKeys2(maxKeys2).Marker2(marker2).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.ListObjects``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `ListObjects`: ListObjectsOutput
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.ListObjects`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiListObjectsRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **delimiter** | **string** | A delimiter is a character you use to group keys. | |
| **encodingType** | **string** |  | |
| **marker** | **string** | Marker is where you want IONOS S3 Object Storage to start listing from. IONOS S3 Object Storage starts listing after this specified key. Marker can be any key in the bucket. | |
| **maxKeys** | **int32** | Sets the maximum number of keys returned in the response. By default the operation returns up to 1,000 key names. The response might contain fewer keys but will never contain more.  | |
| **prefix** | **string** | Limits the response to keys that begin with the specified prefix. | |
| **xAmzRequestPayer** | **string** | Confirms that the requester knows that she or he will be charged for the list objects request. Bucket owners need not specify this parameter in their requests. | |
| **maxKeys2** | **string** | Pagination limit | |
| **marker2** | **string** | Pagination token | |

### Return type

[**ListObjectsOutput**](../models/ListObjectsOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.ListObjects"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.ListObjects": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.ListObjects": {
    "port": "8443",
},
})
```


## ListObjectsV2

```go
var result ListBucketResultV2 = ListObjectsV2(ctx, bucket)
                      .Delimiter(delimiter)
                      .EncodingType(encodingType)
                      .MaxKeys(maxKeys)
                      .Prefix(prefix)
                      .ContinuationToken(continuationToken)
                      .FetchOwner(fetchOwner)
                      .StartAfter(startAfter)
                      .Execute()
```

ListObjectsV2



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    delimiter := "/" // string | A delimiter is a character you use to group keys. (optional)
    encodingType := "encodingType_example" // string | Encoding type used by IONOS S3 Object Storage to encode object keys in the response. (optional)
    maxKeys := int32(56) // int32 | Sets the maximum number of keys returned in the response. By default the operation returns up to 1000 key names. The response might contain fewer keys but will never contain more. (optional) (default to 1000)
    prefix := "folder/subfolder/" // string | Limits the response to keys that begin with the specified prefix. (optional)
    continuationToken := "continuationToken_example" // string | ContinuationToken indicates IONOS S3 Object Storage that the list is being continued on this bucket with a token. ContinuationToken is obfuscated and is not a real key. (optional)
    fetchOwner := true // bool | The owner field is not present in listV2 by default, if you want to return owner field with each key in the result then set the fetch owner field to true. (optional) (default to false)
    startAfter := "startAfter_example" // string | StartAfter is where you want to start listing from. IONOS S3 Object Storage starts listing after this specified key. StartAfter can be any key in the bucket. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.ListObjectsV2(context.Background(), bucket).Delimiter(delimiter).EncodingType(encodingType).MaxKeys(maxKeys).Prefix(prefix).ContinuationToken(continuationToken).FetchOwner(fetchOwner).StartAfter(startAfter).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.ListObjectsV2``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `ListObjectsV2`: ListBucketResultV2
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.ListObjectsV2`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiListObjectsV2Request struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **delimiter** | **string** | A delimiter is a character you use to group keys. | |
| **encodingType** | **string** | Encoding type used by IONOS S3 Object Storage to encode object keys in the response. | |
| **maxKeys** | **int32** | Sets the maximum number of keys returned in the response. By default the operation returns up to 1000 key names. The response might contain fewer keys but will never contain more. | [default to 1000]|
| **prefix** | **string** | Limits the response to keys that begin with the specified prefix. | |
| **continuationToken** | **string** | ContinuationToken indicates IONOS S3 Object Storage that the list is being continued on this bucket with a token. ContinuationToken is obfuscated and is not a real key. | |
| **fetchOwner** | **bool** | The owner field is not present in listV2 by default, if you want to return owner field with each key in the result then set the fetch owner field to true. | [default to false]|
| **startAfter** | **string** | StartAfter is where you want to start listing from. IONOS S3 Object Storage starts listing after this specified key. StartAfter can be any key in the bucket. | |

### Return type

[**ListBucketResultV2**](../models/ListBucketResultV2.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.ListObjectsV2"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.ListObjectsV2": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.ListObjectsV2": {
    "port": "8443",
},
})
```


## OPTIONSObject

```go
var result  = OPTIONSObject(ctx, bucket)
                      .Origin(origin)
                      .AccessControlRequestMethod(accessControlRequestMethod)
                      .AccessControlRequestHeaders(accessControlRequestHeaders)
                      .Execute()
```

OPTIONSObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    origin := "origin_example" // string | <p>Identifies the origin of the cross-origin request to the IONOS S3 Object Storage. </p>
    accessControlRequestMethod := "accessControlRequestMethod_example" // string |  Identifies what HTTP method will be used in the actual request.
    accessControlRequestHeaders := "accessControlRequestHeaders_example" // string | <p> A comma-delimited list of HTTP headers that will be sent in the actual request. </p> <p> For example, to put an object with server-side encryption, this preflight request  will determine if it can include the `x-amz-server-side-encryption` header with the request. </p> (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.OPTIONSObject(context.Background(), bucket).Origin(origin).AccessControlRequestMethod(accessControlRequestMethod).AccessControlRequestHeaders(accessControlRequestHeaders).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.OPTIONSObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiOPTIONSObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **origin** | **string** | &lt;p&gt;Identifies the origin of the cross-origin request to the IONOS S3 Object Storage. &lt;/p&gt; | |
| **accessControlRequestMethod** | **string** |  Identifies what HTTP method will be used in the actual request. | |
| **accessControlRequestHeaders** | **string** | &lt;p&gt; A comma-delimited list of HTTP headers that will be sent in the actual request. &lt;/p&gt; &lt;p&gt; For example, to put an object with server-side encryption, this preflight request  will determine if it can include the &#x60;x-amz-server-side-encryption&#x60; header with the request. &lt;/p&gt; | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.OPTIONSObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.OPTIONSObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.OPTIONSObject": {
    "port": "8443",
},
})
```


## POSTObject

```go
var result map[string]interface{} = POSTObject(ctx, bucket, key)
                      .POSTObjectRequest(pOSTObjectRequest)
                      .CacheControl(cacheControl)
                      .ContentDisposition(contentDisposition)
                      .ContentEncoding(contentEncoding)
                      .ContentLanguage(contentLanguage)
                      .ContentLength(contentLength)
                      .ContentMD5(contentMD5)
                      .ContentType(contentType)
                      .Expires(expires)
                      .XAmzServerSideEncryption(xAmzServerSideEncryption)
                      .XAmzStorageClass(xAmzStorageClass)
                      .XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .XAmzServerSideEncryptionContext(xAmzServerSideEncryptionContext)
                      .XAmzServerSideEncryptionBucketKeyEnabled(xAmzServerSideEncryptionBucketKeyEnabled)
                      .XAmzRequestPayer(xAmzRequestPayer)
                      .XAmzTagging(xAmzTagging)
                      .XAmzObjectLockMode(xAmzObjectLockMode)
                      .XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate)
                      .XAmzObjectLockLegalHold(xAmzObjectLockLegalHold)
                      .Execute()
```

POSTObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    key := "key_example" // string | Key name of the object to post.
    pOSTObjectRequest := *openapiclient.NewPOSTObjectRequest() // POSTObjectRequest | 
    cacheControl := "cacheControl_example" // string |  Can be used to specify caching behavior along the request/reply chain. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9\">Cache-Control</a>. (optional)
    contentDisposition := "contentDisposition_example" // string | Specifies presentational information for the object. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1\">Content-Disposition</a>. (optional)
    contentEncoding := "contentEncoding_example" // string | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11\">Content-Encoding</a>. (optional)
    contentLanguage := "contentLanguage_example" // string | The language the content is in. (optional)
    contentLength := int32(56) // int32 | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.13\">Content-Length</a>. (optional)
    contentMD5 := "contentMD5_example" // string |  (optional)
    contentType := "contentType_example" // string | A standard MIME type describing the format of the contents. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.17\">Content-Type</a>. (optional)
    expires := time.Now() // time.Time | The date and time at which the object is no longer cacheable. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21\">Expires</a>. (optional)
    xAmzServerSideEncryption := "xAmzServerSideEncryption_example" // string | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). (optional)
    xAmzStorageClass := "xAmzStorageClass_example" // string | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. (optional)
    xAmzWebsiteRedirectLocation := "xAmzWebsiteRedirectLocation_example" // string | <p>If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata.</p> <p>In the following example, the request header sets the redirect to an object (anotherPage.html) in the same bucket:</p> <p> `x-amz-website-redirect-location: /anotherPage.html` </p> <p>In the following example, the request header sets the object redirect to another website:</p> <p> `x-amz-website-redirect-location: http://www.example.com/` </p> (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm` header. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    xAmzServerSideEncryptionContext := "xAmzServerSideEncryptionContext_example" // string | Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs. (optional)
    xAmzServerSideEncryptionBucketKeyEnabled := true // bool | <p>Specifies whether IONOS S3 Object Storage should use an S3 Bucket Key for object encryption with server-side encryption. Setting this header to `true` causes IONOS S3 Object Storage to use an S3 Bucket Key for object encryption.</p> <p>Specifying this header with a PUT operation doesn’t affect bucket-level settings for S3 Bucket Key.</p> (optional)
    xAmzRequestPayer := "xAmzRequestPayer_example" // string |  (optional)
    xAmzTagging := "xAmzTagging_example" // string | The tag-set for the object. The tag-set must be encoded as URL Query parameters. (For example, \"Key1=Value1\") (optional)
    xAmzObjectLockMode := "xAmzObjectLockMode_example" // string | The Object Lock mode that you want to apply to this object. (optional)
    xAmzObjectLockRetainUntilDate := time.Now() // time.Time | The date and time when you want this object's Object Lock to expire. Must be formatted as a timestamp parameter. (optional)
    xAmzObjectLockLegalHold := "xAmzObjectLockLegalHold_example" // string | Specifies whether a legal hold will be applied to this object. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.POSTObject(context.Background(), bucket, key).POSTObjectRequest(pOSTObjectRequest).CacheControl(cacheControl).ContentDisposition(contentDisposition).ContentEncoding(contentEncoding).ContentLanguage(contentLanguage).ContentLength(contentLength).ContentMD5(contentMD5).ContentType(contentType).Expires(expires).XAmzServerSideEncryption(xAmzServerSideEncryption).XAmzStorageClass(xAmzStorageClass).XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).XAmzServerSideEncryptionContext(xAmzServerSideEncryptionContext).XAmzServerSideEncryptionBucketKeyEnabled(xAmzServerSideEncryptionBucketKeyEnabled).XAmzRequestPayer(xAmzRequestPayer).XAmzTagging(xAmzTagging).XAmzObjectLockMode(xAmzObjectLockMode).XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate).XAmzObjectLockLegalHold(xAmzObjectLockLegalHold).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.POSTObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `POSTObject`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `ObjectsApi.POSTObject`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Key name of the object to post. | |

### Other Parameters

Other parameters are passed through a pointer to an apiPOSTObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **pOSTObjectRequest** | [**POSTObjectRequest**](../models/POSTObjectRequest.md) |  | |
| **cacheControl** | **string** |  Can be used to specify caching behavior along the request/reply chain. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9\&quot;&gt;Cache-Control&lt;/a&gt;. | |
| **contentDisposition** | **string** | Specifies presentational information for the object. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1\&quot;&gt;Content-Disposition&lt;/a&gt;. | |
| **contentEncoding** | **string** | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11\&quot;&gt;Content-Encoding&lt;/a&gt;. | |
| **contentLanguage** | **string** | The language the content is in. | |
| **contentLength** | **int32** | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.13\&quot;&gt;Content-Length&lt;/a&gt;. | |
| **contentMD5** | **string** |  | |
| **contentType** | **string** | A standard MIME type describing the format of the contents. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.17\&quot;&gt;Content-Type&lt;/a&gt;. | |
| **expires** | **time.Time** | The date and time at which the object is no longer cacheable. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21\&quot;&gt;Expires&lt;/a&gt;. | |
| **xAmzServerSideEncryption** | **string** | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). | |
| **xAmzStorageClass** | **string** | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. | |
| **xAmzWebsiteRedirectLocation** | **string** | &lt;p&gt;If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata.&lt;/p&gt; &lt;p&gt;In the following example, the request header sets the redirect to an object (anotherPage.html) in the same bucket:&lt;/p&gt; &lt;p&gt; &#x60;x-amz-website-redirect-location: /anotherPage.html&#x60; &lt;/p&gt; &lt;p&gt;In the following example, the request header sets the object redirect to another website:&lt;/p&gt; &lt;p&gt; &#x60;x-amz-website-redirect-location: http://www.example.com/&#x60; &lt;/p&gt; | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm&#x60; header. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **xAmzServerSideEncryptionContext** | **string** | Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs. | |
| **xAmzServerSideEncryptionBucketKeyEnabled** | **bool** | &lt;p&gt;Specifies whether IONOS S3 Object Storage should use an S3 Bucket Key for object encryption with server-side encryption. Setting this header to &#x60;true&#x60; causes IONOS S3 Object Storage to use an S3 Bucket Key for object encryption.&lt;/p&gt; &lt;p&gt;Specifying this header with a PUT operation doesn’t affect bucket-level settings for S3 Bucket Key.&lt;/p&gt; | |
| **xAmzRequestPayer** | **string** |  | |
| **xAmzTagging** | **string** | The tag-set for the object. The tag-set must be encoded as URL Query parameters. (For example, \&quot;Key1&#x3D;Value1\&quot;) | |
| **xAmzObjectLockMode** | **string** | The Object Lock mode that you want to apply to this object. | |
| **xAmzObjectLockRetainUntilDate** | **time.Time** | The date and time when you want this object&#39;s Object Lock to expire. Must be formatted as a timestamp parameter. | |
| **xAmzObjectLockLegalHold** | **string** | Specifies whether a legal hold will be applied to this object. | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.POSTObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.POSTObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.POSTObject": {
    "port": "8443",
},
})
```


## PutObject

```go
var result  = PutObject(ctx, bucket, key)
                      .Body(body)
                      .CacheControl(cacheControl)
                      .ContentDisposition(contentDisposition)
                      .ContentEncoding(contentEncoding)
                      .ContentLanguage(contentLanguage)
                      .ContentLength(contentLength)
                      .ContentMD5(contentMD5)
                      .ContentType(contentType)
                      .Expires(expires)
                      .XAmzServerSideEncryption(xAmzServerSideEncryption)
                      .XAmzStorageClass(xAmzStorageClass)
                      .XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .XAmzServerSideEncryptionContext(xAmzServerSideEncryptionContext)
                      .XAmzRequestPayer(xAmzRequestPayer)
                      .XAmzTagging(xAmzTagging)
                      .XAmzObjectLockMode(xAmzObjectLockMode)
                      .XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate)
                      .XAmzObjectLockLegalHold(xAmzObjectLockLegalHold)
                      .XAmzMeta(xAmzMeta)
                      .Execute()
```

PutObject



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    ionoscloud "github.com/ionos-cloud/ionoscloud"
)

func main() {
    bucket := "bucket_example" // string | 
    key := "key_example" // string | Object key for which the PUT operation was initiated.
    body := os.NewFile(1234, "some_file") // *os.File | 
    cacheControl := "cacheControl_example" // string |  Can be used to specify caching behavior along the request/reply chain. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9\">Cache-Control</a>. (optional)
    contentDisposition := "contentDisposition_example" // string | Specifies presentational information for the object. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1\">Content-Disposition</a>. (optional)
    contentEncoding := "contentEncoding_example" // string | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11\">Content-Encoding</a>. (optional)
    contentLanguage := "contentLanguage_example" // string | The language the content is in. (optional)
    contentLength := int32(56) // int32 | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.13\">Content-Length</a>. (optional)
    contentMD5 := "contentMD5_example" // string |  (optional)
    contentType := "contentType_example" // string | A standard MIME type describing the format of the contents. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.17\">Content-Type</a>. (optional)
    expires := time.Now() // time.Time | The date and time at which the object is no longer cacheable. For more information, see <a href=\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21\">Expires</a>. (optional)
    xAmzServerSideEncryption := "xAmzServerSideEncryption_example" // string | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). (optional)
    xAmzStorageClass := "xAmzStorageClass_example" // string | The valid value is `STANDARD`. (optional)
    xAmzWebsiteRedirectLocation := "xAmzWebsiteRedirectLocation_example" // string | <p>If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata.</p> <p>In the following example, the request header sets the redirect to an object (anotherPage.html) in the same bucket:</p> <p> `x-amz-website-redirect-location: /anotherPage.html` </p> <p>In the following example, the request header sets the object redirect to another website:</p> <p> `x-amz-website-redirect-location: http://www.example.com/` </p> (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object. The valid option is `AES256`. (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data. For example, `4ZRNYBCCvL0YZeqo3f2+9qDyIfnLdbg5S99R2XWr0aw=`. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. For example, `bPU7G1zD2MlOi5gqnkRqZg==`. (optional)
    xAmzServerSideEncryptionContext := "xAmzServerSideEncryptionContext_example" // string | Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs. (optional)
    xAmzRequestPayer := "xAmzRequestPayer_example" // string |  (optional)
    xAmzTagging := "xAmzTagging_example" // string | The tag-set for the object. The tag-set must be encoded as URL Query parameters. (For example, \"Key1=Value1\") (optional)
    xAmzObjectLockMode := "xAmzObjectLockMode_example" // string | The Object Lock mode that you want to apply to this object. (optional)
    xAmzObjectLockRetainUntilDate := time.Now() // time.Time | The date and time when you want this object's Object Lock to expire. Must be formatted as a timestamp parameter. (optional)
    xAmzObjectLockLegalHold := "xAmzObjectLockLegalHold_example" // string | Specifies whether a legal hold will be applied to this object. (optional)
    xAmzMeta := map[string]string{"key": "Inner_example"} // map[string]string | A map of metadata to store with the object in S3. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ObjectsApi.PutObject(context.Background(), bucket, key).Body(body).CacheControl(cacheControl).ContentDisposition(contentDisposition).ContentEncoding(contentEncoding).ContentLanguage(contentLanguage).ContentLength(contentLength).ContentMD5(contentMD5).ContentType(contentType).Expires(expires).XAmzServerSideEncryption(xAmzServerSideEncryption).XAmzStorageClass(xAmzStorageClass).XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).XAmzServerSideEncryptionContext(xAmzServerSideEncryptionContext).XAmzRequestPayer(xAmzRequestPayer).XAmzTagging(xAmzTagging).XAmzObjectLockMode(xAmzObjectLockMode).XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate).XAmzObjectLockLegalHold(xAmzObjectLockLegalHold).XAmzMeta(xAmzMeta).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ObjectsApi.PutObject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the PUT operation was initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiPutObjectRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | ***os.File** |  | |
| **cacheControl** | **string** |  Can be used to specify caching behavior along the request/reply chain. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9\&quot;&gt;Cache-Control&lt;/a&gt;. | |
| **contentDisposition** | **string** | Specifies presentational information for the object. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1\&quot;&gt;Content-Disposition&lt;/a&gt;. | |
| **contentEncoding** | **string** | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11\&quot;&gt;Content-Encoding&lt;/a&gt;. | |
| **contentLanguage** | **string** | The language the content is in. | |
| **contentLength** | **int32** | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.13\&quot;&gt;Content-Length&lt;/a&gt;. | |
| **contentMD5** | **string** |  | |
| **contentType** | **string** | A standard MIME type describing the format of the contents. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.17\&quot;&gt;Content-Type&lt;/a&gt;. | |
| **expires** | **time.Time** | The date and time at which the object is no longer cacheable. For more information, see &lt;a href&#x3D;\&quot;http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21\&quot;&gt;Expires&lt;/a&gt;. | |
| **xAmzServerSideEncryption** | **string** | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). | |
| **xAmzStorageClass** | **string** | The valid value is &#x60;STANDARD&#x60;. | |
| **xAmzWebsiteRedirectLocation** | **string** | &lt;p&gt;If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata.&lt;/p&gt; &lt;p&gt;In the following example, the request header sets the redirect to an object (anotherPage.html) in the same bucket:&lt;/p&gt; &lt;p&gt; &#x60;x-amz-website-redirect-location: /anotherPage.html&#x60; &lt;/p&gt; &lt;p&gt;In the following example, the request header sets the object redirect to another website:&lt;/p&gt; &lt;p&gt; &#x60;x-amz-website-redirect-location: http://www.example.com/&#x60; &lt;/p&gt; | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object. The valid option is &#x60;AES256&#x60;. | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data. For example, &#x60;4ZRNYBCCvL0YZeqo3f2+9qDyIfnLdbg5S99R2XWr0aw&#x3D;&#x60;. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. For example, &#x60;bPU7G1zD2MlOi5gqnkRqZg&#x3D;&#x3D;&#x60;. | |
| **xAmzServerSideEncryptionContext** | **string** | Specifies the IONOS S3 Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs. | |
| **xAmzRequestPayer** | **string** |  | |
| **xAmzTagging** | **string** | The tag-set for the object. The tag-set must be encoded as URL Query parameters. (For example, \&quot;Key1&#x3D;Value1\&quot;) | |
| **xAmzObjectLockMode** | **string** | The Object Lock mode that you want to apply to this object. | |
| **xAmzObjectLockRetainUntilDate** | **time.Time** | The date and time when you want this object&#39;s Object Lock to expire. Must be formatted as a timestamp parameter. | |
| **xAmzObjectLockLegalHold** | **string** | Specifies whether a legal hold will be applied to this object. | |
| **xAmzMeta** | [**map[string]string**](../models/string.md) | A map of metadata to store with the object in S3. | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: text/plain
- **Accept**: text/plain


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ObjectsApiService.PutObject"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ObjectsApiService.PutObject": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ObjectsApiService.PutObject": {
    "port": "8443",
},
})
```

