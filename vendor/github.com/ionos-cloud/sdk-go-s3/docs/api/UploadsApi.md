# \UploadsApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**AbortMultipartUpload**](UploadsApi.md#AbortMultipartUpload) | **Delete** /{Bucket}/{Key}?uploadId | AbortMultipartUpload|
|[**CompleteMultipartUpload**](UploadsApi.md#CompleteMultipartUpload) | **Post** /{Bucket}/{Key}?uploadId | CompleteMultipartUpload|
|[**CreateMultipartUpload**](UploadsApi.md#CreateMultipartUpload) | **Post** /{Bucket}/{Key}?uploads | CreateMultipartUpload|
|[**ListMultipartUploads**](UploadsApi.md#ListMultipartUploads) | **Get** /{Bucket}?uploads | ListMultipartUploads|
|[**ListParts**](UploadsApi.md#ListParts) | **Get** /{Bucket}/{Key}?uploadId | ListParts|
|[**UploadPart**](UploadsApi.md#UploadPart) | **Put** /{Bucket}/{Key}?uploadId | UploadPart|
|[**UploadPartCopy**](UploadsApi.md#UploadPartCopy) | **Put** /{Bucket}/{Key}?x-amz-copy-source&amp;partNumber&amp;uploadId | UploadPartCopy|



## AbortMultipartUpload

```go
var result map[string]interface{} = AbortMultipartUpload(ctx, bucket, key)
                      .UploadId(uploadId)
                      .Execute()
```

AbortMultipartUpload



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
    key := "key_example" // string | Key of the object for which the multipart upload was initiated. <p> **Possible values:** length ≥ 1 </p>
    uploadId := "uploadId_example" // string | Upload ID that identifies the multipart upload.

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.AbortMultipartUpload(context.Background(), bucket, key).UploadId(uploadId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.AbortMultipartUpload``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `AbortMultipartUpload`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.AbortMultipartUpload`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Key of the object for which the multipart upload was initiated. &lt;p&gt; **Possible values:** length ≥ 1 &lt;/p&gt; | |

### Other Parameters

Other parameters are passed through a pointer to an apiAbortMultipartUploadRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **uploadId** | **string** | Upload ID that identifies the multipart upload. | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xmlaplication/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.AbortMultipartUpload"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.AbortMultipartUpload": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.AbortMultipartUpload": {
    "port": "8443",
},
})
```


## CompleteMultipartUpload

```go
var result CompleteMultipartUploadOutput = CompleteMultipartUpload(ctx, bucket, key)
                      .UploadId(uploadId)
                      .Example(example)
                      .Execute()
```

CompleteMultipartUpload



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
    key := "key_example" // string | Object key for which the multipart upload was initiated.
    uploadId := "uploadId_example" // string | ID for the initiated multipart upload.
    example := *openapiclient.NewExample() // Example | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.CompleteMultipartUpload(context.Background(), bucket, key).UploadId(uploadId).Example(example).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.CompleteMultipartUpload``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `CompleteMultipartUpload`: CompleteMultipartUploadOutput
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.CompleteMultipartUpload`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the multipart upload was initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiCompleteMultipartUploadRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **uploadId** | **string** | ID for the initiated multipart upload. | |
| **example** | [**Example**](../models/Example.md) |  | |

### Return type

[**CompleteMultipartUploadOutput**](../models/CompleteMultipartUploadOutput.md)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.CompleteMultipartUpload"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.CompleteMultipartUpload": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.CompleteMultipartUpload": {
    "port": "8443",
},
})
```


## CreateMultipartUpload

```go
var result CreateMultipartUploadOutput = CreateMultipartUpload(ctx, bucket, key)
                      .Uploads(uploads)
                      .CacheControl(cacheControl)
                      .ContentDisposition(contentDisposition)
                      .ContentEncoding(contentEncoding)
                      .ContentType(contentType)
                      .Expires(expires)
                      .XAmzServerSideEncryption(xAmzServerSideEncryption)
                      .XAmzStorageClass(xAmzStorageClass)
                      .XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .XAmzObjectLockMode(xAmzObjectLockMode)
                      .XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate)
                      .XAmzObjectLockLegalHold(xAmzObjectLockLegalHold)
                      .XAmzMeta(xAmzMeta)
                      .Execute()
```

CreateMultipartUpload



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
    key := "key_example" // string | Object key for which the multipart upload is to be initiated.
    uploads := true // bool | 
    cacheControl := "cacheControl_example" // string | Specifies caching behavior along the request/reply chain. (optional)
    contentDisposition := "contentDisposition_example" // string | Specifies presentational information for the object. (optional)
    contentEncoding := "contentEncoding_example" // string | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. (optional)
    contentType := "contentType_example" // string | A standard MIME type describing the format of the object data. (optional)
    expires := time.Now() // time.Time | The date and time at which the object is no longer cacheable. (optional)
    xAmzServerSideEncryption := "xAmzServerSideEncryption_example" // string | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). (optional)
    xAmzStorageClass := "xAmzStorageClass_example" // string | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. (optional)
    xAmzWebsiteRedirectLocation := "xAmzWebsiteRedirectLocation_example" // string | If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata. (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm` header. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)
    xAmzObjectLockMode := "xAmzObjectLockMode_example" // string | Specifies the Object Lock mode that you want to apply to the uploaded object. (optional)
    xAmzObjectLockRetainUntilDate := time.Now() // time.Time | Specifies the date and time when you want the Object Lock to expire. (optional)
    xAmzObjectLockLegalHold := "xAmzObjectLockLegalHold_example" // string | Specifies whether you want to apply a Legal Hold to the uploaded object. (optional)
    xAmzMeta := map[string]string{"key": "Inner_example"} // map[string]string | A map of metadata to store with the object in S3. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.CreateMultipartUpload(context.Background(), bucket, key).Uploads(uploads).CacheControl(cacheControl).ContentDisposition(contentDisposition).ContentEncoding(contentEncoding).ContentType(contentType).Expires(expires).XAmzServerSideEncryption(xAmzServerSideEncryption).XAmzStorageClass(xAmzStorageClass).XAmzWebsiteRedirectLocation(xAmzWebsiteRedirectLocation).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).XAmzObjectLockMode(xAmzObjectLockMode).XAmzObjectLockRetainUntilDate(xAmzObjectLockRetainUntilDate).XAmzObjectLockLegalHold(xAmzObjectLockLegalHold).XAmzMeta(xAmzMeta).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.CreateMultipartUpload``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `CreateMultipartUpload`: CreateMultipartUploadOutput
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.CreateMultipartUpload`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the multipart upload is to be initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiCreateMultipartUploadRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **uploads** | **bool** |  | |
| **cacheControl** | **string** | Specifies caching behavior along the request/reply chain. | |
| **contentDisposition** | **string** | Specifies presentational information for the object. | |
| **contentEncoding** | **string** | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. | |
| **contentType** | **string** | A standard MIME type describing the format of the object data. | |
| **expires** | **time.Time** | The date and time at which the object is no longer cacheable. | |
| **xAmzServerSideEncryption** | **string** | The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage (AES256). | |
| **xAmzStorageClass** | **string** | IONOS S3 Object Storage uses the STANDARD Storage Class to store newly created objects. The STANDARD storage class provides high durability and high availability. | |
| **xAmzWebsiteRedirectLocation** | **string** | If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS S3 Object Storage stores the value of this header in the object metadata. | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm&#x60; header. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |
| **xAmzObjectLockMode** | **string** | Specifies the Object Lock mode that you want to apply to the uploaded object. | |
| **xAmzObjectLockRetainUntilDate** | **time.Time** | Specifies the date and time when you want the Object Lock to expire. | |
| **xAmzObjectLockLegalHold** | **string** | Specifies whether you want to apply a Legal Hold to the uploaded object. | |
| **xAmzMeta** | [**map[string]string**](../models/string.md) | A map of metadata to store with the object in S3. | |

### Return type

[**CreateMultipartUploadOutput**](../models/CreateMultipartUploadOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.CreateMultipartUpload"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.CreateMultipartUpload": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.CreateMultipartUpload": {
    "port": "8443",
},
})
```


## ListMultipartUploads

```go
var result ListMultipartUploadsOutput = ListMultipartUploads(ctx, bucket)
                      .Uploads(uploads)
                      .Delimiter(delimiter)
                      .EncodingType(encodingType)
                      .KeyMarker(keyMarker)
                      .MaxUploads(maxUploads)
                      .Prefix(prefix)
                      .UploadIdMarker(uploadIdMarker)
                      .MaxUploads2(maxUploads2)
                      .KeyMarker2(keyMarker2)
                      .UploadIdMarker2(uploadIdMarker2)
                      .Execute()
```

ListMultipartUploads



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
    uploads := true // bool | 
    delimiter := "delimiter_example" // string | <p>Character you use to group keys.</p> <p>All keys that contain the same string between the prefix, if specified, and the first occurrence of the delimiter after the prefix are grouped under a single result element, `CommonPrefixes`. If you don't specify the prefix parameter, then the substring starts at the beginning of the key. The keys that are grouped under `CommonPrefixes` result element are not returned elsewhere in the response.</p> (optional)
    encodingType := "encodingType_example" // string |  (optional)
    keyMarker := "keyMarker_example" // string | <p>Together with upload-id-marker, this parameter specifies the multipart upload after which listing should begin.</p> <p>If `upload-id-marker` is not specified, only the keys lexicographically greater than the specified `key-marker` will be included in the list.</p> <p>If `upload-id-marker` is specified, any multipart uploads for a key equal to the `key-marker` might also be included, provided those multipart uploads have upload IDs lexicographically greater than the specified `upload-id-marker`.</p> (optional)
    maxUploads := int32(56) // int32 | Sets the maximum number of multipart uploads, from 1 to 1,000, to return in the response body. 1,000 is the maximum number of uploads that can be returned in a response. (optional)
    prefix := "prefix_example" // string | Lists in-progress uploads only for those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different grouping of keys. (You can think of using prefix to make groups in the same way you'd use a folder in a file system.) (optional)
    uploadIdMarker := "uploadIdMarker_example" // string | Together with key-marker, specifies the multipart upload after which listing should begin. If key-marker is not specified, the upload-id-marker parameter is ignored. Otherwise, any multipart uploads for a key equal to the key-marker might be included in the list only if they have an upload ID lexicographically greater than the specified `upload-id-marker`. (optional)
    maxUploads2 := "maxUploads_example" // string | Pagination limit (optional)
    keyMarker2 := "keyMarker_example" // string | Pagination token (optional)
    uploadIdMarker2 := "uploadIdMarker_example" // string | Pagination token (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.ListMultipartUploads(context.Background(), bucket).Uploads(uploads).Delimiter(delimiter).EncodingType(encodingType).KeyMarker(keyMarker).MaxUploads(maxUploads).Prefix(prefix).UploadIdMarker(uploadIdMarker).MaxUploads2(maxUploads2).KeyMarker2(keyMarker2).UploadIdMarker2(uploadIdMarker2).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.ListMultipartUploads``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `ListMultipartUploads`: ListMultipartUploadsOutput
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.ListMultipartUploads`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiListMultipartUploadsRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **uploads** | **bool** |  | |
| **delimiter** | **string** | &lt;p&gt;Character you use to group keys.&lt;/p&gt; &lt;p&gt;All keys that contain the same string between the prefix, if specified, and the first occurrence of the delimiter after the prefix are grouped under a single result element, &#x60;CommonPrefixes&#x60;. If you don&#39;t specify the prefix parameter, then the substring starts at the beginning of the key. The keys that are grouped under &#x60;CommonPrefixes&#x60; result element are not returned elsewhere in the response.&lt;/p&gt; | |
| **encodingType** | **string** |  | |
| **keyMarker** | **string** | &lt;p&gt;Together with upload-id-marker, this parameter specifies the multipart upload after which listing should begin.&lt;/p&gt; &lt;p&gt;If &#x60;upload-id-marker&#x60; is not specified, only the keys lexicographically greater than the specified &#x60;key-marker&#x60; will be included in the list.&lt;/p&gt; &lt;p&gt;If &#x60;upload-id-marker&#x60; is specified, any multipart uploads for a key equal to the &#x60;key-marker&#x60; might also be included, provided those multipart uploads have upload IDs lexicographically greater than the specified &#x60;upload-id-marker&#x60;.&lt;/p&gt; | |
| **maxUploads** | **int32** | Sets the maximum number of multipart uploads, from 1 to 1,000, to return in the response body. 1,000 is the maximum number of uploads that can be returned in a response. | |
| **prefix** | **string** | Lists in-progress uploads only for those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different grouping of keys. (You can think of using prefix to make groups in the same way you&#39;d use a folder in a file system.) | |
| **uploadIdMarker** | **string** | Together with key-marker, specifies the multipart upload after which listing should begin. If key-marker is not specified, the upload-id-marker parameter is ignored. Otherwise, any multipart uploads for a key equal to the key-marker might be included in the list only if they have an upload ID lexicographically greater than the specified &#x60;upload-id-marker&#x60;. | |
| **maxUploads2** | **string** | Pagination limit | |
| **keyMarker2** | **string** | Pagination token | |
| **uploadIdMarker2** | **string** | Pagination token | |

### Return type

[**ListMultipartUploadsOutput**](../models/ListMultipartUploadsOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.ListMultipartUploads"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.ListMultipartUploads": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.ListMultipartUploads": {
    "port": "8443",
},
})
```


## ListParts

```go
var result ListPartsOutput = ListParts(ctx, bucket, key)
                      .UploadId(uploadId)
                      .MaxParts(maxParts)
                      .PartNumberMarker(partNumberMarker)
                      .PartNumberMarker2(partNumberMarker2)
                      .Execute()
```

ListParts



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
    key := "key_example" // string | Object key for which the multipart upload was initiated.
    uploadId := "uploadId_example" // string | Upload ID identifying the multipart upload whose parts are being listed.
    maxParts := int32(56) // int32 | Sets the maximum number of parts to return. (optional)
    partNumberMarker := int32(56) // int32 | Specifies the part after which listing should begin. Only parts with higher part numbers will be listed. (optional)
    partNumberMarker2 := "partNumberMarker_example" // string | Pagination token (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.ListParts(context.Background(), bucket, key).UploadId(uploadId).MaxParts(maxParts).PartNumberMarker(partNumberMarker).PartNumberMarker2(partNumberMarker2).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.ListParts``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `ListParts`: ListPartsOutput
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.ListParts`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the multipart upload was initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiListPartsRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **uploadId** | **string** | Upload ID identifying the multipart upload whose parts are being listed. | |
| **maxParts** | **int32** | Sets the maximum number of parts to return. | |
| **partNumberMarker** | **int32** | Specifies the part after which listing should begin. Only parts with higher part numbers will be listed. | |
| **partNumberMarker2** | **string** | Pagination token | |

### Return type

[**ListPartsOutput**](../models/ListPartsOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.ListParts"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.ListParts": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.ListParts": {
    "port": "8443",
},
})
```


## UploadPart

```go
var result map[string]interface{} = UploadPart(ctx, bucket, key)
                      .PartNumber(partNumber)
                      .UploadId(uploadId)
                      .UploadPartRequest(uploadPartRequest)
                      .ContentLength(contentLength)
                      .ContentMD5(contentMD5)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey)
                      .XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5)
                      .Execute()
```

UploadPart



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
    key := "key_example" // string | Object key for which the multipart upload was initiated.
    partNumber := int32(56) // int32 | Part number of part being uploaded. This is a positive integer between 1 and 10,000.
    uploadId := "uploadId_example" // string | Upload ID identifying the multipart upload whose part is being uploaded.
    uploadPartRequest := *openapiclient.NewUploadPartRequest() // UploadPartRequest | 
    contentLength := int32(56) // int32 | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. (optional)
    contentMD5 := "contentMD5_example" // string |  (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)
    xAmzServerSideEncryptionCustomerKey := "xAmzServerSideEncryptionCustomerKey_example" // string | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the `x-amz-server-side-encryption-customer-algorithm header`. This must be the same encryption key specified in the initiate multipart upload request. (optional)
    xAmzServerSideEncryptionCustomerKeyMD5 := "xAmzServerSideEncryptionCustomerKeyMD5_example" // string | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.UploadPart(context.Background(), bucket, key).PartNumber(partNumber).UploadId(uploadId).UploadPartRequest(uploadPartRequest).ContentLength(contentLength).ContentMD5(contentMD5).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).XAmzServerSideEncryptionCustomerKey(xAmzServerSideEncryptionCustomerKey).XAmzServerSideEncryptionCustomerKeyMD5(xAmzServerSideEncryptionCustomerKeyMD5).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.UploadPart``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `UploadPart`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.UploadPart`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the multipart upload was initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiUploadPartRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **partNumber** | **int32** | Part number of part being uploaded. This is a positive integer between 1 and 10,000. | |
| **uploadId** | **string** | Upload ID identifying the multipart upload whose part is being uploaded. | |
| **uploadPartRequest** | [**UploadPartRequest**](../models/UploadPartRequest.md) |  | |
| **contentLength** | **int32** | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. | |
| **contentMD5** | **string** |  | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |
| **xAmzServerSideEncryptionCustomerKey** | **string** | Specifies the customer-provided encryption key for IONOS S3 Object Storage to use in encrypting data. This value is used to store the object and then it is discarded; IONOS S3 Object Storage does not store the encryption key. The key must be appropriate for use with the algorithm specified in the &#x60;x-amz-server-side-encryption-customer-algorithm header&#x60;. This must be the same encryption key specified in the initiate multipart upload request. | |
| **xAmzServerSideEncryptionCustomerKeyMD5** | **string** | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Storage uses this header for a message integrity check to ensure that the encryption key was transmitted without error. | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.UploadPart"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.UploadPart": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.UploadPart": {
    "port": "8443",
},
})
```


## UploadPartCopy

```go
var result UploadPartCopyOutput = UploadPartCopy(ctx, bucket, key)
                      .XAmzCopySource(xAmzCopySource)
                      .PartNumber(partNumber)
                      .UploadId(uploadId)
                      .XAmzCopySourceIfMatch(xAmzCopySourceIfMatch)
                      .XAmzCopySourceIfModifiedSince(xAmzCopySourceIfModifiedSince)
                      .XAmzCopySourceIfNoneMatch(xAmzCopySourceIfNoneMatch)
                      .XAmzCopySourceIfUnmodifiedSince(xAmzCopySourceIfUnmodifiedSince)
                      .XAmzCopySourceRange(xAmzCopySourceRange)
                      .XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm)
                      .Execute()
```

UploadPartCopy



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
    xAmzCopySource := "xAmzCopySource_example" // string | <p>Specifies the source object for the copy operation. </p>
    key := "key_example" // string | Object key for which the multipart upload was initiated.
    partNumber := int32(56) // int32 | Part number of part being copied. This is a positive integer between 1 and 10,000.
    uploadId := "uploadId_example" // string | Upload ID identifying the multipart upload whose part is being copied.
    xAmzCopySourceIfMatch := "xAmzCopySourceIfMatch_example" // string | Copies the object if its entity tag (ETag) matches the specified tag. (optional)
    xAmzCopySourceIfModifiedSince := time.Now() // time.Time | Copies the object if it has been modified since the specified time. (optional)
    xAmzCopySourceIfNoneMatch := "xAmzCopySourceIfNoneMatch_example" // string | Copies the object if its entity tag (ETag) is different than the specified ETag. (optional)
    xAmzCopySourceIfUnmodifiedSince := time.Now() // time.Time | Copies the object if it hasn't been modified since the specified time. (optional)
    xAmzCopySourceRange := "xAmzCopySourceRange_example" // string | The range of bytes to copy from the source object. The range value must use the form bytes=first-last, where the first and last are the zero-based byte offsets to copy. For example, bytes=0-9 indicates that you want to copy the first 10 bytes of the source. You can copy a range only if the source object is greater than 5 MB. (optional)
    xAmzServerSideEncryptionCustomerAlgorithm := "xAmzServerSideEncryptionCustomerAlgorithm_example" // string | Specifies the algorithm to use to when encrypting the object (AES256). (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.UploadsApi.UploadPartCopy(context.Background(), bucket, key).XAmzCopySource(xAmzCopySource).PartNumber(partNumber).UploadId(uploadId).XAmzCopySourceIfMatch(xAmzCopySourceIfMatch).XAmzCopySourceIfModifiedSince(xAmzCopySourceIfModifiedSince).XAmzCopySourceIfNoneMatch(xAmzCopySourceIfNoneMatch).XAmzCopySourceIfUnmodifiedSince(xAmzCopySourceIfUnmodifiedSince).XAmzCopySourceRange(xAmzCopySourceRange).XAmzServerSideEncryptionCustomerAlgorithm(xAmzServerSideEncryptionCustomerAlgorithm).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UploadsApi.UploadPartCopy``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `UploadPartCopy`: UploadPartCopyOutput
    fmt.Fprintf(os.Stdout, "Response from `UploadsApi.UploadPartCopy`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which the multipart upload was initiated. | |

### Other Parameters

Other parameters are passed through a pointer to an apiUploadPartCopyRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **xAmzCopySource** | **string** | &lt;p&gt;Specifies the source object for the copy operation. &lt;/p&gt; | |
| **partNumber** | **int32** | Part number of part being copied. This is a positive integer between 1 and 10,000. | |
| **uploadId** | **string** | Upload ID identifying the multipart upload whose part is being copied. | |
| **xAmzCopySourceIfMatch** | **string** | Copies the object if its entity tag (ETag) matches the specified tag. | |
| **xAmzCopySourceIfModifiedSince** | **time.Time** | Copies the object if it has been modified since the specified time. | |
| **xAmzCopySourceIfNoneMatch** | **string** | Copies the object if its entity tag (ETag) is different than the specified ETag. | |
| **xAmzCopySourceIfUnmodifiedSince** | **time.Time** | Copies the object if it hasn&#39;t been modified since the specified time. | |
| **xAmzCopySourceRange** | **string** | The range of bytes to copy from the source object. The range value must use the form bytes&#x3D;first-last, where the first and last are the zero-based byte offsets to copy. For example, bytes&#x3D;0-9 indicates that you want to copy the first 10 bytes of the source. You can copy a range only if the source object is greater than 5 MB. | |
| **xAmzServerSideEncryptionCustomerAlgorithm** | **string** | Specifies the algorithm to use to when encrypting the object (AES256). | |

### Return type

[**UploadPartCopyOutput**](../models/UploadPartCopyOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"UploadsApiService.UploadPartCopy"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "UploadsApiService.UploadPartCopy": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "UploadsApiService.UploadPartCopy": {
    "port": "8443",
},
})
```

