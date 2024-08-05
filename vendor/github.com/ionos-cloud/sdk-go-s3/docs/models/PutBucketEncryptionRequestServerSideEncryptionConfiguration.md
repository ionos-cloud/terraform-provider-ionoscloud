# PutBucketEncryptionRequestServerSideEncryptionConfiguration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Rules** | Pointer to [**[]ServerSideEncryptionRule**](ServerSideEncryptionRule.md) |  | [optional] |

## Methods

### NewPutBucketEncryptionRequestServerSideEncryptionConfiguration

`func NewPutBucketEncryptionRequestServerSideEncryptionConfiguration() *PutBucketEncryptionRequestServerSideEncryptionConfiguration`

NewPutBucketEncryptionRequestServerSideEncryptionConfiguration instantiates a new PutBucketEncryptionRequestServerSideEncryptionConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutBucketEncryptionRequestServerSideEncryptionConfigurationWithDefaults

`func NewPutBucketEncryptionRequestServerSideEncryptionConfigurationWithDefaults() *PutBucketEncryptionRequestServerSideEncryptionConfiguration`

NewPutBucketEncryptionRequestServerSideEncryptionConfigurationWithDefaults instantiates a new PutBucketEncryptionRequestServerSideEncryptionConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRules

`func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) GetRules() []ServerSideEncryptionRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) GetRulesOk() (*[]ServerSideEncryptionRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) SetRules(v []ServerSideEncryptionRule)`

SetRules sets Rules field to given value.

### HasRules

`func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) HasRules() bool`

HasRules returns a boolean if a field has been set.


