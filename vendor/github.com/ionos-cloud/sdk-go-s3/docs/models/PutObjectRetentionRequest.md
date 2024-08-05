# PutObjectRetentionRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Retention** | Pointer to [**PutObjectRetentionRequestRetention**](PutObjectRetentionRequestRetention.md) |  | [optional] |

## Methods

### NewPutObjectRetentionRequest

`func NewPutObjectRetentionRequest() *PutObjectRetentionRequest`

NewPutObjectRetentionRequest instantiates a new PutObjectRetentionRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectRetentionRequestWithDefaults

`func NewPutObjectRetentionRequestWithDefaults() *PutObjectRetentionRequest`

NewPutObjectRetentionRequestWithDefaults instantiates a new PutObjectRetentionRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRetention

`func (o *PutObjectRetentionRequest) GetRetention() PutObjectRetentionRequestRetention`

GetRetention returns the Retention field if non-nil, zero value otherwise.

### GetRetentionOk

`func (o *PutObjectRetentionRequest) GetRetentionOk() (*PutObjectRetentionRequestRetention, bool)`

GetRetentionOk returns a tuple with the Retention field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetention

`func (o *PutObjectRetentionRequest) SetRetention(v PutObjectRetentionRequestRetention)`

SetRetention sets Retention field to given value.

### HasRetention

`func (o *PutObjectRetentionRequest) HasRetention() bool`

HasRetention returns a boolean if a field has been set.


