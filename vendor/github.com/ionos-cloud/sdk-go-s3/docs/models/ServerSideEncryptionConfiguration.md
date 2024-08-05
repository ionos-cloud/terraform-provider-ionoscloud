# ServerSideEncryptionConfiguration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Rules** | [**[]ServerSideEncryptionRule**](ServerSideEncryptionRule.md) |  | |

## Methods

### NewServerSideEncryptionConfiguration

`func NewServerSideEncryptionConfiguration(rules []ServerSideEncryptionRule, ) *ServerSideEncryptionConfiguration`

NewServerSideEncryptionConfiguration instantiates a new ServerSideEncryptionConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerSideEncryptionConfigurationWithDefaults

`func NewServerSideEncryptionConfigurationWithDefaults() *ServerSideEncryptionConfiguration`

NewServerSideEncryptionConfigurationWithDefaults instantiates a new ServerSideEncryptionConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRules

`func (o *ServerSideEncryptionConfiguration) GetRules() []ServerSideEncryptionRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *ServerSideEncryptionConfiguration) GetRulesOk() (*[]ServerSideEncryptionRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *ServerSideEncryptionConfiguration) SetRules(v []ServerSideEncryptionRule)`

SetRules sets Rules field to given value.



