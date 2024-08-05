# \LifecycleApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**DeleteBucketLifecycle**](LifecycleApi.md#DeleteBucketLifecycle) | **Delete** /{Bucket}?lifecycle | DeleteBucketLifecycle|
|[**GetBucketLifecycle**](LifecycleApi.md#GetBucketLifecycle) | **Get** /{Bucket}?lifecycle | GetBucketLifecycle|
|[**PutBucketLifecycle**](LifecycleApi.md#PutBucketLifecycle) | **Put** /{Bucket}?lifecycle | PutBucketLifecycle|



## DeleteBucketLifecycle

```go
var result  = DeleteBucketLifecycle(ctx, bucket)
                      .Lifecycle(lifecycle)
                      .Execute()
```

DeleteBucketLifecycle



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
    lifecycle := true // bool | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.LifecycleApi.DeleteBucketLifecycle(context.Background(), bucket).Lifecycle(lifecycle).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `LifecycleApi.DeleteBucketLifecycle``: %v\n", err)
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

Other parameters are passed through a pointer to an apiDeleteBucketLifecycleRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **lifecycle** | **bool** |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"LifecycleApiService.DeleteBucketLifecycle"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "LifecycleApiService.DeleteBucketLifecycle": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "LifecycleApiService.DeleteBucketLifecycle": {
    "port": "8443",
},
})
```


## GetBucketLifecycle

```go
var result GetBucketLifecycleOutput = GetBucketLifecycle(ctx, bucket)
                      .Lifecycle(lifecycle)
                      .Execute()
```

GetBucketLifecycle



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
    lifecycle := true // bool | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.LifecycleApi.GetBucketLifecycle(context.Background(), bucket).Lifecycle(lifecycle).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `LifecycleApi.GetBucketLifecycle``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetBucketLifecycle`: GetBucketLifecycleOutput
    fmt.Fprintf(os.Stdout, "Response from `LifecycleApi.GetBucketLifecycle`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetBucketLifecycleRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **lifecycle** | **bool** |  | |

### Return type

[**GetBucketLifecycleOutput**](../models/GetBucketLifecycleOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"LifecycleApiService.GetBucketLifecycle"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "LifecycleApiService.GetBucketLifecycle": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "LifecycleApiService.GetBucketLifecycle": {
    "port": "8443",
},
})
```


## PutBucketLifecycle

```go
var result  = PutBucketLifecycle(ctx, bucket)
                      .ContentMD5(contentMD5)
                      .Lifecycle(lifecycle)
                      .PutBucketLifecycleRequest(putBucketLifecycleRequest)
                      .Execute()
```

PutBucketLifecycle



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
    contentMD5 := "contentMD5_example" // string | 
    lifecycle := true // bool | 
    putBucketLifecycleRequest := *openapiclient.NewPutBucketLifecycleRequest() // PutBucketLifecycleRequest | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.LifecycleApi.PutBucketLifecycle(context.Background(), bucket).ContentMD5(contentMD5).Lifecycle(lifecycle).PutBucketLifecycleRequest(putBucketLifecycleRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `LifecycleApi.PutBucketLifecycle``: %v\n", err)
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

Other parameters are passed through a pointer to an apiPutBucketLifecycleRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **contentMD5** | **string** |  | |
| **lifecycle** | **bool** |  | |
| **putBucketLifecycleRequest** | [**PutBucketLifecycleRequest**](../models/PutBucketLifecycleRequest.md) |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"LifecycleApiService.PutBucketLifecycle"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "LifecycleApiService.PutBucketLifecycle": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "LifecycleApiService.PutBucketLifecycle": {
    "port": "8443",
},
})
```

