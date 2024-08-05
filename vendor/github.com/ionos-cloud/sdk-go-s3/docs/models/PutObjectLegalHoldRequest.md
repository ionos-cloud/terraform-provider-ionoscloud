# PutObjectLegalHoldRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**LegalHold** | Pointer to [**PutObjectLegalHoldRequestLegalHold**](PutObjectLegalHoldRequestLegalHold.md) |  | [optional] |

## Methods

### NewPutObjectLegalHoldRequest

`func NewPutObjectLegalHoldRequest() *PutObjectLegalHoldRequest`

NewPutObjectLegalHoldRequest instantiates a new PutObjectLegalHoldRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectLegalHoldRequestWithDefaults

`func NewPutObjectLegalHoldRequestWithDefaults() *PutObjectLegalHoldRequest`

NewPutObjectLegalHoldRequestWithDefaults instantiates a new PutObjectLegalHoldRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLegalHold

`func (o *PutObjectLegalHoldRequest) GetLegalHold() PutObjectLegalHoldRequestLegalHold`

GetLegalHold returns the LegalHold field if non-nil, zero value otherwise.

### GetLegalHoldOk

`func (o *PutObjectLegalHoldRequest) GetLegalHoldOk() (*PutObjectLegalHoldRequestLegalHold, bool)`

GetLegalHoldOk returns a tuple with the LegalHold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLegalHold

`func (o *PutObjectLegalHoldRequest) SetLegalHold(v PutObjectLegalHoldRequestLegalHold)`

SetLegalHold sets LegalHold field to given value.

### HasLegalHold

`func (o *PutObjectLegalHoldRequest) HasLegalHold() bool`

HasLegalHold returns a boolean if a field has been set.


