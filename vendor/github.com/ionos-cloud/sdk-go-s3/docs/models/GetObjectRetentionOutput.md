# GetObjectRetentionOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Retention** | Pointer to [**ObjectLockRetention**](ObjectLockRetention.md) |  | [optional] |

## Methods

### NewGetObjectRetentionOutput

`func NewGetObjectRetentionOutput() *GetObjectRetentionOutput`

NewGetObjectRetentionOutput instantiates a new GetObjectRetentionOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetObjectRetentionOutputWithDefaults

`func NewGetObjectRetentionOutputWithDefaults() *GetObjectRetentionOutput`

NewGetObjectRetentionOutputWithDefaults instantiates a new GetObjectRetentionOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRetention

`func (o *GetObjectRetentionOutput) GetRetention() ObjectLockRetention`

GetRetention returns the Retention field if non-nil, zero value otherwise.

### GetRetentionOk

`func (o *GetObjectRetentionOutput) GetRetentionOk() (*ObjectLockRetention, bool)`

GetRetentionOk returns a tuple with the Retention field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetention

`func (o *GetObjectRetentionOutput) SetRetention(v ObjectLockRetention)`

SetRetention sets Retention field to given value.

### HasRetention

`func (o *GetObjectRetentionOutput) HasRetention() bool`

HasRetention returns a boolean if a field has been set.


