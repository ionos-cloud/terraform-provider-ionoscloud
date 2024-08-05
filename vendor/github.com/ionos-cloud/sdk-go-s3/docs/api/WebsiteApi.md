# \WebsiteApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**DeleteBucketWebsite**](WebsiteApi.md#DeleteBucketWebsite) | **Delete** /{Bucket}?website | DeleteBucketWebsite|
|[**GetBucketWebsite**](WebsiteApi.md#GetBucketWebsite) | **Get** /{Bucket}?website | GetBucketWebsite|
|[**PutBucketWebsite**](WebsiteApi.md#PutBucketWebsite) | **Put** /{Bucket}?website | PutBucketWebsite|



## DeleteBucketWebsite

```go
var result  = DeleteBucketWebsite(ctx, bucket)
                      .Website(website)
                      .Execute()
```

DeleteBucketWebsite



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
    website := true // bool | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.WebsiteApi.DeleteBucketWebsite(context.Background(), bucket).Website(website).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebsiteApi.DeleteBucketWebsite``: %v\n", err)
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

Other parameters are passed through a pointer to an apiDeleteBucketWebsiteRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **website** | **bool** |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"WebsiteApiService.DeleteBucketWebsite"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "WebsiteApiService.DeleteBucketWebsite": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "WebsiteApiService.DeleteBucketWebsite": {
    "port": "8443",
},
})
```


## GetBucketWebsite

```go
var result GetBucketWebsiteOutput = GetBucketWebsite(ctx, bucket)
                      .Website(website)
                      .Execute()
```

GetBucketWebsite



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
    website := true // bool | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.WebsiteApi.GetBucketWebsite(context.Background(), bucket).Website(website).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebsiteApi.GetBucketWebsite``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetBucketWebsite`: GetBucketWebsiteOutput
    fmt.Fprintf(os.Stdout, "Response from `WebsiteApi.GetBucketWebsite`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetBucketWebsiteRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **website** | **bool** |  | |

### Return type

[**GetBucketWebsiteOutput**](../models/GetBucketWebsiteOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"WebsiteApiService.GetBucketWebsite"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "WebsiteApiService.GetBucketWebsite": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "WebsiteApiService.GetBucketWebsite": {
    "port": "8443",
},
})
```


## PutBucketWebsite

```go
var result  = PutBucketWebsite(ctx, bucket)
                      .Website(website)
                      .PutBucketWebsiteRequest(putBucketWebsiteRequest)
                      .ContentMD5(contentMD5)
                      .Execute()
```

PutBucketWebsite



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
    website := true // bool | 
    putBucketWebsiteRequest := *openapiclient.NewPutBucketWebsiteRequest(*openapiclient.NewPutBucketWebsiteRequestWebsiteConfiguration()) // PutBucketWebsiteRequest | 
    contentMD5 := "contentMD5_example" // string |  (optional)

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.WebsiteApi.PutBucketWebsite(context.Background(), bucket).Website(website).PutBucketWebsiteRequest(putBucketWebsiteRequest).ContentMD5(contentMD5).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebsiteApi.PutBucketWebsite``: %v\n", err)
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

Other parameters are passed through a pointer to an apiPutBucketWebsiteRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **website** | **bool** |  | |
| **putBucketWebsiteRequest** | [**PutBucketWebsiteRequest**](../models/PutBucketWebsiteRequest.md) |  | |
| **contentMD5** | **string** |  | |

### Return type

 (empty response body)

### HTTP request headers

- **Content-Type**: application/xml
- **Accept**: Not defined


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"WebsiteApiService.PutBucketWebsite"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "WebsiteApiService.PutBucketWebsite": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "WebsiteApiService.PutBucketWebsite": {
    "port": "8443",
},
})
```

