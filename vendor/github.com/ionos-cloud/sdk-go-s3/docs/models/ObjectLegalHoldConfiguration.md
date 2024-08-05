# ObjectLegalHoldConfiguration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Status** | Pointer to **string** | Object Legal Hold status | [optional] |

## Methods

### NewObjectLegalHoldConfiguration

`func NewObjectLegalHoldConfiguration() *ObjectLegalHoldConfiguration`

NewObjectLegalHoldConfiguration instantiates a new ObjectLegalHoldConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObjectLegalHoldConfigurationWithDefaults

`func NewObjectLegalHoldConfigurationWithDefaults() *ObjectLegalHoldConfiguration`

NewObjectLegalHoldConfigurationWithDefaults instantiates a new ObjectLegalHoldConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *ObjectLegalHoldConfiguration) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ObjectLegalHoldConfiguration) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ObjectLegalHoldConfiguration) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ObjectLegalHoldConfiguration) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


