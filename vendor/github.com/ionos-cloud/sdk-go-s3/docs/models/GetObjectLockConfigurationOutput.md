# GetObjectLockConfigurationOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ObjectLockEnabled** | Pointer to **string** |  | [optional] |
|**Rule** | Pointer to [**ObjectLockRule**](ObjectLockRule.md) |  | [optional] |

## Methods

### NewGetObjectLockConfigurationOutput

`func NewGetObjectLockConfigurationOutput() *GetObjectLockConfigurationOutput`

NewGetObjectLockConfigurationOutput instantiates a new GetObjectLockConfigurationOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetObjectLockConfigurationOutputWithDefaults

`func NewGetObjectLockConfigurationOutputWithDefaults() *GetObjectLockConfigurationOutput`

NewGetObjectLockConfigurationOutputWithDefaults instantiates a new GetObjectLockConfigurationOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObjectLockEnabled

`func (o *GetObjectLockConfigurationOutput) GetObjectLockEnabled() string`

GetObjectLockEnabled returns the ObjectLockEnabled field if non-nil, zero value otherwise.

### GetObjectLockEnabledOk

`func (o *GetObjectLockConfigurationOutput) GetObjectLockEnabledOk() (*string, bool)`

GetObjectLockEnabledOk returns a tuple with the ObjectLockEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectLockEnabled

`func (o *GetObjectLockConfigurationOutput) SetObjectLockEnabled(v string)`

SetObjectLockEnabled sets ObjectLockEnabled field to given value.

### HasObjectLockEnabled

`func (o *GetObjectLockConfigurationOutput) HasObjectLockEnabled() bool`

HasObjectLockEnabled returns a boolean if a field has been set.

### GetRule

`func (o *GetObjectLockConfigurationOutput) GetRule() ObjectLockRule`

GetRule returns the Rule field if non-nil, zero value otherwise.

### GetRuleOk

`func (o *GetObjectLockConfigurationOutput) GetRuleOk() (*ObjectLockRule, bool)`

GetRuleOk returns a tuple with the Rule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRule

`func (o *GetObjectLockConfigurationOutput) SetRule(v ObjectLockRule)`

SetRule sets Rule field to given value.

### HasRule

`func (o *GetObjectLockConfigurationOutput) HasRule() bool`

HasRule returns a boolean if a field has been set.


