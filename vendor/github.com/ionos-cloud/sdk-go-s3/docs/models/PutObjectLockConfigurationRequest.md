# PutObjectLockConfigurationRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ObjectLockEnabled** | Pointer to **string** |  | [optional] |
|**Rule** | Pointer to [**PutObjectLockConfigurationRequestRule**](PutObjectLockConfigurationRequestRule.md) |  | [optional] |

## Methods

### NewPutObjectLockConfigurationRequest

`func NewPutObjectLockConfigurationRequest() *PutObjectLockConfigurationRequest`

NewPutObjectLockConfigurationRequest instantiates a new PutObjectLockConfigurationRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectLockConfigurationRequestWithDefaults

`func NewPutObjectLockConfigurationRequestWithDefaults() *PutObjectLockConfigurationRequest`

NewPutObjectLockConfigurationRequestWithDefaults instantiates a new PutObjectLockConfigurationRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObjectLockEnabled

`func (o *PutObjectLockConfigurationRequest) GetObjectLockEnabled() string`

GetObjectLockEnabled returns the ObjectLockEnabled field if non-nil, zero value otherwise.

### GetObjectLockEnabledOk

`func (o *PutObjectLockConfigurationRequest) GetObjectLockEnabledOk() (*string, bool)`

GetObjectLockEnabledOk returns a tuple with the ObjectLockEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectLockEnabled

`func (o *PutObjectLockConfigurationRequest) SetObjectLockEnabled(v string)`

SetObjectLockEnabled sets ObjectLockEnabled field to given value.

### HasObjectLockEnabled

`func (o *PutObjectLockConfigurationRequest) HasObjectLockEnabled() bool`

HasObjectLockEnabled returns a boolean if a field has been set.

### GetRule

`func (o *PutObjectLockConfigurationRequest) GetRule() PutObjectLockConfigurationRequestRule`

GetRule returns the Rule field if non-nil, zero value otherwise.

### GetRuleOk

`func (o *PutObjectLockConfigurationRequest) GetRuleOk() (*PutObjectLockConfigurationRequestRule, bool)`

GetRuleOk returns a tuple with the Rule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRule

`func (o *PutObjectLockConfigurationRequest) SetRule(v PutObjectLockConfigurationRequestRule)`

SetRule sets Rule field to given value.

### HasRule

`func (o *PutObjectLockConfigurationRequest) HasRule() bool`

HasRule returns a boolean if a field has been set.


