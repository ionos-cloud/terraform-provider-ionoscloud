# ServerSideEncryptionRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ApplyServerSideEncryptionByDefault** | Pointer to [**ServerSideEncryptionByDefault**](ServerSideEncryptionByDefault.md) |  | [optional] |

## Methods

### NewServerSideEncryptionRule

`func NewServerSideEncryptionRule() *ServerSideEncryptionRule`

NewServerSideEncryptionRule instantiates a new ServerSideEncryptionRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerSideEncryptionRuleWithDefaults

`func NewServerSideEncryptionRuleWithDefaults() *ServerSideEncryptionRule`

NewServerSideEncryptionRuleWithDefaults instantiates a new ServerSideEncryptionRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApplyServerSideEncryptionByDefault

`func (o *ServerSideEncryptionRule) GetApplyServerSideEncryptionByDefault() ServerSideEncryptionByDefault`

GetApplyServerSideEncryptionByDefault returns the ApplyServerSideEncryptionByDefault field if non-nil, zero value otherwise.

### GetApplyServerSideEncryptionByDefaultOk

`func (o *ServerSideEncryptionRule) GetApplyServerSideEncryptionByDefaultOk() (*ServerSideEncryptionByDefault, bool)`

GetApplyServerSideEncryptionByDefaultOk returns a tuple with the ApplyServerSideEncryptionByDefault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplyServerSideEncryptionByDefault

`func (o *ServerSideEncryptionRule) SetApplyServerSideEncryptionByDefault(v ServerSideEncryptionByDefault)`

SetApplyServerSideEncryptionByDefault sets ApplyServerSideEncryptionByDefault field to given value.

### HasApplyServerSideEncryptionByDefault

`func (o *ServerSideEncryptionRule) HasApplyServerSideEncryptionByDefault() bool`

HasApplyServerSideEncryptionByDefault returns a boolean if a field has been set.


