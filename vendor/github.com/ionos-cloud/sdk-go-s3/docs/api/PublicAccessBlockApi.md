# \PublicAccessBlockApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**DeletePublicAccessBlock**](PublicAccessBlockApi.md#DeletePublicAccessBlock) | **Delete** /{Bucket}?publicAccessBlock | DeletePublicAccessBlock|
|[**GetPublicAccessBlock**](PublicAccessBlockApi.md#GetPublicAccessBlock) | **Get** /{Bucket}?publicAccessBlock | GetPublicAccessBlock|
|[**PutPublicAccessBlock**](PublicAccessBlockApi.md#PutPublicAccessBlock) | **Put** /{Bucket}?publicAccessBlock | PutPublicAccessBlock|



## DeletePublicAccessBlock

```go
var result  = DeletePublicAccessBlock(ctx, bucket)
                      .Execute()
```

DeletePublicAccessBlock



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
    resource, resp, err := apiClient.PublicAccessBlockApi.DeletePublicAccessBlock(context.Background(), bucket).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicAccessBlockApi.DeletePublicAccessBlock``: %v\n", err)
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

Other parameters are passed through a pointer to an apiDeletePublicAccessBlockRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"PublicAccessBlockApiService.DeletePublicAccessBlock"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "PublicAccessBlockApiService.DeletePublicAccessBlock": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "PublicAccessBlockApiService.DeletePublicAccessBlock": {
    "port": "8443",
},
})
```


## GetPublicAccessBlock

```go
var result BlockPublicAccessOutput = GetPublicAccessBlock(ctx, bucket)
                      .Execute()
```

GetPublicAccessBlock



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
    resource, resp, err := apiClient.PublicAccessBlockApi.GetPublicAccessBlock(context.Background(), bucket).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicAccessBlockApi.GetPublicAccessBlock``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetPublicAccessBlock`: BlockPublicAccessOutput
    fmt.Fprintf(os.Stdout, "Response from `PublicAccessBlockApi.GetPublicAccessBlock`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetPublicAccessBlockRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|

### Return type

[**BlockPublicAccessOutput**](../models/BlockPublicAccessOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/jsonapplication/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"PublicAccessBlockApiService.GetPublicAccessBlock"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "PublicAccessBlockApiService.GetPublicAccessBlock": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "PublicAccessBlockApiService.GetPublicAccessBlock": {
    "port": "8443",
},
})
```


## PutPublicAccessBlock

```go
var result  = PutPublicAccessBlock(ctx, bucket)
                      .BlockPublicAccessPayload(blockPublicAccessPayload)
                      .ContentMD5(contentMD5)
                      .Execute()
```

PutPublicAccessBlock



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
    blockPublicAccessPayload := *openapiclient.NewBlockPublicAccessPayload() // BlockPublicAccessPayload | 
    contentMD5 := "contentMD5_example" // string |  (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.PublicAccessBlockApi.PutPublicAccessBlock(context.Background(), bucket).BlockPublicAccessPayload(blockPublicAccessPayload).ContentMD5(contentMD5).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PublicAccessBlockApi.PutPublicAccessBlock``: %v\n", err)
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

Other parameters are passed through a pointer to an apiPutPublicAccessBlockRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **blockPublicAccessPayload** | [**BlockPublicAccessPayload**](../models/BlockPublicAccessPayload.md) |  | |
| **contentMD5** | **string** |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"PublicAccessBlockApiService.PutPublicAccessBlock"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "PublicAccessBlockApiService.PutPublicAccessBlock": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "PublicAccessBlockApiService.PutPublicAccessBlock": {
    "port": "8443",
},
})
```

