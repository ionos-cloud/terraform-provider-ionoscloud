# PutObjectLockConfigurationRequestRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**DefaultRetention** | Pointer to [**DefaultRetention**](DefaultRetention.md) |  | [optional] |

## Methods

### NewPutObjectLockConfigurationRequestRule

`func NewPutObjectLockConfigurationRequestRule() *PutObjectLockConfigurationRequestRule`

NewPutObjectLockConfigurationRequestRule instantiates a new PutObjectLockConfigurationRequestRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectLockConfigurationRequestRuleWithDefaults

`func NewPutObjectLockConfigurationRequestRuleWithDefaults() *PutObjectLockConfigurationRequestRule`

NewPutObjectLockConfigurationRequestRuleWithDefaults instantiates a new PutObjectLockConfigurationRequestRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefaultRetention

`func (o *PutObjectLockConfigurationRequestRule) GetDefaultRetention() DefaultRetention`

GetDefaultRetention returns the DefaultRetention field if non-nil, zero value otherwise.

### GetDefaultRetentionOk

`func (o *PutObjectLockConfigurationRequestRule) GetDefaultRetentionOk() (*DefaultRetention, bool)`

GetDefaultRetentionOk returns a tuple with the DefaultRetention field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultRetention

`func (o *PutObjectLockConfigurationRequestRule) SetDefaultRetention(v DefaultRetention)`

SetDefaultRetention sets DefaultRetention field to given value.

### HasDefaultRetention

`func (o *PutObjectLockConfigurationRequestRule) HasDefaultRetention() bool`

HasDefaultRetention returns a boolean if a field has been set.


