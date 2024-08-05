# \VersionsApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**ListObjectVersions**](VersionsApi.md#ListObjectVersions) | **Get** /{Bucket}?versions | ListObjectVersions|



## ListObjectVersions

```go
var result ListObjectVersionsOutput = ListObjectVersions(ctx, bucket)
                      .Versions(versions)
                      .Delimiter(delimiter)
                      .EncodingType(encodingType)
                      .KeyMarker(keyMarker)
                      .MaxKeys(maxKeys)
                      .Prefix(prefix)
                      .VersionIdMarker(versionIdMarker)
                      .MaxKeys2(maxKeys2)
                      .KeyMarker2(keyMarker2)
                      .VersionIdMarker2(versionIdMarker2)
                      .Execute()
```

ListObjectVersions



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
    versions := true // bool | 
    delimiter := "delimiter_example" // string | A delimiter is a character that you specify to group keys. All keys that contain the same string between the `prefix` and the first occurrence of the delimiter are grouped under a single result element in CommonPrefixes. These groups are counted as one result against the max-keys limitation. These keys are not returned elsewhere in the response. (optional)
    encodingType := "encodingType_example" // string |  (optional)
    keyMarker := "keyMarker_example" // string | Specifies the key to start with when listing objects in a bucket. (optional)
    maxKeys := int32(56) // int32 | Sets the maximum number of keys returned in the response. By default the operation returns up to 1,000 key names. The response might contain fewer keys but will never contain more. If additional keys satisfy the search criteria, but were not returned because max-keys was exceeded, the response contains &lt;isTruncated&gt;true&lt;/isTruncated&gt;. To return the additional keys, see key-marker and version-id-marker. (optional)
    prefix := "prefix_example" // string | Use this parameter to select only those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different groupings of keys. (You can think of using prefix to make groups in the same way you'd use a folder in a file system.) You can use prefix with delimiter to roll up numerous objects into a single result under CommonPrefixes.  (optional)
    versionIdMarker := "versionIdMarker_example" // string | Specifies the object version you want to start listing from. (optional)
    maxKeys2 := "maxKeys_example" // string | Pagination limit (optional)
    keyMarker2 := "keyMarker_example" // string | Pagination token (optional)
    versionIdMarker2 := "versionIdMarker_example" // string | Pagination token (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.VersionsApi.ListObjectVersions(context.Background(), bucket).Versions(versions).Delimiter(delimiter).EncodingType(encodingType).KeyMarker(keyMarker).MaxKeys(maxKeys).Prefix(prefix).VersionIdMarker(versionIdMarker).MaxKeys2(maxKeys2).KeyMarker2(keyMarker2).VersionIdMarker2(versionIdMarker2).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `VersionsApi.ListObjectVersions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `ListObjectVersions`: ListObjectVersionsOutput
    fmt.Fprintf(os.Stdout, "Response from `VersionsApi.ListObjectVersions`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiListObjectVersionsRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **versions** | **bool** |  | |
| **delimiter** | **string** | A delimiter is a character that you specify to group keys. All keys that contain the same string between the &#x60;prefix&#x60; and the first occurrence of the delimiter are grouped under a single result element in CommonPrefixes. These groups are counted as one result against the max-keys limitation. These keys are not returned elsewhere in the response. | |
| **encodingType** | **string** |  | |
| **keyMarker** | **string** | Specifies the key to start with when listing objects in a bucket. | |
| **maxKeys** | **int32** | Sets the maximum number of keys returned in the response. By default the operation returns up to 1,000 key names. The response might contain fewer keys but will never contain more. If additional keys satisfy the search criteria, but were not returned because max-keys was exceeded, the response contains &amp;lt;isTruncated&amp;gt;true&amp;lt;/isTruncated&amp;gt;. To return the additional keys, see key-marker and version-id-marker. | |
| **prefix** | **string** | Use this parameter to select only those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different groupings of keys. (You can think of using prefix to make groups in the same way you&#39;d use a folder in a file system.) You can use prefix with delimiter to roll up numerous objects into a single result under CommonPrefixes.  | |
| **versionIdMarker** | **string** | Specifies the object version you want to start listing from. | |
| **maxKeys2** | **string** | Pagination limit | |
| **keyMarker2** | **string** | Pagination token | |
| **versionIdMarker2** | **string** | Pagination token | |

### Return type

[**ListObjectVersionsOutput**](../models/ListObjectVersionsOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"VersionsApiService.ListObjectVersions"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "VersionsApiService.ListObjectVersions": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "VersionsApiService.ListObjectVersions": {
    "port": "8443",
},
})
```

