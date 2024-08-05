# PutObjectLegalHoldRequestLegalHold

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Status** | Pointer to [**ObjectLegalHoldConfiguration**](ObjectLegalHoldConfiguration.md) |  | [optional] |

## Methods

### NewPutObjectLegalHoldRequestLegalHold

`func NewPutObjectLegalHoldRequestLegalHold() *PutObjectLegalHoldRequestLegalHold`

NewPutObjectLegalHoldRequestLegalHold instantiates a new PutObjectLegalHoldRequestLegalHold object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectLegalHoldRequestLegalHoldWithDefaults

`func NewPutObjectLegalHoldRequestLegalHoldWithDefaults() *PutObjectLegalHoldRequestLegalHold`

NewPutObjectLegalHoldRequestLegalHoldWithDefaults instantiates a new PutObjectLegalHoldRequestLegalHold object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *PutObjectLegalHoldRequestLegalHold) GetStatus() ObjectLegalHoldConfiguration`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PutObjectLegalHoldRequestLegalHold) GetStatusOk() (*ObjectLegalHoldConfiguration, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PutObjectLegalHoldRequestLegalHold) SetStatus(v ObjectLegalHoldConfiguration)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *PutObjectLegalHoldRequestLegalHold) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


