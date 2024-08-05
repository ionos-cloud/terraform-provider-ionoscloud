# \ReplicationApi

All URIs are relative to *https://s3.eu-central-3.ionoscloud.com*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**GetBucketReplication**](ReplicationApi.md#GetBucketReplication) | **Get** /{Bucket}?replication | GetBucketReplication|



## GetBucketReplication

```go
var result GetBucketReplicationOutput = GetBucketReplication(ctx, bucket)
                      .Replication(replication)
                      .Execute()
```

GetBucketReplication



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
    replication := true // bool | 

    configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
    apiClient := ionoscloud.NewAPIClient(configuration)
    resource, resp, err := apiClient.ReplicationApi.GetBucketReplication(context.Background(), bucket).Replication(replication).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ReplicationApi.GetBucketReplication``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
    }
    // response from `GetBucketReplication`: GetBucketReplicationOutput
    fmt.Fprintf(os.Stdout, "Response from `ReplicationApi.GetBucketReplication`: %v\n", resource)
}
```

### Path Parameters


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
|**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.|
|**bucket** | **string** |  | |

### Other Parameters

Other parameters are passed through a pointer to an apiGetBucketReplicationRequest struct via the builder pattern


|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **replication** | **bool** |  | |

### Return type

[**GetBucketReplicationOutput**](../models/GetBucketReplicationOutput.md)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/xml


### URLs Configuration per Operation
Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"ReplicationApiService.GetBucketReplication"` string.
Similar rules for overriding default operation server index and variables apply by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), {packageName}.ContextOperationServerIndices, map[string]int{
    "ReplicationApiService.GetBucketReplication": 2,
})
ctx = context.WithValue(context.Background(), {packageName}.ContextOperationServerVariables, map[string]map[string]string{
    "ReplicationApiService.GetBucketReplication": {
    "port": "8443",
},
})
```

