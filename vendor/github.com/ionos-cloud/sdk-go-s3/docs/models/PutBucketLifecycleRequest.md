# PutBucketLifecycleRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**LifecycleConfiguration** | Pointer to [**PutBucketLifecycleRequestLifecycleConfiguration**](PutBucketLifecycleRequestLifecycleConfiguration.md) |  | [optional] |

## Methods

### NewPutBucketLifecycleRequest

`func NewPutBucketLifecycleRequest() *PutBucketLifecycleRequest`

NewPutBucketLifecycleRequest instantiates a new PutBucketLifecycleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutBucketLifecycleRequestWithDefaults

`func NewPutBucketLifecycleRequestWithDefaults() *PutBucketLifecycleRequest`

NewPutBucketLifecycleRequestWithDefaults instantiates a new PutBucketLifecycleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLifecycleConfiguration

`func (o *PutBucketLifecycleRequest) GetLifecycleConfiguration() PutBucketLifecycleRequestLifecycleConfiguration`

GetLifecycleConfiguration returns the LifecycleConfiguration field if non-nil, zero value otherwise.

### GetLifecycleConfigurationOk

`func (o *PutBucketLifecycleRequest) GetLifecycleConfigurationOk() (*PutBucketLifecycleRequestLifecycleConfiguration, bool)`

GetLifecycleConfigurationOk returns a tuple with the LifecycleConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLifecycleConfiguration

`func (o *PutBucketLifecycleRequest) SetLifecycleConfiguration(v PutBucketLifecycleRequestLifecycleConfiguration)`

SetLifecycleConfiguration sets LifecycleConfiguration field to given value.

### HasLifecycleConfiguration

`func (o *PutBucketLifecycleRequest) HasLifecycleConfiguration() bool`

HasLifecycleConfiguration returns a boolean if a field has been set.


