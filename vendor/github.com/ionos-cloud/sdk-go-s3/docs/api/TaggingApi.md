# \TaggingApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**DeleteBucketTagging**](TaggingApi.md#DeleteBucketTagging) | **Delete** /{Bucket}?tagging | DeleteBucketTagging|
|[**DeleteObjectTagging**](TaggingApi.md#DeleteObjectTagging) | **Delete** /{Bucket}/{Key}?tagging | DeleteObjectTagging|
|[**GetBucketTagging**](TaggingApi.md#GetBucketTagging) | **Get** /{Bucket}?tagging | GetBucketTagging|
|[**GetObjectTagging**](TaggingApi.md#GetObjectTagging) | **Get** /{Bucket}/{Key}?tagging | GetObjectTagging|
|[**PutBucketTagging**](TaggingApi.md#PutBucketTagging) | **Put** /{Bucket}?tagging | PutBucketTagging|
|[**PutObjectTagging**](TaggingApi.md#PutObjectTagging) | **Put** /{Bucket}/{Key}?tagging | PutObjectTagging|



## DeleteBucketTagging

```go
var result  = DeleteBucketTagging(ctx, bucket)
                      .Execute()
```

DeleteBucketTagging



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

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.DeleteBucketTagging(context.Background(), bucket).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.DeleteBucketTagging``: %v\n", err)
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

Other parameters are passed through a pointer to an apiDeleteBucketTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.DeleteBucketTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.DeleteBucketTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.DeleteBucketTagging": {
    "port": "8443",
},
})
```


## DeleteObjectTagging

```go
var result map[string]interface{} = DeleteObjectTagging(ctx, bucket, key)
                      .VersionId(versionId)
                      .Execute()
```

DeleteObjectTagging



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
    key := "key_example" // string | The key that identifies the object in the bucket from which to remove all tags.
    versionId := "versionId_example" // string | The versionId of the object that the tag-set will be removed from. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.DeleteObjectTagging(context.Background(), bucket, key).VersionId(versionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.DeleteObjectTagging``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `DeleteObjectTagging`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `TaggingApi.DeleteObjectTagging`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | The key that identifies the object in the bucket from which to remove all tags. | |

### Other Parameters

Other parameters are passed through a pointer to an apiDeleteObjectTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **versionId** | **string** | The versionId of the object that the tag-set will be removed from. | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.DeleteObjectTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.DeleteObjectTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.DeleteObjectTagging": {
    "port": "8443",
},
})
```


## GetBucketTagging

```go
var result GetBucketTaggingOutput = GetBucketTagging(ctx, bucket)
                      .Execute()
```

GetBucketTagging



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

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.GetBucketTagging(context.Background(), bucket).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.GetBucketTagging``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetBucketTagging`: GetBucketTaggingOutput
    fmt.Fprintf(os.Stdout, "Response from `TaggingApi.GetBucketTagging`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetBucketTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|

### Return type

[**GetBucketTaggingOutput**](../models/GetBucketTaggingOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.GetBucketTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.GetBucketTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.GetBucketTagging": {
    "port": "8443",
},
})
```


## GetObjectTagging

```go
var result GetObjectTaggingOutput = GetObjectTagging(ctx, bucket, key)
                      .VersionId(versionId)
                      .Execute()
```

GetObjectTagging



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
    key := "key_example" // string | Object key for which to get the tagging information.
    versionId := "versionId_example" // string | The versionId of the object for which to get the tagging information. (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.GetObjectTagging(context.Background(), bucket, key).VersionId(versionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.GetObjectTagging``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetObjectTagging`: GetObjectTaggingOutput
    fmt.Fprintf(os.Stdout, "Response from `TaggingApi.GetObjectTagging`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Object key for which to get the tagging information. | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetObjectTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **versionId** | **string** | The versionId of the object for which to get the tagging information. | |

### Return type

[**GetObjectTaggingOutput**](../models/GetObjectTaggingOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.GetObjectTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.GetObjectTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.GetObjectTagging": {
    "port": "8443",
},
})
```


## PutBucketTagging

```go
var result  = PutBucketTagging(ctx, bucket)
                      .PutBucketTaggingRequest(putBucketTaggingRequest)
                      .ContentMD5(contentMD5)
                      .Execute()
```

PutBucketTagging



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
    putBucketTaggingRequest := *openapiclient.NewPutBucketTaggingRequest(*openapiclient.NewPutBucketTaggingRequestTagging()) // PutBucketTaggingRequest | 
    contentMD5 := "contentMD5_example" // string |  (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.PutBucketTagging(context.Background(), bucket).PutBucketTaggingRequest(putBucketTaggingRequest).ContentMD5(contentMD5).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.PutBucketTagging``: %v\n", err)
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

Other parameters are passed through a pointer to an apiPutBucketTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **putBucketTaggingRequest** | [**PutBucketTaggingRequest**](../models/PutBucketTaggingRequest.md) |  | |
| **contentMD5** | **string** |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: Not defined


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.PutBucketTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.PutBucketTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.PutBucketTagging": {
    "port": "8443",
},
})
```


## PutObjectTagging

```go
var result map[string]interface{} = PutObjectTagging(ctx, bucket, key)
                      .PutObjectTaggingRequest(putObjectTaggingRequest)
                      .VersionId(versionId)
                      .ContentMD5(contentMD5)
                      .Execute()
```

PutObjectTagging



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
    key := "key_example" // string | Name of the object key.
    putObjectTaggingRequest := *openapiclient.NewPutObjectTaggingRequest(*openapiclient.NewPutObjectTaggingRequestTagging()) // PutObjectTaggingRequest | 
    versionId := "versionId_example" // string | The versionId of the object that the tag-set will be added to. (optional)
    contentMD5 := "contentMD5_example" // string |  (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.TaggingApi.PutObjectTagging(context.Background(), bucket, key).PutObjectTaggingRequest(putObjectTaggingRequest).VersionId(versionId).ContentMD5(contentMD5).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TaggingApi.PutObjectTagging``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `PutObjectTagging`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `TaggingApi.PutObjectTagging`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |
|**key** | **string** | Name of the object key. | |

### Other Parameters

Other parameters are passed through a pointer to an apiPutObjectTaggingRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **putObjectTaggingRequest** | [**PutObjectTaggingRequest**](../models/PutObjectTaggingRequest.md) |  | |
| **versionId** | **string** | The versionId of the object that the tag-set will be added to. | |
| **contentMD5** | **string** |  | |

### Return type

**map[string]interface{}**

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"TaggingApiService.PutObjectTagging"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "TaggingApiService.PutObjectTagging": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "TaggingApiService.PutObjectTagging": {
    "port": "8443",
},
})
```

