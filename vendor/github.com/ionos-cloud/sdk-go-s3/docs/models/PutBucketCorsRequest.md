# PutBucketCorsRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**CORSConfiguration** | [**PutBucketCorsRequestCORSConfiguration**](PutBucketCorsRequestCORSConfiguration.md) |  | |

## Methods

### NewPutBucketCorsRequest

`func NewPutBucketCorsRequest(cORSConfiguration PutBucketCorsRequestCORSConfiguration, ) *PutBucketCorsRequest`

NewPutBucketCorsRequest instantiates a new PutBucketCorsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutBucketCorsRequestWithDefaults

`func NewPutBucketCorsRequestWithDefaults() *PutBucketCorsRequest`

NewPutBucketCorsRequestWithDefaults instantiates a new PutBucketCorsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCORSConfiguration

`func (o *PutBucketCorsRequest) GetCORSConfiguration() PutBucketCorsRequestCORSConfiguration`

GetCORSConfiguration returns the CORSConfiguration field if non-nil, zero value otherwise.

### GetCORSConfigurationOk

`func (o *PutBucketCorsRequest) GetCORSConfigurationOk() (*PutBucketCorsRequestCORSConfiguration, bool)`

GetCORSConfigurationOk returns a tuple with the CORSConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCORSConfiguration

`func (o *PutBucketCorsRequest) SetCORSConfiguration(v PutBucketCorsRequestCORSConfiguration)`

SetCORSConfiguration sets CORSConfiguration field to given value.



